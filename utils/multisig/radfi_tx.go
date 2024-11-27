package multisig

import (
	"bytes"
	"fmt"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/bxelab/runestone"
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
	txFee int64,
) (*wire.MsgTx, error) {
	if msg.InitPrice.Lo == 0 {
		return nil, fmt.Errorf("init price should be > 0")
	}
	if msg.Token1Id == BitcoinRuneId() {
		return nil, fmt.Errorf("the second token in the pair can only be rune")
	}
	// the inputs should be from trading wallet
	for idx, input := range inputs {
		if !bytes.Equal(input.PkScript, userPkScript) {
			return nil, fmt.Errorf("the input %v should be from trading wallet", idx)
		}
	}

	radfiScript, _ := CreateProvideLiquidityScript(msg)

	userChangeOutput := uint32(3)
	runeOutput := &runestone.Runestone{
		Edicts: []runestone.Edict{},
		Pointer: &userChangeOutput,
	}

	sequenceNumberAmount := DUST_UTXO_AMOUNT
	if msg.Token0Id == BitcoinRuneId() {
		sequenceNumberAmount = int64(msg.Amount0Desired.Lo)
	} else {
		runeOutput.Edicts = append(runeOutput.Edicts, runestone.Edict{
			ID:	msg.Token0Id,
			Amount: msg.Amount0Desired,
			Output: 0,
		})
	}
	runeOutput.Edicts = append(runeOutput.Edicts, runestone.Edict{
		ID:	msg.Token1Id,
		Amount: msg.Amount1Desired,
		Output: 0,
	})
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

func CreateRadFiTxProvideLiquidity(
	msg *RadFiProvideLiquidityMsg,
	inputs []*Input,
	relayerPkScript []byte,
	userPkScript []byte,
	txFee int64,
) (*wire.MsgTx, error) {
	if msg.InitPrice.Lo != 0 {
		return nil, fmt.Errorf("init price should be 0")
	}
	if msg.Token1Id == BitcoinRuneId() {
		return nil, fmt.Errorf("the second token in the pair can only be rune")
	}
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

	radfiScript, _ := CreateProvideLiquidityScript(msg)

	userChangeOutput := uint32(3)
	runeOutput := &runestone.Runestone{
		Edicts: []runestone.Edict{},
		Pointer: &userChangeOutput,
	}

	sequenceNumberAmount := DUST_UTXO_AMOUNT
	if msg.Token0Id == BitcoinRuneId() {
		if len(inputs[0].Runes) != 1 {
			return nil, fmt.Errorf("bitcoin-rune pool sequence number UTXO should hold exactly 1 rune")
		}
		sequenceNumberAmount = inputs[0].OutputAmount + int64(msg.Amount0Desired.Lo)
		runeOutput.Edicts = append(runeOutput.Edicts, runestone.Edict{
			ID:	msg.Token1Id,
			Amount: inputs[0].Runes[0].Amount.Add(msg.Amount1Desired),
			Output: 0,
		})
	} else {
		if len(inputs[0].Runes) != 2 {
			return nil, fmt.Errorf("rune-rune pool sequence number UTXO should hold exactly 2 rune")
		}
		runeOutput.Edicts = append(runeOutput.Edicts, runestone.Edict{
			ID:	msg.Token0Id,
			Amount: inputs[0].Runes[0].Amount.Add(msg.Amount0Desired),
			Output: 0,
		})
		runeOutput.Edicts = append(runeOutput.Edicts, runestone.Edict{
			ID:	msg.Token1Id,
			Amount: inputs[0].Runes[1].Amount.Add(msg.Amount1Desired),
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
