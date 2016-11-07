package models

import (
	"fmt"
	"golang.org/x/crypto/scrypt"
	"log"
	"math/rand"
	"regexp"
	"strings"
)

// The User struct defines a user, and stores their login data,
// domain name, etc. The subdomain name is the unique key for finding
// the user.
type User struct {
	Name         string `json:"name"`
	PasswordHash string `json:"password_hash"`
	PasswordSalt string `json:"password_salt"`
	Email        string `json:"email"`
	PhoneNumber  string `json:"phone_number"`
	Bio          string `json:"bio"`
	Domain       string `json:"domain"`
}

// Export converts a user into fields which are "safe" to export to
// the web (i.e. excluding sensitive fields like password hashes).
func (u *User) Export() map[string]interface{} {
	return map[string]interface{}{
		"name":         u.Name,
		"email":        u.Email,
		"phone_number": u.PhoneNumber,
		"domain":       u.Domain,
	}
}

// NewUser creates a new instance of the user object, presetting
// any fields which need to be set initially, such as the hash salt.
func NewUser() *User {
	u := new(User)

	// Create the password salt.
	const saltLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 32)
	for i := range b {
		b[i] = saltLetters[rand.Int63()%int64(len(saltLetters))]
	}
	u.PasswordSalt = string(b)

	return u
}

// SetPassword sets the password for a user, using scrypt.
func (u *User) SetPassword(password string) {
	dk, err := scrypt.Key([]byte(password), []byte(u.PasswordSalt), 16384, 8, 1, 32)
	if err != nil {
		log.Fatal("Unexpected error setting password.")
	}
	u.PasswordHash = string(dk)
}

// CheckPassword takes a string and checks if that matches the hashed password.
func (u *User) CheckPassword(password string) bool {
	dk, err := scrypt.Key([]byte(password), []byte(u.PasswordSalt), 16384, 8, 1, 32)
	if err != nil {
		log.Fatal("Unexpected error setting password.")
	}
	return string(dk) == u.PasswordHash
}

// MakeDefault sets the default fields for the user. If it has some
// partial information to work with (in this case the name) it'll try
// to fill in the domain name with a best guess.
func (u *User) MakeDefault() {
	// Try to create a domain based upon the given name.
	if u.Domain == "" && u.Name != "" {
		u.Domain = strings.ToLower(u.Name)
		u.Domain = regexp.MustCompile("[^A-Za-z0-9]+").ReplaceAllString(u.Domain, "")
	}
}

// Key returns a unique key for use in the redis database. The user database
// always has a prefix of user: before all keys. The suffix is the domain name
// of the user.
func (u *User) Key() string {
	return fmt.Sprintf("user:%s", u.Domain)
}

// RegistrationKey returns a key which is used to store a set of all sibling
// elements. In this case, all users are siblings to each other.
func (u *User) RegistrationKey() string {
	return "users"
}

var domainValidator = regexp.MustCompile("^[A-Za-z0-9]+$")
var emailValidator = regexp.MustCompile(`^([a-zA-Z0-9_\-\.]+)@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.)|(([a-zA-Z0-9\-]+\.)+))([a-zA-Z]{2,4}|[0-9]{1,3})(\]?)$`)

// Validate checks all the fields of the user object to ensure that they
// follow the correct rules (e.g. make sure that the email address looks right,
// and that the domain name isn't messed up.)
func (u *User) Validate() bool {
	// Check that the domain name is valid.
	if !domainValidator.MatchString(u.Domain) {
		log.Printf("Validation failed on user %s, illegal domain '%s'\n", u.Name, u.Domain)
		return false
	}

	// Check that the email address is valid.
	if !emailValidator.MatchString(u.Email) {
		log.Printf("Validation failed on user %s, illegal email '%s'\n", u.Name, u.Email)
		return false
	}

	if len(u.PasswordSalt) < 32 {
		log.Printf("Validation failed: insufficiently long salt.")
		return false
	}

	return true
}
