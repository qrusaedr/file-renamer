package rename

import (
	"os"
	"path/filepath"
	"strings"
)

type Renamer struct {
	Suffix  string
	Prefix  string
	Replace string
	With    string
}

func (r Renamer) Rename(changeList map[string]string) error {
	for k, v := range changeList {
		if err := os.Rename(k, v); err != nil {
			lerr := err.(*os.LinkError)
			return lerr.Err
		}
	}
	return nil
}

func (r Renamer) CompileAll(filenames []string) map[string]string {
	replacements := make(map[string]string)
	for _, name := range filenames {
		replacements[name] = r.Compile(name)
	}
	return replacements
}

// Compile builds a new filename, replacement of old value in filename
// occurs before pre-/-suffixes to avoid potentially replacing them
func (r Renamer) Compile(filename string) string {
	dir := filepath.Dir(filename)
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filepath.Base(filename), ext)

	if r.Replace != "" {
		if r.With != "" {
			name = strings.ReplaceAll(name, r.Replace, r.With)
		}
	}

	if r.Prefix != "" {
		name = r.Prefix + name
	}

	if r.Suffix != "" {
		name = name + r.Suffix
	}

	return filepath.Join(dir, name+ext)
}
