// (function ($) {

$(document).ready(function() {

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

    function postStartRound(cont) {
        function f(data, status, _jqXHR) {
            console.log("post /player response -- status " + status);
            cont(status === 'success', data);
        }
        $.post({
            'url': '/TODO?',
            'data': JSON.stringify({'Player': name}),
            'dataType': 'json',
            'success': f,
            'error': f,
            'contentType': 'application/json'
        });
        console.log("fired off /player call");
    }

    // model

    console.log("hello");

    function Model() {
        this.me = "";
        this.players = [];
        this.listeners = {};
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

    Model.prototype.setPlayers = function(players) {
        console.log("settings players to " + JSON.stringify(players));
        this.players = players;
        this.notify("players");
    };

    Model.prototype.updateFromServer = function(ok, data) {
        if ( ok ) {
            this.setPlayers(data.Players);
        }
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

    let model = new Model();

    // views

    let me = $("#me");

    let game = $("#game");
    let gamePlayers = $("#game-players");
    let gameCardsPerPlayer = $("#game-cards-per-player");

    let round = $("#round");
    let hand = $("#hand");

    function configSetName() {
        let name = model.me;
        $("#me-show-name").append(name);
        $("#me-get-name").hide();
    }

    function gameSetPlayers() {
        // TODO kind of janky to call in to model instead of getting passed the data ... ?
        let names = model.players;
        //
        gamePlayers.empty();
        names.forEach(function(player) {
            gamePlayers.append("<li>" + player + "</li>")
        });

        if ( names.length === 0 ) {
            return;
        }
        let maxCards = 52 / names.length;
        gameCardsPerPlayer.empty();
        for (let i = 1; i < maxCards; i++) {
            gameCardsPerPlayer.append("<option value='" + i + "'>" + i + "</option>");
        }
        // TODO use gameCardsPerPlayer.val(prevSelection) to not blow away the selection unless necessary
    }

    // TODO these are sources, not sinks -- so maybe don't need these?
    // although: deck and number of players determines cards per player max
    // function gameSetCardsPerPlayer(count) {
    //
    // }
    //
    // function gameSetDeck(deckType) {
    //     // TODO
    // }

    function roundStart() {

    }

    // tie user actions to model

    $("#me-join").click(function join() {
        console.log("me-join click");
        let name = $("#me-input-name").val();
        model.join(name);
    });

    $("#round-start").click(function join() {
        console.log("submit-name click");
        model.startRound();
    });


    // initialization: poll server for model, tie model to ui updates

    model.listen("players", gameSetPlayers);
    model.listen("me", configSetName);

    function pollServer() {
        getModel((ok, data) => model.updateFromServer(ok, data));
        setTimeout(pollServer, 5000);
    }

    pollServer();

});
// })($);
