//go:build linux && !cuda

package wrapper

/*
#cgo CFLAGS: -std=c11
#cgo CXXFLAGS: -std=c++17
#cgo CFLAGS: -I${SRCDIR}/../core/include
#cgo CXXFLAGS: -I${SRCDIR}/../core/include
#cgo LDFLAGS: -L${SRCDIR}/../build/lib -lllama_core -lcommon -lllama -lwhisper -lwhisper-common -lmtmd -lggml -lggml-base -lggml-cpu -lstdc++ -lm
#include <stdlib.h>
#include "core.h"
*/
import "C"
