package main

import (
	"context"
	"fmt"

	ps "github.com/libp2p/go-floodsub"
	"github.com/libp2p/go-libp2p"
	//dht "github.com/libp2p/go-libp2p-kad-dht"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background()) // Initialize context

	defer cancel() // Cancel

	host, err := libp2p.New(
		ctx,
		libp2p.NATPortMap(),
		libp2p.ListenAddrStrings(
			"/ip4/0.0.0.0/tcp/1111",
			"/ip6/::1/tcp/1111",
		),
	)
	if err != nil { // Check for errors
		panic(err) // Panic
	}

	//dht, err := dht.New(ctx, host) // Initialize dht
	//fmt.Println(dht)

	if err != nil { // Check for errors
		panic(err) // Panic
	}
	Frouter := ps.FloodSubRouter{}
	pubs, err1 := ps.NewFloodSub(ctx, host)

	if err1 != nil {
		panic(err1)
	}

	Frouter.Attach(pubs)
	//fmt.Println(host.ID)
	Frouter.Join("New Topic")
	Frouter.Publish(host.ID, []byte("hello"))
	fmt.Println(pubs.ListPeers("New topic"))

	select {}
}
