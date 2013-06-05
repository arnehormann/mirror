Plan for CanConvertUnsafe
==========
### Add tests
Lots of tests. Lots and lots of tests... You get it...
- [ ] each possible type, first with single-type structs
- [ ] non exported version in `test/from` and `test/to`
- [ ] interfaces in struct
- [ ] recursion

### Fine tune behavior
- [ ] Don't check type names top level; that can also be done by caller or by wrapping it into an anonymous struct.
- [ ] Make struct length comparison optional with `ignoresize` arg in another func; breaks for slices of the type!
- [ ] Skip field name check and type check for `_` by default (padding)

## Make behavior per field configurable by field tags
- [ ] add field tag handling on `to` (**`mirror-check:...,...`**)
- [ ] don't check at all (dangerous but available with tag `ignore`)
- [ ] match field name (default is **do not**; tag `fieldname` / `nofieldname`)
- [ ] match type name (default is **do** unless name is `_`; tag name `typename` / `notypename`)
- [ ] match type (default is **do** unless name is `_`; tag name `type` / `notype`)

# Structs
- [ ] recurse into structs (default is **do not**; tag name `follow` / `nofollow`)

#Interfaces
- [ ] by type identity (tag name `same`)
- [ ] by matching method set (tag name `match`)
- [ ] by assignability - `from` has to contain `to` (tag name `assignable`)

### Other
- [ ] extract logging into another function (common "walk struct" function?)

## Add subpackages?
- [ ] JSON-export of type layout by HTTP with blocking channel to feed it, requires GET by client to unblock
