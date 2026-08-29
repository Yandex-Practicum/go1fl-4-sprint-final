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

	info := strings.Split(data, ",")
	if len(info) != 3 {
		return 0, "0", 0, fmt.Errorf("ошибка преобразования данных")
	}
	steps, err := strconv.Atoi(info[0])
	if err != nil {
		return 0, "0", 0, fmt.Errorf("ошибка преобразования шагов: %v", err)
	}
	if steps <= 0 {
		return 0, "0", 0, fmt.Errorf("количество шагов должно быть положительным")
	}
	trType := info[1]
	duration, err := time.ParseDuration(info[2])
	if err != nil {
		return 0, "0", 0, fmt.Errorf("ошибка преобразования времени: %v", err)
	}
	if duration <= 0 {
		return 0, "0", 0, fmt.Errorf("длительность тренировки должна быть положительной")
	}
	return steps, trType, duration, nil
}

func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	dist := (float64(steps) * stepLength) / float64(mInKm)
	return dist
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	speed := distance(steps, height) / duration.Hours()
	return speed
}

// Инфоромация о проведенной тренировке.
func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, trType, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
	}
	runCal, err := RunningSpentCalories(steps, weight, height, duration)
	walkCal, err := WalkingSpentCalories(steps, weight, height, duration)
	var answer string
	switch trType {
	case "Ходьба":
		answer = fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", trType, duration.Hours(), distance(steps, height), meanSpeed(steps, height, duration), walkCal)
		return answer, nil
	case "Бег":
		answer = fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", trType, duration.Hours(), distance(steps, height), meanSpeed(steps, height, duration), runCal)
	default:
		return answer, fmt.Errorf("неизвестный тип тренировки")
	}
	return answer, nil
}

// Количество израсходованных калорий при беге.
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("данные некорректны")
	}
	calories := (weight * meanSpeed(steps, height, duration) * duration.Minutes()) / float64(minInH)
	return calories, nil
}

// Количество израсходованных калорий при ходьбе.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("данные некорректны")
	}
	calories := (duration.Minutes() * weight * meanSpeed(steps, height, duration)) / float64(minInH)
	return calories * walkingCaloriesCoefficient, nil
}
