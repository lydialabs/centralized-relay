package multisig

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/bxelab/runestone"
	"lukechampine.com/uint128"
)

const(
	TX_FEE = 10000
	RELAYER_MULTISIG_ADDRESS = "tb1pv5j5j0dmq2c8d0vnehrlsgrwr9g95m849dl5v0tal8chfdgzqxfskv0w8u"
	RELAYER_MULTISIG_PK_SCRIPT = "51206525493dbb02b076bd93cdc7f8206e19505a6cf52b7f463d7df9f174b5020193"
	USER_MULTISIG_ADDRESS = "tb1pgzx880yfr7q8dgz8dqhw50sncu4f4hmw5cn3800354tuzcy9jx5shvv7su"
	USER_MULTISIG_PK_SCRIPT = "5120408c73bc891f8076a047682eea3e13c72a9adf6ea62713bdf1a557c1608591a9"
)

func TestGenerateKeys(t *testing.T) {
	chainParam := &chaincfg.TestNet3Params

	for i := 0; i < 3; i++ {
		privKey := GeneratePrivateKeyFromSeed([]byte{byte(i)}, chainParam)
		wif, _ := btcutil.NewWIF(privKey, chainParam, true)
		pubKey := wif.SerializePubKey();
		witnessProg := btcutil.Hash160(pubKey)
		p2wpkh, _ := btcutil.NewAddressWitnessPubKeyHash(witnessProg, chainParam)

		fmt.Printf("Account %v:\n Private Key: %v\n Public Key: %v\n Address: %v\n", i, wif.String(), hex.EncodeToString(pubKey), p2wpkh)
	}
}

func TestLoadWalletFromPrivateKey(t *testing.T) {
	chainParam := &chaincfg.TestNet3Params

	wif, _ := btcutil.DecodeWIF("cTYRscQxVhtsGjHeV59RHQJbzNnJHbf3FX4eyX5JkpDhqKdhtRvy")
	pubKey := wif.SerializePubKey();
	witnessProg := btcutil.Hash160(pubKey)
	p2wpkh, _ := btcutil.NewAddressWitnessPubKeyHash(witnessProg, chainParam)

	fmt.Printf("Account:\n Private Key: %v\n Public Key: %v\n Address: %v\n", string(wif.String()), hex.EncodeToString(pubKey), p2wpkh)
}

func TestRandomKeys(t *testing.T) {
	randomKeys(3, &chaincfg.TestNet3Params, []int{0, 1, 2})
}

func TestBuildMultisigTapScript(t *testing.T) {
	chainParam := &chaincfg.TestNet3Params

	relayerPrivKeys, relayersMultisigInfo := RandomMultisigInfo(3, 3, chainParam, []int{0, 1, 2}, 0, 1)
	relayersMultisigWallet, _ := BuildMultisigWallet(relayersMultisigInfo)
	_, userMultisigInfo := RandomMultisigInfo(2, 2, chainParam, []int{0, 3}, 1, 1)
	userMultisigWallet, _ := BuildMultisigWallet(userMultisigInfo)

	fmt.Println("relayersMultisigPKScript: ", hex.EncodeToString(relayersMultisigWallet.PKScript))
	relayersMultisigAddress, err := AddressOnChain(chainParam, relayersMultisigWallet)
	fmt.Println("relayersMultisigAddress, err : ", relayersMultisigAddress, err)
	fmt.Println("relayersPubKey Master : ", hex.EncodeToString(relayersMultisigInfo.PubKeys[0]))
	fmt.Println("relayersPubKey Slave 1 : ", hex.EncodeToString(relayersMultisigInfo.PubKeys[1]))
	fmt.Println("relayersPubKey Slave 2 : ", hex.EncodeToString(relayersMultisigInfo.PubKeys[2]))
	fmt.Println("relayersPrivKey Master : ", relayerPrivKeys[0])
	fmt.Println("relayersPrivKey Slave 1 : ", relayerPrivKeys[1])
	fmt.Println("relayersPrivKey Slave 2 : ", relayerPrivKeys[2])

	fmt.Println("userMultisigPKScript: ", hex.EncodeToString(userMultisigWallet.PKScript))
	userMultisigAddress, err := AddressOnChain(chainParam, userMultisigWallet)
	fmt.Println("userMultisigAddress, err : ", userMultisigAddress, err)
}

func TestSignTaproot(t *testing.T) {
	// TODO: REMOVE LATER --- test sign
	chainParam := &chaincfg.TestNet3Params

	_, userMultisigInfo := RandomMultisigInfo(2, 2, chainParam, []int{0, 3}, 1, 1)
	userMultisigInfo.PubKeys[0], _ = hex.DecodeString("d1fdca976c0ff461c501057d9c43bff16aa8145d5c0e0117984eadac11ce8973")
	userMultisigWallet, _ := BuildMultisigWallet(userMultisigInfo)

	fmt.Println("address: ", userMultisigWallet.SharedPublicKey)
	wif, _ := btcutil.DecodeWIF("cVb3ZUwwXoM5L3f8eFbg4NHU3VFqdy6nmcdEk7Nukc5huRAiHW7m")
	fmt.Println("PRIV KEY: ", wif.String())
	fmt.Println("pub key: ", hex.EncodeToString(wif.SerializePubKey()))
	signature, err := schnorr.Sign(wif.PrivKey, []byte("asdfghjklasdfghjklasdfghjklasdfg"))
	fmt.Println("err: ", err)
	fmt.Println("sign: ", hex.EncodeToString(signature.Serialize()))
}

func TestTransferBitcoin(t *testing.T) {
	chainParam := &chaincfg.TestNet3Params

	inputs := []*UTXO{
		{
			IsRelayersMultisig: false,
			TxHash:        "4933e04e3d9320df6e9f046ff83cfc3e9f884d8811df0539af7aaca0218189aa",
			OutputIdx:     0,
			OutputAmount:  4000000,
		},
	}

	outputs := []*OutputTx{
		{
			ReceiverAddress: "tb1pf0atpt2d3zel6udws38pkrh2e49vqd3c5jcud3a82srphnmpe55q0ecrzk",
			Amount:          1000000,
		},
	}

	userPrivKeys, userMultisigInfo := RandomMultisigInfo(2, 2, chainParam, []int{0, 3}, 1, 1)
	userMultisigWallet, _ := BuildMultisigWallet(userMultisigInfo)

	changeReceiverAddress := "tb1pgzx880yfr7q8dgz8dqhw50sncu4f4hmw5cn3800354tuzcy9jx5shvv7su"
	msgTx, _, txSigHashes, _ := CreateMultisigTx(inputs, outputs, 1000, &MultisigWallet{}, userMultisigWallet, chainParam, changeReceiverAddress, 0)

	tapSigParams := TapSigParams {
		TxSigHashes: txSigHashes,
		RelayersPKScript: []byte{},
		RelayersTapLeaf: txscript.TapLeaf{},
		UserPKScript: userMultisigWallet.PKScript,
		UserTapLeaf: userMultisigWallet.TapLeaves[0],
	}

	totalSigs := [][][]byte{}

	// USER SIGN TX
	userSigs1, _ := PartSignOnRawExternalTx(userPrivKeys[0], msgTx, inputs, tapSigParams, chainParam, true)
	totalSigs = append(totalSigs, userSigs1)
	userSigs2, _ := PartSignOnRawExternalTx(userPrivKeys[1], msgTx, inputs, tapSigParams, chainParam, true)
	totalSigs = append(totalSigs, userSigs2)

	// COMBINE SIGNS
	signedMsgTx, err := CombineMultisigSigs(msgTx, inputs, userMultisigWallet, 0, userMultisigWallet, 0, totalSigs)

	signedMsgTx.TxIn[0].Witness = signedMsgTx.TxIn[0].Witness[2:]

	var signedTx bytes.Buffer
	signedMsgTx.Serialize(&signedTx)
	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	signedMsgTxID := signedMsgTx.TxHash().String()

	fmt.Println("userPrivKeys 0 : ", userPrivKeys[0])
	fmt.Println("hexSignedTx: ", hexSignedTx)
	fmt.Println("signedMsgTxID: ", signedMsgTxID)
	fmt.Println("err sign: ", err)
}

func TestGenSharedInternalPubKey(t *testing.T) {
	b := make([]byte, 32)
	rand.Read(b)
	bHex := hex.EncodeToString(b)
	fmt.Printf("bHex: %v\n", bHex)
	sharedRandom := new(big.Int).SetBytes(b)
	genSharedInternalPubKey(sharedRandom)
}

