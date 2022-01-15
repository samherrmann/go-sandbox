// Package relpath provides the ability to unmarshal file paths defined in JSON
// and evaluate them relative to a given base path. A common use case is the
// consumption of a configuration file that contains relative paths to other
// files. Relative paths in configuration files are most often expected to be
// relative to the configuration file itself and not relative to the working
// directory of the process making use of the configuration file.
//
// Pros of this package:
//  - Does not use or need a custom JSON decoder.
//    - Works with the encoding/json package in the Go standard library.
//  - Does not use reflection.
//  - No "processing" code to maintain when adding/removing path fields from the
//    object.
//  - Does not put the burden on the user of the object to evaluate paths.
//    - Paths are evaluated during the JSON encoding/decoding process.
//
// Cons of this package:
//  - This package uses a mutex that results in concurrent JSON encode/decode
//    operations to execute sequentially.
//    - For the configuration file use case, this should not be an issue because
//      most likey there are not multiple configuration files that are processed
//      concurrently.
//    - JSON encode/decode opeations that do not use this package can run
//      concurrently with one JSON encode/decode operation that does use this
//      package.
package relpath

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

var (
	basepath = ""
	mu       = &sync.Mutex{}
)

// Lock locks b as the basepath for relative paths that are unmarshaled into the
// Path type. If b is a path to a file, then the directory that the file is
// located in is used as the basepath.
func Lock(b string) error {
	mu.Lock()
	fileInfo, err := os.Stat(b)
	if err != nil {
		mu.Unlock()
		return err
	}
	if !fileInfo.IsDir() {
		b = filepath.Dir(b)
	}
	basepath = filepath.FromSlash(b)
	return nil
}

// Unlock unlocks basepath.
func Unlock() {
	basepath = ""
	mu.Unlock()
}

type Path string

// UnmarshalJSON parses the JSON-encoded data and stores the result in Path.
// UnmarshalJSON is automatically called by json.Unmarshal.
func (p *Path) UnmarshalJSON(data []byte) error {
	path := ""
	if err := json.Unmarshal(data, &path); err != nil {
		return err
	}
	// Join path with basepath if it's not an absolute path.
	if !filepath.IsAbs(path) {
		path = filepath.Join(basepath, path)
	}
	*p = Path(path)
	return nil
}

// MarshalJSON returns the JSON encoding of Path.
func (p *Path) MarshalJSON() ([]byte, error) {
	path := p.String()
	// Remove basepath from path if path is not an absolute path.
	if !filepath.IsAbs(path) {
		var err error
		path, err = filepath.Rel(basepath, path)
		if err != nil {
			return nil, err
		}
	}
	return json.Marshal(path)
}

// String returns Path as a string.
func (p *Path) String() string {
	return string(*p)
}
