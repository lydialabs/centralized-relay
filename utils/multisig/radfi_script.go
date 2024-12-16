package multisig

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/studyzy/runestone"
	"lukechampine.com/uint128"
)

const (
	OP_RADFI_IDENT				= txscript.OP_12

	OP_RADFI_PROVIDE_LIQUIDITY	= 0x01
	OP_RADFI_SWAP				= 0x02
	OP_RADFI_WITHDRAW_LIQUIDITY	= 0x03
	OP_RADFI_COLLECT_FEES		= 0x04
	OP_RADFI_INCREASE_LIQUIDITY	= 0x05
)

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
	Min0			uint16
	Min1			uint16
	Amount0			uint128.Uint128
	Amount1			uint128.Uint128
	// for init price only
	InitPrice		uint128.Uint128
	Token0Decimal   uint8
	Token1Decimal   uint8

	SequenceNumber	uint128.Uint128
	Fee				uint32
	Token0Id		runestone.RuneId
	Token1Id		runestone.RuneId
}

type RadFiSwapMsg struct {
	// OP_RETURN output data
	IsExactIn		bool
	PoolsCount		uint8
	AmountIn 		uint128.Uint128
	AmountOut 		uint128.Uint128
	SequenceNumber	uint128.Uint128
	Fees 			[]uint32
	Tokens			[]*runestone.RuneId
}

type RadFiWithdrawLiquidityMsg struct {
	// OP_RETURN output data
	LiquidityValue	uint128.Uint128
	NftId			uint128.Uint128
	Amount0			uint128.Uint128
	Amount1			uint128.Uint128
	SequenceNumber	uint128.Uint128
	Fee				uint32
	Token0Id		runestone.RuneId
	Token1Id		runestone.RuneId
}

type RadFiCollectFeesMsg struct {
	// OP_RETURN output data
	NftId			uint128.Uint128
	Amount0			uint128.Uint128
	Amount1			uint128.Uint128
	SequenceNumber	uint128.Uint128
	Fee				uint32
	Token0Id		runestone.RuneId
	Token1Id		runestone.RuneId
}

type RadFiIncreaseLiquidityMsg struct {
	// OP_RETURN output data
	Min0	uint16
	Min1	uint16
	NftId	uint128.Uint128
	Amount0	uint128.Uint128
	Amount1	uint128.Uint128
	SequenceNumber	uint128.Uint128
	Fee				uint32
	Token0Id		runestone.RuneId
	Token1Id		runestone.RuneId
}

type RadFiDecodedMsg struct {
	Flag					byte
	ProvideLiquidityMsg		*RadFiProvideLiquidityMsg
	SwapMsg					*RadFiSwapMsg
	WithdrawLiquidityMsg	*RadFiWithdrawLiquidityMsg
	CollectFeesMsg			*RadFiCollectFeesMsg
	IncreaseLiquidityMsg	*RadFiIncreaseLiquidityMsg
}

func BitcoinRuneId() (runestone.RuneId) {
	return runestone.RuneId{
		Block:	0,
		Tx:		0,
	}
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
	builder.AddOp(OP_RADFI_IDENT)
	// encode message content
	buf := new(bytes.Buffer)
	ticksData := msg.Ticks
	err := binary.Write(buf, binary.BigEndian, ticksData)
	if err != nil {
		return nil, fmt.Errorf("could not encode data - Error %v", err)
	}

	data := append([]byte{OP_RADFI_PROVIDE_LIQUIDITY}, buf.Bytes()...)

	data = append(data, EncodeUint16(msg.Min0)...)
	data = append(data, EncodeUint16(msg.Min1)...)
	data = append(data, runestone.EncodeUint128(msg.Amount0)...)
	data = append(data, runestone.EncodeUint128(msg.Amount1)...)
	data = append(data, runestone.EncodeUint128(msg.InitPrice)...)
	data = append(data, runestone.EncodeUint8(msg.Token0Decimal)...)
	data = append(data, runestone.EncodeUint8(msg.Token1Decimal)...)
	data = append(data, runestone.EncodeUint128(msg.SequenceNumber)...)
	data = append(data, runestone.EncodeUint32(msg.Fee)...)
	data = append(data, runestone.EncodeUint64(msg.Token0Id.Block)...)
	data = append(data, runestone.EncodeUint32(msg.Token0Id.Tx)...)
	data = append(data, runestone.EncodeUint64(msg.Token1Id.Block)...)
	data = append(data, runestone.EncodeUint32(msg.Token1Id.Tx)...)

	return builder.AddData(data).Script()
}