func TestParseTx(t *testing.T) {
	hexSignedTx := "01000000000101bcbbb24bd5953d424debb9a24c8009298771eecd3ac0d3c4b219d906a319dfa80000000000e803000002e803000000000000225120d5254f2c52e2672daea941a86c99232693149fd0423ef523fe4e0dcb12a68d53401f000000000000225120d5254f2c52e2672daea941a86c99232693149fd0423ef523fe4e0dcb12a68d530540f4085e4f85eb81b8bd6afd77f728ea75716108cb29cd02aa031def6be65e97e98db40430554669d7b64476d76fd9ae6646529b7abfeee1ac4ad67de0bce9608040f3fc057a9ad0e4a0132040826e2c8e3ca0678ebd515146b8825f527f31195e1966d8424cdf9e963b7335178cab820534e1bd4ede4e8addf47c1bc449a764cec400962c7b22626173655661756c7441646472657373223a22222c22726563656976657241646472657373223a22227d7520fe44ec9f26b97ed30bd33898cf22de726e05389bde632d3aa6ad6746e15221d2ac2030edd881db1bc32b94f83ea5799c2e959854e0f99427d07c211206abd876d052ba201e83d56728fde393b41b74f2b859381661025f2ecec567cf392da7372de47833ba529c21c0636e6671d0135074f83177c5e456191043de9bd54744423b88d6b1ab4751650f00000000"
	msgTx, err := ParseTx(hexSignedTx)
	if err != nil {
		fmt.Printf("Err parse tx: %v", err)
		return
	}

	for _, txIn := range msgTx.TxIn {
		fmt.Printf("txIn: %+v\n ", txIn)
	}

	for _, txOut := range msgTx.TxOut {
		fmt.Printf("txOut: %+v\n ", txOut)
	}
}

func TestUserRecoveryTimeLock(t *testing.T) {
	chainParam := &chaincfg.TestNet3Params

	inputs := []*UTXO{
		{
			IsRelayersMultisig: false,
			TxHash:        "d316231a8aa1f74472ed9cc0f1ed0e36b9b290254cf6b2c377f0d92b299868bf",
			OutputIdx:     4,
			OutputAmount:  1929000,
		},
	}
	script1, _ := CreateRuneTransferScript(Rune{BlockNumber: 840000, TxIndex: 3}, big.NewInt(1000000000), 0)
	outputs := []*OutputTx{
		{
			ReceiverAddress: RELAYER_MULTISIG_ADDRESS,
			Amount:          1000000,
		},
		{
			OpReturnScript: script1,
		},
	}

	userPrivKeys, userMultisigInfo := RandomMultisigInfo(2, 2, chainParam, []int{0, 3}, 1, 1)
	userMultisigWallet, _ := BuildMultisigWallet(userMultisigInfo)

	changeReceiverAddress := USER_MULTISIG_ADDRESS
	msgTx, _, txSigHashes, _ := CreateMultisigTx(inputs, outputs, TX_FEE, &MultisigWallet{}, userMultisigWallet, chainParam, changeReceiverAddress, 1)
	tapSigParams := TapSigParams {
		TxSigHashes: txSigHashes,
		RelayersPKScript: []byte{},
		RelayersTapLeaf: txscript.TapLeaf{},
		UserPKScript: userMultisigWallet.PKScript,
		UserTapLeaf: userMultisigWallet.TapLeaves[1],
	}

	totalSigs := [][][]byte{}

	// USER SIGN TX
	userSigs, _ := PartSignOnRawExternalTx(userPrivKeys[1], msgTx, inputs, tapSigParams, chainParam, true)
	totalSigs = append(totalSigs, userSigs)

	// COMBINE SIGNS
	signedMsgTx, err := CombineMultisigSigs(msgTx, inputs, userMultisigWallet, 0, userMultisigWallet, 1, totalSigs)

	var signedTx bytes.Buffer
	signedMsgTx.Serialize(&signedTx)
	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	signedMsgTxID := signedMsgTx.TxHash().String()

	fmt.Println("hexSignedTx: ", hexSignedTx)
	fmt.Println("signedMsgTxID: ", signedMsgTxID)
	fmt.Println("err sign: ", err)
}

func TestTransferRune(t *testing.T) {
	chainParam := &chaincfg.TestNet3Params

	inputs := []*UTXO{
		{
			IsRelayersMultisig: false,
			TxHash:        "d316231a8aa1f74472ed9cc0f1ed0e36b9b290254cf6b2c377f0d92b299868bf",
			OutputIdx:     4,
			OutputAmount:  1929000,
		},
	}

	// Add Rune transfering
	// rune id 840000:3, amount 10000 (5 decimals), to output id 0
	runeId, _ := runestone.NewRuneId(840000, 3)
	defaultReceiver := uint32(0)
	runeStone := &runestone.Runestone{
		Edicts: []runestone.Edict{
			{
				ID: *runeId,
				Amount: uint128.From64(1000000000),
				Output: 0,
			},
		},
		Pointer: &defaultReceiver,
	}
	script, _ := runeStone.Encipher()
	outputs := []*OutputTx{
		{
			ReceiverAddress: RELAYER_MULTISIG_ADDRESS,
			Amount:          1000,
		},
		{
			OpReturnScript: script,
		},
	}

	userPrivKeys, userMultisigInfo := RandomMultisigInfo(2, 2, chainParam, []int{0, 3}, 1, 1)
	userMultisigWallet, _ := BuildMultisigWallet(userMultisigInfo)

	changeReceiverAddress := USER_MULTISIG_ADDRESS
	msgTx, _, txSigHashes, _ := CreateMultisigTx(inputs, outputs, TX_FEE, &MultisigWallet{}, userMultisigWallet, chainParam, changeReceiverAddress, 1)
	tapSigParams := TapSigParams {
		TxSigHashes: txSigHashes,
		RelayersPKScript: []byte{},
		RelayersTapLeaf: txscript.TapLeaf{},
		UserPKScript: userMultisigWallet.PKScript,
		UserTapLeaf: userMultisigWallet.TapLeaves[1],
	}

	totalSigs := [][][]byte{}

	// USER SIGN TX
	userSigs, _ := PartSignOnRawExternalTx(userPrivKeys[1], msgTx, inputs, tapSigParams, chainParam, true)
	totalSigs = append(totalSigs, userSigs)

	// COMBINE SIGNS
	signedMsgTx, err := CombineMultisigSigs(msgTx, inputs, userMultisigWallet, 0, userMultisigWallet, 1, totalSigs)

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
}

