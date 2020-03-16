// Search the formatted string with every information, if found return X
function search_formatted_string(TEXT) {
	// Remove useless
	TEXT = TEXT.replace(/\./g, " ")
		.replace(/\:/g, " ")
		.replace(/\|/g, " ")
		.replace(/\[/g, " ")
		.replace(/\,/g, " ")
		.replace(/\-/g, " ")
		.replace(/\_/g, " ")
		.replace(/\‘/g, " ")
		.replace(/\'/g, " ")
		.replace(/\`/g, " ")
		.replace(/\s+/g, " "); // One whitespace

	// Regex to find THE STRING
	var regex = "[0-9]{2}[ ][0-9]{2,3}[ ][0-9]{2,3}[ ][0-9][ ].+";

	var match = TEXT.match(regex);

	if (!!match) {
		return "X";
	}

	return null;
}

// Search Cod 70
function search_doc_type(TEXT) {
	// Regex to find Cod 70
	var regex = "(.|)(0|O)(D|B)([.]|-|,|[~]| )*(7|])(0|O)";

	var match = TEXT.match(regex);

	if (!!match) {
		$console.Log("70");
		return "70";
	}

	return null;
}

// Search 2 digits after "ENTE"
function search_ente(TEXT) {
	var array = TEXT.split("\n");

	var line = null;

	// RegExp Explained: try to match something similar to ENTE
	var regex = "(E|3|£|€)(N|W|M)(T|1|!|R)(E|3|£|€)";

	for (var i = 0; i < array.length; i++) {
		var match = array[i].match(regex);

		if (!!match) {
			line = array[i];
			break;
		}
	}

	var doc_type = null;

	if (!!line) {
		// Format to find 2 digits
		match = substring.match("[0-9]{2}");

		if (!!match) {
			doc_type = match[0];
		}
	}

	return doc_type;
}

(function() {
	var TEXT = VAR_text;

	TEXT = TEXT.toUpperCase();

	// First, search the string and if found return null, to pass to the next model
	var doc_type = search_formatted_string(TEXT);
	if (!!doc_type) {
		return null;
	}

	// Then, search for Cod 70
	doc_type = search_doc_type(TEXT);
	if (!!doc_type) {
		return doc_type;
	}

	// Then, search for Ente 70
	doc_type = search_ente(TEXT);

	if (!!doc_type) {
		return doc_type;
	}

	return "";
})();
