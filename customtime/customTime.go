package customtime

import (
	"database/sql"
	"strconv"
	"time"
)

// TODO: A custom type will mean comparisons like .After() will not work without casting.
// Would it be better to wrap instead of subtype?

type CustomTime time.Time

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	t := time.Time(ct)
	s := strconv.FormatInt(t.UnixMilli(), 10)
	return []byte(s), nil
}

func (ct *CustomTime) UnmarshalJSON(d []byte) error {
	u, err := strconv.ParseInt(string(d), 10, 64)
	if err != nil {
		return err
	}
	t := time.UnixMilli(u)
	*ct = CustomTime(t)
	return nil
}

// NullTimeToCustomTime helps convert NullTime to a CustomTime pointer or nil.
func NullTimeToCustomTime(nt sql.NullTime) *CustomTime {
	if nt.Valid {
		c := CustomTime(nt.Time)
		return &c
	}
	return nil
}
