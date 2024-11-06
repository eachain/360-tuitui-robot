package webhook

import (
	"encoding/json"
	"net/http"
)

type Options struct {
	// Errorf用于输出错误日志，由调用方提供，可以为log.Printf。
	// 默认为nil，表示不输出日志。
	Errorf func(string, ...any)
}

// Callback是推推webhook回调函数列表。
//
// 只需要注册感兴趣的事件回调即可，没有注册的事件将被忽略。
//
// 如果事件由收回调的本机器人触发，不会有回调产生。
type Callback struct {
	// 单聊消息回调，对应事件single_chat。
	OnReceiveSingleMessage func(SingleMessageEvent)

	// 群聊消息回调，对应事件group_chat。
	OnReceiveGroupMessage func(GroupMessageEvent)
	// 建群事件，对应事件group_create。
	OnCreateGroup func(GroupMemberEvent)
	// 新成员进群回调，对应事件group_invite。
	OnNewMemberJoinGroup func(GroupMemberEvent)
	// 踢群成员回调，对应事件group_kick。
	OnGroupKickMember func(GroupMemberEvent)

	// 团队创建帖子回调，对应事件teams_post_create。
	OnCreateTeamsPost func(TeamsPostEvent)
	// 团队修改帖子回调，对应事件teams_post_modify。
	OnModifyTeamsPost func(TeamsPostEvent)
	// 团队添加成员回调，对应事件teams_member_add。
	OnAddTeamsMember func(TeamsMemberEvent)
	// 团队移除成员回调，对应事件teams_member_remove。
	OnRemoveTeamsMember func(TeamsMemberEvent)
	// 团队添加频道回调，对应事件teams_channel_create。
	OnCreateTeamsChannel func(TeamsChannelEvent)
	// 团队删除频道回调，对应事件teams_channel_delete。
	OnDeleteTeamsChannel func(TeamsChannelEvent)
	// 团队频道添加选项卡回调，对应事件teams_channel_tab_create。
	OnCreateTeamsChannelTab func(TeamsChannelTabEvent)
	// 团队频道删除选项卡回调，对应事件teams_channel_tab_create。
	OnDeleteTeamsChannelTab func(TeamsChannelTabEvent)
}

type eventRequest struct {
	// 事件发起人
	Cid  string `json:"cid"`
	Uid  string `json:"uid"`
	User string `json:"user_account"` // 域账号
	Name string `json:"user_name"`

	Timestamp int64 `json:"timestamp,string"` // 事件发起时间

	// 具体事件
	Event string `json:"event"`

	// 事件数据，事件不同，结构不同
	Data json.RawMessage `json:"data"`

	// 解析Data字段回调
	decode func(any) error
}

// raiser: 事件发起人
func (er *eventRequest) raiser() User {
	return User{
		// Cid:     er.Cid,
		Uid:     er.Uid,
		Account: er.User,
		Name:    er.Name,
	}
}

// NewHandler将推推webhook转为对应回调。
// 上层业务注册感兴趣的事件回调函数，没有注册的事件将被忽略。
func NewHandler(cb Callback, opts *Options) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := new(eventRequest)
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			if opts != nil && opts.Errorf != nil {
				opts.Errorf("webhook: json decode request body: %v", err)
			}
			return
		}

		req.decode = func(event any) error {
			err = json.Unmarshal(req.Data, event)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				if opts != nil && opts.Errorf != nil {
					opts.Errorf("webhook: json decode event %q data: %v", req.Event, err)
				}
			}
			return err
		}

		switch req.Event {
		case "single_chat":
			onReceiveSingleMessage(req, cb.OnReceiveSingleMessage)

		case "group_chat":
			onReceiveGroupMessage(req, cb.OnReceiveGroupMessage)
		case "group_create":
			onCreateGroup(req, cb.OnCreateGroup)
		case "group_invite":
			onNewMemberJoinGroup(req, cb.OnNewMemberJoinGroup)
		case "group_kick":
			onGroupKickMember(req, cb.OnGroupKickMember)

		case "teams_post_create":
			onCreateTeamsPost(req, cb.OnCreateTeamsPost)
		case "teams_post_modify":
			onModifyTeamsPost(req, cb.OnModifyTeamsPost)
		case "teams_member_add":
			onAddTeamsMember(req, cb.OnAddTeamsMember)
		case "teams_member_remove":
			onRemoveTeamsMember(req, cb.OnRemoveTeamsMember)
		case "teams_channel_create":
			onCreateTeamsChannel(req, cb.OnCreateTeamsChannel)
		case "teams_channel_delete":
			onDeleteTeamsChannel(req, cb.OnDeleteTeamsChannel)
		case "teams_channel_tab_create":
			onCreateTeamsChannelTab(req, cb.OnCreateTeamsChannelTab)
		case "teams_channel_tab_delete":
			onDeleteTeamsChannelTab(req, cb.OnDeleteTeamsChannelTab)
		}
	})
}

// some common structures

type User struct {
	// Cid     string `json:"cid,omitempty"`
	Uid     string `json:"uid,omitempty"`
	Account string `json:"user,omitempty"` // 域账号
	Name    string `json:"name,omitempty"`
}

type File struct {
	Name    string `json:"name"`     // 文件名
	MediaId string `json:"media_id"` // 对应发消息用到的media_id
	URL     string `json:"url"`      // 临时下载URL，有过期时间，防盗链
}

type Image struct {
	MediaId string `json:"media_id"` // 对应发消息用到的media_id
	URL     string `json:"url"`      // 临时下载URL，有过期时间，防盗链
}

type Voice struct {
	MediaId string `json:"media_id"` // 对应发消息用到的media_id
	URL     string `json:"url"`      // 临时下载URL，有过期时间，防盗链
}

type Message struct {
	MsgId     string   `json:"msgid"`               // 消息id
	MsgType   string   `json:"msgtype"`             // 消息类型：文本"text", 引用"reference", 图文混排"mixed", 文件"file", 图片"image", 语音"voice"
	Text      string   `json:"text,omitempty"`      // valid if MsgType == "text" || MsgType == "reference" || MsgType == "mixed"
	Reference string   `json:"reference,omitempty"` // valid if MsgType == "reference"
	File      *File    `json:"file,omitempty"`      // valid if MsgType == "file"
	Images    []*Image `json:"images,omitempty"`    // valid if MsgType == "image" || MsgType == "mixed"
	Voice     *Voice   `json:"voice,omitempty"`     // valid if MsgType == "voice"
}
