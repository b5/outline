# outline

I spend a lot of time sketching out software packages. These sketches are often programming-language agnostic and written in lots of weird places (mainly my own notes, github issues, documentation, and source code). A while back @alandonovan posted what I felt is a clear, concise definition of a package: https://github.com/google/starlark-go/issues/19#issuecomment-336013683. Since reading that package outline I've started to think and write in something like that format.

While working on a standard library for starlark I ran into the issue of needing a way to embed documentation about a package that's written in go, but targets another language (starlark). I figure with a little rigor it'd be easiest to formalize my preferred outlining format in a way that it can be embedded in a comment, riding with the source code itself. Using the "template" commands we can generate documentation markdown for our website

### Project Status
use-at-your-own-risk alpha. I'm working on this with Qri's [rfc](https://github.com/qri-io/rfcs) and [starlib](https://github.com/qri-io/starlib) projects as concrete use-cases to drive development.

### Example
```shell
go install github.com/b5/outline
```

make a file: `outline.txt`:
```
  outline: geo
    geo defines geographic operations

    functions:
      point(lat,lng)
        Point constructor takes an x(longitude) and y(latitude) value and returns a Point object
        params
          lat float
          lng float
      within(geomA,geomB)
        Returns True if geometry A is entirely contained by geometry B
        params:
          geomA [point,line,polygon]
            maybe-inner geometry
          geomB [point,line,polygon]
            maybe-outer geometery
      intersects(geomA,geomB)
        Similar to within but part of geometry B can lie outside of geometry A and it will still return True

    types:
      point
        methods:
          buffer(x int)
            Generates a buffered region of x units around a point
          distance(p2 point)
            Euclidian Distance
          distanceGeodesic(p2 point)
            Distance on the surface of a sphere with the same radius as Earth
          KNN()
            Given a target point T and an array of other points, return the K nearest points to T
          greatCircle(p2 point)
            Returns the great circle line segment to point 2
      line
        methods:
          buffer()
          length()
          geodesicLength()
      polygon
```

Currently the only thing you can do out-of-the box with outline is parse documents & template them:
```
outline template ./outline.txt
```

This will by default spit out a markdown version of your outline. dope. The real upside of this format is it's designed to survive in weird places. try running the same command against this readme file:
```
git clone git@github.com:b5/outline.git
cd outline
outline template ./readme.md
```

And you'll get the same result. Lovely! You can supply custom templates with the `template` flag. The markdown template is [here](/cmd/template.go).


### Maybe someday...
* `outline fmt` <- "golint" style formatter
* `outline .` <- validate any found outline documents in a given filepath
* `outline require .` <- a command that requires at least one outline document present in the given filepath, useful for integration with CI
* `outline starter --language python .` <- generate starter stub code for a given package based on templates