package client_test

import (
	"flag"
	"log"

	"github.com/eachain/360-tuitui-robot/client"
	"github.com/eachain/360-tuitui-robot/message"
)

func ExampleClient_SendMessageToUsers() {
	appid := flag.String("appid", "", "tuitui robot appid")
	secret := flag.String("secret", "", "tuitui robot secret")
	user := flag.String("user", "", "send message to user")
	text := flag.String("text", "Hello tuitui robot", "text message content")
	flag.Parse()

	cli := client.New(*appid, *secret, nil)

	msgs, warn, err := cli.SendMessageToUsers([]string{*user}, message.NewText(*text))
	if err != nil {
		log.Printf("send message to users: %v", err)
		return
	}
	if warn != nil {
		log.Printf("send messages failed users: %v", warn.Explains)
	}
	log.Printf("send messages to users ok, messages: %v", msgs)
}

func ExampleClient_SendMessageToUser() {
	appid := flag.String("appid", "", "tuitui robot appid")
	secret := flag.String("secret", "", "tuitui robot secret")
	user := flag.String("user", "", "send message to user")
	text := flag.String("text", "Hello tuitui robot", "text message content")
	flag.Parse()

	cli := client.New(*appid, *secret, nil)

	msgid, err := cli.SendMessageToUser(*user, message.NewText(*text))
	if err != nil {
		log.Printf("send message to user: %v", err)
		return
	}
	log.Printf("send messages to user ok, message id: %v", msgid)
}

func ExampleClient_ModifyUserMessages() {
	appid := flag.String("appid", "", "tuitui robot appid")
	secret := flag.String("secret", "", "tuitui robot secret")
	user := flag.String("user", "", "modify message to user")
	// msgid := flag.String("msgid", "", "modify message id")
	// file := flag.String("file", "993a9e16b365f4529fc9ccfa", "file message media id")
	flag.Parse()

	cli := client.New(*appid, *secret, nil)

	msgid, err := cli.SendMessageToUser(*user, message.NewText("Hello world"))
	if err != nil {
		log.Printf("send message to user %v: %v", *user, err)
		return
	}

	msgs, warn, err := cli.ModifyUserMessages([]client.UserMsgIdPair{{User: *user, MsgId: msgid}},
		message.NewText("Hello tuitui robot"), nil)
	if err != nil {
		log.Printf("modify messages to users: %v", err)
		return
	}
	if warn != nil {
		log.Printf("modify messages failed users: %v", warn.Explains)
	}
	log.Printf("modify messages to users ok, messages: %v", msgs)
}

func ExampleClient_ModifyUserMessage() {
	appid := flag.String("appid", "", "tuitui robot appid")
	secret := flag.String("secret", "", "tuitui robot secret")
	user := flag.String("user", "", "modify message to user")
	flag.Parse()

	cli := client.New(*appid, *secret, nil)

	msgid, err := cli.SendMessageToUser(*user, message.NewText("Hello world"))
	if err != nil {
		log.Printf("send message to user %v: %v", *user, err)
		return
	}

	err = cli.ModifyUserMessage(client.UserMsgIdPair{User: *user, MsgId: msgid},
		message.NewText("Hello tuitui robot"), nil)
	if err != nil {
		log.Printf("modify message to user: %v", err)
		return
	}
	log.Printf("modify message to user ok")
}

func ExampleClient_SendMessageToGroups() {
	appid := flag.String("appid", "", "tuitui robot appid")
	secret := flag.String("secret", "", "tuitui robot secret")
	group := flag.String("group", "", "send message to group id")
	at := flag.String("at", "", "at member in message")
	text := flag.String("text", "Hello tuitui robot", "text message content")
	flag.Parse()

	cli := client.New(*appid, *secret, nil)

	var atUsers []string
	if *at != "" {
		atUsers = append(atUsers, *at)
	}

	msgs, warn, err := cli.SendMessageToGroups([]string{*group}, atUsers, message.NewText(*text))
	if err != nil {
		log.Printf("send messages to groups: %v", err)
		return
	}
	if warn != nil {
		log.Printf("send messages failed groups: %v", warn.Explains)
	}
	log.Printf("send messages to groups ok, messages: %v", msgs)
}

