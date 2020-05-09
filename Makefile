CC		= go build
BIN		= out


xgo: *.go
	$(CC) -o $@ $^

tmp.s: xgo test/test.go
	./xgo test/test > $(BIN)/$@

tmp.out: $(BIN)/tmp.s
	gcc -o $(BIN)/$@ $^


.PHONY: clean
clean:
	rm -rf out/* xgo
