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
import { IHomeState, IRootState } from '../reducers/reducer';

const style: CSSProperties = {
  margin: 0,
  top: 'auto',
  right: 20,
  bottom: 20,
  left: 'auto',
  position: 'fixed'
};

interface RefreshProps {
  onClick: React.MouseEventHandler<JSX.Element | HTMLElement>;
  isFetching: boolean;
}

const Refresh = (props: RefreshProps) => (
  <FloatingActionButton
    style={style}
    onClick={props.onClick}
    disabled={props.isFetching}
  >
    <AutoRenew />
  </FloatingActionButton>
);

interface IHomeProps {
  global: IHomeState;
  issueTable: IIssueTableProps;
  actions: {
    issueTable: IIssueTableActionCreators;
  };
}

class Home extends React.Component<IHomeProps, any> {
  constructor(props: IHomeProps) {
    super(props);
    this.onClickRefreshButton = this.onClickRefreshButton.bind(this);
  }

  onClickRefreshButton() {
    this.props.actions.issueTable.requestUpdate();
  }

  public render() {
    return (
      <div>
        <Refresh
          onClick={this.onClickRefreshButton}
          isFetching={this.props.global.isFetchingData}
        />
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
