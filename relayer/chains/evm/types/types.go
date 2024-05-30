package types

import (
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type BlockNotification struct {
	Hash   common.Hash
	Height *big.Int
	Header *types.Header
	Logs   []types.Log
}

type Block struct {
	Transactions []string `json:"transactions"`
	GasUsed      string   `json:"gasUsed"`
}

type NonceValue struct {
	Previous *big.Int
	Current  *big.Int
}

type NonceTracker struct {
	address map[common.Address]*NonceValue
	*sync.Mutex
}

type NonceTrackerI interface {
	Get(common.Address) *big.Int
	Set(common.Address, *big.Int)
	Inc(common.Address)
}

// NewNonceTracker
func NewNonceTracker() NonceTrackerI {
	return &NonceTracker{
		address: make(map[common.Address]*NonceValue),
		Mutex:   &sync.Mutex{},
	}
}

func (n *NonceTracker) Get(addr common.Address) *big.Int {
	n.Lock()
	defer n.Unlock()
	nonce := n.address[addr]
	if nonce.Previous == nonce.Current {
		nonce.Current = nonce.Current.Add(nonce.Current, big.NewInt(1))
	}
	return nonce.Current
}

func (n *NonceTracker) Set(addr common.Address, nonce *big.Int) {
	n.Lock()
	defer n.Unlock()
	n.address[addr] = &NonceValue{
		Previous: nonce.Sub(nonce, big.NewInt(1)),
		Current:  nonce,
	}
}

func (n *NonceTracker) Inc(addr common.Address) {
	n.Lock()
	defer n.Unlock()
	nonce := n.address[addr]
	n.address[addr].Current = nonce.Current.Add(nonce.Current, big.NewInt(1))
}
