package multisig

import (
	"bytes"
	"context"

	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"sort"
	"time"

	"github.com/btcsuite/btcd/wire"
	"github.com/studyzy/runestone"
	"lukechampine.com/uint128"
)

const (
	UNISAT_DEFAULT_MAINNET	= "https://open-api.unisat.io"
	UNISAT_DEFAULT_TESTNET	= "https://open-api-testnet.unisat.io"
)

type DataBlockchainInfo struct {
	Chain         string `json:"chain"`
	Blocks        int64  `json:"blocks"`
	Headers       int64  `json:"headers"`
	BestBlockHash string `json:"bestBlockHash"`
	PrevBlockHash string `json:"prevBlockHash"`
	Difficulty    string `json:"difficulty"`
	MedianTime    int64  `json:"medianTime"`
	ChainWork     string `json:"chainwork"`
}

type ResponseBlockchainInfo struct {
	Code    int64  `json:"code"`
	Message string `json:"msg"`

	Data DataBlockchainInfo `json:"data"`
}

type Tx struct {
	TxId             string   `json:"txid"`
	Ins              int      `json:"nIn"`
	Outs             int      `json:"nOut"`
	Size             int      `json:"size"`
	WitOffset        int      `json:"witOffset"`
	Locktime         int      `json:"locktime"`
	InSatoshi        *big.Int `json:"inSatoshi"`
	OutSatoshi       *big.Int `json:"outSatoshi"`
	NewInscriptions  int      `json:"nNewInscription"`
	InInscriptions   int      `json:"nInInscription"`
	OutInscriptions  int      `json:"nOutInscription"`
	LostInscriptions int      `json:"nLostInscription"`
	Timestamp        int64    `json:"timestamp"`
	Height           int64    `json:"height"`
	BlockId          string   `json:"blkid"`
	Index            int      `json:"idx"`
	Confirmations    int      `json:"confirmations"`
}

type ResponseBlockTransactions struct {
	Code    int64  `json:"code"`
	Message string `json:"msg"`

	Data []Tx `json:"data"`
}

type ResponseTxInfo struct {
	Code    int64  `json:"code"`
	Message string `json:"msg"`

	Data Tx `json:"data"`
}

type Inscription struct {
	InscriptionNumber int64  `json:"inscriptionNumber"`
	InscriptionId     string `json:"inscriptionId"`
	Offset            int    `json:"offset"`
	Moved             bool   `json:"moved"`
	IsBRC20           bool   `json:"isBRC20"`
}

type InputData struct {
	Height       int64         `json:"height"`
	TxId         string        `json:"txid"`
	Index        int           `json:"idx"`
	ScriptSig    string        `json:"scriptSig"`
	ScriptWits   string        `json:"scriptWits"`
	Sequence     int           `json:"sequence"`
	HeightTxo    int64         `json:"heightTxo"`
	Utxid        string        `json:"utxid"`
	Vout         int           `json:"vout"`
	Address      string        `json:"address"`
	CodeType     int           `json:"codeType"`
	Satoshi      *big.Int      `json:"satoshi"`
	ScriptType   string        `json:"scriptType"`
	ScriptPk     string        `json:"scriptPk"`
	Inscriptions []Inscription `json:"inscriptions"`
}

type ResponseTxInputs struct {
	Code    int64  `json:"code"`
	Message string `json:"msg"`

	Data []InputData `json:"data"`
}

type Output struct {
	TxId         string        `json:"txid"`
	Vout         int           `json:"vout"`
	Address      string        `json:"address"`
	CodeType     int           `json:"codeType"`
	Satoshi      *big.Int      `json:"satoshi"`
	ScriptType   string        `json:"scriptType"`
	ScriptPk     string        `json:"scriptPk"`
	Height       int64         `json:"height"`
	Index        int           `json:"idx"`
	Inscriptions []Inscription `json:"inscriptions"`
	TxSpent      string        `json:"txidSpent"`
	HeightSpent  int64         `json:"heightSpent"`
}

type ResponseTxOutputs struct {
	Code    int64  `json:"code"`
	Message string `json:"msg"`

	Data []Output `json:"data"`
}

type Balance struct {
	Address string `json:"address"`

	Satoshi        *big.Int `json:"satoshi"`
	PendingSatoshi *big.Int `json:"pendingSatoshi"`
	UtxoCount      int64    `json:"utxoCount"`

	BtcSatoshi        *big.Int `json:"btcSatoshi"`
	BtcPendingSatoshi *big.Int `json:"btcPendingSatoshi"`
	BtcUtxoCount      int64    `json:"btcUtxoCount"`

	InscriptionSatoshi        *big.Int `json:"inscriptionSatoshi"`
	InscriptionPendingSatoshi *big.Int `json:"inscriptionPendingSatoshi"`
	InscriptionUtxoCount      int64    `json:"inscriptionUtxoCount"`
}

type ResponseAddressBalance struct {
	Code    int64  `json:"code"`
	Message string `json:"msg"`

	Data Balance `json:"data"`
}

