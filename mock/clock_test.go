package mock_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Patagonicus/suijin/legacy"
	"github.com/Patagonicus/suijin/mock"
)

const goroutineCount = 20
const iterationCount = 1000

func ExampleClock() {
	c := mock.NewClock(mustParseTime("5:06PM"), time.Minute)
	fmt.Println(c.Now().Format(time.Kitchen))
	fmt.Println(c.Now().Format(time.Kitchen))
	fmt.Println(c.Now().Format(time.Kitchen))
	fmt.Println(c.Now().Format(time.Kitchen))
	// Output:
	// 5:06PM
	// 5:07PM
	// 5:08PM
	// 5:09PM
}

func ExampleClock_Next() {
	c := mock.NewClock(mustParseTime("5:06PM"), time.Minute)
	fmt.Println(c.Next().Format(time.Kitchen))
	fmt.Println(c.Now().Format(time.Kitchen))
	fmt.Println(c.Now().Format(time.Kitchen))
	fmt.Println(c.Next().Format(time.Kitchen))
	fmt.Println(c.Next().Format(time.Kitchen))
	fmt.Println(c.Now().Format(time.Kitchen))
	// Output:
	// 5:06PM
	// 5:06PM
	// 5:07PM
	// 5:08PM
	// 5:08PM
	// 5:08PM
}

func TestConcurrent(t *testing.T) {
	start := mustParseTime("0:00AM")
	c := mock.NewClock(start, time.Minute)
	ch := make(chan []time.Time)

	for i := 0; i < goroutineCount; i++ {
		go func() {
			result := make([]time.Time, 0, iterationCount)
			for i := 0; i < iterationCount; i++ {
				result = append(result, c.Now())
			}
			ch <- result
		}()
	}

	result := make([]time.Time, 0, goroutineCount*iterationCount)
	for i := 0; i < goroutineCount; i++ {
		result = append(result, (<-ch)...)
	}

	legacy.Slice(result, func(i, j int) bool {
		return result[i].Before(result[j])
	})

	if result[0] != start {
		t.Errorf("wrong start time, expected %s but got %s", start, result[0])
	}
	for i := 1; i < len(result); i++ {
		diff := result[i].Sub(result[i-1])
		if diff != time.Minute {
			t.Errorf("difference between %d and %d is %s, not a minute", i-1, i, diff)
			// if there's one problem here, there are probably a lot more. Stop the test so we don't spam the output.
			t.FailNow()
		}
	}
}

func mustParseTime(t string) time.Time {
	result, err := time.Parse(time.Kitchen, t)
	if err != nil {
		panic(err)
	}
	return result
}
