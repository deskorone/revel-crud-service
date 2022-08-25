let ws


function init() {
    ws = new WebSocket("ws://127.0.0.1:8080/app/feed?user=uservsalue")
    ws.onclose= ()=>{    
        console.log("CLOSED");   
    }
    ws.onmessage = (payload) => {
    console.log(payload.data);
}
}

init()



let but = document.getElementById("b")


$('#b').click(()=>{
    console.log("SEND")
    ws.send(JSON.stringify({message : "msg"}))
})


console.log("HELLo")