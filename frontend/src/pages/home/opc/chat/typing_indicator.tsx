import React from 'react';
import { useTranslation } from 'react-i18next';

interface TypingIndicatorProps {
    names: string[];
}

const TypingIndicator: React.FC<TypingIndicatorProps> = ({ names }) => {
    const { t } = useTranslation();

    if (names.length === 0) return null;

    const text = names.length === 1
        ? t('opc.chat.typing', { name: names[0] })
        : t('opc.chat.typingMultiple', { names: names.join('、') });

    return (
        <div className="typing-indicator">
            <div className="typing-dots">
                <span className="dot"></span>
                <span className="dot"></span>
                <span className="dot"></span>
            </div>
            <span className="typing-text">{text}</span>
        </div>
    );
};

export default TypingIndicator;
