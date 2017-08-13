package perf

import (
	"encoding/gob"
	"os"
	"fmt"
	"bufio"
	"strings"
	"strconv"
)

func ScoreMap(w *[]float64, features *map[int]float64 ) float64 {
	ret := 0.0
	for k, v := range *features {
		ret += (*w)[k] * v
	}
	return ret
}

func ScoreList(w *[]float64, features *[]float64) float64 {
	ret := 0.0
	for k, v := range *features {
		ret += (*w)[k] * v
	}
	return ret
}

func WriteData(filename string) {
	store := map[string]string{}
	for i := 0; i < 10000000; i++ {
		store[strconv.Itoa(i)] = strconv.Itoa(i)
	}

	file, err := os.Create(filename)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(store)
	}
	file.Close()
}

func LoadData(filename string) {
	store := map[string]string{}
	file, err := os.Open(filename)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(&store)
		if err != nil {
			panic(err)
		}
	} else {
		panic(err)
	}
	file.Close()
	fmt.Println(len(store))
}

func WriteData2(filename string) {
	file, err := os.Create(filename)
	if err == nil {
		writer := bufio.NewWriter(file)
		for i := 0; i < 10000000; i++ {
			writer.WriteString(strconv.Itoa(i) + "\t" + strconv.Itoa(i) + "\n")
		}
	}
	file.Close()
}

func LoadData2(filename string) {
	store := map[string]string{}
	file, err := os.Open(filename)
	if err == nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			items := strings.Split(scanner.Text(), "\t")
			store[items[0]] = items[1]
		}
		if err = scanner.Err(); err != nil {
			panic(err)
		}
	} else {
		panic(err)
	}
	file.Close()
	fmt.Println(len(store))
}
