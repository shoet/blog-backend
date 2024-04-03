package options_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/shoet/blog/internal/options"
)

func Test_SetDefault(t *testing.T) {
	type args struct {
		v            any
		fieldName    string
		setValue     *bool
		defaultValue bool
	}
	type wants struct {
		wantStruct any
		err        error
	}

	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name: "デフォルト値",
			args: args{
				v: &struct {
					TargetField bool
				}{},
				fieldName:    "TargetField",
				setValue:     nil,
				defaultValue: false,
			},
			wants: wants{
				wantStruct: &struct {
					TargetField bool
				}{
					TargetField: false,
				},
				err: nil,
			},
		},
		{
			name: "ポインタでない",
			args: args{
				v: struct {
					TargetField bool
				}{},
				fieldName:    "TargetField",
				setValue:     nil,
				defaultValue: false,
			},
			wants: wants{
				wantStruct: nil,
				err:        options.ErrNotPointer,
			},
		},
		{
			name: "フィールドが存在しない",
			args: args{
				v: &struct {
					TargetField bool
				}{},
				fieldName:    "NotFoundFiled",
				setValue:     nil,
				defaultValue: false,
			},
			wants: wants{
				wantStruct: nil,
				err:        options.ErrFieldNotFound,
			},
		},
		{
			name: "セット先の型とデフォルト値の型が異なる",
			args: args{
				v: &struct {
					TargetField string
				}{},
				fieldName:    "TargetField",
				setValue:     nil,
				defaultValue: false,
			},
			wants: wants{
				wantStruct: &struct {
					TargetField bool
				}{
					TargetField: false,
				},
				err: options.ErrDefaultValueUnmatchType,
			},
		},
		{
			name: "フィールドが先頭小文字でエクスポートされていない",
			args: args{
				v: &struct {
					targetField bool
				}{},
				fieldName:    "targetField",
				setValue:     nil,
				defaultValue: false,
			},
			wants: wants{
				wantStruct: nil,
				err:        options.ErrFieldNotFound,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := tt.args.v
			err := options.SetDefault(v, tt.args.fieldName, tt.args.setValue, tt.args.defaultValue)
			if err != tt.wants.err {
				t.Errorf("got: %v, want: %v", err, tt.wants.err)
			}

			if err == nil {
				if diff := cmp.Diff(v, tt.wants.wantStruct); diff != "" {
					t.Errorf("differs: (-got +want)\n%s", diff)
				}
			}
		})
	}
}
