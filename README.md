# collectionDB

![image](static/logo.svg)

A simple webapp for managing your collections of physical media written in Go with SQLite.

## Usage

### Configure

The `config.yml` contains basic configuration options:

```
{
"listen":     "127.0.0.1:8080", ### Listening Adress and Port
"database":   "/var/lib/collectionDB/collections.db", ### Location where the SQLite file will be created
"static   ":  "/etc/collectionDB/static", ### Location of the static files
"debug":      false, ### set to true to log debug output
"proxy":      false, ### set to true to only allow traffic from localhost (for running behind a reverse proxy)
"tls":        false, ### set to true to activate HTTPS (required `cert` and `key`)
"tls_listen": "127.0.0.1:8888", ### Listening Adress and Port for TLS
"cert":       "", ### Location of the SSL certifiate
"key":        "" ### Location of the SSL key
}
```

Restart collectionDB after editing config.yml

### From Source

1. `go build main.go`
2. edit `config.yml`
3. `./main`

### RPM Package (Tested in Fedora 42)

1. Download the latest RPM Package.
2. `sudo dnf install collectionDB-<VERSION>.rpm`
3. edit `config.yml` under `/etc/collectionDB/config.yml`
4. `sudo systemctl start collectionDB.service`