#!/bin/sh

project_dir="$HOME/src/miniscrape"
repo_url="https://github.com/pestanko/miniscrape.git"
old_dir=$PWD

if [ ! -e "$project_dir" ]; then
  mkdir -p "$project_dir"
  git clone "$repo_url" "$project_dir"
fi


cd "$project_dir" || die "Unable to change directory to the $project_dir"

git stash
git pull
git stash apply


docker compose up --profile full --build -d --wait

cd "$old_dir" || die "Unable to return the $old_dir"