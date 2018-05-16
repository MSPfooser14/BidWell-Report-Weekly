package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Decarium/go-coinmarketcap/coinmarketcap"
	"github.com/alecthomas/template"
)

const (
	BILLION = 1000000000
)

type MarketCapData struct {
	TotalMarketCap                    string
	MarketCapGrowth                   string
	MarketCapGrowthPercentage         string
	TotalWeeklyVolume                 string
	TotalWeeklyVolumeGrowth           string
	TotalWeeklyVolumeGrowthPercentage string
}

//Create functions for each section

//Quick thought here which is that we are taking the average marketcap of the day, we should move this
//To be the market cap at the end of the day UTC time. This way we don't skew the numbers by average.
//We can fix this when we move this to a more standard thing
func CreateMarketCap() {
	t := template.New("marketcap.tmpl")
	t, err := t.ParseFiles("./templates/marketcap.tmpl") // Parse template file.

	if err != nil {
		log.Fatal(err)
	}

	data := GetMarketCapData()

	f, err := os.Create("./pdf/sections/marketcap/marketcap.pug")
	if err != nil {
		log.Println("create file: ", err)
		return
	}

	err = t.Execute(f, data)

	if err != nil {
		fmt.Println(err)
	}

	f.Close()
}
func GetMarketCapData() MarketCapData {

	//Total Market Cap
	totalMarketCap := GetTotalMarketCap()

	oldMarketCap := GetTotalMarketCapLastWeek()

	// Market Cap Growth

	growth := totalMarketCap - oldMarketCap

	growthPercentage := (growth / oldMarketCap) * 100

	// Total Weekly Volume

	totalWeeklyVolume := GetTotalWeeklyVolume()

	oldWeeklyVolume := GetTotalWeeklyVolumeLastWeek()

	// Weekly Volume Growth

	volumeGrowth := totalWeeklyVolume - oldWeeklyVolume

	volumeGrowthPercentage := (growth / oldWeeklyVolume) * 100

	//Format everything
	totalMarketCapFormatted := fmt.Sprintf("%.2f", (totalMarketCap / BILLION))

	growthFormatted := fmt.Sprintf("%.2f", (growth / BILLION))

	growthPercentageFormatted := fmt.Sprintf("%.2f", growthPercentage)

	totalWeeklyVolumeFormatted := fmt.Sprintf("%.2f", (totalWeeklyVolume / BILLION))

	volumeGrowthFormatted := fmt.Sprintf("%.2f", (volumeGrowth / BILLION))

	volumeGrowthPercentageFormatted := fmt.Sprintf("%.2f", volumeGrowthPercentage)

	//If growth is positive, we add the + sign
	if growth > 0 {
		growthFormatted = "+" + growthFormatted
	}

	mcd := MarketCapData{
		TotalMarketCap:                    totalMarketCapFormatted,
		MarketCapGrowth:                   growthFormatted,
		MarketCapGrowthPercentage:         growthPercentageFormatted,
		TotalWeeklyVolume:                 totalWeeklyVolumeFormatted,
		TotalWeeklyVolumeGrowth:           volumeGrowthFormatted,
		TotalWeeklyVolumeGrowthPercentage: volumeGrowthPercentageFormatted,
	}

	//We want to format this data so that it looks clean on the excel

	return mcd

}

//Done
func GetTotalMarketCap() float64 {

	global, err := coinmarketcap.GetGlobal()

	if err != nil {
		fmt.Println(err)
	}

	return global.TotalMarketCapUsd
}

//Done
func GetTotalMarketCapLastWeek() float64 {

	t := time.Now().AddDate(0, 0, -7)

	//Start date is going to be midnight 7 days ago
	start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	//End will be 1 day forward so we get 24hrs of ticks
	end := start.AddDate(0, 0, 1)

	global, err := coinmarketcap.GetGlobalHistoricalTicksDailyByDate(start, end)

	if err != nil {
		fmt.Println(err)
	}

	return global.MarketCapByAvailableSupply[0].Amount
}

func GetTotalWeeklyVolume() float64 {

	t := time.Now().AddDate(0, 0, -8)

	start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	//So we do last week + 7 days so if this report is released on a Tuesday, we get Monday to Monday
	end := start.AddDate(0, 0, 7)

	global, err := coinmarketcap.GetGlobalHistoricalTicksDailyByDate(start, end)

	if err != nil {
		fmt.Println(err)
	}

	var total float64

	//Global is now 7 days worth of volume so we want to iterate through it
	for _, day := range global.VolumeUsd {
		total += day.Amount
	}

	return total
}

func GetTotalWeeklyVolumeLastWeek() float64 {

	t := time.Now().AddDate(0, 0, -15)

	start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	//So we do last week + 7 days so if this report is released on a Tuesday, we get Monday to Monday
	end := start.AddDate(0, 0, 7)

	global, err := coinmarketcap.GetGlobalHistoricalTicksDailyByDate(start, end)

	if err != nil {
		fmt.Println(err)
	}

	var total float64

	//Global is now 7 days worth of volume so we want to iterate through it
	for _, day := range global.VolumeUsd {
		total += day.Amount
	}

	return total
}