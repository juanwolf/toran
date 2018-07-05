package toran

import "fmt"

// A TranslationTableEntry is the row in memory for the nat
// storing what incoming connection it had and to where.
// The natPort is a random alloacted port to deal with this incoming connection,
// so there's no clash between two packets reaching the same endpoint
type TranslationTableEntry struct {
	dstAddr string
	dstPort int
	srcAddr string
	srcPort int
	natPort int
}

// TranslationTable is a list of entries
// for the translation table
type TranslationTable struct {
	Entries []TranslationTableEntry
}

func newTranslationTable() *TranslationTable {
	return &TranslationTable{
		Entries: make([]TranslationTableEntry, 0),
	}
}

// GetEntries returns the list of entries of this Translation Table
func (t *TranslationTable) GetEntries() []TranslationTableEntry {
	return t.Entries
}

// Print will returns in the standard output what the translation table is made of.
func (t *TranslationTable) Print() {
	fmt.Println("+-----------------------------------------------------------------+")
	fmt.Println("|     srcAddr     | srcPort |     dstAddr     | dstPort | natPort |")
	fmt.Println("+-----------------------------------------------------------------+")

	for _, entry := range t.Entries {
		fmt.Printf("| %v | %v | %v | %v | %v |\n", entry.srcAddr, entry.srcPort, entry.dstAddr, entry.dstPort, entry.natPort)
	}

	fmt.Println("+-----------------------------------------------------------------+")
}

// getRandomPort returns a random port available on the local machine
func getRandomPort() int {
	return 42
}

func newTranslationTableEntry(dstPort, srcPort int, dstAddr, srcAddr string) *TranslationTableEntry {
	randomPort := getRandomPort()
	return &TranslationTableEntry{
		dstAddr: dstAddr,
		dstPort: dstPort,
		srcAddr: srcAddr,
		srcPort: srcPort,
		natPort: randomPort,
	}
}

// ToString returns the string reprensenting this Entry
func (t *TranslationTableEntry) ToString() string {
	return "entry"
}

func main() {
	translationTable := newTranslationTable()
	translationTable.Print()
}
