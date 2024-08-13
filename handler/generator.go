package handler

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

const (
	sequencebits = 2
	timebits     = 32

	maxsequencebits = (1 << sequencebits) - 1
	maxtimebits     = (1 << timebits) - 1

	sequenceshift = timebits
)

type snowflake struct {
	sequence int64
	time     int64
	mux      sync.Mutex
}

func Newsnowflake() *snowflake {
	return &snowflake{
		sequence: 0,
		time:     -1,
	}
}

func (s *snowflake) Generator() int64 {
	s.mux.Lock()
	defer s.mux.Unlock()

	epoch := time.Date(2024, time.August, 1, 0, 0, 0, 0, time.UTC)
	now := time.Now().UnixMilli() - epoch.UnixMilli()

	if now > maxtimebits {
		fmt.Println("Timestamp exceeds 32-bit range!")
		return -1
	}

	if now == s.time {
		s.sequence = (s.sequence + 1) & maxsequencebits
		if s.sequence == 0 {
			for now <= s.time {
				now = time.Now().UnixMilli() - epoch.UnixMilli()
			}
		}
	} else {
		s.sequence = 0
	}

	s.time = now

	id := (s.sequence << sequenceshift) | s.time

	return id
}

func Generateshortkey() string {
	sr := Newsnowflake()
	id := sr.Generator()
	id1 := fmt.Sprintf("%034b", id)
	desimalid, err := strconv.ParseInt(id1, 2, 64)
	if err != nil {
		fmt.Println("err:",err)
	}
	desimalid1 := strconv.FormatInt(desimalid, 36)
	return string(desimalid1)
}
