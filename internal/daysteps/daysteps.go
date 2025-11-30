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

var (
	ErrMissingArguments    = errors.New("daysteps: missing arguments")
	ErrorStringToInt       = errors.New("daysteps: cannot convert string to int")
	ErrorZeroSteps         = errors.New("daysteps: zero steps")
	ErrorParseToDuration   = errors.New("daysteps: cannot parse duration")
	ErrorIncorrectArgument = errors.New("daysteps: incorrect argument")
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	dataSlice := strings.Split(data, ",")
	if len(dataSlice) != 2 {
		return 0, 0, fmt.Errorf("%w: missing input, e.g. steps, duration; parsePackage", ErrMissingArguments)
	}

	steps, err := strconv.Atoi(dataSlice[0])
	if err != nil {
		return 0, 0, fmt.Errorf("%w: something went wrong while converting string to int; parsePackage", ErrorStringToInt)
	}

	if steps <= 0 {
		return 0, 0, fmt.Errorf("%w: zero steps; parsePackage", ErrorZeroSteps)
	}

	duration, err := time.ParseDuration(dataSlice[1])
	if err != nil {
		return 0, 0, fmt.Errorf("%w: something went wrong while parsing string to duration; parsePackage", ErrorParseToDuration)
	}
	if duration <= 0 {
		return 0, 0, fmt.Errorf("%w: duration cannot be less than zero; parsePackage", ErrorIncorrectArgument)
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	if weight <= 0 {
		log.Printf("weight is incorrect; DAyActionInfo")
		return ""
	}
	if height <= 0 {
		log.Printf("height is incorrect; DayActionInfo")
		return ""
	}

	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Print(err)
		return ""
	}

	distance := float64(steps) * stepLength
	distance = distance / mInKm
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Print(err)
		return ""
	}
	str := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distance, calories)
	return str
}
