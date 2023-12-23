import React, { useState, useEffect } from "react";
import { useEditor } from "@tiptap/react";
import StarterKit from "@tiptap/starter-kit";
import Underline from "@tiptap/extension-underline";
import Loader from "./Loader";
import "../CSS/Chatbot.css"; // Import your custom CSS file

const Chatbot = () => {
  const [editorContent, setEditorContent] = useState("");
  const [question, setQuestion] = useState("");
  const [answer, setAnswer] = useState("");
  const [loading, setLoading] = useState("");
  const [plainText, setPlainText] = useState("");

  const editor = useEditor({
    extensions: [StarterKit, Underline],
    content: editorContent,
    onUpdate: ({ editor }) => {
      setEditorContent(editor.getHTML());
    },
  });

  const fetchNoteText = async () => {
    try {
      const apiUrl =
        "https://neko-notesbackendstorage.onrender.com/notebooks/657dfa6a9905d9766bfff39a";

      const response = await fetch(apiUrl);

      if (response.ok) {
        const data = await response.json();
        console.log("Plain Text ----- ", data.notes[0].text);
        setPlainText(data.notes[0].text);
        askQuestion();
      } else {
        console.error("Error fetching note text:", response.statusText);
      }
    } catch (error) {
      console.error("Error fetching note text:", error.message);
    }
  };

  const askQuestion = async () => {
    if (!question.trim()) return;

    setLoading(true);

    try {
      const apiUrl = "https://nekonotes-7gmw.onrender.com/qna";

      const response = await fetch(apiUrl, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          inputs: {
            context: plainText,
            question: question.trim(),
          },
        }),
      });

      if (response.ok) {
        const data = await response.json();
        console.log(data);
        setAnswer(data.answer);
      } else {
        console.error("Error asking question:", response.statusText);
      }
    } catch (error) {
      console.error("Error asking question:", error.message);
    }

    setLoading(false);
  };

  useEffect(() => {
    if (editor) {
      setEditorContent(editor.getHTML());
    }
  }, [editor]);

  return (
    <div className="custom-flex custom-fill-viewport">
      <div className="custom-chatbot-container">
        <div className="custom-chatbot-interaction">
          <div className="custom-chatbot-question">
            <input
              type="text"
              placeholder="Ask a question..."
              value={question}
              onChange={(e) => setQuestion(e.target.value)}
              className="custom-border custom-rounded-md custom-p-2 custom-mr-2 custom-flex-grow"
            />
            <button
              onClick={fetchNoteText}
              disabled={loading}
              className="custom-bg-blue-800 custom-text-white custom-px-4 custom-py-2 custom-rounded-md custom-submit-button"
            >
              {loading ? <Loader /> : "Submit"}
            </button>
          </div>
        </div>
        {answer && (
          <div className="custom-chatbot-answer">
            <h3 className="custom-text-xl custom-font-bold custom-mb-2">
              Answer:
            </h3>
            <p>{answer}</p>
          </div>
        )}
      </div>
    </div>
  );
};

export default Chatbot;
