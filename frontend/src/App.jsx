import React, { Component } from 'react';
import { connect, sendMsg } from "./api/inedx.js";
import ReactAce from "react-ace-editor"

import "brace/theme/monokai"
import "brace/mode/c_cpp"
let editor = null;
let changed = false;

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      text: ""
    };
  }

  componentDidMount() {
    editor = this.ace.editor;
    editor.$blockScrolling = Infinity;
    //editor.keyBinding.$defaultHandler.commandKeyBinding = {};
    document.getElementById("ace-editor").firstChild.addEventListener("keydown", (event) => {
      //console.log(event.key);
      if (event.key === "Shift" ||
        event.key === "Backspace" ||
        event.key === "Control" ||
        event.key === "AltGraph" ||
        event.key === "Alt" ||
        event.key === "ArrowUp" ||
        event.key === "ArrowDown" ||
        event.key === "ArrowLeft" ||
        event.key === "ArrowRight") {
        console.log(event.key);
        return;
      }
      //sendMsg(String(editor.getValue()))
      changed = true;
    });
    connect(msg => {
      let temp = JSON.parse(msg.data).body;
      console.log(temp);
      this.setState(prevState => ({
        text: temp
      }));
      let curs = editor.getCursorPosition();
      editor.setValue(temp);
      editor.clearSelection();
      console.log(curs);
      editor.selection.moveTo(curs.row, curs.column);
    });
  }
  send(val) {
    if (changed) {
      sendMsg(String(val));

      changed = false;
    } else {
      return;
    }
  }

  render() {
    return (
      <ReactAce
        mode="c_cpp"
        theme="monokai"
        name="Code"
        editorProps={{ $blockScrolling: true }}
        onChange={this.send}
        style={{ height: '400px' }}
        ref={instance => { this.ace = instance; }} // Let's put things into scope
      />
    );
  }
}

export default App;


