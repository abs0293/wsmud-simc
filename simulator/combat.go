package simulator

import (
	"math"
	"math/rand"
)

const (
	BODYPART_HEAD int = iota
	BODYPART_NECK
	BODYPART_CHEST
	BODYPART_BACK       // 后背
	BODYPART_BACKCENTER // 后心
	BODYPART_LEFTSHOULDER
	BODYPART_RIGHTSHOULDER
	BODYPART_ABDOMEN // 小腹
	BODYPART_WAIST   //腰
	BODYPART_LEFTHAND
	BODYPART_RIGHTHAND
	BODYPART_LEFTLEG
	BODYPART_RIGHTLEG
	BODYPART_LEFTFOOT
	BODYPART_RIGHTFOOT
)

var BodyPart_DamageModifier = []float64{
	1.2,  //BODYPART_HEAD
	1.1,  //BODYPART_NECK
	0.95, //BODYPART_CHEST
	0.97, //BODYPART_BACK
	1,    //BODYPART_BACKCENTER
	0.85, //BODYPART_LEFTSHOULDER
	0.85, //BODYPART_RIGHTSHOULDER
	0.9,  //BODYPART_ABDOMEN
	0.99, //BODYPART_WAIST
	0.85, //BODYPART_LEFTHAND
	0.85, //BODYPART_RIGHTHAND
	0.85, //BODYPART_LEFTLEG
	0.85, //BODYPART_RIGHTLEG
	0.8,  //BODYPART_LEFTFOOT
	0.8,  //BODYPART_RIGHTFOOT
}

const (
	Hit_Result_Unknown int = iota
	Hit_Result_Hit
	Hit_Result_Dodge
	Hit_Result_Parry
)

var Hit_Result_String = map[int]string{
	Hit_Result_Hit:   "命中",
	Hit_Result_Dodge: "躲闪",
	Hit_Result_Parry: "招架",
}

func CalcChance(x float64, min float64, max float64) float64 {
	if x < min {
		return 0
	}
	if x > max {
		return 1
	}
	return (x - min) / (max - min)
}

// 命中判定
//
//	1.先判定是否闪避, 闪避失败再判定是否招架
//	2.命中<闪避/2必闪避,命中>=闪避必命中
//	3.命中<招架/2必招架,命中>=招架必命中
func HitCheck(roll *rand.Rand, atkHit, tarDodge, tarParry float64) int {
	if atkHit < tarDodge/2 {
		return Hit_Result_Dodge
	}
	if atkHit < tarParry/2 {
		return Hit_Result_Parry
	}
	r := roll.Float64()
	if tarDodge > 0 && atkHit >= tarDodge/2 && atkHit < tarDodge {
		if r <= CalcChance(atkHit, tarDodge/2, tarDodge) {
			return Hit_Result_Dodge
		}
	}
	if tarParry > 0 && atkHit >= tarParry/2 && atkHit < tarParry {
		if r < CalcChance(atkHit, tarParry/2, tarParry) {
			return Hit_Result_Parry
		}
	}
	return Hit_Result_Hit
}

// 防御值 免伤
func CalcDefenceDamageReduction(damage float64, defence float64) float64 {
	return damage / (damage + defence)
}

// 免伤计算
//
//	1.先计算防御值(防御d),再计算百分比免伤(免伤%),最后计算固定免伤(免伤d)
//	2.忽视防御百分比(忽防%)超过100%的部分可以减少百分比免伤(免伤%),加法计算
//	3.百分比免伤最小为0
func CalcDamageReduction(
	damage float64,
	igDefPct float64,
	defence float64,
	dmgRed float64,
	dmgReducePct float64,
) float64 {
	// 防御d
	if igDefPct < 0 {
		igDefPct = 0
	}
	if igDefPct < 1 {
		defence *= 1 - igDefPct
		damage = CalcDefenceDamageReduction(damage, defence)
	}
	// 免伤%
	igDefPct = math.Max(0, igDefPct-1)
	dmgReducePct = math.Max(0, dmgReducePct-igDefPct)
	damage *= 1 - dmgReducePct
	// 免伤d
	damage = math.Max(0, damage-dmgRed)
	return damage
}

