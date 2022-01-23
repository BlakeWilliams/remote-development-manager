# Remote Development Manager

An experimental Go application that allows an SSH session to interact with the
clipboard of the host machine and forward calls to `open`. RDM works by
listening on a unix socket locally that can be forwarded to an SSH session.

So far the server only works on macOS, but the client commands are not OS
specific.

## Usage

Server commands:

* `rdm server` - hosts a server locally (macOS only) so that your machine can receive copy, paste, and open commands.
* `rdm close` - attempts to close a running server.
* `rdm logpath` - returns the path where server logs are located. useful for `tail $(rdm logpath)`

Client commands:

* `rdm copy` - reads stdin and forwards the input to the host machine, adding it to the clipboard. e.g. `echo "hello world" | rdm copy`
* `rdm paste` - reads and prints the host machine's clipboard. `rdm paste`
* `rdm open` - forwards the first argument to `open`. e.g. `rdm open https://github.com/blakewilliams/remote-development-manager`

## TO-DO

So far this is just an experiment and there's a lot to be done to get it to a stable point. Contributions are very welcome.

* Add test coverage
* Daemonize the server process
* Add a configuration file that allows custom commands
