package daysteps

import (
	"fmt"
	"go1fl-4-sprint-final-beh/internal/spentcalories"
	"strconv"
	"strings"
	"time"
	"log"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

// parsePackage парсит строку формата "steps,duration" и возвращает количество шагов и длительность.
func parsePackage(data string) (int, time.Duration, error) {
    parts := strings.Split(data, ",")
    if len(parts) != 2 {
        return 0, 0, fmt.Errorf("invalid data format: expected two parts separated by comma")
    }

    steps, err := strconv.Atoi(parts[0])
    if err != nil {
        return 0, 0, fmt.Errorf("invalid step count: %w", err)
    }
    if steps <= 0 {
        return 0, 0, fmt.Errorf("step count must be greater than zero")
    }

    duration, err := time.ParseDuration(parts[1])
    if err != nil {
        return 0, 0, err // ← Возвращаем оригинальную ошибку
    }
    if duration <= 0 {
        return 0, 0, fmt.Errorf("duration must be greater than zero")
    }

    return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
    steps, duration, err := parsePackage(data)
    if err != nil {
        log.Println(err) // ← Важно! Тест перехватывает именно log.Println
        return ""
    }

    distanceMeters := float64(steps) * stepLength
    distanceKm := distanceMeters / mInKm

    calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
    if err != nil {
        log.Printf("Error calculating calories: %v", err)
        return ""
    }

    return fmt.Sprintf(
        "Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
        steps,
        distanceKm,
        calories,
    )
}