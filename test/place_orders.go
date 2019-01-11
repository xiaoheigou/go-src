package main

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"yuudidi.com/pkg/protocol/response"
	"yuudidi.com/pkg/service"
)

func main() {
	//read test data from configs/../test/data/orders.csv, <PROJ_ROOT>/configs/ as root
	absPath, _ := filepath.Abs("../test/data/orders.csv")
	file, err := os.OpenFile(absPath, os.O_RDONLY, 0755)
	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for n, record := range records {
		if n == 0 {
			continue
		}
		on := time.Now().Unix() // use time stamp as order number
		on = on + int64(n)
		price, _ := strconv.ParseFloat(record[2], 32)
		amount, _ := strconv.ParseFloat(record[3], 32)
		distributor, _ := strconv.ParseInt(record[4], 10, 32)
		ot, _ := strconv.ParseInt(record[5], 10, 8)
		tc, _ := strconv.ParseFloat(record[7], 32)
		pt, _ := strconv.ParseInt(record[8], 10, 8)
		request := response.CreateOrderRequest{
			ApiKey:        record[0],
			OrderNo:       strconv.FormatInt(on, 10),
			Price:         float32(price),
			Amount:        amount,
			DistributorId: distributor,
			OrderType:     int(ot),
			CoinType:      record[6],
			TotalCount:    tc,
			PayType:       uint(pt),
			Name:          record[9],
			BankAccount:   record[10],
			BankBranch:    record[11],
			Phone:         record[12],
			Remark:        record[13],
			QrCode:        record[14],
			PageUrl:       record[15],
			ServerUrl:     record[16],
			CurrencyFiat:  record[17],
			AccountId:     record[18],
		}
		service.PlaceOrder(request)
	}
}
