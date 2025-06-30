package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var ErrIncorrectFormat = errors.New("некорректный формат данных")
var ErrUnknownTraining = errors.New("неизвестный тип тренировки")

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	splitData := strings.Split(data, ",")
	if len(splitData) != 3 {
		return 0, "", 0, fmt.Errorf("ошибка: %w", ErrIncorrectFormat)
	}

	steps, err := strconv.Atoi(splitData[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка: %w", err)
	}
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("ошибка: %w", ErrIncorrectFormat)
	}

	activity := splitData[1]

	durationOfActivity, err := time.ParseDuration(splitData[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка: %w", err)
	}
	if durationOfActivity <= 0 {
		return 0, "", 0, fmt.Errorf("ошибка: %w", ErrIncorrectFormat)
	}

	return steps, activity, durationOfActivity, nil
}

func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	temporaryValue := float64(steps) * stepLength
	distanceInKm := temporaryValue / mInKm
	return distanceInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0.0
	}
	distanceInKm := distance(steps, height)
	averageSpeed := distanceInKm / duration.Hours()
	return averageSpeed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, durationOfActivity, err := parseTraining(data)
	if err != nil {
		return "", fmt.Errorf("ошибка %w", err)
	}
	distanceInKm := distance(steps, height)
	averageSpeed := meanSpeed(steps, height, durationOfActivity)
	hours := durationOfActivity.Hours()
	Info := ""
	switch {
	case activity == "Бег":
		runCalories, err := RunningSpentCalories(steps, weight, height, durationOfActivity)
		if err != nil {
			return "", fmt.Errorf("ошибка: %w", err)
		}
		Info = fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
			activity, hours, distanceInKm, averageSpeed, runCalories)
	case activity == "Ходьба":
		walkCalories, err := WalkingSpentCalories(steps, weight, height, durationOfActivity)
		if err != nil {
			return "", fmt.Errorf("ошибка: %w", err)
		}
		Info = fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
			activity, hours, distanceInKm, averageSpeed, walkCalories)
	default:
		return "", fmt.Errorf("%w", ErrUnknownTraining)
	}
	return Info, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0.0 || height <= 0.0 || duration <= 0 {
		return 0.0, fmt.Errorf("ошибка %w", ErrIncorrectFormat)
	}
	averageSpeed := meanSpeed(steps, height, duration)
	durationMinutes := duration.Minutes()
	calories := (weight * averageSpeed * durationMinutes) / minInH
	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0.0 || height <= 0.0 || duration <= 0 {
		return 0.0, fmt.Errorf("ошибка %w", ErrIncorrectFormat)
	}
	averageSpeed := meanSpeed(steps, height, duration)
	durationMinutes := duration.Minutes()
	calories := (weight * averageSpeed * durationMinutes) / minInH
	calories *= walkingCaloriesCoefficient
	return calories, nil
}
