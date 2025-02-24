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
	lenStep   = 0.65  // средняя длина шага.
	mInKm     = 1000  // количество метров в километре.
	minInH    = 60    // количество минут в часе.
	kmhInMsec = 0.278 // коэффициент для преобразования км/ч в м/с.
	cmInM     = 100   // количество сантиметров в метре.
)

func parseTraining(data string) (int, string, time.Duration, error) {
	// ваш код ниже
	str := strings.Split(data, ",")
	if len(str) != 3 {
		return 0, "", 0, errors.New("error split data")
	}
	count, err := strconv.Atoi(str[0])
	if err != nil {
		return 0, "", 0, err
	}
	duration, err := time.ParseDuration(str[2])
	if err != nil {
		return 0, "", 0, err
	}
	return count, str[1], duration, nil
}

// distance возвращает дистанцию(в километрах), которую преодолел пользователь за время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий (число шагов при ходьбе и беге).
func distance(steps int) float64 {
	// ваш код ниже
	return float64(steps) * lenStep / float64(mInKm)
}

// meanSpeed возвращает значение средней скорости движения во время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий(число шагов при ходьбе и беге).
// duration time.Duration — длительность тренировки.
func meanSpeed(steps int, duration time.Duration) float64 {
	// ваш код ниже
	if duration <= 0 {
		return 0
	}
	lenDistance := distance(steps)
	hours := duration.Hours()
	return lenDistance / hours
}

// ShowTrainingInfo возвращает строку с информацией о тренировке.
//
// Параметры:
//
// data string - строка с данными.
// weight, height float64 — вес и рост пользователя.
func TrainingInfo(data string, weight, height float64) string {
	// ваш код ниже
	count, activity, duration, err := parseTraining(data)
	if err != nil {
	}
	hours := duration.Hours()
	lenDistance := distance(count)
	speed := meanSpeed(count, duration)
	switch {
	case activity == "Бег":
		calories := RunningSpentCalories(count, weight, duration)
		return fmt.Sprintf("Тип тренировки: %s.\nДлительность: %.2f ч. \nДистанция составила %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", activity, hours, lenDistance, speed, calories)
	case activity == "Ходьба":
		calories := WalkingSpentCalories(count, weight, height, duration)
		return fmt.Sprintf("Тип тренировки: %s.\nДлительность: %.2f ч. \nДистанция составила %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", activity, hours, lenDistance, speed, calories)
	}
	return "неизвестный тип тренировки"
}

// Константы для расчета калорий, расходуемых при беге.
const (
	runningCaloriesMeanSpeedMultiplier = 18.0 // множитель средней скорости.
	runningCaloriesMeanSpeedShift      = 20.0 // среднее количество сжигаемых калорий при беге.
)

// RunningSpentCalories возвращает количество потраченных калорий при беге.
//
// Параметры:
//
// steps int - количество шагов.
// weight float64 — вес пользователя.
// duration time.Duration — длительность тренировки.
func RunningSpentCalories(steps int, weight float64, duration time.Duration) float64 {
	// ваш код здесь
	speed := meanSpeed(steps, duration)
	return (runningCaloriesMeanSpeedMultiplier*speed - runningCaloriesMeanSpeedShift) * weight
}

// Константы для расчета калорий, расходуемых при ходьбе.
const (
	walkingCaloriesWeightMultiplier = 0.035 // множитель массы тела.
	walkingSpeedHeightMultiplier    = 0.029 // множитель роста.
)

// WalkingSpentCalories возвращает количество потраченных калорий при ходьбе.
//
// Параметры:
//
// steps int - количество шагов.
// duration time.Duration — длительность тренировки.
// weight float64 — вес пользователя.
// height float64 — рост пользователя.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) float64 {
	// ваш код здесь
	speed := meanSpeed(steps, duration)
	hours := duration.Hours()
	return ((walkingCaloriesWeightMultiplier * weight) + (speed*speed/height)*walkingSpeedHeightMultiplier) * hours * minInH
}
