package client

import (
	"fmt"
	"strings"
)

type explainGroupMemberFailInfo struct {
	Members []string `json:"members"`
	Reason  string   `json:"reason"`
}

func (f explainGroupMemberFailInfo) Error() string {
	return fmt.Sprintf("%v: %v", strings.Join(f.Members, ","), f.Reason)
}

// 机器人建群，参数name群名称，owner群主域账号，members群成员域账号列表（包不包含群主都行），返回群id。
//
// Warning.Fails为拉群失败成员域账号列表（不包含群主）。
//
// 如果机器人无权拉群主进群，将直接返回error建群失败。
func (cli *Client) CreateGroup(name, owner string, members []string) (string, *Warning[string], error) {
	const api = "/group/create"
	args := object{
		"name":    name,
		"owner":   owner,
		"members": members,
	}
	var result struct {
		GroupId string `json:"group_id"`
		warning[string]
	}
	err := cli.call(api, args, &result)
	if err != nil {
		return "", nil, err
	}
	if len(result.Fails) == 0 {
		return result.GroupId, nil, nil
	}
	return result.GroupId, result.parse(explainGroupMemberFailInfo{}), nil
}

// 添加群成员。members为要添加的新成员域账号列表。如果群成员已经在群中，不会报错。
//
// Warning.Fails为添加失败成员域账号列表。
func (cli *Client) AddGroupMembers(groupId string, members []string) (*Warning[string], error) {
	const api = "/group/member/add"
	args := object{
		"group_id": groupId,
		"members":  members,
	}
	warn := new(warning[string])
	err := cli.call(api, args, warn)
	if err != nil {
		return nil, err
	}
	if len(warn.Fails) == 0 {
		return nil, nil
	}
	return warn.parse(explainGroupMemberFailInfo{}), nil
}

// 移除群成员。members为要移除的成员域账号列表。如果成员不在群中，不会报错。
//
// Warning.Fails为移除失败成员域账号列表。
func (cli *Client) RemoveGroupMembers(groupId string, members []string) (*Warning[string], error) {
	const api = "/group/member/remove"
	args := object{
		"group_id": groupId,
		"members":  members,
	}
	warn := new(warning[string])
	err := cli.call(api, args, warn)
	if err != nil {
		return nil, err
	}
	if len(warn.Fails) == 0 {
		return nil, nil
	}
	return warn.parse(explainGroupMemberFailInfo{}), nil
}

type GroupIdNamePair struct {
	GroupId string `json:"group_id"`
	Name    string `json:"name"` // 群名称
}

// 机器人所在群列表。
func (cli *Client) GroupsRobotIn() ([]GroupIdNamePair, error) {
	const api = "/group/robot/in"
	var result struct {
		Groups []GroupIdNamePair `json:"groups"`
	}
	err := cli.call(api, nil, &result)
	if err != nil {
		return nil, err
	}
	return result.Groups, nil
}

// 判断用户是否在群内。返回用户和机器人共同所在群列表。
func (cli *Client) IsUserInGroups(user string, groupIds []string) ([]GroupIdNamePair, error) {
	const api = "/group/user/isin"
	args := object{
		"user":   user,
		"groups": groupIds,
	}
	var result struct {
		Groups []GroupIdNamePair `json:"groups"`
	}
	err := cli.call(api, args, &result)
	if err != nil {
		return nil, err
	}
	return result.Groups, nil
}

type GroupMemberInfo struct {
	Uid     string   `json:"uid"`
	Account string   `json:"account,omitempty"` // 域账号，机器人没有域账号
	Name    string   `json:"name,omitempty"`    // 群成员姓名
	Dept    []string `json:"dept,omitempty"`    // 群成员所在部门，机器人不属于任何部门
}

// 获取所有群成员。返回结果包括机器人。
func (cli *Client) GetGroupMembers(groupId string) ([]GroupMemberInfo, error) {
	const api = "/group/members"
	args := object{
		"group_id": groupId,
	}
	var result struct {
		Members []GroupMemberInfo `json:"members"`
	}
	err := cli.call(api, args, &result)
	if err != nil {
		return nil, err
	}
	return result.Members, nil
}
