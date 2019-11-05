package main

import (
	"bufio"

	"context"

	"flag"

	"fmt"

	"sync"

	proto "github.com/golang/protobuf/proto"

	"github.com/libp2p/go-libp2p"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"

	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/libp2p/go-libp2p-core/protocol"

	discovery "github.com/libp2p/go-libp2p-discovery"

	dht "github.com/libp2p/go-libp2p-kad-dht"

	multiaddr "github.com/multiformats/go-multiaddr"

	"github.com/ipfs/go-log"

	mp "./mylib"

	"./flagt"
)

//var walletmap = make(map[string]string)
var done = make(chan []byte)
var walletmap4 *mp.WalletMap
var recievedmap1 *mp.WalletMap
var host1 host.Host
var hashOfmap string

//var m = make(map[string]string)

var logger = log.Logger("rendezvous")

func sendResponse(str string) {
	//fmt.Println("Inresponse")
	//h := sha256.New()
	//h.Write([]byte(walletmap4.Data))
	msg := Kipcntxtmessage{

		Type:      Kipcntxtmessage_RESPONSE,
		Data:      hashOfmap,
		Sid:       host1.ID().String(),
		Delimiter: "|",
	}
	data, err := proto.Marshal(&msg)
	if err != nil {
		logger.Error("Marshaling error")
		fmt.Println(err)
	}
	ctx := context.Background()
	//config, err := flagt.ParseFlags()
	fmt.Println(str)
	id, err1 := peer.IDB58Decode(str)
	if err1 != nil {
		fmt.Println("panic !!") //go readData(rw, done)
		panic(err1)
	}
	//fmt.Print("Up stream")
	//fmt.Print(host1)
	stream1, err := host1.NewStream(ctx, id, protocol.ID("/chat/1.1.0"))
	if err != nil {
		panic(err)
	}

	rw := bufio.NewReadWriter(bufio.NewReader(stream1), bufio.NewWriter(stream1))
	writeData(rw, data)
	fmt.Println("Response sent to:", str)
}
func sendQuery(rw *bufio.ReadWriter, id string) {
	fmt.Println("Sending Query")
	//rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
	msg := Kipcntxtmessage{

		Type:      Kipcntxtmessage_QUERY,
		Data:      "",
		Sid:       id,
		Delimiter: "|",
	}
	data, err := proto.Marshal(&msg)
	if err != nil {
		logger.Error("Marshaling error")
		fmt.Println(err)
	}
	//fmt.Println(data)
	writeData(rw, data)

	//go readData(rw, done)

}
func processmessage(rw *bufio.ReadWriter, msg []byte) {
	recieved := &Kipcntxtmessage{} //Create a template of message
	proto.Unmarshal(msg, recieved) //unmarshal the message

	var MsgType = recieved.GetType() //Identify the type of message
	var senderid = recieved.GetSid()

	fmt.Println("Processing the message")
	if MsgType == Kipcntxtmessage_QUERY {
		fmt.Println("Its a query from Node: ", senderid)
		fmt.Println("Intiating Response ")
		sendResponse(senderid)

	} else {
		fmt.Println("It is Response message from Node:", senderid)
		fmt.Println("Adding Response data to local Map")
		var temphash = recieved.GetData()
		recievedmap1.Put(temphash, recieved.Sid)
		//walletmap4.PutEntries(MapData)
		fmt.Println("Recieved Hash", temphash)

	}

	//fmt.Println("hello", map1.Data)

}
func handleStream(stream network.Stream) {

	logger.Info("Got a new stream!")

	// Create a buffer stream for non blocking read and write.

	rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
	//rw.

	go readData(rw, done)

}

func readData(rw *bufio.ReadWriter, done chan []byte) {
	fmt.Println("\x1b[32m", "New Message recieved")
	for {
		//rw.read

		//var buf []byte
		data, err := rw.ReadBytes(124)

		if err != nil {

			fmt.Println("Error reading from buffer")

			panic(err)
		}

		//fmt.Println(data)
		processmessage(rw, data)

		//done <- data

		fmt.Print("Outstide read function")
	}
}

