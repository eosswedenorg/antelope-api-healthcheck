# EOSIO API Healthcheck for HAProxy

This program implements EOSIO healthcheck for HAProxy over TCP.

## Compiling

You will need go-lang version `1.14` or later to compile the source.

First you need to install golang dependacies:

```sh
$ make deps
```

Then you can compile with `compile.sh` script

```sh
$ ./compile.sh
```

Execute `./compile.sh --help` to see all available flags to crosscompile for different systems/architectures.

### Package

run `./package.sh` to generate package. Debian (.deb) and FreeBSD are supported.

## TCP Protocol

The protocol is simple and has 4 rules.

1. Data is transmitted in `packets` encoded in ascii and ends with newline `\n`.
2. There are two types of packages: `Request` and `Response`. Each `Request` has exactly one `Response`.
3. Each parameter inside a `Request` is separated by `|`
4. Each response contains exactly one `status code` (see below)


### Request

The following parameters are supported in a request and are ordered from
first to last below:

| # | Name       | Required                | Description                                                    |
| - | ---------- | ----------------------- | -------------------------------------------------------------- |
| 1 | url        | Yes (port default `80`) | http url to the api. `http(s)://<ip-or-domain>(:<port>)`       |
| 2 | num_blocks | No (default `10`)       | Number of blocks the api can drift before reported `down`      |
| 3 | api        | No (default `v1`)       | Type of API to check against, `v1` = standard, `v2` = Hyperion |
| 4 | host       | No (default from `url`) | Value to send in the `HTTP Host Header` to the API             |

### Response

The api can respond with exactly one `status code`.
See [HAproxy documentation](https://cbonte.github.io/haproxy-dconv/1.7/configuration.html#5.2-agent-check) for more information

| code      | Description                                                |
| --------- | ---------------------------------------------------------- |
| `up`      | Api is healthy                                             |
| `down`    | Api is not healthy                                         |
| `failed`  | The program failed to read the status from the api.        |
| `maint`   | Api is set in maintenance mode (not used)                  |
| `ready`   | Api is ready again after being in `maint` state (not used) |
| `stopped` | Api has been stopped manually (not used)                   |

## Author

Henrik Hautakoski - [henrik@eossweden.org](mailto:henrik@eossweden.org)
