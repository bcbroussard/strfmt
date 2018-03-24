// Copyright 2015 go-swagger maintainers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package strfmt

import (
	"database/sql/driver"
	// "errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	// "gopkg.in/mgo.v2/bson"

	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
)

func init() {
	dt := USDateTime{}
	Default.Add("usdatetime", &dt, IsUSDateTime)
}

// IsUSDateTime returns true when the string is a valid date-time
func IsUSDateTime(str string) bool {
	if len(str) < 4 {
		return false
	}
	s := strings.SplitN(str, " ", 2)
	if len(s) == 1 && IsUSDate(s[0]) {
		return true
	}

	if len(s) < 2 || !IsUSDate(s[0]) {
		return false
	}

	matches := rxUSDateTime.FindAllStringSubmatch(s[1], -1)
	if len(matches) == 0 || len(matches[0]) == 0 {
		return false
	}
	return true
}

const (
	USDateTimeFormat      = "01/02/2006 03:04pm"
	AltUSDateTimeFormat   = "2006-01-02 15:04:05"
	USDateTimeFormatZone  = "01/02/2006 03:04pm MST"
	USDateMixedTimeFormat = "01/02/2006 15:04pm"
	// USDateTimePattern pattern to match for the date-time format
	USDateTimePattern = `^((0[0-9]|1[0-2]):([0-5][0-9])(am|pm)( [A-Z]{3})?)|(((?:0[0-9]|1[0-2])|(?:1[0-9])):([0-5][0-9]):([0-5][0-9])?)$`
)

var (
	usDateTimeFormats = []string{USDateTimeFormat, AltUSDateTimeFormat, FullUSDateFormat, AltUSDateFormat, USDateMixedTimeFormat}
	rxUSDateTime      = regexp.MustCompile(USDateTimePattern)
	// USMarshalFormat sets the time resolution format used for marshaling time (set to milliseconds)
	USMarshalFormat = USDateTimeFormat
)

// ParseUSDateTime parses a string that represents an US Formatted PST
func ParseUSDateTime(data string) (USDateTime, error) {
	if data == "" || data == "0000-00-00 00:00:00" || data == "0000-00-00" || data == "00/00/0000" {
		return NewUSDateTime(), nil
	}
	pst, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		return USDateTime{}, err
	}

	var lastError error
	for _, layout := range usDateTimeFormats {
		dd, err := time.ParseInLocation(layout, data, pst)
		if err != nil {
			lastError = err
			continue
		}
		lastError = nil
		return USDateTime(dd), nil
	}
	return USDateTime{}, lastError
}

// USDateTime is a time but it serializes to ISO8601 format with millis
// It knows how to read 3 different variations of a RFC3339 date time.
// Most APIs we encounter want either millisecond or second precision times.
// This just tries to make it worry-free.
//
// swagger:strfmt date-time
type USDateTime time.Time

// NewUSDateTime is a representation of zero value for USDateTime type
func NewUSDateTime() USDateTime {
	return USDateTime(time.Unix(0, 0).UTC())
}

// String converts this time to a string
func (t USDateTime) String() string {
	tt := time.Time(t)
	if tt.IsZero() || tt.Unix() == 0 {
		return ""
	}
	return tt.Format(USMarshalFormat)
}

// MarshalText implements the text marshaller interface
func (t USDateTime) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

// UnmarshalText implements the text unmarshaller interface
func (t *USDateTime) UnmarshalText(text []byte) error {
	tt, err := ParseUSDateTime(string(text))
	if err != nil {
		return err
	}
	*t = tt
	return nil
}

// Scan scans a USDateTime value from database driver type.
func (t *USDateTime) Scan(raw interface{}) error {
	// TODO: case int64: and case float64: ?
	tt := DateTime{}
	switch v := raw.(type) {
	case []byte:
		if err := tt.UnmarshalText(v); err != nil {
			return err
		}
		*t = USDateTime(time.Time(tt))
	case string:
		if err := tt.UnmarshalText([]byte(v)); err != nil {
			return err
		}
		*t = USDateTime(time.Time(tt))
	case time.Time:
		*t = USDateTime(v)
	case nil:
		*t = USDateTime{}
	default:
		return fmt.Errorf("cannot sql.Scan() strfmt.USDateTime from: %#v", v)
	}

	return nil
}

// Value converts USDateTime to a primitive value ready to written to a database.
func (t USDateTime) Value() (driver.Value, error) {
	tt := time.Time(t)
	if tt.IsZero() || tt.Unix() == 0 {
		return nil, nil
	}
	return tt, nil
}

// MarshalJSON returns the USDateTime as JSON
func (t USDateTime) MarshalJSON() ([]byte, error) {
	var w jwriter.Writer
	t.MarshalEasyJSON(&w)
	return w.BuildBytes()
}

// MarshalEasyJSON writes the USDateTime to a easyjson.Writer
func (t USDateTime) MarshalEasyJSON(w *jwriter.Writer) {
	s := t.String()
	if s == "" {
		w.RawString("null")
	} else {
		w.String(s)
	}
}

// UnmarshalJSON sets the USDateTime from JSON
func (t *USDateTime) UnmarshalJSON(data []byte) error {
	l := jlexer.Lexer{Data: data}
	t.UnmarshalEasyJSON(&l)
	return l.Error()
}

// UnmarshalEasyJSON sets the USDateTime from a easyjson.Lexer
func (t *USDateTime) UnmarshalEasyJSON(in *jlexer.Lexer) {
	if data := in.String(); in.Ok() {
		tt, err := ParseUSDateTime(data)
		if err != nil {
			in.AddError(err)
			return
		}
		*t = tt
	}
}
