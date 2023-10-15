# lyswc
My own version of wc command line tool

### Just for Golang practicing! Inspired by [Coding Challenges](https://codingchallenges.fyi/challenges/challenge-wc)

#### To use this command line tool
1. Install Go
2. Run `go build`
3. OR we could run `go build -o <output_name>` to specify the output name
4. Move the `<output_name>` to `usr/local/bin`
5. Example of usage
```
# Mimic wc -c
$ lyswc -c test.txt

# Mimic wc -l
$ lyswc -l test.txt

# Mimic wc -w
$ lyswc -w test.txt

# Mimic wc -m
$ lyswc -m test.txt

# No flag passing
$ lyswc test.txt

# Read from standard input
cat test.txt | lyswc
```

#### Unit testing
1. Run `go test`
