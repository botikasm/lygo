package lygo_num2word_languages

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleIntegerToIrIr() {
	fmt.Println(IntegerToIrIr(42))
	// Output: چهل و دو
}

func TestIntegerToIrIr(t *testing.T) {
	t.Parallel()

	tests := map[int]string{
		-1:        "منفی یک",
		0:         "صفر",
		1:         "یک",
		9:         "نه",
		10:        "ده",
		11:        "یازده",
		19:        "نوزده",
		20:        "بیست",
		21:        "بیست و یک",
		80:        "هشتاد",
		90:        "نود",
		99:        "نود و نه",
		100:       "صد",
		101:       "صد یک",
		111:       "صد یازده",
		120:       "صد بیست",
		121:       "صد بیست و یک",
		900:       "نهصد",
		909:       "نهصد نه",
		919:       "نهصد نوزده",
		990:       "نهصد نود",
		999:       "نهصد نود و نه",
		1000:      "یک هزار",
		2000:      "دو هزار",
		4000:      "چهار هزار",
		5000:      "پنج هزار",
		11000:     "یازده هزار",
		21000:     "بیست و یک هزار",
		999000:    "نهصد نود و نه هزار",
		999999:    "نهصد نود و نه هزار نهصد نود و نه",
		1000000:   "یک میلیون",
		2000000:   "دو میلیون",
		4000000:   "چهار میلیون",
		5000000:   "پنج میلیون",
		100100100: "صد میلیون صد هزار صد",
		500500500: "پانصد میلیون پانصد هزار پانصد",
		606606606: "ششصد شش میلیون ششصد شش هزار ششصد شش",
		999000000: "نهصد نود و نه میلیون",
		999000999: "نهصد نود و نه میلیون نهصد نود و نه",
		999999000: "نهصد نود و نه میلیون نهصد نود و نه هزار",
		999999999: "نهصد نود و نه میلیون نهصد نود و نه هزار نهصد نود و نه",
	}

	for input, expectedOutput := range tests {
		name := fmt.Sprintf("%d", input)
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, expectedOutput, IntegerToIrIr(input))
		})
	}
}
