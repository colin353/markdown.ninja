package models

import (
  "testing"
  "fmt"
)

func init() {
  Connect()
}

func TestSiteValidation(t *testing.T) {
  p := &Page{}
  p.Domain = "test123"
  p.MakeDefault()
  if !p.Validate() {
    t.Error("Default page should pass validation.")
  }

  p.Domain = "google__333.foo"
  if p.Validate() {
    t.Error("Domain validation should have failed, but it didn't.")
  }

  p.Domain = "test123"
  p.Name = "my_file 3.txt/dev/null"
  if p.Validate() {
    t.Error("Filename should not validate, but did.")
  }
}

func TestSiteInsertion(t *testing.T) {
  p := Page{}
  p.MakeDefault()
  p.Domain = "test123"
  p.Name = "test.txt"
  p.Markdown = "testing 1 2 3"

  Save(&p)

  g := Page{}
  g.Domain = "test123"
  g.Name = "test.txt"
  err := Load(&g)
  if err != nil {
    fmt.Printf("Error: %v\n", err.Error())
    t.Fatal("Couldn't load saved file.")
  }

  if g.Markdown != "testing 1 2 3" {
    t.Fatal("Didn't save correctly.")
  }
}
