# postdeploy - The deployment hook server

postdeploy is a http service that listens for deployment requests and executes a predefined command when the request arrives

[![Build Status of "postdeploy"](https://travis-ci.org/andreaskoch/postdeploy.png?branch=master)](https://travis-ci.org/andreaskoch/postdeploy)

## Build

If you have [go installed](http://golang.org/doc/install) you can build **postdeploy** yourself:

```bash
git clone git@github.com:andreaskoch/postdeploy.git && \
cd postdeploy && \
make
```

or with `go get`:

```bash
go get github.com/andreaskoch/postdeploy
```

## Running postdeploy with docker

There is an automated build at [Docker Hub](https://hub.docker.com) for **postdeploy**: [andreaskoch/postdeploy](https://hub.docker.com/r/andreaskoch/postdeploy/).

Example: Run postdeploy on port `80`:

```bash
docker run -p 80:7070 andreaskoch/postdeploy:latest
```

## Cross-Compilation

If you want to cross-compile postdeploy for **macOS** (amd64), **Linux** (arm5, arm6, arm7 and amd64) and **Windows** (amd64) you can use the `crosscompile` action of the make script:

```bash
make crosscompile
```

## Usage

**Get help**

```bash
postdeploy --help
```

**Run postdeploy**

For running **postdeploy** you must specify an **ip binding** (e.g. `:7070`, `127.0.0.1:8080` or `[::]:80`) and the **path to a JSON configurtion file** (e.g. [conf/ping-sample.json](conf/ping-sample.json)):

```bash
postdeploy --binding=:7070 -configfile=conf/ping-sample.json
```

**postdeploy** will then spawn a http server and **listen for POST requests** to `/deploy/<provider>/<route>`. When a request arrives it executes the commands that have been configured for given route.

## The Configuration File

The postdeploy configuration has the following JSON structure:

```json
{
    "hooks": [
        {
            "provider": "<provider-name>",
            "route": "some/route",
            "directory": "/the/working/directory",
            "commands": [
                {
                    "name": "<Some command>",
                    "args": [
                        "arg1",
                        "arg2",
                        "..."
                    ]
                }
            ]
        }
    ]
}
```

Assuming you bind postdeploy to port `7070` a POST request to `http://127.0.0.1:7070/deploy/<provider-name>/some/route` will execute the command `<Some command>` in the specified directory `/the/working/directory`.

An example of a postdeploy configuration files can be found here: [conf/ping-sample.json](conf/ping-sample.json).

### Examples

#### A simple ping

Write the current date and time to a log file every time the ping route is executed:

```json
{
    "hooks": [
        {
            "provider": "generic",
            "route": "ping",
            "directory": "",
            "commands": [
                {
                    "name": "bash",
                    "args": [
                        "-c",
                        "echo $(date) >> ping.log"
                    ]
                }
            ]
        }
    ]
}
```

```bash
postdeploy --binding=:7070 -configfile=conf/ping-sample.json
```

```bash
curl -X POST http://127.0.0.1:7070/deploy/generic/ping
```

### Security

**Do not run this tool on mission critical components!**

I've tried to build the system in a way that the worst thing that could happen is that someone who knows your routes can trigger the deployment hook associated with that route. But this can be bad enough. So please keep in mind that this is just a convenience tool and nothing you should deploy to your production servers.

**Routes**

As of now there is **no security**. Everybody who knows your configured routes will be able to trigger the action.

**Code Execution**

The system will only execute the commands specified in the config file and nothing else. Attackers should not be able trigger any other commands than the ones you specified.

## Roadmap

- Providers
    - Add a github provider
- Security
    - Make sure postdeploy executes every hook only once every x seconds.
    - Block IP addresses that try to find deployment hooks (e.g. block after 3 attempts)

## Contribute

If you have an idea how to make this little tool better please send me a message or a pull request.

All contributions are welcome.
