package libraries

import (
	"math/big"

	"github.com/pkg/errors"

	"github.com/openchat-im/go-uniswap/utils"
)

var (
	ErrTickMisaligned = errors.New("tick misaligned")
)

// TickBitmap represents a bitmap of initialized ticks.
// Each "word" in the bitmap represents 256 ticks.
// A set bit in a word indicates that the corresponding tick is initialized.
type TickBitmap struct {
	words map[int16]*big.Int
}

// NewTickBitmap creates and returns a new TickBitmap instance.
func NewTickBitmap() *TickBitmap {
	return &TickBitmap{
		words: make(map[int16]*big.Int),
	}
}

// compress calculates the "compressed tick" by dividing the tick by tickSpacing.
// For negative ticks not aligned exactly, it decrements the result to maintain consistency.
func compress(tick int32, tickSpacing int32) int32 {
	c := tick / tickSpacing
	if tick < 0 && tick%tickSpacing != 0 {
		c--
	}
	return c
}

// position computes the word index and bit position of a compressed tick in the bitmap.
// wordPos: the index of the 256-bit word containing the tick.
// bitPos: the bit position (0-255) within the word.
func position(tick int32) (wordPos int16, bitPos uint8) {
	wordPos = int16(tick >> 8)
	bitPos = uint8(tick & 0xff)
	return
}


// FlipTick flips the state of a tick in the bitmap (initialized ↔ uninitialized).
//
// Parameters:
//   - tick: the tick to flip
//   - tickSpacing: the spacing between ticks; tick must be a multiple of tickSpacing
//
// Returns:
//   - error: ErrTickMisaligned if tick is not aligned with tickSpacing
func (tb *TickBitmap) FlipTick(tick int32, tickSpacing int32) error {
	if tick%tickSpacing != 0 {
		return ErrTickMisaligned
	}

	compressed := compress(tick, tickSpacing)
	wordPos, bitPos := position(compressed)

	mask := big.NewInt(0).SetUint64(1)
	mask.Lsh(mask, uint(bitPos)) // Create a mask with a 1 at bitPos

	word, ok := tb.words[wordPos]
	if !ok {
		word = big.NewInt(0)
		tb.words[wordPos] = word
	}

	word.Xor(word, mask)
	return nil
}

// NextInitializedTickWithinOneWord finds the next initialized tick within a single 256-tick word.
//
// Parameters:
//   - tick: the reference tick
//   - tickSpacing: the spacing between ticks
//   - lte: if true, search for the next tick <= reference; otherwise, search for next tick > reference
//
// Returns:
//   - next: the next initialized tick (or boundary if none exists)
//   - initialized: true if an initialized tick was found
//   - err: any error occurred during bit scanning
func (tb *TickBitmap) NextInitializedTickWithinOneWord(tick int32, tickSpacing int32, lte bool) (next int32, initialized bool, err error) {
	compressed := compress(tick, tickSpacing)

	if lte {
		wordPos, bitPos := position(compressed)
		word, ok := tb.words[wordPos]
		if !ok {
			word = big.NewInt(0)
		}

		// Mask to include all bits <= bitPos
		mask := new(big.Int).Lsh(big.NewInt(1), uint(bitPos+1)) 
		mask.Sub(mask, big.NewInt(1))                              
		masked := big.NewInt(0).And(word, mask)                   

		if masked.Sign() == 0 {
			// No initialized ticks ≤ current tick in this word
			return (compressed - int32(bitPos)) * tickSpacing, false, nil
		}

		// Find most significant bit (highest initialized tick ≤ compressed)
		msb, err := utils.MostSignificantBit(masked)
		if err != nil {
			return 0, false, err
		}
		return (compressed - int32(int(bitPos)-msb)) * tickSpacing, true, nil
	} else {
		compressed++
		wordPos, bitPos := position(compressed)
		word, ok := tb.words[wordPos]
		if !ok {
			word = big.NewInt(0)
		}

		// Shift word right to start search from bitPos
		mask := new(big.Int).Rsh(new(big.Int).Set(word), uint(bitPos)) 
		if mask.Sign() == 0 {
			// No initialized ticks > current tick in this word
			return (compressed + int32(255-bitPos)) * tickSpacing, false, nil
		}

		// Find least significant bit (lowest initialized tick ≥ compressed)
		lsb, err := utils.LeastSignificantBit(mask)
		if err != nil {
			return 0, false, err
		}
		return (compressed + int32(lsb)) * tickSpacing, true, nil
	}
}