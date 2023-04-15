package simulator

// 绝招:长生诀.混沌诀
type Perform_ChangShengJue_HunDunJue struct {
	BasePerform
}

func (pfm *Perform_ChangShengJue_HunDunJue) Run(ctx *RunContext) {
	pfm.PreFlight(ctx)
	if ctx.IsFail() {
		return
	}
	pfm.Player.AddBuff(BuffRepo.Build("长生诀.混沌", pfm.Player, pfm.Level))
}

func Perform_ChangShengJue_HunDunJue_Builder(player *Player, level int, mixed bool) Perform {
	return &Perform_ChangShengJue_HunDunJue{
		BasePerform: BasePerform{
			Name:     "长生诀.混沌诀",
			Type:     "force",
			Player:   player,
			Level:    level,
			Mixed:    mixed,
			CoolDown: 60000,
			Timer:    NewTimer(),
		},
	}
}

// 绝招:长生诀.天地诀
type Perform_ChangShengJue_TianDiJue struct {
	BasePerform
}

func (pfm *Perform_ChangShengJue_TianDiJue) Run(ctx *RunContext) {
	pfm.PreFlight(ctx)
	if ctx.IsFail() {
		return
	}
	pfm.Player.Buffs.ClearDebuffs()
	pfm.Player.LostHP = 0
	pfm.Player.Skills.ResetAllPerform()
}

func Perform_ChangShengJue_TianDiJue_Builder(player *Player, level int, mixed bool) Perform {
	return &Perform_ChangShengJue_TianDiJue{
		BasePerform: BasePerform{
			Name:     "长生诀.天地诀",
			Type:     "force",
			Player:   player,
			Level:    level,
			Mixed:    mixed,
			CoolDown: 60000,
			Timer:    NewTimer(),
		},
	}
}

// 光环:长生诀.混沌
type Buff_ChangShengJue_HunDun struct {
	BaseBuff
}

func Buff_ChangShengJue_HunDun_Builder(player *Player, level int, args ...interface{}) Buff {
	return &Buff_ChangShengJue_HunDun{
		BaseBuff{
			Name:     "长生诀.混沌",
			Type:     "force",
			Creator:  player,
			Duration: 30000,
			Steady:   true,
		},
	}
}

func init() {
	PerformRepo.Add("长生诀.混沌诀", Perform_ChangShengJue_HunDunJue_Builder)
	PerformRepo.Add("长生诀.天地诀", Perform_ChangShengJue_TianDiJue_Builder)
	BuffRepo.Add("长生诀.混沌", Buff_ChangShengJue_HunDun_Builder)
}
