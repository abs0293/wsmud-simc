package simulator

import (
	"math"
)

// 绝招:阴阳九转.九烛
type Perform_YinYangJiuZhuan_JiuZhu struct {
	BasePerform
}

func (pfm *Perform_YinYangJiuZhuan_JiuZhu) Run(ctx *RunContext) {
	pfm.PreFlight(ctx)
	if ctx.IsFail() {
		return
	}
	pfm.Player.AddBuff(BuffRepo.Build("阴阳九转.九烛", pfm.Player, pfm.Level))
}

func Perform_YinYangJiuZhuan_JiuZhu_Builder(player *Player, level int, mixed bool) Perform {
	return &Perform_YinYangJiuZhuan_JiuZhu{
		BasePerform: BasePerform{
			Name:     "阴阳九转.九烛",
			Type:     "force",
			Player:   player,
			Level:    level,
			Mixed:    mixed,
			InitOnly: true,
			CoolDown: 30000,
			Timer:    NewTimer(),
		},
	}
}

// 绝招:阴阳九转.九幽
type Perform_YinYangJiuZhuan_JiuYou struct {
	BasePerform
}

func (pfm *Perform_YinYangJiuZhuan_JiuYou) Run(ctx *RunContext) {
	pfm.PreFlight(ctx)
	if ctx.IsFail() {
		return
	}
	pfm.Player.AddBuff(BuffRepo.Build("阴阳九转.九幽", pfm.Player, pfm.Level))
}

func Perform_YinYangJiuZhuan_JiuYou_Builder(player *Player, level int, mixed bool) Perform {
	return &Perform_YinYangJiuZhuan_JiuYou{
		BasePerform: BasePerform{
			Name:     "阴阳九转.九幽",
			Type:     "force",
			Player:   player,
			Level:    level,
			InitOnly: true,
			Mixed:    mixed,
			CoolDown: 30000,
			Timer:    NewTimer(),
		},
	}
}

// 绝招:阴阳九转.定乾坤
type Perform_YinYangJiuZhuan_DingQianKun struct {
	BasePerform
}

func (pfm *Perform_YinYangJiuZhuan_DingQianKun) GetDisPerformTime() int {
	t := Sec2Ms(float64(pfm.Level) * 0.002)
	if t > 5000 {
		t = 5000
	}
	return t
}

func (pfm *Perform_YinYangJiuZhuan_DingQianKun) Run(ctx *RunContext) {
	pfm.PreFlight(ctx)
	if ctx.IsFail() {
		return
	}

	var (
		target = ctx.target
	)

	target.PerformCD.Start(pfm.GetDisPerformTime())
	pfm.Attack(ctx)
}

func Perform_YinYangJiuZhuan_DingQianKun_Builder(player *Player, level int, mixed bool) Perform {
	return &Perform_YinYangJiuZhuan_DingQianKun{
		BasePerform: BasePerform{
			Name:     "阴阳九转.定乾坤",
			Type:     "force",
			Player:   player,
			Level:    level,
			Mixed:    mixed,
			CoolDown: 60000,
			Timer:    NewTimer(),
		},
	}
}

// 绝招:阴阳九转.镇天地
type Perform_YinYangJiuZhuan_ZhenTianDi struct {
	BasePerform
	InProgress bool
}

func (pfm *Perform_YinYangJiuZhuan_ZhenTianDi) CanRun(target *Player, args ...interface{}) bool {
	return pfm.BasePerform.CanRun(target, args...) && !pfm.InProgress
}

func (pfm *Perform_YinYangJiuZhuan_ZhenTianDi) Run(ctx *RunContext) {
	pfm.PreFlight(ctx)
	if ctx.IsFail() {
		return
	}
	pfm.InProgress = true
	pfm.Player.AddBuff(BuffRepo.Build("阴阳九转.镇天地", pfm.Player, pfm.Level, pfm))
}

func (pfm *Perform_YinYangJiuZhuan_ZhenTianDi) Update(diff int) {
	if !pfm.InProgress {
		pfm.Timer.Update(diff)
	}
}

func Perform_YinYangJiuZhuan_ZhenTianDi_Builder(player *Player, level int, mixed bool) Perform {
	return &Perform_YinYangJiuZhuan_ZhenTianDi{
		BasePerform: BasePerform{
			Name:     "阴阳九转.镇天地",
			Type:     "force",
			Player:   player,
			Level:    level,
			Mixed:    mixed,
			CoolDown: 60000,
			Timer:    NewTimer(),
		},
	}
}

