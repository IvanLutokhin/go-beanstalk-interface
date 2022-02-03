package embed

import "io/fs"

type FSFunc func(name string) (fs.File, error)

func (f FSFunc) Open(name string) (fs.File, error) {
	return f(name)
}
