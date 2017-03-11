package slicelement

import (
	"reflect"
	"testing"
)

func TestIndex_getIndexInt(t *testing.T) {
	type args struct {
		Data    interface{}
		Element interface{}
	}
	tests := []struct {
		Name      string
		T         *indexT
		Args      args
		WantIndex int
		WantErr   bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			index := &indexT{}
			gotIndex, err := index.getIndexInt(tt.Args.Data, tt.Args.Element)
			if (err != nil) != tt.WantErr {
				t.Errorf("Index.getindexTInt() error = %v, WantErr %v", err, tt.WantErr)
				return
			}
			if gotIndex != tt.WantIndex {
				t.Errorf("indexT.getIndexInt() = %v, Want %v", gotIndex, tt.WantIndex)
			}
		})
	}
}

func TestIndex_getIndexString(t *testing.T) {
	type args struct {
		Data    interface{}
		Element interface{}
	}
	tests := []struct {
		Name      string
		T         *indexT
		Args      args
		WantIndex int
		WantErr   bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			index := &indexT{}
			gotIndex, err := index.getIndexString(tt.Args.Data, tt.Args.Element)
			if (err != nil) != tt.WantErr {
				t.Errorf("indexT.getIndexString() error = %v, WantErr %v", err, tt.WantErr)
				return
			}
			if gotIndex != tt.WantIndex {
				t.Errorf("indexT.getIndexString() = %v, Want %v", gotIndex, tt.WantIndex)
			}
		})
	}
}

func TestIndex_getIndexFloat32(t *testing.T) {
	type args struct {
		Data    interface{}
		Element interface{}
	}
	tests := []struct {
		Name      string
		T         *indexT
		Args      args
		WantIndex int
		WantErr   bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			index := &indexT{}
			gotIndex, err := index.getIndexFloat32(tt.Args.Data, tt.Args.Element)
			if (err != nil) != tt.WantErr {
				t.Errorf("indexT.getIndexFloat32() error = %v, WantErr %v", err, tt.WantErr)
				return
			}
			if gotIndex != tt.WantIndex {
				t.Errorf("Index.getindexTFloat32() = %v, Want %v", gotIndex, tt.WantIndex)
			}
		})
	}
}

func TestIndex_getIndexUint(t *testing.T) {
	type Args struct {
		Data    interface{}
		Element interface{}
	}
	tests := []struct {
		Name      string
		T         *indexT
		Args      Args
		WantIndex int
		WantErr   bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			index := &indexT{}
			gotIndex, err := index.getIndexUint(tt.Args.Data, tt.Args.Element)
			if (err != nil) != tt.WantErr {
				t.Errorf("indexT.getIndexUint() error = %v, WantErr %v", err, tt.WantErr)
				return
			}
			if gotIndex != tt.WantIndex {
				t.Errorf("indexT.getIndexUint() = %v, Want %v", gotIndex, tt.WantIndex)
			}
		})
	}
}

func TestIndex_getIndexStruct(t *testing.T) {
	type Args struct {
		Data    interface{}
		Element interface{}
		Tag     string
	}
	tests := []struct {
		Name      string
		T         *indexT
		Args      Args
		WantIndex int
		WantErr   bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			index := &indexT{}
			gotIndex, err := index.getIndexStruct(tt.Args.Data, tt.Args.Element, tt.Args.Tag)
			if (err != nil) != tt.WantErr {
				t.Errorf("indexT.getIndexStruct() error = %v, WantErr %v", err, tt.WantErr)
				return
			}
			if gotIndex != tt.WantIndex {
				t.Errorf("Index.getIndexStruct() = %v, Want %v", gotIndex, tt.WantIndex)
			}
		})
	}
}

func TestIndex_decodeStruct(t *testing.T) {
	type args struct {
		DataVal reflect.Value
		Element interface{}
		Tag     string
	}
	tests := []struct {
		Name        string
		T           *indexT
		Args        args
		WantIsExist bool
		WantErr     bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			index := &indexT{}
			gotIsExist, err := index.decodeStruct(tt.Args.DataVal, tt.Args.Element, tt.Args.Tag)
			if (err != nil) != tt.WantErr {
				t.Errorf("indexT.decodeStruct() error = %v, WantErr %v", err, tt.WantErr)
				return
			}
			if gotIsExist != tt.WantIsExist {
				t.Errorf("indexT.decodeStruct() = %v, Want %v", gotIsExist, tt.WantIsExist)
			}
		})
	}
}
