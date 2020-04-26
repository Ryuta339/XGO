CC		= go build


xgo: *.go
	$(CC) -o $@ $^


.PHONY: clean
clean:
	rm -rf out/* xgo
