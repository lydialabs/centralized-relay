package cmd

// Transaction holds the entire structure of the transaction JSON.
type Transaction struct {
	TxID     string `json:"txid"`
	Version  int    `json:"version"`
	LockTime int    `json:"locktime"`
	Vin      []Vin  `json:"vin"`
	Vout     []Vout `json:"vout"`
	Size     int    `json:"size"`
	Weight   int    `json:"weight"`
	SigOps   int    `json:"sigops"`
	Fee      int    `json:"fee"`
	Status   Status `json:"status"`
}

// Vin holds the structure for input transactions.
type Vin struct {
	TxID                  string   `json:"txid"`
	Vout                  int      `json:"vout"`
	PrevOut               PrevOut  `json:"prevout"`
	ScriptSig             string   `json:"scriptsig"`
	ScriptSigAsm          string   `json:"scriptsig_asm"`
	Witness               []string `json:"witness"`
	IsCoinbase            bool     `json:"is_coinbase"`
	Sequence              uint32   `json:"sequence"`
	InnerWitnessScriptAsm string   `json:"inner_witnessscript_asm"`
}

// PrevOut holds details about previous output referred by an input.
type PrevOut struct {
	ScriptPubKey        string `json:"scriptpubkey"`
	ScriptPubKeyAsm     string `json:"scriptpubkey_asm"`
	ScriptPubKeyType    string `json:"scriptpubkey_type"`
	ScriptPubKeyAddress string `json:"scriptpubkey_address"`
	Value               int    `json:"value"`
}

// Vout holds the structure for output transactions.
type Vout struct {
	ScriptPubKey        string `json:"scriptpubkey"`
	ScriptPubKeyAsm     string `json:"scriptpubkey_asm"`
	ScriptPubKeyType    string `json:"scriptpubkey_type"`
	ScriptPubKeyAddress string `json:"scriptpubkey_address"`
	Value               int    `json:"value"`
}

// Status holds details about the transaction confirmation.
type Status struct {
	Confirmed   bool   `json:"confirmed"`
	BlockHeight int    `json:"block_height"`
	BlockHash   string `json:"block_hash"`
	BlockTime   int64  `json:"block_time"`
}
