package simulator

// 绝招:黯然销魂掌.无中生有
type Perform_AnRanXiaoHunZhang_WuZhongShenYou struct {
	BasePerform
	dmgMod float64
}

func (pfm *Perform_AnRanXiaoHunZhang_WuZhongShenYou) Run(ctx *RunContext) {
	pfm.PreFlight(ctx)
	if ctx.IsFail() {
		return
	}

	var (
		attacker = ctx.attacker
		target   = ctx.target
	)

	ret := pfm.Attack(ctx, Modifier{"绝招.伤害倍率%", pfm.dmgMod})

	if ret.GetDamageFinal() > 0 {
		buf := target.Buffs.BeStolen(pfm.Player, pfm.Name)
		if buf != nil {
			attacker.AddBuff(buf)
		}
	}
}

func Perform_AnRanXiaoHunZhang_WuZhongShenYou_Builder(player *Player, level int, mixed bool) Perform {
	return &Perform_AnRanXiaoHunZhang_WuZhongShenYou{
		BasePerform: BasePerform{
			Name:     "黯然销魂掌.无中生有",
			Type:     "unarmed",
			Player:   player,
			Level:    level,
			Mixed:    mixed,
			CoolDown: 20000,
			Timer:    NewTimer(),
		},
		dmgMod: 1 + float64(level)*0.001,
	}
}

// 绝招:黯然销魂掌.呆若木鸡
type Perform_AnRanXiaoHunZhang_DaiRuoMuJi struct {
	BasePerform
	dmgMod float64
}

func (pfm *Perform_AnRanXiaoHunZhang_DaiRuoMuJi) Run(ctx *RunContext) {
	pfm.PreFlight(ctx)
	if ctx.IsFail() {
		return
	}

	var (
		attacker = ctx.attacker
		target   = ctx.target
	)

	ret := pfm.Attack(ctx, Modifier{"绝招.伤害倍率%", pfm.dmgMod}, Modifier{"命中%", 1})

	if ret.GetDamageFinal() > 0 {
		target.AddBuff(BuffRepo.Build("黯然销魂掌.迟钝", attacker, pfm.Level))
	}
}

func Perform_AnRanXiaoHunZhang_DaiRuoMuJi_Builder(player *Player, level int, mixed bool) Perform {
	return &Perform_AnRanXiaoHunZhang_DaiRuoMuJi{
		BasePerform: BasePerform{
			Name:     "黯然销魂掌.呆若木鸡",
			Type:     "unarmed",
			Player:   player,
			Level:    level,
			Mixed:    mixed,
			CoolDown: 36000,
			Timer:    NewTimer(),
		},
		dmgMod: 1 + float64(level)/600,
	}
}

// 光环:黯然销魂掌.迟钝
//
//	看看这是正经debuff能有的持续时间吗
type Buff_AnRanXiaoHunZhang_ChiDun struct {
	BaseBuff
}

func (b *Buff_AnRanXiaoHunZhang_ChiDun) GetModifier(n string) float64 {
	return b.BaseBuff.GetModifier(n)
}

func Buff_AnRanXiaoHunZhang_ChiDun_Builder(player *Player, level int, args ...interface{}) Buff {
	mod := 0.5 + float64(level)/10000
	return &Buff_AnRanXiaoHunZhang_ChiDun{
		BaseBuff{
			Name:     "黯然销魂掌.迟钝",
			Type:     "chidun",
			Debuff:   true,
			Creator:  player,
			Duration: 10000 + Sec2Ms(float64(level/200)),
			Modifiers: []Modifier{
				{"攻速%", -mod},
				{"绝招释放%", -mod},
			},
		},
	}
}

func init() {
	PerformRepo.Add("黯然销魂掌.无中生有", Perform_AnRanXiaoHunZhang_WuZhongShenYou_Builder)
	PerformRepo.Add("黯然销魂掌.呆若木鸡", Perform_AnRanXiaoHunZhang_DaiRuoMuJi_Builder)
	BuffRepo.Add("黯然销魂掌.迟钝", Buff_AnRanXiaoHunZhang_ChiDun_Builder)
}
