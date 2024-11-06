package client

// 机器人属性
type RobotProperties struct {
	Name           string `json:"name"`            // 名称
	Avatar         string `json:"avatar"`          // 头像
	AvatarPreview  string `json:"avatar_preview"`  // 临时下载链接，防盗链，带过期时间
	Webhook        string `json:"webhook"`         // webhook回调地址
	InteractiveURL string `json:"interactive_url"` // 可交互式消息交互回调地址
}

// 获取机器人属性。
func (cli *Client) GetRobotProps() (*RobotProperties, error) {
	const api = "/robot/prop/get"
	var result struct {
		Props RobotProperties `json:"properties"`
	}
	err := cli.call(api, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result.Props, nil
}

// 修改机器人名称。
func (cli *Client) ModifyRobotName(name string) error {
	const api = "/robot/name/modify"
	args := object{"name": name}
	return cli.call(api, args, nil)
}

// 修改机器人头像。
func (cli *Client) ModifyRobotAvatar(avatar string) error {
	const api = "/robot/avatar/modify"
	args := object{"avatar": avatar}
	return cli.call(api, args, nil)
}

// 修改机器人收消息回调地址。
func (cli *Client) ModifyRobotWebhook(webhook string) error {
	const api = "/robot/webhook/modify"
	args := object{"url": webhook}
	return cli.call(api, args, nil)
}

// 修改机器人可交互式消息，用户交互回调地址。
func (cli *Client) ModifyRobotInteractiveURL(url string) error {
	const api = "/robot/interactive_url/modify"
	args := object{"url": url}
	return cli.call(api, args, nil)
}
