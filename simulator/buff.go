package simulator

import "github.com/abs0293/wsmud-simc/simulator/log_pb"

type Buff interface {
	GetName() string
	GetType() string
	GetDuration() int
	GetRemaining() int
	SetRemaining(int)
	IsStealable() bool
	IsStackable() bool
	IsExpired() bool
	IsDebuff() bool
	IsIrresistible() bool
	GetStacks() int
	GetModifier(string) float64
	AddStack()
	Update(int)
	OnEnable()
	OnDisable()
	Help() string
	BeStolen(*Player) Buff
	GetOwner() *Player
	SetOwner(*Player)
	GetTarget() *Player
	SetTarget(*Player)
}

type BuffData struct {
	ID        string
	Type      string
	Duration  float64
	Stackable bool
	StackMax  int `yaml:"stack_max"`
	Args      []float64
}

type BaseBuff struct {
	Name          string
	Type          string
	Creator       *Player
	Owner         *Player
	Target        *Player
	Permanent     bool
	Debuff        bool
	Steady        bool // false:可被偷取
	Irresistible  bool
	Duration      int
	RemainingTime int
	Stackable     bool
	StackMax      int
	Stacks        int
	Modifiers     []Modifier
}

func (b *BaseBuff) GetName() string {
	return b.Name
}

func (b *BaseBuff) GetType() string {
	return b.Type
}

func (b *BaseBuff) GetDuration() int {
	return b.Duration
}

func (b *BaseBuff) GetRemaining() int {
	return b.RemainingTime
}

func (b *BaseBuff) SetRemaining(t int) {
	b.RemainingTime = t
}

func (b *BaseBuff) IsStealable() bool {
	return !b.Debuff && !b.Steady
}

func (b *BaseBuff) IsStackable() bool {
	return b.Stackable
}

func (b *BaseBuff) IsIrresistible() bool {
	return b.Irresistible
}

func (b *BaseBuff) IsExpired() bool {
	if b.Permanent {
		return false
	}
	return b.GetRemaining() == 0
}

func (b *BaseBuff) IsDebuff() bool {
	return b.Debuff
}

func (b *BaseBuff) GetStacks() int {
	return b.Stacks
}

func (b *BaseBuff) AddStack() {
	if b.Stacks < 1 {
		return
	}
	if b.Stacks < b.StackMax {
		b.Stacks++
	}
	b.RemainingTime = b.Duration
	b.EventRefreshLog()
}

func (b *BaseBuff) GetModifier(t string) float64 {
	v := 0.0
	for _, m := range b.Modifiers {
		if t == m.Name {
			v += m.Value
		}
	}
	return v
}

func (b *BaseBuff) Update(diff int) {
	b.elasped(diff)
}

func (b *BaseBuff) elasped(diff int) {
	for diff > 0 {
		d := b.RemainingTime - diff
		if d >= 0 {
			b.RemainingTime = d
			return
		}
		if b.Stacks <= 1 {
			b.RemainingTime = 0
			return
		}
		diff -= b.RemainingTime
		b.Stacks--
		b.RemainingTime = b.Duration
	}
}

func (b *BaseBuff) OnEnable() {
	b.EventAddLog()
}
func (b *BaseBuff) OnDisable() {
	b.EventRemoveLog()
}

// TODO:偷取机制未明
func (b *BaseBuff) BeStolen(who *Player) Buff { b.Owner = who; b.Stacks = 1; return b }

func (b BaseBuff) GetCreator() *Player {
	return b.Creator
}

func (b BaseBuff) GetOwner() *Player {
	return b.Owner
}

func (b *BaseBuff) SetOwner(player *Player) {
	b.Owner = player
}

func (b BaseBuff) GetTarget() *Player {
	return b.Target
}

func (b *BaseBuff) SetTarget(target *Player) {
	b.Target = target
}

func (b *BaseBuff) Help() string {
	return ""
}

func (b *BaseBuff) Log(log *log_pb.Log) {
	log.Timestamp = int32(b.Owner.Arena.Ticks)
	log.Player = b.Owner.Name
	log.Aura.Name = b.Name
	log.Aura.Remaining = int32(b.RemainingTime)
	log.Aura.Stacks = int32(b.Stacks)
	log.Aura.Type = b.Type
	for _, mod := range b.Modifiers {
		log.Aura.Modifier = append(log.Aura.Modifier, &log_pb.Modifier{Name: mod.Name, Value: mod.Value})
	}
	b.Owner.Log(log)
}

func (b *BaseBuff) EventAddLog() {
	l := &log_pb.Log{
		Aura: &log_pb.Log_Aura{
			Event: log_pb.Log_Aura_Event_Add,
		},
	}
	b.Log(l)
}

func (b *BaseBuff) EventRemoveLog() {
	l := &log_pb.Log{
		Aura: &log_pb.Log_Aura{
			Event: log_pb.Log_Aura_Event_Remove,
		},
	}
	b.Log(l)
}

