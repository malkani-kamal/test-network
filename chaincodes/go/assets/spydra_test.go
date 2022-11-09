package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"go/f2p/mocks"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/stretchr/testify/require"
)

// go:generate counterfeiter -o mocks/transaction.go -fake-name TransactionContext . transactionContext
type transactionContext interface {
	contractapi.TransactionContextInterface
}

//go:generate counterfeiter -o mocks/chaincodestub.go -fake-name ChaincodeStub . chaincodeStub
type chaincodeStub interface {
	shim.ChaincodeStubInterface
}

func TestAddProductJourneyMapping(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	assetTransfer := SmartContract{}

	_, err := assetTransfer.AddProductJourneyMapping(transactionContext, "12345", "{\"productId\":\"7365567009123\",\"model\":[\"upstreamOrganization\",\"downstreamOrganization\"],\"upstreamOrganization\":[{\"orgName\":\"Land O Lakes\",\"orgId\":\"2\",\"orgType\":\"Supplier\",\"productId\":\"361597367381\"}],\"downstreamOrganization\":[{\"orgName\":\"Natures Harvest\",\"orgId\":\"3\",\"orgType\":\"Retailer\"}],\"createdBy\":\"anbu@paramountsoft.net\"}")
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve asset"))
	_, err = assetTransfer.AddProductJourneyMapping(transactionContext, "12345", "{\"productId\":\"7365567009123\",\"model\":[\"upstreamOrganization\",\"downstreamOrganization\"],\"upstreamOrganization\":[{\"orgName\":\"Land O Lakes\",\"orgId\":\"2\",\"orgType\":\"Supplier\",\"productId\":\"361597367381\"}],\"downstreamOrganization\":[{\"orgName\":\"Natures Harvest\",\"orgId\":\"3\",\"orgType\":\"Retailer\"}],\"createdBy\":\"anbu@paramountsoft.net\"}")
	require.EqualError(t, err, "Error fetching productJourneyMapping data - unable to retrieve asset")

	chaincodeStub.GetStateReturns([]byte{}, nil)
	_, err = assetTransfer.AddProductJourneyMapping(transactionContext, "12345", "{\"productId\":\"7365567009123\",\"model\":[\"upstreamOrganization\",\"downstreamOrganization\"],\"upstreamOrganization\":[{\"orgName\":\"Land O Lakes\",\"orgId\":\"2\",\"orgType\":\"Supplier\",\"productId\":\"361597367381\"}],\"downstreamOrganization\":[{\"orgName\":\"Natures Harvest\",\"orgId\":\"3\",\"orgType\":\"Retailer\"}],\"createdBy\":\"anbu@paramountsoft.net\"}")
	require.EqualError(t, err, "ProductJourneyMapping with id: 12345 already exist.")
}

func TestUpdateProduct(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	editProductTest := SmartContract{}

	editProductTest.UpdateProduct(transactionContext, "F111113", "{\"docType\":\"PRODUCT\",\"productId\":\"F111113\",\"productName\":\"CCCCCCYerba Mate expressograde6\",\"productDescription\":\"CCCCCC\",\"orgId\":\"1\",\"productCategory\":\"CCCCCC Food and Beverage\",\"productUOM\":\"CCCCCC\",\"productGTIN\":\"F2222\",\"productF2PID\":\"\",\"source\":\"CCCCCCfrontend\",\"consumptionGuidelines\":\"cha5nothing\",\"isActive\":0,\"version\":5,\"ingredients\":[{\"ingredientName\":\"cha5wheat floor\"}],\"nutritionalFacts\":[{\"UoM\":\"cha5KG\",\"amount\":\"5555\",\"nutrient\":\"cha5Loon\"},{\"UoM\":\"cha5gGM\",\"amount\":\"5555\",\"nutrient\":\"cha5\"}]}")
	// require.NoError(t, err)

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to update asset"))
	_, err := editProductTest.UpdateProduct(transactionContext, "F111111", "{\"docType\":\"PRODUCT\",\"productId\":\"F111111\",\"productName\":\"Milkmaid\",\"orgId\":\"1\",\"productCategory\":\"Confectionary Items\",\"productImage\":{\"fileName\":\"wheat.jpg\",\"fileURL\":\"c:kdjkjkjjda\",\"contentHash\":\"wwddweeddccccddddxxxxxxxx\"},\"productUOM\":\"KG\",\"productGTIN\":\"F111111\",\"productF2PID\":\"1\",\"source\":\"frontend\",\"consumptionGuidelines\":\"nothing\",\"isActive\":\"0\",\"ingredients\":[{\"ingredientName\":\"wheat floor\"}],\"nutritionalFacts\":[{\"UoM\":\"KG\",\"amount\":\"5\",\"nutrient\":\"Loon\"},{\"UoM\":\"GM\",\"amount\":\"5\",\"nutrient\":\"Aloo\"}],\"updatedDate\":\"12-07-2021\",\"createdDate\":\"12-07-2021\",\"createdBy\":\"anbu@paramountsoft.net\",\"updatedBy\":\"anbu@paramountsoft.net\"}")
	require.EqualError(t, err, "Error in fetching product from world state.")

	chaincodeStub.GetStateReturns([]byte{}, nil)
	_, err = editProductTest.UpdateProduct(transactionContext, "F111113", "{\"docType\":\"PRODUCT\",\"productId\":\"F111113\",\"productName\":\"CCCCCCYerba Mate expressograde6\",\"productDescription\":\"CCCCCC\",\"orgId\":\"1\",\"productCategory\":\"CCCCCC Food and Beverage\",\"productUOM\":\"CCCCCC\",\"productGTIN\":\"F2222\",\"productF2PID\":\"\",\"source\":\"CCCCCCfrontend\",\"consumptionGuidelines\":\"cha5nothing\",\"isActive\":0,\"version\":5,\"ingredients\":[{\"ingredientName\":\"cha5wheat floor\"}],\"nutritionalFacts\":[{\"UoM\":\"cha5KG\",\"amount\":\"5555\",\"nutrient\":\"cha5Loon\"},{\"UoM\":\"cha5gGM\",\"amount\":\"5555\",\"nutrient\":\"cha5\"}]}")
	require.EqualError(t, err, "Error in parsing product to be updated.")

}

