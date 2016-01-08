package main

import (
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"log"
	"net"
	"os"
	"time"
)

var targetIP = ""
var ipaddr = ""

func init() {
	if len(os.Args) > 1 {
		targetIP = os.Args[1]
		fmt.Println("Sending ICMP echo request to ", targetIP)
	} else {
		targetIP = "8.8.8.8"
		fmt.Println("Sending ICMP echo request to ", targetIP)
	}
}

func checkError(err error, message string) {
	if err != nil {
		logMessage := message + " " + err.Error()
		log.Fatalf(logMessage)
	}
}

func main() {
	timeout := 5
	t1 := time.Now()
	c, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		checkError(err, "listen err, ")
	}
	c.SetDeadline(t1.Add(time.Second * time.Duration(timeout)))
	defer c.Close()

	wm := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: 1,
			Data: []byte("HELLO-R-U-THERE"),
		},
	}
	wb, err := wm.Marshal(nil)
	if err != nil {
		log.Fatal(err)
	}
	// Send ICMP echo and start timer
	starttime := time.Now()
	if _, err := c.WriteTo(wb, &net.IPAddr{IP: net.ParseIP(targetIP)}); err != nil {
		log.Fatalf("WriteTo err, %s", err)
	}

	rb := make([]byte, 1500)
	// Look for ICMP response and stop time
	n, peer, err := c.ReadFrom(rb)
	stoptime := time.Now()
	if err != nil {
		log.Fatal(err)
	}
	//rm, err := icmp.ParseMessage(iana.ProtocolICMP, rb[:n])
	rm, err := icmp.ParseMessage(1, rb[:n])
	if err != nil {
		log.Fatal(err)
	}
	switch rm.Type {
	case ipv4.ICMPTypeEchoReply:
		log.Printf("got reflection from %v", peer)
		log.Printf("ICMP echo request sent on ", starttime)
		log.Printf("ICMP echo response received on ", stoptime)
	default:
		log.Printf("got %+v; want echo reply", rm)
	}
}
