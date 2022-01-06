package cyclopes

import (
	"testing"
)

func TestStart(t *testing.T) {
	type args struct {
		configPath string
	}
	tests := []struct {
		name      string
		args      args
		wantPanic bool
	}{
		{name: "Test with panic", args: args{configPath: "../../example-configs/panic.yml"}, wantPanic: true},
		{name: "Test with panic", args: args{configPath: "../../example-configs/pass.yml"}, wantPanic: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("Paniced when wantPanic = %v", tt.wantPanic)
				}
			}()
			Start(tt.args.configPath)
		})
	}
}
