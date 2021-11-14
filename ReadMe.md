[![CircleCI](https://circleci.com/gh/kieran-osgood/update-git-tags/tree/main.svg?style=svg)](https://circleci.com/gh/kieran-osgood/update-git-tags/tree/main)

## update-git-tags

Automates the creation of git-tags for a project based on the value of a version code. Supported file formats:

- `.json`
- `.txt`

Compares the value of the version to a branch `HEAD` vs `HEAD~1`, creating a new git tag if its changed. Will compare
versions via rules of SEMVER, if it is a lower version you will receive a warning.

---

### Usage
Using the executable:
```
./update-git-tags --project_url git@github.com/kieran-osgood/test-repository --branch main
```

Assuming `origin/main@HEAD` commit had a value of

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
would produce a new git tag `v1.1.1` on `origin/main`.

## Arguments

<!-- https://www.tablesgenerator.com/markdown_tables -->

| Name              | Description                                                                              | Required | Default Value   |
|-------------------|------------------------------------------------------------------------------------------|:--------:|-----------------|
| --RepositoryUrl   | SSH url to the repository to be cloned e.g. git@github.com/<user-name>/<repository-name> |     x    | $REPOSITORY_URL |
| --Branch          | Branch to check.                                                                         |          |      "main"     |
| --previous_commit | Commit hash to compare against <branch>@HEAD                                             |     x    |        ""       |
| --SshKey          | Base64 encoded of SSH private key.                                                       |     x    |     $SSH_KEY    |
| --PreviousHash    | Commit hash of the previous commit to HEAD.                                              |          |  "$CIRCLE_SHA1" |
| --PropertyPath    | Property path to the Version code in the json file.                                      |          |    "Version"    |
| --FilePath        | File path to the json file with the Version code.                                        |          |  "package.json" |
| --VersionPrefix   | Prefix for the git tag                                                                   |          |       "v"       |
| --VersionSuffix   | Suffix for the git tag                                                                   |          |        ""       |
|                   |                                                                                          |          |                 |
| --Version         | Check Version of app binary                                                              |          |                 |
| --Help            | View command documentation                                                               |          |                 |

