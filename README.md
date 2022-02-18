# Remote Development Manager

An experimental Go application that allows an SSH session to interact with the
clipboard of the host machine and forward calls to `open`. RDM works by
listening on a unix socket locally that can be forwarded to an SSH session.

So far the server only works on macOS, but the client commands are not OS
specific.

## Installation

The easiest way to install rdm is to [download the latest
release](https://github.com/BlakeWilliams/remote-development-manager/releases)
for your platform. Alternatively, you can build it yourself with `go build main.go`.

e.g. for a Linux server you can use `wget` to download the binary then put it somewhere in your `$PATH`:

```
wget https://github.com/BlakeWilliams/remote-development-manager/releases/latest/download/rdm-linux-amd64
mv rdm-linux-amd64 /usr/local/bin/rdm
chmod +x /usr/local/bin/rdm
```

## Usage

The following is an example of forwarding an rdm server to a remote host: `ssh
-R 127.0.0.1:7391:$(rdm socket) user@mysite.net`. It's worth noting the port
number is not currently configurable and will always attempt to connect to
`7391`.

For Codespaces, `rdm` can be forwarded as part of the `gh cs ssh` command as
arguments to `ssh`, e.g.: `gh cs ssh -- -R 127.0.0.1:7391:$(rdm socket)`

Server commands:

* `rdm server` - hosts a server locally (macOS only) so that your machine can receive copy, paste, and open commands.
* `rdm close` - attempts to close a running server.
* `rdm logpath` - returns the path where server logs are located. Useful for `tail $(rdm logpath)`
* `rdm socket` - returns the path where the server socket lives. Useful for SSH commands, as seen above.

Client commands:

* `rdm copy` - reads stdin and forwards the input to the host machine, adding it to the clipboard. e.g. `echo "hello world" | rdm copy`
* `rdm paste` - reads and prints the host machine's clipboard. `rdm paste`
* `rdm open` - forwards the first argument to `open`. e.g. `rdm open https://github.com/blakewilliams/remote-development-manager`
* `rdm run` - runs a custom command defined in `~/.config/rdm/rdm.json`. e.g. `rdm run forward 80:80`
* `rdm ps` - lists processes started by `run`.
* `rdm kill` - kills a process started by `run`.

## Configuration

`rdm` supports custom commands via a configuration file located in `~/.config/rdm/rdm.json`.

Here's a small example defining two custom commands:

```json
{
  "commands": {
    "host_time": {
      "executablePath": "./hosttime.sh"
    },
    "forward": {
      "executablePath": "./forward_ports.sh"
    }
  }
}
```

The `executablePath` can either be an absolute path to an executable or a
relative path that's relative to the config file.

## Integrations

Here's a few tools you can easily hook `rdm` into:

### Tmux

If you're using macOS and are already delegating copy to `pbcopy` you can
easily use `rdm` in an ssh session by creating an alias.

```shell
alias pbcopy="rdm copy"
```

Alternatively, you can define the commands explicitly for `rdm`:

```
bind-key -T copy-mode-vi Enter send -X copy-pipe-and-cancel "rdm copy"
bind-key -T copy-mode-vi 'y' send -X copy-pipe-and-cancel "rdm copy"
```

### Neovim

Neovim supports custom clipboards out-of-the-box. You can use `rdm` with Neovim
using the following code:

```viml
let g:clipboard = {"name": "rdm", "copy": {}, "paste": {}}
let g:clipboard.copy["+"] = ["rdm", "copy"]
let g:clipboard.paste["+"] = ["rdm", "paste"]
let g:clipboard.copy["*"] = ["rdm", "copy"]
let g:clipboard.paste["*"] = ["rdm", "paste"]
```

For `open` support, add the following to `~/.zshenv` if you're using zsh:

```shell
alias open="rdm open"
alias xdg-open="rdm open"
```

## GitHub CLI

GitHub CLI allows you to configure the browser used to open URL's. We can use
this to set `rdm` as the browser target:

```
$ gh config set browser "rdm open"
```

## TO-DO

So far this is just an experiment and there's a lot to be done to get it to a
stable point. Contributions are very welcome.

* Daemonize the server process
* Add a configuration file that allows custom commands
* Add instructions for vim
* Linux support, if anyone wants to add it
