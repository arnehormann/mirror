#!/bin/bash
<create_tags.template go run create_tags_run.go | gofmt > tags.go
