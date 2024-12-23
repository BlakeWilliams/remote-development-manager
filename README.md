# Remote Development Manager

An experimental Go application that allows an SSH session to interact with the
clipboard of the host machine and forward calls to `open`. RDM works by
listening on a unix socket locally that can be forwarded to an SSH session.

The server works on MacOS and Linux, but the client commands are not OS
specific.

## Installation

### Homebrew Casks:

1. Tap the cask:
   ```bash
   brew tap BlakeWilliams/remote-development-manager
   ```
2. Install `rdm`:
   ```bash
   brew install --cask rdm
   ```

### Manual:

[Download the latest
release](https://github.com/BlakeWilliams/remote-development-manager/releases)
for your platform.

e.g. for a Linux server you can use `wget` to download the binary then put it somewhere in your `$PATH`:

```
wget https://github.com/BlakeWilliams/remote-development-manager/releases/latest/download/rdm-linux-amd64
mv rdm-linux-amd64 /usr/local/bin/rdm
chmod +x /usr/local/bin/rdm
```

### Build from source:

`go build main.go`.

### Mac daemon installation

If you are running the server on MacOS you can set up rdm as a
[launchd](https://www.launchd.info/) service that will automatically start on
system boot:

```
$ rdm service install
Run state:     [Running] done!
Run `launchctl print gui/501/me.blakewilliams.rdm` for more detail.
Configured to start at boot. Uninstall using:
        rdm service uninstall
```

### Linux daemon installation with systemd

If you use systemd you can easily daemonize `rdm` for your user with:

```sh
systemctl edit --user --force --full rdm.service
```

That will open your `$EDITOR` and you should fill the file with something like:

```systemd
[Unit]
Description=Remote Development Manager

[Service]
ExecStart=/path/to/rdm server
ExecStop=/path/to/rdm stop
```

Once that is done you can use `systemctl --user start rdm` or `systemctl --user stop rdm`
when you want to start/stop the daemon or `systemctl --user enable rdm` so that it is
automatically enabled for your user when it logs in.

## Usage

The following is an example of forwarding an rdm server to a remote host: `ssh
-R 127.0.0.1:7391:$(rdm socket) user@mysite.net`. It's worth noting the port
number is not currently configurable and will always attempt to connect to
`7391`.

For Codespaces, `rdm` can be forwarded as part of the `gh cs ssh` command as
arguments to `ssh`, e.g.: `gh cs ssh -- -R 127.0.0.1:7391:$(rdm socket)`

Server commands:

* `rdm server` - hosts a server locally (macOS only) so that your machine can receive copy, paste, and open commands.
* `rdm stop` - attempts to close a running server.
* `rdm logpath` - returns the path where server logs are located. Useful for `tail $(rdm logpath)`
* `rdm socket` - returns the path where the server socket lives. Useful for SSH commands, as seen above.

Client commands:

* `rdm copy` - reads stdin and forwards the input to the host machine, adding it to the clipboard. e.g. `echo "hello world" | rdm copy`
* `rdm paste` - reads and prints the host machine's clipboard. `rdm paste`
* `rdm open` - forwards the first argument to `open`. e.g. `rdm open https://github.com/blakewilliams/remote-development-manager`

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

Or if you use lua:

```lua
vim.g.clipboard = {
  name = "rdm",
  copy = {
    ["+"] = {"rdm", "copy"},
    ["*"] = {"rdm", "copy"}
  },
  paste = {
    ["+"] = {"rdm", "paste"},
    ["*"] = {"rdm", "paste"}
  },
}
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
