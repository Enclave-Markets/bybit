package main

import (
	"time"

	"github.com/Enclave-Markets/bybit/v2"
)

func main() {
	wsClient := bybit.NewTestWebsocketClient()
	svc, err := wsClient.V5().Public(bybit.CategoryV5Linear)
	if err != nil {
		panic(err)
	}
	// A channel to receive the response from the websocket
	c := make(chan any, 100)

	_, err = svc.SubscribeKline(
		bybit.V5WebsocketPublicKlineParamKey{
			Interval: bybit.Interval5,
			Symbol:   bybit.SymbolV5BTCUSDT,
		},
		func(response bybit.V5WebsocketPublicKlineResponse) error {
			c <- response
			return nil
		},
	)
	if err != nil {
		panic(err)
	}

	first := true
	_, err = svc.SubscribeOrderBook(
		bybit.V5WebsocketPublicOrderBookParamKey{
			Depth:  1,
			Symbol: bybit.SymbolV5BTCUSDT,
		},
		func(response bybit.V5WebsocketPublicOrderBookResponse) error {
			if first && response.Type != "snapshot" {
				panic("Orderbook first response should be snapshot")
			}
			first = false
			c <- response
			return nil
		},
	)
	if err != nil {
		panic(err)
	}

	_, err = svc.SubscribeTicker(
		bybit.V5WebsocketPublicTickerParamKey{
			Symbol: bybit.SymbolV5BTCUSDT,
		},
		func(response bybit.V5WebsocketPublicTickerResponse) error {
			c <- response
			return nil
		},
	)

	if err != nil {
		panic(err)
	}

	for i := 0; i < 3; i++ {
		err = svc.Run()
		if err != nil {
			panic(err)
		}
	}

	timeout := time.After(5 * time.Second)
	orderbook := false
	kline := false
	ticker := false
	for !orderbook || !kline || !ticker {
		select {
		case msg := <-c:
			switch msg.(type) {
			case bybit.V5WebsocketPublicOrderBookResponse:
				orderbook = true
			case bybit.V5WebsocketPublicKlineResponse:
				kline = true
			case bybit.V5WebsocketPublicTickerResponse:
				ticker = true
			}

		case <-timeout:
			panic("Sub test timed out")
		}
	}
	println("Sub test passed")
}
