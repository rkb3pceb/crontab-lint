// Package formatter provides output formatting utilities for crontab-lint results.
//
// It supports multiple output formats:
//
//   - FormatText: human-readable plain text output, suitable for terminal use.
//     Each result is prefixed with a status tag ([OK], [WARN], or [ERROR]),
//     followed by the human-readable schedule description and any diagnostics.
//
//   - FormatJSON: machine-readable JSON array output, suitable for integration
//     with editors, CI pipelines, or other tooling.
//
// Example usage:
//
//	f := formatter.New(formatter.FormatText, os.Stdout)
//	f.Write(results)
package formatter
