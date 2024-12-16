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

// RuneFactoryTokenInfo is an auto generated low-level Go binding around an user-defined struct.
type RuneFactoryTokenInfo struct {
	Name     string
	Symbol   string
	Decimals uint8
}

// RunefactoryMetaData contains all meta data concerning the Runefactory contract.
var RunefactoryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"computeTokenAddress\",\"inputs\":[{\"name\":\"tokenName\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deployRune\",\"inputs\":[{\"name\":\"runeName\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"decimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"tokenAddr\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"params\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structRuneFactory.TokenInfo\",\"components\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"symbol\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"decimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"runes\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false}]",
}

// RunefactoryABI is the input ABI used to generate the binding from.
// Deprecated: Use RunefactoryMetaData.ABI instead.
var RunefactoryABI = RunefactoryMetaData.ABI

// Runefactory is an auto generated Go binding around an Ethereum contract.
type Runefactory struct {
	RunefactoryCaller     // Read-only binding to the contract
	RunefactoryTransactor // Write-only binding to the contract
	RunefactoryFilterer   // Log filterer for contract events
}

// RunefactoryCaller is an auto generated read-only Go binding around an Ethereum contract.
type RunefactoryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RunefactoryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RunefactoryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RunefactoryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RunefactoryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RunefactorySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RunefactorySession struct {
	Contract     *Runefactory      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RunefactoryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RunefactoryCallerSession struct {
	Contract *RunefactoryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// RunefactoryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RunefactoryTransactorSession struct {
	Contract     *RunefactoryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// RunefactoryRaw is an auto generated low-level Go binding around an Ethereum contract.
type RunefactoryRaw struct {
	Contract *Runefactory // Generic contract binding to access the raw methods on
}

// RunefactoryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RunefactoryCallerRaw struct {
	Contract *RunefactoryCaller // Generic read-only contract binding to access the raw methods on
}

// RunefactoryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RunefactoryTransactorRaw struct {
	Contract *RunefactoryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRunefactory creates a new instance of Runefactory, bound to a specific deployed contract.
func NewRunefactory(address common.Address, backend bind.ContractBackend) (*Runefactory, error) {
	contract, err := bindRunefactory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Runefactory{RunefactoryCaller: RunefactoryCaller{contract: contract}, RunefactoryTransactor: RunefactoryTransactor{contract: contract}, RunefactoryFilterer: RunefactoryFilterer{contract: contract}}, nil
}

// NewRunefactoryCaller creates a new read-only instance of Runefactory, bound to a specific deployed contract.
func NewRunefactoryCaller(address common.Address, caller bind.ContractCaller) (*RunefactoryCaller, error) {
	contract, err := bindRunefactory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RunefactoryCaller{contract: contract}, nil
}

// NewRunefactoryTransactor creates a new write-only instance of Runefactory, bound to a specific deployed contract.
func NewRunefactoryTransactor(address common.Address, transactor bind.ContractTransactor) (*RunefactoryTransactor, error) {
	contract, err := bindRunefactory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RunefactoryTransactor{contract: contract}, nil
}

// NewRunefactoryFilterer creates a new log filterer instance of Runefactory, bound to a specific deployed contract.
func NewRunefactoryFilterer(address common.Address, filterer bind.ContractFilterer) (*RunefactoryFilterer, error) {
	contract, err := bindRunefactory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RunefactoryFilterer{contract: contract}, nil
}

// bindRunefactory binds a generic wrapper to an already deployed contract.
func bindRunefactory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := RunefactoryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Runefactory *RunefactoryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Runefactory.Contract.RunefactoryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Runefactory *RunefactoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Runefactory.Contract.RunefactoryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Runefactory *RunefactoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Runefactory.Contract.RunefactoryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Runefactory *RunefactoryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Runefactory.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Runefactory *RunefactoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Runefactory.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Runefactory *RunefactoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Runefactory.Contract.contract.Transact(opts, method, params...)
}

