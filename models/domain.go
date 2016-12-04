package models

import (
	"fmt"
	"log"
	"regexp"
)

// A Domain is a custom domain owned by a user, which
// they'll set up a CNAME record to markdown.ninja. When
// we get a request from that domain, we'll use these
// records to figure out which user's page to show.
type Domain struct {
	ExternalDomain string `json:"external_domain"`
	InternalDomain string `json:"internal_domain"`
}

// MakeDefault sets default values for the Domain struct.
func (d *Domain) MakeDefault() {}

// Export returns public fields for the struct which can
// be sent to a user, e.g. via web request.
func (d *Domain) Export() map[string]interface{} {
	return map[string]interface{}{
		"internal_domain": d.InternalDomain,
		"external_domain": d.ExternalDomain,
	}
}

// RegistrationKey defines the set to which this page
// will belong.
func (d *Domain) RegistrationKey() string {
	return "domains"
}

// Key returns a unique key for use in the redis database.
func (d *Domain) Key() string {
	return fmt.Sprintf("domains:%s", d.ExternalDomain)
}

var externalDomainValidator = regexp.MustCompile("^[A-Za-z0-9\\.]*$")

// Validate checks that the domain struct contains
// acceptable data.
func (d *Domain) Validate() bool {
	if len(d.ExternalDomain) < 2 {
		log.Printf("Validation failed on domain %s, too short", d.ExternalDomain)
		return false
	}
	if !domainValidator.MatchString(d.InternalDomain) || !externalDomainValidator.MatchString(d.ExternalDomain) {
		log.Printf("Validation failed on domain %s, illegal domain '%s'", d.InternalDomain, d.ExternalDomain)
		return false
	}

	return true
}
