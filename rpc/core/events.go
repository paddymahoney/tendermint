package core

import (
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	rpctypes "github.com/tendermint/tendermint/rpc/lib/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

func Subscribe(wsCtx rpctypes.WSRPCContext, query string) (*ctypes.ResultSubscribe, error) {
	logger.Info("Subscribe to query", "remote", wsCtx.GetRemoteAddr(), "query", query)
	ch := pubsub.Subscribe(query)
	if err := wsCtx.AddSubscription(query, ch); err != nil {
		pubsub.Unsubscribe(ch)
		return nil, err
	}
	go func() {
		for event := range ch {
			tmResult := &ctypes.ResultEvent{query, event.(tmtypes.TMEventData)}
      if wsCtx.IsRunning() {
        wsCtx.WriteRPCResponse(rpctypes.NewRPCResponse(wsCtx.Request.ID+"#event", tmResult, ""))
      } else {
        pubsub.Unsubscribe(ch)
      }
		}
	}()
	return &ctypes.ResultSubscribe{}, nil
}

func Unsubscribe(wsCtx rpctypes.WSRPCContext, query string) (*ctypes.ResultUnsubscribe, error) {
	logger.Info("Unsubscribe from query", "remote", wsCtx.GetRemoteAddr(), "query", query)
	ch := wsCtx.DeleteSubscription(query)
	if ch != nil {
		pubsub.Unsubscribe(ch)
	}
	return &ctypes.ResultUnsubscribe{}, nil
}

func UnsubscribeAll(wsCtx rpctypes.WSRPCContext) (*ctypes.ResultUnsubscribe, error) {
	logger.Info("Unsubscribe from all", "remote", wsCtx.GetRemoteAddr())
	channels := wsCtx.DeleteAllSubscriptions()
	for _, ch := range channels {
		pubsub.Unsubscribe(ch)
	}
	return &ctypes.ResultUnsubscribe{}, nil
}
