package simulator

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/abs0293/wsmud-simc/simulator/log_pb"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type PlayerState struct {
	Faint bool
	Busy  bool
	Blind bool
	Fixed bool
}

type PlayerData struct {
	Name          string          `mapstructure:"姓名" yaml:"姓名,omitempty"`
	HP            float64         `mapstructure:"气血d" yaml:"气血d,omitempty"`
	HPPct         float64         `mapstructure:"气血%" yaml:"气血%,omitempty"`
	MP            float64         `mapstructure:"内力d" yaml:"内力d,omitempty"`
	MP2HP         float64         `mapstructure:"内力转化%" yaml:"内力转化%,omitempty"`
	Str           float64         `mapstructure:"先天臂力d" yaml:"先天臂力d,omitempty"`
	StrAdd        float64         `mapstructure:"臂力d" yaml:"臂力d,omitempty"`
	Dex           float64         `mapstructure:"先天身法d" yaml:"先天身法d,omitempty"`
	DexAdd        float64         `mapstructure:"身法d" yaml:"身法d,omitempty"`
	Con           float64         `mapstructure:"先天根骨d" yaml:"先天根骨d,omitempty"`
	ConAdd        float64         `mapstructure:"根骨d" yaml:"根骨d,omitempty"`
	Int           float64         `mapstructure:"先天悟性d" yaml:"先天悟性d,omitempty"`
	IntAdd        float64         `mapstructure:"悟性d" yaml:"悟性d,omitempty"`
	Attack        float64         `mapstructure:"攻击d" yaml:"攻击d,omitempty"`
	AttackPct     float64         `mapstructure:"攻击%" yaml:"攻击%,omitempty"`
	Defence       float64         `mapstructure:"防御d" yaml:"防御d,omitempty"`
	DefencePct    float64         `mapstructure:"防御%" yaml:"防御%,omitempty"`
	Hit           float64         `mapstructure:"命中d" yaml:"命中d,omitempty"`
	HitPct        float64         `mapstructure:"命中%" yaml:"命中%,omitempty"`
	Dodge         float64         `mapstructure:"闪避d" yaml:"闪避d,omitempty"`
	DodgePct      float64         `mapstructure:"闪避%" yaml:"闪避%,omitempty"`
	Parry         float64         `mapstructure:"招架d" yaml:"招架d,omitempty"`
	ParryPct      float64         `mapstructure:"招架%" yaml:"招架%,omitempty"`
	Speed         float64         `mapstructure:"攻速d" yaml:"攻速d,omitempty"`
	SpeedPct      float64         `mapstructure:"攻速%" yaml:"攻速%,omitempty"`
	CDR           float64         `mapstructure:"绝招冷却d" yaml:"绝招冷却d,omitempty"`
	CDRPct        float64         `mapstructure:"绝招冷却%" yaml:"绝招冷却%,omitempty"`
	DmR           float64         `mapstructure:"免伤d" yaml:"免伤d,omitempty"`
	DmRPct        float64         `mapstructure:"免伤%" yaml:"免伤%,omitempty"`
	IgDPct        float64         `mapstructure:"忽防%" yaml:"忽防%,omitempty"`
	FiDPct        float64         `mapstructure:"终伤%" yaml:"终伤%,omitempty"`
	CTR           float64         `mapstructure:"绝招释放d" yaml:"绝招释放d,omitempty"`
	CTRPct        float64         `mapstructure:"绝招释放%" yaml:"绝招释放%,omitempty"`
	MPCR          float64         `mapstructure:"内力消耗d" yaml:"内力消耗d,omitempty"`
	MPCRPct       float64         `mapstructure:"内力消耗%" yaml:"内力消耗%,omitempty"`
	NRPct         float64         `mapstructure:"负面抵抗%" yaml:"负面抵抗%,omitempty"`
	WeaponData    WeaponData      `mapstructure:"武器" yaml:"武器,omitempty"`
	SkillData     []SkillData     `mapstructure:"武学" yaml:"武学,omitempty"`
	EquipmentData []EquipmentData `mapstructure:"装备" yaml:"装备,omitempty"`
}

type Player struct {
	PlayerData
	LostHP float64
	LostMP float64
	Roll   *rand.Rand
	Weapon *Weapon
	// 出招时间
	CastTime *Timer
	// 普通攻击CD
	AttackCD *Timer
	// 行动CD，等于人物攻速
	ActionCD *Timer
	// 禁止绝招
	PerformCD *Timer
	// 状态
	State      *PlayerState
	Arena      *Arena
	Target     *Player
	Buffs      *Buffs
	Skills     *Skills
	Equipments *Equipments
}

