package assettypes

import "github.com/goledgerdev/cc-tools/assets"

// Description of a book
var Book = assets.AssetType{
	Tag:         "book",
	Label:       "Book",
	Description: "",

	Props: []assets.AssetProp{
		{
			// Composite Key
			Required: true,
			IsKey:    true,
			Tag:      "title",
			Label:    "Book Title",
			DataType: "string",
			Writers:  []string{`org2MSP`}, // This means only org2 can create the asset (others can edit)
		},
		{
			// Composite Key
			Required: true,
			IsKey:    true,
			Tag:      "author",
			Label:    "Book Author",
			DataType: "string",
			Writers:  []string{`org2MSP`}, // This means only org2 can create the asset (others can edit)
		},
		{
			/// Reference to another asset
			Tag:      "currentTenant",
			Label:    "Current Tenant",
			DataType: "->person",
		},
		{
			// String list
			Tag:      "genres",
			Label:    "Genres",
			DataType: "[]string",
		},
		{
			// Date property
			Tag:      "published",
			Label:    "Publishment Date",
			DataType: "datetime",
		},
	},
}
