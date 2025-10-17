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

func parseTraining(data string) (int, string, time.Duration, error) {

	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, errors.New("неверный формат входных данных")
	}
	// Преобразуем первую часть в количество шагов
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("неверное количество шагов: %w", err)
	}
	if steps <= 0 {
		return 0, "", 0, errors.New("неверное количество шагов!")
	}
	// Преобразуем вторую часть данных в тип тренировки
	activity := parts[1]
	// Преобразуем третью часть в продолжительность тренировки (например, "45m" или "1h30m")
	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка преобразования продолжительности: %w", err)
	}
	if duration <= 0 {
		return 0, "", 0, errors.New("неверная продолжительность тренировки!")
	}

	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {
	// convLength переводит шаги в Км учитывая рост
	convLength := float64(steps) * lenStep * height * stepLengthCoefficient / mInKm
	return convLength
}

// speed
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	speed := distance(steps, height) / duration.Hours()
	return speed
}

func TrainingInfo(data string, weight, height float64) (string, error) {

	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println("Ошибка конвертации данных:", err)
		return "", err
	}
	if weight <= 0 || height <= 0 {
		err := errors.New("вес и рост должны быть больше нуля")
		log.Println("Ошибка данных: %w", err)
		return "", err
	}

	convLength := distance(steps, height)
	speed := meanSpeed(steps, height, duration)
	durationHours := duration.Hours()
	var calories float64
	//определяем тип тренировки и расчитываем калории
	switch activity {
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", errors.New("неизвестный тип тренировки")
	}
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", activity, durationHours, convLength, speed, calories), nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if weight <= 0 || height <= 0 {
		err := errors.New("вес и рост должны быть больше нуля")
		return 0, err
	}
	if steps <= 0 {
		return 0, errors.New("шаги должны быть положительным числом")
	}
	if duration <= 0 {
		return 0, errors.New("продолжительность должна быть больше нуля")
	}
	meanSpeed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()
	calories := (weight * meanSpeed * durationInMinutes) / minInH

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if weight <= 0 || height <= 0 {
		err := errors.New("вес и рост должны быть больше нуля")
		return 0, err
	}
	if steps <= 0 {
		return 0, errors.New("шаги должны быть положительным числом")
	}
	if duration <= 0 {
		return 0, errors.New("продолжительность должна быть больше нуля")
	}
	meanSpeed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()
	calories := (weight * meanSpeed * durationInMinutes) / minInH

	return calories * walkingCaloriesCoefficient, nil
}
