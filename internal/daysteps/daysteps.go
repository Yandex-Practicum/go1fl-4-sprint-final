package daysteps

import (
	"errors"
	"strconv"
	"time"
)

// Константы
const (
	stepLength = 0.7 // длина шага в метрах
	mInKm      = 1000 // метров в километре
)

// parsePackage парсит строку с данными о шагах и продолжительности прогулки.
func parsePackage(data string) (int, time.Duration, error) {
	parts := splitData(data)
	if len(parts) != 2 {
		return 0, 0, errors.New("некорректный формат данных: ожидается два элемента")
	}

	stepsStr, durationStr := parts[0], parts[1]

	// Преобразуем количество шагов в int
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, 0, errors.New("ошибка преобразования количества шагов: " + err.Error())
	}
	if steps <= 0 {
		return 0, 0, errors.New("количество шагов должно быть больше нуля")
	}

	// Преобразуем продолжительность в time.Duration
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, 0, errors.New("ошибка парсинга продолжительности: " + err.Error())
	}

	return steps, duration, nil
}

// DayActionInfo формирует строку с информацией о прогулке.
func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		println(err.Error())
		return ""
	}
	if steps <= 0 {
		return ""
	}

	// Вычисляем дистанцию в километрах
	distanceM := float64(steps) * stepLength
	distanceKm := distanceM / mInKm

	// Вычисляем калории (используем функцию из пакета spentcalories)
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		println("Ошибка вычисления калорий: " + err.Error())
		return ""
	}

	// Формируем итоговую строку
	return fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.",
		steps, distanceKm, calories,
	)
}

// splitData разделяет строку по запятой.
func splitData(data string) []string {
	return strings.Split(data, ",")
}
