package simulator

import (
	"github.com/abs0293/wsmud-simc/simulator/log_pb"
	"google.golang.org/protobuf/proto"
)

type CombatContext struct {
	attacker     *Player
	target       *Player
	log          *log_pb.Log
	done         bool
	extraAttacks []*CombatContext
}

func (ctx *CombatContext) SetAttacker(attacker *Player) *CombatContext {
	ctx.attacker = attacker
	ctx.log.Player = attacker.Name
	return ctx
}

func (ctx *CombatContext) GetAttacker() *Player {
	return ctx.attacker
}

func (ctx *CombatContext) SetTarget(target *Player) *CombatContext {
	ctx.target = target
	ctx.log.Attack.Target = target.Name
	return ctx
}

func (ctx *CombatContext) GetTarget() *Player {
	return ctx.target
}

func (ctx *CombatContext) SetTimestamp(ts int) *CombatContext {
	ctx.log.Timestamp = int32(ts)
	return ctx
}

func (ctx *CombatContext) GetTimestamp() int {
	return int(ctx.log.Timestamp)
}

func (ctx *CombatContext) SetSource(source string) *CombatContext {
	ctx.log.Attack.Source = source
	return ctx
}

func (ctx *CombatContext) GetSource() string {
	return ctx.log.Attack.Source
}

func (ctx *CombatContext) SetMainAttack() *CombatContext {
	ctx.log.Attack.Kind = log_pb.Log_Attack_Kind_Main
	return ctx
}

func (ctx *CombatContext) IsMainAttack() bool {
	return ctx.log.Attack.Kind == log_pb.Log_Attack_Kind_Main
}

func (ctx *CombatContext) SetExtraAttack() *CombatContext {
	ctx.log.Attack.Kind = log_pb.Log_Attack_Kind_Extra
	return ctx
}

func (ctx *CombatContext) IsExtraAttack() bool {
	return ctx.log.Attack.Kind == log_pb.Log_Attack_Kind_Extra
}

func (ctx *CombatContext) SetPerformAttack() *CombatContext {
	ctx.log.Attack.Kind = log_pb.Log_Attack_Kind_Perform
	return ctx
}

func (ctx *CombatContext) IsPerformAttack() bool {
	return ctx.log.Attack.Kind == log_pb.Log_Attack_Kind_Perform
}

func (ctx *CombatContext) SetWeaponAttack() *CombatContext {
	ctx.log.Attack.Class = log_pb.Log_Attack_Class_Weapone
	return ctx
}

func (ctx *CombatContext) IsWeaponAttack() bool {
	return ctx.log.Attack.Class == log_pb.Log_Attack_Class_Weapone
}

func (ctx *CombatContext) SetUnarmedAttack() *CombatContext {
	ctx.log.Attack.Class = log_pb.Log_Attack_Class_Unarmed
	return ctx
}

func (ctx *CombatContext) IsUnarmedAttack() bool {
	return ctx.log.Attack.Class == log_pb.Log_Attack_Class_Unarmed
}

func (ctx *CombatContext) SetThrowingAttack() *CombatContext {
	ctx.log.Attack.Class = log_pb.Log_Attack_Class_Throwing
	return ctx
}

func (ctx *CombatContext) IsThrowingAttack() bool {
	return ctx.log.Attack.Class == log_pb.Log_Attack_Class_Throwing
}

func (ctx *CombatContext) SetForceAttack() *CombatContext {
	ctx.log.Attack.Class = log_pb.Log_Attack_Class_Force
	return ctx
}

func (ctx *CombatContext) IsForceAttack() bool {
	return ctx.log.Attack.Class == log_pb.Log_Attack_Class_Force
}

func (ctx *CombatContext) SetHitCheck(hit, dodge, parry float64) *CombatContext {
	ctx.log.Attack.HitCheck = &log_pb.Log_Attack_HitCheck{
		AttackerHit: hit,
		TargetDodge: dodge,
		TargetParry: parry,
	}
	return ctx
}

func (ctx *CombatContext) SetHit() *CombatContext {
	if ctx.log.Attack.HitCheck == nil {
		ctx.log.Attack.HitCheck = &log_pb.Log_Attack_HitCheck{}
	}
	ctx.log.Attack.HitCheck.Result = log_pb.Log_Attack_HitCheck_Result_Hit
	return ctx
}