func (b *BaseBuff) EventRefreshLog() {
	l := &log_pb.Log{
		Aura: &log_pb.Log_Aura{
			Event: log_pb.Log_Aura_Event_Refresh,
		},
	}
	if b.Stackable {
		l.Aura.Stacks = int32(b.Stacks)
	}
	if !b.Permanent {
		l.Aura.Remaining = int32(b.RemainingTime)
	}
	b.Log(l)
}

type Modifier struct {
	Name  string
	Value float64
}

type Modifiers []Modifier

func (mods Modifiers) GetModifier(t string) float64 {
	v := 0.0
	for _, m := range mods {
		v += m.Value
	}
	return v
}

func GetModifier(target string, mods_arr ...[]Modifier) float64 {
	v := 0.0
	for _, mods := range mods_arr {
		for _, m := range mods {
			if m.Name == target {
				v += m.Value
			}
		}
	}
	return v
}

// Buff管理器
type Buffs struct {
	Player *Player
	buffs  []Buff
}

func (mgr *Buffs) Update(diff int) {
	next := []Buff{}
	for _, b := range mgr.buffs {
		b.Update(diff)
		if !b.IsExpired() {
			next = append(next, b)
		} else {
			b.OnDisable()
			//mgr.Player.Printf("光环:失去%s(%s)\n", b.GetName(), b.GetType())
		}
	}
	mgr.buffs = next
}

func (mgr *Buffs) Remove(id string) {
	i := mgr.IndexOf(id)
	if i >= 0 {
		mgr.buffs = append(mgr.buffs[:i], mgr.buffs[i+1:]...)
	}
}

func (mgr *Buffs) RemoveByType(typ string) {
	tmp := []Buff{}
	for _, buff := range mgr.buffs {
		if buff.GetType() == typ {
			buff.OnDisable()
			continue
		}
		tmp = append(tmp, buff)
	}
	mgr.buffs = tmp
}

func (mgr *Buffs) Add(buf Buff) {
	if buf == nil {
		return
	}
	if buf.IsDebuff() && !buf.IsIrresistible() {
		buf.SetRemaining(CalcDebuffDurationP(buf.GetDuration(), mgr.Player))
	} else {
		buf.SetRemaining(buf.GetDuration())
	}
	// 负面抵抗减少debuff的持续时间,持续时间为0时debuff完全不生效,直接返回。
	if buf.IsExpired() {
		return
	}
	for _, b := range mgr.buffs {
		if b.GetName() == buf.GetName() {
			if b.IsStackable() {
				b.AddStack()
				//mgr.Player.Printf("光环:刷新%s(%s)[%d]\n", b.GetName(), b.GetType(), b.GetStacks())
			} else {
				return
			}
		}
	}
	buf.SetOwner(mgr.Player)
	buf.OnEnable()
	mgr.buffs = append(mgr.buffs, buf)
	//mgr.Player.Printf("光环:获得%s(%s),持续时间:%.3f秒\n", buf.GetName(), buf.GetType(), Ms2Sec(buf.GetRemaining()))
}

func (mgr Buffs) IndexOf(name string) int {
	for i, b := range mgr.buffs {
		if b.GetName() == name {
			return i
		}
	}
	return -1
}

func (mgr Buffs) GetBuff(name string) Buff {
	for _, b := range mgr.buffs {
		if b.GetName() == name {
			return b
		}
	}
	return nil
}

func (mgr *Buffs) BeStolen(who *Player, sid string) Buff {
	tmp := []Buff{}
	for _, b := range mgr.buffs {
		if b.IsStealable() {
			tmp = append(tmp, b)
		}
	}
	if s := len(tmp); s > 0 {
		b := tmp[mgr.Player.Roll.Intn(s)]
		b.OnDisable()
		mgr.Remove(b.GetName())
		return b.BeStolen(who)
	}
	return nil
}

func (mgr *Buffs) ClearDebuffs() {
	tmp := []Buff{}
	for _, b := range mgr.buffs {
		if b.IsDebuff() {
			b.OnDisable()
		} else {
			tmp = append(tmp, b)
		}
	}
	mgr.buffs = tmp
}

func (mgr Buffs) GetModifier(n string) float64 {
	m := 0.0
	for _, b := range mgr.buffs {
		m += b.GetModifier(n)
	}
	return m
}

func (mgr Buffs) GetModifierExclude(n string, exclude ...string) float64 {
	m := 0.0
	for _, b := range mgr.buffs {
		if IsMember(b.GetName(), exclude) {
			continue
		}
		m += b.GetModifier(n)
	}
	return m
}

func NewBuffs(player *Player) *Buffs {
	return &Buffs{Player: player}
}

func IsMember(x string, arr []string) bool {
	for _, m := range arr {
		if x == m {
			return true
		}
	}
	return false
}

var BuffRepo = &BuffRepository{make(map[string]BuffBuilder)}

type BuffBuilder func(*Player, int, ...interface{}) Buff

type BuffRepository struct {
	builders map[string]BuffBuilder
}

func (repo *BuffRepository) Add(name string, builder BuffBuilder) {
	repo.builders[name] = builder
}

func (repo *BuffRepository) Build(name string, player *Player, level int, args ...interface{}) Buff {
	b, ok := repo.builders[name]
	if ok {
		return b(player, level, args...)
	}
	return nil
}
