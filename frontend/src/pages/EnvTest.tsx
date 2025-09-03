import React from 'react';
import { Card, Descriptions, Tag, Typography } from 'antd';

const { Title } = Typography;

const EnvTest: React.FC = () => {
  return (
    <div style={{ padding: '24px', maxWidth: '800px', margin: '0 auto' }}>
      <Title level={2}>环境变量配置测试</Title>
      
      <Card title="环境信息" style={{ marginBottom: '16px' }}>
        <Descriptions column={1} bordered>
          <Descriptions.Item label="运行模式">
            <Tag color={import.meta.env.DEV ? 'orange' : 'purple'}>{import.meta.env.MODE}</Tag>
          </Descriptions.Item>
          <Descriptions.Item label="开发模式">
            <Tag color={import.meta.env.DEV ? 'success' : 'default'}>
              {import.meta.env.DEV ? '是' : '否'}
            </Tag>
          </Descriptions.Item>
          <Descriptions.Item label="生产模式">
            <Tag color={import.meta.env.PROD ? 'success' : 'default'}>
              {import.meta.env.PROD ? '是' : '否'}
            </Tag>
          </Descriptions.Item>
        </Descriptions>
      </Card>

      <Card title="项目信息">
        <Descriptions column={1} bordered>
          <Descriptions.Item label="项目名称">
            <code>Lemon Tea Desktop - UI Only</code>
          </Descriptions.Item>
          <Descriptions.Item label="状态">
            <Tag color="green">纯 UI 演示版本</Tag>
          </Descriptions.Item>
          <Descriptions.Item label="说明">
            <span>已移除所有 HTTP 请求功能，仅保留 UI 组件和模拟数据</span>
          </Descriptions.Item>
        </Descriptions>
      </Card>
    </div>
  );
};

export default EnvTest;