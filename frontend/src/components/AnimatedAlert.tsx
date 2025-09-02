import React, { useState, useEffect } from 'react';
import { Alert } from 'antd';
import type { AlertProps } from 'antd';

interface AnimatedAlertProps extends AlertProps {
  visible: boolean;
  onExitComplete?: () => void;
  exitDuration?: number;
}

const AnimatedAlert: React.FC<AnimatedAlertProps> = ({
  visible,
  onExitComplete,
  exitDuration = 250,
  onClose,
  ...alertProps
}) => {
  const [shouldRender, setShouldRender] = useState(visible);
  const [isExiting, setIsExiting] = useState(false);

  useEffect(() => {
    if (visible) {
      setShouldRender(true);
      setIsExiting(false);
    } else if (shouldRender) {
      // 开始退场动画
      setIsExiting(true);
      const timer = setTimeout(() => {
        setShouldRender(false);
        setIsExiting(false);
        onExitComplete?.();
      }, exitDuration);

      return () => clearTimeout(timer);
    }
  }, [visible, shouldRender, exitDuration, onExitComplete]);

  const handleClose = (e: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
    if (onClose) {
      onClose(e);
    }
  };

  if (!shouldRender) {
    return null;
  }

  return (
    <Alert
      {...alertProps}
      onClose={handleClose}
      style={{
        ...alertProps.style,
        animation: isExiting 
          ? 'slideOutFade 0.25s ease-in forwards' 
          : 'slideInFade 0.3s ease-out forwards',
      }}
    />
  );
};

export default AnimatedAlert;