func CalcDamageReductionP(
	damage float64,
	attacker *Player,
	target *Player,
) float64 {
	return CalcDamageReduction(
		damage,
		attacker.GetIgnoreDefencePercent(),
		target.GetDefence(),
		target.GetDamageReduce(),
		target.GetDamageReducePercent(),
	)
}

/*
测试数据

	原始    衰减后
	85	84.47
	86	84.89
	87	85.29
	88	85.65
	89	86
	90	86.32
	91	86.63
	93	87.21
	94	87.48
	97	88.24
	98	88.48
	100	88.94
	101	89.16
	104	89.79
	106	90.19
	110	90.95
	114	91.66
	119	92.48
	121	92.8
	122	92.96
	123	93.11
	125	93.41
	126	93.56
	134	94.69
	138	95.23
	145	96.12
	146	96.24
	151	96.85
	154	97.2
	158	97.66
	159	97.77
	...
*/
var staticDmgRedPct = map[float64]float64{
	0.85: 0.8447, 0.86: 0.8489, 0.87: 0.8529, 0.88: 0.8565, 0.89: 0.86,
	0.90: 0.8632, 0.91: 0.8663, 0.92: 0.8692, 0.93: 0.8721, 0.94: 0.8748, 0.95: 0.8773, 0.96: 0.8799, 0.97: 0.8824, 0.98: 0.8848, 0.99: 0.8871,
	1: 0.8894, 1.01: 0.8916, 1.02: 0.8937, 1.03: 0.8958, 1.04: 0.8979, 1.06: 0.9019, 1.07: 0.9038, 1.08: 0.9057, 1.09: 0.9076,
	1.1: 0.9095, 1.11: 0.9112, 1.12: 0.9130, 1.13: 0.9148, 1.14: 0.9166, 1.15: 0.9182, 1.16: 0.9199, 1.17: 0.9215, 1.18: 0.9232, 1.19: 0.9248,
	1.21: 0.928, 1.22: 0.9296, 1.23: 0.9311, 1.25: 0.9341, 1.26: 0.9356,
}

// 静态免伤衰减公式未知,暂时用线性插值近似代替
func CalcStaticDamageReducePercent(drpct float64) float64 {
	if drpct <= 0.84 {
		return drpct
	}
	drpct = math.Trunc(drpct*100) / 100
	if drpct <= 1.26 {
		return staticDmgRedPct[drpct]
	}
	delta := staticDmgRedPct[126] - staticDmgRedPct[125]
	return math.Max(staticDmgRedPct[126]+(drpct-1.26)/0.01*delta, 1-1e-5)
}

// 绝招冷却时间最小3秒，神照1秒
func CalcCoolDown(time int, cdr float64, cdrpct float64, min float64) int {
	t := Ms2Sec(time)
	t = math.Max((t-cdr)*(1-cdrpct), min)
	return Sec2Ms(t)
}

func CalcCoolDownP(time int, player *Player) int {
	return CalcCoolDown(time, player.GetCoolDownReduce(), player.GetCoolDownReducePercent(), 3)
}

func CalcCoolDownPWithouts(time int, player *Player, withouts ...string) int {
	return CalcCoolDown(time, player.GetCoolDownReduceExclude(withouts...), player.GetCoolDownReducePercentExclude(withouts...), 3)
}

// 绝招释放时间
func CalcCastTime(time int, ctr float64, ctrpct float64) int {
	t := Ms2Sec(time)
	t = math.Max((t-ctr)*(1-ctrpct), 0)
	return Sec2Ms(t)
}

func CalcCastTimeP(time int, player *Player) int {
	return CalcCastTime(time, player.GetCastTimeReduce(), player.GetCastTimeReducePercent())
}

// 计算负面持续时间
func CalcDebuffDuration(duration int, nr float64) int {
	v := Ms2Sec(duration) * (1 - nr)
	if v < 0 {
		v = 0
	}
	return Sec2Ms(v)
}

func CalcDebuffDurationP(duration int, target *Player) int {
	return CalcDebuffDuration(duration, target.GetNegativeResistPercent())
}

// 计算内力消耗
func CalcMPCost(cost float64, mpcr float64, mpcrpct float64) float64 {
	return math.Max(0, (cost-mpcr)*(1-mpcrpct))
}

