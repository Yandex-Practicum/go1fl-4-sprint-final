package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы
const (
	//lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, errors.New("invalid format: expected 3 parts")
	}

	stepsStr := parts[0]
	activityType := parts[1]
	durationStr := parts[2]

	if stepsStr != strings.TrimSpace(stepsStr) ||
		activityType != strings.TrimSpace(activityType) ||
		durationStr != strings.TrimSpace(durationStr) {
		return 0, "", 0, errors.New("invalid format: whitespace detected")
	}

	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, "", 0, err
	}
	if steps <= 0 {
		return 0, "", 0, errors.New("steps must be positive")
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, "", 0, err
	}
	if duration <= 0 {
		return 0, "", 0, errors.New("duration must be positive")
	}

	return steps, activityType, duration, nil
}

func distance(steps int, height float64) float64 {
	stepLen := height * stepLengthCoefficient
	totalMeters := float64(steps) * stepLen
	return totalMeters / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration.Hours() <= 0 {
		return 0
	}
	distKm := distance(steps, height)
	return distKm / duration.Hours()
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("steps must be positive")
	}
	if weight <= 0 {
		return 0, errors.New("weight must be positive")
	}
	if height <= 0 {
		return 0, errors.New("height must be positive")
	}
	if duration <= 0 {
		return 0, errors.New("duration must be positive")
	}

	speed := meanSpeed(steps, height, duration)
	calories := (weight * speed * duration.Minutes()) / minInH

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("steps must be positive")
	}
	if weight <= 0 {
		return 0, errors.New("weight must be positive")
	}
	if height <= 0 {
		return 0, errors.New("height must be positive")
	}
	if duration <= 0 {
		return 0, errors.New("duration must be positive")
	}

	speed := meanSpeed(steps, height, duration)
	baseCalories := (weight * speed * duration.Minutes()) / minInH
	calories := baseCalories * walkingCaloriesCoefficient

	return calories, nil
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activityType, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	var calories float64

	switch activityType {
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	if err != nil {
		log.Println(err)
		return "", err
	}

	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	result := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activityType,
		duration.Hours(),
		dist,
		speed,
		calories,
	)

	return result, nil
}
