[![CircleCI](https://circleci.com/gh/kieran-osgood/update-git-tags/tree/main.svg?style=svg)](https://circleci.com/gh/kieran-osgood/update-git-tags/tree/main)
## update-git-tag
Compares the `app.expo.version` value of the `app.json` file in a repository for changes and pushes a new git tag for the new value if it has changed.

Example:
If the `origin/main@HEAD` commit had a value of
```json
// app.json 
{
  "expo": {
    "version": "1.1.1"
  }
}
```
with `origin/main@HEAD~1` having any other value for version e.g.
```json
// app.json 
{
  "expo": {
    "version": "1.1.0"
  }
}
```
then running 

```
./update-git-tags --project_url git@github.com/kieran-osgood/test-repository --branch main
```

would produce a new git tag `v1.1.1` on `origin/main`.

## Arguments
<!-- https://www.tablesgenerator.com/markdown_tables -->
| name              | description                                                                              | Required | default value   |
|-------------------|------------------------------------------------------------------------------------------|----------|-----------------|
| --project_url     | SSH url to the repository to be cloned e.g. git@github.com/<user-name>/<repository-name> | true     | $REPOSITORY_URL |
| --branch          | Default branch to checkout                                                               | false    | "master"        |
| --previous_commit | Commit hash to compare against <branch>@HEAD                                             | true     | ""              |
| --ssh_key_path    | SSH git authentication method Note: Base64 encoded private keyfile                       | false    | ""              |
| --username        | HTTPS git authentication method                                                          | false    | ""              |
| --password        | HTTPS git authentication method                                                          | false    | ""              |
|                   |                                                                                          |          |                 |
---
