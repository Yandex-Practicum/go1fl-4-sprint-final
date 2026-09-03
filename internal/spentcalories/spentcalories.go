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
	sliceString := strings.Split(data, ",")
	if len(sliceString) != 3 {
		return 0, " ", 0, fmt.Errorf("неверный формат строки: ожидается 2 части, найдено %d", len(sliceString))
	}
	step, err := strconv.Atoi(sliceString[0])
	if err != nil {
		return 0, " ", 0, fmt.Errorf("неверный формат строки %s", sliceString[0])
	}
	if step <= 0 {
		return 0, " ", 0, fmt.Errorf("Количество шагов %s ошибка", sliceString[0])
	}
	duration, ok := time.ParseDuration(sliceString[2])
	if ok != nil {
		return 0, "", 0, fmt.Errorf("неверный формат времени: %w", ok)
	}
	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("неверный формат времени: %w", ok)
	}
	return step, sliceString[1], duration, nil
}
func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	totalMeters := float64(steps) * stepLength
	return totalMeters / mInKm
}
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	if height <= 0 {
		return 0
	}
	if steps <= 0 {
		return 0
	}
	distKm := distance(steps, height)
	hours := duration.Hours()
	return distKm / hours
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	var calories float64
	var err error
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}
	speed := meanSpeed(steps, height, duration)
	switch activity {
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", fmt.Errorf("неизвестный тип тренировки: %s", activity)
	}
	if err != nil {
		return "", err
	}
	hours := float64(duration) / float64(time.Hour)
	report := fmt.Sprintf(
		"Тип тренировки: %s\n"+
			"Длительность: %.2f ч.\n"+
			"Дистанция: %.2f км.\n"+
			"Скорость: %.2f км/ч\n"+
			"Сожгли калорий: %.2f\n",
		activity,                // %s
		hours,                   // %.2f (переведенные часы)
		distance(steps, height), // %.2f (вызываем функцию дистанции)
		speed,                   // %.2f
		calories,                // %.2f
	)
	return report, nil

}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("неправильные данные шагов")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("неправильные данные времени")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("неправильные данные веса")
	}
	if height <= 0 {
		return 0, fmt.Errorf("неправильные данные роста")
	}
	speed := float64(meanSpeed(steps, height, duration))
	durationInMinutes := float64(duration) / float64(time.Minute)
	cal := (weight * speed * durationInMinutes) / minInH
	return cal, nil

}
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("неправильные данные шагов")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("неправильные данные времени")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("неправильные данные веса")
	}
	if height <= 0 {
		return 0, fmt.Errorf("неправильные данные роста")
	}
	speed := float64(meanSpeed(steps, height, duration))
	durationInMin := float64(duration) / float64(time.Minute)
	cal := ((durationInMin * weight * speed) / minInH) * walkingCaloriesCoefficient
	return cal, nil
}
