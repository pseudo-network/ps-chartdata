# ps-chartdata

## Development

```bash
# start application
$ make dev
```

## Deployment

```bash
# create branch
# do work
# bumpversion
$ make bumpversion-patch
# stage, commit, push
$ git add .
$ git commit -m "commit message"
$ git push origin {branch name}
# create PR on the github webapp
# merge PR, CICD only builds and deploys on the master branch
```
