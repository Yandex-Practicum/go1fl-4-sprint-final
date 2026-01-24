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
	stringTemplate             = `Тип тренировки: %s
Длительность: %.2f ч.
Дистанция: %.2f км.
Скорость: %.2f км/ч
Сожгли калорий: %.2f
`
)

func parseTraining(data string) (int, string, time.Duration, error) {
	slice := strings.Split(data, ",")
	if len(slice) != 3 {
		return 0, "", 0, errors.New("Длина элементов не 3")
	}
	steps, err := strconv.Atoi(slice[0])
	if err != nil {
		return 0, "", 0, err
	}
	if steps <= 0 {
		return 0, "", 0, errors.New("Колличество шагов должно быть больше 0")
	}
	duration, err := time.ParseDuration(slice[2])
	if err != nil {
		return 0, "", 0, err
	}
	if duration <= 0 {
		return 0, "", 0, errors.New("Продолжительность тренировки должна быть больше 0")
	}
	return steps, slice[1], duration, nil
}

func distance(steps int, height float64) float64 {
	return float64(steps) * (height * stepLengthCoefficient) / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	d := distance(steps, height)
	speed := d / duration.Hours()
	return speed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, trainType, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
	}
	d := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	switch trainType {
	case "Бег":
		c, err := RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf(stringTemplate, trainType, duration.Hours(), d, speed, c), nil
	case "Ходьба":
		c, err := WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf(stringTemplate, trainType, duration.Hours(), d, speed, c), nil
	default:
		return "", errors.New("неизвестный тип тренировки")
	}
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("Все параметры должны быть больше 0")
	}
	speed := meanSpeed(steps, height, duration)
	return duration.Minutes() * weight * speed / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("Все параметры должны быть больше 0")
	}
	speed := meanSpeed(steps, height, duration)
	return duration.Minutes() * weight * speed / minInH * walkingCaloriesCoefficient, nil
}
