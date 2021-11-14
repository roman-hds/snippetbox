APP?==snippetbox

hello:
	echo "Hello"

clean:
	rm -f ${APP}

build: clean
	go build -o ${APP} ./cmd/web/ 
run: build
	./${APP}