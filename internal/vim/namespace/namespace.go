package namespace

import "github.com/lancekrogers/lazysql/internal/vim/leader"

type Namespace interface {
	Prefix() rune
	Name() string
	Commands() []leader.Command
}
