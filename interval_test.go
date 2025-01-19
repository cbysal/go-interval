package interval

import (
	"maps"
	"slices"
	"testing"
)

func mapKeysToIntervals(vals map[int]struct{}) []Interval[int] {
	sortedVals := slices.Collect(maps.Keys(vals))
	slices.Sort(sortedVals)
	intervals := make([]Interval[int], 0)
	left := 0
	for left < len(vals) {
		right := left + 1
		for right < len(vals) && sortedVals[right-1]+1 == sortedVals[right] {
			right++
		}
		intervals = append(intervals, Interval[int]{sortedVals[left], sortedVals[right-1] + 1})
		left = right
	}
	return intervals
}

var (
	standard     []Interval[int]
	alternatives []Interval[int]
)

func init() {
	for i := 0; i < 4; i++ {
		standard = append(standard, Interval[int]{i * 20, i*20 + 10})
	}
	for i := -5; i <= 75; i += 5 {
		alternatives = append(alternatives, Interval[int]{i - 2, i + 2})
	}
}

func TestIntervalString(t *testing.T) {
	interval := Interval[int]{0, 10}
	output := interval.String()
	expect := "[0 10]"
	if output != expect {
		t.Fatalf("expect %v, got %v", expect, output)
	}
}

func TestIntervalSetAdd(t *testing.T) {
	for begin := -10; begin <= 80; begin++ {
		for end := -10; end <= 80; end++ {
			set := NewIntervalSet[int]()
			for _, interval := range standard {
				set.Add(interval)
			}
			set.Add(Interval[int]{begin, end})
			output := make(map[int]struct{})
			for _, interval := range set.Intervals() {
				for j := interval.Begin; j < interval.End; j++ {
					output[j] = struct{}{}
				}
			}
			if !slices.IsSortedFunc(set.Intervals(), func(a, b Interval[int]) int {
				if a.Begin >= a.End || b.Begin >= b.End {
					return -1
				}
				if a.Begin <= b.End {
					return -1
				}
				return 1
			}) {
				t.Fatalf("testcase [%v %v], invalid intervals %v", begin, end, set.Intervals())
			}
			expect := make(map[int]struct{})
			for _, interval := range standard {
				for j := interval.Begin; j < interval.End; j++ {
					expect[j] = struct{}{}
				}
			}
			for i := begin; i < end; i++ {
				expect[i] = struct{}{}
			}
			if !maps.Equal(output, expect) {
				t.Fatalf("testcase [%v %v], expect %v, got %v", begin, end, mapKeysToIntervals(expect), mapKeysToIntervals(output))
			}
		}
	}
}

func TestIntervalSetRemove(t *testing.T) {
	for begin := -10; begin <= 80; begin++ {
		for end := -10; end <= 80; end++ {
			set := NewIntervalSet[int]()
			for _, interval := range standard {
				set.Add(interval)
			}
			set.Remove(Interval[int]{begin, end})
			output := make(map[int]struct{})
			for _, interval := range set.Intervals() {
				for j := interval.Begin; j < interval.End; j++ {
					output[j] = struct{}{}
				}
			}
			if !slices.IsSortedFunc(set.Intervals(), func(a, b Interval[int]) int {
				if a.Begin >= a.End || b.Begin >= b.End {
					return -1
				}
				if a.Begin <= b.End {
					return -1
				}
				return 1
			}) {
				t.Fatalf("testcase [%v %v], invalid intervals %v", begin, end, set.Intervals())
			}
			expect := make(map[int]struct{})
			for _, interval := range standard {
				for j := interval.Begin; j < interval.End; j++ {
					expect[j] = struct{}{}
				}
			}
			for i := begin; i < end; i++ {
				delete(expect, i)
			}
			if !maps.Equal(output, expect) {
				t.Fatalf("testcase [%v %v], expect %v, got %v", begin, end, mapKeysToIntervals(expect), mapKeysToIntervals(output))
			}
		}
	}
}

func TestIntervalSetDifference(t *testing.T) {
	for i := 0; i < 2<<len(alternatives); i++ {
		input1, input2 := make([]Interval[int], 0), make([]Interval[int], 0)
		for _, interval := range standard {
			input1 = append(input1, interval)
		}
		for j := 0; j < len(alternatives); j++ {
			if i&(1<<j) != 0 {
				input2 = append(input2, alternatives[j])
			}
		}

		expect := make(map[int]struct{})
		for _, interval := range input1 {
			for j := interval.Begin; j < interval.End; j++ {
				expect[j] = struct{}{}
			}
		}
		for _, interval := range input2 {
			for j := interval.Begin; j < interval.End; j++ {
				delete(expect, j)
			}
		}

		output := make(map[int]struct{})
		self := NewIntervalSet[int]()
		for _, interval := range standard {
			self.Add(interval)
		}
		other := NewIntervalSet[int]()
		for j := 0; j < len(alternatives); j++ {
			if i&(1<<j) != 0 {
				other.Add(alternatives[j])
			}
		}
		result := self.Difference(other)
		for _, interval := range result.Intervals() {
			for j := interval.Begin; j < interval.End; j++ {
				output[j] = struct{}{}
			}
		}

		if !slices.IsSortedFunc(result.Intervals(), func(a, b Interval[int]) int {
			if a.Begin >= a.End || b.Begin >= b.End {
				return -1
			}
			if a.Begin <= b.End {
				return -1
			}
			return 1
		}) {
			t.Fatalf("testcase: %v %v, result: %v", self, other, result)
		}
		if !maps.Equal(output, expect) {
			t.Fatalf("testcase: %v %v, expect: %v, result: %v", self, other, mapKeysToIntervals(expect), mapKeysToIntervals(output))
		}
	}
}

func TestIntervalSetContainsAll(t *testing.T) {
	set := NewIntervalSet[int]()
	for _, interval := range standard {
		set.Add(interval)
	}
	values := make(map[int]struct{})
	for _, interval := range set.Intervals() {
		for j := interval.Begin; j < interval.End; j++ {
			values[j] = struct{}{}
		}
	}
	for begin := -10; begin <= 80; begin++ {
		for end := -10; end <= 80; end++ {
			output := set.ContainsAll(Interval[int]{begin, end})
			expect := true
			for i := begin; i < end; i++ {
				if _, ok := values[i]; !ok {
					expect = false
					break
				}
			}
			if output != expect {
				t.Fatalf("testcase [%v %v], expect %v, got %v", begin, end, expect, output)
			}
		}
	}
}

func TestIntervalSetContainsAny(t *testing.T) {
	set := NewIntervalSet[int]()
	for _, interval := range standard {
		set.Add(interval)
	}
	values := make(map[int]struct{})
	for _, interval := range set.Intervals() {
		for j := interval.Begin; j < interval.End; j++ {
			values[j] = struct{}{}
		}
	}
	for begin := -10; begin <= 80; begin++ {
		for end := -10; end <= 80; end++ {
			output := set.ContainsAny(Interval[int]{begin, end})
			expect := false
			for i := begin; i < end; i++ {
				if _, ok := values[i]; ok {
					expect = true
					break
				}
			}
			if output != expect {
				t.Fatalf("testcase [%v %v], expect %v, got %v", begin, end, expect, output)
			}
		}
	}
}

func TestIntervalSetString(t *testing.T) {
	set := NewIntervalSet[int]()
	for _, interval := range standard {
		set.Add(interval)
	}
	output := set.String()
	expect := "[[0 10] [20 30] [40 50] [60 70]]"
	if output != expect {
		t.Fatalf("expect %v, got %v", expect, output)
	}
}
