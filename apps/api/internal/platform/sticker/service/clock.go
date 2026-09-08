package service

import "time"

func nowMinusDay() time.Time { return time.Now().Add(-24 * time.Hour) }
