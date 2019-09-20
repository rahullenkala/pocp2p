package main

import (
	"context"

	"flag"
	
	"net"

	"fmt"

	"github.com/libp2p/go-libp2p"

	crypto "github.com/libp2p/go-libp2p-crypto"

	dht "github.com/libp2p/go-libp2p-kad-dht"

	"github.com/multiformats/go-multiaddr"

	mrand "math/rand"

	"os"
)

var ipaddrss="null"
func main() {

	help := flag.Bool("help", false, "Display Help")

	listenHost := flag.String("host", "0.0.0.0", "The bootstrap node host listen address\n")

	port := flag.Int("port", 4001, "The bootstrap node listen port")

	flag.Parse()

	if *help {

		fmt.Printf("This is a simple bootstrap node for kad-dht application using libp2p\n\n")

		fmt.Printf("Usage: \n   Run './bootnode'\nor Run './bootnode -host [host] -port [port]'\n")

		os.Exit(0)

	}
	addrs,err:=net.InterfaceAddrs()
	if err!=nil{
		os.Stderr.WriteString("Oops:"+err.Error()+"\n")
		os.Exit(1)
	}
	
	for _,a:=range addrs{
		if ipnet,ok:=a.(*net.IPNet);ok&&!ipnet.IP.IsLoopback(){
			if ipnet.IP.To4()!=nil{
				ipaddrss=ipnet.IP.String()
			}
		}
	}
	fmt.Printf("[*] Listening on: %s with port: %d\n", ipaddrss, *port)

	ctx := context.Background()

	r := mrand.New(mrand.NewSource(int64(*port)))

	// Creates a new RSA key pair for this host.

	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)

	if err != nil {

		panic(err)

	}

	// 0.0.0.0 will listen on any interface device.

	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d",ipaddrss, *port))

	// libp2p.New constructs a new libp2p Host.

	// Other options can be added here.

	host, err := libp2p.New(

		ctx,
		
		libp2p.NATPortMap(),

		libp2p.ListenAddrs(sourceMultiAddr),

		libp2p.Identity(prvKey),
		
	)

	if err != nil {

		panic(err)

	}

	_, err = dht.New(ctx, host)

	if err != nil {

		panic(err)

	}

	fmt.Println("")

	fmt.Printf("[*] Your Bootstrap ID Is: /ip4/%s/tcp/%v/p2p/%s\n", ipaddrss, *port, host.ID().Pretty())

	fmt.Println("")

	select {}

}
