package daysteps

import (
	"errors"
	//"errors"
	"fmt"
	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

// Функция принимает строку с данными, которая содержит количество шагов и продолжительность прогулки в формате 3h50m.
// Возвращает три значения:
//
//	int — количество шагов;
//	time.Duration — продолжительность прогулки;
//	error — ошибку, если что-то пошло не так.
func parsePackage(data string) (int, time.Duration, error) {
	// TODO: реализовать функцию

	sliceData := strings.Split(data, ",")
	if len(sliceData) != 2 {
		return 0, 0, errors.New("не корректный состав строки: <>2")
	}
	numberOfSteps, err := strconv.Atoi(sliceData[0])
	if err != nil {
		return 0, 0, err
	} else if numberOfSteps <= 0 {
		err := errors.New("количество шагов <= 0")
		return 0, 0, err
	}

	duration, err := time.ParseDuration(sliceData[1])
	if err != nil {
		return 0, 0, err
	} else if duration <= 0 {
		err := errors.New("заданное время <= 0")
		return 0, 0, err
	}
	return numberOfSteps, duration, nil
}

// Функция парсит строку с данными с помощью parsePackage(),
//
//	 вычисляет дистанцию в километрах и количество потраченных калорий и возвращать строку:
//		Количество шагов: 792.
//		Дистанция составила 0.51 км.
//		Вы сожгли 221.33 ккал.
func DayActionInfo(data string, weight, height float64) string {
	// TODO: реализовать функцию
	var distanceKm float64
	var duration time.Duration
	numberOfSteps, duration, err := parsePackage(data)
	if err != nil || numberOfSteps <= 0 {
		log.Printf("не корректный формат или кол-во шагов: %v", err)
		return ""
	}

	/*
		if numberOfSteps <= 0 {
			log.Printf("не корректный формат: %v", err)
			return ""
		}
	*/

	distanceKm = float64(numberOfSteps) * stepLength / float64(mInKm)

	numberCalStep, err := spentcalories.WalkingSpentCalories(numberOfSteps, weight, height, duration)
	if err != nil {
		log.Printf("не корректный формат: %v", err)
		return ""
	}
	message := fmt.Sprintf("Количество шагов: %d.\n"+
		"Дистанция составила %.2f км.\n"+
		"Вы сожгли %.2f ккал.\n", numberOfSteps, distanceKm, numberCalStep)
	return message
}
