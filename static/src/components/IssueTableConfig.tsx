import * as React from 'react';
import Toggle from 'material-ui/Toggle';
import { CSSProperties } from 'react';
import TextField from 'material-ui/TextField';
import MenuItem from 'material-ui/MenuItem';
import SelectField from 'material-ui/SelectField';
import { ProjectSelectField } from './ProjectSelectField';
import RaisedButton from 'material-ui/RaisedButton';

interface IIssueTableConfigProps {
  style: CSSProperties;
  onEnableManDay: () => void;
  onDisableManDay: () => void;
  onChangeParallels: (parallels: number) => void;
  projectNames: string[];
  onChangeProjectSelectField: (p: string) => void;
}

export const IssueTableConfig = (props: IIssueTableConfigProps) => {
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

  const content = 'あいうえお,かきくけこ,さしすせそ';
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
