// Package digest computes a stable fingerprint for a crontab schedule expression.
//
// A digest captures two complementary representations:
//
//   - Hash: a short (16-character) hexadecimal SHA-256 fingerprint of the
//     normalized schedule fields. Two expressions that are semantically
//     equivalent (e.g. "@hourly" and "0 * * * *") will produce the same hash.
//
//   - Signature: a compact, human-readable key/value string that encodes each
//     of the five cron fields, e.g. "m=0 h=* dom=* mon=* dow=1".
//
// Typical usage:
//
//	r, err := digest.Compute("*/15 * * * * /usr/bin/backup")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(r.Hash)      // e.g. "3f4a1b2c9d8e7f06"
//	fmt.Println(r.Signature) // e.g. "m=*/15 h=* dom=* mon=* dow=*"
package digest
