package bitcoin

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"testing"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/bxelab/runestone"
	"github.com/icon-project/centralized-relay/relayer/events"
	"github.com/icon-project/centralized-relay/relayer/types"
	"github.com/icon-project/centralized-relay/utils/multisig"
	"github.com/stretchr/testify/assert"
	"github.com/syndtr/goleveldb/leveldb"
	"go.uber.org/zap"
)

const (
	TX_FEE                   = 10000
	RELAYER_MULTISIG_ADDRESS = "tb1pf0atpt2d3zel6udws38pkrh2e49vqd3c5jcud3a82srphnmpe55q0ecrzk"
	USER_MULTISIG_ADDRESS    = "tb1pgzx880yfr7q8dgz8dqhw50sncu4f4hmw5cn3800354tuzcy9jx5shvv7su"
)

func TestDecodeWithdrawToMessage(t *testing.T) {
	input := "+QEdAbkBGfkBFrMweDIuaWNvbi9jeGZjODZlZTc2ODdlMWJmNjgxYjU1NDhiMjY2Nzg0NDQ4NWMwZTcxOTK4PnRiMXBneng4ODB5ZnI3cThkZ3o4ZHFodzUwc25jdTRmNGhtdzVjbjM4MDAzNTR0dXpjeTlqeDVzaHZ2N3N1gh6FAbhS+FCKV2l0aGRyYXdUb4MwOjC4PnRiMXBneng4ODB5ZnI3cThkZ3o4ZHFodzUwc25jdTRmNGhtdzVjbjM4MDAzNTR0dXpjeTlqeDVzaHZ2N3N1ZPhIuEYweDIuYnRjL3RiMXBneng4ODB5ZnI3cThkZ3o4ZHFodzUwc25jdTRmNGhtdzVjbjM4MDAzNTR0dXpjeTlqeDVzaHZ2N3N1"
	// Decode base64
	decodedBytes, _ := base64.StdEncoding.DecodeString(input)

	result, data, err := decodeWithdrawToMessage(decodedBytes)

	fmt.Println("data", data)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "tb1pgzx880yfr7q8dgz8dqhw50sncu4f4hmw5cn3800354tuzcy9jx5shvv7su", result.To)
	assert.Equal(t, big.NewInt(100).Bytes(), result.Amount)
	assert.Equal(t, "WithdrawTo", result.Action)
	assert.Equal(t, "0:0", result.TokenAddress)
}

