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

const useCardPicture = true;

function Card(suit, number) {
    if ( useCardPicture ) {
        let symbol = unicodeCards[suit][number];
        let [color, _] = suitToUnicode[suit];
        let klazz = `suit-${color}`;
        return `<div class="card card-picture ${klazz}" uadtr-card-suit="${suit}" uadtr-card-number="${number}">
            <div>${symbol}</div>
        </div>`;
    } else {
        let [color, symbol] = suitToUnicode[suit];
        let klazz = `suit-${color}`;
        return `<div class="card card-text wrapper-vertical ${klazz}" uadtr-card-suit="${suit}" uadtr-card-number="${number}">
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

    for ( let i = 0; i < a.length; ++i ) {
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
    for ( let key in a ) {
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
    this.getName = $("#me-get-name");
    this.showName = $("#me-show-name");
}

Me.prototype.update = function(name) {
    if ( name === this.name ) { return; }
    this.name = name;
    if ( name === "" ) {
        this.showName.hide();
        this.getName.show();
    } else {
        this.showName.empty();
        this.showName.append(name);
        this.showName.show();
        this.getName.hide();
    }
};

// Game

function Game(didClickRemovePlayer, didChangeCardsPerPlayer, didClickStartRound) {
    this.players = [];
    this.cardsPerPlayer = null;

    this.div = $("#game");
    this.startButton = $("#round-start-button");
    this.startButton.click(didClickStartRound);

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

    this.setStateWaitingForPlayers([], 1);
}

Game.prototype.setStateNotJoined = function(players) {
    this.setPlayers(players);
    this.div.show();
    $(".game-remove-player").hide();
    this.startButton.hide();
    this.cardsPerPlayerContainer.hide();
};

Game.prototype.setStateWaitingForPlayers = function(players, cardsPerPlayer) {
    this.setPlayers(players);
    this.setCardsPerPlayer(cardsPerPlayer);

    this.div.show();
    $(".game-remove-player").show();
    this.startButton.prop('disabled', this.players.length < 2);
    this.startButton.show();
    this.cardsPerPlayerContainer.show();
};

Game.prototype.setOtherStates = function() {
    this.div.hide();
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

    $(document).on("click", "#my-cards .card", function() {
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

MyCards.prototype.setOtherStates = function() {
    // do nothing
};

MyCards.prototype.setRoundFinished = function() {
    this.setCards([]);
};

MyCards.prototype.setWagerTurn = function(cards) {
    this.setCards(cards);
};

MyCards.prototype.setPlayCardTurn = function(cards, nextHandPlayer) {
    this.setCards(cards, nextHandPlayer);
    this.setNextHandPlayer(nextHandPlayer);
};

MyCards.prototype.setNextHandPlayer = function(nextHandPlayer) {
    if ( nextHandPlayer === this.nextHandPlayer ) { return; }
    this.nextHandPlayer = nextHandPlayer;
    if ( nextHandPlayer === this.me ) {
        console.log("enable hand clicks");
        $('#my-cards .card').css('pointer-events', 'auto');
        this.cardsTable.addClass("my-cards-my-turn");
    } else {
        console.log("disabling hand clicks");
        $('#my-cards .card').css('pointer-events', 'none');
        this.cardsTable.removeClass("my-cards-my-turn");
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

function Round(didClickFinishRound) {
    this.div = $("#round");
    this.trumpContainer = $("#trump-suit");

    this.finishRoundButton = $("#round-finish-button");
    this.finishRoundButton.click(didClickFinishRound);

    this.trumpSuit = "";

    this.setOtherStates();
}

Round.prototype.setTrumpSuit = function(trumpSuit) {
    if ( this.trumpSuit === trumpSuit ) { return; }
    this.trumpSuit = trumpSuit;
    this.trumpContainer.empty();
    let [color, symbol] = suitToUnicode[trumpSuit];
    let klazz = `suit-${color}`;
    for ( let i = 0; i < 5; i++ ) {
        this.trumpContainer.append(`<div class="${klazz}">${symbol}</div>`);
    }
};

Round.prototype.setOtherStates = function() {
    this.div.hide();
};

Round.prototype.setWagerTurn = function(trumpSuit) {
    this.setTrumpSuit(trumpSuit);
    this.finishRoundButton.hide();
    this.div.show();
};

Round.prototype.setPlayCardTurn = function(trumpSuit) {
    this.setTrumpSuit(trumpSuit);
    this.finishRoundButton.hide();
    this.div.show();
};

Round.prototype.setRoundFinished = function(trumpSuit) {
    this.setTrumpSuit(trumpSuit);
    this.finishRoundButton.show();
    this.div.show();
};

// status

function Status(didChooseWager) {
    this.cont = $("#status");
    this.statusTableBody = $("#status tbody");

    $(document).on("click", "#place-wager-button", function() {
        let wager = parseInt($("#place-wager-select").val(), 10);
        didChooseWager(wager);
    });

    this.setOtherStates();
}

Status.prototype.setOtherStates = function() {
    this.cont.hide();
};

Status.prototype.setWagerTurn = function(statuses, cardsPerPlayer) {
    this.setPlayerStatuses(statuses, cardsPerPlayer);
    this.cont.show();
};

Status.prototype.setPlayCardTurn = function(statuses) {
    this.setPlayerStatuses(statuses, null);
    this.cont.show();
};

Status.prototype.setRoundFinished = function(statuses) {
    this.setPlayerStatuses(statuses, null);
    this.cont.show();
};

function buildStatusTableModel(statuses, cardsPerPlayer) {
    let rows = [];
    statuses.forEach(function(status) {
        let wager = {
            'count': status.Wager,
        };
        if ( status.IsNextWagerer && status.IsMe ) {
            wager.options = [];
            for ( let i = 0; i <= cardsPerPlayer; i++ ) {
                wager.options.push(i);
            }
        }
        let cc = status.CurrentCard;
        let pc = status.PreviousCard;
        let row = {
            'classes': {
                'statuses-me': status.IsMe,
                'statuses-wager-turn': status.IsNextWagerer,
                'statuses-play-card-turn': status.IsNextPlayer,
                'statuses-leader': status.IsCurrentLeader,
                'statuses-previous-hand-winner': status.IsPreviousWinner,
            },
            'name': status.Player,
            'wager': wager,
            'handsWon': status.HandsWon,
            'prevCard': pc ? [pc.Suit, pc.Number] : null,
            'currCard': cc ? [cc.Suit, cc.Number] : null,
        };
        rows.push(row);
    });
    return rows;
}

function statusSuit(suit) {
    if ( !suit ) {
        return "";
    }
    let [color, symbol] = suitToUnicode[suit];
    let klazz = `suit-${color}`;
    let elems = ['<div class="status-suit wrapper-horizontal">'];
    for ( let i = 0; i < 3; i++ ) {
        elems.push(`<div class="${klazz}">${symbol}</div>`);
    }
    elems.push("</div>");
    return elems.join("\n");
}

Status.prototype.setPlayerStatuses = function(status, cardsPerPlayer) {
    let next = [status, cardsPerPlayer];
    if ( equals(this.current, next) ) {
        return;
    }
    this.current = next;

    let model = buildStatusTableModel(status.PlayerStatuses, cardsPerPlayer);
    this.statusTableBody.empty();
    let ph = status.PreviousHand;
    let ch = status.CurrentHand;
    this.statusTableBody.append(`
        <tr>
            <td></td>
            <td></td>
            <td></td>
            <td>${statusSuit(ph ? ph.Suit : null)}</td>
            <td>${statusSuit(ch ? ch.Suit : null)}</td>
        </tr>
    `);
    let self = this;
    model.forEach(function(status) {
        let wager;
        if ( status.wager.options ) {
            wager = [
                `<button id="place-wager-button">Place your wager!</button>`,
                `<select id="place-wager-select">`,
            ].concat(
                status.wager.options.map((i) => `<option value="${i}">${i}</option>`)
            );
            wager.push("</select>");
            wager = wager.join("\n");
        } else {
            wager = (status.wager.count !== null) ? status.wager.count : "";
        }
        let tds = `
            <td class="status-name">${status.name}</td>
            <td class="status-wager">${wager}</td>
            <td class="status-hands-won">${(status.handsWon !== null) ? status.handsWon : ""}</td>
            <td class="status-previous-card">${(status.prevCard !== null) ? Card(status.prevCard[0], status.prevCard[1]) : ""}</td>
            <td class="status-current-card">${(status.currCard !== null) ? Card(status.currCard[0], status.currCard[1]) : ""}</td>
        `;
        let klazzes = [];
        for ( let klazz in status.classes ) {
            if ( status.classes[klazz] ) {
                klazzes.push(klazz);
            }
        }
        let classes = klazzes.join(" ");
        self.statusTableBody.append(`<tr class="${classes}">${tds}</tr>`);
    });
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
    function didClickStartRound() {
        self.startRound();
    }
    this.game = new Game(didClickRemovePlayer, didChangeCardsPerPlayer, didClickStartRound);

    function didClickFinishRound() {
        self.finishRound();
    }
    this.round = new Round(didClickFinishRound);

    function didChooseWager(wager) {
        self.makeWager(wager);
    }
    this.status = new Status(didChooseWager);

    function didClickPlayCard(card) {
        self.playCard(card);
    }
    this.myCards = new MyCards(didClickPlayCard);

    this.pollServer();
}

Model.prototype.pollServer = function() {
    getMyModel(this.me.name, this.updateFromServer.bind(this));
    setTimeout(this.pollServer.bind(this), 2500);
};

Model.prototype.updateFromServer = function(ok, data) {
    if ( !ok ) { return; }

    let me = data.Me;
    let game = data.Game;
    this.me.update(me);
    this.myCards.me = me;
    switch (data.State) {
        case "NotJoined":
            this.game.setStateNotJoined(game.Players);
            this.myCards.setOtherStates();
            this.round.setOtherStates();
            this.status.setOtherStates();
            break;
        case "WaitingForPlayers":
            this.game.setStateWaitingForPlayers(game.Players, game.CardsPerPlayer);
            this.myCards.setOtherStates();
            this.round.setOtherStates();
            this.status.setOtherStates();
            break;
        case "WagerTurn":
            this.game.setOtherStates();
            this.myCards.setWagerTurn(data.MyCards);
            this.round.setWagerTurn(data.Status.TrumpSuit);
            this.status.setWagerTurn(data.Status, game.CardsPerPlayer);
            break;
        case "PlayCardTurn":
            let nextPlayer = data.Status.CurrentHand.NextPlayer;
            this.game.setOtherStates();
            this.myCards.setPlayCardTurn(data.MyCards, nextPlayer);
            this.round.setPlayCardTurn(data.Status.TrumpSuit);
            this.status.setPlayCardTurn(data.Status);
            break;
        case "RoundFinished":
            this.game.setOtherStates();
            this.myCards.setRoundFinished();
            this.round.setRoundFinished(data.Status.TrumpSuit);
            this.status.setRoundFinished(data.Status);
            break;
        default:
            throw new Error(`unrecognized state ${data.State}`);
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
    console.log("finishing round");
    postFinishRound(this.me.name, this.updateFromServer.bind(this));
};


//

let model = new Model();

});
