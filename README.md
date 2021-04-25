# Inplace Environment Substitution for Files

Simple github action that allows inplace environment substitution for files

## How it works

The action scans all your files in your "working directory-directory", based on the "search_input" provided (see bellow
for more details)
If the file files matches your search input, the files are replaced in place based on the environment variables set for
the current workflow / job / step

The environment variables placeholder must follow the following pattern ${PLACEHOLDER_KEY}

Some notes:

- Nested environment placeholders do not work. For `${KEY${KEY}}`, the environment variable key translated will
  be `KEY${KEY}`
- Currently, does not support default values (For example `${KEY:-defaultvalue}`)

## Inputs

| Name                        | Default | Description |
|-----------------------------|---------|-------------|
| working-directory           |         | The root directory where the file search and environment substitution takes place |
| fail_on_missing_variables   |  false  | If set to "true" the job fails if a placeholder could not be resolved to an env variable. If set to false, action will ignore the non resolvable placeholder |
| replace_in_parallel         |  false  | If set to true, for each file matching the regex, the env substitution will be executed in parallel, spawning a go routine for each |
| search_input                |         | Json string defining the search input for the files wished to be substituted. Example `{"patterns": [".+.yaml"],"files": ["replace.me"],"depth": 2}`

## How to use:

Example on how to use this action provided in this workflow: https://github.com/Slidem/inplace-envsubst-action/blob/master/.github/workflows/test.yaml

## The Search Input

The `search_input` JSON has 3 fields:
- patterns: Array of regex patterns to match the files to be substituted. If any of the regex matches the files, the files will be substituted.
- files: Array of strings, for exact match of files to be replaced
- depth: (!mandatory) Defines how many levels deep the recursive search will go, starting from your `working-directory`
