/*
**	
**	Smart Contract
**	
**	Develop By Mahendra
**
*/
package main

import (
	// "bytes"
	"encoding/jsocn"
	"fmt"
	"strconv"
	// "strings"
	"time"
	"github.com/google/uuid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

// Declaration VSmartContract
type VSmartContract struct {
}

//
// <!-- Plan -->
// Schema define attributed mapped with Plan
type plan struct {
	Id				string		`json:"id"`	
	Name      		string 			`json:"name"`
	Partners 		[]string		`json:"partners"`
	Target			float64			`json:"target"`
	Achievement 	float64			`json:"achievement"`
	Budget 			float64			`json:"budget"`
	Expenses 		float64			`json:"expenses"`
	Description 	string 			`json:"description"`
	StartDate   	int64	 		`json:"startDate"`
	EndDate     	int64  			`json:"endDate"`
	TeamMembers 	[]string		`json:"teamMembers"`
	TeamMemberRole	int				`json:"teamMemberRole"`
	CreatedAt		int64			`json:"createdAt"`			// 	time.Now().UnixNano() / 1e6
	UpdatedAt		int64			`json:"updatedAt"`
}

// uuid.NewRandom()
// ===================================================================================
// ChainCode Main Method
// ===================================================================================
func main() {
	err := shim.Start(new(VSmartContract))
	if err != nil {
		fmt.Printf("Found error, During starting VSmartContract chaincode: %s", err)
	}
}

// ==============================================================================
// Init initializes chaincode
// ==============================================================================
func (vSmartContract *VSmartContract) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// ==============================================================================
// Invoke - Contract transaction executed by this function
// ==============================================================================
func (vSmartContract *VSmartContract) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	function, args := stub.GetFunctionAndParameters()
	
	fmt.Println("==========================================================")
	fmt.Println("function executing request : ", function)
	
	//	i******
	if function[0:1] == "i" {
		fmt.Println("==========================================================")
		return VSmartContract.invoke(stub, function, args)
	}

	//	q******
	if function[0:1] == "q" {
		fmt.Println("==========================================================")
		return VSmartContract.query(stub, function, args)
	}

	fmt.Println("==========================================================")

	return shim.Error("Received unknown function invocation - function names begin with a q or i")
}

//==============================================================================================================================
//	Invoke
//==============================================================================================================================
func (vSmartContract *VSmartContract) invoke(stub shim.ChaincodeStubInterface, function string ,args []string) pb.Response {
	
	InvokeRequest := InvokeFunction(function)
	if InvokeRequest != nil {
		response := InvokeRequest(stub, function, args)
		return (response)
	}

	return shim.Error("Received unknown function invocation " + function )
}

//==============================================================================================================================
//	Query	
//==============================================================================================================================
func (vSmartContract *VSmartContract) query(stub shim.ChaincodeStubInterface, function string, args []string) pb.Response {

	// var buff []byte
	var response pb.Response
	fmt.Println("Query() : ID Extracted and Type = ", args[0])
	fmt.Println("Query() : Args supplied : ", args)

	if len(args) < 1 {
		fmt.Println("Query() : Include at least 1 arguments Key ")
		return shim.Error("Query() : Expecting Transation type and Key value for query")
	}

	QueryRequest := QueryFunction(function)

	if QueryRequest != nil {
		response = QueryRequest(stub, function, args)
	} else {
		fmt.Println("Query() Invalid function call : ", function)
		response_str := "Query() : Invalid function call : " + function
		return shim.Error(response_str)
	}

	if response.Status != shim.OK {
		fmt.Println("Query() Object not found : ", args[0])
		response_str := "Query() : Object not found : " + args[0]
		return shim.Error(response_str)
	}

	return response
}

//////////////////////////////////////////////////////////////
// Invoke Functions based on Function name
//////////////////////////////////////////////////////////////
func InvokeFunction(functionName string) func(stub shim.ChaincodeStubInterface, function string, args []string) pb.Response {
	InvokeFunc := map[string]func(stub shim.ChaincodeStubInterface, function string, args []string) pb.Response{
		"iAddPlan":	addPlan,
		"iUpdatePlan": updatePlan,
	}
	return InvokeFunc[functionName]
}

