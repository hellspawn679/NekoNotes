import React, { useState } from 'react';
import { useEditor, EditorContent } from '@tiptap/react';
import StarterKit from '@tiptap/starter-kit';
import Underline from '@tiptap/extension-underline';
import '../CSS/AiText.css';
import Chatbot from "./ChatBot";

function AiText() {
    const [inputText, setInputText] = useState('');
    const [messages, setMessages] = useState([]);

    const editor = useEditor({
        extensions: [StarterKit, Underline],
        content: '',
    });

    const handleSendMessage = () => {
        if (inputText.trim() !== '') {
            setMessages([...messages, { type: 'user', content: inputText }]);
            setInputText('');
        }
    };

    return (
        // <div className="TextAi ">
        //     <div className="qnaText">
        //         {/*<EditorContent editor={editor} />*/}
        //         <div className="input-container">
        //             <textarea
        //                 placeholder="Type your message..."
        //                 value={inputText}
        //                 onChange={(e) => setInputText(e.target.value)}
        //                 className="w-full h-20 border border-black rounded p-2 mt-2"
        //             />
        //             <button onClick={handleSendMessage} className="bg-blue-500 text-white p-2 mt-2 rounded">
        //                 Send
        //             </button>
        //         </div>
        //     </div>
        // </div>
        <>
            {/*<Chatbot/>*/}
        </>
    );
}

export default AiText;
