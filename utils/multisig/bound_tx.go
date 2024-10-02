package multisig

import (
	"bytes"
	"fmt"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/bxelab/runestone"
)

func CreateTx(
	inputs []*Input,
	outputs []*wire.TxOut,
	changePkScript []byte,
	txFee int64,
	lockTime uint32,
) (*wire.MsgTx, error) {
	msgTx := wire.NewMsgTx(2)
	// add TxIns into raw tx
	totalInputAmount := int64(0)
	for _, input := range inputs {
		utxoHash, err := chainhash.NewHashFromStr(input.TxHash)
		if err != nil {
			return nil, err
		}
		outPoint := wire.NewOutPoint(utxoHash, input.OutputIdx)
		txIn := wire.NewTxIn(outPoint, nil, nil)
		txIn.Sequence = lockTime
		msgTx.AddTxIn(txIn)
		totalInputAmount += input.OutputAmount
	}
	// add TxOuts into raw tx
	totalOutputAmount := txFee
	for _, output := range outputs {
		msgTx.AddTxOut(output)

		totalOutputAmount += output.Value
	}
	// check amount of input coins and output coins
	if totalInputAmount < totalOutputAmount {
		return nil, fmt.Errorf("CreateMultisigTx - Total input amount %v is less than total output amount %v", totalInputAmount, totalOutputAmount)
	}
	// calculate the change output
	if totalInputAmount > totalOutputAmount {
		changeAmt := totalInputAmount - totalOutputAmount
		if changeAmt >= MIN_SAT {
			// adding the destination address and the amount to the transaction
			redeemTxOut := wire.NewTxOut(changeAmt, changePkScript)
			msgTx.AddTxOut(redeemTxOut)
		}
	}

	return msgTx, nil
}

func CreateBoundTxRedeemStableCoin(
	msg *BoundRedeemStableCoinMsg,
	inputs []*Input,
	relayerPkScript []byte,
	userPkScript []byte,
	txFee int64,
) (*wire.MsgTx, error) {
	// the inputs should be from user wallet
	for idx, input := range inputs {
		if !bytes.Equal(input.PkScript, userPkScript) {
			return nil, fmt.Errorf("the input %v should be from trading wallet", idx)
		}
	}

	boundScript, _ := CreateBoundRedeemScript(msg)

	userChangeOutput := uint32(3)
	runeOutput := &runestone.Runestone{
		Edicts: []runestone.Edict{
			{
				ID:		msg.BoundRuneId,
				Amount:	msg.Amount,
				Output: 0,
			},
		},
		Pointer: &userChangeOutput,
	}
	runeScript, _ := runeOutput.Encipher()

	outputs := []*wire.TxOut{
		// bound rune send to relayer
		{
			Value: DUST_UTXO_AMOUNT,
			PkScript: relayerPkScript,
		},
		// bound OP_RETURN
		{
			Value: 0,
			PkScript: boundScript,
		},
		// rune OP_RETURN
		{
			Value: 0,
			PkScript: runeScript,
		},
		// rune change output
		{
			Value: DUST_UTXO_AMOUNT,
			PkScript: userPkScript,
		},
	}

	return CreateTx(inputs, outputs, userPkScript, txFee, 0)
}

func CreateBoundTxMintRequest(
	msg *BoundMintRequestMsg,
	inputs []*Input,
	relayerPkScript []byte,
	txFee int64,
) (*wire.MsgTx, error) {
	// the inputs should be from relayer wallet
	for idx, input := range inputs {
		if !bytes.Equal(input.PkScript, relayerPkScript) {
			return nil, fmt.Errorf("the input %v should be from relayer wallet", idx)
		}
	}

	relayerChangeOutput := uint32(2)
	runeOutput := &runestone.Runestone{
		Edicts: []runestone.Edict{
			{
				ID:		msg.BoundRuneId,
				Amount:	msg.Amount,
				Output: 0,
			},
		},
		Pointer: &relayerChangeOutput,
	}
	runeScript, _ := runeOutput.Encipher()

	outputs := []*wire.TxOut{
		// bound rune send to user
		{
			Value: DUST_UTXO_AMOUNT,
			PkScript: msg.ToPkScript,
		},
		// rune OP_RETURN
		{
			Value: 0,
			PkScript: runeScript,
		},
		// rune change output
		{
			Value: DUST_UTXO_AMOUNT,
			PkScript: relayerPkScript,
		},
	}

	return CreateTx(inputs, outputs, relayerPkScript, txFee, 0)
}

func SignTapMultisig(
	privKey string,
	msgTx *wire.MsgTx,
	inputs []*Input,
	multisigWallet *MultisigWallet,
	indexTapLeaf int,
) ([][]byte, error) {
	if len(inputs) != len(msgTx.TxIn) {
		return nil, fmt.Errorf("len of inputs %v and TxIn %v mismatch", len(inputs), len(msgTx.TxIn))
	}
	prevOuts := txscript.NewMultiPrevOutFetcher(nil)
	for _, input := range inputs {
		utxoHash, err := chainhash.NewHashFromStr(input.TxHash)
		if err != nil {
			return nil, err
		}
		outPoint := wire.NewOutPoint(utxoHash, input.OutputIdx)

		prevOuts.AddPrevOut(*outPoint, &wire.TxOut{
			Value:    input.OutputAmount,
			PkScript: input.PkScript,
		})
	}
	txSigHashes := txscript.NewTxSigHashes(msgTx, prevOuts)

	wif, err := btcutil.DecodeWIF(privKey)
	if err != nil {
		return nil, fmt.Errorf("[PartSignOnRawExternalTx] Error when generate btc private key from seed: %v", err)
	}
	// sign on each TxIn
	tapLeaf := multisigWallet.TapLeaves[0]
	sigs := [][]byte{}
	for i, input := range inputs {
		if bytes.Equal(input.PkScript, multisigWallet.PKScript) {
			sig, err := txscript.RawTxInTapscriptSignature(
				msgTx, txSigHashes, i, int64(inputs[i].OutputAmount), multisigWallet.PKScript, tapLeaf, txscript.SigHashDefault, wif.PrivKey)
			if err != nil {
				return nil, fmt.Errorf("fail to sign tx: %v", err)
			}

			sigs = append(sigs, sig)
		} else {
			sigs = append(sigs, []byte{})
		}
	}

	return sigs, nil
}

func CombineTapMultisig(
	totalSigs [][][]byte,
	msgTx *wire.MsgTx,
	inputs []*Input,
	multisigWallet *MultisigWallet,
	indexTapLeaf int,
) (*wire.MsgTx, error) {
	tapLeafScript := multisigWallet.TapLeaves[indexTapLeaf].Script
	multisigControlBlock := multisigWallet.TapScriptTree.LeafMerkleProofs[indexTapLeaf].ToControlBlock(multisigWallet.SharedPublicKey)
	multisigControlBlockBytes, err := multisigControlBlock.ToBytes()
	if err != nil {
		return nil, err
	}

	transposedSigs := TransposeSigs(totalSigs)
	for idx, v := range transposedSigs {
		if bytes.Equal(inputs[idx].PkScript, multisigWallet.PKScript) {
			reverseV := [][]byte{}
			for i := len(v) - 1; i >= 0; i-- {
				if (len(v[i]) != 0) {
					reverseV = append(reverseV, v[i])
				}
			}

			msgTx.TxIn[idx].Witness = append(reverseV, tapLeafScript, multisigControlBlockBytes)
		}
	}

	return msgTx, nil
}
