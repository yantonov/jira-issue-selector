package ui

import (
	"testing"
)

func TestNormalizeEmptyString(t *testing.T) {
	if normalizeTaskName("") != "" {
		t.Errorf("Empty string is expected")
	}
}

func TestInvalidCharactersAreRemove(t *testing.T) {
	if normalizeTaskName("abc!@#def") != "abcdef" {
		t.Errorf("Characters except a-z, A-Z, 0-9, _, space will be removed")
	}
}

func TestSpaceReplacement(t *testing.T) {
	if normalizeTaskName("   abc   def   ") != "abc_def" {
		t.Errorf("Whitespaces will be replaced by underscore")
	}
}

func TestMaximumSize(t *testing.T) {
	expected := "0123456789012345678901234567890123456789012345678901234567890123456789"
	input := "01234567890123456789012345678901234567890123456789012345678901234567890123456789"
	if normalizeTaskName(input) != expected {
		t.Errorf("Whitespaces will be replaced by underscore")
	}
}
