package progressbar

import (
	"fmt"
	"strings"
	"sync"

	"github.com/KonstantinZy/go-console-progressbar/color"
)

var RENDERWIDTH = 20

type ProgresBarPool struct {
	Pool    map[int]*ProgresBar
	maxMsgs int
	writeCh chan string
	mesages []string
	Wg      *sync.WaitGroup
}

func (pbp *ProgresBarPool) AddProgresBar(pb *ProgresBar) {
	pbNum := len(pbp.Pool) + 1
	pb.setPoolNumber(pbNum)
	pb.pool = pbp
	pbp.Pool[pbNum] = pb
	fmt.Println() // spaces for render progresbars
}

func (pbp *ProgresBarPool) Write(msg string) {
	pbp.mesages = append(pbp.mesages, msg)
	if len(pbp.mesages) > pbp.maxMsgs {
		pbp.mesages = pbp.mesages[len(pbp.mesages)-pbp.maxMsgs:]
	}

	curLine := len(pbp.Pool) + pbp.maxMsgs
	resMsg := fmt.Sprintf("\033[%dA", curLine)
	for _, msg := range pbp.mesages {
		resMsg += fmt.Sprintf("%s\n", msg)
		curLine--
	}
	resMsg += fmt.Sprintf("\033[%dB", curLine)

	pbp.writeCh <- resMsg
}

func (pbp *ProgresBarPool) End() {
	close(pbp.writeCh)
}

func NewProgresBarsPool() *ProgresBarPool {
	maxMsgs := 5
	wg := &sync.WaitGroup{}

	wg.Add(1)
	ch := make(chan string)
	go func(ch chan string, wg *sync.WaitGroup) {
		defer wg.Done()
		for msg := range ch {
			fmt.Print(msg)
		}
	}(ch, wg)

	// make space for messages
	fmt.Println(strings.Repeat("=", 50), " Messages ", strings.Repeat("=", 50))
	for i := 1; i <= maxMsgs; i++ {
		fmt.Println()
	}

	return &ProgresBarPool{
		Pool:    make(map[int]*ProgresBar),
		maxMsgs: maxMsgs,
		writeCh: ch,
		mesages: []string{},
		Wg:      wg,
	}
}

type ProgresBar struct {
	Max          int
	Min          int
	Val          int
	Step         int
	Title        string
	numberInPool int
	Complite     bool
	RenderWidth  int
	RenderColor  string
	pool         *ProgresBarPool
}

func (pb *ProgresBar) setPoolNumber(n int) {
	pb.numberInPool = n
}

func (pb *ProgresBar) SetVal(val int) bool {
	pb.Val = val
	if pb.Val >= pb.Max {
		pb.Complite = true
	}

	pb.render(pb.pool) // render brogres bar
	return true
}

func (pb *ProgresBar) render(pool *ProgresBarPool) {
	linePref := fmt.Sprintf("\033[%dA", pb.numberInPool)
	linePostf := fmt.Sprintf("\033[%dB", pb.numberInPool)
	title := pb.Title
	if pb.Complite {
		title += fmt.Sprintf("%s complite%s", color.Red, color.Reset)
	}

	renderW := RENDERWIDTH
	if pb.RenderWidth != 0 {
		renderW = pb.RenderWidth
	}

	rederC := pb.RenderColor
	if rederC == "" {
		rederC = color.Green
	}

	line := getProgresLine(pb.Min, pb.Max, pb.Val, renderW, rederC) + " " + title
	lineBack := fmt.Sprintf("\033[%dD", len(line)) // return cursor back to line begin

	pool.writeCh <- fmt.Sprintf("%s%s%s%s", linePref, line, lineBack, linePostf)
}

func getProgresLine(min int, max int, val int, width int, rColor string) string {
	cnt := int(float64(val-min) / (float64(max-min) / float64(width)))
	spaceCnt := width - cnt
	if spaceCnt < 0 {
		spaceCnt = 0
	}
	title := fmt.Sprintf("%d from [%d - %d]", val, min, max)
	return fmt.Sprintf("%s[ %s ] %-18s %s", rColor, strings.Repeat("#", cnt)+strings.Repeat(" ", spaceCnt), title, color.Reset)
}