type UTXO struct {
	TxId         string        `json:"txid"`
	Vout         int           `json:"vout"`
	Satoshi      *big.Int      `json:"satoshi"`
	ScriptType   string        `json:"scriptType"`
	ScriptPk     string        `json:"scriptPk"`
	CodeType     int           `json:"codeType"`
	Address      string        `json:"address"`
	Height       int64         `json:"height"`
	Index        int           `json:"idx"`
	IsOpInRBF    bool          `json:"isOpInRBF"`
	IsSpent      bool          `json:"isSpent"`
	Inscriptions []Inscription `json:"inscriptions"`
}

type DataUtxoList struct {
	Cursor                int    `json:"cursor"`
	Total                 int    `json:"total"`
	TotalConfirmed        int    `json:"totalConfirmed"`
	TotalUnconfirmed      int    `json:"totalUnconfirmed"`
	TotalUnconfirmedSpent int    `json:"totalUnconfirmedSpent"`
	Utxo                  []UTXO `json:"utxo"`
}

type ResponseBtcUtxo struct {
	Code    int64  `json:"code"`
	Message string `json:"msg"`

	Data DataUtxoList `json:"data"`
}

type RuneDetail struct {
	Amount       string        `json:"amount"`
	RuneId       string        `json:"runeid"`
	Rune         string        `json:"rune"`
	SpacedRune   string        `json:"spacedRune"`
	Symbol       string        `json:"symbol"`
	Divisibility int           `json:"divisibility"`
}

type RuneUTXO struct {
	Height        int64         `json:"height"`
	Confirmations int64         `json:"confirmations"`
	Address       string        `json:"address"`
	Satoshi       *big.Int      `json:"satoshi"`
	ScriptPk      string        `json:"scriptPk"`
	TxId          string        `json:"txid"`
	Vout          int           `json:"vout"`
	Runes 		  []RuneDetail  `json:"runes"`
}

type DataRuneUtxoList struct {
	Height int64     `json:"height"`
	Start int        `json:"start"`
	Total int        `json:"total"`
	Utxo  []RuneUTXO `json:"utxo"`
}

type ResponseRuneUtxo struct {
	Code    int64  `json:"code"`
	Message string `json:"msg"`

	Data DataRuneUtxoList `json:"data"`
}

type ResponseRuneBalanceByUtxo struct {
	Code    int64  `json:"code"`
	Message string `json:"msg"`

	Data []RuneDetail `json:"data"`
}

type ResponseRuneBalanceOfAddress struct {
	Code    int64  `json:"code"`
	Message string `json:"msg"`

	Data RuneDetail `json:"data"`
}

type ResponseBitcoinBalanceByUtxo struct {
	Code    int64  `json:"code"`
	Message string `json:"msg"`

	Data    UTXO `json:"data"`
}

type MempoolRecommendedFeeResponse struct {
	FastestFee  uint64 `json:"fastestFee"`
	HalfHourFee uint64 `json:"halfHourFee"`
	HourFee     uint64 `json:"hourFee"`
	EconomyFee  uint64 `json:"economyFee"`
	MinimumFee  uint64 `json:"minimumFee"`
}

type MempoolGetTransactionStatus struct {
	Confirmed    bool  `json:"confirmed"`
	BlockHeight  int64 `json:"block_height"`
	// block_hash
	// block_time
}

type MempoolGetTransactionResponse struct {
	TxId        string `json:"txid"`
	Version     int    `json:"version"`
	Locktime    int   `json:"locktime"`
	// Vin         [] `json:"vin"`
	// Vout        [] `json:"vout"`
	// Size  int `json:"size"`
	// Weight  int `json:"weight"`
	// Fee int64 `json:"fee"`
	Status      MempoolGetTransactionStatus   `json:"status"`
}

type ResponseUnisatBroadcastTx struct {
	Code    int64  `json:"code"`
	Message string `json:"msg"`

	Data    string `json:"data"`
}

func BtcUtxoUrl(server, address string, offset, limit int64) string {
	return fmt.Sprintf("%s/v1/indexer/address/%s/utxo-data?cursor=%d&size=%d", server, address, offset, limit)
}

func RuneUtxoUrl(server, address, runeId string, offset, limit int64) string {
	return fmt.Sprintf("%s/v1/indexer/address/%s/runes/%s/utxo?start=%d&limit=%d", server, address, runeId, offset, limit)
}

func RuneBalanceByUtxoUrl(server, txid string, index uint32) string {
	return fmt.Sprintf("%s/v1/indexer/runes/utxo/%s/%d/balance", server, txid, index)
}

func RuneBalanceOfAddressUrl(server, address, runeId string) string {
	return fmt.Sprintf("%s/v1/indexer/address/%s/runes/%s/balance", server, address, runeId)
}

func TxInfoUrl(server, txid string) string {
	return fmt.Sprintf("%s/v1/indexer/tx/%s", server, txid)
}

func TxInputsUrl(server, txid string, offset, limit int64) string {
	return fmt.Sprintf("%s/v1/indexer/tx/%s/ins?cursor=%d&size=%d", server, txid, offset, limit)
}

func BitcoinBalanceByUtxoUrl(server, txid string, index uint32) string {
	return fmt.Sprintf("%s/v1/indexer/utxo/%s/%d", server, txid, index)
}

func UnisatWalletBroadcastTxUrl(server string) string {
	return fmt.Sprintf("%s/v5/tx/broadcast", server)
}

