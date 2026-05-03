package daysteps

import (
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
	parts := strings.Split(data, ",")

	// Проверка, чтобы длина слайса равна 2
	if len(parts) != 2 {
		err := fmt.Errorf("Ожидается: шаги, продолжительность. Получено: %s", data)
		log.Printf("Ошибка парсинга: %v", err)
		return 0, 0, err
	}

	// Получаем шаги
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		log.Printf("Ошибка преобразования шагов: %v", err)
		return 0, 0, fmt.Errorf("Невозможно преобразовать шаги: %v", err)
	}

	// Проверка количества шагов
	if steps <= 0 {
		err := fmt.Errorf("Количество шагов должно быть больше 0, получено: %d", steps)
		log.Printf("Ошибка валидации: %v", err)
		return 0, 0, err
	}

	// Получаем время шагов
	duration, err := time.ParseDuration(parts[1])

	if duration <= 0 {
		log.Print("Ошибка времени. Время меньше или равно нулю")
		return 0, 0, fmt.Errorf("Ошибка времени. Время меньше или равно нулю")
	}

	if err != nil {
		log.Printf("Ошибка преобразования времени: %v", err)
		return 0, 0, fmt.Errorf("Ошибка преобразования времени: %v", err)
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {

	// Данные шагов и продолжительности прогулки
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Printf("Ошибка при парсинге данных: %v", err)
		return ""
	}

	// Проверка, чтобы количество шагов больше 0
	if steps <= 0 {
		log.Printf("Ошибка: количество шагов должно быть больше 0, получено: %d", steps)
		return ""
	}

	// Расчет дистанции
	distance := (float64(steps) * stepLength) / mInKm

	// Расчет калорий
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Printf("Ошибка при расчете калорий: %v", err)
		return ""
	}

	// Формирование результата
	result := fmt.Sprintf(
		"Количество шагов: %d.\n"+
			"Дистанция составила %.2f км.\n"+
			"Вы сожгли %.2f ккал.\n",
		steps, distance, calories)

	return result
}
