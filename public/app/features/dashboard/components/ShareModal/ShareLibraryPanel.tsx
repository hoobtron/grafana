import React, { useEffect } from 'react';

import { reportInteraction } from '@grafana/runtime/src';
import { VizPanel } from '@grafana/scenes';
import { Trans } from 'app/core/internationalization';
import { AddLibraryPanelContents } from 'app/features/library-panels/components/AddLibraryPanelModal/AddLibraryPanelModal';

import { ShareModalTabProps } from './types';

interface Props extends ShareModalTabProps {
  initialFolderUid?: string;
  vizPanel?: VizPanel;
}

export const ShareLibraryPanel = ({ panel, initialFolderUid, vizPanel, onDismiss }: Props) => {
  useEffect(() => {
    reportInteraction('grafana_dashboards_library_panel_share_viewed');
  }, []);

  if (!panel) {
    return null;
  }

  return (
    <>
      <p className="share-modal-info-text">
        <Trans i18nKey="share-modal.library.info">Create library panel.</Trans>
      </p>
      <AddLibraryPanelContents
        panel={panel}
        initialFolderUid={initialFolderUid}
        onDismiss={onDismiss}
        vizPanel={vizPanel}
      />
    </>
  );
};
