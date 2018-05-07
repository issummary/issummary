import * as React from 'react';
import {ConnectedIssueTable} from "./IssueTable";

export class Home extends React.Component<{}, any> {
    constructor(props: {}) {
        super(props);
    }

    public render() {
        return (
            <div>
                <h2>Home</h2>
                <ConnectedIssueTable/>
            </div>
        );
    }
}
