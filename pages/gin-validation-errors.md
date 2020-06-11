title: "Gin Validation Errors Handling"
description: |
    Gin's validation system is powerful but it is complex to return meaningful
    JSON errors.
slug: gin-validation-errors-handling
banner: ""
draft: false
date: 2020-06-10 19:00:00
tags: [dev,go,gin,json]

# Introduction

> Gin is a web framework written in Go (Golang). It features a martini-like API
> with much better performance, up to 40 times faster thanks to 
> [httprouter](https://github.com/julienschmidt/httprouter). If 
> you need performance and good productivity, you will love Gin.
> 
> <cite>[Gin's Introduction](https://gin-gonic.com/docs/introduction/)</cite>

[Gin](https://gin-gonic.com/) is a very powerful web framework for Go. In my
opinion it has just the right balance between being really easy to use and the 
performance it provides. It is one if the most popular Go frameworks along with 
[gorilla/mux](https://github.com/gorilla/mux) and
[echo](https://github.com/labstack/echo).

In this post we'll see how Gin's validation works and how to return meaningful
errors to the clients calling your API. Because it's not as simple as it seems.

# Validator

## Struct Tags

As stated in its introduction, Gin is internally using 
[httprouter](https://github.com/julienschmidt/httprouter), but it's also using
[go-playground/validator](https://github.com/go-playground/validator) to 
validate incoming requests. Validator is using struct tag to determine what
to check and how to validate a struct field. It is not uncommon to find 
this kind of code when working with Gin:

```go
type DataRequest struct {
    Email string `json:"email" binding:"required"`
    Name  string `json:"name" binding:"required"`
}

func PostSomeData(c *gin.Context) {
	var q DataRequest

	if err := c.ShouldBind(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "field validation failed"})
		return
	}
	// Data is OK
}
```

In the above example, the struct `DataRequest` has two fields `Email` and `Name`
that are both required using the `binding:"required"` tags. They also have a 
`json` tag to determine the JSON key representation of this field, this tag is
used by `json.Marhsal(…)` and `json.Unmarshal(…)` for example. In short, to be
valid, an incoming request must look like this:

```json
{
    "email": "me@example.com",
    "name": "Me"
}
```

If one of those two fields are missing or empty (more on that later), then the
`c.ShouldBind` method will return an error. So what's happening under the hood?
First of all, Gin will unmarshal the request's body to the given `DataRequest`
variable, and that **can** fail for example, if the body is not even JSON. 
Then, if the unmarshal is successful, Gin will run its validator on the now 
filled struct. 

## Error Representation

In the previous example, if the unmarshalling or validation fails, we simply
responded with a 
[400 Bad Request](https://developer.mozilla.org/fr/docs/Web/HTTP/Status/400),
with no additional information. Let's see what kind of error is returned if 
the validation fails. For that we'll log the error using 
[zerolog](https://github.com/rs/zerolog) as our logging library:

```go
func PostSomeData(c *gin.Context) {
	var q DataRequest

	if err := c.ShouldBind(&q); err != nil {
		log.Info().Err(err).Msg("field validation failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": "field validation failed"})
		return
	}
	// Data is OK
}
```

```
5:23PM INF field validation failed error="Key: 'DataRequest.Email' Error:Field validation for 'Email' failed on the 'required' tag\nKey: 'DataRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag"
```

Right. Here we have the string representation of our error, it describes
everything wrong with our struct's values. We could send that back to our
caller, right? Well not really, there are several issues with that format:

- In that error we have <ins>struct field names</ins>, not the actual JSON tag 
  associated with the field. Meaning, if our JSON struct tag doesn't match with 
  the field name, the error is completely opaque for our users.
- It gives out too many details. Do our end users need to know that the internal 
  struct we unmarshal to is called `DataRequest`? Also, if they're unfamiliar
  with go, the `tag` thing is also opaque.
- Errors are separated with `\n`, that's great for text based responses or even
  logging, but it won't work well with JSON. Imagine if someone was developping 
  a frontend for your API, they would probably have a hard time parsing that and
  making the error meaningful for the end user.

## Validation Error

What we saw earlier was just the string representation of our error. That makes
sense, because in Go errors must implement the `error` interface defined like
so:

```go
type error interface {
    Error() string
}
```

So when we log it using our logging library, it will show us the string
representation of that error, the one that we can read and make sense of. But
that's not the actual error. We know that Gin's validator will return 
[validator.ValidationErrors](https://pkg.go.dev/github.com/go-playground/validator/v10?tab=doc#ValidationErrors)
if a validation error occurs, basically it will just send back the error it
encountered. And what we can do now is called 
[type assertion](https://tour.golang.org/methods/15):

```go
func PostSomeData(c *gin.Context) {
	var q DataRequest

	if err := c.ShouldBind(&q); err != nil {

		if verr, ok := err.(validator.ValidationErrors); ok {
            log.Info().Err(verr).Msg("this is actually a validation error")
        }
        
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	// Data is OK
}
```

> A type assertion provides access to an interface value's underlying concrete 
> value.
>
> <cite>[Tour of Go](https://tour.golang.org/methods/15)</cite>

Basically what this mean is that we know `err` is an error (that implements the
`error` interface) but we also suspect it's a `validator.ValidationErrors`, so 
we try to assert the error type to access the underlying 
`validator.ValidationErrors` that implements `error`. If the type assertion
fails, then it's not a validation issue, and we can define another behavior, 
for example if the request body is an invalid JSON.

Using type assertion, we can now determine the error returned by the 
`ShouldBind` method. And as of Go 1.13 we now have a more "friendly" way of
doing [error type assertion](https://blog.golang.org/go1.13-errors):

```go
func PostSomeData(c *gin.Context) {
	var q DataRequest

	if err := c.ShouldBind(&q); err != nil {

		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			log.Info().Err(verr).Msg("this is actually a validation error")
        }
        
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	// Data is OK
}
```

And now using `verr` we have full access to this error and its associated
methods!

## Making sense of ValidationErrors

We can see that `validator.ValidationErrors` is actually a slice of 
`validator.FieldError` and that `FieldError` is actually an interface with
many methods used to manipulate and retrieve information about the error:

```go
// ValidationErrors is an array of FieldError's
// for use in custom error messages post validation.
type ValidationErrors []FieldError
```

<details class="code"><summary>FieldError interface</summary>

```go
// FieldError contains all functions to get error details
type FieldError interface {

	// returns the validation tag that failed. if the
	// validation was an alias, this will return the
	// alias name and not the underlying tag that failed.
	//
	// eg. alias "iscolor": "hexcolor|rgb|rgba|hsl|hsla"
	// will return "iscolor"
	Tag() string

	// returns the validation tag that failed, even if an
	// alias the actual tag within the alias will be returned.
	// If an 'or' validation fails the entire or will be returned.
	//
	// eg. alias "iscolor": "hexcolor|rgb|rgba|hsl|hsla"
	// will return "hexcolor|rgb|rgba|hsl|hsla"
	ActualTag() string

	// returns the namespace for the field error, with the tag
	// name taking precedence over the fields actual name.
	//
	// eg. JSON name "User.fname"
	//
	// See StructNamespace() for a version that returns actual names.
	//
	// NOTE: this field can be blank when validating a single primitive field
	// using validate.Field(...) as there is no way to extract it's name
	Namespace() string

	// returns the namespace for the field error, with the fields
	// actual name.
	//
	// eq. "User.FirstName" see Namespace for comparison
	//
	// NOTE: this field can be blank when validating a single primitive field
	// using validate.Field(...) as there is no way to extract it's name
	StructNamespace() string

	// returns the fields name with the tag name taking precedence over the
	// fields actual name.
	//
	// eq. JSON name "fname"
	// see StructField for comparison
	Field() string

	// returns the fields actual name from the struct, when able to determine.
	//
	// eq.  "FirstName"
	// see Field for comparison
	StructField() string

	// returns the actual fields value in case needed for creating the error
	// message
	Value() interface{}

	// returns the param value, in string form for comparison; this will also
	// help with generating an error message
	Param() string

	// Kind returns the Field's reflect Kind
	//
	// eg. time.Time's kind is a struct
	Kind() reflect.Kind

	// Type returns the Field's reflect Type
	//
	// // eg. time.Time's type is time.Time
	Type() reflect.Type

	// returns the FieldError's translated error
	// from the provided 'ut.Translator' and registered 'TranslationFunc'
	//
	// NOTE: if no registered translator can be found it returns the same as
	// calling fe.Error()
	Translate(ut ut.Translator) string
}
```

</details>

And that makes sense because we could be dealing with multiple errors since
we're validating multiple fields. So let's iterate over this slice and see
what information we can get:

```go
var verr validator.ValidationErrors
if errors.As(err, &verr) {
    for _, f := range verr {
        log.Info().Str("name", f.Field()).Str("tag", f.Tag()).Msg("field error")
    }
}
```

```
6:20PM INF field error name=Email tag=required
6:20PM INF field error name=Name tag=required
```

Now we could do something with that! Except the `name` is still the struct field
name and not the associated JSON, but we'll get to that later. So `f.Field()`
returns the field that failed and `f.Tag()` returns the tag that triggered the
failure. Now that's useful, as our tags usually make sense, like `required`,
`max`, `min`, etc. Now this snippet is getting quite big so let's handle that
in another function that will take a `validator.ValidationErrors` as its input:

```go
func Simple(verr validator.ValidationErrors) map[string]string {
	errs := make(map[string]string)

	for _, f := range verr {
		err := f.ActualTag()
		if f.Param() != "" {
			err = fmt.Sprintf("%s=%s", err, f.Param())
		}
		errs[f.Field()] = err
	}

	return errs
}
```

Here we're returning a `map[string]string` because that's a type Go's JSON
library can deal with pretty easily. Here we're ranging over the field errors
and for each of them we're setting the field name as the map key and the tag
that matched as the value. Some tags, unlike `required`, are using parameters.
For example `max` takes an integer that will determine which is the maximum
length of the provided string or the maximal value for an integer field.

Let's use that in our handler:

```go
func PostSomeData(c *gin.Context) {
	var q DataRequest

	if err := c.ShouldBind(&q); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			c.JSON(http.StatusBadRequest, gin.H{"errors": Simple(verr)})
			return
		}

        // We now know that this error is not a validation error
        // probably a malformed JSON
		log.Info().Err(err).Msg("unable to bind")
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	// Data is OK
}
```

```json
{
    "errors": {
        "Email": "required",
        "Name": "required"
    }
}
```

Now we need to handle that field name instead of tag situation. And that will
happen at Gin's level.

# Gin's Validator Customization

Gin is using
[go-playground/validator](https://pkg.go.dev/github.com/go-playground/validator/v10) 
internally, but it's possible to access Gin's validator instance to customize 
it. In fact this is even shown in [the documentation](https://github.com/gin-gonic/gin#custom-validators):

```go
func main() {
	route := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
        // ...
	}
    // ...
	route.Run(":8085")
}
```

We don't want to add custom validators though, we want to access the JSON tag
instead of our struct field. Well, lucky for us, [`RegisterTagNameFunc`](https://pkg.go.dev/github.com/go-playground/validator/v10?tab=doc#Validate.RegisterTagNameFunc) exists in the `validator`
lib! So let's register our tag name func:

```go
if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
    v.RegisterTagNameFunc(func(fld reflect.StructField) string {
        name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
        if name == "-" {
            return ""
        }
        return name
    })
}
```

`RegisterTagNameFunc` expects a `func(fld reflect.StructField) string` function.
Here we're basically telling Gin's validator instance that the `f.Field()`
method we used earlier should not return the struct field name, but the
associated JSON tag (omitting everything after the coma if there is one).

And just like that, our API now returns:

```json
{
    "errors": {
        "email": "required",
        "name": "required"
    }
}
```

# Going further

## No dynamic JSON keys

So we've seen how we can customize the JSON response in case there is a
validation error. Now depending on who is going to call our API, we might not
want dynamic keys in our JSON output. For example in JavaScript it's fairly
easy to handle, due to the dynamic nature of this language. 
But if another Go program calls it will make its life difficult. So let's fix
that:

```go
type ValidationError struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

func Descriptive(verr validator.ValidationErrors) []ValidationError {
	errs := []ValidationError{}

	for _, f := range verr {
		err := f.ActualTag()
		if f.Param() != "" {
			err = fmt.Sprintf("%s=%s", err, f.Param())
		}
		errs = append(errs, ValidationError{Field: f.Field(), Reason: err})
	}

	return errs
}
```

And let's simply replace our earlier `Simple(verr)` to `Descriptive(verr)`:

```go
if errors.As(err, &verr) {
    c.JSON(http.StatusBadRequest, gin.H{"errors": Descriptive(verr)})
    return
}
```

<details class="code"><summary>JSON Response</summary>

```json
{
    "errors": [
        {
            "field": "email",
            "reason": "required"
        },
        {
            "field": "name",
            "reason": "required"
        }
    ]
}
```

</details>

## Package

We could even put all that in a package and customize the validator instance
during creation of our formatter. This will tend to get closer to a clean design
architecture:

<details class="code"><summary>formatter.go</summary>

```go
package formatter

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type JSONFormatter struct{}

// NewJSONFormatter will create a new JSON formatter and register a custom tag
// name func to gin's validator
func NewJSONFormatter() *JSONFormatter {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}

	return &JSONFormatter{}
}

type ValidationError struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

func (JSONFormatter) Descriptive(verr validator.ValidationErrors) []ValidationError {
	errs := []ValidationError{}

	for _, f := range verr {
		err := f.ActualTag()
		if f.Param() != "" {
			err = fmt.Sprintf("%s=%s", err, f.Param())
		}
		errs = append(errs, ValidationError{Field: f.Field(), Reason: err})
	}

	return errs
}

func (JSONFormatter) Simple(verr validator.ValidationErrors) map[string]string {
	errs := make(map[string]string)

	for _, f := range verr {
		err := f.ActualTag()
		if f.Param() != "" {
			err = fmt.Sprintf("%s=%s", err, f.Param())
		}
		errs[f.Field()] = err
	}

	return errs
}
```
</details>
