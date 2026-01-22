package spentcalories

import (
	"errors"
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

// Парсит строку с данными о тренировке
func parseTraining(data string) (int, string, time.Duration, error) {

	// Разделяем строку по запятой
	training := strings.Split(data, ",")

	if len(training) != 3 {
		return 0, "", 0, errors.New("Неверный формат")
	}

	// Парсим количество шагов
	number, err := strconv.Atoi(training[0])

	if err != nil {
		return 0, "", 0, err
	}
	if number <= 0 {
		return 0, "", 0, errors.New("Количество шагов должно быть больше 0")
	}

	// Парсим продолжительность
	duration, err := time.ParseDuration(training[2])

	if err != nil {
		return 0, "", 0, err
	}
	// Проверяем продолжительность больше 0
	if duration <= 0 {
		return 0, "", 0, errors.New("Продолжительность должна быть больше 0")
	}

	return number, training[1], duration, nil

}

// Вычисляет дистанцию в километра
func distance(steps int, height float64) float64 {

	lenStep := (height * stepLengthCoefficient) * float64(steps) / mInKm
	return lenStep
}

// Вычисляет среднюю скорость
func meanSpeed(steps int, height float64, duration time.Duration) float64 {

	if duration <= 0 {
		return 0
	}

	dist := distance(steps, height)
	speed := dist / duration.Hours()

	return speed

}

// TrainingInfo парсит данные о тренировке и возвращает информацию о дистанции, скорости и калориях.
func TrainingInfo(data string, weight, height float64) (string, error) {

	// Парсим строку с данными о тренировке
	steps, trainingType, duration, err := parseTraining(data)

	if err != nil {
		log.Println(err)
		return "", err
	}

	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	// Вычисляем калории в зависимости от типа тренировки
	var calories float64

	switch trainingType {
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	if err != nil {
		log.Println(err)
		return "", err
	}

	result := fmt.Sprintf(
		"Тип тренировки: %s\n"+
			"Длительность: %.2f ч.\n"+
			"Дистанция: %.2f км.\n"+
			"Скорость: %.2f км/ч\n"+
			"Сожгли калорий: %.2f\n",
		trainingType,
		duration.Hours(),
		dist,
		speed,
		calories,
	)

	return result, nil

}

// RunningSpentCalories вычисляет количество калорий, потраченных при беге.
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	// Проверка входных параметров
	if steps <= 0 {
		return 0, errors.New("Количество шагов должно быть больше 0")
	}
	if weight <= 0 {
		return 0, errors.New("Вес должен быть больше 0")
	}
	if height <= 0 {
		return 0, errors.New("Рост должен быть больше 0")
	}
	if duration <= 0 {
		return 0, errors.New("Продолжительность должна быть больше 0")
	}

	speed := meanSpeed(steps, height, duration)
	calories := (duration.Minutes() * weight * speed) / minInH

	return calories, nil

}

// WalkingSpentCalories вычисляет количество калорий, потраченных при ходьбе.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	// Проверка входных параметров
	if steps <= 0 {
		return 0, errors.New("Количество шагов должно быть больше 0")
	}
	if weight <= 0 {
		return 0, errors.New("Вес должен быть больше 0")
	}
	if height <= 0 {
		return 0, errors.New("Рост должен быть больше 0")
	}
	if duration <= 0 {
		return 0, errors.New("Продолжительность должна быть больше 0")
	}

	speed := meanSpeed(steps, height, duration)
	calories := (duration.Minutes() * weight * speed) / minInH

	return calories * walkingCaloriesCoefficient, nil

}
