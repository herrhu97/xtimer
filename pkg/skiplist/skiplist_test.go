package skiplist

import (
	"testing"
)

func TestSkiplist(t *testing.T) {
	sl := &Skiplist{&node{nexts: make([]*node, 1)}}
	testdata := []struct {
		key, val int
	}{
		{1, 10},
		{2, 20},
		{3, 30},
		{4, 40},
		{5, 50},
	}

	for _, v := range testdata {
		sl.Put(v.key, v.val)
	}

	for _, v := range testdata {
		if val, ok := sl.Get(v.key); !ok || val != v.val {
			t.Errorf("Skiplist Get failed. key: %d, expected value: %d, got value: %d", v.key, v.val, val)
		}
	}

	notExistKey := 6
	if _, ok := sl.Get(notExistKey); ok {
		t.Errorf("Skiplist Get failed. key: %d should not exist but got a value", notExistKey)
	}

	sl.Del(1)

	if _, ok := sl.Get(1); ok {
		t.Errorf("Skiplist Del failed. Key 1 should have been deleted.")
	}
}

func TestRoll(t *testing.T) {
	sl := &Skiplist{}
	cnt := make(map[int]int)
	samples := 100000
	for i := 0; i < samples; i++ {
		level := sl.roll()
		cnt[level]++
	}
	t.Logf("roll test result: %+v", cnt)
}
