import * as React from 'react';

import { Table, TableBody } from 'material-ui/Table';
import TableHeader from 'material-ui/Table/TableHeader';
import { IBacklogTableActionCreators } from '../actions/backlogTable';
import { IMilestone } from '../models/milestone';
import { IWork } from '../models/work';
import { eachSum } from '../services/util';
import { BacklogTableHeadersRow } from './BacklogTableHeadersRow';
import { BacklogTableRow } from './BacklogTableRow';

export interface IBacklogTableProps {
  works: IWork[];
  milestones: IMilestone[];
  showManDayColumn: boolean;
  showTotalManDayColumn: boolean;
  showSPColumn: boolean;
  showTotalSPColumn: boolean;
  velocityPerManPerDay: number;
  parallels: number;
  selectedProjectName: string;
  maxClassNum: number;
  actions: IBacklogTableActionCreators; // FIXME convert to onXXX methods
}

export class BacklogTable extends React.Component<IBacklogTableProps, any> {
  public componentDidMount() {
    this.props.actions.requestUpdate();
  }

  public render() {
    console.log(this.props.works); // tslint:disable-line
    const totalSPs = eachSum(this.props.works.map(w => w.StoryPoint));

    return (
      <Table fixedHeader={false} style={{ tableLayout: 'auto' }}>
        <TableHeader displaySelectAll={false} adjustForCheckbox={false}>
          <BacklogTableHeadersRow
            showManDayColumn={this.props.showManDayColumn}
            showTotalManDayColumn={this.props.showTotalManDayColumn}
            showSPColumn={this.props.showSPColumn}
            showTotalSPColumn={this.props.showTotalSPColumn}
            maxClassNum={this.props.maxClassNum}
          />
        </TableHeader>

        <TableBody displayRowCheckbox={false}>
          {this.props.works.map((w, i) => (
            <BacklogTableRow
              work={w}
              key={w.Issue.ProjectName + w.Issue.IID}
              totalSP={totalSPs[i]}
              showManDayColumn={this.props.showManDayColumn}
              showTotalManDayColumn={this.props.showTotalManDayColumn}
              showSPColumn={this.props.showSPColumn}
              showTotalSPColumn={this.props.showTotalSPColumn}
              velocityPerManPerDay={this.props.velocityPerManPerDay}
              parallels={this.props.parallels}
              maxClassNum={this.props.maxClassNum}
            />
          ))}
        </TableBody>
      </Table>
    );
  }
}
