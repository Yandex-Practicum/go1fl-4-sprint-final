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

	parseData := strings.Split(data, ",")
	if len(parseData) == 3 {
		step, err := strconv.Atoi(parseData[0])
		if err != nil {
			return 0, "", time.Duration(0), err
		}
		if step <= 0 {
			return 0, "", time.Duration(0), errors.New("The number of steps is less than zero")
		}

		walkingTime, err := time.ParseDuration(parseData[2])
		if err != nil {
			return 0, "", time.Duration(0), err
		}
		if walkingTime <= time.Duration(0) {
			return 0, "", time.Duration(0), errors.New("Duration <= 0")
		}

		return step, parseData[1], walkingTime, nil
	}

	return 0, "", time.Duration(0), errors.New("Incorrect input data")
}

func distance(steps int, height float64) float64 {
	return (float64(steps) * (height * stepLengthCoefficient)) / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= time.Duration(0) {
		return 0
	}
	return distance(steps, height) / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	step, training, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
	}

	switch training {
	case "Ходьба":
		dist := distance(step, height)
		arrSpid := meanSpeed(step, height, duration)
		calories, err := WalkingSpentCalories(step, weight, height, duration)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", training, duration.Hours(), dist, arrSpid, calories), nil

	case "Бег":
		dist := distance(step, height)
		arrSpid := meanSpeed(step, height, duration)
		calories, err := RunningSpentCalories(step, weight, height, duration)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", training, duration.Hours(), dist, arrSpid, calories), nil
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= time.Duration(0) {
		return 0, errors.New("Input parameters are incorrect")
	}
	durationInMinutes := duration.Minutes()
	return (weight * meanSpeed(steps, height, duration) * durationInMinutes) / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= time.Duration(0) {
		return 0, errors.New("Input parameters are incorrect")
	}
	return ((weight * meanSpeed(steps, height, duration) * duration.Minutes()) / minInH) * walkingCaloriesCoefficient, nil

}
