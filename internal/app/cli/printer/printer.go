package printer

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
)

func Print(format string, w io.Writer, v interface{}) error {
	switch format {
	case "json":
		return Json(w, v)

	case "yaml":
		return Yaml(w, v)

	default:
		return Default(w, v)
	}
}

func Json(w io.Writer, v interface{}) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "%s", bytes)

	return err
}

func Yaml(w io.Writer, v interface{}) error {
	bytes, err := yaml.Marshal(v)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "%s", bytes)

	return nil
}

func Default(w io.Writer, v interface{}) error {
	_, err := fmt.Fprintf(w, "%v", v)

	return err
}
