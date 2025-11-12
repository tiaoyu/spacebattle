package sound

import (
	"encoding/binary"
	"math"
	"math/rand"
	"sync"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

// 简易的8-bit风格声音生成与播放（使用16-bit PCM以兼容 Ebiten 音频）

const sampleRate = 22050

// WaveformConfig 可调参数
// Waveform: "square" | "triangle" | "noise"
type WaveformConfig struct {
	DurationSec float64
	BaseFreq    float64
	MinFreq     float64
	SweepFactor float64
	Decay       float64
	Amplitude   float64 // 0..1
	Waveform    string
}

var (
	once sync.Once
	ctx  *audio.Context

	cfgMu sync.RWMutex
	// 预设射击音效
	shootConfig = WaveformConfig{
		DurationSec: 0.2,
		BaseFreq:    600.0,
		MinFreq:     200.0,
		SweepFactor: 1.0,
		Decay:       1.0,
		Amplitude:   0.2,
		Waveform:    "triangle",
	}
	// 预设击中音效
	hitConfig = WaveformConfig{
		DurationSec: 0.2,
		BaseFreq:    600.0,
		MinFreq:     200.0,
		SweepFactor: 1.0,
		Decay:       1.0,
		Amplitude:   0.2,
		Waveform:    "noise",
	}

	// PCM 缓存
	shootPCM      []byte
	shootPCMValid bool
	shootLastCfg  WaveformConfig

	hitPCM      []byte
	hitPCMValid bool
	hitLastCfg  WaveformConfig
)

// Init 初始化音频上下文
func Init() {
	once.Do(func() {
		ctx = audio.NewContext(sampleRate)
	})
}

// GetShootConfig 获取当前配置（拷贝）
func GetShootConfig() WaveformConfig {
	cfgMu.RLock()
	defer cfgMu.RUnlock()
	return shootConfig
}

// SetShootConfig 设置当前配置（覆盖）
func SetShootConfig(c WaveformConfig) {
	cfgMu.Lock()
	shootConfig = c
	shootPCMValid = false // 使缓存失效
	cfgMu.Unlock()
}

// GetHitConfig 获取当前配置（拷贝）
func GetHitConfig() WaveformConfig {
	cfgMu.RLock()
	defer cfgMu.RUnlock()
	return hitConfig
}

// SetHitConfig 设置当前配置（覆盖）
func SetHitConfig(c WaveformConfig) {
	cfgMu.Lock()
	hitConfig = c
	hitPCMValid = false // 使缓存失效
	cfgMu.Unlock()
}

// PlayShoot 播放射击音效（使用缓存；配置变化时重新生成）
func PlayShoot() {
	if ctx == nil {
		Init()
	}
	cfgMu.Lock()
	cfg := shootConfig
	if !shootPCMValid || !equalCfg(shootLastCfg, cfg) {
		shootPCM = generatePCM(cfg)
		shootLastCfg = cfg
		shootPCMValid = true
	}
	data := shootPCM
	cfgMu.Unlock()

	p := ctx.NewPlayerFromBytes(data)
	p.Play()
}

// PlayHit 播放击中音效（使用缓存；配置变化时重新生成）
func PlayHit() {
	if ctx == nil {
		Init()
	}
	cfgMu.Lock()
	cfg := hitConfig
	if !hitPCMValid || !equalCfg(hitLastCfg, cfg) {
		hitPCM = generatePCM(cfg)
		hitLastCfg = cfg
		hitPCMValid = true
	}
	data := hitPCM
	cfgMu.Unlock()

	p := ctx.NewPlayerFromBytes(data)
	p.Play()
}

// generatePCM 生成一个短促的8-bit风格“pew”音（方波/三角/噪声 + 衰减），导出为16-bit PCM
func generatePCM(cfg WaveformConfig) []byte {
	// 采样点数量
	n := max(int(float64(sampleRate)*clamp(cfg.DurationSec, 0.01, 0.5)), 1)
	baseFreq := clamp(cfg.BaseFreq, 50, 4000)
	minFreq := clamp(cfg.MinFreq, 20, baseFreq)
	sweep := clamp(cfg.SweepFactor, 0, 1)
	decay := clamp(cfg.Decay, 1, 80)
	amp := clamp(cfg.Amplitude, 0, 1)

	pcm := make([]byte, n*2)
	for i := range n {
		t := float64(i) / sampleRate
		// 频率从 baseFreq 向下滑至 minFreq
		f := baseFreq * (1.0 - sweep*t)
		if f < minFreq {
			f = minFreq
		}

		// 选择波形
		var v float64
		switch cfg.Waveform {
		case "triangle":
			// 归一化三角波 [-1,1]
			phase := math.Mod(2*math.Pi*f*t, 2*math.Pi)
			v = 2 / math.Pi * math.Asin(math.Sin(phase))
		case "noise":
			v = 2*rand.Float64() - 1
		default: // square
			if math.Sin(2*math.Pi*f*t) >= 0 {
				v = 1
			} else {
				v = -1
			}
		}

		// 指数衰减包络
		env := math.Exp(-decay * t)
		s := v * env
		// 映射到 int16
		amp16 := int16(amp * 16000)
		binary.LittleEndian.PutUint16(pcm[i*2:(i+1)*2], uint16(int16(float64(amp16)*s)))
	}
	return pcm
}

func equalCfg(a, b WaveformConfig) bool {
	return a.DurationSec == b.DurationSec &&
		a.BaseFreq == b.BaseFreq &&
		a.MinFreq == b.MinFreq &&
		a.SweepFactor == b.SweepFactor &&
		a.Decay == b.Decay &&
		a.Amplitude == b.Amplitude &&
		a.Waveform == b.Waveform
}

func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
