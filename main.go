package main

import (
	"fmt"
	"math"
)

func main() {
	n := 1000
	tol := 1e-4

	// solve Laplace equation in parallel
	fmt.Println("Solving Laplace equation...")
	uLaplace := solveEquation(n, tol, "Laplace")
	fmt.Println("Laplace solution:", uLaplace)

	// solve heat conduction equation in parallel
	fmt.Println("Solving heat conduction equation...")
	uHeat := solveEquation(n, tol, "Heat")
	fmt.Println("Heat conduction solution:", uHeat)

	// solve wave propagation equation in parallel
	fmt.Println("Solving wave propagation equation...")
	uWave := solveEquation(n, tol, "Wave")
	fmt.Println("Wave propagation solution:", uWave)
}

func solveEquation(n int, tol float64, eqnType string) [][]float64 {
	// initialize the solution arrays
	u := make([][]float64, n+2)
	uNew := make([][]float64, n+2)
	for i := 0; i < n+2; i++ {
		u[i] = make([]float64, n+2)
		uNew[i] = make([]float64, n+2)
	}

	// set the boundary conditions
	for i := 0; i < n+2; i++ {
		u[i][0] = 1.0
		u[i][n+1] = 1.0
		u[0][i] = 0.0
		u[n+1][i] = 0.0
	}

	// set the initial guess for the solution
	for i := 1; i <= n; i++ {
		for j := 1; j <= n; j++ {
			u[i][j] = 0.0
		}
	}

	// initialize the done channel
	done := make(chan bool, n)

	// define the update function
	update := func(start, end int, u, uNew [][]float64, done chan<- bool) {
		for i := start; i <= end; i++ {
			for j := 1; j <= n; j++ {
				switch eqnType {
				case "Laplace":
					uNew[i][j] = 0.25 * (u[i-1][j] + u[i+1][j] + u[i][j-1] + u[i][j+1])
				case "Heat":
					uNew[i][j] = 0.25 * (u[i-1][j] + u[i+1][j] + u[i][j-1] + u[i][j+1])
					uNew[i][j] += 0.25 * (u[i][j] - 4.0*u[i][j] + u[i][j])
				case "Wave":
					uNew[i][j] = 2.0*u[i][j] - uNew[i][j] + 0.25*(u[i-1][j]+u[i+1][j]+u[i][j-1]+u[i][j]+0.125*(u[i-1][j-1]+u[i+1][j-1]+u[i-1][j+1]+u[i+1][j+1])-u[i][j])
				}
			}
			done <- true
		}
	}

	// start the update loop
	iter := 0
	maxIter := 100
	for iter < maxIter {
		// update the interior points in parallel
		for i := 1; i <= n; i++ {
			start := (i-1)*n/n + 1
			end := i * n / n
			go update(start, end, u, uNew, done)
		}
		// wait for all the updates to finish
		for i := 0; i < n; i++ {
			<-done
		}
		// check the convergence
		diff := 0.0
		for i := 1; i <= n; i++ {
			for j := 1; j <= n; j++ {
				diff += math.Abs(uNew[i][j] - u[i][j])
			}
		}
		if diff/float64(n*n) < tol {
			break
		}
		// update the solution
		for i := 1; i <= n; i++ {
			for j := 1; j <= n; j++ {
				u[i][j] = uNew[i][j]
			}
		}
		iter++
	}
	return u
}
