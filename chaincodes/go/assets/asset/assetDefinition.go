package asset

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"spydra.com/assetManagement/event"
	"spydra.com/assetManagement/metadata"
)

type AssetDefinition struct {
	IdAttribute string                `json:"idAttribute"`
	OwnerOrg    string                `json:"ownerOrgId"`
	DocType     string                `json:"docType,omitempty" metadata:",optional"`
	Type        string                `json:"assetType"`
	IsActive    bool                  `json:"isActive"`
	References  []ReferenceDefinition `json:"references,omitempty" metadata:",optional"`
	CreatedAt   string                `json:"createdAt"`
	UpdatedAt   string                `json:"updatedAt"`
	CreatedBy   metadata.User         `json:"createdBy"`
	UpdatedBy   metadata.User         `json:"updatedBy"`
}

type ReferenceDefinition struct {
	IdAttribute      string `json:"idAttribute"`
	TypeAttribute    string `json:"typeAttribute"`
	EnforceIntegrity bool   `json:"enforceIntegrity"`
}

func (assetDefinition *AssetDefinition) CreateAssetDefinition(ctx contractapi.TransactionContextInterface) (event event.Event, err error) {
	stub := ctx.GetStub()

	if assetDefinition.Type == "" {
		err = fmt.Errorf("invalid asset definition type")
		return
	}

	err = assetDefinition.ReadAssetDefinition(ctx)
	if err == nil {
		err = fmt.Errorf("asset definition already exists")
		return
	}

	key, err := stub.CreateCompositeKey("assetDefinition", []string{assetDefinition.Type})
	if err != nil {
		return
	}

	value, err := json.Marshal(assetDefinition)
	if err != nil {
		return
	}

	err = stub.PutState(key, value)
	if err != nil {
		return
	}

	event.Type = "CreateAssetDefinition"
	event.Message = fmt.Sprintf("Created: %s", assetDefinition.Type)

	return
}

func (assetDefinition *AssetDefinition) ReadAssetDefinition(ctx contractapi.TransactionContextInterface) (err error) {
	stub := ctx.GetStub()

	if assetDefinition.Type == "" {
		err = fmt.Errorf("invalid asset definition type")
		return
	}

	key, err := stub.CreateCompositeKey("assetDefinition", []string{assetDefinition.Type})
	if err != nil {
		return
	}

	value, err := stub.GetState(key)
	if err != nil {
		return
	}

	err = json.Unmarshal(value, assetDefinition)
	if err != nil {
		return
	}

	if assetDefinition.Type == "" {
		err = fmt.Errorf("invalid asset definition")
		return
	}

	return
}

func (assetDefinition *AssetDefinition) UpdateAssetDefinition(ctx contractapi.TransactionContextInterface) (event event.Event, err error) {
	stub := ctx.GetStub()

	oldAssetDefinition := AssetDefinition{Type: assetDefinition.Type}

	err = oldAssetDefinition.ReadAssetDefinition(ctx)
	if err != nil {
		return
	}

	key, err := stub.CreateCompositeKey("assetDefinition", []string{assetDefinition.Type})
	if err != nil {
		return
	}

	value, err := json.Marshal(assetDefinition)
	if err != nil {
		return
	}

	err = stub.PutState(key, value)
	if err != nil {
		return
	}

	event.Type = "UpdateAssetDefinition"
	event.Message = fmt.Sprintf("Updated: %s", assetDefinition.Type)

	return
}
