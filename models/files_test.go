package models

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestFileValidation(t *testing.T) {
	f := File{}
	f.Name = "!illegal-name$"
	f.Domain = "testdomain"

	err := Insert(&f)
	if err == nil {
		t.Fatalf("Validation failed to catch illegal filename.")
	}

	f.SetNameSafely("!illegal-name$")
	err = Insert(&f)
	if err != nil {
		t.Fatalf("Unable to safely set illegal name -> legal one? New name: `%v`. Reason: %s", f.Name, err.Error())
	}

	f.Domain = "!illegal-$domain"
	err = Insert(&f)
	if err == nil {
		t.Fatalf("Shouldn't be allowed to save to illegal domain name `%s`", f.Domain)
	}
}

func TestFileCreation(t *testing.T) {
	f := File{}
	f.Name = "test123.md"
	f.Domain = "testdomain"
	f.Hash = "abc123"

	err := Insert(&f)
	if err != nil {
		t.Fatalf("Failed to save file because: %v", err.Error())
	}

	result := f.Export()
	if result["name"] != "test123.md" {
		t.Fatalf("Didn't export the file correctly.")
	}

	file, err := os.Create(f.GetPath())
	if err != nil {
		t.Fatalf("Couldn't write to the file's path, `%v`", f.GetPath())
	}
	file.Close()

	err = os.Remove(f.GetPath())
	if err != nil {
		t.Fatalf("Unable to delete the file path, `%v`", f.GetPath())
	}

}

func TestFileRename(t *testing.T) {
	f := File{}
	f.Name = "test111.md"
	f.Domain = "testdomain"
	f.Hash = "abcdef012345"
	err := Insert(&f)
	if err != nil {
		t.Fatalf("Unable to create file test111.md in database: %v", err.Error())
	}

	ioutil.WriteFile(f.GetPath(), []byte("this is a test"), 0644)

	err = f.RenameFile("test456.md")
	if err != nil {
		t.Fatalf("Failed to rename the file.")
	}

	g := File{}
	g.Name = "test456.md"
	g.Domain = "testdomain"
	err = Load(&g)
	if err != nil {
		t.Fatalf("Error loading file: %s", err.Error())
	}

	fileContents, _ := ioutil.ReadFile(g.GetPath())
	if string(fileContents) != "this is a test" {
		log.Printf("Contents at %s were: %s", g.GetPath(), string(fileContents))
		t.Fatalf("After renaming file, the data stored in the file wasn't moved.")
	}

	// Now try to rename the file, but to an illegal name. This should prevent
	// the new name record being saved, stop the old one being deleted, and also
	// keep the file saved in the same place.
	err = g.RenameFile("!illegal-$name")
	if err == nil {
		t.Fatalf("Renaming file to an illegal name should fail.")
	}

	// Check that the file exists still.
	g.Name = "test456.md"
	fileContents, _ = ioutil.ReadFile(g.GetPath())
	if string(fileContents) != "this is a test" {
		t.Fatalf("After renaming file, the data stored in the file wasn't moved.")
	}

	// And check that the record wasn't deleted.
	err = Load(&g)
	if err != nil {
		t.Fatalf("After renaming to illegal name, original record was deleted.")
	}

	// And double check that the new record wasn't created.
	g.Name = "!illegal-$name"
	err = Load(&g)
	if err == nil {
		t.Fatalf("After renaming to an illegal name, new record was accidentally created anyway.")
	}

	g.Name = "test456.md"
	Delete(&g)
	err = os.Remove(g.GetPath())
	if err != nil {
		t.Fatalf("Failed to delete the file after renaming was successful.")
	}
}
