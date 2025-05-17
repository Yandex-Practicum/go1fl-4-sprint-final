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
	value := strings.Split(data, ",")
	if len(value) != 2 {
		return 0, 0, errors.New("!= 2")
	}

	valueStep, err := strconv.Atoi(value[0])
	if err != nil {
		return 0, 0, err
	}
	if valueStep <= 0 {
		return 0, 0, errors.New("0 шагов")
	}

	valueDuration, err := time.ParseDuration(value[1])
	if err != nil {
		return 0, 0, err
	}
	if valueDuration <= 0 {
		return 0, 0, errors.New("0 часов")
	}

	return valueStep, valueDuration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	valueStep, valueDuration, err := parsePackage(data)
	if err != nil {
		log.Printf("Ошибка парсинга: %v", err)
		return ""
	}
	if valueStep <= 0 {
		return ""
	}
	if valueDuration <= 0 {
		return ""
	}

	distance := (float64(valueStep) * stepLength) / mInKm
	callories, err := spentcalories.WalkingSpentCalories(valueStep, weight, height, valueDuration)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", valueStep, distance, callories)

}