func TestRadFiInitPoolBitcoinRune(t *testing.T) {
	chainParam := &chaincfg.TestNet3Params
	_, relayersMultisigInfo := RandomMultisigInfo(3, 3, chainParam, []int{0, 1, 2}, 0, 1)
	relayersMultisigWallet, _ := BuildMultisigWallet(relayersMultisigInfo)
	userPrivKeys, userMultisigInfo := RandomMultisigInfo(2, 2, chainParam, []int{0, 3}, 1, 1)
	userMultisigWallet, _ := BuildMultisigWallet(userMultisigInfo)

	radfiMsg := RadFiProvideLiquidityMsg{
		Ticks: 	RadFiProvideLiquidityTicks{ UpperTick: 12345, LowerTick: -12345 },
		Fee:	30,
		Min0:	0,
		Min1:	10000,
		Amount0Desired: uint128.From64(10000),
		Amount1Desired: uint128.From64(1000000),
		InitPrice: uint128.From64(123456),
		Token0Id: runestone.RuneId{ Block: 0, Tx: 0},
		Token1Id: runestone.RuneId{ Block: 678, Tx: 90},
	}

	inputs := []*Input{
		// user bitcoin UTXO to add liquidity and pay tx fee
		{
			TxHash:			"d316231a8aa1f74472ed9cc0f1ed0e36b9b290254cf6b2c377f0d92b299868bf",
			OutputIdx:		4,
			OutputAmount:	1929000,
			PkScript:		userMultisigWallet.PKScript,
		},
		// user rune UTXO to add liquidity
		{
			TxHash:			"647a499a394bdb2a477f29b9f0515ed186e57a469a732be362a172cde4ea67a5",
			OutputIdx:		0,
			OutputAmount:	DUST_UTXO_AMOUNT,
			PkScript:		userMultisigWallet.PKScript,
		},
	}
	// create tx
	msgTx, err := CreateRadFiTxInitPool(
		&radfiMsg,
		inputs,
		relayersMultisigWallet.PKScript,
		userMultisigWallet.PKScript,
		TX_FEE,
	)
	fmt.Println("err: ", err)
	// sign tx
	totalSigs := [][][]byte{}
	// trading wallet admin sign tx
	adminSigs, _ := SignTapMultisig(userPrivKeys[0], msgTx, inputs, userMultisigWallet, 0)
	totalSigs = append(totalSigs, adminSigs)
	// user sign tx
	userSigs, _ := SignTapMultisig(userPrivKeys[1], msgTx, inputs, userMultisigWallet, 0)
	totalSigs = append(totalSigs, userSigs)

	// COMBINE SIGN
	signedMsgTx, _ := CombineTapMultisig(totalSigs, msgTx, inputs, userMultisigWallet, 0)

	var signedTx bytes.Buffer
	signedMsgTx.Serialize(&signedTx)
	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	signedMsgTxID := signedMsgTx.TxHash().String()

	fmt.Println("hexSignedTx: ", hexSignedTx)
	fmt.Println("signedMsgTxID: ", signedMsgTxID)

	// Decipher runestone
	r := &runestone.Runestone{}
	artifact, err := r.Decipher(signedMsgTx)
	if err != nil {
		fmt.Println(err)
		return
	}
	a, _ := json.Marshal(artifact)
	fmt.Printf("Artifact: %s\n", string(a))

	// Decode Radfi message
	decodedRadFiMessage, err := ReadRadFiMessage(signedMsgTx)
	fmt.Println("err decode: ", err)
	fmt.Println("decoded message - Flag     : ", decodedRadFiMessage.Flag)
	fmt.Println("decoded message - UpperTick: ", decodedRadFiMessage.ProvideLiquidityMsg.Ticks.UpperTick)
	fmt.Println("decoded message - LowerTick: ", decodedRadFiMessage.ProvideLiquidityMsg.Ticks.LowerTick)
	fmt.Println("decoded message - Fee      : ", decodedRadFiMessage.ProvideLiquidityMsg.Fee)
	fmt.Println("decoded message - Min0     : ", decodedRadFiMessage.ProvideLiquidityMsg.Min0)
	fmt.Println("decoded message - Min1     : ", decodedRadFiMessage.ProvideLiquidityMsg.Min1)
	fmt.Println("decoded message - Amount0  : ", decodedRadFiMessage.ProvideLiquidityMsg.Amount0Desired)
	fmt.Println("decoded message - Amount1  : ", decodedRadFiMessage.ProvideLiquidityMsg.Amount1Desired)
	fmt.Println("decoded message - InitPrice: ", decodedRadFiMessage.ProvideLiquidityMsg.InitPrice)
	fmt.Println("decoded message - Token0Id: ", decodedRadFiMessage.ProvideLiquidityMsg.Token0Id)
	fmt.Println("decoded message - Token1Id: ", decodedRadFiMessage.ProvideLiquidityMsg.Token1Id)
}

func TestRadFiInitPoolRuneRune(t *testing.T) {
	chainParam := &chaincfg.TestNet3Params
	_, relayersMultisigInfo := RandomMultisigInfo(3, 3, chainParam, []int{0, 1, 2}, 0, 1)
	relayersMultisigWallet, _ := BuildMultisigWallet(relayersMultisigInfo)
	userPrivKeys, userMultisigInfo := RandomMultisigInfo(2, 2, chainParam, []int{0, 3}, 1, 1)
	userMultisigWallet, _ := BuildMultisigWallet(userMultisigInfo)

	radfiMsg := RadFiProvideLiquidityMsg{
		Ticks: 	RadFiProvideLiquidityTicks{ UpperTick: 12345, LowerTick: -12345 },
		Fee:	30,
		Min0:	0,
		Min1:	10000,
		Amount0Desired: uint128.From64(1000000),
		Amount1Desired: uint128.From64(2000000),
		InitPrice: uint128.From64(123456),
		Token0Id: runestone.RuneId{ Block: 123, Tx: 321},
		Token1Id: runestone.RuneId{ Block: 678, Tx: 90},
	}

	inputs := []*Input{
		// user rune0 UTXO to add liquidity
		{
			TxHash:			"3aa9c4b9a71fe19560c467cdddce932eae8a10e28123a09acb27701123f2a8f7",
			OutputIdx:		6,
			OutputAmount:	DUST_UTXO_AMOUNT,
			PkScript:		userMultisigWallet.PKScript,
		},
		// user rune1 UTXO to add liquidity
		{
			TxHash:			"647a499a394bdb2a477f29b9f0515ed186e57a469a732be362a172cde4ea67a5",
			OutputIdx:		0,
			OutputAmount:	DUST_UTXO_AMOUNT,
			PkScript:		userMultisigWallet.PKScript,
		},
		// user bitcoin UTXO to pay tx fee
		{
			TxHash:			"d316231a8aa1f74472ed9cc0f1ed0e36b9b290254cf6b2c377f0d92b299868bf",
			OutputIdx:		4,
			OutputAmount:	1929000,
			PkScript:		userMultisigWallet.PKScript,
		},
	}
	// create tx
	msgTx, err := CreateRadFiTxInitPool(
		&radfiMsg,
		inputs,
		relayersMultisigWallet.PKScript,
		userMultisigWallet.PKScript,
		TX_FEE,
	)
	fmt.Println("err: ", err)
	// sign tx
	totalSigs := [][][]byte{}
	// trading wallet admin sign tx
	adminSigs, _ := SignTapMultisig(userPrivKeys[0], msgTx, inputs, userMultisigWallet, 0)
	totalSigs = append(totalSigs, adminSigs)
	// user sign tx
	userSigs, _ := SignTapMultisig(userPrivKeys[1], msgTx, inputs, userMultisigWallet, 0)
	totalSigs = append(totalSigs, userSigs)

	// COMBINE SIGN
	signedMsgTx, _ := CombineTapMultisig(totalSigs, msgTx, inputs, userMultisigWallet, 0)

	var signedTx bytes.Buffer
	signedMsgTx.Serialize(&signedTx)
	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	signedMsgTxID := signedMsgTx.TxHash().String()

	fmt.Println("hexSignedTx: ", hexSignedTx)
	fmt.Println("signedMsgTxID: ", signedMsgTxID)

	// Decipher runestone
	r := &runestone.Runestone{}
	artifact, err := r.Decipher(signedMsgTx)
	if err != nil {
		fmt.Println(err)
		return
	}
	a, _ := json.Marshal(artifact)
	fmt.Printf("Artifact: %s\n", string(a))

	// Decode Radfi message
	decodedRadFiMessage, err := ReadRadFiMessage(signedMsgTx)
	fmt.Println("err decode: ", err)
	fmt.Println("decoded message - Flag     : ", decodedRadFiMessage.Flag)
	fmt.Println("decoded message - UpperTick: ", decodedRadFiMessage.ProvideLiquidityMsg.Ticks.UpperTick)
	fmt.Println("decoded message - LowerTick: ", decodedRadFiMessage.ProvideLiquidityMsg.Ticks.LowerTick)
	fmt.Println("decoded message - Fee      : ", decodedRadFiMessage.ProvideLiquidityMsg.Fee)
	fmt.Println("decoded message - Min0     : ", decodedRadFiMessage.ProvideLiquidityMsg.Min0)
	fmt.Println("decoded message - Min1     : ", decodedRadFiMessage.ProvideLiquidityMsg.Min1)
	fmt.Println("decoded message - Amount0  : ", decodedRadFiMessage.ProvideLiquidityMsg.Amount0Desired)
	fmt.Println("decoded message - Amount1  : ", decodedRadFiMessage.ProvideLiquidityMsg.Amount1Desired)
	fmt.Println("decoded message - InitPrice: ", decodedRadFiMessage.ProvideLiquidityMsg.InitPrice)
	fmt.Println("decoded message - Token0Id: ", decodedRadFiMessage.ProvideLiquidityMsg.Token0Id)
	fmt.Println("decoded message - Token1Id: ", decodedRadFiMessage.ProvideLiquidityMsg.Token1Id)
}

