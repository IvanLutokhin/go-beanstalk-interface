package security

import "strings"

const (
	ScopeReadServer = "read:server"
	ScopeReadTubes  = "read:tubes"
	ScopeWriteTubes = "write:tubes"
	ScopeReadJobs   = "read:jobs"
	ScopeWriteJobs  = "write:jobs"
)

type Scope string

func (s Scope) String() string {
	return string(s)
}

func GetAvailableScopes() []Scope {
	return []Scope{
		ScopeReadServer,
		ScopeReadTubes,
		ScopeWriteTubes,
		ScopeReadJobs,
		ScopeWriteJobs,
	}
}

func IsAvailableScope(scope Scope) bool {
	availableScopes := GetAvailableScopes()
	for _, availableScope := range availableScopes {
		if availableScope == scope {
			return true
		}
	}

	return false
}

func ParseScopes(values []string) (scopes []Scope) {
	m := make(map[Scope]struct{})

	for _, value := range values {
		scope := Scope(strings.ToLower(value))

		if _, found := m[scope]; !found && IsAvailableScope(scope) {
			m[scope] = struct{}{}
		}
	}

	for scope := range m {
		scopes = append(scopes, scope)
	}

	return
}

func VerifyScopes(actualScopes, expectedScopes []Scope) (requiredScopes []Scope) {
	m := make(map[Scope]struct{}, len(actualScopes))
	for _, scope := range actualScopes {
		m[scope] = struct{}{}
	}

	for _, scope := range expectedScopes {
		if _, found := m[scope]; !found && IsAvailableScope(scope) {
			requiredScopes = append(requiredScopes, scope)
		}
	}

	return requiredScopes
}
