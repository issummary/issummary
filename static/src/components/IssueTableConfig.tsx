import * as React from 'react';
import Toggle from 'material-ui/Toggle';
import { CSSProperties } from 'react';

interface IIssueTableConfigProps {
  style: CSSProperties;
  onEnableManDay: () => void;
  onDisableManDay: () => void;
}

export const IssueTableConfig = (props: IIssueTableConfigProps) => {
  const handleToggle = (event: object, isInputChecked: boolean) => {
    if (isInputChecked) {
      props.onEnableManDay();
    } else {
      props.onDisableManDay();
    }
  };

  return (
    <div style={props.style}>
      <Toggle label="ManDay" onToggle={handleToggle} />
    </div>
  );
};
