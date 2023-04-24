package simulator

import (
	"log"
	"reflect"

	"gopkg.in/yaml.v3"
)

const (
	WeaponTypeUnknown int = iota
	WeaponTypeUnarmed
	WeaponTypeSword
	WeaponTypeBlade
	WeaponTypeStaff
	WeaponTypeClub
)

func GetWeaponTypeValue(name string) int {
	switch name {
	case "空手":
		return 1
	case "剑":
		return 2
	case "刀":
		return 3
	case "杖":
		return 4
	case "棍":
		return 5
	default:
		return 0
	}
}

func GetWeaponTypeName(value int) string {
	switch value {
	case 1:
		return "空手"
	case 2:
		return "剑"
	case 3:
		return "刀"
	case 4:
		return "杖"
	case 5:
		return "棍"
	default:
		return "未知"
	}
}

type WeaponData struct {
	Name        string      `mapstructure:"名称" yaml:"名称,omitempty"`
	Type        string      `mapstructure:"类型" yaml:"类型,omitempty"`
	HP          float64     `mapstructure:"气血d" yaml:"气血d,omitempty"`
	HPPct       float64     `mapstructure:"气血%" yaml:"气血%,omitempty"`
	StrAdd      float64     `mapstructure:"臂力" yaml:"臂力,omitempty"`
	DexAdd      float64     `mapstructure:"身法" yaml:"身法,omitempty"`
	ConAdd      float64     `mapstructure:"根骨" yaml:"根骨,omitempty"`
	IntAdd      float64     `mapstructure:"悟性" yaml:"悟性,omitempty"`
	Attack      float64     `mapstructure:"攻击d" yaml:"攻击d,omitempty"`
	AttackPct   float64     `mapstructure:"攻击%" yaml:"攻击%,omitempty"`
	Defence     float64     `mapstructure:"防御d" yaml:"防御d,omitempty"`
	DefencePct  float64     `mapstructure:"防御%" yaml:"防御%,omitempty"`
	Hit         float64     `mapstructure:"命中d" yaml:"命中d,omitempty"`
	HitPct      float64     `mapstructure:"命中%" yaml:"命中%,omitempty"`
	Dodge       float64     `mapstructure:"闪避d" yaml:"闪避d,omitempty"`
	DodgePct    float64     `mapstructure:"闪避%" yaml:"闪避%,omitempty"`
	Parry       float64     `mapstructure:"招架d" yaml:"招架d,omitempty"`
	ParryPct    float64     `mapstructure:"招架%" yaml:"招架%,omitempty"`
	Speed       float64     `mapstructure:"攻速d" yaml:"攻速d,omitempty"`
	SpeedPct    float64     `mapstructure:"攻速%" yaml:"攻速%,omitempty"`
	CDR         float64     `mapstructure:"绝招冷却d" yaml:"绝招冷却d,omitempty"`
	CDRPct      float64     `mapstructure:"绝招冷却%" yaml:"绝招冷却%,omitempty"`
	DmR         float64     `mapstructure:"免伤d" yaml:"免伤d,omitempty"`
	DmRPct      float64     `mapstructure:"免伤%" yaml:"免伤%,omitempty"`
	IgDPct      float64     `mapstructure:"忽防%" yaml:"忽防%,omitempty"`
	FiDPct      float64     `mapstructure:"终伤%" yaml:"终伤%,omitempty"`
	CTR         float64     `mapstructure:"绝招释放d" yaml:"绝招释放d,omitempty"`
	CTRPct      float64     `mapstructure:"绝招释放%" yaml:"绝招释放%,omitempty"`
	MPCR        float64     `mapstructure:"内力消耗d" yaml:"内力消耗d,omitempty"`
	MPCRPct     float64     `mapstructure:"内力消耗%" yaml:"内力消耗%,omitempty"`
	NRPct       float64     `mapstructure:"负面抵抗%" yaml:"负面抵抗%,omitempty"`
	PassiveData PassiveData `mapstructure:"被动" yaml:"被动,omitempty"`
}

