package models

import (
  "fmt"
  "log"
  "strings"
  "regexp"
  "math/rand"
  "golang.org/x/crypto/scrypt"
)

type User struct {
  Name string `json:"name"`
  PasswordHash string `json:"password_hash"`
  PasswordSalt string `json:"password_salt"`
  Email string `json:"email"`
  PhoneNumber string `json:"phone_number"`
  Bio string `json:"bio"`
  Domain string `json:"domain"`
}

func NewUser() *User {
  u := new(User)

  // Create the password salt.
  const saltLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
  b := make([]byte, 32)
  for i := range b {
      b[i] = saltLetters[rand.Int63() % int64(len(saltLetters))]
  }
  u.PasswordSalt = string(b)

  return  u
}

func (u *User) SetPassword(password string) {
  dk, err := scrypt.Key([]byte(password), []byte(u.PasswordSalt), 16384, 8, 1, 32)
  if err != nil {
    log.Fatal("Unexpected error setting password.")
  }
  u.PasswordHash = string(dk)
}

func (u *User) CheckPassword(password string) bool {
  dk, err := scrypt.Key([]byte(password), []byte(u.PasswordSalt), 16384, 8, 1, 32)
  if err != nil {
    log.Fatal("Unexpected error setting password.")
  }
  return string(dk) == u.PasswordHash
}

func (u *User) MakeDefault() {
  // Try to create a domain based upon the given name.
  if u.Domain == "" && u.Name != "" {
    u.Domain = strings.ToLower(u.Name)
    u.Domain = regexp.MustCompile("[^A-Za-z0-9]+").ReplaceAllString(u.Domain, "")
  }
}

func (u *User) Key() string {
  return fmt.Sprintf("user:%s", u.Domain)
}

var domainValidator = regexp.MustCompile("^[A-Za-z0-9]+$")
var emailValidator = regexp.MustCompile(`^([a-zA-Z0-9_\-\.]+)@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.)|(([a-zA-Z0-9\-]+\.)+))([a-zA-Z]{2,4}|[0-9]{1,3})(\]?)$`)

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
