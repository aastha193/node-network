package main

import (
        "crypto/ecdsa"
        "crypto/sha256"
        "crypto/x509"
        "encoding/json"
        "errors"
        "fmt"
        "math/big"

        "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

//SmartContract provides functions
type SmartContract struct {
        contractapi.Contract
}


func (sc *SmartContract) Init(ctx contractapi.TransactionContextInterface, chainid string, chaincodename string) (string, error) {
        ctx.GetStub().PutState("test", []byte("0"))
        return "", nil
}

// Patient is
type Patient struct {
        ID     string `json:"id"`
        Name   string `json:"name"`
        DOB    int  `json:"dateOfBirth"`
        Gender string `json:"gender"`
        docs  []byte `json:"hash"`
}

//AddPatient ...
func (sc *SmartContract) AddPatient(ctx contractapi.TransactionContextInterface, patientJSON string) (string, error) {
        patientBytes := []byte(patientJSON)

        patient := Patient{}
        err = json.Unmarshal(patientBytes, &patient)
        if err != nil {
                return "", errors.New("patient format is not valid: " + err.Error())
        }

        compositeKey, compositeErr := ctx.GetStub().CreateCompositeKey("patientData", patient.ID)
        if compositeErr != nil {
                return "", errors.New("Could not create a composite key for " + patient.ID + ": " + compositeErr.Error())
        }

        compositePutErr := ctx.GetStub().PutState(compositeKey, patientBytes)
        if compositePutErr != nil {
                return "", errors.New("Could not insert " + patient.ID + "in the ledger: " + compositePutErr.Error())
        }

        return "Successfully added Patient to the ledger", nil
}


//editPatient ...
func (sc *SmartContract) editPatient(ctx contractapi.TransactionContextInterface, patientJSON string, patientId string) (string, error) {
        patientNewBytes := []byte(patientJSON)
    
        patientNew := Patient{}
        err = json.Unmarshal(patientNewBytes, &patientNew)
        if err != nil {
                return "", errors.New("patient format is not valid: " + err.Error())
        }

        patientbytes, err := sc.getRecord(ctx, []string{"patientData", patientId})
        if err != nil {
                return "", errors.New("error in fetching patient record with patientID " + patientId + ": " + err.Error())
        }
        if patientbytes == nil {
                return "", errors.New("patient with id " + patientId + "not found")
        }

        patient := Patient{}
        err = json.Unmarshal(patientbytes, &patient)
        if err != nil {
                return "", errors.New("FATAL: Error in database content, corrupted data.. Error unmarshal: " + err.Error())
        }

        compositeKey, compositeErr := ctx.GetStub().CreateCompositeKey("patientData", patientNew.ID)
        if compositeErr != nil {
                return "", errors.New("Could not create a composite key for " + patientNew.ID + ": " + compositeErr.Error())
        }

        compositePutErr := ctx.GetStub().PutState(compositeKey, patientNewBytes)
        if compositePutErr != nil {
                return "", errors.New("Could not insert " + patientNew.ID + "in the ledger: " + compositePutErr.Error())
        }

        return "Successfully updated Patient to the ledger", nil
}

func (sc *SmartContract) GetPatient(ctx contractapi.TransactionContextInterface, id string) (string, error) {
        userID := id
        record, err := sc.getRecord(ctx, []string{"patientData", userID})
        if err != nil {
                return "", errors.New("Error is " + err.Error())
        }
        return string(record), nil
}

func (sc *SmartContract) getRecord(ctx contractapi.TransactionContextInterface, infos []string) ([]byte, error) {
        result, err := ctx.GetStub().GetStateByPartialCompositeKey(infos)
        if err != nil {
                return nil, fmt.Errorf("Could not retrieve value for %v: %v", infos, err)
        }
        defer result.Close()

        // Check the variable existed
        if !result.HasNext() {
                return nil, nil
        }

        value, nextErr := result.Next()
        if nextErr != nil {
                return nil, nextErr
        }
        return value.GetValue(), nil
}


func main() {
        fmt.Println("Inside main...")
        chaincode, err := contractapi.NewChaincode(new(SmartContract))

        if err != nil {
                fmt.Printf("Error create mycc chaincode: %s", err.Error())
                return
        }

        if err := chaincode.Start(); err != nil {
                fmt.Printf("Error starting mycc chaincode: %s", err.Error())
        }
}

