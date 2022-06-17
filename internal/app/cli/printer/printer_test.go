package printer_test

import (
	"bytes"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/cli/printer"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPrint(t *testing.T) {
	t.Run("json", func(t *testing.T) {
		b := new(bytes.Buffer)

		require.NoError(t, printer.Print("json", b, map[string]interface{}{"string": "test", "int": 1, "bool": false}))
		require.Equal(t, `{"bool":false,"int":1,"string":"test"}`, b.String())
	})

	t.Run("yaml", func(t *testing.T) {
		b := new(bytes.Buffer)

		require.NoError(t, printer.Print("yaml", b, map[string]interface{}{"string": "test", "int": 1, "bool": false}))
		require.Equal(t, "bool: false\nint: 1\nstring: test\n", b.String())
	})

	t.Run("default", func(t *testing.T) {
		b := new(bytes.Buffer)

		require.NoError(t, printer.Print("default", b, map[string]interface{}{"string": "test", "int": 1, "bool": false}))
		require.Equal(t, "map[bool:false int:1 string:test]", b.String())
	})
}
