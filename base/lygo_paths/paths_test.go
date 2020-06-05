package lygo_paths

import (
	"fmt"
	"testing"
)

func TestFileName(t *testing.T) {
	urls := []string{"http://domain.com/dir/file.html", "http://domain.com/"}
	for _, url := range urls {
		name := FileName(url, true)
		fmt.Println(name)
	}
}

func TestAbsolute(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Absolute(tt.args.path); got != tt.want {
				t.Errorf("Absolute() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExists(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Exists(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Exists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Exists() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtension(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Extension(tt.args.path); got != tt.want {
				t.Errorf("Extension() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtensionName(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtensionName(tt.args.path); got != tt.want {
				t.Errorf("ExtensionName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetWorkspacePath(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetWorkspacePath(); got != tt.want {
				t.Errorf("GetWorkspacePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsDir(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsDir(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsDir() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsFile(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsFile(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsSymLink(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsSymLink(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsSymLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsSymLink() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMkdir(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test1", args{"./dir1/dir2/file.log"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Mkdir(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("Mkdir() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetWorkspacePath(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func TestTmpFile(t *testing.T) {
	type args struct {
		extension string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TmpFile(tt.args.extension); got != tt.want {
				t.Errorf("TmpFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorkspacePath(t *testing.T) {
	type args struct {
		partial string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"aaa", args{"path1/path2/file.txt"}, "aa"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WorkspacePath(tt.args.partial); got != tt.want {
				t.Errorf("WorkspacePath() = %v, want %v", got, tt.want)
			}
		})
	}
}
