import React from 'react';
import { Result, Button } from 'antd';
import { useNavigate } from 'react-router-dom';
import { HomeOutlined, ArrowLeftOutlined } from '@ant-design/icons';
import { useTranslation } from 'react-i18next';

const NotFound: React.FC = () => {
  const { t } = useTranslation();
  const navigate = useNavigate();

  React.useEffect(() => {
    document.title = t('app.notFoundTitle');
  }, [t]);

  const handleGoHome = () => {
    navigate('/home');
  };

  const handleGoBack = () => {
    navigate(-1);
  };

  return (
    <div style={{ 
      display: 'flex', 
      justifyContent: 'center', 
      alignItems: 'center', 
      minHeight: '100vh',
      padding: '24px'
    }}>
      <Result
        status="404"
        title="404"
        subTitle={t('notFound.subtitle')}
        extra={
          <div style={{ display: 'flex', gap: '12px', justifyContent: 'center' }}>
            <Button type="primary" icon={<HomeOutlined />} onClick={handleGoHome}>
              {t('notFound.home')}
            </Button>
            <Button icon={<ArrowLeftOutlined />} onClick={handleGoBack}>
              {t('notFound.back')}
            </Button>
          </div>
        }
      />
    </div>
  );
};

export default NotFound;
