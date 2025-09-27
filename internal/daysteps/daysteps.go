package daysteps

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"go1fl-4-sprint-final/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	info := strings.Split(data, ",")

	if len(info) != 2 {
		return 0, 0, fmt.Errorf("не указано количество шагов или продолжительность прогулки")
	}

	steps, err := strconv.Atoi(info[0])
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка преобразования количества шагов")
	}

	if steps <= 0 {
		return 0, 0, fmt.Errorf("отсутствуют шаги")
	}

	duration, err := time.ParseDuration(info[1])
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка преобразования времени")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		fmt.Println("")
	}

	if steps <= 0 {
		fmt.Println("")
	}

	dist := (stepLength * float64(steps)) / float64(mInKm)

	ccal, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		return fmt.Sprintf("ошибка расчета калорий: %v", err)
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %2.f км.\nВы сожгли %2.f ккал.", steps, dist, ccal)
}
