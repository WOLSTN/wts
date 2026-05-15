package printer

import (
	"github.com/wolstn/wts/internal/ast"
	"github.com/wolstn/wts/internal/tspath"
)

type SourceFileMetaDataProvider interface {
	GetSourceFileMetaData(path tspath.Path) *ast.SourceFileMetaData
}
