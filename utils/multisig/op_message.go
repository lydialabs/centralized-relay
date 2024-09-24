package multisig

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/bxelab/runestone"
	"github.com/ethereum/go-ethereum/common"
	"lukechampine.com/uint128"
)

const (
	OP_RADFI_IDENT				= txscript.OP_12
	OP_RUNE_IDENT				= txscript.OP_13
	OP_BRIDGE_IDENT				= txscript.OP_14

	OP_RADFI_PROVIDE_LIQUIDITY	= txscript.OP_1
	OP_RADFI_SWAP				= txscript.OP_2
	OP_RADFI_WITHDRAW_LIQUIDITY	= txscript.OP_3
	OP_RADFI_COLLECT_FEES		= txscript.OP_4
	OP_RADFI_INCREASE_LIQUIDITY	= txscript.OP_5
)

type TokenId struct {
	BlockNumber uint64
	TxIndex     uint32
}

type XCallMessage struct {
	Action       string
	TokenAddress string
	From         string
	To           string
	Amount       []byte
	Data         []byte
}

type RadFiProvideLiquidityTicks struct {
	UpperTick	int32
	LowerTick	int32
}

type RadFiProvideLiquidityMsg struct {
	// OP_RETURN output data
	Ticks			RadFiProvideLiquidityTicks
	Fee				uint32
	Min0			uint16
	Min1			uint16
	Amount0Desired	uint128.Uint128
	Amount1Desired	uint128.Uint128
	InitPrice		uint128.Uint128
	// other outputs data
	Token0Id		string
	Token1Id		string
	// smart contract data
	Token0Addr		common.Address
	Token1Addr		common.Address
}

type RadFiSwapMsg struct {
	// OP_RETURN output data
	IsExactIn		bool
	PoolsCount		uint8
	TokenOutIndex	uint32
	AmountIn 		uint128.Uint128
	AmountOut 		uint128.Uint128
	Fees 			[]uint32
	Tokens			[]TokenId
}

type RadFiWithdrawLiquidityMsg struct {
	// OP_RETURN output data
	RecipientIndex	uint32
	LiquidityValue	uint128.Uint128
	NftId			uint128.Uint128
	Amount0			uint128.Uint128
	Amount1			uint128.Uint128
}

type RadFiCollectFeesMsg struct {
	// OP_RETURN output data
	RecipientIndex	uint32
	NftId			uint128.Uint128
}

type RadFiIncreaseLiquidityMsg struct {
	// OP_RETURN output data
	Min0	uint16
	Min1	uint16
	NftId	uint128.Uint128
	Amount0	uint128.Uint128
	Amount1	uint128.Uint128
}

type RadFiDecodedMsg struct {
	Flag					byte
	ProvideLiquidityMsg		*RadFiProvideLiquidityMsg
	SwapMsg					*RadFiSwapMsg
	WithdrawLiquidityMsg	*RadFiWithdrawLiquidityMsg
	CollectFeesMsg			*RadFiCollectFeesMsg
	IncreaseLiquidityMsg	*RadFiIncreaseLiquidityMsg
}

func integers(payload []byte) ([]uint128.Uint128, error) {
	integers := make([]uint128.Uint128, 0)
	i := 0

	for i < len(payload) {
		integer, length, err := uvarint128(payload[i:])
		if err != nil {
			return nil, err
		}
		integers = append(integers, integer)
		i += length
	}

	return integers, nil
}

func uvarint128(buf []byte) (uint128.Uint128, int, error) {
	n := big.NewInt(0)
	for i, tick := range buf {
		if i > 18 {
			return uint128.Zero, 0, fmt.Errorf("varint too long")
		}
		value := uint64(tick) & 0b0111_1111
		if i == 18 && value&0b0111_1100 != 0 {
			return uint128.Zero, 0, fmt.Errorf("varint too large")
		}
		temp := new(big.Int).SetUint64(value)
		n.Or(n, temp.Lsh(temp, uint(7*i)))
		if tick&0b1000_0000 == 0 {
			return uint128.FromBig(n), i + 1, nil
		}
	}
	return uint128.Zero, 0, fmt.Errorf("varint too short")
}

func EncodeUint16(n uint16) []byte {
	var result []byte
	for n >= 128 {
		result = append(result, byte(n&0x7F|0x80))
		n >>= 7
	}
	result = append(result, byte(n))
	return result
}

