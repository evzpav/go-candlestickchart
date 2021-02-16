package main

import (
	"fmt"
	"net/http"

	"go-candlestickchart/candles"
	"go-candlestickchart/candlestickchart"

	"github.com/gorilla/mux"
)

func main() {

	candleService := candles.NewService()
	h := NewHandler(candleService)

	r := mux.NewRouter()
	r.HandleFunc("/", h.indexHandler).Methods("GET")
	r.HandleFunc("/candlestickchart", h.getCandleStickchart).Methods("GET")

	fmt.Printf("Candlestich chart server running on http://localhost:7777\n")

	http.ListenAndServe(":7777", r)

}

type handler struct {
	candleService *candles.CandleService
}

func NewHandler(candleService *candles.CandleService) *handler {
	return &handler{
		candleService: candleService,
	}
}

func (h *handler) indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Candlestick Charts</h1>"))
	w.Write([]byte(`<p><a href="\candlestickchart">BTCUSDT Daily Chart</a></p>`))
	w.Write([]byte(`<p><a href="\candlestickchart?symbol=BTCUSDT&period=1h">BTCUSDT Houly Chart</a></p>`))
	w.Write([]byte(`<p><a href="\candlestickchart?symbol=ETHUSDT&period=1d">ETHUSDT Daily Chart</a><p>`))
}

func (h *handler) getCandleStickchart(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	period := "1d"
	symbol := "BTCUSDT"

	periodQuery := query.Get("period")
	if periodQuery != "" {
		period = periodQuery
	}

	symbolQuery := query.Get("symbol")
	if symbolQuery != "" {
		symbol = symbolQuery
	}

	fmt.Printf("Retrieving candles for %s in period %s\n", symbol, period)
	candles, err := h.candleService.RetrieveCandles(symbol, period, 100)
	if err != nil {
		fmt.Printf("failed to retrieve candles: %v", err)
	}

	chartCandles := convertCandlesToChartCandles(candles)

	csPage := candlestickchart.NewCandleStickPage()
	csPage.AddCandleStickChart(symbol, chartCandles)
	csPage.AddMarkPoint("buy", "2020/12/15 21:00", 22000.0)
	csPage.AddMarkPoint("sell", "2021/02/07 21:00", 46000.0)
	csPage.Render(w)

}

func convertCandlesToChartCandles(candles []*candles.Candle) []candlestickchart.ChartCandle {
	var chartCandles []candlestickchart.ChartCandle
	for i := 0; i < len(candles); i++ {
		c := candles[i]
		cc := candlestickchart.ChartCandle{
			Symbol:    c.Symbol,
			Timestamp: c.Timestamp,
			Open:      c.Open,
			High:      c.High,
			Low:       c.Low,
			Close:     c.Close,
		}

		chartCandles = append(chartCandles, cc)
	}

	return chartCandles

}
