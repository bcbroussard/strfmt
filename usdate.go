package strfmt

import (
	"database/sql/driver"
	// "errors"
	"fmt"
	"time"

	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
)

func init() {
	d := USDate{}
	// register this format in the default registry
	Default.Add("usdate", &d, IsUSDate)
}

// IsUSDate returns true when the string is a valid date
func IsUSDate(str string) bool {
	if str == "0000-00-00" || str == "00/00/0000" {
		return true
	}
	if _, err := time.Parse(FullUSDateFormat, str); err != nil {
		if _, err := time.Parse(AltUSDateFormat, str); err != nil {
			return false
		}
	}
	return true
}

const (
	// FullUSDateFormat represents a full-date
	FullUSDateFormat = "01/02/2006"
	AltUSDateFormat  = "2006-01-02"
)

// USDate represents a date from the API
//
// swagger:strfmt usdate
type USDate time.Time

// String converts this date into a string
func (d USDate) String() string {
	return time.Time(d).Format(FullUSDateFormat)
}

// UnmarshalText parses a text representation into a date type
func (d *USDate) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		return nil
	}
	sDate := string(text)
	if sDate == "0000-00-00" || sDate == "00/00/0000" {
		return nil
	}

	dd, err := time.Parse(FullUSDateFormat, sDate)
	if err != nil {
		ad, aErr := time.Parse(AltUSDateFormat, sDate)
		if aErr != nil {
			return aErr
		}
		dd = ad
	}
	*d = USDate(dd)
	return nil
}

// MarshalText serializes this date type to string
func (d USDate) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

// Scan scans a USDate value from database driver type.
func (d *USDate) Scan(raw interface{}) error {
	dd := Date{}
	switch v := raw.(type) {
	case []byte:
		if err := dd.UnmarshalText(v); err != nil {
			return err
		}
		*d = USDate(time.Time(dd))
		return nil
	case string:
		if err := dd.UnmarshalText([]byte(v)); err != nil {
			return err
		}
		*d = USDate(time.Time(dd))
		return nil
	case time.Time:
		*d = USDate(v)
		return nil
	case nil:
		*d = USDate{}
		return nil
	default:
		return fmt.Errorf("cannot sql.Scan() strfmt.USDate from: %#v", v)
	}
}

// Value converts USDate to a primitive value ready to written to a database.
func (d USDate) Value() (driver.Value, error) {
	return driver.Value(time.Time(d).Format(RFC3339FullDate)), nil
}

// MarshalJSON returns the USDate as JSON
func (d USDate) MarshalJSON() ([]byte, error) {
	var w jwriter.Writer
	d.MarshalEasyJSON(&w)
	return w.BuildBytes()
}

// MarshalEasyJSON writes the USDate to a easyjson.Writer
func (d USDate) MarshalEasyJSON(w *jwriter.Writer) {
	w.String(time.Time(d).Format(FullUSDateFormat))
}

// UnmarshalJSON sets the USDate from JSON
func (d *USDate) UnmarshalJSON(data []byte) error {
	l := jlexer.Lexer{Data: data}
	d.UnmarshalEasyJSON(&l)
	return l.Error()
}

// UnmarshalEasyJSON sets the USDate from a easyjson.Lexer
func (d *USDate) UnmarshalEasyJSON(in *jlexer.Lexer) {
	if data := in.String(); in.Ok() {
		if data == "0000-00-00" {
			return
		}

		tt, err := time.Parse(FullUSDateFormat, data)
		if err != nil {
			ad, aErr := time.Parse(AltUSDateFormat, data)
			if aErr != nil {
				in.AddError(aErr)
				return
			}
			tt = ad
		}
		*d = USDate(tt)
	}
}

// type Time struct {
// 	time.Time
// }

// const dateFormat = "01/02/2006"

// // const dateTimeFormat = "01/02/2006 03:04pm"

// const dateFormatInput = "01/02/2006 MST"

// // const dateTimeFormatInput = "01/02/2006 03:04pm MST"

// // String converts the unix timestamp into a string
// func (t Time) String() string {
// 	if t.Time.IsZero() {
// 		return ""
// 	}
// 	return t.Format(dateTimeFormat)
// }

// func (t Time) USDateString() string {
// 	if t.Time.IsZero() {
// 		return ""
// 	}
// 	return t.Format(dateFormat)
// }

// func Parse(s string) (Time, error) {
// 	f := dateFormatInput
// 	if len(s) > 10 {
// 		f = dateTimeFormatInput
// 	}
// 	jt, err := time.Parse(f, s+" PST")
// 	if err != nil {
// 		return Time{}, err
// 	}

// 	t := Time{jt}
// 	return t, nil
// }

// func (t *Time) UnmarshalJSON(b []byte) error {
// 	tt, err := Parse(string(b))
// 	if err != nil {
// 		return err
// 	}
// 	t.Time = tt.Time
// 	return nil
// }

// func (t Time) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(t.String())
// }
