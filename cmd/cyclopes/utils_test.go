package cyclopes

import (
	"fmt"
	"testing"
)

var UserDefinedURL = "http://localhost:8080"

var Default = Configuration{Server: true}
var UserServer = Configuration{Server: false, ServerURL: UserDefinedURL}
var Fail = Configuration{Server: false}

func TestExtractServerURL(t *testing.T) {

	t.Parallel()

	type args struct {
		config *Configuration
	}
	tests := []struct {
		name  string
		args  args
		want  string
		error bool
	}{
		{name: "Default server url", want: DEFAULT_URL, args: args{config: &Default}},
		{name: "Defined server url", want: UserDefinedURL, args: args{config: &UserServer}},
		{name: "Defined server url", error: true, args: args{config: &Fail}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := tt.args.config.ExtractServerURL(); got != tt.want || (err == nil) && tt.error {
				t.Errorf("ExtractServerURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConstructServerURL(t *testing.T) {

	t.Parallel()

	type args struct {
		rawURL string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "Localhost without scheme", want: "http://localhost:3000", args: args{rawURL: "localhost:3000"}},
		{name: "Localhost without scheme and trailing slash", want: "http://localhost:3000", args: args{rawURL: "localhost:3000/"}},
		{name: "Localhost with scheme", want: "http://localhost:3000", args: args{rawURL: "http://localhost:3000"}},
		{name: "Local url 'server'", want: "http://server:3000", args: args{rawURL: "server:3000"}},
		{name: "Remote url", want: "https://app.visualeyes.design", args: args{rawURL: "https://app.visualeyes.design"}},
		{name: "Error", wantErr: true, args: args{rawURL: "-randomURL?!0"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConstructServerURL(tt.args.rawURL)
			fmt.Println(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConstructServerURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ConstructServerURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
