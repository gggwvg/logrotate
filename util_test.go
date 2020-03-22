package logrotate

import "testing"

func Test_stringToBytes(t *testing.T) {
	cases := []struct {
		Str  string
		Want int64
	}{
		{Str: "", Want: 0},
		{Str: "200", Want: 200},
		{Str: "400KB", Want: 409600},
		{Str: "500mb", Want: 524288000},
		{Str: "1g", Want: 1073741824},
	}
	for _, c := range cases {
		t.Run(c.Str, func(t *testing.T) {
			bytes, err := stringToBytes(c.Str)
			if err != nil {
				t.FailNow()
			}
			if bytes != c.Want {
				t.Errorf("case:%s failed, want:%d, get:%d", c.Str, c.Want, bytes)
			}
		})
	}
}

func Test_splitFilename(t *testing.T) {
	cases := []struct {
		Name       string
		WantPrefix string
		WantExt    string
	}{
		{Name: "tmp.log", WantPrefix: "tmp", WantExt: ".log"},
		{Name: "/tmp/tmp.log", WantPrefix: "tmp", WantExt: ".log"},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			prefix, ext := splitFilename(c.Name)
			if prefix != c.WantPrefix {
				t.Errorf("wrong prefix:%s, want:%s", prefix, c.WantPrefix)
			}
			if ext != c.WantExt {
				t.Errorf("wrong ext:%s, want:%s", ext, c.WantExt)
			}
		})
	}
}
