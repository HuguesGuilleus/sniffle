# sniffle

Scrap EU act then make a website.

```sh
go run -tags=devmode .
```

## Packages

- tool/ some packages to generated HTML, cache request... It's independant of
  others parts so you can use for other project.
- service/ manage specific EU part
- common/ some lib used by multiples services
- front/ JS, CSS and some component. Like common/ by for the frontend.
