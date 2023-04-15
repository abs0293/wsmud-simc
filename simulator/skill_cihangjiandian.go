package simulator

// 慈航剑典.剑心通明
type Buff_CiHangJianDian_JianXin struct {
	BaseBuff
	useSword bool
}

func (b *Buff_CiHangJianDian_JianXin) GetModifier(n string) float64 {
	if (n == "忽防%" || n == "终伤%") && !b.useSword {
		return 0
	}
	return b.BaseBuff.GetModifier(n)
}

func Buff_CiHangJianDian_JianXin_Builder(player *Player, level int, args ...interface{}) Buff {
	mod := 0.3 + float64(level/100)*0.01
	return &Buff_CiHangJianDian_JianXin{
		BaseBuff{
			Name:     "慈航剑典.剑心",
			Type:     "force",
			Creator:  player,
			Steady:   true,
			Duration: 10000,
			Modifiers: []Modifier{
				{"绝招冷却%", 1},
				{"绝招释放%", 1},
				{"忽防%", mod / 2},
				{"终伤%", mod},
				{"绝对命中", 1},
			},
		},
		player.GetWeaponType() == WeaponTypeSword,
	}
}

type Perform_CiHangJianDian_JianXinTongMing struct {
	*BasePerform
}

func (pfm *Perform_CiHangJianDian_JianXinTongMing) Run(ctx *RunContext) {
	pfm.PreFlight(ctx)
	if ctx.IsFail() {
		return
	}

	var (
		attacker = ctx.attacker
	)
	attacker.AddBuff(BuffRepo.Build("慈航剑典.剑心", attacker, pfm.Level))
}

func Perform_CiHangJianDian_JianXinTongMing_Builder(
	player *Player,
	level int,
	mixed bool,
) Perform {
	return &Perform_CiHangJianDian_JianXinTongMing{
		BasePerform: &BasePerform{
			Name:     "慈航剑典.剑心通明",
			Type:     "force",
			Player:   player,
			Level:    level,
			Mixed:    mixed,
			CoolDown: 60000,
			Timer:    NewTimer(),
		},
	}
}

func init() {
	PerformRepo.Add("慈航剑典.剑心通明", Perform_CiHangJianDian_JianXinTongMing_Builder)
	BuffRepo.Add("慈航剑典.剑心", Buff_CiHangJianDian_JianXin_Builder)
}
