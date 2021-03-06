'use strict';

var esprima = require('esprima/dist/esprima.js');

var lang = "javascript";
var version = "1.0.0";

function isClassDecl(stat) {
	if (!stat || stat === null) {
		return false;
	}
	if (stat.type === "ClassDeclaration") {
		return true;
	}
	return false;
}

function isExport(stat) {
	if (!stat || stat === null) {
		return false;
	}
	if (stat.type === "AssignmentExpression" && stat.left.type === "MemberExpression" &&
			stat.left.object.type === "Identifier" && stat.left.object.name === "module" &&
			stat.left.property.type === "Identifier" && stat.left.property.name === "exports") {
		return true;
	}
	return false;
}

function getExportName(stat) {
	if (stat.right.type != "Identifier") {
		throw "module.exports should be assigned to an identifier";
	}
	return stat.right.name;
}

function isPublicMethod(def) {
	return def.key.type === "Identifier" && def.value.type === "FunctionExpression" && !def.key.name.startsWith("_");
}

function genAbi(def) {
	var abi = {
		"name": def.key.name,
		"args": new Array(def.value.params.length).fill("string")
	};
	return abi;
}

function genAbiArr(stat) {
	var abiArr = [];
	if (!isClassDecl(stat) || stat.body.type !== "ClassBody") {
		console.error("invalid statment for generate abi. stat = " + stat);
		return null;
	}
	var initFound = false;
	for (var i in stat.body.body) {
		var def = stat.body.body[i];
		if (def.type === "MethodDefinition" && isPublicMethod(def)) {
			if (def.key.name == "constructor") {
			} else if (def.key.name == "init") {
				initFound = true;
			} else {
				abiArr.push(genAbi(def));
			}
		}
	}
	if (!initFound) {
		console.error("init not found!");
		return null;
	}
	return abiArr;
}

function processContract(source) {
	var newSource, abi;
    var ast = esprima.parseModule(source, {
		range: true,
		loc: true,
		tokens: true
	});

	var abiArr = [];
	if (!ast || ast === null || !ast.body || ast.body === null || ast.body.length === 0) {
		console.error("invalid source! ast = " + ast);
		return ["", ""]
	}
	var validRange = [];
	var className;
	for (var i in ast.body) {
		var stat = ast.body[i];

		if (isClassDecl(stat)) {
			validRange.push(stat.range);
		}
		else if (stat.type === "ExpressionStatement" && isExport(stat.expression)) {
			validRange.push(stat.range);
			className = getExportName(stat.expression);
		}
	}

	for (var i in ast.body) {
		var stat = ast.body[i];

		if (isClassDecl(stat) && stat.id.type == "Identifier" && stat.id.name == className) {
			abiArr = genAbiArr(stat);
		}
	}

	newSource = 'use strict;\n';
	for (var i in validRange) {
		var r = validRange[i];
    	newSource += source.slice(r[0], r[1]) + "\n";
	}

	var abi = {};
	abi["lang"] = lang;
	abi["version"] = version;
	abi["abi"] = abiArr;
	var abiStr = JSON.stringify(abi, null, 4); 

	return [newSource, abiStr]
}
module.exports = processContract;


var fs = require('fs');

var file = process.argv[2]
fs.readFile(file, 'utf8', function(err, contents) {
	console.log('before calling process, len = ' + contents.length);
	var [newSource, abi] = processContract(contents)
	console.log('after calling process, newSource len = ' + newSource.length + ", abi len = " + abi.length);

	fs.writeFile(file + ".after", newSource, function(err) {
    	if(err) {
    	    return console.log(err);
    	}
    	console.log("The new contract file was saved as " + file + ".after");
	}); 

	fs.writeFile(file + ".abi", abi, function(err) {
    	if(err) {
    	    return console.log(err);
    	}
    	console.log("The new abi file was saved as " + file + ".abi");
	}); 
});
