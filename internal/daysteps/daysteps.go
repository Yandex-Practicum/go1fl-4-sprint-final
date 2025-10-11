package daysteps

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	//"../spentcalories"

	//"fitness/go1fl-4-sprint-final/internal/spentcalories"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

// parsePackage()
// Функция принимает строку с данными, которая содержит количество шагов и продолжительность прогулки в формате 3h50m(3 часа 50 минут). Возвращает три значения:
// int — количество шагов
// time.Duration — продолжительность прогулки.
// error — ошибку, если что-то пошло не так.
// Функция парсит строку формата "678,0h50m", где:
// 678 — количество шагов.
// 0h50m — продолжительность прогулки.
// Переводит эти данные в int и time.Duration соответственно и возвращает их.
func parsePackage(data string) (int, time.Duration, error) {

	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("Incorrect data format ")
	}

	// Парсим количество шагов
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("Incorrect steps: %w", err)
	}
	//Проверить: количество шагов должно быть больше 0. Если это не так, вернуть нули и ошибку.
	if steps <= 0 {
		return 0, 0, fmt.Errorf("Steps cannot be less than zero")
	}
	// Парсим длительность
	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("Incorrect duration: %w", err)
	}
	if duration <= 0 {
		return 0, 0, fmt.Errorf("Duration cannot be less than zero")
	}
	return steps, duration, nil
}

// DayActionInfo
// Функция должна парсить строку с данными с помощью parsePackage(),
// вычислять дистанцию в километрах и количество потраченных калорий и
// возвращать строку в таком виде:
// Количество шагов: 792.
// Дистанция составила 0.51 км.
// Вы сожгли 221.33 ккал.
func DayActionInfo(data string, weight, height float64) string {

	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Printf("Ошибка получения числа шагов и продолжительности %v\n", err)
		return ""
	}
	if steps <= 0 {
		return ""
	}
	distance := (stepLength * float64(steps)) / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distance, calories)

}
