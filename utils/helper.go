package utils

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func SetupLog(name string) *os.File {
	path := Concatenate("logs/", name, ".txt")
	_ = os.Remove(path)

	f, err := os.OpenFile(path, os.O_RDWR | os.O_CREATE, 0666)
	CheckError(err)

	log.SetOutput(f)
	//log.SetPrefix(name + " ")
	//log.Print("#########################################")
	//log.Println("Start server...")
	return f
}

func Concatenate(elem ...interface{}) string {
	var ipAddress []string
	for _, e := range elem {
		switch v := e.(type) {
		case string:
			ipAddress = append(ipAddress, v)
		case int:
			t := strconv.Itoa(v)
			ipAddress = append(ipAddress, t)
		default:
			fmt.Printf("unexpected type %T", v)
		}
	}

	return strings.Join(ipAddress, "")
}


func CheckError(err error) {
	if err != nil {
		_, fn, line, _ := runtime.Caller(1)
		log.Printf("[error] %s:%d %v", fn, line, err)
		os.Exit(1)
	}
}


func Shuffle(vals []int) []int {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ret := make([]int, len(vals))
	perm := r.Perm(len(vals))
	for i, randIndex := range perm {
		ret[i] = vals[randIndex]
	}
	return ret
}

func Arange(start, stop, step int) []int {
	N := (stop - start) / step
	rnge := make([]int, N, N)
	i := 0
	for x := start; x < stop; x += step {
		rnge[i] = x;
		i += 1
	}
	return rnge
}

func StringAddrToIntArr(addr string) []int {
	addrArr := strings.Split(addr, ".")
	var res = [] int {}

	for _, i := range addrArr {
		j, err := strconv.Atoi(i)
		if err != nil {
			panic(err)
		}
		res = append(res, j)
	}

	return res
}