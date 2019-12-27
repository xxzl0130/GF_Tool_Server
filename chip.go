package main

import (
	"strings"
	"encoding/json"
	"sort"
	"time"
	"fmt"
	"encoding/base64"

	//"github.com/xxzl0130/rsyars/pkg/util"
	"github.com/xxzl0130/rsyars/rsyars.adapter/hycdes"
	"github.com/xxzl0130/rsyars/rsyars.x/soc"

	cipher "github.com/xxzl0130/GF_Tool_Server/GF_cipher"
)

// 存储详细游戏数据的结构体，外层以UID为map的key
type UserInfo struct{
	rule string			// 芯片规则
	ruleChanged bool	// 规则是否改变的标志
	name string			// 游戏昵称
	uid  string			// 游戏UID
	chipCode string		// 芯片代码
	chipJson string		// 芯片json
	kalinaLevel string	// 格林娜好感等级
	kalinnaFavor string	// 格林娜好感度
	spendPoint string	// 消费
	allData string		// 完整的数据
	time int64			// 时间戳
}

type User_Info struct{
	User_id		string `json:"user_id"`
	Name		string `json:"name"`
}
type UserRecord struct {
	Spend_point string `json:"spend_point"`
}
type KalinaData struct {
	Level string `json:"level"`
	Favor string `json:"favor"`
}
type GF_Json struct {
	Info User_Info `json:"user_info"`
	Record UserRecord `json:"user_record"`
	Kalina KalinaData `json:"kalina_with_user_info"`
	SoC map[string]*soc.SoC `json:"chip_with_user_info"`
}
type GF_Simple_Json struct{
	Info User_Info `json:"user_info"`
	Record UserRecord `json:"user_record"`
	Kalina KalinaData `json:"kalina_with_user_info"`
}
type GF_Chip_Json struct{
	SoC map[string]*soc.SoC `json:"chip_with_user_info"`
}

func GzipCompress(data []byte) string{
	gzip, err := cipher.GzipCompress([]byte(data))
	if err != nil {
		return "数据处理错误！"
	}
	b64 := base64.StdEncoding
	b64str := "#" + b64.EncodeToString(gzip)
	return b64str
}

func (tool *Tool) buildChips(uid string) string{
	tool.infoMutex.RLock()
	info, isPresent := tool.userinfo[uid]
	tool.infoMutex.RUnlock()
	if !isPresent{
		return "数据不存在！"
	}
	if len(info.chipCode) > 5 && !info.ruleChanged{
		return info.chipCode
	}
	if time.Now().Unix() - info.time < 5{
		return "请勿频繁请求！"
	}
	tool.infoMutex.Lock()
	info.time = time.Now().Unix()
	tool.userinfo[uid] = info // 更新时间避免被清理
	tool.infoMutex.Unlock()

	girls := GF_Json{}
	if err := json.Unmarshal([]byte(info.allData), &girls); err != nil {
		return "数据解析错误！"
	}

	var values []*soc.SoC
	for _, c := range girls.SoC {
		values = append(values, c)
	}

	var targets []*hycdes.SoC
	for _, value := range values {
		if !tool.checkSoc(value, info.rule) {
			continue
		}
		target, err := hycdes.NewSoC(value)
		if err != nil {
			if !strings.Contains(err.Error(), "unknown") {
				return "数据解析错误！"
			} else {
				continue
			}
		}
		targets = append(targets, target)
	}

	if len(targets) == 0 {
		return "数据解析错误！"
	}

	sort.SliceStable(targets, func(i, j int) bool {
		return targets[i].ID < targets[j].ID
	})

	c, err := hycdes.Build(targets)
	if err != nil {
		//rs.log.Errorf("生成芯片代码失败 -> %+v", err)
		return "数据解析错误！"
	}
	info.chipCode = GzipCompress([]byte(c))
	tool.infoMutex.Lock()
	tool.userinfo[uid] = info
	tool.infoMutex.Unlock()

	return info.chipCode
}

func (tool *Tool) buildChipJson(uid string) string{
	tool.infoMutex.RLock()
	info, isPresent := tool.userinfo[uid]
	tool.infoMutex.RUnlock()
	if !isPresent{
		return "数据不存在！"
	}
	if len(info.chipJson) > 5{
		return info.chipJson
	}
	if time.Now().Unix() - info.time < 5{
		return "请勿频繁请求！"
	}
	tool.infoMutex.Lock()
	info.time = time.Now().Unix()
	tool.userinfo[uid] = info // 更新时间避免被清理
	tool.infoMutex.Unlock()

	girls := GF_Json{}
	if err := json.Unmarshal([]byte(info.allData), &girls); err != nil {
		return "数据解析错误！"
	}

	tmp := GF_Chip_Json{
		SoC : girls.SoC,
	}

	data, _ := json.Marshal(tmp)

	info.chipJson = GzipCompress(data) //保存json
	tool.infoMutex.Lock()
	tool.userinfo[uid] = info
	tool.infoMutex.Unlock()

	return info.chipJson;
}

func (tool *Tool) buildKalina(uid string) [2]string{
	var str[2] string
	str[0]="数据解析错误！"
	tool.infoMutex.RLock()
	info, isPresent := tool.userinfo[uid]
	tool.infoMutex.RUnlock()
	if !isPresent{
		return str
	}
	if time.Now().Unix() - info.time < 5{
		return str
	}
	tool.infoMutex.Lock()
	info.time = time.Now().Unix()
	tool.userinfo[uid] = info // 更新时间避免被清理
	tool.infoMutex.Unlock()

	value := GF_Simple_Json{}
	if err := json.Unmarshal([]byte([]byte(info.allData)), &value); err != nil {
		//fmt.Printf("解析JSON数据失败 -> %+v\n", err)
		return str
	}
	str[0] = fmt.Sprintf("%v 元",value.Record.Spend_point)
	str[1] = fmt.Sprintf("%v  Lv.%v/30", value.Kalina.Favor, value.Kalina.Level)
	return str
}

func (tool *Tool) checkSoc(value *soc.SoC, rule string) bool {
	if len(rule) < 2{
		rule = "02"
	}
	switch rule[0] {
	case '0':
		switch rule[1] {
		case '0':
			return true
		case '1':
			return value.SquadWithUserID != "0"
		case '2':
			return value.SquadWithUserID == "0"
		}
	case '1':
		switch rule[1] {
		case '0':
			return value.Locked != "0"
		case '1':
			return value.Locked != "0" && value.SquadWithUserID != "0"
		case '2':
			return value.Locked != "0" && value.SquadWithUserID == "0"
		}
	case '2':
		switch rule[1] {
		case '0':
			return value.Locked == "0"
		case '1':
			return value.Locked == "0" && value.SquadWithUserID != "0"
		case '2':
			return value.Locked == "0" && value.SquadWithUserID == "0"
		}
	}
	return false
}
