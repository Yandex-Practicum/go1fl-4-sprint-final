package daysteps

import (
	"errors"
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
		return 0, 0, errors.New("можно принимать не больше двух значений")
	}
	stepString := parts[0]
	step, err := strconv.Atoi(stepString)
	if err != nil {

		return 0, 0, fmt.Errorf("не удалось преобразовать шаги %v", err)
	}
	if step <= 0 {

		return 0, 0, fmt.Errorf("количество шагов %d", step)
	}
	durationString := parts[1]
	duration, err := time.ParseDuration(durationString)

	if err != nil {
		return 0, 0, errors.New("некорректная длительность")
	}
	return step, duration, nil

}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		fmt.Println("Ошибка:", err.Error())
		return ""
	}
	if duration <= 0 {
		fmt.Println("Ошибка: длительность должна быть > 0, получено:", duration)
		return ""
	}
	if steps <= 0 {
		fmt.Println("Количество шагов должно быть больше 0")
		return ""
	}
	distance := float64(steps) * stepLength
	DisKm := distance / mInKm
	realisation, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		fmt.Println("ошибка", err.Error())
	}

	result := fmt.Sprintf(`Количество шагов: %d
       Дистанция составила %.2f км.
       Вы сожгли %.2f ккал.`,
		steps,
		DisKm,
		realisation,
	)
	return result
}