// 光环:阴阳九转.九烛
type Buff_YinYangJiuZhuan_JiuZhu struct {
	BaseBuff
}

func Buff_YinYangJiuZhuan_JiuZhu_Builder(player *Player, level int, args ...interface{}) Buff {
	mod := 0.5 + float64(level)*0.00005
	return &Buff_YinYangJiuZhuan_JiuZhu{
		BaseBuff{
			Name:      "阴阳九转.九烛",
			Type:      "force",
			Creator:   player,
			Permanent: true,
			Steady:    true,
			Modifiers: []Modifier{
				{"气血%", mod},
				{"静态免伤%", mod},
			},
		},
	}
}

// 光环:阴阳九转.九幽
type Buff_YinYangJiuZhuan_JiuYou struct {
	BaseBuff
}

func Buff_YinYangJiuZhuan_JiuYou_Builder(player *Player, level int, args ...interface{}) Buff {
	mod := math.Min(0.8, 0.2+float64(level)*0.0002)
	return &Buff_YinYangJiuZhuan_JiuYou{
		BaseBuff{
			Name:      "阴阳九转.九幽",
			Type:      "force",
			Creator:   player,
			Permanent: true,
			Steady:    true,
			Modifiers: []Modifier{
				{"攻速%", mod / 2},
				{"忽防%", mod},
				{"攻击%", mod},
			},
		},
	}
}

// 光环:阴阳九转.镇天地
type Buff_YinYangJiuZhuan_ZhenTianDi struct {
	BaseBuff
	atkTimer *Timer
	pfm      *Perform_YinYangJiuZhuan_ZhenTianDi
}

func (b *Buff_YinYangJiuZhuan_ZhenTianDi) Update(diff int) {
	b.atkTimer.Update(diff)
	if b.atkTimer.IsDone() {
		if !b.Owner.State.Faint {
			dmg := (b.Owner.GetRawStaticDamageReducePercent() * 0.025) * b.Target.GetHP()
			ctx := NewCombatContext(
				b.Owner,
				b.Target,
				b.Owner.Arena.Ticks,
				Modifier{"绝招.附加伤害d", dmg},
			)
			ctx.SetForceAttack().SetSource(b.Name)
			GenericAttack(ctx)
			b.Owner.Log(ctx.log)
		}
		b.atkTimer.Start(1000)
	}
	b.BaseBuff.Update(diff)
}

func (b *Buff_YinYangJiuZhuan_ZhenTianDi) OnEnable() {
	b.atkTimer.Start(1000)
}

func (b *Buff_YinYangJiuZhuan_ZhenTianDi) OnDisable() {
	b.pfm.InProgress = false
}

func (b *Buff_YinYangJiuZhuan_ZhenTianDi) GetModifier(n string) float64 {
	return b.BaseBuff.GetModifier(n)
}

func Buff_YinYangJiuZhuan_ZhenTianDi_Builder(player *Player, level int, args ...interface{}) Buff {
	duration := Sec2Ms(math.Min(float64(level/300)+2.0, 15))
	pfm := args[0].(*Perform_YinYangJiuZhuan_ZhenTianDi)
	return &Buff_YinYangJiuZhuan_ZhenTianDi{
		BaseBuff{
			Name:      "阴阳九转.镇天地",
			Type:      "ztd",
			Creator:   player,
			Target:    player.Target,
			Steady:    true,
			Duration:  duration,
			Modifiers: []Modifier{{"免伤%", 0.8}},
		},
		NewTimer(),
		pfm,
	}
}

func init() {
	PerformRepo.Add("阴阳九转.九烛", Perform_YinYangJiuZhuan_JiuZhu_Builder)
	PerformRepo.Add("阴阳九转.九幽", Perform_YinYangJiuZhuan_JiuYou_Builder)
	PerformRepo.Add("阴阳九转.定乾坤", Perform_YinYangJiuZhuan_DingQianKun_Builder)
	PerformRepo.Add("阴阳九转.镇天地", Perform_YinYangJiuZhuan_ZhenTianDi_Builder)
	BuffRepo.Add("阴阳九转.九烛", Buff_YinYangJiuZhuan_JiuZhu_Builder)
	BuffRepo.Add("阴阳九转.九幽", Buff_YinYangJiuZhuan_JiuYou_Builder)
	BuffRepo.Add("阴阳九转.镇天地", Buff_YinYangJiuZhuan_ZhenTianDi_Builder)
}
