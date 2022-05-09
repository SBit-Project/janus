package transformer

import (
	"github.com/SBit-Project/janus/pkg/sbit"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

func GetLogger(proxy ETHProxy, q *sbit.Sbit) log.Logger {
	method := proxy.Method()
	logger := q.Client.GetLogger()
	return log.WithPrefix(level.Info(logger), method)
}

func GetLoggerFromETHCall(proxy *ProxyETHCall) log.Logger {
	return GetLogger(proxy, proxy.Sbit)
}

func GetDebugLogger(proxy ETHProxy, q *sbit.Sbit) log.Logger {
	method := proxy.Method()
	logger := q.Client.GetDebugLogger()
	return log.WithPrefix(level.Debug(logger), method)
}

func GetDebugLoggerFromETHCall(proxy *ProxyETHCall) log.Logger {
	return GetDebugLogger(proxy, proxy.Sbit)
}
