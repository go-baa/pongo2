# pongo2
A helper for using [Pongo2](https://github.com/micate/pongo2) with Baa.

## getting started

```go
package main

import (
    "github.com/go-baa/pongo2"
    "gopkg.in/baa.v1"
)

func main() {
    // new app
    app := baa.New()

    // register pongo2 render
    // render is template DI for baa, must be this name.
    app.SetDI("render", pongo2.New(pongo2.Options{
        Baa:        b,
        Root:       "templates/",
        Extensions: []string{".html"},
        Functions:  map[string]interface{}{},
        Context: map[string]interface{}{
            "SITE_NAME": "Yet another website",
        },
    }))

    // router
    app.Get("/", func(c *baa.Context) {
        c.HTML(200, "index")
    })

    // run app
    app.Run(":1323")
}
```

## usage

### common

#### output

```html
This is {{ name }}.
```

#### include template

```html
{% include "path/to/tpl.html" %}
```

with params:

```html
{% include "relative/path/to/tpl.html" with foo=var %}
{% include "relative/path/to/tpl.html" with foo="bar" %}
```

**note**: nested template reveived param as string type. 

#### if / elif / else / endif

```html
{% if vara %}
{% elif varb %}
{% else %}
{% endif %}
```

#### for

```html
{% for item in items %}
{{ forloop.Counter }} {{ forloop.Counter0 }} {{ forloop.First }} {{ forloop.Last }} {{ forloop.Revcounter }} {{ forloop.Revcounter0 }}
{{ item }}
{% endfor %}
```

#### extends / block / macro and so on ...
see [document](https://docs.djangoproject.com/en/dev/ref/templates/language/).

### builtin filters

* escape
* safe
* escapejs
* add
* addslashes
* capfirst
* center
* cut
* date
* default
* default_if_none
* divisibleby
* first
* floatformat
* get_digit
* iriencode
* join
* last
* length
* length_is
* linebreaks
* linebreaksbr
* linenumbers
* ljust
* lower
* make_list
* phone2numeric
* pluralize
* random
* removetags
* rjust
* slice
* stringformat
* striptags
* time
* title
* truncatechars
* truncatechars_html
* truncatewords
* truncatewords_html
* upper
* urlencode
* urlize
* urlizetrunc
* wordcount
* wordwrap
* yesno
* float
* integer