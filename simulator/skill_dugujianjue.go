package simulator

// 绝招:独孤剑诀.剑意
type Perform_DuGuJianJue_JianYi struct {
	BasePerform
}

func (pfm *Perform_DuGuJianJue_JianYi) Run(ctx *RunContext) {
	pfm.PreFlight(ctx)
	if ctx.IsFail() {
		return
	}
	pfm.Player.AddBuff(BuffRepo.Build("独孤剑诀.剑意", pfm.Player, pfm.Level))
}

func Perform_DuGuJianJue_JianYi_Builder(player *Player, level int, mixed bool) Perform {
	return &Perform_DuGuJianJue_JianYi{
		BasePerform: BasePerform{
			Name:     "独孤剑诀.剑意",
			Type:     "weapon",
			Player:   player,
			Level:    level,
			Mixed:    mixed,
			CoolDown: 30000,
			Timer:    NewTimer(),
		},
	}
}

// 绝招:独孤剑诀.剑来
type Perform_DuGuJianJue_JianLai struct {
	BasePerform
}

func (pfm Perform_DuGuJianJue_JianLai) HitNumber() int {
	return 8
}

func (pfm *Perform_DuGuJianJue_JianLai) Run(ctx *RunContext) {
	pfm.PreFlight(ctx)
	if ctx.IsFail() {
		return
	}
	pDuGu := pfm.Player.Skills.Passives.DuGu

	pfm.Player.AddBuff(BuffRepo.Build("独孤剑诀.剑来", pfm.Player, pfm.Level))
	for i := 0; i < pfm.HitNumber(); i++ {
		pfm.Attack(ctx)
		pDuGu.DoExtraHit(ctx.AddAttack())
	}
}

func Perform_DuGuJianJue_JianLai_Builder(player *Player, level int, mixed bool) Perform {
	return &Perform_DuGuJianJue_JianLai{
		BasePerform: BasePerform{
			Name:     "独孤剑诀.剑来",
			Type:     "weapon",
			Player:   player,
			Level:    level,
			Mixed:    mixed,
			CoolDown: 60000,
			Timer:    NewTimer(),
		},
	}
}

// 光环:独孤剑诀.剑意
type Buff_DuGuJianJue_JianYi struct {
	BaseBuff
}

func (b *Buff_DuGuJianJue_JianYi) GetModifier(n string) float64 {
	return b.BaseBuff.GetModifier(n)
}

func Buff_DuGuJianJue_JianYi_Builder(player *Player, level int, args ...interface{}) Buff {
	mod := 0.5 + float64(level/50)*0.01
	return &Buff_DuGuJianJue_JianYi{
		BaseBuff{
			Name:     "独孤剑诀.剑意",
			Type:     "weapon",
			Creator:  player,
			Duration: 5000 + level*2,
			Modifiers: []Modifier{
				{"命中%", mod},
				{"招架%", mod},
				{"终伤%", mod},
				{"攻速%", 1},
			},
		},
	}
}

// 光环:独孤剑诀.剑来
type Buff_DuGuJianJue_JianLai struct {
	BaseBuff
}

func (b *Buff_DuGuJianJue_JianLai) GetModifier(n string) float64 {
	return b.BaseBuff.GetModifier(n)
}

func Buff_DuGuJianJue_JianLai_Builder(player *Player, level int, args ...interface{}) Buff {
	return &Buff_DuGuJianJue_JianLai{
		BaseBuff{
			Name:      "独孤剑诀.剑来",
			Type:      "feijian",
			Creator:   player,
			Duration:  10000,
			Modifiers: []Modifier{},
		},
	}
}

// 光环:独孤剑诀.无我
type Buff_DuGuJianJue_WuWo struct {
	BaseBuff
	value float64
	timer *Timer
}

func (b *Buff_DuGuJianJue_WuWo) OnEnable() {
	b.timer.Start(5000)
	b.BaseBuff.OnEnable()
}

func (b *Buff_DuGuJianJue_WuWo) Update(diff int) {
	if b.Stacks < b.StackMax {
		b.timer.Update(diff)
		if b.timer.IsDone() {
			b.Stacks++
			b.EventRefreshLog()
			b.timer.Start(5000)
		}
	}
}

func (b *Buff_DuGuJianJue_WuWo) GetModifier(n string) float64 {
	switch n {
	case "命中%", "攻击%", "招架%":
		return b.value * float64(b.Stacks)
	default:
		return 0
	}
}

func Buff_DuGuJianJue_WuWo_Builder(player *Player, level int, args ...interface{}) Buff {
	return &Buff_DuGuJianJue_WuWo{
		BaseBuff{
			Name:      "独孤剑诀.无我",
			Type:      "wuwo",
			Creator:   player,
			Permanent: true,
			Steady:    true,
			Stackable: true,
			StackMax:  10,
			Stacks:    1,
			Modifiers: []Modifier{},
		},
		0.03 + float64(level/1000)*0.01,
		NewTimer(),
	}
}

// 光环:独孤剑诀.无剑
type Buff_DuGuJianJue_WuJian struct {
	BaseBuff
}

func Buff_DuGuJianJue_WuJian_Builder(player *Player, level int, args ...interface{}) Buff {
	return &Buff_DuGuJianJue_WuJian{
		BaseBuff{
			Name:      "独孤剑诀.无剑",
			Type:      "wujian",
			Creator:   player,
			Permanent: true,
			Steady:    true,
			Modifiers: []Modifier{
				{"免疫缴械d", 1},
			},
		},
	}
}

func init() {
	PerformRepo.Add("独孤剑诀.剑意", Perform_DuGuJianJue_JianYi_Builder)
	PerformRepo.Add("独孤剑诀.剑来", Perform_DuGuJianJue_JianLai_Builder)
	BuffRepo.Add("独孤剑诀.剑意", Buff_DuGuJianJue_JianYi_Builder)
	BuffRepo.Add("独孤剑诀.剑来", Buff_DuGuJianJue_JianLai_Builder)
	BuffRepo.Add("独孤剑诀.无我", Buff_DuGuJianJue_WuWo_Builder)
	BuffRepo.Add("独孤剑诀.无剑", Buff_DuGuJianJue_WuJian_Builder)
}
