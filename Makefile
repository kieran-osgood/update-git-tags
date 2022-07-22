docker-executor:
	docker build -t=kieranosgood/update-git-tags-orb:$(version) ./bin
docker-push:
	docker push kieranosgood/update-git-tags-orb:$(version)
b64-ssh-key:
	base64 ~/.ssh/kieran-osgood-GitHub
git-tag:
# run via husky post commit?
	git tag -a "v$(shell cat version.txt)" -m "Release version: v$(shell cat version.txt)"
	git push --tags