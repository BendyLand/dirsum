package main

import (
	"cmp"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

func main() {
	if checkHelpArg(os.Args) || len(os.Args) == 1 {
		printHelp()
		return
	}
	path := getPath(os.Args)
	files := getFiles(path)
	exts := getFileExtensionsList(files)
	extCounts := countExts(files, exts)
	flags := getFlags(os.Args)
	flagStr := strings.Join(flags, " ")
	isVerbose := strings.Contains(flagStr, "v")
	if strings.Contains(flagStr, "n") || isVerbose {
		reverse := !strings.Contains(flagStr, "r")
		displayCountsByNumber(extCounts, reverse)
	} else {
		reverse := strings.Contains(flagStr, "r")
		displayCounts(extCounts, reverse)
	}
	if strings.Contains(flagStr, "t") || isVerbose {
		total := getTotal(extCounts)
		printTotal(total)
	}
}

func getPath(args []string) string {
	if len(args) > 1 {
		return args[1]
	}
	return "."
}

func getFiles(path string) []string {
	result := make([]string, 0)
	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		result = append(result, d.Name())
		return nil
	})
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	return result
}

func getFileExtensionsList(contents []string) []string {
	exts := make([]string, 0)
	for _, content := range contents {
		base := filepath.Base(content)
		if base != content {
			content = base
		}
		idx := strings.LastIndex(content, ".")
		if idx < 0 {
			continue
		}
		ext := content[idx:]
		if slices.Contains(exts, ext) {
			continue
		}
		exts = append(exts, ext)
	}
	return exts
}

func countExts(files []string, exts []string) map[string]int {
	result := make(map[string]int)
	for _, file := range files {
		if !endsWithAny(file, exts) {
			continue
		}
		ext := getFileExt(file)
		if _, ok := result[ext]; ok {
			result[ext]++
		} else {
			result[ext] = 1
		}
	}
	return result
}

func endsWithAny(file string, exts []string) bool {
	for _, ext := range exts {
		if strings.Index(file, ext) != -1 {
			return true
		}
	}
	return false
}

func displayCounts(extCounts map[string]int, reverse bool) {
	keys := make([]string, 0)
	for k, _ := range extCounts {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	if reverse {
		slices.Reverse(keys)
	}
	for _, key := range keys {
		fmt.Printf("%s: %d\n", key, extCounts[key])
	}
}

type Pair struct {
	key   string
	value int
}

func displayCountsByNumber(extCounts map[string]int, reverse bool) {
	pairs := make([]Pair, 0)
	for key, value := range extCounts {
		pairs = append(pairs, Pair{key, value})
	}
	slices.SortFunc(pairs, func(a, b Pair) int {
		return cmp.Compare(a.value, b.value)
	})
	if reverse {
		slices.Reverse(pairs)
	}
	for i := 0; i < len(pairs); i++ {
		fmt.Printf("%s: %d\n", pairs[i].key, pairs[i].value)
	}
}

func getFileExt(file string) string {
	idx := strings.LastIndex(file, ".")
	return file[idx:]
}

func getIdxOfNextHighestValue(values []int, target int) int {
	for i, value := range values {
		if target >= value {
			return i
		}
	}
	return len(values) - 1
}

func getFlags(args []string) []string {
	result := make([]string, 0)
	pattern, err := regexp.Compile("\\-[nrtv]*[nrtv]")
	if err != nil {
		fmt.Println("Unable to compile regex")
		os.Exit(1)
	}
	for _, arg := range args {
		if !strings.Contains(arg, "-") {
			continue
		}
		if pattern.MatchString(arg) {
			result = append(result, arg)
		}
	}
	return result
}

func getTotal(extCounts map[string]int) int {
	total := 0
	for _, v := range extCounts {
		total += v
	}
	return total
}

func printTotal(total int) {
	fmt.Printf("\nTotal files: %d\n", total)
}

func checkHelpArg(args []string) bool {
	for _, arg := range args {
		if arg == "-h" || arg == "--help" {
			return true
		}
	}
	return false
}

func printHelp() {
	fmt.Printf(
		"Welcome to the dirsum help menu!\n" +
			"Usage: dirsum <path> <opt args>\n" +
			"Valid arguments include:\n" +
			"-h or --help (see this menu)\n-n (sort by number)\n-r (reverse sort by number)\n-t (print the total number of files)\n-v (verbose: equal to -nt)\n" +
			"All flags may be chained together (e.g. -nrtv) except for -h and --help.\n",
	)
}
