CC		= go build
BIN		= out


xgo: *.go
	$(CC) -o $@ $^

tmp.s: xgo test/test.go
	./xgo test/test.go > $(BIN)/$@

tmp.out: xgo tmp.s
	gcc -o $(BIN)/$@ $(BIN)/tmp.s


.PHONY: clean
clean:
	rm -rf out/* xgo
