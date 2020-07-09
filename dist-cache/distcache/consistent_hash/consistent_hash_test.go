package consistent_hash

import (
	"strconv"
	"testing"
)

func TestAdd(t *testing.T) {
	hash := New(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})

	// Given the above hash function, this will give replicas with "hashes":
	// 2, 4, 6, 12, 14, 16, 22, 24, 26
	hash.Add("6", "4", "2")

	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}

	for k, v := range testCases {
		if cv := hash.Get(k); cv != v {
			t.Errorf("Asking for %s, should have yielded %s, but got %s", k, v, cv)
		}
	}

	// Adds 8, 18, 28
	hash.Add("8")

	// 27 should now map to 8.
	testCases["27"] = "8"

	for k, v := range testCases {
		if cv := hash.Get(k); cv != v {
			t.Errorf("Asking for %s, should have yielded %s, but got %s", k, v, cv)
		}
	}
}

func TestRemove(t *testing.T) {
	hash := New(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})

	// Given the above hash function, this will give replicas with "hashes":
	// 2, 4, 6, 8, 12, 14, 16, 18, 22, 24, 26, 28
	hash.Add("8", "6", "4", "2")

	// Remove 4
	hash.Remove("4")

	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "6", // 4 -> 6
		"27": "8",
		"30": "2",
	}

	for k, v := range testCases {
		if cv := hash.Get(k); cv != v {
			t.Errorf("Asking for %s, should have yielded %s, but got %s", k, v, cv)
		}
	}

}