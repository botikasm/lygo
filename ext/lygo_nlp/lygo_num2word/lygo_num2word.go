package lygo_num2word

import (
	"github.com/botikasm/lygo/ext/lygo_nlp/lygo_num2word/lygo_num2word_languages"
	"strings"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e
//----------------------------------------------------------------------------------------------------------------------

type Num2Word struct {
	Options *Num2WordOptions
}

type Num2WordOptions struct {
	WordSeparator string
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewNum2Word() *Num2Word {
	instance := new(Num2Word)
	instance.Options = new(Num2WordOptions)
	instance.Options.WordSeparator = " "

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *Num2Word) ConvertDefault(input int) string {
	lang := lygo_num2word_languages.Languages.Default()
	return num2Word(input, &lang, instance.Options)
}

func (instance *Num2Word) Convert(input int, langCode string) string {
	lang := lygo_num2word_languages.Languages.Lookup(langCode)
	if nil == lang {
		lang = lygo_num2word_languages.Languages.Lookup("en-us")
	}
	return num2Word(input, lang, instance.Options)
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func num2Word(input int, lang *lygo_num2word_languages.Language, options *Num2WordOptions) string {
	response:=""
	if len(lang.Name) > 0 && nil != lang.IntegerToWords {
		response = lang.IntegerToWords(input)
	}

	if nil!=options{
		response = strings.Replace(response, " ", options.WordSeparator, -1)
	}

	return response
}
