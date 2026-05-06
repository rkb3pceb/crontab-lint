// Package template provides a catalog of named crontab expression templates
// for common scheduling patterns.
//
// It allows tools and CLI commands to quickly look up well-known schedules by
// name or search for templates by keyword, tag, or description.
//
// # Usage
//
//	// List all available templates
//	templates := template.All()
//
//	// Look up a specific template by name
//	tmpl, err := template.Lookup("hourly")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(tmpl.Expression) // "0 * * * *"
//
//	// Search templates by keyword
//	results := template.Search("daily")
//	for _, r := range results {
//		fmt.Printf("%s\t%s\n", r.Name, r.Expression)
//	}
package template
