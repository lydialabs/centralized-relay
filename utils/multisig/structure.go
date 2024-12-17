package multisig

import (
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/txscript"
	"github.com/studyzy/runestone"
	"lukechampine.com/uint128"
)

type MultisigInfo struct {
	PubKeys            [][]byte
	EcPubKeys          []*btcutil.AddressPubKey
	NumberRequiredSigs int
	RecoveryPubKey     []byte
	RecoveryLockTime uint64
}

type MultisigWallet struct {
	TapScriptTree *txscript.IndexedTapScriptTree
	TapLeaves     []txscript.TapLeaf
	PKScript        []byte
	SharedPublicKey *btcec.PublicKey
}

type Input struct {
	TxHash			string
	OutputIdx		uint32
	OutputAmount	int64				`json:"outputAmount"`
	PkScript		[]byte				`json:"pkScript"`
	Sigs			[][]byte
	Runes			[]*runestone.Edict
}

type PoolUTXOBalance struct {
	Token0Id		runestone.RuneId
	Token0Amount	uint128.Uint128
	Token1Id		runestone.RuneId
	Token1Amount	uint128.Uint128
}
