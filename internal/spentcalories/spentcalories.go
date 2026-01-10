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
	refactoredData := strings.Split(data, ",")

	if len(refactoredData) != 3 {
		return 0, "", 0, errors.New("Ошибка: некорректные данные")
	}

	steps, err := strconv.Atoi(refactoredData[0])
	if err != nil && steps <= 0 {
		fmt.Println("Ошибка:", err)
		return 0, "", 0, err
	}

	duration, err := time.ParseDuration(refactoredData[2])
	if err != nil && duration <= 0 {
		fmt.Println("Ошибка:", err)
		return 0, "", 0, err
	}
	training := refactoredData[1]
	return steps, training, duration, nil

}

func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	return (stepLength * float64(steps)) / mInKm

}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	meanSpeed := distance(steps, height) / duration.Hours()
	return meanSpeed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, training, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	walking, err := WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println(err)
		return "", err
	}

	running, err := RunningSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println(err)
		return "", err
	}

	switch strings.ToLower(training) {
	case "бег":
		return fmt.Sprintf(
			"Тип тренировки: %s\n"+
				"Длительность: %.2f ч.\n"+
				"Дистанция: %.2f км.\n"+
				"Скорость: %.2f км/ч\n"+
				"Сожгли калорий: %.2f",
			training,
			duration.Hours(),
			distance(steps, height),
			meanSpeed(steps, height, duration),
			running,
		), nil

	case "ходьба":
		return fmt.Sprintf(
			"Тип тренировки: %s\n"+
				"Длительность: %.2f ч.\n"+
				"Дистанция: %.2f км.\n"+
				"Скорость: %.2f км/ч\n"+
				"Сожгли калорий: %.2f",
			training,
			duration.Hours(),
			distance(steps, height),
			meanSpeed(steps, height, duration),
			walking,
		), nil
	default:
		return "", errors.New("Неизвестный тип тренировки")

	}

}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 && weight <= 0 && height <= 0 && duration <= 0 {
		return 0, errors.New("Некорректные данные")

	}
	runningSpentCalories := weight * meanSpeed(steps, height, duration) * duration.Minutes() / minInH
	return runningSpentCalories, nil

}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 && weight <= 0 && height <= 0 && duration <= 0 {
		return 0, errors.New("Некорректные данные")
	}
	TotalCalories := (weight * meanSpeed(steps, height, duration) * duration.Minutes() / minInH) * walkingCaloriesCoefficient
	return TotalCalories, nil
}
