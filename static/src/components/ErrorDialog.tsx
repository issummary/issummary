import * as React from 'react';
import Dialog from 'material-ui/Dialog';
import FlatButton from 'material-ui/FlatButton';

interface IDialogAlertProps {
  open: boolean;
  error: string;
  onRequestClose: () => void;
}

export class ErrorDialog extends React.Component<IDialogAlertProps, any> {
  render() {
    const actions = [
      <FlatButton
        label="OK"
        primary={true}
        onClick={this.props.onRequestClose}
      />
    ];

    return (
      <Dialog
        title="Error occurred"
        actions={actions}
        modal={false}
        open={this.props.open}
        onRequestClose={this.props.onRequestClose}
      >
        {this.props.error}
      </Dialog>
    );
  }
}
