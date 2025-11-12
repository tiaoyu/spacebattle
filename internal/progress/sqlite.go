package progress

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"

	_ "modernc.org/sqlite"
)

var (
	db   *sql.DB
	once sync.Once
)

// Init 初始化SQLite存储
func Init(path string) error {
	var err error
	once.Do(func() {
		db, err = sql.Open("sqlite", path)
		if err != nil {
			return
		}
		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS kv (
            key TEXT PRIMARY KEY,
            value TEXT NOT NULL
        );`)
	})
	return err
}

func kvGet(key string) (string, bool, error) {
	if db == nil {
		return "", false, fmt.Errorf("progress DB not initialized")
	}
	var val string
	err := db.QueryRow("SELECT value FROM kv WHERE key=?", key).Scan(&val)
	if err == sql.ErrNoRows {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return val, true, nil
}

func kvSet(key, value string) error {
	if db == nil {
		return fmt.Errorf("progress DB not initialized")
	}
	_, err := db.Exec("INSERT INTO kv(key, value) VALUES(?, ?) ON CONFLICT(key) DO UPDATE SET value=excluded.value", key, value)
	return err
}

// UpgradeData 用于持久化的加点数据
type UpgradeData struct {
	ModFireRateHz     float64 `json:"fire_rate_hz"`
	ModBulletsPerShot int     `json:"bullets_per_shot"`
	ModPenetration    int     `json:"penetration"`
	ModSpreadDeltaDeg float64 `json:"spread_delta_deg"`
	ModBulletSpeed    float64 `json:"bullet_speed"`
	ModBurstChance    float64 `json:"burst_chance"`
	ModEnableHoming   bool    `json:"enable_homing"`
	ModTurnRateRad    float64 `json:"turn_rate_rad"`
}

const (
	keyMerits   = "merits"
	keyUpgrades = "upgrades"
)

// 持久化API

func loadMerits() (int, error) {
	s, ok, err := kvGet(keyMerits)
	if err != nil || !ok {
		return 0, err
	}
	var m int
	_, err = fmt.Sscanf(s, "%d", &m)
	if err != nil {
		return 0, err
	}
	return m, nil
}

func saveMerits(m int) error {
	return kvSet(keyMerits, fmt.Sprintf("%d", m))
}

func GetUpgrades() (UpgradeData, error) {
	var u UpgradeData
	s, ok, err := kvGet(keyUpgrades)
	if err != nil || !ok {
		return u, err
	}
	if err := json.Unmarshal([]byte(s), &u); err != nil {
		return UpgradeData{}, err
	}
	return u, nil
}

func SaveUpgrades(u UpgradeData) error {
	b, err := json.Marshal(u)
	if err != nil {
		return err
	}
	return kvSet(keyUpgrades, string(b))
}
