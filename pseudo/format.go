package pseudo

import (
	_ "embed"
	"math"
	"math/rand"
	"strings"
	"unicode"

	jsoniter "github.com/json-iterator/go"
)

//go:embed data/pseudo_chars.json
var defaultPseudoChars string

const (
	defaultPrependChars = "你好"
	defaultAppendChars  = "世界"
	allCaseAlphaChars   = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	chars0ExpnPercent      = 0.0
	chars1to5ExpnPercent   = 0.80
	chars6to10ExpnPercent  = 0.70
	chars11to25ExpnPercent = 0.50
	chars26to75ExpnPercent = 0.25
)

type FormatOptions struct {
	// The characters to put at the end of output
	AppendChars string

	// The characters to use for the expansion
	ExpandChars string

	// The characters to put at the beginning of output
	PrependChars string

	// Whether to expand the length of text
	PreventExpansion bool

	// The pseudo characters to use when replacing ASCII characters
	PseudoChars map[rune]rune

	// Length of expanded characters appended
	TargetExpansion float64
}

type Format struct {
	Options FormatOptions
}

func New(fo FormatOptions) *Format {
	gen := &Format{Options: fo}

	// Setting the Pseudo Characters if they weren't provided
	if gen.Options.PseudoChars == nil {
		var charMap map[string]string
		_ = jsoniter.UnmarshalFromString(defaultPseudoChars, &charMap)

		gen.Options.PseudoChars = make(map[rune]rune, len(charMap))
		for from, to := range charMap {
			gen.Options.PseudoChars[[]rune(from)[0]] = []rune(to)[0]
		}
	}

	if gen.Options.AppendChars == "" {
		gen.Options.AppendChars = defaultPrependChars
	}
	if gen.Options.ExpandChars == "" {
		gen.Options.ExpandChars = allCaseAlphaChars
	}
	if gen.Options.PrependChars == "" {
		gen.Options.PrependChars = defaultAppendChars
	}

	return gen
}

// Accepts an incoming string and converts it to Pseudo Translation
func (pf Format) Format(input string) string {
	if input == "" {
		return input
	}

	output := input

	// Expand the text unless told otherwise
	if !pf.Options.PreventExpansion {
		output = pf.expand(input)
	}

	// Convert it to pseudo translation
	output = pf.makePseudo(output)

	// Add challenge characters
	return pf.Options.PrependChars + output + pf.Options.AppendChars
}

// Accepts a given input and expands based on length
func (pf Format) expand(input string) string {
	var targetExpansion float64
	inputLen := len(input)

	// Use configured TargetExpansion, otherwise, output will expand by based on incoming length
	if pf.Options.TargetExpansion > 0.0 {
		targetExpansion = pf.Options.TargetExpansion
	} else {
		if inputLen <= 5 {
			targetExpansion = chars1to5ExpnPercent
		} else if inputLen <= 10 {
			targetExpansion = chars6to10ExpnPercent
		} else if inputLen <= 25 {
			targetExpansion = chars11to25ExpnPercent
		} else if inputLen <= 75 {
			targetExpansion = chars26to75ExpnPercent
		} else {
			targetExpansion = chars0ExpnPercent
		}
	}

	addlChars := int(math.RoundToEven(float64(inputLen) * targetExpansion))
	return input + pf.generateRandomExpansion(addlChars)
}

// Generates a random string based on expansion characters
func (pf Format) generateRandomExpansion(addlChars int) string {
	var sb strings.Builder

	for x := 0; x < addlChars; x++ {
		thisIdx := rand.Intn(len(pf.Options.ExpandChars))
		thisChar := []rune(pf.Options.ExpandChars)[thisIdx]
		sb.WriteRune(thisChar)
	}

	return sb.String()
}

// Converts the given text to pseudo characters
func (pf Format) makePseudo(in string) string {
	var sb strings.Builder
	var hasOpenCurlyBracket, hasOpenAngleBracket, hasOpenPercentSign bool

	input := []rune(in)
	for x := 0; x < len(input); x++ {
		thisChar := input[x]

		if thisChar == '{' {
			hasOpenCurlyBracket = true
		} else if thisChar == '}' {
			hasOpenCurlyBracket = false
		} else if thisChar == '<' {
			// Ignoring when input contains a single < by doing a space check against the next character
			if x < len(input)-1 {
				nextChar := input[x+1]
				if !unicode.IsSpace(nextChar) {
					hasOpenAngleBracket = true
				}
			}
		} else if thisChar == '>' {
			hasOpenAngleBracket = false
		} else if thisChar == '%' {
			hasOpenPercentSign = true
		} else if hasOpenPercentSign && x >= 2 {
			// Only the first character following % will not get pseudo translated
			prevSecondChar := input[x-2]
			if prevSecondChar == '%' {
				hasOpenPercentSign = false
			}
		}

		// Skip replacing of pseudo character, if { or < is open.
		if !(hasOpenAngleBracket || hasOpenCurlyBracket || hasOpenPercentSign) {
			// Replace with pseudo variant if exists
			if val, ok := pf.Options.PseudoChars[thisChar]; ok {
				thisChar = val
			}
		}

		sb.WriteRune(thisChar)
	}

	return sb.String()
}
