import * as React from 'react';
import Toggle from 'material-ui/Toggle';
import { CSSProperties } from 'react';
import TextField from 'material-ui/TextField';
import MenuItem from 'material-ui/MenuItem';
import SelectField from 'material-ui/SelectField';
import { ProjectSelectField } from './ProjectSelectField';

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
    </div>
  );
};
