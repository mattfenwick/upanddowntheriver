'use strict';

$(document).ready(function() {

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

// network actions

function getModel(cont) {
    function f(data, status, _jqXHR) {
        console.log("GET model response -- status " + status);
        cont(status === 'success', data);
    }
    $.ajax({
        'url': '/model',
        'dataType': 'json',
        'success': f,
        'error': f,
    });
    console.log("fired off GET to /model");
}

function postAction(payload, cont) {
    function f(data, status, _jqXHR) {
        console.log("post /action response -- status " + status);
        cont(status === 'success', data);
    }
    $.post({
        'url': '/action',
        'data': JSON.stringify(payload),
        'dataType': 'json',
        'success': f,
        'error': f,
        'contentType': 'application/json'
    });
    console.log("fired off POST to /action");
}

function getMyModel(me, cont) {
    postAction({'Me': me, 'GetModel': {}}, cont);
}

function postJoin(me, cont) {
    postAction({'Me': me, 'Join': {}}, cont)
}

function postRemovePlayer(me, name, cont) {
    postAction({'Me': me, 'RemovePlayer': {'Player': name}}, cont)
}

function postSetCardsPerPlayer(me, cardCount, cont) {
    postAction({'Me': me, 'SetCardsPerPlayer': {'Count': cardCount}}, cont)
}

function postStartRound(me, cont) {
    postAction({'Me': me, 'StartRound': {}}, cont)
}

function postWager(me, hands, cont) {
    postAction({'Me': me, 'MakeWager': {'Hands': hands}}, cont)
}

function postStartHand(me, cont) {
    postAction({'Me': me, 'StartHand': {}}, cont);
}

function postPlayCard(me, card, cont) {
    postAction({'Me': me, 'PlayCard': card}, cont)
}

function postFinishRound(me, cont) {
    postAction({'Me': me, 'FinishRound': {}}, cont);
}

// util

function equals(a, b) {
    if ( Array.isArray(a) && Array.isArray(b) ) {
        return arrayEquals(a, b);
    }
    if ( a === null ) { return b === null; }
    if ( b === null ) { return false; }
    if ( (typeof a === "object") && (typeof b === "object") ) {
        return objectEquals(a, b);
    }
    return a === b;
}

function arrayEquals(a, b) {
    if ( a.length !== b.length ) {
        return false;
    }

    for (let i = 0; i < a.length; ++i) {
        if ( !equals(a[i], b[i]) ) {
        // if ( a[i] !== b[i] ) {
            return false;
        }
    }
    return true;
}

function objectEquals(a, b) {
    if ( !arrayEquals(Object.keys(a), Object.keys(b)) ) {
        return false;
    }
    for (let key in a) {
        if ( !equals(a[key], b[key]) ) {
            return false;
        }
    }
    return true;
}

// me

function Me(didClickJoin) {
    $("#me-join").click(function () {
        console.log("me-join click");
        let name = $("#me-input-name").val();
        didClickJoin(name);
    });
    this.name = "";
}

Me.prototype.update = function(name) {
    if ( name === this.name ) { return; }
    this.name = name;
    $("#me-show-name").append(name);
    $("#me-get-name").hide();
};

// Game

function Game(didClickRemovePlayer, didChangeCardsPerPlayer, didClickStart) {
    this.players = [];
    this.cardsPerPlayer = null;

    this.startButton = $("#game-round-start");
    this.startButton.click(didClickStart);

    // can't use `click` because the elements might be ever-changing
    $(document).on("click", ".game-remove-player", function() {
        // TODO apparently we're not supposed to use the `player` attribute here
        let player = $(this).attr('player');
        console.log("game-remove-player: ${player}");
        didClickRemovePlayer(player);
    });

    let self = this;
    this.cardsPerPlayerSelect = $("#game-cards-per-player");
    this.cardsPerPlayerSelect.change(function() {
        let newCardsPerPlayer = parseInt(self.cardsPerPlayerSelect.val(), 10);
        self.cardsPerPlayer = newCardsPerPlayer;
        didChangeCardsPerPlayer(newCardsPerPlayer);
    });

    this.setState("WaitingForPlayers");
}

Game.prototype.update = function(me, state, players, cardsPerPlayer) {
    this.setPlayers(players);
    this.setCardsPerPlayer(cardsPerPlayer);
    this.setState(state);
};

Game.prototype.setState = function(state) {
    if ( this.state === state ) { return; }
    console.log(`updating game state to ${state}`);
    switch (state) {
        case "NotJoined":
            $(".game-remove-player").hide();
            this.startButton.prop('disabled', true);
            this.cardsPerPlayerSelect.prop('disabled', true);
            break;
        case "WaitingForPlayers":
            $(".game-remove-player").show();
            this.startButton.prop('disabled', this.players.length < 2);
            this.cardsPerPlayerSelect.prop('disabled', false);
            break;
        case "RoundWagerTurn":
            $(".game-remove-player").hide();
            this.startButton.prop('disabled', true);
            this.cardsPerPlayerSelect.prop('disabled', true);
            break;
        case "RoundHandReady":
            break;
        case "HandPlayTurn":
            break;
        case "RoundFinished":
            break;
        default:
            throw new Error("unrecognized game state: " + state)
    }
};

Game.prototype.setPlayers = function(players) {
    if ( arrayEquals(this.players, players) ) { return; }
    //
    this.players = players;
    let domPlayers = $("#game-players");
    domPlayers.empty();
    players.forEach(function(player) {
        domPlayers.append(`
        <tr>
            <td>${player}</td>
            <td>
                <button class='game-remove-player' player='${player}'>Remove</button>
            </td>
        </tr>`);
    });

    this.startButton.prop('disabled', players.length < 2);

    this.setCardsPerPlayerOptions();
};

Game.prototype.setCardsPerPlayerOptions = function() {
    let domElem = this.cardsPerPlayerSelect;
    if ( this.players.length === 0 ) {
        domElem.prop('disabled', true);
        return;
    }
    domElem.prop('disabled', false);
    let maxCards = 52 / this.players.length;
    domElem.empty();
    for (let i = 1; i <= maxCards; i++) {
        domElem.append(`<option value='${i}'>${i}</option>`);
    }
};

Game.prototype.setCardsPerPlayer = function(cardsPerPlayer) {
    // if ( this.cardsPerPlayer === cardsPerPlayer ) { return; }
    this.cardsPerPlayerSelect.val(cardsPerPlayer);
    this.cardsPerPlayer = cardsPerPlayer;
};

// round

function Round(didChooseWager, didClickStartHand, didClickFinishRound) {
    // this.roundDiv = $("#round");
    this.finishButton = $("#round-finish");
    this.cardsTable = $("#round-cards");
    this.wagersTable = $("#round-wagers");
    this.wagersTableBody = $("#round-wagers tbody");
    this.wagerSelect = $("#round-wager-select");
    this.wagerButton = $("#round-wager-button");
    this.trumpSuitDiv = $("#round-suit");
    this.trumpContainer = $("#round-suit-container");
    this.wagerContainer = $("#round-wager");
    this.startHandButton = $("#round-start-hand");

    let self = this;
    this.wagerButton.click(function() {
        let wager = parseInt(self.wagerSelect.val(), 10);
        didChooseWager(wager);
    });
    this.finishButton.click(didClickFinishRound);
    this.startHandButton.click(didClickStartHand);

    this.me = "";
    this.nextWagerPlayer = "";
    this.trumpSuit = "";
    this.state = "";

    this.setRoundState("WaitingForPlayers");
    this.setCards([]);
    this.setWagers([], []);
}

Round.prototype.update = function(me, state, trumpSuit, cards, wagers, nextWagerPlayer, wagerSum) {
    // TODO do something wagerSum ?
    this.me = me;
    if ( cards ) {
        this.setCards(cards);
    }
    if ( wagers ) {
        this.setWagers(wagers);
    }
    if ( trumpSuit ) {
        this.setTrumpSuit(trumpSuit);
    }
    this.setRoundState(state);
    if ( nextWagerPlayer ) {
        this.setNextWagerPlayer(nextWagerPlayer);
    }
};

Round.prototype.setTrumpSuit = function(trumpSuit) {
    if ( this.trumpSuit === trumpSuit ) { return; }
    this.trumpSuit = trumpSuit;
    this.trumpSuitDiv.empty();
    this.trumpSuitDiv.append(trumpSuit);
};

Round.prototype.setNextWagerPlayer = function(player) {
    if ( player === this.nextWagerPlayer ) { return; }
    this.nextWagerPlayer = player;
    if ( player === this.me && this.state === "RoundWagerTurn" ) {
        this.wagerSelect.empty();
        for ( let i = 1; i <= this.cards.length; i++ ) {
            this.wagerSelect.append(`<option value="${i}">${i}</option>`);
        }
        this.wagerContainer.show();
    } else {
        this.wagerContainer.hide();
    }
};

Round.prototype.setRoundState = function(state) {
    if ( state === this.state ) { return; }
    this.state = state;
    console.log(`setting round state to ${state}`);
    switch (state) {
        case "NotJoined":
            this.cardsTable.hide();
            this.finishButton.hide();
            this.startHandButton.hide();
            this.trumpContainer.hide();
            this.wagerContainer.hide();
            this.wagersTable.hide();
            break;
        case "WaitingForPlayers":
            this.cardsTable.hide();
            this.finishButton.hide();
            this.startHandButton.hide();
            this.trumpContainer.hide();
            this.wagerContainer.hide();
            this.wagersTable.hide();
            break;
        case "RoundWagerTurn":
            this.cardsTable.show();
            this.finishButton.hide();
            this.startHandButton.hide();
            this.trumpContainer.show();
            this.wagerContainer.hide();
            this.wagersTable.show();
            break;
        case "RoundHandReady":
            this.cardsTable.show();
            this.finishButton.hide();
            this.startHandButton.show();
            this.trumpContainer.show();
            this.wagerContainer.hide();
            this.wagersTable.show();
            break;
        case "HandPlayTurn":
            this.cardsTable.show();
            this.finishButton.hide();
            this.startHandButton.hide();
            this.trumpContainer.show();
            this.wagerContainer.hide();
            this.wagersTable.show();
            break;
        case "RoundFinished":
            this.finishButton.show();
            break;
        default:
            throw new Error(`invalid round state ${state}`);
    }
};

Round.prototype.setCards = function(cards) {
    if ( equals(cards, this.cards) ) {
        return;
    }
    this.cards = cards;
    this.cardsTable.empty();
    let cardsBySuit = {};
    cards.forEach(function(card) {
        if ( !(card.Suit in cardsBySuit) ) {
            cardsBySuit[card.Suit] = [];
        }
        cardsBySuit[card.Suit].push(card.Number);
    });
    let suits = Object.keys(cardsBySuit);
    suits.sort();
    let self = this;
    suits.forEach(function(suit) {
        let suitTds = [];
        cardsBySuit[suit].forEach(function(number) {
            suitTds.push(`<td>${number}</td>`);
        });
        let tr = `<tr><td>${suit}</td>${suitTds.join("\n")}</tr>`;
        self.cardsTable.append(tr);
    });
};

Round.prototype.setWagers = function(wagers) {
    if ( equals(this.wagers, wagers) ) {
        return;
    }
    this.wagers = wagers;
    let self = this;
    self.wagersTableBody.empty();
    wagers.forEach(function(wager) {
        let player = wager.Player;
        let wagerCount = (wager.Count !== null) ? wager.Count : "";
        let wonCount = (wager.HandsWon !== null) ? wager.HandsWon : "";
        let style = (player === self.me) ? 'style="border: 1px dashed; padding: 8px; margin: 4px;"' : '';
        self.wagersTableBody.append(`
            <tr ${style}>
                <td>
                ${player}
                </td>
                <td>
                ${wagerCount}
                </td>
                <td>
                ${wonCount}
                </td>
            </tr>
        `)
    });
};

// hand

function Hand(didClickPlayCard) {
    this.me = "";
    this.state = "";
    this.suit = "";
    this.leader = "";
    this.leaderCard = null;
    this.cardsPlayed = [];
    this.nextPlayer = "";
    this.myCards = [];

    this.suitDiv = $("#hand-suit");
    this.leaderContainer = $("#hand-leader-container");
    this.leaderPlayerDiv = $("#hand-leader-player");
    this.leaderCardDiv = $("#hand-leader-card");
    this.cardsPlayedTable = $("#hand-cards");
    this.cardsPlayedTableBody = $("#hand-cards tbody");
    this.playContainer = $("#hand-play-container");
    this.playSelect = $("#hand-play-select");
    this.playButton = $("#hand-play-button");

    this.setState("NotJoined");

    let self = this;
    this.playButton.click(function() {
        let i = parseInt(self.playSelect.val(), 10);
        didClickPlayCard(self.myCards[i]);
    });

}

Hand.prototype.update = function(me, state, suit, leader, leaderCard, cardsPlayed, nextPlayer, myCards) {
    this.me = me;
    this.setState(state);
    // if ( suit ) {
        this.setSuit(suit);
    // }
    // if ( leader ) {
        this.setLeader(leader);
    // }
    // if ( leaderCard ) {
        this.setLeaderCard(leaderCard);
    // }
    // if ( cardsPlayed ) {
        this.setCardsPlayed(cardsPlayed);
    // }
    // if ( myCards ) {
        this.setMyCards(myCards);
    // }
    // if ( nextPlayer ) {
        this.setNextPlayer(nextPlayer);
    // }
};

Hand.prototype.setMyCards = function(myCards) {
    if ( equals(myCards, this.myCards) ) { return; }
    this.myCards = myCards;
};

Hand.prototype.setNextPlayer = function(player) {
    if ( this.nextPlayer === player ) { return; }
    this.nextPlayer = player;

    if ( player === this.me && this.state === "HandPlayTurn" ) {
        this.playSelect.empty();
        for ( let i = 0; i < this.myCards.length; i++ ) {
            let card = this.myCards[i];
            let desc = `${card.Number} of ${card.Suit}`;
            this.playSelect.append(`<option value="${i}">${desc}</option>`);
        }
        this.playContainer.show();
    } else {
        this.playContainer.hide();
    }
};

Hand.prototype.setSuit = function(suit) {
    if ( this.suit === suit ) { return; }
    this.suit = suit;
    this.suitDiv.empty();
    this.suitDiv.append(`Suit: ${suit}`);
};

Hand.prototype.setLeader = function(leader) {
    if ( this.leader === leader ) { return; }
    this.leader = leader;
    this.leaderPlayerDiv.empty();
    this.leaderPlayerDiv.append(`Leader: ${leader}`);
};

Hand.prototype.setLeaderCard = function(leaderCard) {
    if ( equals(this.leaderCard, leaderCard) ) { return; }
    this.leaderCard = leaderCard;
    this.leaderCardDiv.empty();
    if ( !leaderCard ) { return; }
    this.leaderCardDiv.append(`Leader card: ${leaderCard.Number} of ${leaderCard.Suit}`);
};

Hand.prototype.setCardsPlayed = function(cardsPlayed) {
    if ( equals(this.cardsPlayed, cardsPlayed) ) { return; }
    this.cardsPlayed = cardsPlayed;
    let self = this;
    self.cardsPlayedTableBody.empty();
    cardsPlayed.forEach(function(playedCard) {
        let player = playedCard.Player;
        let card = playedCard.Card;
        let desc = card ? `${playedCard.Card.Number} of ${playedCard.Card.Suit}` : "";
        let style = (player === self.me) ? 'style="border: 1px dashed; padding: 8px; margin: 4px;"' : '';
        self.cardsPlayedTableBody.append(`
            <tr ${style}>
                <td>
                ${player}
                </td>
                <td>
                ${desc}
                </td>
            </tr>
        `)
    });
};

Hand.prototype.setState = function(state) {
    if ( state === this.state ) { return; }
    this.state = state;
    console.log(`setting hand state to ${state}`);
    switch (state) {
        case "NotJoined":
            this.suitDiv.hide();
            this.leaderContainer.hide();
            this.cardsPlayedTable.hide();
            this.playContainer.hide();
            break;
        case "WaitingForPlayers":
            this.suitDiv.hide();
            this.leaderContainer.hide();
            this.cardsPlayedTable.hide();
            this.playContainer.hide();
            break;
        case "RoundWagerTurn":
            this.suitDiv.hide();
            this.leaderContainer.hide();
            this.cardsPlayedTable.hide();
            this.playContainer.hide();
            break;
        case "RoundHandReady":
            this.suitDiv.hide();
            this.leaderContainer.hide();
            this.cardsPlayedTable.hide();
            this.playContainer.hide();
            break;
        case "HandPlayTurn":
            this.suitDiv.show();
            this.leaderContainer.show();
            this.cardsPlayedTable.show();
            // this.playContainer.show(); // this gets handled by setNextPlayer()
            break;
        case "RoundFinished":
            this.suitDiv.hide();
            this.leaderContainer.hide();
            this.cardsPlayedTable.hide();
            this.playContainer.hide();
            break;
        default:
            throw new Error(`invalid hand state ${state}`);
    }
};

// model

function Model() {
    let self = this;

    function didClickJoin(name) {
        self.join(name);
    }
    this.me = new Me(didClickJoin);

    function didClickRemovePlayer(player) {
        self.removePlayer(player);
    }
    function didChangeCardsPerPlayer(cardsPerPlayer) {
        self.changeCardsPerPlayer(cardsPerPlayer);
    }
    function didClickStart() {
        self.startRound();
    }
    this.game = new Game(didClickRemovePlayer, didChangeCardsPerPlayer, didClickStart);

    function didChooseWager(wager) {
        self.makeWager(wager);
    }
    function didClickStartHand() {
        self.startHand();
    }
    function didClickFinishRound() {
        self.finishRound();
    }
    this.round = new Round(didChooseWager, didClickStartHand, didClickFinishRound);

    function didClickPlayCard(card) {
        self.playCard(card);
    }
    this.hand = new Hand(didClickPlayCard);

    this.pollServer();
}

Model.prototype.pollServer = function() {
    if ( this.me.name ) {
        getMyModel(this.me.name, this.updateFromServer.bind(this));
    } else {
        getModel(this.updateFromServerFullModel.bind(this));
    }
    setTimeout(this.pollServer.bind(this), 5000);
};

Model.prototype.updateFromServerFullModel = function(ok, data) {
    this.game.update("", "NotJoined", data.Players, data.CardsPerPlayer);
    this.round.update("", "NotJoined", null, null, null, null, null);
};

Model.prototype.updateFromServer = function(ok, data) {
    if ( !ok ) { return; }

    let me = data.Me;
    this.me.update(me);

    let game = data.Game;
    this.game.update(me, data.State, game.Players, game.CardsPerPlayer);

    let round = data.Round;
    if ( round ) {
        this.round.update(me, data.State, round.TrumpSuit, round.Cards, round.Wagers, round.NextWagerPlayer);
    } else {
        this.round.update(me, data.State, null, null, null, null);
    }

    let hand = data.Hand;
    if ( hand ) {
        this.hand.update(me, data.State, hand.Suit, hand.Leader, hand.LeaderCard, hand.CardsPlayed, hand.NextPlayer, hand.Cards);
    } else {
        this.hand.update(me, data.State, "", "", null, [], "", []);
    }
};

// hitting the server

Model.prototype.join = function(name) {
    console.log(`joining game as ${name}`);
    if ( this.game.players.indexOf(name) >= 0 ) {
        getMyModel(name, this.updateFromServer.bind(this));
    } else {
        postJoin(name, this.updateFromServer.bind(this));
    }
};

Model.prototype.startHand = function() {
    postStartHand(this.me.name, this.updateFromServer.bind(this));
};

Model.prototype.removePlayer = function(player) {
    console.log(`removing player ${player}`);
    postRemovePlayer(this.me.name, player, this.updateFromServer.bind(this));
};

Model.prototype.changeCardsPerPlayer = function(count) {
    console.log(`changing cards per player to ${count}`);
    postSetCardsPerPlayer(this.me.name, count, this.updateFromServer.bind(this));
};

Model.prototype.startRound = function() {
    console.log(`starting round`);
    postStartRound(this.me.name, this.updateFromServer.bind(this));
};

Model.prototype.makeWager = function(hands) {
    console.log(`making wager ${hands}`);
    postWager(this.me.name, hands, this.updateFromServer.bind(this));
};

Model.prototype.playCard = function(card) {
    console.log(`playing card ${card.Number} of ${card.Suit}`);
    postPlayCard(this.me.name, card, this.updateFromServer.bind(this));
};

Model.prototype.finishRound = function() {
    postFinishRound(this.me.name, this.updateFromServer.bind(this));
};


//

let model = new Model();

});
