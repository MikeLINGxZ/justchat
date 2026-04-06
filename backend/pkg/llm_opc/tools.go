package llm_opc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
)

// ==================== GroupInfoTool ====================

type groupInfoTool struct {
	groupName string
	groupDesc string
	members   []IOpcPerson
}

func (t *groupInfoTool) Id() string   { return "opc_group_info" }
func (t *groupInfoTool) Name() string { return "获取群聊信息" }
func (t *groupInfoTool) Description() string {
	return "获取当前群聊的基本信息和成员列表"
}
func (t *groupInfoTool) RequireConfirmation() bool { return false }

type groupInfoResult struct {
	GroupName   string             `json:"group_name"`
	Description string             `json:"description"`
	Members     []groupMemberBrief `json:"members"`
}

type groupMemberBrief struct {
	Name   string `json:"name"`
	Duties string `json:"duties"`
	Desc   string `json:"desc"`
}

func (t *groupInfoTool) Tool() tool.BaseTool {
	return utils.NewTool(
		&schema.ToolInfo{
			Name:        "opc_group_info",
			Desc:        "获取当前群聊的基本信息和成员列表",
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{}),
		},
		func(ctx context.Context, _ emptyParams) (string, error) {
			var members []groupMemberBrief
			for _, m := range t.members {
				members = append(members, groupMemberBrief{
					Name:   m.Name(),
					Duties: m.Duties(),
					Desc:   m.Desc(),
				})
			}
			result := groupInfoResult{
				GroupName:   t.groupName,
				Description: t.groupDesc,
				Members:     members,
			}
			data, err := json.Marshal(result)
			if err != nil {
				return "", fmt.Errorf("failed to marshal group info: %w", err)
			}
			return string(data), nil
		},
	)
}

// ==================== PersonInfoTool ====================

type personInfoTool struct {
	members []IOpcPerson
}

func (t *personInfoTool) Id() string   { return "opc_person_info" }
func (t *personInfoTool) Name() string { return "获取人员信息" }
func (t *personInfoTool) Description() string {
	return "根据姓名获取群聊中某个人员的详细信息"
}
func (t *personInfoTool) RequireConfirmation() bool { return false }

type personInfoParams struct {
	Name string `json:"name"`
}

type personInfoResult struct {
	Name   string   `json:"name"`
	Duties string   `json:"duties"`
	Desc   string   `json:"desc"`
	Skills []string `json:"skills"`
}

func (t *personInfoTool) Tool() tool.BaseTool {
	return utils.NewTool(
		&schema.ToolInfo{
			Name: "opc_person_info",
			Desc: "根据姓名获取群聊中某个人员的详细信息",
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
				"name": {
					Type:     schema.String,
					Desc:     "要查询的人员姓名",
					Required: true,
				},
			}),
		},
		func(ctx context.Context, params personInfoParams) (string, error) {
			for _, m := range t.members {
				if m.Name() == params.Name {
					result := personInfoResult{
						Name:   m.Name(),
						Duties: m.Duties(),
						Desc:   m.Desc(),
						Skills: m.Skills(),
					}
					data, err := json.Marshal(result)
					if err != nil {
						return "", fmt.Errorf("failed to marshal person info: %w", err)
					}
					return string(data), nil
				}
			}
			return fmt.Sprintf("未找到名为 %q 的人员", params.Name), nil
		},
	)
}

// ==================== 构造函数 ====================

type emptyParams struct{}

// NewGroupChatTools 创建群聊会话级内置工具（GroupInfo + PersonInfo）
// 这些工具不注册到全局 ToolRouter，而是每次会话动态创建
func NewGroupChatTools(groupName, groupDesc string, members []IOpcPerson) []tool.BaseTool {
	groupInfo := &groupInfoTool{
		groupName: groupName,
		groupDesc: groupDesc,
		members:   members,
	}
	personInfo := &personInfoTool{
		members: members,
	}
	return []tool.BaseTool{
		groupInfo.Tool(),
		personInfo.Tool(),
	}
}
