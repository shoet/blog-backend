package clocker

import "time"

type Clocker interface {
	Now() time.Time
}

type RealClocker struct{}

func NewRealClocker() *RealClocker {
	return &RealClocker{}
}

func (c *RealClocker) Now() time.Time {
	return time.Now()
}

type FiexedClocker struct{}

func NewFixedClocker() *FiexedClocker {
	return &FiexedClocker{}
}

func (c *FiexedClocker) Now() time.Time {
	return time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
}
