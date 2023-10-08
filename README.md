# Puzzle Solver

Match patterns and solve anagrams - handy for crossword fanatics.

## Standalone Server

```bash
cd standalone
go run main.go
```

## Cloud Function

To test using the Cloud Functions Framework:

```bash
env FUNCTION_TARGET=WordSearch WORDLIST_BUCKET=word-search-1729-assets \
  WORDLIST_PATH=data/wordlist.txt LOCAL_ONLY=true go run ./cmd/fn-framework-test/main.go

curl 'http://localhost:8080?mode=anagrams&pattern=idea'
```
