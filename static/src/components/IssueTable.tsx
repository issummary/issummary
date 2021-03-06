import * as React from 'react';

import { Table, TableBody } from 'material-ui/Table';
import TableHeader from 'material-ui/Table/TableHeader';
import { IIssueTableActionCreators } from '../actions/issueTable';
import { IMilestone } from '../models/milestone';
import { IWork } from '../models/work';
import { eachSum } from '../services/util';
import { IssueTableHeadersRow } from './IssueTableHeadersRow';
import { IssueTableRow } from './IssueTableRow';

export interface IIssueTableProps {
  works: IWork[];
  milestones: IMilestone[];
  showManDayColumn: boolean;
  showTotalManDayColumn: boolean;
  showSPColumn: boolean;
  showTotalSPColumn: boolean;
  velocityPerManPerDay: number;
  velocityPerWeek: number;
  selectedProjectName: string;
  maxClassNum: number;
  actions: IIssueTableActionCreators;
}

export class IssueTable extends React.Component<IIssueTableProps, any> {
  public componentDidMount() {
    this.props.actions.requestUpdate();
  }

  public render() {
    console.log(this.props.works); // tslint:disable-line
    const totalSPs = eachSum(this.props.works.map(w => w.StoryPoint));
    return (
      <Table fixedHeader={false} style={{ tableLayout: 'auto' }}>
        <TableHeader displaySelectAll={false} adjustForCheckbox={false}>
          <IssueTableHeadersRow
            showManDayColumn={this.props.showManDayColumn}
            showTotalManDayColumn={this.props.showTotalManDayColumn}
            showSPColumn={this.props.showSPColumn}
            showTotalSPColumn={this.props.showTotalSPColumn}
            maxClassNum={this.props.maxClassNum}
          />
        </TableHeader>

        <TableBody displayRowCheckbox={false}>
          {this.props.works.map((w, i) => (
            <IssueTableRow
              work={w}
              key={w.Issue.ProjectName + w.Issue.IID}
              totalSP={totalSPs[i]}
              showManDayColumn={this.props.showManDayColumn}
              showTotalManDayColumn={this.props.showTotalManDayColumn}
              showSPColumn={this.props.showSPColumn}
              showTotalSPColumn={this.props.showTotalSPColumn}
              velocityPerManPerDay={this.props.velocityPerManPerDay}
              velocityPerWeek={this.props.velocityPerWeek}
              maxClassNum={this.props.maxClassNum}
            />
          ))}
        </TableBody>
      </Table>
    );
  }
}
