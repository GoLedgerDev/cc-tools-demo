package txdefs

import (
	"encoding/json"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/errors"
	tx "github.com/goledgerdev/cc-tools/transactions"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// Updates the tenant of a Book
// POST Method
var UpdateBookTenant = tx.Transaction{
	Tag:         "updateBookTenant",
	Label:       "Update Book Tenant",
	Description: "Change the tenant of a book",
	Method:      "POST",
	Callers:     []string{`org\dMSP`}, // Any orgs can call this transaction

	Args: []tx.Argument{
		{
			Tag:         "book",
			Label:       "Book",
			Description: "Book",
			DataType:    "book",
			Required:    true,
		},
		{
			Tag:         "tenant",
			Label:       "tenant",
			Description: "New tenant of the book",
			DataType:    "person",
		},
	},
	Routine: func(stub shim.ChaincodeStubInterface, req map[string]interface{}) ([]byte, errors.ICCError) {
		bookKey, ok := req["book"].(assets.Key)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter book must be an asset")
		}
		tenantKey, ok := req["tenant"].(assets.Key)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter tenant must be an asset")
		}

		// Returns Book from channel
		bookAsset, err := bookKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "failed to get asset from the ledger")
		}
		bookMap := (map[string]interface{})(*bookAsset)
		if bookMap["@assetType"].(string) != "book" {
			return nil, errors.WrapError(err, "failed to get solicitacao from the ledger")
		}

		// Returns person from channel
		tenantAsset, err := tenantKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "failed to get asset from the ledger")
		}
		tenantMap := (map[string]interface{})(*tenantAsset)
		if tenantMap["@assetType"].(string) != "person" {
			return nil, errors.WrapError(err, "failed to get solicitacao from the ledger")
		}

		// Update data
		bookMap["tenant"] = tenantMap

		bookMap, nerr := bookAsset.Update(stub, bookMap)
		if nerr != nil {
			return nil, errors.WrapError(err, "failed to update asset")
		}

		// Marshal asset back to JSON format
		bookJSON, nerr := json.Marshal(bookAsset)
		if nerr != nil {
			return nil, errors.WrapError(err, "failed to marshal response")
		}

		return bookJSON, nil
	},
}
