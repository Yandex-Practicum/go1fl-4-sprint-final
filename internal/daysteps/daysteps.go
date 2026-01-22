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

// Парсит строку с данными о шагах и продолжительности прогулки
func parsePackage(data string) (int, time.Duration, error) {

	// Разделяем строку по запятой
	parts := strings.Split(data, ",")

	if len(parts) != 2 {
		return 0, 0, errors.New("Неверный формат")
	}

	// Парсим количество шагов
	number, err := strconv.Atoi(parts[0])

	if err != nil {
		return 0, 0, err
	}
	if number <= 0 {
		return 0, 0, errors.New("Количество шагов должно быть больше 0")
	}

	// Парсим продолжительость
	duration, err := time.ParseDuration(parts[1])

	// Проверяем продолжительность
	if duration <= 0 {
		return 0, 0, errors.New("Продолжительность должна быть больше 0")
	}

	if err != nil {
		return 0, 0, err
	}
	return number, duration, nil
}

// DayActionInfo парсит данные о прогулке и возвращает информацию о шагах, дистанции и калориях.
func DayActionInfo(data string, weight, height float64) string {

	// Парсим данные о количестве шагов и продолжительсти прогулки
	steps, duration, err := parsePackage(data)

	if err != nil {
		log.Println(err)
		return ""
	}
	if steps <= 0 {
		return ""
	}

	// Рассчитываем дистанцию и калории
	distance := (stepLength * float64(steps)) / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distance, calories)
}
