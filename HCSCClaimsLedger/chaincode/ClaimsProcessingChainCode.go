package main
import (
	"errors"
	"fmt"
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


func (self *ClaimsProcessingChainCode) Invoke(stub shim.ChaincodeStubInterface,
	function string, args []string) ([]byte, error) {
	fmt.Println("In Invoke with function  " + function)

	if function == "processClaim" {
		fmt.Println("invoking processClaim " + function)
		testBytes,err := services.ProcessClaim(args,stub)
		if err != nil {
			fmt.Println("Error performing ProcessClaim ")
			return nil, err
		}
		fmt.Println("Processed Claim Update successfully. ")
		return testBytes, nil
	}

	if function == "processClaimAdjust" {
		fmt.Println("invoking processClaimAdjust " + function)

		if (len(args) == 6 ){
			_,err := services.ProcessClaimAdjust(args,stub)
			if err != nil {
				fmt.Println("Error performing ProcessClaim ")
				return nil, err
			}
		}else {
			fmt.Println("Missing inputs. Need 6 input params. ")
			return nil, errors.New("Missing inputs. Need 6 input params. ")
		}
		fmt.Println("Processed Claim Update successfully. ")
		return nil, nil
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
		fmt.Println("invoking AdjustLimitsOfContract " + function)

		testBytes,err := services.AdjustLimitsOfContract(args,stub)
		if err != nil {
			fmt.Println("Error resetting  balances ")
			return nil, err
		}
		fmt.Println("Limits Adusted for ")
		return testBytes, nil
	}
	//AdjustLimitsOfContract

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
