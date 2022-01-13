docker-executor:
	docker build -t=kieranosgood/update-git-tags-orb:$(version) ./bin
docker-push:
	docker push kieranosgood/update-git-tags-orb:$(version)