func ExampleClient_SendMessageToGroupAt() {
	appid := flag.String("appid", "", "tuitui robot appid")
	secret := flag.String("secret", "", "tuitui robot secret")
	group := flag.String("group", "", "send message to group id")
	at := flag.String("at", "", "at member in message")
	text := flag.String("text", "Hello tuitui robot", "text message content")
	flag.Parse()

	cli := client.New(*appid, *secret, nil)

	var atUsers []string
	if *at != "" {
		atUsers = append(atUsers, *at)
	}

	msgid, err := cli.SendMessageToGroupAt(*group, atUsers, message.NewText(*text))
	if err != nil {
		log.Printf("send message to group: %v", err)
		return
	}
	log.Printf("send message to group ok, message id: %v", msgid)
}

func ExampleClient_SendMessageToGroup() {
	appid := flag.String("appid", "", "tuitui robot appid")
	secret := flag.String("secret", "", "tuitui robot secret")
	group := flag.String("group", "", "send message to group id")
	text := flag.String("text", "Hello tuitui robot", "text message content")
	flag.Parse()

	cli := client.New(*appid, *secret, nil)

	msgid, err := cli.SendMessageToGroup(*group, message.NewText(*text))
	if err != nil {
		log.Printf("send message to group: %v", err)
		return
	}
	log.Printf("send message to group ok, message id: %v", msgid)
}

func ExampleClient_ModifyGroupMessages() {
	appid := flag.String("appid", "", "tuitui robot appid")
	secret := flag.String("secret", "", "tuitui robot secret")
	group := flag.String("group", "", "send and modify message to group id")
	flag.Parse()

	cli := client.New(*appid, *secret, nil)

	msgid, err := cli.SendMessageToGroup(*group, message.NewText("Hello world"))
	if err != nil {
		log.Printf("send message to group: %v", err)
		return
	}

	msgs, warn, err := cli.ModifyGroupMessages([]client.GroupMsgIdPair{{Group: *group, MsgId: msgid}},
		nil, message.NewText("Hello tuitui robot"), nil)
	if err != nil {
		log.Printf("modify group messages: %v", err)
		return
	}
	if warn != nil {
		log.Printf("modify group messages failed messages: %v", warn.Explains)
	}
	log.Printf("modify group messages success messages: %v", msgs)
}

func ExampleClient_ModifyGroupMessageAt() {
	appid := flag.String("appid", "", "tuitui robot appid")
	secret := flag.String("secret", "", "tuitui robot secret")
	group := flag.String("group", "", "send and modify message to group id")
	at := flag.String("at", "", "at member in message")
	flag.Parse()

	cli := client.New(*appid, *secret, nil)

	msgid, err := cli.SendMessageToGroup(*group, message.NewText("Hello world"))
	if err != nil {
		log.Printf("send message to group: %v", err)
		return
	}

	var atUsers []string
	if *at != "" {
		atUsers = append(atUsers, *at)
	}

	err = cli.ModifyGroupMessageAt(client.GroupMsgIdPair{Group: *group, MsgId: msgid},
		atUsers, message.NewText("Hello tuitui robot"), nil)
	if err != nil {
		log.Printf("modify group message: %v", err)
		return
	}
	log.Printf("modify group message OK")
}

func ExampleClient_ModifyGroupMessage() {
	appid := flag.String("appid", "", "tuitui robot appid")
	secret := flag.String("secret", "", "tuitui robot secret")
	group := flag.String("group", "", "send and modify message to group id")
	flag.Parse()

	cli := client.New(*appid, *secret, nil)

	msgid, err := cli.SendMessageToGroup(*group, message.NewText("Hello world"))
	if err != nil {
		log.Printf("send message to group: %v", err)
		return
	}

	err = cli.ModifyGroupMessage(client.GroupMsgIdPair{Group: *group, MsgId: msgid},
		message.NewText("Hello tuitui robot"), nil)
	if err != nil {
		log.Printf("modify group message: %v", err)
		return
	}
	log.Printf("modify group message OK")
}

