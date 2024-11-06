package client_test

import (
	"flag"
	"log"

	"github.com/eachain/360-tuitui-robot/client"
)

func ExampleClient_CreateGroup() {
	appid := flag.String("appid", "", "tuitui robot appid")
	secret := flag.String("secret", "", "tuitui robot secret")
	name := flag.String("name", "", "group name")
	owner := flag.String("owner", "", "group owner")
	member := flag.String("member", "", "group member")
	flag.Parse()

	cli := client.New(*appid, *secret, nil)

	groupId, warn, err := cli.CreateGroup(*name, *owner, []string{*member})
	if err != nil {
		log.Printf("create group: %v", err)
		return
	}
	if warn != nil {
		log.Printf("create group failed members: %v", warn.Explains)
	}
	log.Printf("create group ok, group id: %v", groupId)
}

func ExampleClient_AddGroupMembers() {
	appid := flag.String("appid", "", "tuitui robot appid")
	secret := flag.String("secret", "", "tuitui robot secret")
	group := flag.String("group", "", "group id")
	member := flag.String("member", "", "user to add to group")
	flag.Parse()

	cli := client.New(*appid, *secret, nil)

	warn, err := cli.AddGroupMembers(*group, []string{*member})
	if err != nil {
		log.Printf("add group member: %v", err)
		return
	}
	if warn != nil {
		log.Printf("add group member failed users: %v", warn.Explains)
	}
}

func ExampleClient_RemoveGroupMembers() {
	appid := flag.String("appid", "", "tuitui robot appid")
	secret := flag.String("secret", "", "tuitui robot secret")
	group := flag.String("group", "", "group id")
	member := flag.String("member", "", "member to remove from group")
	flag.Parse()

	cli := client.New(*appid, *secret, nil)

	warn, err := cli.AddGroupMembers(*group, []string{*member})
	if err != nil {
		log.Printf("remove group member: %v", err)
		return
	}
	if warn != nil {
		log.Printf("remove group member failed users: %v", warn.Explains)
	}
}