func (p Player) GetConfig() string {
	out, _ := yaml.Marshal(&p.PlayerData)
	return string(out)
}

func (p *Player) GetAttack(mods ...[]Modifier) float64 {
	atk := p.GetStr() + (p.GetStr()/10)*p.GetStrAdd()
	m1 := p.GetModifier("攻击d", mods...)
	m2 := p.GetModifier("攻击%", mods...)
	return (p.Attack + atk + m1) * (1 + p.AttackPct + m2)
}

func (p *Player) GetDefence() float64 {
	m1 := p.GetModifier("防御d")
	m2 := p.GetModifier("防御%")
	return (p.Defence + m1) * (1 + m2 + p.DefencePct)
}

func (p *Player) GetDodge(mods ...[]Modifier) float64 {
	rate := 1 + p.Skills.Passives.DodgeZhuanZhu.GetRate()*2
	dodge := p.GetDex()/2 + (p.GetDex()/10)*p.GetDexAdd()*rate
	m1 := p.GetModifier("闪避d", mods...)
	m2 := p.GetModifier("闪避%", mods...)
	return (p.Dodge + dodge + m1) * (1 + p.DodgePct + m2)
}

func (p *Player) GetHit(mods ...[]Modifier) float64 {
	hit := p.Dex / 2
	m1 := p.GetModifier("命中d", mods...)
	m2 := p.GetModifier("命中%", mods...)
	return (p.Hit + hit + m1) * (1 + p.HitPct + m2)
}

func (p *Player) GetParry(mods ...[]Modifier) float64 {
	rate := 1 + p.Skills.Passives.ParryZhuanZhu.GetRate()
	parry := p.GetStr()/2 + (p.GetStr()/10)*p.GetStrAdd()*rate
	m1 := p.GetModifier("招架d", mods...)
	m2 := p.GetModifier("招架%", mods...)
	return (p.Parry + parry + m1) * (1 + p.ParryPct + m2)
}

func (p *Player) GetCoolDownReduce() float64 {
	return p.CDR + p.GetModifier("绝招冷却d")
}

func (p *Player) GetCoolDownReducePercent() float64 {
	return p.CDRPct + p.GetModifier("绝招冷却%")
}

func (p *Player) GetCoolDownReduceExclude(exclude ...string) float64 {
	return p.CDR + p.GetModifierExclude("绝招冷却d", exclude...)
}

func (p *Player) GetCoolDownReducePercentExclude(exclude ...string) float64 {
	return p.CDRPct + p.GetModifierExclude("绝招冷却%", exclude...)
}

func (p *Player) GetDamageReduce() float64 {
	return p.DmR + p.GetModifier("免伤d")
}

func (p *Player) GetDamageReducePercent() float64 {
	return CalcStaticDamageReducePercent(p.DmRPct+p.GetModifier("静态免伤%")) + p.GetModifier("免伤%")
}

// 未衰减静态免伤
func (p *Player) GetRawStaticDamageReducePercent() float64 {
	return p.DmRPct + p.GetModifier("静态免伤%")
}

func (p *Player) GetIgnoreDefencePercent(mods ...[]Modifier) float64 {
	return p.IgDPct + p.GetModifier("忽防%", mods...)
}

func (p *Player) GetFinalDamagePercent(mods ...[]Modifier) float64 {
	return p.FiDPct + p.GetModifier("终伤%", mods...)
}

func (p *Player) GetMPMax() float64 {
	m1 := p.GetModifier("内力d")
	m2 := p.GetModifier("内力%")
	return (p.MP + m1) * (1 + m2)
}

func (p *Player) GetMP() float64 {
	return p.GetMPMax() - p.LostMP
}

func (p *Player) SubMP(v float64) float64 {
	if v <= 0 {
		return 0
	}
	p.LostMP += v
	max := p.GetMPMax()
	over := math.Max(p.LostMP-max, 0)
	p.LostMP = math.Min(max, p.LostMP)

	return over
}

func (p *Player) AddMP(v float64) float64 {
	if v <= 0 {
		return 0
	}

	p.LostMP += v
	if o := p.LostMP - p.GetMPMax(); o > 0 {
		v = o // 过量
		p.LostMP -= o
	}
	return v
}

func (p *Player) GetHPMax() float64 {
	// 内力转化
	mp2hp := p.GetMPMax() * p.MP2HP
	// 根骨
	con2hp := p.Con * (5 + p.ConAdd)
	m1 := p.GetModifier("气血d")
	m2 := p.GetModifier("气血%")
	return (p.HP + mp2hp + con2hp + m1) * (1 + p.HPPct + m2) * 2 //擂台血量翻倍
}

