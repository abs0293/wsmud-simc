package simulator

// 绝招:神游太虚.凌波
type Perform_ShenYouTaiXu_LingBo struct {
	BasePerform
}

func (pfm *Perform_ShenYouTaiXu_LingBo) Run(ctx *RunContext) {
	pfm.PreFlight(ctx)
	if ctx.IsFail() {
		return
	}
	pfm.Player.AddBuff(BuffRepo.Build("神游太虚.凌波", pfm.Player, pfm.Level))
}

func Perform_ShenYouTaiXu_LingBo_Builder(player *Player, level int, mixed bool) Perform {
	return &Perform_ShenYouTaiXu_LingBo{
		BasePerform: BasePerform{
			Name:     "神游太虚.凌波",
			Type:     "dodge",
			Player:   player,
			Level:    level,
			Mixed:    mixed,
			CoolDown: 60000,
			Timer:    NewTimer(),
		},
	}
}

// 绝招:神游太虚.化蝶
type Perform_ShenYouTaiXu_HuaDie struct {
	BasePerform
}

func (pfm *Perform_ShenYouTaiXu_HuaDie) Run(ctx *RunContext) {
	pfm.PreFlight(ctx)
	if ctx.IsFail() {
		return
	}

	// pfm.BasePerform.Hit(target, Modifier{})
}

func Perform_ShenYouTaiXu_HuaDie_Builder(player *Player, level int, mixed bool) Perform {
	return &Perform_ShenYouTaiXu_HuaDie{
		BasePerform: BasePerform{
			Name:     "神游太虚.化蝶",
			Type:     "dodge",
			Player:   player,
			Level:    level,
			Mixed:    mixed,
			CoolDown: 60000,
			Timer:    NewTimer(),
		},
	}
}

// 光环:神游太虚.凌波
type Buff_ShenYouTaiXu_LingBo struct {
	BaseBuff
}

func Buff_ShenYouTaiXu_LingBo_Builder(player *Player, level int, args ...interface{}) Buff {
	return &Buff_ShenYouTaiXu_LingBo{
		BaseBuff{
			Name:     "神游太虚.凌波",
			Type:     "dodge",
			Creator:  player,
			Duration: 15000,
			Modifiers: []Modifier{
				{"闪避d", 3e6},
				{"命中d", player.GetDodge()},
			},
		},
	}
}

func init() {
	PerformRepo.Add("神游太虚.凌波", Perform_ShenYouTaiXu_LingBo_Builder)
	PerformRepo.Add("神游太虚.化蝶", Perform_ShenYouTaiXu_HuaDie_Builder)
	BuffRepo.Add("神游太虚.凌波", Buff_ShenYouTaiXu_LingBo_Builder)
}
