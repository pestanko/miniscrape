#!/bin/sh

project_dir="$HOME/src/miniscrape"
repo_url="https://github.com/pestanko/miniscrape.git"
old_dir=$PWD

if [ ! -e "$project_dir" ]; then
  mkdir -p "$project_dir"
  git clone "$repo_url" "$project_dir"
fi


cd "$project_dir"

git pull
go build

cd "$old_dir"