func ExampleClient_SendPageToUsers() {
	appid := flag.String("appid", "", "tuitui robot appid")
	secret := flag.String("secret", "", "tuitui robot secret")
	user := flag.String("user", "", "send page message to user")
	flag.Parse()

	cli := client.New(*appid, *secret, nil)

	page := message.NewPage().WithTitle("标题").WithImage("993a9e16b365f4529fc9ccfa").
		WithSummary("简介").
		WithContent("<h2>Hello tuitui robot</h2>")

	pageId, msgs, warn, err := cli.SendPageToUsers([]string{*user}, page)
	if err != nil {
		log.Printf("send page to users: %v", err)
		return
	}
	if warn != nil {
		log.Printf("send page failed users: %v", warn.Explains)
	}
	log.Printf("send page to users messages: %v", msgs)
	log.Printf("send page result page_id: %v", pageId)
}

func ExampleClient_SendPageToGroups() {
	appid := flag.String("appid", "", "tuitui robot appid")
	secret := flag.String("secret", "", "tuitui robot secret")
	group := flag.String("group", "", "send page message to group id")
	flag.Parse()

	cli := client.New(*appid, *secret, nil)

	page := message.NewPage().WithTitle("标题").WithImage("993a9e16b365f4529fc9ccfa").
		WithSummary("简介").
		WithContent("<h2>Hello tuitui robot</h2>")

	pageId, msgs, warn, err := cli.SendPageToGroups([]string{*group}, page)
	if err != nil {
		log.Printf("send page to groups: %v", err)
		return
	}
	if warn != nil {
		log.Printf("send page failed groups: %v", warn.Explains)
	}
	log.Printf("send page to groups messages: %v", msgs)
	log.Printf("send page result page_id: %v", pageId)
}

func ExampleClient_ModifyPageContent() {
	appid := flag.String("appid", "", "tuitui robot appid")
	secret := flag.String("secret", "", "tuitui robot secret")
	pageId := flag.String("page", "", "send page message to page id")
	flag.Parse()

	cli := client.New(*appid, *secret, nil)

	page := message.NewPage().WithPageId(*pageId).WithTitle("标题").
		WithImage("993a9e16b365f4529fc9ccfa").
		WithSummary("简介").
		WithContent("<h2>Hello tuitui robot</h2>").
		WithDebug(true)

	err := cli.ModifyPageContent(page)
	if err != nil {
		log.Printf("modify page content: %v", err)
		return
	}
	log.Printf("modify page content OK")
}

func ExampleClient_SendPostToTeams() {
	appid := flag.String("appid", "", "tuitui robot appid")
	secret := flag.String("secret", "", "tuitui robot secret")
	team := flag.String("team", "", "send post message to team id")
	channel := flag.String("channel", "", "send post message to channel id")
	parent := flag.String("parent", "", "reply post message to post id")
	flag.Parse()

	cli := client.New(*appid, *secret, nil)

	post := message.NewRichTextHTML("Hello tuitui robot")
	results, warn, err := cli.SendPostToTeams([]client.TeamChannel{{
		TeamId:    *team,
		ChannelId: *channel,
		ParentId:  *parent,
	}}, post)
	if err != nil {
		log.Printf("send post to teams: %v", err)
		return
	}
	if warn != nil {
		log.Printf("send post fails: %v", warn.Explains)
	}
	log.Printf("send post results: %v", results)
}

