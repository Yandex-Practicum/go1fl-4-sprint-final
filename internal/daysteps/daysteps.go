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
	// TODO: реализовать функцию

	dataSlice1 := strings.Split(data, ",")

	if len(dataSlice1) != 2 {
		return 0, 0, errors.New("Ошибка в parsePackage:мало данных")
	}

	stepD, err := strconv.Atoi(dataSlice1[0])
	if err != nil {
		return 0, 0, errors.New("Ошибка в parsePackage : преобразование строки в число")
	} else if stepD <= 0 {
		return 0, 0, errors.New("Ошибка в parsePackage : шаги равны 0")
	}

	timeD, err := time.ParseDuration(dataSlice1[1])
	if err != nil {
		return 0, 0, errors.New("Ошибка в parsePackage : перевод строки в часы")
	}

	return stepD, timeD, nil

}

func DayActionInfo(data string, weight, height float64) string {
	// TODO: реализовать функцию

	step, _, err := parsePackage(data)
	if err != nil {
		return fmt.Sprint(err)
	}

	if step <= 0 {
		return ""
	}

	distance := (float64(step) * stepLength) / mInKm

	_, time, _ := parsePackage(data)

	kkal, err := spentcalories.WalkingSpentCalories(step, weight, height, time)

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\n Вы сожгли %.2f ккал.",
		step, distance, kkal)
}
