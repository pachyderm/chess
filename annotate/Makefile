ll: compile docker push

compile:
	CGO_ENABLED=0 go build -o annotate

docker:
	sudo docker build --no-cache=true --force-rm=true -t pachyderm/chess-annotate .

push:
	sudo docker push pachyderm/chess-annotate

