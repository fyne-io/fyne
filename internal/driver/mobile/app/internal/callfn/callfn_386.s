#include "textflag.h"
#include "funcdata.h"

TEXT ·CallFn(SB),$0-4
	MOVL fn+0(FP), AX
	CALL AX
	RET
