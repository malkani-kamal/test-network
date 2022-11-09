package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"spydra.com/assetManagement/asset"
	"spydra.com/assetManagement/event"
	"spydra.com/assetManagement/permission"
)

type AssetType struct {
	AssetId   string `json:"assetId"`
	AssetType string `json:"assetType"`
}

type PaginatedQueryResult struct {
	Records             string `json:"records"`
	FetchedRecordsCount int32  `json:"fetchedRecordsCount"`
	Bookmark            string `json:"bookmark"`
}

func (contract *SmartContract) CreateAssetDefinitions(ctx contractapi.TransactionContextInterface) (err error) {

	stub := ctx.GetStub()
	assetDefinitionsData := stub.GetArgs()
	assetDefinitionsData = assetDefinitionsData[1:]
	events := []event.Event{}

	for _, assetDefinitionData := range assetDefinitionsData {
		assetDefinition := &asset.AssetDefinition{}

		err = json.Unmarshal(assetDefinitionData, assetDefinition)
		if err != nil {
			return err
		}

		event, err := assetDefinition.CreateAssetDefinition(ctx)
		if err != nil {
			return err
		}

		events = append(events, event)
	}

	eventData, err := json.Marshal(events)
	if err != nil {
		return err
	}

	err = stub.SetEvent("CreateAssetDefinition", eventData)
	if err != nil {
		return err
	}

	return
}

func (contract *SmartContract) ReadAssetDefinition(ctx contractapi.TransactionContextInterface, assetType string) (assetDefinition *asset.AssetDefinition, err error) {

	assetDefinition = &asset.AssetDefinition{}
	assetDefinition.Type = assetType

	err = assetDefinition.ReadAssetDefinition(ctx)
	if err != nil {
		return
	}

	return
}

func (contract *SmartContract) UpdateAssetDefinition(ctx contractapi.TransactionContextInterface, assetDefinition asset.AssetDefinition) (err error) {

	stub := ctx.GetStub()

	event, err := assetDefinition.UpdateAssetDefinition(ctx)
	if err != nil {
		return
	}

	eventData, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = stub.SetEvent("UpdateAssetDefinition", eventData)
	if err != nil {
		return err
	}
	return
}

func (contract *SmartContract) CreateAssets(ctx contractapi.TransactionContextInterface) (err error) {

	stub := ctx.GetStub()
	assetsData := stub.GetArgs()
	assetsData = assetsData[1:]
	events := []event.Event{}

	for _, assetData := range assetsData {
		asset := &asset.Asset{}

		err = json.Unmarshal(assetData, asset)
		if err != nil {
			return err
		}

		event, err := asset.CreateAsset(ctx)
		if err != nil {
			return err
		}

		events = append(events, event)
	}

	eventData, err := json.Marshal(events)
	if err != nil {
		return err
	}

	err = stub.SetEvent("CreateAsset", eventData)
	if err != nil {
		return err
	}

	return
}

func (contract *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, assetType string, assetID string) (assetData *asset.Asset, err error) {

	assetData = &asset.Asset{}
	assetData.AssetType = assetType
	assetData.AssetId = assetID

	err = assetData.ReadAsset(ctx)
	if err != nil {
		return
	}

	return
}

func (contract *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, asset asset.Asset) (err error) {
	stub := ctx.GetStub()

	event, err := asset.UpdateAsset(ctx)
	if err != nil {
		return
	}

	eventData, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = stub.SetEvent("UpdateAsset", eventData)
	if err != nil {
		return err
	}
	return
}

func (s *SmartContract) CreatePermissions(ctx contractapi.TransactionContextInterface) (err error) {

	stub := ctx.GetStub()
	permissionData := stub.GetArgs()
	permissionData = permissionData[1:]
	events := []event.Event{}

	for _, permissionData := range permissionData {
		permission := &permission.Permission{}

		err = json.Unmarshal(permissionData, permission)
		if err != nil {
			return err
		}

		event, err := permission.CreatePermission(ctx)
		if err != nil {
			return err
		}

		events = append(events, event)
	}

	eventData, err := json.Marshal(events)
	if err != nil {
		return err
	}

	err = stub.SetEvent("CreatePermissions", eventData)
	if err != nil {
		return err
	}

	return
}

