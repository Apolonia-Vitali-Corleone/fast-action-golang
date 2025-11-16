package main

import "fmt"

func main() {
	arr := []int{99, 5, 2, 4, 6, 1, 3, 85, 66, 37}
	quickSort(arr, 0, len(arr)-1)
	fmt.Println(arr) // [1 2 3 4 5 6]
}

func quickSort(arr []int, left, right int) {
	if left < right {
		pivotIndex := partition(arr, left, right)

		quickSort(arr, left, pivotIndex-1)
		quickSort(arr, pivotIndex+1, right)
	}
}

func partition(arr []int, left, right int) int {
	pivot := arr[right]
	i := left - 1
	for j := left; j < right; j++ {
		if arr[j] < pivot {
			i++
			arr[i], arr[j] = arr[j], arr[i]
		}
	}
	arr[i+1], arr[right] = arr[right], arr[i+1]
	return i + 1
}
