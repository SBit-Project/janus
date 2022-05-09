package eth

import (
	"encoding/json"
	"fmt"
)

// unknown service
// return fmt.Sprintf("The method %s%s%s does not exist/is not available", e.service, serviceMethodSeparator, e.method)
var MethodNotFoundErrorCode = -32601

// invalid request
var InvalidRequestErrorCode = -32600
var InvalidMessageErrorCode = -32700
var InvalidParamsErrorCode = -32602

// logic error
var CallbackErrorCode = -32000

// shutdown error
// "server is shutting down"
var ShutdownErrorCode = -32000
var ShutdownError = NewJSONRPCError(ShutdownErrorCode, "server is shutting down", nil)

func NewMethodNotFoundError(method string) JSONRPCError {
	return NewJSONRPCError(
		MethodNotFoundErrorCode,
		fmt.Sprintf("The method %s does not exist/is not available", method),
		nil,
	)
}

func NewInvalidRequestError(message string) JSONRPCError {
	return NewJSONRPCError(InvalidRequestErrorCode, message, nil)
}

func NewInvalidMessageError(message string) JSONRPCError {
	return NewJSONRPCError(InvalidMessageErrorCode, message, nil)
}

func NewInvalidParamsError(message string) JSONRPCError {
	return NewJSONRPCError(InvalidParamsErrorCode, message, nil)
}

func NewCallbackError(message string) JSONRPCError {
	return NewJSONRPCError(CallbackErrorCode, message, nil)
}

type JSONRPCError interface {
	Code() int
	Message() string
	Error() error
}

func NewJSONRPCError(code int, message string, err error) JSONRPCError {
	return &GenericJSONRPCError{
		code:    code,
		message: message,
		err:     err,
	}
}

// JSONRPCError contains the message and code for an ETH RPC error
type GenericJSONRPCError struct {
	code    int
	message string
	err     error
}

func (err *GenericJSONRPCError) Code() int {
	return err.code
}

func (err *GenericJSONRPCError) Message() string {
	return err.message
}

func (err *GenericJSONRPCError) Error() error {
	return err.err
}

func (err *GenericJSONRPCError) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}{
		Code:    err.code,
		Message: err.message,
	})
}
