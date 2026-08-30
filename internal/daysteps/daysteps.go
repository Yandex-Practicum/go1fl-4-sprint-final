package daysteps

import (
	"fmt"
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

	sliceString := strings.Split(data, ",")
	if len(sliceString) != 2 {
		return 0, 0, fmt.Errorf("неверный формат строки: ожидается 2 части, найдено %d", len(sliceString))
	}
	step, ok := strconv.Atoi(sliceString[0])
	if ok != nil {
		return 0, 0, fmt.Errorf("неверный формат строки %s", sliceString[0])
	}
	if step <= 0 {
		return 0, 0, fmt.Errorf("Количество шагов %s ошибка", sliceString[0])
	}
	duration, err := time.ParseDuration(sliceString[1])
	if duration >= 24*time.Hour || duration <= 0 {
		return 0, 0, fmt.Errorf("ошибка времени")
	}

	if err != nil {
		return 0, 0, fmt.Errorf("ошибка времени")
	}
	return step, duration, nil
}

// TODO: реализовать функцию

func DayActionInfo(data string, weight, height float64) string {
	// TODO: реализовать функцию
	step, duration, err := parsePackage(data)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return ""
	}
	if step <= 0 {
		return ""
	}
	if duration <= 0 {
		return ""
	}
	durationStep := float64(step) * stepLength
	distance := durationStep / mInKm
	dayCalories, ok := spentcalories.WalkingSpentCalories(step, weight, height, duration)
	if ok != nil {
		fmt.Println("Ошибка расчета калорий:", err)
		return "" // Или обрабатываем ошибку как нужно по заданию
	}
	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", step, distance, dayCalories)

}
