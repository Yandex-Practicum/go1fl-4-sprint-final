package daysteps

import (
	"errors"
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
	parts := strings.Split(data, ",")

	// проверяем количество частей
	if len(parts) != 2 {

		return 0, 0, errors.New("две части")
	}

	stepStr := parts[0]
	durationStr := parts[1]

	// Парсим количество шагов (первый элемент)
	step, err := strconv.Atoi(stepStr)
	if err != nil {
		log.Println(err.Error())
		return 0, 0, errors.New("неправильное количество шагов")
	}
	if step <= 0 {

		return 0, 0, errors.New("количество шагов должно быть больше нуля")
	}

	// Парсим продолжительность (второй элемент)
	duration, err := time.ParseDuration(durationStr)
	if err != nil {

		return 0, 0, err
	}
	if duration <= 0 {
		return 0, 0, errors.New("продолжительность должна быть боьше нуля")
	}

	return step, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}
	if steps <= 0 {
		log.Println("error")
		return ""
	}

	distanceMeters := float64(steps) * stepLength
	distanceMeters = distanceMeters / mInKm
	calories, _ := spentcalories.WalkingSpentCalories(steps, weight, height, duration)

	result := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distanceMeters, calories)

	return result
}
