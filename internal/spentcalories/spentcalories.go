package spentcalories

import (
	"errors"
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
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, errors.New("incorrect data set")
	}

	steps, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil || steps <= 0 {
		return 0, "", 0, errors.New("steps incorrect")
	}

	activityType := strings.TrimSpace(parts[1])
	if activityType == "" {
		return 0, "", 0, errors.New("Type activity incorrect")
	}

	duration, err := time.ParseDuration(strings.TrimSpace(parts[2]))
	if err != nil || duration <= 0 {
		return 0, "", 0, errors.New("duration incorrect")
	}

	return steps, activityType, duration, nil
}

func distance(steps int, height float64) float64 {
	// TODO: реализовать функцию
	return ((height * stepLengthCoefficient) * float64(steps)) / mInKm

}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	// TODO: реализовать функцию
	if duration <= 0 {
		return 0
	}
	calculateDistance := distance(steps, height)

	hours := duration.Hours()

	return calculateDistance / hours
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	// TODO: реализовать функцию
	steps, activityType, duration, err := parseTraining(data)
	if steps <= 0 || activityType == "" || duration <= 0 || err != nil || weight <= 0 || height <= 0 {
		return "", errors.New("input data is not correct")
	}

	switch activityType {
	case "Бег":
		trainigDistance := distance(steps, height)
		trainigSpeed := meanSpeed(steps, height, duration)
		trainigRunCalories, err := RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", errors.New("calorie calculation error")
		}
		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
				activityType,
				duration.Hours(),
				trainigDistance,
				trainigSpeed,
				trainigRunCalories),
			nil
	case "Ходьба":
		trainigDistance := distance(steps, height)
		trainigSpeed := meanSpeed(steps, height, duration)
		trainigWalkCalories, err := WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", errors.New("calorie calculation error")
		}
		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
				activityType,
				duration.Hours(),
				trainigDistance,
				trainigSpeed,
				trainigWalkCalories),
			nil
	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}

}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// TODO: реализовать функцию
	if weight <= 0 || height <= 0 || steps <= 0 || duration <= 0 {
		return 0, errors.New("input data is not correct")
	}
	averageSpeed := meanSpeed(steps, height, duration)

	durationInMinutes := duration.Minutes()

	return (weight * averageSpeed * durationInMinutes) / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// TODO: реализовать функцию
	if weight <= 0 || height <= 0 || steps <= 0 || duration <= 0 {
		return 0, errors.New("input data is not correct")
	}

	averageSpeed := meanSpeed(steps, height, duration)

	durationInMinutes := duration.Minutes()

	return ((weight * averageSpeed * durationInMinutes) / 60) * walkingCaloriesCoefficient, nil
}
