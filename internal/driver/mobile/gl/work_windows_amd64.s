// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

TEXT ·fixFloat(SB),NOSPLIT,$0-32
	MOVQ 	x0+0(FP), X0
	MOVQ	x1+8(FP), X1
	MOVQ	x2+16(FP), X2
	MOVQ	x3+24(FP), X3
	RET
