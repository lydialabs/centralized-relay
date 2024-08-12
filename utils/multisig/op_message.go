package multisig

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/holiman/uint256"
)

const (
	OP_RADFI_IDENT				= txscript.OP_12
	OP_RUNE_IDENT				= txscript.OP_13
	OP_BRIDGE_IDENT				= txscript.OP_14

	OP_RADFI_PROVIDE_LIQUIDITY	= txscript.OP_1
	OP_RADFI_SWAP				= txscript.OP_2
	OP_RADFI_WITHDRAW_LIQUIDITY	= txscript.OP_3
	OP_RADFI_COLLECT_FEES		= txscript.OP_4
)

type RadFiProvideLiquidityMsg struct {
	Fee			uint8
	UpperTick	int32
	LowerTick	int32
	Min0		uint16
	Min1		uint16
}

type RadFiSwapMsg struct {
	IsExactInOut	bool
	TokenOutIndex	uint32
	// TokenOutId		*Rune
	// TokenOutAmount	*uint256.Int
}

type RadFiWithdrawLiquidityMsg struct {
	RecipientIndex	uint32
	LiquidityValue	*uint256.Int
	NftId			*uint256.Int
}

type RadFiCollectFeesMsg struct {
	RecipientIndex	uint32
	NftId			*uint256.Int
}
type RadFiDecodedMsg struct {
	Flag					byte
	ProvideLiquidityMsg		*RadFiProvideLiquidityMsg
	SwapMsg					*RadFiSwapMsg
	WithdrawLiquidityMsg	*RadFiWithdrawLiquidityMsg
	CollectFeesMsg			*RadFiCollectFeesMsg
}

func CreateBridgeMessageScripts(payload []byte, partLimit int) ([][]byte, error) {
	var chunk []byte
	chunks := make([][]byte, 0, len(payload)/partLimit+1)
	for len(payload) >= partLimit {
		chunk, payload = payload[:partLimit], payload[partLimit:]
		chunks = append(chunks, chunk)
	}
	if len(payload) > 0 {
		chunks = append(chunks, payload)
	}

	scripts := [][]byte{}
	for _, part := range chunks {
		builder := txscript.NewScriptBuilder()

		builder.AddOp(txscript.OP_RETURN)
		builder.AddOp(OP_BRIDGE_IDENT)
		builder.AddData(part)

		script, err := builder.Script()
		if err != nil {
			return nil, fmt.Errorf("could not build script - Error %v", err)
		}
		scripts = append(scripts, script)
	}

	return scripts, nil
}

func CreateProvideLiquidityScript(msg *RadFiProvideLiquidityMsg) ([]byte, error) {
	builder := txscript.NewScriptBuilder()

	builder.AddOp(txscript.OP_RETURN)
	builder.AddOp(OP_RADFI_IDENT)
	builder.AddOp(OP_RADFI_PROVIDE_LIQUIDITY)
	// encode message content
	buf := new(bytes.Buffer)
	var data = []any{ msg.Fee, msg.UpperTick, msg.LowerTick, msg.Min0, msg.Min1 }
	for _, v := range data {
		err := binary.Write(buf, binary.BigEndian, v)
		if err != nil {
			fmt.Println("CreateProvideLiquidityScript encode data failed:", err)
		}
	}

	return builder.AddData(buf.Bytes()).Script()
}

