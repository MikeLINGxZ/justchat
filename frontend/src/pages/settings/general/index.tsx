import React, { useState } from 'react';
import { Card, Typography, Slider, Button, message } from 'antd';
import { ReloadOutlined, CheckOutlined, FontSizeOutlined } from '@ant-design/icons';
import { useFontSizeStore, FONT_SIZE_OPTIONS, FONT_SIZE_OFFSETS, getFontSizeLabel } from '@/stores/fontSizeStore';
import styles from './index.module.scss';

const { Title, Text } = Typography;

const GeneralSettingsPage: React.FC = () => {
  const { fontSizeOffset, setFontSizeOffset } = useFontSizeStore();
  
  const [previewOffset, setPreviewOffset] = useState(fontSizeOffset);
  const [hasChanges, setHasChanges] = useState(false);

  const handlePreviewChange = (value: number) => {
    setPreviewOffset(value as any);
    setHasChanges(value !== fontSizeOffset);
  };

  const handleApplySettings = () => {
    setFontSizeOffset(previewOffset as any);
    setHasChanges(false);
    message.success('字体设置已应用');
  };

  const handleResetSettings = () => {
    setPreviewOffset(FONT_SIZE_OFFSETS.NORMAL);
    setHasChanges(FONT_SIZE_OFFSETS.NORMAL !== fontSizeOffset);
  };

  const sliderMarks = FONT_SIZE_OPTIONS.reduce((marks, option) => {
    marks[option.value] = {
      style: { fontSize: '11px', color: 'var(--text-color-secondary)' },
      label: option.label,
    };
    return marks;
  }, {} as any);

  return (
    <div className={styles.generalSettings}>
      <div className={styles.content}>
        <Card title="显示设置" className={styles.settingCard}>
          <div className={styles.settingItem}>
            <div className={styles.settingLabel}>
              <Title level={5}>字体大小</Title>
              <Text type="secondary">调整应用中的文字大小，提升阅读体验</Text>
            </div>
            
            <div className={styles.fontSizeControl}>
              <div className={styles.sliderContainer}>
                <div className={styles.sliderHeader}>
                  <Text strong>字体大小</Text>
                  <div className={styles.currentSize}>
                    <Text strong className={styles.sizeLabel}>{getFontSizeLabel(previewOffset)}</Text>
                    <Text type="secondary" className={styles.sizeDesc}>
                      {14 + previewOffset}px
                    </Text>
                  </div>
                </div>
                <Slider
                  min={FONT_SIZE_OFFSETS.VERY_SMALL}
                  max={FONT_SIZE_OFFSETS.EXTRA_LARGE}
                  step={2}
                  value={previewOffset}
                  onChange={handlePreviewChange}
                  marks={sliderMarks}
                />
              </div>

              <div className={styles.presetOptions}>
                <Text strong>快速选择</Text>
                <div className={styles.presetButtons}>
                  {FONT_SIZE_OPTIONS.map((option) => (
                    <button
                      key={option.value}
                      className={`${styles.presetButton} ${previewOffset === option.value ? styles.active : ''}`}
                      onClick={() => handlePreviewChange(option.value)}
                    >
                      <span className={styles.buttonLabel}>{option.label}</span>
                      <span className={styles.buttonSize}>{option.description}</span>
                    </button>
                  ))}
                </div>
              </div>

              <div className={styles.previewArea}>
                <div className={styles.previewHeader}>
                  <FontSizeOutlined className={styles.previewIcon} />
                  <Title level={5}>预览效果</Title>
                </div>
                <div 
                  className={styles.previewContent}
                  style={{
                    fontSize: `${14 + previewOffset}px`,
                    lineHeight: 1.5715 + (previewOffset > 0 ? -0.05 : previewOffset < 0 ? 0.05 : 0)
                  }}
                >
                  <div className={styles.previewText}>
                    这是标准字体大小的文本预览。你可以通过调整上方的设置来改变文字的大小，找到最适合你阅读习惯的字体尺寸。
                  </div>
                  <div className={styles.previewSmall} style={{ fontSize: `${12 + previewOffset}px` }}>
                    小号文字：这是较小的辅助信息文本。
                  </div>
                  <div className={styles.previewBold} style={{ fontSize: `${14 + previewOffset}px`, fontWeight: 600 }}>
                    粗体文字：这是重要的加粗文本。
                  </div>
                  <div className={styles.previewSecondary} style={{ fontSize: `${14 + previewOffset}px`, opacity: 0.65 }}>
                    次要文字：这是次要信息文本。
                  </div>
                </div>
              </div>

              <div className={styles.actions}>
                <Button 
                  icon={<ReloadOutlined />}
                  onClick={handleResetSettings}
                  disabled={previewOffset === FONT_SIZE_OFFSETS.NORMAL}
                  className={styles.resetButton}
                >
                  恢复默认
                </Button>
                <Button 
                  type="primary" 
                  icon={<CheckOutlined />}
                  onClick={handleApplySettings}
                  disabled={!hasChanges}
                  className={styles.applyButton}
                >
                  应用设置
                </Button>
              </div>
            </div>
          </div>
        </Card>
      </div>
    </div>
  );
};

export default GeneralSettingsPage;