func writeData(rw *bufio.ReadWriter, data []byte) {

	//stdReader := bufio.NewReader(os.Stdin)
	//fmt.Println(data)
	_, err := rw.Write(data)

	if err != nil {

		fmt.Println("Error writing to buffer")

		panic(err)

	}
	//fmt.Println("No.of bytes:", leng)
	err = rw.Flush()

	if err != nil {

		fmt.Println("Error flushing buffer")

		panic(err)

	}

}

func main() {
	walletmap4 = mp.New()
	recievedmap1 = mp.New()
	walletmap4.Put("WalletId:21", "Address:21")
	walletmap4.Put("WalletId:22", "Address:123")
	walletmap4.Put("walletId:2345", "Address:1234")
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

	host1, err = libp2p.New(ctx, libp2p.NATPortMap(),
		libp2p.ListenAddrs([]multiaddr.Multiaddr(config.ListenAddresses)...),
	)
	//fmt.Print(host1)
	if err != nil {

		panic(err)

	}

	logger.Info("Host created. We are:", host1.ID())

	logger.Info(host1.Addrs())

	// Set a function as stream handler. This function is called when a peer

	// initiates a connection and starts a stream with this peer.

	host1.SetStreamHandler(protocol.ID(config.ProtocolID), handleStream)

	// Start a DHT, for use ipeern peer discovery. We can't just make a new DHT

	// client because we want each peer to maintain its own local copy of the

	// DHT, so that the bootstrapping node of the DHT can go down without

	// inhibiting future peer discovery.

	kademliaDHT, err := dht.New(ctx, host1)

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
	//pubs, err1 := ps.NewFloodSub(ctx, host)
	/*
		if err1 != nil {
			panic(err1)
		}
	*/
	//fmt.Println(kademliaDHT.RoutingTable().ListPeers())

	var wg sync.WaitGroup

	for _, peerAddr := range config.BootstrapPeers {

		peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)

		wg.Add(1)

		go func() {

			defer wg.Done()

			if err := host1.Connect(ctx, *peerinfo); err != nil {

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
	//fmt.Println("AFTER searching")
	if err != nil {

		panic(err)

	}

	for peer := range peerChan {
		//fmt.Println("hello rahul")
		if peer.ID == host1.ID() {

			continue

		}
		if len(peer.Addrs) <= 0 {

			continue
		}
		logger.Info("Found peer:", peer)

		logger.Info("Connecting to:", peer)

		err4 := host1.Connect(ctx, peer) //we connect to other peer

		if err4 != nil {
			fmt.Println("Connection not possible", peer)
			continue
		}

		logger.Info("Connected to:", peer)
		//fmt.Println(kademliaDHT)

	}
	//	fmt.Print(kademliaDHT.GetPublicKey(ctx, "QmdSyhb8eR9dDSR5jjnRoTDBwpBCSAjT7WueKJ9cQArYoA"))
	//if flagpeer == 0 {
	//	goto search
	//}
	//var wg1 sync.WaitGroup
	//fmt.Println(string(kademliaDHT.PutValue(ctx, "QmdSyhb8eR9dDSR5jjnRoTDBwpBCSAjT7WueKJ9cQArYoA", []byte(""))))
	var ids = kademliaDHT.RoutingTable().ListPeers()
	for x := range ids {
		//sfmt.Println(ids)
		//fmt.Println(ids[x])
		//fmt.Print(kademliaDHT.GetValue(ctx, string(ids[x])))
		stream, err := host1.NewStream(ctx, ids[x], protocol.ID(config.ProtocolID))
		if err != nil {
			logger.Warning("Connection failed:", err)
			continue
		} else {

			rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
			fmt.Println("Sending Query to :", ids[x])
			sendQuery(rw, host1.ID().String())

		}
	}

	fmt.Println("Wallet Data", walletmap4.Data)

	select {}

}
