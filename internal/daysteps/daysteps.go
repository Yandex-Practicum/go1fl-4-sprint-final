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
	// TODO: реализовать функцию
	splitData := strings.Split(data, ",")
	if len(splitData) == 2 {

		steps, err := strconv.Atoi(splitData[0])
		if steps <= 0 {
			return 0, 0, fmt.Errorf("количество шагов должно быть больше 0")
		}

		if err != nil {
			return 0, 0, fmt.Errorf("ошибка при преобразовании количества шагов: %v", err)
		}

		duration, err := time.ParseDuration(splitData[1])
		if duration <= 0 {
			return 0, 0, fmt.Errorf("время должно быть больше 0")
		}
		if err != nil {
			return 0, 0, fmt.Errorf("ошибка при преобразовании времени")
		}
		return steps, duration, nil

	}
	return 0, 0, fmt.Errorf("Данные должны содержать 2 элемента")
}

func DayActionInfo(data string, weight, height float64) string {
	// TODO: реализовать функцию
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println("Ошибка:", err)
		return ""
	}
	if steps <= 0 {

		return ""
	}
	distance := float64(steps) * stepLength
	distance = distance / mInKm
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println("Ошибка:", err)
		return ""
	}
	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distance, calories)
}
