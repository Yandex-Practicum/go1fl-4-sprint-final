package daysteps

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	// 1. Разделяем строку
	parts := strings.Split(data, ",")

	// 2. Проверяем длину
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("неверный формат данных")
	}

	// 3. Парсим шаги и проверяем ошибки
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}
	if steps <= 0 {
		return 0, 0, fmt.Errorf("количество шагов должно быть положительным")
	}

	// 4. Парсим время и проверяем ошибки
	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, err
	}
	// Проверка на нулевое или отрицательное время 🛡️
	if duration <= 0 {
		return 0, 0, fmt.Errorf("время должно быть положительным")
	}

	// 5. Успешный возврат
	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	// 1. Извлекаем количество шагов и время из входной сырой строки
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}

	// 2. Проверяем, что шаги положительные
	if steps <= 0 {
		return ""
	}

	// 3. Рассчитываем пройденную дистанцию в километрах
	dist := (float64(steps) * stepLength) / mInKm

	// 4. Вычисляем потраченные калории при ходьбе и проверяем возможную ошибку
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println(err)
		return ""
	}

	// 5. Формируем и возвращаем итоговую строку с результатами активности за день
	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, dist, calories)
}
