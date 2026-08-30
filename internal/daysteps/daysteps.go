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
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("неверный формат данных: ожидается два элемента, получено %d", len(parts))
	}

	stepsStr := parts[0]
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, 0, fmt.Errorf("не удалось преобразовать количество шагов в число: %v", err)
	}
	if steps <= 0 {

		return 0, 0, fmt.Errorf("количество шагов должно быть больше 0")

	}
	durationStr := parts[1]
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, 0, fmt.Errorf("не удалось преобразовать длительность в time.Duration: %v", err)
	}
	if duration <= 0 {
		return 0, 0, fmt.Errorf("продолжительность должна быть больше 0")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {

	steps, duration, err := parsePackage(data)
	if err != nil {

		message := err.Error()
		fmt.Println(message)
		return ""
	}
	if steps <= 0 {
		fmt.Println("количество шагов должно быть больше 0")
		return ""
	}
	distanceMeters := float64(steps) * stepLength

	disKm := distanceMeters / mInKm

	if disKm <= 0 {
		fmt.Println("дистанция должна быть больше 0")
		return ""
	}

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		errorMessage := err.Error()
		fmt.Println(errorMessage)
		return ""
	}

	if calories <= 0 {
		fmt.Println("количество калорий должно быть больше 0")
		return ""
	}

	result := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, disKm, calories)
	return result
}
