package daysteps

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"github.com/Sergey-81/Fitness-tracker-module/internal/spentcalories"
)

const (
	stepLength = 0.65
	mInKm      = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid data format: expected two parts separated by comma")
	}

	// Жесткая проверка на пробелы в числах
	stepsStr := parts[0]
	if strings.ContainsAny(stepsStr, " \t\n\r") {
		return 0, 0, fmt.Errorf("step count contains spaces")
	}

	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid step count: %v", err)
	}
	if steps <= 0 {
		return 0, 0, fmt.Errorf("step count must be greater than zero")
	}

	durationStr := parts[1]
	if strings.ContainsAny(durationStr, " \t\n\r") {
		return 0, 0, fmt.Errorf("duration contains spaces")
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid duration format: %v", err)
	}
	if duration <= 0 {
		return 0, 0, fmt.Errorf("duration must be greater than zero")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
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