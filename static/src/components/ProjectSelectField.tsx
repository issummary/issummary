import * as React from 'react';
import SelectField from 'material-ui/SelectField';
import MenuItem from 'material-ui/MenuItem';

interface ProjectSelectFieldProps {
  projectNames: string[];
  onChange: (projectName: string) => void;
}

export class ProjectSelectField extends React.Component<
  ProjectSelectFieldProps,
  any
> {
  state = {
    value: 'all'
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
          <MenuItem value={'all'} primaryText="All" key={'all'} />
          {this.props.projectNames.map(pn => {
            return <MenuItem value={pn} primaryText={pn} key={pn} />;
          })}
        </SelectField>
      </div>
    );
  }
}
