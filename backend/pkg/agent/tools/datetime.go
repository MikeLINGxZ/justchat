package tools

import (
	"context"
	"encoding/json"
	"time"

	"trpc.group/trpc-go/trpc-agent-go/tool/function"
)

type dateTimeInput struct {
	Format   string `json:"format" jsonschema:"description=Go time format string (default: RFC3339)"`
	Timezone string `json:"timezone" jsonschema:"description=IANA timezone name (default: Local)"`
}

type dateTimeOutput struct {
	DateTime string `json:"datetime"`
	Unix     int64  `json:"unix"`
	Timezone string `json:"timezone"`
}

func dateTimeFunc(ctx context.Context, input dateTimeInput) (dateTimeOutput, error) {
	loc := time.Local
	if input.Timezone != "" {
		parsed, err := time.LoadLocation(input.Timezone)
		if err != nil {
			return dateTimeOutput{}, err
		}
		loc = parsed
	}

	now := time.Now().In(loc)
	format := time.RFC3339
	if input.Format != "" {
		format = input.Format
	}

	return dateTimeOutput{
		DateTime: now.Format(format),
		Unix:     now.Unix(),
		Timezone: loc.String(),
	}, nil
}

func NewDateTimeTool() *function.FunctionTool[dateTimeInput, dateTimeOutput] {
	return function.NewFunctionTool(
		dateTimeFunc,
		function.WithName("datetime"),
		function.WithDescription("Get the current date and time in a specified format and timezone"),
	)
}

func DateTimeMeta() ToolMeta {
	return ToolMeta{
		Name:            "datetime",
		Description:     "Get the current date and time",
		Category:        CategoryBuiltin,
		RequiresConfirm: false,
		FormatPurpose: func(args json.RawMessage) string {
			return "Get current date and time"
		},
	}
}
