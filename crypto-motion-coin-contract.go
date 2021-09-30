/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// CryptoMotionCoinContract contract for managing CRUD for CryptoMotionCoin
type CryptoMotionCoinContract struct {
	contractapi.Contract
}

func getCollectionName(ctx contractapi.TransactionContextInterface) (string, error) {
	mspid, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", err
	}

	collectionName := "_implicit_org_" + mspid

	return collectionName, nil
}

// CryptoMotionCoinExists returns true when asset with given ID exists in private data collection
func (c *CryptoMotionCoinContract) CryptoMotionCoinExists(ctx contractapi.TransactionContextInterface, cryptoMotionCoinID string) (bool, error) {
	collectionName, collectionNameErr := getCollectionName(ctx)
	if collectionNameErr != nil {
		return false, collectionNameErr
	}

	data, err := ctx.GetStub().GetPrivateDataHash(collectionName, cryptoMotionCoinID)

	if err != nil {
		return false, err
	}

	return data != nil, nil
}

// CreateCryptoMotionCoin creates a new instance of CryptoMotionCoin
func (c *CryptoMotionCoinContract) CreateCryptoMotionCoin(ctx contractapi.TransactionContextInterface, cryptoMotionCoinID string) error {
	exists, err := c.CryptoMotionCoinExists(ctx, cryptoMotionCoinID)
	if err != nil {
		return fmt.Errorf("Could not read from world state. %s", err)
	} else if exists {
		return fmt.Errorf("The asset %s already exists", cryptoMotionCoinID)
	}

	cryptoMotionCoin := new(CryptoMotionCoin)

	transientData, _ := ctx.GetStub().GetTransient()

	privateValue, exists := transientData["privateValue"]

	if len(transientData) == 0 || !exists {
		return fmt.Errorf("The privateValue key was not specified in transient data. Please try again")
	}

	cryptoMotionCoin.PrivateValue = string(privateValue)

	bytes, _ := json.Marshal(cryptoMotionCoin)

	collectionName, collectionNameErr := getCollectionName(ctx)
	if collectionNameErr != nil {
		return collectionNameErr
	}

	return ctx.GetStub().PutPrivateData(collectionName, cryptoMotionCoinID, bytes)
}

// ReadCryptoMotionCoin retrieves an instance of CryptoMotionCoin from the private data collection
func (c *CryptoMotionCoinContract) ReadCryptoMotionCoin(ctx contractapi.TransactionContextInterface, cryptoMotionCoinID string) (*CryptoMotionCoin, error) {
	exists, err := c.CryptoMotionCoinExists(ctx, cryptoMotionCoinID)
	if err != nil {
		return nil, fmt.Errorf("Could not read from world state. %s", err)
	} else if !exists {
		return nil, fmt.Errorf("The asset %s does not exist", cryptoMotionCoinID)
	}

	collectionName, collectionNameErr := getCollectionName(ctx)
	if collectionNameErr != nil {
		return nil, collectionNameErr
	}

	bytes, _ := ctx.GetStub().GetPrivateData(collectionName, cryptoMotionCoinID)

	cryptoMotionCoin := new(CryptoMotionCoin)

	err = json.Unmarshal(bytes, cryptoMotionCoin)

	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal private data collection data to type CryptoMotionCoin")
	}

	return cryptoMotionCoin, nil
}

// UpdateCryptoMotionCoin retrieves an instance of CryptoMotionCoin from the private data collection and updates its value
func (c *CryptoMotionCoinContract) UpdateCryptoMotionCoin(ctx contractapi.TransactionContextInterface, cryptoMotionCoinID string) error {
	exists, err := c.CryptoMotionCoinExists(ctx, cryptoMotionCoinID)
	if err != nil {
		return fmt.Errorf("Could not read from world state. %s", err)
	} else if !exists {
		return fmt.Errorf("The asset %s does not exist", cryptoMotionCoinID)
	}

	transientData, _ := ctx.GetStub().GetTransient()
	newValue, exists := transientData["privateValue"]

	if len(transientData) == 0 || !exists {
		return fmt.Errorf("The privateValue key was not specified in transient data. Please try again")
	}

	cryptoMotionCoin := new(CryptoMotionCoin)
	cryptoMotionCoin.PrivateValue = string(newValue)

	bytes, _ := json.Marshal(cryptoMotionCoin)

	collectionName, collectionNameErr := getCollectionName(ctx)
	if collectionNameErr != nil {
		return collectionNameErr
	}

	return ctx.GetStub().PutPrivateData(collectionName, cryptoMotionCoinID, bytes)
}

// DeleteCryptoMotionCoin deletes an instance of CryptoMotionCoin from the private data collection
func (c *CryptoMotionCoinContract) DeleteCryptoMotionCoin(ctx contractapi.TransactionContextInterface, cryptoMotionCoinID string) error {
	exists, err := c.CryptoMotionCoinExists(ctx, cryptoMotionCoinID)
	if err != nil {
		return fmt.Errorf("Could not read from world state. %s", err)
	} else if !exists {
		return fmt.Errorf("The asset %s does not exist", cryptoMotionCoinID)
	}

	collectionName, collectionNameErr := getCollectionName(ctx)
	if collectionNameErr != nil {
		return collectionNameErr
	}

	return ctx.GetStub().DelPrivateData(collectionName, cryptoMotionCoinID)
}

// VerifyCryptoMotionCoin verifies the hash for an instance of CryptoMotionCoin from the private data collection matches the hash stored in the public ledger //FIXME check this
func (c *CryptoMotionCoinContract) VerifyCryptoMotionCoin(ctx contractapi.TransactionContextInterface, mspid string, cryptoMotionCoinID string, objectToVerify *CryptoMotionCoin) (bool, error) {
	bytes, _ := json.Marshal(objectToVerify)
	hashToVerify := sha256.New()
	hashToVerify.Write(bytes)

	pdHashBytes, err := ctx.GetStub().GetPrivateDataHash("_implicit_org_" + mspid, cryptoMotionCoinID)
	if err != nil {
		return false, err
	} else if len(pdHashBytes) == 0 {
		return false, fmt.Errorf("No private data hash with the Key: %s", cryptoMotionCoinID)
	}

	return hex.EncodeToString(hashToVerify.Sum(nil)) == hex.EncodeToString(pdHashBytes), nil
}
