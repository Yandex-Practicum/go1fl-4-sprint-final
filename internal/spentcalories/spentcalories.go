package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	stepLengthCoefficient      = 0.45
	walkingCaloriesCoefficient = 0.5
	mInKm                      = 1000
	minInH                     = 60
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("неверный формат тренировки")
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil || steps <= 0 {
		return 0, "", 0, fmt.Errorf("некорректные шаги")
	}

	duration, err := time.ParseDuration(parts[2])
	if err != nil || duration <= 0 {
		return 0, "", 0, fmt.Errorf("некорректная продолжительность")
	}

	return steps, parts[1], duration, nil
}

func distance(steps int, height float64) float64 {
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
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("некорректные параметры")
	}

	speed := meanSpeed(steps, height, duration)
	minutes := duration.Minutes()

	return (weight * speed * minutes) / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("некорректные параметры")
	}

	speed := meanSpeed(steps, height, duration)
	minutes := duration.Minutes()

	calories := (weight * speed * minutes) / minInH
	return calories * walkingCaloriesCoefficient, nil
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	var calories float64

	switch activity {
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}

	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activity,
		duration.Hours(),
		dist,
		speed,
		calories,
	), nil
}
