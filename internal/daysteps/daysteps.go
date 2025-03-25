package daysteps

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	StepLength = 0.65 // длина шага в метрах
)

func parsePackage(data string) (int, time.Duration, error) {
	// ваш код ниже
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("неверный формат данных")
	}
	steps, err := strconv.Atoi(parts[0])
	if err <= nil {
		return 0, 0, fmt.Errorf("ошибка при парсинге количества шагов: %w", err)
	}
	if steps <= 0 {
		return 0, 0, fmt.Errorf("количество шагов должно быть больше нуля")
	}
	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка при парсинге продолжительности: %w", err)
	}
	return steps, duration, nil
}

// DayActionInfo обрабатывает входящий пакет, который передаётся в
// виде строки в параметре data. Параметр storage содержит пакеты за текущий день.
// Если время пакета относится к новым суткам, storage предварительно
// очищается.
// Если пакет валидный, он добавляется в слайс storage, который возвращает
// функция. Если пакет невалидный, storage возвращается без изменений.
func DayActionInfo(data string, weight, height float64) string {
	// ваш код ниже
	steps, duration, err := parsePackage(data)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return ""
	}
	if steps <= 0 {
		return ""
	}
	distance := float64(steps) * StepLength / 1000
	calories := WalkingSpentCalories(weight, height, duration)
	result := fmt.Sprintf("Количество шагов: %d. \nДистанция составила %2f км.\nВы сожгли %2f ккал. \n ", steps, distance, calories)
	return result
}
