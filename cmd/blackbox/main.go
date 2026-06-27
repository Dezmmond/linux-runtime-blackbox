package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/yourname/linux-runtime-blackbox/internal/collector"
	"github.com/yourname/linux-runtime-blackbox/internal/model"
	"github.com/yourname/linux-runtime-blackbox/internal/report"
)

const (
	exitGeneric     = 1
	exitUsage       = 2
	exitNotFound    = 3
	exitPermissions = 4
)

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	if len(args) == 0 {
		usage(os.Stderr)
		return exitUsage
	}

	switch args[0] {
	case "inspect":
		return runInspect(args[1:])
	case "explain":
		return runExplain(args[1:])
	case "version":
		fmt.Println("blackbox 0.1")
		return 0
	case "-h", "--help", "help":
		usage(os.Stdout)
		return 0
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", args[0])
		usage(os.Stderr)
		return exitUsage
	}
}

func runInspect(args []string) int {
	fs := flag.NewFlagSet("inspect", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	pid := fs.Int("pid", 0, "target process ID")
	output := fs.String("output", "", "write JSON report to path")
	pretty := fs.Bool("pretty", false, "print a human-readable report")
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	if *pid <= 0 {
		fmt.Fprintln(os.Stderr, "inspect requires --pid <pid>")
		return exitUsage
	}

	r, err := collector.CollectPID(*pid)
	if err != nil {
		return handleCollectError(err)
	}

	if *output != "" {
		if err := writeJSONFile(*output, r); err != nil {
			fmt.Fprintf(os.Stderr, "write output: %v\n", err)
			return exitGeneric
		}
	}

	if *pretty {
		if err := report.WritePretty(os.Stdout, r); err != nil {
			fmt.Fprintf(os.Stderr, "write pretty report: %v\n", err)
			return exitGeneric
		}
		return 0
	}

	if *output == "" {
		if err := report.WriteJSON(os.Stdout, r); err != nil {
			fmt.Fprintf(os.Stderr, "write JSON report: %v\n", err)
			return exitGeneric
		}
	}
	return 0
}

func runExplain(args []string) int {
	fs := flag.NewFlagSet("explain", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	if fs.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "explain requires a report JSON path")
		return exitUsage
	}
	f, err := os.Open(fs.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "open report: %v\n", err)
		return exitGeneric
	}
	defer f.Close()
	r, err := report.ReadJSON(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read report: %v\n", err)
		return exitGeneric
	}
	if err := report.WritePretty(os.Stdout, r); err != nil {
		fmt.Fprintf(os.Stderr, "write pretty report: %v\n", err)
		return exitGeneric
	}
	return 0
}

func writeJSONFile(path string, r model.Report) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return report.WriteJSON(f, r)
}

func handleCollectError(err error) int {
	if errors.Is(err, collector.ErrProcessNotFound) {
		fmt.Fprintln(os.Stderr, "target process not found")
		return exitNotFound
	}
	if errors.Is(err, os.ErrPermission) {
		fmt.Fprintf(os.Stderr, "insufficient permissions: %v\n", err)
		return exitPermissions
	}
	fmt.Fprintf(os.Stderr, "inspect failed: %v\n", err)
	return exitGeneric
}

func usage(out *os.File) {
	fmt.Fprintf(out, `Usage:
  blackbox inspect --pid <pid> [--output report.json] [--pretty]
  blackbox explain report.json
  blackbox version

Examples:
  blackbox inspect --pid %s --pretty
  blackbox inspect --pid %s --output /tmp/blackbox-report.json
`, strconv.Itoa(os.Getpid()), strconv.Itoa(os.Getpid()))
}