func (ctx *CombatContext) IsHit() bool {
	return ctx.log.Attack.HitCheck.Result == log_pb.Log_Attack_HitCheck_Result_Hit
}

func (ctx *CombatContext) SetDodge() *CombatContext {
	if ctx.log.Attack.HitCheck == nil {
		ctx.log.Attack.HitCheck = &log_pb.Log_Attack_HitCheck{}
	}
	ctx.log.Attack.HitCheck.Result = log_pb.Log_Attack_HitCheck_Result_Dodge
	return ctx
}

func (ctx *CombatContext) IsDodge() bool {
	return ctx.log.Attack.HitCheck.Result == log_pb.Log_Attack_HitCheck_Result_Dodge
}

func (ctx *CombatContext) SetParry() *CombatContext {
	if ctx.log.Attack.HitCheck == nil {
		ctx.log.Attack.HitCheck = &log_pb.Log_Attack_HitCheck{}
	}
	ctx.log.Attack.HitCheck.Result = log_pb.Log_Attack_HitCheck_Result_Parry
	return ctx
}

func (ctx *CombatContext) IsParry() bool {
	return ctx.log.Attack.HitCheck.Result == log_pb.Log_Attack_HitCheck_Result_Parry
}

func (ctx *CombatContext) SetTrueDamage() *CombatContext {
	ctx.log.Attack.TrueDamage = true
	return ctx
}

func (ctx *CombatContext) IsTrueDamge() bool {
	return ctx.log.Attack.TrueDamage
}

func (ctx *CombatContext) SetDamageMain(value float64) *CombatContext {
	ctx.log.Attack.DamageMain = value
	return ctx
}

func (ctx *CombatContext) GetDamageMain() float64 {
	return ctx.log.Attack.DamageMain
}

func (ctx *CombatContext) AddDamageAdd(source string, value float64) *CombatContext {
	ctx.log.Attack.DamageAdd = append(
		ctx.log.Attack.DamageAdd,
		&log_pb.SourcedDouble{Source: source, Value: value},
	)
	return ctx
}

func (ctx *CombatContext) GetDamageAdd() float64 {
	add := 0.
	for _, d := range ctx.log.Attack.DamageAdd {
		add += d.Value
	}
	return add
}

func (ctx *CombatContext) AddDamageAbsort(source string, value float64) *CombatContext {
	ctx.log.Attack.DamageAbsort = append(
		ctx.log.Attack.DamageAbsort,
		&log_pb.SourcedDouble{Source: source, Value: value},
	)
	return ctx
}

func (ctx *CombatContext) GetDamageAbsort() float64 {
	ab := 0.
	for _, d := range ctx.log.Attack.DamageAbsort {
		ab += d.Value
	}
	return ab
}

func (ctx *CombatContext) AddDamageReflect(source string, value float64) *CombatContext {
	ctx.log.Attack.DamageReflect = append(
		ctx.log.Attack.DamageReflect,
		&log_pb.SourcedDouble{Source: source, Value: value},
	)
	return ctx
}

func (ctx *CombatContext) GetDamageReflect() float64 {
	rd := 0.
	for _, d := range ctx.log.Attack.DamageReflect {
		rd += d.Value
	}
	return rd
}

func (ctx *CombatContext) AddHpLeech(source string, value float64) *CombatContext {
	ctx.log.Attack.HpLeech = append(
		ctx.log.Attack.HpLeech,
		&log_pb.SourcedDouble{Source: source, Value: value},
	)
	return ctx
}

func (ctx *CombatContext) GetHpLeech() float64 {
	hp := 0.
	for _, d := range ctx.log.Attack.HpLeech {
		hp += d.Value
	}
	return hp
}

func (ctx *CombatContext) AddMpLeech(source string, value float64) *CombatContext {
	ctx.log.Attack.MpLeech = append(
		ctx.log.Attack.MpLeech,
		&log_pb.SourcedDouble{Source: source, Value: value},
	)
	return ctx
}

func (ctx *CombatContext) GetMpLeech() float64 {
	mp := 0.
	for _, d := range ctx.log.Attack.MpLeech {
		mp += d.Value
	}
	return mp
}

func (ctx *CombatContext) SetDamageImmunity(source string, value bool) *CombatContext {
	ctx.log.Attack.DamageImmunity = &log_pb.SourcedBool{
		Source: source,
		Value:  value,
	}
	return ctx
}

