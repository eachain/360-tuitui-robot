package main

import "github.com/eachain/360-tuitui-robot/webhook"

// 单聊消息回调
func OnReceiveSingleMessage(event webhook.SingleMessageEvent) {
	// do your business here
}

// 群聊消息回调
func OnReceiveGroupMessage(event webhook.GroupMessageEvent) {
	// do your business here
}

// 建群事件
func OnCreateGroup(event webhook.GroupMemberEvent) {
	// do your business here
}

// 新成员进群回调
func OnNewMemberJoinGroup(event webhook.GroupMemberEvent) {
	// do your business here
}

// 踢群成员回调
func OnGroupKickMember(event webhook.GroupMemberEvent) {
	// do your business here
}

// 团队创建帖子回调
func OnCreateTeamsPost(event webhook.TeamsPostEvent) {
	// do your business here
}

// 团队修改帖子回调
func OnModifyTeamsPost(event webhook.TeamsPostEvent) {
	// do your business here
}

// 团队添加成员回调
func OnAddTeamsMember(event webhook.TeamsMemberEvent) {
	// do your business here
}

// 团队移除成员回调
func OnRemoveTeamsMember(event webhook.TeamsMemberEvent) {
	// do your business here
}

// 团队添加频道回调
func OnCreateTeamsChannel(event webhook.TeamsChannelEvent) {
	// do your business here
}

// 团队删除频道回调
func OnDeleteTeamsChannel(event webhook.TeamsChannelEvent) {
	// do your business here
}

// 团队频道添加选项卡回调
func OnCreateTeamsChannelTab(event webhook.TeamsChannelTabEvent) {
	// do your business here
}

// 团队频道删除选项卡回调
func OnDeleteTeamsChannelTab(event webhook.TeamsChannelTabEvent) {
	// do your business here
}
