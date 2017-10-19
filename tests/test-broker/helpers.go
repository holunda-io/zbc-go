package testbroker

import (
	"runtime/debug"
	"reflect"
	"testing"
	"math/rand"
	"time"
	"errors"
)

var errClientStartFailed = errors.New("cannot connect to the broker")

const topicName = "default-topic"
const brokerAddr = "0.0.0.0:51015"

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")


func assert(t *testing.T, exp, got interface{}, equal bool) {
	if reflect.DeepEqual(exp, got) != equal {
		debug.PrintStack()
		t.Fatalf("Expecting '%v' got '%v'\n", exp, got)
	}
}

func RandStringBytes(n int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}