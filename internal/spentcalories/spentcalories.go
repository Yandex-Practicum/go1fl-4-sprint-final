package spentcalories

import (
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

func parseTraining(data string) (int, string, time.Duration, error) {
	// TODO: реализовать функцию
	splitData := strings.Split(data, ",")
	if len(splitData) == 3 {
		activity := splitData[1]
		steps, err := strconv.Atoi(splitData[0])
		if steps <= 0 {
			return 0, "", 0, fmt.Errorf("количество шагов должно быть больше 0 %v", err)
		}

		if err != nil {
			return 0, "", 0, fmt.Errorf("ошибка при преобразовании количества шагов: %v", err)
		}
		duration, err := time.ParseDuration(splitData[2])
		if duration <= 0 {
			return 0, "", 0, fmt.Errorf("время должно быть больше 0 %v", err)
		}
		if err != nil {
			return 0, "", 0, fmt.Errorf("ошибка при преобразовании времени %v", err)
		}
		return steps, activity, duration, nil

	}
	return 0, "", 0, fmt.Errorf("Данные должны содержать 3 элемента")
}

func distance(steps int, height float64) float64 {
	// TODO: реализовать функцию
	stepLength := height * stepLengthCoefficient
	wayInMetrs := float64(steps) * stepLength
	way := wayInMetrs / mInKm
	return way
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	// TODO: реализовать функцию
	if duration <= 0 {
		return 0
	} else if steps <= 0 {
		return 0
	} else if height <= 0 {
		return 0
	} else {
		way := distance(steps, height)
		durationInHours := duration.Hours()
		avarageSpeed := way / durationInHours
		return avarageSpeed
	}
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	// TODO: реализовать функцию
	steps, activity, duration, err := parseTraining(data)
	durationInHours := duration.Hours()
	if err != nil {
		log.Println("Ошибка:", err)
		return "", err
	}
	if activity != "Ходьба" && activity != "Бег" {
		return "", fmt.Errorf("неизвестный тип тренировки")
	} else {
		switch activity {
		case "Ходьба":
			spentCalories, err := WalkingSpentCalories(steps, weight, height, duration)
			if err != nil {
				return "", err
			}
			speed := meanSpeed(steps, height, duration)
			distance := distance(steps, height)
			return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", activity, durationInHours, distance, speed, spentCalories), nil
		case "Бег":
			spentCalories, err := RunningSpentCalories(steps, weight, height, duration)
			if err != nil {
				return "", err
			}
			speed := meanSpeed(steps, height, duration)
			distance := distance(steps, height)
			return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", activity, durationInHours, distance, speed, spentCalories), nil
		}
		return "", err
	}

}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// TODO: реализовать функцию
	durationInMinutes := duration.Minutes()
	if durationInMinutes <= 0 {
		return 0, fmt.Errorf("время должен быть больше 0")
	}
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть больше 0")
	}
	if height <= 0 {
		return 0, fmt.Errorf("рост должен быть больше 0")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("вес должно быть больше 0")
	}
	speed := meanSpeed(steps, height, duration)
	spentCalories := (weight * speed * durationInMinutes) / minInH
	return spentCalories, nil

}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// TODO: реализовать функцию

	durationInMinutes := duration.Minutes()
	if durationInMinutes <= 0 {
		return 0, fmt.Errorf("время должен быть больше 0")
	}
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть больше 0")
	}
	if height <= 0 {
		return 0, fmt.Errorf("рост должен быть больше 0")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("вес должно быть больше 0")
	}
	Speed := meanSpeed(steps, height, duration)
	spentCalories := ((weight * Speed * durationInMinutes) / minInH) * walkingCaloriesCoefficient
	return spentCalories, nil

}
