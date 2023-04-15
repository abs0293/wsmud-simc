package simulator

// 绝招:移花接木.移花
type Perform_YiHuaJieMu_YiHua struct {
	BasePerform
}

func (pfm *Perform_YiHuaJieMu_YiHua) Run(ctx *RunContext) {
	pfm.PreFlight(ctx)
	if ctx.IsFail() {
		return
	}

	var (
		attacker = ctx.attacker
		target   = ctx.target
	)

	attacker.AddBuff(BuffRepo.Build("移花接木.移花.绿", pfm.Player, pfm.Level))
	target.AddBuff(BuffRepo.Build("移花接木.移花.红", pfm.Player, pfm.Level))
}

func Perform_YiHuaJieMu_YiHua_Builder(player *Player, level int, mixed bool) Perform {
	return &Perform_YiHuaJieMu_YiHua{
		BasePerform: BasePerform{
			Name:     "移花接木.移花",
			Type:     "parry",
			Player:   player,
			Level:    level,
			Mixed:    mixed,
			CoolDown: 25000,
			Timer:    NewTimer(),
		},
	}
}

// 光环:移花接木.移花.绿
type Buff_YiHuaJieMu_YiHua_Green struct {
	BaseBuff
}

func Buff_YiHuaJieMu_YiHua_Green_Builder(player *Player, level int, args ...interface{}) Buff {
	duration := Sec2Ms(5.0 + float64(level)*0.003)
	return &Buff_YiHuaJieMu_YiHua_Green{
		BaseBuff{
			Name:      "移花接木.移花.绿",
			Type:      "yihua.green",
			Creator:   player,
			Duration:  duration,
			Modifiers: []Modifier{},
		},
	}
}

// 光环:移花接木.移花.红
type Buff_YiHuaJieMu_YiHua_Red struct {
	BaseBuff
}

func Buff_YiHuaJieMu_YiHua_Red_Builder(player *Player, level int, args ...interface{}) Buff {
	duration := Sec2Ms(5.0 + float64(level)*0.003)
	mod := 0.3 + float64(level/300)*0.01
	return &Buff_YiHuaJieMu_YiHua_Red{
		BaseBuff{
			Name:     "移花接木.移花.红",
			Type:     "yihua.red",
			Creator:  player,
			Duration: duration,
			Debuff:   true,
			Modifiers: []Modifier{
				{"闪避%", -mod},
				{"招架%", -mod},
			},
		},
	}
}

func init() {
	PerformRepo.Add("移花接木.移花", Perform_YiHuaJieMu_YiHua_Builder)
	BuffRepo.Add("移花接木.移花.绿", Buff_YiHuaJieMu_YiHua_Green_Builder)
	BuffRepo.Add("移花接木.移花.红", Buff_YiHuaJieMu_YiHua_Red_Builder)
}
