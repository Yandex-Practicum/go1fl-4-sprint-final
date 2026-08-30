package spentcalories

import (
	"errors"
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
		return 0, "", 0, fmt.Errorf("неверный формат строки: ожидается 3 части (шаги,активность,длительность), получено %d частей", len(parts))
	}
	stepsStr := parts[0]
	activity := parts[1]
	durationStr := parts[2]

	steps, err := strconv.Atoi(stepsStr)
	if err != nil || steps <= 0 {
		return 0, "", 0, errors.New("количество шагов должно быть больше 0")
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil || duration <= 0 {
		return 0, "", 0, errors.New("длительность должна быть больше 0")
	}
	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {
	if steps < 0 || height <= 0 {
		return 0
	}

	stepsFloat := float64(steps)
	stepLength := height * stepLengthCoefficient
	totalMeters := stepsFloat * stepLength
	return totalMeters / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if steps < 0 || height <= 0 || duration <= 0 {
		return 0
	}
	distKm := distance(steps, height)
	hours := duration.Hours()
	return distKm / hours
}
func TrainingInfo(data string, weight, height float64) (string, error) {

	steps, activityType, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
	}

	switch activityType {

	case "Ходьба":
		distance := distance(steps, height)
		meanSpeed := meanSpeed(steps, height, duration)
		calories, err := WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
		}

		message := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", activityType, duration.Hours(), distance, meanSpeed, calories)
		return message, nil

	case "Бег":
		distance := distance(steps, height)
		meanSpeed := meanSpeed(steps, height, duration)
		calories, err := RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
		}

		message := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", activityType, duration.Hours(), distance, meanSpeed, calories)
		return message, nil

	default:
		return "", errors.New("неизвестный тип тренировки")

	}

}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps < 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("некорректные параметры тренировки")
	}

	speed := meanSpeed(steps, height, duration)

	calories := (weight * speed * duration.Minutes()) / minInH
	return calories, nil
}
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps < 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("некорректные параметры тренировки")
	}

	speed := meanSpeed(steps, height, duration)
	calories := (weight * speed * duration.Minutes()) / minInH * walkingCaloriesCoefficient

	return calories, nil
}
