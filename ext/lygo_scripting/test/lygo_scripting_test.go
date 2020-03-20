package test

import (
	"fmt"
	"github.com/botikasm/lygo/base/lygo_io"
	"github.com/botikasm/lygo/ext/lygo_scripting"
	"github.com/botikasm/lygo/ext/lygo_scripting/goja_nodejs/require"
	"github.com/dop251/goja"
	"testing"
)

func TestSimple(t *testing.T) {

	vm := lygo_scripting.New()

	EXPECTED := int64(30)

	vm.Set("VAR_x", 26)
	v, err := vm.RunString("2 + 2 + VAR_x")

	if err != nil {
		panic(err)
	}
	if num := v.Export().(int64); num != EXPECTED {
		panic(num)
	} else {
		fmt.Println(num)
	}
}

func TestExpression(t *testing.T) {

	vm := lygo_scripting.New()

	TEXT := "426"

	vm.Set("x", "26")
	v, err := vm.RunString("2 + 2 + x")

	if err != nil {
		panic(err)
	}
	if num := v.Export().(string); num != TEXT {
		panic(num)
	} else {
		fmt.Println(num)
	}
}

func TestFunc(t *testing.T) {

	vm := lygo_scripting.New()

	want := int64(4)

	vm.Set("myfunc", func(call goja.FunctionCall) goja.Value {
		val_1 := call.Argument(0).ToInteger()
		val_2 := call.Argument(1).ToInteger()
		panic("errore")
		return vm.ToValue(val_1 + val_2)
	})
	v, err := vm.RunString("myfunc(2,2)")

	if err != nil {
		panic(err)
	}
	if num := v.ToInteger(); num != want {
		panic(num)
	} else {
		fmt.Println(num)
	}
}

func TestTool(t *testing.T) {

	TEXT := "this is a text\ncod. 80 is a matching value! Cod. 80"

	vm := lygo_scripting.New()
	vm.SetToolContext("$strings", TEXT)
	vm.SetToolContext("$regexps", TEXT)
	// vm.SetToolContext("$strings", 1234)

	want := true

	v, err := vm.RunString("$regexps.HasMatch('?od??80') || $regexps.HasMatch('?od?80')")

	if err != nil {
		t.Error(err)
	}
	if b := v.ToBoolean(); b != want {
		t.Error("Unexpected response", b)
	} else {
		fmt.Println(b)
	}

	// MatchAll
	v, err = vm.RunString("$regexps.MatchAll('?od??80')")
	if err != nil {
		t.Error(err)
	} else {
		array, _ := v.Export().([]string)
		if nil == array {
			t.Error("Unexpected response: ", v)
		} else {
			fmt.Println(v)
		}
	}

	// MatchAt
	v, err = vm.RunString("$regexps.MatchAt('?od??80', 1)")
	if err != nil {
		t.Error(err)
	} else {
		s, _ := v.Export().(string)
		if s != "Cod. 80" {
			t.Error("Unexpected response: ", v)
		} else {
			fmt.Println(v)
		}
	}

	// MatchFirst
	v, err = vm.RunString("$regexps.MatchFirst('?od??80')")
	if err != nil {
		t.Error(err)
	} else {
		s, _ := v.Export().(string)
		if s != "cod. 80" {
			t.Error("Unexpected response: ", v)
		} else {
			fmt.Println("MatchFirst: ", v)
		}
	}

	// MatchFirstStartingAt
	v, err = vm.RunString("$regexps.MatchFirst('?od??80')")
	if err != nil {
		t.Error(err)
	} else {
		s, _ := v.Export().(string)
		if s != "cod. 80" {
			t.Error("Unexpected response: ", v)
		} else {
			fmt.Println("MatchFirst: ", v)
		}
	}

	// MatchLast complex
	v, err = vm.RunString("(function(){var x = $regexps.MatchLast('?od??80');\nif (!!x){return x;} else {return 'NOT FOUND'}})()")
	if err != nil {
		t.Error(err)
	} else {
		s, _ := v.Export().(string)
		if s != "Cod. 80" {
			t.Error("Unexpected response: ", v)
		} else {
			fmt.Println("MatchLast: ", v)
		}
	}

	// Split
	v, err = vm.RunString("$strings.Split(' \\n', 'this is a phrase')") // passing last parameter
	if err != nil {
		t.Error(err)
	} else {
		array, _ := v.Export().([]string)
		if nil == array {
			t.Error("Unexpected response: ", v)
		} else {
			fmt.Println(v)
		}
	}

	// combine Split with $arrays.GetAt
	v, err = vm.RunString("$arrays.GetAt(1, $strings.Split(' \\n', 'this \\nis a \\nphrase'))") // passing last parameter
	if err != nil {
		t.Error(err)
	} else {
		s, _ := v.Export().(string)
		if s != "is" {
			t.Error("Unexpected response: ", v)
		} else {
			fmt.Println(v)
		}
	}

	// Split and get a word
	v, err = vm.RunString("$strings.SplitBySpaceWordAt(4)")
	if err != nil {
		t.Error(err)
	}
	if v.Export().(string) != "cod." {
		t.Error("Unexpected response: ", v)
	} else {
		fmt.Println(v)
	}
}