//////////////////////////////////////////////////////////////
// Query Functions based on Function name
//////////////////////////////////////////////////////////////
func QueryFunction(fname string) func(stub shim.ChaincodeStubInterface, function string, args []string) pb.Response {
	QueryFunc := map[string]func(stub shim.ChaincodeStubInterface, function string, args []string) pb.Response {
		"qGetPlanById":  getPlanById,
		"qGetHistoryForPlan":  getHistoryForPlan,
	}
	return QueryFunc[fname]
}

//==============================================================================================================================
//	 Add - Plan
//==============================================================================================================================
func addPlan(stub shim.ChaincodeStubInterface, function string, args []string) pb.Response {

	fmt.Printf("Invoke : Initialization of The Rebate Contract\n")
	
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Plan Data Not Found")
	}

	planInstance := plan{}
	err := json.Unmarshal([]byte(args[0]), &planInstance)
	if err != nil {
		fmt.Printf("Add Plan : error while unmarshaling plan data %s\n", err)
		return shim.Error(err.Error())
	}
	
	planInstance.Id 		= uuid.New().String()
	planInstance.createdAt 	= time.Now().UnixNano() / 1e6 
	planInstance.updatedAt 	= time.Now().UnixNano() / 1e6 

	// Marshal contract object to bytes
	planBytes, err := json.Marshal(planInstance)
	if err != nil {
		fmt.Printf("Add Plan : error while marshaling plan instance %s", err)
		return shim.Error(err.Error())
	}

	//
	err = stub.PutState(planInstance.Id, planBytes);
	if err != nil {
		fmt.Printf("Add Plan : error while storing plan data, due to %s", err)
		return shim.Error(err.Error())
	}

	//
	response := map[string]string{"planId": planInstance.Id}
    responseInBytes, err := json.Marshal(response)
	if err != nil {
		fmt.Printf("Add Plan : Cannot marshall PlanId " + planInstance.Id + ": %s, getting error ", err )
		return shim.Error("Add Plan : Cannot marshall PlanId ")
	}
	
	//
	return shim.Success(nil)
}

//==============================================================================================================================
//	 Update - Plan
//==============================================================================================================================
func updatePlan(stub shim.ChaincodeStubInterface, function string, args []string) pb.Response {

	fmt.Printf("Invoke : Initialization of The Rebate Contract\n")
	
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Plan Data Not Found")
	}

	planInstance := plan{}
	err := json.Unmarshal([]byte(args[0]), &planInstance)
	if err != nil {
		fmt.Printf("Add Plan : error while unmarshaling plan data %s\n", err)
		return shim.Error(err.Error())
	}
	
	planInstance.updatedAt 	= time.Now().UnixNano() / 1e6 	

	// Marshal contract object to bytes
	planBytes, err := json.Marshal(planInstance)
	if err != nil {
		fmt.Printf("Add Plan : error while marshaling plan instance %s", err)
		return shim.Error(err.Error())
	}

	//
	err = stub.PutState(planInstance.Id.String(), planBytes);
	if err != nil {
		fmt.Printf("Add Plan : error while updating plan %s", err)
		return shim.Error(err.Error())
	}

	//
	return shim.Success(nil)
}


//==============================================================================================================================
//	 Get Plan Contract By Id
//==============================================================================================================================
func getPlanById(stub shim.ChaincodeStubInterface, function string, args []string) pb.Response {
	fmt.Println("Query : Get Plan Contract By Id")

	if len(args) != 1 {
		return shim.Error("Check there should be single argument contract details")
	}

	// initialize plan id
	planId := args[0]
	
	// Retrieving plan state from state db 
	planAsBytes, err := stub.GetState(planId)
	if err != nil {
		fmt.Printf("Query : Fetch Plan contract By Id, failed due to, %s", err)
		return shim.Error(err.Error())
	}

	return shim.Success(planAsBytes)
}

//==============================================================================================================================
//	 Get Plan Contract By Id
//==============================================================================================================================
func getHistoryForPlan(stub shim.ChaincodeStubInterface, function string, args []string) pb.Response {
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	planId := args[0]

	fmt.Printf("- start getHistoryForPlan: %s\n", planId)

	resultsIterator, err := stub.GetHistoryForKey(planId)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the marble
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON marble)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getHistoryForPlan returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}
