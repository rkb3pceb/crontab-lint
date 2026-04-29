// Package diff provides utilities for comparing two crontab schedule
// expressions and identifying which fields have changed between them.
//
// It is useful for code-review tooling, changelog generation, or any
// workflow where a crontab expression is being updated and the operator
// wants a clear, field-level summary of what changed.
//
// Example:
//
//	r, err := diff.Compare("0 9 * * 1 /bin/backup", "0 18 * * 5 /bin/backup")
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, d := range r.Diffs {
//		fmt.Printf("field %s: %s -> %s\n", d.Field, d.From, d.To)
//	}
//
// The Compare function validates both expressions via the parser package
// before performing the comparison, so callers do not need to validate
// inputs separately.
package diff