func CalcMPCostP(cost float64, caster *Player) float64 {
	return CalcMPCost(cost, caster.GetMPCostReduce(), caster.GetMPCostReducePercent())
}

const (
	ATTACK_MUST_HIT int = iota
	ATTACK_IGNORE_MAIN_DAMAGE
	ATTACK_IGNORE_FORCE_DAMAGE_ADD
	ATTACK_IGNORE_FORCE_SHOCK
	ATTACK_IGNORE_DAMAGE_REFLECTION
	ATTACK_IGNORE_LIFE_LEECHING
	ATTACK_DO_TRUE_DAMAGE
	ATTACK_IS_MAIN_HIT
	ATTACK_IS_EXTRA_HIT
	ATTACK_IS_PERFORM
	ATTACK_IS_WEAPON_HIT
	ATTACK_IS_UNARMED_HIT
	ATTACK_IS_THROWING_HIT
	ATTACK_FLAG_COUNT
)

func GenericAttack(ctx *CombatContext) {
	var (
		hitResult int
		rawDamage float64 = 0
		attacker          = ctx.GetAttacker()
		target            = ctx.GetTarget()
	)

	defer func() {
		if attacker.IsAlive() && ctx.IsMainAttack() {
			attacker.AttackCD.Start(attacker.GetAttackSpeed())
		}
	}()

	// 命中判定
	attackerHit := attacker.GetHit(ctx.GetModifiers())
	if ctx.IsPerformAttack() {
		attackerHit *= ctx.GetPerformHitRate()
	}
	targetDodge := target.GetDodge()
	targetParry := target.GetParry()
	if ctx.GetModifier("绝对命中") > 0 {
		hitResult = Hit_Result_Hit
	} else if attacker.GetModifier("绝对命中") > 0 {
		hitResult = Hit_Result_Hit
	} else if target.State.Faint {
		hitResult = Hit_Result_Hit
	} else {
		if attacker.State.Blind {
			attackerHit = 0
		}
		if target.State.Fixed {
			targetDodge = 0
		}
		if target.State.Busy {
			targetParry = 0
		}
		hitResult = HitCheck(attacker.Roll, attackerHit, targetDodge, targetParry)
	}
	ctx.SetHitCheck(attackerHit, targetDodge, targetParry)

	// 闪避
	if hitResult == Hit_Result_Dodge {
		ctx.SetDodge()
	}

	// 招架
	if hitResult == Hit_Result_Parry {
		ctx.SetParry()
	}

	// 命中
	if hitResult == Hit_Result_Hit {
		ctx.SetHit()
		ProcessHitTrigger(ctx)

		dmg := attacker.GetAttack()
		// 绝招增伤
		if ctx.IsPerformAttack() {
			dmg += ctx.GetPerformDamageAdd()
			dmg *= ctx.GetPerformDamageRate()
		}
		// 部位
		partHit := attacker.Roll.Intn(len(BodyPart_DamageModifier))
		dmg *= BodyPart_DamageModifier[partHit]
		// 暴击
		// isCrit := false
		// 终伤
		dmg *= 1 + attacker.GetFinalDamagePercent(ctx.GetModifiers())
		// 绝招附加伤害
		if ctx.IsPerformAttack() {
			dmg += ctx.GetPerformDamageAppend()
		}
		rawDamage += dmg
		ctx.SetDamageMain(rawDamage)
	}

	// 内力附加
	rawDamage += CalcForceDamageAdd(ctx)
	// 伤害反射
	ProcessDamageReflection(ctx)
	if !attacker.IsAlive() {
		return
	}
	// 造成伤害
	ProcessDamageApply(ctx)
	if !attacker.IsAlive() || !target.IsAlive() {
		return
	}

	// 吸血
	// 自创战神
	if pZhanShen := attacker.Skills.Passives.ZhanShen; pZhanShen != nil {
		leech := pZhanShen.Leech(ctx.GetDamageFinal())
		ctx.AddHpLeech("自创.战神", leech)
		ctx.AddMpLeech("自创.战神", leech)
	}
	// 招架反击
}

