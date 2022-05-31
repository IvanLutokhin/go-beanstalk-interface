package env

import (
	"errors"
	"os"
	"strconv"
)

var ErrVarNotExists = errors.New("env: variable does not exists")

func GetString(key string) (string, error) {
	if v, ok := os.LookupEnv(key); ok {
		return v, nil
	}

	return "", ErrVarNotExists
}

func MustGetString(key string, defaultValue string) string {
	if value, err := GetString(key); err == nil {
		return value
	} else if errors.Is(err, ErrVarNotExists) {
		return defaultValue
	} else {
		panic(err)
	}
}

func GetInt(key string) (int, error) {
	if v, ok := os.LookupEnv(key); ok {
		return strconv.Atoi(v)
	}

	return 0, ErrVarNotExists
}

func MustGetInt(key string, defaultValue int) int {
	if value, err := GetInt(key); err == nil {
		return value
	} else if errors.Is(err, ErrVarNotExists) {
		return defaultValue
	} else {
		panic(err)
	}
}
