# 技术设计文档

## 1. 架构概述

### 1.1 整体架构
```
┌─────────────────┐
│   Main Entry    │ (cmd/game/main.go)
└─────────┬───────┘
          │
┌─────────▼───────┐
│   Game Core     │ (internal/game/)
└─────────┬───────┘
          │
┌─────────▼───────┐
│   Scenes        │ (internal/scenes/)
└─────────┬───────┘
          │
┌─────────▼───────┐
│   Entities      │ (internal/entities/)
└─────────────────┘
```

### 1.2 模块职责
- **Game Core**: 游戏主循环、场景管理
- **Scenes**: 不同游戏场景的实现
- **Entities**: 游戏实体（玩家、敌人、道具等）
- **Utils**: 工具函数（输入处理、数学计算等）
- **Config**: 游戏配置管理

## 2. 核心系统设计

### 2.1 场景系统
```go
type Scene interface {
    Update() error
    Draw(screen *ebiten.Image)
}
```

**职责**:
- 管理不同游戏状态
- 处理场景切换
- 维护场景生命周期

### 2.2 实体系统
```go
type Entity interface {
    Update()
    Draw(screen *ebiten.Image)
    GetPosition() (float64, float64)
    SetPosition(x, y float64)
}
```

**职责**:
- 管理游戏对象
- 处理实体交互
- 维护实体状态

### 2.3 输入系统
```go
type InputManager struct {
    keys map[ebiten.Key]bool
}
```

**职责**:
- 处理键盘输入
- 管理输入状态
- 提供输入查询接口

### 2.4 武器与射击系统
- 数据结构：复用 `internal/scenes/shooter_types.go` 中的 `FireSkillConfig`、`Bullet`
- 关键能力：连发、散射、追踪、穿透、弹速
- 主要接口（现状）：
```
shoot()                 // 由玩家输入触发
processScheduledShots() // 处理计划连发，在 Update 中调用
updateBullets()         // 子弹更新（含追踪转向）
```
- 追踪限幅：按 `HomingTurnRateRad` 限制每帧旋转，避免瞬折。
- 数值上限：为射速/穿透/散射角设置软硬上限，兼顾性能与可读性。

### 2.9 刷怪与波次系统
- 波次：`[1,2,4,...,2^10]`；同屏上限 `maxSimultaneous`；`batchSize` 批量生成。
- 冷却：随波次缩短，最低 250ms；清场→下一波；全部完成后生成 Boss。
- 敌机血量：`1 + waveIndex/3` 上限 6。

### 2.10 敌人AI与Boss阶段（预留）
- 小怪：直线/蛇形/冲锋/召唤/护盾；Boss 多阶段、招式切换、输出窗口。

### 2.11 GM 调试面板
- Tab 切页；上下选择、左右调整；支持射击音效参数即调即听。
- 可扩展：刷怪/清场/波次跳转/碰撞可视化/JSON 导入导出数值。

## 3. 数据流设计

### 3.1 更新循环
```
Input → Update → Draw → Present
```

### 3.2 场景切换流程
```
Current Scene → Transition → New Scene

### 3.3 机制时序
- 射击与连发：`SPACE → shoot() → schedule → Update: processScheduledShots()`
- 波次推进：`spawnEnemies() → clear → nextWave() → all waves → spawnBoss()`
```

## 4. 资源管理

### 4.1 资源类型
- **图像资源**: PNG格式，存储在 `assets/images/`
- **音频资源**: OGG格式，存储在 `assets/sounds/`
- **字体资源**: TTF格式，存储在 `assets/fonts/`

### 4.2 资源加载策略
- 预加载关键资源
- 按需加载场景资源
- 资源缓存管理

## 5. 性能优化

### 5.1 渲染优化
- 对象池模式
- 批量绘制
- 视锥剔除

### 5.2 内存管理
- 及时释放不用的资源
- 避免内存泄漏
- 合理使用缓存

### 5.3 逻辑优化（建议）
- 空间划分：网格/四叉树分桶，降低碰撞复杂度
- 对象池：`Bullet`、`Explosion` 循环复用，避免 GC 峰值
- 节流合并：音效/飘字/震屏频率限制，UI 文本合并绘制

## 6. 调试系统

### 6.1 调试信息
- FPS显示
- 内存使用情况
- 实体数量统计

### 6.2 调试工具
- 场景切换快捷键
- 实体位置显示
- 碰撞检测可视化

### 6.3 GM 面板扩展（建议）
- 一键刷怪/清场、波次跳转、Boss 直达
- 数值面板导入/导出 JSON；显示当前 `FireSkillConfig`
- 性能计数：子弹/敌人/GC 次数/分配

## 7. 配置管理

### 7.1 配置文件
- 游戏设置 (config.json)
- 键位设置 (controls.json)
- 关卡数据 (levels.json)

### 7.2 配置热重载
- 运行时配置更新
- 设置持久化

### 7.3 数据Schema（示例）
ships.json
```
{ "ships": [
  {"id":"scout","hp":2,"speed":6.0,"size":0.8,"passive":"overheat_delay"},
  {"id":"tank","hp":5,"speed":3.5,"size":1.3,"passive":"+1_penetration"}
]}
```
skills.json（对应 FireSkillConfig）
```
{ "skills": [
  {"id":"basic_gun","fireRateHz":5.0,"bulletsPerShot":1,"spreadDeg":2,
   "bulletSpeed":8,"burstChance":0.0,"penetrationCount":0,
   "enableHoming":false,"homingTurnRateRad":0.08,"burstIntervalMs":60}
]}
```
waves.json
```
{ "waves":[1,2,4,8,16,32,64,128,256,512,1024],
  "maxSimultaneous":20,"batchSize":3,
  "minEnemyDelayMs":250,"baseEnemyDelayMs":800 }
```
loot.json（可选）
```
{ "drops": [
  {"id":"score_small","chance":0.5,"score":5},
  {"id":"heal_small","chance":0.1,"heal":1}
]}
```

## 8. 测试策略

### 8.1 单元测试
- 实体行为测试
- 工具函数测试
- 配置加载测试

### 8.2 集成测试
- 场景切换测试
- 输入处理测试
- 渲染流程测试

### 8.3 逻辑与数值测试（建议）
- 射速与冷却一致性、连发计划时序
- 追踪子弹限幅与目标选择
- 升级阈值与冷却生效
- 刷怪上限、批量与波切换

### 8.4 集成/回归
- 波间升级→属性变化→战斗体感回归
- Boss 击杀后的胜利与重开（R）
- GM 面板调参即时生效校验

## 9. 部署考虑

### 9.1 构建配置
- 跨平台编译
- 资源打包
- 版本管理

### 9.2 发布流程
- 自动化构建
- 测试验证
- 发布打包
