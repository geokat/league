package main

import (
	"testing"
)

func TestHandleEcho(t *testing.T) {
	tests := []formFileTestCase{
		{
			"smoke-test",
			[]byte("1,2,3\n4,5,6\n7,8,9"),
			"1,2,3\n4,5,6\n7,8,9\n",
			200,
		},
		{
			"empty-csv",
			[]byte{},
			"\n",
			200,
		},
		{
			"one-element-csv",
			[]byte("0"),
			"0\n",
			200,
		},
		{
			"only-empty-lines",
			[]byte("\n\n\n\n\r\n"),
			"\n",
			200,
		},
		{
			"extraneous-whitespace",
			[]byte("\t\u20011, 2, -3\n\t\v \t\u0085\t1, -2,\t3\n1,2,-3\n"),
			"1,2,-3\n1,-2,3\n1,2,-3\n",
			200,
		},
		{
			"non-integer-literals",
			[]byte("1, 2, -3.234\n1,-2.121,3\n1.983,2,-3\n"),
			"Error: parsing CSV: record on line 1: parsing \"-3.234\": invalid syntax\n",
			400,
		},
		{
			"non-numeric-literals",
			[]byte("1&fl-, 2,3\n1fl-, 2,3\n1fl-,2,3\n"),
			"Error: parsing CSV: record on line 1: parsing \"1&fl-\": invalid syntax\n",
			400,
		},
		{
			"only-commas",
			[]byte(",,\n,,\n,,\n"),
			"Error: parsing CSV: record on line 1: parsing \"\": invalid syntax\n",
			400,
		},
		{
			"empty-lines",
			[]byte("1,2,-3\n\n\n1,-2,3\n\n1,2,-3\n"),
			"1,2,-3\n1,-2,3\n1,2,-3\n",
			200,
		},
		{
			"ensure-trailing-new-line",
			[]byte("1,2,-3\n\n\n1,-2,3\n\n1,2,-3"),
			"1,2,-3\n1,-2,3\n1,2,-3\n",
			200,
		},
		{
			"ensure-only-one-trailing-new-line",
			[]byte("1,2,-3\n\n\n1,-2,3\n\n1,2,-3\n\n\n\n"),
			"1,2,-3\n1,-2,3\n1,2,-3\n",
			200,
		},
	}

	h := webApiMiddleware(handleEcho)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runFormFileTestCase(t, h, tt.payload, tt.wantBody, tt.wantStatus)
		})
	}
}

func TestHandleFlatten(t *testing.T) {
	tests := []formFileTestCase{
		{
			"smoke-test",
			[]byte("1,2,3\n4,5,6\n7,8,9"),
			"1,2,3,4,5,6,7,8,9\n",
			200,
		},
		{
			"empty-csv",
			[]byte{},
			"\n",
			200,
		},
		{
			"one-element-csv",
			[]byte("0"),
			"0\n",
			200,
		},
		{
			"only-empty-lines",
			[]byte("\n\n\n\n\r\n"),
			"\n",
			200,
		},
		{
			"extraneous-whitespace",
			[]byte("\t\u20011, 2, -3\n\t\v \t\u0085\t1, -2,\t3\n1,2,-3\n"),
			"1,2,-3,1,-2,3,1,2,-3\n",
			200,
		},
		{
			"non-integer-literals",
			[]byte("1, 2, -3.234\n1,-2.121,3\n1.983,2,-3\n"),
			"Error: parsing CSV: record on line 1: parsing \"-3.234\": invalid syntax\n",
			400,
		},
		{
			"non-numeric-literals",
			[]byte("1&fl-, 2,3\n1fl-, 2,3\n1fl-,2,3\n"),
			"Error: parsing CSV: record on line 1: parsing \"1&fl-\": invalid syntax\n",
			400,
		},
		{
			"only-commas",
			[]byte(",,\n,,\n,,\n"),
			"Error: parsing CSV: record on line 1: parsing \"\": invalid syntax\n",
			400,
		},
		{
			"empty-lines",
			[]byte("1,2,-3\n\n\n1,-2,3\n\n1,2,-3\n"),
			"1,2,-3,1,-2,3,1,2,-3\n",
			200,
		},
		{
			"ensure-trailing-new-line",
			[]byte("1,2,-3\n\n\n1,-2,3\n\n1,2,-3"),
			"1,2,-3,1,-2,3,1,2,-3\n",
			200,
		},
		{
			"ensure-only-one-trailing-new-line",
			[]byte("1,2,-3\n\n\n1,-2,3\n\n1,2,-3\n\n\n\n"),
			"1,2,-3,1,-2,3,1,2,-3\n",
			200,
		},
	}

	h := webApiMiddleware(handleFlatten)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runFormFileTestCase(t, h, tt.payload, tt.wantBody, tt.wantStatus)
		})
	}
}

