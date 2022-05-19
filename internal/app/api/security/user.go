package security

import "sync"

type User struct {
	name           string
	hashedPassword []byte
	scopes         []Scope
}

func NewUser(name string, hashedPassword []byte, scopes []Scope) *User {
	return &User{
		name:           name,
		hashedPassword: hashedPassword,
		scopes:         scopes,
	}
}

func (u *User) Name() string {
	return u.name
}

func (u *User) HashedPassword() []byte {
	return u.hashedPassword
}

func (u *User) Scopes() []Scope {
	return u.scopes
}

type UserProvider struct {
	users sync.Map
}

func NewUserProvider() *UserProvider {
	return &UserProvider{}
}

func (p *UserProvider) Get(name string) *User {
	if item, ok := p.users.Load(name); ok {
		return item.(*User)
	}

	return nil
}

func (p *UserProvider) Set(name string, user *User) {
	p.users.Store(name, user)
}
