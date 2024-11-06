package client

import (
	"encoding/json"
	"fmt"
	"strings"
)

// 消息格式定义，发消息和修改消息时，用来组装请求参数。
//
// 具体实现参考github.com/eachain/360-tuitui-robot/message。
type Message interface {
	Type() string  // 消息类型，用于填充发消息时"msgtype"参数
	Index() string // 消息内容在发消息时所用的键。如文件消息用"text"键，值用Message本身
}

type object = map[string]any

type UserMsgIdPair struct {
	User  string `json:"user"`
	MsgId string `json:"msgid"`
}

type explainSendUserMsgsFailInfo struct {
	ToUsers []string `json:"tousers"`
	Reason  string   `json:"reason"`
}

func (f explainSendUserMsgsFailInfo) Error() string {
	return fmt.Sprintf("%v: %v", strings.Join(f.ToUsers, ","), f.Reason)
}

// 批量发送单聊消息。返回发送成功的消息id列表。Warning.Fails为发送失败域账号列表。
func (cli *Client) SendMessageToUsers(users []string, msg Message) ([]UserMsgIdPair, *Warning[string], error) {
	m := object{
		"tousers":   users,
		"msgtype":   msg.Type(),
		msg.Index(): msg,
	}
	return send[UserMsgIdPair, string, explainSendUserMsgsFailInfo](cli, m)
}

// 发送单聊消息。返回发送成功的消息id。
func (cli *Client) SendMessageToUser(user string, msg Message) (string, error) {
	pairs, warn, err := cli.SendMessageToUsers([]string{user}, msg)
	if err != nil {
		return "", err
	}
	if len(pairs) == 0 {
		if warn != nil {
			return "", warn.Explains
		}
		return "", fmt.Errorf("send message to user %v failed", user)
	}
	return pairs[0].MsgId, nil
}

type ModifyOptions struct {
	WithoutPush bool // 是否需要推送厂商推送，可用于修改消息时避免手机响铃
}

type explainModifyUserMsgsFailInfo struct {
	ToUsers []UserMsgIdPair `json:"tousers"`
	Reason  string          `json:"reason"`
}

func (f explainModifyUserMsgsFailInfo) Error() string {
	msgs := make([]string, len(f.ToUsers))
	for i := range msgs {
		msgs[i] = fmt.Sprintf("%v.%v", f.ToUsers[i].User, f.ToUsers[i].MsgId)
	}
	return fmt.Sprintf("%v: %v", strings.Join(msgs, ","), f.Reason)
}

// 修改单聊消息。返回修改成功的消息id列表。Warning.Fails为修改失败消息列表。
func (cli *Client) ModifyUserMessages(msgids []UserMsgIdPair, msg Message, opt *ModifyOptions) ([]UserMsgIdPair, *Warning[UserMsgIdPair], error) {
	m := object{
		"tousers":   msgids,
		"msgtype":   msg.Type(),
		msg.Index(): msg,
	}
	if opt != nil && opt.WithoutPush {
		m["without_push"] = true
	}

	return modify[UserMsgIdPair, explainModifyUserMsgsFailInfo](cli, m)
}

// 修改单聊消息。
func (cli *Client) ModifyUserMessage(msgid UserMsgIdPair, msg Message, opt *ModifyOptions) error {
	oks, warn, err := cli.ModifyUserMessages([]UserMsgIdPair{msgid}, msg, opt)
	if err != nil {
		return err
	}
	if len(oks) == 0 {
		if warn != nil {
			return warn.Explains
		}
		return fmt.Errorf("modify user message failed: %v.%v", msgid.User, msgid.MsgId)
	}
	return nil
}

type GroupMsgIdPair struct {
	Group string `json:"group"`
	MsgId string `json:"msgid"`
}

type explainSendGroupMsgsFailInfo struct {
	ToGroups []string `json:"togroups"`
	Reason   string   `json:"reason"`
}

func (f explainSendGroupMsgsFailInfo) Error() string {
	return fmt.Sprintf("%v: %v", strings.Join(f.ToGroups, ","), f.Reason)
}

