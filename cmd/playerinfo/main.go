package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/abiosoft/ishell/v2"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

var (
	httpClient = http.DefaultClient
	wsConn     *websocket.Conn
	shell      = ishell.New()
	token      = ""
	roleID     string
	roleName   string

	info = viper.New()
)

type Servers struct {
	Servers []ServerInfo
}

func (svrs Servers) List() []string {
	list := []string{}
	for _, svr := range svrs.Servers {
		list = append(list, svr.Name)
	}
	return list
}

type ServerInfo struct {
	ID     int
	Name   string
	Port   int
	IP     string
	IsTest bool
	IsRcd  bool
}

type RoleMessage struct {
	Type  string
	Roles []RoleInfo
}

func (msg RoleMessage) List() []string {
	list := []string{}
	for _, r := range msg.Roles {
		list = append(list, r.Name)
	}
	return list
}

func (msg *RoleMessage) Marshal(in []byte) error {
	err := json.Unmarshal(preMarshal(in), msg)
	return err
}

type RoleInfo struct {
	Name  string
	Title string
	ID    string
}

func (i ServerInfo) WsUrl() string {
	return fmt.Sprintf("ws://%s:%d", i.IP, i.Port)
}

type SkillMessage struct {
	Type   string
	Dialog string
	Items  []SkillInfo
}

func (msg *SkillMessage) Marshal(in []byte) error {
	err := json.Unmarshal(preMarshal(in), msg)
	return err
}

type SkillInfo struct {
	ID     string
	Enable string `json:"enable_skill"`
}

var auraSkill = []string{
	"taijijian", "taijijian2",
	"yijinjing", "yijinjing2",
	"dugujiujian", "dugujiujian2",
	"yitianjianfa", "yitianjianfa2",
	"lingboweibu", "lingboweibu2",
	"xianglongzhang", "xianglongzhang2",
	"mantianhuayu", "mantianhuayu2",
}

func (i SkillInfo) HasAura() bool {
	for _, as := range auraSkill {
		if as == i.ID {
			return true
		}
	}
	return false
}

type SkillDescMessage struct {
	Type   string
	Dialog string
	Desc   string
}

func (msg *SkillDescMessage) Marshal(in []byte) error {
	err := json.Unmarshal(preMarshal(in), msg)
	return err
}

type ItemDescMessage struct {
	Type   string
	Dialog string
	Desc   string
}

func (msg *ItemDescMessage) Marshal(in []byte) error {
	err := json.Unmarshal(preMarshal(in), msg)
	return err
}

type ScoreMessage struct {
	Type   string
	Dialog string
	MP     int `json:"max_mp"`
}

func (msg *ScoreMessage) Marshal(in []byte) error {
	err := json.Unmarshal(preMarshal(in), msg)
	return err
}

func main() {
	userLogin()
	svrs, err := getServerInfo()
	if err != nil {
		log.Panic(err)
	}

	idx := shell.MultiChoice(svrs.List(), "选择大区")
	svr := svrs.Servers[idx]

	connect(svr.WsUrl())
	selectRole()
	getSkillInfo()
	getEquipInfo()
	getMaxMP()

	info.Set("姓名", roleName)
	info.AddConfigPath(".")
	info.SetConfigFile(roleName + ".yaml")
	_, err = os.Create(roleName + ".yaml")
	if err != nil {
		log.Panic(err)
	}
	err = info.WriteConfig()
	if err != nil {
		log.Panic(err)
	}
}

func getMaxMP() {
	wsConn.WriteMessage(1, []byte("score"))
	Done := false
	for !Done {
		_, msg, err := wsConn.ReadMessage()
		if err != nil {
			log.Panic(err)
		}
		data := string(msg)
		if strings.Contains(data, `"dialog":"score"`) {
			desc := &ScoreMessage{}
			err := desc.Marshal(msg)
			if err != nil {
				log.Panic(err)
			}
			info.Set("内力d", desc.MP)
			Done = true
		}
	}
}

func getEquipInfo() {
	var (
		data string
		reg  = regexp.MustCompile(`<blk>[^<]+</blk>`)
		reg2 = regexp.MustCompile(`(\p{Han}+)：[\+|\-]([0-9\.%]+)`)
	)

	wsConn.WriteMessage(1, []byte("setting hide_equip 0"))

	for i := 0; i <= 10; i++ {
		fmt.Printf("装备(%d):\n", i)
		cmd := fmt.Sprintf("look %d of %s", i, roleID)
		wsConn.WriteMessage(1, []byte(cmd))
		time.Sleep(time.Second)

		Done := false
		for !Done {
			_, msg, err := wsConn.ReadMessage()
			if err != nil {
				log.Panic(err)
			}
			data = string(msg)
			if strings.HasPrefix(data, `{"type":"item"`) {
				desc := &ItemDescMessage{}
				err := desc.Marshal(msg)
				if err != nil {
					log.Panic(err)
				}
				d := reg.ReplaceAllString(desc.Desc, "")
				for _, s := range reg2.FindAllStringSubmatch(d, -1) {
					if i == 0 {
						s[1] = "武器." + s[1]
					}
					writeKv(s)
				}
				Done = true
			}
		}
	}

	wsConn.WriteMessage(1, []byte("setting hide_equip 1"))
}

