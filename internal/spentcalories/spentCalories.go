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

func parseTraining(data string) (int, string, time.Duration, error) {
	// ваш код ниже
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("неверный формат данных: %s", data)
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil || steps <= 0 {
		return 0, "", 0, fmt.Errorf("ошибка при парсинге количества шагов: %v", err)
	}
	activity := strings.TrimSpace(parts[1])
	if activity == "" {
		return 0, "", 0, fmt.Errorf("вид активности не может быть пустым")
	}
	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка при парсинге продолжительности: %v", err)
	}

	return steps, activity, duration, nil
}

// distance возвращает дистанцию(в километрах), которую преодолел пользователь за время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий (число шагов при ходьбе и беге).
func distance(steps int) float64 {
	// ваш код ниже
	return float64(steps) * lenStep / mInKm
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
	dist := distance(steps)
	hours := duration.Hours()
	return dist / hours
}

// ShowTrainingInfo возвращает строку с информацией о тренировке.
//
// Параметры:
//
// data string - строка с данными.
// weight, height float64 — вес и рост пользователя.
func TrainingInfo(data string, weight, height float64) string {
	// ваш код ниже
	steps, activityType, duration, err := parseTraining(data)
	if err != nil {
		return fmt.Sprintf("Ошибка парсинга данных: %v", err)
	}
	var (
		distance float64
		speed    float64
		calories float64
		result   string
	)
	switch activityType {
	case "Ходьба":
		distance = distance(steps)
		speed = meanSpeed(steps, duration)
		calories = WalkingSpentCalories(steps, weight, height, duration)
		result = fmt.Sprintf(
			"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км\nСожгли калорий: %.2f",
			activityType, duration.Hours(), distance, speed, calories)
	case "Бег":
		distance = distance(steps)
		speed = meanSpeed(steps, duration)
		calories = RunningSpentCalories(steps, weight, duration)
		result = fmt.Sprintf(
			"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f",
			activityType, duration.Hours(), distance, speed, calories)
	default:
		result = "Неизвестный тип тренировки"
	}
	return result
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
	meanSpeed := meanSpeed(steps, duration)
	calories := ((runningCaloriesMeanSpeedMultiplier * meanSpeed) - runningCaloriesMeanSpeedShift) * weight
	if calories < 0 {
		return 0
	}
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
	// ваш код здесь
	meanSpeed := meanSpeed(steps, duration)
	hours := duration.Hours()
	calories := ((walkingCaloriesWeightMultiplier * weight) + ((meanSpeed * meanSpeed / height) * walkingSpeedHeightMultiplier)) * hours * minInH
	if calories < 0 {
		return 0
	}
	return calories
}
