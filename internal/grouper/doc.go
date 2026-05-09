// Package grouper clusters a collection of crontab expressions by their
// normalised schedule, making it easy to identify duplicate or equivalent
// jobs within a crontab file or across multiple files.
//
// # Usage
//
//	exprs := []string{
//		"@hourly",
//		"0 * * * *",
//		"*/5 * * * *",
//	}
//
//	groups := grouper.Cluster(exprs)
//	for _, g := range groups {
//		fmt.Printf("Schedule %q — %d job(s)\n", g.Key, len(g.Expressions))
//		for _, e := range g.Expressions {
//			fmt.Printf("  %s\n", e)
//		}
//	}
//
// Normalisation is performed via the normalizer package so that aliases such
// as @hourly are expanded before comparison. Expressions that fail
// normalisation are each placed in their own singleton group.
package grouper
