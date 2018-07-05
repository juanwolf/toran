package toran

import (
	"testing"
)

func TestNewTranslationTable(t *testing.T) {
	table := NewTranslationTable()
	if table == nil {
		t.Error("NewTranslationTable returned nil instead of a new table")
	}
	if len(table.Entries) > 0 {
		t.Error("NewTranslationTable created a table in the incorrect state")
	}
}
