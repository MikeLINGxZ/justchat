import React, { useState } from 'react';
import { Button, Card, Typography, Space, Alert, Divider } from 'antd';
import { handleGlobalError, isLoginExpiredError } from '../utils/globalErrorHandler';
import { chatClient } from '../api/chatClient';

const { Title, Text } = Typography;

const TestLoginExpiry: React.FC = () => {
  const [testResult, setTestResult] = useState<string>('');
  const [isLoading, setIsLoading] = useState(false);

  // 模拟登录过期错误
  const simulateLoginExpiredError = () => {
    const mockError = {
      code: 2,
      message: 'ErrCodeLoginExpired',
      details: []
    };
    
    console.log('模拟登录过期错误:', mockError);
    const handled = handleGlobalError(mockError);
    setTestResult(`登录过期错误处理结果: ${handled ? '已处理' : '未处理'}`);
  };

  // 测试错误检测函数
  const testErrorDetection = () => {
    const testCases = [
      { code: 2, message: 'ErrCodeLoginExpired', details: [] },
      { error: { message: 'ErrCodeLoginExpired' } },
      { data: { message: 'ErrCodeLoginExpired' } },
      { message: 'ErrCodeLoginExpired' },
      { code: 1, message: 'OtherError', details: [] }
    ];
    
    const results = testCases.map((testCase, index) => {
      const isExpired = isLoginExpiredError(testCase);
      return `测试用例 ${index + 1}: ${isExpired ? '✓ 检测到登录过期' : '✗ 未检测到登录过期'}`;
    });
    
    setTestResult(results.join('\n'));
  };

  // 测试实际API调用（如果后端返回登录过期错误）
  const testRealApiCall = async () => {
    setIsLoading(true);
    try {
      // 尝试调用需要认证的API
      await chatClient.listChats({
        limit: '10',
        offset: '0'
      });
      setTestResult('API调用成功，未触发登录过期错误');
    } catch (error) {
      console.log('API调用错误:', error);
      const isExpired = isLoginExpiredError(error);
      if (isExpired) {
        setTestResult('✓ 检测到登录过期错误，全局处理器应该已经处理');
      } else {
        setTestResult('✗ API调用失败，但不是登录过期错误: ' + JSON.stringify(error, null, 2));
      }
    } finally {
      setIsLoading(false);
    }
  };

  // 清除当前token（模拟登录过期）
  const clearToken = () => {
    localStorage.removeItem('token');
    setTestResult('已清除token，下次API调用可能会触发登录过期错误');
  };

  // 设置无效token
  const setInvalidToken = () => {
    localStorage.setItem('token', 'invalid_token_12345');
    setTestResult('已设置无效token，下次API调用可能会触发登录过期错误');
  };

  return (
    <div style={{ padding: '24px', maxWidth: '800px', margin: '0 auto' }}>
      <Card>
        <Title level={2}>登录过期处理测试</Title>
        <Text type="secondary">
          此页面用于测试全局登录过期错误处理机制是否正常工作。
        </Text>
        
        <Divider />
        
        <Space direction="vertical" size="large" style={{ width: '100%' }}>
          <div>
            <Title level={4}>1. 错误检测测试</Title>
            <Text>测试 isLoginExpiredError 函数是否能正确识别登录过期错误</Text>
            <br />
            <Button onClick={testErrorDetection} style={{ marginTop: '8px' }}>
              运行错误检测测试
            </Button>
          </div>
          
          <div>
            <Title level={4}>2. 模拟登录过期</Title>
            <Text>直接调用全局错误处理器，模拟登录过期情况</Text>
            <br />
            <Button onClick={simulateLoginExpiredError} style={{ marginTop: '8px' }}>
              模拟登录过期错误
            </Button>
          </div>
          
          <div>
            <Title level={4}>3. Token管理</Title>
            <Text>管理本地存储的认证token</Text>
            <br />
            <Space style={{ marginTop: '8px' }}>
              <Button onClick={clearToken}>清除Token</Button>
              <Button onClick={setInvalidToken}>设置无效Token</Button>
            </Space>
          </div>
          
          <div>
            <Title level={4}>4. 实际API测试</Title>
            <Text>调用真实API接口，测试登录过期处理</Text>
            <br />
            <Button 
              onClick={testRealApiCall} 
              loading={isLoading}
              style={{ marginTop: '8px' }}
            >
              测试API调用
            </Button>
          </div>
          
          {testResult && (
            <Alert
              message="测试结果"
              description={<pre style={{ whiteSpace: 'pre-wrap' }}>{testResult}</pre>}
              type="info"
              showIcon
            />
          )}
        </Space>
      </Card>
    </div>
  );
};

export default TestLoginExpiry;