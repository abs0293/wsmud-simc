package simulator

import (
	"log"
	"math"
)

// 被动效果
type SkillPassives struct {
	FanZhen       *Passive_ZiChuang_FanZhen
	JianXin       *Passive_ZiChuang_JianXin
	ZhanShen      *Passive_ZiChuang_ZhanShen
	BuMie         *Passive_ZiChuang_BuMie
	DodgeZhuanZhu *Passive_ZiChuang_Dodge_ZhuanZhu
	ParryZhuanZhu *Passive_ZiChuang_Parry_ZhuanZhu
	RuMo          *Passive_ZiChuang_Rumo
	WanDaoRuMo    *Passive_YuanYueWanDao_Rumo
	FuYu          *Passive_FuYuJianFa_FuYu
	DuGu          *Passive_DuGuJianJue_DuGu
}

func (mgr *SkillPassives) Load(player *Player, sdata SkillData) {
	data := sdata.PassiveData
	if data.Name == "" {
		return
	}
	switch data.Name {
	case "自创.反震":
		mgr.FanZhen = Passive_ZiChuang_FanZhen_Builder(player, data)
	case "自创.剑心":
		mgr.JianXin = Passive_ZiChuang_JianXin_Builder(player, data)
	case "自创.战神":
		mgr.ZhanShen = Passive_ZiChuang_ZhanShen_Builder(player, data)
	case "自创.不灭":
		mgr.BuMie = Passive_ZiChuang_BuMie_Builder(player, data)
	case "自创.专注(轻功)":
		mgr.DodgeZhuanZhu = Passive_ZiChuang_Dodge_ZhuanZhu_Builder(player, data)
	case "自创.专注(招架)":
		mgr.ParryZhuanZhu = Passive_ZiChuang_Parry_ZhuanZhu_Builder(player, data)
	case "自创.入魔":
		isWeapon := true
		if sdata.Name == "拳脚" {
			isWeapon = false
		}
		mgr.RuMo = Passive_ZiChuang_Rumo_Builder(player, isWeapon, data)
	case "圆月弯刀.入魔":
		mgr.WanDaoRuMo = Passive_YuanYueWanDao_Rumo_Builder(player, data)
	case "覆雨剑法.覆雨":
		mgr.FuYu = Passive_FuYuJianFa_FuYu_Builder(player, data)
	case "独孤剑诀.独孤":
		mgr.DuGu = Passive_DuGuJianJue_DuGu_Builder(player, data)
	default:
		log.Println("不支持被动:", data.Name)
	}
}

func (p *SkillPassives) Update(diff int) {
	p.BuMie.Update(diff)
}

func NewSkillPassives(player *Player, datas ...SkillData) *SkillPassives {
	passives := &SkillPassives{}
	for _, sData := range datas {
		passives.Load(player, sData)
	}
	return passives
}

// 反震
type Passive_ZiChuang_FanZhen struct {
	player *Player
	rate   float64
}

func (p *Passive_ZiChuang_FanZhen) Damage() float64 {
	return p.player.GetMPMax() * p.rate
}

func Passive_ZiChuang_FanZhen_Builder(player *Player, data PassiveData) *Passive_ZiChuang_FanZhen {
	return &Passive_ZiChuang_FanZhen{
		player: player,
		rate:   0.005 + float64(data.Level/2)*0.001,
	}
}

// 剑心
type Passive_ZiChuang_JianXin struct {
	player    *Player
	hitNumber int
}

func (p *Passive_ZiChuang_JianXin) HitNumber() int {
	return p.hitNumber
}

func (p *Passive_ZiChuang_JianXin) DoExtraHit(trigger *CombatContext) {
	var (
		attacker = trigger.GetAttacker()
		target   = trigger.GetTarget()
	)

	for i := 0; i < p.HitNumber(); i++ {
		if !attacker.IsAlive() || !target.IsAlive() {
			break
		}
		ctx := trigger.AddExtraAttack()
		if attacker.GetWeaponType() == 1 {
			ctx.SetUnarmedAttack()
		} else {
			ctx.SetWeaponAttack()
		}
		ctx.SetSource("自创.剑心")
		GenericAttack(ctx)
	}
}

