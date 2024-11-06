package client_test

import (
	"encoding/json"
	"flag"
	"log"

	"github.com/eachain/360-tuitui-robot/client"
)

func ExampleClient_Call() {
	appid := flag.String("appid", "", "tuitui robot appid")
	secret := flag.String("secret", "", "tuitui robot secret")
	flag.Parse()

	cli := client.New(*appid, *secret, nil)

	const api = "/teams/team/list"
	var reply json.RawMessage
	err := cli.Call(api, nil, &reply)
	if err != nil {
		log.Printf("call %v: %v", api, err)
		return
	}
	log.Printf("call %v: %s", api, reply)
}
