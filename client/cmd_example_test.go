package client_test

import (
	"flag"
	"log"

	"github.com/eachain/360-tuitui-robot/client"
)

func ExampleClient_SetShortcutCommands() {
	appid := flag.String("appid", "", "tuitui robot appid")
	secret := flag.String("secret", "", "tuitui robot secret")
	flag.Parse()

	cli := client.New(*appid, *secret, nil)

	err := cli.SetShortcutCommands([]client.ShortcutCommand{{
		Name: "基础模型_国风",
		Desc: "2.5D华丽国风风格，擅长人物画",
	}})
	if err != nil {
		log.Printf("set shortcut commands: %v", err)
		return
	}
	log.Printf("set shortcut commands ok")
}
