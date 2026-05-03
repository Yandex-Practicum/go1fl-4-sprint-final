package spentcalories

import (
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
	// Разделить строку на слайс по запятой
	parts := strings.Split(data, ",")

	// Проверка, что длина слайса равна 3
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("Неверный формат строки: ожидается 'шаги,активность,продолжительность', получено: %s", data)
	}

	// Получаем шаги
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("Неверный формат количества шагов: %v", err)
	}

	// Проверка, что количество шагов больше 0
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("Количество шагов должно быть больше 0, получено: %d", steps)
	}

	// Получить вид активности
	activity := parts[1]

	// Получаем время тренировки
	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("Неверный формат продолжительности: %v", err)
	}

	// Проверка, что продолжительность положительная
	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("Продолжительность должна быть больше 0, получено: %v", duration)
	}

	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {
	// Рассчитываем длину шага
	stepLength := height * stepLengthCoefficient

	// Рассчитываем дистанцию в метрах
	distanceMeters := float64(steps) * stepLength

	// Переводим дистанцию в километры
	distanceKm := distanceMeters / mInKm

	return distanceKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	distKm := distance(steps, height)

	durationHours := duration.Hours()

	speed := distKm / durationHours

	return speed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	// Получаем значения из строки данных
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Printf("Ошибка при парсинге данных тренировки: %v\n", err)
		return "", err
	}

	var distanceKm float64
	var speed float64
	var calories float64
	var calcErr error

	switch activity {
	case "Ходьба":
		distanceKm = distance(steps, height)

		speed = meanSpeed(steps, height, duration)

		calories, calcErr = WalkingSpentCalories(steps, weight, height, duration)
		if calcErr != nil {
			log.Printf("Ошибка при расчете калорий для ходьбы: %v\n", calcErr)
			return "", calcErr
		}

	case "Бег":
		distanceKm = distance(steps, height)

		speed = meanSpeed(steps, height, duration)

		calories, calcErr = RunningSpentCalories(steps, weight, height, duration)
		if calcErr != nil {
			log.Printf("Ошибка при расчете калорий для бега: %v\n", calcErr)
			return "", calcErr
		}

		// Неизвестный тип тренировки
	default:
		log.Printf("неизвестный тип тренировки: %s\n", activity)
		return "", fmt.Errorf("неизвестный тип тренировки: %s", activity)
	}

	durationHours := duration.Hours()

	// Сформировать строку с информацией о тренировке
	result := fmt.Sprintf(
		"Тип тренировки: %s\n"+
			"Длительность: %.2f ч.\n"+
			"Дистанция: %.2f км.\n"+
			"Скорость: %.2f км/ч\n"+
			"Сожгли калорий: %.2f\n",
		activity, durationHours, distanceKm, speed, calories)

	return result, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("Количество шагов должно быть больше 0, получено: %d", steps)
	}

	if weight <= 0 {
		return 0, fmt.Errorf("Вес должен быть больше 0, получено: %.2f кг", weight)
	}

	if height <= 0 {
		return 0, fmt.Errorf("Рост должен быть больше 0, получено: %.2f м", height)
	}

	if duration <= 0 {
		return 0, fmt.Errorf("Продолжительность должна быть больше 0, получено: %v", duration)
	}

	// Средняя скорость
	avgSpeed := meanSpeed(steps, height, duration)

	// Перевести продолжительность в минуты
	durationMinutes := duration.Minutes()

	// Рассчитать количество калорий
	calories := (weight * avgSpeed * durationMinutes) / minInH

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверка входных параметров на корректность
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть больше 0, получено: %d", steps)
	}

	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть больше 0, получено: %.2f кг", weight)
	}

	if height <= 0 {
		return 0, fmt.Errorf("рост должен быть больше 0, получено: %.2f м", height)
	}

	if duration <= 0 {
		return 0, fmt.Errorf("продолжительность должна быть больше 0, получено: %v", duration)
	}

	// Средняя скорость
	avgSpeed := meanSpeed(steps, height, duration)

	// Перевести продолжительность в минуты
	durationMinutes := duration.Minutes()

	// Рассчитать базовое количество калорий (как для бега)
	baseCalories := (weight * avgSpeed * durationMinutes) / minInH

	// Применить корректирующий коэффициент для ходьбы
	calories := baseCalories * walkingCaloriesCoefficient

	return calories, nil
}
