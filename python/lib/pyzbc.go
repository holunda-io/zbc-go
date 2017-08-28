package main

// #cgo pkg-config: python3
/*
#define Py_LIMITED_API
#include <Python.h>
int PyArg_ParseTuple_LL(PyObject *, long long *, long long *);
int PyArg_ParseTuple_String(PyObject * args, char** a, char** b, char** c);
PyObject* Py_String(char *pystring);
*/
import "C"

import (
	"fmt"
	"github.com/zeebe-io/zbc-go/zbc"
)

//export NewClient
func NewClient(addr string) *C.PyObject {
	client, err := zbc.NewClient(addr)
	if err != nil {

	}
	return C.PyLong_FromLong(client)
}

//export Add
func Add(x, y int) int {
	fmt.Printf("Go says: adding %v and %v\n", x, y)
	return x + y
}

func main() {}