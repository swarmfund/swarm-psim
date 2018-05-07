// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package internal

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

// MultisigWalletABI is the input ABI used to generate the binding from.
const MultisigWalletABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"MAX_PENDING_TRANSFERS\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"MAX_SIGNERS\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_signers\",\"type\":\"address[]\"},{\"name\":\"_requiredSignatures\",\"type\":\"uint8\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"signers\",\"type\":\"address[]\"},{\"indexed\":false,\"name\":\"requiredSignatures\",\"type\":\"uint8\"}],\"name\":\"SignersChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"creator\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"transferId\",\"type\":\"uint256\"}],\"name\":\"PendingEtherTransferCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"creator\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"transferId\",\"type\":\"uint256\"}],\"name\":\"PendingTokenTransferCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"transferId\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"signer\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"signaturesCount\",\"type\":\"uint256\"}],\"name\":\"PendingTransferConfirmed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"transferId\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"EtherTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"transferId\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TokensTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"limit\",\"type\":\"uint256\"}],\"name\":\"PendingTransfersLimitReached\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Pause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Unpause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"getSigners\",\"outputs\":[{\"name\":\"_signers\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getRequiredSignatures\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getPendingTransfers\",\"outputs\":[{\"name\":\"_pendingTransfersIds\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"transferId\",\"type\":\"uint256\"}],\"name\":\"getPendingTransfer\",\"outputs\":[{\"name\":\"id\",\"type\":\"uint256\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"transferType\",\"type\":\"uint8\"},{\"name\":\"numberOrSignatures\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_signers\",\"type\":\"address[]\"},{\"name\":\"_requiredSignatures\",\"type\":\"uint8\"}],\"name\":\"changeSigners\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"createEtherTransfer\",\"outputs\":[{\"name\":\"_transferId\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_token\",\"type\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"createTokenTransfer\",\"outputs\":[{\"name\":\"_transferId\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_transferId\",\"type\":\"uint256\"}],\"name\":\"confirmTransfer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// MultisigWallet is an auto generated Go binding around an Ethereum contract.
type MultisigWallet struct {
	MultisigWalletCaller     // Read-only binding to the contract
	MultisigWalletTransactor // Write-only binding to the contract
	MultisigWalletFilterer   // Log filterer for contract events
}