func (s *SmartContract) ReadPermission(ctx contractapi.TransactionContextInterface, assetType string, orgID string) (permissionData *permission.Permission, err error) {

	permissionData = &permission.Permission{}
	permissionData.AssetType = assetType
	permissionData.OrgID = orgID

	err = permissionData.ReadPermission(ctx)
	if err != nil {
		return
	}

	return
}

func (s *SmartContract) UpdatePermission(ctx contractapi.TransactionContextInterface, permission permission.Permission) (err error) {

	stub := ctx.GetStub()

	event, err := permission.UpdatePermission(ctx)
	if err != nil {
		return
	}

	eventData, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = stub.SetEvent("UpdatePermission", eventData)
	if err != nil {
		return err
	}
	return
}

func (t *SmartContract) GetAssetWithPagination(ctx contractapi.TransactionContextInterface, queryString string, size int, bookmark string) (*PaginatedQueryResult, error) {

	if len(queryString) == 0 {
		return nil, fmt.Errorf("provide correct query string")
	}

	if size < 1 {
		// Seting default pagesize
		size = 10
	}

	queryResults, err := getQueryResultForQueryStringWithPagination(ctx, queryString, int32(size), bookmark)
	if err != nil {
		return nil, fmt.Errorf("error occurred while getting query result and error is : %s", err.Error())
	}
	return queryResults, nil
}

func getQueryResultForQueryStringWithPagination(ctx contractapi.TransactionContextInterface, queryString string, pageSize int32, bookmark string) (*PaginatedQueryResult, error) {

	resultsIterator, responseMetadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, pageSize, bookmark)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten {
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

	return &PaginatedQueryResult{
		Records:             buffer.String(),
		FetchedRecordsCount: responseMetadata.FetchedRecordsCount,
		Bookmark:            responseMetadata.Bookmark,
	}, nil
}

func (s *SmartContract) GetAssetByQueryString(ctx contractapi.TransactionContextInterface, queryString string) (string, error) {

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		spydraLogger.Errorf("Error in fetching data from world state - %s", err.Error())
		return "", err
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			spydraLogger.Errorf("Error iterating data - %s", err.Error())
			return "", err
		}
		if bArrayMemberAlreadyWritten {
			// Add a comma before array members, suppress it for the first array member
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

	return buffer.String(), nil
}

func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, assetData string) (string, error) {

	var eventKey string

	if len(assetData) == 0 {
		spydraLogger.Errorf("Invalid doc data")
		return "", fmt.Errorf("please provide correct document data")
	}

	asset := &asset.Asset{}

	err := json.Unmarshal([]byte(assetData), &asset)
	if err != nil {
		spydraLogger.Errorf("Failed while un-marshling document. %s", err.Error())
		return "", fmt.Errorf("failed while un-marshling document. %s", err.Error())
	}

	eventKey = asset.AssetType + asset.AssetId
	return ctx.GetStub().GetTxID(), ctx.GetStub().PutState(eventKey, []byte(assetData))
}

func (s *SmartContract) GetAssetById(ctx contractapi.TransactionContextInterface, assetId string) (string, error) {

	if len(assetId) == 0 {
		spydraLogger.Errorf("Please provide correct assets Id.")
		return "", fmt.Errorf("please provide correct assets Id")
	}

	assetAsBytes, err := ctx.GetStub().GetState(assetId)

	if err != nil {
		spydraLogger.Errorf("Error fetching asset from world state. %s", err.Error())
		return "", fmt.Errorf("error fetching asset from world state. %s", err.Error())
	}

	if assetAsBytes == nil {
		spydraLogger.Errorf("%s does not exist", assetId)
		return "", fmt.Errorf("%s does not exist", assetId)
	}

	return string(assetAsBytes), nil

}

