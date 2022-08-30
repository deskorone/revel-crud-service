let arr;

let list = document.getElementById("list")

$.get("/hotels",
    function (data, textStatus) {
        for (let i in data) {
            $('#list').append(
                `<div class="cell">
                <div class="name">${data[i].name}</div>
                <div class="info">
                    <div class="price">${data[i].price}</div>
                    <div class="rating">${data[i].rating}</div>
                </div>
            </div>`)
        }

    },
);


let ws = new WebSocket("ws://localhost:8080/hotels/ws")

ws.onopen = ()=>{
    console.log("OPEN WSOCKET")
}

ws.onmessage = (payload)=>{
    console.log("HEllo")
    // data = JSON.parse(payload.data)
    // $('#list').append(
    //     `<div class="cell">
    //             <div class="name">${data.name}</div>
    //             <div class="info">
    //                 <div class="price">${data.price}</div>
    //                 <div class="rating">${data.rating}</div>
    //             </div>
    //         </div>`)
}

ws.onclose = ()=>{
    console.log("WS CLOSE")
}




