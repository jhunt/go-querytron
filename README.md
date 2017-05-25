Querytron
=========

Following hot on the successes of [Envirotron][env], **Querytron**
is here to save the day!

Ever needed to deal with a remote system that dealt in both
`application/json` _AND_ `application/x-www-form-urlencoded` data?

No?  I see you've never had the misfortune of integrating with
OAuth2 providers!  Good on you then.

For the rest of us, I wrote Querytron.  It works a lot like
Envirotron:

```
package thing

import (
  "fmt"
  qs "github.com/jhunt/go-querytron"
)

type Response struct {
  Error string `qs:"error"`
  URI   string `qs:"error_uri"`
}

func main() {
  url := SomeFunction()
  var r Response
  qs.Override(&r, url.Query())

  fmt.Printf("error %s (see also %s)\n", c.Error, c.URI)
}
```

Happy Hacking!

[env]: https://github.com/jhunt/go-envirotron
