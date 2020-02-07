# EOSIO API Healthcheck for HAProxy

This program implements EOSIO healthcheck for HAProxy over TCP.

## TCP Protocol

The protocol is simple and has 4 rules.

1. Data is transmitted in `packets` encoded in ascii and ends with newline `\n`.
2. There are two types of packages: `Request` and `Response`. Each `Request` has exactly one `Response`.
3. Each parameter inside a `Request` is separated by `:`
4. Each response contains exactly one `status code` (see below)


### Request

The following parameters are supported in a request and are ordered from
first to last below:

| # | Name       | Required | Description                                                                   |
| - | ---------- | -------- | ----------------------------------------------------------------------------- |
| 1 | Host       | Yes      | IP/Hostname to the api.                                                       |
| 2 | Port       | No       | Port number to the api (default `80`)                                         |
| 3 | num_blocks | No       | Number of blocks the api can drift before reported `down` (default 10)        |
| 4 | version    | No       | API Version to check against, `v1` = standard, `v2` = Hyperion (default `v1`) |

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
