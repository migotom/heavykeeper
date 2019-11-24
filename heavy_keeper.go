package heavykeeper

import (
	"math"
	"math/rand"

	"github.com/OneOfOne/xxhash"
	"github.com/migotom/heavykeeper/pkg/minheap"
)

type TopK struct {
	k     uint32
	width uint32
	depth uint32
	decay float64

	buckets [][]*bucket
	minHeap *minheap.Heap
}

func New(k, width, depth uint32, decay float64) *TopK {
	// err check...
	arrays := make([][]*bucket, depth)
	for i := range arrays {
		arrays[i] = make([]*bucket, width)
		for j := range arrays[i] {
			arrays[i][j] = new(bucket)
		}
	}

	return &TopK{
		k:       k,
		width:   width,
		depth:   depth,
		decay:   decay,
		buckets: arrays,
		minHeap: minheap.NewHeap(k),
	}
}

func (topk *TopK) Query(item string) (exist bool) {
	_, exist = topk.minHeap.Find(item)
	return
}

func (topk *TopK) Count(item string) (uint64, bool) {
	if id, exist := topk.minHeap.Find(item); exist {
		return topk.minHeap.Nodes[id].Count, true
	}
	return 0, false
}

func (topk *TopK) List() []minheap.Node {
	return topk.minHeap.Sorted()
}

func (topk *TopK) Add(item string) uint64 {
	var maxCount uint64
	var itemHeapIdx int
	var itemHeapExist, minHeapSearched bool

	heapMin := topk.minHeap.Min()
	itemFingerprint := xxhash.Checksum64([]byte(item))

	// compute d hashes
	for i := uint32(0); i < topk.depth; i++ {

		bucketNumber := xxhash.Checksum64S([]byte(item), uint64(i)) % uint64(topk.width)
		bucket := topk.buckets[i][bucketNumber]

		count := &bucket.count

		if *count == 0 {
			bucket.fingerprint = itemFingerprint
			maxCount = max(maxCount, 1)
			(*count) = 1

		} else if bucket.fingerprint == itemFingerprint {

			if !minHeapSearched && (*count) >= heapMin {
				itemHeapIdx, itemHeapExist = topk.minHeap.Find(item)
				minHeapSearched = true
			}
			if itemHeapExist || (*count) <= heapMin {
				(*count)++
				maxCount = max(maxCount, (*count))
			}

		} else {
			decay := math.Pow(topk.decay, float64(*count))
			if rand.Float64() < decay {
				(*count)--
				if (*count) == 0 {
					(*count) = 1
					maxCount = max(maxCount, (*count))
					bucket.fingerprint = itemFingerprint
				}
			}
		}
	}

	// update heap
	if itemHeapExist {
		topk.minHeap.Nodes[itemHeapIdx].Count = maxCount
		topk.minHeap.Fix(itemHeapIdx)
	} else {
		topk.minHeap.Add(minheap.Node{
			Count: maxCount,
			Item:  item,
		})
	}

	return maxCount
}

type bucket struct {
	fingerprint uint64
	count       uint64
}

func max(x, y uint64) uint64 {
	if x > y {
		return x
	}
	return y
}
