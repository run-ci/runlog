package http

import "testing"

func TestGetTaskID(t *testing.T) {
	errtests := []struct {
		in  string
		val string
		err error
	}{
		{"/log/1/", "1", nil},
		{"/log/1", "1", nil},
		{"/log/", "", errNoTaskID},
		{"/log", "", errNoTaskID},
	}

	for _, test := range errtests {
		t.Run(test.in, func(t *testing.T) {
			val, err := getTaskID(test.in)
			if err != test.err {
				t.Errorf("\n\nexpected error: %v\n\ngot: %v\n\n", test.err, err)
			}
			if val != test.val {
				t.Errorf("\n\nexpected val: %v\n\ngot: %v\n\n", test.val, val)
			}
		})
	}
}
