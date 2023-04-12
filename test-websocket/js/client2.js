import WebSocket from 'ws';

const ws = new WebSocket('ws://127.0.0.1:8081/join-room')

ws.on("error", function (err) {
    console.log(err)
})

ws.on("open", function open() {
    const obj = {
        roomId: "123",
        Player: {
            displayName: "neo-player"
        }
    }

    var strObj = JSON.stringify(obj)
    ws.send(strObj)
    console.log("ok")
})

ws.on("message", function message(data) {
    var obj = JSON.parse(data.toString())
    console.log(obj)
})