package recommend

import (
	"bookkeeper/internal/model"
	"container/heap"
)

type PriorityQueue []model.ScoredDocument

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Score < pq[j].Score
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x any) {
	*pq = append(*pq, x.(model.ScoredDocument))
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

func InsertIntoQueue(pq *PriorityQueue, doc model.ScoredDocument, maxSize int) {
	if pq.Len() < maxSize {
		heap.Push(pq, doc)
	} else if (*pq)[0].Score < doc.Score {
		(*pq)[0] = doc
		heap.Fix(pq, 0)
	}
}
