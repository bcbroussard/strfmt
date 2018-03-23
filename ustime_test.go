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
	// "bytes"
	"testing"
	"time"

	// "gopkg.in/mgo.v2/bson"

	"github.com/stretchr/testify/assert"
)

var (
	pst, _ = time.LoadLocation("America/Los_Angeles")

	usTestCases = []struct {
		in    []byte    // externally sourced data -- to be unmarshalled
		time  time.Time // its representation in time.Time
		dbstr string    // its db representation
		str   string    // its marshalled representation
	}{
		{[]byte("04/15/2014 08:00am"), time.Date(2014, 4, 15, 8, 0, 0, 0, pst), "2014-04-15T08:00:00.000-07:00", "04/15/2014 08:00am"},
		{[]byte("08/18/2011 07:03pm"), time.Date(2011, 8, 18, 19, 3, 00, 0, pst), "2011-08-18T19:03:00.000-07:00", "08/18/2011 07:03pm"},
		{[]byte("12/05/2014 07:30pm"), time.Date(2014, 12, 05, 19, 30, 00, 0, pst), "2014-12-05T19:30:00.000-08:00", "12/05/2014 07:30pm"},
		{[]byte("2014-12-05 19:30:00"), time.Date(2014, 12, 05, 19, 30, 00, 0, pst), "2014-12-05T19:30:00.000-08:00", "12/05/2014 07:30pm"},
		{[]byte("01/01/2001"), time.Date(2001, 1, 1, 0, 0, 0, 0, pst), "2001-01-01T00:00:00.000-08:00", "01/01/2001 12:00am"},
		{[]byte("2001-01-01"), time.Date(2001, 1, 1, 0, 0, 0, 0, pst), "2001-01-01T00:00:00.000-08:00", "01/01/2001 12:00am"},
		// {[]byte("0000-00-00 00:00:00"), time.Unix(0, 0), "1970-01-01T00:00:00.000-00:00", "01/01/1970 00:00am"},

		// {[]byte(""), time.Unix(0, 0).In(pst), "1970-01-01T00:00:00.000-00:00", "01/01/1970 00:00am"},
		// {[]byte(nil), time.Unix(0, 0).In(pst), "1970-01-01T00:00:00.000-00:00", "01/01/1970 00:00am"},
	}
)

func TestNewUSDateTime(t *testing.T) {
	assert.EqualValues(t, time.Unix(0, 0).UTC(), NewUSDateTime())
}

func TestParseUSDateTime_errorCases(t *testing.T) {
	_, err := ParseUSDateTime("yada")
	assert.Error(t, err)
}

// TestParseUSDateTime tests the full cycle:
// parsing -> marshalling -> unmarshalling / scanning
func TestParseUSDateTime_fullCycle(t *testing.T) {
	for caseNum, example := range usTestCases {
		t.Logf("Case #%d", caseNum)

		sin := string(example.in)

		parsed, err := ParseUSDateTime(sin)
		assert.NoError(t, err)
		assert.EqualValues(t, USDateTime(example.time).StringWithZone(), parsed.StringWithZone())

		mt, err := parsed.MarshalText()
		assert.NoError(t, err)
		assert.Equal(t, []byte(example.str), mt)

		if sin != "" {
			v := IsUSDateTime(sin)
			assert.True(t, v, sin)

		} else {
			t.Logf("IsUSDateTime() skipped for empty testcases")
		}

		pp := NewUSDateTime()
		err = pp.UnmarshalText([]byte(sin))
		assert.NoError(t, err)
		assert.EqualValues(t, example.time.String(), time.Time(pp).String())

		pp = NewUSDateTime()
		err = pp.Scan([]byte(example.dbstr))
		assert.NoError(t, err)
		assert.EqualValues(t, USDateTime(example.time).StringWithZone(), pp.StringWithZone())
	}
}

func TestUSDateTime_IsUSDateTime_errorCases(t *testing.T) {
	v := IsUSDateTime("zor")
	assert.False(t, v)

	v = IsUSDateTime("zorg")
	assert.False(t, v)

	v = IsUSDateTime("zorgTx")
	assert.False(t, v)

	v = IsUSDateTime("1972-12-31Tx")
	assert.False(t, v)

	v = IsUSDateTime("1972-12-31T24:40:00.000Z")
	assert.False(t, v)

	v = IsUSDateTime("1972-12-31T23:63:00.000Z")
	assert.False(t, v)

	v = IsUSDateTime("1972-12-31T23:59:60.000Z")
	assert.False(t, v)

}
func TestUSDateTime_UnmarshalText_errorCases(t *testing.T) {
	pp := NewUSDateTime()
	err := pp.UnmarshalText([]byte("yada"))
	assert.Error(t, err)
	err = pp.UnmarshalJSON([]byte("yada"))
	assert.Error(t, err)
}

