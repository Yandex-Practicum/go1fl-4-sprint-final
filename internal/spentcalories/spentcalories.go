package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

// Функция принимает строку с данными формата "3456,Ходьба,3h00m",
//
//	которая содержит количество шагов, вид активности и продолжительность активности.
//	Функция возвращает четыре значения:
//		int — количество шагов.
//		string — вид активности.
//		time.Duration — продолжительность активности.
//		error — ошибку, если что-то пошло не так.
func parseTraining(data string) (int, string, time.Duration, error) {
	// TODO: реализовать функцию
	sliceData := strings.Split(data, ",")
	if len(sliceData) != 3 {
		return 0, "", 0, errors.New("не корректный состав строки: <>3")
	}

	numberOfSteps, err := strconv.Atoi(sliceData[0])
	if err != nil {
		return 0, "", 0, err
	} else if numberOfSteps <= 0 {
		err := errors.New("количество шагов <= 0")
		return 0, "", 0, err
	}

	duration, err := time.ParseDuration(sliceData[2])
	if err != nil {
		return 0, "", 0, err
	} else if duration <= 0 {
		err := errors.New("заданное время <= 0")
		return 0, "", 0, err
	}

	return numberOfSteps, sliceData[1], duration, nil
}

// Функция принимает количество шагов и рост пользователя в метрах,
//
//	а возвращает дистанцию в километрах(float64).
func distance(steps int, height float64) float64 {
	// TODO: реализовать функцию
	return height * stepLengthCoefficient * float64(steps) / float64(mInKm)
}

// Функция принимает количество шагов steps,
//
//	рост пользователя height и продолжительность активности duration
//	и возвращает среднюю скорость.
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	// TODO: реализовать функцию
	if duration.Seconds() <= 0 {
		return 0
	}
	return distance(steps, height) / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	// TODO: реализовать функцию
	var (
		typeActiv        string
		numberOfSteps    int
		duration         time.Duration
		distanceActive   float64
		speedActive      float64
		numberOfCalories float64
	)

	numberOfSteps, typeActiv, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}
	switch typeActiv {
	case "Ходьба":
		numberOfCalories, err = WalkingSpentCalories(numberOfSteps, weight, height, duration)
		distanceActive = distance(numberOfSteps, height)
		speedActive = meanSpeed(numberOfSteps, height, duration)
	case "Бег":
		numberOfCalories, err = RunningSpentCalories(numberOfSteps, weight, height, duration)
		distanceActive = distance(numberOfSteps, height)
		speedActive = meanSpeed(numberOfSteps, height, duration)
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	message := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\n"+
		"Дистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		typeActiv, duration.Hours(), distanceActive, speedActive, numberOfCalories)
	return message, nil

}

// Функция рассчитывает количество калорий по входным данным:
//
//	кол-во шагов, вес, рост, время активности
//
// Возвращает: кол-во калорий (float64) при беге, ошибку
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// TODO: реализовать функцию
	if duration <= 0 {
		return 0, errors.New("RunningSpentCalories: не корректное время")
	} else if steps <= 0 || height <= 0 || weight <= 0 {
		return 0, errors.New("RunningSpentCalories: шаги, вес, рост <= 0")
	}

	averageSpeed := meanSpeed(steps, height, duration)
	return weight * averageSpeed * duration.Minutes() / float64(minInH), nil
}

// Функция рассчитывает количество калорий по входным данным:
//
//	кол-во шагов, вес, рост, время активности
//
// Возвращает: кол-во калорий (float64) при ходьбе, ошибку
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// TODO: реализовать функцию
	if duration <= 0 {
		return 0, errors.New("WalkingSpentCalories: не корректное время")
	} else if steps <= 0 || height <= 0 || weight <= 0 {
		return 0, errors.New("WalkingSpentCalories: шаги, вес, рост <= 0")
	}

	averageSpeed := meanSpeed(steps, height, duration)
	numberCal := weight * averageSpeed * duration.Minutes() / float64(minInH)
	return numberCal * walkingCaloriesCoefficient, nil
}
