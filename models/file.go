/*
  file.go

  The file model describes a file that the user has uploaded.
*/

package models

import (
	"fmt"
	"log"
	"os"
	"regexp"
)

// A File is a file that a user has uploaded, such
// as a resume or an image.
type File struct {
	Name   string `json:"name"`
	Hash   string `json:"hash"`
	Size   int    `json:"size"`
	Domain string `json:"domain"`
}

// MakeDefault initializes the file and sets defaults.
func (f *File) MakeDefault() {}

// Export returns the fields which are acceptable to send
// to the client as JSON.
func (f *File) Export() map[string]interface{} {
	return map[string]interface{}{
		"name": f.Name,
	}
}

// Key returns a unique key for use in the redis database.
func (f *File) Key() string {
	return fmt.Sprintf("files:%s:%s", f.Domain, f.Name)
}

// RegistrationKey defines the set to which this file belongs.
func (f *File) RegistrationKey() string {
	return fmt.Sprintf("files:%s", f.Domain)
}

// Validate checks the fields of the page to make sure they are
// acceptable to be inserted into the database.
func (f *File) Validate() bool {
	if !domainValidator.MatchString(f.Domain) {
		log.Printf("Validation failed on page %s, illegal domain '%s'\n", f.Name, f.Domain)
		return false
	}

	if !filenameValidator.MatchString(f.Name) {
		log.Printf("Validation failed on page %s, illegal filename '%s'\n", f.Name, f.Name)
		return false
	}

	return true
}

// GetPath returns the path that the file should be uploaded
// to, or can be accessed from.
func (f *File) GetPath() string {
	return fmt.Sprintf("./data/%s-%s-%s", f.Domain, f.Hash, f.Name)
}

var filenameReplacer = regexp.MustCompile("[^A-Za-z0-9_\\.]+")

// SetNameSafely strips illegal characters from the filename and
// is guaranteed to result in a valid name.
func (f *File) SetNameSafely(name string) {
	f.Name = string(filenameReplacer.ReplaceAll([]byte(name), []byte("")))
}

// RenameFile takes an existing file and renames it. It's a bit tricky to rename the
// file, because the file name defines the key, which is required in lookups. So you can't
// just load the record, change the name, and save it.
func (f *File) RenameFile(newName string) error {
	pool, err := connectionPool.Get()
	if err != nil {
		log.Fatal("Couldn't connect to the redis database.")
		return err
	}

	oldPath := f.GetPath()
	oldKey := f.Key()
	f.Name = newName
	newPath := f.GetPath()

	// Need to check key validation, in case the new name is not valid.
	if !f.Validate() {
		return fmt.Errorf("Tried to rename file to invalid name `%s`", newName)
	}

	// Step one: rename the old key to the new key.
	pool.Cmd("RENAME", oldKey, f.Key())

	// Take the registration pool and delete the old key
	// and add a new key.
	pool.Cmd("SREM", f.RegistrationKey(), oldKey)
	pool.Cmd("SADD", f.RegistrationKey(), f.Key())

	// Rename the associated file.
	err = os.Rename(oldPath, newPath)
	if err != nil {
		log.Printf("Failed during file rename operation, tried `%v` -> `%v`", oldPath, newPath)
		return err
	}

	// Save the object with the new parameters.
	err = Save(f)
	if err != nil {
		log.Printf("Tried to rename file to `%v`, but couldn't save it.", newName)
		return err
	}

	return nil
}
