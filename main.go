package main

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	URL          = "https://gist.githubusercontent.com/huned/1d1c076eac47b34a5e2b0b0aab4aea65/raw/dc187a07d2899bb06d71c06c2008f7aa7406d6d1/meals.json.txt"
	RESULT_FILE  = "results.json"
	RESULT_FILE2 = "results2.json"
	TIME_LAYOUT  = `2006-01-02 15:04`
)

type MealTime struct {
	time.Time
}

type Meal struct {
	Meals []Result `json:"results"`
}

type Result struct {
	VendorId int      `json:"vendor_id"`
	ClientId int      `json:"client_id"`
	Datetime MealTime `json:"datetime"`
}

type VendorReq struct {
	VendorId  int      `json:"vendor_id"`
	Timestamp MealTime `json:"timestamp"`
}

func (m *MealTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	m.Time, err = time.Parse(TIME_LAYOUT, s)
	return
}

func (v *VendorReq) Available() (bool, error) {
	var meal Meal
	// resp, err := http.Get(URL)
	// if err != nil {
	// 	return false, err
	// }
	data, err := ioutil.ReadFile(RESULT_FILE)
	// data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	json.Unmarshal(data, &meal)
	for _, res := range meal.Meals {
		if res.VendorId != v.VendorId {
			continue
		}
		start := res.Datetime.Add(-20 * time.Minute)
		end := res.Datetime.Add(10 * time.Minute)
		if v.Timestamp.Equal(start) || v.Timestamp.Equal(end) {
			return false, nil
		}
		if v.Timestamp.After(start) && v.Timestamp.Before(end) {
			return false, nil
		}
	}
	return true, nil
}

func vendorAvail(c *gin.Context) {
	// ret, err := vendor.Availabe()
}

func setupRouter() {
	router := gin.Default()
	v1 := router.Group("/v1")
	{
		v1.POST("/vendor/available", vendorAvail)
	}
	router.Run(":8000")
}

func main() {
	setupRouter()
}