// 批量发送群聊消息。
//
// atUsers为该群消息需要@的域账号列表如["zhangsan"]（不带@符号）。如果需要@所有人，传["@all"]（带@符号）。
//
// 仅文本(text)、图片(image)和图文混排(mixed)消息支持@，其它消息类型不支持。
//
// 返回发送成功的消息id列表。Warning.Fails为发送失败群id列表。
//
// FAQ：群ID怎么拿到？
//
// 方案1：移动端APP打开群设置中二维码页面，保存截图再拿微信等工具扫一下，能看到里面的json格式的groupId字段。
//
// 方案2：机器人在群里，机器人收一下消息就知道了。机器人收消息见ModifyRobotWebhook方法和github.com/eachain/360-tuitui-robot/webhook。
func (cli *Client) SendMessageToGroups(groupIds, atUsers []string, msg Message) ([]GroupMsgIdPair, *Warning[string], error) {
	m := object{
		"togroups":  groupIds,
		"msgtype":   msg.Type(),
		msg.Index(): msg,
	}
	if len(atUsers) > 0 {
		m["at"] = atUsers
	}
	return send[GroupMsgIdPair, string, explainSendGroupMsgsFailInfo](cli, m)
}

// 发送群聊消息。
//
// atUsers为该群消息需要@的域账号列表如["zhangsan"]（不带@符号）。如果需要@所有人，传["@all"]（带@符号）。
//
// 仅文本(text)、图片(image)和图文混排(mixed)消息支持@，其它消息类型不支持。
//
// 返回消息id。
//
// FAQ：群ID怎么拿到？
//
// 方案1：移动端APP打开群设置中二维码页面，保存截图再拿微信等工具扫一下，能看到里面的json格式的groupId字段。
//
// 方案2：机器人在群里，机器人收一下消息就知道了。机器人收消息见ModifyRobotWebhook方法和github.com/eachain/360-tuitui-robot/webhook。
func (cli *Client) SendMessageToGroupAt(groupId string, atUsers []string, msg Message) (string, error) {
	pairs, warn, err := cli.SendMessageToGroups([]string{groupId}, atUsers, msg)
	if err != nil {
		return "", err
	}
	if len(pairs) == 0 {
		if warn != nil {
			return "", warn.Explains
		}
		return "", fmt.Errorf("send message to group %v failed", groupId)
	}
	return pairs[0].MsgId, nil
}

// 发送群聊消息。返回消息id。
//
// FAQ：群ID怎么拿到？
//
// 方案1：移动端APP打开群设置中二维码页面，保存截图再拿微信等工具扫一下，能看到里面的json格式的groupId字段。
//
// 方案2：机器人在群里，机器人收一下消息就知道了。机器人收消息见ModifyRobotWebhook方法和github.com/eachain/360-tuitui-robot/webhook。
func (cli *Client) SendMessageToGroup(groupId string, msg Message) (string, error) {
	return cli.SendMessageToGroupAt(groupId, nil, msg)
}

type explainModifyGroupMsgsFailInfo struct {
	ToGroups []GroupMsgIdPair `json:"togroups"`
	Reason   string           `json:"reason"`
}

func (f explainModifyGroupMsgsFailInfo) Error() string {
	msgs := make([]string, len(f.ToGroups))
	for i := range msgs {
		msgs[i] = fmt.Sprintf("%v.%v", f.ToGroups[i].Group, f.ToGroups[i].MsgId)
	}
	return fmt.Sprintf("%v: %v", strings.Join(msgs, ","), f.Reason)
}

// 批量修改群聊消息。
//
// atUsers为该群消息需要@的域账号列表如["zhangsan"]（不带@符号）。如果需要@所有人，传["@all"]（带@符号）。
//
// 仅文本(text)、图片(image)和图文混排(mixed)消息支持@，其它消息类型不支持。
//
// 返回修改成功的消息列表。Warning.Fails为修改失败消息列表。
func (cli *Client) ModifyGroupMessages(msgids []GroupMsgIdPair, atUsers []string, msg Message, opt *ModifyOptions) ([]GroupMsgIdPair, *Warning[GroupMsgIdPair], error) {
	m := object{
		"togroups":  msgids,
		"msgtype":   msg.Type(),
		msg.Index(): msg,
	}
	if opt != nil && opt.WithoutPush {
		m["without_push"] = true
	}
	return modify[GroupMsgIdPair, explainModifyGroupMsgsFailInfo](cli, m)
}

// 修改群聊消息。
//
// atUsers为该群消息需要@的域账号列表如["zhangsan"]（不带@符号）。如果需要@所有人，传["@all"]（带@符号）。
//
// 仅文本(text)、图片(image)和图文混排(mixed)消息支持@，其它消息类型不支持。
func (cli *Client) ModifyGroupMessageAt(msgid GroupMsgIdPair, atUsers []string, msg Message, opt *ModifyOptions) error {
	oks, warn, err := cli.ModifyGroupMessages([]GroupMsgIdPair{msgid}, atUsers, msg, opt)
	if err != nil {
		return err
	}
	if len(oks) == 0 {
		if warn != nil {
			return warn.Explains
		}
		return fmt.Errorf("modify group %v message %v failed", msgid.Group, msgid.MsgId)
	}
	return nil
}

