(function () {

    var TEXT = "CAUSALE PAGAMENTO !2CEVUTA";
    var expressions = "C?U?A?? PAG* *U* | ??CEVUTA"; // 2 espressions

    // get best score
    var score = $regexps.Score(TEXT, expressions, "best");

    return score;

})();
