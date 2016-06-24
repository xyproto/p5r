# P5R

This is a fork of [regexp2](https://github.com/dlclark/regexp2).

The motivation is to make Perl 5 compatible regular expressions available to [Otto](https://github.com/robertkrimen/otto), an implementation of JavaScript for Go.

Currently, [my fork of Otto](https://github.com/xyproto/otto) is not building with P5R, but it's pretty close.

Like `regexp2`, this package does not have constant time guarantees like the `regexp` package in the standard library, but it allows backtracking and is compatible with Perl 5 regular expressions.

`regexp2` was inspired by the regular expression implementation in .NET (which is released under an MIT license).

```go
re := p5r.MustCompile(`Your pattern`)
if isMatch := re.MatchString(`Something to match`); isMatch {
    //do something
}
```

For more documentation, take a look at the [regexp2](https://github.com/dlclark/regexp2) README.md.