func CreateSwapScript(msg *RadFiSwapMsg) ([]byte, error) {
	builder := txscript.NewScriptBuilder()
	builder.AddOp(OP_RADFI_IDENT)
	// encode message content
	var isExactInUint8 uint8
	if msg.IsExactIn {
		isExactInUint8 = 1
	}
	singleByte := byte((isExactInUint8 << 7) ^ msg.PoolsCount) // 1 byte contain both IsExactIn and PoolsCount

	data := []byte{OP_RADFI_SWAP, singleByte}
	data = append(data, runestone.EncodeUint128(msg.AmountIn)...)
	data = append(data, runestone.EncodeUint128(msg.AmountOut)...)
	data = append(data, runestone.EncodeUint128(msg.SequenceNumber)...)
	for _, fee := range msg.Fees {
		data = append(data, runestone.EncodeUint32(fee)...)
	}
	for _, token := range msg.Tokens {
		data = append(data, runestone.EncodeUint64(token.Block)...)
		data = append(data, runestone.EncodeUint32(token.Tx)...)
	}

	return builder.AddData(data).Script()
}

func CreateWithdrawLiquidityScript(msg *RadFiWithdrawLiquidityMsg) ([]byte, error) {
	builder := txscript.NewScriptBuilder()
	builder.AddOp(OP_RADFI_IDENT)
	// encode message content
	data := append([]byte{OP_RADFI_WITHDRAW_LIQUIDITY}, runestone.EncodeUint128(msg.LiquidityValue)...)
	data = append(data, runestone.EncodeUint128(msg.NftId)...)
	data = append(data, runestone.EncodeUint128(msg.Amount0)...)
	data = append(data, runestone.EncodeUint128(msg.Amount1)...)
	data = append(data, runestone.EncodeUint128(msg.SequenceNumber)...)
	data = append(data, runestone.EncodeUint32(msg.Fee)...)
	data = append(data, runestone.EncodeUint64(msg.Token0Id.Block)...)
	data = append(data, runestone.EncodeUint32(msg.Token0Id.Tx)...)
	data = append(data, runestone.EncodeUint64(msg.Token1Id.Block)...)
	data = append(data, runestone.EncodeUint32(msg.Token1Id.Tx)...)

	return builder.AddData(data).Script()
}

func CreateCollectFeesScript(msg *RadFiCollectFeesMsg) ([]byte, error) {
	builder := txscript.NewScriptBuilder()
	builder.AddOp(OP_RADFI_IDENT)
	// encode message content
	data := append([]byte{OP_RADFI_COLLECT_FEES}, runestone.EncodeUint128(msg.NftId)...)
	data = append(data, runestone.EncodeUint128(msg.Amount0)...)
	data = append(data, runestone.EncodeUint128(msg.Amount1)...)
	data = append(data, runestone.EncodeUint128(msg.SequenceNumber)...)
	data = append(data, runestone.EncodeUint32(msg.Fee)...)
	data = append(data, runestone.EncodeUint64(msg.Token0Id.Block)...)
	data = append(data, runestone.EncodeUint32(msg.Token0Id.Tx)...)
	data = append(data, runestone.EncodeUint64(msg.Token1Id.Block)...)
	data = append(data, runestone.EncodeUint32(msg.Token1Id.Tx)...)

	return builder.AddData(data).Script()
}

func CreateIncreaseLiquidityScript(msg *RadFiIncreaseLiquidityMsg) ([]byte, error) {
	builder := txscript.NewScriptBuilder()
	builder.AddOp(OP_RADFI_IDENT)
	// encode message content
	data := append([]byte{OP_RADFI_INCREASE_LIQUIDITY}, EncodeUint16(msg.Min0)...)
	data = append(data, EncodeUint16(msg.Min1)...)
	data = append(data, runestone.EncodeUint128(msg.NftId)...)
	data = append(data, runestone.EncodeUint128(msg.Amount0)...)
	data = append(data, runestone.EncodeUint128(msg.Amount1)...)
	data = append(data, runestone.EncodeUint128(msg.SequenceNumber)...)
	data = append(data, runestone.EncodeUint32(msg.Fee)...)
	data = append(data, runestone.EncodeUint64(msg.Token0Id.Block)...)
	data = append(data, runestone.EncodeUint32(msg.Token0Id.Tx)...)
	data = append(data, runestone.EncodeUint64(msg.Token1Id.Block)...)
	data = append(data, runestone.EncodeUint32(msg.Token1Id.Tx)...)

	return builder.AddData(data).Script()
}

