package main

import (
	"fmt"
	"silkroad/m/internal/domain/tour"
)

func main() {
	// Тестируем базовый фильтр без параметров
	filter := tour.TourFilter{
		Limit:  4,
		Offset: 0,
	}

	fmt.Printf("Filter: %+v\n", filter)
	fmt.Printf("Quantity: %v\n", filter.Quantity)
	fmt.Printf("Len(Quantity): %d\n", len(filter.Quantity))
}
