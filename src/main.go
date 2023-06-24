package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	BrainfuckIncrementPtr  = '>'
	BrainfuckDecrementPtr  = '<'
	BrainfuckIncrementData = '+'
	BrainfuckDecrementData = '-'
	BrainfuckOutput        = '.'
	BrainfuckInput         = ','
	BrainfuckLoopStart     = '['
	BrainfuckLoopEnd       = ']'
)

type BrainfuckInterpreter struct {
	code         string
	data         []byte
	dataPointer  int
	loopStack    []int
	outputBuffer strings.Builder
}

func NewBrainfuckInterpreter(code string) *BrainfuckInterpreter {
	return &BrainfuckInterpreter{
		code:        code,
		data:        make([]byte, 30000),
		dataPointer: 0,
		loopStack:   make([]int, 0),
	}
}

func banner() {
	fmt.Printf(`
	
	
██████  ██████   █████  ██ ███    ██ ███████ ██    ██  ██████ ██   ██ 
██   ██ ██   ██ ██   ██ ██ ████   ██ ██      ██    ██ ██      ██  ██  
██████  ██████  ███████ ██ ██ ██  ██ █████   ██    ██ ██      █████   
██   ██ ██   ██ ██   ██ ██ ██  ██ ██ ██      ██    ██ ██      ██  ██  
██████  ██   ██ ██   ██ ██ ██   ████ ██       ██████   ██████ ██   ██ 


`)
}

func (bf *BrainfuckInterpreter) execute() {
	for i := 0; i < len(bf.code); i++ {
		switch bf.code[i] {
		case BrainfuckIncrementPtr:
			bf.dataPointer++
		case BrainfuckDecrementPtr:
			bf.dataPointer--
		case BrainfuckIncrementData:
			bf.data[bf.dataPointer]++
		case BrainfuckDecrementData:
			bf.data[bf.dataPointer]--
		case BrainfuckOutput:
			bf.outputBuffer.WriteByte(bf.data[bf.dataPointer])
		case BrainfuckInput:
		case BrainfuckLoopStart:
			if bf.data[bf.dataPointer] == 0 {
				loopCount := 1
				for loopCount > 0 {
					i++
					if bf.code[i] == BrainfuckLoopStart {
						loopCount++
					} else if bf.code[i] == BrainfuckLoopEnd {
						loopCount--
					}
				}
			} else {
				bf.loopStack = append(bf.loopStack, i)
			}
		case BrainfuckLoopEnd:
			if bf.data[bf.dataPointer] == 0 {
				bf.loopStack = bf.loopStack[:len(bf.loopStack)-1]
			} else {
				i = bf.loopStack[len(bf.loopStack)-1] - 1
			}
		}
	}
}

func (bf *BrainfuckInterpreter) getOutput() string {
	return bf.outputBuffer.String()
}

func main() {
	banner()
	helpFlag := flag.Bool("h", false, "Display help")
	fileFlag := flag.String("f", "", "Brainfuck code file")
	flag.Parse()

	if *helpFlag {
		flag.PrintDefaults()
		fmt.Println("Example usage:")
		fmt.Println("go run main.go -f /path/to/save/location/file.bf")
		return
	}

	if *fileFlag == "" {
		fmt.Print("Enter the location of the Brainfuck code file: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		*fileFlag = scanner.Text()
		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			return
		}
	}

	file, err := os.Open(*fileFlag)
	if err != nil {
		fmt.Printf("Error opening Brainfuck code file: %v\n", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var code strings.Builder
	for scanner.Scan() {
		code.WriteString(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading Brainfuck code file: %v\n", err)
		return
	}

	interpreter := NewBrainfuckInterpreter(code.String())
	interpreter.execute()

	outputFile, err := os.Create("output.txt")
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
	defer outputFile.Close()

	outputWriter := bufio.NewWriter(outputFile)
	_, err = outputWriter.WriteString(interpreter.getOutput())
	if err != nil {
		fmt.Printf("Error writing to output file: %v\n", err)
		return
	}

	err = outputWriter.Flush()
	if err != nil {
		fmt.Printf("Error flushing output file: %v\n", err)
		return
	}

	fmt.Println("Output saved successfully.")
}