func ReadRadFiMessage(transaction *wire.MsgTx) (*RadFiDecodedMsg, error) {
	var flag byte
	var payload []byte
	for _, output := range transaction.TxOut {
		tokenizer := txscript.MakeScriptTokenizer(0, output.PkScript)
		if !tokenizer.Next() || tokenizer.Err() != nil || tokenizer.Opcode() != OP_RADFI_IDENT {
			// Check to ignore non RadFi protocol identifier (Rune or Bridge)
			continue
		}

		// Construct the payload by concatenating remaining data pushes
		for tokenizer.Next() {
			if tokenizer.Err() != nil {
				return nil, tokenizer.Err()
			}
			payload = append(payload, tokenizer.Data()...)
		}

		// only read 1 OP_RADFI_IDENT output for RadFi protocol
		break
	}

	if len(payload) == 0 {
		return nil, fmt.Errorf("could not find radfi data")
	}

	flag = payload[0]
	var integersStart uint
	if flag == OP_RADFI_PROVIDE_LIQUIDITY {
		integersStart = 9
	} else if flag == OP_RADFI_SWAP {
		integersStart = 2
	} else {
		integersStart = 1
	}
	integers, err := integers(payload[integersStart:])
	if err != nil {
		return nil, fmt.Errorf("could not decode data - Error %v", err)
	}

	// Decode RadFi message
	switch flag {
		case OP_RADFI_PROVIDE_LIQUIDITY:
			// OP_RETURN output data
			r := bytes.NewReader(payload[1:9])
			var ticks RadFiProvideLiquidityTicks
			if err := binary.Read(r, binary.BigEndian, &ticks); err != nil {
				return nil, fmt.Errorf("OP_RADFI_PROVIDE_LIQUIDITY could not read ticks data - Error %v", err)
			}
			// TODO: check if integers overflow
			return &RadFiDecodedMsg {
				Flag:					flag,
				ProvideLiquidityMsg:	&RadFiProvideLiquidityMsg{
					Ticks:			ticks,
					Min0:			uint16(integers[0].Lo),
					Min1:			uint16(integers[1].Lo),
					Amount0:		integers[2],
					Amount1:		integers[3],
					InitPrice:		integers[4],
					Token0Decimal:  uint8(integers[5].Lo),
					Token1Decimal:  uint8(integers[6].Lo),
					SequenceNumber: integers[7],
					Fee:			uint32(integers[8].Lo),
					Token0Id:		runestone.RuneId{
										Block: integers[9].Lo,
										Tx: uint32(integers[10].Lo),
									},
					Token1Id:		runestone.RuneId{
										Block: integers[11].Lo,
										Tx: uint32(integers[12].Lo),
									},
				},
			}, nil

		case OP_RADFI_SWAP:
			singleByte := uint8(payload[1])
			isExactIn := (singleByte >> 7) != 0
			poolsCount := singleByte << 1 >> 1
			fees := []uint32{}
			for _, fee := range(integers[3:3+poolsCount]) {
				fees = append(fees, uint32(fee.Lo))
			}
			tokens := []*runestone.RuneId{}
			for i := 3+int(poolsCount); i < len(integers)-1 ; i += 2 {
				tokens = append(tokens, &runestone.RuneId{
					Block: integers[i].Lo,
					Tx: uint32(integers[i+1].Lo),
				})
			}
			// TODO: check if integers overflow
			return &RadFiDecodedMsg {
				Flag:		flag,
				SwapMsg:	&RadFiSwapMsg{
					IsExactIn:		isExactIn,
					PoolsCount:		poolsCount,
					AmountIn:		integers[0],
					AmountOut:		integers[1],
					SequenceNumber: integers[2],
					Fees:			fees,
					Tokens:			tokens,
				},
			}, nil

		case OP_RADFI_WITHDRAW_LIQUIDITY:
			return &RadFiDecodedMsg {
				Flag: 					flag,
				WithdrawLiquidityMsg:	&RadFiWithdrawLiquidityMsg{
					LiquidityValue:	integers[0],
					NftId:			integers[1],
					Amount0:		integers[2],
					Amount1:		integers[3],
					SequenceNumber: integers[4],
					Fee:			uint32(integers[5].Lo),
					Token0Id:		runestone.RuneId{
										Block: integers[6].Lo,
										Tx: uint32(integers[7].Lo),
									},
					Token1Id:		runestone.RuneId{
										Block: integers[8].Lo,
										Tx: uint32(integers[9].Lo),
									},
				},
			}, nil

		case OP_RADFI_COLLECT_FEES:
			return &RadFiDecodedMsg {
				Flag:			flag,
				CollectFeesMsg:	&RadFiCollectFeesMsg{
					NftId:			integers[0],
					Amount0:		integers[1],
					Amount1:		integers[2],
					SequenceNumber: integers[3],
					Fee:			uint32(integers[4].Lo),
					Token0Id:		runestone.RuneId{
										Block: integers[5].Lo,
										Tx: uint32(integers[6].Lo),
									},
					Token1Id:		runestone.RuneId{
										Block: integers[7].Lo,
										Tx: uint32(integers[8].Lo),
									},
				},
			}, nil

		case OP_RADFI_INCREASE_LIQUIDITY:
			return &RadFiDecodedMsg {
				Flag:					flag,
				IncreaseLiquidityMsg:	&RadFiIncreaseLiquidityMsg{
					Min0:		    uint16(integers[0].Lo),
					Min1:		    uint16(integers[1].Lo),
					NftId:		    integers[2],
					Amount0:	    integers[3],
					Amount1:	    integers[4],
					SequenceNumber: integers[5],
					Fee:			uint32(integers[6].Lo),
					Token0Id:		runestone.RuneId{
										Block: integers[7].Lo,
										Tx: uint32(integers[8].Lo),
									},
					Token1Id:		runestone.RuneId{
										Block: integers[9].Lo,
										Tx: uint32(integers[10].Lo),
									},
				},
			}, nil

		default:
			return nil, fmt.Errorf("ReadRadFiMessage - invalid flag")
	}
}
