# firebase-simple-storage

An unofficial Go based client for Firebase storage that preserves security
(this goes through Firebase and not Google Cloud Storage).

You can use this to upload files while pinned to user authentication credentials (upload in
the name of a user) from backend services or background jobs that may run on a user's machine.

A good use case for this is when you distribute an app that isn't Web, mobile or backend (obviously),
and don't want to compromise app-level credentials.


## Quick Start

Initialize a `Storage` client:

```go
s := Storage{
  Token:        "--user access token, from any login method--",
  RefreshToken: "--same, but refresh token--",
  Bucket:       "your-bucket-id.appspot.com",
  APIKey:       "your-project-api-key",
}
```

And now you can do any of these:

```go
err := s.Refresh()
if err != nil {
  log.Fatal(err)
}
res, err := s.Put("main.go", "test/main.go")
if err != nil {
  log.Fatal(err)
}
fmt.Printf("%v", res)

res, err = s.Object("test/main.go")
if err != nil {
  log.Fatal(err)
}
fmt.Printf("object: %v", res)

err = s.Download("test/main.go", "out.go")
if err != nil {
  log.Fatal(err)
}
```


# Contributing

Fork, implement, add tests, pull request, get my everlasting thanks and a respectable place here :).


# Copyright

Copyright (c) 2014 [Dotan Nahum](http://gplus.to/dotan) [@jondot](http://twitter.com/jondot). See MIT-LICENSE for further details.



