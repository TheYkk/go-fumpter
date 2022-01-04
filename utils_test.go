package main

import "testing"

func Test_randomString(t *testing.T) {

	t.Run("Test random", func(t *testing.T) {
		t.Log(randomString(10))
	})

}
