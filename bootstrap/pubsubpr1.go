package main

import (
	"context"

	"flag"

	"fmt"

	"sync"

	"github.com/libp2p/go-libp2p"

	"github.com/libp2p/go-libp2p-core/peer"

	discovery "github.com/libp2p/go-libp2p-discovery"

	dht "github.com/libp2p/go-libp2p-kad-dht"

	multiaddr "github.com/multiformats/go-multiaddr"

	"github.com/ipfs/go-log"

	"practice/flagt"

	ps "github.com/libp2p/go-libp2p-pubsub"
)

var logger = log.Logger("rendezvous")

func readData(cont context.Context, subs *ps.Subscription) {

	for {

		msg, err := subs.Next(cont)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(msg.GetFrom()) + string(msg.Data))

	}

}

func publishdata(topic string, pub *ps.PubSub, data string) {

	err := pub.Publish(topic, []byte(data))

	if err != nil {
		panic(err)
	}
	fmt.Println("Data Published")

}
func main() {

	//log.SetAllLoggers(logging.WARNING)

	log.SetLogLevel("rendezvous", "info")

	help := flag.Bool("h", false, "Display Help")

	config, err := flagt.ParseFlags()

	if err != nil {

		panic(err)

	}

	if *help {

		fmt.Println("This program demonstrates a simple p2p chat application using libp2p")

		fmt.Println()

		fmt.Println("Usage: Run './chat in two different terminals. Let them connect to the bootstrap nodes, announce themselves and connect to the peers")

		flag.PrintDefaults()

		return

	}

	ctx := context.Background()

	// libp2p.New constructs a new libp2p Host. Other options can be added

	// here.

	host, err := libp2p.New(ctx, libp2p.NATPortMap(),
		libp2p.ListenAddrs([]multiaddr.Multiaddr(config.ListenAddresses)...),
	)

	if err != nil {

		panic(err)

	}

	logger.Info("Host created. We are:", host.ID())

	logger.Info(host.Addrs())

	// Set a function as stream handler. This function is called when a peer

	// initiates a connection and starts a stream with this peer.

	//host.SetStreamHandler(protocol.ID(config.ProtocolID), handleStream)
	pubs, err1 := ps.NewFloodSub(ctx, host)

	if err1 != nil {
		panic(err1)
	}
	// Start a DHT, for use in peer discovery. We can't just make a new DHT

	// client because we want each peer to maintain its own local copy of the

	// DHT, so that the bootstrapping node of the DHT can go down without

	// inhibiting future peer discovery.

	kademliaDHT, err := dht.New(ctx, host)

	if err != nil {

		panic(err)

	}

	// Bootstrap the DHT. In the default configuration, this spawns a Background

	// thread that will refresh the peer table every five minutes.

	logger.Debug("Bootstrapping the DHT")

	if err = kademliaDHT.Bootstrap(ctx); err != nil {

		panic(err)

	}

	// Let's connect to the bootstrap nodes first. They will tell us about the

	// other nodes in the network.
	fmt.Println()

	//fmt.Println(kademliaDHT.RoutingTable().ListPeers())

	var wg sync.WaitGroup

	for _, peerAddr := range config.BootstrapPeers {

		peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)

		wg.Add(1)

		go func() {

			defer wg.Done()

			if err := host.Connect(ctx, *peerinfo); err != nil {

				logger.Warning(err)

			} else {

				logger.Info("Connection established with bootstrap node:", *peerinfo)

			}

		}()

	}

	wg.Wait()

	// We use a rendezvous point "meet me here" to announce our location.

	// This is like telling your friends to meet you at the Eiffel Tower.
	fmt.Println(kademliaDHT.RoutingTable().ListPeers())

	logger.Info("Announcing ourselves...")

	routingDiscovery := discovery.NewRoutingDiscovery(kademliaDHT)

	discovery.Advertise(ctx, routingDiscovery, config.RendezvousString)

	logger.Info("Successfully announced!")

	// Now, look for others who have announced

	// This is like your friend telling you the location to meet you.
	//search:
	logger.Info("Searching for other peers...")

	peerChan, err := routingDiscovery.FindPeers(ctx, config.RendezvousString)
	fmt.Println("AFTER searching")
	if err != nil {

		panic(err)

	}

	for peer := range peerChan {
		fmt.Println("hello rahul")
		if peer.ID == host.ID() {

			continue

		}
		if len(peer.Addrs) <= 0 {

			continue
		}
		logger.Info("Found peer:", peer)

		logger.Info("Connecting to:", peer)

		//stream, err := host.NewStream(ctx, peer.ID, protocol.ID(config.ProtocolID))

		if err != nil {
			fmt.Println("log")
			logger.Warning("Connection failed:", err)

			continue

		} else {
			fmt.Println("Before Sub")

			err4 := host.Connect(ctx, peer)

			if err4 != nil {
				panic(err4)
			}

		}

		/*rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

		go writeData(rw)

		go readData(rw)

		//flagpeer = 1
		*/
		logger.Info("Connected to:", peer)
	}

	//fmt.Println(kademliaDHT)

	subs, err2 := pubs.Subscribe("New Topic")
	if err2 != nil {
		panic(err2)
	}
	//time.sleep(500)
	go readData(ctx, subs)

	/*	for i := 0; i < 10; i++ {
		publishdata("New Topic", pubs, string(i))

	}*/
	select {}

}
