.PHONY: build-image
build-image:
	docker build -t http-crash .

.PHONY: run-container
run-container:
	docker run --rm -p 8080:8080 http-crash
