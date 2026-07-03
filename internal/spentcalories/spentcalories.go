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
	dataSlice := strings.Split(data, ",")
	if len(dataSlice) != 3 {
		return 0, "", time.Duration(0), errors.New("Неверный формат данных")
	}

	steps, err := strconv.Atoi(dataSlice[0])
	if err != nil {
		return 0, "", time.Duration(0), errors.New("Нeверный формат количества шагов")
	}
	if steps <= 0 {
		return 0, "", time.Duration(0), errors.New("Количество шагов должно быть больше 0")
	}

	duration, err := time.ParseDuration(dataSlice[2])
	if err != nil {
		return 0, "", time.Duration(0), errors.New("Неверный формат продолжительности тренировки")
	}
	if duration <= 0 {
		return 0, "", time.Duration(0), errors.New("Продолжительность тренировки должна быть больше 0")
	}

	activityType := dataSlice[1]

	return steps, activityType, duration, nil
}

func distance(steps int, height float64) float64 {
	distance := (height * stepLengthCoefficient * float64(steps)) / mInKm
	return distance
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	distance := distance(steps, height)
	speed := distance / duration.Hours()
	return speed
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
	if steps <= 0 {
		return 0, errors.New("Количество шагов должно быть больше 0")
	}
	if duration <= 0 {
		return 0, errors.New("Продолжительность тренировки должна быть больше 0")
	}
	if weight <= 0 || height <= 0 {
		return 0, errors.New("Вес и рост должны быть больше 0")
	}

	meanSpeed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()

	calories := (weight * meanSpeed * durationInMinutes) / minInH

	return calories, nil

}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("Количество шагов должно быть больше 0")
	}
	if duration <= 0 {
		return 0, errors.New("Продолжительность тренировки должна быть больше 0")
	}
	if weight <= 0 || height <= 0 {
		return 0, errors.New("Вес и рост должны быть больше 0")
	}

	meanSpeed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()

	calories := (weight * meanSpeed * durationInMinutes) / minInH
	calories *= walkingCaloriesCoefficient

	return calories, nil

}
