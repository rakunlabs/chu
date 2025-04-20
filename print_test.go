package chu

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestPrint(t *testing.T) {
	type args struct {
		v any
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				v: struct {
					Name   string `cfg:"name"`
					Age    int    `cfg:"age"`
					Secret string `log:"-"`
				}{
					Name: "test",
					Age:  18,
				},
			},
			want:    `{"name":"test","age":18}`,
			wantErr: false,
		},
		{
			name: "test internal",
			args: args{
				v: struct {
					Name     string `cfg:"name"`
					Age      int    `cfg:"age"`
					Secret   string `log:"-"`
					Internal map[string]struct {
						APIKey string `cfg:"api_key" log:"false"`
						Test   string `cfg:"test"`
					} `cfg:"internal"`
					Func    func()        `cfg:"func"`
					Complex complex128    `cfg:"complex"`
					Time    time.Duration `cfg:"time"`
				}{
					Name: "test",
					Age:  18,
					Internal: map[string]struct {
						APIKey string `cfg:"api_key" log:"false"`
						Test   string `cfg:"test"`
					}{
						"test": {
							APIKey: "123456",
							Test:   "test",
						},
					},
					Func:    func() {},
					Complex: complex(1, 2),
					Time:    time.Hour * 2,
				},
			},
			want:    `{"name":"test","age":18,"internal":{"test":{"test":"test"}},"time":"2h0m0s"}`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PrintE(t.Context(), tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("Print() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			var gotMap, wantMap map[string]any
			if err := json.Unmarshal([]byte(got), &gotMap); err != nil {
				t.Errorf("Failed to unmarshal got: %v", err)
				return
			}
			if err := json.Unmarshal([]byte(tt.want), &wantMap); err != nil {
				t.Errorf("Failed to unmarshal want: %v", err)
				return
			}
			if !reflect.DeepEqual(gotMap, wantMap) {
				t.Errorf("Print() = %v, want %v", gotMap, wantMap)
			}
		})
	}
}