func Passive_ZiChuang_JianXin_Builder(player *Player, data PassiveData) *Passive_ZiChuang_JianXin {
	return &Passive_ZiChuang_JianXin{
		player:    player,
		hitNumber: 1 + data.Level/30,
	}
}

// 战神
type Passive_ZiChuang_ZhanShen struct {
	player   *Player
	appended float64
	leech    float64
}

func (p *Passive_ZiChuang_ZhanShen) Append() float64 {
	rate := p.appended
	if p.player.GetWeaponType() == 0 {
		rate *= 2
	}
	return p.player.GetMPMax() * rate
}

func (p *Passive_ZiChuang_ZhanShen) Leech(damage float64) float64 {
	rate := p.leech
	if p.player.GetWeaponType() == 0 {
		rate *= 2
	}
	return damage * rate
}

func Passive_ZiChuang_ZhanShen_Builder(player *Player, data PassiveData) *Passive_ZiChuang_ZhanShen {
	return &Passive_ZiChuang_ZhanShen{
		player:   player,
		appended: 0.01 + float64(data.Level/3)*0.001,
		leech:    0.2,
	}
}

// 不灭
type Passive_ZiChuang_BuMie struct {
	player   *Player
	activate bool
	cd       *Timer
	rate     float64
}

func (p *Passive_ZiChuang_BuMie) IsActivate() bool {
	return p.activate
}

func (p *Passive_ZiChuang_BuMie) Absorb(d float64) (float64, float64) {
	if !p.player.HasBuff("长生诀.混沌") {
		return d, 0
	}
	v := p.player.GetHPMax() * p.rate
	if d < v {
		return d, 0
	}
	return v, math.Min(d-v, p.player.LostHP)
}

func (p *Passive_ZiChuang_BuMie) Activate(dmg float64) bool {
	if p.cd.IsDone() &&
		((p.player.GetHP() - dmg) < p.player.GetHPMax()/10) {
		p.cd.Start(60000)
		p.activate = true
		p.player.LostHP = 0
		return true
	}
	return false
}

func (p *Passive_ZiChuang_BuMie) Update(diff int) {
	if p == nil {
		return
	}
	p.cd.Update(diff)
}

func Passive_ZiChuang_BuMie_Builder(player *Player, data PassiveData) *Passive_ZiChuang_BuMie {
	return &Passive_ZiChuang_BuMie{
		player: player,
		rate:   0.02,
	}
}

// 轻功专注
type Passive_ZiChuang_Dodge_ZhuanZhu struct {
	player *Player
	rate   float64
}

func (p *Passive_ZiChuang_Dodge_ZhuanZhu) GetRate() float64 {
	if p == nil {
		return 0
	}
	return p.rate
}

func Passive_ZiChuang_Dodge_ZhuanZhu_Builder(player *Player, data PassiveData) *Passive_ZiChuang_Dodge_ZhuanZhu {
	return &Passive_ZiChuang_Dodge_ZhuanZhu{
		player: player,
		rate:   0.48 + float64(data.Level)*0.02,
	}
}

// 招架专注
type Passive_ZiChuang_Parry_ZhuanZhu struct {
	player *Player
	rate   float64
}

func (p *Passive_ZiChuang_Parry_ZhuanZhu) GetRate() float64 {
	if p == nil {
		return 0
	}
	return p.rate
}

func Passive_ZiChuang_Parry_ZhuanZhu_Builder(player *Player, data PassiveData) *Passive_ZiChuang_Parry_ZhuanZhu {
	return &Passive_ZiChuang_Parry_ZhuanZhu{
		player: player,
		rate:   0.48 + float64(data.Level)*0.02,
	}
}

