package models

import "testing"

func init() {
	Connect()
}

type TestStructure struct {
	Model
	KeyStr   string `json:"key"`
	LongName string `json:"long_name"`
	Error    bool   `json:"error"`
	Age      int    `json:"age"`
}

func (s *TestStructure) MakeDefault() {
	s.LongName = "Johnny Appleseed"
	s.Error = false
	s.Age = 34
}

func (s *TestStructure) Export() map[string]interface{} {
	return nil
}

func (s *TestStructure) RegistrationKey() string {
	return "tests"
}

func (s *TestStructure) Validate() bool {
	return !s.Error
}

func (s *TestStructure) Key() string {
	if s.KeyStr != "" {
		return s.KeyStr
	}
	return "test:" + MakeKeyForTable("test")
}

func TestKeyCreation(t *testing.T) {
	a := MakeKeyForTable("teststructure")
	b := MakeKeyForTable("teststructure")

	if a == b {
		t.Fatalf(
			`Key creation should always return unique values,
			never the same. Got "%v", then "%v"`, a, b,
		)
	}

}

func TestIllegalLoad(t *testing.T) {
	s := TestStructure{}
	s.KeyStr = "aihsdfliuhasdf"
	err := Load(&s)
	if err == nil {
		t.Error("I tired to load illegal key but it worked??")
	}
}

func TestInsertion(t *testing.T) {
	s := TestStructure{}
	s.KeyStr = "ff8f8f8f8f8"

	// Delete all existing TestStructures from other tests.
	iterator, err := GetList(&s)
	if err != nil {
		t.Fatal("Failed to generate iterator.")
	}

	for iterator.Next() {
		ts := iterator.Value()
		Delete(ts)
	}

	Delete(&s)
	err = Save(&s)
	if err == nil {
		t.Fatal("Saved object when its key didn't exist. This should have failed.")
	}
	err = Insert(&s)
	if err != nil {
		t.Fatal("Should be able to insert object with unused key.")
	}
	err = Insert(&s)
	if err == nil {
		t.Fatal("Should not be able to insert object with the same key twice.")
	}

	// Check that the TestStructure was inserted into its registration set.
	iterator, err = GetList(&s)
	if err != nil {
		t.Fatal("Failed to generate iterator.")
	}

	if iterator.Count() != 1 {
		t.Fatalf("Should have exactly one teststructure after insertion, got %v", iterator.Count())
	}

	// Try to get the TestStructure from the list.
	ok := iterator.Next()
	if !ok {
		t.Fatal("Failed to get the TestStructure from the iterator, even though an item was present.")
	}

	g := iterator.Value()
	if g.Key() != s.Key() {
		t.Fatal("Retrieved object keys don't match.")
	}

	ok = iterator.Next()
	if ok {
		t.Fatal("Shouldn't be able to retrieve beyond the end of the list from the iterator.")
	}

	Delete(&s)

	// Check that the TestStructure was deleted from its registration set.
	iterator, err = GetList(&s)
	if err != nil {
		t.Fatal("Failed to generate iterator.")
	}

	if iterator.Count() != 0 {
		t.Fatalf("Should have exactly zero teststructures after deletion, got %v.", iterator.Count())
	}
}

func TestPartialUpdates(t *testing.T) {
	s := TestStructure{}
	s.KeyStr = "test:002"
	if Delete(&s) != nil {
		t.Fatal("Couldn't delete!")
	}

	s.Age = 34
	s.LongName = "Test Testerson"
	if Insert(&s) != nil {
		t.Fatal("Couldn't save!")
	}

	// Apply a small set of changes.
	err := UpdateWithChanges(&s, map[string]interface{}{"age": 112})
	if err != nil {
		t.Fatal("Got an error while UpdatingWithChanges.")
	}

	if s.Age != 112 {
		t.Fatal("The original struct didn't get updated with the changes.")
	}

	g := TestStructure{}
	g.KeyStr = "test:002"
	Load(&g)

	if g.Age != 112 {
		t.Fatal("The updates didn't get saved correctly.")
	}

	err = UpdateWithChanges(&g, map[string]interface{}{"error": true})
	if err == nil {
		t.Fatal("Attempted to apply illegal changes, should have been caught.")
	}

}

func TestSaving(t *testing.T) {
	s := TestStructure{}
	s.KeyStr = "test:001"

	// Start by deleting the existing record, to make sure that the
	// database is not contaminated from previous runs.
	err := Delete(&s)
	if err != nil {
		t.Fatalf("Failed to delete: %v", err)
	}

	// Now load the value from the key. Since the value was deleted,
	// it should just load the defaults.
	s.MakeDefault()

	if s.LongName != "Johnny Appleseed" {
		t.Errorf("Name not set to default (was %s)", s.LongName)
	}

	s.LongName = "Johnny Carson"
	s.Age = 99
	err = Insert(&s)
	if err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	// Now we'll load the data back from redis into another part of
	// memory, and compare that the updates were preserved.
	g := TestStructure{}
	g.KeyStr = "test:001"
	err = Load(&g)

	if err != nil {
		t.Fatalf("Failed to load: %v", err)
	}
	if g.LongName != "Johnny Carson" {
		t.Fatalf("Didn't save name correctly (retrieved %s)", g.LongName)
	}
	if g.Age != 99 {
		t.Fatalf("Didn't save the age correctly (retrieved %d)", g.Age)
	}

	// Check if the model validation works correctly.
	g.Error = true
	err = Save(&g)
	if err == nil {
		t.Fatalf("Should have raised a validation error, but didn't.")
	}
}
