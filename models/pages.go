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

// GenerateName creates a unique name
// for the page. This function finds a unique key (or returns err if it can't
// do that in a reasonable time) and attaches it to the page.
func (p *Page) GenerateName() error {
	pool, err := connectionPool.Get()
	if err != nil {
		log.Fatal("Couldn't connect to the redis database.")
	}
	for i := 0; i < 10; i++ {
		if i == 0 {
			p.Name = "untitled.md"
		} else {
			p.Name = fmt.Sprintf("untitled_%d.md", i)
		}
		result := pool.Cmd("EXISTS", p.Key())
		exists, err := result.Int()
		if err != nil {
			return err
		}
		if exists == 0 {
			return nil
		}
	}
	return fmt.Errorf("Couldn't generate a new key after significant effort (tried %s)", p.Name)
}

// RenamePage takes an existing page and rename it. It's a bit tricky to rename the
// page, because the page create sthe key, which prevents lookups.
func (p *Page) RenamePage(newName string) error {
	pool, err := connectionPool.Get()
	if err != nil {
		log.Fatal("Couldn't connect to the redis database.")
		return err
	}

	oldKey := p.Key()
	p.Name = newName

	// Need to check key validation, in case the new name is not valid.
	if !p.Validate() {
		return fmt.Errorf("Tried to rename to invalid name `%s`", newName)
	}

	// Step one: rename the old key to the new key.
	pool.Cmd("RENAME", oldKey, p.Key())

	// Take the registration pool and delete the old key
	// and add a new key.
	pool.Cmd("SREM", p.RegistrationKey(), oldKey)
	pool.Cmd("SADD", p.RegistrationKey(), p.Key())

	// Save the object with the new parameters.
	err = Save(p)
	if err != nil {
		log.Printf("Tried to rename page to `%v`, but couldn't save at the end.", newName)
	}

	return nil
}
