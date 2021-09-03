#include "textflag.h"
#include "funcdata.h"

TEXT ·CallFn(SB),$0-8
	MOVQ fn+0(FP), AX
	CALL AX
	RET
