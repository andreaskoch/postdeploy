# postdeploy

postdeploy listens for deployment requests and executes a custom command when the request arrives.

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
            "provider": "bitbucket",
            "route": "magento-modules",
            "directory": "/var/vhosts/example.com/htdocs",
            "command": "modman update-all"
        }
    ]
}
```

When using the above sample configuration, postdeploy will execute the command `modman update-all` in the folder `/var/vhosts/example.com/htdocs` every time a POST request is sent to `http://127.0.0.1:80/deploy/bitbucket/magento-modules` (assuming you bound postdeploy to "127.0.0.1:80").

## Build Status

[![Build Status](https://travis-ci.org/andreaskoch/postdeploy.png?branch=master)](https://travis-ci.org/andreaskoch/postdeploy)

## Contribute

If you have an idea how to make this little tool better please send me a message or a pull request.

All contributions are welcome.