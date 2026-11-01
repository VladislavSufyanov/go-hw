package hw04lrucache

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList[any]()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList[any]()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})
}

func TestListAny(t *testing.T) {
	l := NewList[any]()
	wanted := []any{1, "abc", 2.1, true, struct{}{}}
	for _, v := range wanted {
		l.PushBack(v)
	}
	result := l.(*list[any]).TransformToSlice()
	require.Equal(t, wanted, result)
}

func TestListRemoveOnlyOneItem(t *testing.T) {
	l := NewList[int]()
	i := l.PushFront(1)
	require.Equal(t, l.Len(), 1)
	require.Equal(t, l.Back(), i)
	require.Equal(t, l.Front(), i)
	l.Remove(i)
	require.Nil(t, l.Back())
	require.Nil(t, l.Front())
	require.Equal(t, l.Len(), 0)
}

type ListTestSuite struct {
	suite.Suite
	Base []int
	L    List[int]
}

func (lts *ListTestSuite) SetupTest() {
	lts.Base = []int{1, 2, 3, 4, 5}
	lts.L = NewList[int]()
	for _, v := range lts.Base {
		lts.L.PushBack(v)
	}
}

func (lts *ListTestSuite) TestListRemoveFront() {
	wanted := lts.Base[1:]
	lts.L.Remove(lts.L.Front())
	result := lts.L.(*list[int]).TransformToSlice()
	require.Equal(lts.T(), wanted, result)
}

func (lts *ListTestSuite) TestListRemoveBack() {
	wanted := lts.Base[:len(lts.Base)-1]
	lts.L.Remove(lts.L.Back())
	result := lts.L.(*list[int]).TransformToSlice()
	require.Equal(lts.T(), wanted, result)
}

func (lts *ListTestSuite) TestListRemoveMiddle() {
	index := len(lts.Base) / 2
	wanted := slices.Delete(lts.Base, index, index+1)
	i := lts.L.Front()
	for range index {
		i = i.Next
	}
	lts.L.Remove(i)
	result := lts.L.(*list[int]).TransformToSlice()
	require.Equal(lts.T(), wanted, result)
}

func (lts *ListTestSuite) TestListMoveToFrontFront() {
	lts.L.MoveToFront(lts.L.Front())
	result := lts.L.(*list[int]).TransformToSlice()
	require.Equal(lts.T(), lts.Base, result)
}

func (lts *ListTestSuite) TestListMoveToFrontBack() {
	wanted := []int{5, 1, 2, 3, 4}
	lts.L.MoveToFront(lts.L.Back())
	result := lts.L.(*list[int]).TransformToSlice()
	require.Equal(lts.T(), wanted, result)
}

func (lts *ListTestSuite) TestListMoveToFront() {
	wanted := []int{3, 1, 2, 4, 5}
	i := lts.L.Front()
	for range 2 {
		i = i.Next
	}
	lts.L.MoveToFront(i)
	result := lts.L.(*list[int]).TransformToSlice()
	require.Equal(lts.T(), wanted, result)
}

func TestListSuite(t *testing.T) {
	suite.Run(t, new(ListTestSuite))
}
