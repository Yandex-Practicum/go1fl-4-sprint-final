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
	// 1. Разделяем строку на слайс строк по запятой
	parts := strings.Split(data, ",")

	// 2. Проверяем, чтобы длина слайса была ровно 3
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("неверный формат данных")
	}

	// 3. Преобразуем первый элемент  в тип int
	steps, err := strconv.Atoi(parts[0])
	if err != nil { // Если возникла ошибка, возвращаем нули, пустую строку и саму ошибку
		return 0, "", 0, err
	}

	// 4. Преобразуем третий элемент в time.Duration
	duration, err := time.ParseDuration(parts[2])
	if err != nil { // Если возникла ошибка, снова возвращаем нули и ошибку
		return 0, "", 0, err
	}

	// 5. Если всё прошло без ошибок, возвращаем шаги, вид активности, время и nil
	return steps, parts[1], duration, nil
}

func distance(steps int, height float64) float64 {
	// 1. Рассчитываем длину шага
	stepLen := height * stepLengthCoefficient

	// 2. Вычисляем дистанцию в километрах
	dist := (float64(steps) * stepLen) / mInKm

	// 3. Возвращаем результат
	return dist
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	// 1. Если время меньше или равно нулю, возвращаем 0
	if duration <= 0 {
		return 0
	}

	// 2. Вычисляем дистанцию
	dist := distance(steps, height)

	// 3. Возвращаем среднюю скорость
	return dist / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	// 1. Извлекаем данные из сырой строки
	steps, activityType, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err) // Логируем ошибку, как указано в ТЗ
		return "", err
	}

	// 2. Считаем универсальные показатели один раз
	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	var calories float64

	// 3. Считаем калории в зависимости от типа активности
	switch activityType {
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	default:
		// Если пришло что-то непонятное
		return "", fmt.Errorf("неизвестный тип тренировки")
	}

	// 4. Формируем отчет и возвращаем его
	return fmt.Sprintf("Тип тренировки: %s \nДлительность: %.2f ч. \nДистанция: %.2f км. \nСкорость: %.2f км/ч \nСожгли калорий: %.2f \n", activityType, duration.Hours(), dist, speed, calories), nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// 1. Проверяем параметры на нули и отрицательные значения
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("некорректные параметры тренировки")
	}

	// 2. Вычисляем среднюю скорость
	speed := meanSpeed(steps, height, duration)

	// 3. Рассчитываем и возвращаем калории и nil
	return (weight * speed * duration.Minutes()) / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// 1. Проверяем параметры на нули и отрицательные значения
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("некорректные параметры тренировки")
	}

	// 2. Вычисляем среднюю скорость
	speed := meanSpeed(steps, height, duration)

	// 3. Рассчитываем и возвращаем калории и nil
	return (weight * speed * duration.Minutes()) / minInH * walkingCaloriesCoefficient, nil
}
