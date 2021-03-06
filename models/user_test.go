package models

import (
	"testing"
)

func TestValidation(t *testing.T) {
	// Create a user with valid data, which should
	// actually pass validation.
	u := NewUser()
	u.Name = "Test Testerson"
	u.Domain = "testdomain"
	u.Email = "test@test.com"

	// Set an illegal password hash, check for validation error.
	u.PasswordSalt = ""

	if u.Validate() {
		t.Fatal("Shouldn't be able to save a user with a short password salt.")
	}

	u.PasswordSalt = "1231231231231209120912091209120912093091"
	u.Domain = "###"

	if u.Validate() {
		t.Fatal("Shouldn't be able to save a user with an illegal domain.")
	}

	u.Domain = "asdf123"
	u.Email = "cc@dd"

	if u.Validate() {
		t.Fatal("Shouldn't be able to save a user with an invalid email.")
	}

	u.Email = "cc@dd.ee"

	if !u.Validate() {
		t.Fatal("User should be valid... but validation failed anyway.")
	}

	result := u.Export()
	if result["name"] != "Test Testerson" {
		t.Fatal("Exported user isn't correct.")
	}

	_, ok := result["password"]
	if ok {
		t.Fatal("Exported user shouldn't contain password hash.")
	}
	_, ok = result["password_salt"]
	if ok {
		t.Fatal("Exported user shouldn't contain password salt.")
	}
}

func TestUserCreation(t *testing.T) {
	u := NewUser()
	u.Name = "Test Tester"
	u.Domain = "testdomain"
	u.SetPassword("gluten tag")
	u.Email = "test123@gmail.com"
	u.MakeDefault()

	Delete(u)

	if !u.CheckPassword("gluten tag") {
		t.Fatal("Authentication failed, even with the correct password.")
	}
	if u.CheckPassword("gluten smaug") {
		t.Fatal("Authentication succeeded, even with the wrong password.")
	}

	Insert(u)

	g := User{}
	g.Domain = "testdomain"
	err := Load(&g)
	if err != nil {
		t.Fatal("Couldn't load saved user.")
	}

	if !g.CheckPassword("gluten tag") {
		t.Fatal("Authentication failed, even with the correct password.")
	}
	if g.CheckPassword("gluten smaug") {
		t.Fatal("Authentication succeeded, even with the wrong password.")
	}
}
