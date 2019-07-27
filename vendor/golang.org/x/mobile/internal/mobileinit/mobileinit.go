// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package mobileinit contains common initialization logic for mobile platforms
// that is relevant to both all-Go apps and gobind-based apps.
//
// Long-term, some code in this package should consider moving into Go stdlib.
package mobileinit

import "C"
