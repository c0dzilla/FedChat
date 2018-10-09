window.onload = function() {
    this.ws = new WebSocket('ws://' + window.location.host + '/ws');
    this.ws.addEventListener('message', function(e) {
        var msg = JSON.parse(e.data);
        // do something with received message
        document.getElementById("stream").appendChild(newMsg(msg));
    })
    document.getElementById("msg-submit").addEventListener('submit', function() {
        sendMessage()
    });
}

function sendMessage() {
    var message = document.getElementById("msg").innerHTML; // todo parse
    this.email = this.email || document.getElementById("email").innerHTML;
    this.username = this.username || document.getElementById("username").innerHTML;
    if (!email || !name) {
        alert('Invalid email and/or username');
    }
    this.ws.send(JSON.stringify(
        email,
        username,
        message
    ));
}

function newMsg(msg) {
    
}