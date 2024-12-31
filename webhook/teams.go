package webhook

// 团队标签，在团队帖子中可@标签
type TeamsTag struct {
	Id      string `json:"id,omitempty"`      // 标签id
	Name    string `json:"name,omitempty"`    // 标签名称
	Desc    string `json:"desc,omitempty"`    // 标签描述
	Members []User `json:"members,omitempty"` // 该标签包含哪些团队成员
}

type PostAtType string

const (
	AtAll  PostAtType = "all"
	AtTag  PostAtType = "tag"
	AtUser PostAtType = "user"
)

type TeamsPostAt struct {
	// @类型，有以下三种：
	// AtAll("all"): @所有人，Tag字段和User字段均为nil
	// AtTag("tag"): @标签，Tag字段不为空，User字段为nil
	// AtUser("user"): @用户，User字段不为空，Tag字段为nil
	Type PostAtType `json:"type"`
	Tag  *TeamsTag  `json:"tag,omitempty"`  // @团队标签
	User *User      `json:"user,omitempty"` // @用户
}

type TeamsPostEvent struct {
	User User `json:"user"` // post sender or modifier

	TeamId       string        `json:"team_id"`             // 团队id
	TeamName     string        `json:"team_name"`           // 团队名称
	TeamDesc     string        `json:"team_desc"`           // 团队描述
	ChannelId    string        `json:"channel_id"`          // 频道id
	ChannelName  string        `json:"channel_name"`        // 频道名称
	ChannelDesc  string        `json:"channel_desc"`        // 频道描述
	IsReply      bool          `json:"is_reply"`            // 是否回帖，如果为true，表示本帖是回复的ParentId的帖子；如果为false，表示本帖是主帖
	ParentId     string        `json:"parent_id,omitempty"` // 父帖，当且仅当IsReply=true时有效
	PostId       string        `json:"post_id"`             // 帖子id
	Content      string        `json:"content"`             // 帖子文本内容
	RichTextType string        `json:"rich_text_type"`      // 富文本类型，可选值有：json/v1, html/v1
	RichText     string        `json:"rich_text"`           // 富文本，不建议业务解析该字段
	At           []TeamsPostAt `json:"at,omitempty"`        // 帖子@列表
	AtMe         bool          `json:"at_me"`               // 是否@机器人，当且仅当At中有明确@机器人时为true。当@所有人和@标签时为false
	Files        []*File       `json:"files,omitempty"`     // 帖子内容中的文件
	Images       []*Image      `json:"images,omitempty"`    // 帖子内容中的图片
}

type teamsPost struct {
	TeamId       string        `json:"team_id"`
	TeamName     string        `json:"team_name"`
	TeamDesc     string        `json:"team_desc"`
	ChannelId    string        `json:"channel_id"`
	ChannelName  string        `json:"channel_name"`
	ChannelDesc  string        `json:"channel_desc"`
	IsReply      bool          `json:"is_reply"`
	ParentId     string        `json:"parent_id,omitempty"`
	PostId       string        `json:"post_id"`
	Content      string        `json:"content"`
	RichTextType string        `json:"rich_text_type"` // json/v1, html/v1
	RichText     string        `json:"rich_text"`
	At           []TeamsPostAt `json:"at,omitempty"`
	AtMe         bool          `json:"at_me"` // whether @robot
	Files        []*fileOutput `json:"files,omitempty"`
	Images       []string      `json:"images,omitempty"`
	ImageIds     []string      `json:"image_ids,omitempty"`
}

func onTeamsPostEvent(req *eventRequest, cb func(TeamsPostEvent)) {
	if cb == nil {
		return
	}

	var tp teamsPost
	err := req.decode(&tp)
	if err != nil {
		return
	}

	post := TeamsPostEvent{
		User:         req.raiser.toUser(),
		TeamId:       tp.TeamId,
		TeamName:     tp.TeamName,
		TeamDesc:     tp.TeamDesc,
		ChannelId:    tp.ChannelId,
		ChannelName:  tp.ChannelName,
		ChannelDesc:  tp.ChannelDesc,
		IsReply:      tp.IsReply,
		ParentId:     tp.ParentId,
		PostId:       tp.PostId,
		Content:      tp.Content,
		RichTextType: tp.RichTextType,
		RichText:     tp.RichText,
		At:           tp.At,
		AtMe:         tp.AtMe,
	}

	for _, file := range tp.Files {
		post.Files = append(post.Files, &File{
			Name:    file.Name,
			MediaId: file.Fid,
			URL:     file.URL,
		})
	}

	if len(tp.Images) > 0 && len(tp.ImageIds) == len(tp.Images) {
		for i := range tp.ImageIds {
			post.Images = append(post.Images, &Image{
				MediaId: tp.ImageIds[i],
				URL:     tp.Images[i],
			})
		}
	}

	go cb(post) // 避免阻塞推推业务
}

