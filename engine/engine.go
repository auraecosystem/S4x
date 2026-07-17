package engine

/*
#cgo LDFLAGS: -L../build/native -laurabridge
#include <stdint.h>

// Explicit C-ABI signature matching our exported Zig/Assembly library symbol
int32_t aura_process_payload(const char* in_ptr, char* out_ptr, uint64_t length);
*/
import "C"
import (
	"errors"
	"runtime"
	"unsafe"
)

var (
	ErrEmptyBuffer       = errors.New("cannot process an empty byte vector slice")
	ErrExecutionFailure  = errors.New("underlying assembly execution pathway returned non-zero error")
)

// ProcessContainerData bridges the Go runtime core down to the 21.7% Assembly layer
func ProcessContainerData(inputSlice []byte) ([]byte, error) {
	// 1. Boundary Safety Check
	if len(inputSlice) == 0 {
		return nil, ErrEmptyBuffer
	}

	// 2. Allocate an identical output slice buffer block
	outputSlice := make([]byte, len(inputSlice))

	// 3. Instantiate a Go Runtime Pinner object 
	// This forces the Garbage Collector to freeze these memory addresses in place
	var pinner runtime.Pinner
	pinner.Pin(&inputSlice[0])
	pinner.Pin(&outputSlice[0])
	defer pinner.Unpin() // Always unpin memory blocks immediately after execution completes

	// 4. Extract raw, unsafe memory addresses to comply with C-ABI rules
	inputPointer  := (*C.char)(unsafe.Pointer(&inputSlice[0]))
	outputPointer := (*C.char)(unsafe.Pointer(&outputSlice[0]))
	bufferLength  := (C.uint64_t)(len(inputSlice))

	// 5. Cross the execution boundary (registers are loaded: RDI=input, RSI=output, RDX=length)
	statusCode := C.aura_process_payload(inputPointer, outputPointer, bufferLength)

	// 6. Evaluate native return code statuses
	if statusCode != 0 {
		return nil, ErrExecutionFailure
	}

	return outputSlice, nil
}
