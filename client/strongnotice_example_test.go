package client_test

import (
	"flag"
	"log"

	"github.com/eachain/360-tuitui-robot/client"
)

func ExampleClient_SendSingleStrongNotice() {
	appid := flag.String("appid", "", "tuitui robot appid")
	secret := flag.String("secret", "", "tuitui robot secret")
	touser := flag.String("user", "", "send tuitui strong notice to whom")
	flag.Parse()

	cli := client.New(*appid, *secret, nil)

	err := cli.SendSingleStrongNotice(*touser, "你好，我有急事找你")
	if err != nil {
		log.Printf("send single strong notice to %v: %v", *touser, err)
		return
	}
	log.Printf("send single strong notice to %v OK", *touser)

	err = cli.SendSingleStrongNotice(
		*touser,
		"你好，我有急事找你。如果超过1分钟未接收强通知，将发短信提醒你。",
		client.WithSMSNotice(),
	)
	if err != nil {
		log.Printf("send single strong notice with sms to %v: %v", *touser, err)
		return
	}
	log.Printf("send single strong notice with sms to %v OK", *touser)
}