// 修改群聊消息。
func (cli *Client) ModifyGroupMessage(msgid GroupMsgIdPair, msg Message, opt *ModifyOptions) error {
	return cli.ModifyGroupMessageAt(msgid, nil, msg, opt)
}

// 发送推推页面消息。返回page_id，可用于后续修改等操作。
//
// 文档：https://easydoc.qihoo.net/doc?project=1d414e4d0ce730bec9b805b12ca28509&doc=3596913a227ae858de8e1dcca7dae3d6&config=toc#h2-%E6%8E%A8%E6%8E%A8%E9%A1%B5%E9%9D%A2%E6%B6%88%E6%81%AF。
func (cli *Client) SendPageToUsers(users []string, msg Message) (string, []UserMsgIdPair, *Warning[string], error) {
	m := object{
		"tousers":   users,
		"msgtype":   msg.Type(),
		msg.Index(): msg,
	}
	var result struct {
		MsgIds []UserMsgIdPair `json:"msgids,omitempty"`
		PageId string          `json:"page_id,omitempty"` // 当且仅当msgtype=="page"时返回该值
		warning[string]
	}
	const api = "/message/custom/send"
	err := cli.call(api, m, &result)
	if err != nil {
		return "", nil, nil, err
	}
	if len(result.Fails) == 0 {
		return result.PageId, result.MsgIds, nil, nil
	}
	return result.PageId, result.MsgIds, result.parse(explainSendUserMsgsFailInfo{}), nil
}

// 发送推推页面消息。返回page_id，可用于后续修改等操作。
//
// 文档：https://easydoc.qihoo.net/doc?project=1d414e4d0ce730bec9b805b12ca28509&doc=3596913a227ae858de8e1dcca7dae3d6&config=toc#h2-%E6%8E%A8%E6%8E%A8%E9%A1%B5%E9%9D%A2%E6%B6%88%E6%81%AF。
func (cli *Client) SendPageToGroups(groupIds []string, msg Message) (string, []GroupMsgIdPair, *Warning[string], error) {
	m := object{
		"togroups":  groupIds,
		"msgtype":   msg.Type(),
		msg.Index(): msg,
	}
	var result struct {
		MsgIds []GroupMsgIdPair `json:"msgids,omitempty"`
		PageId string           `json:"page_id,omitempty"` // 当且仅当msgtype=="page"时返回该值
		warning[string]
	}
	const api = "/message/custom/send"
	err := cli.call(api, m, &result)
	if err != nil {
		return "", nil, nil, err
	}
	if len(result.Fails) == 0 {
		return result.PageId, result.MsgIds, nil, nil
	}
	return result.PageId, result.MsgIds, result.parse(explainSendGroupMsgsFailInfo{}), nil
}

// 修改推推页面消息。
//
// 文档：https://easydoc.qihoo.net/doc?project=1d414e4d0ce730bec9b805b12ca28509&doc=3596913a227ae858de8e1dcca7dae3d6&config=toc#h2-%E4%BF%AE%E6%94%B9%E6%8E%A8%E6%8E%A8%E9%A1%B5%E9%9D%A2%E6%B6%88%E6%81%AF。
func (cli *Client) ModifyPageContent(msg Message) error {
	m := object{
		"msgtype":   msg.Type(),
		msg.Index(): msg,
	}
	_, _, err := modify[string, explainModifyUserMsgsFailInfo](cli, m)
	return err
}

type TeamChannel struct {
	TeamId    string   `json:"team_id"`
	ChannelId string   `json:"channel_id"`
	ParentId  string   `json:"parent_id,omitempty"` // 回贴
	Tags      []string `json:"tags,omitempty"`
}

type TeamPost struct {
	TeamId    string `json:"team_id"`
	ChannelId string `json:"channel_id"`
	ParentId  string `json:"parent_id,omitempty"` // 回帖
	PostId    string `json:"post_id"`             // 帖子id
}

type explainSendPostsFailInfo struct {
	ToTeams []TeamChannel `json:"toteams"`
	Reason  string        `json:"reason"`
}

