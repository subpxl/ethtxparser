// Package utils provides utility functions for common type conversions

package utils

import (
    "strconv"
)

// String2Int64 converts a string to an int64 using the specified base.
func String2Int64(s string, base int) (int64, error) {
    return strconv.ParseInt(s, base, 64)
}