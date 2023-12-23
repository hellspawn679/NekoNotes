// // Tiptap.jsx
// import React from 'react';
// import { EditorProvider, FloatingMenu, BubbleMenu } from '@tiptap/react';
// import StarterKit from '@tiptap/starter-kit';
// import { Bold, Italic, Underline } from 'lucide-react';
// import '../CSS/Tiptap.css';
//
// const extensions = [
//     StarterKit,
//     Bold,
//     Italic,
//     Underline,
// ];
//
// const Tiptap = ({ content, onContentChange }) => {
//     return (
//         <EditorProvider extensions={extensions} content={content} onUpdate={(data) => onContentChange(data.editor.getHTML())}>
//             <div className="tiptap-navbar">
//                 <button onClick={() => window.tiptap.commands.toggleBold()}><Bold size={20} /></button>
//                 <button onClick={() => window.tiptap.commands.toggleItalic()}><Italic size={20} /></button>
//                 <button onClick={() => window.tiptap.commands.toggleUnderline()}><Underline size={20} /></button>
//             </div>
//             <FloatingMenu>
//                 {/* Additional floating menu items if needed */}
//             </FloatingMenu>
//             <BubbleMenu>This is the bubble menu</BubbleMenu>
//         </EditorProvider>
//     );
// };
//
// export default Tiptap;


import { useState } from "react";
import "../CSS/Tiptap.css";
import { Tiptap } from "../NotesComponents/TipTap";

function TipTap() {
    const [description, setDescription] = useState("");

    return (
        <div className="App">
            <Tiptap setDescription={setDescription} />
        </div>
    );
}

export default TipTap;