func TestToolRegexps(t *testing.T) {

	TEXT := "this is a text\ncod. 80 is a matching value! Cod. 80"

	fmt.Println("TEXT: ", TEXT)

	vm := lygo_scripting.New()
	vm.SetToolContext("$regexps", TEXT)

	var v goja.Value
	var err error

	// IndexFirst
	v, err = vm.RunString("$regexps.IndexFirst('?od??80')")
	if err != nil {
		t.Error(err)
	} else {
		i, _ := v.Export().(int)
		if i == -1 {
			t.Error("Unexpected response: ", v)
		} else {
			fmt.Println("IndexFirst: ", v)
		}
	}

	// Index
	v, err = vm.RunString("$regexps.Index('?od??80')")
	if err != nil {
		t.Error(err)
	} else {
		a, _ := v.Export().([]int)
		if len(a) == 0 {
			t.Error("Unexpected response: ", v)
		} else {
			fmt.Println("Index: ", v)
		}
	}

	// IndexLenPair
	v, err = vm.RunString("$regexps.IndexLenPair('?od??80')")
	if err != nil {
		t.Error(err)
	} else {
		a, _ := v.Export().([][]int)
		if len(a) == 0 {
			t.Error("Unexpected response: ", v)
		} else {
			fmt.Println("IndexLenPair: ", v)
		}
	}

}

func TestToolArrays(t *testing.T) {

	ARRAY := []interface{}{1, "hello", 34.7}

	fmt.Println("ARRAY: ", ARRAY)

	vm := lygo_scripting.New()
	vm.SetToolContext("$arrays", ARRAY)

	var v goja.Value
	var err error

	// GetAt
	v, err = vm.RunString("$arrays.GetAt(1)")
	if err != nil {
		t.Error(err)
	} else {
		s, _ := v.Export().(string)
		if s != "hello" {
			t.Error("Unexpected response: ", v)
		} else {
			fmt.Println("GetAt: ", v)
		}
	}

	// GetLast
	v, err = vm.RunString("$arrays.GetLast(1)")
	if err != nil {
		t.Error(err)
	} else {
		f, _ := v.Export().(float64)
		if f != 34.7 {
			t.Error("Unexpected response: ", v)
		} else {
			fmt.Println("GetLast: ", v)
		}
	}

	// GetSub
	v, err = vm.RunString("$arrays.GetSub(1,2)")
	if err != nil {
		t.Error(err)
	} else {
		a, _ := v.Export().([]interface{})
		if len(a) < 2 {
			t.Error("Unexpected response: ", v)
		} else {
			fmt.Println("GetSub: ", v)
		}
	}

	// GetSub all values
	v, err = vm.RunString("$arrays.GetSub()")
	if err != nil {
		t.Error(err)
	} else {
		a, _ := v.Export().([]interface{})
		if len(a) < 3 {
			t.Error("Unexpected response: ", v)
		} else {
			fmt.Println("GetSub: ", v)
		}
	}
}