func TestUSDateTime_UnmarshalText(t *testing.T) {
	for caseNum, example := range usTestCases {
		t.Logf("Case #%d", caseNum)
		pp := NewUSDateTime()
		err := pp.UnmarshalText(example.in)
		assert.NoError(t, err)
		assert.EqualValues(t, example.time, pp)

		// Other way around
		val, erv := pp.Value()
		assert.NoError(t, erv)
		assert.EqualValues(t, example.dbstr, val)

	}
}
func TestUSDateTime_UnmarshalJSON(t *testing.T) {
	for caseNum, example := range usTestCases {
		t.Logf("Case #%d", caseNum)
		pp := NewUSDateTime()
		err := pp.UnmarshalJSON(esc(example.in))
		assert.NoError(t, err)
		assert.EqualValues(t, example.time, pp)
	}

	// Check UnmarshalJSON failure with no lexed items
	pp := NewUSDateTime()
	err := pp.UnmarshalJSON([]byte("zorg emperor"))
	assert.Error(t, err)

	// Check lexer failure
	err = pp.UnmarshalJSON([]byte(`"zorg emperor"`))
	assert.Error(t, err)
}

// func esc(v []byte) []byte {
// 	var buf bytes.Buffer
// 	buf.WriteByte('"')
// 	buf.Write(v)
// 	buf.WriteByte('"')
// 	return buf.Bytes()
// }

func TestUSDateTime_MarshalText(t *testing.T) {
	for caseNum, example := range usTestCases {
		t.Logf("Case #%d", caseNum)
		dt := USDateTime(example.time)
		mt, err := dt.MarshalText()
		assert.NoError(t, err)
		assert.Equal(t, []byte(example.str), mt)
	}
}
func TestUSDateTime_MarshalJSON(t *testing.T) {
	for caseNum, example := range usTestCases {
		t.Logf("Case #%d", caseNum)
		dt := USDateTime(example.time)
		bb, err := dt.MarshalJSON()
		assert.NoError(t, err)
		assert.EqualValues(t, esc([]byte(example.str)), bb)
	}
}

func TestUSDateTime_Scan(t *testing.T) {
	for caseNum, example := range usTestCases {
		t.Logf("Case #%d", caseNum)

		pp := NewUSDateTime()
		err := pp.Scan(example.dbstr)
		assert.NoError(t, err)
		assert.Equal(t, USDateTime(example.time).StringWithZone(), pp.StringWithZone())

		pp = NewUSDateTime()
		err = pp.Scan(string(example.dbstr))
		assert.NoError(t, err)
		assert.Equal(t, USDateTime(example.time).StringWithZone(), pp.StringWithZone())

		pp = NewUSDateTime()
		err = pp.Scan(example.time)
		assert.NoError(t, err)
		assert.Equal(t, USDateTime(example.time).StringWithZone(), pp.StringWithZone())
	}
}

func TestUSDateTime_Scan_Failed(t *testing.T) {
	pp := NewUSDateTime()
	zero := NewUSDateTime()

	err := pp.Scan(nil)
	assert.NoError(t, err)
	// Zero values differ...
	//assert.Equal(t, zero, pp)
	assert.Equal(t, USDateTime{}, pp)

	err = pp.Scan("")
	assert.NoError(t, err)
	assert.Equal(t, zero, pp)

	err = pp.Scan(int64(0))
	assert.Error(t, err)

	err = pp.Scan(float64(0))
	assert.Error(t, err)
}

// func TestUSDateTime_BSON(t *testing.T) {
// 	for caseNum, example := range usTestCases {
// 		t.Logf("Case #%d", caseNum)
// 		dt := USDateTime(example.time)

// 		bsonData, err := bson.Marshal(&dt)
// 		assert.NoError(t, err)

// 		var dtCopy USDateTime
// 		err = bson.Unmarshal(bsonData, &dtCopy)
// 		assert.NoError(t, err)
// 		assert.Equal(t, dt, dtCopy)
// 	}
// }
