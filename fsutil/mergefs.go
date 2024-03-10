package fsutil

import (
	"errors"
	"io/fs"
)

// MergeFS returns a merged file system of the given file systems.
func MergeFS(filesys ...fs.FS) *mergedFS {
	merged := mergedFS{}
	for _, f := range filesys {
		merged = append(merged, f)
	}
	return &merged
}

type mergedFS []fs.FS

func (mfs *mergedFS) Open(name string) (fs.File, error) {
	for _, filesys := range *mfs {
		f, err := filesys.Open(name)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		return f, err
	}
	return nil, fs.ErrNotExist
}
