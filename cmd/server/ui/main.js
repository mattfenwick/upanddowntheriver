'use strict';

$(document).ready(function() {

// javascript card characters

let unicodeCards = {
    "Clubs": {
        "2" : "\uD83C\uDCD2",
        "3" : "\uD83C\uDCD3",
        "4" : "\uD83C\uDCD4",
        "5" : "\uD83C\uDCD5",
        "6" : "\uD83C\uDCD6",
        "7" : "\uD83C\uDCD7",
        "8" : "\uD83C\uDCD8",
        "9" : "\uD83C\uDCD9",
        "10": "\uD83C\uDCDA",
        "J" : "\uD83C\uDCDB",
        "Q" : "\uD83C\uDCDD",
        "K" : "\uD83C\uDCDE",
        "A" : "\uD83C\uDCD1",
    },
    "Diamonds": {
        "2" : "\uD83C\uDCC2",
        "3" : "\uD83C\uDCC3",
        "4" : "\uD83C\uDCC4",
        "5" : "\uD83C\uDCC5",
        "6" : "\uD83C\uDCC6",
        "7" : "\uD83C\uDCC7",
        "8" : "\uD83C\uDCC8",
        "9" : "\uD83C\uDCC9",
        "10": "\uD83C\uDCCA",
        "J" : "\uD83C\uDCCB",
        "Q" : "\uD83C\uDCCD",
        "K" : "\uD83C\uDCCE",
        "A" : "\uD83C\uDCC1",
    },
    "Hearts": {
        "2" : "\uD83C\uDCB2",
        "3" : "\uD83C\uDCB3",
        "4" : "\uD83C\uDCB4",
        "5" : "\uD83C\uDCB5",
        "6" : "\uD83C\uDCB6",
        "7" : "\uD83C\uDCB7",
        "8" : "\uD83C\uDCB8",
        "9" : "\uD83C\uDCB9",
        "10": "\uD83C\uDCBA",
        "J" : "\uD83C\uDCBB",
        "Q" : "\uD83C\uDCBD",
        "K" : "\uD83C\uDCBE",
        "A" : "\uD83C\uDCB1",
    },
    "Spades": {
        "2" : "\uD83C\uDCA2",
        "3" : "\uD83C\uDCA3",
        "4" : "\uD83C\uDCA4",
        "5" : "\uD83C\uDCA5",
        "6" : "\uD83C\uDCA6",
        "7" : "\uD83C\uDCA7",
        "8" : "\uD83C\uDCA8",
        "9" : "\uD83C\uDCA9",
        "10": "\uD83C\uDCAA",
        "J" : "\uD83C\uDCAB",
        "Q" : "\uD83C\uDCAD",
        "K" : "\uD83C\uDCAE",
        "A" : "\uD83C\uDCA1",
    }
};

let suitToUnicode = {
    'Diamonds': ['red', '\u2666'],
    'Hearts': ['red', '\u2665'],
    'Clubs': ['black', '\u2663'],
    'Spades': ['black', '\u2660'],
};

const useUnicodeCard = true;

function Card(suit, number) {
    if ( useUnicodeCard ) {
        let symbol = unicodeCards[suit][number];
        let [color, _] = suitToUnicode[suit];
        let klazz = `suit-${color}`;
        return `<div class="round-card round-card-unicode-style ${klazz}" uadtr-card-suit="${suit}" uadtr-card-number="${number}">
        <div>${symbol}</div>
    </div>`;
    } else {
        let [color, symbol] = suitToUnicode[suit];
        let klazz = `suit-${color}`;
        return `<div class="round-card round-card-style wrapper-vertical ${klazz}" uadtr-card-suit="${suit}" uadtr-card-number="${number}">
        <div>${number}</div>
        <div>${symbol}</div>
    </div>`;
    }
}

// network actions

// function getModel(cont) {
//     function f(data, status, _jqXHR) {
//         console.log("GET model response -- status " + status);
//         cont(status === 'success', data);
//     }
//     $.ajax({
//         'url': '/model',
//         'dataType': 'json',
//         'success': f,
//         'error': f,
//     });
//     console.log("fired off GET to /model");
// }

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

function postFinishHand(me, cont) {
    postAction({'Me': me, 'FinishHand': {}}, cont);
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

    this.div = $("#game");
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
    this.cardsPerPlayerContainer = $("#card-per-player-container");

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
            this.div.show();
            $(".game-remove-player").hide();
            this.startButton.hide();
            this.cardsPerPlayerContainer.hide();
            break;
        case "WaitingForPlayers":
            this.div.show();
            $(".game-remove-player").show();
            this.startButton.prop('disabled', this.players.length < 2);
            this.startButton.show();
            this.cardsPerPlayerContainer.show();
            break;
        case "RoundWagerTurn":
            this.div.hide();
            break;
        case "RoundHandReady":
            this.div.hide();
            break;
        case "HandPlayTurn":
            this.div.hide();
            break;
        case "HandFinished":
            this.div.hide();
            break;
        case "RoundFinished":
            this.div.hide();
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
    // TODO pull out this 52 constant -- have server report number of cards
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

// my cards

function MyCards(didClickPlayCard) {
    this.cardsTable = $("#my-cards");

    $(document).on("click", "#my-cards .round-card", function() {
        let suit = $(this).attr('uadtr-card-suit');
        let number = $(this).attr('uadtr-card-number');
        let card = {'Suit': suit, 'Number': number};
        console.log(`round-play-card: ${card}`);
        didClickPlayCard(card);
    });

    this.me = "";
    this.cards = [];
    this.nextHandPlayer = "";
}

MyCards.prototype.update = function(me, cards, nextHandPlayer) {
    this.me = me;
    if ( cards ) {
        this.setCards(cards);
    }
    if ( nextHandPlayer ) {
        this.setNextHandPlayer(nextHandPlayer);
    }
};

MyCards.prototype.setNextHandPlayer = function(nextHandPlayer) {
    if ( nextHandPlayer === this.nextHandPlayer ) { return; }
    this.nextHandPlayer = nextHandPlayer;
    if ( nextHandPlayer === this.me ) {
        $('.round-card').css('pointer-events', 'auto');
        this.cardsTable.addClass("round-my-turn");
    } else {
        $('.round-card').css('pointer-events', 'none');
        this.cardsTable.removeClass("round-my-turn");
    }
};

MyCards.prototype.setCards = function(cards) {
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
        let suitDivs = [];
        cardsBySuit[suit].forEach(function(number) {
            suitDivs.push(Card(suit, number));
        });
        let row = `<div class="wrapper-horizontal">${suitDivs.join("\n")}</div>`;
        self.cardsTable.append(row);
    });
};

// round

function Round(didChooseWager, didClickPlayCard, didClickStartHand, didClickFinishHand, didClickFinishRound) {
    this.div = $("#round");
    this.wagersTableBody = $("#wagers tbody");
    this.trumpContainer = $("#trump-suit");

    this.myCards = new MyCards(didClickPlayCard);
    this.hand = new Hand(didClickStartHand, didClickFinishHand, didClickFinishRound);

    $(document).on("click", "#place-wager-button", function() {
        let wager = parseInt($("#place-wager-select").val(), 10);
        didChooseWager(wager);
    });

    this.me = "";
    this.nextWagerPlayer = "";
    this.trumpSuit = "";
    this.state = "";

    this.setRoundState("WaitingForPlayers");
    this.setWagers([], "");
}

Round.prototype.update = function(me, state, trumpSuit, cards, wagers, nextWagerPlayer, nextHandPlayer, cardsPerPlayer, hand) {
    this.me = me;
    this.myCards.update(me, cards, nextHandPlayer);
    this.cardsPerPlayer = cardsPerPlayer;
    if ( wagers ) {
        this.setWagers(wagers, nextWagerPlayer);
    }
    if ( trumpSuit ) {
        this.setTrumpSuit(trumpSuit);
    }
    this.setRoundState(state);
    if ( hand ) {
        this.hand.update(me, state, hand.Suit, hand.Leader, hand.CardsPlayed);
    } else {
        this.hand.update(me, state, null, null, null);
    }
};

Round.prototype.setTrumpSuit = function(trumpSuit) {
    if ( this.trumpSuit === trumpSuit ) { return; }
    this.trumpSuit = trumpSuit;
    this.trumpContainer.empty();
    let [color, symbol] = suitToUnicode[trumpSuit];
    let klazz = `suit-${color}`;
    for ( let i = 0; i < 5; i++ ) {
        this.trumpContainer.append(`<div class="trump-suit ${klazz}">${symbol}</div>`);
    }
};

Round.prototype.setRoundState = function(state) {
    if ( state === this.state ) { return; }
    this.state = state;
    console.log(`setting round state to ${state}`);
    switch (state) {
        case "NotJoined":
        case "WaitingForPlayers":
            this.div.hide();
            break;
        case "RoundWagerTurn":
        case "RoundHandReady":
        case "HandPlayTurn":
        case "HandFinished":
        case "RoundFinished":
            this.div.show();
            break;
        default:
            throw new Error(`invalid round state ${state}`);
    }
};

Round.prototype.setWagers = function(wagers, nextWagerPlayer) {
    if ( equals(this.wagers, wagers) && ( nextWagerPlayer === this.nextWagerPlayer ) ) {
        return;
    }
    this.nextWagerPlayer = nextWagerPlayer;
    this.wagers = wagers;

    this.wagersTableBody.empty();
    let self = this;
    wagers.forEach(function(wager) {
        let player = wager.Player;
        let wagerHtml = "";
        if ( nextWagerPlayer === self.me && nextWagerPlayer === player ) {
            let elems = [];
            elems.push(`<button id="place-wager-button">Place your wager!</button>`);
            elems.push(`<select id="place-wager-select">`);
            for ( let i = 0; i <= self.cardsPerPlayer; i++ ) {
                elems.push(`<option value="${i}">${i}</option>`);
            }
            elems.push("</select>");
            wagerHtml = elems.join("\n");
        } else {
            wagerHtml = (wager.Count !== null) ? wager.Count : "";
        }
        let wonCount = (wager.HandsWon !== null) ? wager.HandsWon : "";
        let klazz = (player === self.me) ? 'wager-me' : 'wager-not-me';
        self.wagersTableBody.append(`
            <tr class="${klazz}">
                <td>
                ${player}
                </td>
                <td>
                ${wagerHtml}
                </td>
                <td>
                ${wonCount}
                </td>
            </tr>
        `)
    });
};

// hand

function Hand(didClickStartHand, didClickFinishHand, didClickFinishRound) {
    this.me = "";
    this.state = "";
    this.suit = "";
    this.leader = "";
    this.cardsPlayed = [];

    this.div = $("#hand");
    this.suitDiv = $("#hand-suit");
    this.cardsPlayedTable = $("#hand-cards");
    this.cardsPlayedTableBody = $("#hand-cards tbody");

    this.startHandButton = $("#hand-start-button");
    this.finishHandButton = $("#hand-finish-button");
    this.finishRoundButton = $("#round-finish-button");

    this.setState("NotJoined");

    this.startHandButton.click(didClickStartHand);
    this.finishHandButton.click(didClickFinishHand);
    this.finishRoundButton.click(didClickFinishRound);
}

Hand.prototype.update = function(me, state, suit, leader, cardsPlayed) {
    this.me = me;
    this.setState(state);
    // if ( suit ) {
        this.setSuit(suit);
    // }
    // if ( leader ) {
        this.setLeader(leader);
    // }
    if ( cardsPlayed ) {
        this.setCardsPlayed(cardsPlayed);
    }
};

Hand.prototype.setSuit = function(suit) {
    if ( this.suit === suit ) { return; }
    this.suit = suit;
    this.suitDiv.empty();
    if ( !suit ) { return; }
    let [color, symbol] = suitToUnicode[suit];
    let klazz = `suit-${color}`;
    for ( let i = 0; i < 3; i++ ) {
        this.suitDiv.append(`<div class="hand-suit ${klazz}">${symbol}</div>`);
    }
};

Hand.prototype.setLeader = function(leader) {
    if ( this.leader === leader ) { return; }
    this.leader = leader;
};

Hand.prototype.setCardsPlayed = function(cardsPlayed) {
    if ( equals(this.cardsPlayed, cardsPlayed) ) { return; }
    this.cardsPlayed = cardsPlayed;
    let self = this;
    self.cardsPlayedTableBody.empty();
    cardsPlayed.forEach(function(playedCard) {
        let player = playedCard.Player;
        let card = playedCard.Card;
        let desc = card ? Card(playedCard.Card.Suit, playedCard.Card.Number) : "";
        let klazz = (player === self.leader) ? 'hand-leader' : 'hand-not-leader';
        self.cardsPlayedTableBody.append(`
            <tr class="${klazz}">
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
            this.div.hide();
            break;
        case "WaitingForPlayers":
            this.div.hide();
            break;
        case "RoundWagerTurn":
            this.div.hide();
            break;
        case "RoundHandReady":
            this.div.show();
            this.startHandButton.show();
            this.finishHandButton.hide();
            this.finishRoundButton.hide();
            break;
        case "HandPlayTurn":
            this.div.show();
            this.suitDiv.show();
            this.cardsPlayedTable.show();
            this.startHandButton.hide();
            this.finishHandButton.hide();
            this.finishRoundButton.hide();
            break;
        case "HandFinished":
            this.div.show();
            this.suitDiv.show();
            this.cardsPlayedTable.show();
            this.startHandButton.hide();
            this.finishHandButton.show();
            this.finishRoundButton.hide();
            break;
        case "RoundFinished":
            this.div.show();
            this.suitDiv.hide();
            this.cardsPlayedTable.hide();
            this.startHandButton.hide();
            this.finishHandButton.hide();
            this.finishRoundButton.show();
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
    function didClickPlayCard(card) {
        self.playCard(card);
    }
    function didClickStartHand() {
        self.startHand();
    }
    function didClickFinishHand() {
        self.finishHand();
    }
    function didClickFinishRound() {
        self.finishRound();
    }
    this.round = new Round(didChooseWager, didClickPlayCard, didClickStartHand, didClickFinishHand, didClickFinishRound);

    this.pollServer();
}

Model.prototype.pollServer = function() {
    getMyModel(this.me.name, this.updateFromServer.bind(this));
    setTimeout(this.pollServer.bind(this), 2500);
};

Model.prototype.updateFromServer = function(ok, data) {
    if ( !ok ) { return; }

    let me = data.Me;
    this.me.update(me);

    let game = data.Game;
    this.game.update(me, data.State, game.Players, game.CardsPerPlayer);

    let round = data.Round;
    if ( round ) {
        this.round.update(
            me,
            data.State,
            round.TrumpSuit,
            round.Cards,
            round.Wagers,
            round.NextWagerPlayer,
            hand ? hand.NextPlayer : null,
            game.CardsPerPlayer,
            data.Hand);
    } else {
        this.round.update(me, data.State, null, null, null, null, null, null);
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

Model.prototype.finishHand = function() {
    console.log("finishing hand");
    postFinishHand(this.me.name, this.updateFromServer.bind(this));
};

Model.prototype.finishRound = function() {
    console.log("finishing round");
    postFinishRound(this.me.name, this.updateFromServer.bind(this));
};


//

let model = new Model();

});
