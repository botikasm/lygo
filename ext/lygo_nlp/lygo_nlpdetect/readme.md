
Detect the language of text.


## Supported languages

Supports 175 “languages”. For a complete list, check out [languages list](languages.md)


## Usage

```go
response := lygo_nlpdetect.DetectOne("Votre temps est limité, ne le gâchez pas en menant une existence qui n’est pas la vôtre.")
// response == {Code:"fra" Count:1}

response := lygo_nlpdetect.Detect("Votre temps est limité, ne le gâchez pas en menant une existence qui n’est pas la vôtre.")
// response == [{Code:"fra" Count:1},{spa 0.7709821779068855},{cat 0.7656434011148622},{src 0.7274083379131664}...]
```

