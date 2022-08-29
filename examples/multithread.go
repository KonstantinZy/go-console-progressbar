package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	progressbar "github.com/KonstantinZy/go-console-progressbar"
	"github.com/KonstantinZy/go-console-progressbar/color"
)

func main() {
	fmt.Println("... multiple progresbars with number up to 10")

	rand.Seed(time.Hour.Milliseconds())

	if len(os.Args) == 1 {
		fmt.Println("Usage multithread <int>")
		os.Exit(0)
	}

	threadsNmb, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(200)
	}

	if threadsNmb > 10 {
		threadsNmb = 10
		fmt.Println("Too much threads - using 10")
	}

	pbp := progressbar.NewProgresBarsPool()
	step := 5
	wg := &sync.WaitGroup{}
	pb := &progressbar.ProgresBar{}

	colors := []string{
		color.Blue,
		color.Cyan,
		color.Purple,
	}

	// add progressbars to pool
	for i := 1; i <= threadsNmb; i++ {
		min := getRand(0, 10)
		max := getRand(min+10, 150)
		title := fmt.Sprintf("%sProgresbar from %d to %d with step %d%s", color.Yellow, min, max, step, color.Reset)
		pb = &progressbar.ProgresBar{Max: max, Val: min, Min: min, Step: step, Title: title, Complite: false, RenderColor: colors[rand.Intn(len(colors))]}
		pbp.AddProgresBar(pb)
	}

	// work gorutines
	wg.Add(len(pbp.Pool))
	for _, pb := range pbp.Pool {
		go func(pb *progressbar.ProgresBar, pbp *progressbar.ProgresBarPool, wg *sync.WaitGroup) {
			defer wg.Done()
			for i := pb.Val; i < pb.Max; i += pb.Step {
				pb.SetVal(i)
				time.Sleep(time.Second)
			}
			pb.SetVal(pb.Max) // finish progres
			pbp.Write(fmt.Sprintf("Fifnish worker of progres bar %d", rand.Intn(20)))
		}(pb, pbp, wg)
	}
	wg.Wait()

	// close pool channel and wait for last message
	pbp.End()
	pbp.Wg.Wait()

	fmt.Println("... All progresbars complited")
}

func getRand(min int, max int) int {
	return rand.Intn(max-min) + min
}
