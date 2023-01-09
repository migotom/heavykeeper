package heavykeeper

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"
)

const (
	HashSeed = 123456890
)

func TestAdd(t *testing.T) {
	type config struct {
		k     uint32
		width uint32
		depth uint32
		decay float64
		seed  int
	}
	cases := []struct {
		Name                     string
		Config                   config
		Items                    []string
		MinAccuracy, MaxAccuracy float64
	}{
		{
			Name:        "uniform stream",
			Config:      config{2, 64, 6, 0.9, HashSeed},
			Items:       []string{"a", "a", "b", "b", "b", "b", "b", "c", "c", "c", "d", "e", "e", "e", "e", "e", "e", "f", "f", "g", "g"},
			MinAccuracy: 0.5,
			MaxAccuracy: 1.0,
		},
		{
			Name:        "mixed stream",
			Config:      config{2, 64, 6, 0.9, HashSeed},
			Items:       []string{"a", "d", "f", "a", "z", "b", "c", "b", "d", "d", "c", "d", "b", "e", "e", "f", "f", "x", "f", "y"},
			MinAccuracy: 0.5,
			MaxAccuracy: 1.0,
		},
		{
			Name:        "dominated by one value stream",
			Config:      config{2, 32, 6, 0.9, HashSeed},
			Items:       []string{"a", "d", "f", "a", "z", "a", "c", "b", "a", "d", "c", "d", "a", "e", "a", "f", "a", "f", "x", "f", "y", "a"},
			MinAccuracy: 0.5,
			MaxAccuracy: 1.0,
		},
		{
			Name:        "dominated by two values streams",
			Config:      config{2, 64, 6, 0.9, HashSeed},
			Items:       []string{"a", "d", "f", "a", "b", "a", "c", "b", "a", "b", "c", "d", "a", "b", "a", "f", "b", "a", "b", "x", "f", "b", "y", "a"},
			MinAccuracy: 0.5,
			MaxAccuracy: 1.0,
		},
		{
			Name:   `Tolstoy's "War and Peace" stream`,
			Config: config{20, 8192, 6, 0.9, HashSeed},
			Items: func(t *testing.T) (items []string) {
				f, err := os.Open("fixtures/war_and_peace.txt")
				if err != nil {
					t.Fatal(err)
				}
				scanner := bufio.NewScanner(f)
				for scanner.Scan() {
					line := scanner.Text()
					items = append(items, strings.Fields(line)...)
				}
				if err := scanner.Err(); err != nil {
					t.Fatal("error during scan: ", err)
				}
				return
			}(t),
			MinAccuracy: 0.8,
			MaxAccuracy: 1.0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			// create heavy keeper TopK pseudo-service
			heavykeeper := New(4, tc.Config.k, tc.Config.width, tc.Config.depth, tc.Config.decay, HashSeed)

			// iterate over stream and build assistant metrics
			frequencies := frequencies{counts: make(map[string]uint64)}
			for _, item := range tc.Items {
				heavykeeper.Add(item)
				frequencies.counts[item]++
			}

			// sort assistant metrics
			frequencies.Sort()

			// compare heavy keeper Top K results with sorted
			var errors []error
			var accuracy []float64

			heavykeeper.Wait()

			topkList := heavykeeper.List()
			for i := uint32(0); i < tc.Config.k; i++ {
				hotspot := heavykeeper.Query(frequencies.keys[i])
				count, _ := heavykeeper.Count(frequencies.keys[i])
				accuracy = append(accuracy, calculateAccuracy(frequencies.counts[frequencies.keys[i]], count))
				// item position test
				if !bytes.Equal(topkList[i].Item, []byte(frequencies.keys[i])) {
					errors = append(errors, fmt.Errorf("TopK Key mismatch expected idx=%d key=%s count=%d, got key=%s count=%d", i, frequencies.keys[i], frequencies.counts[frequencies.keys[i]], topkList[i].Item, count))
				}
				// item counter test
				if !hotspot || count != frequencies.counts[frequencies.keys[i]] {
					errors = append(errors, fmt.Errorf("TopK Counter mismatch expected idx=%d key=%s count=%d, got key=%s count=%d", i, frequencies.keys[i], frequencies.counts[frequencies.keys[i]], topkList[i].Item, count))
				}
			}
			a := overallAccuracy(accuracy)

			if a < tc.MinAccuracy || a > tc.MaxAccuracy {
				t.Errorf("TopK expected accuracy %v, got %v", tc.MinAccuracy, a)
				for _, e := range errors {
					t.Error(e)
				}
			}
		})
	}
}

func calculateAccuracy(expected, got uint64) float64 {
	return 1.0 - (float64(expected)-float64(got))/float64(expected)
}

func overallAccuracy(a []float64) float64 {
	var overall float64
	for _, accuracy := range a {
		overall += accuracy
	}
	return (overall / float64(len(a)))
}

type frequencies struct {
	keys   []string
	counts map[string]uint64
}

func (f frequencies) Len() int {
	return len(f.keys)
}
func (f *frequencies) Less(i, j int) bool {
	return f.counts[f.keys[i]] > f.counts[f.keys[j]] || f.counts[f.keys[i]] == f.counts[f.keys[j]] && f.keys[i] < f.keys[j]
}
func (f *frequencies) Swap(i, j int) {
	f.keys[i], f.keys[j] = f.keys[j], f.keys[i]
}
func (f *frequencies) Sort() {
	for key := range f.counts {
		f.keys = append(f.keys, key)
	}
	sort.Sort(f)
}
