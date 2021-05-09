package datatypes

import (
	"fmt"
	"math"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/errors"
)

// Book Rating (0, 1, 2, ..., 9, 10)
var bookRating = assets.DataType{
	Parse: func(data interface{}) (string, interface{}, errors.ICCError) {
		bookRating, ok := data.(float64)
		if !ok {
			return "", nil, errors.NewCCError("property must be a number", 400)
		}

		if bookRating < 0 || bookRating > 10 {
			return "", nil, errors.NewCCError("Book Rating must be 0, 1, 2, ..., 9, 10", 400)
		}

		if math.Mod(bookRating, 10) != 0 {
			return "", nil, errors.NewCCError("Book Rating must be 0, 1, 2, ..., 9, 10", 400)
		}

		return fmt.Sprintf("%d", int(bookRating)), bookRating, nil
	},
}
