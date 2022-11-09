package permission

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"spydra.com/assetManagement/event"
	"spydra.com/assetManagement/metadata"
)

type Permission struct {
	AssetType string        `json:"assetType"`
	OrgID     string        `json:"forOrgId"`
	Role      []string      `json:"role"`
	CreatedAt string        `json:"createdAt"`
	UpdatedAt string        `json:"updatedAt"`
	CreatedBy metadata.User `json:"createdBy"`
	UpdatedBy metadata.User `json:"updatedBy"`
}

func (permission *Permission) CreatePermission(ctx contractapi.TransactionContextInterface) (event event.Event, err error) {
	stub := ctx.GetStub()

	if permission.AssetType == "" {
		err = fmt.Errorf("invalid asset type")
		return
	}

	if permission.OrgID == "" {
		err = fmt.Errorf("invalid org ID")
		return
	}

	err = permission.ReadPermission(ctx)
	if err == nil {
		err = fmt.Errorf("permission already exists")
		return
	}

	key, err := stub.CreateCompositeKey("permission", []string{permission.AssetType, permission.OrgID})
	if err != nil {
		return
	}

	value, err := json.Marshal(permission)
	if err != nil {
		return
	}

	err = stub.PutState(key, value)
	if err != nil {
		return
	}

	event.Type = "Permission"
	event.Message = fmt.Sprintf("Created: %s of type %s", permission.AssetType, permission.OrgID)

	return
}

func (permission *Permission) ReadPermission(ctx contractapi.TransactionContextInterface) (err error) {
	stub := ctx.GetStub()

	if permission.AssetType == "" {
		err = fmt.Errorf("invalid asset type")
		return
	}

	if permission.OrgID == "" {
		err = fmt.Errorf("invalid org ID")
		return
	}

	key, err := stub.CreateCompositeKey("permission", []string{permission.AssetType, permission.OrgID})
	if err != nil {
		return
	}

	value, err := stub.GetState(key)
	if err != nil {
		return
	}

	err = json.Unmarshal(value, permission)
	if err != nil {
		return
	}

	return
}

func (permission *Permission) UpdatePermission(ctx contractapi.TransactionContextInterface) (event event.Event, err error) {
	stub := ctx.GetStub()

	oldAsset := Permission{AssetType: permission.AssetType, OrgID: permission.OrgID}

	err = oldAsset.ReadPermission(ctx)
	if err != nil {
		return
	}

	key, err := stub.CreateCompositeKey("permission", []string{permission.AssetType, permission.OrgID})
	if err != nil {
		return
	}

	value, err := json.Marshal(permission)
	if err != nil {
		return
	}

	err = stub.PutState(key, value)
	if err != nil {
		return
	}

	event.Type = "Permission"
	event.Message = fmt.Sprintf("Updated: %s of type %s", permission.AssetType, permission.OrgID)

	return
}
