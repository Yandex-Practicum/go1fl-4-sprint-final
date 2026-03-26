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

// parsePackage разбивает на отдельные позиции строку с данными о количестве шагов и продолжительности прогулки.
func parsePackage(data string) (int, time.Duration, error) {
	walking := strings.Split(data, ",")

	if len(walking) != 2 {
		return 0, 0, errors.New("ошибка: входящая строка не соответствует требованию о наличии 2 позиций")
	}
	stepsQuantity, err := strconv.Atoi(walking[0])
	if err != nil {
		return 0, 0, err
	}

	if stepsQuantity <= 0 {
		return 0, 0, errors.New("ошибка, значение количества шагов равно или менее нуля!")
	}

	walkDuration, err := time.ParseDuration(walking[1])
	if err != nil {
		return 0, 0, err
	}

	if walkDuration <= 0 {
		return 0, 0, errors.New("ошибка, значение длительности прогулки равно или менее нуля!")
	}

	return stepsQuantity, walkDuration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, walkDuration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}
	if steps <= 0 {
		return ""
	}

	distanceInKm := (float64(steps) * stepLength) / float64(mInKm)

	spentCalories, err := spentcalories.WalkingSpentCalories(steps, weight, height, walkDuration)
	if err != nil {
		fmt.Println(err.Error())
	}

	res := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distanceInKm, spentCalories)
	return res

}
