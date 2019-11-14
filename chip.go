package main

import (
	"strings"
	"encoding/json"
	"sort"
	"time"

	//"github.com/xxzl0130/rsyars/pkg/util"
	"github.com/xxzl0130/rsyars/rsyars.adapter/hycdes"
	"github.com/xxzl0130/rsyars/rsyars.x/soc"
)

// 存储详细游戏数据的结构体，外层以UID为map的key
type UserInfo struct{
	rule string			// 芯片规则
	name string			// 游戏昵称
	uid  string			// 游戏UID
	chipCode string		// 芯片代码
	kalinaLevel string	// 格林娜好感等级
	kalinnaFavor string	// 格林娜好感度
	spendPoint string	// 消费
	allData string		// 完整的数据
	time int64			// 时间戳
}

func (tool *Tool) buildChips(uid string) bool{
	type Girls struct {
		SoC map[string]*soc.SoC `json:"chip_with_user_info"`
	}

	tool.infoMutex.RLock()
	info, isPresent := tool.userinfo[uid]
	tool.infoMutex.RUnlock()
	if !isPresent{
		return false
	}
	if len(info.chipCode) > 5{
		return true
	}
	tool.infoMutex.Lock()
	info.time = time.Now().Unix()
	tool.userinfo[uid] = info // 更新时间避免被清理
	tool.infoMutex.Unlock()

	girls := Girls{}
	if err := json.Unmarshal([]byte(info.allData), &girls); err != nil {
		return false
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
				return false
			} else {
				continue
			}
		}
		targets = append(targets, target)
	}

	if len(targets) == 0 {
		return false
	}

	sort.SliceStable(targets, func(i, j int) bool {
		return targets[i].ID < targets[j].ID
	})

	c, err := hycdes.Build(targets)
	if err != nil {
		//rs.log.Errorf("生成芯片代码失败 -> %+v", err)
		return false
	}
	info.chipCode = c
	tool.infoMutex.Lock()
	tool.userinfo[uid] = info
	tool.infoMutex.Unlock()

	return true
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
