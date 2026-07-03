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
	dataSlice := strings.Split(data, ",")
	if len(dataSlice) != 2 {
		return 0, time.Duration(0), errors.New("Неверный формат данных")
	}

	steps, err := strconv.Atoi(dataSlice[0])
	if err != nil {
		return 0, time.Duration(0), errors.New("Неверный формат количества шагов")
	}
	if steps <= 0 {
		return 0, time.Duration(0), errors.New("Количество шагов должно быть больше 0")
	}

	duration, err := time.ParseDuration(dataSlice[1])
	if err != nil {
		return 0, time.Duration(0), errors.New("Неверный формат продолжительности прогулки")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		return fmt.Sprintf("Не получилось получить информацию о прогулке: %v", err)
	}
	if steps <= 0 {
		return ""
	}

	distance := float64(steps) * stepLength
	distance = distance / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		return fmt.Sprintf("Не получилось получить информацию о калориях: %v", err)
	}

	message := fmt.Sprintf("Количество шагов: %d./nДистанция составила %.2f км./nВы сожгли %.2f ккал.", steps, distance, calories)
	return message
}
