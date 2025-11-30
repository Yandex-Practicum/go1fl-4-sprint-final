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
	ErrorMissingArguments  = errors.New("spentcalories: missing arguments")
	ErrorStringToInt       = errors.New("spentcalories: cannot convert string to int")
	ErrorZeroSteps         = errors.New("spentcalories: zero steps")
	ErrorParseToDuration   = errors.New("spentcalories: cannot parse duration")
	ErrorIncorrectArgument = errors.New("spentcalories: incorrect argument")
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
		return 0, "", 0, fmt.Errorf("%w: not enough data arguments: something is missing parseTraining", ErrorMissingArguments)
	}
	steps, err := strconv.Atoi(dataSlice[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("%w: something went wrong while converting string to int; parseTrraining", ErrorStringToInt)
	}
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("%w: steps are equal or lower than zero; parsetTraining", ErrorIncorrectArgument)
	}
	duration, err := time.ParseDuration(dataSlice[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("%w: something went wrong while parsing string to duration; parsetraining", ErrorParseToDuration)
	}
	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("%w: duration is equal or lower than zero; parseTraining", ErrorIncorrectArgument)
	}
	if dataSlice[1] != "Бег" && dataSlice[1] != "Ходьба" {
		return 0, "", 0, fmt.Errorf("неизвестный тип тренировки")
	}

	return steps, dataSlice[1], duration, nil
}

func distance(steps int, height float64) float64 {
	if steps <= 0 {
		return 0
	}

	if height <= 0 {
		return 0
	}
	stepLength := stepLengthCoefficient * height
	distance := stepLength * float64(steps)
	distance = distance / mInKm

	return distance
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if height <= 0 {
		return 0
	}
	if duration <= 0 {
		return 0
	}
	distance := distance(steps, height)
	if distance <= 0 {
		return 0
	}

	avgSpeed := distance / duration.Hours()

	return avgSpeed

}

func TrainingInfo(data string, weight, height float64) (string, error) {
	var str string
	steps, activityType, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}
	if weight <= 0 {
		return "", fmt.Errorf("%w: weight <= 0; TrainingInfo", ErrorIncorrectArgument)
	}
	if height <= 0 {
		return "", fmt.Errorf("%w: height <= 0; TrainingInfo", ErrorIncorrectArgument)
	}

	switch activityType {
	case "Ходьба":
		calories, err := WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}

		avgSpeed := meanSpeed(steps, height, duration)
		distance := distance(steps, height)

		str = fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", activityType, duration.Hours(), distance, avgSpeed, calories)
	case "Бег":
		calories, err := RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}

		avgSpeed := meanSpeed(steps, height, duration)
		distance := distance(steps, height)

		str = fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", activityType, duration.Hours(), distance, avgSpeed, calories)
	}
	return str, err
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("%w: steps equal to zero", ErrorMissingArguments)
	}
	if duration <= 0 {
		return 0, fmt.Errorf("%w: duration equal to zero", ErrorMissingArguments)
	}
	if weight <= 0 {
		return 0, fmt.Errorf("%w: weight <= 0; RunningSpentCalories", ErrorIncorrectArgument)
	}
	if height <= 0 {
		return 0, fmt.Errorf("%w: height <= 0; RunningSpentCalories", ErrorIncorrectArgument)
	}
	avgSpeed := meanSpeed(steps, height, duration)
	if avgSpeed <= 0 {
		return 0, fmt.Errorf("%w; meanSpeed equals to 0; RunningSpentCalories", ErrorIncorrectArgument)
	}

	calories := (weight * avgSpeed * duration.Minutes()) / minInH

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("%w: steps equal to zero", ErrorMissingArguments)
	}
	if duration <= 0 {
		return 0, fmt.Errorf("%w: duration equal to zero", ErrorMissingArguments)
	}
	avgSpeed := meanSpeed(steps, height, duration)
	if avgSpeed <= 0 {
		return 0, fmt.Errorf("%w; meanSpeed equals to 0; WalkingSpentCalories", ErrorIncorrectArgument)
	}
	if weight <= 0 {
		return 0, fmt.Errorf("%w: weight <= 0; WalkingSpentCalories", ErrorIncorrectArgument)
	}
	if height <= 0 {
		return 0, fmt.Errorf("%w: height <= 0; WalkingSpentCalories", ErrorIncorrectArgument)
	}

	calories := (weight * avgSpeed * duration.Minutes() * walkingCaloriesCoefficient) / minInH

	return calories, nil
}
