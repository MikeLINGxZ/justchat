import React from "react";
import { useTranslation } from 'react-i18next';

interface AppsPageProps {
    className?: string;
}

const AppsPage: React.FC<AppsPageProps> = ({className}) => {
    const { t } = useTranslation();

    return (
        <div>
            {t('apps.placeholder')}
        </div>
    );
}

export default AppsPage;
