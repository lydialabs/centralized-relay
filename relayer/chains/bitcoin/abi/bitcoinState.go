// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package abi

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// BitcoinstateMetaData contains all meta data concerning the Bitcoinstate contract.
var BitcoinstateMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"accountBalances\",\"inputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"addConnection\",\"inputs\":[{\"name\":\"connection_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"bitcoinNid\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"claimTokens\",\"inputs\":[{\"name\":\"token0\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"token1\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"connections\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"connectionsEndpoints\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"handleCallMessage\",\"inputs\":[{\"name\":\"_from\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_protocols\",\"type\":\"string[]\",\"internalType\":\"string[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initPool\",\"inputs\":[{\"name\":\"data_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"xcall_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"uinswapV3Router_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nonfungiblePositionManager_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"connections\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"migrate\",\"inputs\":[{\"name\":\"from_\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"migrateComplete\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"nftOwners\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"nonFungibleManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"removeConnection\",\"inputs\":[{\"name\":\"connection_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"removeLiquidity\",\"inputs\":[{\"name\":\"data_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"requestCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"routerV2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"runeFactory\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setRuneFactory\",\"inputs\":[{\"name\":\"_runeFactory\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"simulateRequest\",\"inputs\":[{\"name\":\"sequenceNumber\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_from\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"parseData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"tokens\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawTo\",\"inputs\":[{\"name\":\"token_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"receiver_\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"xcallService\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"AddConnection\",\"inputs\":[{\"name\":\"connection_\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RemoveConnection\",\"inputs\":[{\"name\":\"connection_\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RequestExecuted\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"stateRoot\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"data\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false}]",
}

// BitcoinstateABI is the input ABI used to generate the binding from.
// Deprecated: Use BitcoinstateMetaData.ABI instead.
var BitcoinstateABI = BitcoinstateMetaData.ABI

// Bitcoinstate is an auto generated Go binding around an Ethereum contract.
type Bitcoinstate struct {
	BitcoinstateCaller     // Read-only binding to the contract
	BitcoinstateTransactor // Write-only binding to the contract
	BitcoinstateFilterer   // Log filterer for contract events
}

