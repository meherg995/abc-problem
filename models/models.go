package models

import "time"

type Class struct {
	ClassName string
	StartDate time.Time
	EndDate   time.Time
	Capacity  int
}

type Booking struct {
	Name string
	Date time.Time
}
