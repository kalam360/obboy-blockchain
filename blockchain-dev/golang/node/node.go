package node

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/kalam360/obboy-blockchain-golang/database"
)

const DefaultMiningDifficulty = 3
const HttpSSLPort = 443

type PeerNode struct {
	IP          string         `json:"ip"`
	Port        uint64         `json:"port"`
	IsBootstrap bool           `json:"is_bootstrap"`
	Account     common.Address `json:"account"`
	NodeVersion string         `json:"node_version"`

	connected bool
}

func (pn PeerNode) TcpAddress() string {
	return fmt.Sprintf("%s:%d", pn.IP, pn.Port)
}

func (pn PeerNode) ApiProtocol() string {
	if pn.Port == HttpSSLPort {
		return "https"
	}

	return "http"
}

type Node struct {
	dataDir string
	info    PeerNode

	state *database.State

	pendingState *database.State

	knownPeers      map[string]PeerNode
	pendingTXs      map[string]database.SignedTx
	archivedTXs     map[string]database.SignedTx
	newSyncedBlocks chan database.Block
	newPendingTXs   chan database.SignedTx
	nodeVersion     string

	miningDifficulty uint
	isMining         bool
}

func New(dataDir string, ip string, port uint64, acc common.Address, bootstrap PeerNode, version string, miningDifficulty uint) *Node {
	knownPeers := make(map[string]PeerNode)

	n := &Node{
		dataDir:          dataDir,
		info:             NewPeerNode(ip, port, false, acc, true, version),
		knownPeers:       knownPeers,
		pendingTXs:       make(map[string]database.SignedTx),
		archivedTXs:      make(map[string]database.SignedTx),
		newSyncedBlocks:  make(chan database.Block),
		newPendingTXs:    make(chan database.SignedTx, 10000),
		nodeVersion:      version,
		isMining:         false,
		miningDifficulty: miningDifficulty,
	}

	n.AddPeer(bootstrap)

	return n
}

func NewPeerNode(ip string, port uint64, isBootstrap bool, acc common.Address, connected bool, version string) PeerNode {
	return PeerNode{ip, port, isBootstrap, acc, version, connected}
}

func (n *Node) Run(ctx context.Context, isSSLDisabled bool, sslEmail string) error {

}

func (n *Node) AddPeer(peer PeerNode) {
	n.knownPeers[peer.TcpAddress()] = peer
}