func TestHandleSum(t *testing.T) {
	tests := []formFileTestCase{
		{
			"smoke-test",
			[]byte("1,2,3\n4,5,6\n7,8,9"),
			"45\n",
			200,
		},
		{
			"large-integers",
			[]byte("-12345678901234567890,0,0\n0,12345678901234567890,0\n0,0,12345678901234567890"),
			"12345678901234567890\n",
			200,
		},
		{
			"empty-csv",
			[]byte{},
			"0\n",
			200,
		},
		{
			"one-element-csv",
			[]byte("-1"),
			"-1\n",
			200,
		},
		{
			"one-zero-csv",
			[]byte("0"),
			"0\n",
			200,
		},
		{
			"only-empty-lines",
			[]byte("\n\n\n\n\r\n"),
			"0\n",
			200,
		},
		{
			"extraneous-whitespace",
			[]byte("\t\u20011, 2, -3\n\t\v \t\u0085\t1, -2,\t3\n1,2,-3\n"),
			"2\n",
			200,
		},
		{
			"non-integer-literals",
			[]byte("1, 2, -3.234\n1,-2.121,3\n1.983,2,-3\n"),
			"Error: parsing CSV: record on line 1: parsing \"-3.234\": invalid syntax\n",
			400,
		},
		{
			"non-numeric-literals",
			[]byte("1&fl-, 2,3\n1fl-, 2,3\n1fl-,2,3\n"),
			"Error: parsing CSV: record on line 1: parsing \"1&fl-\": invalid syntax\n",
			400,
		},
		{
			"only-commas",
			[]byte(",,\n,,\n,,\n"),
			"Error: parsing CSV: record on line 1: parsing \"\": invalid syntax\n",
			400,
		},
		{
			"empty-lines",
			[]byte("1,2,-3\n\n\n1,-2,3\n\n1,2,-3\n"),
			"2\n",
			200,
		},
		{
			"ensure-trailing-new-line",
			[]byte("1,2,-3\n\n\n1,-2,3\n\n1,2,-3"),
			"2\n",
			200,
		},
		{
			"ensure-only-one-trailing-new-line",
			[]byte("1,2,-3\n\n\n1,-2,3\n\n1,2,-3\n\n\n\n"),
			"2\n",
			200,
		},
	}

	h := webApiMiddleware(handleSum)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runFormFileTestCase(t, h, tt.payload, tt.wantBody, tt.wantStatus)
		})
	}
}

