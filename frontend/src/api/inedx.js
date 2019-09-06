var socket = new WebSocket("ws://192.168.1.136:8080/ws");

let connect = cb =>{
    console.log("connecting");
    socket.onopen = () =>{
        console.log("Connection successfull");
    };

    socket.onmessage = msg => {
        //console.log(msg);
        cb(msg);
    };
    socket.onclose = event => {
        console.log("Socket cloesed: ", event);
    };

    socket.onerror = error =>{
        console.log("Error: ", error);
    };    
};

let sendMsg = msg =>{
    console.log("Sending:", msg);
    socket.send(msg);
};

export {connect, sendMsg};