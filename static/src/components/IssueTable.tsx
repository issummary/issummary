import * as React from 'react';

import {
    Table,
    TableBody,
    TableHeader,
    TableHeaderColumn,
    TableRow,
    TableRowColumn,
} from 'material-ui/Table';
import {connect, Dispatch} from "react-redux";
import {issueTableActionCreator} from "../actionCreators";
import {IRootState, IIssueTableRowProps, IIssueTableProps} from "../reducer";
import {bindActionCreators} from "redux";

const IssueTableRow = (props: IIssueTableRowProps) => (
    <TableRow key={props.iid}>
        <TableRowColumn>{props.iid}</TableRowColumn>
        <TableRowColumn>{props.classes.large}</TableRowColumn>
        <TableRowColumn>{props.classes.middle}</TableRowColumn>
        <TableRowColumn>{props.classes.small}</TableRowColumn>
        <TableRowColumn>{props.title}</TableRowColumn>
        <TableRowColumn style={{wordWrap: 'break-word', whiteSpace: 'normal'}}>{props.summary}</TableRowColumn>
        <TableRowColumn style={{wordWrap: 'break-word', whiteSpace: 'normal'}}>{props.note}</TableRowColumn>
    </TableRow>
);

const IssueTable = (props: IIssueTableProps) => (
    <Table>
        {console.log(props)}
        <TableHeader>
            <TableRow>
                <TableHeaderColumn>IID</TableHeaderColumn>
                <TableHeaderColumn>Large Class</TableHeaderColumn>
                <TableHeaderColumn>Middle Class</TableHeaderColumn>
                <TableHeaderColumn>Small Class</TableHeaderColumn>
                <TableHeaderColumn>Title</TableHeaderColumn>
                <TableHeaderColumn>Summary</TableHeaderColumn>
                <TableHeaderColumn>Note</TableHeaderColumn>
            </TableRow>
        </TableHeader>
        <TableBody>
            {props.rowProps.map(rowProp => IssueTableRow(rowProp))}
        </TableBody>
    </Table>
);

function mapStateToProps(state: IRootState) {
    return  state.issueTable;
}

function mapDispatchToProps(dispatch: Dispatch<any>) {
    return { actions: bindActionCreators(issueTableActionCreator, dispatch) };
}

// tslint:disable-next-line variable-name
export const ConnectedIssueTable = connect(mapStateToProps, mapDispatchToProps)(IssueTable as any);
