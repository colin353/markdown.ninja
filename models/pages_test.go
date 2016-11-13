package models

import (
	"fmt"
	"testing"
)

func init() {
	Connect()

	// Delete all file records under the domain "testdomain"
	p := Page{}
	p.Domain = "testdomain"
	iterator, err := GetList(&p)
	if err != nil {
		panic("Unable to get a list of page records under `testdomain`")
	}

	for iterator.Next() {
		Delete(iterator.Value())
	}
}

func TestSiteValidation(t *testing.T) {
	p := &Page{}
	p.Domain = "testdomain"
	p.MakeDefault()
	if !p.Validate() {
		t.Error("Default page should pass validation.")
	}

	p.Domain = "google__333.foo"
	if p.Validate() {
		t.Error("Domain validation should have failed, but it didn't.")
	}

	p.Domain = "testdomain"
	p.Name = "my_file 3.txt/dev/null"
	if p.Validate() {
		t.Error("Filename should not validate, but did.")
	}
}

func TestPageNameGeneration(t *testing.T) {
	p := Page{}
	p.Domain = "testdomain"
	p.GenerateName()
	err := Insert(&p)
	if err != nil {
		t.Fatalf("Didn't generate valid name: got `%v`", p.Name)
	}

	for i := 0; i < 9; i++ {
		p.GenerateName()
		err = Insert(&p)
		if err != nil {
			t.Fatalf("Didn't generate valid name: got `%v`", p.Name)
		}
	}

	// Will fail to generate an additional name, but that's okay.
	p.GenerateName()
	Insert(&p)
}

func TestPageRename(t *testing.T) {
	p := Page{}
	p.Domain = "testdomain"
	p.Name = "originalName.md"
	p.Markdown = "hello world"
	p.HTML = "<p>hello world</p>"

	err := Insert(&p)
	if err != nil {
		t.Fatalf("Failed to insert page.")
	}

	err = p.RenamePage("newName.md")
	if err != nil {
		t.Fatalf("Failed to rename page: %v", err.Error())
	}

	g := Page{}
	err = Load(&g)
	if err != nil {
		t.Fatalf("Failed to load newly renamed page: %v", err.Error())
	}

	if g.Markdown != "hello world" {
		t.Fatalf("Renamed page didn't contain all data.")
	}

}

func TestSiteInsertion(t *testing.T) {
	p := Page{}
	p.MakeDefault()
	p.Domain = "testdomain"
	p.Name = "test.txt"
	p.HTML = "<p>test</p>"
	p.Markdown = "testing 1 2 3"

	err := Insert(&p)
	if err != nil {
		t.Fatalf("Error inserting new page record.")
	}

	g := Page{}
	g.Domain = "testdomain"
	g.Name = "test.txt"
	err = Load(&g)
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		t.Fatal("Couldn't load saved file.")
	}

	if g.Markdown != "testing 1 2 3" {
		t.Fatal("Didn't save correctly.")
	}

	result := g.Export()
	for _, s := range []string{"name", "markdown"} {
		_, ok := result[s]
		if !ok {
			t.Fatalf("Didn't find required field `%v` in exported page.", s)
		}
	}
}
