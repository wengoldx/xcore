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

import "reflect"

// Empty struct regard as map value occupy 0 memery size to simulate sets type.
//
//	var sets = map[string]EmptyStruct{
//		"set value 1": utils.E_, "set value 2": struct{}{},
//	}
type EmptyStruct struct{}

// Empty map value for simulate Set type.
var E_ EmptyStruct

// Sets type for cache different types and unique value datas.
type Sets struct {
	sets map[any]EmptyStruct
}

// Create a sets instance.
func NewSets() *Sets {
	return &Sets{sets: make(map[any]EmptyStruct)}
}

// Add the value to sets if not exist.
func (s *Sets) Add(value any) *Sets {
	s.sets[value] = EmptyStruct{}
	return s
}

// Add the values to sets if not exist.
//
//	utils.NewSets().Adds(1, "2", byte(3))
//
//	vs := []any{1, "2", byte(3)}
//	utils.NewSets().Adds(vs...)
func (s *Sets) Adds(values ...any) *Sets {
	for _, value := range values {
		s.sets[value] = EmptyStruct{}
	}
	return s
}

// Add the sets values to self sets if not exist.
func (s *Sets) AddSets(other *Sets) *Sets {
	for _, value := range other.sets {
		s.sets[value] = EmptyStruct{}
	}
	return s
}

// Remove the value from sets.
func (s *Sets) Remove(value any) *Sets {
	delete(s.sets, value)
	return s
}

// Remove the values from sets.
//
//	utils.NewSets().Adds(1, "2", byte(3)).Removes(1, "2")
//
//	vs := []any{"2", byte(3)}
//	utils.NewSets().Adds(1, "2", byte(3)).Removes(vs...)
func (s *Sets) Removes(values ...any) *Sets {
	for _, value := range values {
		delete(s.sets, value)
	}
	return s
}

// Remove the sets values from self sets.
func (s *Sets) RemoveSets(other *Sets) *Sets {
	for _, value := range other.sets {
		delete(s.sets, value)
	}
	return s
}

// Clear sets all values.
func (s *Sets) Clear() *Sets {
	s.sets = make(map[any]EmptyStruct)
	return s
}

// Check the value if exist in sets.
func (s *Sets) Contain(value any) bool {
	_, exist := s.sets[value]
	return exist
}

// Check the values if contain in self sets.
//
//	ok := utils.NewSets().Adds(1, "2", byte(3)).Contains(1, "2")
//
//	vs := []any{"2", byte(3)}
//	ok := utils.NewSets().Adds(1, "2", byte(3)).Contains(vs...)
func (s *Sets) Contains(values ...any) bool {
	for _, ov := range values {
		if _, exist := s.sets[ov]; !exist {
			return false
		}
	}
	return true
}

