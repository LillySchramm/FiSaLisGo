# FiSaLisGo

## Description

This is a Go wrapper for the [FiSaLis](https://www.finanz-sanktionsliste.de/fisalis/?) Website. It allows to search for companies and persons in the sanctions list.

## Installation

```bash
go get -u github.com/LillySchramm/FiSaLisGo
```

## Usage

```go
package main

import (
	"context"

	fisalisgo "github.com/LillySchramm/FiSaLisGo.git"
)

func main() {
	ctx := context.Background()

	res, err := fisalisgo.Search(ctx, "Vladimir Vladimirovich Putin")
	if err != nil {
		panic(err)
	}

    [...]
}
```

## Disclaimer

This is an unofficial wrapper and is not affiliated with the service in any way.