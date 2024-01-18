build:
	go build -o packss

run: build
	./run.sh

clean:
	rm -rf sharder-blocks