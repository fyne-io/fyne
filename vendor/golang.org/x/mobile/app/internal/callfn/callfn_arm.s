// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"
#include "funcdata.h"

TEXT ·CallFn(SB),$0-4
	MOVW fn+0(FP), R0
	BL (R0)
	RET
