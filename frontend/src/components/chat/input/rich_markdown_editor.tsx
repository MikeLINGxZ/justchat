import React, { useEffect, useMemo, useRef } from "react";
import { useTranslation } from 'react-i18next';
import DOMPurify from "dompurify";
import { marked } from "marked";
import TurndownService from "turndown";
import { EditorContent, useEditor } from "@tiptap/react";
import type { EditorView } from "@tiptap/pm/view";
import StarterKit from "@tiptap/starter-kit";
import Placeholder from "@tiptap/extension-placeholder";
import styles from "./rich_markdown_editor.module.scss";

interface RichMarkdownEditorProps {
    value?: string;
    placeholder?: string;
    onChange: (markdown: string) => void;
    onSend: () => void;
}

const turndownService = new TurndownService({
    bulletListMarker: "-",
    codeBlockStyle: "fenced",
    emDelimiter: "*",
    headingStyle: "atx",
});

turndownService.addRule("strikethrough", {
    filter: ["del", "s", "strike"],
    replacement(content) {
        return `~~${content}~~`;
    },
});

function markdownToHtml(markdown: string): string {
    const parsed = marked.parse(markdown, {
        breaks: true,
        gfm: true,
    }) as string;
    return DOMPurify.sanitize(parsed);
}

function htmlToMarkdown(html: string): string {
    return turndownService
        .turndown(html)
        .replace(/\n{3,}/g, "\n\n")
        .trim();
}

function isComposingKeyboardEvent(event: KeyboardEvent): boolean {
    return event.isComposing || event.keyCode === 229;
}

function insertHardBreak(view: EditorView): boolean {
    const hardBreak = view.state.schema.nodes.hardBreak;
    if (!hardBreak) {
        return false;
    }

    const transaction = view.state.tr.replaceSelectionWith(hardBreak.create()).scrollIntoView();
    view.dispatch(transaction);
    return true;
}

const RichMarkdownEditor: React.FC<RichMarkdownEditorProps> = ({
    value = "",
    placeholder,
    onChange,
    onSend,
}) => {
    const { t } = useTranslation();
    const resolvedPlaceholder = placeholder ?? t('chat.input.editorPlaceholder');
    const initialContent = useMemo(() => markdownToHtml(value), [value]);
    const lastMarkdownRef = useRef(value);

    const editor = useEditor({
        immediatelyRender: false,
        extensions: [
            StarterKit.configure({
                codeBlock: true,
            }),
            Placeholder.configure({
                placeholder: resolvedPlaceholder,
            }),
        ],
        content: initialContent,
        editorProps: {
            attributes: {
                class: "markdown-editor-content",
            },
            handleKeyDown(view, event) {
                if (isComposingKeyboardEvent(event)) {
                    return false;
                }

                if (event.key !== "Enter") {
                    return false;
                }

                const parentNodeName = view.state.selection.$from.parent.type.name;
                const inCodeBlock = parentNodeName === "codeBlock";
                const inListItem = parentNodeName === "listItem";

                if (event.shiftKey) {
                    if (inCodeBlock) {
                        return false;
                    }

                    event.preventDefault();
                    return insertHardBreak(view);
                }

                if (!inCodeBlock && !inListItem) {
                    event.preventDefault();
                    onSend();
                    return true;
                }

                return false;
            },
        },
        onUpdate({ editor: currentEditor }) {
            const nextMarkdown = htmlToMarkdown(currentEditor.getHTML());
            lastMarkdownRef.current = nextMarkdown;
            onChange(nextMarkdown);
        },
    });

    useEffect(() => {
        if (!editor) {
            return;
        }
        if (value === lastMarkdownRef.current) {
            return;
        }
        editor.commands.setContent(value ? markdownToHtml(value) : "", false);
        lastMarkdownRef.current = value;
    }, [editor, value]);

    return (
        <div className={styles.editorShell}>
            <EditorContent editor={editor} className={styles.editorContent} />
        </div>
    );
};

export default RichMarkdownEditor;
