package handler

import "fmt"

func InvalidArgumentError(arg string) error {
	return fmt.Errorf("argument '%s' is invalid", arg)
}

func RequiredArgumentError(arg string) error {
	return fmt.Errorf("argument '%s' is required", arg)
}
