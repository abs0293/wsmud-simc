package simulator

// 绝招:天魔策.补天道
type Perform_TianMoCe_BuTianDao struct {
	BasePerform
}

func (pfm *Perform_TianMoCe_BuTianDao) Run(ctx *RunContext) {
	pfm.PreFlight(ctx)
	if ctx.IsFail() {
		return
	}

	pfm.Player.AddBuff(BuffRepo.Build("天魔策.补天", pfm.Player, pfm.Level))
}

func Perform_TianMoCe_BuTianDao_Builder(player *Player, level int, mixed bool) Perform {
	return &Perform_TianMoCe_BuTianDao{
		BasePerform: BasePerform{
			Name:     "天魔策.补天道",
			Type:     "parry",
			Player:   player,
			Level:    level,
			Mixed:    mixed,
			CoolDown: 45000,
			Timer:    NewTimer(),
		},
	}
}

// 绝招:天魔策.种魔
type Perform_TianMoCe_ZhongMo struct {
	BasePerform
}

func (pfm *Perform_TianMoCe_ZhongMo) Run(ctx *RunContext) {
	pfm.PreFlight(ctx)
	if ctx.IsFail() {
		return
	}

	ctx.target.RemoveBuffByType("force")
	ctx.target.AddBuff(BuffRepo.Build("天魔策.种魔", pfm.Player, pfm.Level))
}

func Perform_TianMoCe_ZhongMo_Builder(player *Player, level int, mixed bool) Perform {
	return &Perform_TianMoCe_ZhongMo{
		BasePerform: BasePerform{
			Name:     "天魔策.种魔",
			Type:     "force",
			Player:   player,
			Level:    level,
			Mixed:    mixed,
			CoolDown: 45000,
			Timer:    NewTimer(),
		},
	}
}

// 绝招:天魔策.道心
type Perform_TianMoCe_DaoXin struct {
	BasePerform
}

func (pfm *Perform_TianMoCe_DaoXin) Run(ctx *RunContext) {
	pfm.PreFlight(ctx)
	if ctx.IsFail() {
		return
	}

	// pfm.BasePerform.Hit(target, Modifier{})
}

func Perform_TianMoCe_DaoXin_Builder(player *Player, level int, mixed bool) Perform {
	return &Perform_TianMoCe_DaoXin{
		BasePerform: BasePerform{
			Name:     "天魔策.道心",
			Type:     "force",
			Player:   player,
			Level:    level,
			Mixed:    mixed,
			CoolDown: 45000,
			Timer:    NewTimer(),
		},
	}
}

// 绝招:天魔策.鬼影
type Perform_TianMoCe_GuiYing struct {
	BasePerform
}

func (pfm *Perform_TianMoCe_GuiYing) Run(ctx *RunContext) {
	pfm.PreFlight(ctx)
	if ctx.IsFail() {
		return
	}

	// pfm.BasePerform.Hit(target, Modifier{})
}

func Perform_TianMoCe_GuiYing_Builder(player *Player, level int, mixed bool) Perform {
	return &Perform_TianMoCe_GuiYing{
		BasePerform: BasePerform{
			Name:     "天魔策.鬼影",
			Type:     "dodge",
			Player:   player,
			Level:    level,
			Mixed:    mixed,
			CoolDown: 40000,
			Timer:    NewTimer(),
		},
	}
}

// 绝招:天魔策.拳罡
type Perform_TianMoCe_QuanGang struct {
	BasePerform
}

func (pfm *Perform_TianMoCe_QuanGang) Run(ctx *RunContext) {
	pfm.PreFlight(ctx)
	if ctx.IsFail() {
		return
	}

	if actx := pfm.Attack(ctx); actx.GetDamageFinal() > 0 {
		ctx.attacker.AddBuff(BuffRepo.Build("天魔策.拳罡", pfm.Player, pfm.Level))
	}
}

func Perform_TianMoCe_QuanGang_Builder(player *Player, level int, mixed bool) Perform {
	return &Perform_TianMoCe_QuanGang{
		BasePerform: BasePerform{
			Name:     "天魔策.拳罡",
			Type:     "unarmed",
			Player:   player,
			Level:    level,
			Mixed:    mixed,
			CoolDown: 10000,
			Timer:    NewTimer(),
		},
	}
}

