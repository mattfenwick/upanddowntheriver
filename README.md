# Up And Down The River (Merde!)

A simple, silly, and fun card game!



## How to run locally

Get set up with [golang on your platform](https://golang.org/dl/)

Clone the repo:
```
git clone git@github.com:mattfenwick/upanddowntheriver.git
```

Run the server:

```
cd upanddowntheriver/cmd/server

go run server.go conf.json
```

[Visit the UI](http://localhost:5932/main.html)


## Build and push an image

```
make server

docker push <whatever-tag-was-built>
```


## Hacking on the server

```
make fmt

make vet

make test
```


## How to deploy to kubernetes

Prereqs: `kubectl` set up and talking to a kubernetes cluster

Run the deploy script:

```
cd upanddowntheriver/deploy

NAMESPACE=my-fave-namespace
IMAGE_TAG=v0.3.0
./deploy.sh $NAMESPACE $IMAGE_TAG
```

Note: this exposes an `up-and-down-the-river` service, you may want to expose it differently depending on your setup!



## Components

Server: [golang](./cmd/server/server.go)

UI: [javascript](./cmd/server/ui)


# TODOs

A modern, idiomatic UI that isn't ugly :)


# Rules

This uses a standard, 52-card deck.

At the beginning of each round, each player is dealt the same number of cards
and a trump suit is chosen.

Starting with the player to the left of the dealer and continuing to the left,
each player wagers the number of tricks they expect to win.  Players may wager
any amount between 0 and the number of cards dealt to each player, *except*
that the dealer's wager may not cause the sum of all wagers to equal the number
of cards.  This ensures that *at least* one player will not be able to hit
their wager.

Once all wagers have been made, the player to the left of the dealer starts the first
trick.  The first card played for each trick determines the suit of that trick.
All players *must* follow suit if they are able; however, if a player is unable to
follow suit, they may play any card of any suit.  For example, if the first player
plays a Diamond, then any player with Diamonds in their hand must also play a Diamond
for that trick; a player with no Diamonds may play Clubs, Hearts, or Spades.

The winner of a trick is the player who played the highest card (Aces high) of the
trick's suit.  **However**, any card of the chosen trump suit beats any card of
another suit.  If two or more trump cards are played, the highest wins.

A player wins a round if the number of tricks taken equals their wager.  If a player
wins more or fewer tricks than wagered, they lose that round.  Thus, multiple players
can win in each round, but at least one player will lose (and possibly all :) ).

If you have a lot of time on your hands, you can go up and down the river.  This means
you play a round with 1 card to each player, then 2, then 3, etc. until you run out of
cards.  For example, with 4 players you could go all the way up to 13 cards apiece.
Then you go back down the river, playing rounds with 13, 12, 11, all the way back down
to 1.  This takes a *lot* of time :) 