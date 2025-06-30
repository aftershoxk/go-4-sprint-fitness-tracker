package daysteps

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

var ErrIncorrectFormat = errors.New("некорректный формат данных")
var ErrZeroSteps = errors.New("количество шагов должно быть больше нуля")

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	splitData := strings.Split(data, ",")
	if len(splitData) != 2 {
		return 0, 0, fmt.Errorf("%w", ErrIncorrectFormat)
	}
	steps, err := strconv.Atoi(splitData[0])
	if err != nil {
		log.Println(err)
		return 0, 0, fmt.Errorf("ошибка: %w", err)
	}
	if steps <= 0 {
		return 0, 0, fmt.Errorf("%w", ErrZeroSteps)
	}
	walkingTime, err := time.ParseDuration(splitData[1])
	if err != nil {
		log.Println(err)
		return 0, 0, fmt.Errorf("ошибка: %w", err)
	}
	if walkingTime <= 0 {
		return 0, 0, fmt.Errorf("ошибка: %w", ErrIncorrectFormat)
	}
	return steps, walkingTime, nil
}

func DayActionInfo(data string, weight, height float64) string {
	// TODO: реализовать функцию
	steps, walkingTime, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}
	if steps <= 0 {
		log.Println(ErrZeroSteps)
		return ""
	}
	distance := float64(steps) * stepLength
	distance /= mInKm
	//Дописать количество калорий ниже
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, walkingTime)
	if err != nil {
		log.Println(err)
		return ""
	}
	Info := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps, distance, calories)
	return Info
}
