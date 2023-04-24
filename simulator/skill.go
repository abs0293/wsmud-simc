package simulator

import (
	"log"
)

const (
	Skill_Type_Unarmed int = iota
	Skill_Type_Weapon
	Skill_Type_Force
	Skill_Type_Dodge
	Skill_Type_Parry
	Skill_Type_Throwing
)

var SkillTypeValue = map[string]int{
	"拳脚": Skill_Type_Unarmed,
	"内功": Skill_Type_Force,
	"武器": Skill_Type_Weapon,
	"轻功": Skill_Type_Dodge,
	"招架": Skill_Type_Parry,
	"暗器": Skill_Type_Throwing,
}

var SkillTypeName = map[int]string{
	Skill_Type_Unarmed:  "拳脚",
	Skill_Type_Force:    "内功",
	Skill_Type_Weapon:   "武器",
	Skill_Type_Dodge:    "轻功",
	Skill_Type_Parry:    "招架",
	Skill_Type_Throwing: "暗器",
}

type PerformData struct {
	Name  string `mapstructure:"名称" yaml:"名称,omitempty"`
	Mixed bool   `mapstructure:"融合" yaml:"融合,omitempty"`
	Level int
}

type PassiveData struct {
	Name      string    `mapstructure:"名称" yaml:"名称,omitempty"`
	Level     int       `mapstructure:"等级" yaml:"等级,omitempty"`
	Arguments []float64 `mapstructure:"参数" yaml:"参数,omitempty"`
}

type SkillData struct {
	Name        string        `mapstructure:"部位" yaml:"部位,omitempty"`
	Level       int           `mapstructure:"等级" yaml:"等级,omitempty"`
	PerformData []PerformData `mapstructure:"绝招" yaml:"绝招,omitempty"`
	PassiveData PassiveData   `mapstructure:"被动" yaml:"被动,omitempty"`
}

type Runable interface {
	Run(*RunContext)
}

type RunFunc func(*RunContext)

func (fn RunFunc) Run(ctx *RunContext) {
	if fn != nil {
		fn(ctx)
	}
}

// 武学管理器
type Skills struct {
	Owner     *Player
	Levels    [6]int
	Performs  []Perform
	InitOnlys []Perform
	Passives  *SkillPassives
}

func (mgr *Skills) GetLevel(part string) int {
	i, ok := SkillTypeValue[part]
	if ok {
		return mgr.Levels[i]
	}
	return 0
}

func (mgr *Skills) Update(diff int) {
	for _, p := range mgr.Performs {
		p.Update(diff)
	}
	mgr.Passives.Update(diff)
}

func (mgr *Skills) CanRun(target *Player, args ...interface{}) []Runable {
	out := []Runable{}
	for _, p := range mgr.Performs {
		if p.CanRun(target, args...) {
			out = append(out, p)
		}
	}
	return out
}

func (mgr *Skills) ResetAllPerform() {
	for _, e := range mgr.Performs {
		if e.GetName() == "太极真意.大道无极" ||
			e.GetName() == "慈航剑典.剑心通明" ||
			e.GetName() == "长生诀.天地诀" {
			continue
		}
		e.Reset()
	}
}

func NewSkills(owner *Player, datas ...SkillData) *Skills {
	skills := &Skills{
		Owner: owner,
	}

	for _, data := range datas {
		sType, ok := SkillTypeValue[data.Name]
		if !ok {
			log.Println("无效的武学部位:", data.Name)
			continue
		}
		skills.Levels[sType] = data.Level
		for _, pData := range data.PerformData {
			p := PerformRepo.Build(pData.Name, owner, data.Level, pData.Mixed)
			if p == nil {
				log.Println("不支持绝招:", pData.Name)
				continue
			}
			if p.IsInitOnly() {
				skills.InitOnlys = append(skills.InitOnlys, p)
			} else {
				skills.Performs = append(skills.Performs, p)
			}
		}
	}

	skills.Passives = NewSkillPassives(owner, datas...)

	return skills
}