func MempoolRecommendedFeeUrl(server string) string {
	return fmt.Sprintf("%s/v1/fees/recommended", server)
}

func MempoolBroadcastTxUrl(server string) string {
	return fmt.Sprintf("%s/tx", server)
}

func MempoolGetTxUrl(server, txid string) string {
	return fmt.Sprintf("%s/tx/%s", server, txid)
}

func GetWithHeader(ctx context.Context, url string, header map[string]string, response interface{}) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}
	for k, v := range header {
		req.Header.Add(k, v)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(body, response); err != nil {
		return err
	}

	return nil
}

func GetWithBear(ctx context.Context, url, bear string, response interface{}) error {
	header := make(map[string]string)
	header["Authorization"] = fmt.Sprintf("Bearer %s", bear)

	return GetWithHeader(ctx, url, header, response)
}

func GetBtcUtxo(ctx context.Context, server, bear, address string, offset, limit int64) (ResponseBtcUtxo, error) {
	var resp ResponseBtcUtxo
	url := BtcUtxoUrl(server, address, offset, limit)
	err := GetWithBear(ctx, url, bear, &resp)

	return resp, err
}

func GetRuneUtxo(ctx context.Context, server, bear, address, runeId string, offset, limit int64) (ResponseRuneUtxo, error) {
	var resp ResponseRuneUtxo
	url := RuneUtxoUrl(server, address, runeId, offset, limit)
	err := GetWithBear(ctx, url, bear, &resp)

	return resp, err
}

func GetRuneBalanceByUtxo(ctx context.Context, server, bear, txid string, index uint32) (ResponseRuneBalanceByUtxo, error) {
	var resp ResponseRuneBalanceByUtxo
	url := RuneBalanceByUtxoUrl(server, txid, index)
	err := GetWithBear(ctx, url, bear, &resp)

	return resp, err
}

func GetRuneBalanceOfAddress(ctx context.Context, server, bear, address, runeId string) (ResponseRuneBalanceOfAddress, error) {
	var resp ResponseRuneBalanceOfAddress
	url := RuneBalanceOfAddressUrl(server, address, runeId)
	err := GetWithBear(ctx, url, bear, &resp)

	return resp, err
}

func GetBitcoinBalanceByUtxo(ctx context.Context, server, bear, txid string, index uint32) (ResponseBitcoinBalanceByUtxo, error) {
	var resp ResponseBitcoinBalanceByUtxo
	url := BitcoinBalanceByUtxoUrl(server, txid, index)
	err := GetWithBear(ctx, url, bear, &resp)

	return resp, err
}

func GetTxInfo(ctx context.Context, server, bear, txid string) (ResponseTxInfo, error) {
	var resp ResponseTxInfo
	url := TxInfoUrl(server, txid)
	err := GetWithBear(ctx, url, bear, &resp)

	return resp, err
}

func GetTxInputs(ctx context.Context, server, bear, txid string, offset, limit int64) (ResponseTxInputs, error) {
	var resp ResponseTxInputs
	url := TxInputsUrl(server, txid, offset, limit)
	err := GetWithBear(ctx, url, bear, &resp)

	return resp, err
}

func GetBtcInputs(amountRequired int64, addressPkScript []byte, utxos []UTXO) ([]*Input, error) {
	// sort utxos by amount in descending order
	sort.Slice(utxos, func(i, j int) bool {
		return utxos[i].Satoshi.Cmp(utxos[j].Satoshi) == 1
	})

	inputs := []*Input{}
	var totalAmount int64
	for _, utxo := range utxos {
		if totalAmount >= amountRequired {
			break
		}
		// check if utxo is already spent in a broadcasted tx
		if utxo.IsSpent {
			continue
		}
		outputAmount := utxo.Satoshi.Int64()
		if outputAmount <= MAX_DUST_UTXO_AMOUNT {
			continue
		}
		inputs = append(inputs, &Input{
			TxHash:             utxo.TxId,
			OutputIdx:          uint32(utxo.Vout),
			OutputAmount:       outputAmount,
			PkScript:			addressPkScript,
		})
		totalAmount += outputAmount
	}

	if totalAmount < amountRequired {
		return nil, fmt.Errorf("insufficient btc balance")
	}

	return inputs, nil
}

