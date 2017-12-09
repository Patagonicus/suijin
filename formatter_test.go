package suijin_test

import (
	"fmt"
	"net"
	"time"

	"github.com/Patagonicus/suijin"
	"github.com/Patagonicus/suijin/mock"
)

func ExampleFormatter() {
	formatter := suijin.NewTextFormatter(suijin.FormatterConfig{
		Clock: mock.NewClock(mustParse("2006-01-02T15:04:05Z"), time.Second),
	})

	b, err := formatter.Format(suijin.Message{
		suijin.InfoLevel,
		"application started succesfully",
		nil,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))

	b, err = formatter.Format(suijin.Message{
		suijin.DebugLevel,
		"client connected",
		suijin.Fields{
			"user": "foobar",
			"ip":   net.IPv4(127, 0, 0, 1),
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))

	// Output:
	// 2006-01-02T15:04:05Z info: application started succesfully
	// 2006-01-02T15:04:06Z debug: client connected "ip"="127.0.0.1" "user"="foobar"
}

func mustParse(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}
