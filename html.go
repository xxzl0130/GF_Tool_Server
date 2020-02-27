package main

import (
	"net/http"
	"fmt"
	"text/template"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/gobuffalo/packr"
)

func (tool *Tool)getChip(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	tool.html["chip"].Execute(w,tool.proxyInfo)
}

func (tool *Tool)postChip(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	r.ParseForm()

	chip := ""
	info, isPresent := tool.userinfo[r.PostForm["uid"][0]]
	if isPresent{
		if r.PostForm["name"][0] == info.name{
			rule := r.PostForm["locked"][0] + r.PostForm["equipped"][0]
			if rule != info.rule || info.ruleChanged{
				info.rule = rule
				info.ruleChanged = true
			}else{
				info.ruleChanged = false
			}
			tool.infoMutex.Lock()
			tool.userinfo[r.PostForm["uid"][0]] = info;
			tool.infoMutex.Unlock()
			chip = tool.buildChips(info.uid)
		}else{
			chip = "数据不存在！"
		}
	}else{
		chip = "数据不存在！"
	}

	w.Header().Set("Access-Control-Allow-Origin","*")
	fmt.Fprintln(w, chip)
}

func (tool *Tool)getKalina(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	tool.html["kalina"].Execute(w,tool.proxyInfo)
}

func (tool *Tool)postChipJson(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	r.ParseForm()

	rule := r.PostForm["locked"][0] + r.PostForm["equipped"][0]
	json := ""
	info, isPresent := tool.userinfo[r.PostForm["uid"][0]]
	if isPresent && r.PostForm["name"][0] == info.name{
			if rule != info.rule || info.ruleChanged{
				info.rule = rule
				info.ruleChanged = true
			}else{
				info.ruleChanged = false
			}
			tool.infoMutex.Lock()
			tool.userinfo[r.PostForm["uid"][0]] = info;
			tool.infoMutex.Unlock()
			json = tool.buildChipJson(info.uid)
	}else{
		flag := false
		for k,v := range(tool.userinfo){
			v.rule = rule
			info.ruleChanged = true
			tool.infoMutex.Lock()
			tool.userinfo[k] = v
			tool.infoMutex.Unlock()
			json = tool.buildChipJson(k)
			flag = true
			break
		}
		if !flag{
			json = "数据不存在！"
		}
	}

	w.Header().Set("Access-Control-Allow-Origin","*")
	fmt.Fprintln(w, json)
}

func (tool *Tool)postJson(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	r.ParseForm()

	json := ""
	info, isPresent := tool.userinfo[r.PostForm["uid"][0]]
	if isPresent{
		if r.PostForm["name"][0] == info.name{
			tool.infoMutex.Lock()
			info.time = time.Now().Unix()
			tool.userinfo[info.uid] = info // 更新时间避免被清理
			tool.infoMutex.Unlock()
			json = GzipCompress([]byte(info.allData))
		}else{
			json = "数据不存在！"
		}
	}else{
		json = "数据不存在！"
	}

	w.Header().Set("Access-Control-Allow-Origin","*")
	fmt.Fprintln(w, json)
}

func (tool *Tool)postKalina(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	r.ParseForm()

	w.Header().Set("Access-Control-Allow-Origin","*")
	info, isPresent := tool.userinfo[r.PostForm["uid"][0]]
	if isPresent && r.PostForm["name"][0] == info.name{
		res := tool.buildKalina(info.uid)
		fmt.Fprintln(w,res[1] + ";" + res[0])
	}else{
		fmt.Fprintln(w,"数据不存在！;请检查！")
	}
}

func (tool *Tool)loadHtml(key, file_name string) {
	data,err := tool.readFile(file_name)
	if err != nil {
		fmt.Print(err)
		return
	}
	t, err := template.New("html").Parse(data)
	if err != nil {
		fmt.Print(err)
		return
	}
	tool.html[key] = t
}

func (tool *Tool)readFile(file_name string) (string, error) {
	box := packr.NewBox("./HTML")
	data,err := box.FindString(file_name)
	return data,err
}