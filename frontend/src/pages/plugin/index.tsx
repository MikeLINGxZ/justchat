import React from 'react';
import { useParams } from 'react-router-dom';
import PluginIframe from '@/components/plugin/PluginIframe';

const PluginPage: React.FC = () => {
  const { pluginId, pageId } = useParams<{ pluginId: string; pageId: string }>();

  if (!pluginId || !pageId) {
    return <div>Plugin not found</div>;
  }

  return (
    <div style={{ width: '100%', height: '100vh' }}>
      <PluginIframe
        pluginId={pluginId}
        entry={pageId}
        style={{ width: '100%', height: '100%' }}
      />
    </div>
  );
};

export default PluginPage;