func ExampleClient_SendPostToTeam() {
	appid := flag.String("appid", "", "tuitui robot appid")
	secret := flag.String("secret", "", "tuitui robot secret")
	team := flag.String("team", "", "send post message to team id")
	channel := flag.String("channel", "", "send post message to channel id")
	parent := flag.String("parent", "", "reply post message to post id")
	flag.Parse()

	cli := client.New(*appid, *secret, nil)

	post := message.NewRichTextHTML("Hello tuitui robot")
	postId, err := cli.SendPostToTeam(client.TeamChannel{
		TeamId:    *team,
		ChannelId: *channel,
		ParentId:  *parent,
	}, post)
	if err != nil {
		log.Printf("send post to teams: %v", err)
		return
	}
	log.Printf("send result post id %v", postId)
}

func ExampleClient_ModifyTeamPosts() {
	appid := flag.String("appid", "", "tuitui robot appid")
	secret := flag.String("secret", "", "tuitui robot secret")
	team := flag.String("team", "", "modify post message to team id")
	channel := flag.String("channel", "", "modify post message to channel id")
	flag.Parse()

	cli := client.New(*appid, *secret, nil)

	post := message.NewRichTextHTML("Hello world")
	postId, err := cli.SendPostToTeam(client.TeamChannel{
		TeamId:    *team,
		ChannelId: *channel,
	}, post)
	if err != nil {
		log.Printf("send post to teams: %v", err)
		return
	}

	post = message.NewRichTextHTML("Hello tuitui robot")
	results, warn, err := cli.ModifyTeamPosts([]client.ModifyTeamPostRequest{{
		TeamId:    *team,
		ChannelId: *channel,
		PostId:    postId,
	}}, post)
	if err != nil {
		log.Printf("modify post to teams: %v", err)
		return
	}
	if warn != nil {
		log.Printf("modify post fails: %v", warn.Explains)
	}
	log.Printf("modify post results: %v", results)
}

func ExampleClient_ModifyTeamPost() {
	appid := flag.String("appid", "", "tuitui robot appid")
	secret := flag.String("secret", "", "tuitui robot secret")
	team := flag.String("team", "", "modify post message to team id")
	channel := flag.String("channel", "", "modify post message to channel id")
	flag.Parse()

	cli := client.New(*appid, *secret, nil)

	post := message.NewRichTextHTML("Hello world")
	postId, err := cli.SendPostToTeam(client.TeamChannel{
		TeamId:    *team,
		ChannelId: *channel,
	}, post)
	if err != nil {
		log.Printf("send post to teams: %v", err)
		return
	}

	post = message.NewRichTextHTML("<p>Hello tuitui robot</p>")
	err = cli.ModifyTeamPost(client.ModifyTeamPostRequest{
		TeamId:    *team,
		ChannelId: *channel,
		PostId:    postId,
	}, post)
	if err != nil {
		log.Printf("modify post to team: %v", err)
		return
	}
	log.Printf("modify post OK")
}

func ExampleClient_SendVoiceToUsers() {
	appid := flag.String("appid", "", "tuitui robot appid")
	secret := flag.String("secret", "", "tuitui robot secret")
	user := flag.String("user", "", "send voice to user")
	flag.Parse()

	cli := client.New(*appid, *secret, nil)

	// message.NewVoice("用户服务")对应文案为：
	// "服务端报警，报警资源名用户服务，详情请查看相关信息。"
	results, err := cli.SendVoiceToUsers([]string{*user}, message.NewVoice("用户服务"))
	if err != nil {
		log.Printf("send voice to users: %v", err)
		return
	}
	log.Printf("send voice to users results: %v", results)
}

func ExampleClient_QueryVoiceDetail() {
	appid := flag.String("appid", "", "tuitui robot appid")
	secret := flag.String("secret", "", "tuitui robot secret")
	callId := flag.String("call", "", "voice call_id")
	flag.Parse()

	cli := client.New(*appid, *secret, nil)

	// message.NewVoice("用户服务")对应文案为：
	// "服务端报警，报警资源名用户服务，详情请查看相关信息。"
	detail, err := cli.QueryVoiceDetail(*callId)
	if err != nil {
		log.Printf("query voice detail: %v", err)
		return
	}
	log.Printf("query voice detail: %v", detail)
}
