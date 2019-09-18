var socket = new WebSocket("ws://typelias.tk:80/ws");


let connect = (codeMessage, compileMessage) => {
    console.log("connecting");
    socket.onopen = () => {
        console.log("Connection successfull");
    };

    socket.onmessage = msg => {
        let temp = JSON.parse(msg.data);
        if (temp.type === 1) {
            codeMessage(msg);
        } else if (temp.type === 2) {
            compileMessage(msg);
        } else {
            console.log("JSON type error");
        }

    };
    socket.onclose = event => {
        console.log("Socket cloesed: ", event);
    };

    socket.onerror = error => {
        console.log("Error: ", error);
    };
};

let sendMsg = msg => {
    console.log("Sending:", msg);
    socket.send(msg);

};

export { connect, sendMsg };
