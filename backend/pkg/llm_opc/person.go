package llm_opc

import (
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/llm_provider/agents"
)

// IOpcPerson 定义 OPC 人员的抽象接口
type IOpcPerson interface {
	Uuid() string
	Name() string
	Duties() string   // 职责（对应现有 Role 字段）
	Desc() string     // 人类可读描述
	Prompt() string   // agent 的系统提示词
	Tools() []string  // 工具 ID 列表
	Skills() []string // 技能 ID 列表
	Avatar() string
}

// OpcPersonAdapter 适配器，组合 OPCPerson + CustomAgentDef 实现 IOpcPerson
type OpcPersonAdapter struct {
	person *data_models.OPCPerson
	agent  *agents.CustomAgentDef
}

func NewOpcPerson(person *data_models.OPCPerson, agent *agents.CustomAgentDef) IOpcPerson {
	return &OpcPersonAdapter{person: person, agent: agent}
}

func (a *OpcPersonAdapter) Uuid() string     { return a.person.Uuid }
func (a *OpcPersonAdapter) Name() string     { return a.person.Name }
func (a *OpcPersonAdapter) Duties() string   { return a.person.Role }
func (a *OpcPersonAdapter) Desc() string     { return a.agent.Description }
func (a *OpcPersonAdapter) Prompt() string   { return a.agent.PromptText }
func (a *OpcPersonAdapter) Tools() []string  { return a.agent.ToolIDs }
func (a *OpcPersonAdapter) Skills() []string { return a.agent.SkillIDs }
func (a *OpcPersonAdapter) Avatar() string   { return a.person.Avatar }
