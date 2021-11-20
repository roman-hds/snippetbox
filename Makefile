APP?==snippetbox

clean:
	rm -f ${APP}

build: clean
	go build -o ${APP} ./cmd/web/ 
run: build
	./${APP}