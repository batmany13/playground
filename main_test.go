package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var FAIL_TIME = []string{"2017-01-01 11:10", "2017-01-01 11:40", "2017-01-01 12:10", "2017-01-01 12:40"}
var FAIL_TIME2 = []string{"2017-01-01 11:40", "2017-01-01 12:10"}
var SUCCESS_TIME = []string{"2017-01-01 10:30"}

func TestAvailable(t *testing.T) {
	assert := assert.New(t)
	for _, fail := range FAIL_TIME {
		mt := MealTime{}
		check, _ := time.Parse(TIME_LAYOUT, fail)
		mt.Time = check
		req := VendorReq{VendorId: 1, Timestamp: mt}
		ret, err := req.Available()
		assert.Nil(err)
		assert.True(ret)
	}
	for _, fail := range FAIL_TIME2 {
		mt := MealTime{}
		check, _ := time.Parse(TIME_LAYOUT, fail)
		mt.Time = check
		req := VendorReq{VendorId: 3, Timestamp: mt}
		ret, err := req.Available()
		assert.Nil(err)
		assert.False(ret)
	}

	for _, fail := range SUCCESS_TIME {
		mt := MealTime{}
		check, _ := time.Parse(TIME_LAYOUT, fail)
		mt.Time = check
		req := VendorReq{VendorId: 1, Timestamp: mt}
		ret, err := req.Available()
		assert.Nil(err)
		assert.True(ret)
	}
}
