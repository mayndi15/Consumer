package model

import (
	"encoding/json"
	"fmt"
	"time"
)

type Gender string

const (
	Male   Gender = "MALE"
	Female Gender = "FEMALE"
	Other  Gender = "OTHER"
)

type CustomTime struct {
	time.Time
}

func (c *CustomTime) UnmarshalJSON(b []byte) error {
	var timestamp []int
	if err := json.Unmarshal(b, &timestamp); err != nil {
		return err
	}
	if len(timestamp) != 7 {
		return fmt.Errorf("invalid timestamp format")
	}
	c.Time = time.Date(timestamp[0], time.Month(timestamp[1]), timestamp[2], timestamp[3], timestamp[4], timestamp[5], timestamp[6], time.UTC)
	return nil
}

type Person struct {
	ID              int        `json:"id"`
	FirstName       string     `json:"firstName"`
	LastName        string     `json:"lastName"`
	Gender          Gender     `json:"gender"`
	CreatedDate     CustomTime `json:"createdDate"`
	UpdatedDate     CustomTime `json:"updatedDate"`
}


var Service string = "person"
