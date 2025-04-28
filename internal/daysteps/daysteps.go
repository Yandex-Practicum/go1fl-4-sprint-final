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
	// TODO: реализовать функцию
	if data == "" {
		return 0, 0, errors.New("no data to process")
	}
	slData := strings.Split(data, ",") // конвертация строки в слайс slDate
	if len(slData) != 2 {
		return 0, 0, errors.New("incorrect amount of data")
	}

	steps, err := strconv.Atoi(slData[0]) // преобразование в int
	if err != nil {
		return 0, 0, errors.New("step count conversion error")
	}
	if steps <= 0 {
		return 0, 0, errors.New("the number of steps must be greater than 0")
	}

	duration, err := time.ParseDuration(slData[1]) // конвертация в time
	if err != nil {
		return 0, 0, errors.New("duration conversion error")
	}
	if duration <= 0 {
		return 0, 0, errors.New("duration must be greater than 0")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	// TODO: реализовать функцию
	steps, duration, err := parsePackage(data) //
	if err != nil {
		log.Println("error reading data: ", err)
		return ""
	}

	distance := (stepLength * float64(steps)) / float64(mInKm)

	countKcal, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println("error when counting calories: ", err)
		return ""
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distance, countKcal)
}
