docker-executor:
	docker build -t=kieranosgood/update-git-tags-orb:$(version) ./bin
docker-push:
	docker push kieranosgood/update-git-tags-orb:$(version)
b64-ssh-key:
	base64 ~/.ssh/kieran-osgood-GitHub