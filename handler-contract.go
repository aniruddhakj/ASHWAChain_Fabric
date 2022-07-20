/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// HandlerContract contract for managing CRUD for Handler
type HandlerContract struct {
	contractapi.Contract
}

// HandlerExists returns true when asset with given ID exists in world state
func (c *HandlerContract) HandlerExists(ctx contractapi.TransactionContextInterface, handlerID string) (bool, error) {
	data, err := ctx.GetStub().GetState(handlerID)

	if err != nil {
		return false, err
	}

	return data != nil, nil
}

// CreateHandler creates a new instance of Handler
func (c *HandlerContract) CreateHandler(ctx contractapi.TransactionContextInterface, handlerID string, value string) error {
	exists, err := c.HandlerExists(ctx, handlerID)
	if err != nil {
		return fmt.Errorf("Could not read from world state. %s", err)
	} else if exists {
		return fmt.Errorf("The asset %s already exists", handlerID)
	}

	handler := new(Handler)
	handler.Value = value

	bytes, _ := json.Marshal(handler)

	return ctx.GetStub().PutState(handlerID, bytes)
}

// ReadHandler retrieves an instance of Handler from the world state
func (c *HandlerContract) ReadHandler(ctx contractapi.TransactionContextInterface, handlerID string) (*Handler, error) {
	exists, err := c.HandlerExists(ctx, handlerID)
	if err != nil {
		return nil, fmt.Errorf("Could not read from world state. %s", err)
	} else if !exists {
		return nil, fmt.Errorf("The asset %s does not exist", handlerID)
	}

	bytes, _ := ctx.GetStub().GetState(handlerID)

	handler := new(Handler)

	err = json.Unmarshal(bytes, handler)

	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal world state data to type Handler")
	}

	return handler, nil
}

// UpdateHandler retrieves an instance of Handler from the world state and updates its value
func (c *HandlerContract) UpdateHandler(ctx contractapi.TransactionContextInterface, handlerID string, newValue string) error {
	exists, err := c.HandlerExists(ctx, handlerID)
	if err != nil {
		return fmt.Errorf("Could not read from world state. %s", err)
	} else if !exists {
		return fmt.Errorf("The asset %s does not exist", handlerID)
	}

	handler := new(Handler)
	handler.Value = newValue

	bytes, _ := json.Marshal(handler)

	return ctx.GetStub().PutState(handlerID, bytes)
}

// DeleteHandler deletes an instance of Handler from the world state
func (c *HandlerContract) DeleteHandler(ctx contractapi.TransactionContextInterface, handlerID string) error {
	exists, err := c.HandlerExists(ctx, handlerID)
	if err != nil {
		return fmt.Errorf("Could not read from world state. %s", err)
	} else if !exists {
		return fmt.Errorf("The asset %s does not exist", handlerID)
	}

	return ctx.GetStub().DelState(handlerID)
}