func (d WeaponData) buildModifiers() []Modifier {
	out := []Modifier{}
	rtData := reflect.TypeOf(d)
	rvData := reflect.ValueOf(d)

	for i := 0; i < rtData.NumField(); i++ {
		if rvData.Field(i).Kind() == reflect.Float64 {
			v := rvData.Field(i).Float()
			if v == 0 {
				continue
			}
			out = append(out, Modifier{
				Name:  string(rtData.Field(i).Tag.Get("mapstructure")),
				Value: rvData.Field(i).Float(),
			})
		}
	}
	return out
}

type WeaponBuff struct {
	BaseBuff
}

type Weapon struct {
	Player   *Player
	Data     WeaponData
	Wielded  bool
	Passives *WeaponPassives
	buff     *WeaponBuff
}

func (w Weapon) ToYaml() string {
	out, err := yaml.Marshal(w.Data)
	if err != nil {
		return err.Error()
	}
	return string(out)
}

func (w Weapon) GetType() int {
	if !w.Wielded {
		return WeaponTypeUnarmed
	}
	return GetWeaponTypeValue(w.Data.Type)
}

func (w *Weapon) Wield() {
	w.Wielded = true
	if w.buff == nil {
		w.buff = &WeaponBuff{
			BaseBuff: BaseBuff{
				Name:      "武器.增益",
				Type:      "weapon.buff",
				Steady:    true,
				Permanent: true,
				Modifiers: w.Data.buildModifiers(),
			},
		}
	}
	w.Player.AddBuff(w.buff)
}

func (w *Weapon) Unwield() {
	w.Wielded = false
	if w.buff != nil {
		w.Player.RemoveBuff(w.buff.Name)
	}
}

func (w *Weapon) Rewield() Runable {
	return RunFunc(func(ctx *RunContext) {
		ctx.SetName("装备武器")
		p := ctx.attacker
		p.Weapon.Wield()
		p.AttackCD.Start(3000)
		p.CastTime.Start(3000)
	})
}

func NewWeapon(player *Player, data WeaponData) *Weapon {
	w := &Weapon{
		Player:   player,
		Data:     data,
		Wielded:  false,
		Passives: &WeaponPassives{},
	}
	if data.PassiveData.Name != "" {
		switch data.PassiveData.Name {
		case "鹰刀":
			w.Passives.YinDao = WeaponPassives_YinDao_Builder(player, data.PassiveData.Level)
		default:
			log.Println("不支持武器被动:", data.PassiveData.Name)
		}
	}
	return w
}

// 武器被动
type WeaponPassives struct {
	YinDao *WeaponPassives_YinDao
}

// 武器被动:鹰刀
type WeaponPassives_YinDao struct {
	Player   *Player
	Level    int
	CoolDown int
	Timer    *Timer
}

func (p *WeaponPassives_YinDao) Stun(target *Player) {
	if !p.Timer.IsDone() {
		return
	}
	p.Timer.Start(p.CoolDown)
	target.AddBuff(BuffRepo.Build("鹰刀.昏迷", p.Player, p.Level))
}

func WeaponPassives_YinDao_Builder(player *Player, level int) *WeaponPassives_YinDao {
	cd := 20000
	if level == 3 {
		cd = 15000
	}
	return &WeaponPassives_YinDao{
		Player:   player,
		Level:    level,
		CoolDown: cd,
		Timer:    NewTimer(),
	}
}

// 光环:鹰刀.昏迷
type Buff_YinDao_Faint struct {
	BaseBuff
}

func (b *Buff_YinDao_Faint) OnEnable() {
	b.Owner.State.Faint = true
	b.BaseBuff.OnEnable()
}

func (b *Buff_YinDao_Faint) OnDisable() {
	b.Owner.State.Faint = false
	b.BaseBuff.OnDisable()
}

func Buff_YinDao_Faint_Builder(player *Player, level int, args ...interface{}) Buff {
	irres := false
	if level >= 2 {
		irres = true
	}
	return &Buff_YinDao_Faint{
		BaseBuff{
			Name:         "鹰刀.昏迷",
			Type:         "faint",
			Debuff:       true,
			Creator:      player,
			Irresistible: irres,
			Duration:     3000,
		},
	}
}

func init() {
	BuffRepo.Add("鹰刀.昏迷", Buff_YinDao_Faint_Builder)
}
