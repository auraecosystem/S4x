package engine

/*
#cgo LDFLAGS: -L../build/native -laurabridge
#include <stdint.h>

// Forward declaration matching our Zig/Assembly exported C-ABI symbol
int32_t aura_process_payload(const char* in_ptr, char* out_ptr, uint64_t length);
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// DispatchOptimizedPayload shifts data down into the C-ABI execution stack
func DispatchOptimizedPayload(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty transactional byte vector")
	}

	outputBuffer := make([]byte, len(data))

	// Pin pointers to avoid Go garbage collector movement during native execution
	cInput := (*C.char)(unsafe.Pointer(&data[0]))
	cOutput := (*C.char)(unsafe.Pointer(&outputBuffer[0]))
	cLength := (C.uint64_t)(len(data))

	// Direct transition through the C-ABI boundary
	resultCode := C.aura_process_payload(cInput, cOutput, cLength)
	if resultCode != 0 {
		return nil, fmt.Errorf("underlying C-ABI execution pathway returned failure: %d", resultCode)
	}

	return outputBuffer, nil
}
