# newsbot

A twitter streaming, website-scraping<sup>1</sup> websocket-streaming news delivery webapp written in golang and JavaScript<sup>2</sup>

## Installation

These docs are shitty.  They are also new, and will be less shitty very soon.

```bash
$ git clone https://github.com/anaxagoras/newsbot.git
$ cd wherever/you/just/put/that
$ go build
```

## Running

```bash
$ ./newsbot
```

## Configuration

The bot is configured by the newsbot.conf file, which is written in [TOML](https://github.com/mojombo/toml/blob/master/versions/toml-v0.2.0.md),
a sample of which has been provided in newsbot.conf.example (you'll need to copy it to **newsbot.conf** and fill
in your own data to get it to work). It should be pretty self explanatory if you [know what you're doing](https://dev.twitter.com/apps).  If not, I'll document it later, I promise. It's not that hard.

```toml
LogLevel = "debug"
# Twitter auth settings
User = "yourusername"
ConsumerKey = "<consumer key here>"
ConsumerSecret = "<consumer secret here>"
OAuthToken = "<OAuth token here>"
OAuthSecret = "<OAuth secret here>"
# Twitter users to follow
Users = [
    1652541,   # @Reuters
    51241574,  # @AP
    18424289,  # @AJELive
    5402612,   # @BBCBreaking
    742143,    # @BBCWorld
    362051343, # @breakingstorm
    1068831    # @slashdot
]
# Keywords are temporarily unsupported in order to keep the bandwidth down and
# because the message culling algorithm currently ignores them.
#Keywords = [
#    "#YOLO",
#]
Port = ":8080"
```
----
<sup>1</sup>: Complete, but still merging. Watch out for it in the next couple days.

<sup>2</sup>: Soon to be dart
