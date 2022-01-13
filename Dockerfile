FROM cimg/go:1.17.3

USER root

COPY update-git-tags /usr/bin/update-git-tags
