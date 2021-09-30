/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-contract-api-go/metadata"
)

func main() {
	cryptoMotionCoinContract := new(CryptoMotionCoinContract)
	cryptoMotionCoinContract.Info.Version = "0.0.1"
	cryptoMotionCoinContract.Info.Description = "My Private Data Smart Contract"
	cryptoMotionCoinContract.Info.License = new(metadata.LicenseMetadata)
	cryptoMotionCoinContract.Info.License.Name = "Apache-2.0"
	cryptoMotionCoinContract.Info.Contact = new(metadata.ContactMetadata)
	cryptoMotionCoinContract.Info.Contact.Name = "John Doe"

	chaincode, err := contractapi.NewChaincode(cryptoMotionCoinContract)
	chaincode.Info.Title = "Create-BlockchainNetwork-IBPV20 chaincode"
	chaincode.Info.Version = "0.0.1"

	if err != nil {
		panic("Could not create chaincode from CryptoMotionCoinContract." + err.Error())
	}

	err = chaincode.Start()

	if err != nil {
		panic("Failed to start chaincode. " + err.Error())
	}
}
