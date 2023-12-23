import React from "react";
import NotesSideBar from "./NotesSideBar";
import Input from "./Input";
import AiText from "./AiText";

import "../CSS/Notes.css";

function Notes() {
  return (
    <div className="notes-section">
      <div className="right">
        <Input />
        <AiText />
      </div>
      {/* <div className="left">
        <NotesSideBar />
      </div> */}
    </div>
  );
}

export default Notes;
