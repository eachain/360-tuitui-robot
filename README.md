# 360推推IM机器人golang库

开发文档：[推推IM机器人开发文档](https://easydoc.qihoo.net/doc?project=1d414e4d0ce730bec9b805b12ca28509&doc=3596913a227ae858de8e1dcca7dae3d6&config=toc)。

## 功能模块

- client: API客户端，主要实现机器人API行为封装
  - [接口鉴权](https://easydoc.qihoo.net/doc?project=1d414e4d0ce730bec9b805b12ca28509&doc=3596913a227ae858de8e1dcca7dae3d6&config=toc#h2-%E6%8E%A5%E5%8F%A3%E4%B8%8E%E9%89%B4%E6%9D%83)
  - [发消息（消息格式定义见github.com/eachain/360-tuitui-robot/message）](https://easydoc.qihoo.net/doc?project=1d414e4d0ce730bec9b805b12ca28509&doc=3596913a227ae858de8e1dcca7dae3d6&config=toc#h1-2%E3%80%81%E6%9C%BA%E5%99%A8%E4%BA%BA%E5%8F%91%E6%B6%88%E6%81%AF)
  - [修改消息](https://easydoc.qihoo.net/doc?project=1d414e4d0ce730bec9b805b12ca28509&doc=3596913a227ae858de8e1dcca7dae3d6&config=toc#h1-4%E3%80%81%E4%BF%AE%E6%94%B9%E6%B6%88%E6%81%AF)
  - [上传图片、文件](https://easydoc.qihoo.net/doc?project=1d414e4d0ce730bec9b805b12ca28509&doc=3596913a227ae858de8e1dcca7dae3d6&config=toc#h2-%E4%B8%8A%E4%BC%A0%E6%96%87%E4%BB%B6)
  - [电话报警](https://easydoc.qihoo.net/doc?project=1d414e4d0ce730bec9b805b12ca28509&doc=3596913a227ae858de8e1dcca7dae3d6&config=toc#h2-%E7%94%B5%E8%AF%9D%E6%8A%A5%E8%AD%A6)
  - [群管理：建群/拉群成员/踢群成员等](https://easydoc.qihoo.net/doc?project=1d414e4d0ce730bec9b805b12ca28509&doc=3596913a227ae858de8e1dcca7dae3d6&config=toc#h1-6%E3%80%81%E7%BE%A4%E7%AE%A1%E7%90%86%E5%8A%9F%E8%83%BD)
  - [修改机器人属性](https://easydoc.qihoo.net/doc?project=1d414e4d0ce730bec9b805b12ca28509&doc=3596913a227ae858de8e1dcca7dae3d6&config=toc#h1-7%E3%80%81%E6%9C%BA%E5%99%A8%E4%BA%BA%E8%87%AA%E5%8A%A9%E4%BF%AE%E6%94%B9%E5%B1%9E%E6%80%A7)
  - [机器人快捷指令](https://easydoc.qihoo.net/doc?project=1d414e4d0ce730bec9b805b12ca28509&doc=3596913a227ae858de8e1dcca7dae3d6&config=toc#h1-8%E3%80%81%E6%9C%BA%E5%99%A8%E4%BA%BA%E5%BF%AB%E6%8D%B7%E6%8C%87%E4%BB%A4)

- message: 消息类型
  - [text(文本)](https://easydoc.qihoo.net/doc?project=1d414e4d0ce730bec9b805b12ca28509&doc=3596913a227ae858de8e1dcca7dae3d6&config=toc#h2-%E6%96%87%E6%9C%AC%E6%B6%88%E6%81%AF)
  - [image(图片)](https://easydoc.qihoo.net/doc?project=1d414e4d0ce730bec9b805b12ca28509&doc=3596913a227ae858de8e1dcca7dae3d6&config=toc#h2-%E5%9B%BE%E7%89%87%E6%B6%88%E6%81%AF)
  - [mixed(图文混排)](https://easydoc.qihoo.net/doc?project=1d414e4d0ce730bec9b805b12ca28509&doc=3596913a227ae858de8e1dcca7dae3d6&config=toc#h2-%E5%9B%BE%E6%96%87%E6%B7%B7%E6%8E%92%E6%B6%88%E6%81%AF)
  - [file(文件)](https://easydoc.qihoo.net/doc?project=1d414e4d0ce730bec9b805b12ca28509&doc=3596913a227ae858de8e1dcca7dae3d6&config=toc#h2-%E9%99%84%E4%BB%B6%E6%B6%88%E6%81%AF)
  - [link(链接)](https://easydoc.qihoo.net/doc?project=1d414e4d0ce730bec9b805b12ca28509&doc=3596913a227ae858de8e1dcca7dae3d6&config=toc#h2-%E9%93%BE%E6%8E%A5%E6%B6%88%E6%81%AF)
  - [page(页面)](https://easydoc.qihoo.net/doc?project=1d414e4d0ce730bec9b805b12ca28509&doc=3596913a227ae858de8e1dcca7dae3d6&config=toc#h2-%E6%8E%A8%E6%8E%A8%E9%A1%B5%E9%9D%A2%E6%B6%88%E6%81%AF)
  - [recall(撤回，仅可用于修改消息)](https://easydoc.qihoo.net/doc?project=1d414e4d0ce730bec9b805b12ca28509&doc=3596913a227ae858de8e1dcca7dae3d6&config=toc#h2-%E6%92%A4%E5%9B%9E%E6%B6%88%E6%81%AF)
  - [voice(电话报警)](https://easydoc.qihoo.net/doc?project=1d414e4d0ce730bec9b805b12ca28509&doc=3596913a227ae858de8e1dcca7dae3d6&config=toc#h2-%E7%94%B5%E8%AF%9D%E6%8A%A5%E8%AD%A6)
  - [richtext(富文本，仅可用于团队帖子)](https://easydoc.qihoo.net/doc?project=1d414e4d0ce730bec9b805b12ca28509&doc=3596913a227ae858de8e1dcca7dae3d6&config=toc#h2-%E5%9B%A2%E9%98%9F%E5%B8%96%E5%AD%90(HTML))

- webhook: [机器人收消息](https://easydoc.qihoo.net/doc?project=1d414e4d0ce730bec9b805b12ca28509&doc=3596913a227ae858de8e1dcca7dae3d6&config=toc#h1-5%E3%80%81%E6%9C%BA%E5%99%A8%E4%BA%BA%E6%94%B6%E6%B6%88%E6%81%AF)
  - [回调注册](https://easydoc.qihoo.net/doc?project=1d414e4d0ce730bec9b805b12ca28509&doc=3596913a227ae858de8e1dcca7dae3d6&config=toc#h2-%E6%94%B6%E6%B6%88%E6%81%AF%E6%A0%BC%E5%BC%8F)
  - [安全身份验证](https://easydoc.qihoo.net/doc?project=1d414e4d0ce730bec9b805b12ca28509&doc=3596913a227ae858de8e1dcca7dae3d6&config=toc#h2-%E5%AE%89%E5%85%A8%E8%BA%AB%E4%BB%BD%E9%AA%8C%E8%AF%81)

- interactive: [可交互式消息](https://easydoc.soft.360.cn/doc?project=38ed795130e25371ef319aeb60d5b4fa&doc=0750ce7dcf9b9f7589a558a857bc7cb9&config=title_menu_toc#h1-5%E5%8F%AF%E4%BA%A4%E4%BA%92%E5%BC%8F%E6%B6%88%E6%81%AF%28%E5%BE%85%E5%AE%8C%E5%96%84%29)
  - [发消息类型](https://easydoc.soft.360.cn/doc?project=38ed795130e25371ef319aeb60d5b4fa&doc=0750ce7dcf9b9f7589a558a857bc7cb9&config=title_menu_toc#h2-3.%20%E5%AD%97%E6%AE%B5%E8%AF%B4%E6%98%8E)
  - [回调注册](https://easydoc.soft.360.cn/doc?project=38ed795130e25371ef319aeb60d5b4fa&doc=0750ce7dcf9b9f7589a558a857bc7cb9&config=title_menu_toc#h2-5.%20%E6%8C%89%E9%92%AE%E5%9B%9E%E8%B0%83)

- util: 工具包
  - cache: webhook分布式防重放
  - chain: 将多个webhook.Callback合成一个，按顺序调用，每个Callback只注册自己感兴趣的事件
  - logcb: 记录所有webhook.Callback事件日志
  - qa: 机器人自动回复webhook.Callback
  - transport: 将所有client请求及响应记录日志

- example: 使用示例
  - cmder: 简易机器人执行命令
  - grafana_alert: 简易grafana报警机器人（仅以[Grafana v11](https://grafana.com/docs/grafana/v11.1/alerting/configure-notifications/manage-contact-points/integrations/webhook-notifier/)示例，其它版本不保证解析正确性）
  - webhook_dev: 开发阶段用来查看webhook事件及参数

## 使用

示例：单聊发文本消息。

```go
package main

import (
	"flag"
	"log"

	"github.com/eachain/360-tuitui-robot/client"
	"github.com/eachain/360-tuitui-robot/message"
)

func main() {
	appid := flag.String("appid", "", "tuitui robot appid")
	secret := flag.String("secret", "", "tuitui robot secret")
	user := flag.String("user", "", "send message to user")
	text := flag.String("text", "Hello tuitui robot", "text message content")
	flag.Parse()

	// 初始化client
	cli := client.New(*appid, *secret, nil)

	// 发文本消息给*user，返回消息id
	msgid, err := cli.SendMessageToUser(*user, message.NewText(*text))
	if err != nil {
		log.Printf("send message to user: %v", err)
		return
	}
	log.Printf("send messages to user %v ok, message id: %v", *user, msgid)
}
```

更多使用案例见`github.com/eachain/360-tuitui-robot/example`。
