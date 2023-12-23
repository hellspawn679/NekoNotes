import React, { useCallback } from 'react';
import { useDropzone } from 'react-dropzone';
import * as pdfjsLib from 'pdfjs-dist';

const FileUploadButton = () => {
    const onDrop = useCallback(async (acceptedFiles) => {
        // Iterate through each uploaded file
        for (const file of acceptedFiles) {
            const reader = new FileReader();

            reader.onload = async (e) => {
                const content = e.target.result;

                if (file.name.endsWith('.pdf')) {
                    // Convert PDF to text using pdfjs-dist
                    const pdfData = new Uint8Array(content);
                    const pdfText = await getPdfText(pdfData);
                    console.log(`Text from ${file.name}: `, pdfText);
                } else if (file.name.endsWith('.txt')) {
                    // For TXT files, use plain text
                    console.log(`Text from ${file.name}: `, content);
                } else {
                    // Unsupported file type
                    console.log(`Unsupported file type: ${file.name}`);
                }
            };

            reader.readAsArrayBuffer(file);
        }
    }, []);

    const { getRootProps, getInputProps, isDragActive } = useDropzone({
        onDrop,
        accept: ['.pdf', '.txt'],
    });

    // Helper function to convert PDF to text
    const getPdfText = async (pdfData) => {
        const loadingTask = pdfjsLib.getDocument({ data: pdfData });
        const pdfDocument = await loadingTask.promise;

        const textContent = [];
        for (let pageNum = 1; pageNum <= pdfDocument.numPages; pageNum++) {
            const page = await pdfDocument.getPage(pageNum);
            const pageText = await page.getTextContent();
            textContent.push(pageText.items.map((item) => item.str).join(' '));
        }

        return textContent.join('\n');
    };

    return (
        <div {...getRootProps()} className={`dropzone ${isDragActive ? 'active' : ''}`}>
            <input {...getInputProps()} />
            <p>Upload +</p>
        </div>
    );
};

export default FileUploadButton;
