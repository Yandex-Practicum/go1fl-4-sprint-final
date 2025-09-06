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
	parseData := strings.Split(data, ",")
	if len(parseData) != 2 {
		return 0, time.Duration(0), errors.New("incorrect input data")
	}

	step, err := strconv.Atoi(parseData[0])
	if err != nil {
		return 0, time.Duration(0), err
	}
	if step <= 0 {
		return 0, time.Duration(0), errors.New("the number of steps is less than zero")
	}

	walkingTime, err := time.ParseDuration(parseData[1])
	if err != nil {
		return 0, time.Duration(0), err
	}
	if walkingTime <= time.Duration(0) {
		return 0, time.Duration(0), errors.New("duration <= 0")
	}

	return step, walkingTime, nil

}

func DayActionInfo(data string, weight, height float64) string {
	step, walkingTime, err := parsePackage(data)
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	if step <= 0 {
		return ""
	}
	distance := (float64(step) * stepLength) / mInKm

	calorie, err := spentcalories.WalkingSpentCalories(step, weight, height, walkingTime)
	if err != nil {
		//fmt.Println(err)
		return ""
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", step, distance, calorie)

}