// Check the sets values if contain in self sets.
func (s *Sets) ContainSets(other *Sets) bool {
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
func (s *Sets) Exists(values ...any) bool {
	for _, ov := range values {
		if _, exist := s.sets[ov]; exist {
			return true
		}
	}
	return false
}

// Fetch all sets values by scaner callback.
func (s *Sets) Fetch(scaning func(e any) bool /* true is break */) {
	for value := range s.sets {
		if scaning(value) {
			break
		}
	}
}

// Return sets values as unordered array.
func (s *Sets) Array() []any {
	rst := []any{}
	for value := range s.sets {
		rst = append(rst, value)
	}
	return rst
}

// ----------------------------------------
// For int type
// ----------------------------------------

// Add the int array values to sets if not exist.
func (s *Sets) AddInts(values []int) *Sets {
	for _, value := range values {
		s.sets[value] = EmptyStruct{}
	}
	return s
}

// Remove the int values from self sets.
func (s *Sets) RemoveInts(values []int) *Sets {
	for _, value := range values {
		delete(s.sets, value)
	}
	return s
}

// Check the int array values if contain in self sets.
func (s *Sets) ContainInts(values []int) bool {
	for _, ov := range values {
		if _, exist := s.sets[ov]; !exist {
			return false
		}
	}
	return true
}

// Fetch the given values and remove the numbers which not contained in sets.
func (s *Sets) FilterInts(values []int) []int {
	valids := []int{}
	for _, ov := range values {
		if _, exist := s.sets[ov]; exist {
			valids = append(valids, ov)
		}
	}
	return valids
}

// Check the int array values if any one exist in self sets.
func (s *Sets) ExistInts(values []int) bool {
	for _, ov := range values {
		if _, exist := s.sets[ov]; exist {
			return true
		}
	}
	return false
}

// Return all exist int values as unorderd int array.
func (s *Sets) ArrayInt() []int {
	rst, vtype := []int{}, reflect.TypeOf(int(0))
	for v := range s.sets {
		if reflect.TypeOf(v) == vtype {
			rst = append(rst, v.(int))
		}
	}
	return rst
}

// ----------------------------------------
// For int64 type
// ----------------------------------------

// Add the int64 array values to sets if not exist.
func (s *Sets) AddInt64s(values []int64) *Sets {
	for _, value := range values {
		s.sets[value] = EmptyStruct{}
	}
	return s
}

// Remove the int64 values from self sets.
func (s *Sets) RemoveInt64s(values []int64) *Sets {
	for _, value := range values {
		delete(s.sets, value)
	}
	return s
}

// Check the int64 array values if contain in self sets.
func (s *Sets) ContainInt64s(values []int64) bool {
	for _, ov := range values {
		if _, exist := s.sets[ov]; !exist {
			return false
		}
	}
	return true
}

// Fetch the given values and remove the numbers which not contained in sets.
func (s *Sets) FilterInt64s(values []int64) []int64 {
	valids := []int64{}
	for _, ov := range values {
		if _, exist := s.sets[ov]; exist {
			valids = append(valids, ov)
		}
	}
	return valids
}

// Check the int64 array values if any one exist in self sets.
func (s *Sets) ExistInt64s(values []int64) bool {
	for _, ov := range values {
		if _, exist := s.sets[ov]; exist {
			return true
		}
	}
	return false
}

// Return all exist int64 values as unorderd int64 array.
func (s *Sets) ArrayInt64() []int64 {
	rst, vtype := []int64{}, reflect.TypeOf(int64(0))
	for v := range s.sets {
		if reflect.TypeOf(v) == vtype {
			rst = append(rst, v.(int64))
		}
	}
	return rst
}

// ----------------------------------------
// For string type
// ----------------------------------------

// Add the string array values to sets if not exist.
func (s *Sets) AddStrings(values []string) *Sets {
	for _, value := range values {
		s.sets[value] = EmptyStruct{}
	}
	return s
}

// Remove the string values from self sets.
func (s *Sets) RemoveStrings(values []string) *Sets {
	for _, value := range values {
		delete(s.sets, value)
	}
	return s
}

// Check the string array values if contain in self sets.
func (s *Sets) ContainStrings(values []string) bool {
	for _, ov := range values {
		if _, exist := s.sets[ov]; !exist {
			return false
		}
	}
	return true
}

// Fetch the given values and remove the strings which not contained in sets.
func (s *Sets) FilterStrings(values []string) []string {
	valids := []string{}
	for _, ov := range values {
		if _, exist := s.sets[ov]; exist {
			valids = append(valids, ov)
		}
	}
	return valids
}

// Check the string array values if any one exist in self sets.
func (s *Sets) ExistStrings(values []string) bool {
	for _, ov := range values {
		if _, exist := s.sets[ov]; exist {
			return true
		}
	}
	return false
}

// Return all exist string values as unorderd string array.
func (s *Sets) ArrayString() []string {
	rst, vtype := []string{}, reflect.TypeOf("")
	for v := range s.sets {
		if reflect.TypeOf(v) == vtype {
			rst = append(rst, v.(string))
		}
	}
	return rst
}
