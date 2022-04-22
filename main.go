package main

import (
	"fmt"
	"net/http"
	"time"
	"sync"
	"strconv"
	"regexp"
	"os"
	"math/rand"
	//"io/ioutil"
	"text/template"

	"github.com/xxzl0130/rsyars/pkg/util"
	"github.com/elazarl/goproxy"
	"github.com/julienschmidt/httprouter"
)

// 保存登陆签名用于解析数据，外层以IP为map的key
type SignInfo struct{
	sign string // 登陆签名
	time int64  // 时间戳
}

type Tool struct{
	userinfo map[string]UserInfo	// 详细用户数据map
	sign map[string]SignInfo		// 签名map
	signMutex *sync.RWMutex			// 锁
	infoMutex *sync.RWMutex			// 锁
	port int						// 代理端口
	timeout int64					// 数据超时
	html map[string]*template.Template			// 网页数据
	proxyInfo map[string]string     // 显示在HTML中的代理信息
}

func main(){
	rand.Seed(time.Now().Unix())
	tool := &Tool{
		userinfo : make(map[string]UserInfo),
		sign : make(map[string]SignInfo),
		signMutex : new(sync.RWMutex),
		infoMutex : new(sync.RWMutex),
		port : 8080,
		timeout : 600, // 10分钟超时
		html : make(map[string]*template.Template),
		proxyInfo : make(map[string]string),
	}
	if err := tool.Run(); err != nil {
		panic(fmt.Sprintf("程序启动失败 -> %+v", err))
	}
}

func (tool *Tool) Run() error {
	localhost, err := tool.getLocalhost()
	if err != nil {
		fmt.Printf("获取代理地址失败 -> %+v", err)
		return err
	}
	tool.proxyInfo["proxy"] = fmt.Sprintf("%s:%d", localhost, tool.port + 1)
	tool.proxyInfo["host"] = fmt.Sprintf("\"http://127.0.0.1:%d\"", tool.port)

	// 代理服务
	proxy := goproxy.NewProxyHttpServer()
	proxy.Logger = new(util.NilLogger)
	proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile("sn-game.txwy.tw"))).
		HandleConnect(goproxy.AlwaysMitm)
	proxy.OnResponse(tool.condition()).DoFunc(tool.onResponse)
	srv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         ":" + strconv.Itoa(tool.port + 1),
		Handler:      proxy,
	}
	go func ()  {
		if err:=srv.ListenAndServe(); err != nil{
			// 随机端口尝试3次
			for i := 0; i < 3;i++{
				tool.port = rand.Intn(50000) + 10000
				srv := &http.Server{
					ReadTimeout:  5 * time.Second,
					WriteTimeout: 10 * time.Second,
					Addr:         ":" + strconv.Itoa(tool.port + 1),
					Handler:      proxy,
				}
				if err=srv.ListenAndServe(); err != nil{
					continue
				}
			}
			fmt.Printf("启动代理服务器失败 -> %+v", err)
			os.Exit(1)
		}
	}()

	// 网页服务
	tool.loadHtml("chip","chip.html")
	tool.loadHtml("kalina","kalina.html")
	router := httprouter.New()
	router.POST("/chip", tool.postChip)
	router.GET("/chip", tool.getChip)
	router.GET("/", tool.getChip)
	router.POST("/chipJson", tool.postChipJson)
	router.POST("/kalina", tool.postKalina)
	router.GET("/kalina", tool.getKalina)
	httpSrv := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         ":" + strconv.Itoa(tool.port),
		Handler:      router,
	}
	fmt.Printf("代理地址 -> %s:%d\n", localhost, tool.port + 1)
	fmt.Printf("网页地址 -> %s:%d\n", localhost, tool.port)

	//data,_ := ioutil.ReadFile("response.json")
	//tool.saveUserInfo(string(data))

	//if err := httpSrv.ListenAndServeTLS("./_.xuanxuan.tech_chain.crt","./_.xuanxuan.tech_key.key"); err != nil {
	if err := httpSrv.ListenAndServe(); err != nil {
		fmt.Printf("启动代理服务器失败 -> %+v", err)
	}
	return nil
}
