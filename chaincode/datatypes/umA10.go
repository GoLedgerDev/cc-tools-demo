package datatypes

import (
	"fmt"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/errors"
)

var umA10 = assets.DataType{
	AcceptedFormats: []string{"number"},
	Parse: func(data interface{}) (string, interface{}, errors.ICCError) {
		n, ok := data.(float64)
		if !ok {
			return "", nil, errors.NewCCError("invalid property", 400)
		}

		if n < 1 || n > 10 {
			return "", nil, errors.NewCCError("number must be between 1 to 10", 400)
		}

		nStr := fmt.Sprintf("%f", n)

		return nStr, n, nil
	},
}
