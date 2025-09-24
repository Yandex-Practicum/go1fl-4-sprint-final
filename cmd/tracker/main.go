package main

import (
	"fmt"

	"github.com/Yandex-Practicum/tracker/internal/daysteps"
	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

func main() {
	fmt.Println("== Пример дневной активности ==")
	info := daysteps.DayActionInfo("1000,30m", 70, 1.75)
	fmt.Println(info)

	fmt.Println("\n== Пример тренировки ==")
	training, err := spentcalories.TrainingInfo("3000,Бег,45m", 70, 1.75)
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println(training)
	}
}
