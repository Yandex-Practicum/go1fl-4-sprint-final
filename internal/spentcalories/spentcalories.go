package spentcalories

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// Константы
const (
	stepLengthCoefficient = 0.45 // коэффициент длины шага
	mInKm                = 1000   // метров в километре
	minInH               = 60     // минут в часе
	walkingCaloriesCoefficient = 0.8 // корректирующий коэффициент для ходьбы
)

// parseTraining парсит строку с данными о тренировке.
func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, errors.New("некорректный формат данных: ожидается три элемента")
	}

	stepsStr, activity, durationStr := parts[0], parts[1], parts[2]

	// Преобразуем количество шагов в int
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, "", 0, errors.New("ошибка преобразования количества шагов: " + err.Error())
	}
	if steps <= 0 {
		return 0, "", 0, errors.New("количество шагов должно быть больше нуля")
	}

	// Преобразуем продолжительность в time.Duration
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, "", 0, errors.New("ошибка парсинга продолжительности: " + err.Error())
	}

	return steps, activity, duration, nil
}

// distance вычисляет дистанцию в километрах.
func distance(steps int, height float64) float64 {
	stepLen := height * stepLengthCoefficient
	distanceM := float64(steps) * stepLen
	return distanceM / mInKm
}

// meanSpeed вычисляет среднюю скорость в км/ч.
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	distKm := distance(steps, height)
	durationH := duration.Hours()

	return distKm / durationH
}

// RunningSpentCalories вычисляет калории, потраченные при беге.
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("некорректные входные параметры")
	}

	meanSpeedVal := meanSpeed(steps, height, duration)
	durationM := duration.Minutes()

	calories := (weight * meanSpeedVal * durationM) / minInH
	return calories, nil
}

// WalkingSpentCalories вычисляет калории, потраченные при ходьбе.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("некорректные входные параметры")
	}

	meanSpeedVal := meanSpeed(steps, height, duration)
	durationM := duration.Minutes()

	calories := (weight * meanSpeedVal * durationM) / minInH
	calories *= walkingCaloriesCoefficient

	return calories, nil
}

// TrainingInfo формирует строку с информацией о тренировке.
func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	distKm := distance(steps, height)
	meanSpeedVal := meanSpeed(steps, height, duration)

	var calories float64
	switch activity {
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	if err != nil {
		fmt.Println("Ошибка вычисления калорий: " + err.Error())
		return "", err
	}

	return fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f",
		activity, duration.Hours(), distKm, meanSpeedVal, calories,
	), nil
}
