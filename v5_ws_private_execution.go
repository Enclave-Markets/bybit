package bybit

import (
	"encoding/json"
	"errors"

	"github.com/gorilla/websocket"
)

type V5WebsocketPrivateExecutionResponse struct {
	ID           string                            `json:"id"`
	Topic        V5WebsocketPrivateTopic           `json:"topic"`
	CreationTime int64                             `json:"creationTime"`
	Data         []V5WebsocketPrivateExecutionData `json:"data"`
}

type V5WebsocketPrivateExecutionData struct {
	Category        CategoryV5 `json:"category"`
	Symbol          SymbolV5   `json:"symbol"`
	ExecFee         string     `json:"execFee"`
	ExecID          string     `json:"execId"`
	ExecPrice       string     `json:"execPrice"`
	ExecQty         string     `json:"execQty"`
	ExecType        string     `json:"execType"`
	ExecValue       string     `json:"execValue"`
	IsMaker         bool       `json:"isMaker"`
	FeeRate         string     `json:"feeRate"`
	TradeIv         string     `json:"tradeIv"`
	MarkIv          string     `json:"markIv"`
	BlockTradeID    string     `json:"blockTradeId"`
	MarkPrice       string     `json:"markPrice"`
	IndexPrice      string     `json:"indexPrice"`
	UnderlyingPrice string     `json:"underlyingPrice"`
	LeavesQty       string     `json:"leavesQty"`
	OrderID         string     `json:"orderId"`
	OrderLinkID     string     `json:"orderLinkId"`
	OrderPrice      string     `json:"orderPrice"`
	OrderQty        string     `json:"orderQty"`
	OrderType       OrderType  `json:"orderType"`
	StopOrderType   string     `json:"stopOrderType"`
	Side            Side       `json:"side"`
	ExecTime        string     `json:"execTime"`
	IsLeverage      string     `json:"isLeverage"`
	ClosedSize      string     `json:"closedSize"`
	Seq             int64      `json:"seq"`
}

// SubscribePosition :
func (s *V5WebsocketPrivateService) SubscribeExecutionAllInOne(
	f func(V5WebsocketPrivateExecutionResponse) error,
) (func() error, error) {
	key := V5WebsocketPrivateParamKey{
		Topic: V5WebsocketPrivateTopicPosition,
	}
	if err := s.addParamExecutionFunc(key, f); err != nil {
		return nil, err
	}
	param := struct {
		Op        string        `json:"op"`
		Args      []interface{} `json:"args"`
		RequestID string        `json:"req_id"`
	}{
		Op:        "subscribe",
		Args:      []interface{}{V5WebsocketPrivateTopicPosition},
		RequestID: "private-execution-all-in-one",
	}
	buf, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	if err := s.writeMessage(websocket.TextMessage, buf); err != nil {
		return nil, err
	}
	if err := s.waitForSubConfirmation("private-execution-all-in-one"); err != nil {
		return nil, err
	}

	return func() error {
		param := struct {
			Op   string        `json:"op"`
			Args []interface{} `json:"args"`
		}{
			Op:   "unsubscribe",
			Args: []interface{}{V5WebsocketPrivateTopicExecution},
		}
		buf, err := json.Marshal(param)
		if err != nil {
			return err
		}
		if err := s.writeMessage(websocket.TextMessage, []byte(buf)); err != nil {
			return err
		}
		s.removeParamExecutionFunc(key)
		return nil
	}, nil
}

func (r *V5WebsocketPrivateExecutionResponse) Key() V5WebsocketPrivateParamKey {
	return V5WebsocketPrivateParamKey{
		Topic: r.Topic,
	}
}

func (s *V5WebsocketPrivateService) addParamExecutionFunc(key V5WebsocketPrivateParamKey, f func(V5WebsocketPrivateExecutionResponse) error) error {
	if _, ok := s.paramExecutionMap[key]; ok {
		return errors.New("callback function already set for this key")
	}
	s.paramExecutionMap[key] = f
	return nil
}

func (s *V5WebsocketPrivateService) removeParamExecutionFunc(key V5WebsocketPrivateParamKey) {
	delete(s.paramExecutionMap, key)
}

func (s *V5WebsocketPrivateService) retrieveExecutionFunc(key V5WebsocketPrivateParamKey) (func(V5WebsocketPrivateExecutionResponse) error, error) {
	if f, ok := s.paramExecutionMap[key]; ok {
		return f, nil
	}
	return nil, errors.New("callback function not set")
}
