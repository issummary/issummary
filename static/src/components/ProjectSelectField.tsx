import MenuItem from 'material-ui/MenuItem';
import SelectField from 'material-ui/SelectField';
import * as React from 'react';

interface IProjectSelectFieldProps {
  projectNames: string[];
  onChange: (projectName: string) => void;
}

export class ProjectSelectField extends React.Component<IProjectSelectFieldProps, any> {
  public state = {
    value: 'all'
  };

  public handleChange = (event: object, index: number, value: string) => {
    this.setState({ value });
    this.props.onChange(value);
  };

  public render() {
    return (
      <div>
        <SelectField floatingLabelText="Repository" value={this.state.value} onChange={this.handleChange}>
          <MenuItem value={'All'} primaryText="All" key={'all'} />
          {this.props.projectNames.map(pn => {
            return <MenuItem value={pn} primaryText={pn} key={pn} />;
          })}
        </SelectField>
      </div>
    );
  }
}
