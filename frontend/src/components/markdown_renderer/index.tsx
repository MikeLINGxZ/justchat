import React from "react";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import { tomorrow } from "react-syntax-highlighter/dist/esm/styles/prism";
import styles from "./index.module.scss";

export type MarkdownVariant = "assistant" | "user" | "reasoning" | "trace";

interface MarkdownRendererProps {
    content?: string | null;
    variant?: MarkdownVariant;
    className?: string;
    transformContent?: (content: string) => string;
    decorateText?: (children: React.ReactNode) => React.ReactNode;
}

function decorateChildren(
    children: React.ReactNode,
    decorateText?: (children: React.ReactNode) => React.ReactNode,
): React.ReactNode {
    if (!decorateText) {
        return children;
    }
    return React.Children.map(children, (child) => {
        if (typeof child === "string") {
            return decorateText(child);
        }
        if (React.isValidElement(child) && child.props?.children) {
            return React.cloneElement(child as React.ReactElement<any>, {
                ...child.props,
                children: decorateChildren(child.props.children, decorateText),
            });
        }
        return child;
    });
}

function renderTextContainer(
    Tag: keyof React.JSX.IntrinsicElements,
    className: string,
    children: React.ReactNode,
    decorateText?: (children: React.ReactNode) => React.ReactNode,
) {
    return React.createElement(Tag, { className }, decorateChildren(children, decorateText));
}

const MarkdownRenderer: React.FC<MarkdownRendererProps> = ({
    content,
    variant = "assistant",
    className = "",
    transformContent,
    decorateText,
}) => {
    const rawContent = content?.trim() ?? "";
    if (!rawContent) {
        return null;
    }

    const displayContent = transformContent ? transformContent(rawContent) : rawContent;

    return (
        <div className={[styles.markdownRoot, styles[variant], className].filter(Boolean).join(" ")}>
            <ReactMarkdown
                remarkPlugins={[remarkGfm]}
                components={{
                    code(props: any) {
                        const { inline, className: nodeClassName, children, ...rest } = props;
                        const match = /language-(\w+)/.exec(nodeClassName || "");
                        const language = match ? match[1] : "";

                        return !inline && language ? (
                            <SyntaxHighlighter
                                style={tomorrow}
                                language={language}
                                PreTag="div"
                                customStyle={{
                                    margin: "10px 0",
                                    borderRadius: "10px",
                                    fontSize: "13px",
                                } as any}
                                {...rest}
                            >
                                {String(children).replace(/\n$/, "")}
                            </SyntaxHighlighter>
                        ) : (
                            <code className={styles.inlineCode} {...rest}>
                                {children}
                            </code>
                        );
                    },
                    table: ({ children }) => (
                        <div className={styles.tableWrapper}>
                            <table className={styles.markdownTable}>{children}</table>
                        </div>
                    ),
                    a: ({ children, href }) => (
                        <a
                            href={href}
                            target="_blank"
                            rel="noopener noreferrer"
                            className={styles.markdownLink}
                        >
                            {children}
                        </a>
                    ),
                    blockquote: ({ children }) => (
                        <blockquote className={styles.markdownBlockquote}>
                            {decorateChildren(children, decorateText)}
                        </blockquote>
                    ),
                    ul: ({ children }) => <ul className={styles.markdownList}>{children}</ul>,
                    ol: ({ children }) => <ol className={styles.markdownList}>{children}</ol>,
                    li: ({ children, ...props }) => (
                        <li
                            className={props.className?.includes("task-list-item")
                                ? `${styles.markdownListItem} ${styles.taskListItem}`
                                : styles.markdownListItem}
                        >
                            {decorateChildren(children, decorateText)}
                        </li>
                    ),
                    p: ({ children }) => renderTextContainer("p", styles.paragraph, children, decorateText),
                    h1: ({ children }) => renderTextContainer("h1", `${styles.heading} ${styles.h1}`, children, decorateText),
                    h2: ({ children }) => renderTextContainer("h2", `${styles.heading} ${styles.h2}`, children, decorateText),
                    h3: ({ children }) => renderTextContainer("h3", `${styles.heading} ${styles.h3}`, children, decorateText),
                    h4: ({ children }) => renderTextContainer("h4", `${styles.heading} ${styles.h4}`, children, decorateText),
                    h5: ({ children }) => renderTextContainer("h5", `${styles.heading} ${styles.h5}`, children, decorateText),
                    h6: ({ children }) => renderTextContainer("h6", `${styles.heading} ${styles.h6}`, children, decorateText),
                    strong: ({ children }) => <strong className={styles.strong}>{children}</strong>,
                    em: ({ children }) => <em className={styles.emphasis}>{children}</em>,
                    hr: () => <hr className={styles.divider} />,
                }}
            >
                {displayContent}
            </ReactMarkdown>
        </div>
    );
};

export default MarkdownRenderer;
