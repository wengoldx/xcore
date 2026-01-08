// Copyright (c) 2018-Now Dunyu All Rights Reserved.
//
// Author      : https://www.wengold.net
// Email       : support@wengold.net
//
// Prismy.No | Date       | Modified by. | Description
// -------------------------------------------------------------------
// 00001       2019/05/22   yangping       New version
// -------------------------------------------------------------------

package utils

// Empty struct regard as map value occupy 0 memery size to simulate sets type.
//
//	var sets = map[string]None{
//		"set value 1": utils.NONE, "set value 2": struct{}{},
//	}
type TNone struct{}

// Empty map value for simulate Set value.
var NONE TNone

// Sets type for cache the T typed datas as set container.
//
//	utils.NewSets[int]()               // create a empty int type sets.
//	utils.NewSets(1, -5, 678)          // create a int type sets with init datas.
//	utils.NewSets().Add("abc", "123")  // create a string type sets and add datas.
//	utils.NewSets[int64]().Add(-8, 35) // create a int64 type sets and  add datas.
//	// ...
//
// The sets datas do well in cache unique values, fast find target, and
// filter the given values by Contain, Exist, Filter methods.
type Sets[T any] struct {
	sets map[any]TNone
}

// Create a sets instance.
//
//	utils.NewSets[int]()               // create a empty int type sets.
//	utils.NewSets(1, -5, 678)          // create a int type sets with init datas.
func NewSets[T any](datas ...T) *Sets[T] {
	ins := &Sets[T]{sets: make(map[any]TNone)}
	return ins.Add(datas...)
}

// Return sets counts.
func (s *Sets[T]) Size() int {
	return len(s.sets)
}

// Add the values to sets if not exist.
//
//	utils.NewSets().Add(1, -5, 678)    // create a int type datas sets.
//	utils.NewSets().Add("abc", "123")  // create a string datas sets.
//	utils.NewSets[int64]().Add(-8, 35) // create a int64 type datas sets.
func (s *Sets[T]) Add(values ...T) *Sets[T] {
	for _, value := range values {
		s.sets[value] = TNone{}
	}
	return s
}

// Add the sets values to self sets if not exist.
//
//	a := utils.NewSets(1, -5, 678)
//	b := utils.NewSets[int]().AddSets(a)  // add all a sets values.
func (s *Sets[T]) AddSets(other *Sets[T]) *Sets[T] {
	for value := range other.sets {
		s.sets[value] = TNone{}
	}
	return s
}

// Remove the values from sets.
//
//	utils.NewSets(1, 2, 3).Remove(1, 3)  // remain [2].
//
//	vs := []int{1, 3}
//	utils.NewSets(1, 2, 3).Remove(vs...) // remain [2].
func (s *Sets[T]) Remove(values ...T) *Sets[T] {
	for _, value := range values {
		delete(s.sets, value)
	}
	return s
}

// Remove the sets values from self sets.
//
//	a := utils.NewSets(1, 3)
//	b := utils.NewSets(2, 3, 5).RemoveSets(a) // remain [2, 5].
func (s *Sets[T]) RemoveSets(other *Sets[T]) *Sets[T] {
	for value := range other.sets {
		delete(s.sets, value)
	}
	return s
}

// Clear sets all values.
func (s *Sets[T]) Clear() *Sets[T] {
	clear(s.sets)
	return s
}

// Check the values if contain in self sets.
//
//	as := utils.NewSets(1, 2, 3)
//	ok := as.Contain(1, 3)  // matched all targets.
//
//	vs := []int{-5, 1, 3}
//	ng := as.Contain(vs...) // exist one outof sets.
func (s *Sets[T]) Contain(values ...T) bool {
	for _, ov := range values {
		if _, exist := s.sets[ov]; !exist {
			return false
		}
	}
	return true
}

// Check the sets values if contain in self sets.
//
//	as := utils.NewSets(1, 2, 3)
//
//	oks := utils.NewSets(1, 3)
//	ok := as.ContainSets(oks) // matched all targets.
//
//	ngs := utils.NewSets(-5, 1, 3)
//	ng := as.ContainSets(ngs) // exist one outof sets.
func (s *Sets[T]) ContainSets(other *Sets[T]) bool {
	for ov := range other.sets {
		if _, exist := s.sets[ov]; !exist {
			return false
		}
	}
	return true
}

// Check the values if any one exist in self sets.
//
//	ok := utils.NewSets(1, 2, 3).Exist(1, 2, 5) // exist one at lest.
//	ng := utils.NewSets(1, 2, 3).Exist(-5, 9)   // unmatched all targets.
func (s *Sets[T]) Exist(values ...T) bool {
	for _, ov := range values {
		if _, exist := s.sets[ov]; exist {
			return true
		}
	}
	return false
}

// Return sets values as unordered array.
func (s *Sets[T]) Array() []T {
	rst := []T{}
	for value := range s.sets {
		if t, ok := value.(T); ok {
			rst = append(rst, t)
		}
	}
	return rst
}

// Fetch all sets values by scaner callback.
func (s *Sets[T]) Fetch(scaning func(v T) bool) {
	for value := range s.sets {
		if t, ok := value.(T); ok {
			if interupt := scaning(t); interupt {
				break
			}
		}
	}
}

// Fetch the given values and remove the items which not contained in sets.
//
//	utils.NewSets(1, 2, 3).Filters(2, 3, 6) // remain [2, 3]
func (s *Sets[T]) Filters(values ...T) []T {
	valids := []T{}
	for _, ov := range values {
		if _, exist := s.sets[ov]; exist {
			valids = append(valids, ov)
		}
	}
	return valids
}
