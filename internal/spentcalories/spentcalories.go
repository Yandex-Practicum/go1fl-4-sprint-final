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
	//a := parts[0]
	//b := parts[1]
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("Not right parts")
	}
	a, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("Error conv a %w", err)
	}
	if a <= 0 {
		return 0, "", 0, fmt.Errorf("0 - quantity of steps")
	}
	active_type := (parts[1])

	//a := strconv.ParseInt(parts[0], 10, 64)
	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("Error conv duration %w", err)
	} else {
		return a, active_type, duration, nil

	}
}
func distance(steps int, height float64) float64 {
	//longOfstep := height * stepLengthCoefficient
	return height * stepLengthCoefficient * float64(steps) / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	distanceKm := distance(steps, height) / mInKm

	// Переводим продолжительность в часы
	// time.Duration в Go хранится в наносекундах
	hours := duration.Hours()

	// Вычисляем скорость: км / ч
	if hours == 0 {
		return 0 // избегаем деления на ноль
	}

	speed := distanceKm / hours
	return speed
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	minutes := duration.Minutes()

	// Вычисляем скорость: км / ч
	if minutes == 0 {
		return 0, fmt.Errorf("0 minutes") // избегаем деления на ноль
	}
	speed := meanSpeed(steps, height, duration)
	return speed * weight * duration.Minutes() / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	minutes := duration.Minutes()

	// Вычисляем скорость: км / ч
	if minutes == 0 {
		return 0, fmt.Errorf("0 minutes") // избегаем деления на ноль
	}
	speed := meanSpeed(steps, height, duration)
	return speed * weight * duration.Minutes() / minInH * walkingCaloriesCoefficient, nil
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, action, duration, err := parseTraining(data)
	if steps == 0 || duration == 0 || action == "" {
		log.Printf("zero values %w\n", err)
		return "", fmt.Errorf("zero values %w", err) // избегаем деления на ноль
	}
	var cal float64
	switch action {
	case "Бег":
		cal, err = RunningSpentCalories(steps, weight, height, duration)
	case "Хотьба":
		cal, err = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", fmt.Errorf("Unknown type of action %s", action)
	}
	dis := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %v\nДистанция: %v км.\nСкорость: %v км/ч\n Сожгли калорий: %.2f\n", action, duration.Hours(), dis, speed, cal), nil
}