func TestRadFiProvideLiquidityPoolBitcoinRune(t *testing.T) {
	chainParam := &chaincfg.TestNet3Params
	relayerPrivKeys, relayersMultisigInfo := RandomMultisigInfo(3, 3, chainParam, []int{0, 1, 2}, 0, 1)
	relayersMultisigWallet, _ := BuildMultisigWallet(relayersMultisigInfo)
	userPrivKeys, userMultisigInfo := RandomMultisigInfo(2, 2, chainParam, []int{0, 3}, 1, 1)
	userMultisigWallet, _ := BuildMultisigWallet(userMultisigInfo)

	radfiMsg := RadFiProvideLiquidityMsg{
		Ticks: 	RadFiProvideLiquidityTicks{ UpperTick: 12345, LowerTick: -12345 },
		Fee:	30,
		Min0:	0,
		Min1:	10000,
		Amount0Desired: uint128.From64(1000),
		Amount1Desired: uint128.From64(100000),
		InitPrice: uint128.From64(0),
		Token0Id: runestone.RuneId{ Block: 0, Tx: 0},
		Token1Id: runestone.RuneId{ Block: 678, Tx: 90},
	}

	inputs := []*Input{
		// pool current sequence number
		{
			TxHash:			"dfa3fc22b6436fdfaaf96bca443270cf1b6b50c24f2f2aff9ceaf668e2b1ed26",
			OutputIdx:		0,
			OutputAmount:	10000,
			PkScript:		relayersMultisigWallet.PKScript,
			Runes:			[]*runestone.Edict{
								{
									ID:		runestone.RuneId{ Block: 678, Tx: 90 },
									Amount:	uint128.From64(1000000),
								},
							},
		},
		// user bitcoin UTXO to add liquidity and pay tx fee
		{
			TxHash:			"d316231a8aa1f74472ed9cc0f1ed0e36b9b290254cf6b2c377f0d92b299868bf",
			OutputIdx:		4,
			OutputAmount:	1929000,
			PkScript:		userMultisigWallet.PKScript,
		},
		// user rune UTXO to add liquidity
		{
			TxHash:			"647a499a394bdb2a477f29b9f0515ed186e57a469a732be362a172cde4ea67a5",
			OutputIdx:		0,
			OutputAmount:	DUST_UTXO_AMOUNT,
			PkScript:		userMultisigWallet.PKScript,
		},
	}
	// create tx
	msgTx, err := CreateRadFiTxProvideLiquidity(
		&radfiMsg,
		inputs,
		relayersMultisigWallet.PKScript,
		userMultisigWallet.PKScript,
		TX_FEE,
	)
	fmt.Println("err: ", err)
	// sign tx
	totalSigs := [][][]byte{}
	totalSigsRelayer := [][][]byte{}
	// trading wallet admin sign tx
	adminSigs, _ := SignTapMultisig(userPrivKeys[0], msgTx, inputs, userMultisigWallet, 0)
	totalSigs = append(totalSigs, adminSigs)
	// user sign tx
	userSigs, _ := SignTapMultisig(userPrivKeys[1], msgTx, inputs, userMultisigWallet, 0)
	totalSigs = append(totalSigs, userSigs)
	// relayers sign tx
	masterRelayerSigs, _ := SignTapMultisig(relayerPrivKeys[0], msgTx, inputs, relayersMultisigWallet, 0)
	totalSigsRelayer = append(totalSigsRelayer, masterRelayerSigs)
	slaveRelayer1Sigs, _ := SignTapMultisig(relayerPrivKeys[1], msgTx, inputs, relayersMultisigWallet, 0)
	totalSigsRelayer = append(totalSigsRelayer, slaveRelayer1Sigs)
	slaveRelayer2Sigs, _ := SignTapMultisig(relayerPrivKeys[2], msgTx, inputs, relayersMultisigWallet, 0)
	totalSigsRelayer = append(totalSigsRelayer, slaveRelayer2Sigs)
	// COMBINE SIGN
	signedMsgTx, _ := CombineTapMultisig(totalSigs, msgTx, inputs, userMultisigWallet, 0)
	signedMsgTx, _ = CombineTapMultisig(totalSigsRelayer, signedMsgTx, inputs, relayersMultisigWallet, 0)

	var signedTx bytes.Buffer
	signedMsgTx.Serialize(&signedTx)
	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	signedMsgTxID := signedMsgTx.TxHash().String()

	fmt.Println("hexSignedTx: ", hexSignedTx)
	fmt.Println("signedMsgTxID: ", signedMsgTxID)

	// Decipher runestone
	r := &runestone.Runestone{}
	artifact, err := r.Decipher(signedMsgTx)
	if err != nil {
		fmt.Println(err)
		return
	}
	a, _ := json.Marshal(artifact)
	fmt.Printf("Artifact: %s\n", string(a))

	// Decode Radfi message
	decodedRadFiMessage, err := ReadRadFiMessage(signedMsgTx)
	fmt.Println("err decode: ", err)
	fmt.Println("decoded message - Flag     : ", decodedRadFiMessage.Flag)
	fmt.Println("decoded message - UpperTick: ", decodedRadFiMessage.ProvideLiquidityMsg.Ticks.UpperTick)
	fmt.Println("decoded message - LowerTick: ", decodedRadFiMessage.ProvideLiquidityMsg.Ticks.LowerTick)
	fmt.Println("decoded message - Fee      : ", decodedRadFiMessage.ProvideLiquidityMsg.Fee)
	fmt.Println("decoded message - Min0     : ", decodedRadFiMessage.ProvideLiquidityMsg.Min0)
	fmt.Println("decoded message - Min1     : ", decodedRadFiMessage.ProvideLiquidityMsg.Min1)
	fmt.Println("decoded message - Amount0  : ", decodedRadFiMessage.ProvideLiquidityMsg.Amount0Desired)
	fmt.Println("decoded message - Amount1  : ", decodedRadFiMessage.ProvideLiquidityMsg.Amount1Desired)
	fmt.Println("decoded message - InitPrice: ", decodedRadFiMessage.ProvideLiquidityMsg.InitPrice)
	fmt.Println("decoded message - Token0Id: ", decodedRadFiMessage.ProvideLiquidityMsg.Token0Id)
	fmt.Println("decoded message - Token1Id: ", decodedRadFiMessage.ProvideLiquidityMsg.Token1Id)
}

