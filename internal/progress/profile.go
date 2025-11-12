package progress

// 简易全局档案：存储功勋等进度

type Profile struct {
	Merits int
}

var globalProfile = &Profile{}

func GetProfile() *Profile { return globalProfile }

func AddMerits(n int) {
	if n <= 0 {
		return
	}
	globalProfile.Merits += n
	_ = saveMerits(globalProfile.Merits)
}

func GetMerits() int { return globalProfile.Merits }

// SpendMerits 消耗功勋，成功返回 true
func SpendMerits(n int) bool {
	if n <= 0 {
		return true
	}
	if globalProfile.Merits < n {
		return false
	}
	globalProfile.Merits -= n
	_ = saveMerits(globalProfile.Merits)
	return true
}

// Load 从SQLite加载当前档案
func Load() error {
	m, err := loadMerits()
	if err == nil {
		globalProfile.Merits = m
	}
	return nil
}
