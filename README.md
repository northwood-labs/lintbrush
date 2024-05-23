# Lint Brush

A linter which focuses on the low-level things that are common across many repositories. Leaves the language-specific linting to more appropriate tools.

I'm moving from [pre-commit] to [lefthook], and don't want to have to deal with a Python installation and particular ways of needing to set up a virtualenv. So I'm rewriting the most relevant rules in Go so they can be compiled into a standalone binary.

> [!CAUTION]
> This is a brand-new project and most of this isn't working yet. Check back later.

## Plan

### Usage

1. Pass files as arguments.

    ```bash
    lint-brush [--flags] $(git ls-files -m)
    ```

1. Pipe a list of files.

    ```bash
    git ls-files -m | lint-brush [--flags]
    ```

1. Pipe a list of files using `NUL` as the separator.

    ```bash
    git ls-files -m -z | lint-brush -0 [--flags]
    ```

    Or all files in a directory, except for the Git directory (using [fd])…

    ```bash
    fd --unrestricted --exclude=.git --prune --type=file . | lint-brush [--flags]
    ```

### Checks/fixers

Will be porting a [number of checks](https://github.com/pre-commit/pre-commit-hooks/tree/main/pre_commit_hooks) from Python → Go for this effort.

* [ ] check: binary files should be stored in Git LFS
* [X] [check: executable text files should have shebangs](https://github.com/pre-commit/pre-commit-hooks/blob/main/pre_commit_hooks/check_executables_have_shebangs.py)
* [ ] [check: find broken synlinks](https://github.com/pre-commit/pre-commit-hooks/blob/main/pre_commit_hooks/check_symlinks.py)
* [ ] [check: GitHub code links should be permalinks](https://github.com/pre-commit/pre-commit-hooks/blob/main/pre_commit_hooks/check_vcs_permalinks.py)
* [ ] [check: JSON/JSONC/JSON5 is well-formed](https://github.com/pre-commit/pre-commit-hooks/blob/main/pre_commit_hooks/check_json.py)
* [ ] [check: large files added](https://github.com/pre-commit/pre-commit-hooks/blob/main/pre_commit_hooks/check_added_large_files.py)
* [ ] [check: merge conflict markers](https://github.com/pre-commit/pre-commit-hooks/blob/main/pre_commit_hooks/check_merge_conflict.py)
* [ ] [check: problems with case-sensitive file systems](https://github.com/pre-commit/pre-commit-hooks/blob/main/pre_commit_hooks/check_case_conflict.py)
* [ ] [check: TOML is well-formed](https://github.com/pre-commit/pre-commit-hooks/blob/main/pre_commit_hooks/check_toml.py)
* [ ] [check: XML is well-formed](https://github.com/pre-commit/pre-commit-hooks/blob/main/pre_commit_hooks/check_xml.py)
* [ ] [check: YAML is well-formed](https://github.com/pre-commit/pre-commit-hooks/blob/main/pre_commit_hooks/check_yaml.py)
* [ ] [fix: destroyed symlinks](https://github.com/pre-commit/pre-commit-hooks/blob/main/pre_commit_hooks/destroyed_symlinks.py)
* [ ] [fix: files should end with a final linebreak](https://github.com/pre-commit/pre-commit-hooks/blob/main/pre_commit_hooks/end_of_file_fixer.py)
* [ ] [fix: mixed line endings](https://github.com/pre-commit/pre-commit-hooks/blob/main/pre_commit_hooks/mixed_line_ending.py)
* [ ] [fix: remove trailing whitespace on lines](https://github.com/pre-commit/pre-commit-hooks/blob/main/pre_commit_hooks/trailing_whitespace_fixer.py)
* [ ] [fix: remove UTF-8 BOM](https://github.com/pre-commit/pre-commit-hooks/blob/main/pre_commit_hooks/fix_byte_order_marker.py)

[fd]: https://github.com/sharkdp/fd
[lefthook]: https://github.com/evilmartians/lefthook
[pre-commit]: https://pre-commit.com
