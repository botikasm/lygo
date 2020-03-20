package lygo_regex

import (
	"fmt"
	"github.com/botikasm/lygo/base/lygo_strings"
	"reflect"
	"testing"
)

func TestIsValidJSON(t *testing.T) {
	obj := "{\"name\": \"foo\"}"
	array := "[{\"name\": \"foo\"}]"

	if !IsValidJsonObject(obj) {
		t.Errorf("Not a valid object %v", obj)
	}
	if !IsValidJsonArray(array) {
		t.Errorf("Not a valid array %v", array)
	}
	if IsValidJsonObject("obj") {
		t.Error("NOT A JSON")
	}
	if IsValidJsonArray("obj") {
		t.Error("NOT A JSON")
	}
}

func TestIndexStartAt(t *testing.T) {
	text := "this is sample text with som 'is' inside"
	pattern := "is"
	offset := 1

	result := WildcardIndex(text, pattern, offset)
	if len(result) != 2 {
		t.Error("Expected 2 matching")
	} else {
		if result[0] != 5 {
			t.Error("Expected 5", result[0])
		}
		if result[1] != 30 {
			t.Error("Expected 30", result[1])
		}
	}
}

func TestWildCardScore(t *testing.T) {
	text := "this is sample text with som 'is' inside and more"
	expressions := []string{"thi? is", "tex* * so?"}
	fmt.Println("all:", WildcardScoreAll(text, expressions),
		"any:", WildcardScoreAny(text, expressions),
		"best:", WildcardScoreBest(text, expressions),
		"\t", text, " ", lygo_strings.ConcatSep(", ", expressions))

	text = "this os sample | ssddfgr good"
	expressions = []string{"thi? ?s", "sampl? * * good"}
	fmt.Println("all:", WildcardScoreAll(text, expressions),
		"any:", WildcardScoreAny(text, expressions),
		"best:", WildcardScoreBest(text, expressions),
		"\t", text, " ", lygo_strings.ConcatSep(", ", expressions))

	text = "this oss sample | ssddfgr good"
	expressions = []string{"thi? ?s", "sampl? * * good"}
	fmt.Println("all:", WildcardScoreAll(text, expressions),
		"any:", WildcardScoreAny(text, expressions),
		"best:", WildcardScoreBest(text, expressions),
		"\t", text, " ", lygo_strings.ConcatSep(", ", expressions))

	text = "this oss sample | ssddfgr good"
	expressions = []string{"thi? ?s"}
	fmt.Println("all:", WildcardScoreAll(text, expressions),
		"any:", WildcardScoreAny(text, expressions),
		"best:", WildcardScoreBest(text, expressions),
		"\t", text, " ", lygo_strings.ConcatSep(", ", expressions))

	text = "this oss sample | ssddfgr good"
	expressions = []string{"sampl? * * good"}
	fmt.Println("all:", WildcardScoreAll(text, expressions),
		"any:", WildcardScoreAny(text, expressions),
		"best:", WildcardScoreBest(text, expressions),
		"\t", text, " ", lygo_strings.ConcatSep(", ", expressions))

}

func TestIWildcardMatchBetween(t *testing.T) {
	text := "this is sample \ntext with som 'is' inside\n and more"
	patternStart := "is"
	patternEnd := "\n"
	offset := 1

	result := WildcardMatchBetween(text, offset, patternStart, patternEnd, " '")
	if len(result) != 2 {
		t.Error("Expected 2 matching")
	} else {
		fmt.Println(result)
	}
}

func TestIndexLenPair(t *testing.T) {
	text := "this is sample text with som 'is' inside"
	pattern := "is"
	offset := 1

	result := WildcardIndexLenPair(text, pattern, offset)
	if len(result) != 2 {
		t.Error("Expected 2 matching")
	} else {
		if result[0][0] != 5 {
			t.Error("Expected 5", result[0])
		}
		if result[1][0] != 30 {
			t.Error("Expected 30", result[1])
		}
	}
}

func TestBtcAddresses(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BtcAddresses(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BtcAddresses() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreditCards(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreditCards(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreditCards() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"email valid", args{"Questa è una data valida 23 Mar 2017"}, []string{"23 Mar 2017"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Date(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Date() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEmails(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"email valid", args{"Questa è una email valida angelo.geminiani@gmail.com"}, []string{"angelo.geminiani@gmail.com"}},
		{"email invalid", args{"Questa è una email valida angelo.geminiani@gmail"}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Emails(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Emails() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGUIDs(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GUIDs(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GUIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGitRepos(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GitRepos(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GitRepos() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHexColors(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HexColors(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HexColors() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIBANs(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IBANs(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IBANs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIPs(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IPs(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IPs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIPv4s(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IPv4s(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IPv4s() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIPv6s(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IPv6s(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IPv6s() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestISBN10s(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ISBN10s(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ISBN10s() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestISBN13s(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ISBN13s(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ISBN13s() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLinks(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Links(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Links() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMACAddresses(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MACAddresses(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MACAddresses() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMCCreditCards(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MCCreditCards(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MCCreditCards() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMD5Hexes(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MD5Hexes(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MD5Hexes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotKnownPorts(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NotKnownPorts(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NotKnownPorts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPhones(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Phones(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Phones() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPhonesWithExts(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PhonesWithExts(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PhonesWithExts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPoBoxes(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PoBoxes(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PoBoxes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrices(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"pricex", args{"€3,200 3,200 3.200,12 3,200.00 1,245,123.123"},
			[]string{"€3,200", "3,200", "3.200,12", "3,200.00", "1,245,123.123"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Prices(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Prices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSHA1Hexes(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SHA1Hexes(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SHA1Hexes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSHA256Hexes(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SHA256Hexes(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SHA256Hexes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSSNs(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SSNs(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SSNs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStreetAddresses(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StreetAddresses(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StreetAddresses() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTime(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Time(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Time() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVISACreditCards(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := VISACreditCards(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VISACreditCards() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZipCodes(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ZipCodes(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ZipCodes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumbers(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"3 numbers", args{"here are some numbers 123,3 123.3 and 1234 1.250,345 1,250.345"}, []string{"123,3", "123.3", "1234", "1.250,345", "1,250.345"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Numbers(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Numbers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWildCard(t *testing.T) {
	expression := "cod?ce"
	text := "codice codoce codice\ncodice coaudace, codice, some other text here       "
	want := []string{"codice", "codoce", "codice", "codice", "codice"}
	got := WildcardMatch(text, expression)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("WildcardMatch() = %v, want %v", got, want)
	}

	expression = "cod ce"
	text = "codice codoce cod ce\ncodice"
	want = []string{"cod ce"}
	got = WildcardMatch(text, expression)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("WildcardMatch() = %v, want %v", got, want)
	}

	expression = "cod? 80"
	text = "cod  80 Cod. 80"
	want = []string{"cod  80"}
	got = WildcardMatch(text, expression)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("WildcardMatch() = %v, want %v", got, want)
	}

	expression = "?od*80"
	text = "cod  80  e novanta Cod. 80"
	want = []string{"cod  80", "Cod. 80"}
	got = WildcardMatch(text, expression)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("WildcardMatch() = %v, want %v", got, want)
	}

	expression = "*od*80"
	text = "cod  80  e novanta Cod. 90"
	want = []string{"cod  80"}
	got = WildcardMatch(text, expression)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("WildcardMatch() = %v, want %v", got, want)
	}
}
