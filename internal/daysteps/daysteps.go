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
		return 0, 0, fmt.Errorf("количество шагов %s ошибка", sliceString[0]) //1
	}
	duration, err := time.ParseDuration(sliceString[1])
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка времени")
	}
	if duration >= 24*time.Hour || duration <= 0 {
		return 0, 0, fmt.Errorf("ошибка времени")
	}
	return step, duration, nil
}
func DayActionInfo(data string, weight, height float64) string {
	step, duration, err := parsePackage(data)
	if err != nil {
		fmt.Println("ошибка:", err)
		return ""
	}
	if step <= 0 { //1
		return ""
	}
	if duration <= 0 {
		return ""
	}
	durationStep := float64(step) * stepLength
	distance := durationStep / mInKm
	dayCalories, ok := spentcalories.WalkingSpentCalories(step, weight, height, duration)
	if ok != nil {
		fmt.Println("ошибка расчета калорий:", ok)
		return ""
	}
	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", step, distance, dayCalories)
}