func CreateSwapScript(msg *RadFiSwapMsg) ([]byte, error) {
	builder := txscript.NewScriptBuilder()

	builder.AddOp(txscript.OP_RETURN)
	builder.AddOp(OP_RADFI_IDENT)
	builder.AddOp(OP_RADFI_SWAP)
	// encode message content
	var isExactInOutUint8 uint8
	if msg.IsExactInOut {
		isExactInOutUint8 = 1
	}

	// tokenOutIdBlockNumberByte := make([]byte, 8)
	// binary.BigEndian.PutUint64(tokenOutIdBlockNumberByte, msg.TokenOutId.BlockNumber)
	// tokenOutIdBlockNumberLen := uint8(bits.Len64(msg.TokenOutId.BlockNumber))

	// tokenOutIdTxIndexByte := make([]byte, 4)
	// binary.BigEndian.PutUint32(tokenOutIdTxIndexByte, msg.TokenOutId.TxIndex)
	// tokenOutIdTxIndexLen := uint8(bits.Len32(msg.TokenOutId.TxIndex))

	// fmt.Println("tokenOutIdBlockNumberLen: ", tokenOutIdBlockNumberLen)
	// singleByte := byte((isExactInOutUint8 << 7) ^ (tokenOutIdBlockNumberLen << 3) ^ tokenOutIdTxIndexLen)
	// fmt.Println("singleByte: ", singleByte)
	// data := append([]byte{singleByte}, tokenOutIdBlockNumberByte[8-tokenOutIdBlockNumberLen:]...)
	// data = append(data, tokenOutIdTxIndexByte[4-tokenOutIdTxIndexLen:]...)
	// data = append(data, msg.TokenOutAmount.Bytes()...)

	TokenOutIndexByte := make([]byte, 4)
	binary.BigEndian.PutUint32(TokenOutIndexByte, msg.TokenOutIndex)
	data := append([]byte{isExactInOutUint8}, TokenOutIndexByte...)

	return builder.AddData(data).Script()
}

func CreateWithdrawLiquidityScript(msg *RadFiWithdrawLiquidityMsg) ([]byte, error) {
	builder := txscript.NewScriptBuilder()

	builder.AddOp(txscript.OP_RETURN)
	builder.AddOp(OP_RADFI_IDENT)
	builder.AddOp(OP_RADFI_WITHDRAW_LIQUIDITY)
	// encode message content
	recipientIndexByte := make([]byte, 4)
	binary.BigEndian.PutUint32(recipientIndexByte, msg.RecipientIndex)
	// recipientIndexLen := uint8(bits.Len32(msg.RecipientIndex))
	liquidityValueBytes := msg.LiquidityValue.Bytes()
	liquidityValueBytesLen := uint8(len(liquidityValueBytes))
	// singleByte := byte((recipientIndexLen << 5) ^ liquidityValueBytesLen)
	singleByte := byte(liquidityValueBytesLen)
	// data := append([]byte{singleByte}, recipientIndexByte[4-recipientIndexLen:]...)
	data := append([]byte{singleByte}, recipientIndexByte...)
	data = append(data, liquidityValueBytes...)
	data = append(data, msg.NftId.Bytes()...)
	// fmt.Println("data: ", data)
	// fmt.Println("singleByte: ", singleByte)
	// fmt.Println("recipientIndexByte: ", recipientIndexByte)
	// fmt.Println("liquidityValueBytes: ", liquidityValueBytes)
	return builder.AddData(data).Script()
}

func CreateCollectFeesScript(msg *RadFiCollectFeesMsg) ([]byte, error) {
	builder := txscript.NewScriptBuilder()

	builder.AddOp(txscript.OP_RETURN)
	builder.AddOp(OP_RADFI_IDENT)
	builder.AddOp(OP_RADFI_COLLECT_FEES)
	// encode message content
	recipientIndexByte := make([]byte, 4)
	binary.BigEndian.PutUint32(recipientIndexByte, msg.RecipientIndex)
	data := append(recipientIndexByte, msg.NftId.Bytes()...)

	return builder.AddData(data).Script()
}

func ReadBridgeMessage(transaction *wire.MsgTx) ([]byte, error) {
	payload := []byte{}
	for _, output := range transaction.TxOut {
		tokenizer := txscript.MakeScriptTokenizer(0, output.PkScript)
		if !tokenizer.Next() || tokenizer.Err() != nil || tokenizer.Opcode() != txscript.OP_RETURN {
			// Check for OP_RETURN
			continue
		}
		if !tokenizer.Next() || tokenizer.Err() != nil || tokenizer.Opcode() != OP_BRIDGE_IDENT {
			// Check to ignore non Bridge protocol identifier (Rune or RadFi)
			continue
		}

		// Construct the payload by concatenating remaining data pushes
		for tokenizer.Next() {
			if tokenizer.Err() != nil {
				return nil, tokenizer.Err()
			}
			payload = append(payload, tokenizer.Data()...)
		}
	}

	if len(payload) == 0 {
		return nil, fmt.Errorf("ReadBridgeMessage - no Bridge message found")
	}

	return payload, nil
}