func (ctx *CombatContext) IsDamageImmunity() bool {
	if ctx.log.Attack.DamageImmunity == nil {
		return false
	}
	return ctx.log.Attack.DamageImmunity.Value
}

func (ctx *CombatContext) SetDamageFinal(value float64) *CombatContext {
	ctx.log.Attack.DamageFinal = value
	return ctx
}

func (ctx *CombatContext) GetDamageFinal() float64 {
	return ctx.log.Attack.DamageFinal
}

func (ctx *CombatContext) IsDamaged() bool {
	return ctx.log.Attack.DamageFinal > 0
}

func (ctx *CombatContext) SetCooldown(value int) *CombatContext {
	ctx.log.Attack.Cooldown = int32(value)
	return ctx
}

func (ctx *CombatContext) GetCoolDown() int {
	return int(ctx.log.Attack.Cooldown)
}

func (ctx *CombatContext) AddModifier(mods ...Modifier) *CombatContext {
	for _, mod := range mods {
		ctx.log.Attack.Modifier = append(ctx.log.Attack.Modifier, &log_pb.Modifier{
			Name:  mod.Name,
			Value: mod.Value,
		})
	}
	return ctx
}

func (ctx *CombatContext) GetModifiers() []Modifier {
	mods := []Modifier{}
	for _, mod := range ctx.log.Attack.Modifier {
		mods = append(mods, Modifier{mod.Name, mod.Value})
	}
	return mods
}

func (ctx *CombatContext) SetPerformDamageRate(value float64) *CombatContext {
	ctx.log.Attack.Modifier = append(ctx.log.Attack.Modifier, &log_pb.Modifier{
		Name:  "绝招.伤害倍率%",
		Value: value,
	})
	return ctx
}

func (ctx *CombatContext) GetPerformDamageRate() float64 {
	set := false
	rate := 0.
	for _, mod := range ctx.log.Attack.Modifier {
		if mod.Name == "绝招.伤害倍率%" {
			set = true
			rate += mod.Value
		}
	}
	if !set {
		return 1
	}
	return rate
}

func (ctx *CombatContext) SetPerformDamageAdd(value float64) *CombatContext {
	ctx.log.Attack.Modifier = append(ctx.log.Attack.Modifier, &log_pb.Modifier{
		Name:  "绝招.基础伤害d",
		Value: value,
	})
	return ctx
}

func (ctx *CombatContext) GetPerformDamageAdd() float64 {
	return ctx.GetModifier("绝招.基础伤害d")
}

func (ctx *CombatContext) SetPerformDamageAppend(value float64) *CombatContext {
	ctx.log.Attack.Modifier = append(ctx.log.Attack.Modifier, &log_pb.Modifier{
		Name:  "绝招.附加伤害d",
		Value: value,
	})
	return ctx
}

func (ctx *CombatContext) GetPerformDamageAppend() float64 {
	return ctx.GetModifier("绝招.附加伤害d")
}

func (ctx *CombatContext) SetPerformHitRate(value float64) *CombatContext {
	ctx.log.Attack.Modifier = append(ctx.log.Attack.Modifier, &log_pb.Modifier{
		Name:  "绝招.命中倍率%",
		Value: value,
	})
	return ctx
}

func (ctx *CombatContext) GetPerformHitRate() float64 {
	set := false
	rate := 0.
	for _, mod := range ctx.log.Attack.Modifier {
		if mod.Name == "绝招.命中倍率%" {
			set = true
			rate += mod.Value
		}
	}
	if !set {
		return 1
	}
	return rate
}

func (ctx *CombatContext) SetMustHit() *CombatContext {
	ctx.log.Attack.Modifier = append(ctx.log.Attack.Modifier, &log_pb.Modifier{
		Name:  "绝对命中",
		Value: 1,
	})
	return ctx
}

func (ctx *CombatContext) IsMustHit() bool {
	return ctx.GetModifier("绝对命中") > 0
}

func (ctx *CombatContext) GetModifier(name string) float64 {
	v := 0.
	for _, m := range ctx.log.Attack.Modifier {
		if m.Name == name {
			v += m.Value
		}
	}
	return v
}

func (ctx *CombatContext) ToProtobuf() ([]byte, error) {
	return proto.Marshal(ctx.log)
}

func (ctx *CombatContext) AddExtraAttack() *CombatContext {
	ext := NewCombatContext(ctx.attacker, ctx.target, ctx.GetTimestamp())
	ctx.extraAttacks = append(ctx.extraAttacks, ext)
	ext.SetExtraAttack()
	return ext
}

