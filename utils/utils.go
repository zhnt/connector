package utils

import (
	"flag"
	"fmt"
	mrand "math/rand"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/multiformats/go-multiaddr"
)

const (
	// PROTOCAL_FILE_TRANS 传输协议 - 文件
	PROTOCAL_FILE_TRANS protocol.ID = "/trans/file/0.0.1"
	// PROTOCAL_HTTP_PROXY 传输协议 - http
	PROTOCAL_HTTP_PROXY protocol.ID = "/proxy/http/0.0.1"

	PROTOCAL_CONNECTOR protocol.ID = "/connector/1.0.0"
)

type config struct {
	RendezvousString string
	listenHost       string
	listenPort       int
}

func parseFlags() *config {
	c := &config{}

	flag.StringVar(&c.RendezvousString, "rendezvous", "meetme", "Unique string to identify group of nodes. Share this with your friends to let them connect with you")
	flag.StringVar(&c.listenHost, "host", "0.0.0.0", "The bootstrap node host listen address\n")
	flag.IntVar(&c.listenPort, "port", 0, "node listen port (0 pick a random unused port)")

	flag.Parse()
	return c
}

func CreateHost() (host.Host, error) {
	return CreateHostWithPort(0)
}

func CreateHostWithPort(port int) (host.Host, error) {
	cfg := parseFlags()

	// 为了调试：生成统一的r，保证host id不变
	r := mrand.New(mrand.NewSource(int64(port)))
	//r = rand.Reader

	// Creates a new RSA key pair for this host.
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		return nil, err
	}

	// 0.0.0.0 will listen on any interface device.
	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", cfg.listenHost, port))

	// libp2p.New constructs a new libp2p Host.
	// Other options can be added here.
	h, err := libp2p.New(
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(prvKey),
	)
	if err != nil {
		return h, err
	}

	for _, la := range h.Addrs() {
		fmt.Printf("[*] Your Multiaddress Is: %v/p2p/%s\n", la, h.ID())
	}

	fmt.Println("create p2p node success")

	return h, nil
}
