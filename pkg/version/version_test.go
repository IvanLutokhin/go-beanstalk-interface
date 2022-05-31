package version

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestString(t *testing.T) {
	title = "Test App"
	tag = "0.0.0"
	commit = "000000"

	require.Equal(t, "Test App 0.0.0 (Build: 000000)", String())
}
