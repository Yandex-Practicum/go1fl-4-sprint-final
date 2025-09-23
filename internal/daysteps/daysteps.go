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
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("parsePackage: expected 2 fields (steps,duration), got %d", len(parts))
	}

	stepsStr := strings.TrimSpace(parts[0])
	steps, err := strconv.Atoi(stepsStr)
	if err != nil || steps <= 0 {
		return 0, 0, fmt.Errorf("parsePackage: invalid steps %q", stepsStr)
	}

	durStr := strings.TrimSpace(parts[1])
	dur, err := time.ParseDuration(durStr)
	if err != nil || dur <= 0 {
		return 0, 0, fmt.Errorf("parsePackage: invalid duration %q", durStr)
	}

	return steps, dur, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, dur, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}

	distKm := (float64(steps) * stepLength) / float64(mInKm)

	cals, err := spentcalories.WalkingSpentCalories(steps, weight, height, dur)
	if err != nil {
		log.Println(err)
		return ""
	}

	return fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.",
		steps, distKm, cals,
	)
}
