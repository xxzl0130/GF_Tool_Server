package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"fmt"

	"github.com/julienschmidt/httprouter"
)

func (tool *Tool)getChip(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	fmt.Fprintln(w,tool.html["chip"])
}

func (tool *Tool)postChip(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	r.ParseForm()

	chip := ""
	info, isPresent := tool.userinfo[r.PostForm["uid"][0]]
	fmt.Printf("uid:%v, name:%v\n",r.PostForm["uid"][0],r.PostForm["name"][0])
	if isPresent{
		if r.PostForm["name"][0] == info.name{
			rule := r.PostForm["locked"][0] + r.PostForm["equipped"][0]
			if rule != info.rule{
				info.rule = rule
				info.ruleChanged = true
			}else{
				info.ruleChanged = false
			}
			chip = tool.buildChips(info.uid)
		}else{
			chip = info.name
		}
	}else{
		chip = "数据不存在！"
	}

	w.Header().Set("Access-Control-Allow-Origin","*")
	fmt.Fprintln(w, chip)
}

func (tool *Tool)getKalina(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	fmt.Fprintln(w,tool.html["kalina"])
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
	data,_ := tool.readFile(file_name)
	tool.html[key] = string(data)
}
func (tool *Tool)readFile(file_name string) ([]byte, error) {
	fi, err := os.Open(file_name)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	return ioutil.ReadAll(fi)
}