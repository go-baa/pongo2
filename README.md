# pongo2
A helper for using [Pongo2](https://godoc.org/github.com/flosch/pongo2) with Baa.

## init

### standalone

```go
pongo2.New(pongo2.Options{
    Root:       "templates/",
    Extensions: []string{".html"},
    Filters:    map[string]pongo2.FilterFunction{},
    Functions:  map[string]interface{}{},
    Context: map[string]interface{}{
        "SITE_NAME": "Yet another website",
    },
})
```

### with baa

```go
b.SetDI("render", pongo2.New(pongo2.Options{
    Baa:        b,
    Root:       "templates/",
    Extensions: []string{".html"},
    Filters:    map[string]pongo2.FilterFunction{},
    Functions:  map[string]interface{}{},
    Context: map[string]interface{}{
        "SITE_NAME": "Yet another website",
    },
}))
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

#### if / else

```html
{% if var %}
{% else %}
{% end %}
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