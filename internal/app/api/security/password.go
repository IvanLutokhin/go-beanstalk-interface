package security

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strings"
)

const (
	PasswordTypePlain   = "plain"
	PasswordTypeEnv     = "env"
	PasswordTypeEncrypt = "encrypt"
)

type HashedPasswordParser interface {
	Parse(value string, bcryptCost int) ([]byte, error)
}

type plainParser struct{}

func (p *plainParser) Parse(value string, bcryptCost int) ([]byte, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(value), bcryptCost)
	if err != nil {
		return nil, err
	}

	return h, nil
}

type envParser struct{}

func (p *envParser) Parse(value string, bcryptCost int) ([]byte, error) {
	ev, exists := os.LookupEnv(value)
	if !exists {
		return nil, errors.New("environment variable does not exists")
	}

	h, err := bcrypt.GenerateFromPassword([]byte(ev), bcryptCost)
	if err != nil {
		return nil, err
	}

	return h, nil
}

type encryptParser struct{}

func (p *encryptParser) Parse(value string, bcryptCost int) ([]byte, error) {
	h := []byte(value)

	c, err := bcrypt.Cost(h)

	if err != nil {
		return nil, err
	}

	if c != bcryptCost {
		return nil, errors.New("bcrypt cost is not equals")
	}

	return h, nil
}

var parsers = map[string]HashedPasswordParser{
	PasswordTypePlain:   &plainParser{},
	PasswordTypeEnv:     &envParser{},
	PasswordTypeEncrypt: &encryptParser{},
}

func MustGetPasswordParser(key string) HashedPasswordParser {
	if parser, found := parsers[key]; found {
		return parser
	}

	panic("password parser is not provided")
}

func ParseHashedPassword(value string, bcryptCost int) ([]byte, bool) {
	if len(value) == 0 {
		return nil, false
	}

	i := strings.Index(value, ":")
	if i == -1 {
		return nil, false
	}

	t, v := value[:i], value[i+1:]

	if len(v) == 0 {
		return nil, false
	}

	if !strings.HasPrefix(t, "!") {
		return nil, false
	}

	pt := t[1:]

	if parser, found := parsers[pt]; found {
		if h, err := parser.Parse(v, bcryptCost); err == nil {
			return h, true
		}

		return nil, false
	}

	return nil, false
}

func VerifyPassword(hashedPassword, password []byte) bool {
	if err := bcrypt.CompareHashAndPassword(hashedPassword, password); err != nil {
		return false
	}

	return true
}
