#!/usr/bin/env sh

find  ./values/  -name "*.yaml" -exec  cue vet ./values.schema.cue {} \;
