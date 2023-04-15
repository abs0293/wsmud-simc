package simulator

import "log"

type EquipmentData struct {
	Name  string `mapstructure:"名称" yaml:"名称,omitempty"`
	Level int    `mapstructure:"等级" yaml:"等级,omitempty"`
}

type Equipment interface {
	GetName() string
	Update(int)
	CanRun(*Player, ...interface{}) bool
	Run(*RunContext)
}

type BaseEquipment struct {
	Player      *Player
	Name        string
	Level       int
	CoolDown    int
	Passive     bool
	IgnoreFaint bool
	IgnoreBusy  bool
	Timer       *Timer
}

func (equip *BaseEquipment) GetName() string {
	return equip.Name
}

func (equip *BaseEquipment) Update(diff int) {
	equip.Timer.Update(diff)
}

func (equip *BaseEquipment) CanRun(target *Player, args ...interface{}) bool {
	if equip.Passive {
		return false
	}
	if equip.Player == nil || !equip.Player.IsAlive() {
		return false
	}
	if !equip.Timer.IsDone() {
		return false
	}
	if equip.Player.State.Faint && !equip.IgnoreFaint {
		return false
	}
	if equip.Player.State.Busy && !equip.IgnoreBusy {
		return false
	}
	if target != nil && !target.IsAlive() {
		return false
	}
	return true
}

func (equip *BaseEquipment) Run(ctx *RunContext) {
	ctx.SetName(equip.Name)
	equip.Timer.Start(equip.CoolDown)
}

// 装备管理器
type Equipments struct {
	Equipments []Equipment
}

func (mgr *Equipments) Update(diff int) {
	for _, e := range mgr.Equipments {
		e.Update(diff)
	}
}

func (mgr *Equipments) CanRun(target *Player, args ...interface{}) []Runable {
	out := []Runable{}
	for _, e := range mgr.Equipments {
		if e.CanRun(target, args...) {
			out = append(out, e)
		}
	}
	return out
}

func NewEquipments(owner *Player, datas ...EquipmentData) *Equipments {
	equips := &Equipments{}
	for _, data := range datas {
		e := EquipmentRepo.Build(data.Name, owner, data.Level)
		if e != nil {
			equips.Equipments = append(equips.Equipments, e)
		} else {
			log.Println("不支持装备:", data.Name)
		}
	}
	return equips
}

// 装备:太极图
type Equipment_TaiJiTu struct {
	BaseEquipment
}

func (equip *Equipment_TaiJiTu) Run(ctx *RunContext) {
	equip.BaseEquipment.Run(ctx)
	equip.Player.AddBuff(BuffRepo.Build("太极图.太极", equip.Player, equip.Level))
}

func Equipment_TaiJiTu_Builder(player *Player, level int) Equipment {
	return &Equipment_TaiJiTu{
		BaseEquipment: BaseEquipment{
			Name:     "太极图",
			Level:    level,
			Player:   player,
			CoolDown: 60000,
			Timer:    NewTimer(),
		},
	}
}

// 光环:太极图.太极
type Buff_TaiJiTu_Taiji struct {
	BaseBuff
}

func Buff_TaiJiTu_Taiji_Builder(player *Player, level int, args ...interface{}) Buff {
	return &Buff_TaiJiTu_Taiji{
		BaseBuff{
			Name:     "太极图.太极",
			Type:     "taiji",
			Creator:  player,
			Duration: 8000,
		},
	}
}

// 装备:太阴
type Equipment_TaiYin struct {
	BaseEquipment
}

func (equip *Equipment_TaiYin) Run(ctx *RunContext) {
	equip.BaseEquipment.Run(ctx)
	ctx.target.AddBuff(BuffRepo.Build("太阴.太阴", equip.Player, equip.Level))
}

func Equipment_TaiYin_Builder(player *Player, level int) Equipment {
	return &Equipment_TaiYin{
		BaseEquipment: BaseEquipment{
			Name:     "太阴",
			Level:    level,
			Player:   player,
			CoolDown: 60000,
			Timer:    NewTimer(),
		},
	}
}

// 光环:太阴.太阴
type Buff_TaiYin_TaiYin struct {
	BaseBuff
}

func Buff_TaiYin_TaiYin_Builder(player *Player, level int, args ...interface{}) Buff {
	mod := 0.2 + 0.01*float64(level)
	return &Buff_TaiYin_TaiYin{
		BaseBuff{
			Name:         "太阴.太阴",
			Type:         "taiyin",
			Creator:      player,
			Irresistible: true,
			Debuff:       true,
			Duration:     7000 + 200*level,
			Modifiers: []Modifier{
				{"攻速%", -mod},
				{"绝招冷却%", -mod},
			},
		},
	}
}

// 装备:缚神索
type Equipment_FuShenSuo struct {
	BaseEquipment
}

func (equip *Equipment_FuShenSuo) Run(ctx *RunContext) {
	equip.BaseEquipment.Run(ctx)
	ctx.target.AddBuff(BuffRepo.Build("缚神索.束缚", equip.Player, equip.Level))
}

func Equipment_FuShenSuo_Builder(player *Player, level int) Equipment {
	return &Equipment_FuShenSuo{
		BaseEquipment: BaseEquipment{
			Name:     "缚神索",
			Level:    level,
			Player:   player,
			CoolDown: 60000,
			Timer:    NewTimer(),
		},
	}
}

// 光环:缚神索.束缚
type Buff_FuShenSuo_ShuFu struct {
	BaseBuff
}

func (buff *Buff_FuShenSuo_ShuFu) OnEnable() {
	buff.Owner.State.Busy = true
	buff.BaseBuff.OnEnable()
}

func (buff *Buff_FuShenSuo_ShuFu) OnDisable() {
	buff.Owner.State.Busy = false
	buff.BaseBuff.OnDisable()
}

func Buff_FuShenSuo_ShuFu_Builder(player *Player, level int, args ...interface{}) Buff {
	return &Buff_FuShenSuo_ShuFu{
		BaseBuff{
			Name:      "缚神索.束缚",
			Type:      "busy",
			Creator:   player,
			Debuff:    true,
			Duration:  6000 + 1000*level,
			Modifiers: []Modifier{},
		},
	}
}

// equipment repo
type EquipmentBuilder func(*Player, int) Equipment

var EquipmentRepo = &EquipmentRepository{make(map[string]EquipmentBuilder)}

type EquipmentRepository struct {
	builders map[string]EquipmentBuilder
}

func (repo *EquipmentRepository) Add(name string, builder EquipmentBuilder) {
	repo.builders[name] = builder
}

func (repo *EquipmentRepository) Build(name string, player *Player, level int) Equipment {
	b, ok := repo.builders[name]
	if ok {
		return b(player, level)
	}
	return nil
}

func init() {
	EquipmentRepo.Add("太极图", Equipment_TaiJiTu_Builder)
	EquipmentRepo.Add("太阴", Equipment_TaiYin_Builder)
	EquipmentRepo.Add("缚神索", Equipment_FuShenSuo_Builder)

	BuffRepo.Add("太极图.太极", Buff_TaiJiTu_Taiji_Builder)
	BuffRepo.Add("太阴.太阴", Buff_TaiYin_TaiYin_Builder)
	BuffRepo.Add("缚神索.束缚", Buff_FuShenSuo_ShuFu_Builder)
}