// 光环:天魔策.补天
type Buff_TianMoCe_BuTian struct {
	BaseBuff
}

func Buff_TianMoCe_BuTian_Builder(player *Player, level int, args ...interface{}) Buff {
	return &Buff_TianMoCe_BuTian{
		BaseBuff{
			Name:     "天魔策.补天",
			Type:     "parry",
			Creator:  player,
			Duration: 15000,
			Modifiers: []Modifier{
				{"招架%", 0.8},
			},
		},
	}
}

// 光环:天魔策.种魔
type Buff_TianMoCe_ZhongMo struct {
	BaseBuff
}

func Buff_TianMoCe_ZhongMo_Builder(player *Player, level int, args ...interface{}) Buff {
	return &Buff_TianMoCe_ZhongMo{
		BaseBuff{
			Name:      "天魔策.种魔",
			Type:      "force",
			Creator:   player,
			Debuff:    true,
			Duration:  15000,
			Modifiers: []Modifier{},
		},
	}
}

// 光环:天魔策.鬼影
type Buff_TianMoCe_GuiYing struct {
	BaseBuff
}

func (b *Buff_TianMoCe_GuiYing) GetModifier(n string) float64 {
	return b.BaseBuff.GetModifier(n)
}

func Buff_TianMoCe_GuiYing_Builder(player *Player, level int, args ...interface{}) Buff {
	return &Buff_TianMoCe_GuiYing{
		BaseBuff{
			Name:      "天魔策.鬼影",
			Type:      "dodge",
			Creator:   player,
			Duration:  5000 + (level/2000)*1000,
			Modifiers: []Modifier{{"绝对命中", 1}},
		},
	}
}

// 光环:天魔策.道心
type Buff_TianMoCe_DaoXin struct {
	BaseBuff
}

func (b *Buff_TianMoCe_DaoXin) GetModifier(n string) float64 {
	return b.BaseBuff.GetModifier(n)
}

func Buff_TianMoCe_DaoXin_Builder(player *Player, level int, args ...interface{}) Buff {
	return &Buff_TianMoCe_DaoXin{
		BaseBuff{
			Name:      "天魔策.道心",
			Type:      "force",
			Creator:   player,
			Modifiers: []Modifier{},
		},
	}
}

// 光环:天魔策.拳罡
type Buff_TianMoCe_QuanGang struct {
	BaseBuff
}

func Buff_TianMoCe_QuanGang_Builder(player *Player, level int, args ...interface{}) Buff {
	mod := 0.01 + float64(level/200)*0.01
	return &Buff_TianMoCe_QuanGang{
		BaseBuff{
			Name:      "天魔策.拳罡",
			Type:      "unarmed",
			Creator:   player,
			Duration:  6000,
			Stackable: true,
			StackMax:  5,
			Modifiers: []Modifier{
				{"终伤%", mod},
				{"忽防%", mod},
			},
		},
	}
}

func init() {
	PerformRepo.Add("天魔策.补天道", Perform_TianMoCe_BuTianDao_Builder)
	PerformRepo.Add("天魔策.种魔", Perform_TianMoCe_ZhongMo_Builder)
	PerformRepo.Add("天魔策.道心", Perform_TianMoCe_DaoXin_Builder)
	PerformRepo.Add("天魔策.鬼影", Perform_TianMoCe_GuiYing_Builder)
	PerformRepo.Add("天魔策.拳罡", Perform_TianMoCe_QuanGang_Builder)
	BuffRepo.Add("天魔策.补天", Buff_TianMoCe_BuTian_Builder)
	BuffRepo.Add("天魔策.种魔", Buff_TianMoCe_ZhongMo_Builder)
	BuffRepo.Add("天魔策.鬼影", Buff_TianMoCe_GuiYing_Builder)
	BuffRepo.Add("天魔策.道心", Buff_TianMoCe_DaoXin_Builder)
	BuffRepo.Add("天魔策.拳罡", Buff_TianMoCe_QuanGang_Builder)
}
