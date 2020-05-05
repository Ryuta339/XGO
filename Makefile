CC		= go build


xgo: *.go
	$(CC) -o $@ $^

out/tmp.s: xgo test/test.go
	./xgo test/test > $@

out/tmp.out: out/tmp.s
	gcc -o $@ $^


.PHONY: clean
clean:
	rm -rf out/* xgo
