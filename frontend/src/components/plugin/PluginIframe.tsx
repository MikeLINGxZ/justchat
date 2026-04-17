import React, { useRef, useEffect, useState } from 'react';
import { Service } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service';

interface PluginIframeProps {
  pluginId: string;
  entry: string;     // relative path to HTML file within plugin dir
  style?: React.CSSProperties;
  className?: string;
  onMessage?: (data: any) => void;
}

const PluginIframe: React.FC<PluginIframeProps> = ({ pluginId, entry, style, className, onMessage }) => {
  const iframeRef = useRef<HTMLIFrameElement>(null);
  const [src, setSrc] = useState<string>('');

  // Resolve the plugin asset path via the backend
  useEffect(() => {
    Service.GetPluginAssetPath(pluginId, entry).then((path) => {
      setSrc(`file://${path}`);
    }).catch((err) => {
      console.error('Failed to resolve plugin asset path:', err);
    });
  }, [pluginId, entry]);

  useEffect(() => {
    const handler = (event: MessageEvent) => {
      if (iframeRef.current && event.source === iframeRef.current.contentWindow) {
        if (event.data?.type === 'resize' && event.data.height) {
          if (iframeRef.current) {
            iframeRef.current.style.height = `${event.data.height}px`;
          }
        }
        onMessage?.(event.data);
      }
    };
    window.addEventListener('message', handler);
    return () => window.removeEventListener('message', handler);
  }, [onMessage]);

  if (!src) return null;

  return (
    <iframe
      ref={iframeRef}
      src={src}
      sandbox="allow-scripts allow-forms"
      style={{
        width: '100%',
        height: '100%',
        border: 'none',
        ...style,
      }}
      className={className}
    />
  );
};

export default PluginIframe;
export type { PluginIframeProps };
