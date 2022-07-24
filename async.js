'use strict';
//1.5 sec
const start = Math.floor(new Date().getTime());

const path = require('node:path');
const fs = require('node:fs');
const fsPromises = require('node:fs/promises');
const readlinePromises = require('node:readline/promises');

let wordCount = 0;
let words = [];

function processResult(result) {
	wordCount++;
	if (result) {
		words.push(result);
	}
}

async function processDirectory(fullPathDirectory, processingFunction) {
	const files = await fsPromises.readdir(fullPathDirectory);
	return Promise.all(files.map(
		async (file) => {
			const fullPathFile = path.join(fullPathDirectory, file);
			let result = await processingFunction(fullPathFile);
			processResult(result);
		}
	));
}

async function readFile(fullPathFile) {
	const stats = await fsPromises.stat(fullPathFile);
	if (stats.isFile()) {
		return fsPromises.readFile(fullPathFile);
	} 
	return Promise.resolve(false);
}

async function readFirstString(fullPathFile) {
	const stats = await fsPromises.stat(fullPathFile);
	if (stats.isFile()) {
		const readInterface = readlinePromises.createInterface({
			input: fs.createReadStream(fullPathFile),
		});

		for await (const line of readInterface) {
			readInterface.close();
			return line;
		}
	}
	return false;
}

async function getWord(fullPathDirectory) {
	if (798 > await readFile(path.join(fullPathDirectory, '2015-2017-spoken-frequency.txt'))) {
		return '';
	}

	let word = '' + await readFirstString(path.join(fullPathDirectory, 'translation.txt'));
	if (word) {
		word = word.trim();
	} else {
		word = path.basename(fullPathDirectory);
	}
	//console.log(word);
	return word;
}

const directory = path.join(__dirname, 'data', 'a');
(async () => {
	await processDirectory(directory, getWord);
	console.log('selected ' + words.length + ' words from ' + wordCount + ', took ' + (Math.floor(new Date().getTime()) - start) / 1000 + ' seconds');
	/*
	let i = 0;
	for (let word of words.sort()) {
		i++;
		console.log('' + i + ' ' + word);
	}
	*/
})();
