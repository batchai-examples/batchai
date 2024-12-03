package comm

import (
	"testing"
)

func TestSet2Strings(t *testing.T) {
	// Test case: Empty set should return an empty slice
	set := hashset.New()
	result := Set2Strings(set)
	if len(result) != 0 {
		t.Errorf("Expected empty slice, got %v", result)
	}

	// Test case: Set with one element
	set.Add("one")
	result = Set2Strings(set)
	if len(result) != 1 || result[0] != "one" {
		t.Errorf("Expected slice with one element 'one', got %v", result)
	}

	// Test case: Set with multiple elements
	set.Add("two")
	set.Add("three")
	result = Set2Strings(set)
	expected := []string{"one", "two", "three"}
	if !SliceEquals(result, expected) {
		t.Errorf("Expected slice with elements 'one', 'two', 'three', got %v", result)
	}
}

func TestSlice2Set(t *testing.T) {
	// Test case: Empty slice should return an empty set
	slice := []string{}
	result := Slice2Set(slice...)
	if result.Size() != 0 {
		t.Errorf("Expected empty set, got size %d", result.Size())
	}

	// Test case: Slice with one element
	slice = []string{"one"}
	result = Slice2Set(slice...)
	if result.Size() != 1 || !result.Contains("one") {
		t.Errorf("Expected set with one element 'one', got size %d, contains %v", result.Size(), result.Values())
	}

	// Test case: Slice with multiple elements
	slice = []string{"one", "two", "three"}
	result = Slice2Set(slice...)
	expected := hashset.New()
	expected.Add("one")
	expected.Add("two")
	expected.Add("three")
	if !result.Equals(expected) {
		t.Errorf("Expected set with elements 'one', 'two', 'three', got %v", result.Values())
	}
}

func TestSlice2Map(t *testing.T) {
	// Test case: Empty slice should return an empty map
	slice := []string{}
	keyFunc := func(v string) string { return v }
	result := Slice2Map(slice, keyFunc)
	if len(result) != 0 {
		t.Errorf("Expected empty map, got %v", result)
	}

	// Test case: Slice with one element
	slice = []string{"one"}
	result = Slice2Map(slice, keyFunc)
	expected := map[string]string{"one": "one"}
	if !SliceEquals(result, expected) {
		t.Errorf("Expected map with one element 'one': 'one', got %v", result)
	}

	// Test case: Slice with multiple elements
	slice = []string{"one", "two", "three"}
	result = Slice2Map(slice, keyFunc)
	expected = map[string]string{"one": "one", "two": "two", "three": "three"}
	if !SliceEquals(result, expected) {
		t.Errorf("Expected map with elements 'one': 'one', 'two': 'two', 'three': 'three', got %v", result)
	}
}

func TestDowncastMap(t *testing.T) {
	// Test case: Empty map should return an empty map
	m := map[string]interface{}{}
	result := DowncastMap(m)
	if len(result) != 0 {
		t.Errorf("Expected empty map, got %v", result)
	}

	// Test case: Map with one element
	m = map[string]interface{}{"one": "one"}
	result = DowncastMap(m)
	expected := map[string]any{"one": "one"}
	if !SliceEquals(result, expected) {
		t.Errorf("Expected map with one element 'one': 'one', got %v", result)
	}

	// Test case: Map with multiple elements
	m = map[string]interface{}{"one": "one", "two": 2}
	result = DowncastMap(m)
	expected = map[string]any{"one": "one", "two": 2}
	if !SliceEquals(result, expected) {
		t.Errorf("Expected map with elements 'one': 'one', 'two': 2, got %v", result)
	}
}

