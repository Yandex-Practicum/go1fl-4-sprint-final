package daysteps

import (
	"testing"
	"time"
)

func TestParsePackage(t *testing.T) {
	steps, dur, err := parsePackage("1000,30m")
	if err != nil {
		t.Errorf("не ожидали ошибку, получили: %v", err)
	}
	if steps != 1000 {
		t.Errorf("ожидали 1000 шагов, получили %d", steps)
	}
	if dur != 30*time.Minute {
		t.Errorf("ожидали 30m, получили %v", dur)
	}
}

func TestDayActionInfo(t *testing.T) {
	result := DayActionInfo("1000,30m", 70, 1.75)
	if result == "" {
		t.Error("ожидали непустой результат, получили пустую строку")
	}
}
