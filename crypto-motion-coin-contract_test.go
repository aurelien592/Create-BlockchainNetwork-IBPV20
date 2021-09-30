/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const getStateError = "private data get error"

var transient map[string][]byte

type MockStub struct {
	shim.ChaincodeStubInterface
	mock.Mock
}

func (ms *MockStub) GetPrivateData(collection string, key string) ([]byte, error) {
	args := ms.Called(collection, key)

	return args.Get(0).([]byte), args.Error(1)
}

func (ms *MockStub) GetPrivateDataHash(collection string, key string) ([]byte, error) {
	args := ms.Called(collection, key)

	return args.Get(0).([]byte), args.Error(1)
}

func (ms *MockStub) GetTransient() (map[string][]byte, error) {

	return transient, nil
}

func (ms *MockStub) PutPrivateData(collection string, key string, value []byte) error {
	args := ms.Called(collection, key, value)

	return args.Error(0)
}

func (ms *MockStub) DelPrivateData(collection string, key string) error {
	args := ms.Called(collection, key)

	return args.Error(0)
}

type MockClientIdentity struct {
	cid.ClientIdentity
	mock.Mock
}

func (mci *MockClientIdentity) GetMSPID() (string, error) {
	args := mci.Called()
	return args.Get(0).(string), args.Error(1)
}

type MockContext struct {
	contractapi.TransactionContextInterface
	mock.Mock
}

func (mc *MockContext) GetStub() shim.ChaincodeStubInterface {
	args := mc.Called()

	return args.Get(0).(*MockStub)
}

func (mc *MockContext) GetClientIdentity() cid.ClientIdentity {
	args := mc.Called()

	return args.Get(0).(*MockClientIdentity)
}

func configureStub() (*MockContext, *MockStub) {
	var nilBytes []byte
	transient = make(map[string][]byte)

	testCryptoMotionCoin := new(CryptoMotionCoin)
	testCryptoMotionCoin.PrivateValue = "set value"
	cryptoMotionCoinBytes, _ := json.Marshal(testCryptoMotionCoin)
	hashToVerify := sha256.New()
	hashToVerify.Write(cryptoMotionCoinBytes)

	ms := new(MockStub)
	ms.On("GetPrivateData", mock.AnythingOfType("string"), "statebad").Return(nilBytes, errors.New(getStateError))
	ms.On("GetPrivateData", mock.AnythingOfType("string"), "missingkey").Return(nilBytes, nil)
	ms.On("GetPrivateData", mock.AnythingOfType("string"), "existingkey").Return([]byte("some value"), nil)
	ms.On("GetPrivateData", mock.AnythingOfType("string"), "cryptoMotionCoinkey").Return(cryptoMotionCoinBytes, nil)
	ms.On("PutPrivateData", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).Return(nil)
	ms.On("DelPrivateData", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)
	ms.On("GetPrivateDataHash", mock.AnythingOfType("string"), "statebad").Return(nilBytes, errors.New(getStateError))
	ms.On("GetPrivateDataHash", mock.AnythingOfType("string"), "missingkey").Return(nilBytes, nil)
	ms.On("GetPrivateDataHash", mock.AnythingOfType("string"), "existingkey").Return([]byte("some hash value"), nil)
	ms.On("GetPrivateDataHash", mock.AnythingOfType("string"), "cryptoMotionCoinkey").Return(hashToVerify.Sum(nil), nil)

	mci := new(MockClientIdentity)
	mci.On("GetMSPID").Return("Org1MSP", nil)

	mc := new(MockContext)
	mc.On("GetStub").Return(ms)
	mc.On("GetClientIdentity").Return(mci)

	return mc, ms
}

func TestCryptoMotionCoinExists(t *testing.T) {
	var exists bool
	var err error

	ctx, _ := configureStub()
	c := new(CryptoMotionCoinContract)

	exists, err = c.CryptoMotionCoinExists(ctx, "statebad")
	assert.EqualError(t, err, getStateError)
	assert.False(t, exists, "should return false on error")

	exists, err = c.CryptoMotionCoinExists(ctx, "missingkey")
	assert.Nil(t, err, "should not return error when can read from world state but no value for key")
	assert.False(t, exists, "should return false when no value for key in world state")

	exists, err = c.CryptoMotionCoinExists(ctx, "existingkey")
	assert.Nil(t, err, "should not return error when can read from world state and value exists for key")
	assert.True(t, exists, "should return true when value for key in world state")
}

func TestCreateCryptoMotionCoin(t *testing.T) {
	var err error

	ctx, stub := configureStub()
	c := new(CryptoMotionCoinContract)

	err = c.CreateCryptoMotionCoin(ctx, "statebad")
	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors")

	err = c.CreateCryptoMotionCoin(ctx, "existingkey")
	assert.EqualError(t, err, "The asset existingkey already exists", "should error when exists returns true")

	err = c.CreateCryptoMotionCoin(ctx, "missingkey")
	assert.EqualError(t, err, "The privateValue key was not specified in transient data. Please try again")

	transient["privateValue"] = []byte("some value")
	err = c.CreateCryptoMotionCoin(ctx, "missingkey")
	assert.Nil(t, err, "should not return error when transaction data provided")
	stub.AssertCalled(t, "PutPrivateData", "_implicit_org_Org1MSP", "missingkey", []byte("{\"privateValue\":\"some value\"}"))
}

