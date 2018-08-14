// Copyright 2018 The FailWhale Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

const (
	Scale   = time.Microsecond
	Failure = -1.0
	Unknown = 0.0
	Success = 1.0
	Depth   = 8
)

type Record struct {
	Weight float64
	Stamp  time.Time
}

type History struct {
	Records []Record
	Index   int
}

func NewHistory(size int) *History {
	return &History{
		Records: make([]Record, size),
	}
}

var clock = time.Unix(1533930608, 0)

func now() time.Time {
	now := clock
	clock = clock.Add(time.Duration(rand.Intn(10)+1) * Scale)
	return now
}

func (h *History) Add(weight float64) {
	h.Records[h.Index].Weight = weight
	h.Records[h.Index].Stamp = now()
	h.Index = (h.Index + 1) % len(h.Records)
}

func (h *History) Probability() float64 {
	records, now, factor, sum := h.Records, now(), 0.0, 0.0
	max := float64(now.Sub(records[h.Index].Stamp))
	for _, record := range records {
		a := math.Exp(-float64(now.Sub(record.Stamp)) / max)
		factor += a
		sum += record.Weight * a
	}
	sum /= factor

	return 1 / (1 + math.Exp(8*sum))
}

func main() {
	rand.Seed(1)

	history := NewHistory(Depth)
	for i := 0; i < 8; i++ {
		if rand.Intn(2) == 0 {
			history.Add(Failure)
		} else {
			history.Add(Success)
		}
	}
	history.Add(Success)
	fmt.Println(history.Probability())

	/*clock = clock.Add(1000000 * Scale)
	fmt.Println(history.Probability())

	total := 0
	for _, record := range history.Records {
		if record.Weight == Success {
			total++
		}
	}
	fmt.Println(total)*/

	history = NewHistory(Depth)
	probability := make([]float64, 0, 100)
	for i := 0; i < Depth; i++ {
		history.Add(Success)
		probability = append(probability, history.Probability())
	}
	for i := 0; i < Depth; i++ {
		history.Add(Failure)
		probability = append(probability, history.Probability())
	}
	for i := 0; i < 4*Depth; i++ {
		p := history.Probability()
		if rand.Float64() > p {
			if rand.Intn(2) == 0 {
				history.Add(Failure)
			} else {
				history.Add(Success)
			}
		} else {
			history.Add(Unknown)
		}
		probability = append(probability, p)
	}
	for i := 0; i < Depth; i++ {
		history.Add(Success)
		probability = append(probability, history.Probability())
	}

	pts := make(plotter.XYs, len(probability))
	for i, p := range probability {
		pts[i].X = float64(i)
		pts[i].Y = p
	}
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Probability over time"
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "Probability"

	err = plotutil.AddLinePoints(p, "", pts)
	if err != nil {
		panic(err)
	}

	err = p.Save(8*vg.Inch, 8*vg.Inch, "probability_vs_time.png")
	if err != nil {
		panic(err)
	}
}
