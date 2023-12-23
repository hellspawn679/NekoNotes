import React, { useEffect, useState } from "react";
import { useEditor, EditorContent } from "@tiptap/react";
import StarterKit from "@tiptap/starter-kit";
import Underline from "@tiptap/extension-underline";
import { FontSize } from "./textStyleExtensions";
import {
  FaBold,
  FaItalic,
  FaListOl,
  FaListUl,
  FaRedo,
  FaStrikethrough,
  FaUnderline,
  FaUndo,
} from "react-icons/fa";
import Select from "react-select"; // Import react-select
import "../CSS/Tiptap.css";
import Loader from "./Loader";
import html2pdf from "html2pdf.js";

const PlainTextLogger = () => {
  return {
    name: "plainTextLogger",
    onCreate(editor) {
      editor.on("update", ({ state }) => {
        const plainTextContent = state.doc.textBetween(
          0,
          state.doc.content.size,
          "\n\n"
        );
        console.log("Plain Text Content:", plainTextContent);
        // Send plain text content to your backend or perform any other actions
      });
    },
  };
};

export const SummarizeButton = ({ editor, summarizeContent, loading }) => {
  return (
    <button
      onClick={() => summarizeContent(editor)}
      style={{
        padding: "5px 10px",
        fontSize: "16px",
        backgroundColor: "mediumvioletred",
        color: "rgb(255,255,255)",
        border: "1px solid black",
        borderRadius: "4px",
        cursor: "pointer",
        marginRight: "10px", // Adjust spacing
      }}
    >
      {loading ? <Loader /> : "Summary"}
    </button>
  );
};

export const TranslateButton = ({ editor, translateContent, loading }) => {
  const [selectedLanguage, setSelectedLanguage] = useState(null);

  // Language options for the dropdown
  const languageOptions = [
    { value: "hi_XX", label: "Hindi" },
    { value: "fr_XX", label: "French" },
    { value: "it_XX", label: "Italian" },
    { value: "de_XX", label: "German" },
    { value: "zh_XX", label: "Chinese" },
    { value: "es_XX", label: "Spanish" },
    { value: "ru_XX", label: "Russian" },
    { value: "ja_XX", label: "Japanese" },
    { value: "pt_XX", label: "Portuguese" },
    { value: "ko_XX", label: "Korean" },
    { value: "ar_XX", label: "Arabic" },
    { value: "pl_XX", label: "Polish" },
    { value: "vi_XX", label: "Vietnamese" },
    { value: "nl_XX", label: "Dutch" },
    { value: "id_XX", label: "Indonesian" },
    { value: "tr_XX", label: "Turkish" },
    { value: "sv_XX", label: "Swedish" },
    { value: "bn_XX", label: "Bengali" },
    { value: "th_XX", label: "Thai" },
    { value: "da_XX", label: "Danish" },

    // Add more languages as needed
  ];

  const handleLanguageChange = (selectedOption) => {
    setSelectedLanguage(selectedOption);
  };

  return (
    <div className="flex">
      <Select
        value={selectedLanguage}
        onChange={handleLanguageChange}
        options={languageOptions}
        placeholder="Select language"
        className="pt-1.5"
      />
      <button
        onClick={() => translateContent(editor, selectedLanguage)}
        style={{
          padding: "5px 10px",
          fontSize: "16px",
          backgroundColor: "blue",
          color: "#ffffff",
          border: "1px solid black",
          borderRadius: "4px",
          cursor: "pointer",
          marginLeft: "10px", // Adjust spacing
        }}
      >
        {loading ? <Loader /> : "Translate"}
      </button>
    </div>
  );
};

const MenuBar = ({
  editor,
  summarizeContent,
  translateContent,
  loadingSummarize,
  loadingTranslate,
}) => {
  const [menuBarState, setMenuBarState] = useState({
    isBoldActive: false,
    isItalicActive: false,
    // Add more states for other menu items as needed
  });

  if (!editor) {
    return null;
  }
  const handleDownloadTxt = () => {
    const content = editor.getHTML();

    html2pdf(content, {
      margin: 10,
      filename: "editor_content.pdf",
      html2canvas: { scale: 2 },
      jsPDF: { unit: "mm", format: "a4", orientation: "portrait" },
    });
  };

  const handleToggleBold = () => {
    editor.chain().focus().toggleBold().run();
    setMenuBarState((prevState) => ({
      ...prevState,
      isBoldActive: !prevState.isBoldActive,
    }));
  };

  const handleToggleItalic = () => {
    editor.chain().focus().toggleItalic().run();
    setMenuBarState((prevState) => ({
      ...prevState,
      isItalicActive: !prevState.isItalicActive,
    }));
  };

  // Add similar functions for other menu items

  return (
    <div className="menuBar">
      <div>
        <button
          onClick={handleToggleBold}
          className={menuBarState.isBoldActive ? "is_active" : ""}
        >
          <FaBold />
        </button>
        <button
          onClick={handleToggleItalic}
          className={menuBarState.isItalicActive ? "is_active" : ""}
        >
          <FaItalic />
        </button>
        <button
          onClick={() => editor.chain().focus().toggleUnderline().run()}
          className={editor.isActive("underline") ? "is_active" : ""}
        >
          <FaUnderline />
        </button>
        <button
          onClick={() => editor.chain().focus().toggleStrike().run()}
          className={editor.isActive("strike") ? "is_active" : ""}
        >
          <FaStrikethrough />
        </button>
        <button
          onClick={() => editor.chain().focus().toggleBulletList().run()}
          className={editor.isActive("bulletList") ? "is_active" : ""}
        >
          <FaListUl />
        </button>
        <button
          onClick={() => editor.chain().focus().toggleOrderedList().run()}
          className={editor.isActive("orderedList") ? "is_active" : ""}
        >
          <FaListOl />
        </button>

        {/* Summarization Button */}
      </div>
      <div>
        <SummarizeButton
          editor={editor}
          summarizeContent={summarizeContent}
          loading={loadingSummarize}
        />
      </div>
      <div>
        <TranslateButton
          editor={editor}
          translateContent={translateContent}
          loading={loadingTranslate}
        />
      </div>
      <div>
        <button onClick={() => editor.chain().focus().undo().run()}>
          <FaUndo />
        </button>
        <button onClick={() => editor.chain().focus().redo().run()}>
          <FaRedo />
        </button>
        <button
          onClick={handleDownloadTxt}
          style={{
            padding: "5px 10px",
            fontSize: "16px",
            backgroundColor: "Black",
            color: "rgb(255,255,255)",
            border: "1px solid black",
            borderRadius: "4px",
            cursor: "pointer",
            marginRight: "10px", // Adjust spacing
          }}
        >
          Download
        </button>
      </div>
    </div>
  );
};

