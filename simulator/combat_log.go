package simulator

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/abs0293/wsmud-simc/simulator/log_pb"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/protobuf/proto"
)

type CombatLog struct {
	Timestamp      float64
	Attacker       *Player                    `json:",omitempty" yaml:",omitempty"`
	Target         *Player                    `json:",omitempty" yaml:",omitempty"`
	Source         string                     `json:",omitempty" yaml:",omitempty"`
	IsExtraHit     bool                       `json:",omitempty" yaml:",omitempty"`
	MustHit        bool                       `json:",omitempty" yaml:",omitempty"`
	IsTrueDamage   bool                       `json:",omitempty" yaml:",omitempty"`
	HitResult      int                        `json:",omitempty" yaml:",omitempty"`
	Damage         CombatLog_Damage           `json:",omitempty" yaml:",omitempty"`
	HitTriggers    []CombatLog_Trigger        `json:",omitempty" yaml:",omitempty"`
	ReflectDamages []CombatLog_Reflect_Damage `json:",omitempty" yaml:",omitempty"`
	LifeLeecheds   []CombatLog_Life_Leeched   `json:",omitempty" yaml:",omitempty"`
	ExtraHits      []CombatLog                `json:",omitempty" yaml:",omitempty"`
}

func (l CombatLog) Dump() string {
	out := map[string]interface{}{}
	mapstructure.Decode(&l, &out)
	for k, v := range out {
		vv := reflect.ValueOf(v)
		if vv.Kind() == reflect.Slice || vv.Kind() == reflect.Array {
			if vv.Len() == 0 {
				delete(out, k)
			}
		}
	}
	b, _ := json.MarshalIndent(out, "", "  ")
	// b, _ := yaml.Marshal(out)
	return string(b)
}

func (l CombatLog) hitSourceString() string {
	if l.IsExtraHit {
		return fmt.Sprintf("额外攻击(%s)", l.Source)
	}
	return l.Source
}

func (l CombatLog) String() string {
	str := []string{
		fmt.Sprintf("%s:->%s",
			l.hitSourceString(), l.Target.Name,
		),
		"," + Hit_Result_String[l.HitResult],
	}

	if l.HitResult != Hit_Result_Dodge {
		str = append(str, l.Damage.String())
		if len(l.ReflectDamages) > 0 {
			v := 0.0
			for _, r := range l.ReflectDamages {
				v += r.Reflected
			}
			str = append(str, fmt.Sprintf(",受到反弹伤害:%.3f", v))
		}
		if len(l.LifeLeecheds) > 0 {
			v := 0.0
			for _, ll := range l.LifeLeecheds {
				v += ll.Leeched
			}
			str = append(str, fmt.Sprintf(",吸取生命:%.3f", v))
		}
		if len(l.ExtraHits) > 0 {
			str = append(str, "\n")
			for _, e := range l.ExtraHits {
				str = append(str, fmt.Sprintf("    %s\n", e.String()))
			}
		}
	}
	if !l.IsExtraHit && len(l.ExtraHits) == 0 {
		str = append(str, "\n")
	}
	return strings.Join(str, "")
}

type CombatLog_Damage struct {
	BaseDamage  float64
	Appendeds   []CombatLog_Damage_Appended
	Immunizeds  []CombatLog_Damage_Immunized
	Absorbeds   []CombatLog_Damage_Absorbed
	FinalDamage float64
}

func (d CombatLog_Damage) String() string {
	if len(d.Immunizeds) > 0 {
		return fmt.Sprintf(",伤害免疫:%s", d.Immunizeds[0].Source)
	}

	v := d.BaseDamage
	for _, a := range d.Appendeds {
		v += a.Appended
	}

	return fmt.Sprintf(",造成伤害:%.3f,原始伤害:%.3f", d.FinalDamage, v)
}

type CombatLog_Damage_Appended struct {
	Source   string
	Appended float64
}

type CombatLog_Damage_Absorbed struct {
	Source   string
	Absorbed float64
}

type CombatLog_Damage_Immunized struct {
	Source string
}

type CombatLog_Reflect_Damage struct {
	Source    string
	Reflected float64
}

type CombatLog_Life_Leeched struct {
	Source  string
	Leeched float64
}

type CombatLog_Trigger struct {
	Source string
}

func ReadLog(data []byte, verbose bool) (string, error) {
	var (
		log = &log_pb.Log{}
	)

	err := proto.Unmarshal(data, log)
	if err != nil {
		return "", err
	}

	if log.Attack != nil {
		return ReadAttackLog(log, verbose), nil
	}

	return "", nil
}

func ReadAttackLog(log *log_pb.Log, verbose bool) string {
	attack := log.Attack

	str := []string{
		fmt.Sprintf("[%4d][%d][%s]%s->%s", log.SerialNumber, log.GetTimestamp(), attack.Source, log.Player, attack.Target),
	}
	aclass := "攻击类型:"
	switch attack.Kind {
	case log_pb.Log_Attack_Kind_Main:
		aclass += "普通"
	case log_pb.Log_Attack_Kind_Extra:
		aclass += "额外"
	case log_pb.Log_Attack_Kind_Perform:
		aclass += "绝招"
	default:
		aclass += "未知"
	}
	akind := "类型:"
	switch attack.Class {
	case log_pb.Log_Attack_Class_Force:
		akind += "内力"
	case log_pb.Log_Attack_Class_Throwing:
		akind += "暗器"
	case log_pb.Log_Attack_Class_Weapone:
		akind += "武器"
	case log_pb.Log_Attack_Class_Unarmed:
		akind += "拳脚"
	default:
		akind += "未知"
	}
	if verbose {
		str = append(str, aclass, akind)
	}

	hit := "结果:"
	switch log.Attack.HitCheck.Result {
	case log_pb.Log_Attack_HitCheck_Result_Hit:
		hit += "命中"
	case log_pb.Log_Attack_HitCheck_Result_Parry:
		hit += "招架"
	case log_pb.Log_Attack_HitCheck_Result_Dodge:
		hit += "躲闪"
	default:
		hit += "未知"
	}
	str = append(str, hit)

	if attack.DamageImmunity != nil {
		imm := "免疫伤害"
		if verbose {
			imm += ":" + attack.DamageImmunity.Source
		}
		str = append(str, imm)
	} else {
		str = append(str, fmt.Sprintf("伤害:%.3f", attack.DamageFinal))
		if verbose {
			dmg := attack.DamageMain
			for _, add := range attack.DamageAdd {
				dmg += add.Value
			}
			str = append(str, fmt.Sprintf("原始伤害:%.3f", dmg))
		}
	}

	if attack.DamageReflect != nil {
		d := 0.
		for _, r := range attack.DamageReflect {
			d += r.Value
		}
		str = append(str, fmt.Sprintf("反伤:%.3f", d))
	}

	hpL := 0.
	for _, l := range attack.HpLeech {
		hpL += l.Value
	}
	if hpL > 0 {
		str = append(str, fmt.Sprintf("吸血:%.3f", hpL))
	}

	mpL := 0.
	for _, l := range attack.MpLeech {
		mpL += l.Value
	}
	if mpL > 0 {
		str = append(str, fmt.Sprintf("吸内:%.3f", mpL))
	}

	if verbose {
		if len(attack.Modifier) > 0 {
			mod := "修饰器:"
			for _, m := range attack.Modifier {
				mod += fmt.Sprintf("{%s:%.3f}", m.Name, m.Value)
			}
			str = append(str, mod)
		}
	}

	return strings.Join(str, ",")
}

func DumpAttackLog(log *log_pb.Log) string {
	return log.String()
}