func TestReadCryptoMotionCoin(t *testing.T) {
	var cryptoMotionCoin *CryptoMotionCoin
	var err error

	ctx, _ := configureStub()
	c := new(CryptoMotionCoinContract)

	cryptoMotionCoin, err = c.ReadCryptoMotionCoin(ctx, "statebad")
	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors when reading")
	assert.Nil(t, cryptoMotionCoin, "should not return CryptoMotionCoin when exists errors when reading")

	cryptoMotionCoin, err = c.ReadCryptoMotionCoin(ctx, "missingkey")
	assert.EqualError(t, err, "The asset missingkey does not exist", "should error when exists returns true when reading")
	assert.Nil(t, cryptoMotionCoin, "should not return CryptoMotionCoin when key does not exist in private data collection when reading")

	cryptoMotionCoin, err = c.ReadCryptoMotionCoin(ctx, "existingkey")
	assert.EqualError(t, err, "Could not unmarshal private data collection data to type CryptoMotionCoin", "should error when data in key is not CryptoMotionCoin")
	assert.Nil(t, cryptoMotionCoin, "should not return CryptoMotionCoin when data in key is not of type CryptoMotionCoin")

	cryptoMotionCoin, err = c.ReadCryptoMotionCoin(ctx, "cryptoMotionCoinkey")
	expectedCryptoMotionCoin := new(CryptoMotionCoin)
	expectedCryptoMotionCoin.PrivateValue = "set value"
	assert.Nil(t, err, "should not return error when CryptoMotionCoin exists in private data collection when reading")
	assert.Equal(t, expectedCryptoMotionCoin, cryptoMotionCoin, "should return deserialized CryptoMotionCoin from private data collection")
}

func TestUpdateCryptoMotionCoin(t *testing.T) {
	var err error

	ctx, stub := configureStub()
	c := new(CryptoMotionCoinContract)

	err = c.UpdateCryptoMotionCoin(ctx, "statebad")
	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors when updating")

	err = c.UpdateCryptoMotionCoin(ctx, "missingkey")
	assert.EqualError(t, err, "The asset missingkey does not exist", "should error when exists is false when updating")

	transient["privateValue"] = []byte("new value")
	err = c.UpdateCryptoMotionCoin(ctx, "cryptoMotionCoinkey")
	expectedCryptoMotionCoin := new(CryptoMotionCoin)
	expectedCryptoMotionCoin.PrivateValue = "new value"
	expectedCryptoMotionCoinBytes, _ := json.Marshal(expectedCryptoMotionCoin)
	assert.Nil(t, err, "should not return error when CryptoMotionCoin exists in private data collection when updating")
	stub.AssertCalled(t, "PutPrivateData", "_implicit_org_Org1MSP", "cryptoMotionCoinkey", expectedCryptoMotionCoinBytes)
}

func TestDeleteCryptoMotionCoin(t *testing.T) {
	var err error

	ctx, stub := configureStub()
	c := new(CryptoMotionCoinContract)

	err = c.DeleteCryptoMotionCoin(ctx, "statebad")
	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors")

	err = c.DeleteCryptoMotionCoin(ctx, "missingkey")
	assert.EqualError(t, err, "The asset missingkey does not exist", "should error when exists returns false when deleting")

	err = c.DeleteCryptoMotionCoin(ctx, "cryptoMotionCoinkey")
	assert.Nil(t, err, "should not return error when CryptoMotionCoin exists in private data collection when deleting")
	stub.AssertCalled(t, "DelPrivateData", "_implicit_org_Org1MSP", "cryptoMotionCoinkey")
}

func TestVerifyCryptoMotionCoin(t *testing.T) {
	var cryptoMotionCoin *CryptoMotionCoin
	var exists bool
	var err error

	ctx, stub := configureStub()
	c := new(CryptoMotionCoinContract)

	cryptoMotionCoin = new(CryptoMotionCoin)
	cryptoMotionCoin.PrivateValue = "set value"

	exists, err = c.VerifyCryptoMotionCoin(ctx, "Org1MSP", "statebad", cryptoMotionCoin)
	assert.False(t, exists, "should return false when unable to read the hash")
	assert.EqualError(t, err, getStateError)

	exists, err = c.VerifyCryptoMotionCoin(ctx, "Org1MSP", "missingkey", cryptoMotionCoin)
	assert.False(t, exists, "should return false when key does not exist")
	assert.EqualError(t, err, "No private data hash with the Key: missingkey", "should error when key does not exist")

	exists, err = c.VerifyCryptoMotionCoin(ctx, "Org1MSP", "cryptoMotionCoinkey", cryptoMotionCoin)
	assert.True(t, exists, "should return true when hash in world state matched hash from data collection")
	assert.Nil(t, err, "should not return error when hash in world state matched hash from data collection")
	stub.AssertCalled(t, "GetPrivateDataHash", "_implicit_org_Org1MSP", "cryptoMotionCoinkey")
}
