package simulator

import (
	"math/rand"

	"github.com/abs0293/wsmud-simc/simulator/log_pb"
)

var Silent = true

type Arena struct {
	Name    string
	PlayerA *Player
	PlayerB *Player
	Roll    *rand.Rand

	logSN int
	logs  []*log_pb.Log

	Ticks    int
	Duration int
}

func (a *Arena) Start(pa, pb *Player) int {
	pa.Target = pb
	pb.Target = pa
	pa.Arena = a
	pb.Arena = a

	a.logSN = 0

	pa.OnCombatStart()
	pb.OnCombatStart()
	for {
		if a.GameOver() {
			if a.PlayerA.IsAlive() {

				return 0
			} else if a.PlayerB.IsAlive() {
				return 1
			} else {
				return 2
			}
		}

		if a.Roll.Float64() < 0.5 {
			pa.Update(TICK)
			pb.Update(TICK)

			pa.Action()
			pb.Action()
		} else {
			pb.Update(TICK)
			pa.Update(TICK)

			pb.Action()
			pa.Action()
		}

		a.Ticks++
	}
}

func (a Arena) TimeOver() bool {
	return a.Ticks*TICK >= a.Duration
}

func (a Arena) Timestamp() float64 {
	return Ms2Sec(a.Ticks * TICK)
}

func (a Arena) GameOver() bool {
	return a.TimeOver() || !a.PlayerA.IsAlive() || !a.PlayerB.IsAlive()
}

func (a *Arena) Log(logs ...*log_pb.Log) {
	for _, l := range logs {
		l.SerialNumber = int32(a.logSN)
		a.logs = append(a.logs, l)
		a.logSN++
	}
}

func (a *Arena) GetLogs() []string {
	out := []string{}
	for _, l := range a.logs {
		out = append(out, l.String())
	}
	return out
}
