package cache

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"sync"
	"testing"
)

func TestCache_Set_NewElementInEmptyCacheSync_ShouldAddElementInQueueAndMap(t *testing.T) {
	// Arrange
	cache := NewLRUCache(10)

	// Act
	cache.Set("key", 15)

	val, ok := cache.records["key"]

	it := val.Value.(Record)

	// Assert
	assert.Equal(t, it.Value.(int), 15)
	assert.True(t, ok)

	assert.Equal(t, 15, cache.queue.Back().Value.(Record).Value.(int))
	assert.Equal(t, 15, cache.queue.Front().Value.(Record).Value.(int))

	assert.Equal(t, 1, cache.queue.Len())
}

func TestCache_Set_NewElementInFilledCacheSync_ShouldAddElementAndDeleteOldestOne(t *testing.T) {
	// Arrange
	cache := NewLRUCache(3)

	cache.Set("key", 15)
	cache.Set("key1", 20)
	cache.Set("key2", 30)

	// Act
	cache.Set("key3", 50)

	val1, ok1 := cache.records["key3"]
	_, ok2 := cache.records["key"]

	it := val1.Value.(Record)

	// Assert
	assert.Equal(t, it.Value.(int), 50)
	assert.True(t, ok1)

	assert.False(t, ok2)

	assert.Equal(t, 20, cache.queue.Back().Value.(Record).Value.(int))
	assert.Equal(t, 50, cache.queue.Front().Value.(Record).Value.(int))

	assert.Equal(t, 3, cache.queue.Len())
}

func TestCache_Set_AddElementThatWasExistsAndWasTheOldestSync_ShouldReplaceElementAndMoveItToTheTopOfList(t *testing.T) {
	// Arrange
	cache := NewLRUCache(10)
	cache.Set("key", 15)
	cache.Set("key2", 13)
	cache.Set("key1", 20)

	// Act
	cache.Set("key", 12)
	val, ok := cache.records["key"]

	it := val.Value.(Record)

	// Assert
	assert.Equal(t, 12, it.Value.(int))
	assert.True(t, ok)

	assert.Equal(t, 13, cache.queue.Back().Value.(Record).Value.(int))
	assert.Equal(t, 12, cache.queue.Front().Value.(Record).Value.(int))

	assert.Equal(t, 3, cache.queue.Len())
}

func TestCache_Set_AddElementsAsync_ShouldSaveChangesFromDifferentGoroutines(t *testing.T) {
	// Arrange
	cache := NewLRUCache(6)

	// Act
	wg := sync.WaitGroup{}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(x int) {
			cache.Set(strconv.Itoa(x), x)
			wg.Done()
		}(i)
	}
	wg.Wait()

	val1, ok1 := cache.records["0"]
	it1 := val1.Value.(Record)
	val2, ok2 := cache.records["1"]
	it2 := val2.Value.(Record)
	val3, ok3 := cache.records["2"]
	it3 := val3.Value.(Record)
	val4, ok4 := cache.records["3"]
	it4 := val4.Value.(Record)
	val5, ok5 := cache.records["4"]
	it5 := val5.Value.(Record)

	// Assert
	assert.True(t, ok1)
	assert.Equal(t, 0, it1.Value.(int))

	assert.True(t, ok2)
	assert.Equal(t, 1, it2.Value.(int))

	assert.True(t, ok3)
	assert.Equal(t, 2, it3.Value.(int))

	assert.True(t, ok4)
	assert.Equal(t, 3, it4.Value.(int))

	assert.True(t, ok5)
	assert.Equal(t, 4, it5.Value)

	assert.Equal(t, 5, cache.queue.Len())
}

func TestCache_Get_GetExistedElementFromEmptyListSync_ShouldReturnElement(t *testing.T) {
	// Arrange
	cache := NewLRUCache(3)
	cache.Set("key", 15)

	// Act
	val, ok := cache.Get("key")

	assert.Equal(t, 15, val.Value.(int))
	assert.True(t, ok)

	assert.Equal(t, 15, cache.queue.Front().Value.(Record).Value.(int))
	assert.Equal(t, 15, cache.queue.Back().Value.(Record).Value.(int))
}

func TestCache_Get_GetExistedElementFromFilledListSync_ShouldReturnElement(t *testing.T) {
	// Arrange
	cache := NewLRUCache(5)
	cache.Set("key", 15)
	cache.Set("key1", 16)
	cache.Set("key2", 17)
	cache.Set("key3", 18)

	// Act
	val, ok := cache.Get("key")

	assert.Equal(t, 15, val.Value.(int))
	assert.True(t, ok)

	assert.Equal(t, 15, cache.queue.Front().Value.(Record).Value.(int))
	assert.Equal(t, 16, cache.queue.Back().Value.(Record).Value.(int))
}

