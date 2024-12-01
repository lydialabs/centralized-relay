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

// BitcoinStateMetaData contains all meta data concerning the BitcoinState contract.
var BitcoinStateMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"accountBalances\",\"inputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"addConnection\",\"inputs\":[{\"name\":\"connection_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"bitcoinNid\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"claimTokens\",\"inputs\":[{\"name\":\"token0\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"token1\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"connections\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"connectionsEndpoints\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"handleCallMessage\",\"inputs\":[{\"name\":\"_from\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_data\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"_protocols\",\"type\":\"string[]\",\"internalType\":\"string[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initPool\",\"inputs\":[{\"name\":\"data_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"xcall_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"uinswapV3Router_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"nonfungiblePositionManager_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"connections\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"migrate\",\"inputs\":[{\"name\":\"from_\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"migrateComplete\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"nftOwners\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"nonFungibleManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"removeConnection\",\"inputs\":[{\"name\":\"connection_\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"removeLiquidity\",\"inputs\":[{\"name\":\"data_\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"requestCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"routerV2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"runeFactory\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setRuneFactory\",\"inputs\":[{\"name\":\"_runeFactory\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"tokens\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawTo\",\"inputs\":[{\"name\":\"token_\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"receiver_\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"xcallService\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"AddConnection\",\"inputs\":[{\"name\":\"connection_\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RemoveConnection\",\"inputs\":[{\"name\":\"connection_\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RequestExecuted\",\"inputs\":[{\"name\":\"id\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"stateRoot\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"data\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false}]",
}

// BitcoinStateABI is the input ABI used to generate the binding from.
// Deprecated: Use BitcoinStateMetaData.ABI instead.
var BitcoinStateABI = BitcoinStateMetaData.ABI

// BitcoinState is an auto generated Go binding around an Ethereum contract.
type BitcoinState struct {
	BitcoinStateCaller     // Read-only binding to the contract
	BitcoinStateTransactor // Write-only binding to the contract
	BitcoinStateFilterer   // Log filterer for contract events
}

