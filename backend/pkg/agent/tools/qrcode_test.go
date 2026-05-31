package tools

import (
	"strings"
	"testing"
)

func TestTerminalQRCodeRendersScannableBlock(t *testing.T) {
	out, err := terminalQRCode("https://example.com")
	if err != nil {
		t.Fatalf("terminalQRCode: %v", err)
	}
	if !strings.ContainsAny(out, "█▀▄") {
		t.Fatalf("expected terminal qr block, got %q", out)
	}
}
