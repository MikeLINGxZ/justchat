package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	qrcode "github.com/skip2/go-qrcode"
	"trpc.group/trpc-go/trpc-agent-go/tool/function"
)

const QRCodeToolName = "QRCode"

type qrCodeInput struct {
	Text string `json:"text" jsonschema:"description=text or URL to encode as a QR code,required"`
}

type qrCodeOutput struct {
	Text                string `json:"text"`
	QR                  string `json:"qr"`
	InteractiveTerminal bool   `json:"interactive_terminal"`
	TerminalStatus      string `json:"terminal_status"`
	TerminalOutput      string `json:"terminal_output"`
}

func BuildQRCodeTool() ToolMeta {
	return ToolMeta{
		Name:        QRCodeToolName,
		Description: "Convert text or a URL into a scannable terminal QR code. Use this when the user needs to scan a code or open a URL from another device.",
		Category:    CategoryBuiltin,
		FormatPurpose: func(args json.RawMessage) string {
			var input qrCodeInput
			_ = json.Unmarshal(args, &input)
			if strings.TrimSpace(input.Text) == "" {
				return "Generate QR code"
			}
			return "Generate QR code for text"
		},
	}
}

func NewQRCodeTool() *function.FunctionTool[qrCodeInput, qrCodeOutput] {
	meta := BuildQRCodeTool()
	return function.NewFunctionTool(
		func(ctx context.Context, input qrCodeInput) (qrCodeOutput, error) {
			text := strings.TrimSpace(input.Text)
			if text == "" {
				return qrCodeOutput{}, errors.New("text is required")
			}
			qr, err := terminalQRCode(text)
			if err != nil {
				return qrCodeOutput{}, err
			}
			return qrCodeOutput{
				Text:                text,
				QR:                  qr,
				InteractiveTerminal: true,
				TerminalStatus:      "active",
				TerminalOutput:      qr,
			}, nil
		},
		function.WithName(QRCodeToolName),
		function.WithDescription(meta.Description),
	)
}

func terminalQRCode(text string) (string, error) {
	qr, err := qrcode.New(text, qrcode.Medium)
	if err != nil {
		return "", fmt.Errorf("encode qr: %w", err)
	}
	bitmap := qr.Bitmap()
	if len(bitmap) == 0 {
		return "", nil
	}

	const quiet = 2
	size := len(bitmap)
	var out strings.Builder
	for y := -quiet; y < size+quiet; y += 2 {
		for x := -quiet; x < size+quiet; x++ {
			top := qrCell(bitmap, x, y)
			bottom := qrCell(bitmap, x, y+1)
			switch {
			case top && bottom:
				out.WriteString("█")
			case top:
				out.WriteString("▀")
			case bottom:
				out.WriteString("▄")
			default:
				out.WriteString(" ")
			}
		}
		out.WriteByte('\n')
	}
	return out.String(), nil
}

func qrCell(bitmap [][]bool, x, y int) bool {
	if y < 0 || y >= len(bitmap) || x < 0 || x >= len(bitmap[y]) {
		return false
	}
	return bitmap[y][x]
}
