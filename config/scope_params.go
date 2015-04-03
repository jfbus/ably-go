package config

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// ScoreParams are parameters used to define the scope of messages to return
// e.g. for history queries.
type ScopeParams struct {
	Start time.Time
	End   time.Time
	Unit  string
}

// EncodeValues encode ScopeParams to url values. Start and End are converted to int64
// timestamps.
func (s *ScopeParams) EncodeValues(out *url.Values) error {
	if !s.Start.IsZero() && !s.End.IsZero() && s.Start.After(s.End) {
		return fmt.Errorf("ScopeParams.Start can not be after ScopeParams.End")
	}

	if !s.Start.IsZero() {
		out.Set("start", strconv.FormatInt(toInt(s.Start), 10))
	}

	if !s.End.IsZero() {
		out.Set("end", strconv.FormatInt(toInt(s.End), 10))
	}

	if s.Unit != "" {
		out.Set("unit", s.Unit)
	}

	return nil
}

func toInt(t time.Time) int64 {
	return int64(time.Nanosecond * time.Duration(t.UnixNano()) / time.Millisecond)
}
