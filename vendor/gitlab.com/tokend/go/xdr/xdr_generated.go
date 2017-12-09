// Package xdr is generated from:
//
//  xdr/raw/Stellar-ledger-entries-account-limits.x
//  xdr/raw/Stellar-ledger-entries-account-type-limits.x
//  xdr/raw/Stellar-ledger-entries-account.x
//  xdr/raw/Stellar-ledger-entries-asset-pair.x
//  xdr/raw/Stellar-ledger-entries-asset.x
//  xdr/raw/Stellar-ledger-entries-balance.x
//  xdr/raw/Stellar-ledger-entries-coins-emission-request.x
//  xdr/raw/Stellar-ledger-entries-exchange.x
//  xdr/raw/Stellar-ledger-entries-fee.x
//  xdr/raw/Stellar-ledger-entries-invoice.x
//  xdr/raw/Stellar-ledger-entries-offer.x
//  xdr/raw/Stellar-ledger-entries-payment-request.x
//  xdr/raw/Stellar-ledger-entries-payment.x
//  xdr/raw/Stellar-ledger-entries-statistics.x
//  xdr/raw/Stellar-ledger-entries.x
//  xdr/raw/Stellar-ledger.x
//  xdr/raw/Stellar-operation-create-account.x
//  xdr/raw/Stellar-operation-demurrage.x
//  xdr/raw/Stellar-operation-direct-debit.x
//  xdr/raw/Stellar-operation-forfeit.x
//  xdr/raw/Stellar-operation-manage-account.x
//  xdr/raw/Stellar-operation-manage-asset-pair.x
//  xdr/raw/Stellar-operation-manage-asset.x
//  xdr/raw/Stellar-operation-manage-balance.x
//  xdr/raw/Stellar-operation-manage-coins-emission-request.x
//  xdr/raw/Stellar-operation-manage-forfeit-request.x
//  xdr/raw/Stellar-operation-manage-invoice.x
//  xdr/raw/Stellar-operation-manage-offer.x
//  xdr/raw/Stellar-operation-payment.x
//  xdr/raw/Stellar-operation-recover.x
//  xdr/raw/Stellar-operation-review-coins-emission-request.x
//  xdr/raw/Stellar-operation-review-payment-request.x
//  xdr/raw/Stellar-operation-set-fees.x
//  xdr/raw/Stellar-operation-set-limits.x
//  xdr/raw/Stellar-operation-set-options.x
//  xdr/raw/Stellar-operation-upload-preemissions.x
//  xdr/raw/Stellar-overlay.x
//  xdr/raw/Stellar-SCP.x
//  xdr/raw/Stellar-transaction.x
//  xdr/raw/Stellar-types.x
//
// DO NOT EDIT or your changes may be overwritten
package xdr

import (
	"fmt"
	"io"

	"github.com/nullstyle/go-xdr/xdr3"
)

// Unmarshal reads an xdr element from `r` into `v`.
func Unmarshal(r io.Reader, v interface{}) (int, error) {
	// delegate to xdr package's Unmarshal
	return xdr.Unmarshal(r, v)
}

// Marshal writes an xdr element `v` into `w`.
func Marshal(w io.Writer, v interface{}) (int, error) {
	// delegate to xdr package's Marshal
	return xdr.Marshal(w, v)
}

// AccountLimitsEntryExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type AccountLimitsEntryExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u AccountLimitsEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of AccountLimitsEntryExt
func (u AccountLimitsEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewAccountLimitsEntryExt creates a new  AccountLimitsEntryExt.
func NewAccountLimitsEntryExt(v LedgerVersion, value interface{}) (result AccountLimitsEntryExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// AccountLimitsEntry is an XDR Struct defines as:
//
//   struct AccountLimitsEntry
//    {
//        AccountID accountID;
//        Limits limits;
//
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type AccountLimitsEntry struct {
	AccountId AccountId             `json:"accountID,omitempty"`
	Limits    Limits                `json:"limits,omitempty"`
	Ext       AccountLimitsEntryExt `json:"ext,omitempty"`
}

// AccountTypeLimitsEntryExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type AccountTypeLimitsEntryExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u AccountTypeLimitsEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of AccountTypeLimitsEntryExt
func (u AccountTypeLimitsEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewAccountTypeLimitsEntryExt creates a new  AccountTypeLimitsEntryExt.
func NewAccountTypeLimitsEntryExt(v LedgerVersion, value interface{}) (result AccountTypeLimitsEntryExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// AccountTypeLimitsEntry is an XDR Struct defines as:
//
//   struct AccountTypeLimitsEntry
//    {
//    	AccountType accountType;
//        Limits limits;
//
//    	// reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type AccountTypeLimitsEntry struct {
	AccountType AccountType               `json:"accountType,omitempty"`
	Limits      Limits                    `json:"limits,omitempty"`
	Ext         AccountTypeLimitsEntryExt `json:"ext,omitempty"`
}

// SignerType is an XDR Enum defines as:
//
//   enum SignerType
//    {
//    	SIGNER_READER = 1,                  // can only read data from API and Horizon
//    	SIGNER_NOT_VERIFIED_ACC_MANAGER = 2,// can manage not verified account and block/unblock general
//    	SIGNER_GENERAL_ACC_MANAGER = 4,     // allowed to create account, block/unblock, change limits for particular general account
//        SIGNER_EXCHANGE_ACC_MANAGER = 8,     // allowed to create exchange account, block/unblock
//    	SIGNER_DEMURRAGE_OPERATOR = 16, // allowed to initiate demurrage
//    	SIGNER_DIRECT_DEBIT_OPERATOR = 32, // allowed to perform direct debit operation
//    	SIGNER_FORFEIT_OPERATOR = 64, // allowed to perform forfeit opertion over account
//    	SIGNER_ASSET_MANAGER = 128, // allowed to create assets/asset pairs and update policies, set fees
//    	SIGNER_ASSET_RATE_MANAGER = 256, // allowed to set physical asset price
//    	SIGNER_BALANCE_MANAGER = 512, // allowed to create balances, spend assets from balances
//    	SIGNER_EMISSION_MANAGER = 1024, // allowed to make emission requests, review emission, upload preemission
//    	SIGNER_INVOICE_MANAGER = 2048, // allowed to create payment requests to other accounts
//    	SIGNER_PAYMENT_OPERATOR = 4096, // allowed to review payment requests
//    	SIGNER_LIMITS_MANAGER = 8192, // allowed to change limits
//    	SIGNER_ACCOUNT_MANAGER = 16384, // allowed to add/delete signers and trust
//    	SIGNER_COMMISSION_BALANCE_MANAGER  = 32768,// allowed to spend from commission balances
//    	SIGNER_OPERATIONAL_BALANCE_MANAGER = 65536,// allowed to spend from operational balances
//    	SIGNER_STORAGE_FEE_BALANCE_MANAGER = 131072// allowed to spend from storage fee balances
//    };
//
type SignerType int32

const (
	SignerTypeSignerReader                    SignerType = 1
	SignerTypeSignerNotVerifiedAccManager     SignerType = 2
	SignerTypeSignerGeneralAccManager         SignerType = 4
	SignerTypeSignerExchangeAccManager        SignerType = 8
	SignerTypeSignerDemurrageOperator         SignerType = 16
	SignerTypeSignerDirectDebitOperator       SignerType = 32
	SignerTypeSignerForfeitOperator           SignerType = 64
	SignerTypeSignerAssetManager              SignerType = 128
	SignerTypeSignerAssetRateManager          SignerType = 256
	SignerTypeSignerBalanceManager            SignerType = 512
	SignerTypeSignerEmissionManager           SignerType = 1024
	SignerTypeSignerInvoiceManager            SignerType = 2048
	SignerTypeSignerPaymentOperator           SignerType = 4096
	SignerTypeSignerLimitsManager             SignerType = 8192
	SignerTypeSignerAccountManager            SignerType = 16384
	SignerTypeSignerCommissionBalanceManager  SignerType = 32768
	SignerTypeSignerOperationalBalanceManager SignerType = 65536
	SignerTypeSignerStorageFeeBalanceManager  SignerType = 131072
)

var SignerTypeAll = []SignerType{
	SignerTypeSignerReader,
	SignerTypeSignerNotVerifiedAccManager,
	SignerTypeSignerGeneralAccManager,
	SignerTypeSignerExchangeAccManager,
	SignerTypeSignerDemurrageOperator,
	SignerTypeSignerDirectDebitOperator,
	SignerTypeSignerForfeitOperator,
	SignerTypeSignerAssetManager,
	SignerTypeSignerAssetRateManager,
	SignerTypeSignerBalanceManager,
	SignerTypeSignerEmissionManager,
	SignerTypeSignerInvoiceManager,
	SignerTypeSignerPaymentOperator,
	SignerTypeSignerLimitsManager,
	SignerTypeSignerAccountManager,
	SignerTypeSignerCommissionBalanceManager,
	SignerTypeSignerOperationalBalanceManager,
	SignerTypeSignerStorageFeeBalanceManager,
}

var signerTypeMap = map[int32]string{
	1:      "SignerTypeSignerReader",
	2:      "SignerTypeSignerNotVerifiedAccManager",
	4:      "SignerTypeSignerGeneralAccManager",
	8:      "SignerTypeSignerExchangeAccManager",
	16:     "SignerTypeSignerDemurrageOperator",
	32:     "SignerTypeSignerDirectDebitOperator",
	64:     "SignerTypeSignerForfeitOperator",
	128:    "SignerTypeSignerAssetManager",
	256:    "SignerTypeSignerAssetRateManager",
	512:    "SignerTypeSignerBalanceManager",
	1024:   "SignerTypeSignerEmissionManager",
	2048:   "SignerTypeSignerInvoiceManager",
	4096:   "SignerTypeSignerPaymentOperator",
	8192:   "SignerTypeSignerLimitsManager",
	16384:  "SignerTypeSignerAccountManager",
	32768:  "SignerTypeSignerCommissionBalanceManager",
	65536:  "SignerTypeSignerOperationalBalanceManager",
	131072: "SignerTypeSignerStorageFeeBalanceManager",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for SignerType
func (e SignerType) ValidEnum(v int32) bool {
	_, ok := signerTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e SignerType) String() string {
	name, _ := signerTypeMap[int32(e)]
	return name
}

func (e SignerType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// SignerExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//    	case SIGNER_NAME:
//    		string256 name;
//        }
//
type SignerExt struct {
	V    LedgerVersion `json:"v,omitempty"`
	Name *String256    `json:"name,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u SignerExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of SignerExt
func (u SignerExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	case LedgerVersionSignerName:
		return "Name", true
	}
	return "-", false
}

// NewSignerExt creates a new  SignerExt.
func NewSignerExt(v LedgerVersion, value interface{}) (result SignerExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	case LedgerVersionSignerName:
		tv, ok := value.(String256)
		if !ok {
			err = fmt.Errorf("invalid value, must be String256")
			return
		}
		result.Name = &tv
	}
	return
}

// MustName retrieves the Name value from the union,
// panicing if the value is not set.
func (u SignerExt) MustName() String256 {
	val, ok := u.GetName()

	if !ok {
		panic("arm Name is not set")
	}

	return val
}

// GetName retrieves the Name value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u SignerExt) GetName() (result String256, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.V))

	if armName == "Name" {
		result = *u.Name
		ok = true
	}

	return
}

// Signer is an XDR Struct defines as:
//
//   struct Signer
//    {
//        AccountID pubKey;
//        uint32 weight; // really only need 1byte
//    	uint32 signerType;
//    	uint32 identity;
//
//    	 // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//    	case SIGNER_NAME:
//    		string256 name;
//        }
//        ext;
//    };
//
type Signer struct {
	PubKey     AccountId `json:"pubKey,omitempty"`
	Weight     Uint32    `json:"weight,omitempty"`
	SignerType Uint32    `json:"signerType,omitempty"`
	Identity   Uint32    `json:"identity,omitempty"`
	Ext        SignerExt `json:"ext,omitempty"`
}

// TrustEntryExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type TrustEntryExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u TrustEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of TrustEntryExt
func (u TrustEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewTrustEntryExt creates a new  TrustEntryExt.
func NewTrustEntryExt(v LedgerVersion, value interface{}) (result TrustEntryExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// TrustEntry is an XDR Struct defines as:
//
//   struct TrustEntry
//    {
//        AccountID allowedAccount;
//        BalanceID balanceToUse;
//
//    	 // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type TrustEntry struct {
	AllowedAccount AccountId     `json:"allowedAccount,omitempty"`
	BalanceToUse   BalanceId     `json:"balanceToUse,omitempty"`
	Ext            TrustEntryExt `json:"ext,omitempty"`
}

// LimitsExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type LimitsExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LimitsExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LimitsExt
func (u LimitsExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewLimitsExt creates a new  LimitsExt.
func NewLimitsExt(v LedgerVersion, value interface{}) (result LimitsExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// Limits is an XDR Struct defines as:
//
//   struct Limits
//    {
//        int64 dailyOut;
//    	int64 weeklyOut;
//    	int64 monthlyOut;
//        int64 annualOut;
//
//    	 // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//
//    };
//
type Limits struct {
	DailyOut   Int64     `json:"dailyOut,omitempty"`
	WeeklyOut  Int64     `json:"weeklyOut,omitempty"`
	MonthlyOut Int64     `json:"monthlyOut,omitempty"`
	AnnualOut  Int64     `json:"annualOut,omitempty"`
	Ext        LimitsExt `json:"ext,omitempty"`
}

// AccountPolicies is an XDR Enum defines as:
//
//   enum AccountPolicies
//    {
//    	NO_PERMISSIONS = 0,
//    	ALLOW_TO_CREATE_USER_VIA_API = 1,
//    	ALLOW_TO_TRANSFER_TOKENS = 2
//    };
//
type AccountPolicies int32

const (
	AccountPoliciesNoPermissions           AccountPolicies = 0
	AccountPoliciesAllowToCreateUserViaApi AccountPolicies = 1
	AccountPoliciesAllowToTransferTokens   AccountPolicies = 2
)

var AccountPoliciesAll = []AccountPolicies{
	AccountPoliciesNoPermissions,
	AccountPoliciesAllowToCreateUserViaApi,
	AccountPoliciesAllowToTransferTokens,
}

var accountPoliciesMap = map[int32]string{
	0: "AccountPoliciesNoPermissions",
	1: "AccountPoliciesAllowToCreateUserViaApi",
	2: "AccountPoliciesAllowToTransferTokens",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for AccountPolicies
func (e AccountPolicies) ValidEnum(v int32) bool {
	_, ok := accountPoliciesMap[v]
	return ok
}

// String returns the name of `e`
func (e AccountPolicies) String() string {
	name, _ := accountPoliciesMap[int32(e)]
	return name
}

func (e AccountPolicies) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// AccountType is an XDR Enum defines as:
//
//   enum AccountType
//    {
//    	OPERATIONAL = 1,       // operational account of the system
//    	GENERAL = 2,           // general account can perform payments, setoptions, be source account for tx, etc.
//    	COMMISSION = 3,        // commission account
//    	MASTER = 4,            // master account
//        EXCHANGE = 5,
//        NOT_VERIFIED = 6,
//    	STORAGE_FEE_MANAGER = 7 // manages storage fee
//    };
//
type AccountType int32

const (
	AccountTypeOperational       AccountType = 1
	AccountTypeGeneral           AccountType = 2
	AccountTypeCommission        AccountType = 3
	AccountTypeMaster            AccountType = 4
	AccountTypeExchange          AccountType = 5
	AccountTypeNotVerified       AccountType = 6
	AccountTypeStorageFeeManager AccountType = 7
)

var AccountTypeAll = []AccountType{
	AccountTypeOperational,
	AccountTypeGeneral,
	AccountTypeCommission,
	AccountTypeMaster,
	AccountTypeExchange,
	AccountTypeNotVerified,
	AccountTypeStorageFeeManager,
}

var accountTypeMap = map[int32]string{
	1: "AccountTypeOperational",
	2: "AccountTypeGeneral",
	3: "AccountTypeCommission",
	4: "AccountTypeMaster",
	5: "AccountTypeExchange",
	6: "AccountTypeNotVerified",
	7: "AccountTypeStorageFeeManager",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for AccountType
func (e AccountType) ValidEnum(v int32) bool {
	_, ok := accountTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e AccountType) String() string {
	name, _ := accountTypeMap[int32(e)]
	return name
}

func (e AccountType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// BlockReasons is an XDR Enum defines as:
//
//   enum BlockReasons
//    {
//    	RECOVERY_REQUEST = 1,
//    	KYC_UPDATE = 2,
//    	SUSPICIOUS_BEHAVIOR = 4
//    };
//
type BlockReasons int32

const (
	BlockReasonsRecoveryRequest    BlockReasons = 1
	BlockReasonsKycUpdate          BlockReasons = 2
	BlockReasonsSuspiciousBehavior BlockReasons = 4
)

var BlockReasonsAll = []BlockReasons{
	BlockReasonsRecoveryRequest,
	BlockReasonsKycUpdate,
	BlockReasonsSuspiciousBehavior,
}

var blockReasonsMap = map[int32]string{
	1: "BlockReasonsRecoveryRequest",
	2: "BlockReasonsKycUpdate",
	4: "BlockReasonsSuspiciousBehavior",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for BlockReasons
func (e BlockReasons) ValidEnum(v int32) bool {
	_, ok := blockReasonsMap[v]
	return ok
}

// String returns the name of `e`
func (e BlockReasons) String() string {
	name, _ := blockReasonsMap[int32(e)]
	return name
}

func (e BlockReasons) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// AccountEntryExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//    	case ACCOUNT_POLICIES:
//    		int32 policies;
//        }
//
type AccountEntryExt struct {
	V        LedgerVersion `json:"v,omitempty"`
	Policies *Int32        `json:"policies,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u AccountEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of AccountEntryExt
func (u AccountEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	case LedgerVersionAccountPolicies:
		return "Policies", true
	}
	return "-", false
}

// NewAccountEntryExt creates a new  AccountEntryExt.
func NewAccountEntryExt(v LedgerVersion, value interface{}) (result AccountEntryExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	case LedgerVersionAccountPolicies:
		tv, ok := value.(Int32)
		if !ok {
			err = fmt.Errorf("invalid value, must be Int32")
			return
		}
		result.Policies = &tv
	}
	return
}

// MustPolicies retrieves the Policies value from the union,
// panicing if the value is not set.
func (u AccountEntryExt) MustPolicies() Int32 {
	val, ok := u.GetPolicies()

	if !ok {
		panic("arm Policies is not set")
	}

	return val
}

// GetPolicies retrieves the Policies value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u AccountEntryExt) GetPolicies() (result Int32, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.V))

	if armName == "Policies" {
		result = *u.Policies
		ok = true
	}

	return
}

// AccountEntry is an XDR Struct defines as:
//
//   struct AccountEntry
//    {
//        AccountID accountID;      // master public key for this account
//
//        // fields used for signatures
//        // thresholds stores unsigned bytes: [weight of master|low|medium|high]
//        Thresholds thresholds;
//
//        Signer signers<>; // possible signers for this account
//        Limits* limits;
//
//    	uint32 blockReasons;
//        AccountType accountType; // type of the account
//
//        // Referral marketing
//        AccountID* referrer;     // parent account
//        int64 shareForReferrer; // share of fee to pay parent
//
//        // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//    	case ACCOUNT_POLICIES:
//    		int32 policies;
//        }
//        ext;
//    };
//
type AccountEntry struct {
	AccountId        AccountId       `json:"accountID,omitempty"`
	Thresholds       Thresholds      `json:"thresholds,omitempty"`
	Signers          []Signer        `json:"signers,omitempty"`
	Limits           *Limits         `json:"limits,omitempty"`
	BlockReasons     Uint32          `json:"blockReasons,omitempty"`
	AccountType      AccountType     `json:"accountType,omitempty"`
	Referrer         *AccountId      `json:"referrer,omitempty"`
	ShareForReferrer Int64           `json:"shareForReferrer,omitempty"`
	Ext              AccountEntryExt `json:"ext,omitempty"`
}

// AssetPairPolicy is an XDR Enum defines as:
//
//   enum AssetPairPolicy
//    {
//    	ASSET_PAIR_TRADEABLE = 1, // if not set pair can not be traided
//    	ASSET_PAIR_PHYSICAL_PRICE_RESTRICTION = 2, // if set, then prices for new offers must be greater then physical price with correction
//    	ASSET_PAIR_CURRENT_PRICE_RESTRICTION = 4, // if set, then price for new offers must be in interval of (1 +- maxPriceStep)*currentPrice
//    	ASSET_PAIR_DIRECT_BUY_ALLOWED = 8
//    };
//
type AssetPairPolicy int32

const (
	AssetPairPolicyAssetPairTradeable                AssetPairPolicy = 1
	AssetPairPolicyAssetPairPhysicalPriceRestriction AssetPairPolicy = 2
	AssetPairPolicyAssetPairCurrentPriceRestriction  AssetPairPolicy = 4
	AssetPairPolicyAssetPairDirectBuyAllowed         AssetPairPolicy = 8
)

var AssetPairPolicyAll = []AssetPairPolicy{
	AssetPairPolicyAssetPairTradeable,
	AssetPairPolicyAssetPairPhysicalPriceRestriction,
	AssetPairPolicyAssetPairCurrentPriceRestriction,
	AssetPairPolicyAssetPairDirectBuyAllowed,
}

var assetPairPolicyMap = map[int32]string{
	1: "AssetPairPolicyAssetPairTradeable",
	2: "AssetPairPolicyAssetPairPhysicalPriceRestriction",
	4: "AssetPairPolicyAssetPairCurrentPriceRestriction",
	8: "AssetPairPolicyAssetPairDirectBuyAllowed",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for AssetPairPolicy
func (e AssetPairPolicy) ValidEnum(v int32) bool {
	_, ok := assetPairPolicyMap[v]
	return ok
}

// String returns the name of `e`
func (e AssetPairPolicy) String() string {
	name, _ := assetPairPolicyMap[int32(e)]
	return name
}

func (e AssetPairPolicy) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// AssetPairEntryExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type AssetPairEntryExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u AssetPairEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of AssetPairEntryExt
func (u AssetPairEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewAssetPairEntryExt creates a new  AssetPairEntryExt.
func NewAssetPairEntryExt(v LedgerVersion, value interface{}) (result AssetPairEntryExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// AssetPairEntry is an XDR Struct defines as:
//
//   struct AssetPairEntry
//    {
//        AssetCode base;
//    	AssetCode quote;
//
//        int64 currentPrice;
//        int64 physicalPrice;
//
//    	int64 physicalPriceCorrection; // correction of physical price in percents. If physical price is set and restriction by physical price set, mininal price for offer for this pair will be physicalPrice * physicalPriceCorrection
//    	int64 maxPriceStep; // max price step in percent. User is allowed to set offer with price < (1 - maxPriceStep)*currentPrice and > (1 + maxPriceStep)*currentPrice
//
//
//    	int32 policies;
//
//        // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type AssetPairEntry struct {
	Base                    AssetCode         `json:"base,omitempty"`
	Quote                   AssetCode         `json:"quote,omitempty"`
	CurrentPrice            Int64             `json:"currentPrice,omitempty"`
	PhysicalPrice           Int64             `json:"physicalPrice,omitempty"`
	PhysicalPriceCorrection Int64             `json:"physicalPriceCorrection,omitempty"`
	MaxPriceStep            Int64             `json:"maxPriceStep,omitempty"`
	Policies                Int32             `json:"policies,omitempty"`
	Ext                     AssetPairEntryExt `json:"ext,omitempty"`
}

// AssetPolicy is an XDR Enum defines as:
//
//   enum AssetPolicy
//    {
//    	ASSET_TRANSFERABLE = 1,
//        ASSET_EMITTABLE_PRIMARY = 2,
//        ASSET_EMITTABLE_SECONDARY = 4,
//        ASSET_NOT_SHOW_EMPTY_DEMURRAGE = 8,
//    	ASSET_TOKEN_SHAREABLE_WITH_REFERRER = 16 // means that asset is a token and can be shared with referrer
//    };
//
type AssetPolicy int32

const (
	AssetPolicyAssetTransferable               AssetPolicy = 1
	AssetPolicyAssetEmittablePrimary           AssetPolicy = 2
	AssetPolicyAssetEmittableSecondary         AssetPolicy = 4
	AssetPolicyAssetNotShowEmptyDemurrage      AssetPolicy = 8
	AssetPolicyAssetTokenShareableWithReferrer AssetPolicy = 16
)

var AssetPolicyAll = []AssetPolicy{
	AssetPolicyAssetTransferable,
	AssetPolicyAssetEmittablePrimary,
	AssetPolicyAssetEmittableSecondary,
	AssetPolicyAssetNotShowEmptyDemurrage,
	AssetPolicyAssetTokenShareableWithReferrer,
}

var assetPolicyMap = map[int32]string{
	1:  "AssetPolicyAssetTransferable",
	2:  "AssetPolicyAssetEmittablePrimary",
	4:  "AssetPolicyAssetEmittableSecondary",
	8:  "AssetPolicyAssetNotShowEmptyDemurrage",
	16: "AssetPolicyAssetTokenShareableWithReferrer",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for AssetPolicy
func (e AssetPolicy) ValidEnum(v int32) bool {
	_, ok := assetPolicyMap[v]
	return ok
}

// String returns the name of `e`
func (e AssetPolicy) String() string {
	name, _ := assetPolicyMap[int32(e)]
	return name
}

func (e AssetPolicy) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// AssetFormExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type AssetFormExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u AssetFormExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of AssetFormExt
func (u AssetFormExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewAssetFormExt creates a new  AssetFormExt.
func NewAssetFormExt(v LedgerVersion, value interface{}) (result AssetFormExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// AssetForm is an XDR Struct defines as:
//
//   struct AssetForm {
//    	int64 unit;
//    	string64 name;
//    	 // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type AssetForm struct {
	Unit Int64        `json:"unit,omitempty"`
	Name String64     `json:"name,omitempty"`
	Ext  AssetFormExt `json:"ext,omitempty"`
}

// AssetEntryExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type AssetEntryExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u AssetEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of AssetEntryExt
func (u AssetEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewAssetEntryExt creates a new  AssetEntryExt.
func NewAssetEntryExt(v LedgerVersion, value interface{}) (result AssetEntryExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// AssetEntry is an XDR Struct defines as:
//
//   struct AssetEntry
//    {
//        AssetCode code;
//
//        AssetCode* token;
//        int32 policies;
//    	AssetForm assetForms<>;
//
//        // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type AssetEntry struct {
	Code       AssetCode     `json:"code,omitempty"`
	Token      *AssetCode    `json:"token,omitempty"`
	Policies   Int32         `json:"policies,omitempty"`
	AssetForms []AssetForm   `json:"assetForms,omitempty"`
	Ext        AssetEntryExt `json:"ext,omitempty"`
}

// BalanceEntryExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type BalanceEntryExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u BalanceEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of BalanceEntryExt
func (u BalanceEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewBalanceEntryExt creates a new  BalanceEntryExt.
func NewBalanceEntryExt(v LedgerVersion, value interface{}) (result BalanceEntryExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// BalanceEntry is an XDR Struct defines as:
//
//   struct BalanceEntry
//    {
//        BalanceID balanceID;
//        AssetCode asset;
//        AccountID accountID;
//        AccountID exchange;
//        int64 amount;
//        int64 locked;
//    	int64 feesPaid;
//    	int64 storageFee;
//    	uint64 storageFeeLastCalculated;
//    	uint64 storageFeeLastCharged;
//
//        // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type BalanceEntry struct {
	BalanceId                BalanceId       `json:"balanceID,omitempty"`
	Asset                    AssetCode       `json:"asset,omitempty"`
	AccountId                AccountId       `json:"accountID,omitempty"`
	Exchange                 AccountId       `json:"exchange,omitempty"`
	Amount                   Int64           `json:"amount,omitempty"`
	Locked                   Int64           `json:"locked,omitempty"`
	FeesPaid                 Int64           `json:"feesPaid,omitempty"`
	StorageFee               Int64           `json:"storageFee,omitempty"`
	StorageFeeLastCalculated Uint64          `json:"storageFeeLastCalculated,omitempty"`
	StorageFeeLastCharged    Uint64          `json:"storageFeeLastCharged,omitempty"`
	Ext                      BalanceEntryExt `json:"ext,omitempty"`
}

// CoinsEmissionRequestEntryExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type CoinsEmissionRequestEntryExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u CoinsEmissionRequestEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of CoinsEmissionRequestEntryExt
func (u CoinsEmissionRequestEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewCoinsEmissionRequestEntryExt creates a new  CoinsEmissionRequestEntryExt.
func NewCoinsEmissionRequestEntryExt(v LedgerVersion, value interface{}) (result CoinsEmissionRequestEntryExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// CoinsEmissionRequestEntry is an XDR Struct defines as:
//
//   struct CoinsEmissionRequestEntry
//    {
//    	uint64 requestID;
//        string64 reference;
//        BalanceID receiver;
//    	AccountID issuer;
//        int64 amount;
//        AssetCode asset;
//    	bool isApproved;
//
//    	// reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type CoinsEmissionRequestEntry struct {
	RequestId  Uint64                       `json:"requestID,omitempty"`
	Reference  String64                     `json:"reference,omitempty"`
	Receiver   BalanceId                    `json:"receiver,omitempty"`
	Issuer     AccountId                    `json:"issuer,omitempty"`
	Amount     Int64                        `json:"amount,omitempty"`
	Asset      AssetCode                    `json:"asset,omitempty"`
	IsApproved bool                         `json:"isApproved,omitempty"`
	Ext        CoinsEmissionRequestEntryExt `json:"ext,omitempty"`
}

// CoinsEmissionEntryExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type CoinsEmissionEntryExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u CoinsEmissionEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of CoinsEmissionEntryExt
func (u CoinsEmissionEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewCoinsEmissionEntryExt creates a new  CoinsEmissionEntryExt.
func NewCoinsEmissionEntryExt(v LedgerVersion, value interface{}) (result CoinsEmissionEntryExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// CoinsEmissionEntry is an XDR Struct defines as:
//
//   struct CoinsEmissionEntry
//    {
//    	string64 serialNumber;
//        int64 amount;
//        AssetCode asset;
//
//    	// reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type CoinsEmissionEntry struct {
	SerialNumber String64              `json:"serialNumber,omitempty"`
	Amount       Int64                 `json:"amount,omitempty"`
	Asset        AssetCode             `json:"asset,omitempty"`
	Ext          CoinsEmissionEntryExt `json:"ext,omitempty"`
}

// ExchangeAssetPolicy is an XDR Enum defines as:
//
//   enum ExchangeAssetPolicy
//    {
//    	CAN_OPERATE_BALANCE = 1,
//        CAN_REQUEST_EMISSION = 2,
//        EMIT_TOKENS_ON_OUT_TRANSFER = 4,
//        IS_NO_RECIPIENT_FEE = 8
//    };
//
type ExchangeAssetPolicy int32

const (
	ExchangeAssetPolicyCanOperateBalance       ExchangeAssetPolicy = 1
	ExchangeAssetPolicyCanRequestEmission      ExchangeAssetPolicy = 2
	ExchangeAssetPolicyEmitTokensOnOutTransfer ExchangeAssetPolicy = 4
	ExchangeAssetPolicyIsNoRecipientFee        ExchangeAssetPolicy = 8
)

var ExchangeAssetPolicyAll = []ExchangeAssetPolicy{
	ExchangeAssetPolicyCanOperateBalance,
	ExchangeAssetPolicyCanRequestEmission,
	ExchangeAssetPolicyEmitTokensOnOutTransfer,
	ExchangeAssetPolicyIsNoRecipientFee,
}

var exchangeAssetPolicyMap = map[int32]string{
	1: "ExchangeAssetPolicyCanOperateBalance",
	2: "ExchangeAssetPolicyCanRequestEmission",
	4: "ExchangeAssetPolicyEmitTokensOnOutTransfer",
	8: "ExchangeAssetPolicyIsNoRecipientFee",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ExchangeAssetPolicy
func (e ExchangeAssetPolicy) ValidEnum(v int32) bool {
	_, ok := exchangeAssetPolicyMap[v]
	return ok
}

// String returns the name of `e`
func (e ExchangeAssetPolicy) String() string {
	name, _ := exchangeAssetPolicyMap[int32(e)]
	return name
}

func (e ExchangeAssetPolicy) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// ExchangePoliciesEntryExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type ExchangePoliciesEntryExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ExchangePoliciesEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ExchangePoliciesEntryExt
func (u ExchangePoliciesEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewExchangePoliciesEntryExt creates a new  ExchangePoliciesEntryExt.
func NewExchangePoliciesEntryExt(v LedgerVersion, value interface{}) (result ExchangePoliciesEntryExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ExchangePoliciesEntry is an XDR Struct defines as:
//
//   struct ExchangePoliciesEntry {
//        AccountID accountID;
//        AssetCode asset;
//        int32 policies;
//
//    	 // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type ExchangePoliciesEntry struct {
	AccountId AccountId                `json:"accountID,omitempty"`
	Asset     AssetCode                `json:"asset,omitempty"`
	Policies  Int32                    `json:"policies,omitempty"`
	Ext       ExchangePoliciesEntryExt `json:"ext,omitempty"`
}

// ExchangeDataEntryExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type ExchangeDataEntryExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ExchangeDataEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ExchangeDataEntryExt
func (u ExchangeDataEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewExchangeDataEntryExt creates a new  ExchangeDataEntryExt.
func NewExchangeDataEntryExt(v LedgerVersion, value interface{}) (result ExchangeDataEntryExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ExchangeDataEntry is an XDR Struct defines as:
//
//   struct ExchangeDataEntry
//    {
//        AccountID accountID;
//        string64 name;
//        bool requireReview;
//
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type ExchangeDataEntry struct {
	AccountId     AccountId            `json:"accountID,omitempty"`
	Name          String64             `json:"name,omitempty"`
	RequireReview bool                 `json:"requireReview,omitempty"`
	Ext           ExchangeDataEntryExt `json:"ext,omitempty"`
}

// FeeType is an XDR Enum defines as:
//
//   enum FeeType
//    {
//        PAYMENT_FEE = 0,
//    	STORAGE_FEE = 1,
//        REFERRAL_FEE = 2,
//    	OFFER_FEE = 3,
//        FORFEIT_FEE = 4,
//        EMISSION_FEE = 5
//    };
//
type FeeType int32

const (
	FeeTypePaymentFee  FeeType = 0
	FeeTypeStorageFee  FeeType = 1
	FeeTypeReferralFee FeeType = 2
	FeeTypeOfferFee    FeeType = 3
	FeeTypeForfeitFee  FeeType = 4
	FeeTypeEmissionFee FeeType = 5
)

var FeeTypeAll = []FeeType{
	FeeTypePaymentFee,
	FeeTypeStorageFee,
	FeeTypeReferralFee,
	FeeTypeOfferFee,
	FeeTypeForfeitFee,
	FeeTypeEmissionFee,
}

var feeTypeMap = map[int32]string{
	0: "FeeTypePaymentFee",
	1: "FeeTypeStorageFee",
	2: "FeeTypeReferralFee",
	3: "FeeTypeOfferFee",
	4: "FeeTypeForfeitFee",
	5: "FeeTypeEmissionFee",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for FeeType
func (e FeeType) ValidEnum(v int32) bool {
	_, ok := feeTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e FeeType) String() string {
	name, _ := feeTypeMap[int32(e)]
	return name
}

func (e FeeType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// EmissionFeeType is an XDR Enum defines as:
//
//   enum EmissionFeeType
//    {
//    	PRIMARY_MARKET = 1,
//    	SECONDARY_MARKET = 2
//    };
//
type EmissionFeeType int32

const (
	EmissionFeeTypePrimaryMarket   EmissionFeeType = 1
	EmissionFeeTypeSecondaryMarket EmissionFeeType = 2
)

var EmissionFeeTypeAll = []EmissionFeeType{
	EmissionFeeTypePrimaryMarket,
	EmissionFeeTypeSecondaryMarket,
}

var emissionFeeTypeMap = map[int32]string{
	1: "EmissionFeeTypePrimaryMarket",
	2: "EmissionFeeTypeSecondaryMarket",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for EmissionFeeType
func (e EmissionFeeType) ValidEnum(v int32) bool {
	_, ok := emissionFeeTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e EmissionFeeType) String() string {
	name, _ := emissionFeeTypeMap[int32(e)]
	return name
}

func (e EmissionFeeType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// FeeEntryExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type FeeEntryExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u FeeEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of FeeEntryExt
func (u FeeEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewFeeEntryExt creates a new  FeeEntryExt.
func NewFeeEntryExt(v LedgerVersion, value interface{}) (result FeeEntryExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// FeeEntry is an XDR Struct defines as:
//
//   struct FeeEntry
//    {
//        FeeType feeType;
//        AssetCode asset;
//        int64 fixedFee; // fee paid for operation
//    	int64 percentFee; // percent of transfer amount to be charged
//
//        AccountID* accountID;
//        AccountType* accountType;
//        int64 subtype; // for example, different withdrawals — bars or coins
//
//        int64 lowerBound;
//        int64 upperBound;
//
//        Hash hash;
//
//    	// reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//
//    };
//
type FeeEntry struct {
	FeeType     FeeType      `json:"feeType,omitempty"`
	Asset       AssetCode    `json:"asset,omitempty"`
	FixedFee    Int64        `json:"fixedFee,omitempty"`
	PercentFee  Int64        `json:"percentFee,omitempty"`
	AccountId   *AccountId   `json:"accountID,omitempty"`
	AccountType *AccountType `json:"accountType,omitempty"`
	Subtype     Int64        `json:"subtype,omitempty"`
	LowerBound  Int64        `json:"lowerBound,omitempty"`
	UpperBound  Int64        `json:"upperBound,omitempty"`
	Hash        Hash         `json:"hash,omitempty"`
	Ext         FeeEntryExt  `json:"ext,omitempty"`
}

// InvoiceState is an XDR Enum defines as:
//
//   enum InvoiceState
//    {
//        INVOICE_NEEDS_PAYMENT = 0,
//        INVOICE_NEEDS_PAYMENT_REVIEW = 1
//    };
//
type InvoiceState int32

const (
	InvoiceStateInvoiceNeedsPayment       InvoiceState = 0
	InvoiceStateInvoiceNeedsPaymentReview InvoiceState = 1
)

var InvoiceStateAll = []InvoiceState{
	InvoiceStateInvoiceNeedsPayment,
	InvoiceStateInvoiceNeedsPaymentReview,
}

var invoiceStateMap = map[int32]string{
	0: "InvoiceStateInvoiceNeedsPayment",
	1: "InvoiceStateInvoiceNeedsPaymentReview",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for InvoiceState
func (e InvoiceState) ValidEnum(v int32) bool {
	_, ok := invoiceStateMap[v]
	return ok
}

// String returns the name of `e`
func (e InvoiceState) String() string {
	name, _ := invoiceStateMap[int32(e)]
	return name
}

func (e InvoiceState) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// InvoiceEntryExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type InvoiceEntryExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u InvoiceEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of InvoiceEntryExt
func (u InvoiceEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewInvoiceEntryExt creates a new  InvoiceEntryExt.
func NewInvoiceEntryExt(v LedgerVersion, value interface{}) (result InvoiceEntryExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// InvoiceEntry is an XDR Struct defines as:
//
//   struct InvoiceEntry
//    {
//        uint64 invoiceID;
//        AccountID receiverAccount;
//        BalanceID receiverBalance;
//    	AccountID sender;
//        int64 amount;
//
//        InvoiceState state;
//
//        // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type InvoiceEntry struct {
	InvoiceId       Uint64          `json:"invoiceID,omitempty"`
	ReceiverAccount AccountId       `json:"receiverAccount,omitempty"`
	ReceiverBalance BalanceId       `json:"receiverBalance,omitempty"`
	Sender          AccountId       `json:"sender,omitempty"`
	Amount          Int64           `json:"amount,omitempty"`
	State           InvoiceState    `json:"state,omitempty"`
	Ext             InvoiceEntryExt `json:"ext,omitempty"`
}

// OfferEntryExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type OfferEntryExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u OfferEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of OfferEntryExt
func (u OfferEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewOfferEntryExt creates a new  OfferEntryExt.
func NewOfferEntryExt(v LedgerVersion, value interface{}) (result OfferEntryExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// OfferEntry is an XDR Struct defines as:
//
//   struct OfferEntry
//    {
//        uint64 offerID;
//    	AccountID ownerID;
//    	bool isBuy;
//        AssetCode base; // A
//        AssetCode quote;  // B
//    	BalanceID baseBalance;
//    	BalanceID quoteBalance;
//        int64 baseAmount;
//    	int64 quoteAmount;
//    	uint64 createdAt;
//    	int64 fee;
//
//        int64 percentFee;
//
//    	// price of A in terms of B
//        int64 price;
//
//        // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type OfferEntry struct {
	OfferId      Uint64        `json:"offerID,omitempty"`
	OwnerId      AccountId     `json:"ownerID,omitempty"`
	IsBuy        bool          `json:"isBuy,omitempty"`
	Base         AssetCode     `json:"base,omitempty"`
	Quote        AssetCode     `json:"quote,omitempty"`
	BaseBalance  BalanceId     `json:"baseBalance,omitempty"`
	QuoteBalance BalanceId     `json:"quoteBalance,omitempty"`
	BaseAmount   Int64         `json:"baseAmount,omitempty"`
	QuoteAmount  Int64         `json:"quoteAmount,omitempty"`
	CreatedAt    Uint64        `json:"createdAt,omitempty"`
	Fee          Int64         `json:"fee,omitempty"`
	PercentFee   Int64         `json:"percentFee,omitempty"`
	Price        Int64         `json:"price,omitempty"`
	Ext          OfferEntryExt `json:"ext,omitempty"`
}

// RequestType is an XDR Enum defines as:
//
//   enum RequestType
//    {
//        REQUEST_TYPE_SALE = 0,
//        REQUEST_TYPE_WITHDRAWAL = 1,
//        REQUEST_TYPE_REDEEM = 2,
//        REQUEST_TYPE_PAYMENT = 3
//    };
//
type RequestType int32

const (
	RequestTypeRequestTypeSale       RequestType = 0
	RequestTypeRequestTypeWithdrawal RequestType = 1
	RequestTypeRequestTypeRedeem     RequestType = 2
	RequestTypeRequestTypePayment    RequestType = 3
)

var RequestTypeAll = []RequestType{
	RequestTypeRequestTypeSale,
	RequestTypeRequestTypeWithdrawal,
	RequestTypeRequestTypeRedeem,
	RequestTypeRequestTypePayment,
}

var requestTypeMap = map[int32]string{
	0: "RequestTypeRequestTypeSale",
	1: "RequestTypeRequestTypeWithdrawal",
	2: "RequestTypeRequestTypeRedeem",
	3: "RequestTypeRequestTypePayment",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for RequestType
func (e RequestType) ValidEnum(v int32) bool {
	_, ok := requestTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e RequestType) String() string {
	name, _ := requestTypeMap[int32(e)]
	return name
}

func (e RequestType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// PaymentRequestEntryExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type PaymentRequestEntryExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u PaymentRequestEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of PaymentRequestEntryExt
func (u PaymentRequestEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewPaymentRequestEntryExt creates a new  PaymentRequestEntryExt.
func NewPaymentRequestEntryExt(v LedgerVersion, value interface{}) (result PaymentRequestEntryExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// PaymentRequestEntry is an XDR Struct defines as:
//
//   struct PaymentRequestEntry
//    {
//        uint64 paymentID;
//    	AccountID exchange;
//        BalanceID sourceBalance;
//        BalanceID* destinationBalance;
//        int64 sourceSend;
//        int64 sourceSendUniversal;
//        int64 destinationReceive;
//
//        uint64 createdAt;
//
//        uint64* invoiceID;
//
//    	// reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type PaymentRequestEntry struct {
	PaymentId           Uint64                 `json:"paymentID,omitempty"`
	Exchange            AccountId              `json:"exchange,omitempty"`
	SourceBalance       BalanceId              `json:"sourceBalance,omitempty"`
	DestinationBalance  *BalanceId             `json:"destinationBalance,omitempty"`
	SourceSend          Int64                  `json:"sourceSend,omitempty"`
	SourceSendUniversal Int64                  `json:"sourceSendUniversal,omitempty"`
	DestinationReceive  Int64                  `json:"destinationReceive,omitempty"`
	CreatedAt           Uint64                 `json:"createdAt,omitempty"`
	InvoiceId           *Uint64                `json:"invoiceID,omitempty"`
	Ext                 PaymentRequestEntryExt `json:"ext,omitempty"`
}

// PaymentEntryExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type PaymentEntryExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u PaymentEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of PaymentEntryExt
func (u PaymentEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewPaymentEntryExt creates a new  PaymentEntryExt.
func NewPaymentEntryExt(v LedgerVersion, value interface{}) (result PaymentEntryExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// PaymentEntry is an XDR Struct defines as:
//
//   struct PaymentEntry
//    {
//    	AccountID exchange;
//        string64 reference;
//
//    	// reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type PaymentEntry struct {
	Exchange  AccountId       `json:"exchange,omitempty"`
	Reference String64        `json:"reference,omitempty"`
	Ext       PaymentEntryExt `json:"ext,omitempty"`
}

// StatisticsEntryExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type StatisticsEntryExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u StatisticsEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of StatisticsEntryExt
func (u StatisticsEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewStatisticsEntryExt creates a new  StatisticsEntryExt.
func NewStatisticsEntryExt(v LedgerVersion, value interface{}) (result StatisticsEntryExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// StatisticsEntry is an XDR Struct defines as:
//
//   struct StatisticsEntry
//    {
//    	AccountID accountID;
//
//    	int64 dailyOutcome;
//    	int64 weeklyOutcome;
//    	int64 monthlyOutcome;
//    	int64 annualOutcome;
//
//    	int64 updatedAt;
//
//        // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type StatisticsEntry struct {
	AccountId      AccountId          `json:"accountID,omitempty"`
	DailyOutcome   Int64              `json:"dailyOutcome,omitempty"`
	WeeklyOutcome  Int64              `json:"weeklyOutcome,omitempty"`
	MonthlyOutcome Int64              `json:"monthlyOutcome,omitempty"`
	AnnualOutcome  Int64              `json:"annualOutcome,omitempty"`
	UpdatedAt      Int64              `json:"updatedAt,omitempty"`
	Ext            StatisticsEntryExt `json:"ext,omitempty"`
}

// ThresholdIndexes is an XDR Enum defines as:
//
//   enum ThresholdIndexes
//    {
//        THRESHOLD_MASTER_WEIGHT = 0,
//        THRESHOLD_LOW = 1,
//        THRESHOLD_MED = 2,
//        THRESHOLD_HIGH = 3
//    };
//
type ThresholdIndexes int32

const (
	ThresholdIndexesThresholdMasterWeight ThresholdIndexes = 0
	ThresholdIndexesThresholdLow          ThresholdIndexes = 1
	ThresholdIndexesThresholdMed          ThresholdIndexes = 2
	ThresholdIndexesThresholdHigh         ThresholdIndexes = 3
)

var ThresholdIndexesAll = []ThresholdIndexes{
	ThresholdIndexesThresholdMasterWeight,
	ThresholdIndexesThresholdLow,
	ThresholdIndexesThresholdMed,
	ThresholdIndexesThresholdHigh,
}

var thresholdIndexesMap = map[int32]string{
	0: "ThresholdIndexesThresholdMasterWeight",
	1: "ThresholdIndexesThresholdLow",
	2: "ThresholdIndexesThresholdMed",
	3: "ThresholdIndexesThresholdHigh",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ThresholdIndexes
func (e ThresholdIndexes) ValidEnum(v int32) bool {
	_, ok := thresholdIndexesMap[v]
	return ok
}

// String returns the name of `e`
func (e ThresholdIndexes) String() string {
	name, _ := thresholdIndexesMap[int32(e)]
	return name
}

func (e ThresholdIndexes) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// LedgerEntryType is an XDR Enum defines as:
//
//   enum LedgerEntryType
//    {
//        ACCOUNT = 0,
//    	COINS_EMISSION_REQUEST = 1,
//        FEE = 2,
//        COINS_EMISSION = 3,
//        BALANCE = 4,
//        PAYMENT_REQUEST = 5,
//        ASSET = 6,
//        PAYMENT_ENTRY = 7,
//        ACCOUNT_TYPE_LIMITS = 8,
//        STATISTICS = 9,
//        EXCHANGE_DATA = 10,
//        EXCHANGE_POLICIES = 11,
//        TRUST = 12,
//        ACCOUNT_LIMITS = 13,
//    	ASSET_PAIR = 14,
//    	OFFER_ENTRY = 15,
//        INVOICE = 16
//    };
//
type LedgerEntryType int32

const (
	LedgerEntryTypeAccount              LedgerEntryType = 0
	LedgerEntryTypeCoinsEmissionRequest LedgerEntryType = 1
	LedgerEntryTypeFee                  LedgerEntryType = 2
	LedgerEntryTypeCoinsEmission        LedgerEntryType = 3
	LedgerEntryTypeBalance              LedgerEntryType = 4
	LedgerEntryTypePaymentRequest       LedgerEntryType = 5
	LedgerEntryTypeAsset                LedgerEntryType = 6
	LedgerEntryTypePaymentEntry         LedgerEntryType = 7
	LedgerEntryTypeAccountTypeLimits    LedgerEntryType = 8
	LedgerEntryTypeStatistics           LedgerEntryType = 9
	LedgerEntryTypeExchangeData         LedgerEntryType = 10
	LedgerEntryTypeExchangePolicies     LedgerEntryType = 11
	LedgerEntryTypeTrust                LedgerEntryType = 12
	LedgerEntryTypeAccountLimits        LedgerEntryType = 13
	LedgerEntryTypeAssetPair            LedgerEntryType = 14
	LedgerEntryTypeOfferEntry           LedgerEntryType = 15
	LedgerEntryTypeInvoice              LedgerEntryType = 16
)

var LedgerEntryTypeAll = []LedgerEntryType{
	LedgerEntryTypeAccount,
	LedgerEntryTypeCoinsEmissionRequest,
	LedgerEntryTypeFee,
	LedgerEntryTypeCoinsEmission,
	LedgerEntryTypeBalance,
	LedgerEntryTypePaymentRequest,
	LedgerEntryTypeAsset,
	LedgerEntryTypePaymentEntry,
	LedgerEntryTypeAccountTypeLimits,
	LedgerEntryTypeStatistics,
	LedgerEntryTypeExchangeData,
	LedgerEntryTypeExchangePolicies,
	LedgerEntryTypeTrust,
	LedgerEntryTypeAccountLimits,
	LedgerEntryTypeAssetPair,
	LedgerEntryTypeOfferEntry,
	LedgerEntryTypeInvoice,
}

var ledgerEntryTypeMap = map[int32]string{
	0:  "LedgerEntryTypeAccount",
	1:  "LedgerEntryTypeCoinsEmissionRequest",
	2:  "LedgerEntryTypeFee",
	3:  "LedgerEntryTypeCoinsEmission",
	4:  "LedgerEntryTypeBalance",
	5:  "LedgerEntryTypePaymentRequest",
	6:  "LedgerEntryTypeAsset",
	7:  "LedgerEntryTypePaymentEntry",
	8:  "LedgerEntryTypeAccountTypeLimits",
	9:  "LedgerEntryTypeStatistics",
	10: "LedgerEntryTypeExchangeData",
	11: "LedgerEntryTypeExchangePolicies",
	12: "LedgerEntryTypeTrust",
	13: "LedgerEntryTypeAccountLimits",
	14: "LedgerEntryTypeAssetPair",
	15: "LedgerEntryTypeOfferEntry",
	16: "LedgerEntryTypeInvoice",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for LedgerEntryType
func (e LedgerEntryType) ValidEnum(v int32) bool {
	_, ok := ledgerEntryTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e LedgerEntryType) String() string {
	name, _ := ledgerEntryTypeMap[int32(e)]
	return name
}

func (e LedgerEntryType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// LedgerEntryData is an XDR NestedUnion defines as:
//
//   union switch (LedgerEntryType type)
//        {
//        case ACCOUNT:
//            AccountEntry account;
//    	case COINS_EMISSION_REQUEST:
//    		CoinsEmissionRequestEntry coinsEmissionRequest;
//        case FEE:
//            FeeEntry feeState;
//        case COINS_EMISSION:
//    		CoinsEmissionEntry coinsEmission;
//        case BALANCE:
//            BalanceEntry balance;
//        case PAYMENT_REQUEST:
//            PaymentRequestEntry paymentRequest;
//        case ASSET:
//            AssetEntry asset;
//        case PAYMENT_ENTRY:
//            PaymentEntry payment;
//        case ACCOUNT_TYPE_LIMITS:
//            AccountTypeLimitsEntry accountTypeLimits;
//        case STATISTICS:
//            StatisticsEntry stats;
//        case EXCHANGE_DATA:
//            ExchangeDataEntry exchangeData;
//        case EXCHANGE_POLICIES:
//            ExchangePoliciesEntry exchangePolicies;
//        case TRUST:
//            TrustEntry trust;
//        case ACCOUNT_LIMITS:
//            AccountLimitsEntry accountLimits;
//    	case ASSET_PAIR:
//    		AssetPairEntry assetPair;
//    	case OFFER_ENTRY:
//    		OfferEntry offer;
//        case INVOICE:
//            InvoiceEntry invoice;
//        }
//
type LedgerEntryData struct {
	Type                 LedgerEntryType            `json:"type,omitempty"`
	Account              *AccountEntry              `json:"account,omitempty"`
	CoinsEmissionRequest *CoinsEmissionRequestEntry `json:"coinsEmissionRequest,omitempty"`
	FeeState             *FeeEntry                  `json:"feeState,omitempty"`
	CoinsEmission        *CoinsEmissionEntry        `json:"coinsEmission,omitempty"`
	Balance              *BalanceEntry              `json:"balance,omitempty"`
	PaymentRequest       *PaymentRequestEntry       `json:"paymentRequest,omitempty"`
	Asset                *AssetEntry                `json:"asset,omitempty"`
	Payment              *PaymentEntry              `json:"payment,omitempty"`
	AccountTypeLimits    *AccountTypeLimitsEntry    `json:"accountTypeLimits,omitempty"`
	Stats                *StatisticsEntry           `json:"stats,omitempty"`
	ExchangeData         *ExchangeDataEntry         `json:"exchangeData,omitempty"`
	ExchangePolicies     *ExchangePoliciesEntry     `json:"exchangePolicies,omitempty"`
	Trust                *TrustEntry                `json:"trust,omitempty"`
	AccountLimits        *AccountLimitsEntry        `json:"accountLimits,omitempty"`
	AssetPair            *AssetPairEntry            `json:"assetPair,omitempty"`
	Offer                *OfferEntry                `json:"offer,omitempty"`
	Invoice              *InvoiceEntry              `json:"invoice,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerEntryData) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerEntryData
func (u LedgerEntryData) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerEntryType(sw) {
	case LedgerEntryTypeAccount:
		return "Account", true
	case LedgerEntryTypeCoinsEmissionRequest:
		return "CoinsEmissionRequest", true
	case LedgerEntryTypeFee:
		return "FeeState", true
	case LedgerEntryTypeCoinsEmission:
		return "CoinsEmission", true
	case LedgerEntryTypeBalance:
		return "Balance", true
	case LedgerEntryTypePaymentRequest:
		return "PaymentRequest", true
	case LedgerEntryTypeAsset:
		return "Asset", true
	case LedgerEntryTypePaymentEntry:
		return "Payment", true
	case LedgerEntryTypeAccountTypeLimits:
		return "AccountTypeLimits", true
	case LedgerEntryTypeStatistics:
		return "Stats", true
	case LedgerEntryTypeExchangeData:
		return "ExchangeData", true
	case LedgerEntryTypeExchangePolicies:
		return "ExchangePolicies", true
	case LedgerEntryTypeTrust:
		return "Trust", true
	case LedgerEntryTypeAccountLimits:
		return "AccountLimits", true
	case LedgerEntryTypeAssetPair:
		return "AssetPair", true
	case LedgerEntryTypeOfferEntry:
		return "Offer", true
	case LedgerEntryTypeInvoice:
		return "Invoice", true
	}
	return "-", false
}

// NewLedgerEntryData creates a new  LedgerEntryData.
func NewLedgerEntryData(aType LedgerEntryType, value interface{}) (result LedgerEntryData, err error) {
	result.Type = aType
	switch LedgerEntryType(aType) {
	case LedgerEntryTypeAccount:
		tv, ok := value.(AccountEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be AccountEntry")
			return
		}
		result.Account = &tv
	case LedgerEntryTypeCoinsEmissionRequest:
		tv, ok := value.(CoinsEmissionRequestEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be CoinsEmissionRequestEntry")
			return
		}
		result.CoinsEmissionRequest = &tv
	case LedgerEntryTypeFee:
		tv, ok := value.(FeeEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be FeeEntry")
			return
		}
		result.FeeState = &tv
	case LedgerEntryTypeCoinsEmission:
		tv, ok := value.(CoinsEmissionEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be CoinsEmissionEntry")
			return
		}
		result.CoinsEmission = &tv
	case LedgerEntryTypeBalance:
		tv, ok := value.(BalanceEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be BalanceEntry")
			return
		}
		result.Balance = &tv
	case LedgerEntryTypePaymentRequest:
		tv, ok := value.(PaymentRequestEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be PaymentRequestEntry")
			return
		}
		result.PaymentRequest = &tv
	case LedgerEntryTypeAsset:
		tv, ok := value.(AssetEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be AssetEntry")
			return
		}
		result.Asset = &tv
	case LedgerEntryTypePaymentEntry:
		tv, ok := value.(PaymentEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be PaymentEntry")
			return
		}
		result.Payment = &tv
	case LedgerEntryTypeAccountTypeLimits:
		tv, ok := value.(AccountTypeLimitsEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be AccountTypeLimitsEntry")
			return
		}
		result.AccountTypeLimits = &tv
	case LedgerEntryTypeStatistics:
		tv, ok := value.(StatisticsEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be StatisticsEntry")
			return
		}
		result.Stats = &tv
	case LedgerEntryTypeExchangeData:
		tv, ok := value.(ExchangeDataEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be ExchangeDataEntry")
			return
		}
		result.ExchangeData = &tv
	case LedgerEntryTypeExchangePolicies:
		tv, ok := value.(ExchangePoliciesEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be ExchangePoliciesEntry")
			return
		}
		result.ExchangePolicies = &tv
	case LedgerEntryTypeTrust:
		tv, ok := value.(TrustEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be TrustEntry")
			return
		}
		result.Trust = &tv
	case LedgerEntryTypeAccountLimits:
		tv, ok := value.(AccountLimitsEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be AccountLimitsEntry")
			return
		}
		result.AccountLimits = &tv
	case LedgerEntryTypeAssetPair:
		tv, ok := value.(AssetPairEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be AssetPairEntry")
			return
		}
		result.AssetPair = &tv
	case LedgerEntryTypeOfferEntry:
		tv, ok := value.(OfferEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be OfferEntry")
			return
		}
		result.Offer = &tv
	case LedgerEntryTypeInvoice:
		tv, ok := value.(InvoiceEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be InvoiceEntry")
			return
		}
		result.Invoice = &tv
	}
	return
}

// MustAccount retrieves the Account value from the union,
// panicing if the value is not set.
func (u LedgerEntryData) MustAccount() AccountEntry {
	val, ok := u.GetAccount()

	if !ok {
		panic("arm Account is not set")
	}

	return val
}

// GetAccount retrieves the Account value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryData) GetAccount() (result AccountEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Account" {
		result = *u.Account
		ok = true
	}

	return
}

// MustCoinsEmissionRequest retrieves the CoinsEmissionRequest value from the union,
// panicing if the value is not set.
func (u LedgerEntryData) MustCoinsEmissionRequest() CoinsEmissionRequestEntry {
	val, ok := u.GetCoinsEmissionRequest()

	if !ok {
		panic("arm CoinsEmissionRequest is not set")
	}

	return val
}

// GetCoinsEmissionRequest retrieves the CoinsEmissionRequest value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryData) GetCoinsEmissionRequest() (result CoinsEmissionRequestEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "CoinsEmissionRequest" {
		result = *u.CoinsEmissionRequest
		ok = true
	}

	return
}

// MustFeeState retrieves the FeeState value from the union,
// panicing if the value is not set.
func (u LedgerEntryData) MustFeeState() FeeEntry {
	val, ok := u.GetFeeState()

	if !ok {
		panic("arm FeeState is not set")
	}

	return val
}

// GetFeeState retrieves the FeeState value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryData) GetFeeState() (result FeeEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "FeeState" {
		result = *u.FeeState
		ok = true
	}

	return
}

// MustCoinsEmission retrieves the CoinsEmission value from the union,
// panicing if the value is not set.
func (u LedgerEntryData) MustCoinsEmission() CoinsEmissionEntry {
	val, ok := u.GetCoinsEmission()

	if !ok {
		panic("arm CoinsEmission is not set")
	}

	return val
}

// GetCoinsEmission retrieves the CoinsEmission value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryData) GetCoinsEmission() (result CoinsEmissionEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "CoinsEmission" {
		result = *u.CoinsEmission
		ok = true
	}

	return
}

// MustBalance retrieves the Balance value from the union,
// panicing if the value is not set.
func (u LedgerEntryData) MustBalance() BalanceEntry {
	val, ok := u.GetBalance()

	if !ok {
		panic("arm Balance is not set")
	}

	return val
}

// GetBalance retrieves the Balance value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryData) GetBalance() (result BalanceEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Balance" {
		result = *u.Balance
		ok = true
	}

	return
}

// MustPaymentRequest retrieves the PaymentRequest value from the union,
// panicing if the value is not set.
func (u LedgerEntryData) MustPaymentRequest() PaymentRequestEntry {
	val, ok := u.GetPaymentRequest()

	if !ok {
		panic("arm PaymentRequest is not set")
	}

	return val
}

// GetPaymentRequest retrieves the PaymentRequest value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryData) GetPaymentRequest() (result PaymentRequestEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "PaymentRequest" {
		result = *u.PaymentRequest
		ok = true
	}

	return
}

// MustAsset retrieves the Asset value from the union,
// panicing if the value is not set.
func (u LedgerEntryData) MustAsset() AssetEntry {
	val, ok := u.GetAsset()

	if !ok {
		panic("arm Asset is not set")
	}

	return val
}

// GetAsset retrieves the Asset value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryData) GetAsset() (result AssetEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Asset" {
		result = *u.Asset
		ok = true
	}

	return
}

// MustPayment retrieves the Payment value from the union,
// panicing if the value is not set.
func (u LedgerEntryData) MustPayment() PaymentEntry {
	val, ok := u.GetPayment()

	if !ok {
		panic("arm Payment is not set")
	}

	return val
}

// GetPayment retrieves the Payment value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryData) GetPayment() (result PaymentEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Payment" {
		result = *u.Payment
		ok = true
	}

	return
}

// MustAccountTypeLimits retrieves the AccountTypeLimits value from the union,
// panicing if the value is not set.
func (u LedgerEntryData) MustAccountTypeLimits() AccountTypeLimitsEntry {
	val, ok := u.GetAccountTypeLimits()

	if !ok {
		panic("arm AccountTypeLimits is not set")
	}

	return val
}

// GetAccountTypeLimits retrieves the AccountTypeLimits value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryData) GetAccountTypeLimits() (result AccountTypeLimitsEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "AccountTypeLimits" {
		result = *u.AccountTypeLimits
		ok = true
	}

	return
}

// MustStats retrieves the Stats value from the union,
// panicing if the value is not set.
func (u LedgerEntryData) MustStats() StatisticsEntry {
	val, ok := u.GetStats()

	if !ok {
		panic("arm Stats is not set")
	}

	return val
}

// GetStats retrieves the Stats value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryData) GetStats() (result StatisticsEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Stats" {
		result = *u.Stats
		ok = true
	}

	return
}

// MustExchangeData retrieves the ExchangeData value from the union,
// panicing if the value is not set.
func (u LedgerEntryData) MustExchangeData() ExchangeDataEntry {
	val, ok := u.GetExchangeData()

	if !ok {
		panic("arm ExchangeData is not set")
	}

	return val
}

// GetExchangeData retrieves the ExchangeData value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryData) GetExchangeData() (result ExchangeDataEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ExchangeData" {
		result = *u.ExchangeData
		ok = true
	}

	return
}

// MustExchangePolicies retrieves the ExchangePolicies value from the union,
// panicing if the value is not set.
func (u LedgerEntryData) MustExchangePolicies() ExchangePoliciesEntry {
	val, ok := u.GetExchangePolicies()

	if !ok {
		panic("arm ExchangePolicies is not set")
	}

	return val
}

// GetExchangePolicies retrieves the ExchangePolicies value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryData) GetExchangePolicies() (result ExchangePoliciesEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ExchangePolicies" {
		result = *u.ExchangePolicies
		ok = true
	}

	return
}

// MustTrust retrieves the Trust value from the union,
// panicing if the value is not set.
func (u LedgerEntryData) MustTrust() TrustEntry {
	val, ok := u.GetTrust()

	if !ok {
		panic("arm Trust is not set")
	}

	return val
}

// GetTrust retrieves the Trust value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryData) GetTrust() (result TrustEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Trust" {
		result = *u.Trust
		ok = true
	}

	return
}

// MustAccountLimits retrieves the AccountLimits value from the union,
// panicing if the value is not set.
func (u LedgerEntryData) MustAccountLimits() AccountLimitsEntry {
	val, ok := u.GetAccountLimits()

	if !ok {
		panic("arm AccountLimits is not set")
	}

	return val
}

// GetAccountLimits retrieves the AccountLimits value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryData) GetAccountLimits() (result AccountLimitsEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "AccountLimits" {
		result = *u.AccountLimits
		ok = true
	}

	return
}

// MustAssetPair retrieves the AssetPair value from the union,
// panicing if the value is not set.
func (u LedgerEntryData) MustAssetPair() AssetPairEntry {
	val, ok := u.GetAssetPair()

	if !ok {
		panic("arm AssetPair is not set")
	}

	return val
}

// GetAssetPair retrieves the AssetPair value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryData) GetAssetPair() (result AssetPairEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "AssetPair" {
		result = *u.AssetPair
		ok = true
	}

	return
}

// MustOffer retrieves the Offer value from the union,
// panicing if the value is not set.
func (u LedgerEntryData) MustOffer() OfferEntry {
	val, ok := u.GetOffer()

	if !ok {
		panic("arm Offer is not set")
	}

	return val
}

// GetOffer retrieves the Offer value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryData) GetOffer() (result OfferEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Offer" {
		result = *u.Offer
		ok = true
	}

	return
}

// MustInvoice retrieves the Invoice value from the union,
// panicing if the value is not set.
func (u LedgerEntryData) MustInvoice() InvoiceEntry {
	val, ok := u.GetInvoice()

	if !ok {
		panic("arm Invoice is not set")
	}

	return val
}

// GetInvoice retrieves the Invoice value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryData) GetInvoice() (result InvoiceEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Invoice" {
		result = *u.Invoice
		ok = true
	}

	return
}

// LedgerEntryExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type LedgerEntryExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerEntryExt
func (u LedgerEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewLedgerEntryExt creates a new  LedgerEntryExt.
func NewLedgerEntryExt(v LedgerVersion, value interface{}) (result LedgerEntryExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// LedgerEntry is an XDR Struct defines as:
//
//   struct LedgerEntry
//    {
//        uint32 lastModifiedLedgerSeq; // ledger the LedgerEntry was last changed
//
//        union switch (LedgerEntryType type)
//        {
//        case ACCOUNT:
//            AccountEntry account;
//    	case COINS_EMISSION_REQUEST:
//    		CoinsEmissionRequestEntry coinsEmissionRequest;
//        case FEE:
//            FeeEntry feeState;
//        case COINS_EMISSION:
//    		CoinsEmissionEntry coinsEmission;
//        case BALANCE:
//            BalanceEntry balance;
//        case PAYMENT_REQUEST:
//            PaymentRequestEntry paymentRequest;
//        case ASSET:
//            AssetEntry asset;
//        case PAYMENT_ENTRY:
//            PaymentEntry payment;
//        case ACCOUNT_TYPE_LIMITS:
//            AccountTypeLimitsEntry accountTypeLimits;
//        case STATISTICS:
//            StatisticsEntry stats;
//        case EXCHANGE_DATA:
//            ExchangeDataEntry exchangeData;
//        case EXCHANGE_POLICIES:
//            ExchangePoliciesEntry exchangePolicies;
//        case TRUST:
//            TrustEntry trust;
//        case ACCOUNT_LIMITS:
//            AccountLimitsEntry accountLimits;
//    	case ASSET_PAIR:
//    		AssetPairEntry assetPair;
//    	case OFFER_ENTRY:
//    		OfferEntry offer;
//        case INVOICE:
//            InvoiceEntry invoice;
//        }
//        data;
//
//        // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type LedgerEntry struct {
	LastModifiedLedgerSeq Uint32          `json:"lastModifiedLedgerSeq,omitempty"`
	Data                  LedgerEntryData `json:"data,omitempty"`
	Ext                   LedgerEntryExt  `json:"ext,omitempty"`
}

// EnvelopeType is an XDR Enum defines as:
//
//   enum EnvelopeType
//    {
//        ENVELOPE_TYPE_SCP = 1,
//        ENVELOPE_TYPE_TX = 2,
//        ENVELOPE_TYPE_AUTH = 3
//    };
//
type EnvelopeType int32

const (
	EnvelopeTypeEnvelopeTypeScp  EnvelopeType = 1
	EnvelopeTypeEnvelopeTypeTx   EnvelopeType = 2
	EnvelopeTypeEnvelopeTypeAuth EnvelopeType = 3
)

var EnvelopeTypeAll = []EnvelopeType{
	EnvelopeTypeEnvelopeTypeScp,
	EnvelopeTypeEnvelopeTypeTx,
	EnvelopeTypeEnvelopeTypeAuth,
}

var envelopeTypeMap = map[int32]string{
	1: "EnvelopeTypeEnvelopeTypeScp",
	2: "EnvelopeTypeEnvelopeTypeTx",
	3: "EnvelopeTypeEnvelopeTypeAuth",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for EnvelopeType
func (e EnvelopeType) ValidEnum(v int32) bool {
	_, ok := envelopeTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e EnvelopeType) String() string {
	name, _ := envelopeTypeMap[int32(e)]
	return name
}

func (e EnvelopeType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// UpgradeType is an XDR Typedef defines as:
//
//   typedef opaque UpgradeType<128>;
//
type UpgradeType []byte

// StellarValueExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type StellarValueExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u StellarValueExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of StellarValueExt
func (u StellarValueExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewStellarValueExt creates a new  StellarValueExt.
func NewStellarValueExt(v LedgerVersion, value interface{}) (result StellarValueExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// StellarValue is an XDR Struct defines as:
//
//   struct StellarValue
//    {
//        Hash txSetHash;   // transaction set to apply to previous ledger
//        uint64 closeTime; // network close time
//
//        // upgrades to apply to the previous ledger (usually empty)
//        // this is a vector of encoded 'LedgerUpgrade' so that nodes can drop
//        // unknown steps during consensus if needed.
//        // see notes below on 'LedgerUpgrade' for more detail
//        // max size is dictated by number of upgrade types (+ room for future)
//        UpgradeType upgrades<6>;
//
//        // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type StellarValue struct {
	TxSetHash Hash            `json:"txSetHash,omitempty"`
	CloseTime Uint64          `json:"closeTime,omitempty"`
	Upgrades  []UpgradeType   `json:"upgrades,omitempty" xdrmaxsize:"6"`
	Ext       StellarValueExt `json:"ext,omitempty"`
}

// LedgerHeaderExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type LedgerHeaderExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerHeaderExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerHeaderExt
func (u LedgerHeaderExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewLedgerHeaderExt creates a new  LedgerHeaderExt.
func NewLedgerHeaderExt(v LedgerVersion, value interface{}) (result LedgerHeaderExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// LedgerHeader is an XDR Struct defines as:
//
//   struct LedgerHeader
//    {
//        uint32 ledgerVersion;    // the protocol version of the ledger
//        Hash previousLedgerHash; // hash of the previous ledger header
//        StellarValue scpValue;   // what consensus agreed to
//        Hash txSetResultHash;    // the TransactionResultSet that led to this ledger
//        Hash bucketListHash;     // hash of the ledger state
//
//        uint32 ledgerSeq; // sequence number of this ledger
//
//        int64 storageFeePeriod;
//    	int64 payoutsPeriod;
//
//        uint64 idPool; // last used global ID, used for generating objects
//
//        uint32 baseFee;     // base fee per operation in stroops
//        uint32 baseReserve; // account base reserve in stroops
//
//        uint32 maxTxSetSize; // maximum size a transaction set can be
//
//        PublicKey issuanceKeys<>;
//        int64 txExpirationPeriod;
//
//        Hash skipList[4]; // hashes of ledgers in the past. allows you to jump back
//                          // in time without walking the chain back ledger by ledger
//                          // each slot contains the oldest ledger that is mod of
//                          // either 50  5000  50000 or 500000 depending on index
//                          // skipList[0] mod(50), skipList[1] mod(5000), etc
//
//        // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type LedgerHeader struct {
	LedgerVersion      Uint32          `json:"ledgerVersion,omitempty"`
	PreviousLedgerHash Hash            `json:"previousLedgerHash,omitempty"`
	ScpValue           StellarValue    `json:"scpValue,omitempty"`
	TxSetResultHash    Hash            `json:"txSetResultHash,omitempty"`
	BucketListHash     Hash            `json:"bucketListHash,omitempty"`
	LedgerSeq          Uint32          `json:"ledgerSeq,omitempty"`
	StorageFeePeriod   Int64           `json:"storageFeePeriod,omitempty"`
	PayoutsPeriod      Int64           `json:"payoutsPeriod,omitempty"`
	IdPool             Uint64          `json:"idPool,omitempty"`
	BaseFee            Uint32          `json:"baseFee,omitempty"`
	BaseReserve        Uint32          `json:"baseReserve,omitempty"`
	MaxTxSetSize       Uint32          `json:"maxTxSetSize,omitempty"`
	IssuanceKeys       []PublicKey     `json:"issuanceKeys,omitempty"`
	TxExpirationPeriod Int64           `json:"txExpirationPeriod,omitempty"`
	SkipList           [4]Hash         `json:"skipList,omitempty"`
	Ext                LedgerHeaderExt `json:"ext,omitempty"`
}

// LedgerUpgradeType is an XDR Enum defines as:
//
//   enum LedgerUpgradeType
//    {
//        LEDGER_UPGRADE_VERSION = 1,
//        LEDGER_UPGRADE_MAX_TX_SET_SIZE = 2,
//        LEDGER_UPGRADE_ISSUANCE_KEYS = 3,
//        LEDGER_UPGRADE_TX_EXPIRATION_PERIOD = 4
//    };
//
type LedgerUpgradeType int32

const (
	LedgerUpgradeTypeLedgerUpgradeVersion            LedgerUpgradeType = 1
	LedgerUpgradeTypeLedgerUpgradeMaxTxSetSize       LedgerUpgradeType = 2
	LedgerUpgradeTypeLedgerUpgradeIssuanceKeys       LedgerUpgradeType = 3
	LedgerUpgradeTypeLedgerUpgradeTxExpirationPeriod LedgerUpgradeType = 4
)

var LedgerUpgradeTypeAll = []LedgerUpgradeType{
	LedgerUpgradeTypeLedgerUpgradeVersion,
	LedgerUpgradeTypeLedgerUpgradeMaxTxSetSize,
	LedgerUpgradeTypeLedgerUpgradeIssuanceKeys,
	LedgerUpgradeTypeLedgerUpgradeTxExpirationPeriod,
}

var ledgerUpgradeTypeMap = map[int32]string{
	1: "LedgerUpgradeTypeLedgerUpgradeVersion",
	2: "LedgerUpgradeTypeLedgerUpgradeMaxTxSetSize",
	3: "LedgerUpgradeTypeLedgerUpgradeIssuanceKeys",
	4: "LedgerUpgradeTypeLedgerUpgradeTxExpirationPeriod",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for LedgerUpgradeType
func (e LedgerUpgradeType) ValidEnum(v int32) bool {
	_, ok := ledgerUpgradeTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e LedgerUpgradeType) String() string {
	name, _ := ledgerUpgradeTypeMap[int32(e)]
	return name
}

func (e LedgerUpgradeType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// LedgerUpgrade is an XDR Union defines as:
//
//   union LedgerUpgrade switch (LedgerUpgradeType type)
//    {
//    case LEDGER_UPGRADE_VERSION:
//        uint32 newLedgerVersion; // update ledgerVersion
//    case LEDGER_UPGRADE_MAX_TX_SET_SIZE:
//        uint32 newMaxTxSetSize; // update maxTxSetSize
//    case LEDGER_UPGRADE_ISSUANCE_KEYS:
//        PublicKey newIssuanceKeys<>;
//    case LEDGER_UPGRADE_TX_EXPIRATION_PERIOD:
//        int64 newTxExpirationPeriod;
//    };
//
type LedgerUpgrade struct {
	Type                  LedgerUpgradeType `json:"type,omitempty"`
	NewLedgerVersion      *Uint32           `json:"newLedgerVersion,omitempty"`
	NewMaxTxSetSize       *Uint32           `json:"newMaxTxSetSize,omitempty"`
	NewIssuanceKeys       *[]PublicKey      `json:"newIssuanceKeys,omitempty"`
	NewTxExpirationPeriod *Int64            `json:"newTxExpirationPeriod,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerUpgrade) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerUpgrade
func (u LedgerUpgrade) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerUpgradeType(sw) {
	case LedgerUpgradeTypeLedgerUpgradeVersion:
		return "NewLedgerVersion", true
	case LedgerUpgradeTypeLedgerUpgradeMaxTxSetSize:
		return "NewMaxTxSetSize", true
	case LedgerUpgradeTypeLedgerUpgradeIssuanceKeys:
		return "NewIssuanceKeys", true
	case LedgerUpgradeTypeLedgerUpgradeTxExpirationPeriod:
		return "NewTxExpirationPeriod", true
	}
	return "-", false
}

// NewLedgerUpgrade creates a new  LedgerUpgrade.
func NewLedgerUpgrade(aType LedgerUpgradeType, value interface{}) (result LedgerUpgrade, err error) {
	result.Type = aType
	switch LedgerUpgradeType(aType) {
	case LedgerUpgradeTypeLedgerUpgradeVersion:
		tv, ok := value.(Uint32)
		if !ok {
			err = fmt.Errorf("invalid value, must be Uint32")
			return
		}
		result.NewLedgerVersion = &tv
	case LedgerUpgradeTypeLedgerUpgradeMaxTxSetSize:
		tv, ok := value.(Uint32)
		if !ok {
			err = fmt.Errorf("invalid value, must be Uint32")
			return
		}
		result.NewMaxTxSetSize = &tv
	case LedgerUpgradeTypeLedgerUpgradeIssuanceKeys:
		tv, ok := value.([]PublicKey)
		if !ok {
			err = fmt.Errorf("invalid value, must be []PublicKey")
			return
		}
		result.NewIssuanceKeys = &tv
	case LedgerUpgradeTypeLedgerUpgradeTxExpirationPeriod:
		tv, ok := value.(Int64)
		if !ok {
			err = fmt.Errorf("invalid value, must be Int64")
			return
		}
		result.NewTxExpirationPeriod = &tv
	}
	return
}

// MustNewLedgerVersion retrieves the NewLedgerVersion value from the union,
// panicing if the value is not set.
func (u LedgerUpgrade) MustNewLedgerVersion() Uint32 {
	val, ok := u.GetNewLedgerVersion()

	if !ok {
		panic("arm NewLedgerVersion is not set")
	}

	return val
}

// GetNewLedgerVersion retrieves the NewLedgerVersion value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerUpgrade) GetNewLedgerVersion() (result Uint32, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "NewLedgerVersion" {
		result = *u.NewLedgerVersion
		ok = true
	}

	return
}

// MustNewMaxTxSetSize retrieves the NewMaxTxSetSize value from the union,
// panicing if the value is not set.
func (u LedgerUpgrade) MustNewMaxTxSetSize() Uint32 {
	val, ok := u.GetNewMaxTxSetSize()

	if !ok {
		panic("arm NewMaxTxSetSize is not set")
	}

	return val
}

// GetNewMaxTxSetSize retrieves the NewMaxTxSetSize value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerUpgrade) GetNewMaxTxSetSize() (result Uint32, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "NewMaxTxSetSize" {
		result = *u.NewMaxTxSetSize
		ok = true
	}

	return
}

// MustNewIssuanceKeys retrieves the NewIssuanceKeys value from the union,
// panicing if the value is not set.
func (u LedgerUpgrade) MustNewIssuanceKeys() []PublicKey {
	val, ok := u.GetNewIssuanceKeys()

	if !ok {
		panic("arm NewIssuanceKeys is not set")
	}

	return val
}

// GetNewIssuanceKeys retrieves the NewIssuanceKeys value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerUpgrade) GetNewIssuanceKeys() (result []PublicKey, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "NewIssuanceKeys" {
		result = *u.NewIssuanceKeys
		ok = true
	}

	return
}

// MustNewTxExpirationPeriod retrieves the NewTxExpirationPeriod value from the union,
// panicing if the value is not set.
func (u LedgerUpgrade) MustNewTxExpirationPeriod() Int64 {
	val, ok := u.GetNewTxExpirationPeriod()

	if !ok {
		panic("arm NewTxExpirationPeriod is not set")
	}

	return val
}

// GetNewTxExpirationPeriod retrieves the NewTxExpirationPeriod value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerUpgrade) GetNewTxExpirationPeriod() (result Int64, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "NewTxExpirationPeriod" {
		result = *u.NewTxExpirationPeriod
		ok = true
	}

	return
}

// LedgerKeyAccountExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type LedgerKeyAccountExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerKeyAccountExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerKeyAccountExt
func (u LedgerKeyAccountExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewLedgerKeyAccountExt creates a new  LedgerKeyAccountExt.
func NewLedgerKeyAccountExt(v LedgerVersion, value interface{}) (result LedgerKeyAccountExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// LedgerKeyAccount is an XDR NestedStruct defines as:
//
//   struct
//        {
//            AccountID accountID;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        }
//
type LedgerKeyAccount struct {
	AccountId AccountId           `json:"accountID,omitempty"`
	Ext       LedgerKeyAccountExt `json:"ext,omitempty"`
}

// LedgerKeyCoinsEmissionRequestExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type LedgerKeyCoinsEmissionRequestExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerKeyCoinsEmissionRequestExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerKeyCoinsEmissionRequestExt
func (u LedgerKeyCoinsEmissionRequestExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewLedgerKeyCoinsEmissionRequestExt creates a new  LedgerKeyCoinsEmissionRequestExt.
func NewLedgerKeyCoinsEmissionRequestExt(v LedgerVersion, value interface{}) (result LedgerKeyCoinsEmissionRequestExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// LedgerKeyCoinsEmissionRequest is an XDR NestedStruct defines as:
//
//   struct
//    	{
//    		AccountID issuer;
//    		uint64 requestID;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//    	}
//
type LedgerKeyCoinsEmissionRequest struct {
	Issuer    AccountId                        `json:"issuer,omitempty"`
	RequestId Uint64                           `json:"requestID,omitempty"`
	Ext       LedgerKeyCoinsEmissionRequestExt `json:"ext,omitempty"`
}

// LedgerKeyCoinsEmissionExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type LedgerKeyCoinsEmissionExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerKeyCoinsEmissionExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerKeyCoinsEmissionExt
func (u LedgerKeyCoinsEmissionExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewLedgerKeyCoinsEmissionExt creates a new  LedgerKeyCoinsEmissionExt.
func NewLedgerKeyCoinsEmissionExt(v LedgerVersion, value interface{}) (result LedgerKeyCoinsEmissionExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// LedgerKeyCoinsEmission is an XDR NestedStruct defines as:
//
//   struct
//    	{
//    		string64 serialNumber;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//    	}
//
type LedgerKeyCoinsEmission struct {
	SerialNumber String64                  `json:"serialNumber,omitempty"`
	Ext          LedgerKeyCoinsEmissionExt `json:"ext,omitempty"`
}

// LedgerKeyFeeStateExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type LedgerKeyFeeStateExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerKeyFeeStateExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerKeyFeeStateExt
func (u LedgerKeyFeeStateExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewLedgerKeyFeeStateExt creates a new  LedgerKeyFeeStateExt.
func NewLedgerKeyFeeStateExt(v LedgerVersion, value interface{}) (result LedgerKeyFeeStateExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// LedgerKeyFeeState is an XDR NestedStruct defines as:
//
//   struct {
//            Hash hash;
//    		int64 lowerBound;
//    		int64 upperBound;
//    		 union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        }
//
type LedgerKeyFeeState struct {
	Hash       Hash                 `json:"hash,omitempty"`
	LowerBound Int64                `json:"lowerBound,omitempty"`
	UpperBound Int64                `json:"upperBound,omitempty"`
	Ext        LedgerKeyFeeStateExt `json:"ext,omitempty"`
}

// LedgerKeyBalanceExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type LedgerKeyBalanceExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerKeyBalanceExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerKeyBalanceExt
func (u LedgerKeyBalanceExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewLedgerKeyBalanceExt creates a new  LedgerKeyBalanceExt.
func NewLedgerKeyBalanceExt(v LedgerVersion, value interface{}) (result LedgerKeyBalanceExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// LedgerKeyBalance is an XDR NestedStruct defines as:
//
//   struct
//        {
//    		BalanceID balanceID;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        }
//
type LedgerKeyBalance struct {
	BalanceId BalanceId           `json:"balanceID,omitempty"`
	Ext       LedgerKeyBalanceExt `json:"ext,omitempty"`
}

// LedgerKeyPaymentRequestExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type LedgerKeyPaymentRequestExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerKeyPaymentRequestExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerKeyPaymentRequestExt
func (u LedgerKeyPaymentRequestExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewLedgerKeyPaymentRequestExt creates a new  LedgerKeyPaymentRequestExt.
func NewLedgerKeyPaymentRequestExt(v LedgerVersion, value interface{}) (result LedgerKeyPaymentRequestExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// LedgerKeyPaymentRequest is an XDR NestedStruct defines as:
//
//   struct
//        {
//    		uint64 paymentID;
//            AccountID exchange;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        }
//
type LedgerKeyPaymentRequest struct {
	PaymentId Uint64                     `json:"paymentID,omitempty"`
	Exchange  AccountId                  `json:"exchange,omitempty"`
	Ext       LedgerKeyPaymentRequestExt `json:"ext,omitempty"`
}

// LedgerKeyAssetExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type LedgerKeyAssetExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerKeyAssetExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerKeyAssetExt
func (u LedgerKeyAssetExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewLedgerKeyAssetExt creates a new  LedgerKeyAssetExt.
func NewLedgerKeyAssetExt(v LedgerVersion, value interface{}) (result LedgerKeyAssetExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// LedgerKeyAsset is an XDR NestedStruct defines as:
//
//   struct
//        {
//    		AssetCode code;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        }
//
type LedgerKeyAsset struct {
	Code AssetCode         `json:"code,omitempty"`
	Ext  LedgerKeyAssetExt `json:"ext,omitempty"`
}

// LedgerKeyPaymentExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type LedgerKeyPaymentExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerKeyPaymentExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerKeyPaymentExt
func (u LedgerKeyPaymentExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewLedgerKeyPaymentExt creates a new  LedgerKeyPaymentExt.
func NewLedgerKeyPaymentExt(v LedgerVersion, value interface{}) (result LedgerKeyPaymentExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// LedgerKeyPayment is an XDR NestedStruct defines as:
//
//   struct
//        {
//    		string64 reference;
//            AccountID exchange;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        }
//
type LedgerKeyPayment struct {
	Reference String64            `json:"reference,omitempty"`
	Exchange  AccountId           `json:"exchange,omitempty"`
	Ext       LedgerKeyPaymentExt `json:"ext,omitempty"`
}

// LedgerKeyAccountTypeLimitsExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type LedgerKeyAccountTypeLimitsExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerKeyAccountTypeLimitsExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerKeyAccountTypeLimitsExt
func (u LedgerKeyAccountTypeLimitsExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewLedgerKeyAccountTypeLimitsExt creates a new  LedgerKeyAccountTypeLimitsExt.
func NewLedgerKeyAccountTypeLimitsExt(v LedgerVersion, value interface{}) (result LedgerKeyAccountTypeLimitsExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// LedgerKeyAccountTypeLimits is an XDR NestedStruct defines as:
//
//   struct {
//            AccountType accountType;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        }
//
type LedgerKeyAccountTypeLimits struct {
	AccountType AccountType                   `json:"accountType,omitempty"`
	Ext         LedgerKeyAccountTypeLimitsExt `json:"ext,omitempty"`
}

// LedgerKeyStatsExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type LedgerKeyStatsExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerKeyStatsExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerKeyStatsExt
func (u LedgerKeyStatsExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewLedgerKeyStatsExt creates a new  LedgerKeyStatsExt.
func NewLedgerKeyStatsExt(v LedgerVersion, value interface{}) (result LedgerKeyStatsExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// LedgerKeyStats is an XDR NestedStruct defines as:
//
//   struct {
//            AccountID accountID;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        }
//
type LedgerKeyStats struct {
	AccountId AccountId         `json:"accountID,omitempty"`
	Ext       LedgerKeyStatsExt `json:"ext,omitempty"`
}

// LedgerKeyExchangeDataExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type LedgerKeyExchangeDataExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerKeyExchangeDataExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerKeyExchangeDataExt
func (u LedgerKeyExchangeDataExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewLedgerKeyExchangeDataExt creates a new  LedgerKeyExchangeDataExt.
func NewLedgerKeyExchangeDataExt(v LedgerVersion, value interface{}) (result LedgerKeyExchangeDataExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// LedgerKeyExchangeData is an XDR NestedStruct defines as:
//
//   struct {
//            AccountID accountID;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        }
//
type LedgerKeyExchangeData struct {
	AccountId AccountId                `json:"accountID,omitempty"`
	Ext       LedgerKeyExchangeDataExt `json:"ext,omitempty"`
}

// LedgerKeyExchangePoliciesExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type LedgerKeyExchangePoliciesExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerKeyExchangePoliciesExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerKeyExchangePoliciesExt
func (u LedgerKeyExchangePoliciesExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewLedgerKeyExchangePoliciesExt creates a new  LedgerKeyExchangePoliciesExt.
func NewLedgerKeyExchangePoliciesExt(v LedgerVersion, value interface{}) (result LedgerKeyExchangePoliciesExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// LedgerKeyExchangePolicies is an XDR NestedStruct defines as:
//
//   struct {
//            AccountID accountID;
//            AssetCode asset;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        }
//
type LedgerKeyExchangePolicies struct {
	AccountId AccountId                    `json:"accountID,omitempty"`
	Asset     AssetCode                    `json:"asset,omitempty"`
	Ext       LedgerKeyExchangePoliciesExt `json:"ext,omitempty"`
}

// LedgerKeyTrustExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type LedgerKeyTrustExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerKeyTrustExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerKeyTrustExt
func (u LedgerKeyTrustExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewLedgerKeyTrustExt creates a new  LedgerKeyTrustExt.
func NewLedgerKeyTrustExt(v LedgerVersion, value interface{}) (result LedgerKeyTrustExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// LedgerKeyTrust is an XDR NestedStruct defines as:
//
//   struct {
//            AccountID allowedAccount;
//            BalanceID balanceToUse;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        }
//
type LedgerKeyTrust struct {
	AllowedAccount AccountId         `json:"allowedAccount,omitempty"`
	BalanceToUse   BalanceId         `json:"balanceToUse,omitempty"`
	Ext            LedgerKeyTrustExt `json:"ext,omitempty"`
}

// LedgerKeyAccountLimitsExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type LedgerKeyAccountLimitsExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerKeyAccountLimitsExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerKeyAccountLimitsExt
func (u LedgerKeyAccountLimitsExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewLedgerKeyAccountLimitsExt creates a new  LedgerKeyAccountLimitsExt.
func NewLedgerKeyAccountLimitsExt(v LedgerVersion, value interface{}) (result LedgerKeyAccountLimitsExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// LedgerKeyAccountLimits is an XDR NestedStruct defines as:
//
//   struct {
//            AccountID accountID;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        }
//
type LedgerKeyAccountLimits struct {
	AccountId AccountId                 `json:"accountID,omitempty"`
	Ext       LedgerKeyAccountLimitsExt `json:"ext,omitempty"`
}

// LedgerKeyAssetPairExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type LedgerKeyAssetPairExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerKeyAssetPairExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerKeyAssetPairExt
func (u LedgerKeyAssetPairExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewLedgerKeyAssetPairExt creates a new  LedgerKeyAssetPairExt.
func NewLedgerKeyAssetPairExt(v LedgerVersion, value interface{}) (result LedgerKeyAssetPairExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// LedgerKeyAssetPair is an XDR NestedStruct defines as:
//
//   struct {
//             AssetCode base;
//    		 AssetCode quote;
//    		 union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        }
//
type LedgerKeyAssetPair struct {
	Base  AssetCode             `json:"base,omitempty"`
	Quote AssetCode             `json:"quote,omitempty"`
	Ext   LedgerKeyAssetPairExt `json:"ext,omitempty"`
}

// LedgerKeyOffer is an XDR NestedStruct defines as:
//
//   struct {
//    		uint64 offerID;
//    		AccountID ownerID;
//    	}
//
type LedgerKeyOffer struct {
	OfferId Uint64    `json:"offerID,omitempty"`
	OwnerId AccountId `json:"ownerID,omitempty"`
}

// LedgerKeyInvoiceExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type LedgerKeyInvoiceExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerKeyInvoiceExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerKeyInvoiceExt
func (u LedgerKeyInvoiceExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewLedgerKeyInvoiceExt creates a new  LedgerKeyInvoiceExt.
func NewLedgerKeyInvoiceExt(v LedgerVersion, value interface{}) (result LedgerKeyInvoiceExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// LedgerKeyInvoice is an XDR NestedStruct defines as:
//
//   struct {
//            uint64 invoiceID;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        }
//
type LedgerKeyInvoice struct {
	InvoiceId Uint64              `json:"invoiceID,omitempty"`
	Ext       LedgerKeyInvoiceExt `json:"ext,omitempty"`
}

// LedgerKey is an XDR Union defines as:
//
//   union LedgerKey switch (LedgerEntryType type)
//    {
//    case ACCOUNT:
//        struct
//        {
//            AccountID accountID;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        } account;
//    case COINS_EMISSION_REQUEST:
//    	struct
//    	{
//    		AccountID issuer;
//    		uint64 requestID;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//    	} coinsEmissionRequest;
//    case COINS_EMISSION:
//    	struct
//    	{
//    		string64 serialNumber;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//    	} coinsEmission;
//    case FEE:
//        struct {
//            Hash hash;
//    		int64 lowerBound;
//    		int64 upperBound;
//    		 union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        } feeState;
//    case BALANCE:
//        struct
//        {
//    		BalanceID balanceID;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        } balance;
//    case PAYMENT_REQUEST:
//        struct
//        {
//    		uint64 paymentID;
//            AccountID exchange;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        } paymentRequest;
//    case ASSET:
//        struct
//        {
//    		AssetCode code;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        } asset;
//    case PAYMENT_ENTRY:
//        struct
//        {
//    		string64 reference;
//            AccountID exchange;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        } payment;
//    case ACCOUNT_TYPE_LIMITS:
//        struct {
//            AccountType accountType;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        } accountTypeLimits;
//    case STATISTICS:
//        struct {
//            AccountID accountID;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        } stats;
//    case EXCHANGE_DATA:
//        struct {
//            AccountID accountID;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        } exchangeData;
//    case EXCHANGE_POLICIES:
//        struct {
//            AccountID accountID;
//            AssetCode asset;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        } exchangePolicies;
//    case TRUST:
//        struct {
//            AccountID allowedAccount;
//            BalanceID balanceToUse;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        } trust;
//    case ACCOUNT_LIMITS:
//        struct {
//            AccountID accountID;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        } accountLimits;
//    case ASSET_PAIR:
//    	struct {
//             AssetCode base;
//    		 AssetCode quote;
//    		 union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        } assetPair;
//    case OFFER_ENTRY:
//    	struct {
//    		uint64 offerID;
//    		AccountID ownerID;
//    	} offer;
//    case INVOICE:
//        struct {
//            uint64 invoiceID;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        } invoice;
//    };
//
type LedgerKey struct {
	Type                 LedgerEntryType                `json:"type,omitempty"`
	Account              *LedgerKeyAccount              `json:"account,omitempty"`
	CoinsEmissionRequest *LedgerKeyCoinsEmissionRequest `json:"coinsEmissionRequest,omitempty"`
	CoinsEmission        *LedgerKeyCoinsEmission        `json:"coinsEmission,omitempty"`
	FeeState             *LedgerKeyFeeState             `json:"feeState,omitempty"`
	Balance              *LedgerKeyBalance              `json:"balance,omitempty"`
	PaymentRequest       *LedgerKeyPaymentRequest       `json:"paymentRequest,omitempty"`
	Asset                *LedgerKeyAsset                `json:"asset,omitempty"`
	Payment              *LedgerKeyPayment              `json:"payment,omitempty"`
	AccountTypeLimits    *LedgerKeyAccountTypeLimits    `json:"accountTypeLimits,omitempty"`
	Stats                *LedgerKeyStats                `json:"stats,omitempty"`
	ExchangeData         *LedgerKeyExchangeData         `json:"exchangeData,omitempty"`
	ExchangePolicies     *LedgerKeyExchangePolicies     `json:"exchangePolicies,omitempty"`
	Trust                *LedgerKeyTrust                `json:"trust,omitempty"`
	AccountLimits        *LedgerKeyAccountLimits        `json:"accountLimits,omitempty"`
	AssetPair            *LedgerKeyAssetPair            `json:"assetPair,omitempty"`
	Offer                *LedgerKeyOffer                `json:"offer,omitempty"`
	Invoice              *LedgerKeyInvoice              `json:"invoice,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerKey) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerKey
func (u LedgerKey) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerEntryType(sw) {
	case LedgerEntryTypeAccount:
		return "Account", true
	case LedgerEntryTypeCoinsEmissionRequest:
		return "CoinsEmissionRequest", true
	case LedgerEntryTypeCoinsEmission:
		return "CoinsEmission", true
	case LedgerEntryTypeFee:
		return "FeeState", true
	case LedgerEntryTypeBalance:
		return "Balance", true
	case LedgerEntryTypePaymentRequest:
		return "PaymentRequest", true
	case LedgerEntryTypeAsset:
		return "Asset", true
	case LedgerEntryTypePaymentEntry:
		return "Payment", true
	case LedgerEntryTypeAccountTypeLimits:
		return "AccountTypeLimits", true
	case LedgerEntryTypeStatistics:
		return "Stats", true
	case LedgerEntryTypeExchangeData:
		return "ExchangeData", true
	case LedgerEntryTypeExchangePolicies:
		return "ExchangePolicies", true
	case LedgerEntryTypeTrust:
		return "Trust", true
	case LedgerEntryTypeAccountLimits:
		return "AccountLimits", true
	case LedgerEntryTypeAssetPair:
		return "AssetPair", true
	case LedgerEntryTypeOfferEntry:
		return "Offer", true
	case LedgerEntryTypeInvoice:
		return "Invoice", true
	}
	return "-", false
}

// NewLedgerKey creates a new  LedgerKey.
func NewLedgerKey(aType LedgerEntryType, value interface{}) (result LedgerKey, err error) {
	result.Type = aType
	switch LedgerEntryType(aType) {
	case LedgerEntryTypeAccount:
		tv, ok := value.(LedgerKeyAccount)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerKeyAccount")
			return
		}
		result.Account = &tv
	case LedgerEntryTypeCoinsEmissionRequest:
		tv, ok := value.(LedgerKeyCoinsEmissionRequest)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerKeyCoinsEmissionRequest")
			return
		}
		result.CoinsEmissionRequest = &tv
	case LedgerEntryTypeCoinsEmission:
		tv, ok := value.(LedgerKeyCoinsEmission)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerKeyCoinsEmission")
			return
		}
		result.CoinsEmission = &tv
	case LedgerEntryTypeFee:
		tv, ok := value.(LedgerKeyFeeState)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerKeyFeeState")
			return
		}
		result.FeeState = &tv
	case LedgerEntryTypeBalance:
		tv, ok := value.(LedgerKeyBalance)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerKeyBalance")
			return
		}
		result.Balance = &tv
	case LedgerEntryTypePaymentRequest:
		tv, ok := value.(LedgerKeyPaymentRequest)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerKeyPaymentRequest")
			return
		}
		result.PaymentRequest = &tv
	case LedgerEntryTypeAsset:
		tv, ok := value.(LedgerKeyAsset)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerKeyAsset")
			return
		}
		result.Asset = &tv
	case LedgerEntryTypePaymentEntry:
		tv, ok := value.(LedgerKeyPayment)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerKeyPayment")
			return
		}
		result.Payment = &tv
	case LedgerEntryTypeAccountTypeLimits:
		tv, ok := value.(LedgerKeyAccountTypeLimits)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerKeyAccountTypeLimits")
			return
		}
		result.AccountTypeLimits = &tv
	case LedgerEntryTypeStatistics:
		tv, ok := value.(LedgerKeyStats)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerKeyStats")
			return
		}
		result.Stats = &tv
	case LedgerEntryTypeExchangeData:
		tv, ok := value.(LedgerKeyExchangeData)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerKeyExchangeData")
			return
		}
		result.ExchangeData = &tv
	case LedgerEntryTypeExchangePolicies:
		tv, ok := value.(LedgerKeyExchangePolicies)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerKeyExchangePolicies")
			return
		}
		result.ExchangePolicies = &tv
	case LedgerEntryTypeTrust:
		tv, ok := value.(LedgerKeyTrust)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerKeyTrust")
			return
		}
		result.Trust = &tv
	case LedgerEntryTypeAccountLimits:
		tv, ok := value.(LedgerKeyAccountLimits)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerKeyAccountLimits")
			return
		}
		result.AccountLimits = &tv
	case LedgerEntryTypeAssetPair:
		tv, ok := value.(LedgerKeyAssetPair)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerKeyAssetPair")
			return
		}
		result.AssetPair = &tv
	case LedgerEntryTypeOfferEntry:
		tv, ok := value.(LedgerKeyOffer)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerKeyOffer")
			return
		}
		result.Offer = &tv
	case LedgerEntryTypeInvoice:
		tv, ok := value.(LedgerKeyInvoice)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerKeyInvoice")
			return
		}
		result.Invoice = &tv
	}
	return
}

// MustAccount retrieves the Account value from the union,
// panicing if the value is not set.
func (u LedgerKey) MustAccount() LedgerKeyAccount {
	val, ok := u.GetAccount()

	if !ok {
		panic("arm Account is not set")
	}

	return val
}

// GetAccount retrieves the Account value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerKey) GetAccount() (result LedgerKeyAccount, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Account" {
		result = *u.Account
		ok = true
	}

	return
}

// MustCoinsEmissionRequest retrieves the CoinsEmissionRequest value from the union,
// panicing if the value is not set.
func (u LedgerKey) MustCoinsEmissionRequest() LedgerKeyCoinsEmissionRequest {
	val, ok := u.GetCoinsEmissionRequest()

	if !ok {
		panic("arm CoinsEmissionRequest is not set")
	}

	return val
}

// GetCoinsEmissionRequest retrieves the CoinsEmissionRequest value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerKey) GetCoinsEmissionRequest() (result LedgerKeyCoinsEmissionRequest, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "CoinsEmissionRequest" {
		result = *u.CoinsEmissionRequest
		ok = true
	}

	return
}

// MustCoinsEmission retrieves the CoinsEmission value from the union,
// panicing if the value is not set.
func (u LedgerKey) MustCoinsEmission() LedgerKeyCoinsEmission {
	val, ok := u.GetCoinsEmission()

	if !ok {
		panic("arm CoinsEmission is not set")
	}

	return val
}

// GetCoinsEmission retrieves the CoinsEmission value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerKey) GetCoinsEmission() (result LedgerKeyCoinsEmission, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "CoinsEmission" {
		result = *u.CoinsEmission
		ok = true
	}

	return
}

// MustFeeState retrieves the FeeState value from the union,
// panicing if the value is not set.
func (u LedgerKey) MustFeeState() LedgerKeyFeeState {
	val, ok := u.GetFeeState()

	if !ok {
		panic("arm FeeState is not set")
	}

	return val
}

// GetFeeState retrieves the FeeState value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerKey) GetFeeState() (result LedgerKeyFeeState, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "FeeState" {
		result = *u.FeeState
		ok = true
	}

	return
}

// MustBalance retrieves the Balance value from the union,
// panicing if the value is not set.
func (u LedgerKey) MustBalance() LedgerKeyBalance {
	val, ok := u.GetBalance()

	if !ok {
		panic("arm Balance is not set")
	}

	return val
}

// GetBalance retrieves the Balance value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerKey) GetBalance() (result LedgerKeyBalance, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Balance" {
		result = *u.Balance
		ok = true
	}

	return
}

// MustPaymentRequest retrieves the PaymentRequest value from the union,
// panicing if the value is not set.
func (u LedgerKey) MustPaymentRequest() LedgerKeyPaymentRequest {
	val, ok := u.GetPaymentRequest()

	if !ok {
		panic("arm PaymentRequest is not set")
	}

	return val
}

// GetPaymentRequest retrieves the PaymentRequest value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerKey) GetPaymentRequest() (result LedgerKeyPaymentRequest, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "PaymentRequest" {
		result = *u.PaymentRequest
		ok = true
	}

	return
}

// MustAsset retrieves the Asset value from the union,
// panicing if the value is not set.
func (u LedgerKey) MustAsset() LedgerKeyAsset {
	val, ok := u.GetAsset()

	if !ok {
		panic("arm Asset is not set")
	}

	return val
}

// GetAsset retrieves the Asset value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerKey) GetAsset() (result LedgerKeyAsset, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Asset" {
		result = *u.Asset
		ok = true
	}

	return
}

// MustPayment retrieves the Payment value from the union,
// panicing if the value is not set.
func (u LedgerKey) MustPayment() LedgerKeyPayment {
	val, ok := u.GetPayment()

	if !ok {
		panic("arm Payment is not set")
	}

	return val
}

// GetPayment retrieves the Payment value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerKey) GetPayment() (result LedgerKeyPayment, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Payment" {
		result = *u.Payment
		ok = true
	}

	return
}

// MustAccountTypeLimits retrieves the AccountTypeLimits value from the union,
// panicing if the value is not set.
func (u LedgerKey) MustAccountTypeLimits() LedgerKeyAccountTypeLimits {
	val, ok := u.GetAccountTypeLimits()

	if !ok {
		panic("arm AccountTypeLimits is not set")
	}

	return val
}

// GetAccountTypeLimits retrieves the AccountTypeLimits value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerKey) GetAccountTypeLimits() (result LedgerKeyAccountTypeLimits, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "AccountTypeLimits" {
		result = *u.AccountTypeLimits
		ok = true
	}

	return
}

// MustStats retrieves the Stats value from the union,
// panicing if the value is not set.
func (u LedgerKey) MustStats() LedgerKeyStats {
	val, ok := u.GetStats()

	if !ok {
		panic("arm Stats is not set")
	}

	return val
}

// GetStats retrieves the Stats value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerKey) GetStats() (result LedgerKeyStats, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Stats" {
		result = *u.Stats
		ok = true
	}

	return
}

// MustExchangeData retrieves the ExchangeData value from the union,
// panicing if the value is not set.
func (u LedgerKey) MustExchangeData() LedgerKeyExchangeData {
	val, ok := u.GetExchangeData()

	if !ok {
		panic("arm ExchangeData is not set")
	}

	return val
}

// GetExchangeData retrieves the ExchangeData value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerKey) GetExchangeData() (result LedgerKeyExchangeData, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ExchangeData" {
		result = *u.ExchangeData
		ok = true
	}

	return
}

// MustExchangePolicies retrieves the ExchangePolicies value from the union,
// panicing if the value is not set.
func (u LedgerKey) MustExchangePolicies() LedgerKeyExchangePolicies {
	val, ok := u.GetExchangePolicies()

	if !ok {
		panic("arm ExchangePolicies is not set")
	}

	return val
}

// GetExchangePolicies retrieves the ExchangePolicies value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerKey) GetExchangePolicies() (result LedgerKeyExchangePolicies, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ExchangePolicies" {
		result = *u.ExchangePolicies
		ok = true
	}

	return
}

// MustTrust retrieves the Trust value from the union,
// panicing if the value is not set.
func (u LedgerKey) MustTrust() LedgerKeyTrust {
	val, ok := u.GetTrust()

	if !ok {
		panic("arm Trust is not set")
	}

	return val
}

// GetTrust retrieves the Trust value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerKey) GetTrust() (result LedgerKeyTrust, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Trust" {
		result = *u.Trust
		ok = true
	}

	return
}

// MustAccountLimits retrieves the AccountLimits value from the union,
// panicing if the value is not set.
func (u LedgerKey) MustAccountLimits() LedgerKeyAccountLimits {
	val, ok := u.GetAccountLimits()

	if !ok {
		panic("arm AccountLimits is not set")
	}

	return val
}

// GetAccountLimits retrieves the AccountLimits value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerKey) GetAccountLimits() (result LedgerKeyAccountLimits, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "AccountLimits" {
		result = *u.AccountLimits
		ok = true
	}

	return
}

// MustAssetPair retrieves the AssetPair value from the union,
// panicing if the value is not set.
func (u LedgerKey) MustAssetPair() LedgerKeyAssetPair {
	val, ok := u.GetAssetPair()

	if !ok {
		panic("arm AssetPair is not set")
	}

	return val
}

// GetAssetPair retrieves the AssetPair value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerKey) GetAssetPair() (result LedgerKeyAssetPair, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "AssetPair" {
		result = *u.AssetPair
		ok = true
	}

	return
}

// MustOffer retrieves the Offer value from the union,
// panicing if the value is not set.
func (u LedgerKey) MustOffer() LedgerKeyOffer {
	val, ok := u.GetOffer()

	if !ok {
		panic("arm Offer is not set")
	}

	return val
}

// GetOffer retrieves the Offer value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerKey) GetOffer() (result LedgerKeyOffer, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Offer" {
		result = *u.Offer
		ok = true
	}

	return
}

// MustInvoice retrieves the Invoice value from the union,
// panicing if the value is not set.
func (u LedgerKey) MustInvoice() LedgerKeyInvoice {
	val, ok := u.GetInvoice()

	if !ok {
		panic("arm Invoice is not set")
	}

	return val
}

// GetInvoice retrieves the Invoice value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerKey) GetInvoice() (result LedgerKeyInvoice, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Invoice" {
		result = *u.Invoice
		ok = true
	}

	return
}

// BucketEntryType is an XDR Enum defines as:
//
//   enum BucketEntryType
//    {
//        LIVEENTRY = 0,
//        DEADENTRY = 1
//    };
//
type BucketEntryType int32

const (
	BucketEntryTypeLiveentry BucketEntryType = 0
	BucketEntryTypeDeadentry BucketEntryType = 1
)

var BucketEntryTypeAll = []BucketEntryType{
	BucketEntryTypeLiveentry,
	BucketEntryTypeDeadentry,
}

var bucketEntryTypeMap = map[int32]string{
	0: "BucketEntryTypeLiveentry",
	1: "BucketEntryTypeDeadentry",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for BucketEntryType
func (e BucketEntryType) ValidEnum(v int32) bool {
	_, ok := bucketEntryTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e BucketEntryType) String() string {
	name, _ := bucketEntryTypeMap[int32(e)]
	return name
}

func (e BucketEntryType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// BucketEntry is an XDR Union defines as:
//
//   union BucketEntry switch (BucketEntryType type)
//    {
//    case LIVEENTRY:
//        LedgerEntry liveEntry;
//
//    case DEADENTRY:
//        LedgerKey deadEntry;
//    };
//
type BucketEntry struct {
	Type      BucketEntryType `json:"type,omitempty"`
	LiveEntry *LedgerEntry    `json:"liveEntry,omitempty"`
	DeadEntry *LedgerKey      `json:"deadEntry,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u BucketEntry) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of BucketEntry
func (u BucketEntry) ArmForSwitch(sw int32) (string, bool) {
	switch BucketEntryType(sw) {
	case BucketEntryTypeLiveentry:
		return "LiveEntry", true
	case BucketEntryTypeDeadentry:
		return "DeadEntry", true
	}
	return "-", false
}

// NewBucketEntry creates a new  BucketEntry.
func NewBucketEntry(aType BucketEntryType, value interface{}) (result BucketEntry, err error) {
	result.Type = aType
	switch BucketEntryType(aType) {
	case BucketEntryTypeLiveentry:
		tv, ok := value.(LedgerEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerEntry")
			return
		}
		result.LiveEntry = &tv
	case BucketEntryTypeDeadentry:
		tv, ok := value.(LedgerKey)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerKey")
			return
		}
		result.DeadEntry = &tv
	}
	return
}

// MustLiveEntry retrieves the LiveEntry value from the union,
// panicing if the value is not set.
func (u BucketEntry) MustLiveEntry() LedgerEntry {
	val, ok := u.GetLiveEntry()

	if !ok {
		panic("arm LiveEntry is not set")
	}

	return val
}

// GetLiveEntry retrieves the LiveEntry value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u BucketEntry) GetLiveEntry() (result LedgerEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "LiveEntry" {
		result = *u.LiveEntry
		ok = true
	}

	return
}

// MustDeadEntry retrieves the DeadEntry value from the union,
// panicing if the value is not set.
func (u BucketEntry) MustDeadEntry() LedgerKey {
	val, ok := u.GetDeadEntry()

	if !ok {
		panic("arm DeadEntry is not set")
	}

	return val
}

// GetDeadEntry retrieves the DeadEntry value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u BucketEntry) GetDeadEntry() (result LedgerKey, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "DeadEntry" {
		result = *u.DeadEntry
		ok = true
	}

	return
}

// TransactionSet is an XDR Struct defines as:
//
//   struct TransactionSet
//    {
//        Hash previousLedgerHash;
//        TransactionEnvelope txs<>;
//    };
//
type TransactionSet struct {
	PreviousLedgerHash Hash                  `json:"previousLedgerHash,omitempty"`
	Txs                []TransactionEnvelope `json:"txs,omitempty"`
}

// TransactionResultPair is an XDR Struct defines as:
//
//   struct TransactionResultPair
//    {
//        Hash transactionHash;
//        TransactionResult result; // result for the transaction
//    };
//
type TransactionResultPair struct {
	TransactionHash Hash              `json:"transactionHash,omitempty"`
	Result          TransactionResult `json:"result,omitempty"`
}

// TransactionResultSet is an XDR Struct defines as:
//
//   struct TransactionResultSet
//    {
//        TransactionResultPair results<>;
//    };
//
type TransactionResultSet struct {
	Results []TransactionResultPair `json:"results,omitempty"`
}

// TransactionHistoryEntryExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type TransactionHistoryEntryExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u TransactionHistoryEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of TransactionHistoryEntryExt
func (u TransactionHistoryEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewTransactionHistoryEntryExt creates a new  TransactionHistoryEntryExt.
func NewTransactionHistoryEntryExt(v LedgerVersion, value interface{}) (result TransactionHistoryEntryExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// TransactionHistoryEntry is an XDR Struct defines as:
//
//   struct TransactionHistoryEntry
//    {
//        uint32 ledgerSeq;
//        TransactionSet txSet;
//
//        // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type TransactionHistoryEntry struct {
	LedgerSeq Uint32                     `json:"ledgerSeq,omitempty"`
	TxSet     TransactionSet             `json:"txSet,omitempty"`
	Ext       TransactionHistoryEntryExt `json:"ext,omitempty"`
}

// TransactionHistoryResultEntryExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type TransactionHistoryResultEntryExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u TransactionHistoryResultEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of TransactionHistoryResultEntryExt
func (u TransactionHistoryResultEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewTransactionHistoryResultEntryExt creates a new  TransactionHistoryResultEntryExt.
func NewTransactionHistoryResultEntryExt(v LedgerVersion, value interface{}) (result TransactionHistoryResultEntryExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// TransactionHistoryResultEntry is an XDR Struct defines as:
//
//   struct TransactionHistoryResultEntry
//    {
//        uint32 ledgerSeq;
//        TransactionResultSet txResultSet;
//
//        // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type TransactionHistoryResultEntry struct {
	LedgerSeq   Uint32                           `json:"ledgerSeq,omitempty"`
	TxResultSet TransactionResultSet             `json:"txResultSet,omitempty"`
	Ext         TransactionHistoryResultEntryExt `json:"ext,omitempty"`
}

// LedgerHeaderHistoryEntryExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type LedgerHeaderHistoryEntryExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerHeaderHistoryEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerHeaderHistoryEntryExt
func (u LedgerHeaderHistoryEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewLedgerHeaderHistoryEntryExt creates a new  LedgerHeaderHistoryEntryExt.
func NewLedgerHeaderHistoryEntryExt(v LedgerVersion, value interface{}) (result LedgerHeaderHistoryEntryExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// LedgerHeaderHistoryEntry is an XDR Struct defines as:
//
//   struct LedgerHeaderHistoryEntry
//    {
//        Hash hash;
//        LedgerHeader header;
//
//        // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type LedgerHeaderHistoryEntry struct {
	Hash   Hash                        `json:"hash,omitempty"`
	Header LedgerHeader                `json:"header,omitempty"`
	Ext    LedgerHeaderHistoryEntryExt `json:"ext,omitempty"`
}

// LedgerScpMessages is an XDR Struct defines as:
//
//   struct LedgerSCPMessages
//    {
//        uint32 ledgerSeq;
//        SCPEnvelope messages<>;
//    };
//
type LedgerScpMessages struct {
	LedgerSeq Uint32        `json:"ledgerSeq,omitempty"`
	Messages  []ScpEnvelope `json:"messages,omitempty"`
}

// ScpHistoryEntryV0 is an XDR Struct defines as:
//
//   struct SCPHistoryEntryV0
//    {
//        SCPQuorumSet quorumSets<>; // additional quorum sets used by ledgerMessages
//        LedgerSCPMessages ledgerMessages;
//    };
//
type ScpHistoryEntryV0 struct {
	QuorumSets     []ScpQuorumSet    `json:"quorumSets,omitempty"`
	LedgerMessages LedgerScpMessages `json:"ledgerMessages,omitempty"`
}

// ScpHistoryEntry is an XDR Union defines as:
//
//   union SCPHistoryEntry switch (LedgerVersion v)
//    {
//    case EMPTY_VERSION:
//        SCPHistoryEntryV0 v0;
//    };
//
type ScpHistoryEntry struct {
	V  LedgerVersion      `json:"v,omitempty"`
	V0 *ScpHistoryEntryV0 `json:"v0,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ScpHistoryEntry) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ScpHistoryEntry
func (u ScpHistoryEntry) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "V0", true
	}
	return "-", false
}

// NewScpHistoryEntry creates a new  ScpHistoryEntry.
func NewScpHistoryEntry(v LedgerVersion, value interface{}) (result ScpHistoryEntry, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		tv, ok := value.(ScpHistoryEntryV0)
		if !ok {
			err = fmt.Errorf("invalid value, must be ScpHistoryEntryV0")
			return
		}
		result.V0 = &tv
	}
	return
}

// MustV0 retrieves the V0 value from the union,
// panicing if the value is not set.
func (u ScpHistoryEntry) MustV0() ScpHistoryEntryV0 {
	val, ok := u.GetV0()

	if !ok {
		panic("arm V0 is not set")
	}

	return val
}

// GetV0 retrieves the V0 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ScpHistoryEntry) GetV0() (result ScpHistoryEntryV0, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.V))

	if armName == "V0" {
		result = *u.V0
		ok = true
	}

	return
}

// LedgerEntryChangeType is an XDR Enum defines as:
//
//   enum LedgerEntryChangeType
//    {
//        LEDGER_ENTRY_CREATED = 0, // entry was added to the ledger
//        LEDGER_ENTRY_UPDATED = 1, // entry was modified in the ledger
//        LEDGER_ENTRY_REMOVED = 2, // entry was removed from the ledger
//        LEDGER_ENTRY_STATE = 3    // value of the entry
//    };
//
type LedgerEntryChangeType int32

const (
	LedgerEntryChangeTypeLedgerEntryCreated LedgerEntryChangeType = 0
	LedgerEntryChangeTypeLedgerEntryUpdated LedgerEntryChangeType = 1
	LedgerEntryChangeTypeLedgerEntryRemoved LedgerEntryChangeType = 2
	LedgerEntryChangeTypeLedgerEntryState   LedgerEntryChangeType = 3
)

var LedgerEntryChangeTypeAll = []LedgerEntryChangeType{
	LedgerEntryChangeTypeLedgerEntryCreated,
	LedgerEntryChangeTypeLedgerEntryUpdated,
	LedgerEntryChangeTypeLedgerEntryRemoved,
	LedgerEntryChangeTypeLedgerEntryState,
}

var ledgerEntryChangeTypeMap = map[int32]string{
	0: "LedgerEntryChangeTypeLedgerEntryCreated",
	1: "LedgerEntryChangeTypeLedgerEntryUpdated",
	2: "LedgerEntryChangeTypeLedgerEntryRemoved",
	3: "LedgerEntryChangeTypeLedgerEntryState",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for LedgerEntryChangeType
func (e LedgerEntryChangeType) ValidEnum(v int32) bool {
	_, ok := ledgerEntryChangeTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e LedgerEntryChangeType) String() string {
	name, _ := ledgerEntryChangeTypeMap[int32(e)]
	return name
}

func (e LedgerEntryChangeType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// LedgerEntryChange is an XDR Union defines as:
//
//   union LedgerEntryChange switch (LedgerEntryChangeType type)
//    {
//    case LEDGER_ENTRY_CREATED:
//        LedgerEntry created;
//    case LEDGER_ENTRY_UPDATED:
//        LedgerEntry updated;
//    case LEDGER_ENTRY_REMOVED:
//        LedgerKey removed;
//    case LEDGER_ENTRY_STATE:
//        LedgerEntry state;
//    };
//
type LedgerEntryChange struct {
	Type    LedgerEntryChangeType `json:"type,omitempty"`
	Created *LedgerEntry          `json:"created,omitempty"`
	Updated *LedgerEntry          `json:"updated,omitempty"`
	Removed *LedgerKey            `json:"removed,omitempty"`
	State   *LedgerEntry          `json:"state,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerEntryChange) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerEntryChange
func (u LedgerEntryChange) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerEntryChangeType(sw) {
	case LedgerEntryChangeTypeLedgerEntryCreated:
		return "Created", true
	case LedgerEntryChangeTypeLedgerEntryUpdated:
		return "Updated", true
	case LedgerEntryChangeTypeLedgerEntryRemoved:
		return "Removed", true
	case LedgerEntryChangeTypeLedgerEntryState:
		return "State", true
	}
	return "-", false
}

// NewLedgerEntryChange creates a new  LedgerEntryChange.
func NewLedgerEntryChange(aType LedgerEntryChangeType, value interface{}) (result LedgerEntryChange, err error) {
	result.Type = aType
	switch LedgerEntryChangeType(aType) {
	case LedgerEntryChangeTypeLedgerEntryCreated:
		tv, ok := value.(LedgerEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerEntry")
			return
		}
		result.Created = &tv
	case LedgerEntryChangeTypeLedgerEntryUpdated:
		tv, ok := value.(LedgerEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerEntry")
			return
		}
		result.Updated = &tv
	case LedgerEntryChangeTypeLedgerEntryRemoved:
		tv, ok := value.(LedgerKey)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerKey")
			return
		}
		result.Removed = &tv
	case LedgerEntryChangeTypeLedgerEntryState:
		tv, ok := value.(LedgerEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be LedgerEntry")
			return
		}
		result.State = &tv
	}
	return
}

// MustCreated retrieves the Created value from the union,
// panicing if the value is not set.
func (u LedgerEntryChange) MustCreated() LedgerEntry {
	val, ok := u.GetCreated()

	if !ok {
		panic("arm Created is not set")
	}

	return val
}

// GetCreated retrieves the Created value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryChange) GetCreated() (result LedgerEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Created" {
		result = *u.Created
		ok = true
	}

	return
}

// MustUpdated retrieves the Updated value from the union,
// panicing if the value is not set.
func (u LedgerEntryChange) MustUpdated() LedgerEntry {
	val, ok := u.GetUpdated()

	if !ok {
		panic("arm Updated is not set")
	}

	return val
}

// GetUpdated retrieves the Updated value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryChange) GetUpdated() (result LedgerEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Updated" {
		result = *u.Updated
		ok = true
	}

	return
}

// MustRemoved retrieves the Removed value from the union,
// panicing if the value is not set.
func (u LedgerEntryChange) MustRemoved() LedgerKey {
	val, ok := u.GetRemoved()

	if !ok {
		panic("arm Removed is not set")
	}

	return val
}

// GetRemoved retrieves the Removed value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryChange) GetRemoved() (result LedgerKey, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Removed" {
		result = *u.Removed
		ok = true
	}

	return
}

// MustState retrieves the State value from the union,
// panicing if the value is not set.
func (u LedgerEntryChange) MustState() LedgerEntry {
	val, ok := u.GetState()

	if !ok {
		panic("arm State is not set")
	}

	return val
}

// GetState retrieves the State value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryChange) GetState() (result LedgerEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "State" {
		result = *u.State
		ok = true
	}

	return
}

// LedgerEntryChanges is an XDR Typedef defines as:
//
//   typedef LedgerEntryChange LedgerEntryChanges<>;
//
type LedgerEntryChanges []LedgerEntryChange

// OperationMeta is an XDR Struct defines as:
//
//   struct OperationMeta
//    {
//        LedgerEntryChanges changes;
//    };
//
type OperationMeta struct {
	Changes LedgerEntryChanges `json:"changes,omitempty"`
}

// TransactionMeta is an XDR Union defines as:
//
//   union TransactionMeta switch (LedgerVersion v)
//    {
//    case EMPTY_VERSION:
//        OperationMeta operations<>;
//    };
//
type TransactionMeta struct {
	V          LedgerVersion    `json:"v,omitempty"`
	Operations *[]OperationMeta `json:"operations,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u TransactionMeta) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of TransactionMeta
func (u TransactionMeta) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "Operations", true
	}
	return "-", false
}

// NewTransactionMeta creates a new  TransactionMeta.
func NewTransactionMeta(v LedgerVersion, value interface{}) (result TransactionMeta, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		tv, ok := value.([]OperationMeta)
		if !ok {
			err = fmt.Errorf("invalid value, must be []OperationMeta")
			return
		}
		result.Operations = &tv
	}
	return
}

// MustOperations retrieves the Operations value from the union,
// panicing if the value is not set.
func (u TransactionMeta) MustOperations() []OperationMeta {
	val, ok := u.GetOperations()

	if !ok {
		panic("arm Operations is not set")
	}

	return val
}

// GetOperations retrieves the Operations value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u TransactionMeta) GetOperations() (result []OperationMeta, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.V))

	if armName == "Operations" {
		result = *u.Operations
		ok = true
	}

	return
}

// CreateExchangeAction is an XDR Enum defines as:
//
//   enum CreateExchangeAction
//    {
//        EXCHANGE_CREATE = 0,
//        EXCHANGE_UPDATE_POLICIES = 1
//    };
//
type CreateExchangeAction int32

const (
	CreateExchangeActionExchangeCreate         CreateExchangeAction = 0
	CreateExchangeActionExchangeUpdatePolicies CreateExchangeAction = 1
)

var CreateExchangeActionAll = []CreateExchangeAction{
	CreateExchangeActionExchangeCreate,
	CreateExchangeActionExchangeUpdatePolicies,
}

var createExchangeActionMap = map[int32]string{
	0: "CreateExchangeActionExchangeCreate",
	1: "CreateExchangeActionExchangeUpdatePolicies",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for CreateExchangeAction
func (e CreateExchangeAction) ValidEnum(v int32) bool {
	_, ok := createExchangeActionMap[v]
	return ok
}

// String returns the name of `e`
func (e CreateExchangeAction) String() string {
	name, _ := createExchangeActionMap[int32(e)]
	return name
}

func (e CreateExchangeAction) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// ExchangePoliciesExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type ExchangePoliciesExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ExchangePoliciesExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ExchangePoliciesExt
func (u ExchangePoliciesExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewExchangePoliciesExt creates a new  ExchangePoliciesExt.
func NewExchangePoliciesExt(v LedgerVersion, value interface{}) (result ExchangePoliciesExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ExchangePolicies is an XDR Struct defines as:
//
//   struct ExchangePolicies {
//        AssetCode asset;
//        int32 policies;
//
//    	 // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type ExchangePolicies struct {
	Asset    AssetCode           `json:"asset,omitempty"`
	Policies Int32               `json:"policies,omitempty"`
	Ext      ExchangePoliciesExt `json:"ext,omitempty"`
}

// ExchangeDataExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type ExchangeDataExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ExchangeDataExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ExchangeDataExt
func (u ExchangeDataExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewExchangeDataExt creates a new  ExchangeDataExt.
func NewExchangeDataExt(v LedgerVersion, value interface{}) (result ExchangeDataExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ExchangeData is an XDR Struct defines as:
//
//   struct ExchangeData {
//        bool requireReview;
//        string64 name;
//        ExchangePolicies* policies;
//        CreateExchangeAction action;
//
//    	 // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type ExchangeData struct {
	RequireReview bool                 `json:"requireReview,omitempty"`
	Name          String64             `json:"name,omitempty"`
	Policies      *ExchangePolicies    `json:"policies,omitempty"`
	Action        CreateExchangeAction `json:"action,omitempty"`
	Ext           ExchangeDataExt      `json:"ext,omitempty"`
}

// CreateAccountOpDetails is an XDR NestedUnion defines as:
//
//   union switch (AccountType accountType)
//        {
//        case EXCHANGE:
//            ExchangeData exchangeData;
//        default:
//            void;
//        }
//
type CreateAccountOpDetails struct {
	AccountType  AccountType   `json:"accountType,omitempty"`
	ExchangeData *ExchangeData `json:"exchangeData,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u CreateAccountOpDetails) SwitchFieldName() string {
	return "AccountType"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of CreateAccountOpDetails
func (u CreateAccountOpDetails) ArmForSwitch(sw int32) (string, bool) {
	switch AccountType(sw) {
	case AccountTypeExchange:
		return "ExchangeData", true
	default:
		return "", true
	}
}

// NewCreateAccountOpDetails creates a new  CreateAccountOpDetails.
func NewCreateAccountOpDetails(accountType AccountType, value interface{}) (result CreateAccountOpDetails, err error) {
	result.AccountType = accountType
	switch AccountType(accountType) {
	case AccountTypeExchange:
		tv, ok := value.(ExchangeData)
		if !ok {
			err = fmt.Errorf("invalid value, must be ExchangeData")
			return
		}
		result.ExchangeData = &tv
	default:
		// void
	}
	return
}

// MustExchangeData retrieves the ExchangeData value from the union,
// panicing if the value is not set.
func (u CreateAccountOpDetails) MustExchangeData() ExchangeData {
	val, ok := u.GetExchangeData()

	if !ok {
		panic("arm ExchangeData is not set")
	}

	return val
}

// GetExchangeData retrieves the ExchangeData value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u CreateAccountOpDetails) GetExchangeData() (result ExchangeData, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.AccountType))

	if armName == "ExchangeData" {
		result = *u.ExchangeData
		ok = true
	}

	return
}

// CreateAccountOpExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//    	case ACCOUNT_POLICIES:
//    		uint32 policies;
//        }
//
type CreateAccountOpExt struct {
	V        LedgerVersion `json:"v,omitempty"`
	Policies *Uint32       `json:"policies,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u CreateAccountOpExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of CreateAccountOpExt
func (u CreateAccountOpExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	case LedgerVersionAccountPolicies:
		return "Policies", true
	}
	return "-", false
}

// NewCreateAccountOpExt creates a new  CreateAccountOpExt.
func NewCreateAccountOpExt(v LedgerVersion, value interface{}) (result CreateAccountOpExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	case LedgerVersionAccountPolicies:
		tv, ok := value.(Uint32)
		if !ok {
			err = fmt.Errorf("invalid value, must be Uint32")
			return
		}
		result.Policies = &tv
	}
	return
}

// MustPolicies retrieves the Policies value from the union,
// panicing if the value is not set.
func (u CreateAccountOpExt) MustPolicies() Uint32 {
	val, ok := u.GetPolicies()

	if !ok {
		panic("arm Policies is not set")
	}

	return val
}

// GetPolicies retrieves the Policies value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u CreateAccountOpExt) GetPolicies() (result Uint32, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.V))

	if armName == "Policies" {
		result = *u.Policies
		ok = true
	}

	return
}

// CreateAccountOp is an XDR Struct defines as:
//
//   struct CreateAccountOp
//    {
//        AccountID destination; // account to create
//        AccountID* referrer;     // parent account
//
//    	union switch (AccountType accountType)
//        {
//        case EXCHANGE:
//            ExchangeData exchangeData;
//        default:
//            void;
//        } details;
//
//    	 // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//    	case ACCOUNT_POLICIES:
//    		uint32 policies;
//        }
//        ext;
//    };
//
type CreateAccountOp struct {
	Destination AccountId              `json:"destination,omitempty"`
	Referrer    *AccountId             `json:"referrer,omitempty"`
	Details     CreateAccountOpDetails `json:"details,omitempty"`
	Ext         CreateAccountOpExt     `json:"ext,omitempty"`
}

// CreateAccountResultCode is an XDR Enum defines as:
//
//   enum CreateAccountResultCode
//    {
//        // codes considered as "success" for the operation
//        CREATE_ACCOUNT_SUCCESS = 0, // account was created
//
//        // codes considered as "failure" for the operation
//        CREATE_ACCOUNT_MALFORMED = -1,       // invalid destination
//    	CREATE_ACCOUNT_ACCOUNT_TYPE_MISMATCHED = -2, // account already exist and change of account type is not allowed
//    	CREATE_ACCOUNT_TYPE_NOT_ALLOWED = -3, // master or commission account types are not allowed
//        CREATE_ACCOUNT_NAME_DUPLICATION = -4,
//        CREATE_ACCOUNT_REFERRER_NOT_FOUND = -5,
//    	CREATE_ACCOUNT_INVALID_ACCOUNT_VERSION = -6 // if account version is higher than ledger version
//    };
//
type CreateAccountResultCode int32

const (
	CreateAccountResultCodeCreateAccountSuccess               CreateAccountResultCode = 0
	CreateAccountResultCodeCreateAccountMalformed             CreateAccountResultCode = -1
	CreateAccountResultCodeCreateAccountAccountTypeMismatched CreateAccountResultCode = -2
	CreateAccountResultCodeCreateAccountTypeNotAllowed        CreateAccountResultCode = -3
	CreateAccountResultCodeCreateAccountNameDuplication       CreateAccountResultCode = -4
	CreateAccountResultCodeCreateAccountReferrerNotFound      CreateAccountResultCode = -5
	CreateAccountResultCodeCreateAccountInvalidAccountVersion CreateAccountResultCode = -6
)

var CreateAccountResultCodeAll = []CreateAccountResultCode{
	CreateAccountResultCodeCreateAccountSuccess,
	CreateAccountResultCodeCreateAccountMalformed,
	CreateAccountResultCodeCreateAccountAccountTypeMismatched,
	CreateAccountResultCodeCreateAccountTypeNotAllowed,
	CreateAccountResultCodeCreateAccountNameDuplication,
	CreateAccountResultCodeCreateAccountReferrerNotFound,
	CreateAccountResultCodeCreateAccountInvalidAccountVersion,
}

var createAccountResultCodeMap = map[int32]string{
	0:  "CreateAccountResultCodeCreateAccountSuccess",
	-1: "CreateAccountResultCodeCreateAccountMalformed",
	-2: "CreateAccountResultCodeCreateAccountAccountTypeMismatched",
	-3: "CreateAccountResultCodeCreateAccountTypeNotAllowed",
	-4: "CreateAccountResultCodeCreateAccountNameDuplication",
	-5: "CreateAccountResultCodeCreateAccountReferrerNotFound",
	-6: "CreateAccountResultCodeCreateAccountInvalidAccountVersion",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for CreateAccountResultCode
func (e CreateAccountResultCode) ValidEnum(v int32) bool {
	_, ok := createAccountResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e CreateAccountResultCode) String() string {
	name, _ := createAccountResultCodeMap[int32(e)]
	return name
}

func (e CreateAccountResultCode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// CreateAccountSuccessExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type CreateAccountSuccessExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u CreateAccountSuccessExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of CreateAccountSuccessExt
func (u CreateAccountSuccessExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewCreateAccountSuccessExt creates a new  CreateAccountSuccessExt.
func NewCreateAccountSuccessExt(v LedgerVersion, value interface{}) (result CreateAccountSuccessExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// CreateAccountSuccess is an XDR Struct defines as:
//
//   struct CreateAccountSuccess
//    {
//    	int64 referrerFee;
//    	 // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type CreateAccountSuccess struct {
	ReferrerFee Int64                   `json:"referrerFee,omitempty"`
	Ext         CreateAccountSuccessExt `json:"ext,omitempty"`
}

// CreateAccountResult is an XDR Union defines as:
//
//   union CreateAccountResult switch (CreateAccountResultCode code)
//    {
//    case CREATE_ACCOUNT_SUCCESS:
//        CreateAccountSuccess success;
//    default:
//        void;
//    };
//
type CreateAccountResult struct {
	Code    CreateAccountResultCode `json:"code,omitempty"`
	Success *CreateAccountSuccess   `json:"success,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u CreateAccountResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of CreateAccountResult
func (u CreateAccountResult) ArmForSwitch(sw int32) (string, bool) {
	switch CreateAccountResultCode(sw) {
	case CreateAccountResultCodeCreateAccountSuccess:
		return "Success", true
	default:
		return "", true
	}
}

// NewCreateAccountResult creates a new  CreateAccountResult.
func NewCreateAccountResult(code CreateAccountResultCode, value interface{}) (result CreateAccountResult, err error) {
	result.Code = code
	switch CreateAccountResultCode(code) {
	case CreateAccountResultCodeCreateAccountSuccess:
		tv, ok := value.(CreateAccountSuccess)
		if !ok {
			err = fmt.Errorf("invalid value, must be CreateAccountSuccess")
			return
		}
		result.Success = &tv
	default:
		// void
	}
	return
}

// MustSuccess retrieves the Success value from the union,
// panicing if the value is not set.
func (u CreateAccountResult) MustSuccess() CreateAccountSuccess {
	val, ok := u.GetSuccess()

	if !ok {
		panic("arm Success is not set")
	}

	return val
}

// GetSuccess retrieves the Success value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u CreateAccountResult) GetSuccess() (result CreateAccountSuccess, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "Success" {
		result = *u.Success
		ok = true
	}

	return
}

// DemurrageOpExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type DemurrageOpExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u DemurrageOpExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of DemurrageOpExt
func (u DemurrageOpExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewDemurrageOpExt creates a new  DemurrageOpExt.
func NewDemurrageOpExt(v LedgerVersion, value interface{}) (result DemurrageOpExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// DemurrageOp is an XDR Struct defines as:
//
//   struct DemurrageOp
//    {
//        AssetCode asset;
//    	 // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type DemurrageOp struct {
	Asset AssetCode      `json:"asset,omitempty"`
	Ext   DemurrageOpExt `json:"ext,omitempty"`
}

// DemurrageResultCode is an XDR Enum defines as:
//
//   enum DemurrageResultCode
//    {
//        // codes considered as "success" for the operation
//        DEMURRAGE_SUCCESS = 0,
//
//        // codes considered as "failure" for the operation
//        DEMURRAGE_ASSET_NOT_FOUND = -1,
//        DEMURRAGE_INVALID_ASSET = -2,
//        DEMURRAGE_NOT_REQUIRED = -3,
//        DEMURRAGE_STATS_OVERFLOW = -4,
//        DEMURRAGE_LIMITS_EXCEEDED = -5
//    };
//
type DemurrageResultCode int32

const (
	DemurrageResultCodeDemurrageSuccess        DemurrageResultCode = 0
	DemurrageResultCodeDemurrageAssetNotFound  DemurrageResultCode = -1
	DemurrageResultCodeDemurrageInvalidAsset   DemurrageResultCode = -2
	DemurrageResultCodeDemurrageNotRequired    DemurrageResultCode = -3
	DemurrageResultCodeDemurrageStatsOverflow  DemurrageResultCode = -4
	DemurrageResultCodeDemurrageLimitsExceeded DemurrageResultCode = -5
)

var DemurrageResultCodeAll = []DemurrageResultCode{
	DemurrageResultCodeDemurrageSuccess,
	DemurrageResultCodeDemurrageAssetNotFound,
	DemurrageResultCodeDemurrageInvalidAsset,
	DemurrageResultCodeDemurrageNotRequired,
	DemurrageResultCodeDemurrageStatsOverflow,
	DemurrageResultCodeDemurrageLimitsExceeded,
}

var demurrageResultCodeMap = map[int32]string{
	0:  "DemurrageResultCodeDemurrageSuccess",
	-1: "DemurrageResultCodeDemurrageAssetNotFound",
	-2: "DemurrageResultCodeDemurrageInvalidAsset",
	-3: "DemurrageResultCodeDemurrageNotRequired",
	-4: "DemurrageResultCodeDemurrageStatsOverflow",
	-5: "DemurrageResultCodeDemurrageLimitsExceeded",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for DemurrageResultCode
func (e DemurrageResultCode) ValidEnum(v int32) bool {
	_, ok := demurrageResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e DemurrageResultCode) String() string {
	name, _ := demurrageResultCodeMap[int32(e)]
	return name
}

func (e DemurrageResultCode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// BalanceDemurrageExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type BalanceDemurrageExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u BalanceDemurrageExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of BalanceDemurrageExt
func (u BalanceDemurrageExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewBalanceDemurrageExt creates a new  BalanceDemurrageExt.
func NewBalanceDemurrageExt(v LedgerVersion, value interface{}) (result BalanceDemurrageExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// BalanceDemurrage is an XDR Struct defines as:
//
//   struct BalanceDemurrage {
//        BalanceID balance;
//        AccountID account;
//        uint64 amount;
//        uint64 period;
//    	 // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type BalanceDemurrage struct {
	Balance BalanceId           `json:"balance,omitempty"`
	Account AccountId           `json:"account,omitempty"`
	Amount  Uint64              `json:"amount,omitempty"`
	Period  Uint64              `json:"period,omitempty"`
	Ext     BalanceDemurrageExt `json:"ext,omitempty"`
}

// PaymentRequestInfoExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type PaymentRequestInfoExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u PaymentRequestInfoExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of PaymentRequestInfoExt
func (u PaymentRequestInfoExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewPaymentRequestInfoExt creates a new  PaymentRequestInfoExt.
func NewPaymentRequestInfoExt(v LedgerVersion, value interface{}) (result PaymentRequestInfoExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// PaymentRequestInfo is an XDR Struct defines as:
//
//   struct PaymentRequestInfo {
//        PaymentRequestEntry paymentRequest;
//        AccountID source;
//        AccountID destination;
//    	 // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type PaymentRequestInfo struct {
	PaymentRequest PaymentRequestEntry   `json:"paymentRequest,omitempty"`
	Source         AccountId             `json:"source,omitempty"`
	Destination    AccountId             `json:"destination,omitempty"`
	Ext            PaymentRequestInfoExt `json:"ext,omitempty"`
}

// DemurrageInfoExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type DemurrageInfoExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u DemurrageInfoExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of DemurrageInfoExt
func (u DemurrageInfoExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewDemurrageInfoExt creates a new  DemurrageInfoExt.
func NewDemurrageInfoExt(v LedgerVersion, value interface{}) (result DemurrageInfoExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// DemurrageInfo is an XDR Struct defines as:
//
//   struct DemurrageInfo {
//        PaymentRequestInfo paymentRequests<>;
//        BalanceDemurrage demurrages<>;
//    	 // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type DemurrageInfo struct {
	PaymentRequests []PaymentRequestInfo `json:"paymentRequests,omitempty"`
	Demurrages      []BalanceDemurrage   `json:"demurrages,omitempty"`
	Ext             DemurrageInfoExt     `json:"ext,omitempty"`
}

// DemurrageResult is an XDR Union defines as:
//
//   union DemurrageResult switch (DemurrageResultCode code)
//    {
//    case DEMURRAGE_SUCCESS:
//        DemurrageInfo demurrageInfo;
//    default:
//        void;
//    };
//
type DemurrageResult struct {
	Code          DemurrageResultCode `json:"code,omitempty"`
	DemurrageInfo *DemurrageInfo      `json:"demurrageInfo,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u DemurrageResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of DemurrageResult
func (u DemurrageResult) ArmForSwitch(sw int32) (string, bool) {
	switch DemurrageResultCode(sw) {
	case DemurrageResultCodeDemurrageSuccess:
		return "DemurrageInfo", true
	default:
		return "", true
	}
}

// NewDemurrageResult creates a new  DemurrageResult.
func NewDemurrageResult(code DemurrageResultCode, value interface{}) (result DemurrageResult, err error) {
	result.Code = code
	switch DemurrageResultCode(code) {
	case DemurrageResultCodeDemurrageSuccess:
		tv, ok := value.(DemurrageInfo)
		if !ok {
			err = fmt.Errorf("invalid value, must be DemurrageInfo")
			return
		}
		result.DemurrageInfo = &tv
	default:
		// void
	}
	return
}

// MustDemurrageInfo retrieves the DemurrageInfo value from the union,
// panicing if the value is not set.
func (u DemurrageResult) MustDemurrageInfo() DemurrageInfo {
	val, ok := u.GetDemurrageInfo()

	if !ok {
		panic("arm DemurrageInfo is not set")
	}

	return val
}

// GetDemurrageInfo retrieves the DemurrageInfo value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u DemurrageResult) GetDemurrageInfo() (result DemurrageInfo, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "DemurrageInfo" {
		result = *u.DemurrageInfo
		ok = true
	}

	return
}

// DirectDebitOpExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type DirectDebitOpExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u DirectDebitOpExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of DirectDebitOpExt
func (u DirectDebitOpExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewDirectDebitOpExt creates a new  DirectDebitOpExt.
func NewDirectDebitOpExt(v LedgerVersion, value interface{}) (result DirectDebitOpExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// DirectDebitOp is an XDR Struct defines as:
//
//   struct DirectDebitOp
//    {
//        AccountID from;
//        PaymentOp paymentOp;
//    	 // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type DirectDebitOp struct {
	From      AccountId        `json:"from,omitempty"`
	PaymentOp PaymentOp        `json:"paymentOp,omitempty"`
	Ext       DirectDebitOpExt `json:"ext,omitempty"`
}

// DirectDebitResultCode is an XDR Enum defines as:
//
//   enum DirectDebitResultCode
//    {
//        // codes considered as "success" for the operation
//        DIRECT_DEBIT_SUCCESS = 0, // payment successfuly completed
//
//        // codes considered as "failure" for the operation
//        DIRECT_DEBIT_MALFORMED = -1,       // bad input
//        DIRECT_DEBIT_UNDERFUNDED = -2,     // not enough funds in source account
//        DIRECT_DEBIT_LINE_FULL = -3,       // destination would go above their limit
//    	DIRECT_DEBIT_FEE_MISMATCHED = -4,   // fee is not equal to expected fee
//        DIRECT_DEBIT_BALANCE_NOT_FOUND = -5, // destination balance not found
//        DIRECT_DEBIT_BALANCE_ACCOUNT_MISMATCHED = -6,
//        DIRECT_DEBIT_BALANCE_ASSETS_MISMATCHED = -7,
//    	DIRECT_DEBIT_SRC_BALANCE_NOT_FOUND = -8, // source balance not found
//        DIRECT_DEBIT_REFERENCE_DUPLICATION = -9,
//        DIRECT_DEBIT_STATS_OVERFLOW = -10,
//        DIRECT_DEBIT_LIMITS_EXCEEDED = -11,
//        DIRECT_DEBIT_NOT_ALLOWED_BY_ASSET_POLICY = -12,
//        DIRECT_DEBIT_NO_TRUST = -13
//    };
//
type DirectDebitResultCode int32

const (
	DirectDebitResultCodeDirectDebitSuccess                  DirectDebitResultCode = 0
	DirectDebitResultCodeDirectDebitMalformed                DirectDebitResultCode = -1
	DirectDebitResultCodeDirectDebitUnderfunded              DirectDebitResultCode = -2
	DirectDebitResultCodeDirectDebitLineFull                 DirectDebitResultCode = -3
	DirectDebitResultCodeDirectDebitFeeMismatched            DirectDebitResultCode = -4
	DirectDebitResultCodeDirectDebitBalanceNotFound          DirectDebitResultCode = -5
	DirectDebitResultCodeDirectDebitBalanceAccountMismatched DirectDebitResultCode = -6
	DirectDebitResultCodeDirectDebitBalanceAssetsMismatched  DirectDebitResultCode = -7
	DirectDebitResultCodeDirectDebitSrcBalanceNotFound       DirectDebitResultCode = -8
	DirectDebitResultCodeDirectDebitReferenceDuplication     DirectDebitResultCode = -9
	DirectDebitResultCodeDirectDebitStatsOverflow            DirectDebitResultCode = -10
	DirectDebitResultCodeDirectDebitLimitsExceeded           DirectDebitResultCode = -11
	DirectDebitResultCodeDirectDebitNotAllowedByAssetPolicy  DirectDebitResultCode = -12
	DirectDebitResultCodeDirectDebitNoTrust                  DirectDebitResultCode = -13
)

var DirectDebitResultCodeAll = []DirectDebitResultCode{
	DirectDebitResultCodeDirectDebitSuccess,
	DirectDebitResultCodeDirectDebitMalformed,
	DirectDebitResultCodeDirectDebitUnderfunded,
	DirectDebitResultCodeDirectDebitLineFull,
	DirectDebitResultCodeDirectDebitFeeMismatched,
	DirectDebitResultCodeDirectDebitBalanceNotFound,
	DirectDebitResultCodeDirectDebitBalanceAccountMismatched,
	DirectDebitResultCodeDirectDebitBalanceAssetsMismatched,
	DirectDebitResultCodeDirectDebitSrcBalanceNotFound,
	DirectDebitResultCodeDirectDebitReferenceDuplication,
	DirectDebitResultCodeDirectDebitStatsOverflow,
	DirectDebitResultCodeDirectDebitLimitsExceeded,
	DirectDebitResultCodeDirectDebitNotAllowedByAssetPolicy,
	DirectDebitResultCodeDirectDebitNoTrust,
}

var directDebitResultCodeMap = map[int32]string{
	0:   "DirectDebitResultCodeDirectDebitSuccess",
	-1:  "DirectDebitResultCodeDirectDebitMalformed",
	-2:  "DirectDebitResultCodeDirectDebitUnderfunded",
	-3:  "DirectDebitResultCodeDirectDebitLineFull",
	-4:  "DirectDebitResultCodeDirectDebitFeeMismatched",
	-5:  "DirectDebitResultCodeDirectDebitBalanceNotFound",
	-6:  "DirectDebitResultCodeDirectDebitBalanceAccountMismatched",
	-7:  "DirectDebitResultCodeDirectDebitBalanceAssetsMismatched",
	-8:  "DirectDebitResultCodeDirectDebitSrcBalanceNotFound",
	-9:  "DirectDebitResultCodeDirectDebitReferenceDuplication",
	-10: "DirectDebitResultCodeDirectDebitStatsOverflow",
	-11: "DirectDebitResultCodeDirectDebitLimitsExceeded",
	-12: "DirectDebitResultCodeDirectDebitNotAllowedByAssetPolicy",
	-13: "DirectDebitResultCodeDirectDebitNoTrust",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for DirectDebitResultCode
func (e DirectDebitResultCode) ValidEnum(v int32) bool {
	_, ok := directDebitResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e DirectDebitResultCode) String() string {
	name, _ := directDebitResultCodeMap[int32(e)]
	return name
}

func (e DirectDebitResultCode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// DirectDebitSuccessExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type DirectDebitSuccessExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u DirectDebitSuccessExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of DirectDebitSuccessExt
func (u DirectDebitSuccessExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewDirectDebitSuccessExt creates a new  DirectDebitSuccessExt.
func NewDirectDebitSuccessExt(v LedgerVersion, value interface{}) (result DirectDebitSuccessExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// DirectDebitSuccess is an XDR Struct defines as:
//
//   struct DirectDebitSuccess {
//    	PaymentResponse paymentResponse;
//    	// reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type DirectDebitSuccess struct {
	PaymentResponse PaymentResponse       `json:"paymentResponse,omitempty"`
	Ext             DirectDebitSuccessExt `json:"ext,omitempty"`
}

// DirectDebitResult is an XDR Union defines as:
//
//   union DirectDebitResult switch (DirectDebitResultCode code)
//    {
//    case DIRECT_DEBIT_SUCCESS:
//        DirectDebitSuccess success;
//    default:
//        void;
//    };
//
type DirectDebitResult struct {
	Code    DirectDebitResultCode `json:"code,omitempty"`
	Success *DirectDebitSuccess   `json:"success,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u DirectDebitResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of DirectDebitResult
func (u DirectDebitResult) ArmForSwitch(sw int32) (string, bool) {
	switch DirectDebitResultCode(sw) {
	case DirectDebitResultCodeDirectDebitSuccess:
		return "Success", true
	default:
		return "", true
	}
}

// NewDirectDebitResult creates a new  DirectDebitResult.
func NewDirectDebitResult(code DirectDebitResultCode, value interface{}) (result DirectDebitResult, err error) {
	result.Code = code
	switch DirectDebitResultCode(code) {
	case DirectDebitResultCodeDirectDebitSuccess:
		tv, ok := value.(DirectDebitSuccess)
		if !ok {
			err = fmt.Errorf("invalid value, must be DirectDebitSuccess")
			return
		}
		result.Success = &tv
	default:
		// void
	}
	return
}

// MustSuccess retrieves the Success value from the union,
// panicing if the value is not set.
func (u DirectDebitResult) MustSuccess() DirectDebitSuccess {
	val, ok := u.GetSuccess()

	if !ok {
		panic("arm Success is not set")
	}

	return val
}

// GetSuccess retrieves the Success value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u DirectDebitResult) GetSuccess() (result DirectDebitSuccess, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "Success" {
		result = *u.Success
		ok = true
	}

	return
}

// ForfeitOpExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type ForfeitOpExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ForfeitOpExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ForfeitOpExt
func (u ForfeitOpExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewForfeitOpExt creates a new  ForfeitOpExt.
func NewForfeitOpExt(v LedgerVersion, value interface{}) (result ForfeitOpExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ForfeitOp is an XDR Struct defines as:
//
//   struct ForfeitOp
//    {
//        BalanceID balance; // balance to withdraw from
//        int64 amount; // amount they end up with
//        RequestType type;
//    	 // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type ForfeitOp struct {
	Balance BalanceId    `json:"balance,omitempty"`
	Amount  Int64        `json:"amount,omitempty"`
	Type    RequestType  `json:"type,omitempty"`
	Ext     ForfeitOpExt `json:"ext,omitempty"`
}

// ForfeitResultCode is an XDR Enum defines as:
//
//   enum ForfeitResultCode
//    {
//        // codes considered as "success" for the operation
//        FORFEIT_SUCCESS = 0, // money was withdrawn
//
//        // codes considered as "failure" for the operation
//        FORFEIT_MALFORMED = -1,       // invalid amount and unknown errors
//        FORFEIT_UNDERFUNDED = -2,     // not enough funds in related account
//        FORFEIT_BALANCE_NOT_FOUND = -3,     // account not found
//        FORFEIT_STATS_OVERFLOW = -4,
//        FORFEIT_LIMITS_EXCEEDED = -5
//    };
//
type ForfeitResultCode int32

const (
	ForfeitResultCodeForfeitSuccess         ForfeitResultCode = 0
	ForfeitResultCodeForfeitMalformed       ForfeitResultCode = -1
	ForfeitResultCodeForfeitUnderfunded     ForfeitResultCode = -2
	ForfeitResultCodeForfeitBalanceNotFound ForfeitResultCode = -3
	ForfeitResultCodeForfeitStatsOverflow   ForfeitResultCode = -4
	ForfeitResultCodeForfeitLimitsExceeded  ForfeitResultCode = -5
)

var ForfeitResultCodeAll = []ForfeitResultCode{
	ForfeitResultCodeForfeitSuccess,
	ForfeitResultCodeForfeitMalformed,
	ForfeitResultCodeForfeitUnderfunded,
	ForfeitResultCodeForfeitBalanceNotFound,
	ForfeitResultCodeForfeitStatsOverflow,
	ForfeitResultCodeForfeitLimitsExceeded,
}

var forfeitResultCodeMap = map[int32]string{
	0:  "ForfeitResultCodeForfeitSuccess",
	-1: "ForfeitResultCodeForfeitMalformed",
	-2: "ForfeitResultCodeForfeitUnderfunded",
	-3: "ForfeitResultCodeForfeitBalanceNotFound",
	-4: "ForfeitResultCodeForfeitStatsOverflow",
	-5: "ForfeitResultCodeForfeitLimitsExceeded",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ForfeitResultCode
func (e ForfeitResultCode) ValidEnum(v int32) bool {
	_, ok := forfeitResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e ForfeitResultCode) String() string {
	name, _ := forfeitResultCodeMap[int32(e)]
	return name
}

func (e ForfeitResultCode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// ForfeitSuccessExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type ForfeitSuccessExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ForfeitSuccessExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ForfeitSuccessExt
func (u ForfeitSuccessExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewForfeitSuccessExt creates a new  ForfeitSuccessExt.
func NewForfeitSuccessExt(v LedgerVersion, value interface{}) (result ForfeitSuccessExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ForfeitSuccess is an XDR Struct defines as:
//
//   struct ForfeitSuccess {
//    	uint64 paymentID;
//    	// reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type ForfeitSuccess struct {
	PaymentId Uint64            `json:"paymentID,omitempty"`
	Ext       ForfeitSuccessExt `json:"ext,omitempty"`
}

// ForfeitResult is an XDR Union defines as:
//
//   union ForfeitResult switch (ForfeitResultCode code)
//    {
//    case FORFEIT_SUCCESS:
//        ForfeitSuccess success;
//    default:
//        void;
//    };
//
type ForfeitResult struct {
	Code    ForfeitResultCode `json:"code,omitempty"`
	Success *ForfeitSuccess   `json:"success,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ForfeitResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ForfeitResult
func (u ForfeitResult) ArmForSwitch(sw int32) (string, bool) {
	switch ForfeitResultCode(sw) {
	case ForfeitResultCodeForfeitSuccess:
		return "Success", true
	default:
		return "", true
	}
}

// NewForfeitResult creates a new  ForfeitResult.
func NewForfeitResult(code ForfeitResultCode, value interface{}) (result ForfeitResult, err error) {
	result.Code = code
	switch ForfeitResultCode(code) {
	case ForfeitResultCodeForfeitSuccess:
		tv, ok := value.(ForfeitSuccess)
		if !ok {
			err = fmt.Errorf("invalid value, must be ForfeitSuccess")
			return
		}
		result.Success = &tv
	default:
		// void
	}
	return
}

// MustSuccess retrieves the Success value from the union,
// panicing if the value is not set.
func (u ForfeitResult) MustSuccess() ForfeitSuccess {
	val, ok := u.GetSuccess()

	if !ok {
		panic("arm Success is not set")
	}

	return val
}

// GetSuccess retrieves the Success value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ForfeitResult) GetSuccess() (result ForfeitSuccess, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "Success" {
		result = *u.Success
		ok = true
	}

	return
}

// ManageAccountOpExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type ManageAccountOpExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageAccountOpExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageAccountOpExt
func (u ManageAccountOpExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewManageAccountOpExt creates a new  ManageAccountOpExt.
func NewManageAccountOpExt(v LedgerVersion, value interface{}) (result ManageAccountOpExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ManageAccountOp is an XDR Struct defines as:
//
//   struct ManageAccountOp
//    {
//        AccountID account; // account to manage
//        AccountType accountType;
//        uint32 blockReasonsToAdd;
//        uint32 blockReasonsToRemove;
//    	 // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type ManageAccountOp struct {
	Account              AccountId          `json:"account,omitempty"`
	AccountType          AccountType        `json:"accountType,omitempty"`
	BlockReasonsToAdd    Uint32             `json:"blockReasonsToAdd,omitempty"`
	BlockReasonsToRemove Uint32             `json:"blockReasonsToRemove,omitempty"`
	Ext                  ManageAccountOpExt `json:"ext,omitempty"`
}

// ManageAccountResultCode is an XDR Enum defines as:
//
//   enum ManageAccountResultCode
//    {
//        // codes considered as "success" for the operation
//        MANAGE_ACCOUNT_SUCCESS = 0, // account was created
//
//        // codes considered as "failure" for the operation
//        MANAGE_ACCOUNT_NOT_FOUND = -1,         // account does not exists
//        MANAGE_ACCOUNT_MALFORMED = -2,
//    	MANAGE_ACCOUNT_NOT_ALLOWED = -3,         // manage account operation is not allowed on this account
//        MANAGE_ACCOUNT_TYPE_MISMATCH = -4
//    };
//
type ManageAccountResultCode int32

const (
	ManageAccountResultCodeManageAccountSuccess      ManageAccountResultCode = 0
	ManageAccountResultCodeManageAccountNotFound     ManageAccountResultCode = -1
	ManageAccountResultCodeManageAccountMalformed    ManageAccountResultCode = -2
	ManageAccountResultCodeManageAccountNotAllowed   ManageAccountResultCode = -3
	ManageAccountResultCodeManageAccountTypeMismatch ManageAccountResultCode = -4
)

var ManageAccountResultCodeAll = []ManageAccountResultCode{
	ManageAccountResultCodeManageAccountSuccess,
	ManageAccountResultCodeManageAccountNotFound,
	ManageAccountResultCodeManageAccountMalformed,
	ManageAccountResultCodeManageAccountNotAllowed,
	ManageAccountResultCodeManageAccountTypeMismatch,
}

var manageAccountResultCodeMap = map[int32]string{
	0:  "ManageAccountResultCodeManageAccountSuccess",
	-1: "ManageAccountResultCodeManageAccountNotFound",
	-2: "ManageAccountResultCodeManageAccountMalformed",
	-3: "ManageAccountResultCodeManageAccountNotAllowed",
	-4: "ManageAccountResultCodeManageAccountTypeMismatch",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ManageAccountResultCode
func (e ManageAccountResultCode) ValidEnum(v int32) bool {
	_, ok := manageAccountResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e ManageAccountResultCode) String() string {
	name, _ := manageAccountResultCodeMap[int32(e)]
	return name
}

func (e ManageAccountResultCode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// ManageAccountSuccessExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type ManageAccountSuccessExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageAccountSuccessExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageAccountSuccessExt
func (u ManageAccountSuccessExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewManageAccountSuccessExt creates a new  ManageAccountSuccessExt.
func NewManageAccountSuccessExt(v LedgerVersion, value interface{}) (result ManageAccountSuccessExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ManageAccountSuccess is an XDR Struct defines as:
//
//   struct ManageAccountSuccess {
//    	uint32 blockReasons;
//     // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type ManageAccountSuccess struct {
	BlockReasons Uint32                  `json:"blockReasons,omitempty"`
	Ext          ManageAccountSuccessExt `json:"ext,omitempty"`
}

// ManageAccountResult is an XDR Union defines as:
//
//   union ManageAccountResult switch (ManageAccountResultCode code)
//    {
//    case MANAGE_ACCOUNT_SUCCESS:
//        ManageAccountSuccess success;
//    default:
//        void;
//    };
//
type ManageAccountResult struct {
	Code    ManageAccountResultCode `json:"code,omitempty"`
	Success *ManageAccountSuccess   `json:"success,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageAccountResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageAccountResult
func (u ManageAccountResult) ArmForSwitch(sw int32) (string, bool) {
	switch ManageAccountResultCode(sw) {
	case ManageAccountResultCodeManageAccountSuccess:
		return "Success", true
	default:
		return "", true
	}
}

// NewManageAccountResult creates a new  ManageAccountResult.
func NewManageAccountResult(code ManageAccountResultCode, value interface{}) (result ManageAccountResult, err error) {
	result.Code = code
	switch ManageAccountResultCode(code) {
	case ManageAccountResultCodeManageAccountSuccess:
		tv, ok := value.(ManageAccountSuccess)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageAccountSuccess")
			return
		}
		result.Success = &tv
	default:
		// void
	}
	return
}

// MustSuccess retrieves the Success value from the union,
// panicing if the value is not set.
func (u ManageAccountResult) MustSuccess() ManageAccountSuccess {
	val, ok := u.GetSuccess()

	if !ok {
		panic("arm Success is not set")
	}

	return val
}

// GetSuccess retrieves the Success value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ManageAccountResult) GetSuccess() (result ManageAccountSuccess, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "Success" {
		result = *u.Success
		ok = true
	}

	return
}

// ManageAssetPairAction is an XDR Enum defines as:
//
//   enum ManageAssetPairAction
//    {
//        MANAGE_ASSET_PAIR_CREATE = 0,
//        MANAGE_ASSET_PAIR_UPDATE_PRICE = 1,
//        MANAGE_ASSET_PAIR_UPDATE_POLICIES = 2
//    };
//
type ManageAssetPairAction int32

const (
	ManageAssetPairActionManageAssetPairCreate         ManageAssetPairAction = 0
	ManageAssetPairActionManageAssetPairUpdatePrice    ManageAssetPairAction = 1
	ManageAssetPairActionManageAssetPairUpdatePolicies ManageAssetPairAction = 2
)

var ManageAssetPairActionAll = []ManageAssetPairAction{
	ManageAssetPairActionManageAssetPairCreate,
	ManageAssetPairActionManageAssetPairUpdatePrice,
	ManageAssetPairActionManageAssetPairUpdatePolicies,
}

var manageAssetPairActionMap = map[int32]string{
	0: "ManageAssetPairActionManageAssetPairCreate",
	1: "ManageAssetPairActionManageAssetPairUpdatePrice",
	2: "ManageAssetPairActionManageAssetPairUpdatePolicies",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ManageAssetPairAction
func (e ManageAssetPairAction) ValidEnum(v int32) bool {
	_, ok := manageAssetPairActionMap[v]
	return ok
}

// String returns the name of `e`
func (e ManageAssetPairAction) String() string {
	name, _ := manageAssetPairActionMap[int32(e)]
	return name
}

func (e ManageAssetPairAction) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// ManageAssetPairOpExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type ManageAssetPairOpExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageAssetPairOpExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageAssetPairOpExt
func (u ManageAssetPairOpExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewManageAssetPairOpExt creates a new  ManageAssetPairOpExt.
func NewManageAssetPairOpExt(v LedgerVersion, value interface{}) (result ManageAssetPairOpExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ManageAssetPairOp is an XDR Struct defines as:
//
//   struct ManageAssetPairOp
//    {
//        ManageAssetPairAction action;
//    	AssetCode base;
//    	AssetCode quote;
//
//        int64 physicalPrice;
//
//    	int64 physicalPriceCorrection; // correction of physical price in percents. If physical price is set and restriction by physical price set, mininal price for offer for this pair will be physicalPrice * physicalPriceCorrection
//    	int64 maxPriceStep;
//
//    	int32 policies;
//
//    	 // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type ManageAssetPairOp struct {
	Action                  ManageAssetPairAction `json:"action,omitempty"`
	Base                    AssetCode             `json:"base,omitempty"`
	Quote                   AssetCode             `json:"quote,omitempty"`
	PhysicalPrice           Int64                 `json:"physicalPrice,omitempty"`
	PhysicalPriceCorrection Int64                 `json:"physicalPriceCorrection,omitempty"`
	MaxPriceStep            Int64                 `json:"maxPriceStep,omitempty"`
	Policies                Int32                 `json:"policies,omitempty"`
	Ext                     ManageAssetPairOpExt  `json:"ext,omitempty"`
}

// ManageAssetPairResultCode is an XDR Enum defines as:
//
//   enum ManageAssetPairResultCode
//    {
//        // codes considered as "success" for the operation
//        MANAGE_ASSET_PAIR_SUCCESS = 0,
//
//        // codes considered as "failure" for the operation
//    	MANAGE_ASSET_PAIR_NOT_FOUND = -1,           // failed to find asset with such code
//    	MANAGE_ASSET_PAIR_ALREADY_EXISTS = -2,
//        MANAGE_ASSET_PAIR_MALFORMED = -3,
//    	MANAGE_ASSET_PAIR_INVALID_ASSET = -4,
//    	MANAGE_ASSET_PAIR_INVALID_ACTION = -5,
//    	MANAGE_ASSET_PAIR_INVALID_POLICIES = -6,
//    	MANAGE_ASSET_PAIR_ASSET_NOT_FOUND = -7
//    };
//
type ManageAssetPairResultCode int32

const (
	ManageAssetPairResultCodeManageAssetPairSuccess         ManageAssetPairResultCode = 0
	ManageAssetPairResultCodeManageAssetPairNotFound        ManageAssetPairResultCode = -1
	ManageAssetPairResultCodeManageAssetPairAlreadyExists   ManageAssetPairResultCode = -2
	ManageAssetPairResultCodeManageAssetPairMalformed       ManageAssetPairResultCode = -3
	ManageAssetPairResultCodeManageAssetPairInvalidAsset    ManageAssetPairResultCode = -4
	ManageAssetPairResultCodeManageAssetPairInvalidAction   ManageAssetPairResultCode = -5
	ManageAssetPairResultCodeManageAssetPairInvalidPolicies ManageAssetPairResultCode = -6
	ManageAssetPairResultCodeManageAssetPairAssetNotFound   ManageAssetPairResultCode = -7
)

var ManageAssetPairResultCodeAll = []ManageAssetPairResultCode{
	ManageAssetPairResultCodeManageAssetPairSuccess,
	ManageAssetPairResultCodeManageAssetPairNotFound,
	ManageAssetPairResultCodeManageAssetPairAlreadyExists,
	ManageAssetPairResultCodeManageAssetPairMalformed,
	ManageAssetPairResultCodeManageAssetPairInvalidAsset,
	ManageAssetPairResultCodeManageAssetPairInvalidAction,
	ManageAssetPairResultCodeManageAssetPairInvalidPolicies,
	ManageAssetPairResultCodeManageAssetPairAssetNotFound,
}

var manageAssetPairResultCodeMap = map[int32]string{
	0:  "ManageAssetPairResultCodeManageAssetPairSuccess",
	-1: "ManageAssetPairResultCodeManageAssetPairNotFound",
	-2: "ManageAssetPairResultCodeManageAssetPairAlreadyExists",
	-3: "ManageAssetPairResultCodeManageAssetPairMalformed",
	-4: "ManageAssetPairResultCodeManageAssetPairInvalidAsset",
	-5: "ManageAssetPairResultCodeManageAssetPairInvalidAction",
	-6: "ManageAssetPairResultCodeManageAssetPairInvalidPolicies",
	-7: "ManageAssetPairResultCodeManageAssetPairAssetNotFound",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ManageAssetPairResultCode
func (e ManageAssetPairResultCode) ValidEnum(v int32) bool {
	_, ok := manageAssetPairResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e ManageAssetPairResultCode) String() string {
	name, _ := manageAssetPairResultCodeMap[int32(e)]
	return name
}

func (e ManageAssetPairResultCode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// ManageAssetPairSuccessExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type ManageAssetPairSuccessExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageAssetPairSuccessExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageAssetPairSuccessExt
func (u ManageAssetPairSuccessExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewManageAssetPairSuccessExt creates a new  ManageAssetPairSuccessExt.
func NewManageAssetPairSuccessExt(v LedgerVersion, value interface{}) (result ManageAssetPairSuccessExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ManageAssetPairSuccess is an XDR Struct defines as:
//
//   struct ManageAssetPairSuccess
//    {
//    	int64 currentPrice;
//    	// reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type ManageAssetPairSuccess struct {
	CurrentPrice Int64                     `json:"currentPrice,omitempty"`
	Ext          ManageAssetPairSuccessExt `json:"ext,omitempty"`
}

// ManageAssetPairResult is an XDR Union defines as:
//
//   union ManageAssetPairResult switch (ManageAssetPairResultCode code)
//    {
//    case MANAGE_ASSET_PAIR_SUCCESS:
//        ManageAssetPairSuccess success;
//    default:
//        void;
//    };
//
type ManageAssetPairResult struct {
	Code    ManageAssetPairResultCode `json:"code,omitempty"`
	Success *ManageAssetPairSuccess   `json:"success,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageAssetPairResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageAssetPairResult
func (u ManageAssetPairResult) ArmForSwitch(sw int32) (string, bool) {
	switch ManageAssetPairResultCode(sw) {
	case ManageAssetPairResultCodeManageAssetPairSuccess:
		return "Success", true
	default:
		return "", true
	}
}

// NewManageAssetPairResult creates a new  ManageAssetPairResult.
func NewManageAssetPairResult(code ManageAssetPairResultCode, value interface{}) (result ManageAssetPairResult, err error) {
	result.Code = code
	switch ManageAssetPairResultCode(code) {
	case ManageAssetPairResultCodeManageAssetPairSuccess:
		tv, ok := value.(ManageAssetPairSuccess)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageAssetPairSuccess")
			return
		}
		result.Success = &tv
	default:
		// void
	}
	return
}

// MustSuccess retrieves the Success value from the union,
// panicing if the value is not set.
func (u ManageAssetPairResult) MustSuccess() ManageAssetPairSuccess {
	val, ok := u.GetSuccess()

	if !ok {
		panic("arm Success is not set")
	}

	return val
}

// GetSuccess retrieves the Success value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ManageAssetPairResult) GetSuccess() (result ManageAssetPairSuccess, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "Success" {
		result = *u.Success
		ok = true
	}

	return
}

// ManageAssetAction is an XDR Enum defines as:
//
//   enum ManageAssetAction
//    {
//        MANAGE_ASSET_CREATE = 0,
//        MANAGE_ASSET_UPDATE_POLICIES = 1,
//    	MANAGE_ASSET_ADD_TOKEN = 2
//    };
//
type ManageAssetAction int32

const (
	ManageAssetActionManageAssetCreate         ManageAssetAction = 0
	ManageAssetActionManageAssetUpdatePolicies ManageAssetAction = 1
	ManageAssetActionManageAssetAddToken       ManageAssetAction = 2
)

var ManageAssetActionAll = []ManageAssetAction{
	ManageAssetActionManageAssetCreate,
	ManageAssetActionManageAssetUpdatePolicies,
	ManageAssetActionManageAssetAddToken,
}

var manageAssetActionMap = map[int32]string{
	0: "ManageAssetActionManageAssetCreate",
	1: "ManageAssetActionManageAssetUpdatePolicies",
	2: "ManageAssetActionManageAssetAddToken",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ManageAssetAction
func (e ManageAssetAction) ValidEnum(v int32) bool {
	_, ok := manageAssetActionMap[v]
	return ok
}

// String returns the name of `e`
func (e ManageAssetAction) String() string {
	name, _ := manageAssetActionMap[int32(e)]
	return name
}

func (e ManageAssetAction) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// ManageAssetOpExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type ManageAssetOpExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageAssetOpExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageAssetOpExt
func (u ManageAssetOpExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewManageAssetOpExt creates a new  ManageAssetOpExt.
func NewManageAssetOpExt(v LedgerVersion, value interface{}) (result ManageAssetOpExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ManageAssetOp is an XDR Struct defines as:
//
//   struct ManageAssetOp
//    {
//        ManageAssetAction action;
//    	AssetCode code;
//
//        AssetCode* token; // if set it is a token of some asset
//
//    	AssetForm assetForms<>;
//        int32 policies;
//
//    	 // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type ManageAssetOp struct {
	Action     ManageAssetAction `json:"action,omitempty"`
	Code       AssetCode         `json:"code,omitempty"`
	Token      *AssetCode        `json:"token,omitempty"`
	AssetForms []AssetForm       `json:"assetForms,omitempty"`
	Policies   Int32             `json:"policies,omitempty"`
	Ext        ManageAssetOpExt  `json:"ext,omitempty"`
}

// ManageAssetResultCode is an XDR Enum defines as:
//
//   enum ManageAssetResultCode
//    {
//        // codes considered as "success" for the operation
//        MANAGE_ASSET_SUCCESS = 0,
//
//        // codes considered as "failure" for the operation
//    	MANAGE_ASSET_NOT_FOUND = -1,           // failed to find asset with such code
//    	MANAGE_ASSET_ALREADY_EXISTS = -2,
//        MANAGE_ASSET_MALFORMED = -3,
//        MANAGE_ASSET_TOKEN_ALREADY_EXISTS = -4,
//    	MANAGE_ASSET_TOKEN_TOKEN_ALREDY_SET = -5,
//    	MANAGE_ASSET_IS_NOT_TOKEN = -6
//    };
//
type ManageAssetResultCode int32

const (
	ManageAssetResultCodeManageAssetSuccess             ManageAssetResultCode = 0
	ManageAssetResultCodeManageAssetNotFound            ManageAssetResultCode = -1
	ManageAssetResultCodeManageAssetAlreadyExists       ManageAssetResultCode = -2
	ManageAssetResultCodeManageAssetMalformed           ManageAssetResultCode = -3
	ManageAssetResultCodeManageAssetTokenAlreadyExists  ManageAssetResultCode = -4
	ManageAssetResultCodeManageAssetTokenTokenAlredySet ManageAssetResultCode = -5
	ManageAssetResultCodeManageAssetIsNotToken          ManageAssetResultCode = -6
)

var ManageAssetResultCodeAll = []ManageAssetResultCode{
	ManageAssetResultCodeManageAssetSuccess,
	ManageAssetResultCodeManageAssetNotFound,
	ManageAssetResultCodeManageAssetAlreadyExists,
	ManageAssetResultCodeManageAssetMalformed,
	ManageAssetResultCodeManageAssetTokenAlreadyExists,
	ManageAssetResultCodeManageAssetTokenTokenAlredySet,
	ManageAssetResultCodeManageAssetIsNotToken,
}

var manageAssetResultCodeMap = map[int32]string{
	0:  "ManageAssetResultCodeManageAssetSuccess",
	-1: "ManageAssetResultCodeManageAssetNotFound",
	-2: "ManageAssetResultCodeManageAssetAlreadyExists",
	-3: "ManageAssetResultCodeManageAssetMalformed",
	-4: "ManageAssetResultCodeManageAssetTokenAlreadyExists",
	-5: "ManageAssetResultCodeManageAssetTokenTokenAlredySet",
	-6: "ManageAssetResultCodeManageAssetIsNotToken",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ManageAssetResultCode
func (e ManageAssetResultCode) ValidEnum(v int32) bool {
	_, ok := manageAssetResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e ManageAssetResultCode) String() string {
	name, _ := manageAssetResultCodeMap[int32(e)]
	return name
}

func (e ManageAssetResultCode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// ManageAssetSuccessExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type ManageAssetSuccessExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageAssetSuccessExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageAssetSuccessExt
func (u ManageAssetSuccessExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewManageAssetSuccessExt creates a new  ManageAssetSuccessExt.
func NewManageAssetSuccessExt(v LedgerVersion, value interface{}) (result ManageAssetSuccessExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ManageAssetSuccess is an XDR Struct defines as:
//
//   struct ManageAssetSuccess {
//     // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type ManageAssetSuccess struct {
	Ext ManageAssetSuccessExt `json:"ext,omitempty"`
}

// ManageAssetResult is an XDR Union defines as:
//
//   union ManageAssetResult switch (ManageAssetResultCode code)
//    {
//    case MANAGE_ASSET_SUCCESS:
//        ManageAssetSuccess success;
//    default:
//        void;
//    };
//
type ManageAssetResult struct {
	Code    ManageAssetResultCode `json:"code,omitempty"`
	Success *ManageAssetSuccess   `json:"success,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageAssetResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageAssetResult
func (u ManageAssetResult) ArmForSwitch(sw int32) (string, bool) {
	switch ManageAssetResultCode(sw) {
	case ManageAssetResultCodeManageAssetSuccess:
		return "Success", true
	default:
		return "", true
	}
}

// NewManageAssetResult creates a new  ManageAssetResult.
func NewManageAssetResult(code ManageAssetResultCode, value interface{}) (result ManageAssetResult, err error) {
	result.Code = code
	switch ManageAssetResultCode(code) {
	case ManageAssetResultCodeManageAssetSuccess:
		tv, ok := value.(ManageAssetSuccess)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageAssetSuccess")
			return
		}
		result.Success = &tv
	default:
		// void
	}
	return
}

// MustSuccess retrieves the Success value from the union,
// panicing if the value is not set.
func (u ManageAssetResult) MustSuccess() ManageAssetSuccess {
	val, ok := u.GetSuccess()

	if !ok {
		panic("arm Success is not set")
	}

	return val
}

// GetSuccess retrieves the Success value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ManageAssetResult) GetSuccess() (result ManageAssetSuccess, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "Success" {
		result = *u.Success
		ok = true
	}

	return
}

// ManageBalanceAction is an XDR Enum defines as:
//
//   enum ManageBalanceAction
//    {
//        MANAGE_BALANCE_CREATE = 0,
//        MANAGE_BALANCE_DELETE = 1
//    };
//
type ManageBalanceAction int32

const (
	ManageBalanceActionManageBalanceCreate ManageBalanceAction = 0
	ManageBalanceActionManageBalanceDelete ManageBalanceAction = 1
)

var ManageBalanceActionAll = []ManageBalanceAction{
	ManageBalanceActionManageBalanceCreate,
	ManageBalanceActionManageBalanceDelete,
}

var manageBalanceActionMap = map[int32]string{
	0: "ManageBalanceActionManageBalanceCreate",
	1: "ManageBalanceActionManageBalanceDelete",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ManageBalanceAction
func (e ManageBalanceAction) ValidEnum(v int32) bool {
	_, ok := manageBalanceActionMap[v]
	return ok
}

// String returns the name of `e`
func (e ManageBalanceAction) String() string {
	name, _ := manageBalanceActionMap[int32(e)]
	return name
}

func (e ManageBalanceAction) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// ManageBalanceOpExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type ManageBalanceOpExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageBalanceOpExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageBalanceOpExt
func (u ManageBalanceOpExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewManageBalanceOpExt creates a new  ManageBalanceOpExt.
func NewManageBalanceOpExt(v LedgerVersion, value interface{}) (result ManageBalanceOpExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ManageBalanceOp is an XDR Struct defines as:
//
//   struct ManageBalanceOp
//    {
//        BalanceID balanceID;
//        ManageBalanceAction action;
//        AccountID destination;
//        AssetCode asset;
//    	union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type ManageBalanceOp struct {
	BalanceId   BalanceId           `json:"balanceID,omitempty"`
	Action      ManageBalanceAction `json:"action,omitempty"`
	Destination AccountId           `json:"destination,omitempty"`
	Asset       AssetCode           `json:"asset,omitempty"`
	Ext         ManageBalanceOpExt  `json:"ext,omitempty"`
}

// ManageBalanceResultCode is an XDR Enum defines as:
//
//   enum ManageBalanceResultCode
//    {
//        // codes considered as "success" for the operation
//        MANAGE_BALANCE_SUCCESS = 0,
//
//        // codes considered as "failure" for the operation
//        MANAGE_BALANCE_MALFORMED = -1,       // invalid destination
//        MANAGE_BALANCE_NOT_FOUND = -2,
//        MANAGE_BALANCE_DESTINATION_NOT_FOUND = -3,
//        MANAGE_BALANCE_ALREADY_EXISTS = -4,
//        MANAGE_BALANCE_ANOTHER_EXCHANGE = -5,
//        MANAGE_BALANCE_ASSET_NOT_FOUND = -6,
//        MANAGE_BALANCE_INVALID_ASSET = -7,
//        MANAGE_BALANCE_NOT_ALLOWED_BY_EXCHANGE_POLICY = -8
//    };
//
type ManageBalanceResultCode int32

const (
	ManageBalanceResultCodeManageBalanceSuccess                    ManageBalanceResultCode = 0
	ManageBalanceResultCodeManageBalanceMalformed                  ManageBalanceResultCode = -1
	ManageBalanceResultCodeManageBalanceNotFound                   ManageBalanceResultCode = -2
	ManageBalanceResultCodeManageBalanceDestinationNotFound        ManageBalanceResultCode = -3
	ManageBalanceResultCodeManageBalanceAlreadyExists              ManageBalanceResultCode = -4
	ManageBalanceResultCodeManageBalanceAnotherExchange            ManageBalanceResultCode = -5
	ManageBalanceResultCodeManageBalanceAssetNotFound              ManageBalanceResultCode = -6
	ManageBalanceResultCodeManageBalanceInvalidAsset               ManageBalanceResultCode = -7
	ManageBalanceResultCodeManageBalanceNotAllowedByExchangePolicy ManageBalanceResultCode = -8
)

var ManageBalanceResultCodeAll = []ManageBalanceResultCode{
	ManageBalanceResultCodeManageBalanceSuccess,
	ManageBalanceResultCodeManageBalanceMalformed,
	ManageBalanceResultCodeManageBalanceNotFound,
	ManageBalanceResultCodeManageBalanceDestinationNotFound,
	ManageBalanceResultCodeManageBalanceAlreadyExists,
	ManageBalanceResultCodeManageBalanceAnotherExchange,
	ManageBalanceResultCodeManageBalanceAssetNotFound,
	ManageBalanceResultCodeManageBalanceInvalidAsset,
	ManageBalanceResultCodeManageBalanceNotAllowedByExchangePolicy,
}

var manageBalanceResultCodeMap = map[int32]string{
	0:  "ManageBalanceResultCodeManageBalanceSuccess",
	-1: "ManageBalanceResultCodeManageBalanceMalformed",
	-2: "ManageBalanceResultCodeManageBalanceNotFound",
	-3: "ManageBalanceResultCodeManageBalanceDestinationNotFound",
	-4: "ManageBalanceResultCodeManageBalanceAlreadyExists",
	-5: "ManageBalanceResultCodeManageBalanceAnotherExchange",
	-6: "ManageBalanceResultCodeManageBalanceAssetNotFound",
	-7: "ManageBalanceResultCodeManageBalanceInvalidAsset",
	-8: "ManageBalanceResultCodeManageBalanceNotAllowedByExchangePolicy",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ManageBalanceResultCode
func (e ManageBalanceResultCode) ValidEnum(v int32) bool {
	_, ok := manageBalanceResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e ManageBalanceResultCode) String() string {
	name, _ := manageBalanceResultCodeMap[int32(e)]
	return name
}

func (e ManageBalanceResultCode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// ManageBalanceSuccessExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type ManageBalanceSuccessExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageBalanceSuccessExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageBalanceSuccessExt
func (u ManageBalanceSuccessExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewManageBalanceSuccessExt creates a new  ManageBalanceSuccessExt.
func NewManageBalanceSuccessExt(v LedgerVersion, value interface{}) (result ManageBalanceSuccessExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ManageBalanceSuccess is an XDR Struct defines as:
//
//   struct ManageBalanceSuccess {
//    	// reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type ManageBalanceSuccess struct {
	Ext ManageBalanceSuccessExt `json:"ext,omitempty"`
}

// ManageBalanceResult is an XDR Union defines as:
//
//   union ManageBalanceResult switch (ManageBalanceResultCode code)
//    {
//    case MANAGE_BALANCE_SUCCESS:
//        ManageBalanceSuccess success;
//    default:
//        void;
//    };
//
type ManageBalanceResult struct {
	Code    ManageBalanceResultCode `json:"code,omitempty"`
	Success *ManageBalanceSuccess   `json:"success,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageBalanceResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageBalanceResult
func (u ManageBalanceResult) ArmForSwitch(sw int32) (string, bool) {
	switch ManageBalanceResultCode(sw) {
	case ManageBalanceResultCodeManageBalanceSuccess:
		return "Success", true
	default:
		return "", true
	}
}

// NewManageBalanceResult creates a new  ManageBalanceResult.
func NewManageBalanceResult(code ManageBalanceResultCode, value interface{}) (result ManageBalanceResult, err error) {
	result.Code = code
	switch ManageBalanceResultCode(code) {
	case ManageBalanceResultCodeManageBalanceSuccess:
		tv, ok := value.(ManageBalanceSuccess)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageBalanceSuccess")
			return
		}
		result.Success = &tv
	default:
		// void
	}
	return
}

// MustSuccess retrieves the Success value from the union,
// panicing if the value is not set.
func (u ManageBalanceResult) MustSuccess() ManageBalanceSuccess {
	val, ok := u.GetSuccess()

	if !ok {
		panic("arm Success is not set")
	}

	return val
}

// GetSuccess retrieves the Success value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ManageBalanceResult) GetSuccess() (result ManageBalanceSuccess, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "Success" {
		result = *u.Success
		ok = true
	}

	return
}

// ManageCoinsEmissionRequestAction is an XDR Enum defines as:
//
//   enum ManageCoinsEmissionRequestAction
//    {
//        MANAGE_COINS_EMISSION_REQUEST_CREATE = 0,
//        MANAGE_COINS_EMISSION_REQUEST_DELETE = 1
//    };
//
type ManageCoinsEmissionRequestAction int32

const (
	ManageCoinsEmissionRequestActionManageCoinsEmissionRequestCreate ManageCoinsEmissionRequestAction = 0
	ManageCoinsEmissionRequestActionManageCoinsEmissionRequestDelete ManageCoinsEmissionRequestAction = 1
)

var ManageCoinsEmissionRequestActionAll = []ManageCoinsEmissionRequestAction{
	ManageCoinsEmissionRequestActionManageCoinsEmissionRequestCreate,
	ManageCoinsEmissionRequestActionManageCoinsEmissionRequestDelete,
}

var manageCoinsEmissionRequestActionMap = map[int32]string{
	0: "ManageCoinsEmissionRequestActionManageCoinsEmissionRequestCreate",
	1: "ManageCoinsEmissionRequestActionManageCoinsEmissionRequestDelete",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ManageCoinsEmissionRequestAction
func (e ManageCoinsEmissionRequestAction) ValidEnum(v int32) bool {
	_, ok := manageCoinsEmissionRequestActionMap[v]
	return ok
}

// String returns the name of `e`
func (e ManageCoinsEmissionRequestAction) String() string {
	name, _ := manageCoinsEmissionRequestActionMap[int32(e)]
	return name
}

func (e ManageCoinsEmissionRequestAction) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// ManageCoinsEmissionRequestOpExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type ManageCoinsEmissionRequestOpExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageCoinsEmissionRequestOpExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageCoinsEmissionRequestOpExt
func (u ManageCoinsEmissionRequestOpExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewManageCoinsEmissionRequestOpExt creates a new  ManageCoinsEmissionRequestOpExt.
func NewManageCoinsEmissionRequestOpExt(v LedgerVersion, value interface{}) (result ManageCoinsEmissionRequestOpExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ManageCoinsEmissionRequestOp is an XDR Struct defines as:
//
//   struct ManageCoinsEmissionRequestOp
//    {
//    	// 0=create a new request, otherwise edit an existing offer
//        ManageCoinsEmissionRequestAction action;
//    	uint64 requestID;
//        int64 amount;        // amount being issued. if set to 0, delete the offer
//        BalanceID receiver;
//        AssetCode asset;
//        string64 reference;
//    	// reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type ManageCoinsEmissionRequestOp struct {
	Action    ManageCoinsEmissionRequestAction `json:"action,omitempty"`
	RequestId Uint64                           `json:"requestID,omitempty"`
	Amount    Int64                            `json:"amount,omitempty"`
	Receiver  BalanceId                        `json:"receiver,omitempty"`
	Asset     AssetCode                        `json:"asset,omitempty"`
	Reference String64                         `json:"reference,omitempty"`
	Ext       ManageCoinsEmissionRequestOpExt  `json:"ext,omitempty"`
}

// ManageCoinsEmissionRequestResultCode is an XDR Enum defines as:
//
//   enum ManageCoinsEmissionRequestResultCode
//    {
//        // codes considered as "success" for the operation
//        MANAGE_COINS_EMISSION_REQUEST_SUCCESS = 0,
//
//        // codes considered as "failure" for the operation
//        MANAGE_COINS_EMISSION_REQUEST_INVALID_AMOUNT = -1,      // amount is negative
//    	MANAGE_COINS_EMISSION_REQUEST_INVALID_REQUEST_ID = -2, // not 0 for delete etc
//    	MANAGE_COINS_EMISSION_REQUEST_NOT_FOUND = -3,           // failed to find emission request with such ID
//    	MANAGE_COINS_EMISSION_REQUEST_ALREADY_REVIEWED = -4,    // emission request have been already reviewed - can't edit
//        MANAGE_COINS_EMISSION_REQUEST_ASSET_NOT_FOUND = -5,
//        MANAGE_COINS_EMISSION_REQUEST_BALANCE_NOT_FOUND = -6,
//        MANAGE_COINS_EMISSION_REQUEST_ASSET_MISMATCH = -7,
//        MANAGE_COINS_EMISSION_REQUEST_INVALID_ASSET = -8,
//        MANAGE_COINS_EMISSION_REQUEST_REFERENCE_DUPLICATION = -9,
//        MANAGE_COINS_EMISSION_REQUEST_LINE_FULL = -10,
//        MANAGE_COINS_EMISSION_REQUEST_INVALID_REFERENCE = -11,
//        MANAGE_COINS_EMISSION_REQUEST_NOT_ALLOWED_BY_EXCHANGE_POLICY = -12
//    };
//
type ManageCoinsEmissionRequestResultCode int32

const (
	ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestSuccess                    ManageCoinsEmissionRequestResultCode = 0
	ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestInvalidAmount              ManageCoinsEmissionRequestResultCode = -1
	ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestInvalidRequestId           ManageCoinsEmissionRequestResultCode = -2
	ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestNotFound                   ManageCoinsEmissionRequestResultCode = -3
	ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestAlreadyReviewed            ManageCoinsEmissionRequestResultCode = -4
	ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestAssetNotFound              ManageCoinsEmissionRequestResultCode = -5
	ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestBalanceNotFound            ManageCoinsEmissionRequestResultCode = -6
	ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestAssetMismatch              ManageCoinsEmissionRequestResultCode = -7
	ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestInvalidAsset               ManageCoinsEmissionRequestResultCode = -8
	ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestReferenceDuplication       ManageCoinsEmissionRequestResultCode = -9
	ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestLineFull                   ManageCoinsEmissionRequestResultCode = -10
	ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestInvalidReference           ManageCoinsEmissionRequestResultCode = -11
	ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestNotAllowedByExchangePolicy ManageCoinsEmissionRequestResultCode = -12
)

var ManageCoinsEmissionRequestResultCodeAll = []ManageCoinsEmissionRequestResultCode{
	ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestSuccess,
	ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestInvalidAmount,
	ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestInvalidRequestId,
	ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestNotFound,
	ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestAlreadyReviewed,
	ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestAssetNotFound,
	ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestBalanceNotFound,
	ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestAssetMismatch,
	ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestInvalidAsset,
	ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestReferenceDuplication,
	ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestLineFull,
	ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestInvalidReference,
	ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestNotAllowedByExchangePolicy,
}

var manageCoinsEmissionRequestResultCodeMap = map[int32]string{
	0:   "ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestSuccess",
	-1:  "ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestInvalidAmount",
	-2:  "ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestInvalidRequestId",
	-3:  "ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestNotFound",
	-4:  "ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestAlreadyReviewed",
	-5:  "ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestAssetNotFound",
	-6:  "ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestBalanceNotFound",
	-7:  "ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestAssetMismatch",
	-8:  "ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestInvalidAsset",
	-9:  "ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestReferenceDuplication",
	-10: "ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestLineFull",
	-11: "ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestInvalidReference",
	-12: "ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestNotAllowedByExchangePolicy",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ManageCoinsEmissionRequestResultCode
func (e ManageCoinsEmissionRequestResultCode) ValidEnum(v int32) bool {
	_, ok := manageCoinsEmissionRequestResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e ManageCoinsEmissionRequestResultCode) String() string {
	name, _ := manageCoinsEmissionRequestResultCodeMap[int32(e)]
	return name
}

func (e ManageCoinsEmissionRequestResultCode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// ManageCoinsEmissionRequestResultManageRequestInfoExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type ManageCoinsEmissionRequestResultManageRequestInfoExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageCoinsEmissionRequestResultManageRequestInfoExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageCoinsEmissionRequestResultManageRequestInfoExt
func (u ManageCoinsEmissionRequestResultManageRequestInfoExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewManageCoinsEmissionRequestResultManageRequestInfoExt creates a new  ManageCoinsEmissionRequestResultManageRequestInfoExt.
func NewManageCoinsEmissionRequestResultManageRequestInfoExt(v LedgerVersion, value interface{}) (result ManageCoinsEmissionRequestResultManageRequestInfoExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ManageCoinsEmissionRequestResultManageRequestInfo is an XDR NestedStruct defines as:
//
//   struct {
//            uint64 requestID;
//            bool fulfilled;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        }
//
type ManageCoinsEmissionRequestResultManageRequestInfo struct {
	RequestId Uint64                                               `json:"requestID,omitempty"`
	Fulfilled bool                                                 `json:"fulfilled,omitempty"`
	Ext       ManageCoinsEmissionRequestResultManageRequestInfoExt `json:"ext,omitempty"`
}

// ManageCoinsEmissionRequestResult is an XDR Union defines as:
//
//   union ManageCoinsEmissionRequestResult switch (ManageCoinsEmissionRequestResultCode code)
//    {
//    case MANAGE_COINS_EMISSION_REQUEST_SUCCESS:
//        struct {
//            uint64 requestID;
//            bool fulfilled;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        } manageRequestInfo;
//    default:
//        void;
//    };
//
type ManageCoinsEmissionRequestResult struct {
	Code              ManageCoinsEmissionRequestResultCode               `json:"code,omitempty"`
	ManageRequestInfo *ManageCoinsEmissionRequestResultManageRequestInfo `json:"manageRequestInfo,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageCoinsEmissionRequestResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageCoinsEmissionRequestResult
func (u ManageCoinsEmissionRequestResult) ArmForSwitch(sw int32) (string, bool) {
	switch ManageCoinsEmissionRequestResultCode(sw) {
	case ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestSuccess:
		return "ManageRequestInfo", true
	default:
		return "", true
	}
}

// NewManageCoinsEmissionRequestResult creates a new  ManageCoinsEmissionRequestResult.
func NewManageCoinsEmissionRequestResult(code ManageCoinsEmissionRequestResultCode, value interface{}) (result ManageCoinsEmissionRequestResult, err error) {
	result.Code = code
	switch ManageCoinsEmissionRequestResultCode(code) {
	case ManageCoinsEmissionRequestResultCodeManageCoinsEmissionRequestSuccess:
		tv, ok := value.(ManageCoinsEmissionRequestResultManageRequestInfo)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageCoinsEmissionRequestResultManageRequestInfo")
			return
		}
		result.ManageRequestInfo = &tv
	default:
		// void
	}
	return
}

// MustManageRequestInfo retrieves the ManageRequestInfo value from the union,
// panicing if the value is not set.
func (u ManageCoinsEmissionRequestResult) MustManageRequestInfo() ManageCoinsEmissionRequestResultManageRequestInfo {
	val, ok := u.GetManageRequestInfo()

	if !ok {
		panic("arm ManageRequestInfo is not set")
	}

	return val
}

// GetManageRequestInfo retrieves the ManageRequestInfo value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ManageCoinsEmissionRequestResult) GetManageRequestInfo() (result ManageCoinsEmissionRequestResultManageRequestInfo, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "ManageRequestInfo" {
		result = *u.ManageRequestInfo
		ok = true
	}

	return
}

// ManageForfeitRequestOpExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type ManageForfeitRequestOpExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageForfeitRequestOpExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageForfeitRequestOpExt
func (u ManageForfeitRequestOpExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewManageForfeitRequestOpExt creates a new  ManageForfeitRequestOpExt.
func NewManageForfeitRequestOpExt(v LedgerVersion, value interface{}) (result ManageForfeitRequestOpExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ManageForfeitRequestOp is an XDR Struct defines as:
//
//   struct ManageForfeitRequestOp
//    {
//        BalanceID balance;
//        int64 amount;
//        string details<>;
//    	AccountID* reviewer;
//    	union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//
//    };
//
type ManageForfeitRequestOp struct {
	Balance  BalanceId                 `json:"balance,omitempty"`
	Amount   Int64                     `json:"amount,omitempty"`
	Details  string                    `json:"details,omitempty"`
	Reviewer *AccountId                `json:"reviewer,omitempty"`
	Ext      ManageForfeitRequestOpExt `json:"ext,omitempty"`
}

// ManageForfeitRequestResultCode is an XDR Enum defines as:
//
//   enum ManageForfeitRequestResultCode
//    {
//        // codes considered as "success" for the operation
//        MANAGE_FORFEIT_REQUEST_SUCCESS = 0,
//
//        // codes considered as "failure" for the operation
//    	MANAGE_FORFEIT_REQUEST_UNDERFUNDED = -1,
//        MANAGE_FORFEIT_REQUEST_INVALID_AMOUNT = -2,
//        MANAGE_FORFEIT_REQUEST_LINE_FULL = -3,
//        MANAGE_FORFEIT_REQUEST_BALANCE_MISMATCH = -4,
//        MANAGE_FORFEIT_REQUEST_STATS_OVERFLOW = -5,
//        MANAGE_FORFEIT_REQUEST_LIMITS_EXCEEDED = -6,
//        MANAGE_FORFEIT_REQUEST_REVIEWER_NOT_FOUND = -7,
//        MANAGE_FORFEIT_REQUEST_INVALID_DETAILS = -8,
//    	MANAGE_FORFEIT_REQUEST_BALANCE_REQUIRES_REVIEW = -9 // if reviewer different from exchange holding the balance and exchange holding the balance requires review
//    };
//
type ManageForfeitRequestResultCode int32

const (
	ManageForfeitRequestResultCodeManageForfeitRequestSuccess               ManageForfeitRequestResultCode = 0
	ManageForfeitRequestResultCodeManageForfeitRequestUnderfunded           ManageForfeitRequestResultCode = -1
	ManageForfeitRequestResultCodeManageForfeitRequestInvalidAmount         ManageForfeitRequestResultCode = -2
	ManageForfeitRequestResultCodeManageForfeitRequestLineFull              ManageForfeitRequestResultCode = -3
	ManageForfeitRequestResultCodeManageForfeitRequestBalanceMismatch       ManageForfeitRequestResultCode = -4
	ManageForfeitRequestResultCodeManageForfeitRequestStatsOverflow         ManageForfeitRequestResultCode = -5
	ManageForfeitRequestResultCodeManageForfeitRequestLimitsExceeded        ManageForfeitRequestResultCode = -6
	ManageForfeitRequestResultCodeManageForfeitRequestReviewerNotFound      ManageForfeitRequestResultCode = -7
	ManageForfeitRequestResultCodeManageForfeitRequestInvalidDetails        ManageForfeitRequestResultCode = -8
	ManageForfeitRequestResultCodeManageForfeitRequestBalanceRequiresReview ManageForfeitRequestResultCode = -9
)

var ManageForfeitRequestResultCodeAll = []ManageForfeitRequestResultCode{
	ManageForfeitRequestResultCodeManageForfeitRequestSuccess,
	ManageForfeitRequestResultCodeManageForfeitRequestUnderfunded,
	ManageForfeitRequestResultCodeManageForfeitRequestInvalidAmount,
	ManageForfeitRequestResultCodeManageForfeitRequestLineFull,
	ManageForfeitRequestResultCodeManageForfeitRequestBalanceMismatch,
	ManageForfeitRequestResultCodeManageForfeitRequestStatsOverflow,
	ManageForfeitRequestResultCodeManageForfeitRequestLimitsExceeded,
	ManageForfeitRequestResultCodeManageForfeitRequestReviewerNotFound,
	ManageForfeitRequestResultCodeManageForfeitRequestInvalidDetails,
	ManageForfeitRequestResultCodeManageForfeitRequestBalanceRequiresReview,
}

var manageForfeitRequestResultCodeMap = map[int32]string{
	0:  "ManageForfeitRequestResultCodeManageForfeitRequestSuccess",
	-1: "ManageForfeitRequestResultCodeManageForfeitRequestUnderfunded",
	-2: "ManageForfeitRequestResultCodeManageForfeitRequestInvalidAmount",
	-3: "ManageForfeitRequestResultCodeManageForfeitRequestLineFull",
	-4: "ManageForfeitRequestResultCodeManageForfeitRequestBalanceMismatch",
	-5: "ManageForfeitRequestResultCodeManageForfeitRequestStatsOverflow",
	-6: "ManageForfeitRequestResultCodeManageForfeitRequestLimitsExceeded",
	-7: "ManageForfeitRequestResultCodeManageForfeitRequestReviewerNotFound",
	-8: "ManageForfeitRequestResultCodeManageForfeitRequestInvalidDetails",
	-9: "ManageForfeitRequestResultCodeManageForfeitRequestBalanceRequiresReview",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ManageForfeitRequestResultCode
func (e ManageForfeitRequestResultCode) ValidEnum(v int32) bool {
	_, ok := manageForfeitRequestResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e ManageForfeitRequestResultCode) String() string {
	name, _ := manageForfeitRequestResultCodeMap[int32(e)]
	return name
}

func (e ManageForfeitRequestResultCode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// ForfeitItemExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type ForfeitItemExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ForfeitItemExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ForfeitItemExt
func (u ForfeitItemExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewForfeitItemExt creates a new  ForfeitItemExt.
func NewForfeitItemExt(v LedgerVersion, value interface{}) (result ForfeitItemExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ForfeitItem is an XDR Struct defines as:
//
//   struct ForfeitItem
//    {
//        AssetForm form;
//        int64 quantity;
//    	// reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type ForfeitItem struct {
	Form     AssetForm      `json:"form,omitempty"`
	Quantity Int64          `json:"quantity,omitempty"`
	Ext      ForfeitItemExt `json:"ext,omitempty"`
}

// ManageForfeitRequestResultForfeitRequestDetailsFees is an XDR NestedStruct defines as:
//
//   struct {
//    			int64 percentFee;
//    			int64 fixedFee;
//    		}
//
type ManageForfeitRequestResultForfeitRequestDetailsFees struct {
	PercentFee Int64 `json:"percentFee,omitempty"`
	FixedFee   Int64 `json:"fixedFee,omitempty"`
}

// ManageForfeitRequestResultForfeitRequestDetailsExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		case FORFEIT_RESULT_FEES:
//    		struct {
//    			int64 percentFee;
//    			int64 fixedFee;
//    		} fees;
//
//    		}
//
type ManageForfeitRequestResultForfeitRequestDetailsExt struct {
	V    LedgerVersion                                        `json:"v,omitempty"`
	Fees *ManageForfeitRequestResultForfeitRequestDetailsFees `json:"fees,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageForfeitRequestResultForfeitRequestDetailsExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageForfeitRequestResultForfeitRequestDetailsExt
func (u ManageForfeitRequestResultForfeitRequestDetailsExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	case LedgerVersionForfeitResultFees:
		return "Fees", true
	}
	return "-", false
}

// NewManageForfeitRequestResultForfeitRequestDetailsExt creates a new  ManageForfeitRequestResultForfeitRequestDetailsExt.
func NewManageForfeitRequestResultForfeitRequestDetailsExt(v LedgerVersion, value interface{}) (result ManageForfeitRequestResultForfeitRequestDetailsExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	case LedgerVersionForfeitResultFees:
		tv, ok := value.(ManageForfeitRequestResultForfeitRequestDetailsFees)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageForfeitRequestResultForfeitRequestDetailsFees")
			return
		}
		result.Fees = &tv
	}
	return
}

// MustFees retrieves the Fees value from the union,
// panicing if the value is not set.
func (u ManageForfeitRequestResultForfeitRequestDetailsExt) MustFees() ManageForfeitRequestResultForfeitRequestDetailsFees {
	val, ok := u.GetFees()

	if !ok {
		panic("arm Fees is not set")
	}

	return val
}

// GetFees retrieves the Fees value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ManageForfeitRequestResultForfeitRequestDetailsExt) GetFees() (result ManageForfeitRequestResultForfeitRequestDetailsFees, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.V))

	if armName == "Fees" {
		result = *u.Fees
		ok = true
	}

	return
}

// ManageForfeitRequestResultForfeitRequestDetails is an XDR NestedStruct defines as:
//
//   struct {
//            uint64 paymentID;
//            AccountID exchange;
//            AssetCode asset;
//            ForfeitItem items<>;
//    		// reserved for future use
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		case FORFEIT_RESULT_FEES:
//    		struct {
//    			int64 percentFee;
//    			int64 fixedFee;
//    		} fees;
//
//    		}
//    		ext;
//        }
//
type ManageForfeitRequestResultForfeitRequestDetails struct {
	PaymentId Uint64                                             `json:"paymentID,omitempty"`
	Exchange  AccountId                                          `json:"exchange,omitempty"`
	Asset     AssetCode                                          `json:"asset,omitempty"`
	Items     []ForfeitItem                                      `json:"items,omitempty"`
	Ext       ManageForfeitRequestResultForfeitRequestDetailsExt `json:"ext,omitempty"`
}

// ManageForfeitRequestResult is an XDR Union defines as:
//
//   union ManageForfeitRequestResult switch (ManageForfeitRequestResultCode code)
//    {
//    case MANAGE_FORFEIT_REQUEST_SUCCESS:
//        struct {
//            uint64 paymentID;
//            AccountID exchange;
//            AssetCode asset;
//            ForfeitItem items<>;
//    		// reserved for future use
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		case FORFEIT_RESULT_FEES:
//    		struct {
//    			int64 percentFee;
//    			int64 fixedFee;
//    		} fees;
//
//    		}
//    		ext;
//        } forfeitRequestDetails;
//    default:
//        void;
//    };
//
type ManageForfeitRequestResult struct {
	Code                  ManageForfeitRequestResultCode                   `json:"code,omitempty"`
	ForfeitRequestDetails *ManageForfeitRequestResultForfeitRequestDetails `json:"forfeitRequestDetails,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageForfeitRequestResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageForfeitRequestResult
func (u ManageForfeitRequestResult) ArmForSwitch(sw int32) (string, bool) {
	switch ManageForfeitRequestResultCode(sw) {
	case ManageForfeitRequestResultCodeManageForfeitRequestSuccess:
		return "ForfeitRequestDetails", true
	default:
		return "", true
	}
}

// NewManageForfeitRequestResult creates a new  ManageForfeitRequestResult.
func NewManageForfeitRequestResult(code ManageForfeitRequestResultCode, value interface{}) (result ManageForfeitRequestResult, err error) {
	result.Code = code
	switch ManageForfeitRequestResultCode(code) {
	case ManageForfeitRequestResultCodeManageForfeitRequestSuccess:
		tv, ok := value.(ManageForfeitRequestResultForfeitRequestDetails)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageForfeitRequestResultForfeitRequestDetails")
			return
		}
		result.ForfeitRequestDetails = &tv
	default:
		// void
	}
	return
}

// MustForfeitRequestDetails retrieves the ForfeitRequestDetails value from the union,
// panicing if the value is not set.
func (u ManageForfeitRequestResult) MustForfeitRequestDetails() ManageForfeitRequestResultForfeitRequestDetails {
	val, ok := u.GetForfeitRequestDetails()

	if !ok {
		panic("arm ForfeitRequestDetails is not set")
	}

	return val
}

// GetForfeitRequestDetails retrieves the ForfeitRequestDetails value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ManageForfeitRequestResult) GetForfeitRequestDetails() (result ManageForfeitRequestResultForfeitRequestDetails, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "ForfeitRequestDetails" {
		result = *u.ForfeitRequestDetails
		ok = true
	}

	return
}

// ManageInvoiceOpExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type ManageInvoiceOpExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageInvoiceOpExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageInvoiceOpExt
func (u ManageInvoiceOpExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewManageInvoiceOpExt creates a new  ManageInvoiceOpExt.
func NewManageInvoiceOpExt(v LedgerVersion, value interface{}) (result ManageInvoiceOpExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ManageInvoiceOp is an XDR Struct defines as:
//
//   struct ManageInvoiceOp
//    {
//        BalanceID receiverBalance;
//    	AccountID sender;
//        int64 amount; // if set to 0, delete the invoice
//
//        // 0=create a new invoice, otherwise edit an existing invoice
//        uint64 invoiceID;
//
//    	// reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type ManageInvoiceOp struct {
	ReceiverBalance BalanceId          `json:"receiverBalance,omitempty"`
	Sender          AccountId          `json:"sender,omitempty"`
	Amount          Int64              `json:"amount,omitempty"`
	InvoiceId       Uint64             `json:"invoiceID,omitempty"`
	Ext             ManageInvoiceOpExt `json:"ext,omitempty"`
}

// ManageInvoiceResultCode is an XDR Enum defines as:
//
//   enum ManageInvoiceResultCode
//    {
//        // codes considered as "success" for the operation
//        MANAGE_INVOICE_SUCCESS = 0,
//
//        // codes considered as "failure" for the operation
//        MANAGE_INVOICE_MALFORMED = -1,
//        MANAGE_INVOICE_BALANCE_NOT_FOUND = -2,
//    	MANAGE_INVOICE_OVERFLOW = -3,
//
//        MANAGE_INVOICE_NOT_FOUND = -4,
//        MANAGE_INVOICE_TOO_MANY_INVOICES = -5,
//        MANAGE_INVOICE_CAN_NOT_DELETE_IN_PROGRESS = -6
//    };
//
type ManageInvoiceResultCode int32

const (
	ManageInvoiceResultCodeManageInvoiceSuccess                ManageInvoiceResultCode = 0
	ManageInvoiceResultCodeManageInvoiceMalformed              ManageInvoiceResultCode = -1
	ManageInvoiceResultCodeManageInvoiceBalanceNotFound        ManageInvoiceResultCode = -2
	ManageInvoiceResultCodeManageInvoiceOverflow               ManageInvoiceResultCode = -3
	ManageInvoiceResultCodeManageInvoiceNotFound               ManageInvoiceResultCode = -4
	ManageInvoiceResultCodeManageInvoiceTooManyInvoices        ManageInvoiceResultCode = -5
	ManageInvoiceResultCodeManageInvoiceCanNotDeleteInProgress ManageInvoiceResultCode = -6
)

var ManageInvoiceResultCodeAll = []ManageInvoiceResultCode{
	ManageInvoiceResultCodeManageInvoiceSuccess,
	ManageInvoiceResultCodeManageInvoiceMalformed,
	ManageInvoiceResultCodeManageInvoiceBalanceNotFound,
	ManageInvoiceResultCodeManageInvoiceOverflow,
	ManageInvoiceResultCodeManageInvoiceNotFound,
	ManageInvoiceResultCodeManageInvoiceTooManyInvoices,
	ManageInvoiceResultCodeManageInvoiceCanNotDeleteInProgress,
}

var manageInvoiceResultCodeMap = map[int32]string{
	0:  "ManageInvoiceResultCodeManageInvoiceSuccess",
	-1: "ManageInvoiceResultCodeManageInvoiceMalformed",
	-2: "ManageInvoiceResultCodeManageInvoiceBalanceNotFound",
	-3: "ManageInvoiceResultCodeManageInvoiceOverflow",
	-4: "ManageInvoiceResultCodeManageInvoiceNotFound",
	-5: "ManageInvoiceResultCodeManageInvoiceTooManyInvoices",
	-6: "ManageInvoiceResultCodeManageInvoiceCanNotDeleteInProgress",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ManageInvoiceResultCode
func (e ManageInvoiceResultCode) ValidEnum(v int32) bool {
	_, ok := manageInvoiceResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e ManageInvoiceResultCode) String() string {
	name, _ := manageInvoiceResultCodeMap[int32(e)]
	return name
}

func (e ManageInvoiceResultCode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// ManageInvoiceSuccessResultExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type ManageInvoiceSuccessResultExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageInvoiceSuccessResultExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageInvoiceSuccessResultExt
func (u ManageInvoiceSuccessResultExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewManageInvoiceSuccessResultExt creates a new  ManageInvoiceSuccessResultExt.
func NewManageInvoiceSuccessResultExt(v LedgerVersion, value interface{}) (result ManageInvoiceSuccessResultExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ManageInvoiceSuccessResult is an XDR Struct defines as:
//
//   struct ManageInvoiceSuccessResult
//    {
//    	uint64 invoiceID;
//    	AssetCode asset;
//    	BalanceID senderBalance;
//
//    	union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type ManageInvoiceSuccessResult struct {
	InvoiceId     Uint64                        `json:"invoiceID,omitempty"`
	Asset         AssetCode                     `json:"asset,omitempty"`
	SenderBalance BalanceId                     `json:"senderBalance,omitempty"`
	Ext           ManageInvoiceSuccessResultExt `json:"ext,omitempty"`
}

// ManageInvoiceResult is an XDR Union defines as:
//
//   union ManageInvoiceResult switch (ManageInvoiceResultCode code)
//    {
//    case MANAGE_INVOICE_SUCCESS:
//        ManageInvoiceSuccessResult success;
//    default:
//        void;
//    };
//
type ManageInvoiceResult struct {
	Code    ManageInvoiceResultCode     `json:"code,omitempty"`
	Success *ManageInvoiceSuccessResult `json:"success,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageInvoiceResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageInvoiceResult
func (u ManageInvoiceResult) ArmForSwitch(sw int32) (string, bool) {
	switch ManageInvoiceResultCode(sw) {
	case ManageInvoiceResultCodeManageInvoiceSuccess:
		return "Success", true
	default:
		return "", true
	}
}

// NewManageInvoiceResult creates a new  ManageInvoiceResult.
func NewManageInvoiceResult(code ManageInvoiceResultCode, value interface{}) (result ManageInvoiceResult, err error) {
	result.Code = code
	switch ManageInvoiceResultCode(code) {
	case ManageInvoiceResultCodeManageInvoiceSuccess:
		tv, ok := value.(ManageInvoiceSuccessResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageInvoiceSuccessResult")
			return
		}
		result.Success = &tv
	default:
		// void
	}
	return
}

// MustSuccess retrieves the Success value from the union,
// panicing if the value is not set.
func (u ManageInvoiceResult) MustSuccess() ManageInvoiceSuccessResult {
	val, ok := u.GetSuccess()

	if !ok {
		panic("arm Success is not set")
	}

	return val
}

// GetSuccess retrieves the Success value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ManageInvoiceResult) GetSuccess() (result ManageInvoiceSuccessResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "Success" {
		result = *u.Success
		ok = true
	}

	return
}

// ManageOfferOpExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type ManageOfferOpExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageOfferOpExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageOfferOpExt
func (u ManageOfferOpExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewManageOfferOpExt creates a new  ManageOfferOpExt.
func NewManageOfferOpExt(v LedgerVersion, value interface{}) (result ManageOfferOpExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ManageOfferOp is an XDR Struct defines as:
//
//   struct ManageOfferOp
//    {
//        BalanceID baseBalance; // balance for base asset
//    	BalanceID quoteBalance; // balance for quote asset
//    	bool isBuy;
//        int64 amount; // if set to 0, delete the offer
//        int64 price;  // price of base asset in terms of quote
//
//        int64 fee;
//
//    	bool isDirect;
//
//        // 0=create a new offer, otherwise edit an existing offer
//        uint64 offerID;
//    	// reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type ManageOfferOp struct {
	BaseBalance  BalanceId        `json:"baseBalance,omitempty"`
	QuoteBalance BalanceId        `json:"quoteBalance,omitempty"`
	IsBuy        bool             `json:"isBuy,omitempty"`
	Amount       Int64            `json:"amount,omitempty"`
	Price        Int64            `json:"price,omitempty"`
	Fee          Int64            `json:"fee,omitempty"`
	IsDirect     bool             `json:"isDirect,omitempty"`
	OfferId      Uint64           `json:"offerID,omitempty"`
	Ext          ManageOfferOpExt `json:"ext,omitempty"`
}

// ManageOfferResultCode is an XDR Enum defines as:
//
//   enum ManageOfferResultCode
//    {
//        // codes considered as "success" for the operation
//        MANAGE_OFFER_SUCCESS = 0,
//
//        // codes considered as "failure" for the operation
//        MANAGE_OFFER_MALFORMED = -1,     // generated offer would be invalid
//        MANAGE_OFFER_PAIR_NOT_TRADED = -2, // it's not allowed to trage with this pair
//        MANAGE_OFFER_BALANCE_NOT_FOUND = -3,  // does not own balance for buying or selling
//        MANAGE_OFFER_REQUIRES_REVIEW = -5, // exchange requires review for one of the balances
//        MANAGE_OFFER_UNDERFUNDED = -6,    // doesn't hold what it's trying to sell
//        MANAGE_OFFER_CROSS_SELF = -7,     // would cross an offer from the same user
//    	MANAGE_OFFER_OVERFLOW = -8,
//    	MANAGE_OFFER_ASSET_PAIR_NOT_TRADABLE = -9,
//    	MANAGE_OFFER_PHYSICAL_PRICE_RESTRICTION = -10, // offer price violates physical price restriction
//    	MAANGE_OFFER_CURRENT_PRICE_RESTRICTION = -11,
//        MANAGE_OFFER_NOT_FOUND = -12, // offerID does not match an existing offer
//        MANAGE_OFFER_INVALID_PERCENT_FEE = -13,
//    	MANAGE_OFFER_DIRECT_BUY_NOT_ALLOWED = -14,
//    	MANAGE_OFFER_INSUFFISIENT_PRICE = -15
//
//
//    };
//
type ManageOfferResultCode int32

const (
	ManageOfferResultCodeManageOfferSuccess                  ManageOfferResultCode = 0
	ManageOfferResultCodeManageOfferMalformed                ManageOfferResultCode = -1
	ManageOfferResultCodeManageOfferPairNotTraded            ManageOfferResultCode = -2
	ManageOfferResultCodeManageOfferBalanceNotFound          ManageOfferResultCode = -3
	ManageOfferResultCodeManageOfferRequiresReview           ManageOfferResultCode = -5
	ManageOfferResultCodeManageOfferUnderfunded              ManageOfferResultCode = -6
	ManageOfferResultCodeManageOfferCrossSelf                ManageOfferResultCode = -7
	ManageOfferResultCodeManageOfferOverflow                 ManageOfferResultCode = -8
	ManageOfferResultCodeManageOfferAssetPairNotTradable     ManageOfferResultCode = -9
	ManageOfferResultCodeManageOfferPhysicalPriceRestriction ManageOfferResultCode = -10
	ManageOfferResultCodeMaangeOfferCurrentPriceRestriction  ManageOfferResultCode = -11
	ManageOfferResultCodeManageOfferNotFound                 ManageOfferResultCode = -12
	ManageOfferResultCodeManageOfferInvalidPercentFee        ManageOfferResultCode = -13
	ManageOfferResultCodeManageOfferDirectBuyNotAllowed      ManageOfferResultCode = -14
	ManageOfferResultCodeManageOfferInsuffisientPrice        ManageOfferResultCode = -15
)

var ManageOfferResultCodeAll = []ManageOfferResultCode{
	ManageOfferResultCodeManageOfferSuccess,
	ManageOfferResultCodeManageOfferMalformed,
	ManageOfferResultCodeManageOfferPairNotTraded,
	ManageOfferResultCodeManageOfferBalanceNotFound,
	ManageOfferResultCodeManageOfferRequiresReview,
	ManageOfferResultCodeManageOfferUnderfunded,
	ManageOfferResultCodeManageOfferCrossSelf,
	ManageOfferResultCodeManageOfferOverflow,
	ManageOfferResultCodeManageOfferAssetPairNotTradable,
	ManageOfferResultCodeManageOfferPhysicalPriceRestriction,
	ManageOfferResultCodeMaangeOfferCurrentPriceRestriction,
	ManageOfferResultCodeManageOfferNotFound,
	ManageOfferResultCodeManageOfferInvalidPercentFee,
	ManageOfferResultCodeManageOfferDirectBuyNotAllowed,
	ManageOfferResultCodeManageOfferInsuffisientPrice,
}

var manageOfferResultCodeMap = map[int32]string{
	0:   "ManageOfferResultCodeManageOfferSuccess",
	-1:  "ManageOfferResultCodeManageOfferMalformed",
	-2:  "ManageOfferResultCodeManageOfferPairNotTraded",
	-3:  "ManageOfferResultCodeManageOfferBalanceNotFound",
	-5:  "ManageOfferResultCodeManageOfferRequiresReview",
	-6:  "ManageOfferResultCodeManageOfferUnderfunded",
	-7:  "ManageOfferResultCodeManageOfferCrossSelf",
	-8:  "ManageOfferResultCodeManageOfferOverflow",
	-9:  "ManageOfferResultCodeManageOfferAssetPairNotTradable",
	-10: "ManageOfferResultCodeManageOfferPhysicalPriceRestriction",
	-11: "ManageOfferResultCodeMaangeOfferCurrentPriceRestriction",
	-12: "ManageOfferResultCodeManageOfferNotFound",
	-13: "ManageOfferResultCodeManageOfferInvalidPercentFee",
	-14: "ManageOfferResultCodeManageOfferDirectBuyNotAllowed",
	-15: "ManageOfferResultCodeManageOfferInsuffisientPrice",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ManageOfferResultCode
func (e ManageOfferResultCode) ValidEnum(v int32) bool {
	_, ok := manageOfferResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e ManageOfferResultCode) String() string {
	name, _ := manageOfferResultCodeMap[int32(e)]
	return name
}

func (e ManageOfferResultCode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// ManageOfferEffect is an XDR Enum defines as:
//
//   enum ManageOfferEffect
//    {
//        MANAGE_OFFER_CREATED = 0,
//        MANAGE_OFFER_UPDATED = 1,
//        MANAGE_OFFER_DELETED = 2
//    };
//
type ManageOfferEffect int32

const (
	ManageOfferEffectManageOfferCreated ManageOfferEffect = 0
	ManageOfferEffectManageOfferUpdated ManageOfferEffect = 1
	ManageOfferEffectManageOfferDeleted ManageOfferEffect = 2
)

var ManageOfferEffectAll = []ManageOfferEffect{
	ManageOfferEffectManageOfferCreated,
	ManageOfferEffectManageOfferUpdated,
	ManageOfferEffectManageOfferDeleted,
}

var manageOfferEffectMap = map[int32]string{
	0: "ManageOfferEffectManageOfferCreated",
	1: "ManageOfferEffectManageOfferUpdated",
	2: "ManageOfferEffectManageOfferDeleted",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ManageOfferEffect
func (e ManageOfferEffect) ValidEnum(v int32) bool {
	_, ok := manageOfferEffectMap[v]
	return ok
}

// String returns the name of `e`
func (e ManageOfferEffect) String() string {
	name, _ := manageOfferEffectMap[int32(e)]
	return name
}

func (e ManageOfferEffect) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// ClaimOfferAtomExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type ClaimOfferAtomExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ClaimOfferAtomExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ClaimOfferAtomExt
func (u ClaimOfferAtomExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewClaimOfferAtomExt creates a new  ClaimOfferAtomExt.
func NewClaimOfferAtomExt(v LedgerVersion, value interface{}) (result ClaimOfferAtomExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ClaimOfferAtom is an XDR Struct defines as:
//
//   struct ClaimOfferAtom
//    {
//        // emitted to identify the offer
//        AccountID bAccountID; // Account that owns the offer
//        uint64 offerID;
//    	int64 baseAmount;
//    	int64 quoteAmount;
//    	int64 bFeePaid;
//    	int64 aFeePaid;
//    	BalanceID baseBalance;
//    	BalanceID quoteBalance;
//
//    	int64 currentPrice;
//
//    	union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type ClaimOfferAtom struct {
	BAccountId   AccountId         `json:"bAccountID,omitempty"`
	OfferId      Uint64            `json:"offerID,omitempty"`
	BaseAmount   Int64             `json:"baseAmount,omitempty"`
	QuoteAmount  Int64             `json:"quoteAmount,omitempty"`
	BFeePaid     Int64             `json:"bFeePaid,omitempty"`
	AFeePaid     Int64             `json:"aFeePaid,omitempty"`
	BaseBalance  BalanceId         `json:"baseBalance,omitempty"`
	QuoteBalance BalanceId         `json:"quoteBalance,omitempty"`
	CurrentPrice Int64             `json:"currentPrice,omitempty"`
	Ext          ClaimOfferAtomExt `json:"ext,omitempty"`
}

// ManageOfferSuccessResultOffer is an XDR NestedUnion defines as:
//
//   union switch (ManageOfferEffect effect)
//        {
//        case MANAGE_OFFER_CREATED:
//        case MANAGE_OFFER_UPDATED:
//            OfferEntry offer;
//        default:
//            void;
//        }
//
type ManageOfferSuccessResultOffer struct {
	Effect ManageOfferEffect `json:"effect,omitempty"`
	Offer  *OfferEntry       `json:"offer,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageOfferSuccessResultOffer) SwitchFieldName() string {
	return "Effect"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageOfferSuccessResultOffer
func (u ManageOfferSuccessResultOffer) ArmForSwitch(sw int32) (string, bool) {
	switch ManageOfferEffect(sw) {
	case ManageOfferEffectManageOfferCreated:
		return "Offer", true
	case ManageOfferEffectManageOfferUpdated:
		return "Offer", true
	default:
		return "", true
	}
}

// NewManageOfferSuccessResultOffer creates a new  ManageOfferSuccessResultOffer.
func NewManageOfferSuccessResultOffer(effect ManageOfferEffect, value interface{}) (result ManageOfferSuccessResultOffer, err error) {
	result.Effect = effect
	switch ManageOfferEffect(effect) {
	case ManageOfferEffectManageOfferCreated:
		tv, ok := value.(OfferEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be OfferEntry")
			return
		}
		result.Offer = &tv
	case ManageOfferEffectManageOfferUpdated:
		tv, ok := value.(OfferEntry)
		if !ok {
			err = fmt.Errorf("invalid value, must be OfferEntry")
			return
		}
		result.Offer = &tv
	default:
		// void
	}
	return
}

// MustOffer retrieves the Offer value from the union,
// panicing if the value is not set.
func (u ManageOfferSuccessResultOffer) MustOffer() OfferEntry {
	val, ok := u.GetOffer()

	if !ok {
		panic("arm Offer is not set")
	}

	return val
}

// GetOffer retrieves the Offer value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ManageOfferSuccessResultOffer) GetOffer() (result OfferEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Effect))

	if armName == "Offer" {
		result = *u.Offer
		ok = true
	}

	return
}

// ManageOfferSuccessResultExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type ManageOfferSuccessResultExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageOfferSuccessResultExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageOfferSuccessResultExt
func (u ManageOfferSuccessResultExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewManageOfferSuccessResultExt creates a new  ManageOfferSuccessResultExt.
func NewManageOfferSuccessResultExt(v LedgerVersion, value interface{}) (result ManageOfferSuccessResultExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ManageOfferSuccessResult is an XDR Struct defines as:
//
//   struct ManageOfferSuccessResult
//    {
//
//        // offers that got claimed while creating this offer
//        ClaimOfferAtom offersClaimed<>;
//    	AssetCode baseAsset;
//    	AssetCode quoteAsset;
//
//        union switch (ManageOfferEffect effect)
//        {
//        case MANAGE_OFFER_CREATED:
//        case MANAGE_OFFER_UPDATED:
//            OfferEntry offer;
//        default:
//            void;
//        }
//        offer;
//
//    	union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type ManageOfferSuccessResult struct {
	OffersClaimed []ClaimOfferAtom              `json:"offersClaimed,omitempty"`
	BaseAsset     AssetCode                     `json:"baseAsset,omitempty"`
	QuoteAsset    AssetCode                     `json:"quoteAsset,omitempty"`
	Offer         ManageOfferSuccessResultOffer `json:"offer,omitempty"`
	Ext           ManageOfferSuccessResultExt   `json:"ext,omitempty"`
}

// ManageOfferResultPhysicalPriceRestrictionExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type ManageOfferResultPhysicalPriceRestrictionExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageOfferResultPhysicalPriceRestrictionExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageOfferResultPhysicalPriceRestrictionExt
func (u ManageOfferResultPhysicalPriceRestrictionExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewManageOfferResultPhysicalPriceRestrictionExt creates a new  ManageOfferResultPhysicalPriceRestrictionExt.
func NewManageOfferResultPhysicalPriceRestrictionExt(v LedgerVersion, value interface{}) (result ManageOfferResultPhysicalPriceRestrictionExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ManageOfferResultPhysicalPriceRestriction is an XDR NestedStruct defines as:
//
//   struct {
//    		int64 physicalPrice;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//    	}
//
type ManageOfferResultPhysicalPriceRestriction struct {
	PhysicalPrice Int64                                        `json:"physicalPrice,omitempty"`
	Ext           ManageOfferResultPhysicalPriceRestrictionExt `json:"ext,omitempty"`
}

// ManageOfferResultCurrentPriceRestrictionExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type ManageOfferResultCurrentPriceRestrictionExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageOfferResultCurrentPriceRestrictionExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageOfferResultCurrentPriceRestrictionExt
func (u ManageOfferResultCurrentPriceRestrictionExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewManageOfferResultCurrentPriceRestrictionExt creates a new  ManageOfferResultCurrentPriceRestrictionExt.
func NewManageOfferResultCurrentPriceRestrictionExt(v LedgerVersion, value interface{}) (result ManageOfferResultCurrentPriceRestrictionExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ManageOfferResultCurrentPriceRestriction is an XDR NestedStruct defines as:
//
//   struct {
//    		int64 currentPrice;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//    	}
//
type ManageOfferResultCurrentPriceRestriction struct {
	CurrentPrice Int64                                       `json:"currentPrice,omitempty"`
	Ext          ManageOfferResultCurrentPriceRestrictionExt `json:"ext,omitempty"`
}

// ManageOfferResult is an XDR Union defines as:
//
//   union ManageOfferResult switch (ManageOfferResultCode code)
//    {
//    case MANAGE_OFFER_SUCCESS:
//        ManageOfferSuccessResult success;
//    case MANAGE_OFFER_PHYSICAL_PRICE_RESTRICTION:
//    	struct {
//    		int64 physicalPrice;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//    	} physicalPriceRestriction;
//    case MAANGE_OFFER_CURRENT_PRICE_RESTRICTION:
//    	struct {
//    		int64 currentPrice;
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//    	} currentPriceRestriction;
//
//    default:
//        void;
//    };
//
type ManageOfferResult struct {
	Code                     ManageOfferResultCode                      `json:"code,omitempty"`
	Success                  *ManageOfferSuccessResult                  `json:"success,omitempty"`
	PhysicalPriceRestriction *ManageOfferResultPhysicalPriceRestriction `json:"physicalPriceRestriction,omitempty"`
	CurrentPriceRestriction  *ManageOfferResultCurrentPriceRestriction  `json:"currentPriceRestriction,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ManageOfferResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ManageOfferResult
func (u ManageOfferResult) ArmForSwitch(sw int32) (string, bool) {
	switch ManageOfferResultCode(sw) {
	case ManageOfferResultCodeManageOfferSuccess:
		return "Success", true
	case ManageOfferResultCodeManageOfferPhysicalPriceRestriction:
		return "PhysicalPriceRestriction", true
	case ManageOfferResultCodeMaangeOfferCurrentPriceRestriction:
		return "CurrentPriceRestriction", true
	default:
		return "", true
	}
}

// NewManageOfferResult creates a new  ManageOfferResult.
func NewManageOfferResult(code ManageOfferResultCode, value interface{}) (result ManageOfferResult, err error) {
	result.Code = code
	switch ManageOfferResultCode(code) {
	case ManageOfferResultCodeManageOfferSuccess:
		tv, ok := value.(ManageOfferSuccessResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageOfferSuccessResult")
			return
		}
		result.Success = &tv
	case ManageOfferResultCodeManageOfferPhysicalPriceRestriction:
		tv, ok := value.(ManageOfferResultPhysicalPriceRestriction)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageOfferResultPhysicalPriceRestriction")
			return
		}
		result.PhysicalPriceRestriction = &tv
	case ManageOfferResultCodeMaangeOfferCurrentPriceRestriction:
		tv, ok := value.(ManageOfferResultCurrentPriceRestriction)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageOfferResultCurrentPriceRestriction")
			return
		}
		result.CurrentPriceRestriction = &tv
	default:
		// void
	}
	return
}

// MustSuccess retrieves the Success value from the union,
// panicing if the value is not set.
func (u ManageOfferResult) MustSuccess() ManageOfferSuccessResult {
	val, ok := u.GetSuccess()

	if !ok {
		panic("arm Success is not set")
	}

	return val
}

// GetSuccess retrieves the Success value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ManageOfferResult) GetSuccess() (result ManageOfferSuccessResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "Success" {
		result = *u.Success
		ok = true
	}

	return
}

// MustPhysicalPriceRestriction retrieves the PhysicalPriceRestriction value from the union,
// panicing if the value is not set.
func (u ManageOfferResult) MustPhysicalPriceRestriction() ManageOfferResultPhysicalPriceRestriction {
	val, ok := u.GetPhysicalPriceRestriction()

	if !ok {
		panic("arm PhysicalPriceRestriction is not set")
	}

	return val
}

// GetPhysicalPriceRestriction retrieves the PhysicalPriceRestriction value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ManageOfferResult) GetPhysicalPriceRestriction() (result ManageOfferResultPhysicalPriceRestriction, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "PhysicalPriceRestriction" {
		result = *u.PhysicalPriceRestriction
		ok = true
	}

	return
}

// MustCurrentPriceRestriction retrieves the CurrentPriceRestriction value from the union,
// panicing if the value is not set.
func (u ManageOfferResult) MustCurrentPriceRestriction() ManageOfferResultCurrentPriceRestriction {
	val, ok := u.GetCurrentPriceRestriction()

	if !ok {
		panic("arm CurrentPriceRestriction is not set")
	}

	return val
}

// GetCurrentPriceRestriction retrieves the CurrentPriceRestriction value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ManageOfferResult) GetCurrentPriceRestriction() (result ManageOfferResultCurrentPriceRestriction, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "CurrentPriceRestriction" {
		result = *u.CurrentPriceRestriction
		ok = true
	}

	return
}

// InvoiceReferenceExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type InvoiceReferenceExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u InvoiceReferenceExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of InvoiceReferenceExt
func (u InvoiceReferenceExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewInvoiceReferenceExt creates a new  InvoiceReferenceExt.
func NewInvoiceReferenceExt(v LedgerVersion, value interface{}) (result InvoiceReferenceExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// InvoiceReference is an XDR Struct defines as:
//
//   struct InvoiceReference {
//        uint64 invoiceID;
//        bool accept;
//    	// reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type InvoiceReference struct {
	InvoiceId Uint64              `json:"invoiceID,omitempty"`
	Accept    bool                `json:"accept,omitempty"`
	Ext       InvoiceReferenceExt `json:"ext,omitempty"`
}

// FeeDataExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type FeeDataExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u FeeDataExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of FeeDataExt
func (u FeeDataExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewFeeDataExt creates a new  FeeDataExt.
func NewFeeDataExt(v LedgerVersion, value interface{}) (result FeeDataExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// FeeData is an XDR Struct defines as:
//
//   struct FeeData {
//        int64 paymentFee;
//        int64 fixedFee;
//    	// reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type FeeData struct {
	PaymentFee Int64      `json:"paymentFee,omitempty"`
	FixedFee   Int64      `json:"fixedFee,omitempty"`
	Ext        FeeDataExt `json:"ext,omitempty"`
}

// PaymentFeeDataExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type PaymentFeeDataExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u PaymentFeeDataExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of PaymentFeeDataExt
func (u PaymentFeeDataExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewPaymentFeeDataExt creates a new  PaymentFeeDataExt.
func NewPaymentFeeDataExt(v LedgerVersion, value interface{}) (result PaymentFeeDataExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// PaymentFeeData is an XDR Struct defines as:
//
//   struct PaymentFeeData {
//        FeeData sourceFee;
//        FeeData destinationFee;
//        bool sourcePaysForDest;    // if true source account pays fee, else destination
//    	// reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type PaymentFeeData struct {
	SourceFee         FeeData           `json:"sourceFee,omitempty"`
	DestinationFee    FeeData           `json:"destinationFee,omitempty"`
	SourcePaysForDest bool              `json:"sourcePaysForDest,omitempty"`
	Ext               PaymentFeeDataExt `json:"ext,omitempty"`
}

// PaymentOpExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type PaymentOpExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u PaymentOpExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of PaymentOpExt
func (u PaymentOpExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewPaymentOpExt creates a new  PaymentOpExt.
func NewPaymentOpExt(v LedgerVersion, value interface{}) (result PaymentOpExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// PaymentOp is an XDR Struct defines as:
//
//   struct PaymentOp
//    {
//        BalanceID sourceBalanceID;
//        BalanceID destinationBalanceID;
//        int64 amount;          // amount they end up with
//
//        PaymentFeeData feeData;
//
//        string256 subject;
//        string64 reference;
//
//        InvoiceReference* invoiceReference;
//
//    	// reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type PaymentOp struct {
	SourceBalanceId      BalanceId         `json:"sourceBalanceID,omitempty"`
	DestinationBalanceId BalanceId         `json:"destinationBalanceID,omitempty"`
	Amount               Int64             `json:"amount,omitempty"`
	FeeData              PaymentFeeData    `json:"feeData,omitempty"`
	Subject              String256         `json:"subject,omitempty"`
	Reference            String64          `json:"reference,omitempty"`
	InvoiceReference     *InvoiceReference `json:"invoiceReference,omitempty"`
	Ext                  PaymentOpExt      `json:"ext,omitempty"`
}

// PaymentResultCode is an XDR Enum defines as:
//
//   enum PaymentResultCode
//    {
//        // codes considered as "success" for the operation
//        PAYMENT_SUCCESS = 0, // payment successfuly completed
//
//        // codes considered as "failure" for the operation
//        PAYMENT_MALFORMED = -1,       // bad input
//        PAYMENT_UNDERFUNDED = -2,     // not enough funds in source account
//        PAYMENT_LINE_FULL = -3,       // destination would go above their limit
//    	PAYMENT_FEE_MISMATCHED = -4,   // fee is not equal to expected fee
//        PAYMENT_BALANCE_NOT_FOUND = -5, // destination balance not found
//        PAYMENT_BALANCE_ACCOUNT_MISMATCHED = -6,
//        PAYMENT_BALANCE_ASSETS_MISMATCHED = -7,
//    	PAYMENT_SRC_BALANCE_NOT_FOUND = -8, // source balance not found
//        PAYMENT_REFERENCE_DUPLICATION = -9,
//        PAYMENT_STATS_OVERFLOW = -10,
//        PAYMENT_LIMITS_EXCEEDED = -11,
//        PAYMENT_NOT_ALLOWED_BY_ASSET_POLICY = -12,
//        PAYMENT_INVOICE_NOT_FOUND = -13,
//        PAYMENT_INVOICE_WRONG_AMOUNT = -14,
//        PAYMENT_INVOICE_BALANCE_MISMATCH = -15,
//        PAYMENT_INVOICE_ACCOUNT_MISMATCH = -16,
//        PAYMENT_INVOICE_ALREADY_PAID = -17
//    };
//
type PaymentResultCode int32

const (
	PaymentResultCodePaymentSuccess                  PaymentResultCode = 0
	PaymentResultCodePaymentMalformed                PaymentResultCode = -1
	PaymentResultCodePaymentUnderfunded              PaymentResultCode = -2
	PaymentResultCodePaymentLineFull                 PaymentResultCode = -3
	PaymentResultCodePaymentFeeMismatched            PaymentResultCode = -4
	PaymentResultCodePaymentBalanceNotFound          PaymentResultCode = -5
	PaymentResultCodePaymentBalanceAccountMismatched PaymentResultCode = -6
	PaymentResultCodePaymentBalanceAssetsMismatched  PaymentResultCode = -7
	PaymentResultCodePaymentSrcBalanceNotFound       PaymentResultCode = -8
	PaymentResultCodePaymentReferenceDuplication     PaymentResultCode = -9
	PaymentResultCodePaymentStatsOverflow            PaymentResultCode = -10
	PaymentResultCodePaymentLimitsExceeded           PaymentResultCode = -11
	PaymentResultCodePaymentNotAllowedByAssetPolicy  PaymentResultCode = -12
	PaymentResultCodePaymentInvoiceNotFound          PaymentResultCode = -13
	PaymentResultCodePaymentInvoiceWrongAmount       PaymentResultCode = -14
	PaymentResultCodePaymentInvoiceBalanceMismatch   PaymentResultCode = -15
	PaymentResultCodePaymentInvoiceAccountMismatch   PaymentResultCode = -16
	PaymentResultCodePaymentInvoiceAlreadyPaid       PaymentResultCode = -17
)

var PaymentResultCodeAll = []PaymentResultCode{
	PaymentResultCodePaymentSuccess,
	PaymentResultCodePaymentMalformed,
	PaymentResultCodePaymentUnderfunded,
	PaymentResultCodePaymentLineFull,
	PaymentResultCodePaymentFeeMismatched,
	PaymentResultCodePaymentBalanceNotFound,
	PaymentResultCodePaymentBalanceAccountMismatched,
	PaymentResultCodePaymentBalanceAssetsMismatched,
	PaymentResultCodePaymentSrcBalanceNotFound,
	PaymentResultCodePaymentReferenceDuplication,
	PaymentResultCodePaymentStatsOverflow,
	PaymentResultCodePaymentLimitsExceeded,
	PaymentResultCodePaymentNotAllowedByAssetPolicy,
	PaymentResultCodePaymentInvoiceNotFound,
	PaymentResultCodePaymentInvoiceWrongAmount,
	PaymentResultCodePaymentInvoiceBalanceMismatch,
	PaymentResultCodePaymentInvoiceAccountMismatch,
	PaymentResultCodePaymentInvoiceAlreadyPaid,
}

var paymentResultCodeMap = map[int32]string{
	0:   "PaymentResultCodePaymentSuccess",
	-1:  "PaymentResultCodePaymentMalformed",
	-2:  "PaymentResultCodePaymentUnderfunded",
	-3:  "PaymentResultCodePaymentLineFull",
	-4:  "PaymentResultCodePaymentFeeMismatched",
	-5:  "PaymentResultCodePaymentBalanceNotFound",
	-6:  "PaymentResultCodePaymentBalanceAccountMismatched",
	-7:  "PaymentResultCodePaymentBalanceAssetsMismatched",
	-8:  "PaymentResultCodePaymentSrcBalanceNotFound",
	-9:  "PaymentResultCodePaymentReferenceDuplication",
	-10: "PaymentResultCodePaymentStatsOverflow",
	-11: "PaymentResultCodePaymentLimitsExceeded",
	-12: "PaymentResultCodePaymentNotAllowedByAssetPolicy",
	-13: "PaymentResultCodePaymentInvoiceNotFound",
	-14: "PaymentResultCodePaymentInvoiceWrongAmount",
	-15: "PaymentResultCodePaymentInvoiceBalanceMismatch",
	-16: "PaymentResultCodePaymentInvoiceAccountMismatch",
	-17: "PaymentResultCodePaymentInvoiceAlreadyPaid",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for PaymentResultCode
func (e PaymentResultCode) ValidEnum(v int32) bool {
	_, ok := paymentResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e PaymentResultCode) String() string {
	name, _ := paymentResultCodeMap[int32(e)]
	return name
}

func (e PaymentResultCode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// PaymentResponseExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type PaymentResponseExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u PaymentResponseExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of PaymentResponseExt
func (u PaymentResponseExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewPaymentResponseExt creates a new  PaymentResponseExt.
func NewPaymentResponseExt(v LedgerVersion, value interface{}) (result PaymentResponseExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// PaymentResponse is an XDR Struct defines as:
//
//   struct PaymentResponse {
//        AccountID exchanges<>;
//        AccountID destination;
//        uint64 paymentID;
//        AssetCode asset;
//        // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type PaymentResponse struct {
	Exchanges   []AccountId        `json:"exchanges,omitempty"`
	Destination AccountId          `json:"destination,omitempty"`
	PaymentId   Uint64             `json:"paymentID,omitempty"`
	Asset       AssetCode          `json:"asset,omitempty"`
	Ext         PaymentResponseExt `json:"ext,omitempty"`
}

// PaymentResult is an XDR Union defines as:
//
//   union PaymentResult switch (PaymentResultCode code)
//    {
//    case PAYMENT_SUCCESS:
//        PaymentResponse paymentResponse;
//    default:
//        void;
//    };
//
type PaymentResult struct {
	Code            PaymentResultCode `json:"code,omitempty"`
	PaymentResponse *PaymentResponse  `json:"paymentResponse,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u PaymentResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of PaymentResult
func (u PaymentResult) ArmForSwitch(sw int32) (string, bool) {
	switch PaymentResultCode(sw) {
	case PaymentResultCodePaymentSuccess:
		return "PaymentResponse", true
	default:
		return "", true
	}
}

// NewPaymentResult creates a new  PaymentResult.
func NewPaymentResult(code PaymentResultCode, value interface{}) (result PaymentResult, err error) {
	result.Code = code
	switch PaymentResultCode(code) {
	case PaymentResultCodePaymentSuccess:
		tv, ok := value.(PaymentResponse)
		if !ok {
			err = fmt.Errorf("invalid value, must be PaymentResponse")
			return
		}
		result.PaymentResponse = &tv
	default:
		// void
	}
	return
}

// MustPaymentResponse retrieves the PaymentResponse value from the union,
// panicing if the value is not set.
func (u PaymentResult) MustPaymentResponse() PaymentResponse {
	val, ok := u.GetPaymentResponse()

	if !ok {
		panic("arm PaymentResponse is not set")
	}

	return val
}

// GetPaymentResponse retrieves the PaymentResponse value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u PaymentResult) GetPaymentResponse() (result PaymentResponse, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "PaymentResponse" {
		result = *u.PaymentResponse
		ok = true
	}

	return
}

// RecoverOpExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type RecoverOpExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u RecoverOpExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of RecoverOpExt
func (u RecoverOpExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewRecoverOpExt creates a new  RecoverOpExt.
func NewRecoverOpExt(v LedgerVersion, value interface{}) (result RecoverOpExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// RecoverOp is an XDR Struct defines as:
//
//   struct RecoverOp
//    {
//        AccountID account;
//        PublicKey oldSigner;
//        PublicKey newSigner;
//    	// reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type RecoverOp struct {
	Account   AccountId    `json:"account,omitempty"`
	OldSigner PublicKey    `json:"oldSigner,omitempty"`
	NewSigner PublicKey    `json:"newSigner,omitempty"`
	Ext       RecoverOpExt `json:"ext,omitempty"`
}

// RecoverResultCode is an XDR Enum defines as:
//
//   enum RecoverResultCode
//    {
//        // codes considered as "success" for the operation
//        RECOVER_SUCCESS = 0,
//
//        // codes considered as "failure" for the operation
//
//        RECOVER_MALFORMED = -1,
//        RECOVER_OLD_SIGNER_NOT_FOUND = -2,
//        RECOVER_SIGNER_ALREADY_EXISTS = -3
//    };
//
type RecoverResultCode int32

const (
	RecoverResultCodeRecoverSuccess             RecoverResultCode = 0
	RecoverResultCodeRecoverMalformed           RecoverResultCode = -1
	RecoverResultCodeRecoverOldSignerNotFound   RecoverResultCode = -2
	RecoverResultCodeRecoverSignerAlreadyExists RecoverResultCode = -3
)

var RecoverResultCodeAll = []RecoverResultCode{
	RecoverResultCodeRecoverSuccess,
	RecoverResultCodeRecoverMalformed,
	RecoverResultCodeRecoverOldSignerNotFound,
	RecoverResultCodeRecoverSignerAlreadyExists,
}

var recoverResultCodeMap = map[int32]string{
	0:  "RecoverResultCodeRecoverSuccess",
	-1: "RecoverResultCodeRecoverMalformed",
	-2: "RecoverResultCodeRecoverOldSignerNotFound",
	-3: "RecoverResultCodeRecoverSignerAlreadyExists",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for RecoverResultCode
func (e RecoverResultCode) ValidEnum(v int32) bool {
	_, ok := recoverResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e RecoverResultCode) String() string {
	name, _ := recoverResultCodeMap[int32(e)]
	return name
}

func (e RecoverResultCode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// RecoverResultSuccessExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type RecoverResultSuccessExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u RecoverResultSuccessExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of RecoverResultSuccessExt
func (u RecoverResultSuccessExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewRecoverResultSuccessExt creates a new  RecoverResultSuccessExt.
func NewRecoverResultSuccessExt(v LedgerVersion, value interface{}) (result RecoverResultSuccessExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// RecoverResultSuccess is an XDR NestedStruct defines as:
//
//   struct {
//    		// reserved for future use
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//    	}
//
type RecoverResultSuccess struct {
	Ext RecoverResultSuccessExt `json:"ext,omitempty"`
}

// RecoverResult is an XDR Union defines as:
//
//   union RecoverResult switch (RecoverResultCode code)
//    {
//    case RECOVER_SUCCESS:
//        struct {
//    		// reserved for future use
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//    	} success;
//    default:
//        void;
//    };
//
type RecoverResult struct {
	Code    RecoverResultCode     `json:"code,omitempty"`
	Success *RecoverResultSuccess `json:"success,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u RecoverResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of RecoverResult
func (u RecoverResult) ArmForSwitch(sw int32) (string, bool) {
	switch RecoverResultCode(sw) {
	case RecoverResultCodeRecoverSuccess:
		return "Success", true
	default:
		return "", true
	}
}

// NewRecoverResult creates a new  RecoverResult.
func NewRecoverResult(code RecoverResultCode, value interface{}) (result RecoverResult, err error) {
	result.Code = code
	switch RecoverResultCode(code) {
	case RecoverResultCodeRecoverSuccess:
		tv, ok := value.(RecoverResultSuccess)
		if !ok {
			err = fmt.Errorf("invalid value, must be RecoverResultSuccess")
			return
		}
		result.Success = &tv
	default:
		// void
	}
	return
}

// MustSuccess retrieves the Success value from the union,
// panicing if the value is not set.
func (u RecoverResult) MustSuccess() RecoverResultSuccess {
	val, ok := u.GetSuccess()

	if !ok {
		panic("arm Success is not set")
	}

	return val
}

// GetSuccess retrieves the Success value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u RecoverResult) GetSuccess() (result RecoverResultSuccess, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "Success" {
		result = *u.Success
		ok = true
	}

	return
}

// ReviewCoinsEmissionRequestOpExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type ReviewCoinsEmissionRequestOpExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ReviewCoinsEmissionRequestOpExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ReviewCoinsEmissionRequestOpExt
func (u ReviewCoinsEmissionRequestOpExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewReviewCoinsEmissionRequestOpExt creates a new  ReviewCoinsEmissionRequestOpExt.
func NewReviewCoinsEmissionRequestOpExt(v LedgerVersion, value interface{}) (result ReviewCoinsEmissionRequestOpExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ReviewCoinsEmissionRequestOp is an XDR Struct defines as:
//
//   struct ReviewCoinsEmissionRequestOp
//    {
//    	CoinsEmissionRequestEntry request;  // request to be reviewed
//    	bool approve;
//    	string64 reason;
//    	// reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type ReviewCoinsEmissionRequestOp struct {
	Request CoinsEmissionRequestEntry       `json:"request,omitempty"`
	Approve bool                            `json:"approve,omitempty"`
	Reason  String64                        `json:"reason,omitempty"`
	Ext     ReviewCoinsEmissionRequestOpExt `json:"ext,omitempty"`
}

// ReviewCoinsEmissionRequestResultCode is an XDR Enum defines as:
//
//   enum ReviewCoinsEmissionRequestResultCode
//    {
//        // codes considered as "success" for the operation
//        REVIEW_COINS_EMISSION_REQUEST_SUCCESS = 0,
//
//        // codes considered as "failure" for the operation
//        REVIEW_COINS_EMISSION_REQUEST_INVALID_REASON = -1,        // reason must be null if approving
//    	REVIEW_COINS_EMISSION_REQUEST_NOT_FOUND = -2,             // failed to find emission request with such ID
//    	REVIEW_COINS_EMISSION_REQUEST_NOT_EQUAL = -3,             // stored emission request is not equal to request provided in op
//    	REVIEW_COINS_EMISSION_REQUEST_ALREADY_REVIEWED = -4,      // emission request have been already reviewed
//    	REVIEW_COINS_EMISSION_REQUEST_MALFORMED = -5,             // emission request is malformed
//        REVIEW_COINS_EMISSION_REQUEST_NOT_ENOUGH_PREEMISSIONS = -6,    // serial is already used in another review
//    	REVIEW_COINS_EMISSION_REQUEST_LINE_FULL = -9,             // balance will overflow
//        REVIEW_COINS_EMISSION_REQUEST_ASSET_NOT_FOUND = -10,
//        REVIEW_COINS_EMISSION_REQUEST_BALANCE_NOT_FOUND = -11,
//    	REVIEW_COINS_EMISSION_REQUEST_REFERENCE_DUPLICATION = -12
//    };
//
type ReviewCoinsEmissionRequestResultCode int32

const (
	ReviewCoinsEmissionRequestResultCodeReviewCoinsEmissionRequestSuccess               ReviewCoinsEmissionRequestResultCode = 0
	ReviewCoinsEmissionRequestResultCodeReviewCoinsEmissionRequestInvalidReason         ReviewCoinsEmissionRequestResultCode = -1
	ReviewCoinsEmissionRequestResultCodeReviewCoinsEmissionRequestNotFound              ReviewCoinsEmissionRequestResultCode = -2
	ReviewCoinsEmissionRequestResultCodeReviewCoinsEmissionRequestNotEqual              ReviewCoinsEmissionRequestResultCode = -3
	ReviewCoinsEmissionRequestResultCodeReviewCoinsEmissionRequestAlreadyReviewed       ReviewCoinsEmissionRequestResultCode = -4
	ReviewCoinsEmissionRequestResultCodeReviewCoinsEmissionRequestMalformed             ReviewCoinsEmissionRequestResultCode = -5
	ReviewCoinsEmissionRequestResultCodeReviewCoinsEmissionRequestNotEnoughPreemissions ReviewCoinsEmissionRequestResultCode = -6
	ReviewCoinsEmissionRequestResultCodeReviewCoinsEmissionRequestLineFull              ReviewCoinsEmissionRequestResultCode = -9
	ReviewCoinsEmissionRequestResultCodeReviewCoinsEmissionRequestAssetNotFound         ReviewCoinsEmissionRequestResultCode = -10
	ReviewCoinsEmissionRequestResultCodeReviewCoinsEmissionRequestBalanceNotFound       ReviewCoinsEmissionRequestResultCode = -11
	ReviewCoinsEmissionRequestResultCodeReviewCoinsEmissionRequestReferenceDuplication  ReviewCoinsEmissionRequestResultCode = -12
)

var ReviewCoinsEmissionRequestResultCodeAll = []ReviewCoinsEmissionRequestResultCode{
	ReviewCoinsEmissionRequestResultCodeReviewCoinsEmissionRequestSuccess,
	ReviewCoinsEmissionRequestResultCodeReviewCoinsEmissionRequestInvalidReason,
	ReviewCoinsEmissionRequestResultCodeReviewCoinsEmissionRequestNotFound,
	ReviewCoinsEmissionRequestResultCodeReviewCoinsEmissionRequestNotEqual,
	ReviewCoinsEmissionRequestResultCodeReviewCoinsEmissionRequestAlreadyReviewed,
	ReviewCoinsEmissionRequestResultCodeReviewCoinsEmissionRequestMalformed,
	ReviewCoinsEmissionRequestResultCodeReviewCoinsEmissionRequestNotEnoughPreemissions,
	ReviewCoinsEmissionRequestResultCodeReviewCoinsEmissionRequestLineFull,
	ReviewCoinsEmissionRequestResultCodeReviewCoinsEmissionRequestAssetNotFound,
	ReviewCoinsEmissionRequestResultCodeReviewCoinsEmissionRequestBalanceNotFound,
	ReviewCoinsEmissionRequestResultCodeReviewCoinsEmissionRequestReferenceDuplication,
}

var reviewCoinsEmissionRequestResultCodeMap = map[int32]string{
	0:   "ReviewCoinsEmissionRequestResultCodeReviewCoinsEmissionRequestSuccess",
	-1:  "ReviewCoinsEmissionRequestResultCodeReviewCoinsEmissionRequestInvalidReason",
	-2:  "ReviewCoinsEmissionRequestResultCodeReviewCoinsEmissionRequestNotFound",
	-3:  "ReviewCoinsEmissionRequestResultCodeReviewCoinsEmissionRequestNotEqual",
	-4:  "ReviewCoinsEmissionRequestResultCodeReviewCoinsEmissionRequestAlreadyReviewed",
	-5:  "ReviewCoinsEmissionRequestResultCodeReviewCoinsEmissionRequestMalformed",
	-6:  "ReviewCoinsEmissionRequestResultCodeReviewCoinsEmissionRequestNotEnoughPreemissions",
	-9:  "ReviewCoinsEmissionRequestResultCodeReviewCoinsEmissionRequestLineFull",
	-10: "ReviewCoinsEmissionRequestResultCodeReviewCoinsEmissionRequestAssetNotFound",
	-11: "ReviewCoinsEmissionRequestResultCodeReviewCoinsEmissionRequestBalanceNotFound",
	-12: "ReviewCoinsEmissionRequestResultCodeReviewCoinsEmissionRequestReferenceDuplication",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ReviewCoinsEmissionRequestResultCode
func (e ReviewCoinsEmissionRequestResultCode) ValidEnum(v int32) bool {
	_, ok := reviewCoinsEmissionRequestResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e ReviewCoinsEmissionRequestResultCode) String() string {
	name, _ := reviewCoinsEmissionRequestResultCodeMap[int32(e)]
	return name
}

func (e ReviewCoinsEmissionRequestResultCode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// ReviewCoinsEmissionRequestResultSuccessExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type ReviewCoinsEmissionRequestResultSuccessExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ReviewCoinsEmissionRequestResultSuccessExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ReviewCoinsEmissionRequestResultSuccessExt
func (u ReviewCoinsEmissionRequestResultSuccessExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewReviewCoinsEmissionRequestResultSuccessExt creates a new  ReviewCoinsEmissionRequestResultSuccessExt.
func NewReviewCoinsEmissionRequestResultSuccessExt(v LedgerVersion, value interface{}) (result ReviewCoinsEmissionRequestResultSuccessExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ReviewCoinsEmissionRequestResultSuccess is an XDR NestedStruct defines as:
//
//   struct {
//    		uint64 requestID;
//    		// reserved for future use
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//    	}
//
type ReviewCoinsEmissionRequestResultSuccess struct {
	RequestId Uint64                                     `json:"requestID,omitempty"`
	Ext       ReviewCoinsEmissionRequestResultSuccessExt `json:"ext,omitempty"`
}

// ReviewCoinsEmissionRequestResult is an XDR Union defines as:
//
//   union ReviewCoinsEmissionRequestResult switch (ReviewCoinsEmissionRequestResultCode code)
//    {
//    case REVIEW_COINS_EMISSION_REQUEST_SUCCESS:
//    	struct {
//    		uint64 requestID;
//    		// reserved for future use
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//    	} success;
//    default:
//        void;
//    };
//
type ReviewCoinsEmissionRequestResult struct {
	Code    ReviewCoinsEmissionRequestResultCode     `json:"code,omitempty"`
	Success *ReviewCoinsEmissionRequestResultSuccess `json:"success,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ReviewCoinsEmissionRequestResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ReviewCoinsEmissionRequestResult
func (u ReviewCoinsEmissionRequestResult) ArmForSwitch(sw int32) (string, bool) {
	switch ReviewCoinsEmissionRequestResultCode(sw) {
	case ReviewCoinsEmissionRequestResultCodeReviewCoinsEmissionRequestSuccess:
		return "Success", true
	default:
		return "", true
	}
}

// NewReviewCoinsEmissionRequestResult creates a new  ReviewCoinsEmissionRequestResult.
func NewReviewCoinsEmissionRequestResult(code ReviewCoinsEmissionRequestResultCode, value interface{}) (result ReviewCoinsEmissionRequestResult, err error) {
	result.Code = code
	switch ReviewCoinsEmissionRequestResultCode(code) {
	case ReviewCoinsEmissionRequestResultCodeReviewCoinsEmissionRequestSuccess:
		tv, ok := value.(ReviewCoinsEmissionRequestResultSuccess)
		if !ok {
			err = fmt.Errorf("invalid value, must be ReviewCoinsEmissionRequestResultSuccess")
			return
		}
		result.Success = &tv
	default:
		// void
	}
	return
}

// MustSuccess retrieves the Success value from the union,
// panicing if the value is not set.
func (u ReviewCoinsEmissionRequestResult) MustSuccess() ReviewCoinsEmissionRequestResultSuccess {
	val, ok := u.GetSuccess()

	if !ok {
		panic("arm Success is not set")
	}

	return val
}

// GetSuccess retrieves the Success value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ReviewCoinsEmissionRequestResult) GetSuccess() (result ReviewCoinsEmissionRequestResultSuccess, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "Success" {
		result = *u.Success
		ok = true
	}

	return
}

// ReviewPaymentRequestOpExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    	{
//    	case EMPTY_VERSION:
//    		void;
//    	}
//
type ReviewPaymentRequestOpExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ReviewPaymentRequestOpExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ReviewPaymentRequestOpExt
func (u ReviewPaymentRequestOpExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewReviewPaymentRequestOpExt creates a new  ReviewPaymentRequestOpExt.
func NewReviewPaymentRequestOpExt(v LedgerVersion, value interface{}) (result ReviewPaymentRequestOpExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ReviewPaymentRequestOp is an XDR Struct defines as:
//
//   struct ReviewPaymentRequestOp
//    {
//        uint64 paymentID;
//
//    	bool accept;
//        string256* rejectReason;
//    	// reserved for future use
//    	union switch (LedgerVersion v)
//    	{
//    	case EMPTY_VERSION:
//    		void;
//    	}
//    	ext;
//    };
//
type ReviewPaymentRequestOp struct {
	PaymentId    Uint64                    `json:"paymentID,omitempty"`
	Accept       bool                      `json:"accept,omitempty"`
	RejectReason *String256                `json:"rejectReason,omitempty"`
	Ext          ReviewPaymentRequestOpExt `json:"ext,omitempty"`
}

// ReviewPaymentRequestResultCode is an XDR Enum defines as:
//
//   enum ReviewPaymentRequestResultCode
//    {
//        // codes considered as "success" for the operation
//        REVIEW_PAYMENT_REQUEST_SUCCESS = 0,
//
//        // codes considered as "failure" for the operation
//    	REVIEW_PAYMENT_REQUEST_NOT_FOUND = -1,           // failed to find Recovery request with such ID
//        REVIEW_PAYMENT_REQUEST_LINE_FULL = -2,
//        REVIEW_PAYMENT_DEMURRAGE_REJECTION_NOT_ALLOWED = -3
//
//    };
//
type ReviewPaymentRequestResultCode int32

const (
	ReviewPaymentRequestResultCodeReviewPaymentRequestSuccess               ReviewPaymentRequestResultCode = 0
	ReviewPaymentRequestResultCodeReviewPaymentRequestNotFound              ReviewPaymentRequestResultCode = -1
	ReviewPaymentRequestResultCodeReviewPaymentRequestLineFull              ReviewPaymentRequestResultCode = -2
	ReviewPaymentRequestResultCodeReviewPaymentDemurrageRejectionNotAllowed ReviewPaymentRequestResultCode = -3
)

var ReviewPaymentRequestResultCodeAll = []ReviewPaymentRequestResultCode{
	ReviewPaymentRequestResultCodeReviewPaymentRequestSuccess,
	ReviewPaymentRequestResultCodeReviewPaymentRequestNotFound,
	ReviewPaymentRequestResultCodeReviewPaymentRequestLineFull,
	ReviewPaymentRequestResultCodeReviewPaymentDemurrageRejectionNotAllowed,
}

var reviewPaymentRequestResultCodeMap = map[int32]string{
	0:  "ReviewPaymentRequestResultCodeReviewPaymentRequestSuccess",
	-1: "ReviewPaymentRequestResultCodeReviewPaymentRequestNotFound",
	-2: "ReviewPaymentRequestResultCodeReviewPaymentRequestLineFull",
	-3: "ReviewPaymentRequestResultCodeReviewPaymentDemurrageRejectionNotAllowed",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ReviewPaymentRequestResultCode
func (e ReviewPaymentRequestResultCode) ValidEnum(v int32) bool {
	_, ok := reviewPaymentRequestResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e ReviewPaymentRequestResultCode) String() string {
	name, _ := reviewPaymentRequestResultCodeMap[int32(e)]
	return name
}

func (e ReviewPaymentRequestResultCode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// PaymentState is an XDR Enum defines as:
//
//   enum PaymentState
//    {
//        PAYMENT_PENDING = 0,
//        PAYMENT_PROCESSED = 1,
//        PAYMENT_REJECTED = 2
//    };
//
type PaymentState int32

const (
	PaymentStatePaymentPending   PaymentState = 0
	PaymentStatePaymentProcessed PaymentState = 1
	PaymentStatePaymentRejected  PaymentState = 2
)

var PaymentStateAll = []PaymentState{
	PaymentStatePaymentPending,
	PaymentStatePaymentProcessed,
	PaymentStatePaymentRejected,
}

var paymentStateMap = map[int32]string{
	0: "PaymentStatePaymentPending",
	1: "PaymentStatePaymentProcessed",
	2: "PaymentStatePaymentRejected",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for PaymentState
func (e PaymentState) ValidEnum(v int32) bool {
	_, ok := paymentStateMap[v]
	return ok
}

// String returns the name of `e`
func (e PaymentState) String() string {
	name, _ := paymentStateMap[int32(e)]
	return name
}

func (e PaymentState) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// ReviewPaymentResponseExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    	{
//    	case EMPTY_VERSION:
//    		void;
//    	}
//
type ReviewPaymentResponseExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ReviewPaymentResponseExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ReviewPaymentResponseExt
func (u ReviewPaymentResponseExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewReviewPaymentResponseExt creates a new  ReviewPaymentResponseExt.
func NewReviewPaymentResponseExt(v LedgerVersion, value interface{}) (result ReviewPaymentResponseExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// ReviewPaymentResponse is an XDR Struct defines as:
//
//   struct ReviewPaymentResponse {
//        PaymentState state;
//
//        uint64* relatedInvoiceID;
//    	// reserved for future use
//    	union switch (LedgerVersion v)
//    	{
//    	case EMPTY_VERSION:
//    		void;
//    	}
//    	ext;
//    };
//
type ReviewPaymentResponse struct {
	State            PaymentState             `json:"state,omitempty"`
	RelatedInvoiceId *Uint64                  `json:"relatedInvoiceID,omitempty"`
	Ext              ReviewPaymentResponseExt `json:"ext,omitempty"`
}

// ReviewPaymentRequestResult is an XDR Union defines as:
//
//   union ReviewPaymentRequestResult switch (ReviewPaymentRequestResultCode code)
//    {
//    case REVIEW_PAYMENT_REQUEST_SUCCESS:
//        ReviewPaymentResponse reviewPaymentResponse;
//    default:
//        void;
//    };
//
type ReviewPaymentRequestResult struct {
	Code                  ReviewPaymentRequestResultCode `json:"code,omitempty"`
	ReviewPaymentResponse *ReviewPaymentResponse         `json:"reviewPaymentResponse,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ReviewPaymentRequestResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ReviewPaymentRequestResult
func (u ReviewPaymentRequestResult) ArmForSwitch(sw int32) (string, bool) {
	switch ReviewPaymentRequestResultCode(sw) {
	case ReviewPaymentRequestResultCodeReviewPaymentRequestSuccess:
		return "ReviewPaymentResponse", true
	default:
		return "", true
	}
}

// NewReviewPaymentRequestResult creates a new  ReviewPaymentRequestResult.
func NewReviewPaymentRequestResult(code ReviewPaymentRequestResultCode, value interface{}) (result ReviewPaymentRequestResult, err error) {
	result.Code = code
	switch ReviewPaymentRequestResultCode(code) {
	case ReviewPaymentRequestResultCodeReviewPaymentRequestSuccess:
		tv, ok := value.(ReviewPaymentResponse)
		if !ok {
			err = fmt.Errorf("invalid value, must be ReviewPaymentResponse")
			return
		}
		result.ReviewPaymentResponse = &tv
	default:
		// void
	}
	return
}

// MustReviewPaymentResponse retrieves the ReviewPaymentResponse value from the union,
// panicing if the value is not set.
func (u ReviewPaymentRequestResult) MustReviewPaymentResponse() ReviewPaymentResponse {
	val, ok := u.GetReviewPaymentResponse()

	if !ok {
		panic("arm ReviewPaymentResponse is not set")
	}

	return val
}

// GetReviewPaymentResponse retrieves the ReviewPaymentResponse value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ReviewPaymentRequestResult) GetReviewPaymentResponse() (result ReviewPaymentResponse, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "ReviewPaymentResponse" {
		result = *u.ReviewPaymentResponse
		ok = true
	}

	return
}

// SetFeesOpExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type SetFeesOpExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u SetFeesOpExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of SetFeesOpExt
func (u SetFeesOpExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewSetFeesOpExt creates a new  SetFeesOpExt.
func NewSetFeesOpExt(v LedgerVersion, value interface{}) (result SetFeesOpExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// SetFeesOp is an XDR Struct defines as:
//
//   struct SetFeesOp
//        {
//            FeeEntry* fee;
//    		bool isDelete;
//    		int64* storageFeePeriod;
//    		int64* payoutsPeriod;
//    		// reserved for future use
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//        };
//
type SetFeesOp struct {
	Fee              *FeeEntry    `json:"fee,omitempty"`
	IsDelete         bool         `json:"isDelete,omitempty"`
	StorageFeePeriod *Int64       `json:"storageFeePeriod,omitempty"`
	PayoutsPeriod    *Int64       `json:"payoutsPeriod,omitempty"`
	Ext              SetFeesOpExt `json:"ext,omitempty"`
}

// SetFeesResultCode is an XDR Enum defines as:
//
//   enum SetFeesResultCode
//        {
//            // codes considered as "success" for the operation
//            SET_FEES_SUCCESS = 0,
//
//            // codes considered as "failure" for the operation
//            SET_FEES_INVALID_AMOUNT = -1,      // amount is negative
//    		SET_FEES_INVALID_FEE_TYPE = -2,     // operation type is invalid
//            SET_FEES_ASSET_NOT_FOUND = -3,
//            SET_FEES_INVALID_ASSET = -4,
//            SET_FEES_MALFORMED = -5,
//    		SET_FEES_MALFORMED_RANGE = -6,
//    		SET_FEES_RANGE_OVERLAP = -7,
//    		SET_FEES_NOT_FOUND = -8,
//    		SET_FEES_SUB_TYPE_NOT_EXIST = -9
//        };
//
type SetFeesResultCode int32

const (
	SetFeesResultCodeSetFeesSuccess         SetFeesResultCode = 0
	SetFeesResultCodeSetFeesInvalidAmount   SetFeesResultCode = -1
	SetFeesResultCodeSetFeesInvalidFeeType  SetFeesResultCode = -2
	SetFeesResultCodeSetFeesAssetNotFound   SetFeesResultCode = -3
	SetFeesResultCodeSetFeesInvalidAsset    SetFeesResultCode = -4
	SetFeesResultCodeSetFeesMalformed       SetFeesResultCode = -5
	SetFeesResultCodeSetFeesMalformedRange  SetFeesResultCode = -6
	SetFeesResultCodeSetFeesRangeOverlap    SetFeesResultCode = -7
	SetFeesResultCodeSetFeesNotFound        SetFeesResultCode = -8
	SetFeesResultCodeSetFeesSubTypeNotExist SetFeesResultCode = -9
)

var SetFeesResultCodeAll = []SetFeesResultCode{
	SetFeesResultCodeSetFeesSuccess,
	SetFeesResultCodeSetFeesInvalidAmount,
	SetFeesResultCodeSetFeesInvalidFeeType,
	SetFeesResultCodeSetFeesAssetNotFound,
	SetFeesResultCodeSetFeesInvalidAsset,
	SetFeesResultCodeSetFeesMalformed,
	SetFeesResultCodeSetFeesMalformedRange,
	SetFeesResultCodeSetFeesRangeOverlap,
	SetFeesResultCodeSetFeesNotFound,
	SetFeesResultCodeSetFeesSubTypeNotExist,
}

var setFeesResultCodeMap = map[int32]string{
	0:  "SetFeesResultCodeSetFeesSuccess",
	-1: "SetFeesResultCodeSetFeesInvalidAmount",
	-2: "SetFeesResultCodeSetFeesInvalidFeeType",
	-3: "SetFeesResultCodeSetFeesAssetNotFound",
	-4: "SetFeesResultCodeSetFeesInvalidAsset",
	-5: "SetFeesResultCodeSetFeesMalformed",
	-6: "SetFeesResultCodeSetFeesMalformedRange",
	-7: "SetFeesResultCodeSetFeesRangeOverlap",
	-8: "SetFeesResultCodeSetFeesNotFound",
	-9: "SetFeesResultCodeSetFeesSubTypeNotExist",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for SetFeesResultCode
func (e SetFeesResultCode) ValidEnum(v int32) bool {
	_, ok := setFeesResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e SetFeesResultCode) String() string {
	name, _ := setFeesResultCodeMap[int32(e)]
	return name
}

func (e SetFeesResultCode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// SetFeesResultSuccessExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    				{
//    				case EMPTY_VERSION:
//    					void;
//    				}
//
type SetFeesResultSuccessExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u SetFeesResultSuccessExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of SetFeesResultSuccessExt
func (u SetFeesResultSuccessExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewSetFeesResultSuccessExt creates a new  SetFeesResultSuccessExt.
func NewSetFeesResultSuccessExt(v LedgerVersion, value interface{}) (result SetFeesResultSuccessExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// SetFeesResultSuccess is an XDR NestedStruct defines as:
//
//   struct {
//    				// reserved for future use
//    				union switch (LedgerVersion v)
//    				{
//    				case EMPTY_VERSION:
//    					void;
//    				}
//    				ext;
//    			}
//
type SetFeesResultSuccess struct {
	Ext SetFeesResultSuccessExt `json:"ext,omitempty"`
}

// SetFeesResult is an XDR Union defines as:
//
//   union SetFeesResult switch (SetFeesResultCode code)
//        {
//            case SET_FEES_SUCCESS:
//                struct {
//    				// reserved for future use
//    				union switch (LedgerVersion v)
//    				{
//    				case EMPTY_VERSION:
//    					void;
//    				}
//    				ext;
//    			} success;
//            default:
//                void;
//        };
//
type SetFeesResult struct {
	Code    SetFeesResultCode     `json:"code,omitempty"`
	Success *SetFeesResultSuccess `json:"success,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u SetFeesResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of SetFeesResult
func (u SetFeesResult) ArmForSwitch(sw int32) (string, bool) {
	switch SetFeesResultCode(sw) {
	case SetFeesResultCodeSetFeesSuccess:
		return "Success", true
	default:
		return "", true
	}
}

// NewSetFeesResult creates a new  SetFeesResult.
func NewSetFeesResult(code SetFeesResultCode, value interface{}) (result SetFeesResult, err error) {
	result.Code = code
	switch SetFeesResultCode(code) {
	case SetFeesResultCodeSetFeesSuccess:
		tv, ok := value.(SetFeesResultSuccess)
		if !ok {
			err = fmt.Errorf("invalid value, must be SetFeesResultSuccess")
			return
		}
		result.Success = &tv
	default:
		// void
	}
	return
}

// MustSuccess retrieves the Success value from the union,
// panicing if the value is not set.
func (u SetFeesResult) MustSuccess() SetFeesResultSuccess {
	val, ok := u.GetSuccess()

	if !ok {
		panic("arm Success is not set")
	}

	return val
}

// GetSuccess retrieves the Success value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u SetFeesResult) GetSuccess() (result SetFeesResultSuccess, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "Success" {
		result = *u.Success
		ok = true
	}

	return
}

// SetLimitsOpExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    	{
//    	case EMPTY_VERSION:
//    		void;
//    	}
//
type SetLimitsOpExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u SetLimitsOpExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of SetLimitsOpExt
func (u SetLimitsOpExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewSetLimitsOpExt creates a new  SetLimitsOpExt.
func NewSetLimitsOpExt(v LedgerVersion, value interface{}) (result SetLimitsOpExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// SetLimitsOp is an XDR Struct defines as:
//
//   struct SetLimitsOp
//    {
//        AccountID* account;
//        AccountType* accountType;
//
//        Limits limits;
//    	// reserved for future use
//    	union switch (LedgerVersion v)
//    	{
//    	case EMPTY_VERSION:
//    		void;
//    	}
//    	ext;
//    };
//
type SetLimitsOp struct {
	Account     *AccountId     `json:"account,omitempty"`
	AccountType *AccountType   `json:"accountType,omitempty"`
	Limits      Limits         `json:"limits,omitempty"`
	Ext         SetLimitsOpExt `json:"ext,omitempty"`
}

// SetLimitsResultCode is an XDR Enum defines as:
//
//   enum SetLimitsResultCode
//    {
//        // codes considered as "success" for the operation
//        SET_LIMITS_SUCCESS = 0,
//        // codes considered as "failure" for the operation
//        SET_LIMITS_MALFORMED = -1
//    };
//
type SetLimitsResultCode int32

const (
	SetLimitsResultCodeSetLimitsSuccess   SetLimitsResultCode = 0
	SetLimitsResultCodeSetLimitsMalformed SetLimitsResultCode = -1
)

var SetLimitsResultCodeAll = []SetLimitsResultCode{
	SetLimitsResultCodeSetLimitsSuccess,
	SetLimitsResultCodeSetLimitsMalformed,
}

var setLimitsResultCodeMap = map[int32]string{
	0:  "SetLimitsResultCodeSetLimitsSuccess",
	-1: "SetLimitsResultCodeSetLimitsMalformed",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for SetLimitsResultCode
func (e SetLimitsResultCode) ValidEnum(v int32) bool {
	_, ok := setLimitsResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e SetLimitsResultCode) String() string {
	name, _ := setLimitsResultCodeMap[int32(e)]
	return name
}

func (e SetLimitsResultCode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// SetLimitsResultSuccessExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type SetLimitsResultSuccessExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u SetLimitsResultSuccessExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of SetLimitsResultSuccessExt
func (u SetLimitsResultSuccessExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewSetLimitsResultSuccessExt creates a new  SetLimitsResultSuccessExt.
func NewSetLimitsResultSuccessExt(v LedgerVersion, value interface{}) (result SetLimitsResultSuccessExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// SetLimitsResultSuccess is an XDR NestedStruct defines as:
//
//   struct {
//    		// reserved for future use
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//    	}
//
type SetLimitsResultSuccess struct {
	Ext SetLimitsResultSuccessExt `json:"ext,omitempty"`
}

// SetLimitsResult is an XDR Union defines as:
//
//   union SetLimitsResult switch (SetLimitsResultCode code)
//    {
//    case SET_LIMITS_SUCCESS:
//        struct {
//    		// reserved for future use
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//    	} success;
//    default:
//        void;
//    };
//
type SetLimitsResult struct {
	Code    SetLimitsResultCode     `json:"code,omitempty"`
	Success *SetLimitsResultSuccess `json:"success,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u SetLimitsResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of SetLimitsResult
func (u SetLimitsResult) ArmForSwitch(sw int32) (string, bool) {
	switch SetLimitsResultCode(sw) {
	case SetLimitsResultCodeSetLimitsSuccess:
		return "Success", true
	default:
		return "", true
	}
}

// NewSetLimitsResult creates a new  SetLimitsResult.
func NewSetLimitsResult(code SetLimitsResultCode, value interface{}) (result SetLimitsResult, err error) {
	result.Code = code
	switch SetLimitsResultCode(code) {
	case SetLimitsResultCodeSetLimitsSuccess:
		tv, ok := value.(SetLimitsResultSuccess)
		if !ok {
			err = fmt.Errorf("invalid value, must be SetLimitsResultSuccess")
			return
		}
		result.Success = &tv
	default:
		// void
	}
	return
}

// MustSuccess retrieves the Success value from the union,
// panicing if the value is not set.
func (u SetLimitsResult) MustSuccess() SetLimitsResultSuccess {
	val, ok := u.GetSuccess()

	if !ok {
		panic("arm Success is not set")
	}

	return val
}

// GetSuccess retrieves the Success value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u SetLimitsResult) GetSuccess() (result SetLimitsResultSuccess, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "Success" {
		result = *u.Success
		ok = true
	}

	return
}

// ManageTrustAction is an XDR Enum defines as:
//
//   enum ManageTrustAction
//    {
//        TRUST_ADD = 0,
//        TRUST_REMOVE = 1
//    };
//
type ManageTrustAction int32

const (
	ManageTrustActionTrustAdd    ManageTrustAction = 0
	ManageTrustActionTrustRemove ManageTrustAction = 1
)

var ManageTrustActionAll = []ManageTrustAction{
	ManageTrustActionTrustAdd,
	ManageTrustActionTrustRemove,
}

var manageTrustActionMap = map[int32]string{
	0: "ManageTrustActionTrustAdd",
	1: "ManageTrustActionTrustRemove",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ManageTrustAction
func (e ManageTrustAction) ValidEnum(v int32) bool {
	_, ok := manageTrustActionMap[v]
	return ok
}

// String returns the name of `e`
func (e ManageTrustAction) String() string {
	name, _ := manageTrustActionMap[int32(e)]
	return name
}

func (e ManageTrustAction) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// TrustDataExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    	{
//    	case EMPTY_VERSION:
//    		void;
//    	}
//
type TrustDataExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u TrustDataExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of TrustDataExt
func (u TrustDataExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewTrustDataExt creates a new  TrustDataExt.
func NewTrustDataExt(v LedgerVersion, value interface{}) (result TrustDataExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// TrustData is an XDR Struct defines as:
//
//   struct TrustData {
//        TrustEntry trust;
//        ManageTrustAction action;
//    	// reserved for future use
//    	union switch (LedgerVersion v)
//    	{
//    	case EMPTY_VERSION:
//    		void;
//    	}
//    	ext;
//    };
//
type TrustData struct {
	Trust  TrustEntry        `json:"trust,omitempty"`
	Action ManageTrustAction `json:"action,omitempty"`
	Ext    TrustDataExt      `json:"ext,omitempty"`
}

// SetOptionsOpExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    	{
//    	case EMPTY_VERSION:
//    		void;
//    	}
//
type SetOptionsOpExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u SetOptionsOpExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of SetOptionsOpExt
func (u SetOptionsOpExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewSetOptionsOpExt creates a new  SetOptionsOpExt.
func NewSetOptionsOpExt(v LedgerVersion, value interface{}) (result SetOptionsOpExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// SetOptionsOp is an XDR Struct defines as:
//
//   struct SetOptionsOp
//    {
//        // account threshold manipulation
//        uint32* masterWeight; // weight of the master account
//        uint32* lowThreshold;
//        uint32* medThreshold;
//        uint32* highThreshold;
//
//        // Add, update or remove a signer for the account
//        // signer is deleted if the weight is 0
//        Signer* signer;
//
//        TrustData* trustData;
//    	// reserved for future use
//    	union switch (LedgerVersion v)
//    	{
//    	case EMPTY_VERSION:
//    		void;
//    	}
//    	ext;
//
//    };
//
type SetOptionsOp struct {
	MasterWeight  *Uint32         `json:"masterWeight,omitempty"`
	LowThreshold  *Uint32         `json:"lowThreshold,omitempty"`
	MedThreshold  *Uint32         `json:"medThreshold,omitempty"`
	HighThreshold *Uint32         `json:"highThreshold,omitempty"`
	Signer        *Signer         `json:"signer,omitempty"`
	TrustData     *TrustData      `json:"trustData,omitempty"`
	Ext           SetOptionsOpExt `json:"ext,omitempty"`
}

// SetOptionsResultCode is an XDR Enum defines as:
//
//   enum SetOptionsResultCode
//    {
//        // codes considered as "success" for the operation
//        SET_OPTIONS_SUCCESS = 0,
//        // codes considered as "failure" for the operation
//        SET_OPTIONS_TOO_MANY_SIGNERS = -1, // max number of signers already reached
//        SET_OPTIONS_THRESHOLD_OUT_OF_RANGE = -2, // bad value for weight/threshold
//        SET_OPTIONS_BAD_SIGNER = -3,             // signer cannot be masterkey
//        SET_OPTIONS_BALANCE_NOT_FOUND = -4,
//        SET_OPTIONS_TRUST_MALFORMED = -5,
//    	SET_OPTIONS_TRUST_TOO_MANY = -6,
//    	SET_OPTIONS_INVALID_SIGNER_VERSION = -7 // if signer version is higher than ledger version
//    };
//
type SetOptionsResultCode int32

const (
	SetOptionsResultCodeSetOptionsSuccess              SetOptionsResultCode = 0
	SetOptionsResultCodeSetOptionsTooManySigners       SetOptionsResultCode = -1
	SetOptionsResultCodeSetOptionsThresholdOutOfRange  SetOptionsResultCode = -2
	SetOptionsResultCodeSetOptionsBadSigner            SetOptionsResultCode = -3
	SetOptionsResultCodeSetOptionsBalanceNotFound      SetOptionsResultCode = -4
	SetOptionsResultCodeSetOptionsTrustMalformed       SetOptionsResultCode = -5
	SetOptionsResultCodeSetOptionsTrustTooMany         SetOptionsResultCode = -6
	SetOptionsResultCodeSetOptionsInvalidSignerVersion SetOptionsResultCode = -7
)

var SetOptionsResultCodeAll = []SetOptionsResultCode{
	SetOptionsResultCodeSetOptionsSuccess,
	SetOptionsResultCodeSetOptionsTooManySigners,
	SetOptionsResultCodeSetOptionsThresholdOutOfRange,
	SetOptionsResultCodeSetOptionsBadSigner,
	SetOptionsResultCodeSetOptionsBalanceNotFound,
	SetOptionsResultCodeSetOptionsTrustMalformed,
	SetOptionsResultCodeSetOptionsTrustTooMany,
	SetOptionsResultCodeSetOptionsInvalidSignerVersion,
}

var setOptionsResultCodeMap = map[int32]string{
	0:  "SetOptionsResultCodeSetOptionsSuccess",
	-1: "SetOptionsResultCodeSetOptionsTooManySigners",
	-2: "SetOptionsResultCodeSetOptionsThresholdOutOfRange",
	-3: "SetOptionsResultCodeSetOptionsBadSigner",
	-4: "SetOptionsResultCodeSetOptionsBalanceNotFound",
	-5: "SetOptionsResultCodeSetOptionsTrustMalformed",
	-6: "SetOptionsResultCodeSetOptionsTrustTooMany",
	-7: "SetOptionsResultCodeSetOptionsInvalidSignerVersion",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for SetOptionsResultCode
func (e SetOptionsResultCode) ValidEnum(v int32) bool {
	_, ok := setOptionsResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e SetOptionsResultCode) String() string {
	name, _ := setOptionsResultCodeMap[int32(e)]
	return name
}

func (e SetOptionsResultCode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// SetOptionsResultSuccessExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type SetOptionsResultSuccessExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u SetOptionsResultSuccessExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of SetOptionsResultSuccessExt
func (u SetOptionsResultSuccessExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewSetOptionsResultSuccessExt creates a new  SetOptionsResultSuccessExt.
func NewSetOptionsResultSuccessExt(v LedgerVersion, value interface{}) (result SetOptionsResultSuccessExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// SetOptionsResultSuccess is an XDR NestedStruct defines as:
//
//   struct {
//    		// reserved for future use
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//    	}
//
type SetOptionsResultSuccess struct {
	Ext SetOptionsResultSuccessExt `json:"ext,omitempty"`
}

// SetOptionsResult is an XDR Union defines as:
//
//   union SetOptionsResult switch (SetOptionsResultCode code)
//    {
//    case SET_OPTIONS_SUCCESS:
//        struct {
//    		// reserved for future use
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//    	} success;
//    default:
//        void;
//    };
//
type SetOptionsResult struct {
	Code    SetOptionsResultCode     `json:"code,omitempty"`
	Success *SetOptionsResultSuccess `json:"success,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u SetOptionsResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of SetOptionsResult
func (u SetOptionsResult) ArmForSwitch(sw int32) (string, bool) {
	switch SetOptionsResultCode(sw) {
	case SetOptionsResultCodeSetOptionsSuccess:
		return "Success", true
	default:
		return "", true
	}
}

// NewSetOptionsResult creates a new  SetOptionsResult.
func NewSetOptionsResult(code SetOptionsResultCode, value interface{}) (result SetOptionsResult, err error) {
	result.Code = code
	switch SetOptionsResultCode(code) {
	case SetOptionsResultCodeSetOptionsSuccess:
		tv, ok := value.(SetOptionsResultSuccess)
		if !ok {
			err = fmt.Errorf("invalid value, must be SetOptionsResultSuccess")
			return
		}
		result.Success = &tv
	default:
		// void
	}
	return
}

// MustSuccess retrieves the Success value from the union,
// panicing if the value is not set.
func (u SetOptionsResult) MustSuccess() SetOptionsResultSuccess {
	val, ok := u.GetSuccess()

	if !ok {
		panic("arm Success is not set")
	}

	return val
}

// GetSuccess retrieves the Success value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u SetOptionsResult) GetSuccess() (result SetOptionsResultSuccess, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "Success" {
		result = *u.Success
		ok = true
	}

	return
}

// PreEmissionExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    	{
//    	case EMPTY_VERSION:
//    		void;
//    	}
//
type PreEmissionExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u PreEmissionExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of PreEmissionExt
func (u PreEmissionExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewPreEmissionExt creates a new  PreEmissionExt.
func NewPreEmissionExt(v LedgerVersion, value interface{}) (result PreEmissionExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// PreEmission is an XDR Struct defines as:
//
//   struct PreEmission
//    {
//        string64 serialNumber;
//        AssetCode asset;
//        int64 amount;
//        DecoratedSignature signatures<20>;
//    	// reserved for future use
//    	union switch (LedgerVersion v)
//    	{
//    	case EMPTY_VERSION:
//    		void;
//    	}
//    	ext;
//    };
//
type PreEmission struct {
	SerialNumber String64             `json:"serialNumber,omitempty"`
	Asset        AssetCode            `json:"asset,omitempty"`
	Amount       Int64                `json:"amount,omitempty"`
	Signatures   []DecoratedSignature `json:"signatures,omitempty" xdrmaxsize:"20"`
	Ext          PreEmissionExt       `json:"ext,omitempty"`
}

// UploadPreemissionsOpExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    	{
//    	case EMPTY_VERSION:
//    		void;
//    	}
//
type UploadPreemissionsOpExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u UploadPreemissionsOpExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of UploadPreemissionsOpExt
func (u UploadPreemissionsOpExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewUploadPreemissionsOpExt creates a new  UploadPreemissionsOpExt.
func NewUploadPreemissionsOpExt(v LedgerVersion, value interface{}) (result UploadPreemissionsOpExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// UploadPreemissionsOp is an XDR Struct defines as:
//
//   struct UploadPreemissionsOp
//    {
//        PreEmission preEmissions<>;
//    	// reserved for future use
//    	union switch (LedgerVersion v)
//    	{
//    	case EMPTY_VERSION:
//    		void;
//    	}
//    	ext;
//    };
//
type UploadPreemissionsOp struct {
	PreEmissions []PreEmission           `json:"preEmissions,omitempty"`
	Ext          UploadPreemissionsOpExt `json:"ext,omitempty"`
}

// UploadPreemissionsResultCode is an XDR Enum defines as:
//
//   enum UploadPreemissionsResultCode
//    {
//        // codes considered as "success" for the operation
//        UPLOAD_PREEMISSIONS_SUCCESS = 0,
//
//        // codes considered as "failure" for the operation
//        UPLOAD_PREEMISSIONS_MALFORMED = -1,
//        UPLOAD_PREEMISSIONS_SERIAL_DUPLICATION = -2,    // serial is already used
//        UPLOAD_PREEMISSIONS_MALFORMED_PREEMISSIONS = -3, // if pre-emissions has empty signatures or zero amount etc
//        UPLOAD_PREEMISSIONS_ASSET_NOT_FOUND = -4,
//        UPLOAD_PREEMISSIONS_LINE_FULL = -5
//    };
//
type UploadPreemissionsResultCode int32

const (
	UploadPreemissionsResultCodeUploadPreemissionsSuccess               UploadPreemissionsResultCode = 0
	UploadPreemissionsResultCodeUploadPreemissionsMalformed             UploadPreemissionsResultCode = -1
	UploadPreemissionsResultCodeUploadPreemissionsSerialDuplication     UploadPreemissionsResultCode = -2
	UploadPreemissionsResultCodeUploadPreemissionsMalformedPreemissions UploadPreemissionsResultCode = -3
	UploadPreemissionsResultCodeUploadPreemissionsAssetNotFound         UploadPreemissionsResultCode = -4
	UploadPreemissionsResultCodeUploadPreemissionsLineFull              UploadPreemissionsResultCode = -5
)

var UploadPreemissionsResultCodeAll = []UploadPreemissionsResultCode{
	UploadPreemissionsResultCodeUploadPreemissionsSuccess,
	UploadPreemissionsResultCodeUploadPreemissionsMalformed,
	UploadPreemissionsResultCodeUploadPreemissionsSerialDuplication,
	UploadPreemissionsResultCodeUploadPreemissionsMalformedPreemissions,
	UploadPreemissionsResultCodeUploadPreemissionsAssetNotFound,
	UploadPreemissionsResultCodeUploadPreemissionsLineFull,
}

var uploadPreemissionsResultCodeMap = map[int32]string{
	0:  "UploadPreemissionsResultCodeUploadPreemissionsSuccess",
	-1: "UploadPreemissionsResultCodeUploadPreemissionsMalformed",
	-2: "UploadPreemissionsResultCodeUploadPreemissionsSerialDuplication",
	-3: "UploadPreemissionsResultCodeUploadPreemissionsMalformedPreemissions",
	-4: "UploadPreemissionsResultCodeUploadPreemissionsAssetNotFound",
	-5: "UploadPreemissionsResultCodeUploadPreemissionsLineFull",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for UploadPreemissionsResultCode
func (e UploadPreemissionsResultCode) ValidEnum(v int32) bool {
	_, ok := uploadPreemissionsResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e UploadPreemissionsResultCode) String() string {
	name, _ := uploadPreemissionsResultCodeMap[int32(e)]
	return name
}

func (e UploadPreemissionsResultCode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// UploadPreemissionsResultSuccessExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//
type UploadPreemissionsResultSuccessExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u UploadPreemissionsResultSuccessExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of UploadPreemissionsResultSuccessExt
func (u UploadPreemissionsResultSuccessExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewUploadPreemissionsResultSuccessExt creates a new  UploadPreemissionsResultSuccessExt.
func NewUploadPreemissionsResultSuccessExt(v LedgerVersion, value interface{}) (result UploadPreemissionsResultSuccessExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// UploadPreemissionsResultSuccess is an XDR NestedStruct defines as:
//
//   struct {
//    		// reserved for future use
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//    	}
//
type UploadPreemissionsResultSuccess struct {
	Ext UploadPreemissionsResultSuccessExt `json:"ext,omitempty"`
}

// UploadPreemissionsResult is an XDR Union defines as:
//
//   union UploadPreemissionsResult switch (UploadPreemissionsResultCode code)
//    {
//    case UPLOAD_PREEMISSIONS_SUCCESS:
//        struct {
//    		// reserved for future use
//    		union switch (LedgerVersion v)
//    		{
//    		case EMPTY_VERSION:
//    			void;
//    		}
//    		ext;
//    	} success;
//    default:
//        void;
//    };
//
type UploadPreemissionsResult struct {
	Code    UploadPreemissionsResultCode     `json:"code,omitempty"`
	Success *UploadPreemissionsResultSuccess `json:"success,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u UploadPreemissionsResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of UploadPreemissionsResult
func (u UploadPreemissionsResult) ArmForSwitch(sw int32) (string, bool) {
	switch UploadPreemissionsResultCode(sw) {
	case UploadPreemissionsResultCodeUploadPreemissionsSuccess:
		return "Success", true
	default:
		return "", true
	}
}

// NewUploadPreemissionsResult creates a new  UploadPreemissionsResult.
func NewUploadPreemissionsResult(code UploadPreemissionsResultCode, value interface{}) (result UploadPreemissionsResult, err error) {
	result.Code = code
	switch UploadPreemissionsResultCode(code) {
	case UploadPreemissionsResultCodeUploadPreemissionsSuccess:
		tv, ok := value.(UploadPreemissionsResultSuccess)
		if !ok {
			err = fmt.Errorf("invalid value, must be UploadPreemissionsResultSuccess")
			return
		}
		result.Success = &tv
	default:
		// void
	}
	return
}

// MustSuccess retrieves the Success value from the union,
// panicing if the value is not set.
func (u UploadPreemissionsResult) MustSuccess() UploadPreemissionsResultSuccess {
	val, ok := u.GetSuccess()

	if !ok {
		panic("arm Success is not set")
	}

	return val
}

// GetSuccess retrieves the Success value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u UploadPreemissionsResult) GetSuccess() (result UploadPreemissionsResultSuccess, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "Success" {
		result = *u.Success
		ok = true
	}

	return
}

// ErrorCode is an XDR Enum defines as:
//
//   enum ErrorCode
//    {
//        ERR_MISC = 0, // Unspecific error
//        ERR_DATA = 1, // Malformed data
//        ERR_CONF = 2, // Misconfiguration error
//        ERR_AUTH = 3, // Authentication failure
//        ERR_LOAD = 4  // System overloaded
//    };
//
type ErrorCode int32

const (
	ErrorCodeErrMisc ErrorCode = 0
	ErrorCodeErrData ErrorCode = 1
	ErrorCodeErrConf ErrorCode = 2
	ErrorCodeErrAuth ErrorCode = 3
	ErrorCodeErrLoad ErrorCode = 4
)

var ErrorCodeAll = []ErrorCode{
	ErrorCodeErrMisc,
	ErrorCodeErrData,
	ErrorCodeErrConf,
	ErrorCodeErrAuth,
	ErrorCodeErrLoad,
}

var errorCodeMap = map[int32]string{
	0: "ErrorCodeErrMisc",
	1: "ErrorCodeErrData",
	2: "ErrorCodeErrConf",
	3: "ErrorCodeErrAuth",
	4: "ErrorCodeErrLoad",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ErrorCode
func (e ErrorCode) ValidEnum(v int32) bool {
	_, ok := errorCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e ErrorCode) String() string {
	name, _ := errorCodeMap[int32(e)]
	return name
}

func (e ErrorCode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// Error is an XDR Struct defines as:
//
//   struct Error
//    {
//        ErrorCode code;
//        string msg<100>;
//    };
//
type Error struct {
	Code ErrorCode `json:"code,omitempty"`
	Msg  string    `json:"msg,omitempty" xdrmaxsize:"100"`
}

// AuthCert is an XDR Struct defines as:
//
//   struct AuthCert
//    {
//        Curve25519Public pubkey;
//        uint64 expiration;
//        Signature sig;
//    };
//
type AuthCert struct {
	Pubkey     Curve25519Public `json:"pubkey,omitempty"`
	Expiration Uint64           `json:"expiration,omitempty"`
	Sig        Signature        `json:"sig,omitempty"`
}

// Hello is an XDR Struct defines as:
//
//   struct Hello
//    {
//        uint32 ledgerVersion;
//        uint32 overlayVersion;
//        uint32 overlayMinVersion;
//        Hash networkID;
//        string versionStr<100>;
//        int listeningPort;
//        NodeID peerID;
//        AuthCert cert;
//        uint256 nonce;
//    };
//
type Hello struct {
	LedgerVersion     Uint32   `json:"ledgerVersion,omitempty"`
	OverlayVersion    Uint32   `json:"overlayVersion,omitempty"`
	OverlayMinVersion Uint32   `json:"overlayMinVersion,omitempty"`
	NetworkId         Hash     `json:"networkID,omitempty"`
	VersionStr        string   `json:"versionStr,omitempty" xdrmaxsize:"100"`
	ListeningPort     int32    `json:"listeningPort,omitempty"`
	PeerId            NodeId   `json:"peerID,omitempty"`
	Cert              AuthCert `json:"cert,omitempty"`
	Nonce             Uint256  `json:"nonce,omitempty"`
}

// Auth is an XDR Struct defines as:
//
//   struct Auth
//    {
//        // Empty message, just to confirm
//        // establishment of MAC keys.
//        int unused;
//    };
//
type Auth struct {
	Unused int32 `json:"unused,omitempty"`
}

// IpAddrType is an XDR Enum defines as:
//
//   enum IPAddrType
//    {
//        IPv4 = 0,
//        IPv6 = 1
//    };
//
type IpAddrType int32

const (
	IpAddrTypeIPv4 IpAddrType = 0
	IpAddrTypeIPv6 IpAddrType = 1
)

var IpAddrTypeAll = []IpAddrType{
	IpAddrTypeIPv4,
	IpAddrTypeIPv6,
}

var ipAddrTypeMap = map[int32]string{
	0: "IpAddrTypeIPv4",
	1: "IpAddrTypeIPv6",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for IpAddrType
func (e IpAddrType) ValidEnum(v int32) bool {
	_, ok := ipAddrTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e IpAddrType) String() string {
	name, _ := ipAddrTypeMap[int32(e)]
	return name
}

func (e IpAddrType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// PeerAddressIp is an XDR NestedUnion defines as:
//
//   union switch (IPAddrType type)
//        {
//        case IPv4:
//            opaque ipv4[4];
//        case IPv6:
//            opaque ipv6[16];
//        }
//
type PeerAddressIp struct {
	Type IpAddrType `json:"type,omitempty"`
	Ipv4 *[4]byte   `json:"ipv4,omitempty"`
	Ipv6 *[16]byte  `json:"ipv6,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u PeerAddressIp) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of PeerAddressIp
func (u PeerAddressIp) ArmForSwitch(sw int32) (string, bool) {
	switch IpAddrType(sw) {
	case IpAddrTypeIPv4:
		return "Ipv4", true
	case IpAddrTypeIPv6:
		return "Ipv6", true
	}
	return "-", false
}

// NewPeerAddressIp creates a new  PeerAddressIp.
func NewPeerAddressIp(aType IpAddrType, value interface{}) (result PeerAddressIp, err error) {
	result.Type = aType
	switch IpAddrType(aType) {
	case IpAddrTypeIPv4:
		tv, ok := value.([4]byte)
		if !ok {
			err = fmt.Errorf("invalid value, must be [4]byte")
			return
		}
		result.Ipv4 = &tv
	case IpAddrTypeIPv6:
		tv, ok := value.([16]byte)
		if !ok {
			err = fmt.Errorf("invalid value, must be [16]byte")
			return
		}
		result.Ipv6 = &tv
	}
	return
}

// MustIpv4 retrieves the Ipv4 value from the union,
// panicing if the value is not set.
func (u PeerAddressIp) MustIpv4() [4]byte {
	val, ok := u.GetIpv4()

	if !ok {
		panic("arm Ipv4 is not set")
	}

	return val
}

// GetIpv4 retrieves the Ipv4 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u PeerAddressIp) GetIpv4() (result [4]byte, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Ipv4" {
		result = *u.Ipv4
		ok = true
	}

	return
}

// MustIpv6 retrieves the Ipv6 value from the union,
// panicing if the value is not set.
func (u PeerAddressIp) MustIpv6() [16]byte {
	val, ok := u.GetIpv6()

	if !ok {
		panic("arm Ipv6 is not set")
	}

	return val
}

// GetIpv6 retrieves the Ipv6 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u PeerAddressIp) GetIpv6() (result [16]byte, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Ipv6" {
		result = *u.Ipv6
		ok = true
	}

	return
}

// PeerAddress is an XDR Struct defines as:
//
//   struct PeerAddress
//    {
//        union switch (IPAddrType type)
//        {
//        case IPv4:
//            opaque ipv4[4];
//        case IPv6:
//            opaque ipv6[16];
//        }
//        ip;
//        uint32 port;
//        uint32 numFailures;
//    };
//
type PeerAddress struct {
	Ip          PeerAddressIp `json:"ip,omitempty"`
	Port        Uint32        `json:"port,omitempty"`
	NumFailures Uint32        `json:"numFailures,omitempty"`
}

// MessageType is an XDR Enum defines as:
//
//   enum MessageType
//    {
//        ERROR_MSG = 0,
//        AUTH = 2,
//        DONT_HAVE = 3,
//
//        GET_PEERS = 4, // gets a list of peers this guy knows about
//        PEERS = 5,
//
//        GET_TX_SET = 6, // gets a particular txset by hash
//        TX_SET = 7,
//
//        TRANSACTION = 8, // pass on a tx you have heard about
//
//        // SCP
//        GET_SCP_QUORUMSET = 9,
//        SCP_QUORUMSET = 10,
//        SCP_MESSAGE = 11,
//        GET_SCP_STATE = 12,
//
//        // new messages
//        HELLO = 13
//    };
//
type MessageType int32

const (
	MessageTypeErrorMsg        MessageType = 0
	MessageTypeAuth            MessageType = 2
	MessageTypeDontHave        MessageType = 3
	MessageTypeGetPeers        MessageType = 4
	MessageTypePeers           MessageType = 5
	MessageTypeGetTxSet        MessageType = 6
	MessageTypeTxSet           MessageType = 7
	MessageTypeTransaction     MessageType = 8
	MessageTypeGetScpQuorumset MessageType = 9
	MessageTypeScpQuorumset    MessageType = 10
	MessageTypeScpMessage      MessageType = 11
	MessageTypeGetScpState     MessageType = 12
	MessageTypeHello           MessageType = 13
)

var MessageTypeAll = []MessageType{
	MessageTypeErrorMsg,
	MessageTypeAuth,
	MessageTypeDontHave,
	MessageTypeGetPeers,
	MessageTypePeers,
	MessageTypeGetTxSet,
	MessageTypeTxSet,
	MessageTypeTransaction,
	MessageTypeGetScpQuorumset,
	MessageTypeScpQuorumset,
	MessageTypeScpMessage,
	MessageTypeGetScpState,
	MessageTypeHello,
}

var messageTypeMap = map[int32]string{
	0:  "MessageTypeErrorMsg",
	2:  "MessageTypeAuth",
	3:  "MessageTypeDontHave",
	4:  "MessageTypeGetPeers",
	5:  "MessageTypePeers",
	6:  "MessageTypeGetTxSet",
	7:  "MessageTypeTxSet",
	8:  "MessageTypeTransaction",
	9:  "MessageTypeGetScpQuorumset",
	10: "MessageTypeScpQuorumset",
	11: "MessageTypeScpMessage",
	12: "MessageTypeGetScpState",
	13: "MessageTypeHello",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for MessageType
func (e MessageType) ValidEnum(v int32) bool {
	_, ok := messageTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e MessageType) String() string {
	name, _ := messageTypeMap[int32(e)]
	return name
}

func (e MessageType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// DontHave is an XDR Struct defines as:
//
//   struct DontHave
//    {
//        MessageType type;
//        uint256 reqHash;
//    };
//
type DontHave struct {
	Type    MessageType `json:"type,omitempty"`
	ReqHash Uint256     `json:"reqHash,omitempty"`
}

// StellarMessage is an XDR Union defines as:
//
//   union StellarMessage switch (MessageType type)
//    {
//    case ERROR_MSG:
//        Error error;
//    case HELLO:
//        Hello hello;
//    case AUTH:
//        Auth auth;
//    case DONT_HAVE:
//        DontHave dontHave;
//    case GET_PEERS:
//        void;
//    case PEERS:
//        PeerAddress peers<>;
//
//    case GET_TX_SET:
//        uint256 txSetHash;
//    case TX_SET:
//        TransactionSet txSet;
//
//    case TRANSACTION:
//        TransactionEnvelope transaction;
//
//    // SCP
//    case GET_SCP_QUORUMSET:
//        uint256 qSetHash;
//    case SCP_QUORUMSET:
//        SCPQuorumSet qSet;
//    case SCP_MESSAGE:
//        SCPEnvelope envelope;
//    case GET_SCP_STATE:
//        uint32 getSCPLedgerSeq; // ledger seq requested ; if 0, requests the latest
//    };
//
type StellarMessage struct {
	Type            MessageType          `json:"type,omitempty"`
	Error           *Error               `json:"error,omitempty"`
	Hello           *Hello               `json:"hello,omitempty"`
	Auth            *Auth                `json:"auth,omitempty"`
	DontHave        *DontHave            `json:"dontHave,omitempty"`
	Peers           *[]PeerAddress       `json:"peers,omitempty"`
	TxSetHash       *Uint256             `json:"txSetHash,omitempty"`
	TxSet           *TransactionSet      `json:"txSet,omitempty"`
	Transaction     *TransactionEnvelope `json:"transaction,omitempty"`
	QSetHash        *Uint256             `json:"qSetHash,omitempty"`
	QSet            *ScpQuorumSet        `json:"qSet,omitempty"`
	Envelope        *ScpEnvelope         `json:"envelope,omitempty"`
	GetScpLedgerSeq *Uint32              `json:"getSCPLedgerSeq,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u StellarMessage) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of StellarMessage
func (u StellarMessage) ArmForSwitch(sw int32) (string, bool) {
	switch MessageType(sw) {
	case MessageTypeErrorMsg:
		return "Error", true
	case MessageTypeHello:
		return "Hello", true
	case MessageTypeAuth:
		return "Auth", true
	case MessageTypeDontHave:
		return "DontHave", true
	case MessageTypeGetPeers:
		return "", true
	case MessageTypePeers:
		return "Peers", true
	case MessageTypeGetTxSet:
		return "TxSetHash", true
	case MessageTypeTxSet:
		return "TxSet", true
	case MessageTypeTransaction:
		return "Transaction", true
	case MessageTypeGetScpQuorumset:
		return "QSetHash", true
	case MessageTypeScpQuorumset:
		return "QSet", true
	case MessageTypeScpMessage:
		return "Envelope", true
	case MessageTypeGetScpState:
		return "GetScpLedgerSeq", true
	}
	return "-", false
}

// NewStellarMessage creates a new  StellarMessage.
func NewStellarMessage(aType MessageType, value interface{}) (result StellarMessage, err error) {
	result.Type = aType
	switch MessageType(aType) {
	case MessageTypeErrorMsg:
		tv, ok := value.(Error)
		if !ok {
			err = fmt.Errorf("invalid value, must be Error")
			return
		}
		result.Error = &tv
	case MessageTypeHello:
		tv, ok := value.(Hello)
		if !ok {
			err = fmt.Errorf("invalid value, must be Hello")
			return
		}
		result.Hello = &tv
	case MessageTypeAuth:
		tv, ok := value.(Auth)
		if !ok {
			err = fmt.Errorf("invalid value, must be Auth")
			return
		}
		result.Auth = &tv
	case MessageTypeDontHave:
		tv, ok := value.(DontHave)
		if !ok {
			err = fmt.Errorf("invalid value, must be DontHave")
			return
		}
		result.DontHave = &tv
	case MessageTypeGetPeers:
		// void
	case MessageTypePeers:
		tv, ok := value.([]PeerAddress)
		if !ok {
			err = fmt.Errorf("invalid value, must be []PeerAddress")
			return
		}
		result.Peers = &tv
	case MessageTypeGetTxSet:
		tv, ok := value.(Uint256)
		if !ok {
			err = fmt.Errorf("invalid value, must be Uint256")
			return
		}
		result.TxSetHash = &tv
	case MessageTypeTxSet:
		tv, ok := value.(TransactionSet)
		if !ok {
			err = fmt.Errorf("invalid value, must be TransactionSet")
			return
		}
		result.TxSet = &tv
	case MessageTypeTransaction:
		tv, ok := value.(TransactionEnvelope)
		if !ok {
			err = fmt.Errorf("invalid value, must be TransactionEnvelope")
			return
		}
		result.Transaction = &tv
	case MessageTypeGetScpQuorumset:
		tv, ok := value.(Uint256)
		if !ok {
			err = fmt.Errorf("invalid value, must be Uint256")
			return
		}
		result.QSetHash = &tv
	case MessageTypeScpQuorumset:
		tv, ok := value.(ScpQuorumSet)
		if !ok {
			err = fmt.Errorf("invalid value, must be ScpQuorumSet")
			return
		}
		result.QSet = &tv
	case MessageTypeScpMessage:
		tv, ok := value.(ScpEnvelope)
		if !ok {
			err = fmt.Errorf("invalid value, must be ScpEnvelope")
			return
		}
		result.Envelope = &tv
	case MessageTypeGetScpState:
		tv, ok := value.(Uint32)
		if !ok {
			err = fmt.Errorf("invalid value, must be Uint32")
			return
		}
		result.GetScpLedgerSeq = &tv
	}
	return
}

// MustError retrieves the Error value from the union,
// panicing if the value is not set.
func (u StellarMessage) MustError() Error {
	val, ok := u.GetError()

	if !ok {
		panic("arm Error is not set")
	}

	return val
}

// GetError retrieves the Error value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u StellarMessage) GetError() (result Error, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Error" {
		result = *u.Error
		ok = true
	}

	return
}

// MustHello retrieves the Hello value from the union,
// panicing if the value is not set.
func (u StellarMessage) MustHello() Hello {
	val, ok := u.GetHello()

	if !ok {
		panic("arm Hello is not set")
	}

	return val
}

// GetHello retrieves the Hello value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u StellarMessage) GetHello() (result Hello, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Hello" {
		result = *u.Hello
		ok = true
	}

	return
}

// MustAuth retrieves the Auth value from the union,
// panicing if the value is not set.
func (u StellarMessage) MustAuth() Auth {
	val, ok := u.GetAuth()

	if !ok {
		panic("arm Auth is not set")
	}

	return val
}

// GetAuth retrieves the Auth value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u StellarMessage) GetAuth() (result Auth, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Auth" {
		result = *u.Auth
		ok = true
	}

	return
}

// MustDontHave retrieves the DontHave value from the union,
// panicing if the value is not set.
func (u StellarMessage) MustDontHave() DontHave {
	val, ok := u.GetDontHave()

	if !ok {
		panic("arm DontHave is not set")
	}

	return val
}

// GetDontHave retrieves the DontHave value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u StellarMessage) GetDontHave() (result DontHave, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "DontHave" {
		result = *u.DontHave
		ok = true
	}

	return
}

// MustPeers retrieves the Peers value from the union,
// panicing if the value is not set.
func (u StellarMessage) MustPeers() []PeerAddress {
	val, ok := u.GetPeers()

	if !ok {
		panic("arm Peers is not set")
	}

	return val
}

// GetPeers retrieves the Peers value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u StellarMessage) GetPeers() (result []PeerAddress, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Peers" {
		result = *u.Peers
		ok = true
	}

	return
}

// MustTxSetHash retrieves the TxSetHash value from the union,
// panicing if the value is not set.
func (u StellarMessage) MustTxSetHash() Uint256 {
	val, ok := u.GetTxSetHash()

	if !ok {
		panic("arm TxSetHash is not set")
	}

	return val
}

// GetTxSetHash retrieves the TxSetHash value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u StellarMessage) GetTxSetHash() (result Uint256, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "TxSetHash" {
		result = *u.TxSetHash
		ok = true
	}

	return
}

// MustTxSet retrieves the TxSet value from the union,
// panicing if the value is not set.
func (u StellarMessage) MustTxSet() TransactionSet {
	val, ok := u.GetTxSet()

	if !ok {
		panic("arm TxSet is not set")
	}

	return val
}

// GetTxSet retrieves the TxSet value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u StellarMessage) GetTxSet() (result TransactionSet, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "TxSet" {
		result = *u.TxSet
		ok = true
	}

	return
}

// MustTransaction retrieves the Transaction value from the union,
// panicing if the value is not set.
func (u StellarMessage) MustTransaction() TransactionEnvelope {
	val, ok := u.GetTransaction()

	if !ok {
		panic("arm Transaction is not set")
	}

	return val
}

// GetTransaction retrieves the Transaction value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u StellarMessage) GetTransaction() (result TransactionEnvelope, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Transaction" {
		result = *u.Transaction
		ok = true
	}

	return
}

// MustQSetHash retrieves the QSetHash value from the union,
// panicing if the value is not set.
func (u StellarMessage) MustQSetHash() Uint256 {
	val, ok := u.GetQSetHash()

	if !ok {
		panic("arm QSetHash is not set")
	}

	return val
}

// GetQSetHash retrieves the QSetHash value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u StellarMessage) GetQSetHash() (result Uint256, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "QSetHash" {
		result = *u.QSetHash
		ok = true
	}

	return
}

// MustQSet retrieves the QSet value from the union,
// panicing if the value is not set.
func (u StellarMessage) MustQSet() ScpQuorumSet {
	val, ok := u.GetQSet()

	if !ok {
		panic("arm QSet is not set")
	}

	return val
}

// GetQSet retrieves the QSet value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u StellarMessage) GetQSet() (result ScpQuorumSet, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "QSet" {
		result = *u.QSet
		ok = true
	}

	return
}

// MustEnvelope retrieves the Envelope value from the union,
// panicing if the value is not set.
func (u StellarMessage) MustEnvelope() ScpEnvelope {
	val, ok := u.GetEnvelope()

	if !ok {
		panic("arm Envelope is not set")
	}

	return val
}

// GetEnvelope retrieves the Envelope value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u StellarMessage) GetEnvelope() (result ScpEnvelope, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Envelope" {
		result = *u.Envelope
		ok = true
	}

	return
}

// MustGetScpLedgerSeq retrieves the GetScpLedgerSeq value from the union,
// panicing if the value is not set.
func (u StellarMessage) MustGetScpLedgerSeq() Uint32 {
	val, ok := u.GetGetScpLedgerSeq()

	if !ok {
		panic("arm GetScpLedgerSeq is not set")
	}

	return val
}

// GetGetScpLedgerSeq retrieves the GetScpLedgerSeq value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u StellarMessage) GetGetScpLedgerSeq() (result Uint32, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "GetScpLedgerSeq" {
		result = *u.GetScpLedgerSeq
		ok = true
	}

	return
}

// AuthenticatedMessageV0 is an XDR NestedStruct defines as:
//
//   struct
//    {
//       uint64 sequence;
//       StellarMessage message;
//       HmacSha256Mac mac;
//        }
//
type AuthenticatedMessageV0 struct {
	Sequence Uint64         `json:"sequence,omitempty"`
	Message  StellarMessage `json:"message,omitempty"`
	Mac      HmacSha256Mac  `json:"mac,omitempty"`
}

// AuthenticatedMessage is an XDR Union defines as:
//
//   union AuthenticatedMessage switch (LedgerVersion v)
//    {
//    case EMPTY_VERSION:
//        struct
//    {
//       uint64 sequence;
//       StellarMessage message;
//       HmacSha256Mac mac;
//        } v0;
//    };
//
type AuthenticatedMessage struct {
	V  LedgerVersion           `json:"v,omitempty"`
	V0 *AuthenticatedMessageV0 `json:"v0,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u AuthenticatedMessage) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of AuthenticatedMessage
func (u AuthenticatedMessage) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "V0", true
	}
	return "-", false
}

// NewAuthenticatedMessage creates a new  AuthenticatedMessage.
func NewAuthenticatedMessage(v LedgerVersion, value interface{}) (result AuthenticatedMessage, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		tv, ok := value.(AuthenticatedMessageV0)
		if !ok {
			err = fmt.Errorf("invalid value, must be AuthenticatedMessageV0")
			return
		}
		result.V0 = &tv
	}
	return
}

// MustV0 retrieves the V0 value from the union,
// panicing if the value is not set.
func (u AuthenticatedMessage) MustV0() AuthenticatedMessageV0 {
	val, ok := u.GetV0()

	if !ok {
		panic("arm V0 is not set")
	}

	return val
}

// GetV0 retrieves the V0 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u AuthenticatedMessage) GetV0() (result AuthenticatedMessageV0, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.V))

	if armName == "V0" {
		result = *u.V0
		ok = true
	}

	return
}

// Value is an XDR Typedef defines as:
//
//   typedef opaque Value<>;
//
type Value []byte

// ScpBallot is an XDR Struct defines as:
//
//   struct SCPBallot
//    {
//        uint32 counter; // n
//        Value value;    // x
//    };
//
type ScpBallot struct {
	Counter Uint32 `json:"counter,omitempty"`
	Value   Value  `json:"value,omitempty"`
}

// ScpStatementType is an XDR Enum defines as:
//
//   enum SCPStatementType
//    {
//        SCP_ST_PREPARE = 0,
//        SCP_ST_CONFIRM = 1,
//        SCP_ST_EXTERNALIZE = 2,
//        SCP_ST_NOMINATE = 3
//    };
//
type ScpStatementType int32

const (
	ScpStatementTypeScpStPrepare     ScpStatementType = 0
	ScpStatementTypeScpStConfirm     ScpStatementType = 1
	ScpStatementTypeScpStExternalize ScpStatementType = 2
	ScpStatementTypeScpStNominate    ScpStatementType = 3
)

var ScpStatementTypeAll = []ScpStatementType{
	ScpStatementTypeScpStPrepare,
	ScpStatementTypeScpStConfirm,
	ScpStatementTypeScpStExternalize,
	ScpStatementTypeScpStNominate,
}

var scpStatementTypeMap = map[int32]string{
	0: "ScpStatementTypeScpStPrepare",
	1: "ScpStatementTypeScpStConfirm",
	2: "ScpStatementTypeScpStExternalize",
	3: "ScpStatementTypeScpStNominate",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for ScpStatementType
func (e ScpStatementType) ValidEnum(v int32) bool {
	_, ok := scpStatementTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e ScpStatementType) String() string {
	name, _ := scpStatementTypeMap[int32(e)]
	return name
}

func (e ScpStatementType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// ScpNomination is an XDR Struct defines as:
//
//   struct SCPNomination
//    {
//        Hash quorumSetHash; // D
//        Value votes<>;      // X
//        Value accepted<>;   // Y
//    };
//
type ScpNomination struct {
	QuorumSetHash Hash    `json:"quorumSetHash,omitempty"`
	Votes         []Value `json:"votes,omitempty"`
	Accepted      []Value `json:"accepted,omitempty"`
}

// ScpStatementPrepare is an XDR NestedStruct defines as:
//
//   struct
//            {
//                Hash quorumSetHash;       // D
//                SCPBallot ballot;         // b
//                SCPBallot* prepared;      // p
//                SCPBallot* preparedPrime; // p'
//                uint32 nC;                // c.n
//                uint32 nH;                // h.n
//            }
//
type ScpStatementPrepare struct {
	QuorumSetHash Hash       `json:"quorumSetHash,omitempty"`
	Ballot        ScpBallot  `json:"ballot,omitempty"`
	Prepared      *ScpBallot `json:"prepared,omitempty"`
	PreparedPrime *ScpBallot `json:"preparedPrime,omitempty"`
	NC            Uint32     `json:"nC,omitempty"`
	NH            Uint32     `json:"nH,omitempty"`
}

// ScpStatementConfirm is an XDR NestedStruct defines as:
//
//   struct
//            {
//                SCPBallot ballot;   // b
//                uint32 nPrepared;   // p.n
//                uint32 nCommit;     // c.n
//                uint32 nH;          // h.n
//                Hash quorumSetHash; // D
//            }
//
type ScpStatementConfirm struct {
	Ballot        ScpBallot `json:"ballot,omitempty"`
	NPrepared     Uint32    `json:"nPrepared,omitempty"`
	NCommit       Uint32    `json:"nCommit,omitempty"`
	NH            Uint32    `json:"nH,omitempty"`
	QuorumSetHash Hash      `json:"quorumSetHash,omitempty"`
}

// ScpStatementExternalize is an XDR NestedStruct defines as:
//
//   struct
//            {
//                SCPBallot commit;         // c
//                uint32 nH;                // h.n
//                Hash commitQuorumSetHash; // D used before EXTERNALIZE
//            }
//
type ScpStatementExternalize struct {
	Commit              ScpBallot `json:"commit,omitempty"`
	NH                  Uint32    `json:"nH,omitempty"`
	CommitQuorumSetHash Hash      `json:"commitQuorumSetHash,omitempty"`
}

// ScpStatementPledges is an XDR NestedUnion defines as:
//
//   union switch (SCPStatementType type)
//        {
//        case SCP_ST_PREPARE:
//            struct
//            {
//                Hash quorumSetHash;       // D
//                SCPBallot ballot;         // b
//                SCPBallot* prepared;      // p
//                SCPBallot* preparedPrime; // p'
//                uint32 nC;                // c.n
//                uint32 nH;                // h.n
//            } prepare;
//        case SCP_ST_CONFIRM:
//            struct
//            {
//                SCPBallot ballot;   // b
//                uint32 nPrepared;   // p.n
//                uint32 nCommit;     // c.n
//                uint32 nH;          // h.n
//                Hash quorumSetHash; // D
//            } confirm;
//        case SCP_ST_EXTERNALIZE:
//            struct
//            {
//                SCPBallot commit;         // c
//                uint32 nH;                // h.n
//                Hash commitQuorumSetHash; // D used before EXTERNALIZE
//            } externalize;
//        case SCP_ST_NOMINATE:
//            SCPNomination nominate;
//        }
//
type ScpStatementPledges struct {
	Type        ScpStatementType         `json:"type,omitempty"`
	Prepare     *ScpStatementPrepare     `json:"prepare,omitempty"`
	Confirm     *ScpStatementConfirm     `json:"confirm,omitempty"`
	Externalize *ScpStatementExternalize `json:"externalize,omitempty"`
	Nominate    *ScpNomination           `json:"nominate,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ScpStatementPledges) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ScpStatementPledges
func (u ScpStatementPledges) ArmForSwitch(sw int32) (string, bool) {
	switch ScpStatementType(sw) {
	case ScpStatementTypeScpStPrepare:
		return "Prepare", true
	case ScpStatementTypeScpStConfirm:
		return "Confirm", true
	case ScpStatementTypeScpStExternalize:
		return "Externalize", true
	case ScpStatementTypeScpStNominate:
		return "Nominate", true
	}
	return "-", false
}

// NewScpStatementPledges creates a new  ScpStatementPledges.
func NewScpStatementPledges(aType ScpStatementType, value interface{}) (result ScpStatementPledges, err error) {
	result.Type = aType
	switch ScpStatementType(aType) {
	case ScpStatementTypeScpStPrepare:
		tv, ok := value.(ScpStatementPrepare)
		if !ok {
			err = fmt.Errorf("invalid value, must be ScpStatementPrepare")
			return
		}
		result.Prepare = &tv
	case ScpStatementTypeScpStConfirm:
		tv, ok := value.(ScpStatementConfirm)
		if !ok {
			err = fmt.Errorf("invalid value, must be ScpStatementConfirm")
			return
		}
		result.Confirm = &tv
	case ScpStatementTypeScpStExternalize:
		tv, ok := value.(ScpStatementExternalize)
		if !ok {
			err = fmt.Errorf("invalid value, must be ScpStatementExternalize")
			return
		}
		result.Externalize = &tv
	case ScpStatementTypeScpStNominate:
		tv, ok := value.(ScpNomination)
		if !ok {
			err = fmt.Errorf("invalid value, must be ScpNomination")
			return
		}
		result.Nominate = &tv
	}
	return
}

// MustPrepare retrieves the Prepare value from the union,
// panicing if the value is not set.
func (u ScpStatementPledges) MustPrepare() ScpStatementPrepare {
	val, ok := u.GetPrepare()

	if !ok {
		panic("arm Prepare is not set")
	}

	return val
}

// GetPrepare retrieves the Prepare value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ScpStatementPledges) GetPrepare() (result ScpStatementPrepare, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Prepare" {
		result = *u.Prepare
		ok = true
	}

	return
}

// MustConfirm retrieves the Confirm value from the union,
// panicing if the value is not set.
func (u ScpStatementPledges) MustConfirm() ScpStatementConfirm {
	val, ok := u.GetConfirm()

	if !ok {
		panic("arm Confirm is not set")
	}

	return val
}

// GetConfirm retrieves the Confirm value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ScpStatementPledges) GetConfirm() (result ScpStatementConfirm, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Confirm" {
		result = *u.Confirm
		ok = true
	}

	return
}

// MustExternalize retrieves the Externalize value from the union,
// panicing if the value is not set.
func (u ScpStatementPledges) MustExternalize() ScpStatementExternalize {
	val, ok := u.GetExternalize()

	if !ok {
		panic("arm Externalize is not set")
	}

	return val
}

// GetExternalize retrieves the Externalize value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ScpStatementPledges) GetExternalize() (result ScpStatementExternalize, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Externalize" {
		result = *u.Externalize
		ok = true
	}

	return
}

// MustNominate retrieves the Nominate value from the union,
// panicing if the value is not set.
func (u ScpStatementPledges) MustNominate() ScpNomination {
	val, ok := u.GetNominate()

	if !ok {
		panic("arm Nominate is not set")
	}

	return val
}

// GetNominate retrieves the Nominate value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ScpStatementPledges) GetNominate() (result ScpNomination, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Nominate" {
		result = *u.Nominate
		ok = true
	}

	return
}

// ScpStatement is an XDR Struct defines as:
//
//   struct SCPStatement
//    {
//        NodeID nodeID;    // v
//        uint64 slotIndex; // i
//
//        union switch (SCPStatementType type)
//        {
//        case SCP_ST_PREPARE:
//            struct
//            {
//                Hash quorumSetHash;       // D
//                SCPBallot ballot;         // b
//                SCPBallot* prepared;      // p
//                SCPBallot* preparedPrime; // p'
//                uint32 nC;                // c.n
//                uint32 nH;                // h.n
//            } prepare;
//        case SCP_ST_CONFIRM:
//            struct
//            {
//                SCPBallot ballot;   // b
//                uint32 nPrepared;   // p.n
//                uint32 nCommit;     // c.n
//                uint32 nH;          // h.n
//                Hash quorumSetHash; // D
//            } confirm;
//        case SCP_ST_EXTERNALIZE:
//            struct
//            {
//                SCPBallot commit;         // c
//                uint32 nH;                // h.n
//                Hash commitQuorumSetHash; // D used before EXTERNALIZE
//            } externalize;
//        case SCP_ST_NOMINATE:
//            SCPNomination nominate;
//        }
//        pledges;
//    };
//
type ScpStatement struct {
	NodeId    NodeId              `json:"nodeID,omitempty"`
	SlotIndex Uint64              `json:"slotIndex,omitempty"`
	Pledges   ScpStatementPledges `json:"pledges,omitempty"`
}

// ScpEnvelope is an XDR Struct defines as:
//
//   struct SCPEnvelope
//    {
//        SCPStatement statement;
//        Signature signature;
//    };
//
type ScpEnvelope struct {
	Statement ScpStatement `json:"statement,omitempty"`
	Signature Signature    `json:"signature,omitempty"`
}

// ScpQuorumSet is an XDR Struct defines as:
//
//   struct SCPQuorumSet
//    {
//        uint32 threshold;
//        PublicKey validators<>;
//        SCPQuorumSet innerSets<>;
//    };
//
type ScpQuorumSet struct {
	Threshold  Uint32         `json:"threshold,omitempty"`
	Validators []PublicKey    `json:"validators,omitempty"`
	InnerSets  []ScpQuorumSet `json:"innerSets,omitempty"`
}

// OperationBody is an XDR NestedUnion defines as:
//
//   union switch (OperationType type)
//        {
//        case CREATE_ACCOUNT:
//            CreateAccountOp createAccountOp;
//        case PAYMENT:
//            PaymentOp paymentOp;
//        case SET_OPTIONS:
//            SetOptionsOp setOptionsOp;
//    	case MANAGE_COINS_EMISSION_REQUEST:
//    		ManageCoinsEmissionRequestOp manageCoinsEmissionRequestOp;
//    	case REVIEW_COINS_EMISSION_REQUEST:
//    		ReviewCoinsEmissionRequestOp reviewCoinsEmissionRequestOp;
//        case SET_FEES:
//            SetFeesOp setFeesOp;
//    	case MANAGE_ACCOUNT:
//    		ManageAccountOp manageAccountOp;
//    	case FORFEIT:
//    		ForfeitOp forfeitOp;
//    	case MANAGE_FORFEIT_REQUEST:
//    		ManageForfeitRequestOp manageForfeitRequestOp;
//    	case RECOVER:
//    		RecoverOp recoverOp;
//    	case MANAGE_BALANCE:
//    		ManageBalanceOp manageBalanceOp;
//    	case REVIEW_PAYMENT_REQUEST:
//    		ReviewPaymentRequestOp reviewPaymentRequestOp;
//        case MANAGE_ASSET:
//            ManageAssetOp manageAssetOp;
//        case DEMURRAGE:
//            DemurrageOp demurrageOp;
//        case UPLOAD_PREEMISSIONS:
//            UploadPreemissionsOp uploadPreemissionsOp;
//        case SET_LIMITS:
//            SetLimitsOp setLimitsOp;
//        case DIRECT_DEBIT:
//            DirectDebitOp directDebitOp;
//    	case MANAGE_ASSET_PAIR:
//    		ManageAssetPairOp manageAssetPairOp;
//    	case MANAGE_OFFER:
//    		ManageOfferOp manageOfferOp;
//        case MANAGE_INVOICE:
//            ManageInvoiceOp manageInvoiceOp;
//        }
//
type OperationBody struct {
	Type                         OperationType                 `json:"type,omitempty"`
	CreateAccountOp              *CreateAccountOp              `json:"createAccountOp,omitempty"`
	PaymentOp                    *PaymentOp                    `json:"paymentOp,omitempty"`
	SetOptionsOp                 *SetOptionsOp                 `json:"setOptionsOp,omitempty"`
	ManageCoinsEmissionRequestOp *ManageCoinsEmissionRequestOp `json:"manageCoinsEmissionRequestOp,omitempty"`
	ReviewCoinsEmissionRequestOp *ReviewCoinsEmissionRequestOp `json:"reviewCoinsEmissionRequestOp,omitempty"`
	SetFeesOp                    *SetFeesOp                    `json:"setFeesOp,omitempty"`
	ManageAccountOp              *ManageAccountOp              `json:"manageAccountOp,omitempty"`
	ForfeitOp                    *ForfeitOp                    `json:"forfeitOp,omitempty"`
	ManageForfeitRequestOp       *ManageForfeitRequestOp       `json:"manageForfeitRequestOp,omitempty"`
	RecoverOp                    *RecoverOp                    `json:"recoverOp,omitempty"`
	ManageBalanceOp              *ManageBalanceOp              `json:"manageBalanceOp,omitempty"`
	ReviewPaymentRequestOp       *ReviewPaymentRequestOp       `json:"reviewPaymentRequestOp,omitempty"`
	ManageAssetOp                *ManageAssetOp                `json:"manageAssetOp,omitempty"`
	DemurrageOp                  *DemurrageOp                  `json:"demurrageOp,omitempty"`
	UploadPreemissionsOp         *UploadPreemissionsOp         `json:"uploadPreemissionsOp,omitempty"`
	SetLimitsOp                  *SetLimitsOp                  `json:"setLimitsOp,omitempty"`
	DirectDebitOp                *DirectDebitOp                `json:"directDebitOp,omitempty"`
	ManageAssetPairOp            *ManageAssetPairOp            `json:"manageAssetPairOp,omitempty"`
	ManageOfferOp                *ManageOfferOp                `json:"manageOfferOp,omitempty"`
	ManageInvoiceOp              *ManageInvoiceOp              `json:"manageInvoiceOp,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u OperationBody) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of OperationBody
func (u OperationBody) ArmForSwitch(sw int32) (string, bool) {
	switch OperationType(sw) {
	case OperationTypeCreateAccount:
		return "CreateAccountOp", true
	case OperationTypePayment:
		return "PaymentOp", true
	case OperationTypeSetOptions:
		return "SetOptionsOp", true
	case OperationTypeManageCoinsEmissionRequest:
		return "ManageCoinsEmissionRequestOp", true
	case OperationTypeReviewCoinsEmissionRequest:
		return "ReviewCoinsEmissionRequestOp", true
	case OperationTypeSetFees:
		return "SetFeesOp", true
	case OperationTypeManageAccount:
		return "ManageAccountOp", true
	case OperationTypeForfeit:
		return "ForfeitOp", true
	case OperationTypeManageForfeitRequest:
		return "ManageForfeitRequestOp", true
	case OperationTypeRecover:
		return "RecoverOp", true
	case OperationTypeManageBalance:
		return "ManageBalanceOp", true
	case OperationTypeReviewPaymentRequest:
		return "ReviewPaymentRequestOp", true
	case OperationTypeManageAsset:
		return "ManageAssetOp", true
	case OperationTypeDemurrage:
		return "DemurrageOp", true
	case OperationTypeUploadPreemissions:
		return "UploadPreemissionsOp", true
	case OperationTypeSetLimits:
		return "SetLimitsOp", true
	case OperationTypeDirectDebit:
		return "DirectDebitOp", true
	case OperationTypeManageAssetPair:
		return "ManageAssetPairOp", true
	case OperationTypeManageOffer:
		return "ManageOfferOp", true
	case OperationTypeManageInvoice:
		return "ManageInvoiceOp", true
	}
	return "-", false
}

// NewOperationBody creates a new  OperationBody.
func NewOperationBody(aType OperationType, value interface{}) (result OperationBody, err error) {
	result.Type = aType
	switch OperationType(aType) {
	case OperationTypeCreateAccount:
		tv, ok := value.(CreateAccountOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be CreateAccountOp")
			return
		}
		result.CreateAccountOp = &tv
	case OperationTypePayment:
		tv, ok := value.(PaymentOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be PaymentOp")
			return
		}
		result.PaymentOp = &tv
	case OperationTypeSetOptions:
		tv, ok := value.(SetOptionsOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be SetOptionsOp")
			return
		}
		result.SetOptionsOp = &tv
	case OperationTypeManageCoinsEmissionRequest:
		tv, ok := value.(ManageCoinsEmissionRequestOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageCoinsEmissionRequestOp")
			return
		}
		result.ManageCoinsEmissionRequestOp = &tv
	case OperationTypeReviewCoinsEmissionRequest:
		tv, ok := value.(ReviewCoinsEmissionRequestOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be ReviewCoinsEmissionRequestOp")
			return
		}
		result.ReviewCoinsEmissionRequestOp = &tv
	case OperationTypeSetFees:
		tv, ok := value.(SetFeesOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be SetFeesOp")
			return
		}
		result.SetFeesOp = &tv
	case OperationTypeManageAccount:
		tv, ok := value.(ManageAccountOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageAccountOp")
			return
		}
		result.ManageAccountOp = &tv
	case OperationTypeForfeit:
		tv, ok := value.(ForfeitOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be ForfeitOp")
			return
		}
		result.ForfeitOp = &tv
	case OperationTypeManageForfeitRequest:
		tv, ok := value.(ManageForfeitRequestOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageForfeitRequestOp")
			return
		}
		result.ManageForfeitRequestOp = &tv
	case OperationTypeRecover:
		tv, ok := value.(RecoverOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be RecoverOp")
			return
		}
		result.RecoverOp = &tv
	case OperationTypeManageBalance:
		tv, ok := value.(ManageBalanceOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageBalanceOp")
			return
		}
		result.ManageBalanceOp = &tv
	case OperationTypeReviewPaymentRequest:
		tv, ok := value.(ReviewPaymentRequestOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be ReviewPaymentRequestOp")
			return
		}
		result.ReviewPaymentRequestOp = &tv
	case OperationTypeManageAsset:
		tv, ok := value.(ManageAssetOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageAssetOp")
			return
		}
		result.ManageAssetOp = &tv
	case OperationTypeDemurrage:
		tv, ok := value.(DemurrageOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be DemurrageOp")
			return
		}
		result.DemurrageOp = &tv
	case OperationTypeUploadPreemissions:
		tv, ok := value.(UploadPreemissionsOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be UploadPreemissionsOp")
			return
		}
		result.UploadPreemissionsOp = &tv
	case OperationTypeSetLimits:
		tv, ok := value.(SetLimitsOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be SetLimitsOp")
			return
		}
		result.SetLimitsOp = &tv
	case OperationTypeDirectDebit:
		tv, ok := value.(DirectDebitOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be DirectDebitOp")
			return
		}
		result.DirectDebitOp = &tv
	case OperationTypeManageAssetPair:
		tv, ok := value.(ManageAssetPairOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageAssetPairOp")
			return
		}
		result.ManageAssetPairOp = &tv
	case OperationTypeManageOffer:
		tv, ok := value.(ManageOfferOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageOfferOp")
			return
		}
		result.ManageOfferOp = &tv
	case OperationTypeManageInvoice:
		tv, ok := value.(ManageInvoiceOp)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageInvoiceOp")
			return
		}
		result.ManageInvoiceOp = &tv
	}
	return
}

// MustCreateAccountOp retrieves the CreateAccountOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustCreateAccountOp() CreateAccountOp {
	val, ok := u.GetCreateAccountOp()

	if !ok {
		panic("arm CreateAccountOp is not set")
	}

	return val
}

// GetCreateAccountOp retrieves the CreateAccountOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetCreateAccountOp() (result CreateAccountOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "CreateAccountOp" {
		result = *u.CreateAccountOp
		ok = true
	}

	return
}

// MustPaymentOp retrieves the PaymentOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustPaymentOp() PaymentOp {
	val, ok := u.GetPaymentOp()

	if !ok {
		panic("arm PaymentOp is not set")
	}

	return val
}

// GetPaymentOp retrieves the PaymentOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetPaymentOp() (result PaymentOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "PaymentOp" {
		result = *u.PaymentOp
		ok = true
	}

	return
}

// MustSetOptionsOp retrieves the SetOptionsOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustSetOptionsOp() SetOptionsOp {
	val, ok := u.GetSetOptionsOp()

	if !ok {
		panic("arm SetOptionsOp is not set")
	}

	return val
}

// GetSetOptionsOp retrieves the SetOptionsOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetSetOptionsOp() (result SetOptionsOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "SetOptionsOp" {
		result = *u.SetOptionsOp
		ok = true
	}

	return
}

// MustManageCoinsEmissionRequestOp retrieves the ManageCoinsEmissionRequestOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustManageCoinsEmissionRequestOp() ManageCoinsEmissionRequestOp {
	val, ok := u.GetManageCoinsEmissionRequestOp()

	if !ok {
		panic("arm ManageCoinsEmissionRequestOp is not set")
	}

	return val
}

// GetManageCoinsEmissionRequestOp retrieves the ManageCoinsEmissionRequestOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetManageCoinsEmissionRequestOp() (result ManageCoinsEmissionRequestOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ManageCoinsEmissionRequestOp" {
		result = *u.ManageCoinsEmissionRequestOp
		ok = true
	}

	return
}

// MustReviewCoinsEmissionRequestOp retrieves the ReviewCoinsEmissionRequestOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustReviewCoinsEmissionRequestOp() ReviewCoinsEmissionRequestOp {
	val, ok := u.GetReviewCoinsEmissionRequestOp()

	if !ok {
		panic("arm ReviewCoinsEmissionRequestOp is not set")
	}

	return val
}

// GetReviewCoinsEmissionRequestOp retrieves the ReviewCoinsEmissionRequestOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetReviewCoinsEmissionRequestOp() (result ReviewCoinsEmissionRequestOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ReviewCoinsEmissionRequestOp" {
		result = *u.ReviewCoinsEmissionRequestOp
		ok = true
	}

	return
}

// MustSetFeesOp retrieves the SetFeesOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustSetFeesOp() SetFeesOp {
	val, ok := u.GetSetFeesOp()

	if !ok {
		panic("arm SetFeesOp is not set")
	}

	return val
}

// GetSetFeesOp retrieves the SetFeesOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetSetFeesOp() (result SetFeesOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "SetFeesOp" {
		result = *u.SetFeesOp
		ok = true
	}

	return
}

// MustManageAccountOp retrieves the ManageAccountOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustManageAccountOp() ManageAccountOp {
	val, ok := u.GetManageAccountOp()

	if !ok {
		panic("arm ManageAccountOp is not set")
	}

	return val
}

// GetManageAccountOp retrieves the ManageAccountOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetManageAccountOp() (result ManageAccountOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ManageAccountOp" {
		result = *u.ManageAccountOp
		ok = true
	}

	return
}

// MustForfeitOp retrieves the ForfeitOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustForfeitOp() ForfeitOp {
	val, ok := u.GetForfeitOp()

	if !ok {
		panic("arm ForfeitOp is not set")
	}

	return val
}

// GetForfeitOp retrieves the ForfeitOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetForfeitOp() (result ForfeitOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ForfeitOp" {
		result = *u.ForfeitOp
		ok = true
	}

	return
}

// MustManageForfeitRequestOp retrieves the ManageForfeitRequestOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustManageForfeitRequestOp() ManageForfeitRequestOp {
	val, ok := u.GetManageForfeitRequestOp()

	if !ok {
		panic("arm ManageForfeitRequestOp is not set")
	}

	return val
}

// GetManageForfeitRequestOp retrieves the ManageForfeitRequestOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetManageForfeitRequestOp() (result ManageForfeitRequestOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ManageForfeitRequestOp" {
		result = *u.ManageForfeitRequestOp
		ok = true
	}

	return
}

// MustRecoverOp retrieves the RecoverOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustRecoverOp() RecoverOp {
	val, ok := u.GetRecoverOp()

	if !ok {
		panic("arm RecoverOp is not set")
	}

	return val
}

// GetRecoverOp retrieves the RecoverOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetRecoverOp() (result RecoverOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "RecoverOp" {
		result = *u.RecoverOp
		ok = true
	}

	return
}

// MustManageBalanceOp retrieves the ManageBalanceOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustManageBalanceOp() ManageBalanceOp {
	val, ok := u.GetManageBalanceOp()

	if !ok {
		panic("arm ManageBalanceOp is not set")
	}

	return val
}

// GetManageBalanceOp retrieves the ManageBalanceOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetManageBalanceOp() (result ManageBalanceOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ManageBalanceOp" {
		result = *u.ManageBalanceOp
		ok = true
	}

	return
}

// MustReviewPaymentRequestOp retrieves the ReviewPaymentRequestOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustReviewPaymentRequestOp() ReviewPaymentRequestOp {
	val, ok := u.GetReviewPaymentRequestOp()

	if !ok {
		panic("arm ReviewPaymentRequestOp is not set")
	}

	return val
}

// GetReviewPaymentRequestOp retrieves the ReviewPaymentRequestOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetReviewPaymentRequestOp() (result ReviewPaymentRequestOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ReviewPaymentRequestOp" {
		result = *u.ReviewPaymentRequestOp
		ok = true
	}

	return
}

// MustManageAssetOp retrieves the ManageAssetOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustManageAssetOp() ManageAssetOp {
	val, ok := u.GetManageAssetOp()

	if !ok {
		panic("arm ManageAssetOp is not set")
	}

	return val
}

// GetManageAssetOp retrieves the ManageAssetOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetManageAssetOp() (result ManageAssetOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ManageAssetOp" {
		result = *u.ManageAssetOp
		ok = true
	}

	return
}

// MustDemurrageOp retrieves the DemurrageOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustDemurrageOp() DemurrageOp {
	val, ok := u.GetDemurrageOp()

	if !ok {
		panic("arm DemurrageOp is not set")
	}

	return val
}

// GetDemurrageOp retrieves the DemurrageOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetDemurrageOp() (result DemurrageOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "DemurrageOp" {
		result = *u.DemurrageOp
		ok = true
	}

	return
}

// MustUploadPreemissionsOp retrieves the UploadPreemissionsOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustUploadPreemissionsOp() UploadPreemissionsOp {
	val, ok := u.GetUploadPreemissionsOp()

	if !ok {
		panic("arm UploadPreemissionsOp is not set")
	}

	return val
}

// GetUploadPreemissionsOp retrieves the UploadPreemissionsOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetUploadPreemissionsOp() (result UploadPreemissionsOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "UploadPreemissionsOp" {
		result = *u.UploadPreemissionsOp
		ok = true
	}

	return
}

// MustSetLimitsOp retrieves the SetLimitsOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustSetLimitsOp() SetLimitsOp {
	val, ok := u.GetSetLimitsOp()

	if !ok {
		panic("arm SetLimitsOp is not set")
	}

	return val
}

// GetSetLimitsOp retrieves the SetLimitsOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetSetLimitsOp() (result SetLimitsOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "SetLimitsOp" {
		result = *u.SetLimitsOp
		ok = true
	}

	return
}

// MustDirectDebitOp retrieves the DirectDebitOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustDirectDebitOp() DirectDebitOp {
	val, ok := u.GetDirectDebitOp()

	if !ok {
		panic("arm DirectDebitOp is not set")
	}

	return val
}

// GetDirectDebitOp retrieves the DirectDebitOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetDirectDebitOp() (result DirectDebitOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "DirectDebitOp" {
		result = *u.DirectDebitOp
		ok = true
	}

	return
}

// MustManageAssetPairOp retrieves the ManageAssetPairOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustManageAssetPairOp() ManageAssetPairOp {
	val, ok := u.GetManageAssetPairOp()

	if !ok {
		panic("arm ManageAssetPairOp is not set")
	}

	return val
}

// GetManageAssetPairOp retrieves the ManageAssetPairOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetManageAssetPairOp() (result ManageAssetPairOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ManageAssetPairOp" {
		result = *u.ManageAssetPairOp
		ok = true
	}

	return
}

// MustManageOfferOp retrieves the ManageOfferOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustManageOfferOp() ManageOfferOp {
	val, ok := u.GetManageOfferOp()

	if !ok {
		panic("arm ManageOfferOp is not set")
	}

	return val
}

// GetManageOfferOp retrieves the ManageOfferOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetManageOfferOp() (result ManageOfferOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ManageOfferOp" {
		result = *u.ManageOfferOp
		ok = true
	}

	return
}

// MustManageInvoiceOp retrieves the ManageInvoiceOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustManageInvoiceOp() ManageInvoiceOp {
	val, ok := u.GetManageInvoiceOp()

	if !ok {
		panic("arm ManageInvoiceOp is not set")
	}

	return val
}

// GetManageInvoiceOp retrieves the ManageInvoiceOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetManageInvoiceOp() (result ManageInvoiceOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ManageInvoiceOp" {
		result = *u.ManageInvoiceOp
		ok = true
	}

	return
}

// Operation is an XDR Struct defines as:
//
//   struct Operation
//    {
//        // sourceAccount is the account used to run the operation
//        // if not set, the runtime defaults to "sourceAccount" specified at
//        // the transaction level
//        AccountID* sourceAccount;
//
//        union switch (OperationType type)
//        {
//        case CREATE_ACCOUNT:
//            CreateAccountOp createAccountOp;
//        case PAYMENT:
//            PaymentOp paymentOp;
//        case SET_OPTIONS:
//            SetOptionsOp setOptionsOp;
//    	case MANAGE_COINS_EMISSION_REQUEST:
//    		ManageCoinsEmissionRequestOp manageCoinsEmissionRequestOp;
//    	case REVIEW_COINS_EMISSION_REQUEST:
//    		ReviewCoinsEmissionRequestOp reviewCoinsEmissionRequestOp;
//        case SET_FEES:
//            SetFeesOp setFeesOp;
//    	case MANAGE_ACCOUNT:
//    		ManageAccountOp manageAccountOp;
//    	case FORFEIT:
//    		ForfeitOp forfeitOp;
//    	case MANAGE_FORFEIT_REQUEST:
//    		ManageForfeitRequestOp manageForfeitRequestOp;
//    	case RECOVER:
//    		RecoverOp recoverOp;
//    	case MANAGE_BALANCE:
//    		ManageBalanceOp manageBalanceOp;
//    	case REVIEW_PAYMENT_REQUEST:
//    		ReviewPaymentRequestOp reviewPaymentRequestOp;
//        case MANAGE_ASSET:
//            ManageAssetOp manageAssetOp;
//        case DEMURRAGE:
//            DemurrageOp demurrageOp;
//        case UPLOAD_PREEMISSIONS:
//            UploadPreemissionsOp uploadPreemissionsOp;
//        case SET_LIMITS:
//            SetLimitsOp setLimitsOp;
//        case DIRECT_DEBIT:
//            DirectDebitOp directDebitOp;
//    	case MANAGE_ASSET_PAIR:
//    		ManageAssetPairOp manageAssetPairOp;
//    	case MANAGE_OFFER:
//    		ManageOfferOp manageOfferOp;
//        case MANAGE_INVOICE:
//            ManageInvoiceOp manageInvoiceOp;
//        }
//        body;
//    };
//
type Operation struct {
	SourceAccount *AccountId    `json:"sourceAccount,omitempty"`
	Body          OperationBody `json:"body,omitempty"`
}

// MemoType is an XDR Enum defines as:
//
//   enum MemoType
//    {
//        MEMO_NONE = 0,
//        MEMO_TEXT = 1,
//        MEMO_ID = 2,
//        MEMO_HASH = 3,
//        MEMO_RETURN = 4
//    };
//
type MemoType int32

const (
	MemoTypeMemoNone   MemoType = 0
	MemoTypeMemoText   MemoType = 1
	MemoTypeMemoId     MemoType = 2
	MemoTypeMemoHash   MemoType = 3
	MemoTypeMemoReturn MemoType = 4
)

var MemoTypeAll = []MemoType{
	MemoTypeMemoNone,
	MemoTypeMemoText,
	MemoTypeMemoId,
	MemoTypeMemoHash,
	MemoTypeMemoReturn,
}

var memoTypeMap = map[int32]string{
	0: "MemoTypeMemoNone",
	1: "MemoTypeMemoText",
	2: "MemoTypeMemoId",
	3: "MemoTypeMemoHash",
	4: "MemoTypeMemoReturn",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for MemoType
func (e MemoType) ValidEnum(v int32) bool {
	_, ok := memoTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e MemoType) String() string {
	name, _ := memoTypeMap[int32(e)]
	return name
}

func (e MemoType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// Memo is an XDR Union defines as:
//
//   union Memo switch (MemoType type)
//    {
//    case MEMO_NONE:
//        void;
//    case MEMO_TEXT:
//        string text<28>;
//    case MEMO_ID:
//        uint64 id;
//    case MEMO_HASH:
//        Hash hash; // the hash of what to pull from the content server
//    case MEMO_RETURN:
//        Hash retHash; // the hash of the tx you are rejecting
//    };
//
type Memo struct {
	Type    MemoType `json:"type,omitempty"`
	Text    *string  `json:"text,omitempty" xdrmaxsize:"28"`
	Id      *Uint64  `json:"id,omitempty"`
	Hash    *Hash    `json:"hash,omitempty"`
	RetHash *Hash    `json:"retHash,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u Memo) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of Memo
func (u Memo) ArmForSwitch(sw int32) (string, bool) {
	switch MemoType(sw) {
	case MemoTypeMemoNone:
		return "", true
	case MemoTypeMemoText:
		return "Text", true
	case MemoTypeMemoId:
		return "Id", true
	case MemoTypeMemoHash:
		return "Hash", true
	case MemoTypeMemoReturn:
		return "RetHash", true
	}
	return "-", false
}

// NewMemo creates a new  Memo.
func NewMemo(aType MemoType, value interface{}) (result Memo, err error) {
	result.Type = aType
	switch MemoType(aType) {
	case MemoTypeMemoNone:
		// void
	case MemoTypeMemoText:
		tv, ok := value.(string)
		if !ok {
			err = fmt.Errorf("invalid value, must be string")
			return
		}
		result.Text = &tv
	case MemoTypeMemoId:
		tv, ok := value.(Uint64)
		if !ok {
			err = fmt.Errorf("invalid value, must be Uint64")
			return
		}
		result.Id = &tv
	case MemoTypeMemoHash:
		tv, ok := value.(Hash)
		if !ok {
			err = fmt.Errorf("invalid value, must be Hash")
			return
		}
		result.Hash = &tv
	case MemoTypeMemoReturn:
		tv, ok := value.(Hash)
		if !ok {
			err = fmt.Errorf("invalid value, must be Hash")
			return
		}
		result.RetHash = &tv
	}
	return
}

// MustText retrieves the Text value from the union,
// panicing if the value is not set.
func (u Memo) MustText() string {
	val, ok := u.GetText()

	if !ok {
		panic("arm Text is not set")
	}

	return val
}

// GetText retrieves the Text value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u Memo) GetText() (result string, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Text" {
		result = *u.Text
		ok = true
	}

	return
}

// MustId retrieves the Id value from the union,
// panicing if the value is not set.
func (u Memo) MustId() Uint64 {
	val, ok := u.GetId()

	if !ok {
		panic("arm Id is not set")
	}

	return val
}

// GetId retrieves the Id value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u Memo) GetId() (result Uint64, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Id" {
		result = *u.Id
		ok = true
	}

	return
}

// MustHash retrieves the Hash value from the union,
// panicing if the value is not set.
func (u Memo) MustHash() Hash {
	val, ok := u.GetHash()

	if !ok {
		panic("arm Hash is not set")
	}

	return val
}

// GetHash retrieves the Hash value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u Memo) GetHash() (result Hash, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Hash" {
		result = *u.Hash
		ok = true
	}

	return
}

// MustRetHash retrieves the RetHash value from the union,
// panicing if the value is not set.
func (u Memo) MustRetHash() Hash {
	val, ok := u.GetRetHash()

	if !ok {
		panic("arm RetHash is not set")
	}

	return val
}

// GetRetHash retrieves the RetHash value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u Memo) GetRetHash() (result Hash, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "RetHash" {
		result = *u.RetHash
		ok = true
	}

	return
}

// TimeBounds is an XDR Struct defines as:
//
//   struct TimeBounds
//    {
//        uint64 minTime;
//        uint64 maxTime; // 0 here means no maxTime
//    };
//
type TimeBounds struct {
	MinTime Uint64 `json:"minTime,omitempty"`
	MaxTime Uint64 `json:"maxTime,omitempty"`
}

// TransactionExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type TransactionExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u TransactionExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of TransactionExt
func (u TransactionExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewTransactionExt creates a new  TransactionExt.
func NewTransactionExt(v LedgerVersion, value interface{}) (result TransactionExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// Transaction is an XDR Struct defines as:
//
//   struct Transaction
//    {
//        // account used to run the transaction
//        AccountID sourceAccount;
//
//        Salt salt;
//
//        // validity range (inclusive) for the last ledger close time
//        TimeBounds timeBounds;
//
//        Memo memo;
//
//        Operation operations<100>;
//
//        // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type Transaction struct {
	SourceAccount AccountId      `json:"sourceAccount,omitempty"`
	Salt          Salt           `json:"salt,omitempty"`
	TimeBounds    TimeBounds     `json:"timeBounds,omitempty"`
	Memo          Memo           `json:"memo,omitempty"`
	Operations    []Operation    `json:"operations,omitempty" xdrmaxsize:"100"`
	Ext           TransactionExt `json:"ext,omitempty"`
}

// TransactionEnvelope is an XDR Struct defines as:
//
//   struct TransactionEnvelope
//    {
//        Transaction tx;
//        DecoratedSignature signatures<20>;
//    };
//
type TransactionEnvelope struct {
	Tx         Transaction          `json:"tx,omitempty"`
	Signatures []DecoratedSignature `json:"signatures,omitempty" xdrmaxsize:"20"`
}

// OperationResultCode is an XDR Enum defines as:
//
//   enum OperationResultCode
//    {
//        opINNER = 0, // inner object result is valid
//
//        opBAD_AUTH = -1,      // too few valid signatures / wrong network
//        opNO_ACCOUNT = -2,    // source account was not found
//    	opNOT_ALLOWED = -3,   // operation is not allowed for this type of source account
//    	opACCOUNT_BLOCKED = -4, // account is blocked
//        opNO_COUNTERPARTY = -5,
//        opCOUNTERPARTY_BLOCKED = -6,
//        opCOUNTERPARTY_WRONG_TYPE = -7,
//    	opBAD_AUTH_EXTRA = -8
//    };
//
type OperationResultCode int32

const (
	OperationResultCodeOpInner                 OperationResultCode = 0
	OperationResultCodeOpBadAuth               OperationResultCode = -1
	OperationResultCodeOpNoAccount             OperationResultCode = -2
	OperationResultCodeOpNotAllowed            OperationResultCode = -3
	OperationResultCodeOpAccountBlocked        OperationResultCode = -4
	OperationResultCodeOpNoCounterparty        OperationResultCode = -5
	OperationResultCodeOpCounterpartyBlocked   OperationResultCode = -6
	OperationResultCodeOpCounterpartyWrongType OperationResultCode = -7
	OperationResultCodeOpBadAuthExtra          OperationResultCode = -8
)

var OperationResultCodeAll = []OperationResultCode{
	OperationResultCodeOpInner,
	OperationResultCodeOpBadAuth,
	OperationResultCodeOpNoAccount,
	OperationResultCodeOpNotAllowed,
	OperationResultCodeOpAccountBlocked,
	OperationResultCodeOpNoCounterparty,
	OperationResultCodeOpCounterpartyBlocked,
	OperationResultCodeOpCounterpartyWrongType,
	OperationResultCodeOpBadAuthExtra,
}

var operationResultCodeMap = map[int32]string{
	0:  "OperationResultCodeOpInner",
	-1: "OperationResultCodeOpBadAuth",
	-2: "OperationResultCodeOpNoAccount",
	-3: "OperationResultCodeOpNotAllowed",
	-4: "OperationResultCodeOpAccountBlocked",
	-5: "OperationResultCodeOpNoCounterparty",
	-6: "OperationResultCodeOpCounterpartyBlocked",
	-7: "OperationResultCodeOpCounterpartyWrongType",
	-8: "OperationResultCodeOpBadAuthExtra",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for OperationResultCode
func (e OperationResultCode) ValidEnum(v int32) bool {
	_, ok := operationResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e OperationResultCode) String() string {
	name, _ := operationResultCodeMap[int32(e)]
	return name
}

func (e OperationResultCode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// OperationResultTr is an XDR NestedUnion defines as:
//
//   union switch (OperationType type)
//        {
//        case CREATE_ACCOUNT:
//            CreateAccountResult createAccountResult;
//        case PAYMENT:
//            PaymentResult paymentResult;
//        case SET_OPTIONS:
//            SetOptionsResult setOptionsResult;
//    	case MANAGE_COINS_EMISSION_REQUEST:
//    		ManageCoinsEmissionRequestResult manageCoinsEmissionRequestResult;
//    	case REVIEW_COINS_EMISSION_REQUEST:
//    		ReviewCoinsEmissionRequestResult reviewCoinsEmissionRequestResult;
//        case SET_FEES:
//            SetFeesResult setFeesResult;
//    	case MANAGE_ACCOUNT:
//    		ManageAccountResult manageAccountResult;
//    	case FORFEIT:
//    		ForfeitResult forfeitResult;
//        case MANAGE_FORFEIT_REQUEST:
//    		ManageForfeitRequestResult manageForfeitRequestResult;
//        case RECOVER:
//    		RecoverResult recoverResult;
//        case MANAGE_BALANCE:
//            ManageBalanceResult manageBalanceResult;
//        case REVIEW_PAYMENT_REQUEST:
//            ReviewPaymentRequestResult reviewPaymentRequestResult;
//        case MANAGE_ASSET:
//            ManageAssetResult manageAssetResult;
//        case DEMURRAGE:
//            DemurrageResult demurrageResult;
//        case UPLOAD_PREEMISSIONS:
//            UploadPreemissionsResult uploadPreemissionsResult;
//        case SET_LIMITS:
//            SetLimitsResult setLimitsResult;
//        case DIRECT_DEBIT:
//            DirectDebitResult directDebitResult;
//    	case MANAGE_ASSET_PAIR:
//    		ManageAssetPairResult manageAssetPairResult;
//    	case MANAGE_OFFER:
//    		ManageOfferResult manageOfferResult;
//    	case MANAGE_INVOICE:
//    		ManageInvoiceResult manageInvoiceResult;
//        }
//
type OperationResultTr struct {
	Type                             OperationType                     `json:"type,omitempty"`
	CreateAccountResult              *CreateAccountResult              `json:"createAccountResult,omitempty"`
	PaymentResult                    *PaymentResult                    `json:"paymentResult,omitempty"`
	SetOptionsResult                 *SetOptionsResult                 `json:"setOptionsResult,omitempty"`
	ManageCoinsEmissionRequestResult *ManageCoinsEmissionRequestResult `json:"manageCoinsEmissionRequestResult,omitempty"`
	ReviewCoinsEmissionRequestResult *ReviewCoinsEmissionRequestResult `json:"reviewCoinsEmissionRequestResult,omitempty"`
	SetFeesResult                    *SetFeesResult                    `json:"setFeesResult,omitempty"`
	ManageAccountResult              *ManageAccountResult              `json:"manageAccountResult,omitempty"`
	ForfeitResult                    *ForfeitResult                    `json:"forfeitResult,omitempty"`
	ManageForfeitRequestResult       *ManageForfeitRequestResult       `json:"manageForfeitRequestResult,omitempty"`
	RecoverResult                    *RecoverResult                    `json:"recoverResult,omitempty"`
	ManageBalanceResult              *ManageBalanceResult              `json:"manageBalanceResult,omitempty"`
	ReviewPaymentRequestResult       *ReviewPaymentRequestResult       `json:"reviewPaymentRequestResult,omitempty"`
	ManageAssetResult                *ManageAssetResult                `json:"manageAssetResult,omitempty"`
	DemurrageResult                  *DemurrageResult                  `json:"demurrageResult,omitempty"`
	UploadPreemissionsResult         *UploadPreemissionsResult         `json:"uploadPreemissionsResult,omitempty"`
	SetLimitsResult                  *SetLimitsResult                  `json:"setLimitsResult,omitempty"`
	DirectDebitResult                *DirectDebitResult                `json:"directDebitResult,omitempty"`
	ManageAssetPairResult            *ManageAssetPairResult            `json:"manageAssetPairResult,omitempty"`
	ManageOfferResult                *ManageOfferResult                `json:"manageOfferResult,omitempty"`
	ManageInvoiceResult              *ManageInvoiceResult              `json:"manageInvoiceResult,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u OperationResultTr) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of OperationResultTr
func (u OperationResultTr) ArmForSwitch(sw int32) (string, bool) {
	switch OperationType(sw) {
	case OperationTypeCreateAccount:
		return "CreateAccountResult", true
	case OperationTypePayment:
		return "PaymentResult", true
	case OperationTypeSetOptions:
		return "SetOptionsResult", true
	case OperationTypeManageCoinsEmissionRequest:
		return "ManageCoinsEmissionRequestResult", true
	case OperationTypeReviewCoinsEmissionRequest:
		return "ReviewCoinsEmissionRequestResult", true
	case OperationTypeSetFees:
		return "SetFeesResult", true
	case OperationTypeManageAccount:
		return "ManageAccountResult", true
	case OperationTypeForfeit:
		return "ForfeitResult", true
	case OperationTypeManageForfeitRequest:
		return "ManageForfeitRequestResult", true
	case OperationTypeRecover:
		return "RecoverResult", true
	case OperationTypeManageBalance:
		return "ManageBalanceResult", true
	case OperationTypeReviewPaymentRequest:
		return "ReviewPaymentRequestResult", true
	case OperationTypeManageAsset:
		return "ManageAssetResult", true
	case OperationTypeDemurrage:
		return "DemurrageResult", true
	case OperationTypeUploadPreemissions:
		return "UploadPreemissionsResult", true
	case OperationTypeSetLimits:
		return "SetLimitsResult", true
	case OperationTypeDirectDebit:
		return "DirectDebitResult", true
	case OperationTypeManageAssetPair:
		return "ManageAssetPairResult", true
	case OperationTypeManageOffer:
		return "ManageOfferResult", true
	case OperationTypeManageInvoice:
		return "ManageInvoiceResult", true
	}
	return "-", false
}

// NewOperationResultTr creates a new  OperationResultTr.
func NewOperationResultTr(aType OperationType, value interface{}) (result OperationResultTr, err error) {
	result.Type = aType
	switch OperationType(aType) {
	case OperationTypeCreateAccount:
		tv, ok := value.(CreateAccountResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be CreateAccountResult")
			return
		}
		result.CreateAccountResult = &tv
	case OperationTypePayment:
		tv, ok := value.(PaymentResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be PaymentResult")
			return
		}
		result.PaymentResult = &tv
	case OperationTypeSetOptions:
		tv, ok := value.(SetOptionsResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be SetOptionsResult")
			return
		}
		result.SetOptionsResult = &tv
	case OperationTypeManageCoinsEmissionRequest:
		tv, ok := value.(ManageCoinsEmissionRequestResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageCoinsEmissionRequestResult")
			return
		}
		result.ManageCoinsEmissionRequestResult = &tv
	case OperationTypeReviewCoinsEmissionRequest:
		tv, ok := value.(ReviewCoinsEmissionRequestResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be ReviewCoinsEmissionRequestResult")
			return
		}
		result.ReviewCoinsEmissionRequestResult = &tv
	case OperationTypeSetFees:
		tv, ok := value.(SetFeesResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be SetFeesResult")
			return
		}
		result.SetFeesResult = &tv
	case OperationTypeManageAccount:
		tv, ok := value.(ManageAccountResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageAccountResult")
			return
		}
		result.ManageAccountResult = &tv
	case OperationTypeForfeit:
		tv, ok := value.(ForfeitResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be ForfeitResult")
			return
		}
		result.ForfeitResult = &tv
	case OperationTypeManageForfeitRequest:
		tv, ok := value.(ManageForfeitRequestResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageForfeitRequestResult")
			return
		}
		result.ManageForfeitRequestResult = &tv
	case OperationTypeRecover:
		tv, ok := value.(RecoverResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be RecoverResult")
			return
		}
		result.RecoverResult = &tv
	case OperationTypeManageBalance:
		tv, ok := value.(ManageBalanceResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageBalanceResult")
			return
		}
		result.ManageBalanceResult = &tv
	case OperationTypeReviewPaymentRequest:
		tv, ok := value.(ReviewPaymentRequestResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be ReviewPaymentRequestResult")
			return
		}
		result.ReviewPaymentRequestResult = &tv
	case OperationTypeManageAsset:
		tv, ok := value.(ManageAssetResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageAssetResult")
			return
		}
		result.ManageAssetResult = &tv
	case OperationTypeDemurrage:
		tv, ok := value.(DemurrageResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be DemurrageResult")
			return
		}
		result.DemurrageResult = &tv
	case OperationTypeUploadPreemissions:
		tv, ok := value.(UploadPreemissionsResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be UploadPreemissionsResult")
			return
		}
		result.UploadPreemissionsResult = &tv
	case OperationTypeSetLimits:
		tv, ok := value.(SetLimitsResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be SetLimitsResult")
			return
		}
		result.SetLimitsResult = &tv
	case OperationTypeDirectDebit:
		tv, ok := value.(DirectDebitResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be DirectDebitResult")
			return
		}
		result.DirectDebitResult = &tv
	case OperationTypeManageAssetPair:
		tv, ok := value.(ManageAssetPairResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageAssetPairResult")
			return
		}
		result.ManageAssetPairResult = &tv
	case OperationTypeManageOffer:
		tv, ok := value.(ManageOfferResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageOfferResult")
			return
		}
		result.ManageOfferResult = &tv
	case OperationTypeManageInvoice:
		tv, ok := value.(ManageInvoiceResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be ManageInvoiceResult")
			return
		}
		result.ManageInvoiceResult = &tv
	}
	return
}

// MustCreateAccountResult retrieves the CreateAccountResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustCreateAccountResult() CreateAccountResult {
	val, ok := u.GetCreateAccountResult()

	if !ok {
		panic("arm CreateAccountResult is not set")
	}

	return val
}

// GetCreateAccountResult retrieves the CreateAccountResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetCreateAccountResult() (result CreateAccountResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "CreateAccountResult" {
		result = *u.CreateAccountResult
		ok = true
	}

	return
}

// MustPaymentResult retrieves the PaymentResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustPaymentResult() PaymentResult {
	val, ok := u.GetPaymentResult()

	if !ok {
		panic("arm PaymentResult is not set")
	}

	return val
}

// GetPaymentResult retrieves the PaymentResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetPaymentResult() (result PaymentResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "PaymentResult" {
		result = *u.PaymentResult
		ok = true
	}

	return
}

// MustSetOptionsResult retrieves the SetOptionsResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustSetOptionsResult() SetOptionsResult {
	val, ok := u.GetSetOptionsResult()

	if !ok {
		panic("arm SetOptionsResult is not set")
	}

	return val
}

// GetSetOptionsResult retrieves the SetOptionsResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetSetOptionsResult() (result SetOptionsResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "SetOptionsResult" {
		result = *u.SetOptionsResult
		ok = true
	}

	return
}

// MustManageCoinsEmissionRequestResult retrieves the ManageCoinsEmissionRequestResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustManageCoinsEmissionRequestResult() ManageCoinsEmissionRequestResult {
	val, ok := u.GetManageCoinsEmissionRequestResult()

	if !ok {
		panic("arm ManageCoinsEmissionRequestResult is not set")
	}

	return val
}

// GetManageCoinsEmissionRequestResult retrieves the ManageCoinsEmissionRequestResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetManageCoinsEmissionRequestResult() (result ManageCoinsEmissionRequestResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ManageCoinsEmissionRequestResult" {
		result = *u.ManageCoinsEmissionRequestResult
		ok = true
	}

	return
}

// MustReviewCoinsEmissionRequestResult retrieves the ReviewCoinsEmissionRequestResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustReviewCoinsEmissionRequestResult() ReviewCoinsEmissionRequestResult {
	val, ok := u.GetReviewCoinsEmissionRequestResult()

	if !ok {
		panic("arm ReviewCoinsEmissionRequestResult is not set")
	}

	return val
}

// GetReviewCoinsEmissionRequestResult retrieves the ReviewCoinsEmissionRequestResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetReviewCoinsEmissionRequestResult() (result ReviewCoinsEmissionRequestResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ReviewCoinsEmissionRequestResult" {
		result = *u.ReviewCoinsEmissionRequestResult
		ok = true
	}

	return
}

// MustSetFeesResult retrieves the SetFeesResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustSetFeesResult() SetFeesResult {
	val, ok := u.GetSetFeesResult()

	if !ok {
		panic("arm SetFeesResult is not set")
	}

	return val
}

// GetSetFeesResult retrieves the SetFeesResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetSetFeesResult() (result SetFeesResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "SetFeesResult" {
		result = *u.SetFeesResult
		ok = true
	}

	return
}

// MustManageAccountResult retrieves the ManageAccountResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustManageAccountResult() ManageAccountResult {
	val, ok := u.GetManageAccountResult()

	if !ok {
		panic("arm ManageAccountResult is not set")
	}

	return val
}

// GetManageAccountResult retrieves the ManageAccountResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetManageAccountResult() (result ManageAccountResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ManageAccountResult" {
		result = *u.ManageAccountResult
		ok = true
	}

	return
}

// MustForfeitResult retrieves the ForfeitResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustForfeitResult() ForfeitResult {
	val, ok := u.GetForfeitResult()

	if !ok {
		panic("arm ForfeitResult is not set")
	}

	return val
}

// GetForfeitResult retrieves the ForfeitResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetForfeitResult() (result ForfeitResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ForfeitResult" {
		result = *u.ForfeitResult
		ok = true
	}

	return
}

// MustManageForfeitRequestResult retrieves the ManageForfeitRequestResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustManageForfeitRequestResult() ManageForfeitRequestResult {
	val, ok := u.GetManageForfeitRequestResult()

	if !ok {
		panic("arm ManageForfeitRequestResult is not set")
	}

	return val
}

// GetManageForfeitRequestResult retrieves the ManageForfeitRequestResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetManageForfeitRequestResult() (result ManageForfeitRequestResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ManageForfeitRequestResult" {
		result = *u.ManageForfeitRequestResult
		ok = true
	}

	return
}

// MustRecoverResult retrieves the RecoverResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustRecoverResult() RecoverResult {
	val, ok := u.GetRecoverResult()

	if !ok {
		panic("arm RecoverResult is not set")
	}

	return val
}

// GetRecoverResult retrieves the RecoverResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetRecoverResult() (result RecoverResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "RecoverResult" {
		result = *u.RecoverResult
		ok = true
	}

	return
}

// MustManageBalanceResult retrieves the ManageBalanceResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustManageBalanceResult() ManageBalanceResult {
	val, ok := u.GetManageBalanceResult()

	if !ok {
		panic("arm ManageBalanceResult is not set")
	}

	return val
}

// GetManageBalanceResult retrieves the ManageBalanceResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetManageBalanceResult() (result ManageBalanceResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ManageBalanceResult" {
		result = *u.ManageBalanceResult
		ok = true
	}

	return
}

// MustReviewPaymentRequestResult retrieves the ReviewPaymentRequestResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustReviewPaymentRequestResult() ReviewPaymentRequestResult {
	val, ok := u.GetReviewPaymentRequestResult()

	if !ok {
		panic("arm ReviewPaymentRequestResult is not set")
	}

	return val
}

// GetReviewPaymentRequestResult retrieves the ReviewPaymentRequestResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetReviewPaymentRequestResult() (result ReviewPaymentRequestResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ReviewPaymentRequestResult" {
		result = *u.ReviewPaymentRequestResult
		ok = true
	}

	return
}

// MustManageAssetResult retrieves the ManageAssetResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustManageAssetResult() ManageAssetResult {
	val, ok := u.GetManageAssetResult()

	if !ok {
		panic("arm ManageAssetResult is not set")
	}

	return val
}

// GetManageAssetResult retrieves the ManageAssetResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetManageAssetResult() (result ManageAssetResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ManageAssetResult" {
		result = *u.ManageAssetResult
		ok = true
	}

	return
}

// MustDemurrageResult retrieves the DemurrageResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustDemurrageResult() DemurrageResult {
	val, ok := u.GetDemurrageResult()

	if !ok {
		panic("arm DemurrageResult is not set")
	}

	return val
}

// GetDemurrageResult retrieves the DemurrageResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetDemurrageResult() (result DemurrageResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "DemurrageResult" {
		result = *u.DemurrageResult
		ok = true
	}

	return
}

// MustUploadPreemissionsResult retrieves the UploadPreemissionsResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustUploadPreemissionsResult() UploadPreemissionsResult {
	val, ok := u.GetUploadPreemissionsResult()

	if !ok {
		panic("arm UploadPreemissionsResult is not set")
	}

	return val
}

// GetUploadPreemissionsResult retrieves the UploadPreemissionsResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetUploadPreemissionsResult() (result UploadPreemissionsResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "UploadPreemissionsResult" {
		result = *u.UploadPreemissionsResult
		ok = true
	}

	return
}

// MustSetLimitsResult retrieves the SetLimitsResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustSetLimitsResult() SetLimitsResult {
	val, ok := u.GetSetLimitsResult()

	if !ok {
		panic("arm SetLimitsResult is not set")
	}

	return val
}

// GetSetLimitsResult retrieves the SetLimitsResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetSetLimitsResult() (result SetLimitsResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "SetLimitsResult" {
		result = *u.SetLimitsResult
		ok = true
	}

	return
}

// MustDirectDebitResult retrieves the DirectDebitResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustDirectDebitResult() DirectDebitResult {
	val, ok := u.GetDirectDebitResult()

	if !ok {
		panic("arm DirectDebitResult is not set")
	}

	return val
}

// GetDirectDebitResult retrieves the DirectDebitResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetDirectDebitResult() (result DirectDebitResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "DirectDebitResult" {
		result = *u.DirectDebitResult
		ok = true
	}

	return
}

// MustManageAssetPairResult retrieves the ManageAssetPairResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustManageAssetPairResult() ManageAssetPairResult {
	val, ok := u.GetManageAssetPairResult()

	if !ok {
		panic("arm ManageAssetPairResult is not set")
	}

	return val
}

// GetManageAssetPairResult retrieves the ManageAssetPairResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetManageAssetPairResult() (result ManageAssetPairResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ManageAssetPairResult" {
		result = *u.ManageAssetPairResult
		ok = true
	}

	return
}

// MustManageOfferResult retrieves the ManageOfferResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustManageOfferResult() ManageOfferResult {
	val, ok := u.GetManageOfferResult()

	if !ok {
		panic("arm ManageOfferResult is not set")
	}

	return val
}

// GetManageOfferResult retrieves the ManageOfferResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetManageOfferResult() (result ManageOfferResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ManageOfferResult" {
		result = *u.ManageOfferResult
		ok = true
	}

	return
}

// MustManageInvoiceResult retrieves the ManageInvoiceResult value from the union,
// panicing if the value is not set.
func (u OperationResultTr) MustManageInvoiceResult() ManageInvoiceResult {
	val, ok := u.GetManageInvoiceResult()

	if !ok {
		panic("arm ManageInvoiceResult is not set")
	}

	return val
}

// GetManageInvoiceResult retrieves the ManageInvoiceResult value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResultTr) GetManageInvoiceResult() (result ManageInvoiceResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ManageInvoiceResult" {
		result = *u.ManageInvoiceResult
		ok = true
	}

	return
}

// OperationResult is an XDR Union defines as:
//
//   union OperationResult switch (OperationResultCode code)
//    {
//    case opINNER:
//        union switch (OperationType type)
//        {
//        case CREATE_ACCOUNT:
//            CreateAccountResult createAccountResult;
//        case PAYMENT:
//            PaymentResult paymentResult;
//        case SET_OPTIONS:
//            SetOptionsResult setOptionsResult;
//    	case MANAGE_COINS_EMISSION_REQUEST:
//    		ManageCoinsEmissionRequestResult manageCoinsEmissionRequestResult;
//    	case REVIEW_COINS_EMISSION_REQUEST:
//    		ReviewCoinsEmissionRequestResult reviewCoinsEmissionRequestResult;
//        case SET_FEES:
//            SetFeesResult setFeesResult;
//    	case MANAGE_ACCOUNT:
//    		ManageAccountResult manageAccountResult;
//    	case FORFEIT:
//    		ForfeitResult forfeitResult;
//        case MANAGE_FORFEIT_REQUEST:
//    		ManageForfeitRequestResult manageForfeitRequestResult;
//        case RECOVER:
//    		RecoverResult recoverResult;
//        case MANAGE_BALANCE:
//            ManageBalanceResult manageBalanceResult;
//        case REVIEW_PAYMENT_REQUEST:
//            ReviewPaymentRequestResult reviewPaymentRequestResult;
//        case MANAGE_ASSET:
//            ManageAssetResult manageAssetResult;
//        case DEMURRAGE:
//            DemurrageResult demurrageResult;
//        case UPLOAD_PREEMISSIONS:
//            UploadPreemissionsResult uploadPreemissionsResult;
//        case SET_LIMITS:
//            SetLimitsResult setLimitsResult;
//        case DIRECT_DEBIT:
//            DirectDebitResult directDebitResult;
//    	case MANAGE_ASSET_PAIR:
//    		ManageAssetPairResult manageAssetPairResult;
//    	case MANAGE_OFFER:
//    		ManageOfferResult manageOfferResult;
//    	case MANAGE_INVOICE:
//    		ManageInvoiceResult manageInvoiceResult;
//        }
//        tr;
//    default:
//        void;
//    };
//
type OperationResult struct {
	Code OperationResultCode `json:"code,omitempty"`
	Tr   *OperationResultTr  `json:"tr,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u OperationResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of OperationResult
func (u OperationResult) ArmForSwitch(sw int32) (string, bool) {
	switch OperationResultCode(sw) {
	case OperationResultCodeOpInner:
		return "Tr", true
	default:
		return "", true
	}
}

// NewOperationResult creates a new  OperationResult.
func NewOperationResult(code OperationResultCode, value interface{}) (result OperationResult, err error) {
	result.Code = code
	switch OperationResultCode(code) {
	case OperationResultCodeOpInner:
		tv, ok := value.(OperationResultTr)
		if !ok {
			err = fmt.Errorf("invalid value, must be OperationResultTr")
			return
		}
		result.Tr = &tv
	default:
		// void
	}
	return
}

// MustTr retrieves the Tr value from the union,
// panicing if the value is not set.
func (u OperationResult) MustTr() OperationResultTr {
	val, ok := u.GetTr()

	if !ok {
		panic("arm Tr is not set")
	}

	return val
}

// GetTr retrieves the Tr value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationResult) GetTr() (result OperationResultTr, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "Tr" {
		result = *u.Tr
		ok = true
	}

	return
}

// TransactionResultCode is an XDR Enum defines as:
//
//   enum TransactionResultCode
//    {
//        txSUCCESS = 0, // all operations succeeded
//
//        txFAILED = -1, // one of the operations failed (none were applied)
//
//        txTOO_EARLY = -2,         // ledger closeTime before minTime
//        txTOO_LATE = -3,          // ledger closeTime after maxTime
//        txMISSING_OPERATION = -4, // no operation was specified
//
//        txBAD_AUTH = -5,             // too few valid signatures / wrong network
//        txNO_ACCOUNT = -6,           // source account not found
//        txBAD_AUTH_EXTRA = -7,      // unused signatures attached to transaction
//        txINTERNAL_ERROR = -8,      // an unknown error occured
//    	txACCOUNT_BLOCKED = -9,     // account is blocked and cannot be source of tx
//        txDUPLICATION = -10         // if timing is stored
//    };
//
type TransactionResultCode int32

const (
	TransactionResultCodeTxSuccess          TransactionResultCode = 0
	TransactionResultCodeTxFailed           TransactionResultCode = -1
	TransactionResultCodeTxTooEarly         TransactionResultCode = -2
	TransactionResultCodeTxTooLate          TransactionResultCode = -3
	TransactionResultCodeTxMissingOperation TransactionResultCode = -4
	TransactionResultCodeTxBadAuth          TransactionResultCode = -5
	TransactionResultCodeTxNoAccount        TransactionResultCode = -6
	TransactionResultCodeTxBadAuthExtra     TransactionResultCode = -7
	TransactionResultCodeTxInternalError    TransactionResultCode = -8
	TransactionResultCodeTxAccountBlocked   TransactionResultCode = -9
	TransactionResultCodeTxDuplication      TransactionResultCode = -10
)

var TransactionResultCodeAll = []TransactionResultCode{
	TransactionResultCodeTxSuccess,
	TransactionResultCodeTxFailed,
	TransactionResultCodeTxTooEarly,
	TransactionResultCodeTxTooLate,
	TransactionResultCodeTxMissingOperation,
	TransactionResultCodeTxBadAuth,
	TransactionResultCodeTxNoAccount,
	TransactionResultCodeTxBadAuthExtra,
	TransactionResultCodeTxInternalError,
	TransactionResultCodeTxAccountBlocked,
	TransactionResultCodeTxDuplication,
}

var transactionResultCodeMap = map[int32]string{
	0:   "TransactionResultCodeTxSuccess",
	-1:  "TransactionResultCodeTxFailed",
	-2:  "TransactionResultCodeTxTooEarly",
	-3:  "TransactionResultCodeTxTooLate",
	-4:  "TransactionResultCodeTxMissingOperation",
	-5:  "TransactionResultCodeTxBadAuth",
	-6:  "TransactionResultCodeTxNoAccount",
	-7:  "TransactionResultCodeTxBadAuthExtra",
	-8:  "TransactionResultCodeTxInternalError",
	-9:  "TransactionResultCodeTxAccountBlocked",
	-10: "TransactionResultCodeTxDuplication",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for TransactionResultCode
func (e TransactionResultCode) ValidEnum(v int32) bool {
	_, ok := transactionResultCodeMap[v]
	return ok
}

// String returns the name of `e`
func (e TransactionResultCode) String() string {
	name, _ := transactionResultCodeMap[int32(e)]
	return name
}

func (e TransactionResultCode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// TransactionResultResult is an XDR NestedUnion defines as:
//
//   union switch (TransactionResultCode code)
//        {
//        case txSUCCESS:
//        case txFAILED:
//            OperationResult results<>;
//        default:
//            void;
//        }
//
type TransactionResultResult struct {
	Code    TransactionResultCode `json:"code,omitempty"`
	Results *[]OperationResult    `json:"results,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u TransactionResultResult) SwitchFieldName() string {
	return "Code"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of TransactionResultResult
func (u TransactionResultResult) ArmForSwitch(sw int32) (string, bool) {
	switch TransactionResultCode(sw) {
	case TransactionResultCodeTxSuccess:
		return "Results", true
	case TransactionResultCodeTxFailed:
		return "Results", true
	default:
		return "", true
	}
}

// NewTransactionResultResult creates a new  TransactionResultResult.
func NewTransactionResultResult(code TransactionResultCode, value interface{}) (result TransactionResultResult, err error) {
	result.Code = code
	switch TransactionResultCode(code) {
	case TransactionResultCodeTxSuccess:
		tv, ok := value.([]OperationResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be []OperationResult")
			return
		}
		result.Results = &tv
	case TransactionResultCodeTxFailed:
		tv, ok := value.([]OperationResult)
		if !ok {
			err = fmt.Errorf("invalid value, must be []OperationResult")
			return
		}
		result.Results = &tv
	default:
		// void
	}
	return
}

// MustResults retrieves the Results value from the union,
// panicing if the value is not set.
func (u TransactionResultResult) MustResults() []OperationResult {
	val, ok := u.GetResults()

	if !ok {
		panic("arm Results is not set")
	}

	return val
}

// GetResults retrieves the Results value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u TransactionResultResult) GetResults() (result []OperationResult, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Code))

	if armName == "Results" {
		result = *u.Results
		ok = true
	}

	return
}

// TransactionResultExt is an XDR NestedUnion defines as:
//
//   union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//
type TransactionResultExt struct {
	V LedgerVersion `json:"v,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u TransactionResultExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of TransactionResultExt
func (u TransactionResultExt) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerVersion(sw) {
	case LedgerVersionEmptyVersion:
		return "", true
	}
	return "-", false
}

// NewTransactionResultExt creates a new  TransactionResultExt.
func NewTransactionResultExt(v LedgerVersion, value interface{}) (result TransactionResultExt, err error) {
	result.V = v
	switch LedgerVersion(v) {
	case LedgerVersionEmptyVersion:
		// void
	}
	return
}

// TransactionResult is an XDR Struct defines as:
//
//   struct TransactionResult
//    {
//        int64 feeCharged; // actual fee charged for the transaction
//
//        union switch (TransactionResultCode code)
//        {
//        case txSUCCESS:
//        case txFAILED:
//            OperationResult results<>;
//        default:
//            void;
//        }
//        result;
//
//        // reserved for future use
//        union switch (LedgerVersion v)
//        {
//        case EMPTY_VERSION:
//            void;
//        }
//        ext;
//    };
//
type TransactionResult struct {
	FeeCharged Int64                   `json:"feeCharged,omitempty"`
	Result     TransactionResultResult `json:"result,omitempty"`
	Ext        TransactionResultExt    `json:"ext,omitempty"`
}

// Hash is an XDR Typedef defines as:
//
//   typedef opaque Hash[32];
//
type Hash [32]byte

// Uint256 is an XDR Typedef defines as:
//
//   typedef opaque uint256[32];
//
type Uint256 [32]byte

// Uint32 is an XDR Typedef defines as:
//
//   typedef unsigned int uint32;
//
type Uint32 uint32

// Int32 is an XDR Typedef defines as:
//
//   typedef int int32;
//
type Int32 int32

// Uint64 is an XDR Typedef defines as:
//
//   typedef unsigned hyper uint64;
//
type Uint64 uint64

// Int64 is an XDR Typedef defines as:
//
//   typedef hyper int64;
//
type Int64 int64

// CryptoKeyType is an XDR Enum defines as:
//
//   enum CryptoKeyType
//    {
//        KEY_TYPE_ED25519 = 0
//    };
//
type CryptoKeyType int32

const (
	CryptoKeyTypeKeyTypeEd25519 CryptoKeyType = 0
)

var CryptoKeyTypeAll = []CryptoKeyType{
	CryptoKeyTypeKeyTypeEd25519,
}

var cryptoKeyTypeMap = map[int32]string{
	0: "CryptoKeyTypeKeyTypeEd25519",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for CryptoKeyType
func (e CryptoKeyType) ValidEnum(v int32) bool {
	_, ok := cryptoKeyTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e CryptoKeyType) String() string {
	name, _ := cryptoKeyTypeMap[int32(e)]
	return name
}

func (e CryptoKeyType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// PublicKeyType is an XDR Enum defines as:
//
//   enum PublicKeyType
//    {
//    	PUBLIC_KEY_TYPE_ED25519 = KEY_TYPE_ED25519
//    };
//
type PublicKeyType int32

const (
	PublicKeyTypePublicKeyTypeEd25519 PublicKeyType = 0
)

var PublicKeyTypeAll = []PublicKeyType{
	PublicKeyTypePublicKeyTypeEd25519,
}

var publicKeyTypeMap = map[int32]string{
	0: "PublicKeyTypePublicKeyTypeEd25519",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for PublicKeyType
func (e PublicKeyType) ValidEnum(v int32) bool {
	_, ok := publicKeyTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e PublicKeyType) String() string {
	name, _ := publicKeyTypeMap[int32(e)]
	return name
}

func (e PublicKeyType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// PublicKey is an XDR Union defines as:
//
//   union PublicKey switch (CryptoKeyType type)
//    {
//    case KEY_TYPE_ED25519:
//        uint256 ed25519;
//    };
//
type PublicKey struct {
	Type    CryptoKeyType `json:"type,omitempty"`
	Ed25519 *Uint256      `json:"ed25519,omitempty"`
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u PublicKey) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of PublicKey
func (u PublicKey) ArmForSwitch(sw int32) (string, bool) {
	switch CryptoKeyType(sw) {
	case CryptoKeyTypeKeyTypeEd25519:
		return "Ed25519", true
	}
	return "-", false
}

// NewPublicKey creates a new  PublicKey.
func NewPublicKey(aType CryptoKeyType, value interface{}) (result PublicKey, err error) {
	result.Type = aType
	switch CryptoKeyType(aType) {
	case CryptoKeyTypeKeyTypeEd25519:
		tv, ok := value.(Uint256)
		if !ok {
			err = fmt.Errorf("invalid value, must be Uint256")
			return
		}
		result.Ed25519 = &tv
	}
	return
}

// MustEd25519 retrieves the Ed25519 value from the union,
// panicing if the value is not set.
func (u PublicKey) MustEd25519() Uint256 {
	val, ok := u.GetEd25519()

	if !ok {
		panic("arm Ed25519 is not set")
	}

	return val
}

// GetEd25519 retrieves the Ed25519 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u PublicKey) GetEd25519() (result Uint256, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Ed25519" {
		result = *u.Ed25519
		ok = true
	}

	return
}

// LedgerVersion is an XDR Enum defines as:
//
//   enum LedgerVersion {
//    	EMPTY_VERSION = 0,
//    	FORFEIT_RESULT_FEES = 3,
//    	IMPROVED_STATS_CALCULATION = 4,
//    	EMISSION_REQUEST_BALANCE_ID = 5,
//    	SIGNER_NAME = 6,
//    	ACCOUNT_POLICIES = 7,
//    	IMPROVED_TRANSFER_FEES_CALC = 8,
//    	USE_IMPROVED_SIGNATURE_CHECK = 9,
//    	TOKEN_REFERRAL_SHARE = 10
//    };
//
type LedgerVersion int32

const (
	LedgerVersionEmptyVersion              LedgerVersion = 0
	LedgerVersionForfeitResultFees         LedgerVersion = 3
	LedgerVersionImprovedStatsCalculation  LedgerVersion = 4
	LedgerVersionEmissionRequestBalanceId  LedgerVersion = 5
	LedgerVersionSignerName                LedgerVersion = 6
	LedgerVersionAccountPolicies           LedgerVersion = 7
	LedgerVersionImprovedTransferFeesCalc  LedgerVersion = 8
	LedgerVersionUseImprovedSignatureCheck LedgerVersion = 9
	LedgerVersionTokenReferralShare        LedgerVersion = 10
)

var LedgerVersionAll = []LedgerVersion{
	LedgerVersionEmptyVersion,
	LedgerVersionForfeitResultFees,
	LedgerVersionImprovedStatsCalculation,
	LedgerVersionEmissionRequestBalanceId,
	LedgerVersionSignerName,
	LedgerVersionAccountPolicies,
	LedgerVersionImprovedTransferFeesCalc,
	LedgerVersionUseImprovedSignatureCheck,
	LedgerVersionTokenReferralShare,
}

var ledgerVersionMap = map[int32]string{
	0:  "LedgerVersionEmptyVersion",
	3:  "LedgerVersionForfeitResultFees",
	4:  "LedgerVersionImprovedStatsCalculation",
	5:  "LedgerVersionEmissionRequestBalanceId",
	6:  "LedgerVersionSignerName",
	7:  "LedgerVersionAccountPolicies",
	8:  "LedgerVersionImprovedTransferFeesCalc",
	9:  "LedgerVersionUseImprovedSignatureCheck",
	10: "LedgerVersionTokenReferralShare",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for LedgerVersion
func (e LedgerVersion) ValidEnum(v int32) bool {
	_, ok := ledgerVersionMap[v]
	return ok
}

// String returns the name of `e`
func (e LedgerVersion) String() string {
	name, _ := ledgerVersionMap[int32(e)]
	return name
}

func (e LedgerVersion) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// Signature is an XDR Typedef defines as:
//
//   typedef opaque Signature<64>;
//
type Signature []byte

// SignatureHint is an XDR Typedef defines as:
//
//   typedef opaque SignatureHint[4];
//
type SignatureHint [4]byte

// NodeId is an XDR Typedef defines as:
//
//   typedef PublicKey NodeID;
//
type NodeId PublicKey

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u NodeId) SwitchFieldName() string {
	return PublicKey(u).SwitchFieldName()
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of PublicKey
func (u NodeId) ArmForSwitch(sw int32) (string, bool) {
	return PublicKey(u).ArmForSwitch(sw)
}

// NewNodeId creates a new  NodeId.
func NewNodeId(aType CryptoKeyType, value interface{}) (result NodeId, err error) {
	u, err := NewPublicKey(aType, value)
	result = NodeId(u)
	return
}

// MustEd25519 retrieves the Ed25519 value from the union,
// panicing if the value is not set.
func (u NodeId) MustEd25519() Uint256 {
	return PublicKey(u).MustEd25519()
}

// GetEd25519 retrieves the Ed25519 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u NodeId) GetEd25519() (result Uint256, ok bool) {
	return PublicKey(u).GetEd25519()
}

// Curve25519Secret is an XDR Struct defines as:
//
//   struct Curve25519Secret
//    {
//            opaque key[32];
//    };
//
type Curve25519Secret struct {
	Key [32]byte `json:"key,omitempty"`
}

// Curve25519Public is an XDR Struct defines as:
//
//   struct Curve25519Public
//    {
//            opaque key[32];
//    };
//
type Curve25519Public struct {
	Key [32]byte `json:"key,omitempty"`
}

// HmacSha256Key is an XDR Struct defines as:
//
//   struct HmacSha256Key
//    {
//            opaque key[32];
//    };
//
type HmacSha256Key struct {
	Key [32]byte `json:"key,omitempty"`
}

// HmacSha256Mac is an XDR Struct defines as:
//
//   struct HmacSha256Mac
//    {
//            opaque mac[32];
//    };
//
type HmacSha256Mac struct {
	Mac [32]byte `json:"mac,omitempty"`
}

// AccountId is an XDR Typedef defines as:
//
//   typedef PublicKey AccountID;
//
type AccountId PublicKey

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u AccountId) SwitchFieldName() string {
	return PublicKey(u).SwitchFieldName()
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of PublicKey
func (u AccountId) ArmForSwitch(sw int32) (string, bool) {
	return PublicKey(u).ArmForSwitch(sw)
}

// NewAccountId creates a new  AccountId.
func NewAccountId(aType CryptoKeyType, value interface{}) (result AccountId, err error) {
	u, err := NewPublicKey(aType, value)
	result = AccountId(u)
	return
}

// MustEd25519 retrieves the Ed25519 value from the union,
// panicing if the value is not set.
func (u AccountId) MustEd25519() Uint256 {
	return PublicKey(u).MustEd25519()
}

// GetEd25519 retrieves the Ed25519 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u AccountId) GetEd25519() (result Uint256, ok bool) {
	return PublicKey(u).GetEd25519()
}

// BalanceId is an XDR Typedef defines as:
//
//   typedef PublicKey BalanceID;
//
type BalanceId PublicKey

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u BalanceId) SwitchFieldName() string {
	return PublicKey(u).SwitchFieldName()
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of PublicKey
func (u BalanceId) ArmForSwitch(sw int32) (string, bool) {
	return PublicKey(u).ArmForSwitch(sw)
}

// NewBalanceId creates a new  BalanceId.
func NewBalanceId(aType CryptoKeyType, value interface{}) (result BalanceId, err error) {
	u, err := NewPublicKey(aType, value)
	result = BalanceId(u)
	return
}

// MustEd25519 retrieves the Ed25519 value from the union,
// panicing if the value is not set.
func (u BalanceId) MustEd25519() Uint256 {
	return PublicKey(u).MustEd25519()
}

// GetEd25519 retrieves the Ed25519 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u BalanceId) GetEd25519() (result Uint256, ok bool) {
	return PublicKey(u).GetEd25519()
}

// Thresholds is an XDR Typedef defines as:
//
//   typedef opaque Thresholds[4];
//
type Thresholds [4]byte

// String32 is an XDR Typedef defines as:
//
//   typedef string string32<32>;
//
type String32 string

// XDRMaxSize implements the Sized interface for String32
func (e String32) XDRMaxSize() int {
	return 32
}

// String64 is an XDR Typedef defines as:
//
//   typedef string string64<64>;
//
type String64 string

// XDRMaxSize implements the Sized interface for String64
func (e String64) XDRMaxSize() int {
	return 64
}

// String256 is an XDR Typedef defines as:
//
//   typedef string string256<256>;
//
type String256 string

// XDRMaxSize implements the Sized interface for String256
func (e String256) XDRMaxSize() int {
	return 256
}

// Longstring is an XDR Typedef defines as:
//
//   typedef string longstring<>;
//
type Longstring string

// AssetCode is an XDR Typedef defines as:
//
//   typedef string AssetCode<16>;
//
type AssetCode string

// XDRMaxSize implements the Sized interface for AssetCode
func (e AssetCode) XDRMaxSize() int {
	return 16
}

// Salt is an XDR Typedef defines as:
//
//   typedef uint64 Salt;
//
type Salt Uint64

// DataValue is an XDR Typedef defines as:
//
//   typedef opaque DataValue<64>;
//
type DataValue []byte

// OperationType is an XDR Enum defines as:
//
//   enum OperationType
//    {
//        CREATE_ACCOUNT = 0,
//        PAYMENT = 1,
//        SET_OPTIONS = 2,
//        MANAGE_COINS_EMISSION_REQUEST = 3,
//        REVIEW_COINS_EMISSION_REQUEST = 4,
//        SET_FEES = 5,
//    	MANAGE_ACCOUNT = 6,
//        FORFEIT = 7,
//        MANAGE_FORFEIT_REQUEST = 8,
//        RECOVER = 9,
//        MANAGE_BALANCE = 10,
//        REVIEW_PAYMENT_REQUEST = 11,
//        MANAGE_ASSET = 12,
//        DEMURRAGE = 13,
//        UPLOAD_PREEMISSIONS = 14,
//        SET_LIMITS = 15,
//        DIRECT_DEBIT = 16,
//    	MANAGE_ASSET_PAIR=17,
//    	MANAGE_OFFER=18,
//        MANAGE_INVOICE = 19
//    };
//
type OperationType int32

const (
	OperationTypeCreateAccount              OperationType = 0
	OperationTypePayment                    OperationType = 1
	OperationTypeSetOptions                 OperationType = 2
	OperationTypeManageCoinsEmissionRequest OperationType = 3
	OperationTypeReviewCoinsEmissionRequest OperationType = 4
	OperationTypeSetFees                    OperationType = 5
	OperationTypeManageAccount              OperationType = 6
	OperationTypeForfeit                    OperationType = 7
	OperationTypeManageForfeitRequest       OperationType = 8
	OperationTypeRecover                    OperationType = 9
	OperationTypeManageBalance              OperationType = 10
	OperationTypeReviewPaymentRequest       OperationType = 11
	OperationTypeManageAsset                OperationType = 12
	OperationTypeDemurrage                  OperationType = 13
	OperationTypeUploadPreemissions         OperationType = 14
	OperationTypeSetLimits                  OperationType = 15
	OperationTypeDirectDebit                OperationType = 16
	OperationTypeManageAssetPair            OperationType = 17
	OperationTypeManageOffer                OperationType = 18
	OperationTypeManageInvoice              OperationType = 19
)

var OperationTypeAll = []OperationType{
	OperationTypeCreateAccount,
	OperationTypePayment,
	OperationTypeSetOptions,
	OperationTypeManageCoinsEmissionRequest,
	OperationTypeReviewCoinsEmissionRequest,
	OperationTypeSetFees,
	OperationTypeManageAccount,
	OperationTypeForfeit,
	OperationTypeManageForfeitRequest,
	OperationTypeRecover,
	OperationTypeManageBalance,
	OperationTypeReviewPaymentRequest,
	OperationTypeManageAsset,
	OperationTypeDemurrage,
	OperationTypeUploadPreemissions,
	OperationTypeSetLimits,
	OperationTypeDirectDebit,
	OperationTypeManageAssetPair,
	OperationTypeManageOffer,
	OperationTypeManageInvoice,
}

var operationTypeMap = map[int32]string{
	0:  "OperationTypeCreateAccount",
	1:  "OperationTypePayment",
	2:  "OperationTypeSetOptions",
	3:  "OperationTypeManageCoinsEmissionRequest",
	4:  "OperationTypeReviewCoinsEmissionRequest",
	5:  "OperationTypeSetFees",
	6:  "OperationTypeManageAccount",
	7:  "OperationTypeForfeit",
	8:  "OperationTypeManageForfeitRequest",
	9:  "OperationTypeRecover",
	10: "OperationTypeManageBalance",
	11: "OperationTypeReviewPaymentRequest",
	12: "OperationTypeManageAsset",
	13: "OperationTypeDemurrage",
	14: "OperationTypeUploadPreemissions",
	15: "OperationTypeSetLimits",
	16: "OperationTypeDirectDebit",
	17: "OperationTypeManageAssetPair",
	18: "OperationTypeManageOffer",
	19: "OperationTypeManageInvoice",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for OperationType
func (e OperationType) ValidEnum(v int32) bool {
	_, ok := operationTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e OperationType) String() string {
	name, _ := operationTypeMap[int32(e)]
	return name
}

func (e OperationType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + e.String() + "\""), nil
}

// DecoratedSignature is an XDR Struct defines as:
//
//   struct DecoratedSignature
//    {
//        SignatureHint hint;  // last 4 bytes of the public key, used as a hint
//        Signature signature; // actual signature
//    };
//
type DecoratedSignature struct {
	Hint      SignatureHint `json:"hint,omitempty"`
	Signature Signature     `json:"signature,omitempty"`
}

var fmtTest = fmt.Sprint("this is a dummy usage of fmt")
