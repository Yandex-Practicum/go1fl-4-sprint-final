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

// parseTraining разбивает на отдельные позиции строку с данными о количестве шагов, виде и продолжительности активности.
func parseTraining(data string) (int, string, time.Duration, error) {
	activity := strings.Split(data, ",")

	if len(activity) != 3 {
		return 0, "", 0, errors.New("Входящая строка не соответствует требованию о наличии 3 позиций")
	}

	stepsQuantity, err := strconv.Atoi(activity[0])
	if err != nil {
		return 0, "", 0, err
	}

	if stepsQuantity <= 0 {
		return 0, "", 0, errors.New("ошибка, количество шагов менее либо равно нулю!")
	}

	activityDuration, err := time.ParseDuration(activity[2])
	if err != nil {
		return 0, "", 0, err
	}

	if activityDuration <= 0 {
		return 0, "", 0, errors.New("ошибка, длительность активности менее либо равна нулю!")
	}

	return stepsQuantity, activity[1], activityDuration, nil
}

func distance(steps int, height float64) float64 {
	oneStepLength := height * stepLengthCoefficient
	result := (float64(steps) * oneStepLength) / mInKm
	return result
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	avgSpeed := distance(steps, height) / duration.Hours()
	return avgSpeed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	stepsQuantity, activityType, activityDuration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
	}
	switch activityType {
	case "Бег":
		distance := distance(stepsQuantity, height)
		speedAvg := meanSpeed(stepsQuantity, height, activityDuration)
		caloriesRunning, err := RunningSpentCalories(stepsQuantity, weight, height, activityDuration)
		if err != nil {
			return "", err
		}
		res := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", activityType, activityDuration.Hours(), distance, speedAvg, caloriesRunning)
		return res, err
	case "Ходьба":
		distance := distance(stepsQuantity, height)
		speedAvg := meanSpeed(stepsQuantity, height, activityDuration)
		caloriesWalking, err := WalkingSpentCalories(stepsQuantity, weight, height, activityDuration)
		if err != nil {
			return "", err
		}
		res := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", activityType, activityDuration.Hours(), distance, speedAvg, caloriesWalking)
		return res, err
	default:
		return "", errors.New("неизвестный тип тренировки")
	}
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if weight <= 0 {
		return 0, errors.New("ошибка, вес указан как отрицательный либо нулевой!")
	}
	if height <= 0 {
		return 0, errors.New("ошибка, рост указан как отрицательный либо нулевой!")
	}
	if duration <= 0 {
		return 0, errors.New("ошибка, длительность менее либо равна нулю!")
	}
	if steps <= 0 {
		return 0, errors.New("ошибка, количество шагов менее либо равно нулю")
	}

	avgSpeed := meanSpeed(steps, height, duration)
	spentCaloriesRunning := (weight * avgSpeed * duration.Minutes()) / minInH

	return spentCaloriesRunning, nil

}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if weight <= 0 {
		return 0, errors.New("ошибка, вес указан как отрицательный либо нулевой!")
	}
	if height <= 0 {
		return 0, errors.New("ошибка, рост указан как отрицательный либо нулевой!")
	}
	if duration <= 0 {
		return 0, errors.New("ошибка, длительность менее либо равна нулю!")
	}
	if steps <= 0 {
		return 0, errors.New("ошибка, количество шагов менее либо равно нулю")
	}

	avgSpeed := meanSpeed(steps, height, duration)
	spentCaloriesWalking := (weight * avgSpeed * duration.Minutes()) / minInH * walkingCaloriesCoefficient
	return spentCaloriesWalking, nil
}
