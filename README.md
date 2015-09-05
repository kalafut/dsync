(**note**: almost none of these things are implemented. The README is a function scoping document.)

# dsync

## Commands
```
add
update

```

## Commands

### add
Add a path to the catalog. Paths are assigned a short name that is used in dsync commands. This name will default to the last portion of the path and may be specified. Duplicate names are not allowed and will generate an error.

**add** only registers the path. It does not automatically invoke **update** to scan the path.

```
dsync add ~/Pictures pics     # add path, alias to "pics"
dsync add d:\Users\foo\docs   # alias to "docs"
dsync add /tmp/docs           # error: "docs" already exists
```

### update
Recursively scan a path and update the catalog. Paths must be first registered using **add**.

```
dsync update pics    # scan the path associated with "pics"
dsync update p*      # scan names starting with "p"
dsync update *       # scan all names
dsync update wrong   # error: "wrong" does not exist
```

### ls
List things, either files or roots. Without any parameters `ls` will list all roots in the catalog.

### select
Select a default catalog
