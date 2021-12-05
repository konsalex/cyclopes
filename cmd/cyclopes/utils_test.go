package cyclopes

import (
	"testing"
)

var UserDefinedURL = "https://www.visualeyes.design"

var defaultVisual = VisualTesting{}
var remote = VisualTesting{RemoteURL: "https://www.visualeyes.design"}
var Default = Configuration{VisualTesting: defaultVisual}
var RemoteURL = Configuration{VisualTesting: remote}

func TestExtractServerURsL(t *testing.T) {

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
		{name: "Defined server url", want: UserDefinedURL, args: args{config: &RemoteURL}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.config.ExtractServerURL(); got != tt.want {
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
