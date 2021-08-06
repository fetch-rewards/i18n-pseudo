package pseudo

import (
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Format(t *testing.T) {
	gen := New(FormatOptions{})

	t.Run("Format - Alphanumeric characters are converted", func(t *testing.T) {
		input := "Aaron Iz Kewl 1234567890"
		results := gen.Format(input)

		assert.True(t, strings.Contains(results, "①②③④⑤⑥⑦⑧⑨⓪"))
		assert.True(t, !strings.Contains(results, "Aaron Iz Kewl"))
		assert.True(t, !strings.Contains(results, "1234567890"))
	})

	t.Run("Format - Ignore <>", func(t *testing.T) {
		input := "This is a <bold>bolded text</bold> and this is not."

		results := gen.Format(input)
		assert.True(t, len(results) > len(input))
		assert.True(t, !strings.Contains(results, "This is a"))
		assert.True(t, strings.Contains(results, "<bold>"))
		assert.True(t, !strings.Contains(results, "bolded text"))
		assert.True(t, strings.Contains(results, "</bold>"))
		assert.True(t, !strings.Contains(results, "and this is not"))
		assert.True(t, strings.Contains(results, "."))
	})

	t.Run("Format - Ignore {}", func(t *testing.T) {
		input := "This user {userId} is not good!"

		results := gen.Format(input)
		assert.True(t, len(results) > len(input))
		assert.True(t, !strings.Contains(results, "This user"))
		assert.True(t, strings.Contains(results, "{userId}"))
		assert.True(t, !strings.Contains(results, "is not good"))
		assert.True(t, strings.Contains(results, "!"))
	})

	t.Run("Format - Ignore {} and <>", func(t *testing.T) {
		input := "This user <italic> ID {userId} is not good</italic> when they are a fraudster.  {receiptId} is bad!"

		results := gen.Format(input)
		fmt.Println(results)
		assert.True(t, len(results) > len(input))
		assert.True(t, !strings.Contains(results, "This user"))
		assert.True(t, strings.Contains(results, "<italic>"))
		assert.True(t, !strings.Contains(results, "ID"))
		assert.True(t, strings.Contains(results, "{userId}"))
		assert.True(t, !strings.Contains(results, "is not good"))
		assert.True(t, strings.Contains(results, "</italic>"))
		assert.True(t, !strings.Contains(results, "when they are a fraudster"))
		assert.True(t, strings.Contains(results, "."))
		assert.True(t, strings.Contains(results, "{receiptId}"))
		assert.True(t, !strings.Contains(results, "is bad"))
		assert.True(t, strings.Contains(results, "!"))
	})

	t.Run("Format - Ignore characters following only one bracket > || <", func(t *testing.T) {
		input := "Postgres > Oracle and you that is <bold>true</bold> right?"

		results := gen.Format(input)
		assert.True(t, len(results) > len(input))
		assert.True(t, !strings.Contains(results, "Postgres"))
		assert.True(t, strings.Contains(results, ">"))
		assert.True(t, !strings.Contains(results, "Oracle"))
		assert.True(t, !strings.Contains(results, "and you that is"))
		assert.True(t, strings.Contains(results, "<bold>"))
		assert.True(t, !strings.Contains(results, "true"))
		assert.True(t, strings.Contains(results, "</bold>"))
		assert.True(t, !strings.Contains(results, "right"))
		assert.True(t, strings.Contains(results, "?"))
	})

	t.Run("Format - Ignore characters following %", func(t *testing.T) {
		input := "This %s has %d kids!"

		results := gen.Format(input)
		assert.True(t, len(results) > len(input))
		assert.True(t, !strings.Contains(results, "This"))
		assert.True(t, strings.Contains(results, "%s"))
		assert.True(t, !strings.Contains(results, "has"))
		assert.True(t, strings.Contains(results, "%d"))
		assert.True(t, !strings.Contains(results, "kids"))
		assert.True(t, strings.Contains(results, "!"))
	})

	t.Run("Format - Do not parse when no characters follow %", func(t *testing.T) {
		input := "Do you know what 5 % 2 is?"

		results := gen.Format(input)
		assert.True(t, len(results) > len(input))
		assert.True(t, !strings.Contains(results, "Do you know what 5 "))
		assert.True(t, strings.Contains(results, "%"))
		assert.True(t, !strings.Contains(results, " 2"))
		assert.True(t, !strings.Contains(results, "is"))
		assert.True(t, strings.Contains(results, "?"))

		input = "50% of the population do not!"

		results = gen.Format(input)
		assert.True(t, len(results) > len(input))
		assert.True(t, !strings.Contains(results, "50"))
		assert.True(t, strings.Contains(results, "%"))
		assert.True(t, !strings.Contains(results, " of the population do not"))
		assert.True(t, strings.Contains(results, "!"))
	})
}

func Test_TargetExpansion(t *testing.T) {
	gen := New(FormatOptions{})

	t.Run("Target Expansion - 0 characters does not get formatted", func(t *testing.T) {
		input := ""
		results := gen.Format(input)
		assert.Equal(t, input, results)
	})

	t.Run("Target Expansion - 1 character expands by 80%", func(t *testing.T) {
		input := "1"
		results := gen.Format(input)
		fmt.Println(results)
		assert.Equal(t, getExpectedResultLen(input, chars1to5ExpnPercent), len([]rune(results)))
	})

	t.Run("Target Expansion - 5 character expands by 80%", func(t *testing.T) {
		input := "12345"
		results := gen.Format(input)
		assert.Equal(t, getExpectedResultLen(input, chars1to5ExpnPercent), len([]rune(results)))
	})

	t.Run("Target Expansion - 6 chars expands by 70%", func(t *testing.T) {
		input := "123456"
		results := gen.Format(input)
		assert.Equal(t, getExpectedResultLen(input, chars6to10ExpnPercent), len([]rune(results)))
	})

	t.Run("Target Expansion - 10 chars expands by 70%", func(t *testing.T) {
		input := "1234567890"
		results := gen.Format(input)
		assert.Equal(t, getExpectedResultLen(input, chars6to10ExpnPercent), len([]rune(results)))
	})

	t.Run("Target Expansion - 11 chars expands by 50%", func(t *testing.T) {
		input := "1234567890a"
		results := gen.Format(input)
		assert.Equal(t, getExpectedResultLen(input, chars11to25ExpnPercent), len([]rune(results)))
	})

	t.Run("Target Expansion - 25 chars expands by 50%", func(t *testing.T) {
		input := "1234567890abcdefghijklmno"
		results := gen.Format(input)
		assert.Equal(t, getExpectedResultLen(input, chars11to25ExpnPercent), len([]rune(results)))
	})

	t.Run("Target Expansion - 26 chars expands by 25%", func(t *testing.T) {
		input := "1234567890abcdefghijklmnop"
		results := gen.Format(input)
		assert.Equal(t, getExpectedResultLen(input, chars26to75ExpnPercent), len([]rune(results)))
	})

	t.Run("Target Expansion - 75 chars expands by 25%", func(t *testing.T) {
		input := "1234567890abcdefghijklmno1234567890abcdefghijklmno1234567890abcdefghijklmno"
		results := gen.Format(input)
		assert.Equal(t, getExpectedResultLen(input, chars26to75ExpnPercent), len([]rune(results)))
	})

	t.Run("Target Expansion - 76 chars expands by 0%", func(t *testing.T) {
		input := "1234567890abcdefghijklmno1234567890abcdefghijklmno1234567890abcdefghijklmno1"
		results := gen.Format(input)
		assert.Equal(t, getExpectedResultLen(input, chars0ExpnPercent), len([]rune(results)))
	})
}

func Test_FormatOptions(t *testing.T) {
	input := "Some Text"

	t.Run("Format Options - AppendChars", func(t *testing.T) {
		appendChars := ">>>"
		gen := New(FormatOptions{
			AppendChars: appendChars,
		})

		results := gen.Format(input)
		assert.Equal(t, results[len(results)-3:], appendChars)
	})

	t.Run("Format Options - ExpandChars", func(t *testing.T) {
		expandChars := "$$"
		gen := New(FormatOptions{
			ExpandChars: expandChars,
		})

		results := gen.Format(input)
		assert.True(t, strings.Contains(results, expandChars))
	})

	t.Run("Format Options - PrependChars", func(t *testing.T) {
		prependChars := "<<<"
		gen := New(FormatOptions{
			PrependChars: prependChars,
		})

		results := gen.Format(input)
		assert.Equal(t, results[:3], prependChars)
	})

	t.Run("Format Options - PreventExpansion", func(t *testing.T) {
		gen := New(FormatOptions{
			PreventExpansion: true,
		})

		results := gen.Format(input)
		expResults := len([]rune(defaultAppendChars)) + len([]rune(defaultPrependChars)) + len([]rune(input))
		assert.Equal(t, expResults, len([]rune(results)))
	})

	t.Run("Format Options - PseudoChars", func(t *testing.T) {
		gen := New(FormatOptions{
			PseudoChars: map[string]string{
				"a": "z",
				"b": "y",
				"c": "x",
				"d": "w",
				"e": "v",
			},
		})
		input := "abcde fghij"

		results := gen.Format(input)
		assert.True(t, len([]rune(results)) > len([]rune(input)))
		assert.True(t, strings.Contains(results, "zyxwv"))
		assert.True(t, strings.Contains(results, "fghij"))
	})

	t.Run("Format Options - TargetExpansion", func(t *testing.T) {
		gen := New(FormatOptions{
			TargetExpansion: 1.2,
		})

		results := gen.Format(input)
		assert.Equal(t, getExpectedResultLen(input, 1.2), len([]rune(results)))
	})
}

func getExpectedResultLen(input string, per float64) int {
	return len([]rune(defaultPrependChars)) +
		len([]rune(input)) +
		int(math.RoundToEven(float64(len(input)) * per)) +
		len([]rune(defaultAppendChars))
}
