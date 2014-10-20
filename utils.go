package main

func sum(arr []int64) int64 {
	var sum int64

	for num := range arr {
		sum = sum + arr[num]
	}

	return sum
}
