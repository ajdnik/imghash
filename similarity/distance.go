// Copyright 2020 Rok Ajdnik. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// Package similarity implements data types and methods used
// to calculate similarities between hashes.
package similarity

import (
	"math"
)

// Distance represents a similarity measure as float64 value.
type Distance float64

// Equal checks if two distances are the same, it uses the epsilon approach for comparing floats.
func (d Distance) Equal(dst Distance) bool {
	eps := 1e-12
	return math.Abs(float64(d)-float64(dst)) <= eps
}
