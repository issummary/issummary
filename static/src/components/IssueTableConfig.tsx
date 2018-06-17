import RaisedButton from 'material-ui/RaisedButton';
import TextField from 'material-ui/TextField';
import Toggle from 'material-ui/Toggle';
import * as moment from 'moment';
import * as React from 'react';
import { CSSProperties } from 'react';
import { IWork } from '../models/work';
import { worksToCSV } from '../services/csv';
import { ProjectSelectField } from './ProjectSelectField';

interface IIssueTableConfigProps {
  works: IWork[];
  velocityPerManPerDay: number;
  parallels: number;
  style: CSSProperties;
  onEnableManDay: () => void;
  onDisableManDay: () => void;
  onChangeParallels: (parallels: number) => void;
  projectNames: string[];
  onChangeProjectSelectField: (p: string) => void;
}

export const IssueTableConfig = (props: IIssueTableConfigProps) => {// tslint:disable-line
  const handleToggle = (event: object, isInputChecked: boolean) => {
    if (isInputChecked) {
      props.onEnableManDay();
    } else {
      props.onDisableManDay();
    }
  };

  const handleParallelsChanging = (event: object, newParallels: string) => {
    const parallels = parseInt(newParallels, 10);
    if (!Number.isNaN(parallels)) {
      props.onChangeParallels(parallels);
    }
  };

  const content = worksToCSV(
    props.works,
    props.velocityPerManPerDay,
    moment(), // FIXME
    props.parallels
  );
  const blob = new Blob([content], { type: 'text/plain' });
  const csvUrl = window.URL.createObjectURL(blob);

  return (
    <div style={props.style}>
      <TextField
        defaultValue="2"
        floatingLabelText="Parallels"
        onChange={handleParallelsChanging}
      />
      <br />
      <Toggle label="ManDay" onToggle={handleToggle} />
      <ProjectSelectField
        projectNames={props.projectNames}
        onChange={props.onChangeProjectSelectField}
      />
      <a href={csvUrl} download="test.csv">
        <RaisedButton label="Export CSV" />
      </a>
    </div>
  );
};