// 绝招
type Perform interface {
	GetName() string
	GetType() string
	IsInitOnly() bool
	Update(int)
	IsReady() bool
	Reset()
	CanRun(*Player, ...interface{}) bool
	Run(*RunContext)
}

// 主动效果.基类
type BasePerform struct {
	Name        string
	Type        string
	InitOnly    bool
	Player      *Player
	Level       int
	CoolDown    int
	IgnoreFaint bool
	IgnoreBusy  bool
	MPCost      float64
	Mixed       bool
	Timer       *Timer
}

func (a *BasePerform) IsReady() bool {
	return a.Timer.IsDone()
}

func (a *BasePerform) GetName() string {
	return a.Name
}

func (a *BasePerform) IsInitOnly() bool {
	return a.InitOnly
}

func (a *BasePerform) GetType() string {
	return a.Type
}

func (a *BasePerform) Update(diff int) {
	a.Timer.Update(diff)
}

func (a *BasePerform) Reset() {
	a.Timer.Done()
}

func (a *BasePerform) CanRun(target *Player, args ...interface{}) bool {
	if a.Player == nil || !a.Player.IsAlive() {
		return false
	}
	if !a.Timer.IsDone() {
		return false
	}
	if !a.Player.PerformCD.IsDone() {
		return false
	}
	if a.Player.State.Faint && !a.IgnoreFaint {
		return false
	}
	if a.Player.State.Busy && !a.IgnoreBusy {
		return false
	}
	if target != nil && !target.IsAlive() {
		return false
	}
	if a.Type == "weapon" {
		if a.Player.Weapon == nil {
			return false
		}
		if !a.Player.Weapon.Wielded {
			return false
		}
	}
	if cost := CalcMPCostP(a.MPCost, a.Player); cost > a.Player.GetMP() {
		return false
	}
	return true
}

func (a *BasePerform) PreFlight(ctx *RunContext) {
	ctx.SetName(a.Name)

	if !a.CanRun(ctx.target) {
		return
	}
	var (
		attacker = ctx.attacker
		target   = ctx.target

		mp = CalcMPCostP(a.MPCost, a.Player)
		ct = CalcCastTimeP(a.Player.GetAttackSpeed(), a.Player)
		cd = CalcCoolDownP(a.CoolDown, a.Player)
	)
	// 武神传说三大bug绝招之移花接木
	if target != nil && target.HasBuff("移花接木.移花.绿") {
		target.RemoveBuff("移花接木.移花.绿")
		ctx.SetFail("移花接木.移花.绿")
	}
	// 失败也要给钱,概不赊账
	if a.Name == "长生诀.天地诀" {
		cd = CalcCoolDownPWithouts(a.CoolDown, attacker, "慈航剑典.剑心")
	}
	attacker.SubMP(mp)
	attacker.CastTime.Start(ct)
	a.Timer.Start(cd)

	ctx.SetMPCost(mp).
		SetCastTime(ct).
		SetCooldown(cd)
}

func (a *BasePerform) Run(ctx *RunContext) {
	a.PreFlight(ctx)
}

func (a *BasePerform) Attack(ctx *RunContext, mods ...Modifier) *CombatContext {
	actx := ctx.AddAttack().AddModifier(mods...)
	switch a.Type {
	case "unarmed":
		actx.SetUnarmedAttack()
	case "weapon":
		actx.SetWeaponAttack()
	case "throwing":
		actx.SetThrowingAttack()
	case "force":
		actx.SetForceAttack()
	}
	GenericAttack(actx)
	return actx
}

// perform repo
type PerformBuilder func(*Player, int, bool) Perform

var PerformRepo = &PerformRepository{make(map[string]PerformBuilder)}

type PerformRepository struct {
	builders map[string]PerformBuilder
}

func (repo *PerformRepository) Add(name string, builder PerformBuilder) {
	repo.builders[name] = builder
}

func (repo *PerformRepository) Build(name string, player *Player, level int, mixed bool) Perform {
	b, ok := repo.builders[name]
	if ok {
		return b(player, level, mixed)
	}
	return nil
}
