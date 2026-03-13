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
	//a := parts[0]
	//b := parts[1]
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("Not right parts")
	}
	a, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("Error conv a %w", err)
	}
	if a <= 0 {
		return 0, 0, fmt.Errorf("0 - quantity of steps")
	}
	//a := strconv.ParseInt(parts[0], 10, 64)
	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("Error conv duration %w", err)
	} else {
		return a, duration, nil
	}

}
func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		fmt.Printf("Error %w", err)
		return ""
	} else if steps <= 0 {
		return ""
	} else {
		cal, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return fmt.Sprintf("Error in dayaction: ", err)
		}
		dis := float64(steps) * stepLength / float64(mInKm)
		return fmt.Sprintf("Количество шагов: %d\nДистанция составила %.2f км.\nВы сожгли %.2f ккал. ", steps, dis, cal)
	}
}
