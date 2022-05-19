# shortlinks

`shortlinks` is a mutable shortlink server.  The idea is that you can run it at
a very short domain and point shortlinks at documentation, dashboards, or whatever
for a given shortlink.

For example if I ran this (I currently do not) at `s.frew.co` and pointed `iam`
at
`https://docs.aws.amazon.com/service-authorization/latest/reference/reference_policies_actions-resources-contextkeys.html`,
`s.frew.co/iam` would redirect to that big URL.

Instead of having a rigid permission model, `shortlinks` is expected to run behind
a VPN.  On top of that, `shortlinks` stores historical versions of shortlinks so
if someone makes an undesired change, you can easily go back to a previous version
of that shortlink.

## Installation

```
$ go install github.com/frioux/shortlinks@main
```

## Usage

The default listen address is :8080 and the default database file is `db.db`.

```
$ shortlinks --listen :8081 --db file:shortlinks.db
```

Then navigate to `http://localhost:8081` and create your first shortlink!

## Custom Drivers

This tool is built to be easy to run using SQLite.  If you want to use some
other relational database or key value store more appropriate for your
infrastructure, all you need to do is create your own implementation of
[shortlinks.DB](https://pkg.go.dev/github.com/frioux/shortlinks/shortlinks#DB).
The only methods that need to actually work are `Shortlink` and
`CreateShortlink`, the rest can be left as stubs.

Similarly, there is a
[shortlinks.Auth](https://pkg.go.dev/github.com/frioux/shortlinks/shortlinks#Auth)
interface that you can use to only allow logged in users to make changes.  I have
an example driver using [tailscale](https://tailscale.com/) and am eager to hear if
the interface is sufficient for other auth methods.

The DynamoDB Storage driver (enabled by the `--dynamodb` flag) works, but is
currently unconfigurable.  Patches welcome to configure region, table name,
endpoint, etc.
