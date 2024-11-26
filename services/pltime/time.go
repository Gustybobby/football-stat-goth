package pltime

import (
	"strconv"
	"time"
)

func GetCurrentSeasonString() string {
	curr_year := time.Now().Year()
	return strconv.Itoa(curr_year) + "/" + strconv.Itoa((curr_year%1000)+1)
}