func TestProvider_Route(t *testing.T) {
	// Setup
	tempDir, err := os.MkdirTemp("", "bitcoin_provider_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "test.db")
	db, err := leveldb.OpenFile(dbPath, nil)
	assert.NoError(t, err)
	defer db.Close()

	logger, _ := zap.NewDevelopment()
	provider := &Provider{
		logger: logger,
		cfg:    &Config{Mode: SlaveMode},
		db:     db,
	}

	// Create a test message
	testMessage := &types.Message{
		Dst:           "destination",
		Src:           "source",
		Sn:            big.NewInt(123),
		Data:          []byte("test data"),
		MessageHeight: 456,
		EventType:     events.EmitMessage,
	}

	// Test storing the message
	err = provider.Route(context.Background(), testMessage, nil)
	assert.NoError(t, err)

	// Test retrieving the message
	key := []byte(fmt.Sprintf("bitcoin_message_%s", testMessage.Sn.String()))
	storedData, err := db.Get(key, nil)
	assert.NoError(t, err)

	var retrievedMessage types.Message
	err = json.Unmarshal(storedData, &retrievedMessage)
	assert.NoError(t, err)

	assert.Equal(t, testMessage.Dst, retrievedMessage.Dst)
	assert.Equal(t, testMessage.Src, retrievedMessage.Src)
	assert.Equal(t, testMessage.Sn.String(), retrievedMessage.Sn.String())
	assert.Equal(t, testMessage.Data, retrievedMessage.Data)
	assert.Equal(t, testMessage.MessageHeight, retrievedMessage.MessageHeight)
	assert.Equal(t, testMessage.EventType, retrievedMessage.EventType)

	// Test deleting the message
	err = db.Delete(key, nil)
	assert.NoError(t, err)

	_, err = db.Get(key, nil)
	assert.Error(t, err) // Should return an error as the key no longer exists
}

func TestDepositBitcoinToIcon(t *testing.T) {
	chainParam := &chaincfg.TestNet3Params
	// relayer multisig
	decodedAddr, _ := btcutil.DecodeAddress(RELAYER_MULTISIG_ADDRESS, chainParam)
	relayerPkScript, _ := txscript.PayToAddrScript(decodedAddr)
	// user key
	userPrivKeys, userMultisigInfo := multisig.RandomMultisigInfo(2, 2, chainParam, []int{0, 3}, 1, 1)
	userMultisigWallet, _ := multisig.BuildMultisigWallet(userMultisigInfo)

	bridgeMsg := multisig.BridgeDecodedMsg{
		Message: &multisig.XCallMessage{
			MessageType:  1,
			Action:       "Deposit",
			TokenAddress: "0:1",
			To:           "0x2.icon/hx452e235f9f1fd1006b1941ed1ad19ef51d1192f6",
			From:         "tb1pgzx880yfr7q8dgz8dqhw50sncu4f4hmw5cn3800354tuzcy9jx5shvv7su",
			Amount:       new(big.Int).SetUint64(100000).Bytes(),
			Data:         []byte(""),
		},
		ChainId:  1,
		Receiver: "cxfc86ee7687e1bf681b5548b2667844485c0e7192",
		Connectors: []string{
			"cx577f5e756abd89cbcba38a58508b60a12754d2f5",
		},
	}

	inputs := []*multisig.Input{
		{
			TxHash:       "e57c10e27f75dbf0856163ca5f825b5af7ffbb3874f606b31330464ddd9df9a1",
			OutputIdx:    4,
			OutputAmount: 2974000,
			PkScript:     userMultisigWallet.PKScript,
		},
	}

	// create tx
	msgTx, err := multisig.CreateBridgeTxSendBitcoin(
		&bridgeMsg,
		inputs,
		userMultisigWallet.PKScript,
		relayerPkScript,
		TX_FEE,
	)
	fmt.Println("err: ", err)
	// sign tx
	totalSigs := [][][]byte{}
	// user key 1 sign tx
	userSigs1, _ := multisig.SignTapMultisig(userPrivKeys[0], msgTx, inputs, userMultisigWallet, 0)
	totalSigs = append(totalSigs, userSigs1)
	// user key 2 sign tx
	userSigs2, _ := multisig.SignTapMultisig(userPrivKeys[1], msgTx, inputs, userMultisigWallet, 0)
	totalSigs = append(totalSigs, userSigs2)
	// COMBINE SIGN
	signedMsgTx, _ := multisig.CombineTapMultisig(totalSigs, msgTx, inputs, userMultisigWallet, 0)

	var signedTx bytes.Buffer
	signedMsgTx.Serialize(&signedTx)
	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	signedMsgTxID := signedMsgTx.TxHash().String()

	fmt.Println("hexSignedTx: ", hexSignedTx)
	fmt.Println("signedMsgTxID: ", signedMsgTxID)
	fmt.Println("err sign: ", err)

	// TODO: test the signedMsgTx
}

func TestDepositBitcoinToIconFail1(t *testing.T) {
	chainParam := &chaincfg.TestNet3Params
	// user key
	userPrivKeys, userMultisigInfo := multisig.RandomMultisigInfo(2, 2, chainParam, []int{0, 3}, 1, 1)
	userMultisigWallet, _ := multisig.BuildMultisigWallet(userMultisigInfo)

	inputs := []*multisig.Input{
		{
			TxHash:       "eeb8c9f79ecfe7c084b2af95bf82acebd130185a0d188283d78abb58d85eddff",
			OutputIdx:    4,
			OutputAmount: 2999000,
			PkScript:     userMultisigWallet.PKScript,
		},
	}

	outputs := []*wire.TxOut{}

	// Add Bridge Message
	scripts, _ := multisig.CreateBridgeMessageScripts(&multisig.BridgeDecodedMsg{
		Message: &multisig.XCallMessage{
			Action:       "Deposit",
			TokenAddress: "0:1",
			To:           "0x2.icon/hx452e235f9f1fd1006b1941ed1ad19ef51d1192f6",
			From:         "tb1pgzx880yfr7q8dgz8dqhw50sncu4f4hmw5cn3800354tuzcy9jx5shvv7su",
			Amount:       new(big.Int).SetUint64(1000000).Bytes(),
			Data:         []byte(""),
		},
		ChainId:  1,
		Receiver: "cx8b52dfea0aa1e548288102df15ad7159f7266106",
		Connectors: []string{
			"cx577f5e756abd89cbcba38a58508b60a12754d2f5",
		},
	})
	for _, script := range scripts {
		outputs = append(outputs, &wire.TxOut{
			Value:    0,
			PkScript: script,
		})
	}

	msgTx, err := multisig.CreateTx(inputs, outputs, userMultisigWallet.PKScript, TX_FEE, 0)
	fmt.Println("err: ", err)
	// sign tx
	totalSigs := [][][]byte{}
	// user key 1 sign tx
	userSigs1, _ := multisig.SignTapMultisig(userPrivKeys[0], msgTx, inputs, userMultisigWallet, 0)
	totalSigs = append(totalSigs, userSigs1)
	// user key 2 sign tx
	userSigs2, _ := multisig.SignTapMultisig(userPrivKeys[1], msgTx, inputs, userMultisigWallet, 0)
	totalSigs = append(totalSigs, userSigs2)
	// COMBINE SIGN
	signedMsgTx, _ := multisig.CombineTapMultisig(totalSigs, msgTx, inputs, userMultisigWallet, 0)

	var signedTx bytes.Buffer
	signedMsgTx.Serialize(&signedTx)
	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	signedMsgTxID := signedMsgTx.TxHash().String()

	fmt.Println("hexSignedTx: ", hexSignedTx)
	fmt.Println("signedMsgTxID: ", signedMsgTxID)
	fmt.Println("err sign: ", err)

	// TODO: test the signedMsgTx
}

func TestDepositBitcoinToIconFail2(t *testing.T) {
	chainParam := &chaincfg.TestNet3Params
	// relayer multisig
	decodedAddr, _ := btcutil.DecodeAddress(RELAYER_MULTISIG_ADDRESS, chainParam)
	relayerPkScript, _ := txscript.PayToAddrScript(decodedAddr)
	// user key
	userPrivKeys, userMultisigInfo := multisig.RandomMultisigInfo(2, 2, chainParam, []int{0, 3}, 1, 1)
	userMultisigWallet, _ := multisig.BuildMultisigWallet(userMultisigInfo)

	inputs := []*multisig.Input{
		{
			TxHash:       "0416795b227e1a6a64eeb7bf7542d15964d18ac4c4732675d3189cda8d38bed7",
			OutputIdx:    3,
			OutputAmount: 2998000,
			PkScript:     userMultisigWallet.PKScript,
		},
	}

	outputs := []*wire.TxOut{}

	// Add Bridge Message
	scripts, _ := multisig.CreateBridgeMessageScripts(&multisig.BridgeDecodedMsg{
		Message: &multisig.XCallMessage{
			Action:       "Deposit",
			TokenAddress: "0:1", // wrong token address
			To:           "0x2.icon/hx452e235f9f1fd1006b1941ed1ad19ef51d1192f6",
			From:         "tb1pgzx880yfr7q8dgz8dqhw50sncu4f4hmw5cn3800354tuzcy9jx5shvv7su",
			Amount:       new(big.Int).SetUint64(1000000).Bytes(), //
			Data:         []byte(""),
		},
		ChainId:  1,
		Receiver: "cx8b52dfea0aa1e548288102df15ad7159f7266106",
		Connectors: []string{
			"cx577f5e756abd89cbcba38a58508b60a12754d2f5",
		},
	})
	for _, script := range scripts {
		outputs = append(outputs, &wire.TxOut{
			Value:    0,
			PkScript: script,
		})
	}

	// Add transfering bitcoin to relayer multisig
	outputs = append(outputs, &wire.TxOut{
		Value:    1000,
		PkScript: relayerPkScript,
	})

	msgTx, err := multisig.CreateTx(inputs, outputs, userMultisigWallet.PKScript, TX_FEE, 0)
	fmt.Println("err: ", err)
	// sign tx
	totalSigs := [][][]byte{}
	// user key 1 sign tx
	userSigs1, _ := multisig.SignTapMultisig(userPrivKeys[0], msgTx, inputs, userMultisigWallet, 0)
	totalSigs = append(totalSigs, userSigs1)
	// user key 2 sign tx
	userSigs2, _ := multisig.SignTapMultisig(userPrivKeys[1], msgTx, inputs, userMultisigWallet, 0)
	totalSigs = append(totalSigs, userSigs2)
	// COMBINE SIGN
	signedMsgTx, _ := multisig.CombineTapMultisig(totalSigs, msgTx, inputs, userMultisigWallet, 0)

	var signedTx bytes.Buffer
	signedMsgTx.Serialize(&signedTx)
	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	signedMsgTxID := signedMsgTx.TxHash().String()

	fmt.Println("hexSignedTx: ", hexSignedTx)
	fmt.Println("signedMsgTxID: ", signedMsgTxID)
	fmt.Println("err sign: ", err)

	// TODO: test the signedMsgTx
}

func TestDepositBitcoinToIconFail3(t *testing.T) {
	chainParam := &chaincfg.TestNet3Params
	// relayer multisig
	decodedAddr, _ := btcutil.DecodeAddress(RELAYER_MULTISIG_ADDRESS, chainParam)
	relayerPkScript, _ := txscript.PayToAddrScript(decodedAddr)
	// user key
	userPrivKeys, userMultisigInfo := multisig.RandomMultisigInfo(2, 2, chainParam, []int{0, 3}, 1, 1)
	userMultisigWallet, _ := multisig.BuildMultisigWallet(userMultisigInfo)

	inputs := []*multisig.Input{
		{
			TxHash:       "dc21f89436d9fbda2cc521ed9b8988c7cbf84cdc67d728b2b2709a5efe7e775a",
			OutputIdx:    4,
			OutputAmount: 2996000,
			PkScript:     userMultisigWallet.PKScript,
		},
	}

	outputs := []*wire.TxOut{}

	// Add Bridge Message
	scripts, _ := multisig.CreateBridgeMessageScripts(&multisig.BridgeDecodedMsg{
		Message: &multisig.XCallMessage{
			Action:       "Deposit",
			TokenAddress: "0:1",
			To:           "0x2.icon/hx452e235f9f1fd1006b1941ed1ad19ef51d1192f6",
			From:         "tb1pgzx880yfr7q8dgz8dqhw50sncu4f4hmw5cn3800354tuzcy9jx5shvv7su",
			Amount:       new(big.Int).SetUint64(1000).Bytes(), // wrong amount
			Data:         []byte(""),
		},
		ChainId:  1,
		Receiver: "cx8b52dfea0aa1e548288102df15ad7159f7266106",
		Connectors: []string{
			"cx577f5e756abd89cbcba38a58508b60a12754d2f5",
		},
	})
	for _, script := range scripts {
		outputs = append(outputs, &wire.TxOut{
			Value:    0,
			PkScript: script,
		})
	}

	// Add transfering bitcoin to relayer multisig
	outputs = append(outputs, &wire.TxOut{
		Value:    10000,
		PkScript: relayerPkScript,
	})

	msgTx, err := multisig.CreateTx(inputs, outputs, userMultisigWallet.PKScript, TX_FEE, 0)
	fmt.Println("err: ", err)
	// sign tx
	totalSigs := [][][]byte{}
	// user key 1 sign tx
	userSigs1, _ := multisig.SignTapMultisig(userPrivKeys[0], msgTx, inputs, userMultisigWallet, 0)
	totalSigs = append(totalSigs, userSigs1)
	// user key 2 sign tx
	userSigs2, _ := multisig.SignTapMultisig(userPrivKeys[1], msgTx, inputs, userMultisigWallet, 0)
	totalSigs = append(totalSigs, userSigs2)
	// COMBINE SIGN
	signedMsgTx, _ := multisig.CombineTapMultisig(totalSigs, msgTx, inputs, userMultisigWallet, 0)

	var signedTx bytes.Buffer
	signedMsgTx.Serialize(&signedTx)
	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	signedMsgTxID := signedMsgTx.TxHash().String()

	fmt.Println("hexSignedTx: ", hexSignedTx)
	fmt.Println("signedMsgTxID: ", signedMsgTxID)
	fmt.Println("err sign: ", err)

	// TODO: test the signedMsgTx
}

func TestDepositBitcoinToIconFail4(t *testing.T) {
	chainParam := &chaincfg.TestNet3Params
	// relayer multisig
	decodedAddr, _ := btcutil.DecodeAddress(RELAYER_MULTISIG_ADDRESS, chainParam)
	relayerPkScript, _ := txscript.PayToAddrScript(decodedAddr)
	// user key
	userPrivKeys, userMultisigInfo := multisig.RandomMultisigInfo(2, 2, chainParam, []int{0, 3}, 1, 1)
	userMultisigWallet, _ := multisig.BuildMultisigWallet(userMultisigInfo)

	inputs := []*multisig.Input{
		{
			TxHash:       "1e29fa62942f92dd4cf688e219641b54229fbc2ec2dc74cc9c1c7f247c7172b2",
			OutputIdx:    4,
			OutputAmount: 2985000,
			PkScript:     userMultisigWallet.PKScript,
		},
	}

	outputs := []*wire.TxOut{}

	// Add Bridge Message
	scripts, _ := multisig.CreateBridgeMessageScripts(&multisig.BridgeDecodedMsg{
		Message: &multisig.XCallMessage{
			Action:       "Withdraw",
			TokenAddress: "0:1",
			To:           "0x2.icon/hx452e235f9f1fd1006b1941ed1ad19ef51d1192f6",
			From:         "tb1pgzx880yfr7q8dgz8dqhw50sncu4f4hmw5cn3800354tuzcy9jx5shvv7su",
			Amount:       new(big.Int).SetUint64(1000).Bytes(),
			Data:         []byte(""),
		},
		ChainId:  1,
		Receiver: "cx8b52dfea0aa1e548288102df15ad7159f7266106",
		Connectors: []string{
			"cx577f5e756abd89cbcba38a58508b60a12754d2f5",
		},
	})
	for _, script := range scripts {
		outputs = append(outputs, &wire.TxOut{
			Value:    0,
			PkScript: script,
		})
	}

	// Add transfering bitcoin to relayer multisig
	outputs = append(outputs, &wire.TxOut{
		Value:    10000,
		PkScript: relayerPkScript,
	})

	msgTx, err := multisig.CreateTx(inputs, outputs, userMultisigWallet.PKScript, TX_FEE, 0)
	fmt.Println("err: ", err)
	// sign tx
	totalSigs := [][][]byte{}
	// user key 1 sign tx
	userSigs1, _ := multisig.SignTapMultisig(userPrivKeys[0], msgTx, inputs, userMultisigWallet, 0)
	totalSigs = append(totalSigs, userSigs1)
	// user key 2 sign tx
	userSigs2, _ := multisig.SignTapMultisig(userPrivKeys[1], msgTx, inputs, userMultisigWallet, 0)
	totalSigs = append(totalSigs, userSigs2)
	// COMBINE SIGN
	signedMsgTx, _ := multisig.CombineTapMultisig(totalSigs, msgTx, inputs, userMultisigWallet, 0)

	var signedTx bytes.Buffer
	signedMsgTx.Serialize(&signedTx)
	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	signedMsgTxID := signedMsgTx.TxHash().String()

	fmt.Println("hexSignedTx: ", hexSignedTx)
	fmt.Println("signedMsgTxID: ", signedMsgTxID)
	fmt.Println("err sign: ", err)

	// TODO: test the signedMsgTx
}

// ... existing code ...

func TestDecodeBitcoinTransaction(t *testing.T) {
	// The raw transaction hex string
	rawTxHex := "02000000000101b272717c247f1c9ccc74dcc22ebc9f22541b6419e288f64cdd922f9462fa291e040000000001000000050000000000000000506a5e4c4cf88588576974686472617783303a31b83e74623170677a7838383079667237713864677a38647168773530736e6375346634686d7735636e3338303033353474757a6379396a7835736876760000000000000000506a5e4c4c377375b33078322e69636f6e2f6878343532653233356639663166643130303662313934316564316164313965663531643131393266368203e88001738b52dfea0aa1e548288102df15ad7100000000000000001d6a5e1a59f726610673577f5e756abd89cbcba38a58508b60a12754d2f510270000000000002251204bfab0ad4d88b3fd71ae844e1b0eeacd4ac03638a4b1c6c7a754061bcf61cd2830612d0000000000225120408c73bc891f8076a047682eea3e13c72a9adf6ea62713bdf1a557c1608591a903405a713aad72e1a2717cab16446e84a0de1c9f908a100c5903086a97910adc1327db65d21c8533c4fe087027195db4ebe7d8595624382e39b7aba829e3b29a29dc2551b275207303e7756826cf3fabc2f9c06978c542039effbb7493627cad22236e3ff10ee4ac41c18a23958fc9bf526c81a09bca529d685adc6900a04e2f520a26c63aa0b61a770f27c71b203eebab28e2e37b992abd8e6e5f9d887bdc2d5ab0efde76df79d3520400000000"

	// Decode the raw transaction
	txBytes, err := hex.DecodeString(rawTxHex)
	assert.NoError(t, err)

	var tx wire.MsgTx
	err = tx.Deserialize(bytes.NewReader(txBytes))
	assert.NoError(t, err)

	// Extract sender information (from the input)
	assert.Equal(t, 1, len(tx.TxIn), "Expected 1 input")
	prevOutHash := tx.TxIn[0].PreviousOutPoint.Hash.String()
	prevOutIndex := tx.TxIn[0].PreviousOutPoint.Index
	fmt.Printf("Sender (Input): %s:%d\n", prevOutHash, prevOutIndex)

	// Extract receiver information (from the outputs)
	assert.Equal(t, 5, len(tx.TxOut), "Expected 5 outputs")

	for i, out := range tx.TxOut {
		fmt.Printf("Output %d:\n", i)
		fmt.Printf("  Amount: %d satoshis\n", out.Value)

		// Attempt to parse the output script
		scriptClass, addresses, _, err := txscript.ExtractPkScriptAddrs(out.PkScript, &chaincfg.TestNet3Params)
		if err != nil {
			fmt.Printf("  Script: Unable to parse (possibly OP_RETURN)\n")
		} else {
			fmt.Printf("  Script Class: %s\n", scriptClass)
			if len(addresses) > 0 {
				fmt.Printf("  Receiver Address: %s\n", addresses[0].String())
			}
		}

		// If it's an OP_RETURN output, print the data
		if scriptClass == txscript.NullDataTy {
			fmt.Printf("  OP_RETURN Data: %x\n", out.PkScript[2:])
		}

		fmt.Println()
	}

	// Add assertions for specific outputs
	assert.Equal(t, txscript.NullDataTy, txscript.GetScriptClass(tx.TxOut[0].PkScript))
	assert.Equal(t, txscript.NullDataTy, txscript.GetScriptClass(tx.TxOut[1].PkScript))
	assert.Equal(t, txscript.NullDataTy, txscript.GetScriptClass(tx.TxOut[2].PkScript))

	_, addresses, _, err := txscript.ExtractPkScriptAddrs(tx.TxOut[3].PkScript, &chaincfg.TestNet3Params)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(addresses))
	assert.Equal(t, "tb1pqqqqp399et2xygdj5xreqhjjvcmzhxw4aywxecjdzew6hylgvsesf3hn0c", addresses[0].String())

	_, addresses, _, err = txscript.ExtractPkScriptAddrs(tx.TxOut[4].PkScript, &chaincfg.TestNet3Params)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(addresses))
	assert.Equal(t, "tb1pqqqqp399et2xygdj5xreqhjjvcmzhxw4aywxecjdzew6hylgvsesf3hn0c", addresses[0].String())
}