func CreateProvideLiquidityScript(msg *RadFiProvideLiquidityMsg) ([]byte, error) {
	builder := txscript.NewScriptBuilder()
	builder.AddOp(txscript.OP_RETURN)
	builder.AddOp(OP_RADFI_IDENT)
	builder.AddOp(OP_RADFI_PROVIDE_LIQUIDITY)
	// encode message content
	buf := new(bytes.Buffer)
	ticksData := msg.Ticks
	err := binary.Write(buf, binary.BigEndian, ticksData)
	if err != nil {
		return nil, fmt.Errorf("could not encode data - Error %v", err)
	}

	data := buf.Bytes()
	data = append(data, runestone.EncodeUint32(msg.Fee)...)
	data = append(data, EncodeUint16(msg.Min0)...)
	data = append(data, EncodeUint16(msg.Min1)...)
	data = append(data, runestone.EncodeUint128(msg.Amount0Desired)...)
	data = append(data, runestone.EncodeUint128(msg.Amount1Desired)...)
	data = append(data, runestone.EncodeUint128(msg.InitPrice)...)

	return builder.AddData(data).Script()
}

func CreateSwapScript(msg *RadFiSwapMsg) ([]byte, error) {
	if msg.PoolsCount > 127 {
		return nil, fmt.Errorf("max pools for swap route is 127")
	}
	if msg.PoolsCount != uint8(len(msg.Fees)) || msg.PoolsCount != uint8(len(msg.Tokens) - 1) {
		return nil, fmt.Errorf("fees and tokens array length mismatch with pools count")
	}

	builder := txscript.NewScriptBuilder()
	builder.AddOp(txscript.OP_RETURN)
	builder.AddOp(OP_RADFI_IDENT)
	builder.AddOp(OP_RADFI_SWAP)
	// encode message content
	var isExactInUint8 uint8
	if msg.IsExactIn {
		isExactInUint8 = 1
	}
	singleByte := byte((isExactInUint8 << 7) ^ msg.PoolsCount) // 1 byte contain both IsExactIn and PoolsCount

	data := []byte{singleByte}
	data = append(data, runestone.EncodeUint32(msg.TokenOutIndex)...)
	data = append(data, runestone.EncodeUint128(msg.AmountIn)...)
	data = append(data, runestone.EncodeUint128(msg.AmountOut)...)
	for _, fee := range msg.Fees {
		data = append(data, runestone.EncodeUint32(fee)...)
	}
	for _, token := range msg.Tokens {
		data = append(data, runestone.EncodeUint64(token.BlockNumber)...)
		data = append(data, runestone.EncodeUint32(token.TxIndex)...)
	}

	return builder.AddData(data).Script()
}

func CreateWithdrawLiquidityScript(msg *RadFiWithdrawLiquidityMsg) ([]byte, error) {
	builder := txscript.NewScriptBuilder()
	builder.AddOp(txscript.OP_RETURN)
	builder.AddOp(OP_RADFI_IDENT)
	builder.AddOp(OP_RADFI_WITHDRAW_LIQUIDITY)
	// encode message content
	data := runestone.EncodeUint32(msg.RecipientIndex)
	data = append(data, runestone.EncodeUint128(msg.LiquidityValue)...)
	data = append(data, runestone.EncodeUint128(msg.NftId)...)
	data = append(data, runestone.EncodeUint128(msg.Amount0)...)
	data = append(data, runestone.EncodeUint128(msg.Amount1)...)

	return builder.AddData(data).Script()
}

func CreateCollectFeesScript(msg *RadFiCollectFeesMsg) ([]byte, error) {
	builder := txscript.NewScriptBuilder()
	builder.AddOp(txscript.OP_RETURN)
	builder.AddOp(OP_RADFI_IDENT)
	builder.AddOp(OP_RADFI_COLLECT_FEES)
	// encode message content
	data := runestone.EncodeUint32(msg.RecipientIndex)
	data = append(data, runestone.EncodeUint128(msg.NftId)...)

	return builder.AddData(data).Script()
}

