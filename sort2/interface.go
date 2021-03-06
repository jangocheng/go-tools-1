// Copyright 2019 xgfone
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sort2

import "sort"

type interfaceSlice struct {
	data []interface{}
	less func(interface{}, interface{}) bool
}

func (s interfaceSlice) Len() int {
	return len(s.data)
}

func (s interfaceSlice) Less(i, j int) bool {
	return s.less(s.data[i], s.data[j])
}

func (s interfaceSlice) Swap(i, j int) {
	tmp := s.data[i]
	s.data[i] = s.data[j]
	s.data[j] = tmp
}

// Interfaces sorts the interface slice, then returns the sorted slice.
//
// The elements in data should have the same type, or it maybe panic, which
// depends on the less function.
//
// less is a function to compare the two elements of the slice data,
// which returns true if the first is less than the second, or returns false.
//
// If giving the second argument and it's true, sort the data in reverse.
func Interfaces(data []interface{}, less func(first, second interface{}) bool) {
	if len(data) > 1 {
		sort.Sort(interfaceSlice{data: data, less: less})
	}
}
