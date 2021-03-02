docker run -it --rm -v "$PWD"/src/app:/go/src/app -p 8080:8080 --env-file=.env go1.15 "$@"