func (p *Player) GetHP() float64 {
	return p.GetHPMax() - p.LostHP
}

func (p *Player) AddHP(v float64) *Player {
	p.LostHP -= v
	if p.LostHP < 0 {
		p.LostHP = 0
	}
	if max := p.GetHPMax(); max < p.LostHP {
		p.LostHP = max
	}
	return p
}

func (p *Player) GetWeaponType() int {
	if !p.Weapon.Wielded {
		return 0
	}
	return p.Weapon.GetType()
}

func (p *Player) Unwield() {
	p.Weapon.Unwield()
}

func (p *Player) GetStr() float64 {
	return p.Str
}

func (p *Player) GetStrAdd() float64 {
	m1 := p.GetModifier("臂力d")
	m2 := p.GetModifier("臂力%")
	return (p.StrAdd + m1) * (1 + m2)
}

func (p *Player) GetDex() float64 {
	return p.Dex
}

func (p *Player) GetDexAdd() float64 {
	m1 := p.GetModifier("身法d")
	m2 := p.GetModifier("身法%")
	return (p.DexAdd + m1) * (1 + m2)
}

func (p *Player) GetAttackSpeed() int {
	m1 := p.GetModifier("攻速d")
	m2 := p.GetModifier("攻速%")
	v := (4 - p.GetDex()/25 - m1 - p.Speed) * (1 - m2 - p.SpeedPct)
	if v < 0.5 {
		v = 0.5
	}
	return Sec2Ms(v)
}

func (p *Player) GetNegativeResistPercent() float64 {
	m1 := p.GetModifier("负面抵抗%")
	return p.NRPct + m1
}

func (p *Player) GetCastTimeReduce() float64 {
	m1 := p.GetModifier("绝招释放d")
	return p.CTR + m1
}

func (p *Player) GetCastTimeReducePercent() float64 {
	m1 := p.GetModifier("绝招释放%")
	return p.CTRPct + m1
}

func (p *Player) GetMPCostReduce() float64 {
	return p.MPCR + p.GetModifier("内力消耗d")
}

func (p *Player) GetMPCostReducePercent() float64 {
	return p.MPCRPct + p.GetModifier("内力消耗%")
}

func (p *Player) IsAlive() bool {
	return p.GetHP() > 0
}

func (p *Player) IsImmuneDisarm() bool {
	return p.GetModifier("免疫缴械d") > 0
}

func (p *Player) TakeDamage(t *Player, d float64) float64 {
	p.LostHP += d
	return d
}

func (p *Player) Update(diff int) {
	p.updateBuff(diff)
	p.updateCD(diff)
}

func (p *Player) OnCombatStart() {
	p.Weapon.Wield()
	if pDuGu := p.Skills.Passives.DuGu; pDuGu != nil {
		if p.GetWeaponType() == WeaponTypeSword {
			p.AddBuff(BuffRepo.Build("独孤剑诀.无我", p, pDuGu.level))
			p.AddBuff(BuffRepo.Build("独孤剑诀.无剑", p, pDuGu.level))
		}
	}
	for _, i := range p.Skills.InitOnlys {
		ctx := NewRunContext(p, p.Target, p.Arena.Ticks)
		ctx.SetName(i.GetName())
		i.Run(ctx)
		ctx.Done()
	}
}

func (p *Player) DebugPerform(name string) {
	for _, p := range p.Skills.Performs {
		if p.GetName() == name {
			p.Run(nil)
			return
		}
	}
	for _, a := range p.Equipments.Equipments {
		if a.GetName() == name {
			a.Run(nil)
			return
		}
	}
}

