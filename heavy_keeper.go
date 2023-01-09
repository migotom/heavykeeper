package heavykeeper

import (
	"encoding/binary"
	"math"
	"math/rand"
	"sync"

	"github.com/OneOfOne/xxhash"
	"github.com/migotom/heavykeeper/pkg/minheap"
)

type TopK struct {
	k     uint32
	width uint32
	depth uint32
	decay float64
	items chan []byte

	seed    int
	buckets [][]bucket
	minHeap *minheap.Heap
	wg      *sync.WaitGroup
}

func New(workers int, k, width, depth uint32, decay float64, seed int) *TopK {
	// err check...

	arrays := make([][]bucket, depth)
	for i := range arrays {
		arrays[i] = make([]bucket, width)
	}

	topk := TopK{
		k:       k,
		width:   width,
		depth:   depth,
		decay:   decay,
		buckets: arrays,
		minHeap: minheap.NewHeap(k),
		items:   make(chan []byte),
		wg:      new(sync.WaitGroup),
		seed:    seed,
	}

	for i := 0; i < workers; i++ {
		topk.wg.Add(1)
		go func() {
			defer topk.wg.Done()

			topk.jobAdder()
		}()
	}

	return &topk
}

func (topk *TopK) Wait() {
	close(topk.items)
	topk.wg.Wait()
}

func (topk *TopK) Query(item string) (exist bool) {
	return topk.QueryBytes([]byte(item))
}

func (topk *TopK) Count(item string) (uint64, bool) {
	return topk.CountBytes([]byte(item))
}

func (topk *TopK) Add(item string) {
	topk.AddBytes([]byte(item))
}

func (topk *TopK) QueryBytes(item []byte) (exist bool) {
	_, exist = topk.minHeap.Find(item)
	return
}

func (topk *TopK) CountBytes(item []byte) (uint64, bool) {
	if id, exist := topk.minHeap.Find(item); exist {
		return topk.minHeap.Nodes[id].Count, true
	}
	return 0, false
}

func (topk *TopK) AddBytes(item []byte) {
	topk.items <- item
}

func (topk *TopK) List() []minheap.Node {
	return topk.minHeap.Sorted()
}

func (topk *TopK) jobAdder() {
	for item := range topk.items {

		itemFingerprint := xxhash.Checksum64S(item, uint64(topk.seed))

		var maxCount uint64
		itemHeapIdx, itemHeapExist := topk.minHeap.Find(item)
		heapMin := topk.minHeap.Min()

		// compute d hashes
		for i := uint32(0); i < topk.depth; i++ {

			bI := make([]byte, 4)
			binary.LittleEndian.PutUint32(bI, uint32(i))

			bucketNumber := xxhash.Checksum64S(append([]byte(item), bI...), uint64(topk.seed)) % uint64(topk.width)

			topk.buckets[i][bucketNumber].Lock()

			fingerprint := topk.buckets[i][bucketNumber].fingerprint
			count := topk.buckets[i][bucketNumber].count

			if count == 0 {
				topk.buckets[i][bucketNumber].fingerprint = itemFingerprint
				topk.buckets[i][bucketNumber].count = 1
				maxCount = max(maxCount, 1)

			} else if fingerprint == itemFingerprint {
				if itemHeapExist || count <= heapMin {
					topk.buckets[i][bucketNumber].count++
					maxCount = max(maxCount, topk.buckets[i][bucketNumber].count)
				}

			} else {
				decay := math.Pow(topk.decay, float64(count))
				if rand.Float64() < decay {
					topk.buckets[i][bucketNumber].count--
					if topk.buckets[i][bucketNumber].count == 0 {
						topk.buckets[i][bucketNumber].fingerprint = itemFingerprint
						topk.buckets[i][bucketNumber].count = 1
						maxCount = max(maxCount, 1)
					}
				}
			}

			topk.buckets[i][bucketNumber].Unlock()
		}

		// update heap
		if itemHeapExist {
			topk.minHeap.Fix(itemHeapIdx, maxCount)
		} else {
			topk.minHeap.Add(minheap.Node{
				Count: maxCount,
				Item:  item,
			})
		}

	}

}

type bucket struct {
	fingerprint uint64
	count       uint64
	sync.Mutex
}

func (b *bucket) Get() (uint64, uint64) {
	// b.RLock()
	// defer b.RUnlock()

	return b.fingerprint, b.count
}

func (b *bucket) Set(fingerprint, count uint64) {
	// b.Lock()
	// defer b.Unlock()

	b.fingerprint = fingerprint
	b.count = count
}

func (b *bucket) Inc(val uint64) uint64 {
	// b.Lock()
	// defer b.Unlock()

	b.count += val
	return b.count
}

func max(x, y uint64) uint64 {
	if x > y {
		return x
	}
	return y
}
