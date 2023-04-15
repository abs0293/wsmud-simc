package simulator

// 绝招:枯木神功.枯木逢春
type Perform_KuMuShenGong_KuMuFengChun struct {
	BasePerform
	rate float64
}

func (pfm *Perform_KuMuShenGong_KuMuFengChun) Run(ctx *RunContext) {
	pfm.PreFlight(ctx)
	if ctx.IsFail() {
		return
	}
	pfm.Player.Buffs.ClearDebuffs()
	pfm.Player.AddHP(pfm.Player.GetHPMax() * pfm.rate)
}

func Perform_KuMuShenGong_KuMuFengChun_Builder(player *Player, level int, mixed bool) Perform {
	rate := float64(level)*0.0001 + 0.1
	return &Perform_KuMuShenGong_KuMuFengChun{
		BasePerform: BasePerform{
			Name:     "枯木神功.枯木逢春",
			Type:     "force",
			Player:   player,
			Level:    level,
			Mixed:    mixed,
			CoolDown: 60000,
			Timer:    NewTimer(),
		},
		rate: rate,
	}
}

func init() {
	PerformRepo.Add("枯木神功.枯木逢春", Perform_KuMuShenGong_KuMuFengChun_Builder)
}
