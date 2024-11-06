package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/eachain/360-tuitui-robot/util/chain"
	"github.com/eachain/360-tuitui-robot/util/logcb"
	"github.com/eachain/360-tuitui-robot/webhook"
)

// webhook_dev可用于根据回调事件，开发对应业务。
func main() {
	appid := flag.String("appid", "", "tuitui robot appid")
	secret := flag.String("secret", "", "tuitui robot secret")
	listen := flag.String("webhook", ":8080", "tuitui robot webhook listen address")
	flag.Parse()

	// 注册业务回调函数，在对应函数中做业务逻辑
	biz := webhook.Callback{
		OnReceiveSingleMessage:  OnReceiveSingleMessage,
		OnReceiveGroupMessage:   OnReceiveGroupMessage,
		OnCreateGroup:           OnCreateGroup,
		OnNewMemberJoinGroup:    OnNewMemberJoinGroup,
		OnGroupKickMember:       OnGroupKickMember,
		OnCreateTeamsPost:       OnCreateTeamsPost,
		OnModifyTeamsPost:       OnModifyTeamsPost,
		OnAddTeamsMember:        OnAddTeamsMember,
		OnRemoveTeamsMember:     OnRemoveTeamsMember,
		OnCreateTeamsChannel:    OnCreateTeamsChannel,
		OnDeleteTeamsChannel:    OnDeleteTeamsChannel,
		OnCreateTeamsChannelTab: OnCreateTeamsChannelTab,
		OnDeleteTeamsChannelTab: OnDeleteTeamsChannelTab,
	}

	panic(http.ListenAndServe(*listen, webhook.WithAuthSign(&webhook.AuthOptions{
		Appid:  *appid,
		Secret: *secret,
		Expire: 10 * time.Second,
		Cache:  webhook.NewMemCache(15 * time.Second),
		Errorf: log.Printf,
	}, webhook.NewHandler(
		chain.Callbacks(logcb.Logged(logcb.Printf), biz),
		&webhook.Options{Errorf: log.Printf},
	))))
}
