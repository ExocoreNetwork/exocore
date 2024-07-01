// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package assets

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
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// AssetsABI is the input ABI used to generate the binding from.
const AssetsABI = "[{\"type\":\"function\",\"name\":\"depositTo\",\"inputs\":[{\"name\":\"clientChainID\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"assetsAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"stakerAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"opAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"success\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"latestAssetState\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getClientChains\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isRegisteredClientChain\",\"inputs\":[{\"name\":\"clientChainID\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"success\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"isRegistered\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"registerOrUpdateClientChain\",\"inputs\":[{\"name\":\"clientChainID\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"addressLength\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"metaInfo\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"signatureType\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"success\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"updated\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"registerOrUpdateTokens\",\"inputs\":[{\"name\":\"clientChainID\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"token\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"decimals\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"tvlLimit\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"metaData\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"success\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"updated\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawPrincipal\",\"inputs\":[{\"name\":\"clientChainID\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"assetsAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"withdrawAddress\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"opAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"success\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"latestAssetState\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"}]"

// Assets is an auto generated Go binding around an Ethereum contract.
type Assets struct {
	AssetsCaller     // Read-only binding to the contract
	AssetsTransactor // Write-only binding to the contract
	AssetsFilterer   // Log filterer for contract events
}

// AssetsCaller is an auto generated read-only Go binding around an Ethereum contract.
type AssetsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AssetsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AssetsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AssetsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AssetsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AssetsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AssetsSession struct {
	Contract     *Assets           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AssetsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AssetsCallerSession struct {
	Contract *AssetsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// AssetsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AssetsTransactorSession struct {
	Contract     *AssetsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AssetsRaw is an auto generated low-level Go binding around an Ethereum contract.
type AssetsRaw struct {
	Contract *Assets // Generic contract binding to access the raw methods on
}

// AssetsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AssetsCallerRaw struct {
	Contract *AssetsCaller // Generic read-only contract binding to access the raw methods on
}

// AssetsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AssetsTransactorRaw struct {
	Contract *AssetsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAssets creates a new instance of Assets, bound to a specific deployed contract.
func NewAssets(address common.Address, backend bind.ContractBackend) (*Assets, error) {
	contract, err := bindAssets(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Assets{AssetsCaller: AssetsCaller{contract: contract}, AssetsTransactor: AssetsTransactor{contract: contract}, AssetsFilterer: AssetsFilterer{contract: contract}}, nil
}

// NewAssetsCaller creates a new read-only instance of Assets, bound to a specific deployed contract.
func NewAssetsCaller(address common.Address, caller bind.ContractCaller) (*AssetsCaller, error) {
	contract, err := bindAssets(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AssetsCaller{contract: contract}, nil
}

// NewAssetsTransactor creates a new write-only instance of Assets, bound to a specific deployed contract.
func NewAssetsTransactor(address common.Address, transactor bind.ContractTransactor) (*AssetsTransactor, error) {
	contract, err := bindAssets(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AssetsTransactor{contract: contract}, nil
}

// NewAssetsFilterer creates a new log filterer instance of Assets, bound to a specific deployed contract.
func NewAssetsFilterer(address common.Address, filterer bind.ContractFilterer) (*AssetsFilterer, error) {
	contract, err := bindAssets(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AssetsFilterer{contract: contract}, nil
}

// bindAssets binds a generic wrapper to an already deployed contract.
func bindAssets(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AssetsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Assets *AssetsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Assets.Contract.AssetsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Assets *AssetsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Assets.Contract.AssetsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Assets *AssetsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Assets.Contract.AssetsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Assets *AssetsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Assets.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Assets *AssetsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Assets.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Assets *AssetsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Assets.Contract.contract.Transact(opts, method, params...)
}

// GetClientChains is a free data retrieval call binding the contract method 0x41a3745b.
//
// Solidity: function getClientChains() view returns(bool, uint32[])
func (_Assets *AssetsCaller) GetClientChains(opts *bind.CallOpts) (bool, []uint32, error) {
	var out []interface{}
	err := _Assets.contract.Call(opts, &out, "getClientChains")

	if err != nil {
		return *new(bool), *new([]uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new([]uint32)).(*[]uint32)

	return out0, out1, err

}

// GetClientChains is a free data retrieval call binding the contract method 0x41a3745b.
//
// Solidity: function getClientChains() view returns(bool, uint32[])
func (_Assets *AssetsSession) GetClientChains() (bool, []uint32, error) {
	return _Assets.Contract.GetClientChains(&_Assets.CallOpts)
}

// GetClientChains is a free data retrieval call binding the contract method 0x41a3745b.
//
// Solidity: function getClientChains() view returns(bool, uint32[])
func (_Assets *AssetsCallerSession) GetClientChains() (bool, []uint32, error) {
	return _Assets.Contract.GetClientChains(&_Assets.CallOpts)
}

// IsRegisteredClientChain is a free data retrieval call binding the contract method 0x6b67d7f7.
//
// Solidity: function isRegisteredClientChain(uint32 clientChainID) view returns(bool success, bool isRegistered)
func (_Assets *AssetsCaller) IsRegisteredClientChain(opts *bind.CallOpts, clientChainID uint32) (struct {
	Success      bool
	IsRegistered bool
}, error) {
	var out []interface{}
	err := _Assets.contract.Call(opts, &out, "isRegisteredClientChain", clientChainID)

	outstruct := new(struct {
		Success      bool
		IsRegistered bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Success = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.IsRegistered = *abi.ConvertType(out[1], new(bool)).(*bool)

	return *outstruct, err

}

// IsRegisteredClientChain is a free data retrieval call binding the contract method 0x6b67d7f7.
//
// Solidity: function isRegisteredClientChain(uint32 clientChainID) view returns(bool success, bool isRegistered)
func (_Assets *AssetsSession) IsRegisteredClientChain(clientChainID uint32) (struct {
	Success      bool
	IsRegistered bool
}, error) {
	return _Assets.Contract.IsRegisteredClientChain(&_Assets.CallOpts, clientChainID)
}

// IsRegisteredClientChain is a free data retrieval call binding the contract method 0x6b67d7f7.
//
// Solidity: function isRegisteredClientChain(uint32 clientChainID) view returns(bool success, bool isRegistered)
func (_Assets *AssetsCallerSession) IsRegisteredClientChain(clientChainID uint32) (struct {
	Success      bool
	IsRegistered bool
}, error) {
	return _Assets.Contract.IsRegisteredClientChain(&_Assets.CallOpts, clientChainID)
}

// DepositTo is a paid mutator transaction binding the contract method 0xfc5b72e2.
//
// Solidity: function depositTo(uint32 clientChainID, bytes assetsAddress, bytes stakerAddress, uint256 opAmount) returns(bool success, uint256 latestAssetState)
func (_Assets *AssetsTransactor) DepositTo(opts *bind.TransactOpts, clientChainID uint32, assetsAddress []byte, stakerAddress []byte, opAmount *big.Int) (*types.Transaction, error) {
	return _Assets.contract.Transact(opts, "depositTo", clientChainID, assetsAddress, stakerAddress, opAmount)
}

// DepositTo is a paid mutator transaction binding the contract method 0xfc5b72e2.
//
// Solidity: function depositTo(uint32 clientChainID, bytes assetsAddress, bytes stakerAddress, uint256 opAmount) returns(bool success, uint256 latestAssetState)
func (_Assets *AssetsSession) DepositTo(clientChainID uint32, assetsAddress []byte, stakerAddress []byte, opAmount *big.Int) (*types.Transaction, error) {
	return _Assets.Contract.DepositTo(&_Assets.TransactOpts, clientChainID, assetsAddress, stakerAddress, opAmount)
}

// DepositTo is a paid mutator transaction binding the contract method 0xfc5b72e2.
//
// Solidity: function depositTo(uint32 clientChainID, bytes assetsAddress, bytes stakerAddress, uint256 opAmount) returns(bool success, uint256 latestAssetState)
func (_Assets *AssetsTransactorSession) DepositTo(clientChainID uint32, assetsAddress []byte, stakerAddress []byte, opAmount *big.Int) (*types.Transaction, error) {
	return _Assets.Contract.DepositTo(&_Assets.TransactOpts, clientChainID, assetsAddress, stakerAddress, opAmount)
}

// RegisterOrUpdateClientChain is a paid mutator transaction binding the contract method 0x1b315b52.
//
// Solidity: function registerOrUpdateClientChain(uint32 clientChainID, uint8 addressLength, string name, string metaInfo, string signatureType) returns(bool success, bool updated)
func (_Assets *AssetsTransactor) RegisterOrUpdateClientChain(opts *bind.TransactOpts, clientChainID uint32, addressLength uint8, name string, metaInfo string, signatureType string) (*types.Transaction, error) {
	return _Assets.contract.Transact(opts, "registerOrUpdateClientChain", clientChainID, addressLength, name, metaInfo, signatureType)
}

// RegisterOrUpdateClientChain is a paid mutator transaction binding the contract method 0x1b315b52.
//
// Solidity: function registerOrUpdateClientChain(uint32 clientChainID, uint8 addressLength, string name, string metaInfo, string signatureType) returns(bool success, bool updated)
func (_Assets *AssetsSession) RegisterOrUpdateClientChain(clientChainID uint32, addressLength uint8, name string, metaInfo string, signatureType string) (*types.Transaction, error) {
	return _Assets.Contract.RegisterOrUpdateClientChain(&_Assets.TransactOpts, clientChainID, addressLength, name, metaInfo, signatureType)
}

// RegisterOrUpdateClientChain is a paid mutator transaction binding the contract method 0x1b315b52.
//
// Solidity: function registerOrUpdateClientChain(uint32 clientChainID, uint8 addressLength, string name, string metaInfo, string signatureType) returns(bool success, bool updated)
func (_Assets *AssetsTransactorSession) RegisterOrUpdateClientChain(clientChainID uint32, addressLength uint8, name string, metaInfo string, signatureType string) (*types.Transaction, error) {
	return _Assets.Contract.RegisterOrUpdateClientChain(&_Assets.TransactOpts, clientChainID, addressLength, name, metaInfo, signatureType)
}

// RegisterOrUpdateTokens is a paid mutator transaction binding the contract method 0xbf49cb71.
//
// Solidity: function registerOrUpdateTokens(uint32 clientChainID, bytes token, uint8 decimals, uint256 tvlLimit, string name, string metaData) returns(bool success, bool updated)
func (_Assets *AssetsTransactor) RegisterOrUpdateTokens(opts *bind.TransactOpts, clientChainID uint32, token []byte, decimals uint8, tvlLimit *big.Int, name string, metaData string) (*types.Transaction, error) {
	return _Assets.contract.Transact(opts, "registerOrUpdateTokens", clientChainID, token, decimals, tvlLimit, name, metaData)
}

// RegisterOrUpdateTokens is a paid mutator transaction binding the contract method 0xbf49cb71.
//
// Solidity: function registerOrUpdateTokens(uint32 clientChainID, bytes token, uint8 decimals, uint256 tvlLimit, string name, string metaData) returns(bool success, bool updated)
func (_Assets *AssetsSession) RegisterOrUpdateTokens(clientChainID uint32, token []byte, decimals uint8, tvlLimit *big.Int, name string, metaData string) (*types.Transaction, error) {
	return _Assets.Contract.RegisterOrUpdateTokens(&_Assets.TransactOpts, clientChainID, token, decimals, tvlLimit, name, metaData)
}

// RegisterOrUpdateTokens is a paid mutator transaction binding the contract method 0xbf49cb71.
//
// Solidity: function registerOrUpdateTokens(uint32 clientChainID, bytes token, uint8 decimals, uint256 tvlLimit, string name, string metaData) returns(bool success, bool updated)
func (_Assets *AssetsTransactorSession) RegisterOrUpdateTokens(clientChainID uint32, token []byte, decimals uint8, tvlLimit *big.Int, name string, metaData string) (*types.Transaction, error) {
	return _Assets.Contract.RegisterOrUpdateTokens(&_Assets.TransactOpts, clientChainID, token, decimals, tvlLimit, name, metaData)
}

// WithdrawPrincipal is a paid mutator transaction binding the contract method 0x6f233c6b.
//
// Solidity: function withdrawPrincipal(uint32 clientChainID, bytes assetsAddress, bytes withdrawAddress, uint256 opAmount) returns(bool success, uint256 latestAssetState)
func (_Assets *AssetsTransactor) WithdrawPrincipal(opts *bind.TransactOpts, clientChainID uint32, assetsAddress []byte, withdrawAddress []byte, opAmount *big.Int) (*types.Transaction, error) {
	return _Assets.contract.Transact(opts, "withdrawPrincipal", clientChainID, assetsAddress, withdrawAddress, opAmount)
}

// WithdrawPrincipal is a paid mutator transaction binding the contract method 0x6f233c6b.
//
// Solidity: function withdrawPrincipal(uint32 clientChainID, bytes assetsAddress, bytes withdrawAddress, uint256 opAmount) returns(bool success, uint256 latestAssetState)
func (_Assets *AssetsSession) WithdrawPrincipal(clientChainID uint32, assetsAddress []byte, withdrawAddress []byte, opAmount *big.Int) (*types.Transaction, error) {
	return _Assets.Contract.WithdrawPrincipal(&_Assets.TransactOpts, clientChainID, assetsAddress, withdrawAddress, opAmount)
}

// WithdrawPrincipal is a paid mutator transaction binding the contract method 0x6f233c6b.
//
// Solidity: function withdrawPrincipal(uint32 clientChainID, bytes assetsAddress, bytes withdrawAddress, uint256 opAmount) returns(bool success, uint256 latestAssetState)
func (_Assets *AssetsTransactorSession) WithdrawPrincipal(clientChainID uint32, assetsAddress []byte, withdrawAddress []byte, opAmount *big.Int) (*types.Transaction, error) {
	return _Assets.Contract.WithdrawPrincipal(&_Assets.TransactOpts, clientChainID, assetsAddress, withdrawAddress, opAmount)
}