// 自创.入魔
type Passive_ZiChuang_Rumo struct {
	player *Player
	weapon bool
	rate   float64
}

func (p *Passive_ZiChuang_Rumo) IsWeapon() bool {
	return p.weapon
}

func (p *Passive_ZiChuang_Rumo) DamageAdd() float64 {
	if p == nil {
		return 0
	}
	m := p.player.GetModifier("入魔.内力附加%")
	v := p.player.GetMPMax() * p.rate * (1 + m)
	cost := CalcMPCostP(v, p.player)
	if p.player.SubMP(cost) > 0 {
		return 0
	}
	return v
}

func Passive_ZiChuang_Rumo_Builder(player *Player, weapon bool, data PassiveData) *Passive_ZiChuang_Rumo {
	return &Passive_ZiChuang_Rumo{
		player: player,
		weapon: weapon,
		rate:   0.001 + float64(data.Level/3)*0.001,
	}
}

// 圆月弯刀.入魔
type Passive_YuanYueWanDao_Rumo struct {
	player *Player
	rate   float64
}

func (p *Passive_YuanYueWanDao_Rumo) DamageAdd(weak bool) float64 {
	if p == nil {
		return 0
	}
	m := p.player.GetModifier("入魔.内力附加%")
	v := p.player.GetMPMax() * p.rate * (1 + m)
	if p.player.SubMP(v) > 0 {
		if weak {
			p.player.AddBuff(BuffRepo.Build("圆月弯刀.虚弱", p.player, 0))
		}
		return 0
	}
	return v
}

func Passive_YuanYueWanDao_Rumo_Builder(player *Player, data PassiveData) *Passive_YuanYueWanDao_Rumo {
	return &Passive_YuanYueWanDao_Rumo{
		player: player,
		rate:   float64(data.Level/1000) / 100,
	}
}

// 覆雨剑法.覆雨
type Passive_FuYuJianFa_FuYu struct {
	player *Player
	level  int
}

func (p *Passive_FuYuJianFa_FuYu) GetChance() float64 {
	return float64(p.level/1000) / 20
}

func (p *Passive_FuYuJianFa_FuYu) GetHitNumber() int {
	return p.level/1000 + 1
}

func (p *Passive_FuYuJianFa_FuYu) DoExtraHit(trigger *CombatContext) {
	var (
		attacker = trigger.GetAttacker()
		target   = trigger.GetTarget()
	)

	if attacker.GetWeaponType() == 1 {
		return
	}

	for i := 0; i < p.GetHitNumber(); i++ {
		if !attacker.IsAlive() || !target.IsAlive() {
			break
		}
		ctx := trigger.AddExtraAttack()
		ctx.SetWeaponAttack()
		ctx.SetSource("覆雨剑法.剑雨")
		GenericAttack(ctx)
		attacker.Log(ctx.log)
	}
}

func Passive_FuYuJianFa_FuYu_Builder(player *Player, data PassiveData) *Passive_FuYuJianFa_FuYu {
	return &Passive_FuYuJianFa_FuYu{
		player: player,
		level:  data.Level,
	}
}

// 独孤剑诀.独孤
type Passive_DuGuJianJue_DuGu struct {
	player *Player
	level  int
}

func (p *Passive_DuGuJianJue_DuGu) Damage() float64 {
	return p.player.GetAttack() * (0.5 + float64(p.level)/1000*0.1)
}

func (p *Passive_DuGuJianJue_DuGu) DoExtraHit(ctx *CombatContext) {
	ProcessDamageApply(
		ctx.SetHit().
			SetWeaponAttack().
			SetDamageFinal(p.Damage()).
			SetTrueDamage(),
	)
}

func Passive_DuGuJianJue_DuGu_Builder(player *Player, data PassiveData) *Passive_DuGuJianJue_DuGu {
	return &Passive_DuGuJianJue_DuGu{
		player: player,
		level:  data.Level,
	}
}
