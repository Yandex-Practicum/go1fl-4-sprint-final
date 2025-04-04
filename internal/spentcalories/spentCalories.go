package spentcalories

import (
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

//Функция принимает строку с данными формата "3456,Ходьба,3h00m", которая содержит количество шагов, вид активности и продолжительность активности. Функция возвращает четыре значения:
//int — количество шагов.
//string — вид активности.
//time.Duration — продолжительность активности.
//error — ошибку, если что-то пошло не так.

func parseTraining(data string) (int, string, time.Duration, error) {
	var period time.Duration
	//	strData := strings.Split(data, ",")
	strData := strings.Split(data, ",")
	if len(strData) != 3 {
		return 0, "0", period, fmt.Errorf("[parseTraining] неверный формат строки: %s", data)
	}
	step, err := strconv.Atoi(strData[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("[parseTraining] неверный формат количества шагов %s: %w", strData[0], err)
	}
	if len(strData) != 3 {
		return 0, "0", period, fmt.Errorf("[parseTraining] неверный формат строки данных")
	}
	if step < 1 {
		return 0, "0", period, fmt.Errorf("[parseTraining] неверное количество шагов: %d", step)
	}

	activity := strData[1]
	period, err = time.ParseDuration(strData[2])
	if err != nil {
		return 0, "0", period, err
	}

	return step, activity, period, nil
}

// distance возвращает дистанцию(в километрах), которую преодолел пользователь за время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий (число шагов при ходьбе и беге).
func distance(steps int) float64 {
	return float64(steps) * lenStep / float64(mInKm)
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
	dist := distance(steps)
	return dist / duration.Hours()
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
	// ваш код здесь
	v := meanSpeed(steps, duration)
	calories := ((runningCaloriesMeanSpeedMultiplier * v) - runningCaloriesMeanSpeedShift) * weight
	return calories
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
	v := meanSpeed(steps, duration)
	calories := ((walkingCaloriesWeightMultiplier * weight) + (v*v/height)*walkingSpeedHeightMultiplier) * duration.Hours() * minInH
	return calories
}

// ShowTrainingInfo возвращает строку с информацией о тренировке.
//
// Параметры:
//
// data string - строка с данными.
// weight, height float64 — вес и рост пользователя.
func TrainingInfo(data string, weight, height float64) string {
	var strInfo string
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		return ""
	}

	switch activity {
	case "Ходьба":
		dist := distance(steps)
		v := meanSpeed(steps, duration)
		calories := WalkingSpentCalories(steps, weight, height, duration)
		strInfo = fmt.Sprintf("Тип тренировки: %s\nДлительность: %0.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f", activity, duration.Hours(), dist, v, calories)
		// воспринимать сотые доли часа неудобно. Вариант вывода в часах и минутах:
		// strInfo = fmt.Sprintf("Тип тренировки: %s\nДлительность: %.0fч %d мин.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f", activity, duration.Hours(), int(duration.Minutes())%60, dist, v, calories)

	case "Бег":
		dist := distance(steps)
		v := meanSpeed(steps, duration)
		calories := RunningSpentCalories(steps, weight, duration)
		strInfo = fmt.Sprintf("Тип тренировки: %s\nДлительность: %0.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f", activity, duration.Hours(), dist, v, calories)

	default:
		fmt.Println("[TrainingInfo(string, float64, float63)] неизвестный тип тренировки: ", activity)
	}
	return strInfo
}
