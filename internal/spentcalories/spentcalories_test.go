package spentcalories

import (
	"testing"
	"time"
)

func TestParseTraining(t *testing.T) {
	steps, activity, dur, err := parseTraining("3000,Бег,45m")
	if err != nil {
		t.Errorf("не ожидали ошибку, получили: %v", err)
	}
	if steps != 3000 {
		t.Errorf("ожидали 3000 шагов, получили %d", steps)
	}
	if activity != "Бег" {
		t.Errorf("ожидали 'Бег', получили %s", activity)
	}
	if dur != 45*time.Minute {
		t.Errorf("ожидали 45m, получили %v", dur)
	}
}

func TestRunningSpentCalories(t *testing.T) {
	calories, err := RunningSpentCalories(3000, 70, 1.75, 45*time.Minute)
	if err != nil {
		t.Errorf("не ожидали ошибку, получили: %v", err)
	}
	if calories <= 0 {
		t.Error("ожидали положительное значение калорий")
	}
}

func TestWalkingSpentCalories(t *testing.T) {
	calories, err := WalkingSpentCalories(3000, 70, 1.75, 45*time.Minute)
	if err != nil {
		t.Errorf("не ожидали ошибку, получили: %v", err)
	}
	if calories <= 0 {
		t.Error("ожидали положительное значение калорий")
	}
}

func TestTrainingInfo(t *testing.T) {
	info, err := TrainingInfo("3000,Бег,45m", 70, 1.75)
	if err != nil {
		t.Errorf("не ожидали ошибку, получили: %v", err)
	}
	if info == "" {
		t.Error("ожидали непустую строку результата")
	}
}
