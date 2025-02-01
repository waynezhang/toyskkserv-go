## v0.0.5 (2025-02-01)

### New feature

- add reload command (`[d96bb6a](https://github.com/waynezhang/toyskkserv/commit/d96bb6ae9dc6ff72bdab056560129a116cd85e9a)`)
- off heap cache to surpress memory usage (`[8574794](https://github.com/waynezhang/toyskkserv/commit/85747944ecfd2803f88519d16e95b24fb72da4b5)`)
- use mapped memory as cache (`[25c7f7f](https://github.com/waynezhang/toyskkserv/commit/25c7f7f8bdb84f7d8eb167351cbcbea745ba8433)`)

### Fix

- crash when multiple reload command received (`[77d4c45](https://github.com/waynezhang/toyskkserv/commit/77d4c4576271baa0a5bebe8ea9b541ce4159cfb8)`)

### Refactor

- move cmd to root directory (`[f5363a8](https://github.com/waynezhang/toyskkserv/commit/f5363a884ccce44ff72b6a66dbec9fa2993554c8)`)



## v0.0.4 (2025-02-01)

### Fix

- [#4](https://github.com/waynezhang/toyskkserv/issues/4) connection lost (`[df2f94f](https://github.com/waynezhang/toyskkserv/commit/df2f94f204029f55575dc77fbb10da07b47f81c5)`)

### Refactor

- make handler more flexible (`[a0ae141](https://github.com/waynezhang/toyskkserv/commit/a0ae14103d1eb5e37eae298111eed81a61e6eddb)`)
- migrate to decoder (`[60db917](https://github.com/waynezhang/toyskkserv/commit/60db9171e7d5fd7b7d7b5330e23e3d98659e96ab)`)
- migrate to google btree (`[ef96e94](https://github.com/waynezhang/toyskkserv/commit/ef96e94f0221e27860397a8f090b865c8ea7a339)`)
- requeest handlers (`[e895f8f](https://github.com/waynezhang/toyskkserv/commit/e895f8fe326ea7f8a902d76508157073e1405f9d)`)
- write response to connection directly to surpress memory usage (`[4d998bd](https://github.com/waynezhang/toyskkserv/commit/4d998bd43c2fdfb4af929d67ea3de8fe5423631d)`)



## v0.0.3 (2025-02-01)

### Refactor

- make handler more flexible (`[a0ae141](https://github.com/waynezhang/toyskkserv/commit/a0ae14103d1eb5e37eae298111eed81a61e6eddb)`)
- migrate to decoder (`[60db917](https://github.com/waynezhang/toyskkserv/commit/60db9171e7d5fd7b7d7b5330e23e3d98659e96ab)`)
- requeest handlers (`[e895f8f](https://github.com/waynezhang/toyskkserv/commit/e895f8fe326ea7f8a902d76508157073e1405f9d)`)
- write response to connection directly to surpress memory usage (`[4d998bd](https://github.com/waynezhang/toyskkserv/commit/4d998bd43c2fdfb4af929d67ea3de8fe5423631d)`)



## v0.0.2 (2025-01-26)

### Bugs fixed:

- local dictionary cannot be loaded correctly([`ae3cc98`](https://github.com/waynezhang/tskks/commit/ae3cc9866d1a02620f96cdfc65990feb01556098))
- fix dictionary update([`9e7dd27`](https://github.com/waynezhang/tskks/commit/9e7dd27b845a7593da0d6d447cbb855e76293f35))


## v0.0.1 (2025-01-26)

### New feature:

- got rid of iconv, introduced a customized EUC decoder, with low([`36c7256`](https://github.com/waynezhang/tskks/commit/36c72566334619524f4ea9376f5266f98b1535be))
- support local file directionary([`0e6f968`](https://github.com/waynezhang/tskks/commit/0e6f968100bd0e691c341f40e04f40120ccb85cd))
- prioritize config file in current directory([`3c74f2e`](https://github.com/waynezhang/tskks/commit/3c74f2ea27a7d1e806e8db128a5bc25452871dd6))
- fix command handling for host, version([`b2e2da8`](https://github.com/waynezhang/tskks/commit/b2e2da8073e680b986d60c83b336a0a032f0cdd8))
- server completion([`937bacb`](https://github.com/waynezhang/tskks/commit/937bacb76bd0e74793469ea899a0b7bef2c7e0d5))
- limited google ime support([`73d38da`](https://github.com/waynezhang/tskks/commit/73d38daaab96f9671939965da8131d263daf5f88))
- be able to disable dict auto update([`def9715`](https://github.com/waynezhang/tskks/commit/def97151b4951cdaafd4115d79cc3e3864e17628))
- error detect on encoding converation([`94500c3`](https://github.com/waynezhang/tskks/commit/94500c395975b063c62ada243f9abb5abed2250e))
- libiconv binding([`8319da3`](https://github.com/waynezhang/tskks/commit/8319da347358bd5f5494b982a8159cecfd226c98))
- first commit([`9f95cfa`](https://github.com/waynezhang/tskks/commit/9f95cfac1b3190471461a40ea8b517e377593e36))

### Bugs fixed:

- iconv test for linux([`34b6e3e`](https://github.com/waynezhang/tskks/commit/34b6e3e934374d92dfe76b820ff936c97d2f97b5))
- expand tilde([`95b2a71`](https://github.com/waynezhang/tskks/commit/95b2a71fd01a9ef580969072df4a2a68997b6312))
- ci on linux([`340a00c`](https://github.com/waynezhang/tskks/commit/340a00cd1ebfee42334852844b132bb17f0fcb73))
- server crash and respond to 0 correctly([`536e320`](https://github.com/waynezhang/tskks/commit/536e3206d5ed3ff7d3145af9a9e15926c81697e4))
