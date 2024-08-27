package core

import "errors"

/*

Let’s take an example with an 8-bit number for clarity:

Suppose typeEncoding = 1011 1101 (which is 189 in decimal).

Right shift by 4 (typeEncoding >> 4):
1011 1101 >> 4  →  0000 1011

Left shift by 4 ((typeEncoding >> 4) << 4):
0000 1011 << 4  →  1011 0000

The result is now 1011 0000, which is 176 in decimal.
The original 4 least significant bits (1101) are now 0000, effectively clearing them.


*/

// calculating first 4 bits
func getType(typeEncoding uint8) uint8 {
	return (typeEncoding >> 4) << 4
}

// calculating last 4 bits
func getEncoding(typeEncoding uint8) uint8 {
	return typeEncoding & 0b00001111
}

func assertType(te uint8, t uint8) error {
	if getType(te) != t {
		return errors.New("the operation is not permitted on this type")
	}
	return nil
}

func assertEncoding(te uint8, e uint8) error {
	if getEncoding(te) != e {
		return errors.New("the operation is not permitted on this encoding")
	}
	return nil
}
