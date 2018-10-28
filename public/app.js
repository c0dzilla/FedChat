this.ws = new WebSocket('ws://' + window.location.host + '/ws');
this.ws.addEventListener('message', function(e) {
    var msg = JSON.parse(e.data);
    // do something with received message
    newMsg(msg);
})
document.getElementById("msg-submit").addEventListener('click', function() {
    sendMessage()
});

var email = null;
var username = null;

function sendMessage() {
    var message = document.getElementById("msg").value; // todo parse
    email = document.getElementById("email").value;
    username = document.getElementById("username").value;
    if (!email || !username) {
        alert('Invalid email and/or username');
        return;
    }
    this.ws.send(JSON.stringify({
        email,
        username,
        message
    }));
}

function newMsg(msg) {
    var msgCon = (document.getElementsByClassName("msgCon")[0]).cloneNode(true);
    msgCon.children[0].innerHTML = msg.email;
    msgCon.children[1].innerHTML = msg.username;
    msgCon.children[2].innerHTML = msg.message;
    document.getElementById("msgList").appendChild(msgCon);
}