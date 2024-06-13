package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const MEMORY_SIZE = 5

func main() {
	flag.Parse()
	var src *os.File
	if flag.NArg() == 1 {
		fp, err := os.Open(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}
		src = fp
	} else {
		log.Fatal("args required")
	}
	defer src.Close()

	var buf bytes.Buffer
	_, err := io.Copy(&buf, src)
	if err != nil {
		log.Fatal(err)
	}

	targetStr := buf.String()
	targetStr = strings.ReplaceAll(targetStr, "\n", "")
	targetStr = strings.ReplaceAll(targetStr, " ", "")

	interpreter := NewInterPreter()

	if err := Execute(targetStr, interpreter); err != nil {
		log.Fatal(err)
	}

	fmt.Println("")
}

func Execute(targetStr string, ip *Interpreter) error {
	i := 0
	for i < len(targetStr) {
		r := targetStr[i]
		switch r {
		case '>':
			ip.Right()
		case '<':
			ip.Left()
		case '+':
			ip.Plus()
		case '-':
			ip.Minus()
		case '.':
			ip.OutPut()
		case ',':
			// TODO
			break
		case '[':
			ip.SaveLoopStartIndex(i - 1)
			endIndex, err := ip.PopLoopEndIndex()
			if err != nil {
				break
			}
			if ip.IsZero() {
				i = endIndex
			}
		case ']':
			ip.SaveLoopEndIndex(i)
			startIndex, err := ip.PopLoopStartIndex()
			if err != nil {
				break
			}
			if !ip.IsZero() {
				i = startIndex
			}
		default:
			return fmt.Errorf("invalid code = %c[%v], index = %d", r, r, i)
		}
		i++
		// ip.Log()
	}
	return nil
}

type Interpreter struct {
	ptr          int
	memory       []byte
	loopStartPtr []int
	loopEndPtr   []int
}

func NewInterPreter() *Interpreter {
	return &Interpreter{
		ptr:    0,
		memory: make([]byte, MEMORY_SIZE),
	}
}

func (i *Interpreter) Plus() {
	i.memory[i.ptr]++
}

func (i *Interpreter) Minus() {
	i.memory[i.ptr]--
}

func (i *Interpreter) Right() {
	i.ptr++
}

func (i *Interpreter) Left() {
	i.ptr--
}

func (i *Interpreter) OutPut() {
	fmt.Printf("%c", i.memory[i.ptr])
}

func (i *Interpreter) Log() {
	for _, v := range i.memory {
		fmt.Printf("\t%v", v)
	}
	fmt.Println("")
}

func (i *Interpreter) IsZero() bool {
	return i.memory[i.ptr] == 0
}

func (i *Interpreter) SaveLoopStartIndex(index int) {
	i.loopStartPtr = append(i.loopStartPtr, index)
}

func (i *Interpreter) PopLoopStartIndex() (int, error) {
	ptr := 0
	if len(i.loopStartPtr) == 0 {
		return 0, fmt.Errorf("Finish")
	}
	ptr = i.loopStartPtr[len(i.loopStartPtr)-1]
	i.loopStartPtr = i.loopStartPtr[:len(i.loopStartPtr)-1]
	return ptr, nil
}

func (i *Interpreter) SaveLoopEndIndex(index int) {
	i.loopEndPtr = append(i.loopEndPtr, index)
}

func (i *Interpreter) PopLoopEndIndex() (int, error) {
	ptr := 0
	if len(i.loopEndPtr) == 0 {
		return 0, fmt.Errorf("Finish")
	}
	ptr = i.loopEndPtr[len(i.loopEndPtr)-1]
	i.loopEndPtr = i.loopEndPtr[:len(i.loopEndPtr)-1]
	return ptr, nil
}
