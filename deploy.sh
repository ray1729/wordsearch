#!/bin/bash

set -e 

V=$(git describe --tags)

zip function.zip go.mod go.sum cloudfn.go match/match.go anagram/anagram.go util/util.go data/*.bin

gsutil cp function.zip gs://word-search-1729-assets/cloudfn/${V}/