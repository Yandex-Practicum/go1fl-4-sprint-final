package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

var (
	ErrInvalidDataFormat = errors.New("некорректный формат входных данных")
	ErrInvalidStepCount  = errors.New("некорректное количество шагов")
	ErrInvalidActivity   = errors.New("некорректный тип тренировки")
	ErrInvalidDuration   = errors.New("некорректная длительность")

	ErrUknownActivity = errors.New("неизвестный тип тренировки")

	ErrInvalidArgs = errors.New("некорректные входные параметры")
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
	splittedData := strings.Split(data, ",")
	splittedDataLen := len(splittedData)

	if splittedDataLen != 3 {
		return 0, "", 0, ErrInvalidDataFormat
	}

	stepData := splittedData[0]
	activityData := splittedData[1]
	durationData := splittedData[2]

	stepCount, err := strconv.Atoi(stepData)
	if err != nil {
		return 0, "", 0, err
	}

	if stepCount <= 0 {
		return 0, "", 0, ErrInvalidStepCount
	}

	if activityData == "" {
		return 0, "", 0, ErrInvalidActivity
	}

	duration, err := time.ParseDuration(durationData)
	if err != nil {
		return 0, "", 0, err
	}

	if duration <= 0 {
		return 0, "", 0, ErrInvalidDuration
	}

	return stepCount, activityData, duration, nil
}

func distance(steps int, height float64) float64 {
	if steps <= 0 || height <= 0 {
		return 0
	}

	return (height * stepLengthCoefficient * float64(steps)) / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	distance := distance(steps, height)
	hourCount := duration.Hours()

	return distance / hourCount
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	stepCount, activity, duration, err := parseTraining(data)

	if err != nil {
		log.Println(err)
		return "", err
	}

	distance := distance(stepCount, height)
	speed := meanSpeed(stepCount, height, duration)

	calories := 0.0

	switch activity {
	case "Ходьба":
		calories, err = WalkingSpentCalories(stepCount, weight, height, duration)
	case "Бег":
		calories, err = RunningSpentCalories(stepCount, weight, height, duration)
	default:
		err = ErrUknownActivity
	}

	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
			`Тип тренировки: %s
Длительность: %.2f ч.
Дистанция: %.2f км.
Скорость: %.2f км/ч
Сожгли калорий: %.2f
`,
			activity,
			duration.Hours(),
			distance,
			speed,
			calories,
		),
		nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, ErrInvalidArgs
	}

	speed := meanSpeed(steps, height, duration)

	return (weight * speed * duration.Minutes()) / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	calories, err := RunningSpentCalories(steps, weight, height, duration)
	if err != nil {
		return 0, err
	}

	return calories * walkingCaloriesCoefficient, nil
}