func TestRadFiProvideLiquidityPoolRuneRune(t *testing.T) {
	chainParam := &chaincfg.TestNet3Params
	relayerPrivKeys, relayersMultisigInfo := RandomMultisigInfo(3, 3, chainParam, []int{0, 1, 2}, 0, 1)
	relayersMultisigWallet, _ := BuildMultisigWallet(relayersMultisigInfo)
	userPrivKeys, userMultisigInfo := RandomMultisigInfo(2, 2, chainParam, []int{0, 3}, 1, 1)
	userMultisigWallet, _ := BuildMultisigWallet(userMultisigInfo)

	radfiMsg := RadFiProvideLiquidityMsg{
		Ticks: 	RadFiProvideLiquidityTicks{ UpperTick: 12345, LowerTick: -12345 },
		Fee:	30,
		Min0:	0,
		Min1:	10000,
		Amount0Desired: uint128.From64(100000),
		Amount1Desired: uint128.From64(200000),
		InitPrice: uint128.Zero,
		Token0Id: runestone.RuneId{ Block: 123, Tx: 321},
		Token1Id: runestone.RuneId{ Block: 678, Tx: 90},
	}

	inputs := []*Input{
		// pool current sequence number
		{
			TxHash:			"a9287ff640a9a06748100f6334f0e50a8c6a055aabff742443a5b1692d3dd0dc",
			OutputIdx:		0,
			OutputAmount:	DUST_UTXO_AMOUNT,
			PkScript:		relayersMultisigWallet.PKScript,
			Runes:			[]*runestone.Edict{
								{
									ID:		runestone.RuneId{ Block: 123, Tx: 321 },
									Amount:	uint128.From64(1000000),
								},
								{
									ID:		runestone.RuneId{ Block: 678, Tx: 90 },
									Amount:	uint128.From64(2000000),
								},
							},
		},
		// user rune0 UTXO to add liquidity
		{
			TxHash:			"3aa9c4b9a71fe19560c467cdddce932eae8a10e28123a09acb27701123f2a8f7",
			OutputIdx:		6,
			OutputAmount:	DUST_UTXO_AMOUNT,
			PkScript:		userMultisigWallet.PKScript,
		},
		// user rune1 UTXO to add liquidity
		{
			TxHash:			"647a499a394bdb2a477f29b9f0515ed186e57a469a732be362a172cde4ea67a5",
			OutputIdx:		0,
			OutputAmount:	DUST_UTXO_AMOUNT,
			PkScript:		userMultisigWallet.PKScript,
		},
		// user bitcoin UTXO to pay tx fee
		{
			TxHash:			"d316231a8aa1f74472ed9cc0f1ed0e36b9b290254cf6b2c377f0d92b299868bf",
			OutputIdx:		4,
			OutputAmount:	1929000,
			PkScript:		userMultisigWallet.PKScript,
		},
	}
	// create tx
	msgTx, err := CreateRadFiTxProvideLiquidity(
		&radfiMsg,
		inputs,
		relayersMultisigWallet.PKScript,
		userMultisigWallet.PKScript,
		TX_FEE,
	)
	fmt.Println("err: ", err)
	// sign tx
	totalSigs := [][][]byte{}
	totalSigsRelayer := [][][]byte{}
	// trading wallet admin sign tx
	adminSigs, _ := SignTapMultisig(userPrivKeys[0], msgTx, inputs, userMultisigWallet, 0)
	totalSigs = append(totalSigs, adminSigs)
	// user sign tx
	userSigs, _ := SignTapMultisig(userPrivKeys[1], msgTx, inputs, userMultisigWallet, 0)
	totalSigs = append(totalSigs, userSigs)
	// relayers sign tx
	masterRelayerSigs, _ := SignTapMultisig(relayerPrivKeys[0], msgTx, inputs, relayersMultisigWallet, 0)
	totalSigsRelayer = append(totalSigsRelayer, masterRelayerSigs)
	slaveRelayer1Sigs, _ := SignTapMultisig(relayerPrivKeys[1], msgTx, inputs, relayersMultisigWallet, 0)
	totalSigsRelayer = append(totalSigsRelayer, slaveRelayer1Sigs)
	slaveRelayer2Sigs, _ := SignTapMultisig(relayerPrivKeys[2], msgTx, inputs, relayersMultisigWallet, 0)
	totalSigsRelayer = append(totalSigsRelayer, slaveRelayer2Sigs)
	// COMBINE SIGN
	signedMsgTx, _ := CombineTapMultisig(totalSigs, msgTx, inputs, userMultisigWallet, 0)
	signedMsgTx, _ = CombineTapMultisig(totalSigsRelayer, signedMsgTx, inputs, relayersMultisigWallet, 0)

	var signedTx bytes.Buffer
	signedMsgTx.Serialize(&signedTx)
	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	signedMsgTxID := signedMsgTx.TxHash().String()

	fmt.Println("hexSignedTx: ", hexSignedTx)
	fmt.Println("signedMsgTxID: ", signedMsgTxID)

	// Decipher runestone
	r := &runestone.Runestone{}
	artifact, err := r.Decipher(signedMsgTx)
	if err != nil {
		fmt.Println(err)
		return
	}
	a, _ := json.Marshal(artifact)
	fmt.Printf("Artifact: %s\n", string(a))

	// Decode Radfi message
	decodedRadFiMessage, err := ReadRadFiMessage(signedMsgTx)
	fmt.Println("err decode: ", err)
	fmt.Println("decoded message - Flag     : ", decodedRadFiMessage.Flag)
	fmt.Println("decoded message - UpperTick: ", decodedRadFiMessage.ProvideLiquidityMsg.Ticks.UpperTick)
	fmt.Println("decoded message - LowerTick: ", decodedRadFiMessage.ProvideLiquidityMsg.Ticks.LowerTick)
	fmt.Println("decoded message - Fee      : ", decodedRadFiMessage.ProvideLiquidityMsg.Fee)
	fmt.Println("decoded message - Min0     : ", decodedRadFiMessage.ProvideLiquidityMsg.Min0)
	fmt.Println("decoded message - Min1     : ", decodedRadFiMessage.ProvideLiquidityMsg.Min1)
	fmt.Println("decoded message - Amount0  : ", decodedRadFiMessage.ProvideLiquidityMsg.Amount0Desired)
	fmt.Println("decoded message - Amount1  : ", decodedRadFiMessage.ProvideLiquidityMsg.Amount1Desired)
	fmt.Println("decoded message - InitPrice: ", decodedRadFiMessage.ProvideLiquidityMsg.InitPrice)
	fmt.Println("decoded message - Token0Id: ", decodedRadFiMessage.ProvideLiquidityMsg.Token0Id)
	fmt.Println("decoded message - Token1Id: ", decodedRadFiMessage.ProvideLiquidityMsg.Token1Id)
}

// func TestRadFiInitPool(t *testing.T) {
// 	chainParam := &chaincfg.TestNet3Params

// 	inputs := []*UTXO{
// 		{
// 			IsRelayersMultisig: false,
// 			TxHash:        "d316231a8aa1f74472ed9cc0f1ed0e36b9b290254cf6b2c377f0d92b299868bf",
// 			OutputIdx:     4,
// 			OutputAmount:  1929000,
// 		},
// 	}

// 	// OP_RETURN RadFi Init pool Message
// 	radfiMsg := RadFiProvideLiquidityMsg {
// 		Detail: &RadFiProvideLiquidityDetail{
// 			Fee:		30,
// 			UpperTick:	12345,
// 			LowerTick: 	-12345,
// 			Min0:		0,
// 			Min1:		10000,
// 		},
// 		InitPrice:  uint256.MustFromDecimal("123456789"),
// 	}
// 	script, _ := CreateProvideLiquidityScript(&radfiMsg)

// 	outputs := []*OutputTx{
// 		{
// 			ReceiverAddress: RELAYER_MULTISIG_ADDRESS,
// 			Amount:          1000,
// 		},
// 		{
// 			OpReturnScript: script,
// 		},
// 	}

// 	userPrivKeys, userMultisigInfo := RandomMultisigInfo(2, 2, chainParam, []int{0, 3}, 1, 1)
// 	userMultisigWallet, _ := BuildMultisigWallet(userMultisigInfo)

// 	changeReceiverAddress := USER_MULTISIG_ADDRESS
// 	msgTx, _, txSigHashes, _ := CreateMultisigTx(inputs, outputs, TX_FEE, &MultisigWallet{}, userMultisigWallet, chainParam, changeReceiverAddress, 1)

// 	tapSigParams := TapSigParams {
// 		TxSigHashes: txSigHashes,
// 		RelayersPKScript: []byte{},
// 		RelayersTapLeaf: txscript.TapLeaf{},
// 		UserPKScript: userMultisigWallet.PKScript,
// 		UserTapLeaf: userMultisigWallet.TapLeaves[1],
// 	}
// 	totalSigs := [][][]byte{}

// 	// USER SIGN TX
// 	userSigs, _ := PartSignOnRawExternalTx(userPrivKeys[1], msgTx, inputs, tapSigParams, chainParam, true)
// 	totalSigs = append(totalSigs, userSigs)

// 	// COMBINE SIGN
// 	signedMsgTx, err := CombineMultisigSigs(msgTx, inputs, userMultisigWallet, 0, userMultisigWallet, 1, totalSigs)

// 	var signedTx bytes.Buffer
// 	signedMsgTx.Serialize(&signedTx)
// 	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
// 	signedMsgTxID := signedMsgTx.TxHash().String()

// 	fmt.Println("hexSignedTx: ", hexSignedTx)
// 	fmt.Println("signedMsgTxID: ", signedMsgTxID)
// 	fmt.Println("err sign: ", err)

// 	// Decode Radfi message
// 	decodedRadFiMessage, err := ReadRadFiMessage(signedMsgTx)
// 	fmt.Println("err decode: ", err)
// 	fmt.Println("decoded message - Flag     : ", decodedRadFiMessage.Flag)
// 	fmt.Println("decoded message - Fee      : ", decodedRadFiMessage.ProvideLiquidityMsg.Detail.Fee)
// 	fmt.Println("decoded message - UpperTick: ", decodedRadFiMessage.ProvideLiquidityMsg.Detail.UpperTick)
// 	fmt.Println("decoded message - LowerTick: ", decodedRadFiMessage.ProvideLiquidityMsg.Detail.LowerTick)
// 	fmt.Println("decoded message - Min0     : ", decodedRadFiMessage.ProvideLiquidityMsg.Detail.Min0)
// 	fmt.Println("decoded message - Min1     : ", decodedRadFiMessage.ProvideLiquidityMsg.Detail.Min1)
// 	fmt.Println("decoded message - InitPrice: ", decodedRadFiMessage.ProvideLiquidityMsg.InitPrice)
// }

// func TestRadFiProvideLiquidity(t *testing.T) {
// 	chainParam := &chaincfg.TestNet3Params

