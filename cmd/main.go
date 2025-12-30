package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/qrusaed/file-renamer/pkg/rename"
	"github.com/qrusaed/file-renamer/pkg/search"
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

	// file options
	replace := flag.String("replace", "", "replaces part of the filename, use with \"-with\" flag")
	with := flag.String("with", "", "provides replacement value, use with -replace flag")
	prefix := flag.String("prefix", "", "adds prefix to filename")
	suffix := flag.String("suffix", "", "adds suffix to filename")

	flag.Parse()
	patternList := strings.Split(*pattern, ",")

	s := &search.Searcher{
		Patterns: patternList,
		Depth:    *depth,
	}

	r := rename.Renamer{
		Suffix:  *suffix,
		Prefix:  *prefix,
		Replace: *replace,
		With:    *with,
	}

	// search the directory provided by dir flag
	fileList, err := s.Search(*dir)
	if err != nil {
		fmt.Println("rn:", err.Error())
		os.Exit(1)
	}

	replacement := r.CompileAll(fileList)
	// show a preview of to be made changes
	previewChanges(replacement)

	if err := r.Rename(replacement); err != nil {
		fmt.Println("rn:", err.Error())
		os.Exit(1)
	}
}

// previewChanges output a previews to the stdout of the changes
// that will be made to the files in the root directory
func previewChanges(changeList map[string]string) {
	for k, v := range changeList {
		fmt.Println(k, "--->", v)
	}
}
