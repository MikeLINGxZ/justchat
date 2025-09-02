import React from 'react';
import { Card, Descriptions, Tag, Typography } from 'antd';
import { env, getApiConfig } from '@/config/env';

const { Title } = Typography;

const EnvTest: React.FC = () => {
  const apiConfig = getApiConfig();

  return (
    <div style={{ padding: '24px', maxWidth: '800px', margin: '0 auto' }}>
      <Title level={2}>环境变量配置测试</Title>
      
      <Card title="API配置" style={{ marginBottom: '16px' }}>
        <Descriptions column={1} bordered>
          <Descriptions.Item label="API基础URL">
            <Tag color="blue">{apiConfig.baseUrl}</Tag>
          </Descriptions.Item>
          <Descriptions.Item label="请求超时时间">
            <Tag color="green">{apiConfig.timeout}ms</Tag>
          </Descriptions.Item>
          <Descriptions.Item label="启用凭证传递">
            <Tag color={apiConfig.withCredentials ? 'success' : 'warning'}>
              {apiConfig.withCredentials ? '是' : '否'}
            </Tag>
          </Descriptions.Item>
        </Descriptions>
      </Card>

      <Card title="环境信息" style={{ marginBottom: '16px' }}>
        <Descriptions column={1} bordered>
          <Descriptions.Item label="运行模式">
            <Tag color={env.isDev ? 'orange' : 'purple'}>{env.mode}</Tag>
          </Descriptions.Item>
          <Descriptions.Item label="开发模式">
            <Tag color={env.isDev ? 'success' : 'default'}>
              {env.isDev ? '是' : '否'}
            </Tag>
          </Descriptions.Item>
          <Descriptions.Item label="生产模式">
            <Tag color={env.isProd ? 'success' : 'default'}>
              {env.isProd ? '是' : '否'}
            </Tag>
          </Descriptions.Item>
        </Descriptions>
      </Card>

      <Card title="环境变量原始值">
        <Descriptions column={1} bordered>
          <Descriptions.Item label="VITE_API_BASE_URL">
            <code>{import.meta.env.VITE_API_BASE_URL || '未设置'}</code>
          </Descriptions.Item>
          <Descriptions.Item label="VITE_API_TIMEOUT">
            <code>{import.meta.env.VITE_API_TIMEOUT || '未设置'}</code>
          </Descriptions.Item>
          <Descriptions.Item label="VITE_API_WITH_CREDENTIALS">
            <code>{import.meta.env.VITE_API_WITH_CREDENTIALS || '未设置'}</code>
          </Descriptions.Item>
        </Descriptions>
      </Card>
    </div>
  );
};

export default EnvTest;