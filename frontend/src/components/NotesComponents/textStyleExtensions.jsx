// textStyleExtensions.js
import { TextStyle } from '@tiptap/extension-text-style';

const FontSize = TextStyle.extend({
    addAttributes() {
        return {
            fontSize: {
                default: null,
                parseHTML: element => element.style.fontSize.replace('px', ''),
                renderHTML: attributes => {
                    if (attributes.fontSize) {
                        return { style: `font-size: ${attributes.fontSize}px` };
                    }
                    return {};
                },
            },
        };
    },
    addCommands() {
        return {
            setFontSize: fontSize => ({ commands }) => {
                return commands.updateAttributes('textStyle', { fontSize });
            },
        };
    },
    addKeyboardShortcuts() {
        return {
            'Mod-Alt-=': ({ commands }) => commands.setFontSize(2),
            'Mod-Alt--': ({ commands }) => commands.setFontSize(-2),
        };
    },
});

export { FontSize };
