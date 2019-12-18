all: hege-front hege-world

clean:
	-rm -f hege-front hege-world hege-ticker

.PHONY: all clean test try world hege-front hege-world

world:
	( cd front-server && go build )
hege-front:
	( cd front-server && go build -o ../$@ )
hege-world:
	( cd world-server && go build -o ../$@ )

test:
	for D in world world-server front-server ; do cd $$D && go test -v && cd - ; done
try: hege-front hege-world
	ci/run.sh $$PWD/ci/bootstrap-empty.json
