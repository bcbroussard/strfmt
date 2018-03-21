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
	"database/sql"
	"database/sql/driver"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var _ sql.Scanner = &USDate{}
var _ driver.Valuer = USDate{}

func TestUSDate(t *testing.T) {
	pp := USDate{}
	err := pp.UnmarshalText([]byte{})
	assert.NoError(t, err)
	err = pp.UnmarshalText([]byte("yada"))
	assert.Error(t, err)

	orig := "12/15/2014"
	b := []byte(orig)
	bj := []byte("\"" + orig + "\"")
	err = pp.UnmarshalText([]byte(orig))
	assert.NoError(t, err)

	txt, err := pp.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, orig, string(txt))

	err = pp.UnmarshalJSON(bj)
	assert.NoError(t, err)
	assert.EqualValues(t, orig, pp.String())

	err = pp.UnmarshalJSON([]byte(`"1972-01-01"`))
	assert.Error(t, err)

	b, err = pp.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, bj, b)

	// dateOriginal := USDate(time.Date(2014, 10, 10, 0, 0, 0, 0, time.UTC))

	// bsonData, err := bson.Marshal(&dateOriginal)
	// assert.NoError(t, err)

	// var dateCopy USDate
	// err = bson.Unmarshal(bsonData, &dateCopy)
	// assert.NoError(t, err)
	// assert.Equal(t, dateOriginal, dateCopy)
}

func TestUSDate_Scan(t *testing.T) {
	ref := time.Now().Truncate(24 * time.Hour).UTC()
	date, str := USDate(ref), ref.Format(RFC3339FullDate)

	values := []interface{}{str, []byte(str), ref}
	for _, value := range values {
		result := USDate{}
		(&result).Scan(value)
		assert.Equal(t, date, result, "value: %#v", value)
	}

	dd := USDate{}
	err := dd.Scan(nil)
	assert.NoError(t, err)
	assert.Equal(t, USDate{}, dd)

	err = dd.Scan(19700101)
	assert.Error(t, err)
}

func TestUSDate_Value(t *testing.T) {
	ref := time.Now().Truncate(24 * time.Hour).UTC()
	date := USDate(ref)
	dbv, err := date.Value()
	assert.NoError(t, err)
	assert.EqualValues(t, dbv, ref.Format("2006-01-02"))
}

func TestUSDate_IsUSDate(t *testing.T) {
	tests := []struct {
		value string
		valid bool
	}{
		{"12/22/2017", true},
		{"1/1/2017", false},
		{"13/22/17", false},
		{"02/29/2017", false}, // not a valid date : 2017 is not a leap year
		{"02/29/1900", false}, // not a valid date : 1900 is not a leap year
		{"02/29/2100", false}, // not a valid date : 2100 is not a leap year
		{"02/29/2000", true},  // a valid date : 2000 is a leap year
		{"02/29/2400", true},  // a valid date : 2000 is a leap year
		{"13/22/2017", false},
		{"12/32/2017", false},
		{"12/32/20171", false},
		{"MM/DD/YYYY", false},
		{"20/17/2017", false},
		{"2017-12-22T01:02:03Z", false},
	}
	for _, test := range tests {
		assert.Equal(t, test.valid, IsUSDate(test.value), "value [%s] should be valid: [%t]", test.value, test.valid)
	}
}
