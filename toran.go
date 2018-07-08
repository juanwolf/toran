package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"strconv"
)

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

// NewTranslationTable creates an empty
// TranslationTable
func NewTranslationTable() *TranslationTable {
	return &TranslationTable{
		Entries: make([]TranslationTableEntry, 0),
	}
}

// GetEntries returns the list of entries of this Translation Table
func (t *TranslationTable) GetEntries() []TranslationTableEntry {
	return t.Entries
}

// AddEntry add a translationTableEntry to the list of entries of this table
func (t *TranslationTable) AddEntry(entry TranslationTableEntry) {
	t.Entries = append(t.Entries, entry)
}

// Print will returns in the standard output what the translation table is made of.
func (t *TranslationTable) Print() {
	fmt.Println("+-----------------------------------------------------------------+")
	fmt.Println("|     srcAddr     | srcPort |     dstAddr     | dstPort | natPort |")
	fmt.Println("+-----------------------------------------------------------------+")

	for _, entry := range t.Entries {
		fmt.Printf("| %-15v | %-7v | %-15v | %-7v | %-7v |\n", entry.srcAddr, entry.srcPort, entry.dstAddr, entry.dstPort, entry.natPort)
	}

	fmt.Println("+-----------------------------------------------------------------+")
}

// getRandomPort returns a random port available on the local machine
func getRandomPort() int {
	return 42
}

func newTranslationTableEntry(srcAddr string, srcPort int, dstAddr string, dstPort int, natPort int) *TranslationTableEntry {
	return &TranslationTableEntry{
		dstAddr: dstAddr,
		dstPort: dstPort,
		srcAddr: srcAddr,
		srcPort: srcPort,
		natPort: natPort,
	}
}

// ToString returns the string reprensenting this Entry
func (t *TranslationTableEntry) ToString() string {
	return "entry"
}

func main() {
	translationTable := NewTranslationTable()
	handle, err := pcap.OpenLive("lo", 1600, true, pcap.BlockForever)
	if err != nil {
		panic(err)
	}
	defer handle.Close()
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		fmt.Println(packet.String()) // Do something with a packet here.
		var networkLayer gopacket.NetworkLayer
		var transportLayer gopacket.TransportLayer
		networkLayer = packet.NetworkLayer()
		transportLayer = packet.TransportLayer()
		networkFlow := networkLayer.NetworkFlow()
		transportFlow := transportLayer.TransportFlow()
		srcPort, _ := strconv.Atoi(transportFlow.Src().String())
		dstPort, _ := strconv.Atoi(transportFlow.Dst().String())

		translationTableEntry := newTranslationTableEntry(networkFlow.Src().String(), srcPort, networkFlow.Dst().String(), dstPort, 42)
		translationTable.AddEntry(*translationTableEntry)
		translationTable.Print()
	}
}