func GetRuneInputs(timeout int64, server, bear, address, runeId string, amountRequired uint128.Uint128, addressPkScript []byte, utxos []UTXO) ([]*Input, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(int64(time.Second)*timeout))
	defer cancel()
	// TODO: loop query until sastified amountRequired
	resp, err := GetRuneUtxo(ctx, server, bear, address, runeId, 0, 500)
	if err != nil {
		return nil, fmt.Errorf("failed to query rune UTXOs from unisat: %v", err)
	}
	runeUtxos := resp.Data.Utxo

	inputs := []*Input{}
	totalAmount := uint128.Zero
	for _, runeUtxo := range runeUtxos {
		if totalAmount.Cmp(amountRequired) >= 0 {
			break
		}
		// check if rune utxo is already spent in a broadcasted tx
		isSpent := true
		for _, utxo := range utxos {
			if runeUtxo.TxId == utxo.TxId && runeUtxo.Vout == utxo.Vout {
				isSpent = utxo.IsSpent
				break
			}
		}
		if isSpent {
			continue
		}

		for _, utxoRune := range runeUtxo.Runes {
			if utxoRune.RuneId == runeId {
				runeAmount, err := uint128.FromString(utxoRune.Amount)
				if err != nil {
					return nil, fmt.Errorf("failed to query rune amount in an UTXO from unisat: %v", err)
				}

				inputs = append(inputs, &Input{
					TxHash:             runeUtxo.TxId,
					OutputIdx:          uint32(runeUtxo.Vout),
					OutputAmount:       runeUtxo.Satoshi.Int64(),
					PkScript:			addressPkScript,
				})
				totalAmount = totalAmount.Add(runeAmount)

				break
			}
		}
	}

	if totalAmount.Cmp(amountRequired) < 0 {
		return nil, fmt.Errorf("insufficient rune balance, need " + amountRequired.String() + ", have " + totalAmount.String())
	}

	return inputs, nil
}

func SelectUTXOs(timeout int64, server, bear, address, runeId string, runeToSend uint128.Uint128, satToSend int64, addressPkScript []byte) ([]*Input, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(int64(time.Second)*timeout))
	defer cancel()
	// TODO: loop query until sastified amountRequired
	resp, err := GetBtcUtxo(ctx, server, bear, address, 0, 1000)
	if err != nil {
		return nil, fmt.Errorf("failed to query bitcoin UTXOs from unisat: %v", err)
	}
	utxos := resp.Data.Utxo

	// add tx fee the the required bitcoin amount
	inputs := []*Input{}
	if !runeToSend.IsZero() {
		// query rune UTXOs from unisat
		runeInputs, err := GetRuneInputs(timeout, server, bear, address, runeId, runeToSend, addressPkScript, utxos)
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, runeInputs...)
	}

	// query bitcoin UTXOs from unisat
	bitcoinInputs, err := GetBtcInputs(satToSend, addressPkScript, utxos)
	if err != nil {
		return nil, err
	}
	inputs = append(inputs, bitcoinInputs...)

	return inputs, nil
}

func GetRunesInUtxo(timeout int64, server, bear, txid string, index uint32) ([]RuneDetail, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(int64(time.Second)*timeout))
	defer cancel()

	resp, err := GetRuneBalanceByUtxo(ctx, server, bear, txid, index)
	if err != nil {
		return nil, fmt.Errorf("failed to query RuneBalanceByUtxo from unisat: %v", err)
	}

	return resp.Data, nil
}

func GetBitcoinInUtxo(timeout int64, server, bear, txid string, index uint32) (*UTXO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(int64(time.Second)*timeout))
	defer cancel()

	resp, err := GetBitcoinBalanceByUtxo(ctx, server, bear, txid, index)
	if err != nil {
		return nil, fmt.Errorf("failed to query RuneBalanceByUtxo from unisat: %v", err)
	}

	return &resp.Data, nil
}

func GetRunesBalanceOf(timeout int64, server, bear, address, runeId string) (uint128.Uint128, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(int64(time.Second)*timeout))
	defer cancel()

	resp, err := GetRuneBalanceOfAddress(ctx, server, bear, address, runeId)
	if err != nil {
		return uint128.Zero, fmt.Errorf("failed to query GetRuneBalanceOfAddress from unisat: %v", err)
	}

	return uint128.FromString(resp.Data.Amount)
}

func GetTxConfirmations(timeout int64, server, bear, txid string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(int64(time.Second)*timeout))
	defer cancel()

	resp, err := GetTxInfo(ctx, server, bear, txid)
	if err != nil {
		return 0, fmt.Errorf("failed to query bitcoin tx info from unisat: %v", err)
	}
	return resp.Data.Confirmations, nil
}

func CheckIfTxExistedOnUnisat(timeout int64, server, bear, txid string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(int64(time.Second)*timeout))
	defer cancel()

	resp, err := GetTxInfo(ctx, server, bear, txid)
	if err != nil {
		return false, fmt.Errorf("failed to query bitcoin tx info from unisat: %v", err)
	}

	if resp.Code == -1 && resp.Message == "get tx failed" {
		return false, nil
	}
	if resp.Code == 0 && resp.Message == "ok" {
		return true, nil
	}

	return false, fmt.Errorf("failed to query bitcoin tx info from unisat: code %v, message %v", resp.Code, resp.Message)
}

func GetTxFirstInputAddr(timeout int64, server, bear, txid string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(int64(time.Second)*timeout))
	defer cancel()

	resp, err := GetTxInputs(ctx, server, bear, txid, 0, 1)
	if err != nil {
		return "", fmt.Errorf("failed to query bitcoin tx input from unisat: %v", err)
	}
	if len(resp.Data) == 0 {
		return "", fmt.Errorf("failed to query bitcoin tx input from unisat: len(resp.Data) = 0")
	}

	return resp.Data[0].Address, nil
}

func GetFeeFromMempool(timeout int64, server string) (MempoolRecommendedFeeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(int64(time.Second)*timeout))
	defer cancel()

	var resp MempoolRecommendedFeeResponse
	url := MempoolRecommendedFeeUrl(server)
	err := GetWithHeader(ctx, url, make(map[string]string), &resp)

	return resp, err
}

