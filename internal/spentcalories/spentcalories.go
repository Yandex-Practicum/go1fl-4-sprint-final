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
	dataToStrings := strings.Split(data, ",")
	if len(dataToStrings) != 3 {
		return 0, "", 0, errors.New("not enough data to parse a training")
	}

	steps, err := strconv.Atoi(strings.TrimSpace(dataToStrings[0]))
	if err != nil {
		return 0, "", 0, errors.New("failed to convert steps to int: " + err.Error())
	}
	if steps < 1 {
		return 0, "", 0, errors.New("can't execute program if there are less than 1 steps")
	}

	activity := strings.TrimSpace(dataToStrings[1])
	if strings.ToLower(activity) != "ходьба" && strings.ToLower(activity) != "бег" {
		return 0, "", 0, errors.New("неизвестный тип тренировки " + activity)
	}

	duration, err := time.ParseDuration(strings.TrimSpace(dataToStrings[2]))
	if err != nil {
		return 0, "", 0, errors.New("failed to parse duration: " + err.Error())
	}
	if duration <= 0 {
		return 0, "", 0, errors.New("duration must be > 0")
	}
	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {
	return (float64(steps) * height * stepLengthCoefficient) / float64(mInKm)
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	return distance(steps, height) / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		return "", errors.New("failed to parse training data: " + err.Error())
	}
	calSpent, err := RunningSpentCalories(steps, weight, height, duration)
	if err != nil {
		return "", errors.New("failed to calculate calories in TrainingInfo: " + err.Error())
	}
	if strings.ToLower(activity) != "бег" {
		calSpent *= walkingCaloriesCoefficient
	}
	totalDistance := distance(steps, height)
	speed := meanSpeed(steps, height, duration)
	durationInFloat := float64(duration.Hours())
	report := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", activity, durationInFloat, totalDistance, speed, calSpent)
	return report, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("invalid input: all parameters must be positive")
	}
	speed := meanSpeed(steps, height, duration)
	calories := (weight * speed * duration.Minutes()) / float64(minInH)
	return calories, nil
}

/*
# Данная функция - лишняя, как по мне. Проще проверить уже непосредственно в функции TrainingInfo
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	calSpentWalking, err := RunningSpentCalories(steps, weight, height, duration)*0.5
	return calSpentWalking, err
}

*/
