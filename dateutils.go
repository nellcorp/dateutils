package dateutils

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type JsonDate time.Time

var timeFormats = []string{time.RFC3339, time.RFC1123Z, time.RFC3339, "2006-01-02 15:04:05", "2006-01-02T15:04:05", "2006-01-02", "1/2/06", "1/2/06 15:05", "1_2_06", "02.01.2006", "20060102"}

func (j *JsonDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := ParseTime(s)
	if err != nil {
		return err
	}
	*j = JsonDate(t)
	return nil
}

func (j JsonDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(j))
}

func (j JsonDate) Format(s string) string {
	t := time.Time(j)
	return t.Format(s)
}

func (j JsonDate) Time() time.Time {
	return time.Time(j)
}

func ParseTimestampString(timestamp string) (date time.Time, err error) {
	i, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return
	}
	date, err = time.Unix(i, 0), nil
	if date.After(time.Now()) && strings.HasSuffix(timestamp, "000") {
		//Check if timestamp was in milisseconds
		date, err = time.Unix(i/1000, 0), nil
	}
	return
}

func ParseTime(t string) (result time.Time, err error) {
	if result, err = ParseTimestampString(t); err == nil {
		return
	}

	for _, f := range timeFormats {
		if result, err = time.Parse(f, t); err == nil {
			return
		}
	}
	err = fmt.Errorf("Could not parse time: %v", err)
	return
}
