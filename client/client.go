package client

import (
	"net/http"
	"net/url"
)

// 机器人api客户端选项，可以提供自定义http.Client和服务端地址。
type Options struct {
	// 自定义http.Client，常见用法有：
	// 1. 自定义http.Client用到的Transport，比如记录详细日志，方便查看机器人客户端的行为（参考robot/util/transport.LoggedTransport）；
	// 2. 自定义http.Client.Timeout，超时控制。
	// 默认使用http.DefaultClient。
	Client *http.Client

	// 机器人api服务端地址，默认为"https://alarm.im.qihoo.net"。
	// 如果业务需要外网访问机器人api，可以将BaseURL设为"https://im.live.360.cn:8282/robot"。
	BaseURL string
}

// 机器人api客户端，将机器人api封装为相应方法，简化业务开发流程。
//
// 文档地址：https://easydoc.qihoo.net/doc?project=1d414e4d0ce730bec9b805b12ca28509&doc=3596913a227ae858de8e1dcca7dae3d6&config=toc。
type Client struct {
	base   string
	appid  string
	secret string
	query  string

	cli *http.Client
}

// 新建客户端，必须提供appid/secret，*Options可以为空（详见Options定义/默认值）。
func New(appid, secret string, opt *Options) *Client {
	query := make(url.Values, 2)
	query.Add("appid", appid)
	query.Add("secret", secret)
	cli := &Client{
		base:   "https://alarm.im.qihoo.net",
		appid:  appid,
		secret: secret,
		query:  query.Encode(),
	}

	if opt != nil {
		if opt.Client != nil {
			cli.cli = opt.Client
		}
		if opt.BaseURL != "" {
			cli.base = opt.BaseURL
		}
	}
	return cli
}
