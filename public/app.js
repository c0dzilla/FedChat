window.onload = function() {
    this.ws = new WebSocket('ws://' + window.location.host + '/ws');
    this.ws.addEventListener('message', function(e) {
        var msg = JSON.parse(e.data);
        // do something with received message
        newMsg(msg);
    })
    document.getElementById("msg-submit").addEventListener('submit', function() {
        sendMessage()
    });
}

function sendMessage() {
    var message = document.getElementById("msg").innerHTML; // todo parse
    this.email = this.email || document.getElementById("email").value;
    this.username = this.username || document.getElementById("username").value;
    if (!email || !name) {
        alert('Invalid email and/or username');
        return;
    }
    this.ws.send(JSON.stringify(
        email,
        username,
        message
    ));
}

function newMsg(msg) {
    var msgCon = (document.getElementsByClassName("msgCon")[0]).cloneNode(true);
    msgCon.children[0].innerHTML = msg.email;
    msgCon.children[1].innerHTML = msg.username;
    msgCon.children[2].innerHTML = msg.email;
    document.getElementById("msgList").appendChild(msgCon);
}