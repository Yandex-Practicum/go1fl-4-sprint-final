package daysteps

//пакет отвечает за учет активности в течении дня
//вывод кол-во шагов
//вывод пройденой дистанции
//вывод потраченных калорий

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/go1fl-4-sprint-final/internal/spentcalories"
)

var (
	StepLength = 0.65 // длина шага в метрах
)

// принимает кол-во шагов и время прогулки в формате 3h50m
// возвращает 3 значения
// int - кол-во шагов
// time.Duration - время прогулки
// error - ошибка если что то пошло не так
// функция парсит строку "678,0h50m" (678 - шаги; 0h50m - время прогулки)
func parsePackage(data string) (int, time.Duration, error) {

	dataslice := strings.Split(data, ",")
	if len(dataslice) != 2 {
		return 0, 0, errors.New("Длина меньше 2")
	}

	step, err := strconv.Atoi(dataslice[0])
	if err != nil {
		return 0, 0, fmt.Errorf("Ошибка при преобразовании в тип int")
	}

	durationTime := dataslice[1]

	parsTimeDuration, err := time.ParseDuration(durationTime)
	if err != nil {
		return 0, 0, fmt.Errorf("Ошибка при преобразовании строки в time.Duration")
	}

	return step, parsTimeDuration, nil

}

// DayActionInfo обрабатывает входящий пакет, который передаётся в
// виде строки в параметре data. Параметр storage содержит пакеты за текущий день.
// Если время пакета относится к новым суткам, storage предварительно
// очищается.
// Если пакет валидный, он добавляется в слайс storage, который возвращает
// функция. Если пакет невалидный, storage возвращается без изменений.

// data - кол-во шагов и время прогулки формат 3h50m
// weight - вес юзера в кг
// height - рост юзера в метрах
func DayActionInfo(data string, weight, height float64) string {
	step, duration, err := parsePackage(data)
	if err != nil {
		return ""
	}
	if step <= 0 {
		return ""
	}

	distance := (float64(step) * (StepLength)) / 1000

	calories := spentcalories.WalkingSpentCalories(step, weight, height, duration)

	return fmt.Sprintf("Количество шагов: %d\nДистанция составила %.2f\nВы сожгли %.2f ккал.", step, distance, calories)

}
