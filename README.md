# dexbot
DeroDex Discord Bot

## Requires go v1.17

## To build / run:

### 1) Clone repository:
	$ git clone https://github.com/deroholic/dexbot.git

### 2) Build:
	$ cd dexbot
	$ go get dexbot
	$ go build

### 3) Configure:
```
{
        "token": "MTAxODE2MTg1NDkyNjgzNTg0NA.G-kSsq.i56LebqZWJdvy9Om-t0P5akP4wyWO77gOsgeVA",
        "prefix": "!",
        "owner": "pieswap#0888",
        "channel": "970710934425333820",
        "daemon": "dero-node-ca.mysrv.cloud:10102"
}
```
	Edit config.json and set "token" to your Discord bot's token & "owner" to your Discord handle.
	If you want to confine the bot to a single channel (other than DM) set "channel" to the ChannelID.

### 4) Run:
	$ ./dexbot
