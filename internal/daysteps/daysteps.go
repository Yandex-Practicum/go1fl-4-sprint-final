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
	// Разделяем строку по запятой
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, errors.New("неверный формат входных данных")
	}

	// Преобразуем первую часть в количество шагов
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка преобразования количества шагов: %w", err)
	}
	if steps <= 0 {
		return 0, 0, errors.New("неверное количество шагов!")
	}

	// Преобразуем вторую часть в продолжительность тренировки (например, "45m" или "1h30m")
	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка преобразования продолжительности: %w", err)
	}
	if duration <= 0 {
		return 0, 0, errors.New("неверная продолжительность тренировки!")
	}

	return steps, duration, nil

}

func DayActionInfo(data string, weight, height float64) string {

	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println("Ошибка конвертации данных:", err)
		return ""
	}

	if weight <= 0 || height <= 0 {
		log.Println("Вес и рост должны быть больше 0")
		return ""
	}
	// convLength переводит шаги в Км
	convLength := float64(steps) * stepLength / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println("ошибка подсчета калорий:", err)
		return ""
	}
	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.", steps, convLength, calories)
}
