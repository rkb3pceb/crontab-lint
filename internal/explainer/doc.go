// Package explainer breaks down crontab expressions into structured,
// field-by-field explanations suitable for display or serialization.
//
// It builds on top of the parser and humanizer packages to provide both
// a per-field description and an overall schedule summary.
//
// Example usage:
//
//	result, err := explainer.Explain("0 9 * * 1 /usr/bin/backup")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(result.Summary)
//	for _, f := range result.Fields {
//		fmt.Printf("  %-14s %s → %s\n", f.Field, f.Value, f.Description)
//	}
package explainer
