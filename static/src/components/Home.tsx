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
import { IssueTableConfig } from './IssueTableConfig';
import { IHomeState } from '../reducers/home';
import { homeActionCreators, IHomeActionCreators } from '../actions/home';

const style: CSSProperties = {
  margin: 0,
  top: 'auto',
  right: 20,
  bottom: 20,
  left: 'auto',
  position: 'fixed'
};

const issueTableConfigStyle: CSSProperties = {
  margin: 10
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
    home: IHomeActionCreators;
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
        <IssueTableConfig
          style={issueTableConfigStyle}
          onEnableManDay={this.props.actions.home.enableManDay}
          onDisableManDay={this.props.actions.home.disableManDay}
        />
        <Refresh
          onClick={this.onClickRefreshButton}
          isFetching={this.props.global.isFetchingData}
        />
        <IssueTable
          works={this.props.issueTable.works}
          actions={this.props.actions.issueTable}
          showManDayColumn={this.props.global.showManDayColumn}
          showTotalManDayColumn={this.props.global.showTotalManDayColumn}
          showSPColumn={this.props.global.showSPColumn}
          showTotalSPColumn={this.props.global.showTotalSPColumn}
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
      home: bindActionCreators(homeActionCreators as {}, dispatch),
      issueTable: bindActionCreators(issueTableActionCreators as {}, dispatch)
    }
  };
}

// tslint:disable-next-line variable-name
export const ConnectedHome = connect(mapStateToProps, mapDispatchToProps)(
  Home as any
);