func TestDeleteProduct(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	deleteProductTest := SmartContract{}

	deleteProductTest.DeleteProduct(transactionContext, "F111113")
	// require.NoError(t, a)

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to update asset"))
	err := deleteProductTest.DeleteProduct(transactionContext, "F111111")
	require.EqualError(t, err, "Error in fetching product from world state.")

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("Error in fetching product from world state."))
	err = deleteProductTest.DeleteProduct(transactionContext, "F1111")
	require.EqualError(t, err, "Error in fetching product from world state.")

}

func TestAddProduct(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	addProductTest := SmartContract{}

	err := addProductTest.AddProduct(transactionContext, "F2", "{\"docType\":\"PRODUCT\",\"productId\":\"F2\",\"productName\":\"Milkmaid\",\"orgId\":\"1\",\"productCategory\":\"Confectionary Items\",\"productImage\":{\"fileName\":\"wheat.jpg\",\"fileURL\":\"c:kdjkjkjjda\",\"contentHash\":\"wwddweeddccccddddxxxxxxxx\"},\"productUOM\":\"KG\",\"productGTIN\":\"F222221\",\"productF2PID\":\"1\",\"source\":\"frontend\",\"consumptionGuidelines\":\"nothing\",\"isActive\":\"0\",\"ingredients\":[{\"ingredientName\":\"wheat floor\"}],\"nutritionalFacts\":[{\"UoM\":\"KG\",\"amount\":\"5\",\"nutrient\":\"Loon\"},{\"UoM\":\"GM\",\"amount\":\"5\",\"nutrient\":\"Aloo\"}],\"updatedDate\":\"12-07-2021\",\"createdDate\":\"12-07-2021\",\"createdBy\":\"anbu@paramountsoft.net\",\"updatedBy\":\"anbu@paramountsoft.net\"}")
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve asset"))
	err = addProductTest.AddProduct(transactionContext, "F2", "{\"docType\":\"PRODUCT\",\"productId\":\"F2\",\"productName\":\"Milkmaid\",\"orgId\":\"1\",\"productCategory\":\"Confectionary Items\",\"productImage\":{\"fileName\":\"wheat.jpg\",\"fileURL\":\"c:kdjkjkjjda\",\"contentHash\":\"wwddweeddccccddddxxxxxxxx\"},\"productUOM\":\"KG\",\"productGTIN\":\"8765438292119\",\"productF2PID\":\"1\",\"source\":\"frontend\",\"consumptionGuidelines\":\"nothing\",\"isActive\":\"0\",\"ingredients\":[{\"ingredientName\":\"wheat floor\"}],\"nutritionalFacts\":[{\"UoM\":\"KG\",\"amount\":\"5\",\"nutrient\":\"Loon\"},{\"UoM\":\"GM\",\"amount\":\"5\",\"nutrient\":\"Aloo\"}],\"updatedDate\":\"12-07-2021\",\"createdDate\":\"12-07-2021\",\"createdBy\":\"anbu@paramountsoft.net\",\"updatedBy\":\"anbu@paramountsoft.net\"}")
	require.EqualError(t, err, "Error in fetching product from world state.")

	chaincodeStub.GetStateReturns([]byte{}, nil)
	err = addProductTest.AddProduct(transactionContext, "FPK01134", "{\"docType\":\"PRODUCT\",\"productId\":\"FPK01134\",\"productName\":\"Milkmaid\",\"orgId\":\"1\",\"productCategory\":\"Confectionary Items\",\"productImage\":\"XXXXXXXXX\",\"productUOM\":\"KG\",\"productGTIN\":\"8765438292119\",\"productF2PID\":\"1\",\"source\":\"frontend\",\"consumptionGuidelines\":\"nothing\",\"isActive\":\"0\",\"ingredients\":[{\"ingredientName\":\"wheat floor\"}],\"nutritionalFacts\":[{\"UoM\":\"KG\",\"amount\":\"5\",\"nutrient\":\"Loon\"},{\"UoM\":\"GM\",\"amount\":\"5\",\"nutrient\":\"Aloo\"}],\"updatedDate\":\"12-07-2021\",\"createdDate\":\"12-07-2021\",\"createdBy\":\"anbu@paramountsoft.net\",\"updatedBy\":\"anbu@paramountsoft.net\"}")
	require.EqualError(t, err, "Product already exist with ID - FPK01134")
}

