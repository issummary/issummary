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
import { MilestoneTable } from './milestoneTable';
import { ErrorDialog } from './ErrorDialog';
import {
  errorDialogActionCreators,
  IErrorDialogActionCreators
} from '../actions/errorDialog';
import { IErrorDialogState } from '../reducers/errorDialog';

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
  errorDialog: IErrorDialogState;
  actions: {
    home: IHomeActionCreators;
    issueTable: IIssueTableActionCreators;
    errorDialog: IErrorDialogActionCreators;
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
        <ErrorDialog
          error={this.props.errorDialog.error}
          onRequestClose={this.props.actions.errorDialog.requestClosing}
          open={this.props.errorDialog.open}
        />
        <IssueTableConfig
          style={issueTableConfigStyle}
          onEnableManDay={this.props.actions.home.enableManDay}
          onDisableManDay={this.props.actions.home.disableManDay}
          onChangeParallels={this.props.actions.home.changeParallels}
          projectNames={this.props.issueTable.works
            .map(w => w.Issue.ProjectName)
            .filter((pn, i, self) => self.indexOf(pn) === i)}
          onChangeProjectSelectField={
            this.props.actions.home.changeProjectTextField
          }
        />
        <Refresh
          onClick={this.onClickRefreshButton}
          isFetching={this.props.global.isFetchingData}
        />
        <IssueTable
          works={this.props.issueTable.works}
          milestones={this.props.issueTable.milestones}
          actions={this.props.actions.issueTable}
          showManDayColumn={this.props.global.showManDayColumn}
          showTotalManDayColumn={this.props.global.showTotalManDayColumn}
          showSPColumn={this.props.global.showSPColumn}
          showTotalSPColumn={this.props.global.showTotalSPColumn}
          velocityPerManPerDay={this.props.global.velocityPerManPerDay}
          parallels={this.props.global.parallels}
          selectedProjectName={this.props.global.selectedProjectName}
        />
        <MilestoneTable milestones={this.props.issueTable.milestones} />
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
      issueTable: bindActionCreators(issueTableActionCreators as {}, dispatch),
      errorDialog: bindActionCreators(errorDialogActionCreators as {}, dispatch)
    }
  };
}

// tslint:disable-next-line variable-name
export const ConnectedHome = connect(mapStateToProps, mapDispatchToProps)(
  Home as any
);
