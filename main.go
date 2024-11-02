package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type StoreAvailability struct {
	AvailableQuantity int    `json:"availableQty"`
	EndDate           string `json:"endDate"`
	StartDate         string `json:"startDate"`
}

var baseUrl = "https://apionline.homedepot.com/product-information/rental/reservations/availability?categoryCode=95&subCategoryCode=001"

func main() {
	storeIdPtr := flag.String("storeId", "", "Home Depot store ID to search for rental availability")

	flag.Parse()

	if len(*storeIdPtr) == 0 {
		log.Fatal("Missing Store ID")
	}

	baseUrl = baseUrl + "&storeList=" + *storeIdPtr

	var wg sync.WaitGroup

	start := time.Now()
	end := start.AddDate(0, 2, 0)

	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		formattedDate := d.Format("2006-01-02")
		wg.Add(1)
		go func(formattedDate string) {
			defer wg.Done()
			searchAvailabilityFor(formattedDate)
		}(formattedDate)
	}

	wg.Wait()
}

func searchAvailabilityFor(date string) {
	res, err := http.Get(baseUrl + "&reservationDate=" + date + "&endDate=" + date)

	if err != nil {
		panic(err.Error())
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		panic(err.Error())
	}

	var data []StoreAvailability
	json.Unmarshal(body, &data)
	fmt.Printf("%+v\n", data)
}
