package daysteps

import (
	"fmt"
	"go1fl-4-sprint-final/internal/spentcalories"
	"strconv"
	"strings"
	"time"
)

const unit = 1000

var (
	StepLength = 0.65 // длина шага в метрах
)

// Функция parsePackage принимает строку с данными, которая содержит количество шагов
// и продолжительность прогулки в формате 3h50m(3 часа 50 минут). Возвращает три значения:
//
// парсит строку формата "678,0h50m", где:
// 678 — количество шагов.
// 0h50m — продолжительность прогулки.
func parsePackage(data string) (int, time.Duration, error) {

	strData := strings.Split(data, ",")
	if len(strData) != 2 {
		return 0, 0, fmt.Errorf("[parsePackage(data string)] неверный формат строки: %s", data)
	}
	step, err := strconv.Atoi(strData[0])
	if err != nil {
		return 0, 0, fmt.Errorf("[parsePackage(data string)] неверный формат количества "+
			"шагов: %s: %w", strData[0], err)
	}
	if step < 1 {
		return 0, 0, fmt.Errorf("[parsePackage(data string)] 0 шагов")
	}

	period, err := time.ParseDuration(strData[1])
	if err != nil {
		return 0, 0, fmt.Errorf("[parsePackage(data string)] неверный формат времени")
	}

	return step, period, nil

}

// DayActionInfo обрабатывает входящий пакет, который передаётся в
// виде строки в параметре data. Параметр storage содержит пакеты за текущий день.
// Если время пакета относится к новым суткам, storage предварительно
// очищается.
// Если пакет валидный, он добавляется в слайс storage, который возвращает
// функция. Если пакет невалидный, storage возвращается без изменений.
func DayActionInfo(data string, weight, height float64) string {
	// ваш код ниже
	step, period, err := parsePackage(data)
	if err != nil {
		fmt.Println("[DayActionInfo(string, float 64, float64)]", err)
		return ""
	}
	if step < 1 {
		fmt.Println("[DayActionInfo(string, float 64, float64)] количество шагов 0")
		return ""
	}
	distance := float64(StepLength) * float64(step) / unit
	kalories := spentcalories.WalkingSpentCalories(step, weight, height, period)
	strInfo := fmt.Sprintf("Количество шагов %d.\n Дистанция составила %.2fкм\n Вы сожгли %.2f ккал.", step, distance, kalories)
	return strInfo
}
