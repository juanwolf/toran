package main

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

func TestNewTranslationTableEntry(t *testing.T) {
	srcAddr := "127.0.0.1"
	srcPort := 1999
	dstAddr := "255.255.255.255"
	dstPort := 25999
	natPort := 19999
	entry := newTranslationTableEntry(srcAddr, srcPort, dstAddr, dstPort, natPort)
	if entry.srcAddr != srcAddr || entry.srcPort != srcPort || entry.dstPort != entry.dstPort || entry.dstAddr != dstAddr {
		t.Error("newTranslationTableEntry did not fill properly all the fields")
	}
}

func TestAddEntry(t *testing.T) {
	table := NewTranslationTable()
	srcAddr := "127.0.0.1"
	srcPort := 1999
	dstAddr := "255.255.255.255"
	dstPort := 25999
	natPort := 19999
	entry := newTranslationTableEntry(srcAddr, srcPort, dstAddr, dstPort, natPort)
	table.AddEntry(*entry)

	if len(table.GetEntries()) != 1 {
		t.Error("AddEntry did not add properly the entry in the translation table.")
	}

	table.Print()

}
