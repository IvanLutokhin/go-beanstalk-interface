package model

import (
	"errors"
	"github.com/IvanLutokhin/go-beanstalk-interface/internal/app/api/security"
)

var ErrNotMappable = errors.New("enum can not be mapped")

var scopes = map[security.Scope]Scope{
	security.ScopeReadServer: ScopeReadServer,
	security.ScopeReadTubes:  ScopeReadTubes,
	security.ScopeReadJobs:   ScopeReadJobs,
	security.ScopeWriteJobs:  ScopeWriteJobs,
}

func MapScope(scope security.Scope) (Scope, error) {
	if s, ok := scopes[scope]; ok {
		return s, nil
	}

	return "", ErrNotMappable
}