func TestHandleMultiply(t *testing.T) {
	tests := []formFileTestCase{
		{
			"smoke-test",
			[]byte("1,2,3\n4,5,6\n7,8,9"),
			"362880\n",
			200,
		},
		{
			"large-integers",
			[]byte("-12345678901234567890,1,1\n1,1,1\n1,1,1"),
			"-12345678901234567890\n",
			200,
		},
		{
			"empty-csv",
			[]byte{},
			"0\n",
			200,
		},
		{
			"one-element-csv",
			[]byte("-1"),
			"-1\n",
			200,
		},
		{
			"one-zero-csv",
			[]byte("0"),
			"0\n",
			200,
		},
		{
			"only-empty-lines",
			[]byte("\n\n\n\n\r\n"),
			"0\n",
			200,
		},
		{
			"extraneous-whitespace",
			[]byte("\t\u20011, 2, -3\n\t\v \t\u0085\t1, -2,\t3\n1,2,-3\n"),
			"-216\n",
			200,
		},
		{
			"non-integer-literals",
			[]byte("1, 2, -3.234\n1,-2.121,3\n1.983,2,-3\n"),
			"Error: parsing CSV: record on line 1: parsing \"-3.234\": invalid syntax\n",
			400,
		},
		{
			"non-numeric-literals",
			[]byte("1&fl-, 2,3\n1fl-, 2,3\n1fl-,2,3\n"),
			"Error: parsing CSV: record on line 1: parsing \"1&fl-\": invalid syntax\n",
			400,
		},
		{
			"only-commas",
			[]byte(",,\n,,\n,,\n"),
			"Error: parsing CSV: record on line 1: parsing \"\": invalid syntax\n",
			400,
		},
		{
			"empty-lines",
			[]byte("1,2,-3\n\n\n1,-2,3\n\n1,2,-3\n"),
			"-216\n",
			200,
		},
		{
			"ensure-trailing-new-line",
			[]byte("1,2,-3\n\n\n1,-2,3\n\n1,2,-3"),
			"-216\n",
			200,
		},
		{
			"ensure-only-one-trailing-new-line",
			[]byte("1,2,-3\n\n\n1,-2,3\n\n1,2,-3\n\n\n\n"),
			"-216\n",
			200,
		},
	}

	h := webApiMiddleware(handleMultiply)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runFormFileTestCase(t, h, tt.payload, tt.wantBody, tt.wantStatus)
		})
	}
}

func TestHandleInvert(t *testing.T) {
	tests := []formFileTestCase{
		{
			"smoke-test",
			[]byte("1,2,3\n4,5,6\n7,8,9"),
			"1,4,7\n2,5,8\n3,6,9\n",
			200,
		},
		{
			"empty-csv",
			[]byte{},
			"\n",
			200,
		},
		{
			"one-element-csv",
			[]byte("0"),
			"0\n",
			200,
		},
		{
			"only-empty-lines",
			[]byte("\n\n\n\n\r\n"),
			"\n",
			200,
		},
		{
			"extraneous-whitespace",
			[]byte("\t\u20011, 2, -3\n\t\v \t\u0085\t1, -2,\t3\n1,2,-3\n"),
			"1,1,1\n2,-2,2\n-3,3,-3\n",
			200,
		},
		{
			"non-integer-literals",
			[]byte("1, 2, -3.234\n1,-2.121,3\n1.983,2,-3\n"),
			"Error: parsing CSV: record on line 1: parsing \"-3.234\": invalid syntax\n",
			400,
		},
		{
			"non-numeric-literals",
			[]byte("1&fl-, 2,3\n1fl-, 2,3\n1fl-,2,3\n"),
			"Error: parsing CSV: record on line 1: parsing \"1&fl-\": invalid syntax\n",
			400,
		},
		{
			"only-commas",
			[]byte(",,\n,,\n,,\n"),
			"Error: parsing CSV: record on line 1: parsing \"\": invalid syntax\n",
			400,
		},
		{
			"empty-lines",
			[]byte("1,2,-3\n\n\n1,-2,3\n\n1,2,-3\n"),
			"1,1,1\n2,-2,2\n-3,3,-3\n",
			200,
		},
		{
			"ensure-trailing-new-line",
			[]byte("1,2,-3\n\n\n1,-2,3\n\n1,2,-3"),
			"1,1,1\n2,-2,2\n-3,3,-3\n",
			200,
		},
		{
			"ensure-only-one-trailing-new-line",
			[]byte("1,2,-3\n\n\n1,-2,3\n\n1,2,-3\n\n\n\n"),
			"1,1,1\n2,-2,2\n-3,3,-3\n",
			200,
		},
	}

	h := webApiMiddleware(handleInvert)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runFormFileTestCase(t, h, tt.payload, tt.wantBody, tt.wantStatus)
		})
	}
}
