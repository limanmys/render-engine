package helpers

import (
	"reflect"
	"testing"
)

func TestStringBetween(t *testing.T) {
	expected := "e"
	received := StringBetween("deneme", "n", "m")
	if expected != received {
		t.Errorf("Expected %s , received %s", expected, received)
	}

	expected = ""
	received = StringBetween("deneme", "a", "b")
	if expected != received {
		t.Errorf("Expected %s , received %s", expected, received)
	}

	expected = ""
	received = StringBetween("deneme", "d", "b")
	if expected != received {
		t.Errorf("Expected %s , received %s", expected, received)
	}

	expected = ""
	received = StringBetween("denem", "m", "m")
	if expected != received {
		t.Errorf("Expected %s , received %s", expected, received)
	}
}

func TestStringBefore(t *testing.T) {
	expected := "de"
	received := StringBefore("deneme", "n")
	if expected != received {
		t.Errorf("Expected %s , received %s", expected, received)
	}

	expected = ""
	received = StringBefore("deneme", "a")
	if expected != received {
		t.Errorf("Expected %s , received %s", expected, received)
	}
}

func TestStringAfter(t *testing.T) {
	expected := "eme"
	received := StringAfter("deneme", "n")
	if expected != received {
		t.Errorf("Expected %s , received %s", expected, received)
	}

	expected = ""
	received = StringAfter("deneme", "a")
	if expected != received {
		t.Errorf("Expected %s , received %s", expected, received)
	}

	expected = ""
	received = StringAfter("denem", "m")
	if expected != received {
		t.Errorf("Expected %s , received %s", expected, received)
	}
}

func TestContains(t *testing.T) {
	values := []string{"a", "b", "c"}
	received := Contains(values, "b")
	if received != true {
		t.Errorf("Expected %v , received %v", true, received)
	}

	received = Contains(values, "e")
	if received != false {
		t.Errorf("Expected %v , received %v", false, received)
	}
}

func TestUniqueStrings(t *testing.T) {
	expected := []string{"a", "b", "c"}
	values := []string{"a", "b", "c", "b"}
	received := UniqueStrings(values)
	if !reflect.DeepEqual(received, expected) {
		t.Errorf("Expected %v , received %v", expected, received)
	}
}

func TestIsValidUUID(t *testing.T) {
	expected := false
	received := IsValidUUID("deneme")
	if expected != received {
		t.Errorf("Expected %v , received %v", expected, received)
	}

	expected = true
	received = IsValidUUID("56dd3724-83ed-4274-a222-71047ee2557d")
	if expected != received {
		t.Errorf("Expected %v , received %v", expected, received)
	}
}

func TestEncodeMessageUTF16(t *testing.T) {
	expected := "006D006500720074"
	received := EncodeMessageUTF16("mert")
	if expected != received {
		t.Errorf("Expected %v , received %v", expected, received)
	}
}