func TestMergeMap(t *testing.T) {
	// Test case: Empty maps should return an empty map
	bases := []map[string]interface{}{}
	result := MergeMap(bases...)
	if len(result) != 0 {
		t.Errorf("Expected empty map, got %v", result)
	}

	// Test case: Single map
	bases = []map[string]interface{}{map[string]interface{}{"one": "one"}}
	result = MergeMap(bases...)
	expected := map[string]interface{}{"one": "one"}
	if !SliceEquals(result, expected) {
		t.Errorf("Expected map with one element 'one': 'one', got %v", result)
	}

	// Test case: Multiple maps
	bases = []map[string]interface{}{
		map[string]interface{}{"one": "one"},
		map[string]interface{}{"two": 2},
	}
	result = MergeMap(bases...)
	expected = map[string]interface{}{"one": "one", "two": 2}
	if !SliceEquals(result, expected) {
		t.Errorf("Expected map with elements 'one': 'one', 'two': 2, got %v", result)
	}

	// Test case: Overlapping keys
	bases = []map[string]interface{}{
		map[string]interface{}{"one": "one"},
		map[string]interface{}{"one": "two"},
	}
	result = MergeMap(bases...)
	expected = map[string]interface{}{"one": "two"}
	if !SliceEquals(result, expected) {
		t.Errorf("Expected map with one element 'one': 'two', got %v", result)
	}
}

func TestDeepCopyMap(t *testing.T) {
	// Test case: Empty map should return an empty map
	m := map[string]interface{}{}
	result := DeepCopyMap(m)
	if len(result) != 0 {
		t.Errorf("Expected empty map, got %v", result)
	}

	// Test case: Map with one element
	m = map[string]interface{}{"one": "one"}
	result = DeepCopyMap(m)
	expected := map[string]interface{}{"one": "one"}
	if !SliceEquals(result, expected) {
		t.Errorf("Expected map with one element 'one': 'one', got %v", result)
	}

	// Test case: Map with multiple elements
	m = map[string]interface{}{"one": "one", "two": 2}
	result = DeepCopyMap(m)
	expected = map[string]interface{}{"one": "one", "two": 2}
	if !SliceEquals(result, expected) {
		t.Errorf("Expected map with elements 'one': 'one', 'two': 2, got %v", result)
	}

	// Test case: Map with nested maps
	m = map[string]interface{}{
		"one": map[string]interface{}{"two": 2},
	}
	result = DeepCopyMap(m)
	expected := map[string]interface{}{
		"one": map[string]interface{}{"two": 2},
	}
	if !SliceEquals(result, expected) {
		t.Errorf("Expected map with nested map 'one': {'two': 2}, got %v", result)
	}
}

func TestSliceEquals(t *testing.T) {
	// Test case: Empty slices should be equal
	slice1 := []interface{}{}
	slice2 := []interface{}{}
	if !SliceEquals(slice1, slice2) {
		t.Errorf("Expected empty slices to be equal")
	}

	// Test case: Slices with same elements in different order should not be equal
	slice1 = []interface{}{"one", "two"}
	slice2 = []interface{}{"two", "one"}
	if SliceEquals(slice1, slice2) {
		t.Errorf("Expected slices with same elements in different order to not be equal")
	}

	// Test case: Slices with different elements should not be equal
	slice1 = []interface{}{"one", "two"}
	slice2 = []interface{}{"three", "four"}
	if SliceEquals(slice1, slice2) {
		t.Errorf("Expected slices with different elements to not be equal")
	}

	// Test case: Slices with same elements should be equal
	slice1 = []interface{}{"one", "two"}
	slice2 = []interface{}{"one", "two"}
	if !SliceEquals(slice1, slice2) {
		t.Errorf("Expected slices with same elements to be equal")
	}
}

func TestMapEquals(t *testing.T) {
	// Test case: Empty maps should be equal
	map1 := map[string]interface{}{}
	map2 := map[string]interface{}{}
	if !MapEquals(map1, map2) {
		t.Errorf("Expected empty maps to be equal")
	}

	// Test case: Maps with same keys and values in different order should not be equal
	map1 = map[string]interface{}{"one": "one", "two": 2}
	map2 = map[string]interface{}{"two": 2, "one": "one"}
	if MapEquals(map1, map2) {
		t.Errorf("Expected maps with same keys and values in different order to not be equal")
	}

	// Test case: Maps with different keys should not be equal
	map1 = map[string]interface{}{"one": "one", "two": 2}
	map2 = map[string]interface{}{"three": "three"}
	if MapEquals(map1, map2) {
		t.Errorf("Expected maps with different keys to not be equal")
	}

	// Test case: Maps with same keys and values should be equal
	map1 = map[string]interface{}{"one": "one", "two": 2}
	map2 = map[string]interface{}{"one": "one", "two": 2}
	if !MapEquals(map1, map2) {
		t.Errorf("Expected maps with same keys and values to be equal")
	}
}

