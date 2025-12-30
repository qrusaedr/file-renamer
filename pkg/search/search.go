package search

import (
	"fmt"
	fsys "io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Searcher struct {
	Patterns []string
	Depth    int
}

// Search searches the root directory, and returns a slice
// containing filenames that matches a list given pattern
func (fs Searcher) Search(rootDir string) ([]string, error) {
	var filter []string
	path := filepath.Clean(rootDir)

	// the depth of root dir relative to from the dir the process is executed on
	rootDepth := strings.Count(path, string(filepath.Separator))
	err := filepath.WalkDir(path, func(path string, d fsys.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// calcutes depth of current directory then skip if equals search depth
		depth := strings.Count(path, string(filepath.Separator)) - rootDepth
		if fs.Depth != -1 && depth >= fs.Depth {
			return filepath.SkipDir
		}

		// filters filename that matches the patterns
		if ok := fs.MatchFilter(d.Name()); ok {
			filter = append(filter, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return filter, nil
}

// filter filename by matching them against glob expression
func (fs Searcher) MatchFilter(name string) bool {
	for _, pattern := range fs.Patterns {
		ok, err := filepath.Match(pattern, name)
		if err != nil {
			fmt.Println("rn: invalid match expression")
			os.Exit(1)
		}

		if ok {
			return ok
		}
	}
	return false
}
