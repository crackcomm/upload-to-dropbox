# upload-to-dropbox

Uploads to dropbox from file or stdin.

## Usage

```sh
export DROPBOX_APP_ID={...}
export DROPBOX_APP_SECRET={...}
export DROPBOX_APP_TOKEN={...}

cat main.go | upload-to-dropbox -dir mydir -filename my_prog_1.go
```

### Environment

```
DROPBOX_DIR         - Dropbox Directory
DROPBOX_APP_ID      - Dropbox Application Id
DROPBOX_APP_SECRET  - Dropbox Application Secret
DROPBOX_APP_TOKEN   - Dropbox Application Token
```

### Flags

```
  -input string
        Input file to upload create (stdin used by default)
  -filename string
        File name to create (required)
  -dir string
        Directory name (required, in env: DROPBOX_DIR)
  -mkdir
        Create the directory
  -token string
        Dropbox App Token (in env: DROPBOX_APP_TOKEN)
  -appId string
        Dropbox App ID (required, in env: DROPBOX_APP_ID)
  -appSecret string
        Dropbox App Secret (required, in env: DROPBOX_APP_SECRET)
  -chunk-size int
        Upload chunk size in megabytes (default 64)
```

## Thanks

Thanks to [github.com/stacktic/dropbox](http://github.com/stacktic/dropbox).

## License

[Apache License Version 2.0](http://www.apache.org/licenses/)
