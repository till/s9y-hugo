# s9y-hugo

A little exporter to move blog posts and categories from [Serendipity (s9y)](https://docs.s9y.org/) to [Hugo](https://gohugo.io/). Comes with lot of assumptions.

Used to learn some sql/sqlx (for Golang). Meant as a POC — published in case anyone finds it useful.

## Usage

`./s9y-hugo`

It will connect to MySQL (running on localhost/socket) and dump your entries into a structure of:

```sh
content
└── posts
    ├── cat-1
    └── cat-2
```

It will not:

- migrate/handle authors
- migrate/handle nested categories
- migrate/handle nested tags
- connect to a remote MySQL server
- configure Hugo (install theme, make changes to `hugo.yaml`, ...)

Config (via environment variables):

- `DB_USER`
- `DB_PASS`
- `DB_NAME`
- `DB_TABLE_PREFIX` (default: `s9y_`)
- `BLOG_URL` (default: `/blog`)

The `BLOG_URL` is used to add an alias to each post so redirects will continue to work.
