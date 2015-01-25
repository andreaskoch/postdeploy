# postdeploy

postdeploy is a http service that listens for deployment requests and executes a predefined command when the request arrives.

## Build Status

[![Build Status](https://travis-ci.org/andreaskoch/postdeploy.png?branch=master)](https://travis-ci.org/andreaskoch/postdeploy)

## Build

If you have [go installed](http://golang.org/doc/install) you can run the `make.go` script to build postdeploy yourself:

```bash
git clone git@github.com:andreaskoch/postdeploy.git && cd postdeploy
go run make.go install
```

## Cross-Compilation

If you want to cross-compile postdeploy for different platforms and architectures you can do so by using the `-crosscompile` flag for the make script (if you have [docker](https://www.docker.com) >= 1.4 installed):

```bash
go run make.go -crosscompile
```

This command will launch a [docker container with go 1.4](https://registry.hub.docker.com/u/library/golang/) in it that is prepared for cross-compilation and build postdeploy for you. The output will be available in the `bin` folder of this project.

## Usage

For running postdeploy you must specify an ip binding (e.g. "127.0.0.1:80") and the path to a JSON configurtion file (e.g. "postdeploy.conf.js"):

```bash
postdeploy -binding ":7070" -config "postdeploy.conf.js"
```

postdeploy will spawn a http server and listen for POST requests to `/deploy/<provider-name>/<route-name>` and will then execute the commands that have been configured for this route.

## The Configuration File

The postdeploy configuration has the following JSON structure:

```json
{
    "hooks": [
        {
            "provider": "<provider-name>",
            "route": "some/route",
            "directory": "/the/working/directory",
            "command": "<Some command>"
        }
    ]
}
```

Assuming you bind postdeploy to port `7070` a POST request to `http://127.0.0.1:7070/deploy/<provider-name>/some/route` will execute the command `<Some command>` in the specified directory `/the/working/directory`.

### Examples

#### Automated Magento Updates using Bitbucket and modman

```json
{
    "hooks": [
        {
            "provider": "bitbucket",
            "route": "magento-modules",
            "directory": "/var/vhosts/example.com/htdocs",
            "command": "modman update-all"
        }
    ]
}
```

When using the above sample configuration, postdeploy will execute the command `modman update-all` in the folder `/var/vhosts/example.com/htdocs` every time a POST request is sent to `http://127.0.0.1:80/deploy/bitbucket/magento-modules` (assuming you bound postdeploy to "127.0.0.1:80").

### Security

**Do not run this tool on mission critical components!**

I've tried to build the system in a way that the worst thing that could happen is that someone who knows your routes can trigger the deployment hook associated with that route. But this can be bad enough. So please keep in mind that this is just a convienience tool and nothing you should deploy to your production servers.

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