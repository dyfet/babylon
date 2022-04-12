// Copyright (C) 2021-2022 David Sugar <tychosoft@gmail.com>.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package lib

import (
	"encoding/json"
	"errors"
	"time"
)

// lib.Duration allows json marshalling
type Duration time.Duration

// marshall as duration string (xhymzs...)
func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).Truncate(time.Second).String())
}

// Parse duration
func (d Duration) Parse(text string) (Duration, error) {
	tmp, err := time.ParseDuration(text)
	return Duration(tmp.Truncate(time.Second)), err
}

// Duration as string
func (d Duration) String() string {
	return time.Duration(d).Truncate(time.Second).String()
}

// unmarshall from float or string form
func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		*d = Duration(time.Duration(value).Truncate(time.Second))
		return nil
	case string:
		tmp, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		*d = Duration(tmp.Truncate(time.Second))
		return nil
	default:
		return errors.New("invalid duration")
	}
}