func TestToolRegexpsArrays(t *testing.T) {

	TEXT := `BANCA CENTRALE jg TiITUTO PER LA SICUREZZA SOCIALE Cod. 70
DELLA REPUBBLICA DI SAN MARINO
San Marino, 30/04/2019
Denominazione o Ragione Sociale:
LA ddddd DEL dddddSRL C.O.E.: ddddddd
Area Causale Versamento -------- Descrizione Versamento
251 300 --- VERSAMENTO MENSILE CONTRIBUTI PREVIDENZIALI LAVORATORI DIPENDENTI
Importo (in cifre): Rif. Mese: Aprile (4) +
1.234,4
Rif. Anno: 2019 °
Importo (in lettere/centesimi in cifre):
Millesessantadue/73
\l pagamento deve essere effettuato addebitando il conto corrente:
ea ee ee. ee RS aa Agenzia: Addebito cic:
Spazio riservato alla quietanza dellaBanca see
eh ee 28 eT Meee Bota RNA
|
`

	fmt.Println("MIXED: \n", TEXT)

	vm := lygo_scripting.New()
	vm.SetToolContext("$regexps", TEXT)
	vm.SetToolContext("$strings", TEXT)

	var v goja.Value
	var err error

	EXPRESSION := ` // lookup price
(function(){
	var START_PATTERN = '?ese';
	var END_PATTERN = '?if? ';
	var pair_start = $arrays.GetFirst($regexps.IndexLenPair(START_PATTERN)); // lookup mese
	if (!!pair_start && pair_start.length>0){
		startIndex = pair_start[0] + pair_start[1] 
		endIndex = $arrays.GetFirst($regexps.IndexStartAt(startIndex, END_PATTERN))			
		
		if (!!endIndex){
			// lookup values in a range
			var sub = $strings.Sub(startIndex, endIndex)
			// deep string analysis
			if (sub.indexOf(')')>-1) {
				sub = sub.substring(sub.indexOf(')'), sub.length-1);
			}			

			// get values
			var values = $regexps.MatchNumbers(sub);				

			return $arrays.GetLast(values);
		}
	}
	return '';
})()
`
	v, err = vm.RunString(EXPRESSION)
	if err != nil {
		t.Error(err)
	} else {
		s, _ := v.Export().(string)
		if len(s) == 0 {
			t.Error("Unexpected response: ", v)
		} else {
			fmt.Println("EXPRESSION: ", v)
		}
	}

	// one line code (same as above)
	EXPRESSION = "(function(){var START_PATTERN='?ese';var END_PATTERN='?if? ';var pair_start=$arrays.GetFirst($regexps.IndexLenPair(START_PATTERN));if(!!pair_start&&pair_start.length>0){startIndex=pair_start[0]+pair_start[1]\nendIndex=$arrays.GetFirst($regexps.IndexStartAt(startIndex,END_PATTERN))\nif(!!endIndex){var sub=$strings.Sub(startIndex,endIndex)\nif(sub.indexOf(')')>-1){sub=sub.substring(sub.indexOf(')'),sub.length-1)}\nvar values=$regexps.MatchNumbers(sub);return $arrays.GetLast(values)}}\nreturn''})()"

	v, err = vm.RunString(EXPRESSION)
	if err != nil {
		t.Error(err)
	} else {
		s, _ := v.Export().(string)
		if len(s) == 0 {
			t.Error("Unexpected response: ", v)
		} else {
			fmt.Println("EXPRESSION (one line code): ", v)
		}
	}
}

func TestNodeJs(t *testing.T) {

	registry := new(require.Registry) // this can be shared by multiple runtimes

	runtime := goja.New()
	req := registry.Enable(runtime)

	v, error := runtime.RunString(`
    var m = require("m.js");
    m.test();
    `)
	if nil != v {
		fmt.Println("VALUE: ", v)
	}
	if nil != error {
		fmt.Println("JS ERROR: ", error)
	}

	m, err := req.Require("m.js")
	_, _ = m, err
	if nil != err {
		fmt.Println("ERROR: ", err)
	}
}

func TestToolCSV(t *testing.T) {

	vm := lygo_scripting.New()

	TEXT, err := lygo_io.ReadTextFromFile("./script_csv.js")
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	v, err := vm.RunString(TEXT)
	if err != nil {
		panic(err)
	}

	//value := v.Export().([]map[string]string)
	fmt.Println(v)

}
