package ghs

import (
	"testing"
)

func TestReadIntPrefix(t *testing.T) {
	var length, remainder, err = readIntPrefix("2something")
	if err != nil {
		t.Error("Received error during test")
	}
	if length != 2 {
		t.Error("Expected length 2, got ", length)
	}
	if remainder != "something" {
		t.Error("Expected remainder 'something', got ", remainder)
	}
}

func TestExtractName(t *testing.T) {
	var name, remainder, err = extractName("9className")
	if err != nil {
		t.Error("Received error during test")
	}
	if name != "className" {
		t.Error("Expected name 'className', got ", name)
	}
	if remainder != "" {
		t.Error("Expected empty remainder, got ", remainder)
	}
}
