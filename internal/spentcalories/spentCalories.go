package spentcalories

//рассчет калорий по виду активности

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep   = 0.65  // средняя длина шага.
	mInKm     = 1000  // количество метров в километре.
	minInH    = 60    // количество минут в часе.
	kmhInMsec = 0.278 // коэффициент для преобразования км/ч в м/с.
	cmInM     = 100   // количество сантиметров в метре.
)

// принимает строку формата "3456,Ходьба,3h00m"
// строка сожержит кол-во шагов, вид активности, время активности
// int - колво шагов
// string - вид активности
// time.Duration - время активности
// error - ошибочка
func parseTraining(data string) (int, string, time.Duration, error) {

	dataslice := strings.Split(data, ",")
	if len(dataslice) != 3 {
		return 0, "", 0, errors.New("Длина меньше 3")
	}

	step, err := strconv.Atoi(dataslice[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("Ошибка в преобразовании в тип int", err)
	}

	activity := dataslice[1]

	durationTime, err := time.ParseDuration(dataslice[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("Ошибка в преобразовани в time.Duration", err)
	}

	return step, activity, durationTime, nil

}

//принимает кол-во шагов и возвразает дистанцию в км
//формула дистанции: (шаги * lenStep)/minKm

func distance(steps int) float64 {
	calcDistance := (float64(steps) * lenStep) / mInKm
	return calcDistance

}

// принимает кол-во шагов и время активности
// возвращает среднюю скорость
func meanSpeed(steps int, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	calcDistance := distance(steps)

	time := duration.Hours()

	calcSpeed := calcDistance / time

	return calcSpeed

}

// data string - строка формата "3456,Ходьба,3h00m"
// weight - вес (кг)
// height - рост (м)
func TrainingInfo(data string, weight, height float64) string {
	// Получаем данные с помощью функции parseTraining
	step, activity, duration, err := parseTraining(data)
	if err != nil {
		return fmt.Sprintf("Ошибка при обработке данных: %v", err)
	}

	// Проверка на количество шагов
	if step <= 0 {
		return "Ошибка: количество шагов должно быть больше 0."
	}

	// Рассчитываем дистанцию и скорость
	dist := distance(step)
	speed := meanSpeed(step, duration)

	// Используем switch для обработки разных типов активности
	switch activity {
	case "Ходьба":
		// Рассчитываем калории для ходьбы
		// В данном примере я предполагаю, что для ходьбы будет отдельная формула для калорий
		calories := 0.04 * weight * dist // Примерная формула для ходьбы
		return fmt.Sprintf("Количество шагов: %d\nВид активности: Ходьба\nДистанция составила %.2f км\nСредняя скорость: %.2f км/ч\nВы сожгли %.2f ккал.", step, dist, speed, calories)
	case "Бег":
		// Рассчитываем калории для бега
		calories := RunningSpentCalories(step, weight, duration)
		return fmt.Sprintf("Количество шагов: %d\nВид активности: Бег\nДистанция составила %.2f км\nСредняя скорость: %.2f км/ч\nВы сожгли %.2f ккал.", step, dist, speed, calories)
	default:
		// Если тип активности неизвестен
		return "Неизвестный тип тренировки."
	}
}

// Константы для расчета калорий, расходуемых при беге.
const (
	runningCaloriesMeanSpeedMultiplier = 18.0 // множитель средней скорости.
	runningCaloriesMeanSpeedShift      = 20.0 // среднее количество сжигаемых калорий при беге.
)

// steps - кол-во шагов
// weight - вес юзера
// duration time.Duration - время бега
// возвращает кол-во каллорий после бега
func RunningSpentCalories(steps int, weight float64, duration time.Duration) float64 {
	calcSpeed := meanSpeed(steps, duration)

	calcKallories := ((runningCaloriesMeanSpeedMultiplier * calcSpeed) - runningCaloriesMeanSpeedShift) * weight

	return calcKallories
}

// Константы для расчета калорий, расходуемых при ходьбе.
const (
	walkingCaloriesWeightMultiplier = 0.035 // множитель массы тела.
	walkingSpeedHeightMultiplier    = 0.029 // множитель роста.
)

// steps - кол-во шагов
// weight - вес (кг)
// height - рост (м)
// duration time.Duration - время ходьбы
// возвращает кол-во калорий во время хлдьбы
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) float64 {
	calcSpeed := meanSpeed(steps, duration)
	time := duration.Hours()

	calcKallories := ((walkingCaloriesWeightMultiplier * weight) + (calcSpeed*calcSpeed/height)*walkingSpeedHeightMultiplier) * time * minInH

	return calcKallories
}
