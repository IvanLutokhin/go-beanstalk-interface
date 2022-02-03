package version

import "fmt"

var (
	title  = "unknown"
	tag    = "unknown"
	commit = "unknown"
)

func Title() string {
	return title
}

func Tag() string {
	return tag
}

func Commit() string {
	return commit
}

func String() string {
	return fmt.Sprintf("%s %s (Build: %s)", Title(), Tag(), Commit())
}
