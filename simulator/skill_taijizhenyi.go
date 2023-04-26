package simulator

// 绝招:太极真义.阴阳无极
type Perform_TaiJiZhenYi_YinYangWuJi struct {
	BasePerform
}

func (pfm *Perform_TaiJiZhenYi_YinYangWuJi) Run(ctx *RunContext) {
	pfm.PreFlight(ctx)
	if ctx.IsFail() {
		return
	}
	var (
		target = ctx.GetTarget()
		cdAdd  = 1000 + 245*(pfm.Level/300)
	)

	target.AddBuff(BuffRepo.Build("太极真义.迟滞", pfm.Player, pfm.Level))
	for _, p := range target.Skills.Performs {
		if !p.IsReady() {
			p.AddCoolDown(cdAdd)
		}
	}
}

func Perform_TaiJiZhenYi_YinYangWuJi_Builder(player *Player, level int, mixed bool) Perform {
	return &Perform_TaiJiZhenYi_YinYangWuJi{
		BasePerform: BasePerform{
			Name:     "太极真义.阴阳无极",
			Type:     "parry",
			Player:   player,
			Level:    level,
			Mixed:    mixed,
			CoolDown: 55000,
			Timer:    NewTimer(),
		},
	}
}

// 绝招:太极真义.大道无极
type Perform_TaiJiZhenYi_DaDaoWuJi struct {
	BasePerform
}

func (pfm *Perform_TaiJiZhenYi_DaDaoWuJi) Run(ctx *RunContext) {
	pfm.PreFlight(ctx)
	if ctx.IsFail() {
		return
	}

	var (
		attacker = ctx.GetAttacker()
	)
	attacker.Buffs.ClearDebuffs()
	attacker.LostHP = 0
	attacker.LostMP = 0
}

func Perform_TaiJiZhenYi_DaDaoWuJi_Builder(player *Player, level int, mixed bool) Perform {
	return &Perform_TaiJiZhenYi_DaDaoWuJi{
		BasePerform: BasePerform{
			Name:     "太极真义.大道无极",
			Type:     "parry",
			Player:   player,
			Level:    level,
			Mixed:    mixed,
			CoolDown: 45000,
			Timer:    NewTimer(),
		},
	}
}

// 光环:太极真义.迟滞
type Buff_TaiJiZhenYi_ChiZhi struct {
	BaseBuff
}

func Buff_TaiJiZhenYi_ChiZhi_Builder(player *Player, level int, args ...interface{}) Buff {
	mod := 0.3 + float64(level/300)*0.01
	return &Buff_TaiJiZhenYi_ChiZhi{
		BaseBuff{
			Name:     "太极真义.迟滞",
			Type:     "chizhi",
			Debuff:   true,
			Creator:  player,
			Duration: 15000,
			Modifiers: []Modifier{
				{"攻速%", -mod},
			},
		},
	}
}

func init() {
	PerformRepo.Add("太极真义.阴阳无极", Perform_TaiJiZhenYi_YinYangWuJi_Builder)
	PerformRepo.Add("太极真义.大道无极", Perform_TaiJiZhenYi_DaDaoWuJi_Builder)
	BuffRepo.Add("太极真义.迟滞", Buff_TaiJiZhenYi_ChiZhi_Builder)
}
