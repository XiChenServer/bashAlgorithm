package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// findCurrentMEX finds the current MEX from a sorted list of present elements
func findCurrentMEX(present []int, n int) int {
	for i := 0; i <= n; i++ {
		if i >= len(present) || present[i] != i {
			return i
		}
	}
	return n + 1
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscanf(reader, "%d\n", &T)

	for t := 0; t < T; t++ {
		var n int
		var k, x int64
		fmt.Fscanf(reader, "%d %d %d\n", &n, &k, &x)

		a := make([]int, n)
		set := make(map[int]struct{})
		for i := 0; i < n; i++ {
			fmt.Fscanf(reader, "%d", &a[i])
			set[a[i]] = struct{}{}
		}
		fmt.Fscanf(reader, "\n") // Handle the newline after the last number

		// Initialize present elements and calculate initial MEX
		present := make([]int, 0, n)
		for key := range set {
			present = append(present, key)
		}
		sort.Ints(present)
		currentMex := findCurrentMEX(present, n)

		// Initialize minimum cost to clear the entire array
		minCost := k * int64(currentMex)

		// Calculate cost by removing elements one by one
		cost := int64(0)
		for i := 0; i < n; i++ {
			cost += x // Accumulate cost of deleting individual elements
			delete(set, a[i])
			// Update the present list and MEX
			present = removeElement(present, a[i])
			currentMex = findCurrentMEX(present, n)
			totalCost := cost + k*int64(currentMex) // Total cost after removing i elements
			if totalCost < minCost {
				minCost = totalCost
			}
		}

		fmt.Fprintf(writer, "%d\n", minCost)
	}
}

// removeElement removes an element from a sorted list
func removeElement(present []int, value int) []int {
	i := sort.SearchInts(present, value)
	if i < len(present) && present[i] == value {
		return append(present[:i], present[i+1:]...)
	}
	return present
}
