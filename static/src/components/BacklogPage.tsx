import * as _ from 'lodash';
import FloatingActionButton from 'material-ui/FloatingActionButton';
import AutoRenew from 'material-ui/svg-icons/action/autorenew';
import * as moment from 'moment';
import { CSSProperties } from 'react';
import * as React from 'react';
import { connect, Dispatch } from 'react-redux';
import { bindActionCreators } from 'redux';
import { ActionCreator } from 'typescript-fsa';
import { backlogActionCreators } from '../actions/backlog';
import { backlogTableActionCreators } from '../actions/backlogTable';
import { errorDialogActionCreators } from '../actions/errorDialog';
import { IBacklogPageState } from '../reducers/backlog';
import { IRootState } from '../reducers/reducer';
import { worksToCSV } from '../services/csv';
import { filterWorksByProjectNames } from '../services/util';
import { BacklogTable, IBacklogTableProps } from './BacklogTable';
import { BacklogTableConfig, IBacklogTableConfigProps } from './BacklogTableConfig';
import { ErrorDialog, IErrorDialogProps } from './ErrorDialog';
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

interface IBacklogPageProps {
  selectedProjectName: string;
  isFetchingData: boolean;
  backlogTable: IBacklogTableProps;
  backlogTableConfig: IBacklogTableConfigProps;
  errorDialog: IErrorDialogProps;
  requestUpdate: ActionCreator<undefined>;
}

class BacklogPage extends React.Component<IBacklogPageProps, any> {
  constructor(props: IBacklogPageProps) {
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
        <BacklogTableConfig {...this.props.backlogTableConfig} />
        <Refresh onClick={this.onClickRefreshButton} isFetching={this.props.isFetchingData} />
        <BacklogTable {...this.props.backlogTable} />
        <MilestoneTable milestones={this.props.backlogTable.milestones} />
      </div>
    );
  }
}

function mapStateToProps(state: IRootState): IBacklogPageState {
  return state.backlogPage;
}

function mapDispatchToProps(dispatch: Dispatch<any>) {
  return {
    backlogPage: bindActionCreators(backlogActionCreators as {}, dispatch),
    backlogTable: bindActionCreators(backlogTableActionCreators as {}, dispatch),
    errorDialog: bindActionCreators(errorDialogActionCreators as {}, dispatch)
  };
}

function mergeProps(stateProps: IBacklogPageState, dispatchProps: any, ownProps: any): IBacklogPageProps {
  const works =
    stateProps.global.selectedProjectName === 'All'
      ? stateProps.backlogTable.works
      : filterWorksByProjectNames(stateProps.backlogTable.works, [stateProps.global.selectedProjectName]);

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

  const backlogTable: IBacklogTableProps = {
    ...stateProps.backlogTable,
    actions: dispatchProps.backlogTable,
    maxClassNum,
    parallels: global.parallels,
    selectedProjectName: global.selectedProjectName,
    showManDayColumn: stateProps.backlogTable.showManDayColumn,
    showSPColumn: stateProps.backlogTable.showSPColumn,
    showTotalManDayColumn: stateProps.backlogTable.showTotalManDayColumn,
    showTotalSPColumn: stateProps.backlogTable.showTotalSPColumn,
    velocityPerManPerDay: global.velocityPerManPerDay,
    works
  };

  const backlogTableConfigStyle: CSSProperties = {
    margin: 10
  };

  const projectNames = stateProps.backlogTable.works
    .map(w => w.Issue.ProjectName)
    .filter((pn, i, self) => self.indexOf(pn) === i);

  const content = worksToCSV(
    works,
    global.velocityPerManPerDay,
    moment(), // FIXME
    global.parallels
  );
  const blob = new Blob([content], { type: 'text/plain' });
  const csvUrl = window.URL.createObjectURL(blob);

  const backlogTableConfig: IBacklogTableConfigProps = {
    csvUrl,
    onChangeParallels: dispatchProps.changeParallels,
    onChangeProjectSelectField: dispatchProps.backlogPage.changeProjectTextField,
    onDisableManDay: dispatchProps.disableManDay,
    onEnableManDay: dispatchProps.backlogPage.enableManDay,
    parallels: stateProps.global.parallels,
    projectNames,
    style: backlogTableConfigStyle,
    velocityPerManPerDay: stateProps.global.velocityPerManPerDay,
    works
  };

  return {
    ...ownProps,
    backlogTable,
    backlogTableConfig,
    errorDialog,
    isFetchingData: stateProps.global.isFetchingData,
    requestUpdate: dispatchProps.backlogTable.requestUpdate,
    selectedProjectName: stateProps.global.selectedProjectName
  };
}

// tslint:disable-next-line variable-name
export const ConnectedBacklogPage = connect(mapStateToProps, mapDispatchToProps, mergeProps)(BacklogPage as any);