func CalcForceDamageAdd(ctx *CombatContext) float64 {
	var (
		attacker = ctx.GetAttacker()
		add      = 0.

		aZhanShen = attacker.Skills.Passives.ZhanShen
		aRuMo     = attacker.Skills.Passives.RuMo
		aWDRuMo   = attacker.Skills.Passives.WanDaoRuMo
	)
	// 自创.战神:所有攻击皆可附加
	if aZhanShen != nil {
		v := aZhanShen.Append()
		ctx.AddDamageAdd("自创.战神", v)
		add += v
	}
	// 自创.入魔:仅作用于入魔被动所属武功类型的普通攻击,例,入魔剑仅剑法普攻可附加入魔内力伤害
	if aRuMo != nil && !ctx.IsPerformAttack() {
		if (ctx.IsUnarmedAttack() && !aRuMo.IsWeapon()) ||
			(ctx.IsWeaponAttack() && aRuMo.IsWeapon()) {
			v := aRuMo.DamageAdd()
			ctx.AddDamageAdd("自创.入魔", v)
			add += v
		}
	}
	// 圆月弯刀.入魔:仅作用于刀法普攻和圆月弯刀.魔刀
	if aWDRuMo != nil {
		v := 0.
		if ctx.IsWeaponAttack() {
			v = aWDRuMo.DamageAdd(true)
		} else if ctx.IsPerformAttack() && ctx.GetSource() == "圆月弯刀.魔刀" {
			v = aWDRuMo.DamageAdd(false)
		}
		if v > 0 {
			ctx.AddDamageAdd("圆月弯刀.入魔", v)
			add += v
		}
	}
	// 先天功.纯阳气:所有攻击皆可附加
	if m := attacker.GetModifier("先天功.内力附加%"); m > 0 {
		v := m * attacker.GetMPMax()
		ctx.AddDamageAdd("先天功.纯阳气", v)
		add += v
	}
	// 无念禅功.无念:未知;需要被动支持,目前擂台无用
	/*
		if m := attacker.GetModifier("无念禅功.内力附加%"); m > 0 {
			v := m * attacker.GetMPMax()
			ctx.AddDamageAdd("无念禅功.无用", v)
			add += v
		}
	*/
	return add
}

func ProcessDamageApply(ctx *CombatContext) {
	var (
		attacker = ctx.GetAttacker()
		victim   = ctx.GetTarget()
		damage   = ctx.GetDamageMain() + ctx.GetDamageAdd()
		final    = 0.

		vBuMie   = victim.Skills.Passives.BuMie
		vJianXin = victim.Skills.Passives.JianXin
	)

	if damage <= 0 {
		return
	}
	// 自创.不灭免疫伤害
	if vBuMie != nil && vBuMie.IsActivate() {
		ctx.SetDamageImmunity("自创.不灭", true)
		return
	}
	// 太极图免疫伤害
	if victim.HasBuff("太极图.太极") {
		ctx.SetDamageImmunity("太极图.太极", true)
		return
	}
	// 真实伤害
	// if ctx.IsTrueDamge() {
	//		victim.TakeDamage(attacker, damage)
	//		ctx.SetDamageFinal(damage)
	//		return
	//}
	// 自创.剑心+慈航剑典.剑心免疫伤害
	if vJianXin != nil && victim.HasBuff("慈航剑典.剑心") {
		// 昏迷时"绝对"招架不生效
		if !victim.State.Faint {
			ctx.SetDamageImmunity("自创.剑心", true)
			return
		}
	}
	// 计算免伤
	final = CalcDamageReduction(
		damage,
		attacker.GetIgnoreDefencePercent(ctx.GetModifiers()), // 覆雨附加1000%忽防
		victim.GetDefence(),
		victim.GetDamageReduce(),
		victim.GetDamageReducePercent(),
	)
	// 自创混沌+长生诀.混沌吸收伤害
	if vBuMie != nil && victim.HasBuff("长生诀.混沌") {
		dmg, absorb := vBuMie.Absorb(final)
		// 混沌工作机制,如果触发混沌,先回复等于(伤害值-阈值)的气血,回复量不超过气血上限,然后造成于阈值的伤害。
		// 这里简化成:如果触发混沌,造成等于(2*阈值-伤害值)的伤害.
		if absorb > 0 {
			dmg -= absorb
			ctx.AddDamageAbsort("长生诀.混沌", absorb)
		}
		final = dmg
	}
	// 般若龙象功.象驱式吸收伤害

	// 自创.不灭触发
	if vBuMie != nil && vBuMie.Activate(final) {
		final = 0
		ctx.SetDamageImmunity("长生诀.不灭", true)
	}
	victim.TakeDamage(attacker, final)
	ctx.SetDamageFinal(final)
}

