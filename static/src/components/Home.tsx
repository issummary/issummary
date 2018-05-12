import * as React from 'react';
import { IIssueTableProps, IssueTable } from './IssueTable';
import FloatingActionButton from 'material-ui/FloatingActionButton';
import AutoRenew from 'material-ui/svg-icons/action/autorenew';
import { CSSProperties } from 'react';
import { bindActionCreators } from 'redux';
import {
  IIssueTableActionCreators,
  issueTableActionCreators
} from '../actions/issueTable';
import { connect, Dispatch } from 'react-redux';
import { IRootState } from '../reducers/reducer';

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

interface IHomeProps {
  issueTable: IIssueTableProps;
  actions: {
    issueTable: IIssueTableActionCreators;
  };
}

class Home extends React.Component<IHomeProps, any> {
  constructor(props: IHomeProps) {
    super(props);
  }

  public render() {
    return (
      <div>
        <Refresh />
        <IssueTable
          works={this.props.issueTable.works}
          actions={this.props.actions.issueTable}
        />
      </div>
    );
  }
}

function mapStateToProps(state: IRootState) {
  return state.home;
}

function mapDispatchToProps(dispatch: Dispatch<any>) {
  return {
    actions: {
      issueTable: bindActionCreators(issueTableActionCreators as {}, dispatch)
    }
  };
}

// tslint:disable-next-line variable-name
export const ConnectedHome = connect(mapStateToProps, mapDispatchToProps)(
  Home as any
);
