package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"text/template"

	"github.com/julienschmidt/httprouter"
)

func (tool *Tool)getChip(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	data := map[string]string{
        "uid":  "\"\"",
		"name": "\"\"",
		"chipcode":"\"\"",
	}
	tool.html["chip"].Execute(w,data)
}

func (tool *Tool)postChip(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	r.ParseForm()

	data := map[string]string{
        "uid":  r.PostForm["uid"][0],
		"name": r.PostForm["name"][0],
		"chipcode":"\"\"",
	}

	info, isPresent := tool.userinfo[r.PostForm["uid"][0]]
	if isPresent{
		if r.PostForm["name"][0] == info.name{
			rule := r.PostForm["locked"][0] + r.PostForm["equipped"][0]
			if rule != info.rule{
				info.rule = rule
				info.ruleChanged = true
			}else{
				info.ruleChanged = false
			}
			data["chipcode"] = tool.buildChips(info.uid)
		}else{
			data["chipcode"] = info.name
		}
	}else{
		data["chipcode"] = "数据不存在！"
	}

	tool.html["chip"].Execute(w,data)
}

func (tool *Tool)getKalina(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	data := map[string]string{
        "uid":  "\"\"",
		"name": "\"\"",
		"kalina":"\"\"",
		"spend":"\"\"",
	}
	tool.html["kalina"].Execute(w,data)
}

func (tool *Tool)postKalina(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	r.ParseForm()

	data := map[string]string{
        "uid":  r.PostForm["uid"][0],
		"name": r.PostForm["name"][0],
		"kalina":"\"\"",
		"spend":"\"\"",
	}

	info, isPresent := tool.userinfo[r.PostForm["uid"][0]]
	if isPresent{
		if r.PostForm["name"][0] == info.name{
			res := tool.buildKalina(info.uid)
			data["kalina"] = res[1]
			data["spend"] = res[0]
		}
	}

	tool.html["kalina"].Execute(w,data)
}

func (tool *Tool)loadHtml(key, file_name string) {
	t, err := template.ParseFiles(file_name)
	if err != nil {
		//fmt.Print(err)
		return
	}
	tool.html[key] = t
}
func (tool *Tool)readFile(file_name string) ([]byte, error) {
	fi, err := os.Open(file_name)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	return ioutil.ReadAll(fi)
}