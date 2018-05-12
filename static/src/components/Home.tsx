import * as React from 'react';
import { ConnectedIssueTable } from './IssueTable';
import FloatingActionButton from 'material-ui/FloatingActionButton';
import AutoRenew from 'material-ui/svg-icons/action/autorenew';
import { CSSProperties } from 'react';

const style: CSSProperties = {
  margin: 0,
  top: 'auto',
  right: 20,
  bottom: 20,
  left: 'auto',
  position: 'fixed'
};

const Refresh = () => (
  <FloatingActionButton style={style}>
    <AutoRenew />
  </FloatingActionButton>
);

export class Home extends React.Component<{}, any> {
  constructor(props: {}) {
    super(props);
  }

  public render() {
    return (
      <div>
        <Refresh />
        <ConnectedIssueTable />
      </div>
    );
  }
}
