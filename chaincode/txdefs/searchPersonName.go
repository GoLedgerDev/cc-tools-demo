package txdefs

import (
	"encoding/json"
	"fmt"

	"github.com/goledgerdev/cc-tools/errors"
	sw "github.com/goledgerdev/cc-tools/stubwrapper"
	tx "github.com/goledgerdev/cc-tools/transactions"
)

// GET method
var SearchPersonName = tx.Transaction{
	Tag:         "searchPersonName",
	Label:       "Search a person name",
	Description: "Find a person by name",
	Method:      "GET",
	Callers:     []string{"$org2MSP"}, // Only org2 can call this transactions

	Args: []tx.Argument{
		{
			Tag:         "name",
			Label:       "Name",
			Description: "Name",
			DataType:    "string",
			Required:    true,
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		name, _ := req["name"].(string)

		queryString := fmt.Sprintf("{\"selector\":{\"@assetType\":\"person\", \"name\":\"%s\"}}", name)

		primaryIterator, err := stub.GetQueryResult(queryString)
		if err != nil {
			return nil, errors.NewCCError("error iterating response", 400)
		}

		defer primaryIterator.Close()

		var primaryAsset map[string]interface{}
		found := false

		for primaryIterator.HasNext() {
			queryResponse, err := primaryIterator.Next()
			if err != nil {
				return nil, errors.NewCCError("error iterating response", 400)
			}
			err = json.Unmarshal(queryResponse.Value, &primaryAsset)
			if err == nil {
				found = true
				break
			}
		}

		if !found {
			return nil, errors.NewCCError("person not found", 400)
		}

		// Marshal asset back to JSON format
		returnJSON, nerr := json.Marshal(primaryAsset)
		if nerr != nil {
			return nil, errors.WrapError(err, "failed to marshal response")
		}

		return returnJSON, nil
	},
}
