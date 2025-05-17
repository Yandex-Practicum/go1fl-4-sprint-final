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
	valueData := strings.Split(data, ",")
	if len(valueData) != 3 {
		return 0, "", 0, errors.New("!= 3")
	}

	valueStep, err := strconv.Atoi(valueData[0])
	if err != nil {
		return 0, "", 0, err
	}
	if valueStep <= 0 {
		return 0, "", 0, errors.New("<= 0")
	}

	valueDuration, err := time.ParseDuration(valueData[2])
	if err != nil {
		return 0, "", 0, err
	}
	if valueDuration <= 0 {
		return 0, "", 0, errors.New("<= 0")
	}

	return valueStep, valueData[1], valueDuration, nil
}

func distance(steps int, height float64) float64 {
	length := height * stepLengthCoefficient
	distance := float64(steps) * length / mInKm

	return distance
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	sectionValue := distance(steps, height)
	valueSpeed := sectionValue / duration.Hours()

	return valueSpeed // TODO: реализовать функцию
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	valueStep, typeActivity, valueDuration, err := parseTraining(data)
	if err != nil {
		return "", err
	}

	var calories float64

	switch {
	case typeActivity == "Бег":
		calories, err = RunningSpentCalories(valueStep, weight, height, valueDuration)
		if err != nil {
			return "", err
		}

	case typeActivity == "Ходьба":
		calories, err = WalkingSpentCalories(valueStep, weight, height, valueDuration)
		if err != nil {
			return "", err
		}
	default:
		return "", errors.New("неизвестный тип тренировки")

	}
	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", typeActivity, float64(valueDuration.Hours()), distance(valueStep, height), meanSpeed(valueStep, height, valueDuration), calories), nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if duration <= 0 {
		return 0, errors.New("<= 0")
	}
	if weight <= 0 {
		return 0, errors.New("<= 0")
	}
	if height <= 0 {
		return 0, errors.New("<= 0")
	}
	if steps <= 0 {
		return 0, errors.New("<= 0")
	}

	averageSpeed := meanSpeed(steps, height, duration)
	numberCalories := (weight * averageSpeed * duration.Minutes()) / minInH

	return numberCalories, nil
	// TODO: реализовать функцию
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if duration <= 0 {
		return 0, errors.New("<= 0")
	}
	if weight <= 0 {
		return 0, errors.New("<= 0")
	}
	if height <= 0 {
		return 0, errors.New("<= 0")
	}
	if steps <= 0 {
		return 0, errors.New("<= 0")
	}

	averageSpeed := meanSpeed(steps, height, duration)
	numberCalories := (weight * averageSpeed * duration.Minutes()) / minInH
	adjustedValue := numberCalories * walkingCaloriesCoefficient

	return adjustedValue, nil
	// TODO: реализовать функцию
}
