package security

import (
	"context"
	"errors"
	"fmt"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/config"
	"go.uber.org/fx"
	"strings"
)

func RegisterHooks(lifecycle fx.Lifecycle, config *config.Config, provider *UserProvider) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if len(config.Security.Users) == 0 {
				return errors.New("security: users not provided")
			}

			for i, user := range config.Security.Users {
				name := strings.TrimSpace(user.Name)
				if u := provider.Get(name); u != nil {
					return fmt.Errorf("security: users[%d].name is already used", i)
				}

				hashedPassword, ok := ParseHashedPassword(user.Password, config.Security.BCryptCost)
				if !ok {
					return fmt.Errorf("security: users[%d].password is illegal", i)
				}

				scopes := ParseScopes(user.Scopes)
				if len(scopes) == 0 {
					return fmt.Errorf("security: users[%d].scopes is not provided", i)
				}

				u := NewUser(name, hashedPassword, scopes)

				provider.Set(name, u)
			}

			return nil
		},
	})
}
