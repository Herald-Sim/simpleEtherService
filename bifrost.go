// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bifrost

import (
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
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// BifrostABI is the input ABI used to generate the binding from.
const BifrostABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"getValue\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setValue\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// Bifrost is an auto generated Go binding around an Ethereum contract.
type Bifrost struct {
	BifrostCaller     // Read-only binding to the contract
	BifrostTransactor // Write-only binding to the contract
	BifrostFilterer   // Log filterer for contract events
}

// BifrostCaller is an auto generated read-only Go binding around an Ethereum contract.
type BifrostCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BifrostTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BifrostTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BifrostFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BifrostFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BifrostSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BifrostSession struct {
	Contract     *Bifrost          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BifrostCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BifrostCallerSession struct {
	Contract *BifrostCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// BifrostTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BifrostTransactorSession struct {
	Contract     *BifrostTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// BifrostRaw is an auto generated low-level Go binding around an Ethereum contract.
type BifrostRaw struct {
	Contract *Bifrost // Generic contract binding to access the raw methods on
}

// BifrostCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BifrostCallerRaw struct {
	Contract *BifrostCaller // Generic read-only contract binding to access the raw methods on
}

// BifrostTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BifrostTransactorRaw struct {
	Contract *BifrostTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBifrost creates a new instance of Bifrost, bound to a specific deployed contract.
func NewBifrost(address common.Address, backend bind.ContractBackend) (*Bifrost, error) {
	contract, err := bindBifrost(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Bifrost{BifrostCaller: BifrostCaller{contract: contract}, BifrostTransactor: BifrostTransactor{contract: contract}, BifrostFilterer: BifrostFilterer{contract: contract}}, nil
}

// NewBifrostCaller creates a new read-only instance of Bifrost, bound to a specific deployed contract.
func NewBifrostCaller(address common.Address, caller bind.ContractCaller) (*BifrostCaller, error) {
	contract, err := bindBifrost(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BifrostCaller{contract: contract}, nil
}

// NewBifrostTransactor creates a new write-only instance of Bifrost, bound to a specific deployed contract.
func NewBifrostTransactor(address common.Address, transactor bind.ContractTransactor) (*BifrostTransactor, error) {
	contract, err := bindBifrost(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BifrostTransactor{contract: contract}, nil
}

// NewBifrostFilterer creates a new log filterer instance of Bifrost, bound to a specific deployed contract.
func NewBifrostFilterer(address common.Address, filterer bind.ContractFilterer) (*BifrostFilterer, error) {
	contract, err := bindBifrost(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BifrostFilterer{contract: contract}, nil
}

// bindBifrost binds a generic wrapper to an already deployed contract.
func bindBifrost(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BifrostABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Bifrost *BifrostRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Bifrost.Contract.BifrostCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Bifrost *BifrostRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bifrost.Contract.BifrostTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Bifrost *BifrostRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Bifrost.Contract.BifrostTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Bifrost *BifrostCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Bifrost.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Bifrost *BifrostTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bifrost.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Bifrost *BifrostTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Bifrost.Contract.contract.Transact(opts, method, params...)
}

// GetValue is a free data retrieval call binding the contract method 0x20965255.
//
// Solidity: function getValue() constant returns(uint256)
func (_Bifrost *BifrostCaller) GetValue(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Bifrost.contract.Call(opts, out, "getValue")
	return *ret0, err
}

// GetValue is a free data retrieval call binding the contract method 0x20965255.
//
// Solidity: function getValue() constant returns(uint256)
func (_Bifrost *BifrostSession) GetValue() (*big.Int, error) {
	return _Bifrost.Contract.GetValue(&_Bifrost.CallOpts)
}

// GetValue is a free data retrieval call binding the contract method 0x20965255.
//
// Solidity: function getValue() constant returns(uint256)
func (_Bifrost *BifrostCallerSession) GetValue() (*big.Int, error) {
	return _Bifrost.Contract.GetValue(&_Bifrost.CallOpts)
}

// SetValue is a paid mutator transaction binding the contract method 0x55241077.
//
// Solidity: function setValue(uint256 value) returns()
func (_Bifrost *BifrostTransactor) SetValue(opts *bind.TransactOpts, value *big.Int) (*types.Transaction, error) {
	return _Bifrost.contract.Transact(opts, "setValue", value)
}

// SetValue is a paid mutator transaction binding the contract method 0x55241077.
//
// Solidity: function setValue(uint256 value) returns()
func (_Bifrost *BifrostSession) SetValue(value *big.Int) (*types.Transaction, error) {
	return _Bifrost.Contract.SetValue(&_Bifrost.TransactOpts, value)
}

// SetValue is a paid mutator transaction binding the contract method 0x55241077.
//
// Solidity: function setValue(uint256 value) returns()
func (_Bifrost *BifrostTransactorSession) SetValue(value *big.Int) (*types.Transaction, error) {
	return _Bifrost.Contract.SetValue(&_Bifrost.TransactOpts, value)
}
