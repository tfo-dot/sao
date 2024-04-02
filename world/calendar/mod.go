package calendar

import "fmt"

type Month int

const (
	JAN Month = iota
	FEB
	MAR
	APR
)

type Season int

const (
	SPRING Season = iota
	SUMMER
	AUTUMN
	WINTER
)

type Calendar struct {
	Day   int
	Month Month
	Year  int
	Time  Time
}

type Time struct {
	Hour int
	Tick int
}

func (c *Calendar) AddDay() {
	c.Day++

	if c.Day > 30 {
		c.Day = 1
		c.AddMonth()
	}
}

func (c *Calendar) AddMonth() {
	c.Month++

	if c.Month > 3 {
		c.Month = 0
		c.Year++
	}
}

func (c *Calendar) GetSeason() Season {
	switch c.Month {
	case JAN:
		return SPRING
	case FEB:
		return SUMMER
	case MAR:
		return AUTUMN
	case APR:
		return WINTER
	}

	//Impossible
	return -1
}

func (c *Calendar) Tick() {
	c.Time.Tick++

	if c.Time.Tick > 12 {
		c.Time.Tick = 0
		c.Time.Hour++
	}

	if c.Time.Hour > 24 {
		c.Time.Hour = 0
		c.AddDay()
	}
}

func StartCalendar() *Calendar {
	return &Calendar{
		Day:   1,
		Month: JAN,
		Year:  1,
		Time:  Time{0, 0},
	}
}

func (c *Calendar) Copy() *Calendar {
	return &Calendar{
		Day:   c.Day,
		Month: c.Month,
		Year:  c.Year,
		Time:  c.Time,
	}
}

func (c *Calendar) String() string {
	return fmt.Sprintf("%d/%d/%d %d:%d", c.Day, c.Month+1, c.Year, c.Time.Hour, c.Time.Tick)
}
