package main

import (
	"fmt"

	simc "github.com/abs0293/wsmud-simc/simulator"
	"github.com/spf13/viper"
)

type SimcConfig struct {
	P1    string `mapstructure:"玩家1" yaml:"玩家1,omitempty"`
	P2    string `mapstructure:"玩家2" yaml:"玩家2,omitempty"`
	Times int    `mapstructure:"模拟次数" yaml:"模拟次数,omitempty"`
	Out   string `mapstructure:"结果输出" yaml:"结果输出,omitempty"`
}

var sConf = &SimcConfig{Times: 1, Out: "stdout"}

func main() {
	err := simc.QuickStart("config.yaml")
	if err != nil {
		fmt.Println(err)
	}
}

func readConfig() error {
	cfg := viper.New()
	cfg.AddConfigPath(".")
	cfg.SetConfigFile("config.yaml")
	err := cfg.ReadInConfig()
	if err != nil {
		return err
	}
	return cfg.Unmarshal(sConf)
}

func scenario1() {
	simc.Silent = false
	pA, _ := simc.NewPlayerFromFile("../../conf/刚哥.yaml")
	pB, _ := simc.NewPlayerFromFile("../../conf/鼎酱.yaml")

	lt := &simc.Arena{
		Name:     "1号擂台",
		PlayerA:  pA,
		PlayerB:  pB,
		Duration: 300000,
	}

	lt.Start(pA, pB)
}

func scenario2() {
	simc.Silent = true
	rets := [3]int{}
	times := 10000
	for i := 0; i < times; i++ {
		if i == times-1 {
			simc.Silent = false
		}
		pA, _ := simc.NewPlayerFromFile("../../conf/刚哥.yaml")
		pB, _ := simc.NewPlayerFromFile("../../conf/鼎酱.yaml")
		lt := &simc.Arena{
			Name:     "1号擂台",
			PlayerA:  pA,
			PlayerB:  pB,
			Duration: 300000,
		}
		rets[lt.Start(pA, pB)]++
	}
	fmt.Printf("%d次结果: 1号胜利:%d次, 2号胜利:%d次, 平局:%d次\n", times, rets[0], rets[1], rets[2])
}
