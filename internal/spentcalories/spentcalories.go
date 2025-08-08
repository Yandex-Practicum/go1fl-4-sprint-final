package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	lenStep                    = 0.65
	mInKm                      = 1000
	minInH                     = 60
	stepLengthCoefficient      = 0.45
	walkingCaloriesCoefficient = 0.5
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, errors.New("неверный формат данных: требуется 3 поля — шаги, тип активности, время")
	}

	stepsStr := strings.TrimSpace(parts[0])
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, "", 0, errors.New("невозможно преобразовать шаги в число")
	}
	if steps <= 0 {
		return 0, "", 0, errors.New("количество шагов должно быть больше нуля")
	}

	activity := strings.TrimSpace(parts[1])
	durationStr := strings.TrimSpace(parts[2])

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, "", 0, errors.New("неверный формат продолжительности")
	}
	if duration <= 0 {
		return 0, "", 0, errors.New("продолжительность должна быть больше нуля")
	}

	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {
	if steps <= 0 || height <= 0 {
		return 0
	}
	stepLength := height * stepLengthCoefficient
	return float64(steps) * stepLength / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	return distance(steps, height) / duration.Hours()
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("количество шагов должно быть больше нуля")
	}
	if weight <= 0 {
		return 0, errors.New("вес должен быть больше нуля")
	}
	if height <= 0 {
		return 0, errors.New("рост должен быть больше нуля")
	}
	if duration <= 0 {
		return 0, errors.New("длительность должна быть больше нуля")
	}

	speed := meanSpeed(steps, height, duration)
	return (weight * speed * duration.Minutes()) / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	calories, err := RunningSpentCalories(steps, weight, height, duration)
	if err != nil {
		return 0, err
	}
	return calories * walkingCaloriesCoefficient, nil
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}

	var calories float64
	var activityName string

	switch activity {
	case "Ходьба":
		activityName = "Ходьба"
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	case "Бег":
		activityName = "Бег"
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}

	if err != nil {
		return "", err
	}

	distance := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	return fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activityName,
		duration.Hours(),
		distance,
		speed,
		calories,
	), nil
}