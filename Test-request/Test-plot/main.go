package main

import (
	TQ "MGA_OJ/Test-request"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func main() {
	Points("pressure1.dat")
	Points("pressure2.dat")
	Points("pressure3.dat")
	Points("pressure4.dat")
	Points("pressure5.dat")
	Points("pressure6.dat")
	Points("pressure7.dat")
	Points("pressure8.dat")
	Points("pressure9.dat")
	Points("pressure10.dat")
}

func Points(dataFile string) {
	var max_spand, ave_spand float64
	userMap := make(map[string]bool)
	max_spand = 0
	ave_spand = 0
	p := plot.New()

	p.Title.Text = "response time"
	p.X.Label.Text = "competition time"
	p.Y.Label.Text = "response time"

	data, err := ioutil.ReadFile("../Test-pressure/" + dataFile)
	if err != nil {
		fmt.Println("读取文件失败：", err)
	}
	var records []TQ.Record
	json.Unmarshal(data, &records)
	points := make(plotter.XYs, len(records))
	for i := range points {
		points[i].X = float64(records[i].Created_at)
		points[i].Y = float64(records[i].Spand)
		max_spand = math.Max(max_spand, points[i].Y)
		ave_spand += points[i].Y
		userMap[records[i].UserId] = true
	}

	ave_spand /= float64(len(records))

	fmt.Printf(dataFile+":max_spand:%v, ave_spand:%v, users:%v, submits:%v, handling_capacity:%v, ave_handling_capacity:%v\n", max_spand, ave_spand, len(userMap), len(records), len(records), float64(len(records))/7500.0)

	err = plotutil.AddLinePoints(p, "response-time", points)
	if err != nil {
		log.Fatal(err)
	}

	if err = p.Save(32*vg.Inch, 16*vg.Inch, "./response-time/"+dataFile+".png"); err != nil {
		log.Fatal(err)
	}

	var handlingCapacity []int64
	index := 0
	handling := 0

	for i := 1; i <= 1500; i++ {
		for index < len(records) && records[index].Created_at < int64(i*5*1000) {
			handling++
			index++
		}
		handlingCapacity = append(handlingCapacity, int64(handling))
		handling = 0
	}
	p = plot.New()
	p.Title.Text = "handling capacity"
	p.X.Label.Text = "competition time"
	p.Y.Label.Text = "handling capacity"
	points = make(plotter.XYs, len(handlingCapacity))
	for i := range points {
		points[i].X = float64(i * 5)
		points[i].Y = float64(handlingCapacity[i])
	}
	err = plotutil.AddLinePoints(p, "handling capacity-time", points)
	if err != nil {
		log.Fatal(err)
	}

	if err = p.Save(32*vg.Inch, 16*vg.Inch, "./handling-capacity-time/"+dataFile+".png"); err != nil {
		log.Fatal(err)
	}

	p = plot.New()
	p.Title.Text = "submits"
	p.X.Label.Text = "response"
	p.Y.Label.Text = "submits"
	points = make(plotter.XYs, int(max_spand/200+1))
	for i := 0; i < int(max_spand/200+1); i++ {
		t := 0
		for j := range records {
			if records[j].Spand >= int64(i)*200 && records[j].Spand < int64(i+1)*100 {
				t++
			}
		}
		points[i].X = float64(i * 200)
		points[i].Y = float64(t)
	}
	err = plotutil.AddLinePoints(p, "response-submits", points)
	if err != nil {
		log.Fatal(err)
	}

	if err = p.Save(32*vg.Inch, 16*vg.Inch, "./response-submits/"+dataFile+".png"); err != nil {
		log.Fatal(err)
	}

}