func TestReadAsset(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	expectedAsset := ProductJourneyMapping{OrgName: "Org1"}
	bytes, err := json.Marshal(expectedAsset)
	require.NoError(t, err)

	assetTransfer := SmartContract{}

	chaincodeStub.GetStateReturns(bytes, nil)

	_, err = assetTransfer.GetAssetById(transactionContext, "")
	require.EqualError(t, err, "Please provide correct assets Id.")

	asset, err := assetTransfer.GetAssetById(transactionContext, "Org1")
	require.NoError(t, err)
	data := ProductJourneyMapping{}
	json.Unmarshal([]byte(asset), &data)

	require.Equal(t, expectedAsset.OrgName, data.OrgName)

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve asset"))
	_, err = assetTransfer.GetAssetById(transactionContext, "org1")
	require.EqualError(t, err, "Error fetching asset from world state. unable to retrieve asset")

	chaincodeStub.GetStateReturns(nil, nil)
	asset, err = assetTransfer.GetAssetById(transactionContext, "org1")
	require.EqualError(t, err, "org1 does not exist")
	require.Equal(t, asset, "")

}

func TestAddEvent(t *testing.T) {

	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	transactionEvent := SmartContract{}

	CommissionData := `{"CTE": "Commission","productId": "248144021951","processStep": "1","processLocation": "Montana","processName": "Manufacturing","processOwnerName": "Paramount software solutions",
		"createdBy": "aparna@paramountsoft.net","CommissioningKDEs": {"Item": "Bread","LGTIN": "2481440219451","SGTINStart": "248144021945101","SGTINEnd": "248144021945105","Quantity": "5","WeighPerUnit" : "",
		"UoM": "kg","OtherAttributes":{"GoodsType": "Food","GoodsType1": "Food"}},
		"When": {"Event Type": "Commission","Record Time": "2020-06-07T14:58:56.591Z"},
		"Where": {"Biz Location": "urn:epc:id:sgln:string.string.integer","Read Point": "urn:epc:id:sgln:string.string.integer"},
		"Why": {"Activity/Transaction Type": "ADD",	"Biz step": "urn:epcglobal:cbv:bizstep:Commission", "Disposition": "urn:epcglobal:cbv:disp:active"}}`

	TransformationData := `{"CTE": "Transformation","productId": "248144021950","processStep": "1","processLocation": "Montana","processName": "Manufacturing",
		"processOwnerName": "Paramount software solutions",	"createdBy": "aparna@paramountsoft.net","InputKDEs": [{	"Item": "Wheat","LGTIN": "2481440219451",
				"SGTIN": "","Quantity": "100","WeighPerUnit" : "","UoM": "KG","OtherAttributes": {"GoodsType": "Food","GoodsType1": "Food"}}],
		"OutputKDEs": {"Item": "Bread","LGTIN": "2481440219501","SGTINStart": "248144021950101","SGTINEnd": "248144021950105","Quantity": "5","WeighPerUnit" : "",
				"UoM": "kg","OtherAttributes":{"GoodsType": "Food","GoodsType1": "Food"}},
		"When": {"Event Type": "Transformation","Record Time": "2020-06-08T14:58:56.591Z"},
		"Where": {"Biz Location": "urn:epc:id:sgln:string.string.integer","Read Point": "urn:epc:id:sgln:string.string.integer"},
		"Why": {"Activity/Transaction Type": "ADD","Biz step": "urn:epcglobal:cbv:bizstep:Transformation","Disposition": "urn:epcglobal:cbv:disp:active"}
	}`

	AggregationData := `{"CTE": "Transformation","productId": "248144021950","processStep": "1","processLocation": "Montana","processName": "Manufacturing","processOwnerName": "Paramount software solutions","createdBy": "aparna@paramountsoft.net",
		"ContainerID": {"SSCC": "1248144021950"},"AggregationKDEs": [{"Item": "Bread","LGTIN": "2481440219501","OtherAttributes": {}}],
		"When": {"Event Type": "Aggregation","Record Time": "2020-06-12T14:58:56.591Z"},
		"Where": {"Biz Location": "urn:epc:id:sgln:string.string.integer","Read Point": "urn:epc:id:sgln:string.string.integer"},
		"Why": {"Activity/Transaction Type": "ADD","Biz step": "urn:epcglobal:cbv:bizstep:Aggregation","Disposition": "urn:epcglobal:cbv:disp:active"}
	}`

	DisaggregationData := `{"CTE": "Transformation","productId": "248144021950","processStep": "1","processLocation": "Montana","processName": "Manufacturing","processOwnerName": "Paramount software solutions","createdBy": "aparna@paramountsoft.net",
		"ContainerID": {"SSCC": "1248144021950"},"DisaggregationKDEs": [{"Item": "Bread","LGTIN": "2481440219501","OtherAttributes": {}}],
		"When": {"Event Type": "Disaggregation","Record Time": "2020-06-15T14:58:56.591Z"},
		"Where": {"Biz Location": "urn:epc:id:sgln:string.string.integer","Read Point": "urn:epc:id:sgln:string.string.integer"},
		"Why": {"Activity/Transaction Type": "ADD","Biz step": "urn:epcglobal:cbv:bizstep:Disaggregation","Disposition": "urn:epcglobal:cbv:disp:active"}
	}`

	ObservationData := `{"CTE": "Transformation","productId": "248144021950","processStep": "1","processLocation": "Montana","processName": "Manufacturing","processOwnerName": "Paramount software solutions","createdBy": "aparna@paramountsoft.net",
		"ObservationKDEs": {"Item": "Bread","SSCC": "1248144021950"},
		"When": {"Event Type": "Observation","Record Time": "2020-06-13T14:58:56.591Z"},
		"Where": {"Biz Location": "urn:epc:id:sgln:string.string.integer","Read Point": "urn:epc:id:sgln:string.string.integer"},
		"Why": {"Activity/Transaction Type": "ADD","Biz step": "urn:epcglobal:cbv:bizstep:Observation","Disposition": "urn:epcglobal:cbv:disp:active"}
	}`

	DecommissionData := `{"CTE": "Commission","productId": "248144021951","processStep": "1","processLocation": "Montana","processName": "Manufacturing","processOwnerName": "Paramount software solutions","createdBy": "aparna@paramountsoft.net",
		"DecommissionKDEs": {"Item": "Bread","SGTIN": "248144021950101"},
		"When": {"Event Type": "Decommission","Record Time": "2020-06-08T14:58:56.591Z"},
		"Where": {"Biz Location": "urn:epc:id:sgln:string.string.integer","Read Point": "urn:epc:id:sgln:string.string.integer"},
		"Why": {"Activity/Transaction Type": "DELETE","Biz step": "urn:epcglobal:cbv:bizstep:Decommission","Disposition": "urn:epcglobal:cbv:disp:active"}
	}`

	// Commission Event test scenario
	_, err := transactionEvent.AddEvent(transactionContext, "12345", CommissionData, "Commission", "Event")
	require.NoError(t, err)

	// Transformation Event test scenario
	_, err = transactionEvent.AddEvent(transactionContext, "123456", TransformationData, "Transformation", "Event")
	require.NoError(t, err)

	// Aggregation Event test scenario
	_, err = transactionEvent.AddEvent(transactionContext, "123456", AggregationData, "Aggregation", "Event")
	require.NoError(t, err)

	// Disaggregation Event test scenario
	_, err = transactionEvent.AddEvent(transactionContext, "123456", DisaggregationData, "Disaggregation", "Event")
	require.NoError(t, err)

	// Observation Event test scenario
	_, err = transactionEvent.AddEvent(transactionContext, "123456", ObservationData, "Observation", "Event")
	require.NoError(t, err)

	// Decommission Event test scenario
	_, err = transactionEvent.AddEvent(transactionContext, "123456", DecommissionData, "Decommission", "Event")
	require.NoError(t, err)

	// Testing with incorrect eventType
	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve asset"))
	_, err = transactionEvent.AddEvent(transactionContext, "12345", CommissionData, "Commissioning", "Event")
	require.EqualError(t, err, "Please pass correct CTE.")

}
