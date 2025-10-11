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

// Функция принимает строку с данными формата "3456,Ходьба,3h00m", которая содержит количество шагов, вид активности и продолжительность активности. Функция возвращает четыре значения:
// int — количество шагов.
// string — вид активности.
// time.Duration — продолжительность активности.
// error — ошибку, если что-то пошло не так.
func parseTraining(data string) (int, string, time.Duration, error) {
	// Делит строку на слайс строк
	parts := strings.Split(data, ",")

	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("неверный формат: ожидалось 3 части, получили %d", len(parts))
	}
	// Парсим шаги
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("не удалось преобразовать шаги: %w", err)
	}
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("количество шагов должно быть положительным числом")
	}
	activity := parts[1]
	// парсим продолжительсть
	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка парсинга длительности: %w", err)
	}
	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("продолжительность должна быть положительной")
	}
	return steps, activity, duration, nil
}

// distance
// Функция принимает количество шагов и рост пользователя в метрах, а возвращает дистанцию в километрах.
func distance(steps int, height float64) float64 {

	lenS := height * stepLengthCoefficient
	distance := (float64(steps) * lenS) / mInKm

	return distance
}

// Функция принимает количество шагов steps, рост пользователя height и продолжительность активности duration  и возвращает среднюю скорость.
func meanSpeed(steps int, height float64, duration time.Duration) float64 {

	if duration <= 0 {
		return 0
	}
	distance := distance(steps, height)
	return distance / duration.Hours()
}

// TrainingInfo
// Функция принимает:
// data string — строку с данными формата "3456,Ходьба,3h00m", которая содержит количество шагов, вид активности и продолжительность активности.
// weight, height float64 — вес (кг.) и рост (м.) пользователя.
// И возвращает два значения:
// string — строка с информацией о тренировке в формате, приведенном ниже.
// error — ошибку, при ее возникновении внутри функции.
func TrainingInfo(data string, weight, height float64) (string, error) {
	var distanceInfo, meanSpeedInfo, caloriesInfo float64
	var err error
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}
	if steps <= 0 {
		return "", fmt.Errorf("количество шагов должно быть положительным числом")
	}

	switch activity {
	case "Бег":
		//stepLength := height * lenStep
		//distanceInfo = (float64(steps) * stepLength) / mInKm
		distanceInfo = distance(steps, height)
		meanSpeedInfo = meanSpeed(steps, height, duration)
		caloriesInfo, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
			return "", err
		}
	case "Ходьба":
		distanceInfo = distance(steps, height)
		meanSpeedInfo = meanSpeed(steps, height, duration)
		caloriesInfo, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
			return "", err
		}
	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}

	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activity, duration.Hours(), distanceInfo, meanSpeedInfo, caloriesInfo), nil

}

// RunningSpentCalories
// Функция принимает:
// steps int — количество шагов.
// weight, height float64 — вес(кг.) и рост(м.) пользователя.
// duration time.Duration — продолжительность бега.
// И возвращает два значения:
// float64 — количество калорий, потраченных при беге.
// error — ошибку, если входные параметры некорректны (подумайте, какие значения параметров имеют смысл).
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("Некорректные входные данные для подсчета")
	}
	meanSpeed := meanSpeed(steps, height, duration)

	durationInMinutes := duration.Minutes()
	calories := (weight * meanSpeed * durationInMinutes) / minInH
	return calories, nil
}

// WalkingSpentCalories
// Функция принимает:
// steps int — количество шагов.
// weight, height float64 — вес(кг.) и рост(м.) пользователя.
// duration time.Duration — продолжительность ходьбы.
// И возвращает два значения:
// float64 — количество калорий, потраченных при ходьбе.
// error — ошибку, если входные параметры некорректны (подумайте, какие значения параметров имеют смысл).
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("Некорректные входные данные для подсчета")
	}
	meanSpeed := meanSpeed(steps, height, duration)

	durationInMinutes := duration.Minutes()
	calories := ((weight * meanSpeed * durationInMinutes) / minInH) * walkingCaloriesCoefficient
	return calories, nil
}