func TestCache_Get_GetNotExistedElementSync_ShouldReturnFalse(t *testing.T) {
	// Arrange
	cache := NewLRUCache(10)
	cache.Set("key", 15)
	cache.Set("key2", 12)
	cache.Set("key1", 20)
	cache.Set("key1", 13)

	// Act
	val, ok := cache.Get("ke")

	// Assert
	assert.Equal(t, Record{}, val)
	assert.False(t, ok)
}

func TestCache_Get_GetExistedElementTwiceSync_ShouldReturnElement(t *testing.T) {
	// Arrange
	cache := NewLRUCache(10)
	cache.Set("key", 15)
	cache.Set("key2", 12)
	cache.Set("key1", 20)
	cache.Set("key3", 13)

	// Act
	val1, ok1 := cache.Get("key")
	val2, ok2 := cache.Get("key")

	// Assert
	assert.Equal(t, 15, val1.Value.(int))
	assert.False(t, ok1)

	assert.Equal(t, 15, val2.Value.(int))
	assert.False(t, ok2)
}

func TestCache_Remove_DeleteLastAddedElementSync_ShouldRemoveElement(t *testing.T) {
	// Arrange
	cache := NewLRUCache(10)
	cache.Set("key", 15)
	cache.Set("key2", 12)
	cache.Set("key1", 20)

	// Act
	cache.Remove("key1")
	_, ok := cache.records["key1"]

	//Assert
	assert.False(t, ok)

	assert.Equal(t, 15, cache.queue.Back().Value.(Record).Value.(int))
	assert.Equal(t, 12, cache.queue.Front().Value.(Record).Value.(int))

	assert.Equal(t, 2, cache.queue.Len())
}

func TestCache_Remove_DeleteTheFirstElementSync_ShouldRemoveElement(t *testing.T) {
	// Arrange
	cache := NewLRUCache(10)
	cache.Set("key", 15)
	cache.Set("key2", 12)
	cache.Set("key1", 20)

	// Act
	cache.Remove("key")
	_, ok := cache.records["key"]

	// Assert
	assert.False(t, ok)

	assert.Equal(t, 12, cache.queue.Back().Value.(Record).Value.(int))
	assert.Equal(t, 20, cache.queue.Front().Value.(Record).Value.(int))

	assert.Equal(t, 2, cache.queue.Len())
}

func TestCache_Remove_DeleteNotExistedElementSync_ShouldNotRemoveElement(t *testing.T) {
	// Arrange
	cache := NewLRUCache(10)
	cache.Set("key", 15)
	cache.Set("key2", 12)
	cache.Set("key1", 20)

	// Act
	cache.Remove("ke")

	// Assert
	assert.Equal(t, 15, cache.queue.Back().Value.(Record).Value.(int))
	assert.Equal(t, 20, cache.queue.Front().Value.(Record).Value.(int))

	assert.Equal(t, 3, cache.queue.Len())
}

func TestCache_Remove_DeleteFewElementAsync_ShouldRemoveElements(t *testing.T) {
	// Arrange
	cache := NewLRUCache(10)
	cache.Set("0", 15)
	cache.Set("1", 12)
	cache.Set("2", 20)
	cache.Set("3", 25)

	// Act
	wg := sync.WaitGroup{}
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(x int) {
			cache.Remove(strconv.Itoa(x))
			wg.Done()
		}(i)
	}
	wg.Wait()

	_, ok1 := cache.records["0"]
	_, ok2 := cache.records["1"]
	_, ok3 := cache.records["2"]
	_, ok4 := cache.records["3"]

	// Assert
	assert.False(t, ok1)
	assert.False(t, ok2)
	assert.False(t, ok3)
	assert.False(t, ok4)

	assert.Equal(t, 0, cache.queue.Len())
}

func TestCache_Len_EmptyCache_ReturnZero(t *testing.T) {
	// Arrange
	cache := NewLRUCache(10)

	// Act
	count := cache.Len()

	//Assert
	assert.Zero(t, count)
}

func TestCache_Len_Filled_ReturnLenOfCache(t *testing.T) {
	// Arrange
	cache := NewLRUCache(5)
	cache.Set("0", 0)
	cache.Set("1", 1)
	cache.Set("2", 2)
	cache.Set("3", 3)
	cache.Set("4", 4)

	// Act
	count := cache.Len()

	//Assert
	assert.Equal(t, 5, count)
}
