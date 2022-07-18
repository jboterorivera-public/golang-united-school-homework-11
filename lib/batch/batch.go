package batch

import (
	"sync"
	"time"
)

type user struct {
	ID int64
}

func getOne(id int64) user {
	time.Sleep(time.Millisecond * 100)
	return user{ID: id}
}

func getBatch(n int64, pool int64) (res []user) {
	users := []user{}
	jobs := make(chan int64, n)
	results := make(chan user, n)

	go createJobs(jobs, n)
	go createWorkerPool(jobs, results, pool)

	for r := range results {
		users = append(users, r)
	}

	return users
}

func createJobs(ch chan int64, totalJobs int64) {
	var i int64
	for i = 0; i < totalJobs; i++ {
		ch <- i
	}
	close(ch)
}

func createWorkerPool(jobs chan int64, results chan user, workers int64) {
	var wg sync.WaitGroup

	var i int64
	for i = 0; i < workers; i++ {
		wg.Add(1)
		go worker(jobs, results, &wg)
	}

	wg.Wait()
	close(results)
}

func worker(jobs chan int64, results chan user, wg *sync.WaitGroup) {
	for job := range jobs {
		results <- getOne(job)
	}
	wg.Done()
}
