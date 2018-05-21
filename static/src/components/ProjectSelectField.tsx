import * as React from 'react';
import SelectField from 'material-ui/SelectField';
import MenuItem from 'material-ui/MenuItem';

interface ProjectSelectFieldProps {
  onChange: (projectName: string) => void;
}

export class ProjectSelectField extends React.Component<
  ProjectSelectFieldProps,
  any
> {
  state = {
    value: 1
  };

  handleChange = (event: object, index: number, value: string) => {
    this.setState({ value });
    this.props.onChange(value);
  };

  render() {
    return (
      <div>
        <SelectField
          floatingLabelText="Frequency"
          value={this.state.value}
          onChange={this.handleChange}
        >
          <MenuItem value={'all'} primaryText="All" />
          <MenuItem value={'ProjectA'} primaryText="ProjectA" />
        </SelectField>
      </div>
    );
  }
}
