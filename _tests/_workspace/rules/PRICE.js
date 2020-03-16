function search_price(string) {
	// Prices can be in the following formats
	// RegExp Matches: 1234,00 1234.00 1.000,00
	var regex = "([0-9]+)([.]|,)*([0-9]*|)([.]|,){1}[0-9]{2}";

	var match = string.match(regex);

	var price = null;

	if (!!match) {
		price = match[0];
	}

	if (!!price) {
		// Format standard price: 1234,00 (remove dots)
		price = price.replace(".", "").replace(",", "");
	}

	return price;
}

(function() {
	// -----------------------------------------------------------------
	// Default value if not found
	// -----------------------------------------------------------------
	var DEFAULT = "";

	var TEXT = VAR_text;

	// Remove ":", "|", "[", "_"
	TEXT = TEXT.replace(/\:/g, "")
		.replace(/\|/g, "")
		.replace(/\[/g, "")
		.replace(/\_/g, "")
		.toUpperCase();

	// RegExp Explained: try to match something similar to IMPORTO
	var regex = "(1|I|!)(M|N|H)(P|7|R|F)(0|O)(8|R|P|4|A)(T|1|F|!)(0|O|E) ";

	var match = TEXT.match(regex);

	if (!!match) {
		var index = match.index;

		var price = null;

		if (index > -1) {
			// If match, search for the num pattern in the next lines
			var substring = TEXT.slice(index);

			// select only next 3 lines
			substring = substring.split("\n");
			substring.splice(3);
			substring = substring.join("\n");

			price = search_price(substring);

			if (!!price) {
				return price;
			}
		}
	}

	// If IMPORTO not found, search for a price in all the TEXT
	var price = search_price(TEXT);

	if (!!price) {
		return price;
	}

	return DEFAULT;
})();
