package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type SearchConfig struct {
	patterns []string
	depth    int
}

type FileRenamer struct {
	prefix  string
	suffix  string
	replace string
	with    string
}

func main() {
	dir := flag.String("dir", ".", "specifies the root directory to transverse")
	pattern := flag.String("match", "*", "match filter the files based on given glob pattern")
	depth := flag.Int("depth", -1, "depth of subdirectories to transverse relative to the root directory")

	// file option
	replace := flag.String("replace", "", "replaces part of the filename, use with \"-with\" flag")
	with := flag.String("with", "", "provides replacement value, use with -replace flag")
	prefix := flag.String("prefix", "", "adds prefix to filename")
	suffix := flag.String("suffix", "", "adds suffix to filename")

	flag.Parse()
	patternList := strings.Split(*pattern, ",")

	fileList, err := listFiles(*dir, &SearchConfig{
		patterns: patternList,
		depth:    *depth,
	})

	if err != nil {
		panic(err)
	}

	replacments := buildFileNames(fileList, &FileRenamer{
		replace: *replace,
		with:    *with,
		prefix:  *prefix,
		suffix:  *suffix,
	})

	// show previes of changes to be made
	previewChanges(replacments)
	// rename file to a new filename
	if err := changeFilenames(replacments); err != nil {
		fmt.Println("not all file name changed")
	}
}

// listFiles func search the root directory, and returns a slice
// containing filenames that matches a given pattern
func listFiles(rootDir string, search *SearchConfig) ([]string, error) {
	var filter []string

	if !fs.ValidPath(rootDir) {
		panic("directory transversal not allowed")
	}

	// the depth of root dir relative to from the dir the process is executed on
	rootDeoth := strings.Count(rootDir, string(filepath.Separator))

	dirFs := os.DirFS(".")
	err := fs.WalkDir(dirFs, rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// calcutes depth of current directory then skip if equals search depth
		depth := strings.Count(path, string(filepath.Separator)) - rootDeoth
		if search.depth != -1 && depth >= search.depth {
			return filepath.SkipDir
		}

		// filters filename that matches the patterns
		if ok := matchFilter(d.Name(), search.patterns); ok {
			filter = append(filter, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return filter, nil
}

func changeFilenames(changeList map[string]string) error {
	for k, v := range changeList {
		if err := os.Rename(k, v); err != nil {
			return err
		}
	}
	return nil
}

// filter filename by matching them against a glob expression
func matchFilter(name string, patterns []string) bool {
	for _, pattern := range patterns {
		ok, err := filepath.Match(pattern, name)
		if err != nil {
			panic("invalid match expression")
		}

		if ok {
			return ok
		}
	}

	return false
}

func buildFileNames(filenames []string, replacement *FileRenamer) map[string]string {
	replacements := make(map[string]string)
	for _, name := range filenames {
		replacements[name] = compileFileName(name, replacement)
	}
	return replacements
}

// compileFileName builds a new name, replacement of old value in filename
// occurs before pre/suffixes to avoid replacing
func compileFileName(filename string, replacement *FileRenamer) string {
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filepath.Base(filename), ext)
	dir := filepath.Dir(filename)

	if replacement.replace != "" {
		name = strings.Replace(name, replacement.replace, replacement.with, 1)
	}

	if replacement.prefix != "" {
		name = replacement.prefix + name
	}

	if replacement.suffix != "" {
		name = name + replacement.suffix
	}

	return filepath.Join(dir, name+ext)
}

// previewChanges output a previews to the stdout of the changes
// that will be made to the files in the root directory
func previewChanges(changeList map[string]string) {
	for k, v := range changeList {
		fmt.Println(k, "--->", v)
	}
}
