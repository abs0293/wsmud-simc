package simulator

// 绝招:五毒钩法.金钩锁魂
//
//	机制:
//	1. 你的钩法等级/2大于对方特殊招架武功等级, 必成功
//	2. 对方特殊招架武功等级大于你钩法等级的7/6, 必失败
//	3. 钩法等级处于[0.85*目标特殊招架等级, 2*目标特殊招架等级]区间时:
//		2000vs2000 下武器成功:181次，失败534次, 约25%
//		3000vs2000 下武器成功:459次，失败140次, 约75%
//		3500vs2000 下武器成功:629次，失败74次, 约90%
//	4. 目标没有武器时,不会攻击4次
type Perform_WuDuGouFa_JinGouSuoHun struct {
	BasePerform
}

func (pfm *Perform_WuDuGouFa_JinGouSuoHun) Run(ctx *RunContext) {
	pfm.PreFlight(ctx)
	if ctx.IsFail() {
		return
	}

	var (
		target = ctx.target
	)

	if !target.Weapon.Wielded {
		ctx.SetFail("空手")
		return
	}

	if target.IsImmuneDisarm() {
		ctx.SetFail("免疫缴械")
		return
	}

	attacker := pfm.Player
	atkLvl := float64(pfm.Level)
	tarLvl := float64(target.Skills.GetLevel("招架"))

	if atkLvl*7/6 < tarLvl {
		return
	}

	if atkLvl < tarLvl*2 && attacker.Roll.Float64() > 0.25 {
		return
	}

	target.Unwield()
	for i := 0; i < 4; i++ {
		pfm.Attack(ctx)
	}
}

func Perform_WuDuGouFa_JinGouSuoHun_Builder(player *Player, level int, mixed bool) Perform {
	return &Perform_WuDuGouFa_JinGouSuoHun{
		BasePerform: BasePerform{
			Name:     "五毒钩法.金钩锁魂",
			Type:     "weapon",
			Player:   player,
			Level:    level,
			Mixed:    mixed,
			CoolDown: 25000,
			Timer:    NewTimer(),
		},
	}
}

func init() {
	PerformRepo.Add("五毒钩法.金钩锁魂", Perform_WuDuGouFa_JinGouSuoHun_Builder)
}
