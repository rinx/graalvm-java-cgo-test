package main

// #cgo CFLAGS: -I${SRCDIR}/../target/native
// #cgo LDFLAGS: -L${SRCDIR}/../target/native -ljavacgo
//
// #include <stdlib.h>
// #include <libjavacgo.h>
import "C"
import (
	"fmt"
	"strconv"
	"sync"
	"unsafe"
)

type javaCgo struct {
	isolate *C.graal_isolate_t
}

type JavaCgo interface {
	Str(s string) (string, error)
}

func New() (JavaCgo, error) {
	var isolate *C.graal_isolate_t
	var thread *C.graal_isolatethread_t

	param := &C.graal_create_isolate_params_t{
		reserved_address_space_size: 1024 * 1024 * 500,
	}

	if C.graal_create_isolate(param, &isolate, &thread) != 0 {
		return nil, fmt.Errorf("failed to initialize")
	}

	return &javaCgo{
		isolate: isolate,
	}, nil
}

func (j *javaCgo) attachThread() (*C.graal_isolatethread_t, error) {
	thread := C.graal_get_current_thread(j.isolate)
	if thread != nil {
		return thread, nil
	}

	var newThread *C.graal_isolatethread_t
	if C.graal_attach_thread(j.isolate, &newThread) != 0 {
		return nil, fmt.Errorf("failed to attach thread")
	}

	return newThread, nil
}

func (j *javaCgo) Str(s string) (string, error) {
	thread, err := j.attachThread()
	if err != nil {
		return "", err
	}

	cstr := C.CString(s)
	defer C.free(unsafe.Pointer(cstr))

	return C.GoString(C.java_cgo_str(thread, cstr)), nil
}

func main() {
	javaCgo, err := New()
	if err != nil {
		println(err)
		return
	}

	var wg sync.WaitGroup

	n := 1000

	for t := 0; t < 100; t++ {

		wg.Add(1)
		go func(t int) {
			defer wg.Done()

			for i := 0; i < n; i++ {
				key := "key" + strconv.Itoa(i)
				val, err := javaCgo.Str(key)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("thread " + strconv.Itoa(t) + ": " + val)
			}
		}(t)
	}

	wg.Wait()

	fmt.Println("parallel run finished")
}
