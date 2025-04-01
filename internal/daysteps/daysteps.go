package daysteps

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/go1fl-4-sprint-final/internal/spentcalories"
)

var (
	StepLength = 0.65 // длина шага в метрах
)

func parsePackage(data string) (int, time.Duration, error) {
	twoString := strings.Split(data, ",")
	if len(twoString) != 2 {
		return 0, 0, fmt.Errorf("Недостаточно элементов")
	}
	steps, err := strconv.Atoi(twoString[0])
	if err != nil || steps <= 0 {
		return 0, 0, err
	}
	durationWalk, err2 := time.ParseDuration(twoString[1])
	if err2 != nil {
		return 0, 0, err2
	}
	return steps, durationWalk, nil
}

// DayActionInfo обрабатывает входящий пакет, который передаётся в
// виде строки в параметре data. Параметр storage содержит пакеты за текущий день.
// Если время пакета относится к новым суткам, storage предварительно
// очищается.
// Если пакет валидный, он добавляется в слайс storage, который возвращает
// функция. Если пакет невалидный, storage возвращается без изменений.
func DayActionInfo(data string, weight, height float64) string {
	stepsCount, walkDuration, err := parsePackage(data)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if stepsCount <= 0 {
		return ""
	}
	distant := (float64(stepsCount) * StepLength) / 1000.0
	calories := spentcalories.WalkingSpentCalories(stepsCount, weight, height, walkDuration)
	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.",
		stepsCount,
		distant,
		calories)
}
