package spentcalories

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep = 0.65 // средняя длина шага.
	mInKm   = 1000 // количество метров в километре.
	minInH  = 60   // количество минут в часе.
)

func parseTraining(data string) (int, string, time.Duration, error) {
	slicesOfActiv := strings.Split(data, ",")
	if len(slicesOfActiv) != 3 {
		return 0, "", 0, fmt.Errorf("Not enough elements in entry data")
	}
	stepsCount, err := strconv.Atoi(slicesOfActiv[0])
	if err != nil {
		return 0, "", 0, err
	}
	duratWalk, err2 := time.ParseDuration(slicesOfActiv[2])
	if err2 != nil {
		return 0, "", 0, err2
	}

	return stepsCount, slicesOfActiv[1], duratWalk, nil
}

// distance возвращает дистанцию(в километрах), которую преодолел пользователь за время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий (число шагов при ходьбе и беге).
func distance(steps int) float64 {
	return float64(steps) * lenStep / mInKm
}

// meanSpeed возвращает значение средней скорости движения во время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий(число шагов при ходьбе и беге).
// duration time.Duration — длительность тренировки.
func meanSpeed(steps int, duration time.Duration) float64 {
	if duration <= 0 {
		return 0.0
	}
	return distance(steps) / duration.Hours()
}

// ShowTrainingInfo возвращает строку с информацией о тренировке.
//
// Параметры:
//
// data string - строка с данными.
// weight, height float64 — вес и рост пользователя.
func TrainingInfo(data string, weight, height float64) string {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		fmt.Println(err)
		return " "
	}
	var distanceOfActivity, speed, calories float64
	switch activity {
	case "Бег":
		distanceOfActivity = distance(steps)
		speed = meanSpeed(steps, duration)
		calories = RunningSpentCalories(steps, weight, duration)
	case "Ходьба":
		distanceOfActivity = distance(steps)
		speed = meanSpeed(steps, duration)
		calories = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "invalid type of training"
	}
	return fmt.Sprintf("Тип тренировки: %s \nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч.\nСожгли калорий: %.2f", activity, duration.Hours(), distanceOfActivity, speed, calories)
}

// Константы для расчета калорий, расходуемых при беге.
const (
	runningCaloriesMeanSpeedMultiplier = 18.0 // множитель средней скорости.
	runningCaloriesMeanSpeedShift      = 20.0 // среднее количество сжигаемых калорий при беге.
)

// RunningSpentCalories возвращает количество потраченных колорий при беге.
//
// Параметры:
//
// steps int - количество шагов.
// weight float64 — вес пользователя.
// duration time.Duration — длительность тренировки.
func RunningSpentCalories(steps int, weight float64, duration time.Duration) float64 {
	meanSpeed := meanSpeed(steps, duration)

	return ((runningCaloriesMeanSpeedMultiplier * meanSpeed) - runningCaloriesMeanSpeedShift) * weight
}

/*  */
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
	return ((walkingCaloriesWeightMultiplier * weight) + (meanSpeed(steps, duration)*meanSpeed(steps, duration)/height)*walkingSpeedHeightMultiplier) * duration.Hours() * minInH
}
