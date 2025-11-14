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
	// TODO: реализовать функцию

	dataSlicSpen := strings.Split(data, ",")
	if len(dataSlicSpen) != 3 {
		return 0, "", 0, errors.New("ошибка в spentcalories: мало данных")
	}

	stepSpent, err := strconv.Atoi(dataSlicSpen[0])
	if err != nil {
		return 0, "", 0, err
	}
	if stepSpent <= 0 {
		return 0, "", 0, errors.New("ошибка spentcalories: шаги <= 0")
	}

	timeSpent, err := time.ParseDuration(dataSlicSpen[2])
	if err != nil {
		return 0, "", 0, err
	}

	activy := dataSlicSpen[1]

	return stepSpent, activy, timeSpent, nil

}
func distance(steps int, height float64) float64 {
	// TODO: реализовать функцию

	lenStep := height * stepLengthCoefficient

	distSpen := (lenStep * float64(steps)) / mInKm

	return distSpen

}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	// TODO: реализовать функцию

	if duration <= 0 {
		return 0
	}

	dist := distance(steps, height)

	tim := duration.Hours()

	return dist / tim

}

func TrainingInfo(data string, weight, height float64) (string, error) {
	// TODO: реализовать функцию

	stepT, activy, timeT, err := parseTraining(data)

	if err != nil {
		log.Println(err)
		return "ошибка Training 1", err
	}

	kkalWal, err := WalkingSpentCalories(stepT, weight, height, timeT)
	if err != nil {
		return "ошибка Training 2", err
	}

	kkalRun, err := RunningSpentCalories(stepT, weight, height, timeT)
	if err != nil {
		return "ошибка Training 3", err
	}

	switch activy {
	case "Ходьба":
		fmt.Sprintf("Тип тренировки: Ходьба\nДлительность: %.2f ч.\nДистанция: %.2f км.\n Скорость: %.2f км/ч\nСожгли калорий: %.2f",
			timeT.Hours(), distance(stepT, height), meanSpeed(stepT, height, timeT), kkalWal)

	case "Бег":
		fmt.Sprintf("Тип тренировки: Бег\nДлительность: %.2f ч.\nДистанция: %.2f км.\n Скорость: %.2f км/ч\nСожгли калорий: %.2f",
			timeT.Hours(), distance(stepT, height), meanSpeed(stepT, height, timeT), kkalRun)
	}

	return "", errors.New("ошибка в выводе Spen")

}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// TODO: реализовать функцию
	switch {
	case steps <= 0:
		return 0, errors.New("ошибка RunningSpentCalories: шаги <= 0")

	case weight <= 0:
		return 0, errors.New("ошибка RunningSpentCalories: вес <= 0")

	case height <= 0:
		return 0, errors.New("ошибка RunningSpentCalories: рост <= 0")

	case duration <= 0:
		return 0, errors.New("ошибка RunningSpentCalories: время <= 0")
	}

	averSpeed := meanSpeed(steps, height, duration)

	timeRunning := duration.Minutes()

	kkal := (weight * averSpeed * timeRunning) / minInH

	return kkal, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// TODO: реализовать функцию
	switch {
	case steps <= 0:
		return 0, errors.New("ошибка WalkingSpentCalories: шаги <= 0")

	case weight <= 0:
		return 0, errors.New("ошибка WalkingSpentCalories: вес <= 0")

	case height <= 0:
		return 0, errors.New("ошибка WalkingSpentCalories: рост <= 0")

	case duration <= 0:
		return 0, errors.New("ошибка WalkingSpentCalories: время <= 0")
	}

	averSpeedW := meanSpeed(steps, height, duration)

	timeWalk := duration.Minutes()

	kkal := (weight * averSpeedW * timeWalk) / minInH

	return kkal * walkingCaloriesCoefficient, nil
}
