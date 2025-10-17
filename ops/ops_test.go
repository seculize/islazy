package ops

import (
	"testing"
)

func TestTernary(t *testing.T) {
	tests := []struct {
		name   string
		cond   bool
		posVal interface{}
		negVal interface{}
		want   interface{}
	}{
		{"true condition with strings", true, "yes", "no", "yes"},
		{"false condition with strings", false, "yes", "no", "no"},
		{"true condition with ints", true, 1, 0, 1},
		{"false condition with ints", false, 1, 0, 0},
		{"true condition with nil", true, "value", nil, "value"},
		{"false condition with nil", false, nil, "value", "value"},
		{"true condition both nil", true, nil, nil, nil},
		{"false condition both nil", false, nil, nil, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Ternary(tt.cond, tt.posVal, tt.negVal); got != tt.want {
				t.Errorf("Ternary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTernaryWithComplexTypes(t *testing.T) {
	// Test with slices
	slice1 := []int{1, 2, 3}
	slice2 := []int{4, 5, 6}
	result := Ternary(true, slice1, slice2).([]int)
	if len(result) != len(slice1) {
		t.Errorf("Ternary with slices failed")
	}

	// Test with structs
	type TestStruct struct {
		Value int
	}
	struct1 := TestStruct{Value: 1}
	struct2 := TestStruct{Value: 2}
	result2 := Ternary(false, struct1, struct2).(TestStruct)
	if result2.Value != 2 {
		t.Errorf("Ternary with structs failed, got %v, want %v", result2.Value, 2)
	}
}
