package flagt

import (
	"flag"
	"fmt"

	"strings"

	maddr "github.com/multiformats/go-multiaddr"
)

// A new type we need for writing a custom flag parser

type addrList []maddr.Multiaddr

func (al *addrList) String() string {

	strs := make([]string, len(*al))

	for i, addr := range *al {

		strs[i] = addr.String()

	}

	return strings.Join(strs, ",")

}

func (al *addrList) Set(value string) error {

	addr, err := maddr.NewMultiaddr(value)

	if err != nil {

		return err

	}

	*al = append(*al, addr)

	return nil

}

func StringsToAddrs(addrStrings []string) (maddrs []maddr.Multiaddr, err error) {

	for _, addrString := range addrStrings {

		addr, err := maddr.NewMultiaddr(addrString)

		if err != nil {

			return maddrs, err

		}

		maddrs = append(maddrs, addr)

	}

	return

}

type Config struct {
	RendezvousString string

	BootstrapPeers addrList

	ListenAddresses addrList

	ProtocolID string
}

func ParseFlags() (Config, error) {

	config := Config{}

	flag.StringVar(&config.RendezvousString, "rendezvous", "meet me here",

		"Unique string to identify group of nodes. Share this with your friends to let them connect with you")

	flag.Var(&config.BootstrapPeers, "peer", "Adds a peer multiaddress to the bootstrap list")

	flag.Var(&config.ListenAddresses, "listen", "Adds a multiaddress to the listen list")

	flag.StringVar(&config.ProtocolID, "pid", "/chat/1.1.0", "Sets a protocol id for stream headers")

	flag.Parse()

	if len(config.BootstrapPeers) == 0 {

		//config.BootstrapPeers = dht.DefaultBootstrapPeers
<<<<<<< HEAD
		addr, err := maddr.NewMultiaddr("/ip4/157.245.238.65/tcp/4001/p2p/QmP2C45o2vZfy1JXWFZDUEzrQCigMtd4r3nesvArV8dFKd")
=======
		addr, err := maddr.NewMultiaddr("/ip4/192.168.1.19/tcp/4321/p2p/QmdSyhb8eR9dDSR5jjnRoTDBwpBCSAjT7WueKJ9cQArYoA")
>>>>>>> 865e3cd5783060d22d1eef96943de0d8e22a7ad6
		fmt.Println(err)
		config.BootstrapPeers = append(config.BootstrapPeers, addr)
		//
		fmt.Println(addr)
		fmt.Println(config.BootstrapPeers)
	}

	return config, nil

}
