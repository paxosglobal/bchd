// Copyright (c) 2014-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package rpcclient

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/gcash/bchd/btcjson"
)

// FutureRawResult is a future promise to deliver the result of a RawRequest RPC
// invocation (or an applicable error).
type FutureRawResult chan *response

// Receive waits for the response promised by the future and returns the raw
// response, or an error if the request was unsuccessful.
func (r FutureRawResult) Receive() (json.RawMessage, error) {
	return receiveFuture(r)
}

// RawRequestAsync returns an instance of a type that can be used to get the
// result of a custom RPC request at some future time by invoking the Receive
// function on the returned instance.
//
// See RawRequest for the blocking version and more details.
func (c *Client) RawRequestAsync(method string, params []json.RawMessage) FutureRawResult {
	// Method may not be empty.
	if method == "" {
		return newFutureError(errors.New("no method"))
	}

	// Marshal parameters as "[]" instead of "null" when no parameters
	// are passed.
	if params == nil {
		params = []json.RawMessage{}
	}

	// Create a raw JSON-RPC request using the provided method and params
	// and marshal it.  This is done rather than using the sendCmd function
	// since that relies on marshalling registered btcjson commands rather
	// than custom commands.
	id := c.NextID()
	rawRequest := &btcjson.Request{
		Jsonrpc: "1.0",
		ID:      id,
		Method:  method,
		Params:  params,
	}
	marshalledJSON, err := json.Marshal(rawRequest)
	if err != nil {
		return newFutureError(err)
	}

	// Generate the request and send it along with a channel to respond on.
	responseChan := make(chan *response, 1)
	jReq := &jsonRequest{
		id:             id,
		method:         method,
		cmd:            nil,
		marshalledJSON: marshalledJSON,
		responseChan:   responseChan,
	}
	c.sendRequest(jReq)

	return responseChan
}

// RawRequest allows the caller to send a raw or custom request to the server.
// This method may be used to send and receive requests and responses for
// requests that are not handled by this client package, or to proxy partially
// unmarshaled requests to another JSON-RPC server if a request cannot be
// handled directly.
func (c *Client) RawRequest(method string, params []json.RawMessage) (json.RawMessage, error) {
	return c.RawRequestAsync(method, params).Receive()
}

// RawRequestAsyncWithContext returns an instance of a type that can be used to
// get the result of a custom RPC request at some future time by invoking the
// Receive function on the returned instance. The provided context is attached
// to the underlying HTTP request for cancellation support.
//
// See RawRequestWithContext for the blocking version.
func (c *Client) RawRequestAsyncWithContext(ctx context.Context, method string, params []json.RawMessage) FutureRawResult {
	// Method may not be empty.
	if method == "" {
		return newFutureError(errors.New("no method"))
	}

	// Marshal parameters as "[]" instead of "null" when no parameters
	// are passed.
	if params == nil {
		params = []json.RawMessage{}
	}

	// Create a raw JSON-RPC request using the provided method and params
	// and marshal it.  This is done rather than using the sendCmd function
	// since that relies on marshalling registered btcjson commands rather
	// than custom commands.
	id := c.NextID()
	rawRequest := &btcjson.Request{
		Jsonrpc: "1.0",
		ID:      id,
		Method:  method,
		Params:  params,
	}
	marshalledJSON, err := json.Marshal(rawRequest)
	if err != nil {
		return newFutureError(err)
	}

	// Generate the request and send it along with a channel to respond on.
	responseChan := make(chan *response, 1)
	jReq := &jsonRequest{
		id:             id,
		method:         method,
		cmd:            nil,
		marshalledJSON: marshalledJSON,
		responseChan:   responseChan,
		ctx:            ctx,
	}
	c.sendRequest(jReq)

	return responseChan
}

// RawRequestWithContext allows the caller to send a raw or custom request to
// the server with context support for cancellation. This method may be used to
// send and receive requests and responses for requests that are not handled by
// this client package.
func (c *Client) RawRequestWithContext(ctx context.Context, method string, params []json.RawMessage) (json.RawMessage, error) {
	return c.RawRequestAsyncWithContext(ctx, method, params).Receive()
}
