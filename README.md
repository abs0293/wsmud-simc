# 武神传说擂台模拟器

使用说明:

1. 玩家数据

```
cd cmd/playerinfo
go run .
```
根据提示输入账号密码获取玩家数据

2. 设置绝招/装备/武器被动

参考cmd/simc/conf/反震.yaml
```
武学:
  - 部位: 内功
    等级: 7500
    绝招:
      - 名称: 阴阳九转.九烛
      - 名称: 阴阳九转.定乾坤
      - 名称: 阴阳九转.镇天地

装备:
  - 名称: 太极图
    等级: 12


武器:
  特效:
    名称: 鹰刀
    等级: 3
```

3. 配置模拟参数

cmd/simc/config.yaml
```
玩家: 
  - conf/闪避入魔剑.yaml
  - conf/反震.yaml
模拟次数: 5000
```

4. 运行模拟
```
cd cmd/simc
go run .

# 结果
闪避入魔剑vs反震,模拟次数:5000次,用时:23.789258805s,闪避入魔剑获胜:582次,反震获胜:4417次,平局:1次
```