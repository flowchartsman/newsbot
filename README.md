# newsbot

A twitter streaming, website-scraping, websocket-transporting news delivery webapp written in Go and Javascript<sup>1</sup>

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

The bot is configured by the newsbot.conf file, which is written in [TOML](https://github.com/mojombo/toml/blob/master/versions/toml-v0.2.0.md).

A sample, seen below has been provided in **newsbot.conf.example**, which you'll need to edit and rename to **newsbot.conf** to get it to work. It should be pretty self explanatory if you [know what you're doing](https://dev.twitter.com/apps).  If not, I'll document it later, I promise. It's not that hard.

```
Port = ":8080"
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

[[scrapers]]
name = "Mysite"
# Located in /static
icon = "mysite.png"
# Pull it every 5 seconds
interval = 5
url = "http://www.my.great.site"
# We want any anchor of class "story" directly underneath things of class 'article'
target = ".article > a.story"
# We don't want any links whose text starts with "Catpics"
excluder = '^Catpics'
# This site always puts a link prefix for tracking. We don't want that.
modifier = 'http://links.my.great.site.com/\d+/'
```

## Todo

NewsBot "works" in the loosest sense only at the moment.  Links are delivered via webscraped stories, but not rendered. Twitter links are rendered, but not extracted. The "UI" is laughably bad because I'm still figuring out what's getting pushed up to the user.  So, yeah, I have a few things to work on, but I'm super serial now.

- Integrate RSS feeds
- Fix twitter link extraction, and delivery to the UI
- Make the UI do anything other than completely suck
    - Seriously. Anyone want to do this? I hate writing Javascript.
- Add chat so the users can talk to each other
- Let users see how many users are logged in (and their "names")

----
<sup>1</sup>: Or whatever. Someone pls halp me. [gopherjs](https://github.com/gopherjs/gopherjs) maybe?
