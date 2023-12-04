package main

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"

	// "io"
	"net/http"
	"time"

	"golang.org/x/net/html/charset"
)

type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Date    string   `xml:"Date,attr"`
	Name    string   `xml:"name,attr"`
	Valutes []Valute  `xml:"Valute"`
}

type Valute struct {
	XMLName   xml.Name `xml:"Valute"`
	ID        string   `xml:"ID,attr"`
	NumCode   string   `xml:"NumCode"`
	CharCode  string   `xml:"CharCode"`
	Nominal   string   `xml:"Nominal"`
	Name      string   `xml:"Name"`
	Value     string   `xml:"Value"`
	VunitRate string   `xml:"VunitRate"`
}

func main() {
	var max, min, mid float64
	min = 100000
	counter := 0
	max_date := ""
	max_name := ""
	min_date := ""
	min_name := ""
	// URL вашего API
	startDate := time.Date(2023, time.December, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 90; i++ {
		startDate = startDate.AddDate(0, 0, -1)
		date := startDate.Format("02/01/2006")
		url := "http://www.cbr.ru/scripts/XML_daily_eng.asp?date_req=" + date
		client := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		req.Header.Add("User-Agent", "My-App")
		req.Header.Add("Authorization", "Bearer YourAccessToken")
		
		if err != nil {
			fmt.Println("Ошибка при создании запроса:", err)
			return
		}
		response, err := client.Do(req)
		if err != nil {
			fmt.Println("Ошибка при выполнении GET-запроса:", err)
			return
		}
		fmt.Println(response.Status)
		defer response.Body.Close()
		var valCurs ValCurs
		decoder := xml.NewDecoder(response.Body)
		decoder.CharsetReader = charset.NewReaderLabel
		err = decoder.Decode(&valCurs)
		if err != nil {
			fmt.Println("Ошибка при декодировании XML:", err)
			return
		}
		fmt.Println("Process day ", i)
		for _, valute := range valCurs.Valutes {
			str := strings.Replace(valute.Value, ",", ".", -1)
			val, err := strconv.ParseFloat(str, 64)
			if err != nil {
				fmt.Println(err.Error())
			}
			mid += val
			if val > max {
				max = val
				max_date = valCurs.Date
				max_name = valute.CharCode
			}
			if val < min {
				min = val
				min_date = valCurs.Date
				min_name = valute.CharCode
			}
				mid += val
				counter++
		}
		time.Sleep(time.Second/2)
	}
	fmt.Println("Max curs is", max_name, "with", max, "at", max_date)
	fmt.Println("Min curs is", min_name, "with", min, "at", min_date)
	fmt.Println("Avg curs for RUB is ", mid/float64(counter))
}