// Copyright 2015,2016,2017,2018,2019 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package coupon

import (
	"flag"
	"os"
	"testing"

	"github.com/kasworld/globalgametick"
	"github.com/kasworld/uuidstr"
)

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func TestCoupon_Decode(t *testing.T) {
	cpn, err := New(
		"eb6136fcf2e6ee6e3d5a03d85a82ad33",
		"3600b48daa74fc83b791f08c149e9f63",
	)
	if err != nil {
		t.Logf("%v", err)
		return
	}
	for i := 0; i < 10; i++ {
		src := uuidstr.New()
		dst := uuidstr.New()
		tick := globalgametick.GetGameTick().ToInt64NanoSec()

		t.Logf("%v %v %v %v", src, dst, tick, i)

		cp, _ := cpn.Generate(src, dst, tick, int64(i))
		t.Logf("%v %v", cp, len(cp))

		s, d, i1, i2, err := cpn.Parse(cp)
		t.Logf("%v %v %v %v %v", s, d, i1, i2, err)
	}
}

func TestGetRandomKey(t *testing.T) {
	for i := 0; i < 10; i++ {
		s, err := GenerateRandomString()
		if err != nil {
			t.Logf("%v", err)
		}
		t.Logf("%s", s)
	}
}

func BenchmarkSearch(b *testing.B) {
	cpn, _ := New(
		"eb6136fcf2e6ee6e3d5a03d85a82ad33",
		"3600b48daa74fc83b791f08c149e9f63",
	)
	src := uuidstr.New()
	dst := uuidstr.New()
	tick := globalgametick.GetGameTick().ToInt64NanoSec()
	for i := 0; i < b.N; i++ {
		cp, _ := cpn.Generate(src, dst, tick, int64(i))
		_, _, _, _, _ = cpn.Parse(cp)

	}
}

func BenchmarkSearch2(b *testing.B) {
	src := uuidstr.New()
	dst := uuidstr.New()
	tick := globalgametick.GetGameTick().ToInt64NanoSec()
	for i := 0; i < b.N; i++ {
		cpn, _ := New(
			"eb6136fcf2e6ee6e3d5a03d85a82ad33",
			"3600b48daa74fc83b791f08c149e9f63",
		)
		cp, _ := cpn.Generate(src, dst, tick, int64(i))
		_, _, _, _, _ = cpn.Parse(cp)

	}
}
