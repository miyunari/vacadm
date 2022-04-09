package version

import "testing"

func Test_Version(t *testing.T) {
	tt := []struct {
		name           string
		hash           string
		buildtimestamp string
		want           string
	}{
		{
			name:           "expected",
			hash:           "a",
			buildtimestamp: "b",
			want:           "Version: a from b",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			hash = tc.hash
			buildtimestamp = tc.buildtimestamp
			got := Version()
			if got != tc.want {
				t.Fatalf("want: %s, got: %s", tc.want, got)
			}
		})
	}
}
