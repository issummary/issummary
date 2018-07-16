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
    return (
      <div>
        <ErrorDialog {...this.props.errorDialog} />
        <IssueTableConfig {...this.props.issueTableConfig} />
        <Refresh onClick={this.onClickRefreshButton} isFetching={this.props.isFetchingData} />
        <IssueTable {...this.props.issueTable} />
        <MilestoneTable milestones={this.props.issueTable.milestones} />
      </div>
    );
  }
}

function mapStateToProps(state: IRootState): ICombinedHomeState {
  return state.home;
}

function mapDispatchToProps(dispatch: Dispatch<any>) {
  return {
    errorDialog: bindActionCreators(errorDialogActionCreators as {}, dispatch),
    home: bindActionCreators(homeActionCreators as {}, dispatch),
    issueTable: bindActionCreators(issueTableActionCreators as {}, dispatch)
  };
}

function mergeProps(stateProps: ICombinedHomeState, dispatchProps: any, ownProps: any): IHomeProps {
  const works =
    stateProps.global.selectedProjectName === 'All'
      ? stateProps.issueTable.works
      : filterWorksByProjectNames(stateProps.issueTable.works, [stateProps.global.selectedProjectName]);

  const global = stateProps.global;

  const errorDialog: IErrorDialogProps = {
    ...stateProps.errorDialog,
    onRequestClose: dispatchProps.errorDialog.requestClosing
  };

  const maxClassNumWork = _.maxBy(works, w => (w.Label ? w.Label.ParentNames.length : 0));
  const maxClassNum =
    maxClassNumWork && maxClassNumWork.Label
      ? maxClassNumWork.Label.ParentNames.length + 1 // 1 is work own label
      : 0;

  const issueTable: IIssueTableProps = {
    ...stateProps.issueTable,
    actions: dispatchProps.issueTable,
    maxClassNum,
    parallels: global.parallels,
    selectedProjectName: global.selectedProjectName,
    showManDayColumn: global.showManDayColumn,
    showSPColumn: global.showSPColumn,
    showTotalManDayColumn: global.showTotalManDayColumn,
    showTotalSPColumn: global.showTotalSPColumn,
    velocityPerManPerDay: global.velocityPerManPerDay,
    works
  };

  const issueTableConfigStyle: CSSProperties = {
    margin: 10
  };

  const projectNames = stateProps.issueTable.works
    .map(w => w.Issue.ProjectName)
    .filter((pn, i, self) => self.indexOf(pn) === i);
  const issueTableConfig: IIssueTableConfigProps = {
    onChangeParallels: dispatchProps.changeParallels,
    onChangeProjectSelectField: dispatchProps.home.changeProjectTextField,
    onDisableManDay: dispatchProps.disableManDay,
    onEnableManDay: dispatchProps.home.enableManDay,
    parallels: stateProps.global.parallels,
    projectNames,
    style: issueTableConfigStyle,
    velocityPerManPerDay: stateProps.global.velocityPerManPerDay,
    works
  };

  return {
    ...ownProps,
    errorDialog,
    isFetchingData: stateProps.global.isFetchingData,
    issueTable,
    issueTableConfig,
    requestUpdate: dispatchProps.issueTable.requestUpdate,
    selectedProjectName: stateProps.global.selectedProjectName
  };
}

// tslint:disable-next-line variable-name
export const ConnectedHome = connect(mapStateToProps, mapDispatchToProps, mergeProps)(Home as any);