func GetTransactionFromMempool(timeout int64, server, txid string) (MempoolGetTransactionResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(int64(time.Second)*timeout))
	defer cancel()

	var resp MempoolGetTransactionResponse
	url := MempoolGetTxUrl(server, txid)
	err := GetWithHeader(ctx, url, make(map[string]string), &resp)

	return resp, err
}

// calculateTxSize calculates the size of a transaction given the inputs, outputs, estimated fee, change address, chain parameters, and multisig wallet.
// It returns the size of the transaction in bytes and an error if any occurs during the process.
func CalculateTxSize(inputs []*Input, outputs []*wire.TxOut, PkKScript []byte, estFee int64, witnessSize int) (int, error) {
	msgTx, err := CreateTx(inputs, outputs, PkKScript, estFee, 0)
	if err != nil {
		return 0, err
	}
	var rawTxBytes bytes.Buffer
	err = msgTx.Serialize(&rawTxBytes)
	if err != nil {
		return 0, err
	}
	txSize := len(rawTxBytes.Bytes()) + len(inputs)*witnessSize
	return txSize, nil
}

func calEstFee(txSize int, fastestFeeRate int64) int64 {
	estFee := int64(txSize) * fastestFeeRate

	return estFee
}

func BuildTxInputs(timeout int64, mempoolUrl, unisatUrl, unisatBear, address, runeId string, runeToSend uint128.Uint128, satToSend int64, outputs []*wire.TxOut, addressPkKScript []byte, witnessSize int) ([]*Input, int64, error) {
	feeRate, err := GetFeeFromMempool(timeout, mempoolUrl)
	if err != nil {
		return nil, 0, err
	}
	fastestFeeRate := int64(feeRate.FastestFee)

	inputs, err := SelectUTXOs(timeout, unisatUrl, unisatBear, address, runeId, runeToSend, satToSend, addressPkKScript)
	sumSelectedInputs := SumInputsSat(inputs)
	if err != nil {
		return nil, 0, err
	}

	txSize, err := CalculateTxSize(inputs, outputs, addressPkKScript, 0, witnessSize)
	if err != nil {
		return nil, 0, err
	}

	estFee := calEstFee(txSize, fastestFeeRate)
	count := 0
	for sumSelectedInputs < satToSend+estFee {
		inputs, err = SelectUTXOs(timeout, unisatUrl, unisatBear, address, runeId, runeToSend, satToSend + estFee, addressPkKScript)
		if err != nil {
			return nil, 0, err
		}

		sumSelectedInputs = SumInputsSat(inputs)

		txSize, err := CalculateTxSize(inputs, outputs, addressPkKScript, estFee, witnessSize)
		if err != nil {
			return nil, 0, err
		}
		estFee = calEstFee(txSize, fastestFeeRate)

		count += 1
		if count > 100 {
			return nil, 0, fmt.Errorf("not enough sats for fee")
		}
	}

	return inputs, estFee, nil
}

func UnisatBroadcastTransaction(server string, signedMsgTx *wire.MsgTx) (string, error) {
	var signedTx bytes.Buffer
	err := signedMsgTx.Serialize(&signedTx)
	if err != nil {
		return "", err
	}

	signedTxHex := hex.EncodeToString(signedTx.Bytes())
	rawMsg := []byte(`{"rawtx":"` + signedTxHex + `"}`)

	url := UnisatWalletBroadcastTxUrl(server)
	resp, err := http.Post(url, "application/json;charset=utf-8", bytes.NewReader(rawMsg))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("broadcast failed: %v", string(body))
	}

	var bodyData ResponseUnisatBroadcastTx
	if err = json.Unmarshal(body, &bodyData); err != nil {
		return "", err
	}

	if bodyData.Code == -1 {
		return "", fmt.Errorf("broadcast failed: %v", string(bodyData.Message))
	}

	if bodyData.Code == 0 && bodyData.Message == "ok" {
		return bodyData.Data, nil
	}

	return "", fmt.Errorf("broadcast failed: code %v, message %v", bodyData.Code, bodyData.Message)
}

