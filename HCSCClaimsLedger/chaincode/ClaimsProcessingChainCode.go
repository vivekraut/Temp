package main
import (
	"errors"
	"fmt"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hcsc/claims/services"
	"github.com/hcsc/claims/query"
)
type ClaimsProcessingChainCode struct {
}

func (self *ClaimsProcessingChainCode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("In Init start ")
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	if function == "initializeCustomerContract" {
		customerBytes, err := services.InitializeCustomerContract(args,stub)
		if err != nil {
			fmt.Println("Error receiving  the Customer contract")
			return nil, err
		}
		fmt.Println("Initialization customer complete")
		return customerBytes, nil
	}
	fmt.Println("Initialization No functions found ")
	return nil, nil
}

func (self *ClaimsProcessingChainCode) Invoke(stub shim.ChaincodeStubInterface,function string, args []string) ([]byte, error) {
	fmt.Println("In Invoke with function  " + function)
	if function == "processClaim" {
		if len(args) != 7 {
			return nil, errors.New("Incorrect number of arguments. Expecting 7")
		}
		fmt.Println("invoking processClaim " + function)
		testBytes,err := services.ProcessClaim(args,stub)
		if err != nil {
			fmt.Println("Error performing ProcessClaim ")
			return nil, err
		}
		fmt.Println("Processed Claim Update successfully. ")
		return testBytes, nil
	}
	if function == "resetCustomeBalances" {
		fmt.Println("invoking resetCustomeBalances " + function)
		testBytes,err := services.ResetCustomeBalances(args,stub)
		if err != nil {
			fmt.Println("Error resetting  balances ")
			return nil, err
		}
		fmt.Println("Balnaces got reset. ")
		return testBytes, nil
	}
	if function == "adjustLimitsOfContract" {
		// if len(args) != 4 {
		// 	return nil, errors.New("Incorrect number of arguments. Expecting 4")
		// }
		fmt.Println("invoking AdjustLimitsOfContract " + function)
		testBytes,err := services.AdjustLimitsOfContract(args,stub)
		if err != nil {
			fmt.Println("Error adjusting  Limits ")
			return nil, err
		}
		fmt.Println("Limits Adusted for ")
		return testBytes, nil
	}
	if function == "updateTxnStatus" {
		// if len(args) != 6 {
		// 	return nil, errors.New("Incorrect number of arguments. Expecting 6")
		// }
		fmt.Println("invoking UpdateTxnStatus " + function)
		testBytes,err := services.UpdateTxnStatus(args,stub)
		if err != nil {
			fmt.Println("Error Updating txn status  ")
			return nil, err
		}
		fmt.Println("Updated Txn Status ")
		return testBytes, nil
	}
	if function == "createCustomer" {
		// if len(args) != 1 {
		// 	return nil, errors.New("Incorrect number of arguments. Expecting 1")
		// }
		fmt.Println("invoking CreateCustomer " + function)
		customerID :=args[0]
		testBytes,err := services.CreateCustomer(customerID,stub)
		if err != nil {
			fmt.Println("Error Creating Customer  ")
			return nil, err
		}
		fmt.Println("Customer created ")
		return testBytes, nil
	}
	if function == "createContract" {
		// if len(args) != 3 {
		// 	return nil, errors.New("Incorrect number of arguments. Expecting 3")
		// }
		fmt.Println("invoking createContract " + function)
		contractJSON := args[0]
		iDedLimit,err :=  strconv.ParseFloat(args[1], 64)
		iOOPLimit,err :=  strconv.ParseFloat(args[2], 64)
		testBytes,err := services.CreateContract(contractJSON,iDedLimit,iOOPLimit,stub)
		if err != nil {
			fmt.Println("Error Creating Contract  ")
			return nil, err
		}
		fmt.Println("Contract created ")
		return testBytes, nil
	}
	if function == "createMemberRelation" {
		// if len(args) != 2 {
		// 	return nil, errors.New("Incorrect number of arguments. Expecting 2")
		// }
		fmt.Println("invoking createMemberRelation " + function)
		subscriberID := args[0]
		memberRelationJSON := args[1]
		testBytes,err := services.CreateMemberRelation(subscriberID,memberRelationJSON,stub)
		if err != nil {
			fmt.Println("Error Creating Member relation with customer  ")
			return nil, err
		}
		fmt.Println("Member Relation created ")
		return testBytes, nil
	}
	if function == "addMemberToContract" {
		fmt.Println("invoking addMemberToContract " + function)
		contractID:= args[0]
		memberID 	:= args[1]
		iDedLimit,err :=  strconv.ParseFloat(args[2], 64)
		iOOPLimit,err :=  strconv.ParseFloat(args[3], 64)
		testBytes,err := services.AddMemberToContract(contractID,memberID,iDedLimit,iOOPLimit,stub)
		if err != nil {
			fmt.Printf("Error adding Member %v to Contract  %v ",memberID,contractID )
			return nil, err
		}
		fmt.Println("Member %v added to Contract  %v ",memberID,contractID)
		return testBytes, nil
	}

	if function == "switchPolicy" {
		fmt.Println("invoking addMemberToContract " + function)

		testBytes,err := services.SwitchPolicy(args,stub)
		if err != nil {
			fmt.Println("Error Switching policy" )
			return nil, err
		}
		fmt.Println("Member successfully Switched policy")
		return testBytes, nil
	}
	fmt.Println("invoke did not find func: " + function)
	return nil, errors.New("Received unknown function invocation: " + function)
}

func (self *ClaimsProcessingChainCode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error){
	fmt.Println("In Query with function " + function)
	bytes, err:= query.Query(stub, function,args)
	if err != nil {
		fmt.Println("Error retrieving function  ")
		return nil, err
	}
	return bytes,nil
}

func main() {
	err := shim.Start(new(ClaimsProcessingChainCode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
//
//