func (f explainSendPostsFailInfo) Error() string {
	teams := make([]string, len(f.ToTeams))
	for i := range teams {
		teams[i] = fmt.Sprintf("team %v channel %v",
			f.ToTeams[i].TeamId,
			f.ToTeams[i].ChannelId,
		)
		if f.ToTeams[i].ParentId != "" {
			teams[i] += fmt.Sprintf(" reply %v", f.ToTeams[i].ParentId)
		}
		if len(f.ToTeams[i].Tags) > 0 {
			teams[i] += " " + "with tags " + strings.Join(f.ToTeams[i].Tags, ",")
		}
	}
	return fmt.Sprintf("%v: %v", strings.Join(teams, ", "), f.Reason)
}

// 批量发送帖子到团队。
//
// 返回发送成功的帖子id。Warning.Fails为发送失败的团队列表。
//
// 文档：https://easydoc.qihoo.net/doc?project=1d414e4d0ce730bec9b805b12ca28509&doc=3596913a227ae858de8e1dcca7dae3d6&config=toc#h2-%E5%9B%A2%E9%98%9F%E5%B8%96%E5%AD%90(HTML)。
func (cli *Client) SendPostToTeams(teams []TeamChannel, msg Message) ([]TeamPost, *Warning[TeamChannel], error) {
	m := object{
		"toteams":   teams,
		"msgtype":   msg.Type(),
		msg.Index(): msg,
	}
	return send[TeamPost, TeamChannel, explainSendPostsFailInfo](cli, m)
}

// 发送帖子到团队。返回发送成功的帖子id。
//
// 文档：https://easydoc.qihoo.net/doc?project=1d414e4d0ce730bec9b805b12ca28509&doc=3596913a227ae858de8e1dcca7dae3d6&config=toc#h2-%E5%9B%A2%E9%98%9F%E5%B8%96%E5%AD%90(HTML)。
func (cli *Client) SendPostToTeam(team TeamChannel, msg Message) (string, error) {
	posts, warn, err := cli.SendPostToTeams([]TeamChannel{team}, msg)
	if err != nil {
		return "", err
	}
	if len(posts) == 0 {
		if warn != nil {
			return "", warn.Explains
		}
		tmp, _ := json.Marshal(team)
		return "", fmt.Errorf("send post to team failed: %s", tmp)
	}
	return posts[0].PostId, nil
}

type ModifyTeamPostRequest struct {
	TeamId    string   `json:"team_id"`
	ChannelId string   `json:"channel_id"`
	PostId    string   `json:"post_id"`
	Tags      []string `json:"tags,omitempty"`
}

type explainModifyPostsFailInfo struct {
	ToTeams []ModifyTeamPostRequest `json:"toteams"`
	Reason  string                  `json:"reason"`
}

func (f explainModifyPostsFailInfo) Error() string {
	teams := make([]string, len(f.ToTeams))
	for i := range teams {
		teams[i] = fmt.Sprintf("team %v channel %v post %v",
			f.ToTeams[i].TeamId,
			f.ToTeams[i].ChannelId,
			f.ToTeams[i].PostId,
		)
		if len(f.ToTeams[i].Tags) > 0 {
			teams[i] += " " + "with tags " + strings.Join(f.ToTeams[i].Tags, ",")
		}
	}
	return fmt.Sprintf("%v: %v", strings.Join(teams, ","), f.Reason)
}

// 批量修改团队帖子。
//
// 返回修改成功的帖子列表。Warning.Fails为修改失败帖子列表。
//
// 文档：https://easydoc.qihoo.net/doc?project=1d414e4d0ce730bec9b805b12ca28509&doc=3596913a227ae858de8e1dcca7dae3d6&config=toc#h2-%E4%BF%AE%E6%94%B9%E5%9B%A2%E9%98%9F%E5%B8%96%E5%AD%90。
func (cli *Client) ModifyTeamPosts(posts []ModifyTeamPostRequest, msg Message) ([]TeamPost, *Warning[ModifyTeamPostRequest], error) {
	m := object{
		"toteams":   posts,
		"msgtype":   msg.Type(),
		msg.Index(): msg,
	}
	var result struct {
		Success []struct {
			TeamPost `json:"post_ids"`
		} `json:"success,omitempty"`
		warning[ModifyTeamPostRequest]
	}

	const api = "/message/custom/modify"
	err := cli.call(api, m, &result)
	if err != nil {
		return nil, nil, err
	}
	oks := make([]TeamPost, len(result.Success))
	for i := range oks {
		oks[i] = result.Success[i].TeamPost
	}
	if len(result.Fails) == 0 {
		return oks, nil, nil
	}
	return oks, result.parse(explainModifyPostsFailInfo{}), nil
}