func TestSetEquals(t *testing.T) {
	// Test case: Empty sets should be equal
	set1 := hashset.New()
	set2 := hashset.New()
	if !SetEquals(set1, set2) {
		t.Errorf("Expected empty sets to be equal")
	}

	// Test case: Sets with same elements in different order should not be equal
	set1 = hashset.NewFromSlice([]interface{}{"one", "two"})
	set2 = hashset.NewFromSlice([]interface{}{"two", "one"})
	if SetEquals(set1, set2) {
		t.Errorf("Expected sets with same elements in different order to not be equal")
	}

	// Test case: Sets with different elements should not be equal
	set1 = hashset.NewFromSlice([]interface{}{"one", "two"})
	set2 = hashset.NewFromSlice([]interface{}{"three", "four"})
	if SetEquals(set1, set2) {
		t.Errorf("Expected sets with different elements to not be equal")
	}

	// Test case: Sets with same elements should be equal
	set1 = hashset.NewFromSlice([]interface{}{"one", "two"})
	set2 = hashset.NewFromSlice([]interface{}{"one", "two"})
	if !SetEquals(set1, set2) {
		t.Errorf("Expected sets with same elements to be equal")
	}
}

func TestDeepCopyMap(t *testing.T) {
	// Test case: Empty map should return an empty map
	m := map[string]interface{}{}
	result := DeepCopyMap(m)
	if len(result) != 0 {
		t.Errorf("Expected empty map, got %v", result)
	}

	// Test case: Map with one element should return a copy of the map
	m = map[string]interface{}{"one": "one"}
	result = DeepCopyMap(m)
	if !reflect.DeepEqual(result, m) {
		t.Errorf("Expected copy of map, got %v", result)
	}

	// Test case: Map with multiple elements should return a copy of the map
	m = map[string]interface{}{"one": "one", "two": 2}
	result = DeepCopyMap(m)
	if !reflect.DeepEqual(result, m) {
		t.Errorf("Expected copy of map, got %v", result)
	}

	// Test case: Map with nested maps should return a deep copy of the map
	m = map[string]interface{}{
		"one": map[string]interface{}{"two": 2},
	}
	result = DeepCopyMap(m)
	if !reflect.DeepEqual(result, m) {
		t.Errorf("Expected deep copy of map, got %v", result)
	}

	// Test case: Map with nested slices should return a deep copy of the map
	m = map[string]interface{}{
		"one": []interface{}{"two", "three"},
	}
	result = DeepCopyMap(m)
	if !reflect.DeepEqual(result, m) {
		t.Errorf("Expected deep copy of map, got %v", result)
	}

	// Test case: Map with nested sets should return a deep copy of the map
	m = map[string]interface{}{
		"one": hashset.NewFromSlice([]interface{}{"two", "three"}),
	}
	result = DeepCopyMap(m)
	if !reflect.DeepEqual(result, m) {
		t.Errorf("Expected deep copy of map, got %v", result)
	}

	// Test case: Map with nested maps and slices should return a deep copy of the map
	m = map[string]interface{}{
		"one": map[string]interface{}{
			"two": []interface{}{"three", "four"},
		},
	}
	result = DeepCopyMap(m)
	if !reflect.DeepEqual(result, m) {
		t.Errorf("Expected deep copy of map, got %v", result)
	}

	// Test case: Map with nested maps and sets should return a deep copy of the map
	m = map[string]interface{}{
		"one": map[string]interface{}{
			"two": hashset.NewFromSlice([]interface{}{"three", "four"}),
		},
	}
	result = DeepCopyMap(m)
	if !reflect.DeepEqual(result, m) {
		t.Errorf("Expected deep copy of map, got %v", result)
	}

	// Test case: Map with nested maps, slices, and sets should return a deep copy of the map
	m = map[string]interface{}{
		"one": map[string]interface{}{
			"two": []interface{}{"three", "four"},
			"five": hashset.NewFromSlice([]interface{}{"six", "seven"}),
		},
	}
	result = DeepCopyMap(m)
	if !reflect.DeepEqual(result, m) {
		t.Errorf("Expected deep copy of map, got %v", result)
	}

	// Test case: Map with nested maps, slices, sets, and other types should return a deep copy of the map
	m = map[string]interface{}{
		"one": map[string]interface{}{
			"two": []interface{}{"three", "four"},
			"five": hashset.NewFromSlice([]interface{}{"six", "seven"}),
			"eight": 8,
			"nine": true,
		},
	}
	result = DeepCopyMap(m)
	if !reflect.DeepEqual(result, m) {
		t.Errorf("Expected deep copy of map, got %v", result)
	}

	// Test case: Map with nested maps, slices, sets, other types, and nil values should return a deep copy of the map
	m = map[string]interface{}{
		"one": map[string]interface{}{
			"two": []interface{}{"three", "four"},
			"five": hashset.NewFromSlice([]interface{}{"six", "seven"}),
			"eight": 8,
			"nine": true,
			"ten": nil,
		},
	}
	result = DeepCopyMap(m)
	if !reflect.DeepEqual(result, m) {
		t.Errorf("Expected deep copy of map, got %v", result)
	}

	// Test case: Map with nested maps, slices, sets, other types, nil values, and empty structures should return a deep copy of the map
	m = map[string]interface{}{
		"one": map[string]interface{}{
			"two": []interface{}{"three", "four"},
			"five": hashset.NewFromSlice([]interface{}{"six", "seven"}),
			"eight": 8,
			"nine": true,
			"ten": nil,
			"eleven": struct{}{},
		},
	}
	result = DeepCopyMap(m)
	if !reflect.DeepEqual(result, m) {
		t.Errorf("Expected deep copy of map, got %v", result)
	}

	// Test case: Map with nested maps, slices, sets, other types, nil values, empty structures, and pointers should return a deep copy of the map
	m = map[string]interface{}{
		"one": map[string]interface{}{
			"two": []interface{}{"three", "four"},
			"five": hashset.NewFromSlice([]interface{}{"six", "seven"}),
			"eight": 8,
			"nine": true,
			"ten": nil,
			"eleven": struct{}{},
			"twelve": &struct{}{},
		},
	}
	result = DeepCopyMap(m)
	if !reflect.DeepEqual(result, m) {
		t.Errorf("Expected deep copy of map, got %v", result)
	}

	// Test case: Map with nested maps, slices, sets, other types, nil values, empty structures, pointers, and functions should return a deep copy of the map
	m = map[string]interface{}{
		"one": map[string]interface{}{
			"two": []interface{}{"three", "four"},
			"five": hashset.NewFromSlice([]interface{}{"six", "seven"}),
			"eight": 8,
			"nine": true,
			"ten": nil,
			"eleven": struct{}{},
			"twelve": &struct{}{},
			"thirteen": func() {},
		},
	}
	result = DeepCopyMap(m)
	if !reflect.DeepEqual(result, m) {
		t.Errorf("Expected deep copy of map, got %v", result)
	}

	// Test case: Map with nested maps, slices, sets, other types, nil values, empty structures, pointers, functions, and interfaces should return a deep copy of the map
	m = map[string]interface{}{
		"one": map[string]interface{}{
			"two": []interface{}{"three", "four"},
			"five": hashset.NewFromSlice([]interface{}{"six", "seven"}),
			"eight": 8,
			"nine": true,
			"ten": nil,
			"eleven": struct{}{},
			"twelve": &struct{}{},
			"thirteen": func() {},
			"fourteen": interface{}(nil),
		},
	}
	result = DeepCopyMap(m)
	if !reflect.DeepEqual(result, m) {
		t.Errorf("Expected deep copy of map, got %v", result)
	}

	// Test case: Map with nested maps, slices, sets, other types, nil values, empty structures, pointers, functions, interfaces, and channels should return a deep copy of the map
	m = map[string]interface{}{
		"one": map[string]interface{}{
			"two": []interface{}{"three", "four"},
			"five": hashset.NewFromSlice([]interface{}{"six", "seven"}),
			"eight": 8,
			"nine": true,
			"ten": nil,
			"eleven": struct{}{},
			"twelve": &struct{}{},
			"thirteen": func() {},
			"fourteen": interface{}(nil),
			"fifteen": make(chan int),
		},
	}
	result = DeepCopyMap(m)
	if !reflect.DeepEqual(result, m) {
		t.Errorf("Expected deep copy of map, got %v", result)
	}

	// Test case: Map with nested maps, slices, sets, other types, nil values, empty structures, pointers, functions, interfaces, channels, and maps should return a deep copy of the map
	m = map[string]interface{}{
		"one": map[string]interface{}{
			"two": []interface{}{"three", "four"},
			"five": hashset.NewFromSlice([]interface{}{"six", "seven"}),
			"eight": 8,
			"nine": true,
			"ten": nil,
			"eleven": struct{}{},
			"twelve": &struct{}{},
			"thirteen": func() {},
			"fourteen": interface{}(nil),
			"fifteen": make(chan int),
			"sixteen": map[string]interface{}{
				"seventeen": 18,
			},
		},
	}
	result = DeepCopyMap(m)
	if !reflect.DeepEqual(result, m) {
		t.Errorf("Expected deep copy of map, got %v", result)
	}

	// Test case: Map with nested maps, slices, sets, other types, nil values, empty structures, pointers, functions, interfaces, channels, maps, and slices should return a deep copy of the map
	m = map[string]interface{}{
		"one": map[string]interface{}{
			"two": []interface{}{"three", "four"},
			"five": hashset.NewFromSlice([]interface{}{"six", "seven"}),
			"eight": 8,
			"nine": true,
			"ten": nil,
			"eleven": struct{}{},
			"twelve": &struct{}{},
			"thirteen": func() {},
			"fourteen": interface{}(nil),
			"fifteen": make(chan int),
			"sixteen": map[string]interface{}{
				"seventeen": 18,
			},
			"eighteen": []interface{}{"nineteen", "twenty"},
		},
	}
	result = DeepCopyMap(m)
	if !reflect.DeepEqual(result, m) {
		t.Errorf("Expected deep copy of map, got %v", result)
	}

	// Test case: Map with nested maps, slices, sets, other types, nil values, empty structures, pointers, functions, interfaces, channels, maps, slices, and sets should return a deep copy of the map
	m = map[string]interface{}{
		"one": map[string]interface{}{
			"two": []interface{}{"three", "four"},
			"five": hashset.NewFromSlice([]interface{}{"six", "seven"}),
			"eight": 8,
			"nine": true,
			"ten": nil,
			"eleven": struct{}{},
			"twelve": &struct{}{},
			"thirteen": func() {},
			"fourteen": interface{}(nil),
			"fifteen": make(chan int),
			"sixteen": map[string]interface{}{
				"seventeen": 18,
			},
			"eighteen": []interface{}{"nineteen", "twenty"},
			"nineteen": hashset.NewFromSlice([]interface{}{"twenty", "twenty-one"}),
		},
	}
	result = DeepCopyMap(m)
	if !reflect.DeepEqual(result, m) {
		t.Errorf("Expected deep copy of map, got %v", result)
	}

	// Test case: Map with nested maps, slices, sets, other types, nil values, empty structures, pointers, functions, interfaces, channels, maps, slices, sets, and arrays should return a deep copy of the map
	m = map[string]interface{}{
		"one": map[string]interface{}{
			"two": []interface{}{"three", "four"},
			"five": hashset.NewFromSlice([]interface{}{"six", "seven"}),
			"eight": 8,
			"nine": true,
			"ten": nil,
			"eleven": struct{}{},
			"twelve": &struct{}{},
			"thirteen": func() {},
			"fourteen": interface{}(nil),
			"fifteen": make(chan int),
			"sixteen": map[string]interface{}{
				"seventeen": 18,
			},
			"eighteen": []interface{}{"nineteen", "twenty"},
			"nineteen": hashset.NewFromSlice([]interface{}{"twenty", "twenty-one"}),
			"twenty": [2]int{23, 24},
		},
	}
	result = DeepCopyMap(m)
	if !reflect.DeepEqual(result, m) {
		t.Errorf("Expected deep copy of map, got %v", result)
	}

	// Test case: Map with nested maps, slices, sets, other types, nil values, empty structures, pointers, functions, interfaces, channels, maps, slices, sets, arrays, and structs should return a deep copy of the map
	m = map[string]interface{}{
		"one": map[string]interface{}{
			"two": []interface{}{"three", "four"},
			"five": hashset.NewFromSlice([]interface{}{"six", "seven"}),
			"eight": 8,
			"nine": true,
			"ten": nil,
			"eleven": struct{}{},
			"twelve": &struct{}{},
			"thirteen": func() {},
			"fourteen": interface{}(nil),
			"fifteen": make(chan int),
			"sixteen": map[string]interface{}{
				"seventeen": 18,
			},
			"eighteen": []interface{}{"nineteen", "twenty"},
			"nineteen": hashset.NewFromSlice([]interface{}{"twenty", "twenty-one"}),
			"twenty": [2]int{23, 24},
			"twenty-one": struct {
				TwentyTwo int
			}{
				TwentyTwo: 25,
			},
		},
	}
	result = DeepCopyMap(m)
	if !reflect.DeepEqual(result, m) {
		t.Errorf("Expected deep copy of map, got %v", result)
	}

	// Test case: Map with nested maps, slices, sets, other types, nil values, empty structures, pointers, functions, interfaces, channels, maps, slices, sets, arrays, structs, and pointers should return a deep copy of the map
	m = map[string]interface{}{
		"one": map[string]interface{}{
			"two": []interface{}{"three", "four"},
			"five": hashset.NewFromSlice([]interface{}{"six", "seven"}),
			"eight": 8,
			"nine": true,
			"ten": nil,
			"eleven": struct{}{},
			"twelve": &struct{}{},
			"thirteen": func() {},
			"fourteen": interface{}(nil),
			"fifteen": make(chan int),
			"sixteen": map[string]interface{}{
				"seventeen": 18,
			},
			"eighteen": []interface{}{"nineteen", "twenty"},
			"nineteen": hashset.NewFromSlice([]interface{}{"twenty", "twenty-one"}),
			"twenty": [2]int{23, 24},
			"twenty-one": struct {
				TwentyTwo int
			}{
				TwentyTwo: 25,
			},
			"twenty-two": &struct {
				TwentyThree int
			}{
				TwentyThree: 26,
			},
		},
	}
	result = DeepCopyMap(m)
	if !reflect.DeepEqual(result, m) {
		t.Errorf("Expected deep copy of map, got %v", result)
	}

	// Test case: Map with nested maps, slices, sets, other types, nil values, empty structures, pointers, functions, interfaces, channels, maps, slices, sets, arrays, structs, pointers, and functions should return a deep copy of the map
	m = map[string]interface{}{
		"one": map[string]interface{}{
			"two": []interface{}{"three", "four"},
			"five": hashset.NewFromSlice([]interface{}{"six", "seven"}),
			"eight": 8,
			"nine": true,
			"ten": nil,
			"eleven": struct{}{},
			"twelve": &struct{}{},
			"thirteen": func() {},
			"fourteen": interface{}(nil),
			"fifteen": make(chan int),
			"sixteen": map[string]interface{}{
				"seventeen": 18,
			},
			"eighteen": []interface{}{"nineteen", "twenty"},
			"nineteen": hashset.NewFromSlice([]interface{}{"twenty", "twenty-one"}),
			"twenty": [2]int{23, 24},
			"twenty-one": struct {
				TwentyTwo int
			}{
				TwentyTwo: 25,
			},
			"twenty-two": &struct {
				TwentyThree int
			}{
				TwentyThree: 26,
			},
			"twenty-three": func() {},
		},
	}
	result = DeepCopyMap(m)
	if !reflect.DeepEqual(result, m) {
		t.Errorf("Expected deep copy of map, got %v", result)
	}

	// Test case: Map with nested maps, slices, sets, other types, nil values, empty structures, pointers, functions, interfaces, channels, maps, slices, sets, arrays, structs, pointers, functions, and interfaces should return a deep copy of the map
	m = map[string]interface{}{
		"one": map[string]interface{}{
			"two": []interface{}{"three", "four"},
			"five": hashset.NewFromSlice([]interface{}{"six", "seven"}),
			"eight": 8,
			"nine": true,
			"ten": nil,
			"eleven": struct{}{},
			"twelve": &struct{}{},
			"thirteen": func() {},
			"fourteen": interface{}(nil),
			"fifteen": make(chan int),
			"sixteen": map[string]interface{}{
				"seventeen": 18,
			},
			"eighteen": []interface{}{"nineteen", "twenty"},
			"nineteen": hashset.NewFromSlice([]interface{}{"twenty", "twenty-one"}),
			"twenty": [2]int{23, 24},
			"twenty-one": struct {
				TwentyTwo int
			}{
				TwentyTwo: 25,
			},
