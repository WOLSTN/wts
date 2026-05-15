package sourcemap

import "github.com/wolstn/wts/internal/core"

type Source interface {
	Text() string
	FileName() string
	ECMALineMap() []core.TextPos
}
