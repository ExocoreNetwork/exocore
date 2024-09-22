package types

type DeltaDelegationAmounts DelegationAmounts

const (
	NotBondedPoolName = "not_bonded_tokens_pool"
	BondedPoolName    = "bonded_tokens_pool"
	// TODO: we currently not support redelegation, and operators is not directly related to bonded(need to optIn first), so we use this pool name for now
	DelegatedPoolName = "delegated_tokens_pool"
)
