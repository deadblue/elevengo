package impl

import (
	"math/rand"
	"time"
)

type Valve struct {
	// Is valve enabled?
	enabled bool
	// Last API calling finish time
	lastTime int64
	// Cooldown duration range
	cdMin, cdMax uint
}

func (v *Valve) ClockIn() {
	if !v.enabled {
		return
	}
	// Save last time
	v.lastTime = time.Now().UnixMilli()
}

func (v *Valve) Wait() {
	if !v.enabled {
		return
	}
	cd := v.getCooldownDuration()
	if cd == 0 {
		return
	}
	sleepDuration := v.lastTime + cd - time.Now().UnixMilli()
	if sleepDuration > 0 {
		time.Sleep(time.Duration(sleepDuration) * time.Millisecond)
	}
}

func (v *Valve) getCooldownDuration() int64 {
	// Skip invalid cooldown duration
	if v.cdMax == 0 || v.cdMax < v.cdMin {
		return 0
	}
	// Generate random duration in range
	var duration int64
	if v.cdMax == v.cdMin {
		duration = int64(v.cdMax)
	} else {
		duration = rand.Int63n(int64(v.cdMax-v.cdMin)) + int64(v.cdMin)
	}
	return duration
}
