package models

import (
	"fmt"
	"log"
	"regexp"
)

// A Page is an HTML/markdown file. The HTML is used for
// rendering, and the markdown is used for editing. It's
// stored in the database under the key:
//    pages:[domain]:page_name
// and a list of those keys are stored under the
//    pages:[domain]
// list.
type Page struct {
	Domain   string `json:"domain"`
	Name     string `json:"name"`
	Markdown string `json:"markdown"`
	HTML     string `json:"html"`
}

// MakeDefault returns the default initialized page.
func (p *Page) MakeDefault() {
	if p.Name == "" {
		p.Name = "new_page.md"
	}
	if p.Markdown == "" {
		p.Markdown = "## Default new page\nThis is an example page."
	}
	if p.HTML == "" {
		p.HTML = "<h1>Default new page</h1><p>This is an example page.</p>"
	}
}

// Export returns the fields which are acceptable to send directly to the
// client over the web.
func (p *Page) Export() map[string]interface{} {
	return map[string]interface{}{
		"name":     p.Name,
		"markdown": p.Markdown,
	}
}

// RegistrationKey defines the set to which this page will
// belong. It'll be of the form:
//    pages:[domain]
func (p *Page) RegistrationKey() string {
	return fmt.Sprintf("pages:%s", p.Domain)
}

// Key returns a unique key for use in the redis database.
func (p *Page) Key() string {
	return fmt.Sprintf("pages:%s:%s", p.Domain, p.Name)
}

var filenameValidator = regexp.MustCompile("^[A-Za-z0-9_\\.]+$")

// Validate checks the fields of the page to make sure they are
// acceptable to be inserted into the database.
func (p *Page) Validate() bool {
	if !domainValidator.MatchString(p.Domain) {
		log.Printf("Validation failed on page %s, illegal domain '%s'\n", p.Name, p.Domain)
		return false
	}

	if !filenameValidator.MatchString(p.Name) {
		log.Printf("Validation failed on page %s, illegal filename '%s'\n", p.Name, p.Name)
		return false
	}

	return true
}
