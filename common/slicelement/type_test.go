package slicelement

import (
	"reflect"
	"testing"
)

func Test_checkInputValid(t *testing.T) {
	type args struct {
		Data    interface{}
		Element interface{}
	}
	tests := []struct {
		Name    string
		Args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if err := checkInputValid(tt.Args.Data, tt.Args.Element); (err != nil) != tt.wantErr {
				t.Errorf("checkInputValid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_checkData(t *testing.T) {
	type args struct {
		Data interface{}
	}
	tests := []struct {
		Name    string
		Args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if err := checkData(tt.Args.Data); (err != nil) != tt.wantErr {
				t.Errorf("checkData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_checkElement(t *testing.T) {
	type args struct {
		Element interface{}
	}
	tests := []struct {
		Name    string
		Args    args
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if err := checkElement(tt.Args.Element); (err != nil) != tt.wantErr {
				t.Errorf("checkElement() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getKindByValue(t *testing.T) {
	type args struct {
		val reflect.Value
	}
	tests := []struct {
		Name     string
		Args     args
		wantKind reflect.Kind
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if gotKind := getKindByValue(tt.Args.val); gotKind != tt.wantKind {
				t.Errorf("getKindByValue() = %v, want %v", gotKind, tt.wantKind)
			}
		})
	}
}

func Test_getSliceUnderlyKind(t *testing.T) {
	type args struct {
		Data interface{}
	}
	tests := []struct {
		Name     string
		Args     args
		wantKind reflect.Kind
		wantErr  bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			gotKind, err := getSliceUnderlyKind(tt.Args.Data)
			if (err != nil) != tt.wantErr {
				t.Errorf("getSliceUnderlyKind() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotKind != tt.wantKind {
				t.Errorf("getSliceUnderlyKind() = %v, want %v", gotKind, tt.wantKind)
			}
		})
	}
}

func TestContains(t *testing.T) {
	type args struct {
		Data    interface{}
		Element interface{}
		tag     string
	}
	tests := []struct {
		Name        string
		Args        args
		wantIsExist bool
		wantErr     bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			gotIsExist, err := Contains(tt.Args.Data, tt.Args.Element, tt.Args.tag)
			if (err != nil) != tt.wantErr {
				t.Errorf("Contains() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotIsExist != tt.wantIsExist {
				t.Errorf("Contains() = %v, want %v", gotIsExist, tt.wantIsExist)
			}
		})
	}
}

func TestGetIndex(t *testing.T) {
	type args struct {
		Data    interface{}
		Element interface{}
		tag     string
	}
	tests := []struct {
		Name      string
		Args      args
		wantIndex int
		wantErr   bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			gotIndex, err := GetIndex(tt.Args.Data, tt.Args.Element, tt.Args.tag)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotIndex != tt.wantIndex {
				t.Errorf("GetIndex() = %v, want %v", gotIndex, tt.wantIndex)
			}
		})
	}
}