// 	inputs := []*UTXO{
// 		{
// 			IsRelayersMultisig: true,
// 			TxHash:        "647a499a394bdb2a477f29b9f0515ed186e57a469a732be362a172cde4ea67a5",
// 			OutputIdx:     0,
// 			OutputAmount:  1000,
// 		},
// 		{
// 			IsRelayersMultisig: false,
// 			TxHash:        "647a499a394bdb2a477f29b9f0515ed186e57a469a732be362a172cde4ea67a5",
// 			OutputIdx:     2,
// 			OutputAmount:  1918000,
// 		},
// 		{
// 			IsRelayersMultisig: false,
// 			TxHash:        "96f43d22e98f9abf78422b27bd31412912b7c5e3a2b3f706d61148fc8d0f6550",
// 			OutputIdx:     0,
// 			OutputAmount:  1000,
// 		},
// 	}

// 	// OP_RETURN RadFi Provive Liquidity Message
// 	radfiMsg := RadFiProvideLiquidityMsg {
// 		Detail: &RadFiProvideLiquidityDetail{
// 			Fee:		30,
// 			UpperTick:	12345,
// 			LowerTick: 	-12345,
// 			Min0:		0,
// 			Min1:		10000,
// 		},
// 	}
// 	script1, _ := CreateProvideLiquidityScript(&radfiMsg)
// 	// OP_RETURN Rune transfering
// 	// rune id 840000:3, amount 10000 (5 decimals), to output id 4
// 	script2, _ := CreateRuneTransferScript(Rune{BlockNumber: 840000, TxIndex: 3}, big.NewInt(1000000000), 4)

// 	outputs := []*OutputTx{
// 		{
// 			ReceiverAddress: RELAYER_MULTISIG_ADDRESS,
// 			Amount:          1000,
// 		},
// 		{
// 			OpReturnScript: script1,
// 		},
// 		{
// 			ReceiverAddress: RELAYER_MULTISIG_ADDRESS,
// 			Amount:          20000,
// 		},
// 		{
// 			OpReturnScript: script2,
// 		},
// 		{
// 			ReceiverAddress: RELAYER_MULTISIG_ADDRESS,
// 			Amount:          1000,
// 		},
// 	}

// 	userPrivKeys, userMultisigInfo := RandomMultisigInfo(2, 2, chainParam, []int{0, 3}, 1, 1)
// 	userMultisigWallet, _ := BuildMultisigWallet(userMultisigInfo)

// 	changeReceiverAddress := USER_MULTISIG_ADDRESS
// 	msgTx, _, txSigHashes, _ := CreateMultisigTx(inputs, outputs, TX_FEE, &MultisigWallet{}, userMultisigWallet, chainParam, changeReceiverAddress, 1)

// 	tapSigParams := TapSigParams {
// 		TxSigHashes: txSigHashes,
// 		RelayersPKScript: []byte{},
// 		RelayersTapLeaf: txscript.TapLeaf{},
// 		UserPKScript: userMultisigWallet.PKScript,
// 		UserTapLeaf: userMultisigWallet.TapLeaves[1],
// 	}
// 	totalSigs := [][][]byte{}

// 	// USER SIGN TX
// 	userSigs, _ := PartSignOnRawExternalTx(userPrivKeys[1], msgTx, inputs, tapSigParams, chainParam, true)
// 	totalSigs = append(totalSigs, userSigs)

// 	// COMBINE SIGNS
// 	signedMsgTx, err := CombineMultisigSigs(msgTx, inputs, userMultisigWallet, 0, userMultisigWallet, 1, totalSigs)

// 	var signedTx bytes.Buffer
// 	signedMsgTx.Serialize(&signedTx)
// 	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
// 	signedMsgTxID := signedMsgTx.TxHash().String()

// 	fmt.Println("hexSignedTx: ", hexSignedTx)
// 	fmt.Println("signedMsgTxID: ", signedMsgTxID)
// 	fmt.Println("err sign: ", err)

// 	// Decode Radfi message
// 	decodedRadFiMessage, err := ReadRadFiMessage(signedMsgTx)
// 	fmt.Println("err decode: ", err)
// 	fmt.Println("decoded message - Flag     : ", decodedRadFiMessage.Flag)
// 	fmt.Println("decoded message - Fee      : ", decodedRadFiMessage.ProvideLiquidityMsg.Detail.Fee)
// 	fmt.Println("decoded message - UpperTick: ", decodedRadFiMessage.ProvideLiquidityMsg.Detail.UpperTick)
// 	fmt.Println("decoded message - LowerTick: ", decodedRadFiMessage.ProvideLiquidityMsg.Detail.LowerTick)
// 	fmt.Println("decoded message - Min0     : ", decodedRadFiMessage.ProvideLiquidityMsg.Detail.Min0)
// 	fmt.Println("decoded message - Min1     : ", decodedRadFiMessage.ProvideLiquidityMsg.Detail.Min1)
// 	fmt.Println("decoded message - InitPrice: ", decodedRadFiMessage.ProvideLiquidityMsg.InitPrice)

// 	// get list of Relayer bitcoin an rune UTXOs
// 	relayerScriptAddress := signedMsgTx.TxOut[0].PkScript
// 	sequenceNumberUTXO, bitcoinUTXOs, runeUTXOs, err := GetRelayerReceivedUTXO(signedMsgTx, 0, relayerScriptAddress)
// 	fmt.Println("-------------GetRelayerReceivedUTXO:")
// 	fmt.Println("relayerScriptAddress: ", relayerScriptAddress)
// 	fmt.Println("err: ", err)
// 	fmt.Println("sequenceNumberUTXO: ", sequenceNumberUTXO.IsRelayersMultisig, sequenceNumberUTXO.OutputIdx, sequenceNumberUTXO.OutputAmount, sequenceNumberUTXO.TxHash)
// 	fmt.Println("bitcoinUTXOs: ")
// 	for _, utxo := range bitcoinUTXOs {
// 		fmt.Println("bitcoinUTXO: ", utxo.IsRelayersMultisig, utxo.OutputIdx, utxo.OutputAmount, utxo.TxHash)
// 	}
// 	fmt.Println("runeUTXOs: ")
// 	for _, utxo := range runeUTXOs {
// 		fmt.Println("runeUTXO: ", utxo.edict, utxo.edictUTXO.IsRelayersMultisig, utxo.edictUTXO.OutputIdx, utxo.edictUTXO.OutputAmount, utxo.edictUTXO.TxHash)
// 	}
// }

// func TestRadFiSwap(t *testing.T) {
// 	chainParam := &chaincfg.TestNet3Params

// 	inputs := []*UTXO{
// 		{
// 			IsRelayersMultisig: true,
// 			TxHash:        "4e2b225d732b4123443e8d737e292921ea4a891ddf3a5c9842a4409f2c82805d",
// 			OutputIdx:     0,
// 			OutputAmount:  1000,
// 		},
// 		{
// 			IsRelayersMultisig: false,
// 			TxHash:        "4e2b225d732b4123443e8d737e292921ea4a891ddf3a5c9842a4409f2c82805d",
// 			OutputIdx:     5,
// 			OutputAmount:  1888000,
// 		},
// 		{
// 			IsRelayersMultisig: false,
// 			TxHash:        "4e2b225d732b4123443e8d737e292921ea4a891ddf3a5c9842a4409f2c82805d",
// 			OutputIdx:     4,
// 			OutputAmount:  1000,
// 		},
// 	}

// 	// OP_RETURN RadFi Swap Message
// 	radfiMsg := RadFiSwapMsg {
// 		IsExactInOut: true,
// 		TokenOutIndex: 2,
// 	}
// 	script1, _ := CreateSwapScript(&radfiMsg)
// 	// OP_RETURN Rune transfering
// 	// rune id 840000:3, amount 4000 (5 decimals), to output id 3
// 	script2, _ := CreateRuneTransferScript(Rune{BlockNumber: 840000, TxIndex: 3}, big.NewInt(400000000), 3)
// 	// rune id 840000:3, amount 6000 (5 decimals), to output id 6
// 	script3, _ := CreateRuneTransferScript(Rune{BlockNumber: 840000, TxIndex: 3}, big.NewInt(600000000), 6)

// 	outputs := []*OutputTx{
// 		{
// 			ReceiverAddress: RELAYER_MULTISIG_ADDRESS,
// 			Amount:          1000,
// 		},
// 		{
// 			OpReturnScript: script1,
// 		},
// 		{
// 			OpReturnScript: script2,
// 		},
// 		{
// 			ReceiverAddress: USER_MULTISIG_ADDRESS,
// 			Amount:          1000,
// 		},
// 		{
// 			ReceiverAddress: RELAYER_MULTISIG_ADDRESS,
// 			Amount:          10000,
// 		},
// 		{
// 			OpReturnScript: script3,
// 		},
// 		{
// 			ReceiverAddress: RELAYER_MULTISIG_ADDRESS,
// 			Amount:          1000,
// 		},
// 	}

