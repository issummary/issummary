import * as _ from 'lodash';
import FloatingActionButton from 'material-ui/FloatingActionButton';
import AutoRenew from 'material-ui/svg-icons/action/autorenew';
import * as React from 'react';
import { CSSProperties } from 'react';
import { connect, Dispatch } from 'react-redux';
import { bindActionCreators } from 'redux';
import { ActionCreator } from 'typescript-fsa';
import { errorDialogActionCreators } from '../actions/errorDialog';
import { homeActionCreators } from '../actions/home';
import { issueTableActionCreators } from '../actions/issueTable';
import { ICombinedHomeState, IRootState } from '../reducers/reducer';
import { filterWorksByProjectNames } from '../services/util';
import { ErrorDialog, IErrorDialogProps } from './ErrorDialog';
import { IIssueTableProps, IssueTable } from './IssueTable';
import { IIssueTableConfigProps, IssueTableConfig } from './IssueTableConfig';
import { MilestoneTable } from './milestoneTable';

const style: CSSProperties = {
  bottom: 20,
  left: 'auto',
  margin: 0,
  position: 'fixed',
  right: 20,
  top: 'auto'
};

const issueTableConfigStyle: CSSProperties = {
  margin: 10
};

interface IRefreshProps {
  onClick: React.MouseEventHandler<JSX.Element | HTMLElement>;
  isFetching: boolean;
}

// tslint:disable-next-line
const Refresh = (props: IRefreshProps) => (
  <FloatingActionButton style={style} onClick={props.onClick} disabled={props.isFetching}>
    <AutoRenew />
  </FloatingActionButton>
);

interface IHomeProps {
  selectedProjectName: string;
  isFetchingData: boolean;
  issueTable: IIssueTableProps;
  issueTableConfig: IIssueTableConfigProps;
  errorDialog: IErrorDialogProps;
  requestUpdate: ActionCreator<undefined>;
}

class Home extends React.Component<IHomeProps, any> {
  constructor(props: IHomeProps) {
    super(props);
    this.onClickRefreshButton = this.onClickRefreshButton.bind(this);
  }

  public onClickRefreshButton() {
    this.props.requestUpdate();
  }

  public render() {
    const works =
      this.props.selectedProjectName === 'All'
        ? this.props.issueTable.works
        : filterWorksByProjectNames(this.props.issueTable.works, [this.props.selectedProjectName]);

    return (
      <div>
        <ErrorDialog
          error={this.props.errorDialog.error}
          onRequestClose={this.props.errorDialog.onRequestClose}
          open={this.props.errorDialog.open}
        />
        <IssueTableConfig
          works={works}
          velocityPerManPerDay={this.props.issueTable.velocityPerManPerDay}
          parallels={this.props.issueTable.parallels}
          style={issueTableConfigStyle}
          onEnableManDay={this.props.issueTableConfig.onEnableManDay}
          onDisableManDay={this.props.issueTableConfig.onDisableManDay}
          onChangeParallels={this.props.issueTableConfig.onChangeParallels}
          projectNames={this.props.issueTable.works
            .map(w => w.Issue.ProjectName)
            .filter((pn, i, self) => self.indexOf(pn) === i)}
          onChangeProjectSelectField={this.props.issueTableConfig.onChangeProjectSelectField}
        />
        <Refresh onClick={this.onClickRefreshButton} isFetching={this.props.isFetchingData} />
        <IssueTable
          works={works}
          milestones={this.props.issueTable.milestones}
          actions={this.props.issueTable.actions}
          showManDayColumn={this.props.issueTable.showManDayColumn}
          showTotalManDayColumn={this.props.issueTable.showTotalManDayColumn}
          showSPColumn={this.props.issueTable.showSPColumn}
          showTotalSPColumn={this.props.issueTable.showTotalSPColumn}
          velocityPerManPerDay={this.props.issueTable.velocityPerManPerDay}
          parallels={this.props.issueTable.parallels}
          selectedProjectName={this.props.issueTable.selectedProjectName}
          maxClassNum={this.props.issueTable.maxClassNum}
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
      errorDialog: bindActionCreators(errorDialogActionCreators as {}, dispatch),
      home: bindActionCreators(homeActionCreators as {}, dispatch),
      issueTable: bindActionCreators(issueTableActionCreators as {}, dispatch)
    }
  };
}

function mergeProps(stateProps: ICombinedHomeState, dispatchProps: any, ownProps: any): IHomeProps {
  const actions = dispatchProps.actions;
  const global = stateProps.global;

  const errorDialog: IErrorDialogProps = {
    ...stateProps.errorDialog,
    onRequestClose: actions.errorDialog.requestClosing
  };

  const maxClassNumWork = _.maxBy(stateProps.issueTable.works, w => (w.Label ? w.Label.ParentNames.length : 0));
  const maxClassNum =
    maxClassNumWork && maxClassNumWork.Label
      ? maxClassNumWork.Label.ParentNames.length + 1 // 1 is work own label
      : 0;

  const issueTable: IIssueTableProps = {
    ...stateProps.issueTable,
    actions: actions.issueTable,
    maxClassNum,
    parallels: global.parallels,
    selectedProjectName: global.selectedProjectName,
    showManDayColumn: global.showManDayColumn,
    showSPColumn: global.showSPColumn,
    showTotalManDayColumn: global.showTotalManDayColumn,
    showTotalSPColumn: global.showTotalSPColumn,
    velocityPerManPerDay: global.velocityPerManPerDay
  };

  const issueTableConfig: IIssueTableConfigProps = {
    onChangeParallels: actions.changeParallels,
    onChangeProjectSelectField: actions.home.changeProjectTextField,
    onDisableManDay: actions.disableManDay,
    onEnableManDay: actions.home.enableManDay,
    parallels: stateProps.global.parallels,
    projectNames: stateProps.issueTable.works
      .map(w => w.Issue.ProjectName)
      .filter((pn, i, self) => self.indexOf(pn) === i),
    velocityPerManPerDay: stateProps.global.velocityPerManPerDay,
    works: stateProps.issueTable.works
  };

  return {
    ...ownProps,
    errorDialog,
    isFetchingData: stateProps.global.isFetchingData,
    issueTable,
    issueTableConfig,
    requestUpdate: actions.issueTable.requestUpdate,
    selectedProjectName: stateProps.global.selectedProjectName
  };
}

// tslint:disable-next-line variable-name
export const ConnectedHome = connect(mapStateToProps, mapDispatchToProps, mergeProps)(Home as any);
