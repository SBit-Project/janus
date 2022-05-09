package transformer

import (
	"errors"

	"github.com/SBit-Project/janus/pkg/eth"
	"github.com/labstack/echo"
)

var UnmarshalRequestErr = errors.New("Input is invalid")

type Option func(*Transformer) error

type ETHProxy interface {
	Request(*eth.JSONRPCRequest, echo.Context) (interface{}, eth.JSONRPCError)
	Method() string
}
