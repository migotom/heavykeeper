package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/migotom/heavykeeper"
)

func main() {
	fileName := flag.String("f", "", "file name")
	k := flag.Uint("k", 10, "find k top values")
	width := flag.Uint("w", 2048, "array's width, higher value - more memory used but more accurate results")
	depth := flag.Uint("d", 5, "depth, defined amount of buckets in one array")
	decay := flag.Float64("p", 0.9, "probability decay")

	flag.Parse()

	var reader io.Reader
	if *fileName == "" {
		reader = os.Stdin
	} else {
		var err error
		reader, err = os.Open(*fileName)
		if err != nil {
			log.Fatal(err)
		}
	}

	rand.Seed(time.Now().UnixNano())

	heavykeeper := heavykeeper.New(4, uint32(*k), uint32(*width), uint32(*depth), *decay, rand.Int())

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		for _, item := range strings.Fields(line) {
			heavykeeper.Add(item)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("error during scan: ", err)
	}

	heavykeeper.Wait()
	for _, e := range heavykeeper.List() {
		fmt.Println(e.Item, e.Count)
	}
}
