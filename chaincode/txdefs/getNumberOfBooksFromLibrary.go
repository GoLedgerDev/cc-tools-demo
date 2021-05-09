package txdefs

import (
	"encoding/json"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/errors"
	tx "github.com/goledgerdev/cc-tools/transactions"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// Return the number of books of a library
// GET method
var GetNumberOfBooksFromLibrary = tx.Transaction{
	Tag:         "getNumberOfBooksFromLibrary",
	Label:       "Get Number Of Books From Library",
	Description: "Return the number of books of a library",
	Method:      "GET",
	Callers:     []string{"$org2MSP"}, // Only org2 can call this transactions

	Args: []tx.Argument{
		{
			Tag:         "library",
			Label:       "Library",
			Description: "Library",
			DataType:    "library",
			Required:    true,
		},
	},
	Routine: func(stub shim.ChaincodeStubInterface, req map[string]interface{}) ([]byte, errors.ICCError) {
		libraryKey, ok := req["library"].(assets.Key)
		if !ok {
			return nil, errors.WrapError(nil, "Parameter library must be an asset")
		}

		// Returns Library from channel
		libraryAsset, err := libraryKey.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "failed to get asset from the ledger")
		}
		libraryMap := (map[string]interface{})(*libraryAsset)
		if libraryMap["@assetType"].(string) != "library" {
			return nil, errors.WrapError(err, "failed to get library from the ledger")
		}

		numberOfBooks := 0
		books, ok := libraryMap["books"].([]interface{})
		if ok {
			numberOfBooks = len(books)
		}

		var returnMap map[string]interface{}
		returnMap["numberOfBooks"] = numberOfBooks

		// Marshal asset back to JSON format
		returnJSON, nerr := json.Marshal(returnMap)
		if nerr != nil {
			return nil, errors.WrapError(err, "failed to marshal response")
		}

		return returnJSON, nil
	},
}
