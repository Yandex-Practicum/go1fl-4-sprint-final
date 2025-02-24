package daysteps

import (
	"errors"
	"fmt"
	"github.com/Yandex-Practicum/go1fl-4-sprint-final/internal/spentcalories"
	"strconv"
	"strings"
	"time"
)

var (
	StepLength = 0.65 // длина шага в метрах
)

func parsePackage(data string) (int, time.Duration, error) {
	// ваш код ниже
	str := strings.Split(data, ",")
	if len(str) != 2 {
		return 0, 0, errors.New("error split data")
	}
	count, err := strconv.Atoi(str[0])
	if err != nil {
		return 0, 0, err
	}
	duration, err := time.ParseDuration(str[1])
	if err != nil {
		return 0, 0, err
	}
	return count, duration, nil
}

// DayActionInfo обрабатывает входящий пакет, который передаётся в
// виде строки в параметре data. Параметр storage содержит пакеты за текущий день.
// Если время пакета относится к новым суткам, storage предварительно
// очищается.
// Если пакет валидный, он добавляется в слайс storage, который возвращает
// функция. Если пакет невалидный, storage возвращается без изменений.
func DayActionInfo(data string, weight, height float64) string {
	// ваш код ниже
	count, duration, err := parsePackage(data)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if count <= 0 {
		return ""
	}
	distance := float64(count) * StepLength / 1000
	calories := spentcalories.WalkingSpentCalories(count, weight, height, duration)
	return fmt.Sprintf("Количество шагов: %d.\n Дистанция составила %.2f км.\n Вы сожгли %.2f ккал.\n", count, distance, calories)
}
