package hw04lrucache

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache[any](10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache[any](5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache[int](3)
		keys := []Key{"a", "b", "c"}
		for i, k := range keys {
			wasInCache := c.Set(k, i)
			require.False(t, wasInCache, fmt.Sprintf("key: %s", k))
		}
		for i, k := range keys {
			got, wasInCache := c.Get(k)
			require.Equal(t, i, got, fmt.Sprintf("key: %s", k))
			require.True(t, wasInCache, fmt.Sprintf("key: %s", k))
		}
		c.Clear()
		for _, k := range keys {
			got, wasInCache := c.Get(k)
			require.Equal(t, 0, got, fmt.Sprintf("key: %s", k))
			require.False(t, wasInCache, fmt.Sprintf("key: %s", k))
		}
	})

	t.Run("displace logic", func(t *testing.T) {
		c := NewCache[int](3)
		c.Set("a", 0)
		c.Set("b", 1)
		c.Set("c", 2)
		c.Set("d", 3)
		got, wasInCache := c.Get("a")
		require.Equal(t, got, 0)
		require.False(t, wasInCache)
		for i, k := range []Key{"b", "c", "d"} {
			got, wasInCache := c.Get(k)
			require.Equal(t, i+1, got, fmt.Sprintf("key: %s", k))
			require.True(t, wasInCache, fmt.Sprintf("key: %s", k))
		}
	})

	t.Run("displace old item logic", func(t *testing.T) {
		c := NewCache[int](3)
		c.Set("a", 0)
		c.Set("b", 1)
		c.Set("c", 2)
		for _, k := range []Key{"b", "a"} {
			_, wasInCache := c.Get(k)
			require.True(t, wasInCache, fmt.Sprintf("key: %s", k))
		}
		wasInCache := c.Set("d", 3)
		require.False(t, wasInCache)
		got, wasInCache := c.Get("c")
		require.Equal(t, got, 0)
		require.False(t, wasInCache)
		keys := []Key{"d", "a", "b"}
		for _, k := range keys {
			_, wasInCache := c.Get(k)
			require.True(t, wasInCache, fmt.Sprintf("key: %s", k))
		}
	})

	t.Run("displace old item set logic", func(t *testing.T) {
		c := NewCache[int](3)
		c.Set("a", 0)
		c.Set("b", 1)
		c.Set("c", 2)
		wasInCache := c.Set("a", 100)
		require.True(t, wasInCache)
		v, wasInCache := c.Get("a")
		require.Equal(t, v, 100)
		require.True(t, wasInCache)
		wasInCache = c.Set("d", 3)
		require.False(t, wasInCache)
		got, wasInCache := c.Get("b")
		require.Equal(t, got, 0)
		require.False(t, wasInCache)
		keys := []Key{"d", "a", "c"}
		for _, k := range keys {
			_, wasInCache := c.Get(k)
			require.True(t, wasInCache, fmt.Sprintf("key: %s", k))
		}
	})
}

func TestCacheMultithreading(_ *testing.T) {
	c := NewCache[any](10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
