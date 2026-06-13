package config

import (
	"strconv"
	"time"
)

type Duration time.Duration

func (d *Duration) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		*d = 0
		return nil
	}

	if data[0] == '"' {
		value, err := strconv.Unquote(string(data))
		if err != nil {
			return err
		}
		if value == "" {
			*d = 0
			return nil
		}
		duration, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		*d = Duration(duration)
		return nil
	}

	seconds, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	*d = Duration(time.Duration(seconds) * time.Second)
	return nil
}

func (d *Duration) UnmarshalText(text []byte) error {
	value := string(text)
	if value == "" {
		*d = 0
		return nil
	}
	if seconds, err := strconv.ParseInt(value, 10, 64); err == nil {
		*d = Duration(time.Duration(seconds) * time.Second)
		return nil
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		return err
	}
	*d = Duration(duration)
	return nil
}

func (d Duration) Duration() time.Duration {
	return time.Duration(d)
}
