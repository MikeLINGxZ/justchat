package tools

import (
	"context"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/i18n"
)

type CurrentDate struct {
}

func (c *CurrentDate) Id() string {
	return "get_current_date"
}

func (c *CurrentDate) Name() string {
	return i18n.TCurrent("tool.current_date.name", nil)
}

func (c *CurrentDate) Description() string {
	return i18n.TCurrent("tool.current_date.description", nil)
}

func (c *CurrentDate) RequireConfirmation() bool { return false }

func (c *CurrentDate) Tool() tool.BaseTool {
	t := utils.NewTool(
		&schema.ToolInfo{
			Name:        "get_current_date",
			Desc:        i18n.TCurrent("tool.current_date.description", nil),
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{}),
		},
		func(ctx context.Context, _ emptyParams) (string, error) {
			now := time.Now()
			if i18n.CurrentLocale() == i18n.LocaleEnUS {
				weekdays := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
				return now.Format("2006-01-02") + " " + weekdays[now.Weekday()], nil
			}
			weekdays := []string{"星期日", "星期一", "星期二", "星期三", "星期四", "星期五", "星期六"}
			return now.Format("2006-01-02") + " " + weekdays[now.Weekday()], nil
		},
	)
	return t
}

type CurrentTime struct {
}

func (c *CurrentTime) Id() string {
	return "get_current_time"
}

func (c *CurrentTime) Name() string {
	return i18n.TCurrent("tool.current_time.name", nil)
}

func (c *CurrentTime) Description() string {
	return i18n.TCurrent("tool.current_time.description", nil)
}

func (c *CurrentTime) RequireConfirmation() bool { return false }

func (c *CurrentTime) Tool() tool.BaseTool {
	t := utils.NewTool(
		&schema.ToolInfo{
			Name:        "get_current_time",
			Desc:        i18n.TCurrent("tool.current_time.description", nil),
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{}),
		},
		func(ctx context.Context, _ emptyParams) (string, error) {
			return time.Now().Format("15:04:05"), nil
		},
	)
	return t
}
