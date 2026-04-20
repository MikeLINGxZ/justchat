import React from 'react';
import { useTranslation } from 'react-i18next';
import {
    MessageOutlined,
    FileTextOutlined,
    CodeOutlined,
    BulbOutlined,
} from '@ant-design/icons';
import styles from './index.module.scss';

interface WelcomeEmptyProps {
    onSuggestionClick?: (prompt: string) => void;
}

const WelcomeEmpty: React.FC<WelcomeEmptyProps> = ({ onSuggestionClick }) => {
    const { t } = useTranslation();

    const suggestions = [
        {
            key: 'newChat',
            icon: <MessageOutlined />,
            iconClass: styles.iconTeal,
            title: t('home.welcome.suggestions.newChat.title'),
            desc: t('home.welcome.suggestions.newChat.desc'),
            prompt: t('home.welcome.suggestions.newChat.prompt'),
        },
        {
            key: 'summarize',
            icon: <FileTextOutlined />,
            iconClass: styles.iconCyan,
            title: t('home.welcome.suggestions.summarize.title'),
            desc: t('home.welcome.suggestions.summarize.desc'),
            prompt: t('home.welcome.suggestions.summarize.prompt'),
        },
        {
            key: 'analyzeCode',
            icon: <CodeOutlined />,
            iconClass: styles.iconPurple,
            title: t('home.welcome.suggestions.analyzeCode.title'),
            desc: t('home.welcome.suggestions.analyzeCode.desc'),
            prompt: t('home.welcome.suggestions.analyzeCode.prompt'),
        },
        {
            key: 'brainstorm',
            icon: <BulbOutlined />,
            iconClass: styles.iconAmber,
            title: t('home.welcome.suggestions.brainstorm.title'),
            desc: t('home.welcome.suggestions.brainstorm.desc'),
            prompt: t('home.welcome.suggestions.brainstorm.prompt'),
        },
    ];

    return (
        <div className={styles.welcomeContainer}>
            <div className={styles.welcomeContent}>
                <div className={styles.welcomeHeader}>
                    <div className={styles.welcomeIcon}>🍋</div>
                    <h1 className={styles.welcomeTitle}>
                        {t('home.welcome.title')}
                    </h1>
                    <p className={styles.welcomeSubtitle}>
                        {t('home.welcome.subtitle')}
                    </p>
                </div>

                <div className={styles.suggestionsGrid}>
                    {suggestions.map((item) => (
                        <button
                            type="button"
                            key={item.key}
                            className={styles.suggestionCard}
                            onClick={() => onSuggestionClick?.(item.prompt)}
                        >
                            <div className={`${styles.suggestionIcon} ${item.iconClass}`}>
                                {item.icon}
                            </div>
                            <div className={styles.suggestionBody}>
                                <h3 className={styles.suggestionTitle}>{item.title}</h3>
                                <p className={styles.suggestionDesc}>{item.desc}</p>
                            </div>
                        </button>
                    ))}
                </div>
            </div>
        </div>
    );
};

export default WelcomeEmpty;