export const Tiptap = ({
  setDescription,
  summarizeContent,
  translateContent,
}) => {
  const [states, setStates] = useState("");

  const editor = useEditor({
    extensions: [StarterKit, Underline, FontSize, PlainTextLogger()],
    content: ``,
    onUpdate: ({ editor }) => {
      const html = editor.getHTML();
      setDescription(html);
    },
  });

  return (
    <div className="textEditor">
      <MenuBar
        editor={editor}
        summarizeContent={summarizeContent}
        translateContent={translateContent}
      />
      <EditorContent editor={editor} />
    </div>
  );
};

// Function to strip HTML tags from a string

const stripHtmlTags = (htmlString) => {
  const doc = new DOMParser().parseFromString(htmlString, "text/html");
  return doc.body.textContent || "";
};

const TipTapWithLogger = () => {
  const [description, setDescription] = useState("");
  const [loadingSummarize, setLoadingSummarize] = useState(false);
  const [loadingTranslate, setLoadingTranslate] = useState(false);

  const summarizeContent = async (editor) => {
    setLoadingSummarize(true);
    // Get the selection from the editor
    const focusedText = editor.commands.setTextSelection({
      from: editor.state.selection.from,
      to: editor.state.selection.to,
    });
    console.log(focusedText);
    const ftext = stripHtmlTags(
      editor.getHTML({
        from: editor.state.selection.from,
        to: editor.state.selection.to,
      })
    );
    console.log("ftext ", ftext);
    try {
      const apiUrl = "https://nekonotes-7gmw.onrender.com/summarize";

      const response = await fetch(apiUrl, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ inputs: ftext }),
      });

      if (response.ok) {
        const data = await response.json();
        console.log("API Response:", data);
        console.log("Focus", focusedText);
        // Assuming the API response includes a 'summary_text' field
        const summaryText = data[0].summary_text;

        // Set the editor content to the summary text
        editor.commands.setContent(summaryText);

        // Update the editor description with the summary text
        setDescription(summaryText);
      } else {
        console.error("API Error:", response.statusText);
      }
    } catch (error) {
      console.error("API Request Error:", error.message);
    }
    setLoadingSummarize(false);
  };

  const translateContent = async (editor, selectedLanguage) => {
    setLoadingTranslate(true);
    const plainTextContent = stripHtmlTags(editor.getHTML());
    console.log("Plain_text =", plainTextContent);
    try {
      const apiUrl = "https://nekonotes-7gmw.onrender.com/translate";

      const response = await fetch(apiUrl, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          inputs: plainTextContent,
          parameters: {
            src_lang: "en_XX",
            tgt_lang: selectedLanguage.value,
          },
          // Use the selected language value
        }),
      });

      if (response.ok) {
        const data = await response.json();
        console.log("Translate API Response:", data);

        // Assuming the API response includes a 'translatedText' field
        const translatedText = data[0].translation_text;

        // Set the editor content to the translated text
        editor.commands.setContent(translatedText);

        // Update the editor description with the translated text
        setDescription(translatedText);
      } else {
        console.error("Translate API Error:", response.statusText);
        setDescription(response.statusText);
      }
    } catch (error) {
      console.error("Translate API Request Error:", error.message);
    }
    setLoadingTranslate(false);
  };

  const postNoteApi = async (plainTextContent) => {
    // Send API request to postNoteApi with plain text as payload
    try {
      const apiUrl =
        "https://neko-notesbackendstorage.onrender.com/notebooks/657dfa6a9905d9766bfff39a/notes/abhiram/text";
      const response = await fetch(apiUrl, {
        method: "PATCH",
        headers: {
          "Content-Type": "application/json", // Assuming you are sending JSON data
          // Add other headers as needed
        },
        body: JSON.stringify({
          text: plainTextContent,
        }),
      });

      if (response.ok) {
        // Handle successful API response if needed
      } else {
        console.error("API Error:", response.statusText);
      }
    } catch (error) {
      console.error("API Request Error:", error.message);
    }
  };

  useEffect(() => {
    console.log(stripHtmlTags(description));
    postNoteApi(stripHtmlTags(description));
    console.log("successful");
  }, [description]);

  return (
    <div className="App">
      <Tiptap
        setDescription={setDescription}
        summarizeContent={summarizeContent}
        translateContent={translateContent}
      />
    </div>
  );
};

export default TipTapWithLogger;
