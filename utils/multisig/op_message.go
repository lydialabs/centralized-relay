package multisig

import (
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/bxelab/runestone"
	"lukechampine.com/uint128"
)

const (
	OP_BOUND_IDENT = txscript.OP_15

	OP_BOUND_REDEEM_STABLE_COIN  = txscript.OP_1
)

type BoundRedeemStableCoinMsg struct {
	// OP_RETURN output data
	StableCoinType  uint8
	To           	string
	// other outputs data
	BoundRuneId		runestone.RuneId
	Amount			uint128.Uint128
}

type BoundDecodedMsg struct {
	Flag					byte
	RedeemStableCoinMsg		*BoundRedeemStableCoinMsg
}

func AddressToPayload(address string) ([]byte, error) {
	prefix := address[0:2]
	if prefix != "0x" {
		return nil, fmt.Errorf("address type not supported")
	}
	addressBytes, err := hex.DecodeString(address[2:])
	if err != nil {
		return nil, fmt.Errorf("could decode string address - Error %v", err)
	}

	return addressBytes, nil
}

func PayloadToEvmAddress(payload []byte) string {
	return "0x" + hex.EncodeToString(payload)
}

func CreateBoundRedeemScript(msg *BoundRedeemStableCoinMsg) ([]byte, error) {
	builder := txscript.NewScriptBuilder()

	builder.AddOp(txscript.OP_RETURN)
	builder.AddOp(OP_BOUND_IDENT)
	builder.AddOp(OP_BOUND_REDEEM_STABLE_COIN)

	receiverBytes, err := AddressToPayload(msg.To)
	if err != nil {
		return nil, err
	}

	return builder.AddData(append([]byte{ byte(msg.StableCoinType) }, receiverBytes...)).Script()
}

func ReadRadFiMessage(transaction *wire.MsgTx) (*BoundDecodedMsg, error) {
	var flag byte
	var payload []byte
	for _, output := range transaction.TxOut {
		tokenizer := txscript.MakeScriptTokenizer(0, output.PkScript)
		if !tokenizer.Next() || tokenizer.Err() != nil || tokenizer.Opcode() != txscript.OP_RETURN {
			// Check for OP_RETURN
			continue
		}
		if !tokenizer.Next() || tokenizer.Err() != nil || tokenizer.Opcode() != OP_BOUND_IDENT {
			// Check to ignore non RadFi protocol identifier (Rune or Bridge)
			continue
		}

		if tokenizer.Next() && tokenizer.Err() == nil {
			flag = tokenizer.Opcode()
		}

		// Construct the payload by concatenating remaining data pushes
		for tokenizer.Next() {
			if tokenizer.Err() != nil {
				return nil, tokenizer.Err()
			}
			payload = append(payload, tokenizer.Data()...)
		}

		// only read 1 OP_RETURN output for RadFi protocol
		break
	}

	// Decode RadFi message
	switch flag {
	case OP_BOUND_REDEEM_STABLE_COIN:
		return &BoundDecodedMsg {
			Flag:					flag,
			RedeemStableCoinMsg:	&BoundRedeemStableCoinMsg{
				StableCoinType:	payload[0],
				To:				PayloadToEvmAddress(payload[1:]),
			},
		}, nil

	default:
		return nil, fmt.Errorf("ReadBoundMessage - invalid flag")
	}
}
