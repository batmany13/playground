package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	MEAL_URL     = "https://gist.githubusercontent.com/huned/1d1c076eac47b34a5e2b0b0aab4aea65/raw/dc187a07d2899bb06d71c06c2008f7aa7406d6d1/meals.json.txt"
	VENDOR_URL   = "https://gist.githubusercontent.com/huned/67b0fb14c4190c76abf3796d28f6811b/raw/9f4d846bca8b3ff8c3f2dd02f8d719ceaefdbd93/vendors.json"
	RESULT_FILE  = "results.json"
	RESULT_FILE2 = "results2.json"
	TIME_LAYOUT  = `2006-01-02 15:04`
)

type MealTime struct {
	time.Time
}

type VendorResult struct {
	Vendors []Vendor `json:"results"`
}

type Meal struct {
	Meals []Result `json:"results"`
}

type Vendor struct {
	VendorId int `json:"vendor_id"`
	Drivers  int `json:"drivers"`
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

func DriverAvail() (map[int]int, error) {
	var vendor VendorResult
	vendorAvail := make(map[int]int)
	resp, err := http.Get(VENDOR_URL)
	if err != nil {
		return vendorAvail, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return vendorAvail, err
	}
	json.Unmarshal(data, &vendor)
	for _, ven := range vendor.Vendors {
		vendorAvail[ven.VendorId] = ven.Drivers - 1
	}
	return vendorAvail, nil
}

func (v *VendorReq) Available() (bool, error) {
	var meal Meal
	vendorAvail, _ := DriverAvail()
	data, err := ioutil.ReadFile(RESULT_FILE)
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
			if vendorAvail[v.VendorId] == 0 {
				return false, nil
			}
			vendorAvail[v.VendorId] = vendorAvail[v.VendorId] - 1
			continue
		}
		if v.Timestamp.After(start) && v.Timestamp.Before(end) {
			if vendorAvail[v.VendorId] == 0 {
				return false, nil
			}
			vendorAvail[v.VendorId] = vendorAvail[v.VendorId] - 1
			continue
		}
	}
	return true, nil
}

func vendorAvail(c *gin.Context) {
	var v VendorReq
	c.BindJSON(&v)
	ret, err := v.Available()
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error,
		})
		return
	}
	c.JSON(200, gin.H{
		"result": ret,
	})
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