func getSkillInfo() {
	var (
		data      string
		checkList = map[string]bool{
			"force":    true,
			"unarmed":  true,
			"sword":    true,
			"blade":    true,
			"staff":    true,
			"club":     true,
			"whip":     true,
			"throwing": true,
			"dodge":    true,
			"parry":    true,
		}

		reg  = regexp.MustCompile(`<blk>[^<]+</blk>`)
		reg2 = regexp.MustCompile(`(\p{Han}+)：[\+|\-]([0-9\.%]+)`)
		reg3 = regexp.MustCompile(`将你内力的(\d+)%转化为气血`)
	)
	wsConn.WriteMessage(1, []byte("skills"))
	time.Sleep(time.Second)
	for {
		_, msg, err := wsConn.ReadMessage()
		if err != nil {
			log.Panic(err)
		}
		data = string(msg)
		if strings.Contains(data, `enable_skill`) {
			skills := &SkillMessage{}
			err := skills.Marshal(msg)
			if err != nil {
				log.Panic(err)
			}
			for _, i := range skills.Items {
				if i.Enable != "" {
					checkList[i.Enable] = true
				}
				if i.HasAura() {
					checkList[i.ID] = true
				}
			}
			break
		}
	}

	for skillName := range checkList {
		fmt.Println("武功:", skillName)
		Done := false
		wsConn.WriteMessage(1, []byte("checkskill "+skillName))
		time.Sleep(time.Second)
		for !Done {
			_, msg, err := wsConn.ReadMessage()
			if err != nil {
				log.Panic(err)
			}
			data = string(msg)
			if strings.HasPrefix(data, `{"type":"dialog","dialog":"skills","desc":`) {
				desc := &SkillDescMessage{}
				err := desc.Marshal(msg)
				if err != nil {
					log.Panic(err)
				}
				d := reg.ReplaceAllString(desc.Desc, "")
				for _, s := range reg2.FindAllStringSubmatch(d, -1) {
					writeKv(s)
				}
				for _, s := range reg3.FindAllStringSubmatch(d, -1) {
					fmt.Println(s[0])
					f, err := strconv.ParseFloat(s[1], 64)
					if err != nil {
						log.Panicln(err)
					}
					f = math.Trunc(f*1e5+0.5) / 1e5
					o := info.GetFloat64("内力转化%")
					info.Set("内力转化%", o+f/100)
				}
				Done = true
			}
		}
	}
}

func userLogin() {
	data := url.Values{}
	shell.Printf("账号:")
	data.Add("code", shell.ReadLine())
	shell.Printf("密码:")
	data.Add("pwd", shell.ReadLine())

	resp, err := http.PostForm("http://wamud.com/UserAPI/Login", data)
	if err != nil {
		log.Panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Panic("status code:", resp.StatusCode)
	}

	t := []string{}
	for _, c := range resp.Cookies() {
		t = append(t, c.Value)
	}
	token = strings.Join(t, " ")
}

func getServerInfo() (*Servers, error) {
	resp, err := httpClient.Get("http://wamud.com/Game/GetServer")
	if err != nil {
		return nil, err
	}
	body, err := respBody(resp)
	if err != nil {
		return nil, err
	}

	svrs := &Servers{}
	err = json.Unmarshal(body, &svrs.Servers)
	if err != nil {
		return nil, err
	}
	return svrs, nil
}

func connect(ws string) {
	var (
		err error
	)
	wsConn, _, err = websocket.DefaultDialer.Dial(ws, nil)
	if err != nil {
		log.Panicln(err)
	}
}

func selectRole() {
	err := wsConn.WriteMessage(websocket.TextMessage, []byte(token))
	if err != nil {
		log.Panic(err)
	}
	_, msg, err := wsConn.ReadMessage()
	if err != nil {
		log.Panic(err)
	}

	roles := &RoleMessage{}
	err = roles.Marshal(msg)
	if err != nil {
		log.Panic(err)
	}

	rIdx := shell.MultiChoice(roles.List(), "选择角色")
	roleID = roles.Roles[rIdx].ID
	roleName = roles.Roles[rIdx].Name

	wsConn.WriteMessage(websocket.TextMessage, []byte("login "+roleID))
	time.Sleep(time.Second)
}

func respBody(resp *http.Response) ([]byte, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if bytes.Equal(body[:3], []byte{239, 187, 191}) {
		body = body[3:]
	}
	return body, nil
}

func preMarshal(in []byte) []byte {
	var (
		reg = regexp.MustCompile(`([a-zA-Z]\w+):`)
	)

	data := reg.ReplaceAllString(string(in), `"$1":`)
	data = strings.ReplaceAll(data, "'", "\"")
	return []byte(data)
}

func parseKey(in string) string {
	var (
		p string
		v string
	)

	if strings.HasPrefix(in, "武器.") {
		p = "武器."
		v = strings.TrimPrefix(in, p)
	} else {
		v = in
	}

	switch v {
	case "忽视对方防御":
		v = "忽防"
	case "最终伤害":
		v = "终伤"
	case "伤害减免":
		v = "免伤"
	case "受到的伤害减少":
		v = "免伤"
	case "负面状态抵抗":
		v = "负面抵抗"
	case "绝招冷却时间":
		v = "绝招冷却"
	case "绝招释放时间":
		v = "绝招释放"
	case "攻击速度":
		v = "攻速"
	}
	return p + v
}

func writeKv(in []string) {
	var (
		base = 1.0
	)
	fmt.Println(in[1], in[2])
	key := parseKey(in[1])
	val := strings.TrimLeftFunc(in[2], func(r rune) bool { return r == '+' || r == '-' })

	if strings.HasSuffix(val, "%") {
		key += "%"
		base = 100
		val = strings.TrimSuffix(val, "%")
	} else {
		key += "d"
	}
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		log.Panic(val, err)
	}
	f = f/base + info.GetFloat64(key)
	f = math.Trunc(f*1e5+0.5) / 1e5
	info.Set(key, f)
}
