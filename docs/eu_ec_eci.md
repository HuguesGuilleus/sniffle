---
title: Get ECI
---

Is a kind of petition (with 1 million signature milestone) addressed to European
Commission

- https://citizens-initiative.europa.eu/find-initiative_en
- https://citizens-initiative.europa.eu/initiatives/details/2022/000002_en

# Get the index

Get summary information for a range:

```txt
GET https://register.eci.ec.europa.eu/core/api/register/search/ALL/EN/:begin/:end
=> JSON response
```

And to get all items:

```txt
GET https://register.eci.ec.europa.eu/core/api/register/search/ALL/EN/0/0
=> JSON response
```

# Detail of one ECI

```txt
GET https://register.eci.ec.europa.eu/core/api/register/details/:year/:number
=> JSON response

ex:
https://register.eci.ec.europa.eu/core/api/register/details/2024/000001
```

# Icon

In the `logo` struct in index entry or in detail struct, take the field `id`

```txt
https://register.eci.ec.europa.eu/core/api/register/logo/:id
```
