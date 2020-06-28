/*
Copyright (C) 2019 Graeme Sutherland, Nodestone Limited


This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package lang

// Map of letters to morse - good for the encoder
// and processed below to make a decoding tree
var MorseMap = map[string]string{
	"a": ".-",
	"b": "-...",
	"c": "-.-.",
	"d": "-..",
	"e": ".",
	"f": "..-.",
	"g": "--.",
	"h": "....",
	"i": "..",
	"j": ".---",
	"k": "-.-",
	"l": ".-..",
	"m": "--",
	"n": "-.",
	"o": "---",
	"p": ".--.",
	"q": "--.-",
	"r": ".-.",
	"s": "...",
	"t": "-",
	"u": "..-",
	"v": "...-",
	"w": ".--",
	"x": "-..-",
	"y": "-.--",
	"z": "--..",
	"1": ".----",
	"2": "..---",
	"3": "...--",
	"4": "....-",
	"5": ".....",
	"6": "-....",
	"7": "--...",
	"8": "---..",
	"9": "----.",
	"0": "-----",
	".": ".-.-.-",
	",": "--..--",
	"?": "..--..",
	`'`: ".----.",
	"/": "-..-.",
	"(": "-.--.",
	")": "-.--.-",
	"&": ".-...",
	":": "---...",
	";": "-.-.-",
	"=": "-...-",
	"+": ".-.-.",
	"-": "-....-",
	"_": "..--.-",
	`"`: ".-..-.",
	"$": "...-..-",
	"!": "-.-.--",
	"@": ".--.-.",
}

type Symbol struct {
	Text  string
	Morse string
	Dit   *Symbol
	Dah   *Symbol
}

func MorseTree() *Symbol {
	root := Symbol{"", "", nil, nil}
	var last *Symbol = &root
	var next *Symbol

	for k, v := range MorseMap {
		last = &root
		next = nil

		for i, b := range v {
			if b == '.' {
				next = last.Dit
				if next == nil {
					next = new(Symbol)
					last.Dit = next
				}
			} else {
				next = last.Dah
				if next == nil {
					next = new(Symbol)
					last.Dah = next
				}
			}

			if i+1 >= len(v) {
				next.Text = k
				next.Morse = v
			} else {
				last = next
				next = nil
			}
		}
	}

	return &root
}

func SearchMorseTree(morse string, root *Symbol) string {
	pos := root

	for i, b := range morse {
		if b == '.' {
			pos = pos.Dit
		} else {
			pos = pos.Dah
		}

		if pos == nil {
			return ""
		}

		if i+1 >= len(morse) {
			return pos.Text
		}

	}

	return ""
}
