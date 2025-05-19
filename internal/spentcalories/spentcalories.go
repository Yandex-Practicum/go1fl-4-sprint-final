package spentcalories

import (
	"time"
	"strconv"
	"strings"
	"errors"
	"fmt"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

// parseTraining парсит данные о тренировке.
func parseTraining(data string) (int, string, time.Duration, error) {
    // Разделяем строку на части
    parts := strings.Split(data, ",")
    if len(parts) != 3 {
        return 0, "", 0, errors.New("неверный формат данных: требуется 3 поля — шаги, тип активности, время")
    }

    // Парсим количество шагов
    stepsStr := parts[0]
    steps, err := strconv.Atoi(stepsStr)
    if err != nil {
        return 0, "", 0, errors.New("невозможно преобразовать шаги в число")
    }
    if steps <= 0 {
        return 0, "", 0, errors.New("количество шагов должно быть больше нуля")
    }

    // Получаем тип активности
    activity := parts[1]

    // Парсим продолжительность
    durationStr := parts[2]
    duration, err := time.ParseDuration(durationStr)
    if err != nil {
        return 0, "", 0, errors.New("неверный формат продолжительности")
    }

    // Проверяем, что продолжительность больше нуля
    if duration <= 0 {
        return 0, "", 0, errors.New("продолжительность должна быть больше нуля")
    }

    return steps, activity, duration, nil
}

// distance рассчитывает пройденное расстояние в километрах.
func distance(steps int, height float64) float64 {
    // Длина одного шага
    stepLength := height * stepLengthCoefficient

    // Общее расстояние в метрах
    totalDistanceMeters := float64(steps) * stepLength

    // Переводим в километры
    return totalDistanceMeters / mInKm
}

// meanSpeed рассчитывает среднюю скорость движения в км/ч.
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
    if duration <= 0 {
        return 0
    }

    // Рассчитываем пройденное расстояние в километрах
    distanceKm := distance(steps, height)

    // Переводим длительность в часы
    hours := duration.Hours()

    // Средняя скорость: км / час
    return distanceKm / hours
}

// TrainingInfo возвращает информацию о тренировке в виде строки.
func TrainingInfo(data string, weight, height float64) (string, error) {
    // Парсим данные из строки
    steps, activity, duration, err := parseTraining(data)
    if err != nil {
        return "", err
    }

    // Рассчитываем базовые показатели
    distanceKm := distance(steps, height)
    speed := meanSpeed(steps, height, duration)

    // Переменные для результата
    var calories float64

    // Определяем тип тренировки и считаем калории
    switch activity {
    case "Ходьба":
        calories, err = WalkingSpentCalories(steps, weight, height, duration)
    case "Бег":
        calories, err = RunningSpentCalories(steps, weight, height, duration)
    default:
        return "", fmt.Errorf("неизвестный тип тренировки: %s", activity)
    }

    if err != nil {
        return "", err
    }

    // Формируем строку с результатом
    return fmt.Sprintf(
        "Тип тренировки: %s\n"+
            "Длительность: %.2f ч.\n"+
            "Дистанция: %.2f км.\n"+
            "Скорость: %.2f км/ч\n"+
            "Сожгли калорий: %.2f\n",
        activity,
        duration.Hours(),
        distanceKm,
        speed,
        calories,
    ), nil
}

// RunningSpentCalories рассчитывает количество калорий, потраченных при беге.
// RunningSpentCalories рассчитывает калории, потраченные при беге, с учетом роста.
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
    if steps <= 0 {
        return 0, errors.New("количество шагов должно быть больше нуля")
    }
    if weight <= 0 {
        return 0, errors.New("вес должен быть больше нуля")
    }
    if height <= 0 {
        return 0, errors.New("рост должен быть больше нуля")
    }
    if duration <= 0 {
        return 0, errors.New("длительность должна быть больше нуля")
    }

    // Рассчитываем среднюю скорость через meanSpeed (уже использует height)
    speed := meanSpeed(steps, height, duration)

    // Переводим время в минуты
    durationInMinutes := duration.Minutes()

    // Формула: (вес * скорость * минуты) / минВЧасе
    calories := (weight * speed * durationInMinutes) / minInH

    return calories, nil
}

// WalkingSpentCalories рассчитывает количество калорий, потраченных при ходьбе.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
    if steps <= 0 {
    return 0, errors.New("steps must be greater than zero")
	}
	if weight <= 0 {
		return 0, errors.New("weight must be greater than zero")
	}
	if height <= 0 {
		return 0, errors.New("height must be greater than zero")
	}
	if duration <= 0 {
		return 0, errors.New("duration must be greater than zero")
	}

    speed := meanSpeed(steps, height, duration)

    durationInMinutes := duration.Minutes()

    calories := (weight * speed * durationInMinutes) / minInH
    calories *= walkingCaloriesCoefficient

    return calories, nil
}