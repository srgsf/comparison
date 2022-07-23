<?php
//4.1 sec
declare(strict_types = 1);

$time_start = microtime(true);

//composer require "amphp/amp", "amphp/file"
require __DIR__ . '/vendor/autoload.php';

ini_set('display_errors', '1');
ini_set('error_reporting', (string) E_ALL);
date_default_timezone_set('Europe/Lisbon');

$wordCount = 0;
$words = [];

function processResult($result) {
    global $wordCount, $words;
    $wordCount++;
    if ($result) {
        $words[] = $result;
    }
}

function processDirectory(string $fullPathDirectory, callable $processingFunction) {
    //TODO read directory asynchronously (there was a problem with Amp\File\listFiles)
    $files = glob($fullPathDirectory . DIRECTORY_SEPARATOR . '*');
    return Amp\Promise\all(array_map(
        function ($fullPathFile) use ($processingFunction) {
            return Amp\call(function () use ($processingFunction, $fullPathFile) {
                $result = yield ($processingFunction($fullPathFile));
                processResult($result);
            });
        }, $files
    ));
}

function readTheFile(string $fullPathFile) {
    return Amp\call(function () use ($fullPathFile) {
        if (yield Amp\File\isFile($fullPathFile)) {
            return Amp\File\read($fullPathFile);
        }
        return false;
    });
}

function readFirstLine(string $fullPathFile) {
    return Amp\call(function () use ($fullPathFile) {
        if (yield Amp\File\isFile($fullPathFile)) {
            //TODO make code below asynchronous (no read line by line function)
            if ($filePointer = fopen($fullPathFile, 'rb')) {
                $string = fgets($filePointer);
                fclose($filePointer);
                return $string;
            }
        }
        return false;
    });
}

function getWord($fullPathDirectory) {
    return Amp\call(function () use ($fullPathDirectory) {
        if (798 > yield readTheFile($fullPathDirectory . DIRECTORY_SEPARATOR . '2015-2017-spoken-frequency.txt')) {
            return '';
        }

        $word = yield readFirstLine($fullPathDirectory . DIRECTORY_SEPARATOR . 'translation.txt');
        if ($word) {
            $word = trim($word);
        } else {
            $word = basename($fullPathDirectory);
        }
        //echo $word, PHP_EOL;
        return $word;
    });
}

$fullPathDirectory = realpath(__DIR__ . DIRECTORY_SEPARATOR . 'data' . DIRECTORY_SEPARATOR . 'a');
Amp\Loop::run(function () use ($fullPathDirectory) {
    yield processDirectory($fullPathDirectory, 'getWord');
    Amp\Loop::stop();
});

echo 'selected ', count($words), ' words from ', $wordCount, ', took ', round(microtime(true) - $time_start, 3), ' seconds', PHP_EOL;
/*
sort($words);
$i = 0;
foreach ($words as $word) {
    $i++;
    echo $i, ' ', $word;
}
*/
