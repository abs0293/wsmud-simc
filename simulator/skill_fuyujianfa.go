package simulator

// 绝招:覆雨剑法.剑雨
type Perform_FuYuJianFa_JianYu struct {
	BasePerform
}

func (pfm *Perform_FuYuJianFa_JianYu) Run(ctx *RunContext) {
	pfm.PreFlight(ctx)
	if ctx.IsFail() {
		return
	}

	// pfm.Attack(ctx)
	pfm.Player.AddBuff(BuffRepo.Build("覆雨剑法.剑雨", pfm.Player, pfm.Level))
}

func Perform_FuYuJianFa_JianYu_Builder(player *Player, level int, mixed bool) Perform {
	return &Perform_FuYuJianFa_JianYu{
		BasePerform: BasePerform{
			Name:     "覆雨剑法.剑雨",
			Type:     "weapon",
			Player:   player,
			Level:    level,
			Mixed:    mixed,
			CoolDown: 30000,
			Timer:    NewTimer(),
		},
	}
}

// 绝招:覆雨剑法.剑罡
type Perform_FuYuJianFa_JianGang struct {
	BasePerform
}

func (pfm *Perform_FuYuJianFa_JianGang) Run(ctx *RunContext) {
	pfm.PreFlight(ctx)
	if ctx.IsFail() {
		return
	}

	// pfm.Attack(ctx)
	// pfm.Player.AddBuff(BuffRepo.Build("", pfm.Player, pfm.Level))
}

func Perform_FuYuJianFa_JianGang_Builder(player *Player, level int, mixed bool) Perform {
	return &Perform_FuYuJianFa_JianGang{
		BasePerform: BasePerform{
			Name:     "覆雨剑法.剑罡",
			Type:     "weapon",
			Player:   player,
			Level:    level,
			Mixed:    mixed,
			CoolDown: 20000,
			Timer:    NewTimer(),
		},
	}
}

// 光环:覆雨剑法.剑雨
type Buff_FuYuJianFa_JianYu struct {
	BaseBuff
}

func (b *Buff_FuYuJianFa_JianYu) GetModifier(n string) float64 {
	return b.BaseBuff.GetModifier(n)
}

func Buff_FuYuJianFa_JianYu_Builder(player *Player, level int, args ...interface{}) Buff {
	return &Buff_FuYuJianFa_JianYu{
		BaseBuff{
			Name:      "覆雨剑法.剑雨",
			Type:      "weapon",
			Creator:   player,
			Duration:  10000,
			Modifiers: []Modifier{},
		},
	}
}

func init() {
	PerformRepo.Add("覆雨剑法.剑雨", Perform_FuYuJianFa_JianYu_Builder)
	PerformRepo.Add("覆雨剑法.剑罡", Perform_FuYuJianFa_JianGang_Builder)
	BuffRepo.Add("覆雨剑法.剑雨", Buff_FuYuJianFa_JianYu_Builder)
}
