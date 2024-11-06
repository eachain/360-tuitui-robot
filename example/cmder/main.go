package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/eachain/360-tuitui-robot/client"
	"github.com/eachain/360-tuitui-robot/message"
	"github.com/eachain/360-tuitui-robot/util/chain"
	"github.com/eachain/360-tuitui-robot/util/logcb"
	"github.com/eachain/360-tuitui-robot/util/qa"
	"github.com/eachain/360-tuitui-robot/util/transport"
	"github.com/eachain/360-tuitui-robot/webhook"
)

type greetOptions struct {
	Language string `option:"lang" default:"cn" usage:"语言，支持cn/en"`
}

func greet(opts *greetOptions, args []string) string {
	if len(args) == 0 {
		return "参数错误：至少需要一个姓名"
	}
	if opts.Language == "en" {
		return "Hello " + strings.Join(args, ", ") + "!"
	}
	if len(args) == 1 {
		return args[0] + "你好！"
	}
	return strings.Join(args, "，") + "，你们好！"
}

func now() string {
	now := time.Now()
	return fmt.Sprintf("%v\nunix timestamp: %v",
		now.Format("2006-01-02 15:04:05.000"), now.Unix())
}

type method struct {
	cli *client.Client
}

type sendOptions struct {
	User    string `option:"user" usage:"发给谁"`
	Group   string `option:"group" usage:"发给哪个群"`
	Team    string `option:"team" usage:"发给哪个团队"`
	Channel string `option:"channel" usage:"发给哪个频道"`
	Text    string `option:"text" default:"Hello tuitui robot" usage:"发送内容"`
}

func (m *method) send(opts *sendOptions) string {
	if opts.Text == "" {
		return ""
	}

	var results []string

	if opts.User != "" {
		msgid, err := m.cli.SendMessageToUser(opts.User, message.NewText(opts.Text))
		if err != nil {
			results = append(results, fmt.Sprintf("user %v: %v", opts.User, err))
		} else {
			results = append(results, fmt.Sprintf("user %v: %v", opts.User, msgid))
		}
	}

	if opts.Group != "" {
		msgid, err := m.cli.SendMessageToGroup(opts.Group, message.NewText(opts.Text))
		if err != nil {
			results = append(results, fmt.Sprintf("group %v: %v", opts.Group, err))
		} else {
			results = append(results, fmt.Sprintf("group %v: %v", opts.Group, msgid))
		}
	}

	if opts.Team != "" && opts.Channel != "" {
		msgid, err := m.cli.SendPostToTeam(client.TeamChannel{
			TeamId:    opts.Team,
			ChannelId: opts.Channel,
		}, message.NewRichTextHTML(opts.Text))
		if err != nil {
			results = append(results, fmt.Sprintf("team %v channel %v: %v",
				opts.Team, opts.Channel, err))
		} else {
			results = append(results, fmt.Sprintf("team %v channel %v: %v",
				opts.Team, opts.Channel, msgid))
		}
	}

	if len(results) == 0 {
		return "参数错误：user, group, team/channel不能全为空"
	}

	return strings.Join(results, "\n")
}

func closureBase64(cmder *Cmder) {
	cmder.Register(func(opts *struct {
		Decode bool `option:"d" usage:"decode the text"`
	}, args []string) string {
		if len(args) == 0 {
			return "参数错误，缺少要编/解码的文本"
		}
		for i, arg := range args {
			fmt.Printf("%v: %q\n", i, arg)
		}
		if len(args) == 1 {
			if opts.Decode {
				p, err := base64.StdEncoding.DecodeString(args[0])
				if err != nil {
					return err.Error()
				}
				return string(p)
			}
			return base64.StdEncoding.EncodeToString([]byte(args[0]))
		}

		buf := new(bytes.Buffer)
		for i, arg := range args {
			if i > 0 {
				fmt.Fprintln(buf)
			}
			if opts.Decode {
				p, err := base64.StdEncoding.DecodeString(arg)
				if err != nil {
					fmt.Fprintf(buf, "%v: %v", i+1, err)
				} else {
					fmt.Fprintf(buf, "%v: %s", i+1, p)
				}
			} else {
				fmt.Fprintf(buf, "%v: %s", i+1, base64.StdEncoding.EncodeToString([]byte(arg)))
			}
		}
		return buf.String()
	}, "base64", "base64 encoding text")
}

// 简易机器人执行命令示例。
func main() {
	appid := flag.String("appid", "", "tuitui robot appid")
	secret := flag.String("secret", "", "tuitui robot secret")
	timeout := flag.Duration("timeout", 10*time.Second, "tuitui robot api call timeout")
	listen := flag.String("webhook", ":8080", "tuitui robot webhook listen address")
	flag.Parse()

	// 回复执行结果需要client发消息功能
	cli := client.New(*appid, *secret, &client.Options{
		Client: &http.Client{
			Timeout:   *timeout,
			Transport: (*transport.LoggedTransport)(nil),
		},
	})

	// 命令注册
	cmder := New("/")
	// 注册命令"/greet"
	cmder.Register(greet, "", "greet users")
	// 注册命令"/now"
	cmder.Register(now, "", "now time")
	// 用闭包的方式注册命令"/base64"
	closureBase64(cmder)

	// 注册命令"/echo"
	cmder.Register(func(args []string) string {
		return strings.Join(args, " ")
	}, "echo", "")

	// 以结构体方法的方式注册命令"/send"
	md := &method{cli: cli}
	cmder.Register(md.send, "", "send message or post to user, group or teams")

	props, err := cli.GetRobotProps()
	if err != nil {
		panic(err)
	}
	// 将收到的消息做为question，cmder执行结果作为answer。
	cb := qa.New(cmder.Exec, cli, &qa.Options{
		TrimAtMe:  true,
		RobotName: props.Name,
		Errorf:    log.Printf,
	}).Webhook()

	panic(http.ListenAndServe(*listen, webhook.WithAuthSign(&webhook.AuthOptions{
		Appid:  *appid,
		Secret: *secret,
		Expire: 10 * time.Second,
		Cache:  webhook.NewMemCache(15 * time.Second),
		Errorf: log.Printf,
	}, webhook.NewHandler(chain.Callbacks(logcb.Logged(logcb.Printf), cb), &webhook.Options{
		Errorf: log.Printf,
	}))))
}
