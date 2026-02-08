package spentcalories

import (
	"fmt"
	"log"
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
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("неверный формат строки")
	}

	stepsStr := strings.TrimSpace(parts[0])
	activity := strings.TrimSpace(parts[1])
	durStr := strings.TrimSpace(parts[2])

	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, "", 0, err
	}
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("некорректное количество шагов")
	}

	duration, err := time.ParseDuration(durStr)
	if err != nil {
		return 0, "", 0, err
	}
	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("некорректная продолжительность")
	}

	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {

	stepLength := height * stepLengthCoefficient

	return float64(steps) * stepLength / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	dist := distance(steps, height)

	hours := duration.Hours()

	return dist / hours
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		err := fmt.Errorf("некорректные входные параметры")
		log.Println(err)
		return "", err
	}

	dist := distance(steps, height)             // км
	speed := meanSpeed(steps, height, duration) // км/ч
	hours := duration.Hours()                   // часы

	var calories float64

	switch activity {
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
			return "", err
		}
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
			return "", err
		}
	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}

	out := fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activity, hours, dist, speed, calories,
	)

	return out, nil
}

func RunningSpentCalories(
	steps int,
	weight, height float64,
	duration time.Duration,
) (float64, error) {

	if steps <= 0 {
		return 0, fmt.Errorf("некорректное количество шагов")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("некорректный вес")
	}
	if height <= 0 {
		return 0, fmt.Errorf("некорректный рост")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("некорректная время")
	}

	speed := meanSpeed(steps, height, duration)

	durationMinutes := duration.Minutes()

	calories := (weight * speed * durationMinutes) / minInH

	return calories, nil
}

func WalkingSpentCalories(
	steps int,
	weight, height float64,
	duration time.Duration,
) (float64, error) {

	if steps <= 0 {
		return 0, fmt.Errorf("некорректное количество шагов")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("некорректный вес")
	}
	if height <= 0 {
		return 0, fmt.Errorf("некорректный рост")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("некорректная продолжительность")
	}

	speed := meanSpeed(steps, height, duration)

	durationMinutes := duration.Minutes()

	calories := (weight * speed * durationMinutes) / minInH

	calories *= walkingCaloriesCoefficient

	return calories, nil
}