func ReadRadFiMessage(transaction *wire.MsgTx) (*RadFiDecodedMsg, error) {
	var flag byte
	var payload []byte
	for _, output := range transaction.TxOut {
		tokenizer := txscript.MakeScriptTokenizer(0, output.PkScript)
		if !tokenizer.Next() || tokenizer.Err() != nil || tokenizer.Opcode() != txscript.OP_RETURN {
			// Check for OP_RETURN
			continue
		}
		if !tokenizer.Next() || tokenizer.Err() != nil || tokenizer.Opcode() != OP_RADFI_IDENT {
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
		case OP_RADFI_PROVIDE_LIQUIDITY:
			r := bytes.NewReader(payload)
			var provideLiquidityMsg RadFiProvideLiquidityMsg
			if err := binary.Read(r, binary.BigEndian, &provideLiquidityMsg); err != nil {
				fmt.Println("OP_RADFI_PROVIDE_LIQUIDITY Read failed:", err)
			}

			return &RadFiDecodedMsg {
				Flag				: flag,
				ProvideLiquidityMsg	: &provideLiquidityMsg,
			}, nil

		case OP_RADFI_SWAP:
			// singleByte := uint8(payload[0])
			// isExactInOut := (singleByte >> 7) != 0
			// tokenOutIdBlockNumberLen := singleByte << 1 >> 4
			// tokenOutIdTxIndexLen := singleByte << 5 >> 5

			// payload = payload[1:]
			// tokenOutIdBlockNumber := binary.BigEndian.Uint64(payload[:tokenOutIdBlockNumberLen])

			// payload = payload[tokenOutIdBlockNumberLen:]
			// tokenOutIdTxIndex := binary.BigEndian.Uint32(payload[:tokenOutIdTxIndexLen])

			// TokenOutAmount := new(uint256.Int).SetBytes(payload[tokenOutIdTxIndexLen:])

			return &RadFiDecodedMsg {
				Flag				: flag,
				SwapMsg: &RadFiSwapMsg{
					IsExactInOut	: payload[0] != 0,
					TokenOutIndex	: binary.BigEndian.Uint32(payload[1:]),
					// TokenOutId		: &Rune{
					// 	BlockNumber	: tokenOutIdBlockNumber,
					// 	TxIndex		: tokenOutIdTxIndex,
					// },
					// TokenOutAmount	: TokenOutAmount,
				},
			}, nil

		case OP_RADFI_WITHDRAW_LIQUIDITY:
			// singleByte := uint8(payload[0])
			// recipientIndexLen := singleByte >> 5
			// fmt.Println("recipientIndexLen", recipientIndexLen)
			// liquidityValueBytesLen := singleByte << 3 >> 3
			liquidityValueBytesLen := uint8(payload[0])
			// fmt.Println("liquidityValueBytesLen", liquidityValueBytesLen)
			payload = payload[1:]
			recipientIndex := binary.BigEndian.Uint32(payload[:4])
			payload = payload[4:]
			liquidityValue := new(uint256.Int).SetBytes(payload[:liquidityValueBytesLen])
			nftId := new(uint256.Int).SetBytes(payload[liquidityValueBytesLen:])
			// fmt.Println("OP_RADFI_WITHDRAW_LIQUIDITY", recipientIndex, liquidityValue, NftId)

			return &RadFiDecodedMsg {
				Flag				: flag,
				WithdrawLiquidityMsg: &RadFiWithdrawLiquidityMsg{
					RecipientIndex	: recipientIndex,
					LiquidityValue	: liquidityValue,
					NftId			: nftId,
				},
			}, nil

		case OP_RADFI_COLLECT_FEES:
			recipientIndex := binary.BigEndian.Uint32(payload[:4])
			nftId := new(uint256.Int).SetBytes(payload[4:])

			return &RadFiDecodedMsg {
				Flag				: flag,
				CollectFeesMsg: &RadFiCollectFeesMsg{
					RecipientIndex	: recipientIndex,
					NftId			: nftId,
				},
			}, nil

		default:
			return nil, fmt.Errorf("ReadRadFiMessage - invalid flag")
	}
}
