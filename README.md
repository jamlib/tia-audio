# tia-audio

Download, transcode to mp3, embed album artwork, & tag lossless audio from The Internet Archive (archive.org).

Transcodes to `320 kbps` or `V0` mp3 using `ffmpeg` & `libmp3lame`.

## Usage

```
Usage: tia-audio [--quality QUALITY] [--dir DIR] URL

Positional arguments:
  URL                    archive.org details url

Options:
  --quality QUALITY      mp3 quality: 320, V0
  --dir DIR              directory where files will be saved
  --help, -h             display this help and exit
  --version              display version and exit
```

## Developing

[Install go](https://golang.org/doc/install).

Fork this repo to your github account.

Clone repo to your github src GOPATH, run:

    gituser='YOUR-GITHUB-USERNAME'
    cd $GOPATH/src/github.com && if [ ! -d $gituser ]; then mkdir $gituser; fi && cd $gituser
    git clone git@github.com:$gituser/tia-audio.git && cd tia-audio

### Building

From within source path, run:

    go build

The binary will build to the current directory. To test by displaying usage, run:

    ./tia-audio --help

### Submitting a Pull Request

From within source path, create a new branch to use for development, run:

    git checkout -b new-branch

Make your changes, add, commit and push to Github, then back on Github, submit pull request.

## License

This code is available open source under the terms of the [MIT License](http://opensource.org/licenses/MIT).