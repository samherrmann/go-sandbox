package templates

import (
	"errors"
	"io/fs"
)

func MergeFS(filesys ...fs.FS) *mergedFS {
	merged := make(mergedFS, 0)
	for _, f := range filesys {
		merged = append(merged, f)
	}
	return &merged
}

type mergedFS []fs.FS

func (mfs *mergedFS) Add(filesys ...fs.FS) {
	*mfs = append(*mfs, filesys...)
}

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
