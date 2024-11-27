package multisig

import (
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/txscript"
	"github.com/bxelab/runestone"
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
	OutputAmount	int64
	PkScript		[]byte
	Sigs			[][]byte
	Runes			[]*runestone.Edict
}

type PoolBalance struct {
	Token0Id		runestone.RuneId
	Token1Id		runestone.RuneId
	Token0Amount	uint128.Uint128
	Token1Amount	uint128.Uint128
}
