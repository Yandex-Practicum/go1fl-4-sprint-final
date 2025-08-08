package main

import (
	"fmt"
	"log"

	"github.com/Sergey-81/Fitness-tracker-module/internal/daysteps"
	"github.com/Sergey-81/Fitness-tracker-module/internal/spentcalories"
)

func main() {
	weight := 84.6
	height := 1.87

	// Настройка вывода ошибок
	log.SetFlags(0) // Убираем дату/время в логах

	fmt.Println("Активность в течение дня")
	fmt.Println("=======================")

	processDailyActivity(weight, height)
	
	fmt.Println("\nЖурнал тренировок")
	fmt.Println("===============")

	processWorkouts(weight, height)
}

func processDailyActivity(weight, height float64) {
	input := []string{
		"678,0h50m", "792,1h14m", "1078,1h30m", "7830,2h40m",
		",3456", "12:40:00, 3456", "something is wrong",
	}

	for i, v := range input {
		if info := daysteps.DayActionInfo(v, weight, height); info != "" {
			fmt.Print(info) // fmt.Print вместо fmt.Println, так как DayActionInfo уже добавляет \n
		} else {
			fmt.Printf("⚠️ Ошибка обработки: '%s'\n", v)
		}
		if i < len(input)-1 {
			fmt.Println("-----------------------")
		}
	}
}

func processWorkouts(weight, height float64) {
	trainings := []string{
		"3456,Ходьба,3h00m", "something is wrong", "678,Бег,0h5m",
		"1078,Бег,0h10m", ",3456 Ходьба", "7892,Ходьба,3h10m", 
		"15392,Бег,0h45m",
	}

	for _, v := range trainings {
		if info, err := spentcalories.TrainingInfo(v, weight, height); err != nil {
			fmt.Printf("⚠️ Ошибка: %v\n", err)
		} else {
			fmt.Print(info) // fmt.Print вместо fmt.Println
		}
		fmt.Println("----------------")
	}
}