func CreateIncreaseLiquidityScript(msg *RadFiIncreaseLiquidityMsg) ([]byte, error) {
	builder := txscript.NewScriptBuilder()
	builder.AddOp(txscript.OP_RETURN)
	builder.AddOp(OP_RADFI_IDENT)
	builder.AddOp(OP_RADFI_INCREASE_LIQUIDITY)
	// encode message content
	data := EncodeUint16(msg.Min0)
	data = append(data, EncodeUint16(msg.Min1)...)
	data = append(data, runestone.EncodeUint128(msg.NftId)...)
	data = append(data, runestone.EncodeUint128(msg.Amount0)...)
	data = append(data, runestone.EncodeUint128(msg.Amount1)...)

	return builder.AddData(data).Script()
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

	var integersStart uint
	if flag == OP_RADFI_PROVIDE_LIQUIDITY {
		integersStart = 16
	} else if flag == OP_RADFI_SWAP {
		integersStart = 1
	} else {
		integersStart = 0
	}
	integers, err := integers(payload[integersStart:])
	if err != nil {
		return nil, fmt.Errorf("could not decode data - Error %v", err)
	}

	// Decode RadFi message
	switch flag {
		case OP_RADFI_PROVIDE_LIQUIDITY:
			// OP_RETURN output data
			r := bytes.NewReader(payload[:16])
			var ticks RadFiProvideLiquidityTicks
			if err := binary.Read(r, binary.BigEndian, &ticks); err != nil {
				return nil, fmt.Errorf("OP_RADFI_PROVIDE_LIQUIDITY could not read ticks data - Error %v", err)
			}
			// TODO: check if integers overflow
			return &RadFiDecodedMsg {
				Flag:					flag,
				ProvideLiquidityMsg:	&RadFiProvideLiquidityMsg{
					Ticks:			ticks,
					Fee:			uint32(integers[0].Lo),
					Min0:			uint16(integers[1].Lo),
					Min1:			uint16(integers[2].Lo),
					Amount0Desired:	integers[3],
					Amount1Desired:	integers[4],
					InitPrice:		integers[5],
				},
			}, nil

		case OP_RADFI_SWAP:
			singleByte := uint8(payload[0])
			isExactIn := (singleByte >> 7) != 0
			poolsCount := singleByte << 1 >> 1
			fees := []uint32{}
			for _, fee := range(integers[3:3+poolsCount]) {
				fees = append(fees, uint32(fee.Lo))
			}
			tokens := []TokenId{}
			for i := 3+int(poolsCount); i < len(integers) ; i += 2 {
				tokens = append(tokens, TokenId{
					BlockNumber: integers[i].Lo,
					TxIndex: uint32(integers[i+1].Lo),
				})
			}
			// TODO: check if integers overflow
			return &RadFiDecodedMsg {
				Flag:		flag,
				SwapMsg:	&RadFiSwapMsg{
					IsExactIn:		isExactIn,
					PoolsCount:		poolsCount,
					TokenOutIndex:	uint32(integers[0].Lo),
					AmountIn:		integers[1],
					AmountOut:		integers[2],
					Fees:			fees,
					Tokens:			tokens,
				},
			}, nil

		case OP_RADFI_WITHDRAW_LIQUIDITY:
			return &RadFiDecodedMsg {
				Flag: 					flag,
				WithdrawLiquidityMsg:	&RadFiWithdrawLiquidityMsg{
					RecipientIndex:	uint32(integers[0].Lo),
					LiquidityValue:	integers[1],
					NftId:			integers[2],
					Amount0:		integers[3],
					Amount1:		integers[4],
				},
			}, nil

		case OP_RADFI_COLLECT_FEES:
			return &RadFiDecodedMsg {
				Flag:			flag,
				CollectFeesMsg:	&RadFiCollectFeesMsg{
					RecipientIndex:	uint32(integers[0].Lo),
					NftId:			integers[1],
				},
			}, nil

		case OP_RADFI_INCREASE_LIQUIDITY:
			return &RadFiDecodedMsg {
				Flag:					flag,
				IncreaseLiquidityMsg:	&RadFiIncreaseLiquidityMsg{
					Min0:		uint16(integers[0].Lo),
					Min1:		uint16(integers[0].Lo),
					NftId:		integers[1],
					Amount0:	integers[2],
					Amount1:	integers[3],
				},
			}, nil

		default:
			return nil, fmt.Errorf("ReadRadFiMessage - invalid flag")
	}
}