// MultisigWalletCaller is an auto generated read-only Go binding around an Ethereum contract.
type MultisigWalletCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultisigWalletTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MultisigWalletTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultisigWalletFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MultisigWalletFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultisigWalletSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MultisigWalletSession struct {
	Contract     *MultisigWallet   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MultisigWalletCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MultisigWalletCallerSession struct {
	Contract *MultisigWalletCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// MultisigWalletTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MultisigWalletTransactorSession struct {
	Contract     *MultisigWalletTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// MultisigWalletRaw is an auto generated low-level Go binding around an Ethereum contract.
type MultisigWalletRaw struct {
	Contract *MultisigWallet // Generic contract binding to access the raw methods on
}

// MultisigWalletCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MultisigWalletCallerRaw struct {
	Contract *MultisigWalletCaller // Generic read-only contract binding to access the raw methods on
}

// MultisigWalletTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MultisigWalletTransactorRaw struct {
	Contract *MultisigWalletTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMultisigWallet creates a new instance of MultisigWallet, bound to a specific deployed contract.
func NewMultisigWallet(address common.Address, backend bind.ContractBackend) (*MultisigWallet, error) {
	contract, err := bindMultisigWallet(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MultisigWallet{MultisigWalletCaller: MultisigWalletCaller{contract: contract}, MultisigWalletTransactor: MultisigWalletTransactor{contract: contract}, MultisigWalletFilterer: MultisigWalletFilterer{contract: contract}}, nil
}

// NewMultisigWalletCaller creates a new read-only instance of MultisigWallet, bound to a specific deployed contract.
func NewMultisigWalletCaller(address common.Address, caller bind.ContractCaller) (*MultisigWalletCaller, error) {
	contract, err := bindMultisigWallet(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MultisigWalletCaller{contract: contract}, nil
}

// NewMultisigWalletTransactor creates a new write-only instance of MultisigWallet, bound to a specific deployed contract.
func NewMultisigWalletTransactor(address common.Address, transactor bind.ContractTransactor) (*MultisigWalletTransactor, error) {
	contract, err := bindMultisigWallet(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MultisigWalletTransactor{contract: contract}, nil
}

// NewMultisigWalletFilterer creates a new log filterer instance of MultisigWallet, bound to a specific deployed contract.
func NewMultisigWalletFilterer(address common.Address, filterer bind.ContractFilterer) (*MultisigWalletFilterer, error) {
	contract, err := bindMultisigWallet(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MultisigWalletFilterer{contract: contract}, nil
}

// bindMultisigWallet binds a generic wrapper to an already deployed contract.
func bindMultisigWallet(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(MultisigWalletABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MultisigWallet *MultisigWalletRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _MultisigWallet.Contract.MultisigWalletCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MultisigWallet *MultisigWalletRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MultisigWallet.Contract.MultisigWalletTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MultisigWallet *MultisigWalletRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MultisigWallet.Contract.MultisigWalletTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MultisigWallet *MultisigWalletCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _MultisigWallet.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MultisigWallet *MultisigWalletTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MultisigWallet.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MultisigWallet *MultisigWalletTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MultisigWallet.Contract.contract.Transact(opts, method, params...)
}

// MAXPENDINGTRANSFERS is a free data retrieval call binding the contract method 0x16f4d52d.
//
// Solidity: function MAX_PENDING_TRANSFERS() constant returns(uint256)
func (_MultisigWallet *MultisigWalletCaller) MAXPENDINGTRANSFERS(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _MultisigWallet.contract.Call(opts, out, "MAX_PENDING_TRANSFERS")
	return *ret0, err
}

// MAXPENDINGTRANSFERS is a free data retrieval call binding the contract method 0x16f4d52d.
//
// Solidity: function MAX_PENDING_TRANSFERS() constant returns(uint256)
func (_MultisigWallet *MultisigWalletSession) MAXPENDINGTRANSFERS() (*big.Int, error) {
	return _MultisigWallet.Contract.MAXPENDINGTRANSFERS(&_MultisigWallet.CallOpts)
}

// MAXPENDINGTRANSFERS is a free data retrieval call binding the contract method 0x16f4d52d.
//
// Solidity: function MAX_PENDING_TRANSFERS() constant returns(uint256)
func (_MultisigWallet *MultisigWalletCallerSession) MAXPENDINGTRANSFERS() (*big.Int, error) {
	return _MultisigWallet.Contract.MAXPENDINGTRANSFERS(&_MultisigWallet.CallOpts)
}

// MAXSIGNERS is a free data retrieval call binding the contract method 0x59ecd657.
//
// Solidity: function MAX_SIGNERS() constant returns(uint256)
func (_MultisigWallet *MultisigWalletCaller) MAXSIGNERS(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _MultisigWallet.contract.Call(opts, out, "MAX_SIGNERS")
	return *ret0, err
}

// MAXSIGNERS is a free data retrieval call binding the contract method 0x59ecd657.
//
// Solidity: function MAX_SIGNERS() constant returns(uint256)
func (_MultisigWallet *MultisigWalletSession) MAXSIGNERS() (*big.Int, error) {
	return _MultisigWallet.Contract.MAXSIGNERS(&_MultisigWallet.CallOpts)
}

// MAXSIGNERS is a free data retrieval call binding the contract method 0x59ecd657.
//
// Solidity: function MAX_SIGNERS() constant returns(uint256)
func (_MultisigWallet *MultisigWalletCallerSession) MAXSIGNERS() (*big.Int, error) {
	return _MultisigWallet.Contract.MAXSIGNERS(&_MultisigWallet.CallOpts)
}

// GetPendingTransfer is a free data retrieval call binding the contract method 0xf91a7b89.
//
// Solidity: function getPendingTransfer(transferId uint256) constant returns(id uint256, to address, amount uint256, token address, transferType uint8, numberOrSignatures uint8)
func (_MultisigWallet *MultisigWalletCaller) GetPendingTransfer(opts *bind.CallOpts, transferId *big.Int) (struct {
	Id                 *big.Int
	To                 common.Address
	Amount             *big.Int
	Token              common.Address
	TransferType       uint8
	NumberOrSignatures uint8
}, error) {
	ret := new(struct {
		Id                 *big.Int
		To                 common.Address
		Amount             *big.Int
		Token              common.Address
		TransferType       uint8
		NumberOrSignatures uint8
	})
	out := ret
	err := _MultisigWallet.contract.Call(opts, out, "getPendingTransfer", transferId)
	return *ret, err
}

// GetPendingTransfer is a free data retrieval call binding the contract method 0xf91a7b89.
//
// Solidity: function getPendingTransfer(transferId uint256) constant returns(id uint256, to address, amount uint256, token address, transferType uint8, numberOrSignatures uint8)
func (_MultisigWallet *MultisigWalletSession) GetPendingTransfer(transferId *big.Int) (struct {
	Id                 *big.Int
	To                 common.Address
	Amount             *big.Int
	Token              common.Address
	TransferType       uint8
	NumberOrSignatures uint8
}, error) {
	return _MultisigWallet.Contract.GetPendingTransfer(&_MultisigWallet.CallOpts, transferId)
}

// GetPendingTransfer is a free data retrieval call binding the contract method 0xf91a7b89.
//
// Solidity: function getPendingTransfer(transferId uint256) constant returns(id uint256, to address, amount uint256, token address, transferType uint8, numberOrSignatures uint8)
func (_MultisigWallet *MultisigWalletCallerSession) GetPendingTransfer(transferId *big.Int) (struct {
	Id                 *big.Int
	To                 common.Address
	Amount             *big.Int
	Token              common.Address
	TransferType       uint8
	NumberOrSignatures uint8
}, error) {
	return _MultisigWallet.Contract.GetPendingTransfer(&_MultisigWallet.CallOpts, transferId)
}

// GetPendingTransfers is a free data retrieval call binding the contract method 0x448878fc.
//
// Solidity: function getPendingTransfers() constant returns(_pendingTransfersIds uint256[])
func (_MultisigWallet *MultisigWalletCaller) GetPendingTransfers(opts *bind.CallOpts) ([]*big.Int, error) {
	var (
		ret0 = new([]*big.Int)
	)
	out := ret0
	err := _MultisigWallet.contract.Call(opts, out, "getPendingTransfers")
	return *ret0, err
}

// GetPendingTransfers is a free data retrieval call binding the contract method 0x448878fc.
//
// Solidity: function getPendingTransfers() constant returns(_pendingTransfersIds uint256[])
func (_MultisigWallet *MultisigWalletSession) GetPendingTransfers() ([]*big.Int, error) {
	return _MultisigWallet.Contract.GetPendingTransfers(&_MultisigWallet.CallOpts)
}

// GetPendingTransfers is a free data retrieval call binding the contract method 0x448878fc.
//
// Solidity: function getPendingTransfers() constant returns(_pendingTransfersIds uint256[])
func (_MultisigWallet *MultisigWalletCallerSession) GetPendingTransfers() ([]*big.Int, error) {
	return _MultisigWallet.Contract.GetPendingTransfers(&_MultisigWallet.CallOpts)
}

// GetRequiredSignatures is a free data retrieval call binding the contract method 0xccd93998.
//
// Solidity: function getRequiredSignatures() constant returns(uint8)
func (_MultisigWallet *MultisigWalletCaller) GetRequiredSignatures(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _MultisigWallet.contract.Call(opts, out, "getRequiredSignatures")
	return *ret0, err
}

// GetRequiredSignatures is a free data retrieval call binding the contract method 0xccd93998.
//
// Solidity: function getRequiredSignatures() constant returns(uint8)
func (_MultisigWallet *MultisigWalletSession) GetRequiredSignatures() (uint8, error) {
	return _MultisigWallet.Contract.GetRequiredSignatures(&_MultisigWallet.CallOpts)
}

// GetRequiredSignatures is a free data retrieval call binding the contract method 0xccd93998.
//
// Solidity: function getRequiredSignatures() constant returns(uint8)
func (_MultisigWallet *MultisigWalletCallerSession) GetRequiredSignatures() (uint8, error) {
	return _MultisigWallet.Contract.GetRequiredSignatures(&_MultisigWallet.CallOpts)
}

// GetSigners is a free data retrieval call binding the contract method 0x94cf795e.
//
// Solidity: function getSigners() constant returns(_signers address[])
func (_MultisigWallet *MultisigWalletCaller) GetSigners(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _MultisigWallet.contract.Call(opts, out, "getSigners")
	return *ret0, err
}

// GetSigners is a free data retrieval call binding the contract method 0x94cf795e.
//
// Solidity: function getSigners() constant returns(_signers address[])
func (_MultisigWallet *MultisigWalletSession) GetSigners() ([]common.Address, error) {
	return _MultisigWallet.Contract.GetSigners(&_MultisigWallet.CallOpts)
}

// GetSigners is a free data retrieval call binding the contract method 0x94cf795e.
//
// Solidity: function getSigners() constant returns(_signers address[])
func (_MultisigWallet *MultisigWalletCallerSession) GetSigners() ([]common.Address, error) {
	return _MultisigWallet.Contract.GetSigners(&_MultisigWallet.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_MultisigWallet *MultisigWalletCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _MultisigWallet.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_MultisigWallet *MultisigWalletSession) Owner() (common.Address, error) {
	return _MultisigWallet.Contract.Owner(&_MultisigWallet.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_MultisigWallet *MultisigWalletCallerSession) Owner() (common.Address, error) {
	return _MultisigWallet.Contract.Owner(&_MultisigWallet.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() constant returns(bool)
func (_MultisigWallet *MultisigWalletCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _MultisigWallet.contract.Call(opts, out, "paused")
	return *ret0, err
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() constant returns(bool)
func (_MultisigWallet *MultisigWalletSession) Paused() (bool, error) {
	return _MultisigWallet.Contract.Paused(&_MultisigWallet.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() constant returns(bool)
func (_MultisigWallet *MultisigWalletCallerSession) Paused() (bool, error) {
	return _MultisigWallet.Contract.Paused(&_MultisigWallet.CallOpts)
}

// ChangeSigners is a paid mutator transaction binding the contract method 0x1878a5c9.
//
// Solidity: function changeSigners(_signers address[], _requiredSignatures uint8) returns()
func (_MultisigWallet *MultisigWalletTransactor) ChangeSigners(opts *bind.TransactOpts, _signers []common.Address, _requiredSignatures uint8) (*types.Transaction, error) {
	return _MultisigWallet.contract.Transact(opts, "changeSigners", _signers, _requiredSignatures)
}

// ChangeSigners is a paid mutator transaction binding the contract method 0x1878a5c9.
//
// Solidity: function changeSigners(_signers address[], _requiredSignatures uint8) returns()
func (_MultisigWallet *MultisigWalletSession) ChangeSigners(_signers []common.Address, _requiredSignatures uint8) (*types.Transaction, error) {
	return _MultisigWallet.Contract.ChangeSigners(&_MultisigWallet.TransactOpts, _signers, _requiredSignatures)
}

// ChangeSigners is a paid mutator transaction binding the contract method 0x1878a5c9.
//
// Solidity: function changeSigners(_signers address[], _requiredSignatures uint8) returns()
func (_MultisigWallet *MultisigWalletTransactorSession) ChangeSigners(_signers []common.Address, _requiredSignatures uint8) (*types.Transaction, error) {
	return _MultisigWallet.Contract.ChangeSigners(&_MultisigWallet.TransactOpts, _signers, _requiredSignatures)
}

// ConfirmTransfer is a paid mutator transaction binding the contract method 0x2c48e7db.
//
// Solidity: function confirmTransfer(_transferId uint256) returns()
func (_MultisigWallet *MultisigWalletTransactor) ConfirmTransfer(opts *bind.TransactOpts, _transferId *big.Int) (*types.Transaction, error) {
	return _MultisigWallet.contract.Transact(opts, "confirmTransfer", _transferId)
}

// ConfirmTransfer is a paid mutator transaction binding the contract method 0x2c48e7db.
//
// Solidity: function confirmTransfer(_transferId uint256) returns()
func (_MultisigWallet *MultisigWalletSession) ConfirmTransfer(_transferId *big.Int) (*types.Transaction, error) {
	return _MultisigWallet.Contract.ConfirmTransfer(&_MultisigWallet.TransactOpts, _transferId)
}

// ConfirmTransfer is a paid mutator transaction binding the contract method 0x2c48e7db.
//
// Solidity: function confirmTransfer(_transferId uint256) returns()
func (_MultisigWallet *MultisigWalletTransactorSession) ConfirmTransfer(_transferId *big.Int) (*types.Transaction, error) {
	return _MultisigWallet.Contract.ConfirmTransfer(&_MultisigWallet.TransactOpts, _transferId)
}

// CreateEtherTransfer is a paid mutator transaction binding the contract method 0x413d3cf1.
//
// Solidity: function createEtherTransfer(_to address, _amount uint256) returns(_transferId uint256)
func (_MultisigWallet *MultisigWalletTransactor) CreateEtherTransfer(opts *bind.TransactOpts, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _MultisigWallet.contract.Transact(opts, "createEtherTransfer", _to, _amount)
}

// CreateEtherTransfer is a paid mutator transaction binding the contract method 0x413d3cf1.
//
// Solidity: function createEtherTransfer(_to address, _amount uint256) returns(_transferId uint256)
func (_MultisigWallet *MultisigWalletSession) CreateEtherTransfer(_to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _MultisigWallet.Contract.CreateEtherTransfer(&_MultisigWallet.TransactOpts, _to, _amount)
}

// CreateEtherTransfer is a paid mutator transaction binding the contract method 0x413d3cf1.
//
// Solidity: function createEtherTransfer(_to address, _amount uint256) returns(_transferId uint256)
func (_MultisigWallet *MultisigWalletTransactorSession) CreateEtherTransfer(_to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _MultisigWallet.Contract.CreateEtherTransfer(&_MultisigWallet.TransactOpts, _to, _amount)
}

// CreateTokenTransfer is a paid mutator transaction binding the contract method 0x01c0954b.
//
// Solidity: function createTokenTransfer(_to address, _token address, _amount uint256) returns(_transferId uint256)
func (_MultisigWallet *MultisigWalletTransactor) CreateTokenTransfer(opts *bind.TransactOpts, _to common.Address, _token common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _MultisigWallet.contract.Transact(opts, "createTokenTransfer", _to, _token, _amount)
}

// CreateTokenTransfer is a paid mutator transaction binding the contract method 0x01c0954b.
//
// Solidity: function createTokenTransfer(_to address, _token address, _amount uint256) returns(_transferId uint256)
func (_MultisigWallet *MultisigWalletSession) CreateTokenTransfer(_to common.Address, _token common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _MultisigWallet.Contract.CreateTokenTransfer(&_MultisigWallet.TransactOpts, _to, _token, _amount)
}

// CreateTokenTransfer is a paid mutator transaction binding the contract method 0x01c0954b.
//
// Solidity: function createTokenTransfer(_to address, _token address, _amount uint256) returns(_transferId uint256)
func (_MultisigWallet *MultisigWalletTransactorSession) CreateTokenTransfer(_to common.Address, _token common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _MultisigWallet.Contract.CreateTokenTransfer(&_MultisigWallet.TransactOpts, _to, _token, _amount)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_MultisigWallet *MultisigWalletTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MultisigWallet.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_MultisigWallet *MultisigWalletSession) Pause() (*types.Transaction, error) {
	return _MultisigWallet.Contract.Pause(&_MultisigWallet.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_MultisigWallet *MultisigWalletTransactorSession) Pause() (*types.Transaction, error) {
	return _MultisigWallet.Contract.Pause(&_MultisigWallet.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_MultisigWallet *MultisigWalletTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _MultisigWallet.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_MultisigWallet *MultisigWalletSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _MultisigWallet.Contract.TransferOwnership(&_MultisigWallet.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_MultisigWallet *MultisigWalletTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _MultisigWallet.Contract.TransferOwnership(&_MultisigWallet.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_MultisigWallet *MultisigWalletTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MultisigWallet.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_MultisigWallet *MultisigWalletSession) Unpause() (*types.Transaction, error) {
	return _MultisigWallet.Contract.Unpause(&_MultisigWallet.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_MultisigWallet *MultisigWalletTransactorSession) Unpause() (*types.Transaction, error) {
	return _MultisigWallet.Contract.Unpause(&_MultisigWallet.TransactOpts)
}

// MultisigWalletDepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the MultisigWallet contract.
type MultisigWalletDepositIterator struct {
	Event *MultisigWalletDeposit // Event containing the contract specifics and raw log

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
func (it *MultisigWalletDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultisigWalletDeposit)
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
		it.Event = new(MultisigWalletDeposit)
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
func (it *MultisigWalletDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultisigWalletDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultisigWalletDeposit represents a Deposit event raised by the MultisigWallet contract.
type MultisigWalletDeposit struct {
	From   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c.
//
// Solidity: event Deposit(from indexed address, amount uint256)
func (_MultisigWallet *MultisigWalletFilterer) FilterDeposit(opts *bind.FilterOpts, from []common.Address) (*MultisigWalletDepositIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _MultisigWallet.contract.FilterLogs(opts, "Deposit", fromRule)
	if err != nil {
		return nil, err
	}
	return &MultisigWalletDepositIterator{contract: _MultisigWallet.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c.
//
// Solidity: event Deposit(from indexed address, amount uint256)
func (_MultisigWallet *MultisigWalletFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *MultisigWalletDeposit, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _MultisigWallet.contract.WatchLogs(opts, "Deposit", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultisigWalletDeposit)
				if err := _MultisigWallet.contract.UnpackLog(event, "Deposit", log); err != nil {
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

// MultisigWalletEtherTransferredIterator is returned from FilterEtherTransferred and is used to iterate over the raw logs and unpacked data for EtherTransferred events raised by the MultisigWallet contract.
type MultisigWalletEtherTransferredIterator struct {
	Event *MultisigWalletEtherTransferred // Event containing the contract specifics and raw log

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
func (it *MultisigWalletEtherTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultisigWalletEtherTransferred)
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
		it.Event = new(MultisigWalletEtherTransferred)
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
func (it *MultisigWalletEtherTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultisigWalletEtherTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultisigWalletEtherTransferred represents a EtherTransferred event raised by the MultisigWallet contract.
type MultisigWalletEtherTransferred struct {
	TransferId *big.Int
	To         common.Address
	Amount     *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterEtherTransferred is a free log retrieval operation binding the contract event 0x8d6a64c40ae95e2d8ba6157294f8d842d794fc27c5023b5cc12193e57fa4237b.
//
// Solidity: event EtherTransferred(transferId indexed uint256, to indexed address, amount uint256)
func (_MultisigWallet *MultisigWalletFilterer) FilterEtherTransferred(opts *bind.FilterOpts, transferId []*big.Int, to []common.Address) (*MultisigWalletEtherTransferredIterator, error) {

	var transferIdRule []interface{}
	for _, transferIdItem := range transferId {
		transferIdRule = append(transferIdRule, transferIdItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MultisigWallet.contract.FilterLogs(opts, "EtherTransferred", transferIdRule, toRule)
	if err != nil {
		return nil, err
	}
	return &MultisigWalletEtherTransferredIterator{contract: _MultisigWallet.contract, event: "EtherTransferred", logs: logs, sub: sub}, nil
}

// WatchEtherTransferred is a free log subscription operation binding the contract event 0x8d6a64c40ae95e2d8ba6157294f8d842d794fc27c5023b5cc12193e57fa4237b.
//
// Solidity: event EtherTransferred(transferId indexed uint256, to indexed address, amount uint256)
func (_MultisigWallet *MultisigWalletFilterer) WatchEtherTransferred(opts *bind.WatchOpts, sink chan<- *MultisigWalletEtherTransferred, transferId []*big.Int, to []common.Address) (event.Subscription, error) {

	var transferIdRule []interface{}
	for _, transferIdItem := range transferId {
		transferIdRule = append(transferIdRule, transferIdItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MultisigWallet.contract.WatchLogs(opts, "EtherTransferred", transferIdRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultisigWalletEtherTransferred)
				if err := _MultisigWallet.contract.UnpackLog(event, "EtherTransferred", log); err != nil {
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

// MultisigWalletOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the MultisigWallet contract.
type MultisigWalletOwnershipTransferredIterator struct {
	Event *MultisigWalletOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *MultisigWalletOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultisigWalletOwnershipTransferred)
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
		it.Event = new(MultisigWalletOwnershipTransferred)
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
func (it *MultisigWalletOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultisigWalletOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultisigWalletOwnershipTransferred represents a OwnershipTransferred event raised by the MultisigWallet contract.
type MultisigWalletOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_MultisigWallet *MultisigWalletFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*MultisigWalletOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _MultisigWallet.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &MultisigWalletOwnershipTransferredIterator{contract: _MultisigWallet.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_MultisigWallet *MultisigWalletFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *MultisigWalletOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _MultisigWallet.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultisigWalletOwnershipTransferred)
				if err := _MultisigWallet.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// MultisigWalletPauseIterator is returned from FilterPause and is used to iterate over the raw logs and unpacked data for Pause events raised by the MultisigWallet contract.
type MultisigWalletPauseIterator struct {
	Event *MultisigWalletPause // Event containing the contract specifics and raw log

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
func (it *MultisigWalletPauseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultisigWalletPause)
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
		it.Event = new(MultisigWalletPause)
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
func (it *MultisigWalletPauseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultisigWalletPauseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultisigWalletPause represents a Pause event raised by the MultisigWallet contract.
type MultisigWalletPause struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterPause is a free log retrieval operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
//
// Solidity: event Pause()
func (_MultisigWallet *MultisigWalletFilterer) FilterPause(opts *bind.FilterOpts) (*MultisigWalletPauseIterator, error) {

	logs, sub, err := _MultisigWallet.contract.FilterLogs(opts, "Pause")
	if err != nil {
		return nil, err
	}
	return &MultisigWalletPauseIterator{contract: _MultisigWallet.contract, event: "Pause", logs: logs, sub: sub}, nil
}

// WatchPause is a free log subscription operation binding the contract event 0x6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff625.
//
// Solidity: event Pause()
func (_MultisigWallet *MultisigWalletFilterer) WatchPause(opts *bind.WatchOpts, sink chan<- *MultisigWalletPause) (event.Subscription, error) {

	logs, sub, err := _MultisigWallet.contract.WatchLogs(opts, "Pause")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultisigWalletPause)
				if err := _MultisigWallet.contract.UnpackLog(event, "Pause", log); err != nil {
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

// MultisigWalletPendingEtherTransferCreatedIterator is returned from FilterPendingEtherTransferCreated and is used to iterate over the raw logs and unpacked data for PendingEtherTransferCreated events raised by the MultisigWallet contract.
type MultisigWalletPendingEtherTransferCreatedIterator struct {
	Event *MultisigWalletPendingEtherTransferCreated // Event containing the contract specifics and raw log

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
func (it *MultisigWalletPendingEtherTransferCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultisigWalletPendingEtherTransferCreated)
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
		it.Event = new(MultisigWalletPendingEtherTransferCreated)
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
func (it *MultisigWalletPendingEtherTransferCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultisigWalletPendingEtherTransferCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultisigWalletPendingEtherTransferCreated represents a PendingEtherTransferCreated event raised by the MultisigWallet contract.
type MultisigWalletPendingEtherTransferCreated struct {
	Creator    common.Address
	To         common.Address
	Amount     *big.Int
	TransferId *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterPendingEtherTransferCreated is a free log retrieval operation binding the contract event 0xe7f0c7c6807ccf602c8eeedddf80ca86ab53f66740779a25c3c93abe0fc809d1.
//
// Solidity: event PendingEtherTransferCreated(creator indexed address, to indexed address, amount uint256, transferId uint256)
func (_MultisigWallet *MultisigWalletFilterer) FilterPendingEtherTransferCreated(opts *bind.FilterOpts, creator []common.Address, to []common.Address) (*MultisigWalletPendingEtherTransferCreatedIterator, error) {

	var creatorRule []interface{}
	for _, creatorItem := range creator {
		creatorRule = append(creatorRule, creatorItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MultisigWallet.contract.FilterLogs(opts, "PendingEtherTransferCreated", creatorRule, toRule)
	if err != nil {
		return nil, err
	}
	return &MultisigWalletPendingEtherTransferCreatedIterator{contract: _MultisigWallet.contract, event: "PendingEtherTransferCreated", logs: logs, sub: sub}, nil
}

// WatchPendingEtherTransferCreated is a free log subscription operation binding the contract event 0xe7f0c7c6807ccf602c8eeedddf80ca86ab53f66740779a25c3c93abe0fc809d1.
//
// Solidity: event PendingEtherTransferCreated(creator indexed address, to indexed address, amount uint256, transferId uint256)
func (_MultisigWallet *MultisigWalletFilterer) WatchPendingEtherTransferCreated(opts *bind.WatchOpts, sink chan<- *MultisigWalletPendingEtherTransferCreated, creator []common.Address, to []common.Address) (event.Subscription, error) {

	var creatorRule []interface{}
	for _, creatorItem := range creator {
		creatorRule = append(creatorRule, creatorItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _MultisigWallet.contract.WatchLogs(opts, "PendingEtherTransferCreated", creatorRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultisigWalletPendingEtherTransferCreated)
				if err := _MultisigWallet.contract.UnpackLog(event, "PendingEtherTransferCreated", log); err != nil {
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

// MultisigWalletPendingTokenTransferCreatedIterator is returned from FilterPendingTokenTransferCreated and is used to iterate over the raw logs and unpacked data for PendingTokenTransferCreated events raised by the MultisigWallet contract.
type MultisigWalletPendingTokenTransferCreatedIterator struct {
	Event *MultisigWalletPendingTokenTransferCreated // Event containing the contract specifics and raw log

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
func (it *MultisigWalletPendingTokenTransferCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultisigWalletPendingTokenTransferCreated)
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
		it.Event = new(MultisigWalletPendingTokenTransferCreated)
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
func (it *MultisigWalletPendingTokenTransferCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultisigWalletPendingTokenTransferCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultisigWalletPendingTokenTransferCreated represents a PendingTokenTransferCreated event raised by the MultisigWallet contract.
type MultisigWalletPendingTokenTransferCreated struct {
	Creator    common.Address
	To         common.Address
	Token      common.Address
	Amount     *big.Int
	TransferId *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterPendingTokenTransferCreated is a free log retrieval operation binding the contract event 0x670c51d7c9ae60ee716af92c18abb6284cc1e377e7b32ce4f710d2ffd3a8e17b.
//
// Solidity: event PendingTokenTransferCreated(creator indexed address, to indexed address, token indexed address, amount uint256, transferId uint256)
func (_MultisigWallet *MultisigWalletFilterer) FilterPendingTokenTransferCreated(opts *bind.FilterOpts, creator []common.Address, to []common.Address, token []common.Address) (*MultisigWalletPendingTokenTransferCreatedIterator, error) {

	var creatorRule []interface{}
	for _, creatorItem := range creator {
		creatorRule = append(creatorRule, creatorItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _MultisigWallet.contract.FilterLogs(opts, "PendingTokenTransferCreated", creatorRule, toRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return &MultisigWalletPendingTokenTransferCreatedIterator{contract: _MultisigWallet.contract, event: "PendingTokenTransferCreated", logs: logs, sub: sub}, nil
}

// WatchPendingTokenTransferCreated is a free log subscription operation binding the contract event 0x670c51d7c9ae60ee716af92c18abb6284cc1e377e7b32ce4f710d2ffd3a8e17b.
//
// Solidity: event PendingTokenTransferCreated(creator indexed address, to indexed address, token indexed address, amount uint256, transferId uint256)
func (_MultisigWallet *MultisigWalletFilterer) WatchPendingTokenTransferCreated(opts *bind.WatchOpts, sink chan<- *MultisigWalletPendingTokenTransferCreated, creator []common.Address, to []common.Address, token []common.Address) (event.Subscription, error) {

	var creatorRule []interface{}
	for _, creatorItem := range creator {
		creatorRule = append(creatorRule, creatorItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _MultisigWallet.contract.WatchLogs(opts, "PendingTokenTransferCreated", creatorRule, toRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultisigWalletPendingTokenTransferCreated)
				if err := _MultisigWallet.contract.UnpackLog(event, "PendingTokenTransferCreated", log); err != nil {
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

// MultisigWalletPendingTransferConfirmedIterator is returned from FilterPendingTransferConfirmed and is used to iterate over the raw logs and unpacked data for PendingTransferConfirmed events raised by the MultisigWallet contract.
type MultisigWalletPendingTransferConfirmedIterator struct {
	Event *MultisigWalletPendingTransferConfirmed // Event containing the contract specifics and raw log

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
func (it *MultisigWalletPendingTransferConfirmedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultisigWalletPendingTransferConfirmed)
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
		it.Event = new(MultisigWalletPendingTransferConfirmed)
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
func (it *MultisigWalletPendingTransferConfirmedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultisigWalletPendingTransferConfirmedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultisigWalletPendingTransferConfirmed represents a PendingTransferConfirmed event raised by the MultisigWallet contract.
type MultisigWalletPendingTransferConfirmed struct {
	TransferId      *big.Int
	Signer          common.Address
	SignaturesCount *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterPendingTransferConfirmed is a free log retrieval operation binding the contract event 0x1c5a4e2f8e0eca8132383b7783986b59dd5c1bd78063bb1262cdac1b8df7507f.
//
// Solidity: event PendingTransferConfirmed(transferId indexed uint256, signer indexed address, signaturesCount uint256)
func (_MultisigWallet *MultisigWalletFilterer) FilterPendingTransferConfirmed(opts *bind.FilterOpts, transferId []*big.Int, signer []common.Address) (*MultisigWalletPendingTransferConfirmedIterator, error) {

	var transferIdRule []interface{}
	for _, transferIdItem := range transferId {
		transferIdRule = append(transferIdRule, transferIdItem)
	}
	var signerRule []interface{}
	for _, signerItem := range signer {
		signerRule = append(signerRule, signerItem)
	}

	logs, sub, err := _MultisigWallet.contract.FilterLogs(opts, "PendingTransferConfirmed", transferIdRule, signerRule)
	if err != nil {
		return nil, err
	}
	return &MultisigWalletPendingTransferConfirmedIterator{contract: _MultisigWallet.contract, event: "PendingTransferConfirmed", logs: logs, sub: sub}, nil
}

// WatchPendingTransferConfirmed is a free log subscription operation binding the contract event 0x1c5a4e2f8e0eca8132383b7783986b59dd5c1bd78063bb1262cdac1b8df7507f.
//
// Solidity: event PendingTransferConfirmed(transferId indexed uint256, signer indexed address, signaturesCount uint256)
func (_MultisigWallet *MultisigWalletFilterer) WatchPendingTransferConfirmed(opts *bind.WatchOpts, sink chan<- *MultisigWalletPendingTransferConfirmed, transferId []*big.Int, signer []common.Address) (event.Subscription, error) {

	var transferIdRule []interface{}
	for _, transferIdItem := range transferId {
		transferIdRule = append(transferIdRule, transferIdItem)
	}
	var signerRule []interface{}
	for _, signerItem := range signer {
		signerRule = append(signerRule, signerItem)
	}

	logs, sub, err := _MultisigWallet.contract.WatchLogs(opts, "PendingTransferConfirmed", transferIdRule, signerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultisigWalletPendingTransferConfirmed)
				if err := _MultisigWallet.contract.UnpackLog(event, "PendingTransferConfirmed", log); err != nil {
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

// MultisigWalletPendingTransfersLimitReachedIterator is returned from FilterPendingTransfersLimitReached and is used to iterate over the raw logs and unpacked data for PendingTransfersLimitReached events raised by the MultisigWallet contract.
type MultisigWalletPendingTransfersLimitReachedIterator struct {
	Event *MultisigWalletPendingTransfersLimitReached // Event containing the contract specifics and raw log

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
func (it *MultisigWalletPendingTransfersLimitReachedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultisigWalletPendingTransfersLimitReached)
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
		it.Event = new(MultisigWalletPendingTransfersLimitReached)
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
func (it *MultisigWalletPendingTransfersLimitReachedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultisigWalletPendingTransfersLimitReachedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultisigWalletPendingTransfersLimitReached represents a PendingTransfersLimitReached event raised by the MultisigWallet contract.
type MultisigWalletPendingTransfersLimitReached struct {
	Limit *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterPendingTransfersLimitReached is a free log retrieval operation binding the contract event 0xdaf40e4fc42e3a9470bf3746f72be3cf5dd809cb0d65a03fa4477d25ba008fde.
//
// Solidity: event PendingTransfersLimitReached(limit uint256)
func (_MultisigWallet *MultisigWalletFilterer) FilterPendingTransfersLimitReached(opts *bind.FilterOpts) (*MultisigWalletPendingTransfersLimitReachedIterator, error) {

	logs, sub, err := _MultisigWallet.contract.FilterLogs(opts, "PendingTransfersLimitReached")
	if err != nil {
		return nil, err
	}
	return &MultisigWalletPendingTransfersLimitReachedIterator{contract: _MultisigWallet.contract, event: "PendingTransfersLimitReached", logs: logs, sub: sub}, nil
}

// WatchPendingTransfersLimitReached is a free log subscription operation binding the contract event 0xdaf40e4fc42e3a9470bf3746f72be3cf5dd809cb0d65a03fa4477d25ba008fde.
//
// Solidity: event PendingTransfersLimitReached(limit uint256)
func (_MultisigWallet *MultisigWalletFilterer) WatchPendingTransfersLimitReached(opts *bind.WatchOpts, sink chan<- *MultisigWalletPendingTransfersLimitReached) (event.Subscription, error) {

	logs, sub, err := _MultisigWallet.contract.WatchLogs(opts, "PendingTransfersLimitReached")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultisigWalletPendingTransfersLimitReached)
				if err := _MultisigWallet.contract.UnpackLog(event, "PendingTransfersLimitReached", log); err != nil {
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

// MultisigWalletSignersChangedIterator is returned from FilterSignersChanged and is used to iterate over the raw logs and unpacked data for SignersChanged events raised by the MultisigWallet contract.
type MultisigWalletSignersChangedIterator struct {
	Event *MultisigWalletSignersChanged // Event containing the contract specifics and raw log

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
func (it *MultisigWalletSignersChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultisigWalletSignersChanged)
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
		it.Event = new(MultisigWalletSignersChanged)
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
func (it *MultisigWalletSignersChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultisigWalletSignersChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultisigWalletSignersChanged represents a SignersChanged event raised by the MultisigWallet contract.
type MultisigWalletSignersChanged struct {
	Owner              common.Address
	Signers            []common.Address
	RequiredSignatures uint8
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterSignersChanged is a free log retrieval operation binding the contract event 0xdd9bdfdff1f24d6d3ea9304274494aa100a63b57249f9673a9d7df2d912c2c47.
//
// Solidity: event SignersChanged(owner indexed address, signers address[], requiredSignatures uint8)
func (_MultisigWallet *MultisigWalletFilterer) FilterSignersChanged(opts *bind.FilterOpts, owner []common.Address) (*MultisigWalletSignersChangedIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _MultisigWallet.contract.FilterLogs(opts, "SignersChanged", ownerRule)
	if err != nil {
		return nil, err
	}
	return &MultisigWalletSignersChangedIterator{contract: _MultisigWallet.contract, event: "SignersChanged", logs: logs, sub: sub}, nil
}

// WatchSignersChanged is a free log subscription operation binding the contract event 0xdd9bdfdff1f24d6d3ea9304274494aa100a63b57249f9673a9d7df2d912c2c47.
//
// Solidity: event SignersChanged(owner indexed address, signers address[], requiredSignatures uint8)
func (_MultisigWallet *MultisigWalletFilterer) WatchSignersChanged(opts *bind.WatchOpts, sink chan<- *MultisigWalletSignersChanged, owner []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _MultisigWallet.contract.WatchLogs(opts, "SignersChanged", ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultisigWalletSignersChanged)
				if err := _MultisigWallet.contract.UnpackLog(event, "SignersChanged", log); err != nil {
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

// MultisigWalletTokensTransferredIterator is returned from FilterTokensTransferred and is used to iterate over the raw logs and unpacked data for TokensTransferred events raised by the MultisigWallet contract.
type MultisigWalletTokensTransferredIterator struct {
	Event *MultisigWalletTokensTransferred // Event containing the contract specifics and raw log

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
func (it *MultisigWalletTokensTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultisigWalletTokensTransferred)
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
		it.Event = new(MultisigWalletTokensTransferred)
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
func (it *MultisigWalletTokensTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultisigWalletTokensTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultisigWalletTokensTransferred represents a TokensTransferred event raised by the MultisigWallet contract.
type MultisigWalletTokensTransferred struct {
	TransferId *big.Int
	To         common.Address
	Token      common.Address
	Amount     *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterTokensTransferred is a free log retrieval operation binding the contract event 0x69d21b294c0b932db9bb0dc81144c65ac23a4c078f00589443bb03baa63e3c96.
//
// Solidity: event TokensTransferred(transferId indexed uint256, to indexed address, token indexed address, amount uint256)
func (_MultisigWallet *MultisigWalletFilterer) FilterTokensTransferred(opts *bind.FilterOpts, transferId []*big.Int, to []common.Address, token []common.Address) (*MultisigWalletTokensTransferredIterator, error) {

	var transferIdRule []interface{}
	for _, transferIdItem := range transferId {
		transferIdRule = append(transferIdRule, transferIdItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _MultisigWallet.contract.FilterLogs(opts, "TokensTransferred", transferIdRule, toRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return &MultisigWalletTokensTransferredIterator{contract: _MultisigWallet.contract, event: "TokensTransferred", logs: logs, sub: sub}, nil
}

// WatchTokensTransferred is a free log subscription operation binding the contract event 0x69d21b294c0b932db9bb0dc81144c65ac23a4c078f00589443bb03baa63e3c96.
//
// Solidity: event TokensTransferred(transferId indexed uint256, to indexed address, token indexed address, amount uint256)
func (_MultisigWallet *MultisigWalletFilterer) WatchTokensTransferred(opts *bind.WatchOpts, sink chan<- *MultisigWalletTokensTransferred, transferId []*big.Int, to []common.Address, token []common.Address) (event.Subscription, error) {

	var transferIdRule []interface{}
	for _, transferIdItem := range transferId {
		transferIdRule = append(transferIdRule, transferIdItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _MultisigWallet.contract.WatchLogs(opts, "TokensTransferred", transferIdRule, toRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultisigWalletTokensTransferred)
				if err := _MultisigWallet.contract.UnpackLog(event, "TokensTransferred", log); err != nil {
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

// MultisigWalletUnpauseIterator is returned from FilterUnpause and is used to iterate over the raw logs and unpacked data for Unpause events raised by the MultisigWallet contract.
type MultisigWalletUnpauseIterator struct {
	Event *MultisigWalletUnpause // Event containing the contract specifics and raw log

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
func (it *MultisigWalletUnpauseIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MultisigWalletUnpause)
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
		it.Event = new(MultisigWalletUnpause)
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
func (it *MultisigWalletUnpauseIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MultisigWalletUnpauseIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MultisigWalletUnpause represents a Unpause event raised by the MultisigWallet contract.
type MultisigWalletUnpause struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterUnpause is a free log retrieval operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
//
// Solidity: event Unpause()
func (_MultisigWallet *MultisigWalletFilterer) FilterUnpause(opts *bind.FilterOpts) (*MultisigWalletUnpauseIterator, error) {

	logs, sub, err := _MultisigWallet.contract.FilterLogs(opts, "Unpause")
	if err != nil {
		return nil, err
	}
	return &MultisigWalletUnpauseIterator{contract: _MultisigWallet.contract, event: "Unpause", logs: logs, sub: sub}, nil
}

// WatchUnpause is a free log subscription operation binding the contract event 0x7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b33.
//
// Solidity: event Unpause()
func (_MultisigWallet *MultisigWalletFilterer) WatchUnpause(opts *bind.WatchOpts, sink chan<- *MultisigWalletUnpause) (event.Subscription, error) {

	logs, sub, err := _MultisigWallet.contract.WatchLogs(opts, "Unpause")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MultisigWalletUnpause)
				if err := _MultisigWallet.contract.UnpackLog(event, "Unpause", log); err != nil {
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
