# zshist

History manager for zsh shell with no dependencies.

Merges and deduplicates history entries from
`~/.zsh_history`, `~/.zsh_history.pre-oh-my-zsh`, `~/.bash_history`
and saves back to `~/.zsh_history`.

Creates backup in `~/.zsh_history.bak` in case you want to restore.

The order of commands is maintained so latest ones are suggested first.


## why

It may save disk space and redundancy but most importantly you get a
specific completion only once as you scroll up in the `zsh` prompt.


## install

```sh
go install github.com/adhocore/zshist
```

Or, you can also download latest prebuilt binary from
[release](https://github.com/adhocore/zshist/releases/latest) for platform of your choice.


## usage

```sh
zshist

# custom home dir (no trailing `/`)
zshist -home /home/user
```


## sample output

```
$ zshist

Parsed and merged 3 files with 4072 commands
Saved into ~/.zsh_history with 770 commands
Backed up into ~/.zsh_history.bak
```

> 4072 commands down to 770, not bad!
