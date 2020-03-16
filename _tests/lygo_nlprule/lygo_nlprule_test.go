package lygo_nlprule

import (
	"fmt"
	"github.com/botikasm/lygo/_tests"
	"github.com/botikasm/lygo/base/lygo_io"
	"github.com/botikasm/lygo/base/lygo_paths"
	"github.com/botikasm/lygo/ext/lygo_nlp/lygo_nlprule"
	"testing"
)

func TestSimple(t *testing.T) {

	_tests.InitContext()

	file := lygo_paths.WorkspacePath("./rules/rule.json")
	config, err := lygo_io.ReadTextFromFile(file)
	if nil != err {
		t.Error(err)
		t.Fail()
	}

	text := "XXXX CCCCC CCCCC —_— XXXXX XXXXXXX XXXXX Cod. 70 \nDELLA REPUBBLICA DI SAN MARINO\nSan Marino, 30/04/2019\nDenominazione o Ragione Sociale:\nXXXXXXX C.O.E.: XXXXXX\nArea Causale Versamento -------- Descrizione Versamento\n251 300 --- VERSAMENTO MENSILE CONTRIBUTI PREVIDENZIALI LAVORATORI DIPENDENTI\nImporto (in cifre): . Rif. Mese: Aprile (4) +\n1.062,73\nRif. Anno: 2019 =\nImporto (in lettere/centesimi in cifre):\nMillesessantadue/73\n\\l pagamento deve essere effettuato addebitando il conto corrente:\nae ee. ee CNS Ry ss Agenzia: Addebito c/c:\nSpazio riservato alla quietanza dellaBanca aes\neh oe P88 eT Mee Bota ee: PIMA\n|\n"

	engine := lygo_nlprule.NewRuleEngine(config)
	if engine.HasConfig() {
		context := map[string]interface{}{
		}
		response := engine.Eval(text, context, -1)
		fmt.Println("elapsed:", response.ElapsedMs)
		for _, item := range response.Items {
			fmt.Println("\t", "score: ", item.Score, "entities: ", len(item.Entities))
			for _, entity := range item.Entities {
				if nil == entity.Errors || len(entity.Errors) == 0 {
					fmt.Println("\t\t", "value: ", entity.Values[0])
				} else {
					fmt.Println("\t\t", "error: ", entity.Errors[0])
				}
			}
		}
	}
}
