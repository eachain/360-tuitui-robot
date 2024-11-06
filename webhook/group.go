package webhook

type GroupAtUser struct {
	IsAtAll bool `json:"is_at_all,omitempty"` // if true, User field is empty
	User         // valid only when IsAtAll == false
}

type GroupMessageEvent struct {
	User      User  `json:"user"`      // 消息发送者
	Timestamp int64 `json:"timestamp"` // 秒级时间戳

	GroupId   string        `json:"group_id"`     // 群id
	GroupName string        `json:"group_name"`   // 群名称
	At        []GroupAtUser `json:"at,omitempty"` // @用户列表
	AtMe      bool          `json:"at_me"`        // 当且仅当明确@机器人时，AtMe=true。@所有人但没有@机器人时，AtMe=false

	Message
}

type groupMessage struct {
	GroupId   string        `json:"group_id"`
	GroupName string        `json:"group_name"`
	At        []GroupAtUser `json:"at"`
	AtMe      bool          `json:"at_me"` // 群消息中是否有 @robot
	MsgId     string        `json:"msgid"`
	MsgType   string        `json:"msg_type"` // 消息类型：text(文本), reference(引用), file(文件), image(图片), mixed(图文混排)
	Text      string        `json:"text"`
	Reference string        `json:"reference"` // 被引用的消息
	File      *fileOutput   `json:"file"`
	Images    []string      `json:"images"`
	ImageIds  []string      `json:"image_ids"`
	Voice     string        `json:"voice"`
	VoiceFid  string        `json:"voice_id"`
}

func onReceiveGroupMessage(req *eventRequest, cb func(GroupMessageEvent)) {
	if cb == nil {
		return
	}

	var gm groupMessage
	err := req.decode(&gm)
	if err != nil {
		return
	}

	var msg GroupMessageEvent
	msg.User = req.raiser()
	msg.Timestamp = req.Timestamp

	msg.GroupId = gm.GroupId
	msg.GroupName = gm.GroupName
	msg.At = gm.At
	msg.AtMe = gm.AtMe
	msg.MsgId = gm.MsgId
	msg.MsgType = gm.MsgType

	msg.Text = gm.Text
	msg.Reference = gm.Reference

	if gm.File != nil {
		msg.File = &File{
			Name:    gm.File.Name,
			MediaId: gm.File.Fid,
			URL:     gm.File.URL,
		}
	}

	if len(gm.Images) > 0 && len(gm.ImageIds) == len(gm.Images) {
		for i := range gm.ImageIds {
			msg.Images = append(msg.Images, &Image{
				MediaId: gm.ImageIds[i],
				URL:     gm.Images[i],
			})
		}
	}

	if gm.VoiceFid != "" {
		msg.Voice = &Voice{
			MediaId: gm.VoiceFid,
			URL:     gm.Voice,
		}
	}

	go cb(msg) // 避免阻塞推推业务
}

type GroupMemberEvent struct {
	User       User   `json:"user"`        // 事件发起人，指是谁将群成员拉入、踢出的
	GroupId    string `json:"group_id"`    // 群id
	GroupName  string `json:"group_name"`  // 群名称
	ContainsMe bool   `json:"contains_me"` // 表示本次群成员变动中，是否包含本机器人
	Members    []User `json:"members"`     // 被拉入/踢出的群成员列表
}

type groupMemberEvent struct {
	GroupId    string `json:"group_id"`
	GroupName  string `json:"group_name"`
	ContainsMe bool   `json:"members_contains_me"`
	Members    []User `json:"members"`
}

func onGroupEvent(req *eventRequest, cb func(GroupMemberEvent)) {
	if cb == nil {
		return
	}

	var gme groupMemberEvent
	err := req.decode(&gme)
	if err != nil {
		return
	}

	go cb(GroupMemberEvent{
		User:       req.raiser(),
		GroupId:    gme.GroupId,
		GroupName:  gme.GroupName,
		ContainsMe: gme.ContainsMe,
		Members:    gme.Members,
	}) // 避免阻塞推推业务
}

func onCreateGroup(req *eventRequest, cb func(GroupMemberEvent)) {
	onGroupEvent(req, cb)
}

func onNewMemberJoinGroup(req *eventRequest, cb func(GroupMemberEvent)) {
	onGroupEvent(req, cb)
}

func onGroupKickMember(req *eventRequest, cb func(GroupMemberEvent)) {
	onGroupEvent(req, cb)
}
