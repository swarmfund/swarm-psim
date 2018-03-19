package airdrop

// This constants are stored in the same place in order to avoid duplications in references or Issuance causes.
//
// NEVER! change these values, it will cause to duplicate money issuance!
const (
	EarlybirdIssuanceCause      = "airdrop"
	KYCIssuanceCause            = "airdrop-for-kyc"
	MarchReferralsIssuanceCause = "airdrop-march-referrals"
	March20b20IssuanceCause     = "airdrop-march-20-20"

	// Reference suffixes must be strictly 8 symbols, because the suffix is appended to AccountID(56 length),
	// and the length of reference must be 64 (56 + 8 == 64).
	EarlybirdReferenceSuffix      = "-airdrop"
	KYCReferenceSuffix            = "-air-kyc"
	MarchReferralsReferenceSuffix = "-air-ref"
	March20b20ReferenceSuffix     = "-air-20/"
)
