import RaisedButton from 'material-ui/RaisedButton';
import TextField from 'material-ui/TextField';
import Toggle from 'material-ui/Toggle';
import * as React from 'react';
import { CSSProperties } from 'react';
import { withHandlers } from 'recompose';
import { IWork } from '../models/work';
import { ProjectSelectField } from './ProjectSelectField';

interface ICommonBacklogTableConfigProps {
  csvUrl: string;
  works: IWork[];
  velocityPerManPerDay: number;
  parallels: number;
  style?: CSSProperties;
  projectNames: string[];
  onChangeProjectSelectField: (p: string) => void;
}

interface IBaseBacklogTableConfigProps extends ICommonBacklogTableConfigProps {
  onToggle: (e: object, isInputChecked: boolean) => void;
  onChangeParallels: (e: object, newParallels: string) => void;
}

export interface IBacklogTableConfigProps extends ICommonBacklogTableConfigProps {
  onEnableManDay: () => void;
  onDisableManDay: () => void;
  onChangeParallels: (parallels: number) => void;
}

// tslint:disable-next-line
const BaseBacklogTableConfig = (props: IBaseBacklogTableConfigProps) => {
  return (
    <div style={props.style}>
      <TextField defaultValue="2" floatingLabelText="Parallels" onChange={props.onChangeParallels} />
      <br />
      <Toggle label="ManDay" onToggle={props.onToggle} />
      <ProjectSelectField projectNames={props.projectNames} onChange={props.onChangeProjectSelectField} />
      <a href={props.csvUrl} download="test.csv">
        <RaisedButton label="Export CSV" />
      </a>
    </div>
  );
};

// tslint:disable-next-line
export const BacklogTableConfig = withHandlers({
  onChangeParallels: (props: IBacklogTableConfigProps) => (event: object, newParallels: string) => {
    const parallels = parseInt(newParallels, 10);
    if (!Number.isNaN(parallels)) {
      props.onChangeParallels(parallels);
    }
  },
  onToggle: (props: IBacklogTableConfigProps) => (event: object, isInputChecked: boolean) => {
    if (isInputChecked) {
      props.onEnableManDay();
    } else {
      props.onDisableManDay();
    }
  }
})(BaseBacklogTableConfig);
