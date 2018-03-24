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
	"fmt"
	"strconv"

	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
)

func init() {
	b := BoolString{}
	Default.Add("boolstring", &b, IsBoolString)
}

// BoolString represents the email string format as specified by the json schema spec
//
// swagger:strfmt email
type BoolString struct {
	Bool  bool
	Valid bool
}

// IsBoolString returns true when the string is a valid date
func IsBoolString(str string) bool {
	if str == "" {
		return true
	}

	_, err := strconv.ParseBool(str)
	return err == nil
}

func (b BoolString) String() string {
	if !b.Valid {
		return ""
	}
	return strconv.FormatBool(b.Bool)
}

func ParseBoolString(data string) (BoolString, error) {
	if data == "" {
		return BoolString{}, nil
	}

	b, err := strconv.ParseBool(data)
	if err != nil {
		return BoolString{}, err
	}
	bb := BoolString{
		Bool:  b,
		Valid: true,
	}

	return bb, nil
}

// MarshalText turns this instance into text
func (b BoolString) MarshalText() ([]byte, error) {
	return []byte(b.String()), nil
}

// UnmarshalText hydrates this instance from text
func (b *BoolString) UnmarshalText(data []byte) error { // validation is performed later on
	bb, err := ParseBoolString(string(data))
	if err != nil {
		return err
	}
	*b = bb

	return nil
}

func (b *BoolString) Scan(raw interface{}) error {
	if raw == nil {
		b.Bool, b.Valid = false, false
		return nil
	}

	var s string
	switch v := raw.(type) {
	case []byte:
		s = string(v)
	case string:
		s = v
	default:
		return fmt.Errorf("cannot sql.Scan() strfmt.BoolString from: %#v", v)
	}

	bb, err := ParseBoolString(s)
	if err != nil {
		return err
	}
	*b = bb

	return nil

}

func (b BoolString) Value() (driver.Value, error) {
	if !b.Valid {
		return nil, nil
	}
	return b.Bool, nil
}

// MarshalJSON returns the BoolString as JSON
func (b BoolString) MarshalJSON() ([]byte, error) {
	var w jwriter.Writer
	b.MarshalEasyJSON(&w)
	return w.BuildBytes()
}

// MarshalEasyJSON writes the BoolString to a easyjson.Writer
func (b BoolString) MarshalEasyJSON(w *jwriter.Writer) {
	if b.Valid {
		w.Bool(b.Bool)
	} else {
		w.RawString("null")
	}
}

// UnmarshalJSON sets the BoolString from JSON
func (b *BoolString) UnmarshalJSON(data []byte) error {
	l := jlexer.Lexer{Data: data}
	b.UnmarshalEasyJSON(&l)
	return l.Error()
}

// UnmarshalEasyJSON sets the BoolString from a easyjson.Lexer
func (b *BoolString) UnmarshalEasyJSON(in *jlexer.Lexer) {
	if data := in.String(); in.Ok() {
		bb, err := ParseBoolString(data)
		if err != nil {
			in.AddError(err)
			return
		}
		*b = bb
	}
}
