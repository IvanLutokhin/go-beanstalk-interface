package version

import (
	"strings"
	"testing"
)

func TestString(t *testing.T) {
	title = "Test App"
	tag = "0.0.0"
	commit = "000000"

	s := String()
	if !strings.EqualFold(s, "Test App 0.0.0 (Build: 000000)") {
		t.Errorf("unexpected version %q", s)
	}
}
