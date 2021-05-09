package main

import (
	txdefs "github.com/goledgerdev/cc-tools-demo/chaincode/txdefs"

	tx "github.com/goledgerdev/cc-tools/transactions"
)

var txList = []tx.Transaction{
	txdefs.GetHeader,
	txdefs.CreateNewLibrary,
	txdefs.UpdateBookTenant,
	txdefs.GetNumberOfBooksFromLibrary,
}