func (s *SmartContract) CheckAssetListById(ctx contractapi.TransactionContextInterface, assetIdList string) (string, error) {

	spydraLogger.Info("Assets list ", assetIdList)

	var assetIds []AssetType
	var noAssets string
	// var asset Asset

	err := json.Unmarshal([]byte(assetIdList), &assetIds)
	if err != nil {
		return "", fmt.Errorf("blockchain error while parsing asset names. %s", err.Error())
	}

	spydraLogger.Info("Printing products struct array.")

	count := 0
	for i := range assetIds {
		eventKey := assetIds[i].AssetType + assetIds[i].AssetId
		// eventKey := assetIds[i].AssetType + assetIds[i].AssetId
		spydraLogger.Info("Checking asset available with name - ", eventKey)
		// assetAsBytes, err := ctx.GetStub().GetState(assetIds[i].AssetId)
		assetAsBytes, err := ctx.GetStub().GetState(eventKey)

		if err != nil {
			spydraLogger.Errorf("Error fetching asset from world state. %s", err.Error())
			return "", fmt.Errorf("error fetching asset from world state. %s", err.Error())
		}

		if assetAsBytes == nil {
			count = 1
			// count = count + 1
			// f2pLogger.Info(count)
			noAssets = noAssets + " " + assetIds[i].AssetId
		}
		// else {
		// 	pErr := json.Unmarshal(assetAsBytes, &asset)

		// 	if pErr != nil {
		// 		f2pLogger.Errorf("Error fetching asset from world state. %s", pErr.Error())
		// 		return "", fmt.Errorf("Error fetching asset from world state. %s", pErr.Error())
		// 	}
		// 	f2pLogger.Info("asset.AssetType --" + asset.AssetType)
		// 	f2pLogger.Info("assetIds[i].AssetType --" + assetIds[i].AssetType)

		// 	if asset.AssetType != assetIds[i].AssetType {
		// 		count = 1
		// 		f2pLogger.Info("Asset Type does not match")
		// 		noAssets = noAssets + " " + assetIds[i].AssetId
		// 	}
		// }

		spydraLogger.Info("Asset found." + string(assetAsBytes))

	}

	if count > 0 {
		spydraLogger.Errorf("Either Asset Id not found or Asset Type not matching for asset id - %s.", noAssets)
		return "", fmt.Errorf("either Asset Id not found or Asset Type not matching for asset id - %s.", noAssets)
	}

	return "Success", nil
}

// func (s *SmartContract) GetHistoryForAsset(ctx contractapi.TransactionContextInterface, assetId string) (string, error) {

// 	resultsIterator, err := ctx.GetStub().GetHistoryForKey(assetId)
// 	if err != nil {
// 		spydraLogger.Errorf("Error fetching asset from world state - %s", err.Error())
// 		return "", fmt.Errorf(err.Error())
// 	}
// 	defer resultsIterator.Close()

// 	var buffer bytes.Buffer
// 	buffer.WriteString("[")

// 	bArrayMemberAlreadyWritten := false
// 	for resultsIterator.HasNext() {
// 		response, err := resultsIterator.Next()
// 		if err != nil {
// 			spydraLogger.Errorf("Error iterating results - %s", err.Error())
// 			return "", fmt.Errorf(err.Error())
// 		}
// 		if bArrayMemberAlreadyWritten {
// 			buffer.WriteString(",")
// 		}
// 		buffer.WriteString("{\"TxId\":")
// 		buffer.WriteString("\"")
// 		buffer.WriteString(response.TxId)
// 		buffer.WriteString("\"")

// 		buffer.WriteString(", \"Value\":")
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

// 	return buffer.String(), nil
// }

// func (s *SmartContract) GetAssetById(ctx contractapi.TransactionContextInterface, assetId string) (string, error) {

// 	if len(assetId) == 0 {
// 		return "", fmt.Errorf("please provide correct assets Id")
// 	}

// 	assetAsBytes, err := ctx.GetStub().GetState(assetId)

// 	if err != nil {
// 		return "", fmt.Errorf("error fetching asset from world state %s", err.Error())
// 	}

// 	if assetAsBytes == nil {
// 		return "", fmt.Errorf("%s does not exist", assetId)
// 	}

// 	return string(assetAsBytes), nil

// }