// 	userPrivKeys, userMultisigInfo := RandomMultisigInfo(2, 2, chainParam, []int{0, 3}, 1, 1)
// 	userMultisigWallet, _ := BuildMultisigWallet(userMultisigInfo)

// 	changeReceiverAddress := USER_MULTISIG_ADDRESS
// 	msgTx, _, txSigHashes, _ := CreateMultisigTx(inputs, outputs, TX_FEE, &MultisigWallet{}, userMultisigWallet, chainParam, changeReceiverAddress, 1)

// 	tapSigParams := TapSigParams {
// 		TxSigHashes: txSigHashes,
// 		RelayersPKScript: []byte{},
// 		RelayersTapLeaf: txscript.TapLeaf{},
// 		UserPKScript: userMultisigWallet.PKScript,
// 		UserTapLeaf: userMultisigWallet.TapLeaves[1],
// 	}
// 	totalSigs := [][][]byte{}

// 	// USER SIGN TX
// 	userSigs, _ := PartSignOnRawExternalTx(userPrivKeys[1], msgTx, inputs, tapSigParams, chainParam, true)
// 	totalSigs = append(totalSigs, userSigs)

// 	// COMBINE SIGNS
// 	signedMsgTx, err := CombineMultisigSigs(msgTx, inputs, userMultisigWallet, 0, userMultisigWallet, 1, totalSigs)

// 	var signedTx bytes.Buffer
// 	signedMsgTx.Serialize(&signedTx)
// 	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
// 	signedMsgTxID := signedMsgTx.TxHash().String()

// 	fmt.Println("hexSignedTx: ", hexSignedTx)
// 	fmt.Println("signedMsgTxID: ", signedMsgTxID)
// 	fmt.Println("err sign: ", err)

// 	// Decode Radfi message
// 	decodedRadFiMessage, err := ReadRadFiMessage(signedMsgTx)

// 	fmt.Println("err decode: ", err)
// 	fmt.Println("decoded message - Flag                  : ", decodedRadFiMessage.Flag)
// 	fmt.Println("decoded message - IsExactInOut          : ", decodedRadFiMessage.SwapMsg.IsExactInOut)
// 	fmt.Println("decoded message - TokenOutIndex         : ", decodedRadFiMessage.SwapMsg.TokenOutIndex)
// }

// func TestRadFiWithdrawLiquidity(t *testing.T) {
// 	chainParam := &chaincfg.TestNet3Params

// 	inputs := []*UTXO{
// 		{
// 			IsRelayersMultisig: true,
// 			TxHash:        "3aa9c4b9a71fe19560c467cdddce932eae8a10e28123a09acb27701123f2a8f7",
// 			OutputIdx:     0,
// 			OutputAmount:  1000,
// 		},
// 		{
// 			IsRelayersMultisig: true,
// 			TxHash:        "3aa9c4b9a71fe19560c467cdddce932eae8a10e28123a09acb27701123f2a8f7",
// 			OutputIdx:     4,
// 			OutputAmount:  10000,
// 		},
// 		{
// 			IsRelayersMultisig: true,
// 			TxHash:        "3aa9c4b9a71fe19560c467cdddce932eae8a10e28123a09acb27701123f2a8f7",
// 			OutputIdx:     6,
// 			OutputAmount:  1000,
// 		},
// 		{
// 			IsRelayersMultisig: false,
// 			TxHash:        "4e2b225d732b4123443e8d737e292921ea4a891ddf3a5c9842a4409f2c82805d",
// 			OutputIdx:     5,
// 			OutputAmount:  1867000,
// 		},
// 	}

// 	// OP_RETURN RadFi Withdraw Liquidity Message
// 	radfiMsg := RadFiWithdrawLiquidityMsg {
// 		RecipientIndex:	2,
// 		LiquidityValue: uint256.MustFromDecimal("123456"),
// 		NftId:			uint256.MustFromDecimal("123456789"),
// 	}
// 	script1, _ := CreateWithdrawLiquidityScript(&radfiMsg)
// 	// OP_RETURN Rune transfering
// 	// rune id 840000:3, amount 2500 (5 decimals), to output id 3
// 	script2, _ := CreateRuneTransferScript(Rune{BlockNumber: 840000, TxIndex: 3}, big.NewInt(250000000), 4)
// 	// rune id 840000:3, amount 3500 (5 decimals), to output id 7
// 	script3, _ := CreateRuneTransferScript(Rune{BlockNumber: 840000, TxIndex: 3}, big.NewInt(350000000), 7)

// 	outputs := []*OutputTx{
// 		{
// 			ReceiverAddress: RELAYER_MULTISIG_ADDRESS,
// 			Amount:          1000,
// 		},
// 		{
// 			OpReturnScript: script1,
// 		},
// 		{
// 			ReceiverAddress: USER_MULTISIG_ADDRESS,
// 			Amount:          5000,
// 		},
// 		{
// 			OpReturnScript: script2,
// 		},
// 		{
// 			ReceiverAddress: USER_MULTISIG_ADDRESS,
// 			Amount:          1000,
// 		},
// 		{
// 			ReceiverAddress: RELAYER_MULTISIG_ADDRESS,
// 			Amount:          5000,
// 		},
// 		{
// 			OpReturnScript: script3,
// 		},
// 		{
// 			ReceiverAddress: RELAYER_MULTISIG_ADDRESS,
// 			Amount:          1000,
// 		},
// 	}
// 	userPrivKeys, userMultisigInfo := RandomMultisigInfo(2, 2, chainParam, []int{0, 3}, 1, 1)
// 	userMultisigWallet, _ := BuildMultisigWallet(userMultisigInfo)

// 	changeReceiverAddress := USER_MULTISIG_ADDRESS
// 	msgTx, _, txSigHashes, _ := CreateMultisigTx(inputs, outputs, TX_FEE, &MultisigWallet{}, userMultisigWallet, chainParam, changeReceiverAddress, 1)

// 	tapSigParams := TapSigParams {
// 		TxSigHashes: txSigHashes,
// 		RelayersPKScript: []byte{},
// 		RelayersTapLeaf: txscript.TapLeaf{},
// 		UserPKScript: userMultisigWallet.PKScript,
// 		UserTapLeaf: userMultisigWallet.TapLeaves[1],
// 	}
// 	totalSigs := [][][]byte{}

// 	// USER SIGN TX
// 	userSigs, _ := PartSignOnRawExternalTx(userPrivKeys[1], msgTx, inputs, tapSigParams, chainParam, true)
// 	totalSigs = append(totalSigs, userSigs)

// 	// COMBINE SIGNS
// 	signedMsgTx, err := CombineMultisigSigs(msgTx, inputs, userMultisigWallet, 0, userMultisigWallet, 1, totalSigs)

// 	var signedTx bytes.Buffer
// 	signedMsgTx.Serialize(&signedTx)
// 	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
// 	signedMsgTxID := signedMsgTx.TxHash().String()

// 	fmt.Println("hexSignedTx: ", hexSignedTx)
// 	fmt.Println("signedMsgTxID: ", signedMsgTxID)
// 	fmt.Println("err sign: ", err)

// 	// Decode Radfi message
// 	decodedRadFiMessage, err := ReadRadFiMessage(signedMsgTx)

// 	fmt.Println("err decode: ", err)
// 	fmt.Println("decoded message - Flag           : ", decodedRadFiMessage.Flag)
// 	fmt.Println("decoded message - RecipientIndex : ", decodedRadFiMessage.WithdrawLiquidityMsg.RecipientIndex)
// 	fmt.Println("decoded message - LiquidityValue : ", decodedRadFiMessage.WithdrawLiquidityMsg.LiquidityValue)
// 	fmt.Println("decoded message - NftId          : ", decodedRadFiMessage.WithdrawLiquidityMsg.NftId)
// }