// 批量修改团队帖子。
//
// 文档：https://easydoc.qihoo.net/doc?project=1d414e4d0ce730bec9b805b12ca28509&doc=3596913a227ae858de8e1dcca7dae3d6&config=toc#h2-%E4%BF%AE%E6%94%B9%E5%9B%A2%E9%98%9F%E5%B8%96%E5%AD%90。
func (cli *Client) ModifyTeamPost(post ModifyTeamPostRequest, msg Message) error {
	oks, warn, err := cli.ModifyTeamPosts([]ModifyTeamPostRequest{post}, msg)
	if err != nil {
		return err
	}
	if len(oks) == 0 {
		if warn != nil {
			return warn.Explains
		}
		tmp, _ := json.Marshal(post)
		return fmt.Errorf("modify team post failed: %s", tmp)
	}
	return nil
}

type UserVoiceResult struct {
	Mobile  string `json:"mobile"`
	Success bool   `json:"success"`
	CallId  string `json:"call_id"`
	Error   string `json:"error,omitempty"`
}

// 电话报警。
//
// 文档：https://easydoc.qihoo.net/doc?project=1d414e4d0ce730bec9b805b12ca28509&doc=3596913a227ae858de8e1dcca7dae3d6&config=toc#h2-%E7%94%B5%E8%AF%9D%E6%8A%A5%E8%AD%A6。
func (cli *Client) SendVoiceToUsers(accounts []string, msg Message) ([]UserVoiceResult, error) {
	m := object{
		"tousers":   accounts,
		"msgtype":   msg.Type(),
		msg.Index(): msg,
	}
	var result struct {
		Voice []UserVoiceResult `json:"voice,omitempty"`
	}
	const api = "/message/custom/send"
	err := cli.call(api, m, &result)
	if err != nil {
		return nil, err
	}
	return result.Voice, nil
}

type VoiceDetail struct {
	CallId           string `json:"call_id"`
	Callee           string `json:"callee"`             // 被叫号码
	CalleeShowNumber string `json:"callee_show_number"` // 被叫显号
	State            string `json:"state"`              // 通话状态
	StateDesc        string `json:"state_desc"`         // 通话状态描述
	GmtCreate        string `json:"gmt_create"`         // 通话请求的接收时间
	StartDate        string `json:"start_date"`         // 呼叫开始时间
	EndDate          string `json:"end_date"`           // 呼叫结束时间
	Duration         int64  `json:"duration"`           // 通话时长，单位：秒

	APIDoc   string `json:"api_doc"`   // 语音报警api文档
	StateDoc string `json:"state_doc"` // 通话状态文档
}

// 查询电话报警接听状态。
func (cli *Client) QueryVoiceDetail(callId string) (*VoiceDetail, error) {
	const api = "/message/voice/detail"
	args := object{
		"call_id": callId,
	}
	var result struct {
		Detail   VoiceDetail `json:"detail"`
		APIDoc   string      `json:"api_doc"`
		StateDoc string      `json:"state_doc"`
	}
	err := cli.call(api, args, &result)
	if err != nil {
		return nil, err
	}
	result.Detail.APIDoc = result.APIDoc
	result.Detail.StateDoc = result.StateDoc
	return &result.Detail, nil
}

// 批量发送单聊消息。返回发送成功的消息id列表。Warning.Fails为发送失败域账号列表。
func send[R, F any, E error](cli *Client, req any) ([]R, *Warning[F], error) {
	var result struct {
		MsgIds []R `json:"msgids,omitempty"`
		warning[F]
	}
	const api = "/message/custom/send"
	err := cli.call(api, req, &result)
	if err != nil {
		return nil, nil, err
	}
	if len(result.Fails) == 0 {
		return result.MsgIds, nil, nil
	}
	var errType E
	return result.MsgIds, result.parse(errType), nil
}

func modify[R any, E error](cli *Client, req any) ([]R, *Warning[R], error) {
	var result struct {
		Success []R `json:"success,omitempty"`
		warning[R]
	}
	const api = "/message/custom/modify"
	err := cli.call(api, req, &result)
	if err != nil {
		return nil, nil, err
	}
	if len(result.Fails) == 0 {
		return result.Success, nil, nil
	}
	var errType E
	return result.Success, result.parse(errType), nil
}
