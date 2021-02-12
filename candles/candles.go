package candles

import (
	"fmt"
	"time"

	"github.com/evzpav/crex"
	"github.com/evzpav/crex/exchanges"
)

type Candle struct {
	ID        int       `json:"id"`
	Symbol    string    `json:"symbol"`
	Period    string    `json:"period"`
	Timestamp time.Time `json:"timestamp"`
	Open      float64   `json:"open"`
	High      float64   `json:"high"`
	Low       float64   `json:"low"`
	Close     float64   `json:"close"`
	Volume    float64   `json:"volume"`
	Source    string    `json:"source"`
}

type CandleService struct {
	exchangeAccount crex.Exchange
}

func NewService() *CandleService {
	s := &CandleService{}
	s.setExchange()
	return s
}

func (s *CandleService) setExchange() {
	exchangeAccount := exchanges.NewExchange(exchanges.BinanceFutures,
		// ApiProxyURLOption("socks5://127.0.0.1:1080"), // 使用代理
		crex.ApiAccessKeyOption("key"),
		crex.ApiSecretKeyOption("secret"),
		// ApiTestnetOption(true)
	)
	s.exchangeAccount = exchangeAccount
}

func (s *CandleService) RetrieveCandles(symbol, period string, numberOfCandles int) ([]*Candle, error) {
	records, err := s.exchangeAccount.GetRecords(symbol, period, 0, 0, numberOfCandles)
	if err != nil {
		errMsg := fmt.Errorf("failed to get records/candles: %v", err)
		return nil, errMsg
	}

	var candles []*Candle
	for _, r := range records {
		candles = append(candles, convertToCandle(r, period, s.exchangeAccount.GetName()))
	}

	if len(candles) == 0 {
		return nil, fmt.Errorf("no candles found")
	}

	return candles, nil

}

func convertToCandle(r *crex.Record, period, source string) *Candle {
	return &Candle{
		Symbol:    r.Symbol,
		Timestamp: r.Timestamp,
		Open:      r.Open,
		High:      r.High,
		Low:       r.Low,
		Close:     r.Close,
		Volume:    r.Volume,
		Period:    period,
		Source:    source,
	}
}
