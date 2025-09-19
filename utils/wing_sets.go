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
//	var sets = map[string]EmptyStruct{
//		"set value 1": utils.E_, "set value 2": struct{}{},
//	}
type EmptyStruct struct{}

// Empty map value for simulate Set type.
var E_ EmptyStruct

// Sets type for cache different types and unique value datas.
type Sets[T any] struct {
	sets map[any]EmptyStruct
}

// Create a sets instance.
func NewSets[T any]() *Sets[T] {
	return &Sets[T]{sets: make(map[any]EmptyStruct)}
}

// Return sets counts.
func (s *Sets[T]) Size() int {
	return len(s.sets)
}

// Add the values to sets if not exist.
//
//	utils.NewSets().Adds(1, "2", byte(3))
//
//	vs := []any{1, "2", byte(3)}
//	utils.NewSets().Add(vs...)
func (s *Sets[T]) Add(values ...T) *Sets[T] {
	for _, value := range values {
		s.sets[value] = EmptyStruct{}
	}
	return s
}

// Add the sets values to self sets if not exist.
func (s *Sets[T]) AddSets(other *Sets[T]) *Sets[T] {
	for _, value := range other.sets {
		s.sets[value] = EmptyStruct{}
	}
	return s
}

// Remove the values from sets.
//
//	utils.NewSets().Adds(1, "2", byte(3)).Removes(1, "2")
//
//	vs := []any{"2", byte(3)}
//	utils.NewSets().Adds(1, "2", byte(3)).Removes(vs...)
func (s *Sets[T]) Remove(values ...T) *Sets[T] {
	for _, value := range values {
		delete(s.sets, value)
	}
	return s
}

// Remove the sets values from self sets.
func (s *Sets[T]) RemoveSets(other *Sets[T]) *Sets[T] {
	for _, value := range other.sets {
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
//	ok := utils.NewSets().Adds(1, "2", byte(3)).Contains(1, "2")
//
//	vs := []any{"2", byte(3)}
//	ok := utils.NewSets().Adds(1, "2", byte(3)).Contains(vs...)
func (s *Sets[T]) Contain(values ...T) bool {
	for _, ov := range values {
		if _, exist := s.sets[ov]; !exist {
			return false
		}
	}
	return true
}

// Check the sets values if contain in self sets.
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
//	ok := utils.NewSets().Adds(1, "2", byte(3)).Exists(1, "2")
//
//	vs := []any{"2", byte(3)}
//	ok := utils.NewSets().Adds(1, "2", byte(3)).Exists(vs...)
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
func (s *Sets[T]) Fetch(scaning func(e T) bool /* true is break */) {
	for value := range s.sets {
		if t, ok := value.(T); ok {
			if scaning(t) {
				break
			}
		}
	}
}

// Fetch the given values and remove the items which not contained in sets.
func (s *Sets[T]) Filters(values ...T) []T {
	valids := []T{}
	for _, ov := range values {
		if _, exist := s.sets[ov]; exist {
			valids = append(valids, ov)
		}
	}
	return valids
}
