package webapp

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image/color"
	"log"
	"math"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func Plot(m []Measurement, v string) (string, error) {
	log.Println(v)
	log.Println(m)
	p := plot.New()
	switch v {
	case "temperature":
		p.Y.Label.Text = "Temperatura [°C]"
	case "humidity":
		p.Y.Label.Text = "Wilgotność [%]"
	default:
		return "", errors.New("unknown value v")
	}

	pts := make(plotter.XYs, len(m))

	for i := range pts {
		timestamp := float64(m[i].Timestamp)
		pts[i].X = timestamp
		switch v {
		case "temperature":
			pts[i].Y = m[i].Temperature
		case "humidity":
			pts[i].Y = float64(m[i].Humidity)
		default:
			return "", errors.New("unknown value v")
		}
	}

	line, err := plotter.NewLine(pts)
	if err != nil {
		return "", err
	}
	line.LineStyle.Width = vg.Points(1)
	p.BackgroundColor = color.RGBA{R: 5, G: 47, B: 74, A: 255}
	line.LineStyle.Color = color.RGBA{R: 7, G: 89, B: 133, A: 255}

	p.Title.TextStyle.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	p.X.Label.TextStyle.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	p.Y.Label.TextStyle.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}

	p.X.Tick.Label.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	p.Y.Tick.Label.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}

	p.X.LineStyle.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	p.X.Tick.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	p.X.Tick.Label.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	p.Y.LineStyle.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	p.Y.Tick.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	p.Y.Tick.Label.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}

	p.X.LineStyle.Width = vg.Points(2)
	p.Y.LineStyle.Width = vg.Points(2)
	line.LineStyle.Width = vg.Points(2)

	p.Add(line)

	grid := plotter.NewGrid()
	grid.Horizontal.Color = color.RGBA{R: 220, G: 220, B: 220, A: 255}
	grid.Vertical.Color = color.RGBA{R: 220, G: 220, B: 220, A: 255}
	grid.Horizontal.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}
	grid.Vertical.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}
	p.Add(grid)

	p.X.Tick.Marker = plot.TimeTicks{Format: "|\n2006-01-02\n15:04", Ticker: ticker{}}
	p.X.Tick.Label.Rotation = 45 * (math.Pi / 180)

	buf := bytes.NewBuffer(nil)
	pngWriter, err := p.WriterTo(8*vg.Inch, 4*vg.Inch, "png")
	if err != nil {
		return "", err
	}
	_, err = pngWriter.WriteTo(buf)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

type ticker struct{}

func (t ticker) Ticks(min float64, max float64) []plot.Tick {
	ticks := []plot.Tick{}
	f := "2006-01-02\n15:04"
	diff := max - min
	step := diff / 7
	ticks = append(ticks, plot.Tick{Value: min, Label: time.Unix(int64(min), 0).Format(f)})
	for i := min + step; i < max; i += step {
		ticks = append(ticks, plot.Tick{Value: i, Label: time.Unix(int64(i), 0).Format(f)})
	}
	ticks = append(ticks, plot.Tick{Value: max, Label: time.Unix(int64(max), 0).Format(f)})
	return ticks
}
