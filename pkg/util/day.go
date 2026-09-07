package util

import "time"

func IsWorkingDay() bool {
	return true

	t := time.Now()
	switch t.Weekday() {
	case time.Saturday, time.Sunday:
		return false // Выходные дни исключаем
	default:
		return true // Рабочие дни обрабатываем
	}
}
