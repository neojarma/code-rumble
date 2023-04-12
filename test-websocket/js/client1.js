import WebSocket from 'ws';

const ws = new WebSocket('ws://127.0.0.1:8081/create-room')

ws.on("error", function (err) {
    console.log(err)
})

ws.on("open", async function open() {
    const obj = {
        roomName: "test-room",
        maxPlayer: 2,
        Player: {
            displayName: "neo-master"
        }
    }

    var strObj = JSON.stringify(obj)
    ws.send(strObj)

    console.log("simulate waiting other players, sleep in 7 second")
    const sleep = ms => new Promise(r => setTimeout(r, ms));
    await sleep(7000)
    console.log("try to start the game")

    const startws = new WebSocket('ws://127.0.0.1:8081/start?id=123')
    console.log("simulate doing task")
    await sleep(15000)

    console.log("try to submit code")
    const submitCode = new WebSocket('ws://127.0.0.1:8081/submit-code')
    submitCode.on("open", () => {
        var submitPayload = {
            playerId: "neoj",
            roomId: "123",
            submissionPayload: {
                questionId: "rapoeyztooztfyq",
                submittedCode: "export function birthday(s, d, m) {\r\n    var counter = 0\r\n    // Write your code here\r\n    for(var i = 0; i < s.length; i++) {\r\n        var temp = 0\r\n        for (var j = i; j < (m+i) && j < s.length; j++) {\r\n            temp += s[j]\r\n        }\r\n        \r\n        if (temp == d) {\r\n            counter++\r\n        } \r\n    }\r\n    \r\n    return counter\r\n}",
                customTestCase: false
            }
        }

        var submitString = JSON.stringify(submitPayload)
        submitCode.send(submitString)
    })
})



ws.on("message", function message(data) {
    var obj = JSON.parse(data.toString())
    console.log(obj)
})