func (p *Player) Action() {
	if !p.ActionCD.IsDone() {
		return
	}

	done := false

	if p.State.Faint || p.State.Busy {
		return
	}

	if p.CastTime.IsDone() && p.Roll.Float64() >= 0.3 {
		can := []Runable{}
		can = append(can, p.Skills.CanRun(p.Target)...)
		can = append(can, p.Equipments.CanRun(p.Target)...)
		// 重新装备武器
		if p.Weapon != nil && !p.Weapon.Wielded {
			can = append(can, p.Weapon.Rewield())
		}
		if len(can) > 0 {
			ctx := NewRunContext(p, p.Target, p.Arena.Ticks)
			r := can[p.Roll.Intn(len(can))]
			r.Run(ctx)
			ctx.Done()
		}
		done = true
	}

	if p.CastTime.IsDone() && p.AttackCD.IsDone() {
		ctx := NewCombatContext(p, p.Target, p.Arena.Ticks)
		ctx.SetMainAttack()
		ctx.SetSource("普通攻击")
		if p.GetWeaponType() > 1 {
			ctx.SetWeaponAttack()
		} else if p.GetWeaponType() == 0 {
			ctx.SetUnarmedAttack()
		}
		GenericAttack(ctx)
		// 额外攻击
		if pDuGu := p.Skills.Passives.DuGu; pDuGu != nil &&
			p.HasBuff("独孤剑诀.剑来") &&
			ctx.IsWeaponAttack() {
			pDuGu.DoExtraHit(ctx.AddExtraAttack())
		}
		if pJianXin := p.Skills.Passives.JianXin; pJianXin != nil {
			pJianXin.DoExtraHit(ctx)
		}
		if pFuYu := p.Skills.Passives.FuYu; pFuYu != nil && p.HasBuff("覆雨剑法.剑雨") {
			pFuYu.DoExtraHit(ctx)
		}
		ctx.Done()
		// 目标招架/闪避触发
		if ctx.IsParry() && p.Target.IsAlive() {
			if tWuJi := p.Target.Skills.Passives.WuJi; tWuJi != nil {
				rctx := tWuJi.Run(p)
				if rctx != nil {
					rctx.Done()
				}
			}
		}
		done = true
	}

	if done {
		p.ActionCD.Start(p.GetAttackSpeed())
	}
}

func (p *Player) GetModifier(t string, extMods ...[]Modifier) float64 {
	return p.Buffs.GetModifier(t) + GetModifier(t, extMods...)
}

func (p *Player) GetModifierExclude(t string, exclude ...string) float64 {
	return p.Buffs.GetModifierExclude(t, exclude...)
}

func (p *Player) GetBuff(name string) Buff {
	return p.Buffs.GetBuff(name)
}

func (p *Player) AddBuff(b Buff) {
	if b == nil {
		log.Panic("AddBuff: buff is nil")
	}
	p.Buffs.Add(b)
}

func (p *Player) RemoveBuff(name string) {
	p.Buffs.Remove(name)
}

func (p *Player) RemoveBuffByType(typ string) {
	p.Buffs.RemoveByType(typ)
}

func (p *Player) HasBuff(n string) bool {
	return p.GetBuff(n) != nil
}

func (p *Player) GetEquipment(name string) Equipment {
	for _, e := range p.Equipments.Equipments {
		if e.GetName() == name {
			return e
		}
	}
	return nil
}

func (p *Player) Log(logs ...*log_pb.Log) {
	p.Arena.Log(logs...)
}

func (p *Player) updateCD(diff int) {
	p.CastTime.Update(diff)
	p.PerformCD.Update(diff)
	p.AttackCD.Update(diff)
	p.Skills.Update(diff)
	p.ActionCD.Update(diff)
}

func (p *Player) updateBuff(diff int) {
	p.Buffs.Update(diff)
}

func (p *Player) Printf(format string, args ...interface{}) {
	if Silent {
		return
	}
	var prefix string
	if p.Arena != nil {
		prefix = fmt.Sprintf("[%08.3f][%s]", p.Arena.Timestamp(), p.Name)
	}
	fmt.Printf(prefix+format, args...)
}

func readConfig(fn string) (*viper.Viper, error) {
	cfg := viper.New()
	cfg.AddConfigPath(".")
	cfg.SetConfigFile(fn)
	return cfg, cfg.ReadInConfig()
}

func NewPlayer(data PlayerData) *Player {
	player := &Player{
		PlayerData: data,
		Roll:       rand.New(rand.NewSource(time.Now().UnixNano())),
		ActionCD:   NewTimer(),
		CastTime:   NewTimer(),
		PerformCD:  NewTimer(),
		AttackCD:   NewTimer(),
		State:      &PlayerState{},
	}

	player.Buffs = NewBuffs(player)
	player.Skills = NewSkills(player, data.SkillData...)
	player.Equipments = NewEquipments(player, data.EquipmentData...)
	player.Weapon = NewWeapon(player, data.WeaponData)

	return player
}

func ReadPlayerDataFromFile(fn string) (PlayerData, error) {
	pData := PlayerData{}

	cfg, err := readConfig(fn)
	if err != nil {
		return pData, err
	}

	if err := cfg.Unmarshal(&pData); err != nil {
		return pData, err
	}

	return pData, nil
}

func NewPlayerFromFile(fn string) (*Player, error) {
	pData, err := ReadPlayerDataFromFile(fn)
	if err != nil {
		return nil, err
	}

	p := NewPlayer(pData)
	return p, nil
}
