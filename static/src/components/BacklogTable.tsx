import * as React from 'react';

import { Table, TableBody } from 'material-ui/Table';
import TableHeader from 'material-ui/Table/TableHeader';
import { lifecycle } from 'recompose';
import { ActionCreator } from 'typescript-fsa';
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
  onRequestUpdate: ActionCreator<undefined>;
}

// tslint:disable-next-line variable-name
export const BacklogTable = (props: IBacklogTableProps) => {
  console.log(props.works); // tslint:disable-line no-console
  const totalSPs = eachSum(props.works.map(w => w.StoryPoint));

  return (
    <Table fixedHeader={false} style={{ tableLayout: 'auto' }}>
      <TableHeader displaySelectAll={false} adjustForCheckbox={false}>
        <BacklogTableHeadersRow
          showManDayColumn={props.showManDayColumn}
          showTotalManDayColumn={props.showTotalManDayColumn}
          showSPColumn={props.showSPColumn}
          showTotalSPColumn={props.showTotalSPColumn}
          maxClassNum={props.maxClassNum}
        />
      </TableHeader>

      <TableBody displayRowCheckbox={false}>
        {props.works.map((w, i) => (
          <BacklogTableRow
            work={w}
            key={w.Issue.ProjectName + w.Issue.IID}
            totalSP={totalSPs[i]}
            showManDayColumn={props.showManDayColumn}
            showTotalManDayColumn={props.showTotalManDayColumn}
            showSPColumn={props.showSPColumn}
            showTotalSPColumn={props.showTotalSPColumn}
            velocityPerManPerDay={props.velocityPerManPerDay}
            parallels={props.parallels}
            maxClassNum={props.maxClassNum}
          />
        ))}
      </TableBody>
    </Table>
  );
};

// tslint:disable-next-line variable-name
export const BacklogTableWithRequest = lifecycle<IBacklogTableProps, undefined, {}>({
  componentDidMount() {
    this.props.onRequestUpdate();
  }
})(BacklogTable);