func TestDepositRuneToIcon(t *testing.T) {
	chainParam := &chaincfg.TestNet3Params
	// relayer multisig
	decodedAddr, _ := btcutil.DecodeAddress(RELAYER_MULTISIG_ADDRESS, chainParam)
	relayerPkScript, _ := txscript.PayToAddrScript(decodedAddr)
	// user key
	userPrivKeys, userMultisigInfo := multisig.RandomMultisigInfo(2, 2, chainParam, []int{0, 3}, 1, 1)
	userMultisigWallet, _ := multisig.BuildMultisigWallet(userMultisigInfo)

	bridgeMsg := multisig.BridgeDecodedMsg{
		Message: &multisig.XCallMessage{
			MessageType:  1,
			Action:       "Deposit",
			TokenAddress: "2904354:3119",
			To:           "0x2.icon/hx1493794ba31fa3372bf7903f04030497e7d14800",
			From:         "tb1pgzx880yfr7q8dgz8dqhw50sncu4f4hmw5cn3800354tuzcy9jx5shvv7su",
			Amount:       new(big.Int).SetUint64(10000).Bytes(),
			Data:         []byte(""),
		},
		ChainId:  1,
		Receiver: "cxfc86ee7687e1bf681b5548b2667844485c0e7192",
		Connectors: []string{
			"cx577f5e756abd89cbcba38a58508b60a12754d2f5",
		},
	}

	inputs := []*multisig.Input{
		// user rune UTXOs to spend
		{
			TxHash:       "69deba39f5a0700cc713f67fe8cb5ed1e35a9f0d4a3a437d839103c6e26cb947",
			OutputIdx:    2,
			OutputAmount: 546,
			PkScript:     userMultisigWallet.PKScript,
		},
		// user bitcoin UTXOs to pay tx fee
		{
			TxHash:       "16de7df933dacd95b0d3af7325a5a2e680a1b7dd447a97e7678d8dfa1ac750b4",
			OutputIdx:    4,
			OutputAmount: 445000,
			PkScript:     userMultisigWallet.PKScript,
		},
	}

	// create tx
	msgTx, err := multisig.CreateBridgeTxSendRune(
		&bridgeMsg,
		inputs,
		userMultisigWallet.PKScript,
		relayerPkScript,
		TX_FEE,
	)
	fmt.Println("err: ", err)
	// sign tx
	totalSigs := [][][]byte{}
	// user key 1 sign tx
	userSigs1, _ := multisig.SignTapMultisig(userPrivKeys[0], msgTx, inputs, userMultisigWallet, 0)
	totalSigs = append(totalSigs, userSigs1)
	// user key 2 sign tx
	userSigs2, _ := multisig.SignTapMultisig(userPrivKeys[1], msgTx, inputs, userMultisigWallet, 0)
	totalSigs = append(totalSigs, userSigs2)
	// COMBINE SIGN
	signedMsgTx, _ := multisig.CombineTapMultisig(totalSigs, msgTx, inputs, userMultisigWallet, 0)

	var signedTx bytes.Buffer
	signedMsgTx.Serialize(&signedTx)
	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	signedMsgTxID := signedMsgTx.TxHash().String()

	fmt.Println("hexSignedTx: ", hexSignedTx)
	fmt.Println("signedMsgTxID: ", signedMsgTxID)
	fmt.Println("err sign: ", err)

	// Decipher runestone
	r := &runestone.Runestone{}
	artifact, err := r.Decipher(signedMsgTx)
	if err != nil {
		fmt.Println(err)
		return
	}
	a, _ := json.Marshal(artifact)
	fmt.Printf("Artifact: %s\n", string(a))

	// TODO: test the signedMsgTx
}