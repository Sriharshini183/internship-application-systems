package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

const (
	IPv4 = 1
	IPv6 = 58
)

// Function for handling ipv6 ping requests
func pingipv6(address *net.IPAddr) (time.Duration, error) {
	// creating an icmp Listener connection
	c, err := icmp.ListenPacket("ip6:ipv6-icmp", "::")
	if err != nil {
		return 0, err
	}
	defer c.Close()

	msg := icmp.Message{
		Type: ipv6.ICMPTypeEchoRequest,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: []byte(""),
		},
	}

	checksum, err := msg.Marshal(nil)
	if err != nil {
		return 0, err
	}

	start := time.Now()
	n, err := c.WriteTo(checksum, address)
	if err != nil {
		return 0, err
	} else if n != len(checksum) {
		return 0, fmt.Errorf("Checksum error: got %v; want %v", n, len(checksum))
	}

	reply := make([]byte, 1500)
	err = c.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	if err != nil {
		return 0, err
	}

	_, peer, err := c.ReadFrom(reply)
	if err != nil {
		return 0, err
	}

	rtt := time.Since(start)

	rm, err := icmp.ParseMessage(IPv6, reply)
	if err != nil {
		return 0, err
	}

	if rm.Type == ipv6.ICMPTypeEchoReply {
		return rtt, nil
	}
	return 0, fmt.Errorf("Did not receive echo reply; got %+v from %v", rm, peer)
}

//function for handling ipv4 ping
func pingipv4(address *net.IPAddr) (time.Duration, error) {
	c, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return 0, err
	}
	defer c.Close()

	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: []byte(""),
		},
	}

	checksum, err := msg.Marshal(nil)
	if err != nil {
		return 0, err
	}

	start := time.Now()
	n, err := c.WriteTo(checksum, address)
	if err != nil {
		return 0, err
	} else if n != len(checksum) {
		return 0, fmt.Errorf("got %v; want %v", n, len(checksum))
	}

	reply := make([]byte, 1500)
	err = c.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	if err != nil {
		return 0, err
	}

	_, peer, err := c.ReadFrom(reply)
	if err != nil {
		return 0, err
	}

	rtt := time.Since(start)

	rm, err := icmp.ParseMessage(IPv4, reply)
	if err != nil {
		return 0, err
	}

	if rm.Type == ipv4.ICMPTypeEchoReply {
		return rtt, nil
	}
	return 0, fmt.Errorf("did not recieve echo reply; got %+v from %v", rm, peer)

}

func printStat(ploss int, packetcount int, plosspct float64) {
	fmt.Println("\n !! Process Ping Interrupted !!")
	fmt.Printf("packet loss = %d/%d (%0.2f%%)\n", ploss, packetcount, plosspct)
}

func main() {
	sysargs := os.Args
	if len(sysargs) < 2 {
		fmt.Println("Hostname not provided.! Please provide a hostname as an argument")
		os.Exit(1)
	}
	addr := os.Args[1]

	dst, err := net.ResolveIPAddr("ip:icmp", addr)
	if err != nil {
		panic(err)
	}

	isIPv4 := false
	if len(dst.IP.To4()) == net.IPv4len {
		isIPv4 = true
	}

	fmt.Printf("\n%30s %25s  %20s\n", addr, "RTT", "Packet Loss")

	ploss := 0
	packetcount := 0
	plosspct := 0.0

	// To print final statistics when interrupted with ctrl + c
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			printStat(ploss, packetcount, plosspct)
			os.Exit(0)
		}
	}()

	for {
		var rtt time.Duration
		var err error
		if isIPv4 {
			rtt, err = pingipv4(dst)
		} else {
			rtt, err = pingipv6(dst)
		}
		if err != nil {
			ploss++
			packetcount++
			plosspct = float64(ploss) / float64(packetcount) * 100
			fmt.Printf("%30s: %30s  %10.2f%%\n", dst, err, plosspct)
		} else {
			packetcount++
			plosspct = float64(ploss) / float64(packetcount) * 100
			fmt.Printf("%30s: %30s  %10.2f%%\n", dst, rtt, plosspct)
		}
		time.Sleep(500 * time.Millisecond)
	}
}
