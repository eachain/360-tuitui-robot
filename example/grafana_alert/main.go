package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/eachain/360-tuitui-robot/client"
	"github.com/eachain/360-tuitui-robot/util/transport"
)

// 简易grafana报警示例。
func main() {
	appid := flag.String("appid", "", "tuitui robot appid")
	secret := flag.String("secret", "", "tuitui robot secret")
	timeout := flag.Duration("timeout", 10*time.Second, "tuitui robot api call timeout")
	listen := flag.String("webhook", ":8080", "tuitui robot and grafana webhook both listen address")
	user := flag.String("user", "", "send grafana alert message to user")
	group := flag.String("group", "", "send grafana alert message to group")
	flag.Parse()

	// 用client发消息报警
	cli := client.New(*appid, *secret, &client.Options{
		Client: &http.Client{
			Timeout:   *timeout,
			Transport: (*transport.LoggedTransport)(nil),
		},
	})

	panic(http.ListenAndServe(*listen, NewGrafanaV11Webhook(&Options{
		Client: cli,
		User:   *user,
		Group:  *group,
		Logf:   log.Printf,
	})))
}
