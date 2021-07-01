package main

import (
	"bytes"
	"fmt"
	"net"
	"testing"
)

func TestNETServer_Request(t *testing.T) {
	tt := []struct {
		test    string
		payload []byte
		want    []byte
	}{
		{
			"Sending a simple request returns result",
			[]byte("9\n"),
			[]byte("Square of 9 is 81\n"),
		},
		{
			"Sending another simple request works",
			[]byte("7\n"),
			[]byte("Square of 7 is 49\n"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.test, func(t *testing.T) {
			conn, err := net.Dial("tcp", ":8080")
			if err != nil {
				t.Error("could not connect to TCP server: ", err)
			}
			defer conn.Close()

			if _, err := conn.Write(tc.payload); err != nil {
				t.Error("could not write payload to TCP server:", err)
			}
			out2 := make([]byte, 1024)
			if _, err := conn.Read(out2); err == nil {
				fmt.Println(string(out2))
			}
			out := make([]byte, 1024)
			if _, err := conn.Read(out); err == nil {
				fmt.Println(string(out))
				if bytes.Compare(out, tc.want) == 0 {
					t.Error("response did match expected output")
				}
			} else {
				t.Error("could not read from connection")
			}
		})
	}
}