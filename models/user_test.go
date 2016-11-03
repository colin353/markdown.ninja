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

}

func TestUserCreation(t *testing.T) {
  u := NewUser()
  u.Name = "Colin Merkel"
  u.SetPassword("gluten tag")
  u.Email = "colin353@gmail.com"
  u.MakeDefault()

  Delete(u)

  if !u.CheckPassword("gluten tag") {
    t.Fatal("Authentication failed, even with the correct password.")
  }
  if u.CheckPassword("gluten smaug") {
    t.Fatal("Authentication succeeded, even with the wrong password.")
  }

  Save(u)

  g := User{}
  g.Domain = "colinmerkel"
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
