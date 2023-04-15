package simulator

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

var SyncLimit int = 1000

type SimulationConfig struct {
	Players []string `mapstructure:"玩家" yaml:"玩家,omitempty"`
	Times   int      `mapstructure:"模拟次数" yaml:"模拟次数,omitempty"`
	Out     string   `mapstructure:"结果输出" yaml:"结果输出,omitempty"`
}

func QuickStart(fn string) error {
	cfg := SimulationConfig{}

	v, err := readConfig(fn)
	if err != nil {
		return err
	}

	if err := v.Unmarshal(&cfg); err != nil {
		return err
	}

	if cfg.Times <= 0 {
		cfg.Times = 0
	}

	if len(cfg.Players) < 2 {
		return fmt.Errorf("玩家人数不足")
	}

	p1, err := ReadPlayerDataFromFile(cfg.Players[0])
	if err != nil {
		return err
	}

	p2, err := ReadPlayerDataFromFile(cfg.Players[1])
	if err != nil {
		return err
	}

	t := &SimulationTask{
		P1:    p1,
		P2:    p2,
		Times: cfg.Times,
		ret:   [3]int{},
		limit: make(chan struct{}, SyncLimit),
		lock:  &sync.Mutex{},
		wg:    &sync.WaitGroup{},
	}
	if t.Times == 1 {
		Silent = false
	}
	t.Run()
	fmt.Println(t.Result())

	return nil
}

type SimulationTask struct {
	P1, P2   PlayerData
	Times    int
	ret      [3]int
	limit    chan struct{}
	bar      *progressbar.ProgressBar
	duration time.Duration
	wg       *sync.WaitGroup
	lock     *sync.Mutex
}

func (t *SimulationTask) Run() {
	st := time.Now()
	t.bar = progressbar.NewOptions(t.Times, progressbar.OptionSetPredictTime(true))
	for i := 0; i < t.Times; i++ {
		fn := func(i int) {
			lt := &Arena{
				Name:     t.P1.Name + "vs" + t.P2.Name + "_" + fmt.Sprintf("%d", i),
				PlayerA:  NewPlayer(t.P1),
				PlayerB:  NewPlayer(t.P2),
				Roll:     rand.New(rand.NewSource(time.Now().UnixNano())),
				Duration: 60000,
			}
			r := lt.Start(lt.PlayerA, lt.PlayerB)
			if !Silent {
				fmt.Println(strings.Join(lt.GetLogs(), "\n"))
			}
			t.Done(r)
			t.wg.Done()
		}
		t.Go(i, fn)
	}
	t.wg.Wait()
	t.bar.Clear()
	t.duration = time.Since(st)
}

func (t *SimulationTask) Go(i int, fn func(int)) {
	t.limit <- struct{}{}
	t.wg.Add(1)
	t.bar.Add(1)
	go func() {
		fn(i)
		<-t.limit
	}()
}

func (t *SimulationTask) Done(i int) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.ret[i]++
}

func (t *SimulationTask) Result() string {
	p1 := t.P1.Name
	p2 := t.P2.Name
	return fmt.Sprintf(
		"%svs%s,模拟次数:%d次,用时:%v,%s获胜:%d次,%s获胜:%d次,平局:%d次",
		p1, p2,
		t.Times, t.duration,
		p1, t.ret[0],
		p2, t.ret[1],
		t.ret[2],
	)
}
