package simulator

// 绝招:灵犀步.灵犀
type Perform_LinXiBu_LinXi struct {
	BasePerform
}

func (pfm *Perform_LinXiBu_LinXi) Run(ctx *RunContext) {
	pfm.PreFlight(ctx)
	if ctx.IsFail() {
		return
	}

	pfm.Player.AddBuff(BuffRepo.Build("灵犀步.灵犀", pfm.Player, pfm.Level))
}

func Perform_LinXiBu_LinXi_Builder(player *Player, level int, mixed bool) Perform {
	return &Perform_LinXiBu_LinXi{
		BasePerform: BasePerform{
			Name:     "灵犀步.灵犀",
			Type:     "dodge",
			Player:   player,
			Level:    level,
			Mixed:    mixed,
			CoolDown: 60000,
			Timer:    NewTimer(),
		},
	}
}

// 光环:灵犀步.灵犀
type Buff_LinXiBu_LinXi struct {
	BaseBuff
}

func Buff_LinXiBu_LinXi_Builder(player *Player, level int, args ...interface{}) Buff {
	return &Buff_LinXiBu_LinXi{
		BaseBuff{
			Name:     "灵犀步.灵犀",
			Type:     "dodge",
			Creator:  player,
			Duration: 10000,
			Modifiers: []Modifier{
				{"负面抵抗%", 1},
			},
		},
	}
}

func init() {
	PerformRepo.Add("灵犀步.灵犀", Perform_LinXiBu_LinXi_Builder)
	BuffRepo.Add("灵犀步.灵犀", Buff_LinXiBu_LinXi_Builder)
}
