# chroma [![Build Status](https://travis-ci.org/alindeman/chroma.svg?branch=master)](https://travis-ci.org/alindeman/chroma)

> freedom from dilution with white and hence vivid in hue

**chroma** is a wrapper around the [Philips Hue
API](http://developers.meethue.com/index.html) in go.

It's not entirely complete, but probably doesn't suck. It has tests. I use it.

## Client

The first thing you'll want to know is how to construct a `chroma.Client`.
You'll need to know your bridge's IP address and a username that's either
[already authorized](http://developers.meethue.com/coreconcepts.html) or will
be authorized before use (an example of `Authorize()` is shown later).

```go
import (
  "github.com/alindeman/chroma"
)

func main() {
  client := &chroma.Client{
    BridgeHost: "192.168.1.2",
    Username:   "chroma",
  }
}
```

## Authorizing

Each new username must be authorized once. The link button on the bridge must
be physically pressed; afterward, `Authorize()` must be run within the next 30
seconds.

```go
if err := client.Authorize(); err != nil {
  fmt.Println("An error occurred: ", err)
}
```

The `Authorize()` step only needs to be performed once per username.

## Lights API

```go
lights, err := client.Lights()
if err != nil && len(lights) > 0 {
  if lights[0].State.On {
    fmt.Println("The first light is on!")
  }
}
```

Setting attributes is a bit awkward because I while I wanted to lean on go's
static typing, the Philips Hue API allows you to request only a subset of the
attributes be changed in a request. For those reasons, each member of the
attribute change struct is a pointer. If a member points to a value, it will
be sent. Otherwise, it will be left out.

For instance, to request that a light be turned on and its hue changed to
red, but other attributes (e.g., brightness and saturation) remain the same:

```go
change := &chroma.LightStateChange{
  On:  new(bool),
  Hue: new(int),
}
*change.Hue = 0 // red
*change.On = true

if err := client.SetLightState("1", change); err != nil {
  fmt.Println("Failed: ", err)
}
```

The Groups API is really similar to the Lights API.

## Full API Docs

<http://godoc.org/github.com/alindeman/chroma>

## Contribute

Pull requests welcomed.
