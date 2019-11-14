package main

import (
	"fmt"
	"net/http"
	"regexp"
	"time"
	"path/filepath"
	"sync"
	"strconv"

	"github.com/xxzl0130/rsyars/pkg/util"
	"github.com/elazarl/goproxy"
)

// 保存登陆签名用于解析数据，外层以IP为map的key
type SignInfo struct{
	sign string // 登陆签名
	time int64  // 时间戳
}

type Tool struct{
	userinfo map[string]UserInfo	// 详细用户数据map
	sign map[string]SignInfo		// 签名map
	host *regexp.Regexp				// 地址正则
	url *regexp.Regexp				// 地址正则
	signMutex *sync.RWMutex			// 锁
	infoMutex *sync.RWMutex			// 锁
	port int						// 代理端口
	timeout int64						// 数据超时
}

func main(){
	tool := &Tool{
		userinfo : make(map[string]UserInfo),
		sign : make(map[string]SignInfo),
		host : regexp.MustCompile(".*\\.ppgame\\.com"),
		url : regexp.MustCompile("(index\\.php(\\/.*\\/Index\\/index)*)|(cn_mica_new\\/.*)|(auth)|(xy\\/.*)|(Config\\/.*)|(normal_login)"),
		signMutex : new(sync.RWMutex),
		infoMutex : new(sync.RWMutex),
		port : 8080,
		timeout : 600, // 10分钟超时
	}
	if err := tool.Run(); err != nil {
		panic(fmt.Sprintf("程序启动失败 -> %+v", err))
	}
}

func (tool *Tool) watchdog(){
	for{
		time.Sleep(time.Second * 60)
		now := time.Now().Unix()

		nSignMap := make(map[string]SignInfo)
		tool.signMutex.RLock()
		for k,v := range tool.sign{
			if now - v.time < tool.timeout{
				// 拷贝未超时的值
				nSignMap[k] = v
			}
		}
		tool.signMutex.RUnlock()
		tool.signMutex.Lock()
		// 覆盖
		tool.sign = nSignMap
		tool.signMutex.Unlock()

		nInfoMap := make(map[string]UserInfo)
		tool.infoMutex.RLock()
		for k,v := range tool.userinfo{
			if now - v.time < tool.timeout{
				// 拷贝未超时的值
				nInfoMap[k] = v
			}
		}
		tool.infoMutex.RUnlock()
		tool.infoMutex.Lock()
		// 覆盖
		tool.userinfo = nInfoMap
		tool.infoMutex.Unlock()
	}
}

func (tool *Tool) Run() error {
	localhost, err := tool.getLocalhost()
	if err != nil {
		fmt.Printf("获取代理地址失败 -> %+v", err)
		return err
	}
	fmt.Printf("代理地址 -> %s:%d", localhost, tool.port)

	proxy := goproxy.NewProxyHttpServer()
	proxy.Logger = new(util.NilLogger)
	proxy.OnResponse(tool.condition()).DoFunc(tool.onResponse)
	proxy.OnRequest(tool.block()).DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			return r, goproxy.NewResponse(r, goproxy.ContentTypeText, http.StatusForbidden, "Forbidden!")
		})
	
	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         ":" + strconv.Itoa(tool.port),
		Handler:      proxy,
	}

	go srv.ListenAndServe()

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("."+string(filepath.Separator)+"PACFile")))
	fileSrv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         ":3333",
		Handler:      mux,
	}
	go tool.watchdog()
	fileSrv.ListenAndServe()
	
	return nil
}
