import React from 'react';
import { Card, Radio, Space, Typography, Slider, Row, Col, Button, Divider } from 'antd';
import { ReloadOutlined, CheckOutlined } from '@ant-design/icons';
import { useFontSizeStore, FONT_SIZE_OPTIONS, FONT_SIZE_OFFSETS, getFontSizeLabel } from '@/stores/fontSizeStore';
import styles from './index.module.scss';

const { Title, Text, Paragraph } = Typography;

const GeneralSettingsPage: React.FC = () => {
  const { fontSizeOffset, setFontSizeOffset, resetFontSize } = useFontSizeStore();

  // 处理字体大小变更
  const handleFontSizeChange = (value: number) => {
    setFontSizeOffset(value as any);
  };

  // 滑块标记
  const sliderMarks = FONT_SIZE_OPTIONS.reduce((marks, option) => {
    marks[option.value] = {
      style: { fontSize: '12px', color: '#666' },
      label: option.label,
    };
    return marks;
  }, {} as any);

  return (
    <div className={styles.generalSettings}>
      <div className={styles.header}>
        <Title level={3}>通用设置</Title>
        <Text type="secondary">个性化你的使用体验</Text>
      </div>

      <div className={styles.content}>
        {/* 字体设置区域 */}
        <Card title="显示设置" className={styles.settingCard}>
          <div className={styles.settingItem}>
            <div className={styles.settingLabel}>
              <Title level={5}>字体大小</Title>
              <Text type="secondary">调整应用中的文字大小，提升阅读体验</Text>
            </div>
            
            <div className={styles.fontSizeControl}>
              {/* 滑块控制 */}
              <div className={styles.sliderContainer}>
                <Row gutter={[16, 16]} align="middle">
                  <Col span={4}>
                    <Text strong>字体大小</Text>
                  </Col>
                  <Col span={16}>
                    <Slider
                      min={FONT_SIZE_OFFSETS.VERY_SMALL}
                      max={FONT_SIZE_OFFSETS.EXTRA_LARGE}
                      step={2}
                      value={fontSizeOffset}
                      onChange={handleFontSizeChange}
                      marks={sliderMarks}
                      className={styles.fontSizeSlider}
                    />
                  </Col>
                  <Col span={4}>
                    <div className={styles.currentSize}>
                      <Text strong>{getFontSizeLabel(fontSizeOffset)}</Text>
                      <Text type="secondary" className={styles.sizeDesc}>
                        {14 + fontSizeOffset}px
                      </Text>
                    </div>
                  </Col>
                </Row>
              </div>

              {/* 预设选项 */}
              <div className={styles.presetOptions}>
                <Text strong>快速选择：</Text>
                <Radio.Group 
                  value={fontSizeOffset} 
                  onChange={(e) => handleFontSizeChange(e.target.value)}
                  className={styles.fontSizeRadio}
                >
                  <Space wrap>
                    {FONT_SIZE_OPTIONS.map((option) => (
                      <Radio.Button key={option.value} value={option.value}>
                        <span className={styles.radioLabel}>
                          {option.label}
                          <small>({option.description})</small>
                        </span>
                      </Radio.Button>
                    ))}
                  </Space>
                </Radio.Group>
              </div>

              {/* 预览区域 */}
              <div className={styles.previewArea}>
                <Title level={5}>预览效果</Title>
                <div className={styles.previewContent}>
                  <Paragraph>
                    这是标准字体大小的文本预览。你可以通过调整上方的设置来改变文字的大小，
                    找到最适合你阅读习惯的字体尺寸。
                  </Paragraph>
                  <Text>小号文字：这是较小的辅助信息文本。</Text>
                  <br />
                  <Text strong>粗体文字：这是重要的加粗文本。</Text>
                  <br />
                  <Text type="secondary">次要文字：这是次要信息文本。</Text>
                </div>
              </div>

              {/* 操作按钮 */}
              <div className={styles.actions}>
                <Button 
                  icon={<ReloadOutlined />}
                  onClick={resetFontSize}
                  disabled={fontSizeOffset === FONT_SIZE_OFFSETS.NORMAL}
                >
                  恢复默认
                </Button>
                <Button 
                  type="primary" 
                  icon={<CheckOutlined />}
                  className={styles.applyButton}
                >
                  应用设置
                </Button>
              </div>
            </div>
          </div>
        </Card>

        {/* 其他设置区域预留 */}
        <Card title="其他设置" className={styles.settingCard}>
          <div className={styles.placeholderContent}>
            <Text type="secondary">更多个性化设置功能正在开发中...</Text>
          </div>
        </Card>
      </div>
    </div>
  );
};

export default GeneralSettingsPage;