func ProcessDamageReflection(ctx *CombatContext) {
	var (
		attacker = ctx.GetAttacker()
		target   = ctx.GetTarget()
		damage   = 0.

		vFanZhen = target.Skills.Passives.FanZhen
		aBuMie   = attacker.Skills.Passives.BuMie
	)

	if ctx.IsDodge() {
		return
	}

	// 千蛛万毒手.万蛊噬天

	// 燃木刀法.护体真焰

	// 九阳神功

	// 软猬甲

	// 自创.反震
	// 狮子吼等内功攻击不触发反震?
	if vFanZhen != nil && !ctx.IsForceAttack() {
		dmg := CalcDamageReductionP(vFanZhen.Damage(), target, attacker)

		{
			if aBuMie != nil && aBuMie.IsActivate() {
				dmg = 0
			}
			// 太极图免疫伤害
			if dmg > 0 && attacker.HasBuff("太极图.太极") {
				dmg = 0
			}

			// 自创混沌+长生诀.混沌吸收伤害
			if dmg > 0 && aBuMie != nil && attacker.HasBuff("长生诀.混沌") {
				real, absorb := aBuMie.Absorb(dmg)
				if absorb > 0 {
					attacker.AddHP(absorb)
				}
				dmg = real
			}
			// 般若龙象功.象驱式吸收伤害

			// 自创.不灭触发
			if dmg > 0 && aBuMie != nil && aBuMie.Activate(dmg) {
				dmg = 0
			}
		}

		ctx.AddDamageReflect("自创.反震", dmg)
		damage += dmg
	}

	attacker.TakeDamage(target, damage)
}

func ProcessLeech(ctx *CombatContext) {
	var (
		attacker = ctx.GetAttacker()
		target   = ctx.GetTarget()
		damage   = ctx.GetDamageFinal()

		vZhanShen = target.Skills.Passives.ZhanShen
	)
	// 自创.战神
	if vZhanShen != nil {
		value := vZhanShen.Leech(damage)
		attacker.AddHP(value)
		attacker.AddMP(value)
		ctx.AddHpLeech("自创.战神", value)
		ctx.AddMpLeech("自创.战神", value)
	}
	// 自创.吸血
}

func ProcessHitTrigger(ctx *CombatContext) {
	var (
		attacker = ctx.GetAttacker()
		target   = ctx.GetTarget()

		aYinDao = attacker.Weapon.Passives.YinDao
		aFuYu   = attacker.Skills.Passives.FuYu
	)
	// 鹰刀
	if ctx.IsWeaponAttack() {
		if aYinDao != nil {
			aYinDao.Stun(target)
		}
	}

	if aFuYu != nil && ctx.IsWeaponAttack() {
		chance := 0.
		if ctx.IsMainAttack() {
			chance = aFuYu.GetChance()
		} else if ctx.IsExtraAttack() && ctx.GetSource() == "覆雨剑法.剑雨" {
			chance = aFuYu.GetChance() * 2
		} else if ctx.IsPerformAttack() && ctx.GetSource() == "覆雨剑法.剑罡" {
			chance = 1
		}

		if attacker.Roll.Float64() <= chance {
			ctx.AddModifier(Modifier{"忽防%", 100})
		}
	}
}

func ProcessAttackDone(ctx *CombatContext) {
	var (
		attacker = ctx.GetAttacker()
		//target   = ctx.GetTarget()

		aDuGu = attacker.Skills.Passives.DuGu
	)

	if ctx.IsHit() {
		if aDuGu != nil && attacker.HasBuff("独孤剑诀.剑来") && ctx.IsWeaponAttack() {
			extra := ctx.AddExtraAttack()
			extra.SetMustHit().SetTrueDamage().SetDamageFinal(1)
		}
	}
}
