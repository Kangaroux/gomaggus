# Gomaggus

Passion project to build a 3.3.5 WoW server in Go.

If you have a better name than "gomaggus" I am open to suggestions :^)

## Migrations

SQL migrations are managed using [goose](https://github.com/pressly/goose). A linux CLI is included in `bin/`.

Apply migrations:

```bash
$ bin/goose -dir migration postgres 'postgres://gomaggus:password@localhost:5432/gomaggus?sslmode=disable' up
```

Create a new migration:

```bash
$ bin/goose -dir migration -s create MIGRATION_NAME sql
```

## Resources

- [WoW SRP6 implementation guide](https://gtker.com/implementation-guide-for-the-world-of-warcraft-flavor-of-srp6/) - very comprehensive guide that even includes test inputs
- [Shadowburn](https://gitlab.com/shadowburn/shadowburn) - basic realmd/worldd written in Elixir, has served as a good reference for networking
- [WoW dev wiki](https://wowdev.wiki) - good general resource, lots of info
- [WoW messages](https://gtker.com/wow_messages/) - networking docs that describe every packet payload, absolute lifesaver
- [Wireshark](https://www.wireshark.org/) - packet sniffer for debugging networking
- [How to read from a TCP socket (but were too afraid to ask)](https://incoherency.co.uk/blog/stories/reading-tcp-sockets.html) - great writeup on the gotchas of TCP
