build:
	go build -v -o bin/http cmd/http/*.go

run: build
	bin/http

watch:
	reflex -s -r "\.(go|json|html)$$" --decoration=none make run
