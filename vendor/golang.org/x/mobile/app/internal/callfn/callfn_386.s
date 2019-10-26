// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"
#include "funcdata.h"

TEXT ·CallFn(SB),$0-4
	MOVL fn+0(FP), AX
	CALL AX
	RET
