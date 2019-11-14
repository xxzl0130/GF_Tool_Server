package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	//"github.com/xxzl0130/rsyars/pkg/util"
	//"github.com/xxzl0130/rsyars/rsyars.adapter/hycdes"
	//"github.com/xxzl0130/rsyars/rsyars.x/soc"
)

// 存储详细游戏数据的结构体，外层以UID为map的key
type UserInfo struct{
	rule string // 芯片规则
	name string // 游戏昵称
	uid  string // 游戏UID
	code string // 芯片代码
	time int64  // 时间戳
}

func (tool *Tool) buildChip(data string){
	// TODO
}
