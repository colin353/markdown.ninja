package models

import "testing"

func TestDomainValidation(t *testing.T) {
	d := Domain{}
	d.ExternalDomain = "testing.structures.net"
	d.InternalDomain = "test"

	err := Insert(&d)
	if err != nil {
		t.Fatalf("Should be able to insert valid domain link.")
	}

	iterator, err := GetList(&d)
	if iterator.Count() != 1 {
		t.Fatalf("Should be able to see the domain in the list.")
	}

	iterator.Next()
	dom := iterator.Value().(*Domain)
	if dom.InternalDomain != "test" || dom.ExternalDomain != "testing.structures.net" {
		t.Fatalf("Did not get the correct domain loaded from list.")
	}

	d = Domain{}
	d.ExternalDomain = "invalid$domain"
	d.InternalDomain = "test"

	err = Insert(&d)
	if err == nil {
		t.Fatalf("Should not be able to insert invalid domain.")
	}
}
