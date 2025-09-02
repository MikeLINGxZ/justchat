import React from 'react';
import { Card, Form, Input, Button, Avatar, Upload, message, Space, Divider } from 'antd';
import { UserOutlined, UploadOutlined, EditOutlined } from '@ant-design/icons';
import type { UploadProps } from 'antd';
import { useAuthStore } from '@/stores';

const Profile: React.FC = () => {
  const { user } = useAuthStore();
  const [form] = Form.useForm();
  const [editing, setEditing] = React.useState(false);

  React.useEffect(() => {
    // 设置页面标题
    document.title = '个人资料 - Lemon Tea';
  }, []);

  React.useEffect(() => {
    if (user) {
      form.setFieldsValue({
        username: user.username,
        email: user.email,
        phone: user.phone || '',
        bio: user.bio || '',
      });
    }
  }, [user, form]);

  const handleSubmit = async (values: any) => {
    try {
      // 这里应该调用更新用户信息的API
      console.log('更新用户信息:', values);
      message.success('个人资料更新成功！');
      setEditing(false);
    } catch (error) {
      message.error('更新失败，请重试');
    }
  };

  const uploadProps: UploadProps = {
    name: 'avatar',
    action: '/api/upload/avatar',
    headers: {
      authorization: `Bearer ${localStorage.getItem('token')}`,
    },
    beforeUpload: (file) => {
      const isJpgOrPng = file.type === 'image/jpeg' || file.type === 'image/png';
      if (!isJpgOrPng) {
        message.error('只能上传 JPG/PNG 格式的图片!');
      }
      const isLt2M = file.size / 1024 / 1024 < 2;
      if (!isLt2M) {
        message.error('图片大小不能超过 2MB!');
      }
      return isJpgOrPng && isLt2M;
    },
    onChange: (info) => {
      if (info.file.status === 'done') {
        message.success('头像上传成功!');
      } else if (info.file.status === 'error') {
        message.error('头像上传失败!');
      }
    },
  };

  return (
    <div style={{ padding: '24px', maxWidth: '800px', margin: '0 auto' }}>
      <Card>
        <Space direction="vertical" size="large" style={{ width: '100%' }}>
          {/* 头像部分 */}
          <div style={{ textAlign: 'center' }}>
            <Avatar
              size={120}
              src={user?.avatar}
              icon={<UserOutlined />}
              style={{ marginBottom: '16px' }}
            />
            <div>
              <Upload {...uploadProps} showUploadList={false}>
                <Button icon={<UploadOutlined />}>更换头像</Button>
              </Upload>
            </div>
          </div>

          <Divider />

          {/* 基本信息 */}
          <div>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '16px' }}>
              <h3>基本信息</h3>
              <Button
                type={editing ? 'default' : 'primary'}
                icon={<EditOutlined />}
                onClick={() => setEditing(!editing)}
              >
                {editing ? '取消编辑' : '编辑资料'}
              </Button>
            </div>

            <Form
              form={form}
              layout="vertical"
              onFinish={handleSubmit}
              disabled={!editing}
            >
              <Form.Item
                label="用户名"
                name="username"
                rules={[
                  { required: true, message: '请输入用户名' },
                  { min: 3, message: '用户名至少3个字符' },
                ]}
              >
                <Input placeholder="请输入用户名" />
              </Form.Item>

              <Form.Item
                label="邮箱"
                name="email"
                rules={[
                  { required: true, message: '请输入邮箱' },
                  { type: 'email', message: '请输入有效的邮箱地址' },
                ]}
              >
                <Input placeholder="请输入邮箱" />
              </Form.Item>

              <Form.Item
                label="手机号"
                name="phone"
                rules={[
                  { pattern: /^1[3-9]\d{9}$/, message: '请输入有效的手机号' },
                ]}
              >
                <Input placeholder="请输入手机号" />
              </Form.Item>

              <Form.Item
                label="个人简介"
                name="bio"
              >
                <Input.TextArea
                  rows={4}
                  placeholder="介绍一下自己吧..."
                  maxLength={200}
                  showCount
                />
              </Form.Item>

              {editing && (
                <Form.Item>
                  <Space>
                    <Button type="primary" htmlType="submit">
                      保存更改
                    </Button>
                    <Button onClick={() => setEditing(false)}>
                      取消
                    </Button>
                  </Space>
                </Form.Item>
              )}
            </Form>
          </div>

          <Divider />

          {/* 账户信息 */}
          <div>
            <h3>账户信息</h3>
            <Space direction="vertical" style={{ width: '100%' }}>
              <div>
                <strong>注册时间：</strong>
                {user?.createdAt ? new Date(user.createdAt).toLocaleDateString('zh-CN') : '未知'}
              </div>
              <div>
                <strong>最后更新：</strong>
                {user?.updatedAt ? new Date(user.updatedAt).toLocaleDateString('zh-CN') : '未知'}
              </div>
              <div>
                <strong>用户ID：</strong>
                {user?.id || '未知'}
              </div>
            </Space>
          </div>
        </Space>
      </Card>
    </div>
  );
};

export default Profile;