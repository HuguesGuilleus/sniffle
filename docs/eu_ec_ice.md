---
title: Get ICE
---

Is a kind of petition addressed to European Commission. If more than one million
signtaures collected in one year, the Commission must responds.

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

# Detail of one ICE

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
