/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const getStateError = "world state get error"

type MockStub struct {
	shim.ChaincodeStubInterface
	mock.Mock
}

func (ms *MockStub) GetState(key string) ([]byte, error) {
	args := ms.Called(key)

	return args.Get(0).([]byte), args.Error(1)
}

func (ms *MockStub) PutState(key string, value []byte) error {
	args := ms.Called(key, value)

	return args.Error(0)
}

func (ms *MockStub) DelState(key string) error {
	args := ms.Called(key)

	return args.Error(0)
}

type MockContext struct {
	contractapi.TransactionContextInterface
	mock.Mock
}

func (mc *MockContext) GetStub() shim.ChaincodeStubInterface {
	args := mc.Called()

	return args.Get(0).(*MockStub)
}

func configureStub() (*MockContext, *MockStub) {
	var nilBytes []byte

	testHandler := new(Handler)
	testHandler.Value = "set value"
	handlerBytes, _ := json.Marshal(testHandler)

	ms := new(MockStub)
	ms.On("GetState", "statebad").Return(nilBytes, errors.New(getStateError))
	ms.On("GetState", "missingkey").Return(nilBytes, nil)
	ms.On("GetState", "existingkey").Return([]byte("some value"), nil)
	ms.On("GetState", "handlerkey").Return(handlerBytes, nil)
	ms.On("PutState", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).Return(nil)
	ms.On("DelState", mock.AnythingOfType("string")).Return(nil)

	mc := new(MockContext)
	mc.On("GetStub").Return(ms)

	return mc, ms
}

func TestHandlerExists(t *testing.T) {
	var exists bool
	var err error

	ctx, _ := configureStub()
	c := new(HandlerContract)

	exists, err = c.HandlerExists(ctx, "statebad")
	assert.EqualError(t, err, getStateError)
	assert.False(t, exists, "should return false on error")

	exists, err = c.HandlerExists(ctx, "missingkey")
	assert.Nil(t, err, "should not return error when can read from world state but no value for key")
	assert.False(t, exists, "should return false when no value for key in world state")

	exists, err = c.HandlerExists(ctx, "existingkey")
	assert.Nil(t, err, "should not return error when can read from world state and value exists for key")
	assert.True(t, exists, "should return true when value for key in world state")
}

func TestCreateHandler(t *testing.T) {
	var err error

	ctx, stub := configureStub()
	c := new(HandlerContract)

	err = c.CreateHandler(ctx, "statebad", "some value")
	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors")

	err = c.CreateHandler(ctx, "existingkey", "some value")
	assert.EqualError(t, err, "The asset existingkey already exists", "should error when exists returns true")

	err = c.CreateHandler(ctx, "missingkey", "some value")
	stub.AssertCalled(t, "PutState", "missingkey", []byte("{\"value\":\"some value\"}"))
}

func TestReadHandler(t *testing.T) {
	var handler *Handler
	var err error

	ctx, _ := configureStub()
	c := new(HandlerContract)

	handler, err = c.ReadHandler(ctx, "statebad")
	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors when reading")
	assert.Nil(t, handler, "should not return Handler when exists errors when reading")

	handler, err = c.ReadHandler(ctx, "missingkey")
	assert.EqualError(t, err, "The asset missingkey does not exist", "should error when exists returns true when reading")
	assert.Nil(t, handler, "should not return Handler when key does not exist in world state when reading")

	handler, err = c.ReadHandler(ctx, "existingkey")
	assert.EqualError(t, err, "Could not unmarshal world state data to type Handler", "should error when data in key is not Handler")
	assert.Nil(t, handler, "should not return Handler when data in key is not of type Handler")

	handler, err = c.ReadHandler(ctx, "handlerkey")
	expectedHandler := new(Handler)
	expectedHandler.Value = "set value"
	assert.Nil(t, err, "should not return error when Handler exists in world state when reading")
	assert.Equal(t, expectedHandler, handler, "should return deserialized Handler from world state")
}

func TestUpdateHandler(t *testing.T) {
	var err error

	ctx, stub := configureStub()
	c := new(HandlerContract)

	err = c.UpdateHandler(ctx, "statebad", "new value")
	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors when updating")

	err = c.UpdateHandler(ctx, "missingkey", "new value")
	assert.EqualError(t, err, "The asset missingkey does not exist", "should error when exists returns true when updating")

	err = c.UpdateHandler(ctx, "handlerkey", "new value")
	expectedHandler := new(Handler)
	expectedHandler.Value = "new value"
	expectedHandlerBytes, _ := json.Marshal(expectedHandler)
	assert.Nil(t, err, "should not return error when Handler exists in world state when updating")
	stub.AssertCalled(t, "PutState", "handlerkey", expectedHandlerBytes)
}

func TestDeleteHandler(t *testing.T) {
	var err error

	ctx, stub := configureStub()
	c := new(HandlerContract)

	err = c.DeleteHandler(ctx, "statebad")
	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors")

	err = c.DeleteHandler(ctx, "missingkey")
	assert.EqualError(t, err, "The asset missingkey does not exist", "should error when exists returns true when deleting")

	err = c.DeleteHandler(ctx, "handlerkey")
	assert.Nil(t, err, "should not return error when Handler exists in world state when deleting")
	stub.AssertCalled(t, "DelState", "handlerkey")
}
