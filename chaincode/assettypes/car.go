package assettypes

import (
	"github.com/goledgerdev/cc-tools/assets"
)

var Car = assets.AssetType{
	Tag:         "car",
	Label:       "Car",
	Description: "Car",

	Props: []assets.AssetProp{
		{
			// Primary key
			Required: true,
			IsKey:    true,
			Tag:      "id",
			Label:    "ID Car",
			DataType: "number", // Datatypes are identified at datatypes folder
		},
		{
			// Mandatory property
			Required: true,
			Tag:      "color",
			Label:    "Color",
			DataType: "string",
		},
		{
			Tag:      "stars",
			Label:    "Stars",
			DataType: "umA10",
		},
	},
}
