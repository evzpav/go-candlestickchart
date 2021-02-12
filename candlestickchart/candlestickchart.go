package candlestickchart

import (
	"io"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

type ChartCandle struct {
	Symbol    string
	Timestamp time.Time `json:"timestamp"`
	Open      float64   `json:"open"`
	High      float64   `json:"high"`
	Low       float64   `json:"low"`
	Close     float64   `json:"close"`
	Volume    float64   `json:"volume"`
}

type CandleStick struct {
	Page   *components.Page
	Charts []components.Charter
	Opts   []charts.SeriesOpts
	Klines []*charts.Kline
}

func NewCandleStickPage() *CandleStick {
	return &CandleStick{
		Page:   components.NewPage(),
		Charts: make([]components.Charter, 0),
	}
}

func (cs *CandleStick) AddCandleStickChart(name string, candles []ChartCandle) {
	kl := charts.NewKLine()

	kl.SetGlobalOptions(
		charts.WithToolboxOpts(opts.Toolbox{Show: true}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:      true,
			Trigger:   "axis",
			TriggerOn: "mousemove|click",
		}),
		charts.WithTitleOpts(opts.Title{
			Title: name,
		}),
		charts.WithXAxisOpts(opts.XAxis{
			SplitNumber: 7,
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Scale:     true,
			SplitLine: &opts.SplitLine{Show: true},
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "inside",
			Start:      0,
			End:        100,
			XAxisIndex: []int{0},
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "slider",
			Start:      0,
			End:        100,
			XAxisIndex: []int{0},
		}),
	)

	x := make([]string, 0)
	y := make([]opts.KlineData, 0)
	for _, c := range candles {
		x = append(x, c.Timestamp.Format("2006/01/02 15:04"))
		y = append(y, opts.KlineData{Name: c.Symbol, Value: [4]float64{c.Open, c.Close, c.Low, c.High}})
	}

	styleOpts := charts.WithItemStyleOpts(opts.ItemStyle{
		Color:        "#00da3c",
		Color0:       "#ec0000",
		BorderColor:  "#008F28",
		BorderColor0: "#8A0000",
	})

	cs.Opts = append(cs.Opts, styleOpts)

	kl = kl.SetXAxis(x).AddSeries(candles[0].Symbol, y)

	cs.Klines = append(cs.Klines, kl)

}

func (cs *CandleStick) AddPins(pins ...charts.SeriesOpts) {

}

func (cs *CandleStick) AddBuyPoint(date string, value float64) {
	cs.AddMarkPoint("buy", date, value)
}

func (cs *CandleStick) AddSellPoint(date string, value float64) {
	cs.AddMarkPoint("sell", date, value)
}

func (cs *CandleStick) AddMarkPoint(name, date string, value float64) {
	markPoint := charts.WithMarkPointNameCoordItemOpts(opts.MarkPointNameCoordItem{
		Name:       name,
		Coordinate: []interface{}{date, value},
		Label: &opts.Label{
			Show:      true,
			Formatter: name,
			Color:     "white",
		},
	})
	cs.Opts = append(cs.Opts, markPoint)
}

func (cs *CandleStick) Render(w io.Writer) {
	for _, kl := range cs.Klines {
		kl.SetSeriesOptions(cs.Opts...)
		cs.Charts = append(cs.Charts, kl)
	}

	cs.Page.AddCharts(cs.Charts...)
	cs.Page.Render(w)
}