func GetInputOutputBalance(timeout int64, server, bear string, transaction *wire.MsgTx, runeArtifact *runestone.Artifact) (map[string]map[string]uint128.Uint128, map[string]map[string]uint128.Uint128, string, error) {
	bitcoinId := BitcoinRuneId().String()

	inputsTokenBalance := make(map[string]map[string]uint128.Uint128)
	totalInputsTokenBalance := make(map[string]uint128.Uint128)
	for _, txIn := range(transaction.TxIn) {
		runeBalances, err := GetRunesInUtxo(timeout, server, bear, txIn.PreviousOutPoint.Hash.String(), txIn.PreviousOutPoint.Index)
		if err != nil {
			return nil, nil, "", err
		}
		bitcoinInfo, err := GetBitcoinInUtxo(timeout, server, bear, txIn.PreviousOutPoint.Hash.String(), txIn.PreviousOutPoint.Index)
		if err != nil {
			return nil, nil, "", err
		}
		// make sure output is not spent
		if bitcoinInfo.IsSpent {
			return nil, nil, "", err
		}

		inputsTokenBalance[bitcoinInfo.ScriptPk][bitcoinId] = inputsTokenBalance[bitcoinInfo.ScriptPk][bitcoinId].Add(uint128.FromBig(bitcoinInfo.Satoshi))
		totalInputsTokenBalance[bitcoinId] = totalInputsTokenBalance[bitcoinId].Add(uint128.FromBig(bitcoinInfo.Satoshi))
		for _, runeBalance := range(runeBalances) {
			runeAmount, err := uint128.FromString(runeBalance.Amount)
			if err != nil {
				return nil, nil, "", err
			}
			inputsTokenBalance[bitcoinInfo.ScriptPk][runeBalance.RuneId] = inputsTokenBalance[bitcoinInfo.ScriptPk][runeBalance.RuneId].Add(runeAmount)
			totalInputsTokenBalance[runeBalance.RuneId] = totalInputsTokenBalance[runeBalance.RuneId].Add(runeAmount)
		}
	}

	outputsTokenBalance := make(map[string]map[string]uint128.Uint128)
	totalOutputsTokenBalance := make(map[string]uint128.Uint128)
	for _, txOut := range(transaction.TxOut) {
		outPkScript := hex.EncodeToString(txOut.PkScript)
		outputsTokenBalance[outPkScript][bitcoinId] = outputsTokenBalance[outPkScript][bitcoinId].Add(uint128.From64(uint64(txOut.Value)))
		totalOutputsTokenBalance[bitcoinId] = totalOutputsTokenBalance[bitcoinId].Add(uint128.From64(uint64(txOut.Value)))
	}
	for _, runeEdict := range(runeArtifact.Runestone.Edicts) {
		outPkScript := hex.EncodeToString(transaction.TxOut[runeEdict.Output].PkScript)
		outputsTokenBalance[outPkScript][runeEdict.ID.String()] = outputsTokenBalance[outPkScript][runeEdict.ID.String()].Add(runeEdict.Amount)
		totalOutputsTokenBalance[runeEdict.ID.String()] = totalOutputsTokenBalance[runeEdict.ID.String()].Add(runeEdict.Amount)
	}

	// make sure total token in = total token out (the rune change amount in the output >= 0)
	runeChangePKScript := hex.EncodeToString(transaction.TxOut[*runeArtifact.Runestone.Pointer].PkScript)
	for outputTokenId, outputTokenValue := range(totalOutputsTokenBalance) {
		inputTokenValue := totalInputsTokenBalance[outputTokenId]
		if outputTokenValue.Cmp(inputTokenValue) > 0 {
			return nil, nil, "", fmt.Errorf("outputTokenValue (%v) > inputTokenValue (%v)", outputTokenValue.String(), inputTokenValue.String())
		}
		// add token change amount to output balance (include bitcoin tx fee of user)
		outputsTokenBalance[runeChangePKScript][outputTokenId] = outputsTokenBalance[runeChangePKScript][outputTokenId].Add(inputTokenValue.Sub(outputTokenValue))
	}

	return inputsTokenBalance, outputsTokenBalance, runeChangePKScript, nil
}