func (ctx *CombatContext) Done() {
	if ctx.done {
		return
	}
	ctx.done = true
	if attacker := ctx.attacker; attacker != nil {
		attacker.Log(ctx.log)
		for _, ext := range ctx.extraAttacks {
			ext.Done()
		}
	}
}

func NewCombatContext(attacker *Player, target *Player, ts int, mods ...Modifier) *CombatContext {
	ctx := &CombatContext{
		log: &log_pb.Log{
			Attack: &log_pb.Log_Attack{},
		},
	}
	return ctx.
		SetAttacker(attacker).
		SetTarget(target).
		SetTimestamp(ts).
		AddModifier(mods...)
}

type RunContext struct {
	attacker *Player
	target   *Player
	skipPf   bool
	done     bool
	attack   []*CombatContext
	log      *log_pb.Log
}

func (ctx *RunContext) SetSkipPreflight(skip bool) *RunContext {
	ctx.skipPf = skip
	return ctx
}

func (ctx *RunContext) IsSkipPrefight() bool {
	return ctx.skipPf
}

func (ctx *RunContext) SetAttacker(attacker *Player) *RunContext {
	ctx.attacker = attacker
	ctx.log.Player = attacker.Name
	return ctx
}

func (ctx *RunContext) GetAttacker() *Player {
	return ctx.attacker
}

func (ctx *RunContext) SetTrigger(trigger string) *RunContext {
	ctx.log.Run.Trigger = trigger
	return ctx
}

func (ctx *RunContext) GetTrigger() string {
	return ctx.log.Run.Trigger
}

func (ctx *RunContext) SetTarget(target *Player) *RunContext {
	ctx.target = target
	ctx.log.Run.Target = target.Name
	return ctx
}

func (ctx *RunContext) GetTarget() *Player {
	return ctx.target
}

func (ctx *RunContext) SetTimestamp(ts int) *RunContext {
	ctx.log.Timestamp = int32(ts)
	return ctx
}

func (ctx *RunContext) GetTimestamp() int {
	return int(ctx.log.Timestamp)
}

func (ctx *RunContext) SetName(name string) *RunContext {
	ctx.log.Run.Name = name
	return ctx
}

func (ctx *RunContext) GetName() string {
	return ctx.log.Run.Name
}

func (ctx *RunContext) SetCastTime(value int) *RunContext {
	ctx.log.Run.CastTime = int32(value)
	return ctx
}

func (ctx *RunContext) GetCastTime() int {
	return int(ctx.log.Run.CastTime)
}

func (ctx *RunContext) SetCooldown(value int) *RunContext {
	ctx.log.Run.Cooldown = int32(value)
	return ctx
}

func (ctx *RunContext) GetCooldown() int {
	return int(ctx.log.Run.Cooldown)
}

func (ctx *RunContext) SetMPCost(value float64) *RunContext {
	ctx.log.Run.MpCost = value
	return ctx
}

func (ctx *RunContext) GetMPCost() float64 {
	return ctx.log.Run.MpCost
}

func (ctx *RunContext) SetFail(source string) *RunContext {
	ctx.log.Run.Fail = &log_pb.SourcedBool{
		Source: source,
		Value:  true,
	}
	return ctx
}

func (ctx *RunContext) IsFail() bool {
	if ctx.log.Run.Fail == nil {
		return false
	}
	return ctx.log.Run.Fail.Value
}

func (ctx *RunContext) AddAttack() *CombatContext {
	atk := NewCombatContext(ctx.attacker, ctx.target, ctx.GetTimestamp())
	atk.SetSource(ctx.GetName()).SetPerformAttack()
	ctx.attack = append(ctx.attack, atk)
	return atk
}

func (ctx *RunContext) Done() {
	if ctx.done {
		return
	}
	ctx.done = true
	if attacker := ctx.attacker; attacker != nil {
		attacker.Log(ctx.log)
		for _, atk := range ctx.attack {
			atk.Done()
		}
	}
}

func NewRunContext(attacker *Player, target *Player, ts int) *RunContext {
	ctx := &RunContext{
		log: &log_pb.Log{
			Run: &log_pb.Log_Run{},
		},
	}
	ctx.SetAttacker(attacker).SetTarget(target).SetTimestamp(ts)
	return ctx
}
