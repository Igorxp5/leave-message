var timer = 60;
var timerInterval = setInterval(function() {
    document.getElementById("coutdown").innerHTML = timer + " ";
    timer -= 1;
    if (timer == 0) {
        clearInterval(timerInterval);
        alert("Your session has expired!");
        window.location = "/";
    }
}, 1000);


document.getElementById("form").addEventListener("submit", function(event) {
    let text = document.getElementById("message").value;
    let clientId = document.getElementById("clientId").value;
    fetch("/api/message", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({clientId: clientId, text: text})
    }).then(function(data) {
        if (data.error) {
            alert(data.error);
        }
        window.location = "/";
    });
    event.preventDefault();
});
