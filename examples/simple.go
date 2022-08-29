package main

import (
	"fmt"
	"math/rand"
	"time"

	progressbar "github.com/KonstantinZy/go-console-progressbar"
	"github.com/KonstantinZy/go-console-progressbar/color"
)

func main() {
	fmt.Println("... simple progresbar from 13 to 110 with step 7")
	rand.Seed(time.Now().UnixNano())

	max, min := 110, 13
	pb := &progressbar.ProgresBar{Max: max, Val: min, Min: min, Step: 7, Title: fmt.Sprintf("%sSimple progres from 13 to 110%s", color.Yellow, color.Reset), Complite: false}
	pbp := progressbar.NewProgresBarsPool()
	pbp.AddProgresBar(pb)

	for i := pb.Val; i < pb.Max; i += pb.Step {
		pb.SetVal(i)
		pbp.Write(fmt.Sprintf("Message from loop %d", rand.Intn(100)))
		time.Sleep(time.Second)
	}
	pb.SetVal(pb.Max)
	pbp.End()
	pbp.Wg.Wait()
	fmt.Println("... All progresbars complited")
}
