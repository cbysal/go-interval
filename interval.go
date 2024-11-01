package interval

import (
	"cmp"
	"fmt"
	"slices"
)

type Interval[T cmp.Ordered] struct {
	Begin, End T
}

func (interval Interval[T]) String() string {
	return fmt.Sprintf("[%v %v]", interval.Begin, interval.End)
}

type IntervalSet[T cmp.Ordered] struct {
	intervals []T
}

func NewIntervalSet[T cmp.Ordered]() *IntervalSet[T] {
	return &IntervalSet[T]{
		intervals: make([]T, 0),
	}
}

func (set *IntervalSet[T]) Add(interval Interval[T]) {
	if interval.Begin >= interval.End {
		return
	}
	left, _ := slices.BinarySearch(set.intervals, interval.Begin)
	right, ok := slices.BinarySearch(set.intervals, interval.End)
	if ok {
		right++
	}
	set.intervals = slices.Delete(set.intervals, left, right)
	switch {
	case left%2 == 0 && right%2 == 0:
		set.intervals = slices.Insert(set.intervals, left, interval.Begin, interval.End)
	case left%2 == 0:
		set.intervals = slices.Insert(set.intervals, left, interval.Begin)
	case right%2 == 0:
		set.intervals = slices.Insert(set.intervals, left, interval.End)
	}
}

func (set *IntervalSet[T]) Remove(interval Interval[T]) {
	if interval.Begin >= interval.End {
		return
	}
	left, _ := slices.BinarySearch(set.intervals, interval.Begin)
	right, ok := slices.BinarySearch(set.intervals, interval.End)
	if ok {
		right++
	}
	set.intervals = slices.Delete(set.intervals, left, right)
	switch {
	case left%2 != 0 && right%2 != 0:
		set.intervals = slices.Insert(set.intervals, left, interval.Begin, interval.End)
	case left%2 != 0:
		set.intervals = slices.Insert(set.intervals, left, interval.Begin)
	case right%2 != 0:
		set.intervals = slices.Insert(set.intervals, left, interval.End)
	}
}

func (set *IntervalSet[T]) ContainsAll(interval Interval[T]) bool {
	if interval.Begin >= interval.End {
		return true
	}
	left, ok := slices.BinarySearch(set.intervals, interval.Begin)
	right, _ := slices.BinarySearch(set.intervals, interval.End)
	if ok {
		left++
	}
	return left == right && left%2 != 0
}

func (set *IntervalSet[T]) ContainsAny(interval Interval[T]) bool {
	if interval.Begin >= interval.End {
		return false
	}
	left, ok := slices.BinarySearch(set.intervals, interval.Begin)
	right, _ := slices.BinarySearch(set.intervals, interval.End)
	if ok {
		left++
	}
	return left < right || left%2 != 0
}

func (set *IntervalSet[T]) Intervals() []Interval[T] {
	intervals := make([]Interval[T], len(set.intervals)/2)
	for i := range intervals {
		intervals[i] = Interval[T]{set.intervals[i*2], set.intervals[i*2+1]}
	}
	return intervals
}

func (set *IntervalSet[T]) String() string {
	intervals := make([]Interval[T], len(set.intervals)/2)
	for i := range intervals {
		intervals[i] = Interval[T]{set.intervals[i*2], set.intervals[i*2+1]}
	}
	return fmt.Sprintf("%v", intervals)
}
