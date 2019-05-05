//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"sync"
	"time"
)

const (
	maxServiceTime   = 10
	totalServiceTime = 10
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
	mutex     sync.Mutex
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	if u.IsPremium {
		process()
		return true
	}

	done := make(chan bool)
	go func() {
		process()
		done <- true
	}()

	select {
	case <-done:
		return true
	case <-time.After(maxServiceTime * time.Second):
		return false
	}
}

func HandleRequest2(process func(), u *User) bool {
	if u.IsPremium {
		process()
		return true
	}
	done := make(chan bool)
	go func() {
		process()
		done <- true
	}()
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-done:
			return true
		case <-ticker.C:
			u.mutex.Lock()
			u.TimeUsed++
			u.mutex.Unlock()
			if u.TimeUsed > totalServiceTime {
				return false
			}
		}
	}
}

func main() {
	RunMockServer()
}
