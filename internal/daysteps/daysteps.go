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

var (
	ErrInvalidDataFormat = errors.New("некорректный формат входных данных")
	ErrInvalidStepCount  = errors.New("некорректное количество шагов")
	ErrInvalidDuration   = errors.New("некорректная длительность")
)

func parsePackage(data string) (int, time.Duration, error) {
	splittedData := strings.Split(data, ",")
	splittedDataLen := len(splittedData)

	if splittedDataLen != 2 {
		return 0, 0, ErrInvalidDataFormat
	}

	stepData := splittedData[0]
	durationData := splittedData[1]

	stepCount, err := strconv.Atoi(stepData)
	if err != nil {
		return 0, 0, err
	}

	if stepCount <= 0 {
		return 0, 0, ErrInvalidStepCount
	}

	duration, err := time.ParseDuration(durationData)
	if err != nil {
		return 0, 0, err
	}

	if duration <= 0 {
		return 0, 0, ErrInvalidDuration
	}

	return stepCount, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	stepCount, duration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}

	spentCalories, err := spentcalories.WalkingSpentCalories(stepCount, weight, height, duration)
	if err != nil {
		return ""
	}

	distance := (stepLength * float64(stepCount)) / mInKm

	return fmt.Sprintf(
		`Количество шагов: %d.
Дистанция составила %.2f км.
Вы сожгли %.2f ккал.
`,
		stepCount,
		distance,
		spentCalories,
	)
}
