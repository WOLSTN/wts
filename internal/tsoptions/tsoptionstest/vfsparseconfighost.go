package tsoptionstest

import (
	"github.com/wolstn/wts/internal/tsoptions"
	"github.com/wolstn/wts/internal/tspath"
	"github.com/wolstn/wts/internal/vfs"
	"github.com/wolstn/wts/internal/vfs/vfstest"
)

func fixRoot(path string) string {
	rootLength := tspath.GetRootLength(path)
	if rootLength == 0 {
		return path
	}
	if len(path) == rootLength {
		return "."
	}
	return path[rootLength:]
}

type VfsParseConfigHost struct {
	Vfs              vfs.FS
	CurrentDirectory string
}

var _ tsoptions.ParseConfigHost = (*VfsParseConfigHost)(nil)

func (h *VfsParseConfigHost) FS() vfs.FS {
	return h.Vfs
}

func (h *VfsParseConfigHost) GetCurrentDirectory() string {
	return h.CurrentDirectory
}

func NewVFSParseConfigHost(files map[string]string, currentDirectory string, useCaseSensitiveFileNames bool) *VfsParseConfigHost {
	return &VfsParseConfigHost{
		Vfs:              vfstest.FromMap(files, useCaseSensitiveFileNames),
		CurrentDirectory: currentDirectory,
	}
}
