[![][travis-badge]][travis-link]
[![][coverage-badge]][coverage-link]
![][license-badge]
![][go-version]

# Measures

A Go package to send application metrics using UDP.


## Usage

Create an instance of `Measures` specifying the client name and server address:

```go
Measures = measures.New("measures-example", "0.0.0.0:3593")
```

Count something – a fruit, for example:

```
Measures.Count("fruits", 1, measures.Dimensions{
    "name": "avocado",
    "color": "green",
})
```

You can also measure time spent on something:

```go
func someLengthyOperation() {
    defer M.Time("example", time.Now(), nil) // No Dimensions in this case
    time.Sleep(359e6)
}

```

And you can provide Dimensions at once:

```go
func anotherLengthyOperation() {
    defer M.Time("example", time.Now(), measures.Dimensions{"number": 359})
    time.Sleep(359e6)
}
```

Or as you have them:

```go
func yetanotherLengthyOperation() {
    d := make(measures.Dimensions, 2)
    defer M.Time("example", time.Now(), d)
    d["number"] = 359
    time.Sleep(359e6)
    d["isPrime"] = true
}
```


# License

[Simplified BSD][bsd-2cl] © [Measures authors][authors] et [al][contributors]


[authors]:         https://github.com/scorphus/measures/blob/master/AUTHORS
[bsd-2cl]:         https://opensource.org/licenses/BSD-2-Clause
[contributors]:    https://github.com/scorphus/measures/blob/master/CONTRIBUTORS
[go-version]:      https://img.shields.io/badge/Go->=1.5.1-6DD2F0.svg
[coverage-badge]:  https://img.shields.io/coveralls/scorphus/measures.svg
[coverage-link]:   https://coveralls.io/github/scorphus/measures
[license-badge]:   https://img.shields.io/badge/license-BSD_2‑Clause-007EC7.svg
[travis-badge]:    http://img.shields.io/travis/scorphus/measures.svg
[travis-link]:     https://travis-ci.org/scorphus/measures
