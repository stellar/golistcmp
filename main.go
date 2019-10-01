package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
)

type modInfo struct {
	Path    string
	Version string
}

type change struct {
	Path  string
	Left  string
	Right string
}

func main() {
	var flagHelp bool
	flag.BoolVar(&flagHelp, "help", false, "print this help")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of golistcmp:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  golistcmp <go list before> <go list after>\n")
		fmt.Fprintf(flag.CommandLine.Output(), "\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Example:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  git checkout master\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  go list -m -json all > go.list.master\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  git checkout mybranchwithchanges\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  golistcmp go.list.master <(go list -m -json all)\n")
		fmt.Fprintf(flag.CommandLine.Output(), "\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Flags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flagHelp {
		flag.Usage()
		os.Exit(0)
	}

	args := flag.Args()

	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "Two arguments required. See usage with "+os.Args[0]+" -help.")
		os.Exit(1)
	}

	left, err := os.Open(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	right, err := os.Open(args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	cmp(os.Stdout, left, right)
}

func cmp(w io.Writer, left, right io.Reader) {
	leftMods := decodeMods(left)
	rightMods := decodeMods(right)

	changes := []change{}
	for path, leftVersion := range leftMods {
		rightVersion := rightMods[path]
		if leftVersion == rightVersion {
			continue
		}
		change := change{
			Path:  path,
			Left:  leftVersion,
			Right: rightVersion,
		}
		changes = append(changes, change)
	}
	for path, rightVersion := range rightMods {
		_, ok := leftMods[path]
		if ok {
			continue
		}
		change := change{
			Path:  path,
			Right: rightVersion,
		}
		changes = append(changes, change)
	}

	sort.Slice(changes, func(i, j int) bool {
		return changes[i].Path < changes[j].Path
	})

	tw := new(tabwriter.Writer)

	tw.Init(w, 0, 8, 0, '\t', 0)
	for i, change := range changes {
		diff := ""
		switch {
		case change.Left == "":
			diff = "+"
		case change.Right == "":
			diff = "-"
		default:
			diff = ">"
		}
		fmt.Fprintf(tw, "%3d\t%s\t%s\t%s\t%s\t%s\t\n", i, change.Path, change.Left, diff, change.Right, diffURL(change))
	}
	tw.Flush()
}

func decodeMods(r io.Reader) map[string]string {
	mods := map[string]string{}
	dec := json.NewDecoder(r)
	for {
		info := modInfo{}
		err := dec.Decode(&info)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "error parsing", r, ":", err)
			os.Exit(1)
		}
		mods[info.Path] = normalizeVersion(info.Version)
	}
	return mods
}

func normalizeVersion(version string) string {
	if strings.Contains(version, "+") {
		version = strings.Split(version, "+")[0]
	}
	if strings.Contains(version, "-") {
		version = strings.Split(version, "-")[2]
	}
	return version
}

func diffURL(c change) string {
	if c.Right == "" {
		return ""
	}
	if c.Left == "" {
		return "https://" + c.Path
	}
	if strings.HasPrefix(c.Path, "golang.org/x/") {
		c.Path = strings.Replace(c.Path, "golang.org/x/", "github.com/golang/", 1)
	}
	if strings.HasPrefix(c.Path, "github.com/") {
		return "https://" + c.Path + "/compare/" + c.Left + "..." + c.Right
	}
	return ""
}
