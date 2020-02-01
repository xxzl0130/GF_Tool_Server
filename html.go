package main

import (
	"net/http"
	"fmt"

	"github.com/julienschmidt/httprouter"
)

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

func (tool *Tool)postChipJson(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	r.ParseForm()

	json := ""
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
			json = tool.buildChipJson(info.uid)
		}else{
			json = info.name
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