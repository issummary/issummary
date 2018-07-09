import * as _ from 'lodash';
import FloatingActionButton from 'material-ui/FloatingActionButton';
import AutoRenew from 'material-ui/svg-icons/action/autorenew';
import * as React from 'react';
import { CSSProperties } from 'react';
import { connect, Dispatch } from 'react-redux';
import { bindActionCreators } from 'redux';
import { errorDialogActionCreators, IErrorDialogActionCreators } from '../actions/errorDialog';
import { homeActionCreators, IHomeActionCreators } from '../actions/home';
import { IIssueTableActionCreators, issueTableActionCreators } from '../actions/issueTable';
import { IErrorDialogState } from '../reducers/errorDialog';
import { IHomeState } from '../reducers/home';
import { IRootState } from '../reducers/reducer';
import { filterWorksByProjectNames } from '../services/util';
import { ErrorDialog } from './ErrorDialog';
import { IIssueTableProps, IssueTable } from './IssueTable';
import { IssueTableConfig } from './IssueTableConfig';
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

const Refresh = (
  // tslint:disable-line
  props: IRefreshProps
) => (
  <FloatingActionButton style={style} onClick={props.onClick} disabled={props.isFetching}>
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

  public onClickRefreshButton() {
    this.props.actions.issueTable.requestUpdate();
  }

  public render() {
    const works =
      this.props.global.selectedProjectName === 'All'
        ? this.props.issueTable.works
        : filterWorksByProjectNames(this.props.issueTable.works, [this.props.global.selectedProjectName]);

    const maxClassNumWork = _.maxBy(works, w => (w.Label ? w.Label.ParentNames.length : 0));
    const maxClassNum =
      maxClassNumWork && maxClassNumWork.Label
        ? maxClassNumWork.Label.ParentNames.length + 1 // 1 is work own label
        : 0;

    return (
      <div>
        <ErrorDialog
          error={this.props.errorDialog.error}
          onRequestClose={this.props.actions.errorDialog.requestClosing}
          open={this.props.errorDialog.open}
        />
        <IssueTableConfig
          works={works}
          velocityPerManPerDay={this.props.global.velocityPerManPerDay}
          parallels={this.props.global.parallels}
          style={issueTableConfigStyle}
          onEnableManDay={this.props.actions.home.enableManDay}
          onDisableManDay={this.props.actions.home.disableManDay}
          onChangeParallels={this.props.actions.home.changeParallels}
          projectNames={this.props.issueTable.works
            .map(w => w.Issue.ProjectName)
            .filter((pn, i, self) => self.indexOf(pn) === i)}
          onChangeProjectSelectField={this.props.actions.home.changeProjectTextField}
        />
        <Refresh onClick={this.onClickRefreshButton} isFetching={this.props.global.isFetchingData} />
        <IssueTable
          works={works}
          milestones={this.props.issueTable.milestones}
          actions={this.props.actions.issueTable}
          showManDayColumn={this.props.global.showManDayColumn}
          showTotalManDayColumn={this.props.global.showTotalManDayColumn}
          showSPColumn={this.props.global.showSPColumn}
          showTotalSPColumn={this.props.global.showTotalSPColumn}
          velocityPerManPerDay={this.props.global.velocityPerManPerDay}
          parallels={this.props.global.parallels}
          selectedProjectName={this.props.global.selectedProjectName}
          maxClassNum={maxClassNum}
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

// tslint:disable-next-line variable-name
export const ConnectedHome = connect(mapStateToProps, mapDispatchToProps)(Home as any);
