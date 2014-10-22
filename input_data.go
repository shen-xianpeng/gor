package main

import (
	"os"
	"bufio"
    "fmt"
)


type DataInput struct {
	data	chan []byte
	path	string
}


func NewDataInput(path string) (i *DataInput) {
	i = new(DataInput)
	i.data = make(chan []byte)
	i.path = path
	i.Init(path)

	go i.emit()

	return
}

func (i *DataInput) Init(path string) {


}

func (i *DataInput) Read(data []byte) (int, error) {
	buf := <-i.data
	copy(data, buf)

	return len(buf), nil
}

func (i *DataInput) String() string {
	return "Data input: " + i.path
}

func (i *DataInput) emit() {
	var lineStr string

	file, err := os.Open(i.path)
    if err!=nil {
        fmt.Println("dsd")
    }
    scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		lineStr = scanner.Text()

		i.data <- []byte(lineStr)

	}
}


