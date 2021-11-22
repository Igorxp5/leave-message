const socketProtocol = (window.location.protocol == "https:") ? "wss" : "ws";
const socket = new WebSocket(`${socketProtocol}://${window.location.hostname}:${window.location.port}/ws`);

var clientId = null;

socket.addEventListener("open", function (event) {
    console.log("Connected!");
});

socket.addEventListener("message", function (event) {
    data = JSON.parse(event.data);
    if (data["event"] == "clientId") {
        clientId = data["message"];
    } else if (data["event"] == "error") {
        alert(data["message"]);
    } else if (data["event"] == "queuePosition") {
        document.getElementById("queuePosition").innerHTML = data["message"];
        if (data["message"] == 1) {
            window.location = `/message?clientId=${clientId}`;
        }
    }
});
