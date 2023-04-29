package concurentskiplist

import (
	"sync"
	"testing"
)

func TestConcurrentSkipList_Del(t *testing.T) {
	// 构造一个比较函数，用于比较 key 的大小
	compareFunc := func(key1, key2 any) bool {
		return key1.(int) < key2.(int)
	}

	// 构造一个并发安全的跳表
	skipList := NewConcurrentSkipList(compareFunc)

	// 如何测试一个并发场景？
	var wg sync.WaitGroup
	// 依次插入若干个 kv 对
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			skipList.Put(i, i)
		}(i)
	}

	// 删除若干个节点
	for i := 1; i < 100; i += 2 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			skipList.Del(i)
		}(i)
	}

	wg.Wait()

	// 验证删除操作是否正确
	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			if _, ok := skipList.Get(i); !ok {
				t.Errorf("key %v should exist", i)
			}
		}
		if i%2 == 1 {
			if _, ok := skipList.Get(i); ok {
				t.Errorf("key %v should not exist", i)
			}
		}
	}
}

func BenchmarkConcurrentSkipList_Delete(b *testing.B) {
	// 构造一个比较函数，用于比较 key 的大小
	compareFunc := func(key1, key2 any) bool {
		return key1.(int) < key2.(int)
	}

	// 构造一个并发安全的跳表
	skipList := NewConcurrentSkipList(compareFunc)

	var wg sync.WaitGroup // 依次插入若干个 kv 对
	for i := 0; i < b.N; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			skipList.Put(i, i)
		}()
	}

	// 删除若干个节点
	for i := 1; i < b.N; i += 2 {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			skipList.Del(i)
		}()
	}

	wg.Wait()
}
