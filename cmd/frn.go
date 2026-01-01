package main

import (
	"bufio"
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
	dryRun := flag.Bool("dry-run", false, "preview before running")

	// rename flags
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

	// if -dry-run flag is set to true,
	// shows the changes without applying them until
	// confirmation
	if *dryRun {
		showPreview(replacement) // show a preview of to be made changes
		confirmChanges()         // promt for confirmation
	}

	if err := r.Rename(replacement); err != nil {
		fmt.Println("rn:", err.Error())
		os.Exit(1)
	}
}

// previewChanges output a previews to the stdout of the changes
// that will be made to the files in the root directory
func showPreview(changeList map[string]string) {
	for k, v := range changeList {
		fmt.Println(k, "--->", v)
	}
}

func confirmChanges() {
	fmt.Print("would you like to apply changes (y/n): ")
	in, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		fmt.Println("rn: end of file")
	}

	switch strings.TrimSpace(in) {
	case "y", "yes":
	case "n", "no":
		fmt.Println("rn: renaming discountinued")
		os.Exit(0)
	default:
		confirmChanges()
	}
}
