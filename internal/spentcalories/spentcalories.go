package spentcalories

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	// TODO: реализовать функцию
	parts := strings.Split(strings.TrimSpace(data), ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("invalid format: expected '<steps>,<activity>,<training duration>'")
	}

	steps, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, "", 0, fmt.Errorf("invalid steps data %q: %w", parts[0], err)
	}

	if steps < 0 {
		return 0, "", 0, fmt.Errorf("invalid steps data %q", parts[0])
	}

	activity := strings.TrimSpace(parts[1])
	if activity == "" {
		return 0, "", 0, fmt.Errorf("activity must not be empty")
	}

	durStr := strings.TrimSpace(parts[2])
	if durStr == "" {
		return 0, "", 0, fmt.Errorf("duration part is empty")
	}

	dur, err := time.ParseDuration(durStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("invalid duration %q: %w", durStr, err)
	}

	return steps, activity, dur, nil

}

func distance(steps int, height float64) float64 {
	// TODO: реализовать функцию
	stepLength := height * stepLengthCoefficient
	d := float64(steps) * stepLength
	dist := d / float64(mInKm)
	return dist
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	// TODO: реализовать функцию
	if duration < 0 {
		return 0
	}
	dist := distance(steps, height)
	meanSpeed := dist / duration.Hours()
	return meanSpeed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	// TODO: реализовать функцию
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}

	var calories float64
	switch activity {
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", fmt.Errorf("running calories calc failed: %w", err)
		}
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", fmt.Errorf("walking calories calc failed: %w", err)
		}
	default:
		return "", fmt.Errorf("unknown activity: %q", activity)
	}

	distanceKm := distance(steps, height)
	speedKmH := meanSpeed(steps, height, duration)

	hours := duration.Hours()

	result := fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f",
		activity,
		hours,
		distanceKm,
		speedKmH,
		calories,
	)

	return result, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// TODO: реализовать функцию
	switch {
	case steps < 0:
		return 0, fmt.Errorf("invalid steps: %d", steps)
	case weight < 0:
		return 0, fmt.Errorf("invalid weight: %.2f", weight)
	case height < 0:
		return 0, fmt.Errorf("invalid height: %.2f", height)
	case duration < 0:
		return 0, fmt.Errorf("invalid duration: %v", duration)
	}
	meanSpeed := meanSpeed(steps, height, duration)
	durationMinutes := duration.Minutes()
	calories := (weight * meanSpeed * durationMinutes) / minInH
	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// TODO: реализовать функцию
	switch {
	case steps < 0:
		return 0, fmt.Errorf("invalid steps: %d", steps)
	case weight < 0:
		return 0, fmt.Errorf("invalid weight: %.2f", weight)
	case height < 0:
		return 0, fmt.Errorf("invalid height: %.2f", height)
	case duration < 0:
		return 0, fmt.Errorf("invalid duration: %v", duration)
	}
	meanSpeed := meanSpeed(steps, height, duration)
	durationMinutes := duration.Minutes()
	calories := (weight * meanSpeed * durationMinutes) / minInH
	walkingCalories := calories * walkingCaloriesCoefficient
	return walkingCalories, nil
}