// call testmempoolaccept RPC to validate the UTXO and bitcoin amount
func VerifyRadfiTx(timeout int64, server, bear string, relayersMultisigInfo *MultisigInfo, transaction *wire.MsgTx) (*RadFiDecodedMsg, error) {
	// Decipher runestone
	r := &runestone.Runestone{}
	runeArtifact, err := r.Decipher(transaction)
	if err != nil {
		return nil, fmt.Errorf("could not decipher runestone - Error %v", err)
	}

	inputsTokenBalance, outputsTokenBalance, runeChangePKScript, err := GetInputOutputBalance(timeout, server, bear, transaction, runeArtifact)
	if err != nil {
		return nil, err
	}

	// decode tx
	radFiMessage, err := ReadRadFiMessage(transaction)
	if err != nil {
		return nil, err
	}
	// verify tx data
	switch radFiMessage.Flag {
		case OP_RADFI_PROVIDE_LIQUIDITY:
			plMessage := radFiMessage.ProvideLiquidityMsg
			poolWalletPkScript, err := GetPoolWalletPkScript(relayersMultisigInfo, plMessage.NftId)
			if err != nil {
				return nil, err
			}

			// check if pool liquidity really increased by the amounts in radfiMsg
			if !outputsTokenBalance[poolWalletPkScript][plMessage.Token0Id.String()].Equals(inputsTokenBalance[poolWalletPkScript][plMessage.Token0Id.String()].Add(plMessage.Amount0Desired)) {
				return nil, fmt.Errorf("OP_RADFI_PROVIDE_LIQUIDITY - the amount0 mismatch with the liquidity pool received")
			}
			if !outputsTokenBalance[poolWalletPkScript][plMessage.Token1Id.String()].Equals(inputsTokenBalance[poolWalletPkScript][plMessage.Token1Id.String()].Add(plMessage.Amount1Desired)) {
				return nil, fmt.Errorf("OP_RADFI_PROVIDE_LIQUIDITY - the amount1 mismatch with the liquidity pool received")
			}
			// check if user rune balance really decreased
			for userTokenId, userTokenBalanceOut := range(outputsTokenBalance[runeChangePKScript]) {
				if userTokenId == plMessage.Token0Id.String() {
					if !userTokenBalanceOut.Add(plMessage.Amount0Desired).Equals(inputsTokenBalance[runeChangePKScript][userTokenId]) {
						return nil, fmt.Errorf("OP_RADFI_PROVIDE_LIQUIDITY - the amount0 mismatch with the rune user provided")
					}
				} else if userTokenId == plMessage.Token1Id.String() {
					if !userTokenBalanceOut.Add(plMessage.Amount1Desired).Equals(inputsTokenBalance[runeChangePKScript][userTokenId]) {
						return nil, fmt.Errorf("OP_RADFI_PROVIDE_LIQUIDITY - the amount1 mismatch with the rune user provided")
					}
				} else {
					// the other token balance not for the pool should be unchanged
					if !userTokenBalanceOut.Equals(inputsTokenBalance[runeChangePKScript][userTokenId]) {
						return nil, fmt.Errorf("OP_RADFI_PROVIDE_LIQUIDITY - the non liquidity token balance should not changed")
					}
				}
			}

		case OP_RADFI_SWAP:
			plMessage := radFiMessage.SwapMsg
			poolWalletPkScripts := []string{}
	
			// check if pool liquidity really increased and decreased by the amounts in radfiMsg
			if !outputsTokenBalance[poolWalletPkScripts[0]][plMessage.Tokens[0].String()].Equals(inputsTokenBalance[poolWalletPkScripts[0]][plMessage.Tokens[0].String()].Add(plMessage.AmountIn)) {
				return nil, fmt.Errorf("OP_RADFI_SWAP - the AmountIn mismatch with the amount pool0 received")
			}
			if !outputsTokenBalance[poolWalletPkScripts[len(poolWalletPkScripts)-1]][plMessage.Tokens[len(plMessage.Tokens)-1].String()].Add(plMessage.AmountOut).Equals(inputsTokenBalance[poolWalletPkScripts[len(poolWalletPkScripts)-1]][plMessage.Tokens[len(plMessage.Tokens)-1].String()]) {
				return nil, fmt.Errorf("OP_RADFI_SWAP - the AmountIn mismatch with the amount pool0 received")
			}
			// check if user rune balance really changed
			for userTokenId, userTokenBalanceOut := range(outputsTokenBalance[runeChangePKScript]) {
				if userTokenId == plMessage.Tokens[0].String() {
					if !userTokenBalanceOut.Add(plMessage.AmountIn).Equals(inputsTokenBalance[runeChangePKScript][userTokenId]) {
						return nil, fmt.Errorf("OP_RADFI_SWAP - the amount0 mismatch with the rune user swap")
					}
				} else if userTokenId == plMessage.Tokens[len(plMessage.Tokens)-1].String() {
					if !userTokenBalanceOut.Equals(inputsTokenBalance[runeChangePKScript][userTokenId].Add(plMessage.AmountOut)) {
						return nil, fmt.Errorf("OP_RADFI_SWAP - the amount1 mismatch with the rune user swap")
					}
				} else {
					// the other token balance not for swap should be unchanged
					if !userTokenBalanceOut.Equals(inputsTokenBalance[runeChangePKScript][userTokenId]) {
						return nil, fmt.Errorf("OP_RADFI_SWAP - the non liquidity token balance should not changed")
					}
				}
			}

		case OP_RADFI_WITHDRAW_LIQUIDITY:
			plMessage := radFiMessage.WithdrawLiquidityMsg
			poolWalletPkScript, err := GetPoolWalletPkScript(relayersMultisigInfo, plMessage.NftId)
			if err != nil {
				return nil, err
			}

			// check if pool liquidity really decreased by the amounts in radfiMsg
			if !outputsTokenBalance[poolWalletPkScript][plMessage.Token0Id.String()].Add(plMessage.Amount0).Equals(inputsTokenBalance[poolWalletPkScript][plMessage.Token0Id.String()]) {
				return nil, fmt.Errorf("OP_RADFI_WITHDRAW_LIQUIDITY - the amount0 mismatch with the liquidity pool withdrawed")
			}
			if !outputsTokenBalance[poolWalletPkScript][plMessage.Token1Id.String()].Add(plMessage.Amount1).Equals(inputsTokenBalance[poolWalletPkScript][plMessage.Token1Id.String()]) {
				return nil, fmt.Errorf("OP_RADFI_WITHDRAW_LIQUIDITY - the amount1 mismatch with the liquidity pool withdrawed")
			}
			// check if user rune balance really increased
			for userTokenId, userTokenBalanceOut := range(outputsTokenBalance[runeChangePKScript]) {
				if userTokenId == plMessage.Token0Id.String() {
					if !userTokenBalanceOut.Equals(inputsTokenBalance[runeChangePKScript][userTokenId].Add(plMessage.Amount0)) {
						return nil, fmt.Errorf("OP_RADFI_PROVIDE_LIQUIDITY - the amount0 mismatch with the rune user received")
					}
				} else if userTokenId == plMessage.Token1Id.String() {
					if !userTokenBalanceOut.Equals(inputsTokenBalance[runeChangePKScript][userTokenId].Add(plMessage.Amount1)) {
						return nil, fmt.Errorf("OP_RADFI_PROVIDE_LIQUIDITY - the amount1 mismatch with the rune user received")
					}
				} else {
					// the other token balance not for the pool should be unchanged
					if !userTokenBalanceOut.Equals(inputsTokenBalance[runeChangePKScript][userTokenId]) {
						return nil, fmt.Errorf("OP_RADFI_PROVIDE_LIQUIDITY - the non liquidity token balance should not changed")
					}
				}
			}

		case OP_RADFI_COLLECT_FEES:
			plMessage := radFiMessage.CollectFeesMsg
			poolWalletPkScript, err := GetPoolWalletPkScript(relayersMultisigInfo, plMessage.NftId)
			if err != nil {
				return nil, err
			}

			// check if pool liquidity really decreased by the amounts in radfiMsg
			if !outputsTokenBalance[poolWalletPkScript][plMessage.Token0Id.String()].Add(plMessage.Amount0).Equals(inputsTokenBalance[poolWalletPkScript][plMessage.Token0Id.String()]) {
				return nil, fmt.Errorf("OP_RADFI_COLLECT_FEES - the amount0 mismatch with the liquidity pool withdrawed")
			}
			if !outputsTokenBalance[poolWalletPkScript][plMessage.Token1Id.String()].Add(plMessage.Amount1).Equals(inputsTokenBalance[poolWalletPkScript][plMessage.Token1Id.String()]) {
				return nil, fmt.Errorf("OP_RADFI_COLLECT_FEES - the amount1 mismatch with the liquidity pool withdrawed")
			}
			// check if user rune balance really increased
			for userTokenId, userTokenBalanceOut := range(outputsTokenBalance[runeChangePKScript]) {
				if userTokenId == plMessage.Token0Id.String() {
					if !userTokenBalanceOut.Equals(inputsTokenBalance[runeChangePKScript][userTokenId].Add(plMessage.Amount0)) {
						return nil, fmt.Errorf("OP_RADFI_COLLECT_FEES - the amount0 mismatch with the rune user received")
					}
				} else if userTokenId == plMessage.Token1Id.String() {
					if !userTokenBalanceOut.Equals(inputsTokenBalance[runeChangePKScript][userTokenId].Add(plMessage.Amount1)) {
						return nil, fmt.Errorf("OP_RADFI_COLLECT_FEES - the amount1 mismatch with the rune user received")
					}
				} else {
					// the other token balance not for the pool should be unchanged
					if !userTokenBalanceOut.Equals(inputsTokenBalance[runeChangePKScript][userTokenId]) {
						return nil, fmt.Errorf("OP_RADFI_COLLECT_FEES - the non liquidity token balance should not changed")
					}
				}
			}

		case OP_RADFI_INCREASE_LIQUIDITY:
			plMessage := radFiMessage.IncreaseLiquidityMsg
			poolWalletPkScript, err := GetPoolWalletPkScript(relayersMultisigInfo, plMessage.NftId)
			if err != nil {
				return nil, err
			}

			// check if pool liquidity really increased by the amounts in radfiMsg
			if !outputsTokenBalance[poolWalletPkScript][plMessage.Token0Id.String()].Equals(inputsTokenBalance[poolWalletPkScript][plMessage.Token0Id.String()].Add(plMessage.Amount0)) {
				return nil, fmt.Errorf("OP_RADFI_INCREASE_LIQUIDITY - the amount0 mismatch with the liquidity pool received")
			}
			if !outputsTokenBalance[poolWalletPkScript][plMessage.Token1Id.String()].Equals(inputsTokenBalance[poolWalletPkScript][plMessage.Token1Id.String()].Add(plMessage.Amount1)) {
				return nil, fmt.Errorf("OP_RADFI_INCREASE_LIQUIDITY - the amount1 mismatch with the liquidity pool received")
			}
			// check if user rune balance really decreased
			for userTokenId, userTokenBalanceOut := range(outputsTokenBalance[runeChangePKScript]) {
				if userTokenId == plMessage.Token0Id.String() {
					if !userTokenBalanceOut.Add(plMessage.Amount0).Equals(inputsTokenBalance[runeChangePKScript][userTokenId]) {
						return nil, fmt.Errorf("OP_RADFI_INCREASE_LIQUIDITY - the amount0 mismatch with the rune user provided")
					}
				} else if userTokenId == plMessage.Token1Id.String() {
					if !userTokenBalanceOut.Add(plMessage.Amount1).Equals(inputsTokenBalance[runeChangePKScript][userTokenId]) {
						return nil, fmt.Errorf("OP_RADFI_INCREASE_LIQUIDITY - the amount1 mismatch with the rune user provided")
					}
				} else {
					// the other token balance not for the pool should be unchanged
					if !userTokenBalanceOut.Equals(inputsTokenBalance[runeChangePKScript][userTokenId]) {
						return nil, fmt.Errorf("OP_RADFI_INCREASE_LIQUIDITY - the non liquidity token balance should not changed")
					}
				}
			}

		default:
			return nil, fmt.Errorf("ReadRadFiMessage - invalid flag")
	}

	return radFiMessage, nil
}
