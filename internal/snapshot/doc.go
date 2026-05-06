// Package snapshot provides point-in-time capture of crontab expression
// metadata, enabling change detection and auditing over time.
//
// A Snapshot records the normalized form, per-field breakdown, and computed
// statistics of a crontab expression at the moment Take is called. Snapshots
// can be compared with Compare to produce a per-field diff, making it easy to
// detect what changed between two versions of a schedule.
//
// Example usage:
//
//	before := snapshot.Take("0 9 * * 1 /usr/bin/report")
//	// ... time passes, expression changes ...
//	after := snapshot.Take("0 10 * * 1 /usr/bin/report")
//
//	diffs, err := snapshot.Compare(before, after)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, d := range diffs {
//		if d.Changed {
//			fmt.Printf("%s: %s -> %s\n", d.Field, d.Before, d.After)
//		}
//	}
package snapshot
