// Copyright 2019 Adam Shannon
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package main

import (
	"testing"
	"time"
)

type testcase struct {
	when time.Time
	open bool
}

func TestMarketOpen(t *testing.T) {
	loc, _ := time.LoadLocation("America/New_York")
	cases := []testcase{
		{
			when: time.Date(2019, time.January, 2, 12, 30, 0, 0, loc),
			open: true,
		},
		{
			// open at 9am
			when: time.Date(2019, time.January, 2, 8, 59, 0, 0, loc),
			open: false,
		},
		{
			// open at 9am
			when: time.Date(2019, time.January, 2, 9, 00, 0, 0, loc),
			open: true,
		},
		{
			// close at 4:30pm
			when: time.Date(2019, time.January, 2, 16, 30, 0, 0, loc),
			open: true,
		},
		{
			// close at 4:30pm
			when: time.Date(2019, time.January, 2, 16, 31, 0, 0, loc),
			open: false,
		},
		{
			// 2019-01-05 is a Saturday
			when: time.Date(2019, time.January, 5, 12, 30, 0, 0, loc),
			open: false,
		},
		{
			// 2019-01-06 is a Sunday
			when: time.Date(2019, time.January, 6, 12, 30, 0, 0, loc),
			open: false,
		},
		{
			// Outside market hours
			when: time.Date(2019, time.January, 2, 1, 30, 0, 0, loc),
			open: false,
		},
		{
			// Outside market hours
			when: time.Date(2019, time.January, 2, 18, 30, 0, 0, loc),
			open: false,
		},
	}
	for i := range cases {
		if v := marketOpen(cases[i].when); v != cases[i].open {
			t.Errorf("%v: expected %v got %v", cases[i].when, cases[i].open, v)
		}
	}
}
