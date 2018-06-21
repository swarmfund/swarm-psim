// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ethcontracts

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

// MultisigWalletBin is the compiled bytecode used for deploying new contracts.
const MultisigWalletBin = `0x608060405260008060146101000a81548160ff021916908315150217905550600060055560405162002346380380620023468339810180604052810190808051820192919060200180519060200190929190505050336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550620000af828262000112640100000000026401000000009004565b620000ca82826200015a640100000000026401000000009004565b336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050506200041c565b600182511180156200012657506008825111155b15156200013257600080fd5b60008160ff161180156200014a575081518160ff1611155b15156200015657600080fd5b5050565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16141515620001b857600080fd5b600090505b82518160ff161015620003fc576000838260ff16815181101515620001de57fe5b9060200190602002015173ffffffffffffffffffffffffffffffffffffffff16141515156200020c57600080fd5b6000151560026000858460ff168151811015156200022657fe5b9060200190602002015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff1615151415156200028757600080fd5b600160026000858460ff168151811015156200029f57fe5b9060200190602002015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055508060036000858460ff168151811015156200031157fe5b9060200190602002015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff021916908360ff1602179055506001838260ff168151811015156200038157fe5b9060200190602002015190806001815401808255809150509060018203906000526020600020016000909192909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550508080600101915050620001bd565b81600460006101000a81548160ff021916908360ff160217905550505050565b611f1a806200042c6000396000f3006080604052600436106100db576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806301c0954b1461013557806316f4d52d146101b65780631878a5c9146101e15780632c48e7db146102545780633f4ba83a14610281578063413d3cf114610298578063448878fc146102f957806359ecd657146103655780635c975abb146103905780638456cb59146103bf5780638da5cb5b146103d657806394cf795e1461042d578063ccd9399814610499578063f2fde38b146104ca578063f91a7b891461050d575b6000341115610133573373ffffffffffffffffffffffffffffffffffffffff167fe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c346040518082815260200191505060405180910390a25b005b34801561014157600080fd5b506101a0600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001909291905050506105dd565b6040518082815260200191505060405180910390f35b3480156101c257600080fd5b506101cb610761565b6040518082815260200191505060405180910390f35b3480156101ed57600080fd5b5061025260048036038101908080359060200190820180359060200190808060200260200160405190810160405280939291908181526020018383602002808284378201915050505050509192919290803560ff169060200190929190505050610766565b005b34801561026057600080fd5b5061027f60048036038101908080359060200190929190505050610af2565b005b34801561028d57600080fd5b50610296610ee4565b005b3480156102a457600080fd5b506102e3600480360381019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050610fa2565b6040518082815260200191505060405180910390f35b34801561030557600080fd5b5061030e6110e9565b6040518080602001828103825283818151815260200191508051906020019060200280838360005b83811015610351578082015181840152602081019050610336565b505050509050019250505060405180910390f35b34801561037157600080fd5b5061037a611141565b6040518082815260200191505060405180910390f35b34801561039c57600080fd5b506103a5611146565b604051808215151515815260200191505060405180910390f35b3480156103cb57600080fd5b506103d4611159565b005b3480156103e257600080fd5b506103eb611219565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561043957600080fd5b5061044261123e565b6040518080602001828103825283818151815260200191508051906020019060200280838360005b8381101561048557808201518184015260208101905061046a565b505050509050019250505060405180910390f35b3480156104a557600080fd5b506104ae6112cc565b604051808260ff1660ff16815260200191505060405180910390f35b3480156104d657600080fd5b5061050b600480360381019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506112e3565b005b34801561051957600080fd5b5061053860048036038101908080359060200190929190505050611438565b604051808781526020018673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018581526020018473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018360028111156105b857fe5b60ff1681526020018260ff1660ff168152602001965050505050505060405180910390f35b600080600260003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16151561063857600080fd5b601460068054905010151561064c57600080fd5b600060149054906101000a900460ff1615151561066857600080fd5b60008573ffffffffffffffffffffffffffffffffffffffff161415151561068e57600080fd5b60008473ffffffffffffffffffffffffffffffffffffffff16141515156106b457600080fd5b6000831115156106c357600080fd5b6106d060028686866114e5565b90508373ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167f670c51d7c9ae60ee716af92c18abb6284cc1e377e7b32ce4f710d2ffd3a8e17b8685604051808381526020018281526020019250505060405180910390a4809150509392505050565b601481565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161415156107c357600080fd5b600060149054906101000a900460ff1615156107de57600080fd5b6107e8838361175f565b600060068054905011156108cb57600090505b6006805490508110156108ca576007600060068381548110151561081b57fe5b906000526020600020015481526020019081526020016000206000808201600090556001820160006101000a81549073ffffffffffffffffffffffffffffffffffffffff021916905560028201600090556003820160006101000a81549073ffffffffffffffffffffffffffffffffffffffff02191690556003820160146101000a81549060ff02191690556003820160156101000a81549060ff0219169055505080806001019150506107fb565b5b600660006108d99190611dea565b600090505b600180549050811015610a0857600260006001838154811015156108fe57fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81549060ff02191690556003600060018381548110151561098657fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81549060ff021916905580806001019150506108de565b60016000610a169190611e0b565b610a2083836117a3565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167fdd9bdfdff1f24d6d3ea9304274494aa100a63b57249f9673a9d7df2d912c2c4784600460009054906101000a900460ff1660405180806020018360ff1660ff168152602001828103825284818151815260200191508051906020019060200280838360005b83811015610ad9578082015181840152602081019050610abe565b50505050905001935050505060405180910390a2505050565b600080600080600080600260003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff161515610b5357600080fd5b600060149054906101000a900460ff16151515610b6f57600080fd5b60076000888152602001908152602001600020955060006002811115610b9157fe5b8660030160149054906101000a900460ff166002811115610bae57fe5b14151515610bbb57600080fd5b610bc486611a5b565b610bcd86611b45565b94503373ffffffffffffffffffffffffffffffffffffffff16877f1c5a4e2f8e0eca8132383b7783986b59dd5c1bd78063bb1262cdac1b8df7507f87604051808260ff16815260200191505060405180910390a3600460009054906101000a900460ff1660ff168560ff16101515610edb578560030160149054906101000a900460ff1693508560010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169250856002015491508560030160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050610cb186611bb0565b60016002811115610cbe57fe5b846002811115610cca57fe5b1415610d6b578273ffffffffffffffffffffffffffffffffffffffff166108fc839081150290604051600060405180830381858888f19350505050158015610d16573d6000803e3d6000fd5b508273ffffffffffffffffffffffffffffffffffffffff16877f8d6a64c40ae95e2d8ba6157294f8d842d794fc27c5023b5cc12193e57fa4237b846040518082815260200191505060405180910390a3610eda565b600280811115610d7757fe5b846002811115610d8357fe5b1415610ed9578073ffffffffffffffffffffffffffffffffffffffff1663a9059cbb84846040518363ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200182815260200192505050602060405180830381600087803b158015610e2c57600080fd5b505af1158015610e40573d6000803e3d6000fd5b505050506040513d6020811015610e5657600080fd5b81019080805190602001909291905050501515610e7257600080fd5b8073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16887f69d21b294c0b932db9bb0dc81144c65ac23a4c078f00589443bb03baa63e3c96856040518082815260200191505060405180910390a45b5b5b50505050505050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16141515610f3f57600080fd5b600060149054906101000a900460ff161515610f5a57600080fd5b60008060146101000a81548160ff0219169083151502179055507f7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b3360405160405180910390a1565b600080600260003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff161515610ffd57600080fd5b601460068054905010151561101157600080fd5b600060149054906101000a900460ff1615151561102d57600080fd5b60008473ffffffffffffffffffffffffffffffffffffffff161415151561105357600080fd5b60008311151561106257600080fd5b6110706001856000866114e5565b90508373ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff167fe7f0c7c6807ccf602c8eeedddf80ca86ab53f66740779a25c3c93abe0fc809d18584604051808381526020018281526020019250505060405180910390a38091505092915050565b6060600680548060200260200160405190810160405280929190818152602001828054801561113757602002820191906000526020600020905b815481526020019060010190808311611123575b5050505050905090565b600881565b600060149054906101000a900460ff1681565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161415156111b457600080fd5b600060149054906101000a900460ff161515156111d057600080fd5b6001600060146101000a81548160ff0219169083151502179055507f6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff62560405160405180910390a1565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b606060018054806020026020016040519081016040528092919081815260200182805480156112c257602002820191906000526020600020905b8160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019060010190808311611278575b5050505050905090565b6000600460009054906101000a900460ff16905090565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614151561133e57600080fd5b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415151561137a57600080fd5b8073ffffffffffffffffffffffffffffffffffffffff166000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a3806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b600080600080600080600080600760008a8152602001908152602001600020915061146282611b45565b905081600001548260010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1683600201548460030160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff168560030160149054906101000a900460ff1685829250975097509750975097509750505091939550919395565b6000806114f0611e2c565b600260003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16151561154857600080fd5b611550611d9d565b915060c0604051908101604052808381526020018773ffffffffffffffffffffffffffffffffffffffff1681526020018581526020018673ffffffffffffffffffffffffffffffffffffffff1681526020018860028111156115ae57fe5b8152602001600060ff16815250905060068160000151908060018154018082558091505090600182039060005260206000200160009091929091909150555080600760008481526020019081526020016000206000820151816000015560208201518160010160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506040820151816002015560608201518160030160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060808201518160030160146101000a81548160ff021916908360028111156116c757fe5b021790555060a08201518160030160156101000a81548160ff021916908360ff16021790555090505061170b60076000848152602001908152602001600020611a5b565b60146006805490501415611752577fdaf40e4fc42e3a9470bf3746f72be3cf5dd809cb0d65a03fa4477d25ba008fde60146040518082815260200191505060405180910390a15b8192505050949350505050565b6001825111801561177257506008825111155b151561177d57600080fd5b60008160ff16118015611794575081518160ff1611155b151561179f57600080fd5b5050565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614151561180057600080fd5b600090505b82518160ff161015611a3b576000838260ff1681518110151561182457fe5b9060200190602002015173ffffffffffffffffffffffffffffffffffffffff161415151561185157600080fd5b6000151560026000858460ff1681518110151561186a57fe5b9060200190602002015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff1615151415156118ca57600080fd5b600160026000858460ff168151811015156118e157fe5b9060200190602002015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055508060036000858460ff1681518110151561195257fe5b9060200190602002015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff021916908360ff1602179055506001838260ff168151811015156119c157fe5b9060200190602002015190806001815401808255809150509060018203906000526020600020016000909192909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550508080600101915050611805565b81600460006101000a81548160ff021916908360ff160217905550505050565b6000600260003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff161515611ab557600080fd5b600360003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff1660ff16600160ff169060020a029050808260030160159054906101000a900460ff16178260030160156101000a81548160ff021916908360ff1602179055505050565b6000806000806000925060019150600090505b60088160ff161015611ba55760008160ff168360ff169060020a028660030160159054906101000a900460ff161660ff161115611b985782806001019350505b8080600101915050611b58565b829350505050919050565b600080600260003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff161515611c0b57600080fd5b60009150600090505b600680549050811015611cd8578115611c8357600681815481101515611c3657fe5b9060005260206000200154600660018303815481101515611c5357fe5b9060005260206000200181905550600681815481101515611c7057fe5b9060005260206000200160009055611ccb565b8260000154600682815481101515611c9757fe5b90600052602060002001541415611cca57600681815481101515611cb757fe5b9060005260206000200160009055600191505b5b8080600101915050611c14565b600160068054905003600681611cee9190611e9d565b50811515611cf857fe5b60076000846000015481526020019081526020016000206000808201600090556001820160006101000a81549073ffffffffffffffffffffffffffffffffffffffff021916905560028201600090556003820160006101000a81549073ffffffffffffffffffffffffffffffffffffffff02191690556003820160146101000a81549060ff02191690556003820160156101000a81549060ff02191690555050505050565b600060056000815480929190600101919050555043600554600036604051808581526020018481526020018383808284378201915050945050505050604051809103902060019004905090565b5080546000825590600052602060002090810190611e089190611ec9565b50565b5080546000825590600052602060002090810190611e299190611ec9565b50565b60c06040519081016040528060008152602001600073ffffffffffffffffffffffffffffffffffffffff16815260200160008152602001600073ffffffffffffffffffffffffffffffffffffffff16815260200160006002811115611e8d57fe5b8152602001600060ff1681525090565b815481835581811115611ec457818360005260206000209182019101611ec39190611ec9565b5b505050565b611eeb91905b80821115611ee7576000816000905550600101611ecf565b5090565b905600a165627a7a7230582059ffcba22194bea374a766229ef13061eaeb801a1fb63200befd5f9ab36008d70029`

// DeployMultisigWallet deploys a new Ethereum contract, binding an instance of MultisigWallet to it.
func DeployMultisigWallet(auth *bind.TransactOpts, backend bind.ContractBackend, _signers []common.Address, _requiredSignatures uint8) (common.Address, *types.Transaction, *MultisigWallet, error) {
	parsed, err := abi.JSON(strings.NewReader(MultisigWalletABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(MultisigWalletBin), backend, _signers, _requiredSignatures)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MultisigWallet{MultisigWalletCaller: MultisigWalletCaller{contract: contract}, MultisigWalletTransactor: MultisigWalletTransactor{contract: contract}, MultisigWalletFilterer: MultisigWalletFilterer{contract: contract}}, nil
}

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
