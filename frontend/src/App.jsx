import React, {Component} from 'react';
import {connect, sendMsg} from "./api/inedx.js";



class App extends Component{
  constructor(props){
    super(props);
    this.state = {
      text: ""
    };
  }

  componentDidMount(){
    connect(msg =>{
      let temp = JSON.parse(msg.data).body;
      console.log(temp);
      this.setState(prevState => ({
        text: temp
      }));
      let area = document.getElementById("text");
      area.value = temp;
    });
  }
  send(event){
      //console.log("Hello");
      sendMsg(event.target.value);
  }

  render(){
    return(
    <div>
      <textarea id="text" onChange={this.send}></textarea>
    </div>
    );
  }
}

export default App;


