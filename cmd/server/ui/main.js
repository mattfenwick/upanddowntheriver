'use strict';

// (function ($) {

$(document).ready(function() {

    // util

    function arrayEquals(a, b) {
        if ( a.length != b.length ) {
            return false;
        }

        for (let i = 0; i < a.length; ++i) {
            if ( a[i] !== b[i] ) {
                return false;
            }
        }
        return true;
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

    function postJoin(name, cont) {
        function f(data, status, _jqXHR) {
            console.log("post /player response -- status " + status);
            cont(status === 'success', data);
        }
        $.post({
            'url': '/player',
            'data': JSON.stringify({'Player': name}),
            'dataType': 'json',
            'success': f,
            'error': f,
            'contentType': 'application/json'
        });
        console.log("fired off POST to /player");
    }

    function deletePlayer(name, cont) {
        function f(data, status, _jqXHR) {
            console.log("delete /player response -- status " + status);
            cont(status === 'success', data);
        }
        $.ajax({
            'url': '/player',
            'type': 'DELETE',
            'data': JSON.stringify({'Player': name}),
            'dataType': 'json',
            'success': f,
            'error': f,
            'contentType': 'application/json'
        });
        console.log("fired off DELETE to /player");
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

    function postStartRound(me, cont) {
        postAction({'StartRound': {'Player': me}}, cont)
    }

    function postWager(me, hands, cont) {
        postAction({'MakeWager': {'Player': me, 'Hands': hands}}, cont)
    }

    function postPlayCard(me, number, suit, cont) {
        postAction({'PlayCard': {'Player': me, 'Card': {'Number': number, 'Suit': suit}}}, cont)
    }

    // model

    console.log("hello");

    function Model() {
        this.me = "";
        this.players = [];
        this.listeners = {};
        this.data = null;
    }

    Model.prototype.listen = function(key, action) {
        if ( !(key in this.listeners) ) {
            this.listeners[key] = [];
        }
        this.listeners[key].push(action);
    };

    Model.prototype.notify = function(key) {
        console.log(JSON.stringify(Object.keys(this.listeners)));
        console.log(key);
        this.listeners[key].forEach(function(f) {
            f();
        });
    };

    Model.prototype.updateFromServer = function(ok, data) {
        if ( !ok ) { return; }
        this.data = data;
        this.setPlayers(data.Players);
        this.setRounds(data.Rounds);
    };

    Model.prototype.setPlayers = function(players) {
        console.log("settings players to " + JSON.stringify(players));
        if ( !arrayEquals(this.players, players) ) {
            this.players = players;
            this.notify("players");
        }
    };

    Model.prototype.removePlayer = function(player) {
        console.log(`removing player ${player}`);
        let self = this;
        deletePlayer(player, function (ok, data) {
            if ( ok ) {
                self.players = self.players.filter((p) => p !== player);
                self.notify("players");
            }
        });
    };

    Model.prototype.join = function(name) {
        let self = this;
        postJoin(name, function (ok, data) {
            if ( ok ) {
                self.me = name;
                self.notify("me");
            }
        });
    };

    Model.prototype.setRounds = function(rounds) {

    };

    Model.prototype.startRound = function() {
        let self = this;
        postStartRound(this.me, function (ok, data) {
            self.notify("roundState");
            console.log(`start round: ${ok}`);
        });
    };

    Model.prototype.makeWager = function(hands) {
        postWager(this.me, hands, function() {
            // TODO update ui to show that it's no longer "my" turn
            console.log(`make wager: ${ok}`);
        })
    };

    Model.prototype.currentRound = function() {
        let rounds = this.data.Rounds[this.data.length - 1];
    };

    let model = new Model();

    // views

    let me = $("#me");

    function configSetName() {
        let name = model.me;
        $("#me-show-name").append(name);
        $("#me-get-name").hide();
    }

    let game = $("#game");
    let gamePlayers = $("#game-players");
    let gameCardsPerPlayer = $("#game-cards-per-player");

    function gameSetPlayers() {
        // TODO kind of janky to call in to model instead of getting passed the data ... ?
        let names = model.players;
        //
        gamePlayers.empty();
        names.forEach(function(player) {
            gamePlayers.append(`
                <tr>
                    <td>${player}</td>
                    <td>
                        <button class='remove-player' player='${player}'>Remove</button>
                    </td>
                </tr>`);
        });

        if ( names.length === 0 ) {
            return;
        }
        let maxCards = 52 / names.length;
        gameCardsPerPlayer.empty();
        for (let i = 1; i < maxCards; i++) {
            gameCardsPerPlayer.append(`<option value='${i}'>${i}</option>`);
        }
        // TODO use gameCardsPerPlayer.val(prevSelection) to not blow away the selection unless necessary
    }

    let round = $("#round");
    let roundStartButton = $("#round-start");
    let roundFinishButton = $("#round-finish");
    let roundCards = $("#round-cards");
    let roundWagers = $("#round-wagers");

    function roundState() {
        let round = model.currentRound();
        switch (round.State) {
            case "RoundStateCardsDealt":
                roundStartButton.hide();
                roundCards.show();
                roundWagers.show();
                break;
            case "RoundStateWagersMade":
                break;
            case "RoundStateHandInProgress":
                break;
            case "RoundStateFinished":
                roundFinishButton.show();
                break;
            default:
                throw new Error(`invalid round state ${round.State}`);
        }
        // TODO what about after wrapping up a round, we now need to show the round-start button to kick
        // the next round off?
    }

    let hand = $("#hand");

    // tie user actions to model

    $("#me-join").click(function () {
        console.log("me-join click");
        let name = $("#me-input-name").val();
        model.join(name);
    });

    roundStartButton.click(function () {
        console.log("submit-name click");
        model.startRound();
    });

    // can't use `click` because the elements might be ever-changing
    $(document).on("click", ".remove-player", function() {
    // $(".remove-player").click(function () {
        let player = $(this).attr('player');
        console.log("remove-player: ${player}");
        model.removePlayer(player);
    });

    // initialization: poll server for model, tie model to ui updates

    model.listen("players", gameSetPlayers);
    model.listen("me", configSetName);
    model.listen("startRound", roundStart);

    function pollServer() {
        getModel((ok, data) => model.updateFromServer(ok, data));
        setTimeout(pollServer, 5000);
    }

    pollServer();

});
// })($);