// BitcoinstateCaller is an auto generated read-only Go binding around an Ethereum contract.
type BitcoinstateCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BitcoinstateTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BitcoinstateTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BitcoinstateFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BitcoinstateFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BitcoinstateSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BitcoinstateSession struct {
	Contract     *Bitcoinstate     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BitcoinstateCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BitcoinstateCallerSession struct {
	Contract *BitcoinstateCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// BitcoinstateTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BitcoinstateTransactorSession struct {
	Contract     *BitcoinstateTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// BitcoinstateRaw is an auto generated low-level Go binding around an Ethereum contract.
type BitcoinstateRaw struct {
	Contract *Bitcoinstate // Generic contract binding to access the raw methods on
}

// BitcoinstateCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BitcoinstateCallerRaw struct {
	Contract *BitcoinstateCaller // Generic read-only contract binding to access the raw methods on
}

// BitcoinstateTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BitcoinstateTransactorRaw struct {
	Contract *BitcoinstateTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBitcoinstate creates a new instance of Bitcoinstate, bound to a specific deployed contract.
func NewBitcoinstate(address common.Address, backend bind.ContractBackend) (*Bitcoinstate, error) {
	contract, err := bindBitcoinstate(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Bitcoinstate{BitcoinstateCaller: BitcoinstateCaller{contract: contract}, BitcoinstateTransactor: BitcoinstateTransactor{contract: contract}, BitcoinstateFilterer: BitcoinstateFilterer{contract: contract}}, nil
}

// NewBitcoinstateCaller creates a new read-only instance of Bitcoinstate, bound to a specific deployed contract.
func NewBitcoinstateCaller(address common.Address, caller bind.ContractCaller) (*BitcoinstateCaller, error) {
	contract, err := bindBitcoinstate(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BitcoinstateCaller{contract: contract}, nil
}

// NewBitcoinstateTransactor creates a new write-only instance of Bitcoinstate, bound to a specific deployed contract.
func NewBitcoinstateTransactor(address common.Address, transactor bind.ContractTransactor) (*BitcoinstateTransactor, error) {
	contract, err := bindBitcoinstate(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BitcoinstateTransactor{contract: contract}, nil
}

// NewBitcoinstateFilterer creates a new log filterer instance of Bitcoinstate, bound to a specific deployed contract.
func NewBitcoinstateFilterer(address common.Address, filterer bind.ContractFilterer) (*BitcoinstateFilterer, error) {
	contract, err := bindBitcoinstate(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BitcoinstateFilterer{contract: contract}, nil
}

// bindBitcoinstate binds a generic wrapper to an already deployed contract.
func bindBitcoinstate(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BitcoinstateMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Bitcoinstate *BitcoinstateRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Bitcoinstate.Contract.BitcoinstateCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Bitcoinstate *BitcoinstateRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bitcoinstate.Contract.BitcoinstateTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Bitcoinstate *BitcoinstateRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Bitcoinstate.Contract.BitcoinstateTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Bitcoinstate *BitcoinstateCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Bitcoinstate.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Bitcoinstate *BitcoinstateTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bitcoinstate.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Bitcoinstate *BitcoinstateTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Bitcoinstate.Contract.contract.Transact(opts, method, params...)
}

// AccountBalances is a free data retrieval call binding the contract method 0x69477464.
//
// Solidity: function accountBalances(string , address ) view returns(uint256)
func (_Bitcoinstate *BitcoinstateCaller) AccountBalances(opts *bind.CallOpts, arg0 string, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Bitcoinstate.contract.Call(opts, &out, "accountBalances", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AccountBalances is a free data retrieval call binding the contract method 0x69477464.
//
// Solidity: function accountBalances(string , address ) view returns(uint256)
func (_Bitcoinstate *BitcoinstateSession) AccountBalances(arg0 string, arg1 common.Address) (*big.Int, error) {
	return _Bitcoinstate.Contract.AccountBalances(&_Bitcoinstate.CallOpts, arg0, arg1)
}

// AccountBalances is a free data retrieval call binding the contract method 0x69477464.
//
// Solidity: function accountBalances(string , address ) view returns(uint256)
func (_Bitcoinstate *BitcoinstateCallerSession) AccountBalances(arg0 string, arg1 common.Address) (*big.Int, error) {
	return _Bitcoinstate.Contract.AccountBalances(&_Bitcoinstate.CallOpts, arg0, arg1)
}

// BitcoinNid is a free data retrieval call binding the contract method 0xf3cc4515.
//
// Solidity: function bitcoinNid() view returns(string)
func (_Bitcoinstate *BitcoinstateCaller) BitcoinNid(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Bitcoinstate.contract.Call(opts, &out, "bitcoinNid")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// BitcoinNid is a free data retrieval call binding the contract method 0xf3cc4515.
//
// Solidity: function bitcoinNid() view returns(string)
func (_Bitcoinstate *BitcoinstateSession) BitcoinNid() (string, error) {
	return _Bitcoinstate.Contract.BitcoinNid(&_Bitcoinstate.CallOpts)
}

// BitcoinNid is a free data retrieval call binding the contract method 0xf3cc4515.
//
// Solidity: function bitcoinNid() view returns(string)
func (_Bitcoinstate *BitcoinstateCallerSession) BitcoinNid() (string, error) {
	return _Bitcoinstate.Contract.BitcoinNid(&_Bitcoinstate.CallOpts)
}

// ClaimTokens is a free data retrieval call binding the contract method 0x69ffa08a.
//
// Solidity: function claimTokens(address token0, address token1) pure returns()
func (_Bitcoinstate *BitcoinstateCaller) ClaimTokens(opts *bind.CallOpts, token0 common.Address, token1 common.Address) error {
	var out []interface{}
	err := _Bitcoinstate.contract.Call(opts, &out, "claimTokens", token0, token1)

	if err != nil {
		return err
	}

	return err

}

// ClaimTokens is a free data retrieval call binding the contract method 0x69ffa08a.
//
// Solidity: function claimTokens(address token0, address token1) pure returns()
func (_Bitcoinstate *BitcoinstateSession) ClaimTokens(token0 common.Address, token1 common.Address) error {
	return _Bitcoinstate.Contract.ClaimTokens(&_Bitcoinstate.CallOpts, token0, token1)
}

// ClaimTokens is a free data retrieval call binding the contract method 0x69ffa08a.
//
// Solidity: function claimTokens(address token0, address token1) pure returns()
func (_Bitcoinstate *BitcoinstateCallerSession) ClaimTokens(token0 common.Address, token1 common.Address) error {
	return _Bitcoinstate.Contract.ClaimTokens(&_Bitcoinstate.CallOpts, token0, token1)
}

// Connections is a free data retrieval call binding the contract method 0xc0896578.
//
// Solidity: function connections(address ) view returns(bool)
func (_Bitcoinstate *BitcoinstateCaller) Connections(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _Bitcoinstate.contract.Call(opts, &out, "connections", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Connections is a free data retrieval call binding the contract method 0xc0896578.
//
// Solidity: function connections(address ) view returns(bool)
func (_Bitcoinstate *BitcoinstateSession) Connections(arg0 common.Address) (bool, error) {
	return _Bitcoinstate.Contract.Connections(&_Bitcoinstate.CallOpts, arg0)
}

// Connections is a free data retrieval call binding the contract method 0xc0896578.
//
// Solidity: function connections(address ) view returns(bool)
func (_Bitcoinstate *BitcoinstateCallerSession) Connections(arg0 common.Address) (bool, error) {
	return _Bitcoinstate.Contract.Connections(&_Bitcoinstate.CallOpts, arg0)
}

// ConnectionsEndpoints is a free data retrieval call binding the contract method 0xb687db2e.
//
// Solidity: function connectionsEndpoints(uint256 ) view returns(address)
func (_Bitcoinstate *BitcoinstateCaller) ConnectionsEndpoints(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Bitcoinstate.contract.Call(opts, &out, "connectionsEndpoints", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ConnectionsEndpoints is a free data retrieval call binding the contract method 0xb687db2e.
//
// Solidity: function connectionsEndpoints(uint256 ) view returns(address)
func (_Bitcoinstate *BitcoinstateSession) ConnectionsEndpoints(arg0 *big.Int) (common.Address, error) {
	return _Bitcoinstate.Contract.ConnectionsEndpoints(&_Bitcoinstate.CallOpts, arg0)
}

// ConnectionsEndpoints is a free data retrieval call binding the contract method 0xb687db2e.
//
// Solidity: function connectionsEndpoints(uint256 ) view returns(address)
func (_Bitcoinstate *BitcoinstateCallerSession) ConnectionsEndpoints(arg0 *big.Int) (common.Address, error) {
	return _Bitcoinstate.Contract.ConnectionsEndpoints(&_Bitcoinstate.CallOpts, arg0)
}

// InitPool is a free data retrieval call binding the contract method 0xca38b326.
//
// Solidity: function initPool(bytes data_) pure returns()
func (_Bitcoinstate *BitcoinstateCaller) InitPool(opts *bind.CallOpts, data_ []byte) error {
	var out []interface{}
	err := _Bitcoinstate.contract.Call(opts, &out, "initPool", data_)

	if err != nil {
		return err
	}

	return err

}

// InitPool is a free data retrieval call binding the contract method 0xca38b326.
//
// Solidity: function initPool(bytes data_) pure returns()
func (_Bitcoinstate *BitcoinstateSession) InitPool(data_ []byte) error {
	return _Bitcoinstate.Contract.InitPool(&_Bitcoinstate.CallOpts, data_)
}

// InitPool is a free data retrieval call binding the contract method 0xca38b326.
//
// Solidity: function initPool(bytes data_) pure returns()
func (_Bitcoinstate *BitcoinstateCallerSession) InitPool(data_ []byte) error {
	return _Bitcoinstate.Contract.InitPool(&_Bitcoinstate.CallOpts, data_)
}

// NftOwners is a free data retrieval call binding the contract method 0xbbd94c2f.
//
// Solidity: function nftOwners(uint256 ) view returns(string)
func (_Bitcoinstate *BitcoinstateCaller) NftOwners(opts *bind.CallOpts, arg0 *big.Int) (string, error) {
	var out []interface{}
	err := _Bitcoinstate.contract.Call(opts, &out, "nftOwners", arg0)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// NftOwners is a free data retrieval call binding the contract method 0xbbd94c2f.
//
// Solidity: function nftOwners(uint256 ) view returns(string)
func (_Bitcoinstate *BitcoinstateSession) NftOwners(arg0 *big.Int) (string, error) {
	return _Bitcoinstate.Contract.NftOwners(&_Bitcoinstate.CallOpts, arg0)
}

// NftOwners is a free data retrieval call binding the contract method 0xbbd94c2f.
//
// Solidity: function nftOwners(uint256 ) view returns(string)
func (_Bitcoinstate *BitcoinstateCallerSession) NftOwners(arg0 *big.Int) (string, error) {
	return _Bitcoinstate.Contract.NftOwners(&_Bitcoinstate.CallOpts, arg0)
}

// NonFungibleManager is a free data retrieval call binding the contract method 0xecc28165.
//
// Solidity: function nonFungibleManager() view returns(address)
func (_Bitcoinstate *BitcoinstateCaller) NonFungibleManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Bitcoinstate.contract.Call(opts, &out, "nonFungibleManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// NonFungibleManager is a free data retrieval call binding the contract method 0xecc28165.
//
// Solidity: function nonFungibleManager() view returns(address)
func (_Bitcoinstate *BitcoinstateSession) NonFungibleManager() (common.Address, error) {
	return _Bitcoinstate.Contract.NonFungibleManager(&_Bitcoinstate.CallOpts)
}

// NonFungibleManager is a free data retrieval call binding the contract method 0xecc28165.
//
// Solidity: function nonFungibleManager() view returns(address)
func (_Bitcoinstate *BitcoinstateCallerSession) NonFungibleManager() (common.Address, error) {
	return _Bitcoinstate.Contract.NonFungibleManager(&_Bitcoinstate.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Bitcoinstate *BitcoinstateCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Bitcoinstate.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Bitcoinstate *BitcoinstateSession) Owner() (common.Address, error) {
	return _Bitcoinstate.Contract.Owner(&_Bitcoinstate.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Bitcoinstate *BitcoinstateCallerSession) Owner() (common.Address, error) {
	return _Bitcoinstate.Contract.Owner(&_Bitcoinstate.CallOpts)
}

// RemoveLiquidity is a free data retrieval call binding the contract method 0x028318cf.
//
// Solidity: function removeLiquidity(bytes data_) pure returns()
func (_Bitcoinstate *BitcoinstateCaller) RemoveLiquidity(opts *bind.CallOpts, data_ []byte) error {
	var out []interface{}
	err := _Bitcoinstate.contract.Call(opts, &out, "removeLiquidity", data_)

	if err != nil {
		return err
	}

	return err

}

// RemoveLiquidity is a free data retrieval call binding the contract method 0x028318cf.
//
// Solidity: function removeLiquidity(bytes data_) pure returns()
func (_Bitcoinstate *BitcoinstateSession) RemoveLiquidity(data_ []byte) error {
	return _Bitcoinstate.Contract.RemoveLiquidity(&_Bitcoinstate.CallOpts, data_)
}

// RemoveLiquidity is a free data retrieval call binding the contract method 0x028318cf.
//
// Solidity: function removeLiquidity(bytes data_) pure returns()
func (_Bitcoinstate *BitcoinstateCallerSession) RemoveLiquidity(data_ []byte) error {
	return _Bitcoinstate.Contract.RemoveLiquidity(&_Bitcoinstate.CallOpts, data_)
}

// RequestCount is a free data retrieval call binding the contract method 0x5badbe4c.
//
// Solidity: function requestCount() view returns(uint256)
func (_Bitcoinstate *BitcoinstateCaller) RequestCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Bitcoinstate.contract.Call(opts, &out, "requestCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RequestCount is a free data retrieval call binding the contract method 0x5badbe4c.
//
// Solidity: function requestCount() view returns(uint256)
func (_Bitcoinstate *BitcoinstateSession) RequestCount() (*big.Int, error) {
	return _Bitcoinstate.Contract.RequestCount(&_Bitcoinstate.CallOpts)
}

// RequestCount is a free data retrieval call binding the contract method 0x5badbe4c.
//
// Solidity: function requestCount() view returns(uint256)
func (_Bitcoinstate *BitcoinstateCallerSession) RequestCount() (*big.Int, error) {
	return _Bitcoinstate.Contract.RequestCount(&_Bitcoinstate.CallOpts)
}

// RouterV2 is a free data retrieval call binding the contract method 0x502f7446.
//
// Solidity: function routerV2() view returns(address)
func (_Bitcoinstate *BitcoinstateCaller) RouterV2(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Bitcoinstate.contract.Call(opts, &out, "routerV2")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RouterV2 is a free data retrieval call binding the contract method 0x502f7446.
//
// Solidity: function routerV2() view returns(address)
func (_Bitcoinstate *BitcoinstateSession) RouterV2() (common.Address, error) {
	return _Bitcoinstate.Contract.RouterV2(&_Bitcoinstate.CallOpts)
}

// RouterV2 is a free data retrieval call binding the contract method 0x502f7446.
//
// Solidity: function routerV2() view returns(address)
func (_Bitcoinstate *BitcoinstateCallerSession) RouterV2() (common.Address, error) {
	return _Bitcoinstate.Contract.RouterV2(&_Bitcoinstate.CallOpts)
}

// RuneFactory is a free data retrieval call binding the contract method 0x7d1f0a56.
//
// Solidity: function runeFactory() view returns(address)
func (_Bitcoinstate *BitcoinstateCaller) RuneFactory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Bitcoinstate.contract.Call(opts, &out, "runeFactory")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RuneFactory is a free data retrieval call binding the contract method 0x7d1f0a56.
//
// Solidity: function runeFactory() view returns(address)
func (_Bitcoinstate *BitcoinstateSession) RuneFactory() (common.Address, error) {
	return _Bitcoinstate.Contract.RuneFactory(&_Bitcoinstate.CallOpts)
}

// RuneFactory is a free data retrieval call binding the contract method 0x7d1f0a56.
//
// Solidity: function runeFactory() view returns(address)
func (_Bitcoinstate *BitcoinstateCallerSession) RuneFactory() (common.Address, error) {
	return _Bitcoinstate.Contract.RuneFactory(&_Bitcoinstate.CallOpts)
}

// Tokens is a free data retrieval call binding the contract method 0x904194a3.
//
// Solidity: function tokens(bytes32 ) view returns(address)
func (_Bitcoinstate *BitcoinstateCaller) Tokens(opts *bind.CallOpts, arg0 [32]byte) (common.Address, error) {
	var out []interface{}
	err := _Bitcoinstate.contract.Call(opts, &out, "tokens", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Tokens is a free data retrieval call binding the contract method 0x904194a3.
//
// Solidity: function tokens(bytes32 ) view returns(address)
func (_Bitcoinstate *BitcoinstateSession) Tokens(arg0 [32]byte) (common.Address, error) {
	return _Bitcoinstate.Contract.Tokens(&_Bitcoinstate.CallOpts, arg0)
}

// Tokens is a free data retrieval call binding the contract method 0x904194a3.
//
// Solidity: function tokens(bytes32 ) view returns(address)
func (_Bitcoinstate *BitcoinstateCallerSession) Tokens(arg0 [32]byte) (common.Address, error) {
	return _Bitcoinstate.Contract.Tokens(&_Bitcoinstate.CallOpts, arg0)
}

// XcallService is a free data retrieval call binding the contract method 0x7bf07164.
//
// Solidity: function xcallService() view returns(address)
func (_Bitcoinstate *BitcoinstateCaller) XcallService(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Bitcoinstate.contract.Call(opts, &out, "xcallService")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// XcallService is a free data retrieval call binding the contract method 0x7bf07164.
//
// Solidity: function xcallService() view returns(address)
func (_Bitcoinstate *BitcoinstateSession) XcallService() (common.Address, error) {
	return _Bitcoinstate.Contract.XcallService(&_Bitcoinstate.CallOpts)
}

// XcallService is a free data retrieval call binding the contract method 0x7bf07164.
//
// Solidity: function xcallService() view returns(address)
func (_Bitcoinstate *BitcoinstateCallerSession) XcallService() (common.Address, error) {
	return _Bitcoinstate.Contract.XcallService(&_Bitcoinstate.CallOpts)
}

// AddConnection is a paid mutator transaction binding the contract method 0x677dea1d.
//
// Solidity: function addConnection(address connection_) returns()
func (_Bitcoinstate *BitcoinstateTransactor) AddConnection(opts *bind.TransactOpts, connection_ common.Address) (*types.Transaction, error) {
	return _Bitcoinstate.contract.Transact(opts, "addConnection", connection_)
}

// AddConnection is a paid mutator transaction binding the contract method 0x677dea1d.
//
// Solidity: function addConnection(address connection_) returns()
func (_Bitcoinstate *BitcoinstateSession) AddConnection(connection_ common.Address) (*types.Transaction, error) {
	return _Bitcoinstate.Contract.AddConnection(&_Bitcoinstate.TransactOpts, connection_)
}

// AddConnection is a paid mutator transaction binding the contract method 0x677dea1d.
//
// Solidity: function addConnection(address connection_) returns()
func (_Bitcoinstate *BitcoinstateTransactorSession) AddConnection(connection_ common.Address) (*types.Transaction, error) {
	return _Bitcoinstate.Contract.AddConnection(&_Bitcoinstate.TransactOpts, connection_)
}

// HandleCallMessage is a paid mutator transaction binding the contract method 0x5d6a16f5.
//
// Solidity: function handleCallMessage(string _from, bytes _data, string[] _protocols) returns()
func (_Bitcoinstate *BitcoinstateTransactor) HandleCallMessage(opts *bind.TransactOpts, _from string, _data []byte, _protocols []string) (*types.Transaction, error) {
	return _Bitcoinstate.contract.Transact(opts, "handleCallMessage", _from, _data, _protocols)
}

// HandleCallMessage is a paid mutator transaction binding the contract method 0x5d6a16f5.
//
// Solidity: function handleCallMessage(string _from, bytes _data, string[] _protocols) returns()
func (_Bitcoinstate *BitcoinstateSession) HandleCallMessage(_from string, _data []byte, _protocols []string) (*types.Transaction, error) {
	return _Bitcoinstate.Contract.HandleCallMessage(&_Bitcoinstate.TransactOpts, _from, _data, _protocols)
}

// HandleCallMessage is a paid mutator transaction binding the contract method 0x5d6a16f5.
//
// Solidity: function handleCallMessage(string _from, bytes _data, string[] _protocols) returns()
func (_Bitcoinstate *BitcoinstateTransactorSession) HandleCallMessage(_from string, _data []byte, _protocols []string) (*types.Transaction, error) {
	return _Bitcoinstate.Contract.HandleCallMessage(&_Bitcoinstate.TransactOpts, _from, _data, _protocols)
}

// Initialize is a paid mutator transaction binding the contract method 0xe6bfbfd8.
//
// Solidity: function initialize(address xcall_, address uinswapV3Router_, address nonfungiblePositionManager_, address[] connections) returns()
func (_Bitcoinstate *BitcoinstateTransactor) Initialize(opts *bind.TransactOpts, xcall_ common.Address, uinswapV3Router_ common.Address, nonfungiblePositionManager_ common.Address, connections []common.Address) (*types.Transaction, error) {
	return _Bitcoinstate.contract.Transact(opts, "initialize", xcall_, uinswapV3Router_, nonfungiblePositionManager_, connections)
}

// Initialize is a paid mutator transaction binding the contract method 0xe6bfbfd8.
//
// Solidity: function initialize(address xcall_, address uinswapV3Router_, address nonfungiblePositionManager_, address[] connections) returns()
func (_Bitcoinstate *BitcoinstateSession) Initialize(xcall_ common.Address, uinswapV3Router_ common.Address, nonfungiblePositionManager_ common.Address, connections []common.Address) (*types.Transaction, error) {
	return _Bitcoinstate.Contract.Initialize(&_Bitcoinstate.TransactOpts, xcall_, uinswapV3Router_, nonfungiblePositionManager_, connections)
}

// Initialize is a paid mutator transaction binding the contract method 0xe6bfbfd8.
//
// Solidity: function initialize(address xcall_, address uinswapV3Router_, address nonfungiblePositionManager_, address[] connections) returns()
func (_Bitcoinstate *BitcoinstateTransactorSession) Initialize(xcall_ common.Address, uinswapV3Router_ common.Address, nonfungiblePositionManager_ common.Address, connections []common.Address) (*types.Transaction, error) {
	return _Bitcoinstate.Contract.Initialize(&_Bitcoinstate.TransactOpts, xcall_, uinswapV3Router_, nonfungiblePositionManager_, connections)
}

// Migrate is a paid mutator transaction binding the contract method 0x53f954ea.
//
// Solidity: function migrate(string from_, bytes _data) returns()
func (_Bitcoinstate *BitcoinstateTransactor) Migrate(opts *bind.TransactOpts, from_ string, _data []byte) (*types.Transaction, error) {
	return _Bitcoinstate.contract.Transact(opts, "migrate", from_, _data)
}

// Migrate is a paid mutator transaction binding the contract method 0x53f954ea.
//
// Solidity: function migrate(string from_, bytes _data) returns()
func (_Bitcoinstate *BitcoinstateSession) Migrate(from_ string, _data []byte) (*types.Transaction, error) {
	return _Bitcoinstate.Contract.Migrate(&_Bitcoinstate.TransactOpts, from_, _data)
}

// Migrate is a paid mutator transaction binding the contract method 0x53f954ea.
//
// Solidity: function migrate(string from_, bytes _data) returns()
func (_Bitcoinstate *BitcoinstateTransactorSession) Migrate(from_ string, _data []byte) (*types.Transaction, error) {
	return _Bitcoinstate.Contract.Migrate(&_Bitcoinstate.TransactOpts, from_, _data)
}

// MigrateComplete is a paid mutator transaction binding the contract method 0xf0ad3762.
//
// Solidity: function migrateComplete() returns()
func (_Bitcoinstate *BitcoinstateTransactor) MigrateComplete(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bitcoinstate.contract.Transact(opts, "migrateComplete")
}

// MigrateComplete is a paid mutator transaction binding the contract method 0xf0ad3762.
//
// Solidity: function migrateComplete() returns()
func (_Bitcoinstate *BitcoinstateSession) MigrateComplete() (*types.Transaction, error) {
	return _Bitcoinstate.Contract.MigrateComplete(&_Bitcoinstate.TransactOpts)
}

// MigrateComplete is a paid mutator transaction binding the contract method 0xf0ad3762.
//
// Solidity: function migrateComplete() returns()
func (_Bitcoinstate *BitcoinstateTransactorSession) MigrateComplete() (*types.Transaction, error) {
	return _Bitcoinstate.Contract.MigrateComplete(&_Bitcoinstate.TransactOpts)
}

// RemoveConnection is a paid mutator transaction binding the contract method 0x65301f0d.
//
// Solidity: function removeConnection(address connection_) returns()
func (_Bitcoinstate *BitcoinstateTransactor) RemoveConnection(opts *bind.TransactOpts, connection_ common.Address) (*types.Transaction, error) {
	return _Bitcoinstate.contract.Transact(opts, "removeConnection", connection_)
}

// RemoveConnection is a paid mutator transaction binding the contract method 0x65301f0d.
//
// Solidity: function removeConnection(address connection_) returns()
func (_Bitcoinstate *BitcoinstateSession) RemoveConnection(connection_ common.Address) (*types.Transaction, error) {
	return _Bitcoinstate.Contract.RemoveConnection(&_Bitcoinstate.TransactOpts, connection_)
}

// RemoveConnection is a paid mutator transaction binding the contract method 0x65301f0d.
//
// Solidity: function removeConnection(address connection_) returns()
func (_Bitcoinstate *BitcoinstateTransactorSession) RemoveConnection(connection_ common.Address) (*types.Transaction, error) {
	return _Bitcoinstate.Contract.RemoveConnection(&_Bitcoinstate.TransactOpts, connection_)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Bitcoinstate *BitcoinstateTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bitcoinstate.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Bitcoinstate *BitcoinstateSession) RenounceOwnership() (*types.Transaction, error) {
	return _Bitcoinstate.Contract.RenounceOwnership(&_Bitcoinstate.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Bitcoinstate *BitcoinstateTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Bitcoinstate.Contract.RenounceOwnership(&_Bitcoinstate.TransactOpts)
}

// SetRuneFactory is a paid mutator transaction binding the contract method 0xa7fb534c.
//
// Solidity: function setRuneFactory(address _runeFactory) returns()
func (_Bitcoinstate *BitcoinstateTransactor) SetRuneFactory(opts *bind.TransactOpts, _runeFactory common.Address) (*types.Transaction, error) {
	return _Bitcoinstate.contract.Transact(opts, "setRuneFactory", _runeFactory)
}

// SetRuneFactory is a paid mutator transaction binding the contract method 0xa7fb534c.
//
// Solidity: function setRuneFactory(address _runeFactory) returns()
func (_Bitcoinstate *BitcoinstateSession) SetRuneFactory(_runeFactory common.Address) (*types.Transaction, error) {
	return _Bitcoinstate.Contract.SetRuneFactory(&_Bitcoinstate.TransactOpts, _runeFactory)
}

// SetRuneFactory is a paid mutator transaction binding the contract method 0xa7fb534c.
//
// Solidity: function setRuneFactory(address _runeFactory) returns()
func (_Bitcoinstate *BitcoinstateTransactorSession) SetRuneFactory(_runeFactory common.Address) (*types.Transaction, error) {
	return _Bitcoinstate.Contract.SetRuneFactory(&_Bitcoinstate.TransactOpts, _runeFactory)
}

// SimulateRequest is a paid mutator transaction binding the contract method 0x4e2eeb94.
//
// Solidity: function simulateRequest(uint256 sequenceNumber, string _from, bytes parseData) returns()
func (_Bitcoinstate *BitcoinstateTransactor) SimulateRequest(opts *bind.TransactOpts, sequenceNumber *big.Int, _from string, parseData []byte) (*types.Transaction, error) {
	return _Bitcoinstate.contract.Transact(opts, "simulateRequest", sequenceNumber, _from, parseData)
}

// SimulateRequest is a paid mutator transaction binding the contract method 0x4e2eeb94.
//
// Solidity: function simulateRequest(uint256 sequenceNumber, string _from, bytes parseData) returns()
func (_Bitcoinstate *BitcoinstateSession) SimulateRequest(sequenceNumber *big.Int, _from string, parseData []byte) (*types.Transaction, error) {
	return _Bitcoinstate.Contract.SimulateRequest(&_Bitcoinstate.TransactOpts, sequenceNumber, _from, parseData)
}

// SimulateRequest is a paid mutator transaction binding the contract method 0x4e2eeb94.
//
// Solidity: function simulateRequest(uint256 sequenceNumber, string _from, bytes parseData) returns()
func (_Bitcoinstate *BitcoinstateTransactorSession) SimulateRequest(sequenceNumber *big.Int, _from string, parseData []byte) (*types.Transaction, error) {
	return _Bitcoinstate.Contract.SimulateRequest(&_Bitcoinstate.TransactOpts, sequenceNumber, _from, parseData)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Bitcoinstate *BitcoinstateTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Bitcoinstate.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Bitcoinstate *BitcoinstateSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Bitcoinstate.Contract.TransferOwnership(&_Bitcoinstate.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Bitcoinstate *BitcoinstateTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Bitcoinstate.Contract.TransferOwnership(&_Bitcoinstate.TransactOpts, newOwner)
}

// WithdrawTo is a paid mutator transaction binding the contract method 0x43aeb854.
//
// Solidity: function withdrawTo(address token_, string receiver_, uint256 amount) returns()
func (_Bitcoinstate *BitcoinstateTransactor) WithdrawTo(opts *bind.TransactOpts, token_ common.Address, receiver_ string, amount *big.Int) (*types.Transaction, error) {
	return _Bitcoinstate.contract.Transact(opts, "withdrawTo", token_, receiver_, amount)
}

// WithdrawTo is a paid mutator transaction binding the contract method 0x43aeb854.
//
// Solidity: function withdrawTo(address token_, string receiver_, uint256 amount) returns()
func (_Bitcoinstate *BitcoinstateSession) WithdrawTo(token_ common.Address, receiver_ string, amount *big.Int) (*types.Transaction, error) {
	return _Bitcoinstate.Contract.WithdrawTo(&_Bitcoinstate.TransactOpts, token_, receiver_, amount)
}

// WithdrawTo is a paid mutator transaction binding the contract method 0x43aeb854.
//
// Solidity: function withdrawTo(address token_, string receiver_, uint256 amount) returns()
func (_Bitcoinstate *BitcoinstateTransactorSession) WithdrawTo(token_ common.Address, receiver_ string, amount *big.Int) (*types.Transaction, error) {
	return _Bitcoinstate.Contract.WithdrawTo(&_Bitcoinstate.TransactOpts, token_, receiver_, amount)
}

// BitcoinstateAddConnectionIterator is returned from FilterAddConnection and is used to iterate over the raw logs and unpacked data for AddConnection events raised by the Bitcoinstate contract.
type BitcoinstateAddConnectionIterator struct {
	Event *BitcoinstateAddConnection // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BitcoinstateAddConnectionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BitcoinstateAddConnection)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BitcoinstateAddConnection)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BitcoinstateAddConnectionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BitcoinstateAddConnectionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BitcoinstateAddConnection represents a AddConnection event raised by the Bitcoinstate contract.
type BitcoinstateAddConnection struct {
	Connection common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterAddConnection is a free log retrieval operation binding the contract event 0x5cd3d8da8ef00a8b5228348fe7683dae605751d4867ba7ea9fdc8260b9e2c7d3.
//
// Solidity: event AddConnection(address connection_)
func (_Bitcoinstate *BitcoinstateFilterer) FilterAddConnection(opts *bind.FilterOpts) (*BitcoinstateAddConnectionIterator, error) {

	logs, sub, err := _Bitcoinstate.contract.FilterLogs(opts, "AddConnection")
	if err != nil {
		return nil, err
	}
	return &BitcoinstateAddConnectionIterator{contract: _Bitcoinstate.contract, event: "AddConnection", logs: logs, sub: sub}, nil
}

// WatchAddConnection is a free log subscription operation binding the contract event 0x5cd3d8da8ef00a8b5228348fe7683dae605751d4867ba7ea9fdc8260b9e2c7d3.
//
// Solidity: event AddConnection(address connection_)
func (_Bitcoinstate *BitcoinstateFilterer) WatchAddConnection(opts *bind.WatchOpts, sink chan<- *BitcoinstateAddConnection) (event.Subscription, error) {

	logs, sub, err := _Bitcoinstate.contract.WatchLogs(opts, "AddConnection")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BitcoinstateAddConnection)
				if err := _Bitcoinstate.contract.UnpackLog(event, "AddConnection", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAddConnection is a log parse operation binding the contract event 0x5cd3d8da8ef00a8b5228348fe7683dae605751d4867ba7ea9fdc8260b9e2c7d3.
//
// Solidity: event AddConnection(address connection_)
func (_Bitcoinstate *BitcoinstateFilterer) ParseAddConnection(log types.Log) (*BitcoinstateAddConnection, error) {
	event := new(BitcoinstateAddConnection)
	if err := _Bitcoinstate.contract.UnpackLog(event, "AddConnection", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BitcoinstateInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Bitcoinstate contract.
type BitcoinstateInitializedIterator struct {
	Event *BitcoinstateInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BitcoinstateInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BitcoinstateInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BitcoinstateInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BitcoinstateInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BitcoinstateInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BitcoinstateInitialized represents a Initialized event raised by the Bitcoinstate contract.
type BitcoinstateInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Bitcoinstate *BitcoinstateFilterer) FilterInitialized(opts *bind.FilterOpts) (*BitcoinstateInitializedIterator, error) {

	logs, sub, err := _Bitcoinstate.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &BitcoinstateInitializedIterator{contract: _Bitcoinstate.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Bitcoinstate *BitcoinstateFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *BitcoinstateInitialized) (event.Subscription, error) {

	logs, sub, err := _Bitcoinstate.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BitcoinstateInitialized)
				if err := _Bitcoinstate.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Bitcoinstate *BitcoinstateFilterer) ParseInitialized(log types.Log) (*BitcoinstateInitialized, error) {
	event := new(BitcoinstateInitialized)
	if err := _Bitcoinstate.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BitcoinstateOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Bitcoinstate contract.
type BitcoinstateOwnershipTransferredIterator struct {
	Event *BitcoinstateOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BitcoinstateOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BitcoinstateOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BitcoinstateOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BitcoinstateOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BitcoinstateOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BitcoinstateOwnershipTransferred represents a OwnershipTransferred event raised by the Bitcoinstate contract.
type BitcoinstateOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Bitcoinstate *BitcoinstateFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*BitcoinstateOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Bitcoinstate.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &BitcoinstateOwnershipTransferredIterator{contract: _Bitcoinstate.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Bitcoinstate *BitcoinstateFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BitcoinstateOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Bitcoinstate.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BitcoinstateOwnershipTransferred)
				if err := _Bitcoinstate.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Bitcoinstate *BitcoinstateFilterer) ParseOwnershipTransferred(log types.Log) (*BitcoinstateOwnershipTransferred, error) {
	event := new(BitcoinstateOwnershipTransferred)
	if err := _Bitcoinstate.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BitcoinstateRemoveConnectionIterator is returned from FilterRemoveConnection and is used to iterate over the raw logs and unpacked data for RemoveConnection events raised by the Bitcoinstate contract.
type BitcoinstateRemoveConnectionIterator struct {
	Event *BitcoinstateRemoveConnection // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BitcoinstateRemoveConnectionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BitcoinstateRemoveConnection)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BitcoinstateRemoveConnection)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BitcoinstateRemoveConnectionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BitcoinstateRemoveConnectionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BitcoinstateRemoveConnection represents a RemoveConnection event raised by the Bitcoinstate contract.
type BitcoinstateRemoveConnection struct {
	Connection common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterRemoveConnection is a free log retrieval operation binding the contract event 0xace8d11a44b7aa536cc46a77b519166c001adc485ba8cfa404e1aa252b07db38.
//
// Solidity: event RemoveConnection(address connection_)
func (_Bitcoinstate *BitcoinstateFilterer) FilterRemoveConnection(opts *bind.FilterOpts) (*BitcoinstateRemoveConnectionIterator, error) {

	logs, sub, err := _Bitcoinstate.contract.FilterLogs(opts, "RemoveConnection")
	if err != nil {
		return nil, err
	}
	return &BitcoinstateRemoveConnectionIterator{contract: _Bitcoinstate.contract, event: "RemoveConnection", logs: logs, sub: sub}, nil
}

// WatchRemoveConnection is a free log subscription operation binding the contract event 0xace8d11a44b7aa536cc46a77b519166c001adc485ba8cfa404e1aa252b07db38.
//
// Solidity: event RemoveConnection(address connection_)
func (_Bitcoinstate *BitcoinstateFilterer) WatchRemoveConnection(opts *bind.WatchOpts, sink chan<- *BitcoinstateRemoveConnection) (event.Subscription, error) {

	logs, sub, err := _Bitcoinstate.contract.WatchLogs(opts, "RemoveConnection")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BitcoinstateRemoveConnection)
				if err := _Bitcoinstate.contract.UnpackLog(event, "RemoveConnection", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRemoveConnection is a log parse operation binding the contract event 0xace8d11a44b7aa536cc46a77b519166c001adc485ba8cfa404e1aa252b07db38.
//
// Solidity: event RemoveConnection(address connection_)
func (_Bitcoinstate *BitcoinstateFilterer) ParseRemoveConnection(log types.Log) (*BitcoinstateRemoveConnection, error) {
	event := new(BitcoinstateRemoveConnection)
	if err := _Bitcoinstate.contract.UnpackLog(event, "RemoveConnection", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BitcoinstateRequestExecutedIterator is returned from FilterRequestExecuted and is used to iterate over the raw logs and unpacked data for RequestExecuted events raised by the Bitcoinstate contract.
type BitcoinstateRequestExecutedIterator struct {
	Event *BitcoinstateRequestExecuted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BitcoinstateRequestExecutedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BitcoinstateRequestExecuted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BitcoinstateRequestExecuted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BitcoinstateRequestExecutedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BitcoinstateRequestExecutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BitcoinstateRequestExecuted represents a RequestExecuted event raised by the Bitcoinstate contract.
type BitcoinstateRequestExecuted struct {
	Id        *big.Int
	StateRoot [32]byte
	Data      []byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRequestExecuted is a free log retrieval operation binding the contract event 0x9c343316d67a8e28446ef883ab491ece3ff70d3eeaa9fbd13a362a0afd690721.
//
// Solidity: event RequestExecuted(uint256 id, bytes32 stateRoot, bytes data)
func (_Bitcoinstate *BitcoinstateFilterer) FilterRequestExecuted(opts *bind.FilterOpts) (*BitcoinstateRequestExecutedIterator, error) {

	logs, sub, err := _Bitcoinstate.contract.FilterLogs(opts, "RequestExecuted")
	if err != nil {
		return nil, err
	}
	return &BitcoinstateRequestExecutedIterator{contract: _Bitcoinstate.contract, event: "RequestExecuted", logs: logs, sub: sub}, nil
}

// WatchRequestExecuted is a free log subscription operation binding the contract event 0x9c343316d67a8e28446ef883ab491ece3ff70d3eeaa9fbd13a362a0afd690721.
//
// Solidity: event RequestExecuted(uint256 id, bytes32 stateRoot, bytes data)
func (_Bitcoinstate *BitcoinstateFilterer) WatchRequestExecuted(opts *bind.WatchOpts, sink chan<- *BitcoinstateRequestExecuted) (event.Subscription, error) {

	logs, sub, err := _Bitcoinstate.contract.WatchLogs(opts, "RequestExecuted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BitcoinstateRequestExecuted)
				if err := _Bitcoinstate.contract.UnpackLog(event, "RequestExecuted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRequestExecuted is a log parse operation binding the contract event 0x9c343316d67a8e28446ef883ab491ece3ff70d3eeaa9fbd13a362a0afd690721.
//
// Solidity: event RequestExecuted(uint256 id, bytes32 stateRoot, bytes data)
func (_Bitcoinstate *BitcoinstateFilterer) ParseRequestExecuted(log types.Log) (*BitcoinstateRequestExecuted, error) {
	event := new(BitcoinstateRequestExecuted)
	if err := _Bitcoinstate.contract.UnpackLog(event, "RequestExecuted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
