package multisig

import (
	"bytes"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/studyzy/runestone"
	"lukechampine.com/uint128"
)

func CreateRadFiTx(
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

func CreateRadFiTxInitPool(
	msg *RadFiProvideLiquidityMsg,
	inputs []*Input,
	relayerPkScript []byte,
	userPkScript []byte,
	poolUTXOsCount uint64,
	poolUTXOsContainLiquidityCount uint64,
	txFee int64,
) (*wire.MsgTx, error) {
	if msg.InitPrice.Lo == 0 {
		return nil, fmt.Errorf("init price should be > 0")
	}
	if msg.Token1Id == BitcoinRuneId() {
		return nil, fmt.Errorf("the second token in the pair can only be rune")
	}
	// the first input should be the pool init sequence UTXO
	if !bytes.Equal(inputs[0].PkScript, relayerPkScript) {
		return nil, fmt.Errorf("the first input should be the pool init sequence UTXO")
	}
	// the remain inputs should be from trading wallet
	for idx, input := range inputs[1:] {
		if !bytes.Equal(input.PkScript, userPkScript) {
			return nil, fmt.Errorf("the input %v should be from trading wallet", idx)
		}
	}
	// must have at least 1 pool UTXOs
	if poolUTXOsCount <= 0 || poolUTXOsContainLiquidityCount <= 0 {
		return nil, fmt.Errorf("invalid poolUTXOsCount or poolUTXOsContainLiquidityCount")
	}

	radfiScript, _ := CreateProvideLiquidityScript(msg)

	userChangeOutput := uint32(poolUTXOsCount+2)
	runeOutput := &runestone.Runestone{
		Edicts: []runestone.Edict{},
		Pointer: &userChangeOutput,
	}

	outputs := []*wire.TxOut{
		// pool init sequence UTXO
		{
			Value: DUST_UTXO_AMOUNT,
			PkScript: relayerPkScript,
		},
	}
	// add pool UTXOs to output
	for poolUTXOIdx := range(poolUTXOsCount) {
		bitcoinAmount := DUST_UTXO_AMOUNT
		// only some pool UTXOs have liquidity due to runestone limit
		if poolUTXOIdx < poolUTXOsContainLiquidityCount {
			if msg.Token0Id == BitcoinRuneId() {
					bitcoinAmount += int64(msg.Amount0Desired.Lo / poolUTXOsContainLiquidityCount)
			} else {
				runeOutput.Edicts = append(runeOutput.Edicts, runestone.Edict{
					ID:	msg.Token0Id,
					Amount: msg.Amount0Desired.Div64(poolUTXOsContainLiquidityCount),
					Output: uint32(poolUTXOIdx),
				})
			}
			runeOutput.Edicts = append(runeOutput.Edicts, runestone.Edict{
				ID:	msg.Token1Id,
				Amount: msg.Amount1Desired.Div64(poolUTXOsContainLiquidityCount),
				Output: uint32(poolUTXOIdx),
			})
		}
		// pool UTXO
		outputs =  append(outputs, &wire.TxOut{
			Value: bitcoinAmount,
			PkScript: relayerPkScript,
		})
		fmt.Println("Token1Id: ", msg.Token1Id)
	}
	if msg.Token0Id == BitcoinRuneId() {
		outputs[0].Value += int64(msg.Amount0Desired.Lo % poolUTXOsContainLiquidityCount)
		runeOutput.Edicts[0].Amount = runeOutput.Edicts[0].Amount.Add64(msg.Amount1Desired.Mod64(poolUTXOsContainLiquidityCount))
	} else {
		runeOutput.Edicts[0].Amount = runeOutput.Edicts[0].Amount.Add64(msg.Amount0Desired.Mod64(poolUTXOsContainLiquidityCount))
		runeOutput.Edicts[1].Amount = runeOutput.Edicts[1].Amount.Add64(msg.Amount1Desired.Mod64(poolUTXOsContainLiquidityCount))
	}
	runeScript, _ := runeOutput.Encipher()

	outputs = append(outputs,[]*wire.TxOut{
		// radfi OP
		{
			Value: DUST_UTXO_AMOUNT,
			PkScript: radfiScript,
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
	}...)

	return CreateRadFiTx(inputs, outputs, userPkScript, txFee, 0)
}

func CreateRadFiTxProvideLiquidity(
	msg *RadFiProvideLiquidityMsg,
	inputs []*Input,
	relayerPkScript []byte,
	userPkScript []byte,
	poolUTXOsCount uint64,
	txFee int64,
) (*wire.MsgTx, error) {
	if msg.InitPrice.Lo != 0 {
		return nil, fmt.Errorf("init price should be 0")
	}
	if msg.Token1Id == BitcoinRuneId() {
		return nil, fmt.Errorf("the second token in the pair can only be rune")
	}
	// the first inputs should be the pool's current sequence number that contain pool's liquidity
	for idx, input := range inputs[:poolUTXOsCount] {
		if !bytes.Equal(input.PkScript, relayerPkScript) {
			return nil, fmt.Errorf("the input %v should be from relayer wallet", idx)
		}
	}
	// the remain inputs should be from trading wallet
	for idx, input := range inputs[poolUTXOsCount:] {
		if !bytes.Equal(input.PkScript, userPkScript) {
			return nil, fmt.Errorf("the input %v should be from trading wallet", idx)
		}
	}
	// must have at least 1 pool UTXOs
	if poolUTXOsCount <= 0 {
		return nil, fmt.Errorf("invalid poolUTXOsCount")
	}

	radfiScript, _ := CreateProvideLiquidityScript(msg)

	userChangeOutput := uint32(poolUTXOsCount+2)
	runeOutput := &runestone.Runestone{
		Edicts: []runestone.Edict{},
		Pointer: &userChangeOutput,
	}

	outputs := []*wire.TxOut{}
	// add pool UTXOs to output
	for poolUTXOIdx := range(poolUTXOsCount) {
		bitcoinAmount := inputs[poolUTXOIdx].OutputAmount
		var inputToken1Amount uint128.Uint128
		if msg.Token0Id == BitcoinRuneId() {
				bitcoinAmount += int64(msg.Amount0Desired.Lo / poolUTXOsCount)
				inputToken1Amount = inputs[poolUTXOIdx].Runes[0].Amount
		} else {
			runeOutput.Edicts = append(runeOutput.Edicts, runestone.Edict{
				ID:	msg.Token0Id,
				Amount: inputs[poolUTXOIdx].Runes[0].Amount.Add(msg.Amount0Desired.Div64(poolUTXOsCount)),
				Output: uint32(poolUTXOIdx),
			})
			inputToken1Amount = inputs[poolUTXOIdx].Runes[1].Amount
		}
		runeOutput.Edicts = append(runeOutput.Edicts, runestone.Edict{
			ID:	msg.Token1Id,
			Amount: inputToken1Amount.Add(msg.Amount1Desired.Div64(poolUTXOsCount)),
			Output: uint32(poolUTXOIdx),
		})

		// pool UTXO
		outputs =  append(outputs, &wire.TxOut{
			Value: bitcoinAmount,
			PkScript: relayerPkScript,
		})
	}
	if msg.Token0Id == BitcoinRuneId() {
		outputs[0].Value += int64(msg.Amount0Desired.Lo % poolUTXOsCount)
		runeOutput.Edicts[0].Amount = runeOutput.Edicts[0].Amount.Add64(msg.Amount1Desired.Mod64(poolUTXOsCount))
	} else {
		runeOutput.Edicts[0].Amount = runeOutput.Edicts[0].Amount.Add64(msg.Amount0Desired.Mod64(poolUTXOsCount))
		runeOutput.Edicts[1].Amount = runeOutput.Edicts[1].Amount.Add64(msg.Amount1Desired.Mod64(poolUTXOsCount))
	}
	runeScript, _ := runeOutput.Encipher()

	outputs = append(outputs,[]*wire.TxOut{
		// radfi OP
		{
			Value: DUST_UTXO_AMOUNT,
			PkScript: radfiScript,
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
	}...)

	return CreateRadFiTx(inputs, outputs, userPkScript, txFee, 0)
}

func CreateRadFiTxSwap(
	msg *RadFiSwapMsg,
	inputs []*Input,
	newPoolBalances []*PoolBalance,
	relayerPkScript []byte,
	userPkScript []byte,
	txFee int64,
) (*wire.MsgTx, error) {
	if msg.PoolsCount < 1 || msg.PoolsCount > 127 {
		return nil, fmt.Errorf("pools count should be in the range of [1; 127]")
	}
	if msg.PoolsCount != uint8(len(msg.Fees)) || msg.PoolsCount != uint8(len(msg.Tokens) - 1) || msg.PoolsCount != uint8(len(newPoolBalances)) {
		return nil, fmt.Errorf("params array length mismatch with pools count")
	}
	// the first inputs should be the pool's current sequence number that contain pool's liquidity
	for idx, input := range inputs[:msg.PoolsCount] {
		if !bytes.Equal(input.PkScript, relayerPkScript) {
			return nil, fmt.Errorf("the input %v should be from relayer wallet", idx)
		}
	}
	// the remain inputs should be from trading wallet
	for idx, input := range inputs[msg.PoolsCount:] {
		if !bytes.Equal(input.PkScript, userPkScript) {
			return nil, fmt.Errorf("the input %v should be from trading wallet", int(msg.PoolsCount)+idx)
		}
	}
	radfiScript, _ := CreateSwapScript(msg)

	runeOutput := &runestone.Runestone{
		Edicts: []runestone.Edict{},
		Pointer: &msg.TokenOutIndex,
	}
	// sequence number outputs
	outputs := []*wire.TxOut{}
	for idx, poolBalance := range newPoolBalances {
		sequenceNumberAmount := DUST_UTXO_AMOUNT
		if poolBalance.Token1Id == BitcoinRuneId() {
			return nil, fmt.Errorf("the second token in the pair can only be rune")
		}
		if poolBalance.Token0Id == BitcoinRuneId() {
			sequenceNumberAmount = int64(poolBalance.Token0Amount.Lo)

		} else {
			runeOutput.Edicts = append(runeOutput.Edicts, runestone.Edict{
				ID:	poolBalance.Token0Id,
				Amount: poolBalance.Token0Amount,
				Output: uint32(idx),
			})
		}
		runeOutput.Edicts = append(runeOutput.Edicts, runestone.Edict{
			ID:	poolBalance.Token1Id,
			Amount: poolBalance.Token1Amount,
			Output: uint32(idx),
		})
		outputs = append(outputs, &wire.TxOut{
			Value: sequenceNumberAmount,
			PkScript: relayerPkScript,
		})
	}
	runeScript, _ := runeOutput.Encipher()

	userTokenOutBitcoin := DUST_UTXO_AMOUNT
	if *msg.Tokens[len(msg.Tokens)-1] == BitcoinRuneId() {
		userTokenOutBitcoin = int64(msg.AmountOut.Lo)
	}
	// add remain outputs
	outputs = append(outputs, []*wire.TxOut{
		// radfi OP
		{
			Value: DUST_UTXO_AMOUNT,
			PkScript: radfiScript,
		},
		// rune OP_RETURN
		{
			Value: 0,
			PkScript: runeScript,
		},
		// user token out output
		{
			Value: userTokenOutBitcoin,
			PkScript: userPkScript,
		},
	}...)

	return CreateRadFiTx(inputs, outputs, userPkScript, txFee, 0)
}

func CreateRadFiTxWithdrawLiquidity(
	msg *RadFiWithdrawLiquidityMsg,
	inputs []*Input,
	relayerPkScript []byte,
	userPkScript []byte,
	txFee int64,
) (*wire.MsgTx, error) {
	// the first input should be the pool's current sequence number that contain pool's liquidity
	if !bytes.Equal(inputs[0].PkScript, relayerPkScript) {
		return nil, fmt.Errorf("the first input should be the pool's current sequence number")
	}
	// the remain inputs should be from trading wallet
	for idx, input := range inputs[1:] {
		if !bytes.Equal(input.PkScript, userPkScript) {
			return nil, fmt.Errorf("the input %v should be from trading wallet", idx)
		}
	}

	radfiScript, _ := CreateWithdrawLiquidityScript(msg)

	withdrawedRuneOutput := uint32(3)
	runeOutput := &runestone.Runestone{
		Edicts: []runestone.Edict{},
		Pointer: &withdrawedRuneOutput,
	}

	sequenceNumberAmount := DUST_UTXO_AMOUNT
	if len(inputs[0].Runes) == 1 {
		sequenceNumberAmount = inputs[0].OutputAmount - int64(msg.Amount0.Lo)
		runeOutput.Edicts = append(runeOutput.Edicts, runestone.Edict{
			ID:	inputs[0].Runes[0].ID,
			Amount: inputs[0].Runes[0].Amount.Sub(msg.Amount1),
			Output: 0,
		})
	} else {
		if len(inputs[0].Runes) != 2 {
			return nil, fmt.Errorf("rune-rune pool sequence number UTXO should hold exactly 2 rune")
		}
		runeOutput.Edicts = append(runeOutput.Edicts, runestone.Edict{
			ID:	inputs[0].Runes[0].ID,
			Amount: inputs[0].Runes[0].Amount.Sub(msg.Amount0),
			Output: 0,
		})
		runeOutput.Edicts = append(runeOutput.Edicts, runestone.Edict{
			ID:	inputs[0].Runes[1].ID,
			Amount: inputs[0].Runes[1].Amount.Sub(msg.Amount1),
			Output: 0,
		})
	}
	runeScript, _ := runeOutput.Encipher()

	outputs := []*wire.TxOut{
		// sequence number output
		{
			Value: sequenceNumberAmount,
			PkScript: relayerPkScript,
		},
		// radfi OP
		{
			Value: DUST_UTXO_AMOUNT,
			PkScript: radfiScript,
		},
		// rune OP_RETURN
		{
			Value: 0,
			PkScript: runeScript,
		},
		// withdrawed rune output
		{
			Value: DUST_UTXO_AMOUNT,
			PkScript: userPkScript,
		},
	}
	// withdrawed bitcoin output is included in the change output
	return CreateRadFiTx(inputs, outputs, userPkScript, txFee, 0)
}

func CreateRadFiTxCollectFees(
	msg *RadFiCollectFeesMsg,
	inputs []*Input,
	relayerPkScript []byte,
	userPkScript []byte,
	txFee int64,
) (*wire.MsgTx, error) {
	// the first input should be the pool's current sequence number that contain pool's liquidity
	if !bytes.Equal(inputs[0].PkScript, relayerPkScript) {
		return nil, fmt.Errorf("the first input should be the pool's current sequence number")
	}
	// the remain inputs should be from trading wallet
	for idx, input := range inputs[1:] {
		if !bytes.Equal(input.PkScript, userPkScript) {
			return nil, fmt.Errorf("the input %v should be from trading wallet", idx)
		}
	}

	radfiScript, _ := CreateCollectFeesScript(msg)

	withdrawedRuneOutput := uint32(3)
	runeOutput := &runestone.Runestone{
		Edicts: []runestone.Edict{},
		Pointer: &withdrawedRuneOutput,
	}

	sequenceNumberAmount := DUST_UTXO_AMOUNT
	if len(inputs[0].Runes) == 1 {
		sequenceNumberAmount = inputs[0].OutputAmount - int64(msg.Amount0.Lo)
		runeOutput.Edicts = append(runeOutput.Edicts, runestone.Edict{
			ID:	inputs[0].Runes[0].ID,
			Amount: inputs[0].Runes[0].Amount.Sub(msg.Amount1),
			Output: 0,
		})
	} else {
		if len(inputs[0].Runes) != 2 {
			return nil, fmt.Errorf("rune-rune pool sequence number UTXO should hold exactly 2 rune")
		}
		runeOutput.Edicts = append(runeOutput.Edicts, runestone.Edict{
			ID:	inputs[0].Runes[0].ID,
			Amount: inputs[0].Runes[0].Amount.Sub(msg.Amount0),
			Output: 0,
		})
		runeOutput.Edicts = append(runeOutput.Edicts, runestone.Edict{
			ID:	inputs[0].Runes[1].ID,
			Amount: inputs[0].Runes[1].Amount.Sub(msg.Amount1),
			Output: 0,
		})
	}
	runeScript, _ := runeOutput.Encipher()

	outputs := []*wire.TxOut{
		// sequence number output
		{
			Value: sequenceNumberAmount,
			PkScript: relayerPkScript,
		},
		// radfi OP
		{
			Value: DUST_UTXO_AMOUNT,
			PkScript: radfiScript,
		},
		// rune OP_RETURN
		{
			Value: 0,
			PkScript: runeScript,
		},
		// withdrawed rune output
		{
			Value: DUST_UTXO_AMOUNT,
			PkScript: userPkScript,
		},
	}
	// withdrawed bitcoin output is included in the change output
	return CreateRadFiTx(inputs, outputs, userPkScript, txFee, 0)
}

func CreateRadFiTxIncreaseLiquidity(
	msg *RadFiIncreaseLiquidityMsg,
	inputs []*Input,
	relayerPkScript []byte,
	userPkScript []byte,
	txFee int64,
) (*wire.MsgTx, error) {
	// the first input should be the pool's current sequence number that contain pool's liquidity
	if !bytes.Equal(inputs[0].PkScript, relayerPkScript) {
		return nil, fmt.Errorf("the first input should be the pool's current sequence number")
	}
	// the remain inputs should be from trading wallet
	for idx, input := range inputs[1:] {
		if !bytes.Equal(input.PkScript, userPkScript) {
			return nil, fmt.Errorf("the input %v should be from trading wallet", idx)
		}
	}

	radfiScript, _ := CreateIncreaseLiquidityScript(msg)

	userChangeOutput := uint32(3)
	runeOutput := &runestone.Runestone{
		Edicts: []runestone.Edict{},
		Pointer: &userChangeOutput,
	}

	sequenceNumberAmount := DUST_UTXO_AMOUNT
	if len(inputs[0].Runes) == 1 {
		sequenceNumberAmount = inputs[0].OutputAmount + int64(msg.Amount0.Lo)
		runeOutput.Edicts = append(runeOutput.Edicts, runestone.Edict{
			ID:	inputs[0].Runes[0].ID,
			Amount: inputs[0].Runes[0].Amount.Add(msg.Amount1),
			Output: 0,
		})
	} else {
		if len(inputs[0].Runes) != 2 {
			return nil, fmt.Errorf("rune-rune pool sequence number UTXO should hold exactly 2 rune")
		}
		runeOutput.Edicts = append(runeOutput.Edicts, runestone.Edict{
			ID:	inputs[0].Runes[0].ID,
			Amount: inputs[0].Runes[0].Amount.Add(msg.Amount0),
			Output: 0,
		})
		runeOutput.Edicts = append(runeOutput.Edicts, runestone.Edict{
			ID:	inputs[0].Runes[1].ID,
			Amount: inputs[0].Runes[1].Amount.Add(msg.Amount1),
			Output: 0,
		})
	}
	runeScript, _ := runeOutput.Encipher()

	outputs := []*wire.TxOut{
		// sequence number output
		{
			Value: sequenceNumberAmount,
			PkScript: relayerPkScript,
		},
		// radfi OP
		{
			Value: DUST_UTXO_AMOUNT,
			PkScript: radfiScript,
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

	return CreateRadFiTx(inputs, outputs, userPkScript, txFee, 0)
}
