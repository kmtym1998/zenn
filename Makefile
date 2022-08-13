build:
	cd tools && \
		go build -o ../tmp/zenn-tool

ls-remote:
	./tmp/zenn-tool ls-remote