func onCreateTeamsPost(req *eventRequest, cb func(TeamsPostEvent)) {
	onTeamsPostEvent(req, cb)
}

func onModifyTeamsPost(req *eventRequest, cb func(TeamsPostEvent)) {
	onTeamsPostEvent(req, cb)
}

type TeamsMemberEvent struct {
	User     User   `json:"user"`      // 事件发起人，指是谁将团队成员拉入、移出的
	TeamId   string `json:"team_id"`   // 团队id
	TeamName string `json:"team_name"` // 团队名称
	TeamDesc string `json:"team_desc"` // 团队描述
	Members  []User `json:"members"`   // 本次事件被操作的团队成员
}

func onTeamsMemberEvent(req *eventRequest, cb func(TeamsMemberEvent)) {
	if cb == nil {
		return
	}

	var event TeamsMemberEvent
	err := req.decode(&event)
	if err != nil {
		return
	}
	event.User = req.raiser.toUser()

	go cb(event) // 避免阻塞推推业务
}

func onAddTeamsMember(req *eventRequest, cb func(TeamsMemberEvent)) {
	onTeamsMemberEvent(req, cb)
}

func onRemoveTeamsMember(req *eventRequest, cb func(TeamsMemberEvent)) {
	onTeamsMemberEvent(req, cb)
}

type TeamsChannelEvent struct {
	User User `json:"user"` // 事件发起人，指是谁创建/删除的团队频道

	TeamId      string `json:"team_id"`                // 团队id
	TeamName    string `json:"team_name"`              // 团队名称
	TeamDesc    string `json:"team_desc"`              // 团队描述
	ChannelId   string `json:"channel_id"`             // 频道id
	ChannelName string `json:"channel_name,omitempty"` // 频道名称
	ChannelDesc string `json:"channel_desc,omitempty"` // 频道描述
}

func onTeamsChannelEvent(req *eventRequest, cb func(TeamsChannelEvent)) {
	if cb == nil {
		return
	}

	var event TeamsChannelEvent
	err := req.decode(&event)
	if err != nil {
		return
	}
	event.User = req.raiser.toUser()

	go cb(event) // 避免阻塞推推业务
}

func onCreateTeamsChannel(req *eventRequest, cb func(TeamsChannelEvent)) {
	onTeamsChannelEvent(req, cb)
}

func onDeleteTeamsChannel(req *eventRequest, cb func(TeamsChannelEvent)) {
	onTeamsChannelEvent(req, cb)
}

// 选项卡设置
type TeamsChannelTabSettings struct {
	EntityId         string `json:"entity_id,omitempty"`     // 实体ID，用于从帖子跳转到选项卡
	Subtype          string `json:"subtype,omitempty"`       // "webpage"
	ContentURL       string `json:"content_url,omitempty"`   // PC端地址
	MobileContentURL string `json:"m_content_url,omitempty"` // 移动端地址，如果为空表示不支持移动端
	ConfigURL        string `json:"config_url,omitempty"`
}

type TeamsChannelTabEvent struct {
	User User `json:"user"` // 事件发起人，指是谁添加/删除的选项卡

	TeamId      string                   `json:"team_id"`                // 团队id
	TeamName    string                   `json:"team_name"`              // 团队名称
	TeamDesc    string                   `json:"team_desc"`              // 团队描述
	ChannelId   string                   `json:"channel_id"`             // 频道id
	ChannelName string                   `json:"channel_name"`           // 团队名称
	ChannelDesc string                   `json:"channel_desc"`           // 团队描述
	TabId       string                   `json:"tab_id"`                 // 选项卡id
	TabName     string                   `json:"tab_name,omitempty"`     // 选项卡名称
	TabOrder    string                   `json:"tab_order,omitempty"`    // 选项卡排序
	TabSettings *TeamsChannelTabSettings `json:"tab_settings,omitempty"` // 选项卡设置
}

func onTeamsChannelTabEvent(req *eventRequest, cb func(TeamsChannelTabEvent)) {
	if cb == nil {
		return
	}

	var event TeamsChannelTabEvent
	err := req.decode(&event)
	if err != nil {
		return
	}
	event.User = req.raiser.toUser()

	go cb(event) // 避免阻塞推推业务
}

func onCreateTeamsChannelTab(req *eventRequest, cb func(TeamsChannelTabEvent)) {
	onTeamsChannelTabEvent(req, cb)
}

func onDeleteTeamsChannelTab(req *eventRequest, cb func(TeamsChannelTabEvent)) {
	onTeamsChannelTabEvent(req, cb)
}
