import "../CSS/Input.css";

import TipTapWithLogger from "./TipTap";
import Chatbot from "./ChatBot";

function Input() {
  return (
    <div className="input-container">
      <div className="tiptap-container">
        <TipTapWithLogger />
      </div>
      <div>
        <Chatbot />
      </div>
    </div>
  );
}

export default Input;
