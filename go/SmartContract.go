/*
**	
**	Smart Contract Solution
**	
**	Develop By Siva
**
*/
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	// "time"
	// // "github.com/google/uuid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the connection structure, with 4 properties.  Structure tags are used by encoding/json library
type Contract struct {
	partner_id   	string 		`json:"partner_id"`
	object  		string 		`json:"object"`
	id 				string 		`json:"id"`
	name  			string 		`json:"name"`
	partners 		string		`json:"partners"`
}


/*
 * The Init method is called when the Smart Contract "vikree" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
 func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "vikree"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
 func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) pb.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryContract" {
		return s.queryContract(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "createContract" {
		return s.createContract(APIstub, args)
	} else if function == "queryAllContracts" {
		return s.queryAllContracts(APIstub)
	} else if function == "changeContractName" {
		return s.changeContractName(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}


func (s *SmartContract) queryContract(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	contractAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(contractAsBytes)
}


func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) pb.Response {
	contracts := []Contract{
		Contract{partner_id: "a", object: "p", id: "P001-a", name: "Cisco Logicshore Plan", partners: "Logicshore a"},
		Contract{partner_id: "b", object: "s", id: "S001-b", name: "Cisco Logicshore Solution", partners: "Logicshore b"},
		Contract{partner_id: "c", object: "o", id: "O001-c", name: "Cisco Logicshore Opportunity", partners: "Logicshore c"},
		Contract{partner_id: "d", object: "p", id: "P001-d", name: "IBM Logicshore Plan", partners: "Logicshore d"},
		Contract{partner_id: "e", object: "s", id: "P001-e", name: "IBM Logicshore Solution", partners: "Logicshore e"},
	}

	i := 0
	for i < len(contracts) {
		fmt.Println("i is ", i)
		contractAsBytes, _ := json.Marshal(contracts[i])
		APIstub.PutState("CONTRACT"+strconv.Itoa(i), contractAsBytes)
		fmt.Println("Added", contracts[i])
		i = i + 1
	}

	return shim.Success(nil)
}


func (s *SmartContract) createContract(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}

	var contract = Contract{partner_id: args[1], object: args[2], id: args[3], name: args[4], partners: args[5]}

	contractAsBytes, _ := json.Marshal(contract)
	APIstub.PutState(args[0], contractAsBytes)

	return shim.Success(nil)
}


func (s *SmartContract) queryAllContracts(APIstub shim.ChaincodeStubInterface) pb.Response {

	startKey := "CONTRACT0"
	endKey := "CONTRACT99999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllContracts:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}


func (s *SmartContract) changeContractName(APIstub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	contractAsBytes, _ := APIstub.GetState(args[0])
	contract := Contract{}

	json.Unmarshal(contractAsBytes, &contract)
	contract.name = args[1]

	contractAsBytes, _ = json.Marshal(contract)
	APIstub.PutState(args[0], contractAsBytes)

	return shim.Success(nil)
}




// // ==============================================================================
// // Invoke - Contract transaction executed by this function
// // ==============================================================================
// func (SmartContract *SmartContract) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

// 	function, args := stub.GetFunctionAndParameters()
	
// 	fmt.Println("==========================================================")
// 	fmt.Println("function executing request : ", function)
	
// 	//	o******
// 	if function[0:1] == "o" {
// 		fmt.Println("==========================================================")
// 		return SmartContract.invoke(stub, function, args)
// 	}

// 	//	p******
// 	if function[0:1] == "p" {
// 		fmt.Println("==========================================================")
// 		return SmartContract.query(stub, function, args)
// 	}

// 	fmt.Println("==========================================================")

// 	return shim.Error("Received unknown function invocation - function names begin with a p or o")
// }

// //==============================================================================================================================
// //	Invoke
// //==============================================================================================================================
// func (SmartContract *SmartContract) invoke(stub shim.ChaincodeStubInterface, function string ,args []string) pb.Response {
	
// 	InvokeRequest := InvokeFunction(function)
// 	if InvokeRequest != nil {
// 		response := InvokeRequest(stub, function, args)
// 		return (response)
// 	}

// 	return shim.Error("Received unknown function invocation " + function )
// }

// //==============================================================================================================================
// //	Query	
// //==============================================================================================================================
// func (SmartContract *SmartContract) query(stub shim.ChaincodeStubInterface, function string, args []string) pb.Response {

// 	// var buff []byte
// 	var response pb.Response
// 	fmt.Println("Query() : ID Extracted and Type = ", args[0])
// 	fmt.Println("Query() : Args supplied : ", args)

// 	if len(args) < 1 {
// 		fmt.Println("Query() : Include at least 1 arguments Key ")
// 		return shim.Error("Query() : Expecting Transation type and Key value for query")
// 	}

// 	QueryRequest := QueryFunction(function)

// 	if QueryRequest != nil {
// 		response = QueryRequest(stub, function, args)
// 	} else {
// 		fmt.Println("Query() Invalid function call : ", function)
// 		response_str := "Query() : Invalid function call : " + function
// 		return shim.Error(response_str)
// 	}

// 	if response.Status != shim.OK {
// 		fmt.Println("Query() Object not found : ", args[0])
// 		response_str := "Query() : Object not found : " + args[0]
// 		return shim.Error(response_str)
// 	}

// 	return response
// }

// //////////////////////////////////////////////////////////////
// // Invoke Functions based on Function name
// //////////////////////////////////////////////////////////////
// func InvokeFunction(functionName string) func(stub shim.ChaincodeStubInterface, function string, args []string) pb.Response {
// 	InvokeFunc := map[string]func(stub shim.ChaincodeStubInterface, function string, args []string) pb.Response{
// 		"iAddPlan":	addPlan,
// 		"iUpdatePlan": updatePlan,
		
// 	}
// 	return InvokeFunc[functionName]
// }

// //////////////////////////////////////////////////////////////
// // Query Functions based on Function name
// //////////////////////////////////////////////////////////////
// func QueryFunction(fname string) func(stub shim.ChaincodeStubInterface, function string, args []string) pb.Response {
// 	QueryFunc := map[string]func(stub shim.ChaincodeStubInterface, function string, args []string) pb.Response {
// 		"qGetPlanById":  getPlanById,
// 		"qGetHistoryForPlan":  getHistoryForPlan,
// 	}
// 	return QueryFunc[fname]
// }

// //==============================================================================================================================
// //	 Add - Plan
// //==============================================================================================================================
// func addPlan(stub shim.ChaincodeStubInterface, function string, args []string) pb.Response {

// 	fmt.Printf("Invoke : Initialization of The Rebate Contract\n")
	
// 	if len(args) != 1 {
// 		return shim.Error("Incorrect number of arguments. Plan Data Not Found")
// 	}

// 	planInstance := plan{}
// 	err := json.Unmarshal([]byte(args[0]), &planInstance)
// 	if err != nil {
// 		fmt.Printf("Add Plan : error while unmarshaling plan data %s\n", err)
// 		return shim.Error(err.Error())
// 	}
	
// 	planInstance.Id 		= uuid.New().String()
// 	planInstance.createdAt 	= time.Now().UnixNano() / 1e6 
// 	planInstance.updatedAt 	= time.Now().UnixNano() / 1e6 

// 	// Marshal contract object to bytes
// 	planBytes, err := json.Marshal(planInstance)
// 	if err != nil {
// 		fmt.Printf("Add Plan : error while marshaling plan instance %s", err)
// 		return shim.Error(err.Error())
// 	}

// 	//
// 	err = stub.PutState(planInstance.Id, planBytes);
// 	if err != nil {
// 		fmt.Printf("Add Plan : error while storing plan data, due to %s", err)
// 		return shim.Error(err.Error())
// 	}

// 	//
// 	response := map[string]string{"planId": planInstance.Id}
//     responseInBytes, err := json.Marshal(response)
// 	if err != nil {
// 		fmt.Printf("Add Plan : Cannot marshall PlanId " + planInstance.Id + ": %s, getting error ", err )
// 		return shim.Error("Add Plan : Cannot marshall PlanId ")
// 	}
	
// 	//
// 	return shim.Success(nil)
// }

// //==============================================================================================================================
// //	 Update - Plan
// //==============================================================================================================================
// func updatePlan(stub shim.ChaincodeStubInterface, function string, args []string) pb.Response {

// 	fmt.Printf("Invoke : Initialization of The Rebate Contract\n")
	
// 	if len(args) != 1 {
// 		return shim.Error("Incorrect number of arguments. Plan Data Not Found")
// 	}

// 	planInstance := plan{}
// 	err := json.Unmarshal([]byte(args[0]), &planInstance)
// 	if err != nil {
// 		fmt.Printf("Add Plan : error while unmarshaling plan data %s\n", err)
// 		return shim.Error(err.Error())
// 	}
	
// 	planInstance.updatedAt 	= time.Now().UnixNano() / 1e6 	

// 	// Marshal contract object to bytes
// 	planBytes, err := json.Marshal(planInstance)
// 	if err != nil {
// 		fmt.Printf("Add Plan : error while marshaling plan instance %s", err)
// 		return shim.Error(err.Error())
// 	}

// 	//
// 	err = stub.PutState(planInstance.Id.String(), planBytes);
// 	if err != nil {
// 		fmt.Printf("Add Plan : error while updating plan %s", err)
// 		return shim.Error(err.Error())
// 	}

// 	//
// 	return shim.Success(nil)
// }


// //==============================================================================================================================
// //	 Get Plan Contract By Id
// //==============================================================================================================================
// func getPlanById(stub shim.ChaincodeStubInterface, function string, args []string) pb.Response {
// 	fmt.Println("Query : Get Plan Contract By Id")

// 	if len(args) != 1 {
// 		return shim.Error("Check there should be single argument contract details")
// 	}

// 	// initialize plan id
// 	planId := args[0]
	
// 	// Retrieving plan state from state db 
// 	planAsBytes, err := stub.GetState(planId)
// 	if err != nil {
// 		fmt.Printf("Query : Fetch Plan contract By Id, failed due to, %s", err)
// 		return shim.Error(err.Error())
// 	}

// 	return shim.Success(planAsBytes)
// }

// //==============================================================================================================================
// //	 Get Plan Contract By Id
// //==============================================================================================================================
// func getHistoryForPlan(stub shim.ChaincodeStubInterface, function string, args []string) pb.Response {
// 	if len(args) < 1 {
// 		return shim.Error("Incorrect number of arguments. Expecting 1")
// 	}

// 	planId := args[0]

// 	fmt.Printf("- start getHistoryForPlan: %s\n", planId)

// 	resultsIterator, err := stub.GetHistoryForKey(planId)
// 	if err != nil {
// 		return shim.Error(err.Error())
// 	}
// 	defer resultsIterator.Close()

// 	// buffer is a JSON array containing historic values for the marble
// 	var buffer bytes.Buffer
// 	buffer.WriteString("[")

// 	bArrayMemberAlreadyWritten := false
// 	for resultsIterator.HasNext() {
// 		response, err := resultsIterator.Next()
// 		if err != nil {
// 			return shim.Error(err.Error())
// 		}
// 		// Add a comma before array members, suppress it for the first array member
// 		if bArrayMemberAlreadyWritten == true {
// 			buffer.WriteString(",")
// 		}
// 		buffer.WriteString("{\"TxId\":")
// 		buffer.WriteString("\"")
// 		buffer.WriteString(response.TxId)
// 		buffer.WriteString("\"")

// 		buffer.WriteString(", \"Value\":")
// 		// if it was a delete operation on given key, then we need to set the
// 		//corresponding value null. Else, we will write the response.Value
// 		//as-is (as the Value itself a JSON marble)
// 		if response.IsDelete {
// 			buffer.WriteString("null")
// 		} else {
// 			buffer.WriteString(string(response.Value))
// 		}

// 		buffer.WriteString(", \"Timestamp\":")
// 		buffer.WriteString("\"")
// 		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
// 		buffer.WriteString("\"")

// 		buffer.WriteString(", \"IsDelete\":")
// 		buffer.WriteString("\"")
// 		buffer.WriteString(strconv.FormatBool(response.IsDelete))
// 		buffer.WriteString("\"")

// 		buffer.WriteString("}")
// 		bArrayMemberAlreadyWritten = true
// 	}
// 	buffer.WriteString("]")

// 	fmt.Printf("- getHistoryForPlan returning:\n%s\n", buffer.String())

// 	return shim.Success(buffer.Bytes())
// }

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