func TestRadFiCollectFees(t *testing.T) {
	chainParam := &chaincfg.TestNet3Params

	inputs := []*UTXO{
		{
			IsRelayersMultisig: true,
			TxHash:        "b508b193dd9712406c8d712af00e6d3e34b8d2ad52809b72c63f6463c1555703",
			OutputIdx:     0,
			OutputAmount:  1000,
		},
		{
			IsRelayersMultisig: true,
			TxHash:        "b508b193dd9712406c8d712af00e6d3e34b8d2ad52809b72c63f6463c1555703",
			OutputIdx:     5,
			OutputAmount:  5000,
		},
		{
			IsRelayersMultisig: true,
			TxHash:        "b508b193dd9712406c8d712af00e6d3e34b8d2ad52809b72c63f6463c1555703",
			OutputIdx:     7,
			OutputAmount:  1000,
		},
		{
			IsRelayersMultisig: false,
			TxHash:        "b508b193dd9712406c8d712af00e6d3e34b8d2ad52809b72c63f6463c1555703",
			OutputIdx:     8,
			OutputAmount:  1856000,
		},
	}

	// OP_RETURN RadFi Collect Fees
	radfiMsg := RadFiCollectFeesMsg {
		RecipientIndex:	2,
		NftId:			uint128.From64(123456789),
	}
	script1, _ := CreateCollectFeesScript(&radfiMsg)
	// OP_RETURN Rune transfering
	// rune id 840000:3, amount 500 (5 decimals), to output id 4
	script2, _ := CreateRuneTransferScript(Rune{BlockNumber: 840000, TxIndex: 3}, big.NewInt(50000000), 4)
	// rune id 840000:3, amount 3000 (5 decimals), to output id 7
	script3, _ := CreateRuneTransferScript(Rune{BlockNumber: 840000, TxIndex: 3}, big.NewInt(300000000), 7)

	outputs := []*OutputTx{
		{
			ReceiverAddress: RELAYER_MULTISIG_ADDRESS,
			Amount:          1000,
		},
		{
			OpReturnScript: script1,
		},
		{
			ReceiverAddress: USER_MULTISIG_ADDRESS,
			Amount:          2000,
		},
		{
			OpReturnScript: script2,
		},
		{
			ReceiverAddress: USER_MULTISIG_ADDRESS,
			Amount:          1000,
		},
		{
			ReceiverAddress: RELAYER_MULTISIG_ADDRESS,
			Amount:          3000,
		},
		{
			OpReturnScript: script3,
		},
		{
			ReceiverAddress: RELAYER_MULTISIG_ADDRESS,
			Amount:          1000,
		},
	}

	userPrivKeys, userMultisigInfo := RandomMultisigInfo(2, 2, chainParam, []int{0, 3}, 1, 1)
	userMultisigWallet, _ := BuildMultisigWallet(userMultisigInfo)

	changeReceiverAddress := USER_MULTISIG_ADDRESS
	msgTx, _, txSigHashes, _ := CreateMultisigTx(inputs, outputs, TX_FEE, &MultisigWallet{}, userMultisigWallet, chainParam, changeReceiverAddress, 1)

	tapSigParams := TapSigParams {
		TxSigHashes: txSigHashes,
		RelayersPKScript: []byte{},
		RelayersTapLeaf: txscript.TapLeaf{},
		UserPKScript: userMultisigWallet.PKScript,
		UserTapLeaf: userMultisigWallet.TapLeaves[1],
	}
	totalSigs := [][][]byte{}

	// USER SIGN TX
	userSigs, _ := PartSignOnRawExternalTx(userPrivKeys[1], msgTx, inputs, tapSigParams, chainParam, true)
	totalSigs = append(totalSigs, userSigs)

	// COMBINE SIGNS
	signedMsgTx, err := CombineMultisigSigs(msgTx, inputs, userMultisigWallet, 0, userMultisigWallet, 1, totalSigs)

	var signedTx bytes.Buffer
	signedMsgTx.Serialize(&signedTx)
	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	signedMsgTxID := signedMsgTx.TxHash().String()

	fmt.Println("hexSignedTx: ", hexSignedTx)
	fmt.Println("signedMsgTxID: ", signedMsgTxID)
	fmt.Println("err sign: ", err)

	// Decode Radfi message
	decodedRadFiMessage, err := ReadRadFiMessage(signedMsgTx)

	fmt.Println("err decode: ", err)
	fmt.Println("decoded message - Flag           : ", decodedRadFiMessage.Flag)
	fmt.Println("decoded message - RecipientIndex : ", decodedRadFiMessage.CollectFeesMsg.RecipientIndex)
	fmt.Println("decoded message - NftId          : ", decodedRadFiMessage.CollectFeesMsg.NftId)

	// Decipher runestone
	r := &runestone.Runestone{}
	artifact, err := r.Decipher(signedMsgTx)
	if err != nil {
		fmt.Println(err)
		return
	}
	a, _ := json.Marshal(artifact)
	fmt.Printf("Artifact: %s\n", string(a))
}

func TestRadFiIncreaseLiquidity(t *testing.T) {
	chainParam := &chaincfg.TestNet3Params

	inputs := []*UTXO{
		{
			IsRelayersMultisig: true,
			TxHash:        "546fd5aa7336a82ec31dda54de0ec0e81c9285e42f75768f5de058b390126e51",
			OutputIdx:     0,
			OutputAmount:  1000,
		},
		{
			IsRelayersMultisig: false,
			TxHash:        "546fd5aa7336a82ec31dda54de0ec0e81c9285e42f75768f5de058b390126e51",
			OutputIdx:     8,
			OutputAmount:  1845000,
		},
		{
			IsRelayersMultisig: false,
			TxHash:        "546fd5aa7336a82ec31dda54de0ec0e81c9285e42f75768f5de058b390126e51",
			OutputIdx:     0,
			OutputAmount:  1000,
		},
	}

	// OP_RETURN RadFi Increase Liquidity Message
	radfiMsg := RadFiIncreaseLiquidityMsg {
		Min0:	0,
		Min1:	10000,
		NftId:	uint128.From64(123456789),
	}
	script1, _ := CreateIncreaseLiquidityScript(&radfiMsg)
	// OP_RETURN Rune transfering
	// rune id 840000:3, amount 500 (5 decimals), to output id 4
	script2, _ := CreateRuneTransferScript(Rune{BlockNumber: 840000, TxIndex: 3}, big.NewInt(50000000), 4)

	outputs := []*OutputTx{
		{
			ReceiverAddress: RELAYER_MULTISIG_ADDRESS,
			Amount:          1000,
		},
		{
			OpReturnScript: script1,
		},
		{
			ReceiverAddress: RELAYER_MULTISIG_ADDRESS,
			Amount:          20000,
		},
		{
			OpReturnScript: script2,
		},
		{
			ReceiverAddress: RELAYER_MULTISIG_ADDRESS,
			Amount:          1000,
		},
	}

	userPrivKeys, userMultisigInfo := RandomMultisigInfo(2, 2, chainParam, []int{0, 3}, 1, 1)
	userMultisigWallet, _ := BuildMultisigWallet(userMultisigInfo)

	changeReceiverAddress := USER_MULTISIG_ADDRESS
	msgTx, _, txSigHashes, _ := CreateMultisigTx(inputs, outputs, TX_FEE, &MultisigWallet{}, userMultisigWallet, chainParam, changeReceiverAddress, 1)

	tapSigParams := TapSigParams {
		TxSigHashes: txSigHashes,
		RelayersPKScript: []byte{},
		RelayersTapLeaf: txscript.TapLeaf{},
		UserPKScript: userMultisigWallet.PKScript,
		UserTapLeaf: userMultisigWallet.TapLeaves[1],
	}
	totalSigs := [][][]byte{}

	// USER SIGN TX
	userSigs, _ := PartSignOnRawExternalTx(userPrivKeys[1], msgTx, inputs, tapSigParams, chainParam, true)
	totalSigs = append(totalSigs, userSigs)

	// COMBINE SIGNS
	signedMsgTx, err := CombineMultisigSigs(msgTx, inputs, userMultisigWallet, 0, userMultisigWallet, 1, totalSigs)

	var signedTx bytes.Buffer
	signedMsgTx.Serialize(&signedTx)
	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	signedMsgTxID := signedMsgTx.TxHash().String()

	fmt.Println("hexSignedTx: ", hexSignedTx)
	fmt.Println("signedMsgTxID: ", signedMsgTxID)
	fmt.Println("err sign: ", err)

	// Decode Radfi message
	decodedRadFiMessage, err := ReadRadFiMessage(signedMsgTx)

	fmt.Println("err decode: ", err)
	fmt.Println("decoded message - Flag           : ", decodedRadFiMessage.Flag)
	fmt.Println("decoded message - Min0           : ", decodedRadFiMessage.IncreaseLiquidityMsg.Min0)
	fmt.Println("decoded message - Min1           : ", decodedRadFiMessage.IncreaseLiquidityMsg.Min1)
	fmt.Println("decoded message - NftId          : ", decodedRadFiMessage.IncreaseLiquidityMsg.NftId)
}