// BitcoinStateCaller is an auto generated read-only Go binding around an Ethereum contract.
type BitcoinStateCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BitcoinStateTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BitcoinStateTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BitcoinStateFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BitcoinStateFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BitcoinStateSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BitcoinStateSession struct {
	Contract     *BitcoinState     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BitcoinStateCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BitcoinStateCallerSession struct {
	Contract *BitcoinStateCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// BitcoinStateTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BitcoinStateTransactorSession struct {
	Contract     *BitcoinStateTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// BitcoinStateRaw is an auto generated low-level Go binding around an Ethereum contract.
type BitcoinStateRaw struct {
	Contract *BitcoinState // Generic contract binding to access the raw methods on
}

// BitcoinStateCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BitcoinStateCallerRaw struct {
	Contract *BitcoinStateCaller // Generic read-only contract binding to access the raw methods on
}

// BitcoinStateTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BitcoinStateTransactorRaw struct {
	Contract *BitcoinStateTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBitcoinState creates a new instance of BitcoinState, bound to a specific deployed contract.
func NewBitcoinState(address common.Address, backend bind.ContractBackend) (*BitcoinState, error) {
	contract, err := bindBitcoinState(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BitcoinState{BitcoinStateCaller: BitcoinStateCaller{contract: contract}, BitcoinStateTransactor: BitcoinStateTransactor{contract: contract}, BitcoinStateFilterer: BitcoinStateFilterer{contract: contract}}, nil
}

// NewBitcoinStateCaller creates a new read-only instance of BitcoinState, bound to a specific deployed contract.
func NewBitcoinStateCaller(address common.Address, caller bind.ContractCaller) (*BitcoinStateCaller, error) {
	contract, err := bindBitcoinState(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BitcoinStateCaller{contract: contract}, nil
}

// NewBitcoinStateTransactor creates a new write-only instance of BitcoinState, bound to a specific deployed contract.
func NewBitcoinStateTransactor(address common.Address, transactor bind.ContractTransactor) (*BitcoinStateTransactor, error) {
	contract, err := bindBitcoinState(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BitcoinStateTransactor{contract: contract}, nil
}

// NewBitcoinStateFilterer creates a new log filterer instance of BitcoinState, bound to a specific deployed contract.
func NewBitcoinStateFilterer(address common.Address, filterer bind.ContractFilterer) (*BitcoinStateFilterer, error) {
	contract, err := bindBitcoinState(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BitcoinStateFilterer{contract: contract}, nil
}

// bindBitcoinState binds a generic wrapper to an already deployed contract.
func bindBitcoinState(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BitcoinStateMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BitcoinState *BitcoinStateRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BitcoinState.Contract.BitcoinStateCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BitcoinState *BitcoinStateRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BitcoinState.Contract.BitcoinStateTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BitcoinState *BitcoinStateRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BitcoinState.Contract.BitcoinStateTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BitcoinState *BitcoinStateCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BitcoinState.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BitcoinState *BitcoinStateTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BitcoinState.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BitcoinState *BitcoinStateTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BitcoinState.Contract.contract.Transact(opts, method, params...)
}

// AccountBalances is a free data retrieval call binding the contract method 0x69477464.
//
// Solidity: function accountBalances(string , address ) view returns(uint256)
func (_BitcoinState *BitcoinStateCaller) AccountBalances(opts *bind.CallOpts, arg0 string, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _BitcoinState.contract.Call(opts, &out, "accountBalances", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AccountBalances is a free data retrieval call binding the contract method 0x69477464.
//
// Solidity: function accountBalances(string , address ) view returns(uint256)
func (_BitcoinState *BitcoinStateSession) AccountBalances(arg0 string, arg1 common.Address) (*big.Int, error) {
	return _BitcoinState.Contract.AccountBalances(&_BitcoinState.CallOpts, arg0, arg1)
}

// AccountBalances is a free data retrieval call binding the contract method 0x69477464.
//
// Solidity: function accountBalances(string , address ) view returns(uint256)
func (_BitcoinState *BitcoinStateCallerSession) AccountBalances(arg0 string, arg1 common.Address) (*big.Int, error) {
	return _BitcoinState.Contract.AccountBalances(&_BitcoinState.CallOpts, arg0, arg1)
}

// BitcoinNid is a free data retrieval call binding the contract method 0xf3cc4515.
//
// Solidity: function bitcoinNid() view returns(string)
func (_BitcoinState *BitcoinStateCaller) BitcoinNid(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _BitcoinState.contract.Call(opts, &out, "bitcoinNid")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// BitcoinNid is a free data retrieval call binding the contract method 0xf3cc4515.
//
// Solidity: function bitcoinNid() view returns(string)
func (_BitcoinState *BitcoinStateSession) BitcoinNid() (string, error) {
	return _BitcoinState.Contract.BitcoinNid(&_BitcoinState.CallOpts)
}

// BitcoinNid is a free data retrieval call binding the contract method 0xf3cc4515.
//
// Solidity: function bitcoinNid() view returns(string)
func (_BitcoinState *BitcoinStateCallerSession) BitcoinNid() (string, error) {
	return _BitcoinState.Contract.BitcoinNid(&_BitcoinState.CallOpts)
}

// ClaimTokens is a free data retrieval call binding the contract method 0x69ffa08a.
//
// Solidity: function claimTokens(address token0, address token1) pure returns()
func (_BitcoinState *BitcoinStateCaller) ClaimTokens(opts *bind.CallOpts, token0 common.Address, token1 common.Address) error {
	var out []interface{}
	err := _BitcoinState.contract.Call(opts, &out, "claimTokens", token0, token1)

	if err != nil {
		return err
	}

	return err

}

// ClaimTokens is a free data retrieval call binding the contract method 0x69ffa08a.
//
// Solidity: function claimTokens(address token0, address token1) pure returns()
func (_BitcoinState *BitcoinStateSession) ClaimTokens(token0 common.Address, token1 common.Address) error {
	return _BitcoinState.Contract.ClaimTokens(&_BitcoinState.CallOpts, token0, token1)
}

// ClaimTokens is a free data retrieval call binding the contract method 0x69ffa08a.
//
// Solidity: function claimTokens(address token0, address token1) pure returns()
func (_BitcoinState *BitcoinStateCallerSession) ClaimTokens(token0 common.Address, token1 common.Address) error {
	return _BitcoinState.Contract.ClaimTokens(&_BitcoinState.CallOpts, token0, token1)
}

// Connections is a free data retrieval call binding the contract method 0xc0896578.
//
// Solidity: function connections(address ) view returns(bool)
func (_BitcoinState *BitcoinStateCaller) Connections(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _BitcoinState.contract.Call(opts, &out, "connections", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Connections is a free data retrieval call binding the contract method 0xc0896578.
//
// Solidity: function connections(address ) view returns(bool)
func (_BitcoinState *BitcoinStateSession) Connections(arg0 common.Address) (bool, error) {
	return _BitcoinState.Contract.Connections(&_BitcoinState.CallOpts, arg0)
}

// Connections is a free data retrieval call binding the contract method 0xc0896578.
//
// Solidity: function connections(address ) view returns(bool)
func (_BitcoinState *BitcoinStateCallerSession) Connections(arg0 common.Address) (bool, error) {
	return _BitcoinState.Contract.Connections(&_BitcoinState.CallOpts, arg0)
}

// ConnectionsEndpoints is a free data retrieval call binding the contract method 0xb687db2e.
//
// Solidity: function connectionsEndpoints(uint256 ) view returns(address)
func (_BitcoinState *BitcoinStateCaller) ConnectionsEndpoints(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _BitcoinState.contract.Call(opts, &out, "connectionsEndpoints", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ConnectionsEndpoints is a free data retrieval call binding the contract method 0xb687db2e.
//
// Solidity: function connectionsEndpoints(uint256 ) view returns(address)
func (_BitcoinState *BitcoinStateSession) ConnectionsEndpoints(arg0 *big.Int) (common.Address, error) {
	return _BitcoinState.Contract.ConnectionsEndpoints(&_BitcoinState.CallOpts, arg0)
}

// ConnectionsEndpoints is a free data retrieval call binding the contract method 0xb687db2e.
//
// Solidity: function connectionsEndpoints(uint256 ) view returns(address)
func (_BitcoinState *BitcoinStateCallerSession) ConnectionsEndpoints(arg0 *big.Int) (common.Address, error) {
	return _BitcoinState.Contract.ConnectionsEndpoints(&_BitcoinState.CallOpts, arg0)
}

// InitPool is a free data retrieval call binding the contract method 0xca38b326.
//
// Solidity: function initPool(bytes data_) pure returns()
func (_BitcoinState *BitcoinStateCaller) InitPool(opts *bind.CallOpts, data_ []byte) error {
	var out []interface{}
	err := _BitcoinState.contract.Call(opts, &out, "initPool", data_)

	if err != nil {
		return err
	}

	return err

}

// InitPool is a free data retrieval call binding the contract method 0xca38b326.
//
// Solidity: function initPool(bytes data_) pure returns()
func (_BitcoinState *BitcoinStateSession) InitPool(data_ []byte) error {
	return _BitcoinState.Contract.InitPool(&_BitcoinState.CallOpts, data_)
}

// InitPool is a free data retrieval call binding the contract method 0xca38b326.
//
// Solidity: function initPool(bytes data_) pure returns()
func (_BitcoinState *BitcoinStateCallerSession) InitPool(data_ []byte) error {
	return _BitcoinState.Contract.InitPool(&_BitcoinState.CallOpts, data_)
}

// NftOwners is a free data retrieval call binding the contract method 0xbbd94c2f.
//
// Solidity: function nftOwners(uint256 ) view returns(string)
func (_BitcoinState *BitcoinStateCaller) NftOwners(opts *bind.CallOpts, arg0 *big.Int) (string, error) {
	var out []interface{}
	err := _BitcoinState.contract.Call(opts, &out, "nftOwners", arg0)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// NftOwners is a free data retrieval call binding the contract method 0xbbd94c2f.
//
// Solidity: function nftOwners(uint256 ) view returns(string)
func (_BitcoinState *BitcoinStateSession) NftOwners(arg0 *big.Int) (string, error) {
	return _BitcoinState.Contract.NftOwners(&_BitcoinState.CallOpts, arg0)
}

// NftOwners is a free data retrieval call binding the contract method 0xbbd94c2f.
//
// Solidity: function nftOwners(uint256 ) view returns(string)
func (_BitcoinState *BitcoinStateCallerSession) NftOwners(arg0 *big.Int) (string, error) {
	return _BitcoinState.Contract.NftOwners(&_BitcoinState.CallOpts, arg0)
}

// NonFungibleManager is a free data retrieval call binding the contract method 0xecc28165.
//
// Solidity: function nonFungibleManager() view returns(address)
func (_BitcoinState *BitcoinStateCaller) NonFungibleManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BitcoinState.contract.Call(opts, &out, "nonFungibleManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// NonFungibleManager is a free data retrieval call binding the contract method 0xecc28165.
//
// Solidity: function nonFungibleManager() view returns(address)
func (_BitcoinState *BitcoinStateSession) NonFungibleManager() (common.Address, error) {
	return _BitcoinState.Contract.NonFungibleManager(&_BitcoinState.CallOpts)
}

// NonFungibleManager is a free data retrieval call binding the contract method 0xecc28165.
//
// Solidity: function nonFungibleManager() view returns(address)
func (_BitcoinState *BitcoinStateCallerSession) NonFungibleManager() (common.Address, error) {
	return _BitcoinState.Contract.NonFungibleManager(&_BitcoinState.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_BitcoinState *BitcoinStateCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BitcoinState.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_BitcoinState *BitcoinStateSession) Owner() (common.Address, error) {
	return _BitcoinState.Contract.Owner(&_BitcoinState.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_BitcoinState *BitcoinStateCallerSession) Owner() (common.Address, error) {
	return _BitcoinState.Contract.Owner(&_BitcoinState.CallOpts)
}

// RemoveLiquidity is a free data retrieval call binding the contract method 0x028318cf.
//
// Solidity: function removeLiquidity(bytes data_) pure returns()
func (_BitcoinState *BitcoinStateCaller) RemoveLiquidity(opts *bind.CallOpts, data_ []byte) error {
	var out []interface{}
	err := _BitcoinState.contract.Call(opts, &out, "removeLiquidity", data_)

	if err != nil {
		return err
	}

	return err

}

// RemoveLiquidity is a free data retrieval call binding the contract method 0x028318cf.
//
// Solidity: function removeLiquidity(bytes data_) pure returns()
func (_BitcoinState *BitcoinStateSession) RemoveLiquidity(data_ []byte) error {
	return _BitcoinState.Contract.RemoveLiquidity(&_BitcoinState.CallOpts, data_)
}

// RemoveLiquidity is a free data retrieval call binding the contract method 0x028318cf.
//
// Solidity: function removeLiquidity(bytes data_) pure returns()
func (_BitcoinState *BitcoinStateCallerSession) RemoveLiquidity(data_ []byte) error {
	return _BitcoinState.Contract.RemoveLiquidity(&_BitcoinState.CallOpts, data_)
}

// RequestCount is a free data retrieval call binding the contract method 0x5badbe4c.
//
// Solidity: function requestCount() view returns(uint256)
func (_BitcoinState *BitcoinStateCaller) RequestCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BitcoinState.contract.Call(opts, &out, "requestCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RequestCount is a free data retrieval call binding the contract method 0x5badbe4c.
//
// Solidity: function requestCount() view returns(uint256)
func (_BitcoinState *BitcoinStateSession) RequestCount() (*big.Int, error) {
	return _BitcoinState.Contract.RequestCount(&_BitcoinState.CallOpts)
}

// RequestCount is a free data retrieval call binding the contract method 0x5badbe4c.
//
// Solidity: function requestCount() view returns(uint256)
func (_BitcoinState *BitcoinStateCallerSession) RequestCount() (*big.Int, error) {
	return _BitcoinState.Contract.RequestCount(&_BitcoinState.CallOpts)
}

// RouterV2 is a free data retrieval call binding the contract method 0x502f7446.
//
// Solidity: function routerV2() view returns(address)
func (_BitcoinState *BitcoinStateCaller) RouterV2(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BitcoinState.contract.Call(opts, &out, "routerV2")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RouterV2 is a free data retrieval call binding the contract method 0x502f7446.
//
// Solidity: function routerV2() view returns(address)
func (_BitcoinState *BitcoinStateSession) RouterV2() (common.Address, error) {
	return _BitcoinState.Contract.RouterV2(&_BitcoinState.CallOpts)
}

// RouterV2 is a free data retrieval call binding the contract method 0x502f7446.
//
// Solidity: function routerV2() view returns(address)
func (_BitcoinState *BitcoinStateCallerSession) RouterV2() (common.Address, error) {
	return _BitcoinState.Contract.RouterV2(&_BitcoinState.CallOpts)
}

// RuneFactory is a free data retrieval call binding the contract method 0x7d1f0a56.
//
// Solidity: function runeFactory() view returns(address)
func (_BitcoinState *BitcoinStateCaller) RuneFactory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BitcoinState.contract.Call(opts, &out, "runeFactory")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RuneFactory is a free data retrieval call binding the contract method 0x7d1f0a56.
//
// Solidity: function runeFactory() view returns(address)
func (_BitcoinState *BitcoinStateSession) RuneFactory() (common.Address, error) {
	return _BitcoinState.Contract.RuneFactory(&_BitcoinState.CallOpts)
}

// RuneFactory is a free data retrieval call binding the contract method 0x7d1f0a56.
//
// Solidity: function runeFactory() view returns(address)
func (_BitcoinState *BitcoinStateCallerSession) RuneFactory() (common.Address, error) {
	return _BitcoinState.Contract.RuneFactory(&_BitcoinState.CallOpts)
}

// Tokens is a free data retrieval call binding the contract method 0x904194a3.
//
// Solidity: function tokens(bytes32 ) view returns(address)
func (_BitcoinState *BitcoinStateCaller) Tokens(opts *bind.CallOpts, arg0 [32]byte) (common.Address, error) {
	var out []interface{}
	err := _BitcoinState.contract.Call(opts, &out, "tokens", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Tokens is a free data retrieval call binding the contract method 0x904194a3.
//
// Solidity: function tokens(bytes32 ) view returns(address)
func (_BitcoinState *BitcoinStateSession) Tokens(arg0 [32]byte) (common.Address, error) {
	return _BitcoinState.Contract.Tokens(&_BitcoinState.CallOpts, arg0)
}

// Tokens is a free data retrieval call binding the contract method 0x904194a3.
//
// Solidity: function tokens(bytes32 ) view returns(address)
func (_BitcoinState *BitcoinStateCallerSession) Tokens(arg0 [32]byte) (common.Address, error) {
	return _BitcoinState.Contract.Tokens(&_BitcoinState.CallOpts, arg0)
}

// XcallService is a free data retrieval call binding the contract method 0x7bf07164.
//
// Solidity: function xcallService() view returns(address)
func (_BitcoinState *BitcoinStateCaller) XcallService(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BitcoinState.contract.Call(opts, &out, "xcallService")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// XcallService is a free data retrieval call binding the contract method 0x7bf07164.
//
// Solidity: function xcallService() view returns(address)
func (_BitcoinState *BitcoinStateSession) XcallService() (common.Address, error) {
	return _BitcoinState.Contract.XcallService(&_BitcoinState.CallOpts)
}

// XcallService is a free data retrieval call binding the contract method 0x7bf07164.
//
// Solidity: function xcallService() view returns(address)
func (_BitcoinState *BitcoinStateCallerSession) XcallService() (common.Address, error) {
	return _BitcoinState.Contract.XcallService(&_BitcoinState.CallOpts)
}

// AddConnection is a paid mutator transaction binding the contract method 0x677dea1d.
//
// Solidity: function addConnection(address connection_) returns()
func (_BitcoinState *BitcoinStateTransactor) AddConnection(opts *bind.TransactOpts, connection_ common.Address) (*types.Transaction, error) {
	return _BitcoinState.contract.Transact(opts, "addConnection", connection_)
}

// AddConnection is a paid mutator transaction binding the contract method 0x677dea1d.
//
// Solidity: function addConnection(address connection_) returns()
func (_BitcoinState *BitcoinStateSession) AddConnection(connection_ common.Address) (*types.Transaction, error) {
	return _BitcoinState.Contract.AddConnection(&_BitcoinState.TransactOpts, connection_)
}

// AddConnection is a paid mutator transaction binding the contract method 0x677dea1d.
//
// Solidity: function addConnection(address connection_) returns()
func (_BitcoinState *BitcoinStateTransactorSession) AddConnection(connection_ common.Address) (*types.Transaction, error) {
	return _BitcoinState.Contract.AddConnection(&_BitcoinState.TransactOpts, connection_)
}

// HandleCallMessage is a paid mutator transaction binding the contract method 0x5d6a16f5.
//
// Solidity: function handleCallMessage(string _from, bytes _data, string[] _protocols) returns()
func (_BitcoinState *BitcoinStateTransactor) HandleCallMessage(opts *bind.TransactOpts, _from string, _data []byte, _protocols []string) (*types.Transaction, error) {
	return _BitcoinState.contract.Transact(opts, "handleCallMessage", _from, _data, _protocols)
}

// HandleCallMessage is a paid mutator transaction binding the contract method 0x5d6a16f5.
//
// Solidity: function handleCallMessage(string _from, bytes _data, string[] _protocols) returns()
func (_BitcoinState *BitcoinStateSession) HandleCallMessage(_from string, _data []byte, _protocols []string) (*types.Transaction, error) {
	return _BitcoinState.Contract.HandleCallMessage(&_BitcoinState.TransactOpts, _from, _data, _protocols)
}

// HandleCallMessage is a paid mutator transaction binding the contract method 0x5d6a16f5.
//
// Solidity: function handleCallMessage(string _from, bytes _data, string[] _protocols) returns()
func (_BitcoinState *BitcoinStateTransactorSession) HandleCallMessage(_from string, _data []byte, _protocols []string) (*types.Transaction, error) {
	return _BitcoinState.Contract.HandleCallMessage(&_BitcoinState.TransactOpts, _from, _data, _protocols)
}

// Initialize is a paid mutator transaction binding the contract method 0xe6bfbfd8.
//
// Solidity: function initialize(address xcall_, address uinswapV3Router_, address nonfungiblePositionManager_, address[] connections) returns()
func (_BitcoinState *BitcoinStateTransactor) Initialize(opts *bind.TransactOpts, xcall_ common.Address, uinswapV3Router_ common.Address, nonfungiblePositionManager_ common.Address, connections []common.Address) (*types.Transaction, error) {
	return _BitcoinState.contract.Transact(opts, "initialize", xcall_, uinswapV3Router_, nonfungiblePositionManager_, connections)
}

// Initialize is a paid mutator transaction binding the contract method 0xe6bfbfd8.
//
// Solidity: function initialize(address xcall_, address uinswapV3Router_, address nonfungiblePositionManager_, address[] connections) returns()
func (_BitcoinState *BitcoinStateSession) Initialize(xcall_ common.Address, uinswapV3Router_ common.Address, nonfungiblePositionManager_ common.Address, connections []common.Address) (*types.Transaction, error) {
	return _BitcoinState.Contract.Initialize(&_BitcoinState.TransactOpts, xcall_, uinswapV3Router_, nonfungiblePositionManager_, connections)
}

// Initialize is a paid mutator transaction binding the contract method 0xe6bfbfd8.
//
// Solidity: function initialize(address xcall_, address uinswapV3Router_, address nonfungiblePositionManager_, address[] connections) returns()
func (_BitcoinState *BitcoinStateTransactorSession) Initialize(xcall_ common.Address, uinswapV3Router_ common.Address, nonfungiblePositionManager_ common.Address, connections []common.Address) (*types.Transaction, error) {
	return _BitcoinState.Contract.Initialize(&_BitcoinState.TransactOpts, xcall_, uinswapV3Router_, nonfungiblePositionManager_, connections)
}

// Migrate is a paid mutator transaction binding the contract method 0x53f954ea.
//
// Solidity: function migrate(string from_, bytes _data) returns()
func (_BitcoinState *BitcoinStateTransactor) Migrate(opts *bind.TransactOpts, from_ string, _data []byte) (*types.Transaction, error) {
	return _BitcoinState.contract.Transact(opts, "migrate", from_, _data)
}

// Migrate is a paid mutator transaction binding the contract method 0x53f954ea.
//
// Solidity: function migrate(string from_, bytes _data) returns()
func (_BitcoinState *BitcoinStateSession) Migrate(from_ string, _data []byte) (*types.Transaction, error) {
	return _BitcoinState.Contract.Migrate(&_BitcoinState.TransactOpts, from_, _data)
}

// Migrate is a paid mutator transaction binding the contract method 0x53f954ea.
//
// Solidity: function migrate(string from_, bytes _data) returns()
func (_BitcoinState *BitcoinStateTransactorSession) Migrate(from_ string, _data []byte) (*types.Transaction, error) {
	return _BitcoinState.Contract.Migrate(&_BitcoinState.TransactOpts, from_, _data)
}

// MigrateComplete is a paid mutator transaction binding the contract method 0xf0ad3762.
//
// Solidity: function migrateComplete() returns()
func (_BitcoinState *BitcoinStateTransactor) MigrateComplete(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BitcoinState.contract.Transact(opts, "migrateComplete")
}

// MigrateComplete is a paid mutator transaction binding the contract method 0xf0ad3762.
//
// Solidity: function migrateComplete() returns()
func (_BitcoinState *BitcoinStateSession) MigrateComplete() (*types.Transaction, error) {
	return _BitcoinState.Contract.MigrateComplete(&_BitcoinState.TransactOpts)
}

// MigrateComplete is a paid mutator transaction binding the contract method 0xf0ad3762.
//
// Solidity: function migrateComplete() returns()
func (_BitcoinState *BitcoinStateTransactorSession) MigrateComplete() (*types.Transaction, error) {
	return _BitcoinState.Contract.MigrateComplete(&_BitcoinState.TransactOpts)
}

// RemoveConnection is a paid mutator transaction binding the contract method 0x65301f0d.
//
// Solidity: function removeConnection(address connection_) returns()
func (_BitcoinState *BitcoinStateTransactor) RemoveConnection(opts *bind.TransactOpts, connection_ common.Address) (*types.Transaction, error) {
	return _BitcoinState.contract.Transact(opts, "removeConnection", connection_)
}

// RemoveConnection is a paid mutator transaction binding the contract method 0x65301f0d.
//
// Solidity: function removeConnection(address connection_) returns()
func (_BitcoinState *BitcoinStateSession) RemoveConnection(connection_ common.Address) (*types.Transaction, error) {
	return _BitcoinState.Contract.RemoveConnection(&_BitcoinState.TransactOpts, connection_)
}

// RemoveConnection is a paid mutator transaction binding the contract method 0x65301f0d.
//
// Solidity: function removeConnection(address connection_) returns()
func (_BitcoinState *BitcoinStateTransactorSession) RemoveConnection(connection_ common.Address) (*types.Transaction, error) {
	return _BitcoinState.Contract.RemoveConnection(&_BitcoinState.TransactOpts, connection_)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BitcoinState *BitcoinStateTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BitcoinState.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BitcoinState *BitcoinStateSession) RenounceOwnership() (*types.Transaction, error) {
	return _BitcoinState.Contract.RenounceOwnership(&_BitcoinState.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BitcoinState *BitcoinStateTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _BitcoinState.Contract.RenounceOwnership(&_BitcoinState.TransactOpts)
}

// SetRuneFactory is a paid mutator transaction binding the contract method 0xa7fb534c.
//
// Solidity: function setRuneFactory(address _runeFactory) returns()
func (_BitcoinState *BitcoinStateTransactor) SetRuneFactory(opts *bind.TransactOpts, _runeFactory common.Address) (*types.Transaction, error) {
	return _BitcoinState.contract.Transact(opts, "setRuneFactory", _runeFactory)
}

// SetRuneFactory is a paid mutator transaction binding the contract method 0xa7fb534c.
//
// Solidity: function setRuneFactory(address _runeFactory) returns()
func (_BitcoinState *BitcoinStateSession) SetRuneFactory(_runeFactory common.Address) (*types.Transaction, error) {
	return _BitcoinState.Contract.SetRuneFactory(&_BitcoinState.TransactOpts, _runeFactory)
}

// SetRuneFactory is a paid mutator transaction binding the contract method 0xa7fb534c.
//
// Solidity: function setRuneFactory(address _runeFactory) returns()
func (_BitcoinState *BitcoinStateTransactorSession) SetRuneFactory(_runeFactory common.Address) (*types.Transaction, error) {
	return _BitcoinState.Contract.SetRuneFactory(&_BitcoinState.TransactOpts, _runeFactory)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_BitcoinState *BitcoinStateTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _BitcoinState.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_BitcoinState *BitcoinStateSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _BitcoinState.Contract.TransferOwnership(&_BitcoinState.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_BitcoinState *BitcoinStateTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _BitcoinState.Contract.TransferOwnership(&_BitcoinState.TransactOpts, newOwner)
}

// WithdrawTo is a paid mutator transaction binding the contract method 0x43aeb854.
//
// Solidity: function withdrawTo(address token_, string receiver_, uint256 amount) returns()
func (_BitcoinState *BitcoinStateTransactor) WithdrawTo(opts *bind.TransactOpts, token_ common.Address, receiver_ string, amount *big.Int) (*types.Transaction, error) {
	return _BitcoinState.contract.Transact(opts, "withdrawTo", token_, receiver_, amount)
}

// WithdrawTo is a paid mutator transaction binding the contract method 0x43aeb854.
//
// Solidity: function withdrawTo(address token_, string receiver_, uint256 amount) returns()
func (_BitcoinState *BitcoinStateSession) WithdrawTo(token_ common.Address, receiver_ string, amount *big.Int) (*types.Transaction, error) {
	return _BitcoinState.Contract.WithdrawTo(&_BitcoinState.TransactOpts, token_, receiver_, amount)
}

// WithdrawTo is a paid mutator transaction binding the contract method 0x43aeb854.
//
// Solidity: function withdrawTo(address token_, string receiver_, uint256 amount) returns()
func (_BitcoinState *BitcoinStateTransactorSession) WithdrawTo(token_ common.Address, receiver_ string, amount *big.Int) (*types.Transaction, error) {
	return _BitcoinState.Contract.WithdrawTo(&_BitcoinState.TransactOpts, token_, receiver_, amount)
}

// BitcoinStateAddConnectionIterator is returned from FilterAddConnection and is used to iterate over the raw logs and unpacked data for AddConnection events raised by the BitcoinState contract.
type BitcoinStateAddConnectionIterator struct {
	Event *BitcoinStateAddConnection // Event containing the contract specifics and raw log

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
func (it *BitcoinStateAddConnectionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BitcoinStateAddConnection)
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
		it.Event = new(BitcoinStateAddConnection)
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
func (it *BitcoinStateAddConnectionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BitcoinStateAddConnectionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BitcoinStateAddConnection represents a AddConnection event raised by the BitcoinState contract.
type BitcoinStateAddConnection struct {
	Connection common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterAddConnection is a free log retrieval operation binding the contract event 0x5cd3d8da8ef00a8b5228348fe7683dae605751d4867ba7ea9fdc8260b9e2c7d3.
//
// Solidity: event AddConnection(address connection_)
func (_BitcoinState *BitcoinStateFilterer) FilterAddConnection(opts *bind.FilterOpts) (*BitcoinStateAddConnectionIterator, error) {

	logs, sub, err := _BitcoinState.contract.FilterLogs(opts, "AddConnection")
	if err != nil {
		return nil, err
	}
	return &BitcoinStateAddConnectionIterator{contract: _BitcoinState.contract, event: "AddConnection", logs: logs, sub: sub}, nil
}

// WatchAddConnection is a free log subscription operation binding the contract event 0x5cd3d8da8ef00a8b5228348fe7683dae605751d4867ba7ea9fdc8260b9e2c7d3.
//
// Solidity: event AddConnection(address connection_)
func (_BitcoinState *BitcoinStateFilterer) WatchAddConnection(opts *bind.WatchOpts, sink chan<- *BitcoinStateAddConnection) (event.Subscription, error) {

	logs, sub, err := _BitcoinState.contract.WatchLogs(opts, "AddConnection")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BitcoinStateAddConnection)
				if err := _BitcoinState.contract.UnpackLog(event, "AddConnection", log); err != nil {
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
func (_BitcoinState *BitcoinStateFilterer) ParseAddConnection(log types.Log) (*BitcoinStateAddConnection, error) {
	event := new(BitcoinStateAddConnection)
	if err := _BitcoinState.contract.UnpackLog(event, "AddConnection", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BitcoinStateInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the BitcoinState contract.
type BitcoinStateInitializedIterator struct {
	Event *BitcoinStateInitialized // Event containing the contract specifics and raw log

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
func (it *BitcoinStateInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BitcoinStateInitialized)
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
		it.Event = new(BitcoinStateInitialized)
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
func (it *BitcoinStateInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BitcoinStateInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BitcoinStateInitialized represents a Initialized event raised by the BitcoinState contract.
type BitcoinStateInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_BitcoinState *BitcoinStateFilterer) FilterInitialized(opts *bind.FilterOpts) (*BitcoinStateInitializedIterator, error) {

	logs, sub, err := _BitcoinState.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &BitcoinStateInitializedIterator{contract: _BitcoinState.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_BitcoinState *BitcoinStateFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *BitcoinStateInitialized) (event.Subscription, error) {

	logs, sub, err := _BitcoinState.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BitcoinStateInitialized)
				if err := _BitcoinState.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_BitcoinState *BitcoinStateFilterer) ParseInitialized(log types.Log) (*BitcoinStateInitialized, error) {
	event := new(BitcoinStateInitialized)
	if err := _BitcoinState.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BitcoinStateOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the BitcoinState contract.
type BitcoinStateOwnershipTransferredIterator struct {
	Event *BitcoinStateOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *BitcoinStateOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BitcoinStateOwnershipTransferred)
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
		it.Event = new(BitcoinStateOwnershipTransferred)
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
func (it *BitcoinStateOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BitcoinStateOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BitcoinStateOwnershipTransferred represents a OwnershipTransferred event raised by the BitcoinState contract.
type BitcoinStateOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_BitcoinState *BitcoinStateFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*BitcoinStateOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _BitcoinState.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &BitcoinStateOwnershipTransferredIterator{contract: _BitcoinState.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_BitcoinState *BitcoinStateFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BitcoinStateOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _BitcoinState.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BitcoinStateOwnershipTransferred)
				if err := _BitcoinState.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_BitcoinState *BitcoinStateFilterer) ParseOwnershipTransferred(log types.Log) (*BitcoinStateOwnershipTransferred, error) {
	event := new(BitcoinStateOwnershipTransferred)
	if err := _BitcoinState.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BitcoinStateRemoveConnectionIterator is returned from FilterRemoveConnection and is used to iterate over the raw logs and unpacked data for RemoveConnection events raised by the BitcoinState contract.
type BitcoinStateRemoveConnectionIterator struct {
	Event *BitcoinStateRemoveConnection // Event containing the contract specifics and raw log

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
func (it *BitcoinStateRemoveConnectionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BitcoinStateRemoveConnection)
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
		it.Event = new(BitcoinStateRemoveConnection)
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
func (it *BitcoinStateRemoveConnectionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BitcoinStateRemoveConnectionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BitcoinStateRemoveConnection represents a RemoveConnection event raised by the BitcoinState contract.
type BitcoinStateRemoveConnection struct {
	Connection common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterRemoveConnection is a free log retrieval operation binding the contract event 0xace8d11a44b7aa536cc46a77b519166c001adc485ba8cfa404e1aa252b07db38.
//
// Solidity: event RemoveConnection(address connection_)
func (_BitcoinState *BitcoinStateFilterer) FilterRemoveConnection(opts *bind.FilterOpts) (*BitcoinStateRemoveConnectionIterator, error) {

	logs, sub, err := _BitcoinState.contract.FilterLogs(opts, "RemoveConnection")
	if err != nil {
		return nil, err
	}
	return &BitcoinStateRemoveConnectionIterator{contract: _BitcoinState.contract, event: "RemoveConnection", logs: logs, sub: sub}, nil
}

// WatchRemoveConnection is a free log subscription operation binding the contract event 0xace8d11a44b7aa536cc46a77b519166c001adc485ba8cfa404e1aa252b07db38.
//
// Solidity: event RemoveConnection(address connection_)
func (_BitcoinState *BitcoinStateFilterer) WatchRemoveConnection(opts *bind.WatchOpts, sink chan<- *BitcoinStateRemoveConnection) (event.Subscription, error) {

	logs, sub, err := _BitcoinState.contract.WatchLogs(opts, "RemoveConnection")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BitcoinStateRemoveConnection)
				if err := _BitcoinState.contract.UnpackLog(event, "RemoveConnection", log); err != nil {
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
func (_BitcoinState *BitcoinStateFilterer) ParseRemoveConnection(log types.Log) (*BitcoinStateRemoveConnection, error) {
	event := new(BitcoinStateRemoveConnection)
	if err := _BitcoinState.contract.UnpackLog(event, "RemoveConnection", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BitcoinStateRequestExecutedIterator is returned from FilterRequestExecuted and is used to iterate over the raw logs and unpacked data for RequestExecuted events raised by the BitcoinState contract.
type BitcoinStateRequestExecutedIterator struct {
	Event *BitcoinStateRequestExecuted // Event containing the contract specifics and raw log

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
func (it *BitcoinStateRequestExecutedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BitcoinStateRequestExecuted)
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
		it.Event = new(BitcoinStateRequestExecuted)
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
func (it *BitcoinStateRequestExecutedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BitcoinStateRequestExecutedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BitcoinStateRequestExecuted represents a RequestExecuted event raised by the BitcoinState contract.
type BitcoinStateRequestExecuted struct {
	Id        *big.Int
	StateRoot [32]byte
	Data      []byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRequestExecuted is a free log retrieval operation binding the contract event 0x9c343316d67a8e28446ef883ab491ece3ff70d3eeaa9fbd13a362a0afd690721.
//
// Solidity: event RequestExecuted(uint256 id, bytes32 stateRoot, bytes data)
func (_BitcoinState *BitcoinStateFilterer) FilterRequestExecuted(opts *bind.FilterOpts) (*BitcoinStateRequestExecutedIterator, error) {

	logs, sub, err := _BitcoinState.contract.FilterLogs(opts, "RequestExecuted")
	if err != nil {
		return nil, err
	}
	return &BitcoinStateRequestExecutedIterator{contract: _BitcoinState.contract, event: "RequestExecuted", logs: logs, sub: sub}, nil
}

// WatchRequestExecuted is a free log subscription operation binding the contract event 0x9c343316d67a8e28446ef883ab491ece3ff70d3eeaa9fbd13a362a0afd690721.
//
// Solidity: event RequestExecuted(uint256 id, bytes32 stateRoot, bytes data)
func (_BitcoinState *BitcoinStateFilterer) WatchRequestExecuted(opts *bind.WatchOpts, sink chan<- *BitcoinStateRequestExecuted) (event.Subscription, error) {

	logs, sub, err := _BitcoinState.contract.WatchLogs(opts, "RequestExecuted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BitcoinStateRequestExecuted)
				if err := _BitcoinState.contract.UnpackLog(event, "RequestExecuted", log); err != nil {
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
func (_BitcoinState *BitcoinStateFilterer) ParseRequestExecuted(log types.Log) (*BitcoinStateRequestExecuted, error) {
	event := new(BitcoinStateRequestExecuted)
	if err := _BitcoinState.contract.UnpackLog(event, "RequestExecuted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
