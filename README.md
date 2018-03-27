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

### Install Go on Linux

Download latest go binary from [golang.org/dl](https://golang.org/dl/). In this case, version 1.10.

Extract to `/usr/local`, run:

    sudo tar -C /usr/local -xzf go1.10.linux-amd64.tar.gz

Open ~/.profile for editing, run:

    nano ~/.profile

Append the following, then save/exit:

    export PATH=$PATH:/usr/local/go/bin
    export GOPATH=$HOME/go
    export PATH=$PATH:$GOPATH/bin

Source updated profile and ensure go $HOME/go exists, run:

    source ~/.profile
    if [ ! -d $HOME/go ]; then mkdir $HOME/go; fi

### Install Dep

On Linux, to install `dep` for dependency management, run:

    curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

For other OS/environments, refer to the [Dep Installation Documentation](https://golang.github.io/dep/docs/installation.html).

### Building

Get latest source, run:

    go get github.com/JamTools/tia-audio

Navigate to source path, run:

    cd $GOPATH/src/github.com/JamTools/tia-audio

Ensure dependencies are installed and up-to-date with `dep`, run:

    dep ensure

From within source path, run:

    go build

The binary will build to the current directory. To test by displaying usage, run:

    ./tia-audio --help

### Testing

From within source path, run:

    go test -cover -v ./...

### Submitting a Pull Request

Fork repo on Github.

From within source path, setup new remote, run:

    git remote add myfork git@github.com:$GITHUB-USERNAME/tia-audio.git

Create a new branch to use for development, run:

    git checkout -b new-branch

Make your changes, add, commit and push to your Github fork.

Back on Github, submit pull request.

## License

This code is available open source under the terms of the [MIT License](http://opensource.org/licenses/MIT).