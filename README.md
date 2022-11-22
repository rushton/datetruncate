# datetrunc
Datetrunc is a utility for truncating dates in an input stream.


# Usage
```
datetrunc <scale> [file]
```

# Example
```
#> echo "2020-01-12T12:00:00Z" | daterunc d
2020-01-12T00:00:00Z
```
