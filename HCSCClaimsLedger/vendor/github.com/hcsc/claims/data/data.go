
package data

import (
	"time"
)
// const CBalancePrefix string ="CB"
// const CLimitPrefix  string ="CLI"


const ContractPrefix = "CC:"
const BalancePrefix = "CB:"
const LimitPrefix = "CLI:"
const ClaimPrefix = "CM:"
const ContractsPrefix ="CCLS:"
const ContractClaimsPrefix ="CCMLS:"
const CMemberPrefix ="CMEM:"
const LastTxnsPrefix ="LTXNSLIST:"
const LastTxnsIndPrefix ="LTXNINDLS:"
const LastTxnsFamPrefix ="LTXNFAMLS:"
const ContractMemsPrefix ="CONTMEMS:"
const AllContractPrefix = "ALLCC:"
const NoOfLastTxns int = 20


type Customer struct{
	CustomerID string `json:"CustomerID"`
	FirstName string `json:"FirstName"`
	LastName string `json:"LastName"`
	Gender string `json:"Gender"`
	DOB string `json:"DOB"`
}

type MemberRelation struct{
	SubscriberID string `json:"SubscriberID"`
	MemberID string `json:"MemberID"`
	Relationship string `json:"Relationship"`
}

type ContractDefinition struct{
	ContractID string `json:"ContractID"`
	SubscriberID string `json:"SubscriberID"`
	ContractStartDate string `json:"ContractStartDate"`
	ContractEndDate string `json:"ContractEndDate"`
	ContractType string `json:"ContractType"`
	F_Ded_Limit float64 `json:"F_Ded_Limit"`
	F_OOP_Limit float64 `json:"F_OOP_Limit"`
	Copay_feed_OOP bool `json:"Copay_feed_OOP"`
	Ded_feed_OOP bool `json:"Ded_feed_OOP"`

}

type ContractMembers struct{
	ContractID string `json:"ContractID"`
	MemberIDs []string `json:"MemberIDs"`
}
type CustomerLimit struct{
	MemberID string `json:"MemberID"`
	ContractID string `json:"ContractID"`
	I_Ded_Limit float64 `json:"I_Ded_Limit"`
	I_OOP_Limit float64 `json:"I_OOP_Limit"`

}
type CustomerBalance struct{
	MemberID string `json:"MemberID"`
	ContractID string `json:"ContractID"`
	I_Ded_Balance float64 `json:"I_Ded_Balance"`
	I_OOP_Balance float64 `json:"I_OOP_Balance"`
}

type CustomerFamilyBalance struct{
	ContractID string `json:"ContractID"`
	F_Ded_Balance float64 `json:"F_Ded_Balance"`
	F_OOP_Balance float64 `json:"F_OOP_Balance"`

}



type Claims struct{
	ClaimID string `json:"ClaimID"`
	ContractID string `json:"ContractID"`
	MemberID string `json:"MemberID"`
	SubscriberID string `json:"SubscriberID"`
	CreateDTTM time.Time `json:"CreateDTTM"`
	LastUpdatedDTTM time.Time  `json:"LastUpdatedDTTM"`
	TotalClaimAmount float64 `json:"TotalClaimAmount"`
	Status string `json:"Status"`
	Transactions []Transaction `json:"Transactions"`
	Participant string `json:"Participant"`
}

type Accum struct {
	Type string `json:"Type"`
	Amount float64 `json:"Amount"`
}

type Transaction struct {
	TransactionID int `json:"TransactionID"`
	Overage float64 `json:"Overage"`
	TotalTransactionAmount float64 `json:"TotalTransactionAmount"`
	Status string `json:"Status"`
	AccumType string `json:"AccumType"`
	AccumAmount float64 `json:"AccumAmount"`
	TransactionDate time.Time `json:"TransactionDate"`
	TxnUpdatedDate time.Time `json:"TxnUpdatedDate"`
	AccumBalance float64 `json:"AccumBalance"`

}


type CustomerAllBalance struct{
	MemberID string `json:"MemberID"`
	ContractID string `json:"ContractID"`
	I_Ded_Balance float64 `json:"I_Ded_Balance"`
	I_OOP_Balance float64 `json:"I_OOP_Balance"`
	F_Ded_Balance float64 `json:"F_Ded_Balance"`
	F_OOP_Balance float64 `json:"F_OOP_Balance"`

}

type LastTxnOfContract struct {
	ClaimID string `json:"ClaimID"`
	ContractID string `json:"ContractID"`
	SubscriberID string `json:"SubscriberID"`
	MemberID string `json:"MemberID"`
	Source string `json:"Source"`
	Transaction Transaction `json:"Transaction"`
}

// Added for UI purpose. Not to store in ledger
type SubMemRelation struct {
	SubscriberID string `json:"SubscriberID"`
	MemberID string `json:"MemberID"`
	FirstName string `json:"FirstName"`
	LastName string `json:"LastName"`
	Relationship string `json:"Relationship"`
}
type MemberNContractDetails struct {
	Contract ContractDefinition `json:"Contract"`
	I_Ded_Limit float64 `json:"I_Ded_Limit"`
	I_OOP_Limit float64 `json:"I_OOP_Limit"`
	I_Ded_Balance float64 `json:"I_Ded_Balance"`
	I_OOP_Balance float64 `json:"I_OOP_Balance"`
	F_Ded_Balance float64 `json:"F_Ded_Balance"`
	F_OOP_Balance float64 `json:"F_OOP_Balance"`
	MemberInfo Customer `json:"MemberInfo"`
	SubMemRelation []SubMemRelation `json:"SubMemRelation"`

}
