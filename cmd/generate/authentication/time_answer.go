package authentication

import (
	"fmt"
	"time"
)

type TimeAnswer struct {
	input string
	value time.Time
}

func (ta *TimeAnswer) WriteAnswer(name string, value interface{}) error {
	return nil
}

func (ta *TimeAnswer) Validate(value interface{}) error {
	ta.input = value.(string)
	dr, err := time.ParseDuration(ta.input)
	if err == nil {
		ta.value = time.Now().Add(dr)
		return nil
	}
	t, err := time.Parse("2006-01-02 15:04:05", ta.input)
	if err != nil {
		return fmt.Errorf("answer value must a time duration or valid time format such as '2006-01-02 15:04:05'")
	}
	if t.UnixNano() < time.Now().UnixNano() {
		return fmt.Errorf("answer value time cannot be in the past")
	}
	ta.value = t
	return nil

}
