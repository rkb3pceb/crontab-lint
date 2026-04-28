// Package formatter provides output formatting for crontab lint results.
package formatter

import (
	"fmt"
	"io"
	"strings"

	"github.com/example/crontab-lint/internal/linter"
)

// Format controls the output style.
type Format string

const (
	FormatText Format = "text"
	FormatJSON  Format = "json"
)

// Formatter writes lint results to a writer.
type Formatter struct {
	Format Format
	Writer io.Writer
}

// New creates a new Formatter.
func New(format Format, w io.Writer) *Formatter {
	return &Formatter{Format: format, Writer: w}
}

// Write outputs the lint results according to the configured format.
func (f *Formatter) Write(results []linter.Result) error {
	switch f.Format {
	case FormatJSON:
		return f.writeJSON(results)
	default:
		return f.writeText(results)
	}
}

// Summary returns a brief summary string of the results, e.g. "3 ok, 1 error, 2 warnings".
func (f *Formatter) Summary(results []linter.Result) string {
	ok, errs, warns := 0, 0, 0
	for _, r := range results {
		if !r.Valid {
			errs++
		} else if len(r.Warnings) > 0 {
			warns++
		} else {
			ok++
		}
	}
	return fmt.Sprintf("%d ok, %d error(s), %d warning(s)", ok, errs, warns)
}

func (f *Formatter) writeText(results []linter.Result) error {
	for _, r := range results {
		status := "OK"
		if !r.Valid {
			status = "ERROR"
		} else if len(r.Warnings) > 0 {
			status = "WARN"
		}

		_, err := fmt.Fprintf(f.Writer, "[%s] %s\n", status, r.Expression)
		if err != nil {
			return err
		}

		if r.Human != "" {
			_, err = fmt.Fprintf(f.Writer, "  Schedule: %s\n", r.Human)
			if err != nil {
				return err
			}
		}

		for _, e := range r.Errors {
			_, err = fmt.Fprintf(f.Writer, "  error: %s\n", e)
			if err != nil {
				return err
			}
		}

		for _, w := range r.Warnings {
			_, err = fmt.Fprintf(f.Writer, "  warning: %s\n", w)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (f *Formatter) writeJSON(results []linter.Result) error {
	var sb strings.Builder
	sb.WriteString("[\n")
	for i, r := range results {
		sb.WriteString(resultToJSON(r))
		if i < len(results)-1 {
			sb.WriteString(",")
		}
		sb.WriteString("\n")
	}
	sb.WriteString("]\n")
	_, err := fmt.Fprint(f.Writer, sb.String())
	return err
}

func resultToJSON(r linter.Result) string {
	valid := "true"
	if !r.Valid {
		valid = "false"
	}
	errs := jsonStringArray(r.Errors)
	warns := jsonStringArray(r.Warnings)
	return fmt.Sprintf(
		`  {"expression":%q,"valid":%s,"human":%q,"errors":%s,"warnings":%s}`,
		r.Expression, valid, r.Human, errs, warns,
	)
}

func jsonStringArray(items []string) string {
	if len(items) == 0 {
		return "[]"
	}
	quoted := make([]string, len(items))
	for i, s := range items {
		quoted[i] = fmt.Sprintf("%q", s)
	}
	return "[" + strings.Join(quoted, ",") + "]"
}
