package services

import (
	"container/heap"
	"math"
)

// ============ 1. SORTING ============

// Bubble Sort - O(n^2), sederhana tapi lambat
func BubbleSort(a []int) []int {
	r := append([]int(nil), a...)
	for i := 0; i < len(r)-1; i++ {
		for j := 0; j < len(r)-1-i; j++ {
			if r[j] > r[j+1] {
				r[j], r[j+1] = r[j+1], r[j]
			}
		}
	}
	return r
}

// Quick Sort - rata-rata O(n log n), divide and conquer dengan pivot
func QuickSort(a []int) []int {
	if len(a) <= 1 {
		return a
	}
	pivot := a[len(a)/2]
	var kiri, tengah, kanan []int
	for _, v := range a {
		switch {
		case v < pivot:
			kiri = append(kiri, v)
		case v == pivot:
			tengah = append(tengah, v)
		default:
			kanan = append(kanan, v)
		}
	}
	return append(append(QuickSort(kiri), tengah...), QuickSort(kanan)...)
}

// Merge Sort - selalu O(n log n), stabil
func MergeSort(a []int) []int {
	if len(a) <= 1 {
		return a
	}
	mid := len(a) / 2
	kiri, kanan := MergeSort(a[:mid]), MergeSort(a[mid:])
	hasil := make([]int, 0, len(a))
	i, j := 0, 0
	for i < len(kiri) && j < len(kanan) {
		if kiri[i] <= kanan[j] {
			hasil = append(hasil, kiri[i])
			i++
		} else {
			hasil = append(hasil, kanan[j])
			j++
		}
	}
	hasil = append(hasil, kiri[i:]...)
	return append(hasil, kanan[j:]...)
}

// ============ 2. SEARCHING ============

// Searching - Binary Search - O(log n), data harus sudah terurut
func BinarySearch(a []int, target int) int {
	lo, hi := 0, len(a)-1
	for lo <= hi {
		mid := lo + (hi-lo)/2
		switch {
		case a[mid] == target:
			return mid
		case a[mid] < target:
			lo = mid + 1
		default:
			hi = mid - 1
		}
	}
	return -1
}

// ============ 3. GRAPH ============

// Graph - BFS - menjelajah per level (jalur terpendek pada graf tak berbobot)
func GraphBfs(graph map[int][]int, start int) []int {
	visited := map[int]bool{start: true}
	queue := []int{start}
	var urutan []int
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		urutan = append(urutan, node)
		for _, tetangga := range graph[node] {
			if !visited[tetangga] {
				visited[tetangga] = true
				queue = append(queue, tetangga)
			}
		}
	}
	return urutan
}

// // DFS - menjelajah sedalam mungkin dulu (rekursif)
// func dfs(graph map[int][]int, node int, visited map[int]bool, urutan *[]int) {
// 	visited[node] = true
// 	*urutan = append(*urutan, node)
// 	for _, tetangga := range graph[node] {
// 		if !visited[tetangga] {
// 			dfs(graph, tetangga, visited, urutan)
// 		}
// 	}
// }
// GraphDfs - menjelajah graf sedalam mungkin dulu, mulai dari node start
func GraphDfs(graph map[int][]int, start int) []int {
	visited := map[int]bool{}
	var urutan []int
	dfsVisit(graph, start, visited, &urutan)
	return urutan
}

// dfsVisit - fungsi pembantu rekursif (huruf kecil: hanya dipakai di package ini)
func dfsVisit(graph map[int][]int, node int, visited map[int]bool, urutan *[]int) {
	visited[node] = true
	*urutan = append(*urutan, node)
	for _, tetangga := range graph[node] {
		if !visited[tetangga] {
			dfsVisit(graph, tetangga, visited, urutan)
		}
	}
}

// Dijkstra - jalur terpendek pada graf berbobot positif
type Edge struct{ to, w int }
type Item struct{ node, dist int }
type PQ []Item

func (p PQ) Len() int            { return len(p) }
func (p PQ) Less(i, j int) bool  { return p[i].dist < p[j].dist }
func (p PQ) Swap(i, j int)       { p[i], p[j] = p[j], p[i] }
func (p *PQ) Push(x interface{}) { *p = append(*p, x.(Item)) }
func (p *PQ) Pop() interface{} {
	old := *p
	n := len(old)
	item := old[n-1]
	*p = old[:n-1]
	return item
}

func Dijkstra(graph [][]Edge, src int) []int {
	dist := make([]int, len(graph))
	for i := range dist {
		dist[i] = math.MaxInt64
	}
	dist[src] = 0
	pq := &PQ{{src, 0}}
	for pq.Len() > 0 {
		cur := heap.Pop(pq).(Item)
		if cur.dist > dist[cur.node] {
			continue
		}
		for _, e := range graph[cur.node] {
			if nd := cur.dist + e.w; nd < dist[e.to] {
				dist[e.to] = nd
				heap.Push(pq, Item{e.to, nd})
			}
		}
	}
	return dist
}

// ============ 4. DYNAMIC PROGRAMMING ============

// Fibonacci dengan memoization - O(n)
func Fibonacci(n int, memo map[int]int) int {
	if n <= 1 {
		return n
	}
	if v, ok := memo[n]; ok {
		return v
	}
	memo[n] = Fibonacci(n-1, memo) + Fibonacci(n-2, memo)
	return memo[n]
}

// Knapsack 0/1 - nilai maksimum dengan kapasitas terbatas
func Knapsack(bobot, nilai []int, kapasitas int) int {
	dp := make([]int, kapasitas+1)
	for i := range bobot {
		for c := kapasitas; c >= bobot[i]; c-- {
			if v := dp[c-bobot[i]] + nilai[i]; v > dp[c] {
				dp[c] = v
			}
		}
	}
	return dp[kapasitas]
}

// ============ 5. MATEMATIKA ============

// GCD / FPB - Algoritma Euclid
func AlgoritmaGcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// Sieve of Eratosthenes - semua bilangan prima <= n
func AlgoritmaSieve(n int) []int {
	prima := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		prima[i] = true
	}
	for i := 2; i*i <= n; i++ {
		if prima[i] {
			for j := i * i; j <= n; j += i {
				prima[j] = false
			}
		}
	}
	var hasil []int
	for i, p := range prima {
		if p {
			hasil = append(hasil, i)
		}
	}
	return hasil
}
