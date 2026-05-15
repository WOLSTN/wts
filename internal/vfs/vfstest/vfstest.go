package vfstest

import (
	"io/fs"
	"testing/fstest"
	"time"

	"github.com/wolstn/wts/internal/vfs"
)

func FromMap(files map[string]string, useCaseSensitiveFileNames bool) vfs.FS {
	mapFS := make(fstest.MapFS)
	for path, content := range files {
		mapFS[path] = &fstest.MapFile{
			Data: []byte(content),
		}
	}
	return FromMapFS(mapFS, useCaseSensitiveFileNames)
}

func FromMapFS(mapFS fstest.MapFS, useCaseSensitiveFileNames bool) vfs.FS {
	return &mapFSWrapper{
		MapFS:                     mapFS,
		useCaseSensitiveFileNames: useCaseSensitiveFileNames,
	}
}

type mapFSWrapper struct {
	fstest.MapFS
	useCaseSensitiveFileNames bool
}

func (m *mapFSWrapper) UseCaseSensitiveFileNames() bool {
	return m.useCaseSensitiveFileNames
}

func (m *mapFSWrapper) ReadFile(path string) (string, bool) {
	data, err := fs.ReadFile(m, path)
	if err != nil {
		return "", false
	}
	return string(data), true
}

func (m *mapFSWrapper) FileExists(path string) bool {
	info := m.Stat(path)
	return info != nil && !info.IsDir()
}

func (m *mapFSWrapper) DirectoryExists(path string) bool {
	info := m.Stat(path)
	return info != nil && info.IsDir()
}

func (m *mapFSWrapper) WriteFile(path string, data string) error {
	m.MapFS[path] = &fstest.MapFile{Data: []byte(data)}
	return nil
}

func (m *mapFSWrapper) AppendFile(path string, data string) error {
	existing, ok := m.MapFS[path]
	if !ok {
		return m.WriteFile(path, data)
	}
	existing.Data = append(existing.Data, []byte(data)...)
	return nil
}

func (m *mapFSWrapper) Remove(path string) error {
	delete(m.MapFS, path)
	return nil
}

func (m *mapFSWrapper) Chtimes(path string, aTime time.Time, mTime time.Time) error {
	if f, ok := m.MapFS[path]; ok {
		f.ModTime = mTime
	}
	return nil
}

func (m *mapFSWrapper) GetAccessibleEntries(path string) vfs.Entries {
	entries, _ := fs.ReadDir(m, path)
	var result vfs.Entries
	for _, e := range entries {
		if e.IsDir() {
			result.Directories = append(result.Directories, e.Name())
		} else {
			result.Files = append(result.Files, e.Name())
		}
	}
	return result
}

func (m *mapFSWrapper) Realpath(path string) string {
	return path
}

func (m *mapFSWrapper) Stat(path string) vfs.FileInfo {
	info, err := m.MapFS.Stat(path)
	if err != nil {
		return nil
	}
	return info
}

func (m *mapFSWrapper) WalkDir(root string, walkFn vfs.WalkDirFunc) error {
	return fs.WalkDir(m.MapFS, root, walkFn)
}
