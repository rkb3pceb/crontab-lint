package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/example/crontab-lint/internal/formatter"
	"github.com/example/crontab-lint/internal/linter"
)

const version = "0.1.0"

func main() {
	var (
		format  = flag.String("format", "text", "Output format: text or json")
		showVer = flag.Bool("version", false, "Print version and exit")
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: crontab-lint [flags] <cron expression>\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExample:\n")
		fmt.Fprintf(os.Stderr, "  crontab-lint \"*/5 * * * * /usr/bin/backup.sh\"\n")
	}
	flag.Parse()

	if *showVer {
		fmt.Printf("crontab-lint version %s\n", version)
		os.Exit(0)
	}

	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "error: cron expression required")
		flag.Usage()
		os.Exit(2)
	}

	expr := strings.Join(flag.Args(), " ")
	result := linter.Lint(expr)

	fmt := formatter.New(os.Stdout)
	switch *format {
	case "json":
		if err := fmt.WriteJSON(result); err != nil {
			fmt.Fprintf(os.Stderr, "error writing JSON: %v\n", err)
			os.Exit(1)
		}
	case "text":
		fmt.WriteText(result)
	default:
		fmt.Fprintf(os.Stderr, "unknown format %q, use text or json\n", *format)
		os.Exit(2)
	}

	if !result.Valid {
		os.Exit(1)
	}
}
