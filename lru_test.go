package cache

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	cache := New(0)

	now := time.Now()
	time.Sleep(time.Millisecond)

	for i := 0; i < 1000; i++ {
		cache.AddValue(uint32(i), fmt.Sprintf("SET key#%d %d", i, i), bytes.Repeat([]byte(strconv.Itoa(i)), 1024))
	}

	if cache.CurrentMemory != 2982140 {
		t.Errorf("size is wrong or the CacheValue has changed its format")
	}

	if cache.Len() != 1000 {
		t.Errorf("expect 1000 objects but got %d", cache.Len())
	}

	if cv, ok := cache.Get(fmt.Sprintf("SET key#%d %d", 481, 481)); !ok {
		t.Error("key not found")
	} else {
		if bytes.Compare(cv.Value, bytes.Repeat([]byte(strconv.Itoa(481)), 1024)) != 0 {
			t.Fatalf("expect 1024 * 500 but got difference")
		}

		if cv.TimeStamp() <= now.UnixNano() {
			t.Errorf("timestamp is wrong. expect > %d but got %d. real: %d", now.UnixNano(), cv.TimeStamp(), cv.timestamp)
		}
	}

	if cv, ok := cache.Get(fmt.Sprintf("SET key#%d %d", 999, 999)); !ok {
		t.Error("key not found")
	} else {
		if bytes.Compare(cv.Value, bytes.Repeat([]byte(strconv.Itoa(999)), 1024)) != 0 {
			t.Fatalf("expect 1024 * 500 but got difference")
		}

		if cv.TimeStamp() <= now.UnixNano() {
			t.Errorf("timestamp is wrong. expect > %d but got %d. real: %d", now.UnixNano(), cv.TimeStamp(), cv.timestamp)
		}
	}
}

func TestCache_MaxMemory(t *testing.T) {
	cache := New(2982140 / 2)

	for i := 0; i < 1000; i++ {
		cache.AddValue(uint32(i), fmt.Sprintf("SET key#%d %d", i, i), bytes.Repeat([]byte(strconv.Itoa(i)), 1024))
	}

	if cache.Len() != 481 {
		t.Fatalf("expect 481 elements but got %d", cache.Len())
	}

	if cache.CurrentMemory > 2982140/2 {
		t.Fatalf("expect < %d but got %d", cache.MaxMemory, cache.CurrentMemory)
	}

	for i := 0; i < 519; i++ {
		_, ok := cache.Get(fmt.Sprintf("SET key#%d %d", i, i))
		if ok {
			t.Fatalf("key %d should not exist", i)
		}
	}
	for i := 519; i < 1000; i++ {
		_, ok := cache.Get(fmt.Sprintf("SET key#%d %d", i, i))
		if !ok {
			t.Fatalf("key 999 should exist")
		}
	}

}
