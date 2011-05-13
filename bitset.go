// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	Package bitset implements bitsets.

	It provides methods for making a BitSet of an arbitrary
	upper limit, setting and testing bit locations, and clearing
	bit locations as well as the entire set.

	Example use:

	b := bitset.New(64000)
	b.SetBit(1000)
	b.SetBit(999)
	if b.Bit(1000) {
		b.ClearBit(1000)
	}
	b.Clear()
	
*/
package bitset

import (
	"fmt"
)

// BitSet internal details 
type BitSet struct {
	capacity uint
	set      []uint64
}

// Make a BitSet with an upper limit on size.
func New(capacity uint) *BitSet {
	return &BitSet{capacity, make([]uint64, (capacity+(64-1))>>6)}
}

// Query maximum size of a bit set
func (b *BitSet) Cap() uint {
	return b.capacity
}

/// Test whether bit i is set. 
func (b *BitSet) Bit(i uint) bool {
	if i >= b.capacity {
		panic(fmt.Sprintf("index out of range: %v", i))
	}
	return ((b.set[i>>6] & (1 << (i & (64-1)))) != 0)
}

// Set bit i to 1
func (b *BitSet) SetBit(i uint) {
	if i >= b.capacity {
		panic(fmt.Sprintf("index out of range: %v", i))
	}
	b.set[i>>6] |= (1 << (i & (64-1)))
}

// Clear bit i to 0
func (b *BitSet) ClearBit(i uint) {
	if i >= b.capacity {
		panic(fmt.Sprintf("index out of range: %v", i))
	}	
	b.set[i>>6] &^= 1 << (i & (64-1))
}

// Clear entire BitSet
func (b *BitSet) Clear() {
	if b != nil {
		for i := range b.set {
			b.set[i] = 0
		}
	}
}

// From Wikipedia: http://en.wikipedia.org/wiki/Hamming_weight                                     
const m1  uint64 = 0x5555555555555555 //binary: 0101...
const m2  uint64 = 0x3333333333333333 //binary: 00110011..
const m4  uint64 = 0x0f0f0f0f0f0f0f0f //binary:  4 zeros,  4 ones ...

// From Wikipedia: count number of set bits.
func popcount_2(x uint64) uint64 {
    x -= (x >> 1) & m1;             //put count of each 2 bits into those 2 bits
    x = (x & m2) + ((x >> 2) & m2); //put count of each 4 bits into those 4 bits 
    x = (x + (x >> 4)) & m4;        //put count of each 8 bits into those 8 bits 
    x += x >>  8;  //put count of each 16 bits into their lowest 8 bits
    x += x >> 16;  //put count of each 32 bits into their lowest 8 bits
    x += x >> 32;  //put count of each 64 bits into their lowest 8 bits
    return x & 0x7f;
}

// Count (number of set bits)
func (b *BitSet) Count() uint {
   	if b != nil {
		cnt := uint64(0)
		for _, word := range b.set {
			cnt += popcount_2(word)
		}
		return uint(cnt)
	}
	return 0
}



func (b *BitSet) XorBit(i uint) {
	if i >= b.capacity {
		panic(fmt.Sprintf("index out of range: %v", i))
	}
	b.set[i >> 6] ^= 1 << (i & (64 - 1))
}

func (b *BitSet) ClearAll() {
	for i, _ := range b.set {
		b.set[i] = 0
	}
}

func (b *BitSet) SetAll() {
	for i, _ := range b.set {
		b.set[i] = ^uint64(0)
	}
}

func (b *BitSet) XorAll() {
	for i, _ := range b.set {
		b.set[i] ^= ^uint64(0)
	}
}
func (b *BitSet) Equ(c *BitSet) bool {
	if c == nil {
		return false
	}
	if b.capacity != c.capacity {
		return false
	}
	for p, v := range b.set {
		if c.set[p] != v {
			return false
		}
	}
	return true
}

func (b *BitSet) Clone() *BitSet {
	c := New(b.capacity)
	copy(c.set, b.set)
	return c
}

func (b *BitSet) Copy(c *BitSet) (count uint) {
	if c == nil {
		return
	}
	copy(c.set, b.set)
	count = c.capacity
	if b.capacity < c.capacity {
		count = b.capacity
	}
	return
}

func (b *BitSet) Sub(start, end uint) *BitSet {
	if end <= start || end > b.capacity {
		return nil
	}
	c := New(end - start)
	if start&(64-1) == 0 {
		copy(c.set, b.set[start>>6:(end+63)>>6])
		return c
	}
	ipos := start & (64 - 1)
	ifirst := start >> 6
	icount := end - start
	ilen := icount >> 6
	var i uint
	for i = 0; i < ilen; i++ {
		v := b.set[i+ifirst] >> ipos
		c.set[i] = v | b.set[i+ifirst+1]<<(64-ipos)
	}
	if icount&(64-1) != 0 {
		c.set[i] = b.set[i+ifirst] >> ipos
	}
	return c
}