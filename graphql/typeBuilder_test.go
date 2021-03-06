package graphql

import (
	"reflect"
	"testing"
)

type B struct {
	c bool
}

type A struct {
	a int32
	b B
	d string `gql:"dName"`
}

type C struct {
	a int32
	b *B
	d string `gql:"dName"`
}

func TestBuild(t *testing.T) {
	tests := []struct {
		name       string
		t          reflect.Type
		wantResult string
	}{
		{"", reflect.TypeOf((*A)(nil)).Elem(), "type ResultType {\n\ta: Int\n\tb: B\n\tdName: String\n\n}"},
		{"", reflect.TypeOf((*B)(nil)).Elem(), "type ResultType {\n\tc: Boolean\n\n}"},
		{"", reflect.TypeOf((*C)(nil)).Elem(), "type ResultType {\n\ta: Int\n\tb: B\n\tdName: String\n\n}"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := Build(tt.t, "ResultType"); result != tt.wantResult {
				t.Errorf("Check() result = %v, wantResult %v", result, tt.wantResult)
			}
		})
	}
}
