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
    delete editor.keyBinding.$defaultHandler.commandKeyBinding["backspace"]
    document.getElementById("ace-editor").firstChild.addEventListener("keydown", (event) => {

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
      //console.log(curs);
      editor.selection.moveTo(curs.row, curs.column);
    }, msg => {
      let temp = JSON.parse(msg.data).body;
      let output = document.getElementById("output");
      output.innerHTML = temp;
    });
  }
  send(val) {
    if (changed) {
      sendMsg("1;:" + String(val));

      changed = false;
    } else {
      return;
    }
  }

  compile() {
    let code = String(editor.getValue());
    sendMsg("2;:" + code);
  }

  render() {
    return (
      <React.Fragment>
        <ReactAce
          mode="c_cpp"
          theme="monokai"
          name="Code"
          editorProps={{ $blockScrolling: true }}
          onChange={this.send}
          style={{ height: '400px' }}
          enableBasicAutocompletion="true"
          enableSnippets="true"
          ref={instance => { this.ace = instance; }} // Let's put things into scope
        />
        <button onClick={this.compile}>Compile</button>
        <div>
          <h3>
            Compiler output:
          </h3>
          <p id="output" />
        </div>
      </React.Fragment >
    );
  }
}

export default App;


