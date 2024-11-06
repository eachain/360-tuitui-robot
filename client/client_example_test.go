package client_test

import (
	"flag"
	"net/http"
	"time"

	"github.com/eachain/360-tuitui-robot/client"
	"github.com/eachain/360-tuitui-robot/util/transport"
)

func ExampleNew() {
	appid := flag.String("appid", "", "tuitui robot appid")
	secret := flag.String("secret", "", "tuitui robot secret")
	timeout := flag.Duration("timeout", 10*time.Second, "call tuitui robot api timeout")
	log := flag.Bool("log", false, "print the log of calling tuitui robot api")
	wan := flag.Bool("wan", false, "call tuitui robot api from WAN")
	flag.Parse()

	opt := &client.Options{}
	if *timeout > 0 {
		if opt.Client == nil {
			opt.Client = new(http.Client)
		}
		opt.Client.Timeout = *timeout // 设置请求超时时间
	}
	if *log {
		if opt.Client == nil {
			opt.Client = new(http.Client)
		}
		opt.Client.Transport = (*transport.LoggedTransport)(nil) // 记录请求响应日志，方便排查问题
	}
	if *wan {
		// 如果需要从外网访问推推机器人api，需要将请求url前缀改为以下地址
		opt.BaseURL = "https://im.live.360.cn:8282/robot"
	}

	cli := client.New(*appid, *secret, opt)
	_ = cli
}