// ComputeTokenAddress is a free data retrieval call binding the contract method 0x4aa6af7b.
//
// Solidity: function computeTokenAddress(string tokenName) view returns(address)
func (_Runefactory *RunefactoryCaller) ComputeTokenAddress(opts *bind.CallOpts, tokenName string) (common.Address, error) {
	var out []interface{}
	err := _Runefactory.contract.Call(opts, &out, "computeTokenAddress", tokenName)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ComputeTokenAddress is a free data retrieval call binding the contract method 0x4aa6af7b.
//
// Solidity: function computeTokenAddress(string tokenName) view returns(address)
func (_Runefactory *RunefactorySession) ComputeTokenAddress(tokenName string) (common.Address, error) {
	return _Runefactory.Contract.ComputeTokenAddress(&_Runefactory.CallOpts, tokenName)
}

// ComputeTokenAddress is a free data retrieval call binding the contract method 0x4aa6af7b.
//
// Solidity: function computeTokenAddress(string tokenName) view returns(address)
func (_Runefactory *RunefactoryCallerSession) ComputeTokenAddress(tokenName string) (common.Address, error) {
	return _Runefactory.Contract.ComputeTokenAddress(&_Runefactory.CallOpts, tokenName)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Runefactory *RunefactoryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Runefactory.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Runefactory *RunefactorySession) Owner() (common.Address, error) {
	return _Runefactory.Contract.Owner(&_Runefactory.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Runefactory *RunefactoryCallerSession) Owner() (common.Address, error) {
	return _Runefactory.Contract.Owner(&_Runefactory.CallOpts)
}

// Params is a free data retrieval call binding the contract method 0xcff0ab96.
//
// Solidity: function params() view returns((string,string,uint8))
func (_Runefactory *RunefactoryCaller) Params(opts *bind.CallOpts) (RuneFactoryTokenInfo, error) {
	var out []interface{}
	err := _Runefactory.contract.Call(opts, &out, "params")

	if err != nil {
		return *new(RuneFactoryTokenInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(RuneFactoryTokenInfo)).(*RuneFactoryTokenInfo)

	return out0, err

}

// Params is a free data retrieval call binding the contract method 0xcff0ab96.
//
// Solidity: function params() view returns((string,string,uint8))
func (_Runefactory *RunefactorySession) Params() (RuneFactoryTokenInfo, error) {
	return _Runefactory.Contract.Params(&_Runefactory.CallOpts)
}

// Params is a free data retrieval call binding the contract method 0xcff0ab96.
//
// Solidity: function params() view returns((string,string,uint8))
func (_Runefactory *RunefactoryCallerSession) Params() (RuneFactoryTokenInfo, error) {
	return _Runefactory.Contract.Params(&_Runefactory.CallOpts)
}

// Runes is a free data retrieval call binding the contract method 0x9c11b3a0.
//
// Solidity: function runes(address ) view returns(string)
func (_Runefactory *RunefactoryCaller) Runes(opts *bind.CallOpts, arg0 common.Address) (string, error) {
	var out []interface{}
	err := _Runefactory.contract.Call(opts, &out, "runes", arg0)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Runes is a free data retrieval call binding the contract method 0x9c11b3a0.
//
// Solidity: function runes(address ) view returns(string)
func (_Runefactory *RunefactorySession) Runes(arg0 common.Address) (string, error) {
	return _Runefactory.Contract.Runes(&_Runefactory.CallOpts, arg0)
}

// Runes is a free data retrieval call binding the contract method 0x9c11b3a0.
//
// Solidity: function runes(address ) view returns(string)
func (_Runefactory *RunefactoryCallerSession) Runes(arg0 common.Address) (string, error) {
	return _Runefactory.Contract.Runes(&_Runefactory.CallOpts, arg0)
}

// DeployRune is a paid mutator transaction binding the contract method 0x92b33acc.
//
// Solidity: function deployRune(string runeName, uint8 decimals) returns(address tokenAddr)
func (_Runefactory *RunefactoryTransactor) DeployRune(opts *bind.TransactOpts, runeName string, decimals uint8) (*types.Transaction, error) {
	return _Runefactory.contract.Transact(opts, "deployRune", runeName, decimals)
}

// DeployRune is a paid mutator transaction binding the contract method 0x92b33acc.
//
// Solidity: function deployRune(string runeName, uint8 decimals) returns(address tokenAddr)
func (_Runefactory *RunefactorySession) DeployRune(runeName string, decimals uint8) (*types.Transaction, error) {
	return _Runefactory.Contract.DeployRune(&_Runefactory.TransactOpts, runeName, decimals)
}

// DeployRune is a paid mutator transaction binding the contract method 0x92b33acc.
//
// Solidity: function deployRune(string runeName, uint8 decimals) returns(address tokenAddr)
func (_Runefactory *RunefactoryTransactorSession) DeployRune(runeName string, decimals uint8) (*types.Transaction, error) {
	return _Runefactory.Contract.DeployRune(&_Runefactory.TransactOpts, runeName, decimals)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Runefactory *RunefactoryTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Runefactory.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Runefactory *RunefactorySession) RenounceOwnership() (*types.Transaction, error) {
	return _Runefactory.Contract.RenounceOwnership(&_Runefactory.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Runefactory *RunefactoryTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Runefactory.Contract.RenounceOwnership(&_Runefactory.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Runefactory *RunefactoryTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Runefactory.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Runefactory *RunefactorySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Runefactory.Contract.TransferOwnership(&_Runefactory.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Runefactory *RunefactoryTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Runefactory.Contract.TransferOwnership(&_Runefactory.TransactOpts, newOwner)
}

// RunefactoryOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Runefactory contract.
type RunefactoryOwnershipTransferredIterator struct {
	Event *RunefactoryOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *RunefactoryOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RunefactoryOwnershipTransferred)
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
		it.Event = new(RunefactoryOwnershipTransferred)
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
func (it *RunefactoryOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RunefactoryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RunefactoryOwnershipTransferred represents a OwnershipTransferred event raised by the Runefactory contract.
type RunefactoryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Runefactory *RunefactoryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*RunefactoryOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Runefactory.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &RunefactoryOwnershipTransferredIterator{contract: _Runefactory.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Runefactory *RunefactoryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *RunefactoryOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Runefactory.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RunefactoryOwnershipTransferred)
				if err := _Runefactory.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Runefactory *RunefactoryFilterer) ParseOwnershipTransferred(log types.Log) (*RunefactoryOwnershipTransferred, error) {
	event := new(RunefactoryOwnershipTransferred)
	if err := _Runefactory.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
