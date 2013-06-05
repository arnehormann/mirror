Mirror
============
Make unsafe code safer by seeing it in the mirror.

Mirror allows you to compare before you cast. Use it in `init` or with `sync.Once` to know you can safely cast a struct to another with `unsafe` because their memory layout is the same.

Access unexported fields without a headache. It may even be rather safe if you only read from them.

Documentation lives at [godoc.org](http://godoc.org/github.com/arnehormann/mirror).

License: [MPL2](https://github.com/arnehormann/mirror/blob/master/LICENSE.md).