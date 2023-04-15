package simulator

// 绝招:圆月弯刀.入魔
type Perform_YuanYueWanDao_RuMo struct {
	BasePerform
}

func (pfm *Perform_YuanYueWanDao_RuMo) Run(ctx *RunContext) {
	pfm.PreFlight(ctx)
	if ctx.IsFail() {
		return
	}
	pfm.Player.AddBuff(BuffRepo.Build("圆月弯刀.入魔", pfm.Player, pfm.Level))
}

func Perform_YuanYueWanDao_RuMo_Builder(player *Player, level int, mixed bool) Perform {
	return &Perform_YuanYueWanDao_RuMo{
		BasePerform: BasePerform{
			Name:     "圆月弯刀.入魔",
			Type:     "weapon",
			Player:   player,
			Level:    level,
			Mixed:    mixed,
			CoolDown: 60000,
			Timer:    NewTimer(),
		},
	}
}

// 绝招:圆月弯刀.魔刀
//
//	机制:
//		1.当前内力记为MPc,被动消耗记为MPp,魔刀伤害记为D,D=MPc
//		2.施放魔刀,内力MP1=MPc-MPp:
//			1)如果MP1>=0,D+=被动伤害
//			2)如果MP1<0,MP1=0,不触发虚弱
//		3.造成伤害D
//		4.回复内力MPr=D*2或者D*(等级/3000),不确定,待测
//		5.扣除内力MP2=MPc
//		6.最后MPf=MPr+MP1-MP2
//		7.实战数据:
//			内力: 1958796/20009166
//			弯刀: 6000级, 入魔被动6%
//			伤害: 1958796+20009166*0.06=3159345
//			魔刀后内力: 3159345*2+(1958796-20009166*0.06)-1958796 = 5118140
type Perform_YuanYueWanDao_MoDao struct {
	BasePerform
}

func (pfm *Perform_YuanYueWanDao_MoDao) Run(ctx *RunContext) {
	pfm.PreFlight(ctx)
	if ctx.IsFail() {
		return
	}

	var (
		attacker = ctx.attacker
	)

	dmg := pfm.Player.GetMP()
	cost := CalcMPCostP(dmg, attacker)
	pfm.Attack(ctx, Modifier{"绝招.伤害倍率%", 0}, Modifier{"绝招.伤害附加d", dmg})
	attacker.AddMP(dmg*2 - cost)
}

func Perform_YuanYueWanDao_MoDao_Builder(player *Player, level int, mixed bool) Perform {
	return &Perform_YuanYueWanDao_MoDao{
		BasePerform: BasePerform{
			Name:     "圆月弯刀.魔刀",
			Type:     "weapon",
			Player:   player,
			Level:    level,
			Mixed:    mixed,
			CoolDown: 60000,
			Timer:    NewTimer(),
		},
	}
}

// 光环:圆月弯刀.入魔
type Buff_YuanYueWanDao_RuMo struct {
	BaseBuff
}

func (b *Buff_YuanYueWanDao_RuMo) GetModifier(n string) float64 {
	return b.BaseBuff.GetModifier(n)
}

func Buff_YuanYueWanDao_RuMo_Builder(player *Player, level int, args ...interface{}) Buff {
	mod := float64(level) * 0.0002
	return &Buff_YuanYueWanDao_RuMo{
		BaseBuff{
			Name:     "圆月弯刀.入魔",
			Type:     "weapon",
			Creator:  player,
			Duration: 3000 + Sec2Ms(float64(level/1000)),
			Modifiers: []Modifier{
				{"入魔.内力附加%", mod},
				{"命中%", mod},
			},
		},
	}
}

// 光环:圆月弯刀.虚弱
type Buff_YuanYueWanDao_Weak struct {
	BaseBuff
}

func (b *Buff_YuanYueWanDao_Weak) OnEnable() {
	// 成功虚弱回复10%最大内力,不要问为什么,实际如此
	b.Owner.AddMP(b.Owner.GetMPMax() * 0.1)
}

func Buff_YuanYueWanDao_Weak_Builder(player *Player, level int, args ...interface{}) Buff {
	return &Buff_YuanYueWanDao_Weak{
		BaseBuff{
			Name:         "圆月弯刀.虚弱",
			Type:         "weak",
			Creator:      player,
			Duration:     6000,
			Irresistible: true,
			Debuff:       true,
			Steady:       true,
			Modifiers: []Modifier{
				{"攻击%", -0.5},
				{"命中%", -0.5},
				{"招架%", -0.5},
				{"闪避%", -0.5},
				{"防御%", -0.5},
			},
		},
	}
}

func init() {
	PerformRepo.Add("圆月弯刀.入魔", Perform_YuanYueWanDao_RuMo_Builder)
	PerformRepo.Add("圆月弯刀.魔刀", Perform_YuanYueWanDao_MoDao_Builder)
	BuffRepo.Add("圆月弯刀.入魔", Buff_YuanYueWanDao_RuMo_Builder)
	BuffRepo.Add("圆月弯刀.虚弱", Buff_YuanYueWanDao_Weak_Builder)
}
