package main

import (
	"fmt"
	"net"
	"strconv"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

// MaxAttemptsGetPort is the number of attempt authorized to get a random port
const MaxAttemptsGetPort = 5

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
	content []byte
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

// GetReverseAddrs returns a map of dstAddr => srcAddr.
// This way it's easy to get the destination of the packet coming back.
// Like a reverse proxy you know?
func (t *TranslationTable) GetReverseAddrs() map[string]string {
	var result map[string]string
	for _, entry := range t.Entries {
		result[entry.DstAddr()] = entry.SrcAddr()
	}
	return result
}

// RemoveEntry removes an entry that has been processed.
// As the entries gets randomly a port associated, it's more than enough
// to differentiate it.
// WARNING: If some issue appear we might need to add the DstAddr "just in case"
func (t *TranslationTable) RemoveEntry(natPort int) {
	for index, entry := range t.Entries {
		if entry.natPort == natPort {
			t.Entries[index] = t.Entries[len(t.Entries)-1]
			t.Entries = t.Entries[:len(t.Entries)-1]
		}
	}
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
// attempt is the number of attempts of this function, reaching a specific number will stop
// the function.
func getRandomPort(attempt, maxAttempt int) int {
	if attempt >= maxAttempt {
		panic("[ERROR] getRandomPort: Too many attempts without getting a port.")
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Println("[ERROR]: Could not get a random port... Retrying...")
		getRandomPort(attempt+1, maxAttempt)
	}
	defer listener.Close()

	return listener.Addr().(*net.TCPAddr).Port
}

func newTranslationTableEntry(srcAddr string, srcPort int, dstAddr string, dstPort int, natPort int, content []byte) *TranslationTableEntry {
	return &TranslationTableEntry{
		dstAddr: dstAddr,
		dstPort: dstPort,
		srcAddr: srcAddr,
		srcPort: srcPort,
		natPort: natPort,
		content: content,
	}
}

// ToString returns the string reprensenting this Entry
func (t *TranslationTableEntry) ToString() string {
	return "entry"
}

// Translate will generate a packet  to send to the outside world
// using the nat port as source port
func (t *TranslationTableEntry) translate() []byte {
	ipLayer := &layers.IPv4{
		SrcIP: net.IP{127, 0, 0, 1},
		DstIP: net.ParseIP(t.dstAddr),
	}
	// ethernetLayer := &layers.Ethernet{
	// 	SrcMAC: net.HardwareAddr{0xFF, 0xAA, 0xFA, 0xAA, 0xFF, 0xAA},
	// 	DstMAC: net.HardwareAddr{0xBD, 0xBD, 0xBD, 0xBD, 0xBD, 0xBD},
	// }
	tcpLayer := &layers.TCP{
		SrcPort: layers.TCPPort(t.natPort),
		DstPort: layers.TCPPort(t.dstPort),
	}
	// And create the packet with the layers
	buffer := gopacket.NewSerializeBuffer()
	gopacket.SerializeLayers(buffer, gopacket.SerializeOptions{},
		// ethernetLayer,
		ipLayer,
		tcpLayer,
		gopacket.Payload(t.content),
	)
	return buffer.Bytes()
}

// SrcAddr returns the source address of a translation table entry.
// Example: 7.7.7.7:56
func (t *TranslationTableEntry) SrcAddr() string {
	return fmt.Sprintf("%s:%d", t.srcAddr, t.srcPort)
}

// NatAddr returns the address that the nat translated to
func (t *TranslationTableEntry) NatAddr() string {
	return fmt.Sprintf("%s:%d", "127.0.0.1", t.natPort)
}

// DstAddr returns the destination address of this entry
func (t *TranslationTableEntry) DstAddr() string {
	return fmt.Sprintf("%s:%d", t.dstAddr, t.dstPort)
}

// SendTCP will actually send the  "TableEntry" to the remote server.
// It internally uses translate to modify the the IP and TCP layers
func (t *TranslationTableEntry) SendTCP() error {

	translatedPayload := t.translate()

	remoteTCPAddr, err := net.ResolveTCPAddr("tcp", t.DstAddr())
	if err != nil {
		panic(err)
	}
	localTCPAddr, err := net.ResolveTCPAddr("tcp", t.NatAddr())
	if err != nil {
		panic(err)
	}
	packetConn, err := net.DialTCP("tcp", localTCPAddr, remoteTCPAddr)
	if err != nil {
		panic(err)
	}

	bytesWrote, err := packetConn.Write(translatedPayload)
	fmt.Println(fmt.Sprintf("packet of length %d sent to %s!\n", bytesWrote, remoteTCPAddr.String()))
	return nil
}

func main() {
	translationTable := NewTranslationTable()
	handle, err := pcap.OpenLive("wlp2s0", 1600, true, pcap.BlockForever)
	if err != nil {
		panic(err)
	}
	defer handle.Close()
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		var networkLayer gopacket.NetworkLayer
		var transportLayer gopacket.TransportLayer
		networkLayer = packet.NetworkLayer()
		transportLayer = packet.TransportLayer()
		networkFlow := networkLayer.NetworkFlow()
		transportFlow := transportLayer.TransportFlow()
		srcPort, _ := strconv.Atoi(transportFlow.Src().String())
		dstPort, _ := strconv.Atoi(transportFlow.Dst().String())

		translationTableEntry := newTranslationTableEntry(networkFlow.Src().String(), srcPort, networkFlow.Dst().String(), dstPort, 0, packet.Data())
		if translationTable.GetReverseAddrs()[translationTableEntry.SrcAddr()] != "" {
			natPort := translationTableEntry.dstPort
			translationTableEntry.dstAddr = translationTable.GetReverseAddrs()[translationTableEntry.SrcAddr()]
			translationTableEntry.SendTCP()
			translationTable.RemoveEntry(natPort)

		} else {
			randomPort := getRandomPort(0, MaxAttemptsGetPort)
			translationTableEntry.natPort = randomPort
			translationTable.AddEntry(*translationTableEntry)
			translationTable.Print()
			translationTableEntry.SendTCP()
		}
	}
}
