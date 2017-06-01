package services// main

import (
	// "errors"
	//"fmt"
	// //"strconv"
	// "time"
	// "github.com/hyperledger/fabric/core/chaincode/shim"
	 //"github.com/hcsc/claims/data"
	// "github.com/hcsc/claims/query"
	// "encoding/json"
)

// type ClaimTxns struct {
// 	txnIDs []
// }



func main() {
/*
var contMemCarrierBalance data.ContMemCarrierBalance

contMemCarrierBalance.ContractID = "1";
carrierBalance := contMemCarrierBalance.CarrierBalances["MEDICAL"]
if (carrierBalance.CarrierName == "" ){
	var carrierBalance data.CarrierBalance
	carrierBalance.CarrierName ="MEDICAL"
	contMemCarrierBalance.CarrierBalances = make(map[string]data.CarrierBalance)
	contMemCarrierBalance.CarrierBalances["MEDICAL"] = carrierBalance

	carrierBalance.CarrierName ="DENTAL"
	contMemCarrierBalance.CarrierBalances["DENTAL"] = carrierBalance
}
carrierBalance = contMemCarrierBalance.CarrierBalances["MEDICAL"]
fmt.Println("carrierBalance : ", carrierBalance)
carrierBalance.AccumBalance  = make(map[string]float64)
accum := carrierBalance.AccumBalance["IIDED"]
if (accum == 0 ){
	fmt.Println("Accum : ", accum)

	accum = accum + 50
}
carrierBalance.AccumBalance["IIDED"] = accum


accum = carrierBalance.AccumBalance["IIOOP"]
fmt.Println("Accum : ", accum)
accum = accum + 50

// if (accum.Type== "" ){
// 	fmt.Println("Accum : ", accum)
// 	accum.Type ="IIOOP"
// 	accum.Amount = 100
// }
carrierBalance.AccumBalance["IIOOP"] = accum

contMemCarrierBalance.CarrierBalances["MEDICAL"] = carrierBalance


fmt.Println("Accum IIOOP : ", carrierBalance.AccumBalance["IIOOP"])
carrierBalance = contMemCarrierBalance.CarrierBalances["MEDICAL"]

fmt.Println("contMemCarrierBalance : ", contMemCarrierBalance)
*/



  // claimID :="CL10"
  // txnID := "tx1"
  // x := make(map[string][]string)
	//
  //  x[claimID] = append(x[claimID], txnID)
  //  txnID = "tx2"
  //  x[claimID] = append(x[claimID], txnID)
  //  claimID ="CL11"
  //  txnID = "tx3"
	//
  //  key :="CL12"
  //  value := "tx3"
  //  if (len(x) > 0){
  //    fmt.Println("Length is : ", len(x))
  //    if (x[key] != nil ){
  //      fmt.Println("x[key] is : ", x[key])
  //      txnIDs :=x[key]
  //     // if (txnIDs != nil && len(txnIDs) > 0){
  //        for index,txnID := range txnIDs{
  //          fmt.Println("x[key][index]: ",x[key][index])
  //          fmt.Println("txnID: ",txnID)
  //          if (txnID == value){
  //               fmt.Println("Found match: ",x[key][index])
  //          }
  //        }
  //      //}
  //    }
	//
  //  }
  //  if(x["CL04"] != nil ){
  //     fmt.Println("Found x[CL04]is : ", x["CL04"])
  //  }

// 	var args []string
// 	searchString :=`{"ContractID":"123","SubscriberID":"112200","MemberID":"112200","TxnStartDt":"","TxnEndDt":""}`
// 	args = append(args,searchString)
//
// 	b := []byte(searchString)
// 	var f interface{}
// 	err := json.Unmarshal(b, &f)
// 	fmt.Println("err : ",err)
// 	searchMap := f.(map[string]interface{})
//
// 	fmt.Println("searchMap : ",searchMap)
// 	contractID := searchMap["ContractID"]
// 	subscriberID := searchMap["SubscriberID"]
// 	memberID := searchMap["MemberID"]
// 	txnStartDt := searchMap["TxnStartDt"]
// 	txnEndDt := searchMap["TxnEndDt"]
//
// 	fmt.Printf("contractID %v , subscriberID %v , memberID %v ,txnStartDt %v,txnEndDt %v \n ",contractID,subscriberID,memberID,txnStartDt,txnEndDt)
// 	if (contractID != "" && subscriberID !="" && memberID !="" && txnStartDt != "" && txnEndDt != "" ){
// 		fmt.Println("Search for all Params")
// 	}else if (contractID !="" && subscriberID !="" && memberID !="" ) {
// 		fmt.Println("Search for contractID,  subscriberID, memberID")
//
// 	}else if (contractID !="" && subscriberID !="" ) {
// 		fmt.Println("Search for contractID,  subscriberID")
//
// 	}else if (contractID !="" ) {
// 		fmt.Println("Search for contractID")
//
// 	}
//
// 	for k, v := range searchMap {
// 		fmt.Println("value is  : ",v)
//     switch vv := v.(type) {
//     case string:
//         fmt.Println(k, "is string", vv)
//     case int:
//         fmt.Println(k, "is int", vv)
//     case []interface{}:
//         fmt.Println(k, "is an array:")
//         for i, u := range vv {
//             fmt.Println(i, u)
//         }
//     default:
//         fmt.Println(k, "is of a type I don't know how to handle")
//     }
// }

}
